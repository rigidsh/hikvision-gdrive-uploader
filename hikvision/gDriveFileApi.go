package hikvision

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"log"
	"os"
	"time"
)

type GDriveFileAPI struct {
	gDriveService *drive.Service
	root          File
}

func NewGDriveFileAPI(rootFolderName string, credentialsPath string, tokenPath string) (*GDriveFileAPI, error) {
	ctx := context.Background()
	b, err := os.ReadFile(credentialsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	tok, err := tokenFromFile(tokenPath)
	client := config.Client(context.Background(), tok)

	gDriveService, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	root, err := getOrCreateRootFolder(gDriveService, rootFolderName)
	if err != nil {
		return nil, err
	}

	return &GDriveFileAPI{
		gDriveService: gDriveService,
		root:          root,
	}, nil
}

func (fileAPI *GDriveFileAPI) GetFilesCreatedBefore(beforeTime time.Time) ([]File, error) {
	var pageToken string
	result := make([]File, 0)

	utcLocation, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, err
	}

	for {
		filesList, err := fileAPI.gDriveService.Files.List().
			Q(fmt.Sprintf("createdTime < '%s' AND '%s' in parents", beforeTime.In(utcLocation).Format(time.RFC3339), fileAPI.root.Id)).
			PageSize(100).
			PageToken(pageToken).
			Fields("nextPageToken, files(id, name)").
			Do()

		if err != nil {
			return nil, err
		}

		for _, file := range filesList.Files {
			result = append(result, File{Id: file.Id, Name: file.Name})
		}

		pageToken = filesList.NextPageToken

		if len(pageToken) == 0 {
			break
		}
	}

	return result, nil
}

func (fileAPI *GDriveFileAPI) GetFilesCreatedAfter(afterTime time.Time) ([]File, error) {
	var pageToken string
	result := make([]File, 0)

	utcLocation, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, err
	}

	for {
		filesList, err := fileAPI.gDriveService.Files.List().
			Q(fmt.Sprintf("createdTime > '%s' AND '%s' in parents", afterTime.In(utcLocation).Format(time.RFC3339), fileAPI.root.Id)).
			PageSize(100).
			PageToken(pageToken).
			Fields("nextPageToken, files(id, name)").
			Do()

		if err != nil {
			return nil, err
		}

		for _, file := range filesList.Files {
			result = append(result, File{Id: file.Id, Name: file.Name})
		}

		pageToken = filesList.NextPageToken

		if len(pageToken) == 0 {
			break
		}
	}

	return result, nil
}

func (fileAPI *GDriveFileAPI) DeleteFile(file File) error {
	return fileAPI.gDriveService.Files.Delete(file.Id).Do()
}

func (fileAPI *GDriveFileAPI) GetFileContent(file File) (io.ReadCloser, error) {
	result, err := fileAPI.gDriveService.Files.Get(file.Id).Download()
	if err != nil {
		return nil, err
	}

	return result.Body, nil
}

func (fileAPI *GDriveFileAPI) CreateFile(name string, content io.Reader) (File, error) {
	uploadFile := &drive.File{
		Name:    name,
		Parents: []string{fileAPI.root.Id},
	}

	newFile, err := fileAPI.gDriveService.Files.Create(uploadFile).Media(content).
		ProgressUpdater(func(now, size int64) { log.Printf("%d, %d\r", now, size) }).
		Do()

	if err != nil {
		return File{}, err
	}

	return File{Id: newFile.Id, Name: newFile.Name}, nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getOrCreateRootFolder(service *drive.Service, rootName string) (File, error) {
	root, err := getFileByName(service, rootName)
	if err != nil {
		return File{}, err
	}

	if root != nil {
		return *root, nil
	} else {
		return createFolderInRoot(service, rootName)
	}
}

func createFolderInRoot(service *drive.Service, rootName string) (File, error) {
	rootFile := &drive.File{
		Name:     rootName,
		MimeType: "application/vnd.google-apps.folder",
	}

	root, err := service.Files.Create(rootFile).Do()

	if err != nil {
		return File{}, err
	}

	return File{Id: root.Id, Name: root.Name}, nil
}

func getFileByName(service *drive.Service, rootName string) (*File, error) {
	result, err := service.Files.List().
		Q(fmt.Sprintf("name = '%s'", rootName)).
		PageSize(1).
		Do()
	if err != nil {
		return nil, err
	}

	if len(result.Files) == 0 {
		return nil, err
	} else {
		return &File{Id: result.Files[0].Id, Name: result.Files[0].Name}, nil
	}
}

//// Retrieve a token, saves the token, then returns the generated client.
//func getClient(config *oauth2.Config) *http.Client {
//	// The file token.json stores the user's access and refresh tokens, and is
//	// created automatically when the authorization flow completes for the first
//	// time.
//	tokFile := "token.json"
//	tok, err := tokenFromFile(tokFile)
//	if err != nil {
//		tok = getTokenFromWeb(config)
//		saveToken(tokFile, tok)
//	}
//	return config.Client(context.Background(), tok)
//}
//
//
//func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
//	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
//	fmt.Printf("Go to the following link in your browser then type the "+
//		"authorization code: \n%v\n", authURL)
//
//	var authCode string
//	if _, err := fmt.Scan(&authCode); err != nil {
//		log.Fatalf("Unable to read authorization code %v", err)
//	}
//
//	tok, err := config.Exchange(context.TODO(), authCode)
//	if err != nil {
//		log.Fatalf("Unable to retrieve token from web %v", err)
//	}
//	return tok
//}
//
//// Saves a token to a file path.
//func saveToken(path string, token *oauth2.Token) {
//	fmt.Printf("Saving credential file to: %s\n", path)
//	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
//	if err != nil {
//		log.Fatalf("Unable to cache oauth token: %v", err)
//	}
//	defer f.Close()
//	json.NewEncoder(f).Encode(token)
//}
