package models

import (
	"github.com/google/uuid"
)

type FileMetaData struct {
	Name string `json:"name"`
	Size uint64 `json:"file_size"`
}

type FileUploadStructure struct {
	Meta FileMetaData `json:"meta_data"`

	ID          uuid.UUID `json:"file_id"`
	UserID      string    `json:"user_id"`
	TotalChunks int       `json:"total_chunks"`
	CreatedAt   int64     `json:"created_at"`
}

type FileList struct {
	UsedSize uint64                `json:"used_size"`
	FileArr  []FileUploadStructure `json:"file_list"`
}

type FileChunkStructure struct {
	ID       uuid.UUID `json:"chunk_id"`
	FileID   uuid.UUID `json:"file_id"`
	CheckSum string    `json:"check_sum"`
	Order    int       `json:"order"`
}

type TokenData struct {
	ID          string `json:"sub"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Picture     string `json:"picture"`
	AllotedSize uint64 `json:"alloted_size"`
}

func (m *FileList) CalculateSize() uint64 {
	var totalSize uint64
	for idx := range m.FileArr {
		totalSize += m.FileArr[idx].Meta.Size
	}

	m.UsedSize = totalSize

	return totalSize
}
