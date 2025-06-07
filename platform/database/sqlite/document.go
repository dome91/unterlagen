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

	// Load preview filepaths for each document
	for i := range documents {
		previews, err := d.loadPreviewFilepaths(documents[i].ID)
		if err != nil {
			return nil, err
		}
		documents[i].PreviewFilepaths = previews
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

	// Load preview filepaths for each document
	for i := range documents {
		previews, err := d.loadPreviewFilepaths(documents[i].ID)
		if err != nil {
			return nil, err
		}
		documents[i].PreviewFilepaths = previews
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

	// Load preview filepaths
	previews, err := d.loadPreviewFilepaths(id)
	if err != nil {
		return archive.Document{}, err
	}
	document.PreviewFilepaths = previews

	return document, nil
}

// Save implements archive.DocumentRepository.
func (d *DocumentRepository) Save(document archive.Document) error {
	tx, err := d.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Save document
	_, err = tx.NamedExec(`
		INSERT INTO documents (id, filename, filetype, filesize, text, folder_id, owner, created_at, updated_at, trashed_at)
		VALUES (:id, :filename, :filetype, :filesize, :text, :folder_id, :owner, :created_at, :updated_at, :trashed_at)
		ON CONFLICT(id) DO UPDATE SET
			filename = :filename,
			filetype = :filetype,
			filesize = :filesize,
			text = :text,
			folder_id = :folder_id,
			owner = :owner,
			updated_at = datetime(),
			trashed_at = :trashed_at
	`, document)
	if err != nil {
		return err
	}

	// Save preview filepaths
	err = d.savePreviewFilepaths(tx, document.ID, document.PreviewFilepaths)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (d *DocumentRepository) loadPreviewFilepaths(documentID string) ([]string, error) {
	rows, err := d.Query(`
		SELECT filepath FROM documents_previews 
		WHERE document_id = ? 
		ORDER BY page_number ASC
	`, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var filepaths []string
	for rows.Next() {
		var filepath string
		if err := rows.Scan(&filepath); err != nil {
			return nil, err
		}
		filepaths = append(filepaths, filepath)
	}

	return filepaths, rows.Err()
}

func (d *DocumentRepository) savePreviewFilepaths(tx *sqlx.Tx, documentID string, filepaths []string) error {
	// Use INSERT ... ON CONFLICT for each preview filepath
	for i, filepath := range filepaths {
		_, err := tx.Exec(`
			INSERT INTO documents_previews (document_id, filepath, page_number, created_at)
			VALUES (?, ?, ?, datetime())
			ON CONFLICT(document_id, page_number) DO UPDATE SET
				filepath = excluded.filepath,
				created_at = datetime()
		`, documentID, filepath, i)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewDocumentRepository(db *sqlx.DB) *DocumentRepository {
	return &DocumentRepository{
		db,
	}
}
