package utils

import (
    "context"
    "io/ioutil"

    "google.golang.org/api/docs/v1"
    "google.golang.org/api/drive/v3"
    "google.golang.org/api/option"
	"golang.org/x/oauth2/google"
	
)

type GoogleDriveHelper struct {
    ctx    context.Context
    srv    *drive.Service
    docsSrv *docs.Service
}

func NewGoogleDriveHelper(credentialsFilePath string) (*GoogleDriveHelper, error) {
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
        ctx:    ctx,
        srv:    srv,
        docsSrv: docsSrv,
    }, nil
}

func (h *GoogleDriveHelper) GetFileContent(fileID string) (string, error) {
    doc, err := h.docsSrv.Documents.Get(fileID).Do()
    if err != nil {
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
