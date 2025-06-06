package sqlite

import (
	"unterlagen/features/archive"

	"github.com/jmoiron/sqlx"
)

var _ archive.DocumentRepository = &DocumentRepository{}

type DocumentRepository struct {
	*sqlx.DB
}

// DeleteByID implements archive.DocumentRepository.
func (d *DocumentRepository) DeleteByID(id string) error {
	_, err := d.Exec("DELETE FROM documents WHERE id = ?", id)
	return err
}

// FindAllByFolderID implements archive.DocumentRepository.
func (d *DocumentRepository) FindAllByFolderID(folderID string) ([]archive.Document, error) {
	var documents []archive.Document
	err := d.Select(&documents, "SELECT * FROM documents WHERE folder_id = ?", folderID)
	if err != nil {
		return nil, err
	}
	return documents, nil
}

// FindAllTrashed implements archive.DocumentRepository.
func (d *DocumentRepository) FindAllTrashed() ([]archive.Document, error) {
	var documents []archive.Document
	err := d.Select(&documents, "SELECT * FROM documents WHERE trashed_at IS NOT NULL")
	if err != nil {
		return nil, err
	}
	return documents, nil
}

// FindByID implements archive.DocumentRepository.
func (d *DocumentRepository) FindByID(id string) (archive.Document, error) {
	var document archive.Document
	err := d.Get(&document, "SELECT * FROM documents WHERE id = ?", id)
	if err != nil {
		return archive.Document{}, err
	}
	return document, nil
}

// Save implements archive.DocumentRepository.
func (d *DocumentRepository) Save(document archive.Document) error {
	_, err := d.NamedExec(`
		INSERT INTO documents (id, filename, filetype, filesize, text, folder_id, owner, created_at, updated_at, trashed_at)
		VALUES (:id, :filename, :filetype, :filesize, :text, :folder_id, :owner, :created_at, :updated_at, :trashed_at)
		ON CONFLICT(id) DO UPDATE SET
			filename = :filename,
			filetype = :filetype,
			filesize = :filesize,
			text = :text,
			folder_id = :folder_id,
			owner = :owner,
			updated_at = :updated_at,
			trashed_at = :trashed_at
	`, document)
	return err
}

func NewDocumentRepository(db *sqlx.DB) *DocumentRepository {
	return &DocumentRepository{
		db,
	}
}
