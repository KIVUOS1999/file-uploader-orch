package handler

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/KIVUOS1999/easyApi/app"
	easyError "github.com/KIVUOS1999/easyApi/errors"
	"github.com/KIVUOS1999/easyLogs/pkg/log"
	"github.com/KIVUOS1999/file-uploader-orch/pkg/models"
	"github.com/KIVUOS1999/file-uploader-orch/service"
	"github.com/google/uuid"
)

type handlerStruct struct {
	dataSvc service.IDataSvc
}

func New(dataSvc service.IDataSvc) *handlerStruct {
	return &handlerStruct{
		dataSvc: dataSvc,
	}
}

func (h *handlerStruct) UploadFile(ctx *app.Context) (interface{}, error) {
	ctx.Response.AddCORS()

	fileStructure := models.FileUploadStructure{}

	err := ctx.Bind(&fileStructure)
	if err != nil {
		log.Error("Bind err:", err.Error())

		return nil, err
	}

	currentTime := time.Now().UTC().Unix()
	fileID := uuid.New()

	fileStructure.ID = fileID
	fileStructure.CreatedAt = currentTime

	log.Infof("uploading file data success: file_id = %s", fileID)

	// DB Upload
	err = h.dataSvc.UploadFileDetails(&fileStructure)
	if err != nil {
		log.Error("data_svc:", err.Error())
		return nil, err
	}

	return fileStructure, nil
}

func (h *handlerStruct) UploadChunks(ctx *app.Context) (interface{}, error) {
	ctx.Response.AddCORS()

	err := ctx.Request.Req.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Error("Error in parse multipart", err.Error())

		return nil, err
	}

	file, _, err := ctx.Request.Req.FormFile("chunk")
	if err != nil {
		log.Error("Error in form file", err.Error())

		return nil, err
	}

	defer file.Close()

	chunkNumber, err := strconv.Atoi(ctx.Request.Req.FormValue("chunk_number"))
	if err != nil {
		return nil, err
	}

	fileID := ctx.Request.Req.FormValue("file_id")

	fileUUID, err := uuid.Parse(fileID)
	if err != nil {
		log.Error("Error in parsing file uuid", fileID)

		return nil, err
	}

	chunkID := uuid.New()

	dest := "./temp/" + fileID + "_" + chunkID.String()

	checksum, err := saveFile(dest, file)
	if err != nil {
		log.Error("Failed to save chunk", chunkID, "for file", fileID)

		return nil, err
	}

	fileChunks := models.FileChunkStructure{
		ID:       chunkID,
		FileID:   fileUUID,
		Order:    chunkNumber,
		CheckSum: checksum,
	}

	log.Infof("Chunk created id: %s ord: %d for file: %s", chunkID, chunkNumber, fileChunks.FileID)

	// DB upload
	h.dataSvc.UploadChunkDetails(&fileChunks)
	return fileChunks, nil
}

func (h *handlerStruct) GetFileByUser(ctx *app.Context) (interface{}, error) {
	ctx.Response.AddCORS()

	userID := ctx.PathParam("user-id")
	log.Debug("user", userID)

	return h.dataSvc.GetFilesByUser(userID)
}

func (h *handlerStruct) GetChunks(ctx *app.Context) (interface{}, error) {
	ctx.Response.AddCORS()

	fileID := ctx.PathParam("file-id")
	if fileID == "" {
		return nil, &easyError.CustomError{
			StatusCode: http.StatusBadRequest,
			Response:   "file-id not passed",
		}
	}

	chunks, err := h.dataSvc.GetChunks(fileID)
	if err != nil {
		return nil, err
	}

	if len(chunks) == 0 {
		return nil, &easyError.CustomError{
			StatusCode: http.StatusInternalServerError,
			Response:   "entity not found",
		}
	}

	return chunks, nil
}

func (h *handlerStruct) DownloadChunk(ctx *app.Context) (interface{}, error) {
	ctx.Response.AddCORS()

	chunkName := ctx.PathParam("chunk-name")

	log.Debug("requesting chunk:", chunkName)

	file, err := os.Open("temp/" + chunkName)
	if err != nil {
		log.Error("Error in file open", err.Error())
		return nil, err
	}

	defer file.Close()

	written, err := io.Copy(ctx.Response.Resp, file)
	if err != nil {
		log.Error("Error in file copy", err.Error())
		return nil, err
	}

	log.Info("written - %+v", written)
	return nil, nil
}

func saveFile(destPath string, file multipart.File) (string, error) {
	destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error("Error in opening file", err.Error())

		return "", err
	}

	defer destFile.Close()

	hash := md5.New()

	_, err = io.Copy(io.MultiWriter(destFile, hash), file)
	if err != nil {
		log.Error("Error in copying data to file", err.Error())

		return "", err
	}

	checksum := hex.EncodeToString(hash.Sum(nil))

	return checksum, nil
}
