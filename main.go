package main

import (
	"hikvision-gdrive-uploader/hikvision"
	"os"
	"time"
)

func main() {

	gdriveRoot := os.Getenv("HIKVISION_GDRIVE_ROOT")
	dataPath := os.Getenv("HIKVISION_DATA")
	credentialsPath := os.Getenv("HIKVISION_GOOGLE_API_CREDENTIALS")
	tokenPath := os.Getenv("HIKVISION_GOOGLE_API_TOKEN")
	uploadInterval, _ := time.ParseDuration(os.Getenv("HIKVISION_UPLOAD_INTERVAL"))
	deleteAfter, _ := time.ParseDuration(os.Getenv("HIKVISION_DELETE_AFTER"))

	uploadAPI, err := hikvision.NewGDriveFileAPI(gdriveRoot, credentialsPath, tokenPath)
	if err != nil {
		return
	}
	localContentAPI := hikvision.NewLocalFolderFileAPI(dataPath)

	hikvision.StartUploadingData(localContentAPI, uploadAPI, uploadInterval, deleteAfter)

}
