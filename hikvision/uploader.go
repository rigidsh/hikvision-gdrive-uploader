package hikvision

import (
	"io"
	"log"
	"time"
)

func StartUploadingData(localContentAPI FileAPI, uploadAPI FileAPI, uploadInterval time.Duration, deleteAfter time.Duration) {
	ticker := time.NewTicker(uploadInterval)
	for ; true; <-ticker.C {
		log.Println("Start iteration")
		filesToUpload, err := localContentAPI.GetFilesCreatedBefore(time.Now().Add(-1 * uploadInterval))
		if err != nil {
			return
		}

		log.Printf("Finded %d files to upload\n", len(filesToUpload))

		for _, file := range filesToUpload {
			err := uploadAndDeleteFile(file, localContentAPI, uploadAPI)
			if err != nil {
				log.Printf("Can't upload file %s: %s", file.Name, err)
				return
			}
		}
		log.Println("Uploading done")
		log.Println("Start deleting old files")

		filesToDelete, err := uploadAPI.GetFilesCreatedBefore(time.Now().Add(-1 * deleteAfter))
		if err != nil {
			return
		}
		log.Printf("Finded %d files to delete\n", len(filesToDelete))
		for _, file := range filesToDelete {
			log.Printf("Start deleting file %s\n", file.Name)
			err := uploadAPI.DeleteFile(file)
			if err != nil {
				return
			}
			log.Printf("Done deleting file %s\n", file.Name)
		}
		log.Println("Deleting done")
	}
}

func uploadAndDeleteFile(file File, localContentAPI FileAPI, uploadAPI FileAPI) error {
	log.Printf("Start uploading file %s\n", file.Name)
	fileContent, err := localContentAPI.GetFileContent(file)
	if err != nil {
		return err
	}
	defer func(fileContent io.ReadCloser) {
		log.Printf("Done uploading file %s\n", file.Name)
		_ = fileContent.Close()
		err = localContentAPI.DeleteFile(file)
		if err != nil {
			log.Printf("Can't delete local file %s: %s", file.Name, err)
		}
		log.Printf("File %s was deleted locally\n", file.Name)
	}(fileContent)

	_, err = uploadAPI.CreateFile(file.Name, fileContent)
	if err != nil {
		return err
	}

	return nil
}
