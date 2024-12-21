package handler

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/KIVUOS1999/easyApi/app"
	"github.com/KIVUOS1999/easyApi/constants"
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

/*
alg:RS256
aud:112183527117-8i7nqtfeev3s1ai1fvla0huapm72kls1.apps.googleusercontent.com
azp:112183527117-8i7nqtfeev3s1ai1fvla0huapm72kls1.apps.googleusercontent.com
email:souviksarkar.ronnie@gmail.com
email_verified:true
exp:1734636877
family_name:Sarkar
given_name:Souvik
iat:1734633277
iss:https://accounts.google.com
jti:241497a46155a0050398a823e203bcdb8d2cea52
kid:564feacec3ebdfaa7311b9d8e73c42818f291264
name:Souvik Sarkar
nbf:1734632977
picture:https://lh3.googleusercontent.com/a/ACg8ocLFmT9rKWtEVLcxNqRP8KIvlq9eCF7QjW-ItpccdZ44P5axTSYp=s96-c
sub:116214675523781428407
typ:JWT
*/
func (h *handlerStruct) getUserDetails(idToken string) (*models.TokenData, error) {
	api := "https://oauth2.googleapis.com/tokeninfo?id_token=" + idToken
	req, err := http.NewRequest(http.MethodGet, api, nil)
	if err != nil {
		log.Error("new request:", err.Error())
		return nil, err
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Error("do:", err.Error())
		return nil, err
	}

	tokenData := models.TokenData{}
	json.NewDecoder(resp.Body).Decode(&tokenData)

	if resp.StatusCode != http.StatusOK {
		log.Error("status:", resp.StatusCode)

		return nil, &easyError.CustomError{
			StatusCode: http.StatusInternalServerError,
			Response:   "return status code not status ok:" + strconv.Itoa(resp.StatusCode),
		}
	}

	return &tokenData, nil
}

func (h *handlerStruct) AddUser(ctx *app.Context) (interface{}, error) {
	idToken := ctx.Request.GetHeader("Authorization")

	userData, err := h.getUserDetails(idToken)
	if err != nil {
		return nil, err
	}

	log.Debugf("user details: %+v", userData)

	user, _ := h.dataSvc.GetUser(userData.ID)
	if user != nil {
		log.Infof("user present - %+v", user)
		userData.AllotedSize = user.AllotedSize

		return userData, nil
	}

	log.Info("user not present -", userData.ID)

	err = h.dataSvc.AddUser(userData)
	if err != nil {
		return nil, err
	}

	userData.AllotedSize = constants.DefaultUserAllotment

	return userData, nil
}

func (h *handlerStruct) UploadFile(ctx *app.Context) (interface{}, error) {
	userID := ctx.PathParam("user-id")
	if userID == "" {
		log.Error("empty user id")
		return nil, &easyError.CustomError{
			StatusCode: http.StatusBadRequest,
			Response:   "user-id is empty",
		}
	}

	fileStructure := models.FileUploadStructure{}

	err := ctx.Bind(&fileStructure)
	if err != nil {
		log.Error("Bind err:", err.Error())

		return nil, err
	}

	currentTime := time.Now().UTC().Unix()
	fileID := uuid.New()

	fileStructure.UserID = userID
	fileStructure.ID = fileID
	fileStructure.CreatedAt = currentTime

	log.Infof("uploading file data success: file_id, user_id = %s, %s", fileID.String(), userID)

	// checks for total size.
	ok, err := h.hasUserSizeExceed(ctx, userID, fileStructure.Meta.Size)
	if err != nil {
		return nil, err
	}

	if ok {
		return nil, &easyError.CustomError{
			StatusCode: http.StatusInsufficientStorage,
			Response:   "you have exceeded you limit",
		}
	}

	// DB Upload
	err = h.dataSvc.UploadFileDetails(&fileStructure)
	if err != nil {
		log.Error("data_svc:", err.Error())
		return nil, err
	}

	return fileStructure, nil
}

func (h *handlerStruct) UploadChunks(ctx *app.Context) (interface{}, error) {
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
	userID := ctx.PathParam("user-id")

	return h.dataSvc.GetFilesByUser(userID)
}

func (h *handlerStruct) GetChunks(ctx *app.Context) (interface{}, error) {
	fileID := ctx.PathParam("file-id")
	if fileID == "" {
		return nil, &easyError.CustomError{
			StatusCode: http.StatusBadRequest,
			Response:   "file-id not passed",
		}
	}

	log.Debug("[DOWNLOAD] quering for chunks:", fileID)

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

func (h *handlerStruct) DeleteFile(ctx *app.Context) (interface{}, error) {
	fileID := ctx.PathParam("file-id")

	chunks, err := h.dataSvc.GetChunks(fileID)
	if err != nil {
		return nil, err
	}

	files := []string{}
	for idx := range chunks {
		files = append(files, "./temp/"+fileID+"_"+chunks[idx].ID.String())
	}

	for _, file := range files {
		deleteFile(file)
	}

	return nil, h.dataSvc.DeleteFile(fileID)
}

func deleteFile(destPath string) error {
	log.Debug("delete", destPath)

	err := os.Remove(destPath)
	if err != nil {
		log.Error(destPath, err.Error())
		return err
	}

	return nil
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

func (h *handlerStruct) hasUserSizeExceed(_ *app.Context, userID string, currentFileSize uint64) (bool, error) {
	uploadedFile, err := h.dataSvc.GetFilesByUser(userID)
	if err != nil {
		return false, err
	}

	var size uint64

	for idx := range uploadedFile.FileArr {
		size += uploadedFile.FileArr[idx].Meta.Size
	}

	usedSize := uploadedFile.CalculateSize() + currentFileSize

	if usedSize < size {
		log.Warnf("%s has exceed the size limit : %+v / %+v", userID, size, usedSize)
		return true, nil
	}

	log.Infof("%s size limit : %+v / %+v", userID, size, usedSize)
	return false, nil
}
