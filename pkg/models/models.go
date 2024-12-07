package models

import (
	"github.com/google/uuid"
)

type FileMetaData struct {
	Name string `json:"name"`
	Size int64  `json:"file_size"`
}

type FileUploadStructure struct {
	Meta FileMetaData `json:"meta_data"`

	ID          uuid.UUID `json:"file_id"`
	TotalChunks int       `json:"total_chunks"`
	CreatedAt   int64     `json:"created_at"`
}

type FileChunkStructure struct {
	ID       uuid.UUID `json:"chunk_id"`
	FileID   uuid.UUID `json:"file_id"`
	CheckSum string    `json:"check_sum"`
	Order    int       `json:"order"`
}
