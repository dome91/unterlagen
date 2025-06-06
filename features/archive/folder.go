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

func (f *folders) GetFolderChildren(parentID string) ([]Folder, error) {
	return f.repository.FindAllByParentID(parentID)
}

func (f *folders) GetFolderHierarchy(folderID string) ([]Folder, error) {
	hierarchy, err := f.repository.GetHierarchy(folderID)
	return hierarchy, err
}

func newFolders(repository FolderRepository, userMessages administration.UserMessages) *folders {
	folders := &folders{
		repository: repository,
	}

	userMessages.SubscribeUserCreated(folders.CreateRootFolderFor)
	return folders
}
