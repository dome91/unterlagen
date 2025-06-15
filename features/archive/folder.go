package archive

import (
	"unterlagen/features/administration"
	"unterlagen/features/common"
)

var FolderRootID = "root"

type Folder struct {
	ID       string
	Name     string
	ParentID string
	Owner    string
}

type FolderRepository interface {
	Save(folder Folder) error
	FindAllByParentID(parentID string) ([]Folder, error)
	GetHierarchy(folderID string) ([]Folder, error)
}

type folders struct {
	repository FolderRepository
}

func (f *folders) CreateFolder(name string, parentID string, owner string) error {
	return f.create(common.GenerateID(), name, parentID, owner)
}

func (f *folders) CreateRootFolderFor(user administration.User) error {
	return f.create(FolderRootID, "Root", "", user.Username)
}

func (f *folders) create(id string, name string, parentID string, owner string) error {
	folder := Folder{ID: id, Name: name, ParentID: parentID, Owner: owner}
	return f.repository.Save(folder)
}

func (f *folders) GetFolderChildren(parentID string, owner string) ([]Folder, error) {
	folders, err := f.repository.FindAllByParentID(parentID)
	if err != nil {
		return nil, err
	}

	// Filter folders to only those owned by the user
	var userFolders []Folder
	for _, folder := range folders {
		if folder.Owner == owner {
			userFolders = append(userFolders, folder)
		}
	}

	return userFolders, nil
}

func (f *folders) GetFolderHierarchy(folderID string, owner string) ([]Folder, error) {
	hierarchy, err := f.repository.GetHierarchy(folderID)
	if err != nil {
		return nil, err
	}

	// Verify that the user owns the target folder (last in hierarchy)
	if len(hierarchy) > 0 {
		targetFolder := hierarchy[len(hierarchy)-1]
		if targetFolder.Owner != owner {
			return nil, ErrNotAllowed
		}
	}

	// Filter hierarchy to only folders owned by the user
	var userHierarchy []Folder
	for _, folder := range hierarchy {
		if folder.Owner == owner {
			userHierarchy = append(userHierarchy, folder)
		}
	}

	return userHierarchy, nil
}

func (f *folders) GetFolder(folderID string, owner string) (Folder, error) {
	// Special case for root folder - construct it dynamically
	if folderID == FolderRootID {
		return Folder{
			ID:       FolderRootID,
			Name:     "Root",
			ParentID: "",
			Owner:    owner,
		}, nil
	}

	// For non-root folders, get the hierarchy to find the target folder
	hierarchy, err := f.repository.GetHierarchy(folderID)
	if err != nil {
		return Folder{}, err
	}

	// The target folder is the last in the hierarchy
	if len(hierarchy) == 0 {
		return Folder{}, ErrNotAllowed
	}

	targetFolder := hierarchy[len(hierarchy)-1]
	if targetFolder.Owner != owner {
		return Folder{}, ErrNotAllowed
	}

	return targetFolder, nil
}

func newFolders(repository FolderRepository, userMessages administration.UserMessages) *folders {
	folders := &folders{
		repository: repository,
	}

	userMessages.SubscribeUserCreated(folders.CreateRootFolderFor)
	return folders
}
