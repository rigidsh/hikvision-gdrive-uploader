package hikvision

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

type LocalFolderFileAPI struct {
	rootName string
}

func NewLocalFolderFileAPI(root string) *LocalFolderFileAPI {
	return &LocalFolderFileAPI{
		rootName: root,
	}
}

func (fileAPI *LocalFolderFileAPI) GetFilesCreatedBefore(beforeTime time.Time) ([]File, error) {
	result := make([]File, 0)

	filesList, err := os.ReadDir(fileAPI.rootName)
	if err != nil {
		return nil, err
	}

	for _, file := range filesList {
		filesStat, err := os.Stat(filepath.Join(fileAPI.rootName, file.Name()))
		if err != nil {
			return nil, err
		}

		if filesStat.ModTime().Before(beforeTime) {
			result = append(result, File{Id: filesStat.Name(), Name: filesStat.Name()})
		}
	}

	return result, nil
}

func (fileAPI *LocalFolderFileAPI) GetFilesCreatedAfter(afterTime time.Time) ([]File, error) {
	result := make([]File, 0)

	filesList, err := os.ReadDir(fileAPI.rootName)
	if err != nil {
		return nil, err
	}

	for _, file := range filesList {
		filesStat, err := os.Stat(file.Name())
		if err != nil {
			return nil, err
		}

		if filesStat.ModTime().After(afterTime) {
			result = append(result, File{Id: filesStat.Name(), Name: filesStat.Name()})
		}
	}

	return result, nil
}

func (fileAPI *LocalFolderFileAPI) DeleteFile(file File) error {
	return os.Remove(filepath.Join(fileAPI.rootName, file.Id))
}
func (fileAPI *LocalFolderFileAPI) GetFileContent(file File) (io.ReadCloser, error) {
	return os.Open(filepath.Join(fileAPI.rootName, file.Id))
}
func (fileAPI *LocalFolderFileAPI) CreateFile(name string, content io.Reader) (File, error) {
	newFile, err := os.Create(filepath.Join(fileAPI.rootName, name))
	if err != nil {
		return File{}, err
	}

	_, err = io.Copy(newFile, content)
	if err != nil {
		return File{}, err
	}

	return File{Id: name, Name: name}, nil
}
