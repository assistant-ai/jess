package chat

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
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
	indexMapping := buildMapping()
	index, err := bleve.New("chat_messages.bleve", indexMapping)
	if err != nil {
		return nil, err
	}
	return index, nil
}
