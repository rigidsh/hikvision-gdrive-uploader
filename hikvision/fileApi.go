package hikvision

import (
	"io"
	"time"
)

type File struct {
	Id   string
	Name string
}

type FileAPI interface {
	GetFilesCreatedBefore(beforeTime time.Time) ([]File, error)
	GetFilesCreatedAfter(afterTime time.Time) ([]File, error)
	DeleteFile(file File) error
	GetFileContent(file File) (io.ReadCloser, error)
	CreateFile(name string, content io.Reader) (File, error)
}
