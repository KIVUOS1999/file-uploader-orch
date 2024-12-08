package main

import (
	"github.com/KIVUOS1999/easyApi/app"
	"github.com/KIVUOS1999/file-uploader-orch/handler"
	"github.com/KIVUOS1999/file-uploader-orch/service"
)

func main() {
	app := app.New()

	url := app.Configs.Get("DATA_SVC_HOST")
	dataSvc := service.New(url)
	handler := handler.New(dataSvc)

	app.Post("/upload_file", handler.UploadFile)
	app.Post("/upload_chunks", handler.UploadChunks)

	app.Get("/files/{user-id}", handler.GetFileByUser)
	app.Get("/chunks/{file-id}", handler.GetChunks)
	app.Get("/download/{chunk-name}", handler.DownloadChunk)

	app.Delete("/file/{file-id}", handler.DeleteFile)
	// app.Options("/file/{file-id}", handler.DeleteFile)

	app.Start()
}
