package search

import (
	"log/slog"
	"slices"
	"strings"
	"unterlagen/features/archive"
	"unterlagen/features/common"
)

type SearchResult struct {
	DocumentID string
	Name       string
	Rank       float64
	Snippet    string
}

type SearchRepository interface {
	IndexDocument(document archive.Document) error
	SearchDocuments(query string, owner string, limit int) ([]SearchResult, error)
}

type Search struct {
	repository SearchRepository
}

func (s *Search) SearchDocuments(query string, owner string, limit int) ([]SearchResult, error) {
	// Sanitize and validate query
	query = strings.TrimSpace(query)
	if query == "" {
		return []SearchResult{}, nil
	}

	// Limit results to prevent excessive load
	if limit <= 0 || limit > 50 {
		limit = 50
	}

	results, err := s.repository.SearchDocuments(query, owner, limit)
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		slog.Info(result.Name, "rank", result.Rank)
	}
	slices.SortFunc(results, func(r1, r2 SearchResult) int {
		if r1.Rank > r2.Rank {
			return -1
		} else if r1.Rank < r2.Rank {
			return 1
		}
		return 0
	})

	return results, err
}

func New(repository SearchRepository, documentMessages archive.DocumentMessages, taskScheduler *common.TaskScheduler) *Search {
	taskProcessor := NewSearchTaskProcessor(repository)
	taskScheduler.Register(taskProcessor)

	err := documentMessages.SubscribeDocumentTextExtracted(func(document archive.Document) error {
		return taskScheduler.ScheduleTask(common.TaskTypeIndexDocument, document, 3)
	})
	if err != nil {
		panic(err)
	}

	return &Search{repository: repository}
}
