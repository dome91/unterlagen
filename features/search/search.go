package search

import (
	"encoding/json"
	"slices"
	"strings"
	"unterlagen/features/archive"
	"unterlagen/features/common"
)

type SearchResult struct {
	DocumentID string
	Name       string
	Snippet    string
	Rank       float64
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
		limit = 20
	}

	results, err := s.repository.SearchDocuments(query, owner, limit)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(results, func(a, b SearchResult) int {
		if a.Rank > b.Rank {
			return -1
		} else if a.Rank < b.Rank {
			return 1
		}
		return 0
	})

	return results, err
}

func (s *Search) ProcessTask(task common.Task) error {
	switch task.Type {
	case common.TaskTypeIndexDocument:
		var document archive.Document
		if err := json.Unmarshal(task.Payload, &document); err != nil {
			return err
		}
		return s.repository.IndexDocument(document)
	default:
		return nil

	}
}

func (s *Search) ResponsibleFor() []common.TaskType {
	return []common.TaskType{common.TaskTypeIndexDocument}
}

func New(repository SearchRepository, documentMessages archive.DocumentMessages, taskScheduler *common.TaskScheduler) *Search {
	search := &Search{repository: repository}
	taskScheduler.RegisterWorker(search)
	documentMessages.SubscribeDocumentTextExtracted(func(document archive.Document) error {
		return taskScheduler.ScheduleTask(common.TaskTypeIndexDocument, document, 3)
	})

	return search
}
