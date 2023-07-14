package hikvision

import (
	"fmt"
	"io"
	"time"
)

func StartUploadingData(localContentAPI FileAPI, uploadAPI FileAPI, uploadInterval time.Duration, deleteAfter time.Duration) {
	ticker := time.NewTicker(uploadInterval)
	for ; true; <-ticker.C {
		fmt.Println("Start iteration")
		filesToUpload, err := localContentAPI.GetFilesCreatedBefore(time.Now().Add(-1 * uploadInterval))
		if err != nil {
			return
		}

		fmt.Printf("Finded %d files to upload\n", len(filesToUpload))

		for _, file := range filesToUpload {
			err := uploadAndDeleteFile(file, localContentAPI, uploadAPI)
			if err != nil {
				fmt.Printf("Can't upload file %s: %s", file.Name, err)
				return
			}
		}
		fmt.Println("Uploading done")
		fmt.Println("Start deleting old files")

		filesToDelete, err := uploadAPI.GetFilesCreatedBefore(time.Now().Add(-1 * deleteAfter))
		if err != nil {
			return
		}
		fmt.Printf("Finded %d files to delete\n", len(filesToDelete))
		for _, file := range filesToDelete {
			fmt.Printf("Start deleting file %s\n", file.Name)
			err := uploadAPI.DeleteFile(file)
			if err != nil {
				return
			}
			fmt.Printf("Done deleting file %s\n", file.Name)
		}
		fmt.Println("Deleting done")
	}
}

func uploadAndDeleteFile(file File, localContentAPI FileAPI, uploadAPI FileAPI) error {
	fmt.Printf("Start uploading file %s\n", file.Name)
	fileContent, err := localContentAPI.GetFileContent(file)
	if err != nil {
		return err
	}
	defer func(fileContent io.ReadCloser) {
		fmt.Printf("Done uploading file %s\n", file.Name)
		_ = fileContent.Close()
		err = localContentAPI.DeleteFile(file)
		if err != nil {
			fmt.Printf("Can't delete local file %s: %s", file.Name, err)
		}
		fmt.Printf("File %s was deleted locally\n", file.Name)
	}(fileContent)

	_, err = uploadAPI.CreateFile(file.Name, fileContent)
	if err != nil {
		return err
	}

	return nil
}
