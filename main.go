package main

import (
	"net/http"
	"strconv"

	"github.com/KIVUOS1999/easyApi/app"
	easyError "github.com/KIVUOS1999/easyApi/errors"
	"github.com/KIVUOS1999/easyLogs/pkg/log"
	"github.com/KIVUOS1999/file-uploader-orch/handler"
	"github.com/KIVUOS1999/file-uploader-orch/service"
)

// Function to validate the ID Token
func verifyIDToken(idToken string) error {
	// validating using the google endpoint
	// this is not a recommended way
	// will fix it later

	api := "https://oauth2.googleapis.com/tokeninfo?id_token=" + idToken
	req, err := http.NewRequest(http.MethodGet, api, nil)
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

func validateAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			authToken := r.Header.Get("Authorization")
			if authToken == "" {
				log.Error("[AUTH EMPTY]:", r.URL.Path)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			err := verifyIDToken(authToken)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}

func main() {
	app := app.New()

	url := app.Configs.Get("DATA_SVC_HOST")
	dataSvc := service.New(url)
	handler := handler.New(dataSvc)

	app.Muxx.Use(validateAuth)

	app.Post("/upload_file/{user-id}", handler.UploadFile)
	app.Post("/upload_chunks", handler.UploadChunks)
	app.Post("/user", handler.AddUser)

	app.Get("/files/{user-id}", handler.GetFileByUser)
	app.Get("/chunks/{file-id}", handler.GetChunks)
	app.Get("/download/{chunk-name}", handler.DownloadChunk)

	app.Delete("/file/{file-id}", handler.DeleteFile)

	app.Start()
}
