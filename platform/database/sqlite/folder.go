package sqlite

import (
	"database/sql"
	"slices"
	"unterlagen/features/archive"

	"github.com/jmoiron/sqlx"
)

var _ archive.FolderRepository = &FolderRepository{}

type FolderRepository struct {
	db *sqlx.DB
}

// Save implements archive.FolderRepository.
func (f *FolderRepository) Save(folder archive.Folder) error {
	query := `
        INSERT INTO folders (id, name, parent_id, owner, created_at, updated_at)
        VALUES (:id, :name, :parent_id, :owner, datetime(), datetime())
        ON CONFLICT (id) DO UPDATE
        SET name = :name, parent_id = :parent_id, owner = :owner, updated_at = datetime()
    `

	_, err := f.db.NamedExec(query, f.mapToEntity(folder))
	return err
}

// FindAllByParentID implements archive.FolderRepository.
func (f *FolderRepository) FindAllByParentID(parentID string) ([]archive.Folder, error) {
	var entities []sqlFolderEntity
	query := `SELECT id, name, parent_id, owner FROM folders WHERE parent_id = $1`

	err := f.db.Select(&entities, query, parentID)
	if err != nil {
		return nil, err
	}

	return f.mapToFolders(entities), nil
}

// GetHierarchy implements archive.FolderRepository.
func (f *FolderRepository) GetHierarchy(folderID string) ([]archive.Folder, error) {
	var entities []sqlFolderEntity

	query := `
        WITH RECURSIVE folder_hierarchy AS (
            SELECT id, name, parent_id, owner, created_at, updated_at FROM folders WHERE id = $1
            UNION ALL
            SELECT f.id, f.name, f.parent_id, f.owner, f.created_at, f.updated_at FROM folders f
            INNER JOIN folder_hierarchy fh ON f.id = fh.parent_id
        )
        SELECT id, name, parent_id, owner, created_at, updated_at FROM folder_hierarchy
        ORDER BY id ASC;
    `

	err := f.db.Select(&entities, query, folderID)
	slices.Reverse(entities)
	return f.mapToFolders(entities), err
}

func (f *FolderRepository) mapToFolders(entities []sqlFolderEntity) []archive.Folder {
	var folders []archive.Folder
	for _, entity := range entities {
		folders = append(folders, entity.ToFolder())
	}

	return folders
}

func (f *FolderRepository) mapToEntity(folder archive.Folder) sqlFolderEntity {
	return sqlFolderEntity{
		ID:   folder.ID,
		Name: folder.Name,
		ParentID: sql.NullString{
			String: folder.ParentID,
			Valid:  len(folder.ParentID) > 0,
		},
		Owner: folder.Owner,
	}
}

func NewFolderRepository(db *sqlx.DB) *FolderRepository {
	return &FolderRepository{db: db}
}

type sqlFolderEntity struct {
	ID        string         `db:"id"`
	Name      string         `db:"name"`
	ParentID  sql.NullString `db:"parent_id"`
	Owner     string         `db:"owner"`
	CreatedAt string         `db:"created_at"`
	UpdatedAt string         `db:"updated_at"`
}

func (f *sqlFolderEntity) ToFolder() archive.Folder {
	return archive.Folder{
		ID:       f.ID,
		Name:     f.Name,
		ParentID: f.ParentID.String,
		Owner:    f.Owner,
	}
}
