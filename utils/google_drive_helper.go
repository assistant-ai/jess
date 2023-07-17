package utils

import (
	"context"
	"fmt"
	googleapi "google.golang.org/api/googleapi"
	"io/ioutil"
	"os"
	"regexp"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GoogleDriveHelper struct {
	ctx     context.Context
	srv     *drive.Service
	docsSrv *docs.Service
}

func GetGoogleDocIDFromURL(url string) (string, error) {
	re := regexp.MustCompile(`\/d\/(.+?)\/`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("Could not extract google document ID from URL")
	}
	return matches[1], nil
}

func getGoogleIdFromInput(inputString string) (string, error) {
	urlPattern := `^https?://.*$`
	if matched, _ := regexp.MatchString(urlPattern, inputString); matched {
		parserd_id, err := GetGoogleDocIDFromURL(inputString)
		return parserd_id, err
	}
	return inputString, nil
}

func isGoogleDocsDocumentID(url string) bool {
	googleDocsURLPattern := `^https?://docs.google.com/document/d/[-_a-zA-Z0-9]+/.*$`
	regex := regexp.MustCompile(googleDocsURLPattern)
	return regex.MatchString(url)
}

func NewGoogleDriveHelper(credentialsFilePath string) (*GoogleDriveHelper, error) {
	checkPassed, err := googleHelperChecks(credentialsFilePath)
	if err != nil {
	}
	if !checkPassed {
		PrintlnRed("Google drive helper checks failed")
		os.Exit(1)
	}

	ctx := context.Background()
	data, err := ioutil.ReadFile(credentialsFilePath)
	if err != nil {
		return nil, err
	}

	conf, err := google.JWTConfigFromJSON(data, drive.DriveReadonlyScope, docs.DocumentsReadonlyScope)
	if err != nil {
		return nil, err
	}

	srv, err := drive.NewService(ctx, option.WithHTTPClient(conf.Client(ctx)))
	if err != nil {
		return nil, err
	}

	docsSrv, err := docs.NewService(ctx, option.WithHTTPClient(conf.Client(ctx)))
	if err != nil {
		return nil, err
	}

	return &GoogleDriveHelper{
		ctx:     ctx,
		srv:     srv,
		docsSrv: docsSrv,
	}, nil
}

func checkAccessError(FileID string, error error) {
	googleAPIError, ok := error.(*googleapi.Error)
	if ok {
		msg := fmt.Sprintf("Google API Error Code: %v", googleAPIError.Code)
		PrintlnRed(msg)
		if googleAPIError.Code == 403 {
			msgDescription := "Google API Error Message: " + googleAPIError.Message
			PrintlnRed(msgDescription)
			msgFileIdrecommendation := "Jessica don't have access to provided document with ID: " + FileID
			PrintlnRed(msgFileIdrecommendation)
			msgReccomendationFixAccess := "Try to provide read access to this document for Jessica using this service account: " + GetConfig().JessEmailAccount
			PrintlnYellow(msgReccomendationFixAccess)
			os.Exit(1)
		}
	} else {
		fmt.Println("Error is not a Google API Error.")
	}
}

func (h *GoogleDriveHelper) GetFileContent(fileID string) (string, error) {
	if isValidURL(fileID) {
		realGoogleDocument := isGoogleDocsDocumentID(fileID)
		if !realGoogleDocument {
			errorMsg := fmt.Sprintf("Provided document is not a google document: %v", fileID)
			PrintlnRed(errorMsg)
			os.Exit(1)
			return fileID, nil
		}
	}

	fileID, err := getGoogleIdFromInput(fileID)
	if err != nil {
		checkAccessError(fileID, err)
		return "", err
	}

	doc, err := h.docsSrv.Documents.Get(fileID).Do()
	if err != nil {
		checkAccessError(fileID, err)
		return "", err
	}

	var content string
	for _, cs := range doc.Body.Content {
		if cs.Paragraph != nil {
			for _, pe := range cs.Paragraph.Elements {
				if pe.TextRun != nil {
					content += pe.TextRun.Content
				}
			}
		}
	}

	return content, nil
}

func googleHelperChecks(serviceAccountPath string) (bool, error) {

	serviceAccountFileExist, err := IsServiceAccountJsonFileExists(serviceAccountPath)
	if err != nil {
		return false, err
	}
	if !serviceAccountFileExist {
		PrintlnRed("Service account file not found or it is empty")
		msgRecommendation := "Please create service account json file and save it to " + serviceAccountPath
		PrintlnYellow(msgRecommendation)
		PrintlnYellow("find more info about google service account here : https://cloud.google.com/iam/docs/service-account-overview")
		os.Exit(1)
	}

	return true, nil
}
