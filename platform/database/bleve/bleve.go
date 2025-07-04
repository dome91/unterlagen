package bleve

import (
	"log/slog"
	"os"
	"path/filepath"
	"unterlagen/features/common"
	"unterlagen/platform/configuration"

	"github.com/blevesearch/bleve/v2"
)

func Initialize(configuration *configuration.Configuration, shutdown *common.Shutdown) bleve.Index {
	var err error
	var index bleve.Index
	if configuration.Production {
		err = os.MkdirAll(configuration.Data.Directory, os.ModePerm)
		if err != nil {
			panic(err)
		}

		indexPath := filepath.Join(configuration.Data.Directory, "unterlagen.bleve")
		if _, err := os.Stat(indexPath); err == nil {
			index, err = bleve.Open(indexPath)
			if err != nil {
				panic(err)
			}
		} else {
			indexMapping := bleve.NewIndexMapping()
			index, err = bleve.New(indexPath, indexMapping)
			if err != nil {
				panic(err)
			}
		}
	} else {
		indexMapping := bleve.NewIndexMapping()
		index, err = bleve.NewMemOnly(indexMapping)
		if err != nil {
			panic(err)
		}
	}

	shutdown.AddCallback(func() {
		err := index.Close()
		if err != nil {
			slog.Warn("Failed to close index", "error", err)
		}
	})

	return index
}
