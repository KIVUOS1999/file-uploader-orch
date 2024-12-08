module github.com/KIVUOS1999/file-uploader-orch

go 1.23.2

require (
	github.com/KIVUOS1999/easyApi v0.0.0-20241117070720-954caed24eaa
	github.com/KIVUOS1999/easyLogs v1.0.0
	github.com/KIVUOS1999/file-uploader-db v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
)

replace github.com/KIVUOS1999/easyApi => ../../easyApi

replace github.com/KIVUOS1999/file-uploader-db => ../file-uploader-db

require (
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
)
