package memory

import (
	"fmt"
	"sync"
	"unterlagen/features/archive"
)

var _ archive.DocumentRepository = &DocumentRepository{}

type DocumentRepository struct {
	documents map[string]archive.Document
	mutex     sync.RWMutex
}

func (r *DocumentRepository) Save(document archive.Document) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.documents[document.ID] = document
	return nil
}

func (r *DocumentRepository) FindByID(id string) (archive.Document, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	document, exists := r.documents[id]
	if !exists {
		return archive.Document{}, fmt.Errorf("document not found")
	}
	return document, nil
}

func (r *DocumentRepository) FindAllByFolderID(folderID string) ([]archive.Document, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	var documents []archive.Document
	for _, document := range r.documents {
		if document.FolderID == folderID {
			documents = append(documents, document)
		}
	}
	return documents, nil
}

func (r *DocumentRepository) FindAllTrashed() ([]archive.Document, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	var documents []archive.Document
	for _, document := range r.documents {
		if document.TrashedAt.Valid {
			documents = append(documents, document)
		}
	}
	return documents, nil
}

func (r *DocumentRepository) DeleteByID(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.documents, id)
	return nil
}

func NewDocumentRepository() *DocumentRepository {
	return &DocumentRepository{
		documents: make(map[string]archive.Document),
	}
}
