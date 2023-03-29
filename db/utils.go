package db

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"

	"github.com/b0noi/go-utils/v2/fs"
)

func buildMapping() *mapping.IndexMappingImpl {
	// create a new index mapping
	mapping := bleve.NewIndexMapping()

	// create a new document mapping for ChatMessage
	docMapping := bleve.NewDocumentMapping()

	docMapping.AddFieldMappingsAt("id", bleve.NewTextFieldMapping())
	docMapping.AddFieldMappingsAt("timestamp", bleve.NewDateTimeFieldMapping())
	docMapping.AddFieldMappingsAt("sender", bleve.NewTextFieldMapping())
	docMapping.AddFieldMappingsAt("content", bleve.NewTextFieldMapping())

	// add the document mapping to the index mapping
	mapping.AddDocumentMapping("chatmessage", docMapping)

	return mapping
}

func GetIndex() (bleve.Index, error) {
	appFolder, err := fs.MaybeCreateProgramFolder("assistent")
	indexPath := filepath.Join(appFolder, "chat_messages.bleve")
	pathExists, err := fs.PathExists(indexPath)
	if err != nil {
		return nil, err
	}
	if pathExists {
		index, err := bleve.Open(indexPath)
		if err != nil {
			return nil, err
		} else if index != nil {
			return index, nil
		} else {
			return nil, errors.New("path for DB exist but I can't open it")
		}
	}
	indexMapping := buildMapping()
	index, err := bleve.New(indexPath, indexMapping)
	if err != nil {
		return nil, err
	}
	return index, nil
}

func DeleteMessage(id string) error {
	idx, err := GetIndex()
	if err != nil {
		return err
	}
	fmt.Println("here")
	err = idx.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
