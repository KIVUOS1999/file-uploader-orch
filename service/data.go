package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	easyError "github.com/KIVUOS1999/easyApi/errors"
	"github.com/KIVUOS1999/easyLogs/pkg/log"
	dataModels "github.com/KIVUOS1999/file-uploader-db/models"
	"github.com/KIVUOS1999/file-uploader-orch/pkg/models"
)

type dataSvc struct {
	url string
}

type IDataSvc interface {
	UploadFileDetails(fileDetails *models.FileUploadStructure) error
	UploadChunkDetails(chunkDetails *models.FileChunkStructure) error
	AddUser(userData *models.TokenData) error

	GetFilesByUser(userID string) ([]models.FileUploadStructure, error)
	GetChunks(fileID string) ([]models.FileChunkStructure, error)
	GetUser(userID string) (*dataModels.Users, error)

	DeleteFile(fileID string) error
}

func New(dataURL string) IDataSvc {
	return &dataSvc{
		url: dataURL,
	}
}

func (svc *dataSvc) UploadFileDetails(fileDetails *models.FileUploadStructure) error {
	dataFileUploadStruct := dataModels.FileDetailStructure{
		Meta: dataModels.FileMetaData{
			Name: fileDetails.Meta.Name,
			Size: fileDetails.Meta.Size,
		},
		ID:          fileDetails.ID,
		TotalChunks: fileDetails.TotalChunks,
		UserID:      fileDetails.UserID,
	}

	jsonData, err := json.Marshal(dataFileUploadStruct)
	if err != nil {
		log.Error(fileDetails.ID, "marshal:", err.Error())
		return err
	}

	api := svc.url + "/upload_file"

	log.Debug("Uploading file endpoint:", api)

	req, err := http.NewRequest(http.MethodPost, api, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error(fileDetails.ID, "data_svc:", err.Error())
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(fileDetails.ID, "send:", err.Error())
		return err
	}

	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Error(fileDetails.ID, "readall:", err.Error())
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		log.Error(fileDetails.ID, "status:", resp.StatusCode)
		return &easyError.CustomError{
			StatusCode: http.StatusInternalServerError,
			Response:   "return status code not status created:" + strconv.Itoa(resp.StatusCode),
		}
	}

	return nil
}

func (svc *dataSvc) UploadChunkDetails(chunkDetails *models.FileChunkStructure) error {
	chunkDataStruct := dataModels.FileChunkStructure{
		ID:       chunkDetails.ID,
		FileID:   chunkDetails.FileID,
		CheckSum: chunkDetails.CheckSum,
		Order:    chunkDetails.Order,
	}

	jsonData, err := json.Marshal(chunkDataStruct)
	if err != nil {
		log.Error(chunkDetails.ID, "json:", err.Error())
		return err
	}

	api := svc.url + "/upload_chunks"

	log.Debug("Uploading chunks endpoint:", api)

	req, err := http.NewRequest(http.MethodPost, api, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error(chunkDetails.ID, "new request:", err.Error())
		return err
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Error(chunkDetails.ID, "do:", err.Error())
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		log.Error(chunkDetails.ID, "status:", resp.StatusCode)
		return &easyError.CustomError{
			StatusCode: http.StatusInternalServerError,
			Response:   "return status code not status created:" + strconv.Itoa(resp.StatusCode),
		}
	}

	return nil
}

func (svc *dataSvc) AddUser(userData *models.TokenData) error {
	userStruct := dataModels.Users{
		ID:      userData.ID,
		Name:    userData.Name,
		Email:   userData.Email,
		Picture: userData.Picture,
	}

	jsonData, err := json.Marshal(userStruct)
	if err != nil {
		log.Error(userData.ID, "json:", err.Error())
		return err
	}

	api := svc.url + "/user"

	log.Debug("Uploading chunks endpoint:", api)

	req, err := http.NewRequest(http.MethodPost, api, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error(userData.ID, "new request:", err.Error())
		return err
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Error(userData.ID, "do:", err.Error())
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		log.Error(userData.ID, "status:", resp.StatusCode)
		return &easyError.CustomError{
			StatusCode: http.StatusInternalServerError,
			Response:   "return status code not status created:" + strconv.Itoa(resp.StatusCode),
		}
	}

	return nil
}

func (svc *dataSvc) GetFilesByUser(userID string) ([]models.FileUploadStructure, error) {
	api := svc.url + "/files/" + userID
	log.Debug("files by user:", api)

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

	if resp.StatusCode != http.StatusOK {
		log.Error("status:", resp.StatusCode)
		return nil, &easyError.CustomError{
			StatusCode: http.StatusInternalServerError,
			Response:   "return status code not status ok:" + strconv.Itoa(resp.StatusCode),
		}
	}

	files := []models.FileUploadStructure{}
	err = json.NewDecoder(resp.Body).Decode(&files)
	if err != nil {
		log.Error("decode response body:", err.Error())
		return nil, err
	}

	return files, nil
}

func (svc *dataSvc) GetChunks(fileID string) ([]models.FileChunkStructure, error) {
	api := svc.url + "/chunks/" + fileID

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

	if resp.StatusCode != http.StatusOK {
		log.Error("status:", resp.StatusCode)
		return nil, &easyError.CustomError{
			StatusCode: http.StatusInternalServerError,
			Response:   "return status code not status ok:" + strconv.Itoa(resp.StatusCode),
		}
	}

	files := []models.FileChunkStructure{}
	err = json.NewDecoder(resp.Body).Decode(&files)
	if err != nil {
		log.Error("decode response body:", err.Error())
		return nil, err
	}

	return files, nil
}

func (svc *dataSvc) GetUser(userID string) (*dataModels.Users, error) {
	api := svc.url + "/user/" + userID

	log.Info("api", api)

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

	if resp.StatusCode != http.StatusOK {
		log.Error("status:", resp.StatusCode)
		return nil, &easyError.CustomError{
			StatusCode: http.StatusInternalServerError,
			Response:   "return status code not status ok:" + strconv.Itoa(resp.StatusCode),
		}
	}

	userDetails := dataModels.Users{}
	err = json.NewDecoder(resp.Body).Decode(&userDetails)
	if err != nil {
		log.Error("decode response body:", err.Error())
		return nil, err
	}

	return &userDetails, nil
}

func (svc *dataSvc) DeleteFile(fileID string) error {
	api := svc.url + "/file/" + fileID

	req, err := http.NewRequest(http.MethodDelete, api, nil)
	if err != nil {
		log.Error("new request:", err.Error())
		return err
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Error("do:", err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("status:", resp.StatusCode)
		return &easyError.CustomError{
			StatusCode: http.StatusInternalServerError,
			Response:   "return status code not status ok:" + strconv.Itoa(resp.StatusCode),
		}
	}

	return nil
}
