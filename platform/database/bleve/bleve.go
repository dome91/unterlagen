package bleve

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"unterlagen/platform/configuration"

	"github.com/blevesearch/bleve/v2"
)

func Initialize(configuration *configuration.Configuration) {
	var index bleve.Index
	var err error
	indexPath := filepath.Join(configuration.Data.Directory, "unterlagen.bleve")
	if configuration.Production {
		err = os.MkdirAll(configuration.Data.Directory, os.ModePerm)
		if err != nil {
			panic(err)
		}

		indexMapping := bleve.NewIndexMapping()
		index, err = bleve.New(filepath.Join(configuration.Data.Directory, "unterlagen.bleve"), indexMapping)
		if err != nil {
			panic(err)
		}
	} else {
	}

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		// Create new index
		indexMapping := createIndexMapping()
		index, err = bleve.New(indexPath, indexMapping)
		if err != nil {
			return nil, fmt.Errorf("failed to create new search index: %w", err)
		}
		slog.Info("created new search index", "path", indexPath)
	} else {
		// Open existing index
		index, err = bleve.Open(indexPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open existing search index: %w", err)
		}
		slog.Info("opened existing search index", "path", indexPath)
	}
}
