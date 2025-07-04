package search

import (
	"encoding/json"
	"unterlagen/features/archive"
	"unterlagen/features/common"
)

type SearchTaskProcessor struct {
	repository SearchRepository
}

func (p *SearchTaskProcessor) Name() string {
	return "SearchTaskProcessor"
}

func (p *SearchTaskProcessor) ProcessTask(task common.Task) error {
	switch task.Type {
	case common.TaskTypeIndexDocument:
		var document archive.Document
		if err := json.Unmarshal(task.Payload, &document); err != nil {
			return err
		}
		return p.repository.IndexDocument(document)
	default:
		return nil
	}
}

func (p *SearchTaskProcessor) ResponsibleFor() []common.TaskType {
	return []common.TaskType{common.TaskTypeIndexDocument}
}

func NewSearchTaskProcessor(repository SearchRepository) *SearchTaskProcessor {
	return &SearchTaskProcessor{
		repository: repository,
	}
}
