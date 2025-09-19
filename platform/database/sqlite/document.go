package sqlite

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"
	"unterlagen/features/archive"

	"github.com/jmoiron/sqlx"
)

var _ archive.DocumentRepository = &DocumentRepository{}

// DocumentEntity represents a document in the database layer
type DocumentEntity struct {
	ID        string       `db:"id"`
	Title     string       `db:"title"`
	Filename  string       `db:"filename"`
	Filetype  string       `db:"filetype"`
	Filesize  uint64       `db:"filesize"`
	Text      string       `db:"text"`
	Summary   []byte       `db:"summary"` // JSON stored as bytes
	FolderID  string       `db:"folder_id"`
	Owner     string       `db:"owner"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	TrashedAt sql.NullTime `db:"trashed_at"`
}

// to converts DocumentEntity to archive.Document
func (entity *DocumentEntity) to() (archive.Document, error) {
	// Deserialize summary
	var summary archive.DocumentSummary
	if len(entity.Summary) > 0 {
		err := json.Unmarshal(entity.Summary, &summary)
		if err != nil {
			return archive.Document{}, err
		}
	}

	return archive.Document{
		ID:        entity.ID,
		Title:     entity.Title,
		Filename:  entity.Filename,
		Filetype:  archive.Filetype(entity.Filetype),
		Filesize:  entity.Filesize,
		Text:      entity.Text,
		Summary:   summary,
		FolderID:  entity.FolderID,
		Owner:     entity.Owner,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		TrashedAt: entity.TrashedAt,
	}, nil
}

// from converts archive.Document to DocumentEntity
func (entity *DocumentEntity) from(doc archive.Document) error {
	summaryData, err := json.Marshal(doc.Summary)
	if err != nil {
		return err
	}

	*entity = DocumentEntity{
		ID:        doc.ID,
		Title:     doc.Title,
		Filename:  doc.Filename,
		Filetype:  string(doc.Filetype),
		Filesize:  doc.Filesize,
		Text:      doc.Text,
		Summary:   summaryData,
		FolderID:  doc.FolderID,
		Owner:     doc.Owner,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
		TrashedAt: doc.TrashedAt,
	}

	return nil
}

type DocumentRepository struct {
	*sqlx.DB
}

// FindAllByIDIn implements archive.DocumentRepository.
func (d *DocumentRepository) FindAllByIDIn(ids []string) ([]archive.Document, error) {
	var query string
	var args []any

	if len(ids) == 0 {
		return []archive.Document{}, nil
	}

	// Create placeholders for the IN clause
	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = "?"
		args = append(args, ids[i])
	}

	// Build the query with the IN clause
	query = "SELECT * FROM documents WHERE id IN (?" + strings.Repeat(",?", len(ids)-1) + ")"

	// Execute the query
	var entities []DocumentEntity
	err := d.Select(&entities, query, args...)
	if err != nil {
		return nil, err
	}

	// Convert entities to domain objects
	var documents []archive.Document
	for _, entity := range entities {
		document, err := entity.to()
		if err != nil {
			return nil, err
		}

		// Load preview filepaths
		previews, err := d.loadPreviewFilepaths(document.ID)
		if err != nil {
			return nil, err
		}
		document.PreviewFilepaths = previews

		documents = append(documents, document)
	}

	return documents, nil
}

// FindAllByOwner implements archive.DocumentRepository.
func (d *DocumentRepository) FindAllByOwner(owner string) ([]archive.Document, error) {
	var entities []DocumentEntity
	err := d.Select(&entities, "SELECT * FROM documents WHERE owner = ?", owner)
	if err != nil {
		return nil, err
	}

	// Convert entities to domain objects
	var documents []archive.Document
	for _, entity := range entities {
		document, err := entity.to()
		if err != nil {
			return nil, err
		}

		// Load preview filepaths
		previews, err := d.loadPreviewFilepaths(document.ID)
		if err != nil {
			return nil, err
		}
		document.PreviewFilepaths = previews

		documents = append(documents, document)
	}

	return documents, nil
}

// DeleteByID implements archive.DocumentRepository.
func (d *DocumentRepository) DeleteByID(id string) error {
	_, err := d.Exec("DELETE FROM documents WHERE id = ?", id)
	return err
}

// FindAllByFolderID implements archive.DocumentRepository.
func (d *DocumentRepository) FindAllByFolderID(folderID string) ([]archive.Document, error) {
	var entities []DocumentEntity
	err := d.Select(&entities, "SELECT * FROM documents WHERE folder_id = ?", folderID)
	if err != nil {
		return nil, err
	}

	// Convert entities to domain objects
	var documents []archive.Document
	for _, entity := range entities {
		document, err := entity.to()
		if err != nil {
			return nil, err
		}

		// Load preview filepaths
		previews, err := d.loadPreviewFilepaths(document.ID)
		if err != nil {
			return nil, err
		}
		document.PreviewFilepaths = previews

		documents = append(documents, document)
	}

	return documents, nil
}

// FindAllTrashed implements archive.DocumentRepository.
func (d *DocumentRepository) FindAllTrashed() ([]archive.Document, error) {
	var entities []DocumentEntity
	err := d.Select(&entities, "SELECT * FROM documents WHERE trashed_at IS NOT NULL")
	if err != nil {
		return nil, err
	}

	// Convert entities to domain objects
	var documents []archive.Document
	for _, entity := range entities {
		document, err := entity.to()
		if err != nil {
			return nil, err
		}

		// Load preview filepaths
		previews, err := d.loadPreviewFilepaths(document.ID)
		if err != nil {
			return nil, err
		}
		document.PreviewFilepaths = previews

		documents = append(documents, document)
	}

	return documents, nil
}

// FindByID implements archive.DocumentRepository.
func (d *DocumentRepository) FindByID(id string) (archive.Document, error) {
	var entity DocumentEntity
	err := d.Get(&entity, "SELECT * FROM documents WHERE id = ?", id)
	if err != nil {
		return archive.Document{}, err
	}

	// Convert entity to domain
	document, err := entity.to()
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

	// Convert domain to entity
	var entity DocumentEntity
	err = entity.from(document)
	if err != nil {
		return err
	}

	// Save document using NamedExec for cleaner code
	_, err = tx.NamedExec(`
		INSERT INTO documents (id, title, filename, filetype, filesize, text, summary, folder_id, owner, created_at, updated_at, trashed_at)
		VALUES (:id, :title, :filename, :filetype, :filesize, :text, :summary, :folder_id, :owner, :created_at, :updated_at, :trashed_at)
		ON CONFLICT(id) DO UPDATE SET
			title = excluded.title,
			filename = excluded.filename,
			filetype = excluded.filetype,
			filesize = excluded.filesize,
			text = excluded.text,
			summary = excluded.summary,
			folder_id = excluded.folder_id,
			owner = excluded.owner,
			updated_at = datetime(),
			trashed_at = excluded.trashed_at
	`, entity)
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
