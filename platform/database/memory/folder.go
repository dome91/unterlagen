package memory

import (
	"sync"
	"unterlagen/features/archive"
)

var _ archive.FolderRepository = &FolderRepository{}

type FolderRepository struct {
	folders map[string]archive.Folder
	mutex   sync.RWMutex
}

func (r *FolderRepository) Save(folder archive.Folder) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.folders[folder.ID] = folder
	return nil
}

func (r *FolderRepository) FindAllByParentID(parentID string) ([]archive.Folder, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	var folders []archive.Folder
	for _, folder := range r.folders {
		if folder.ParentID == parentID {
			folders = append(folders, folder)
		}
	}
	return folders, nil
}

func (r *FolderRepository) GetHierarchy(folderID string) ([]archive.Folder, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	var hierarchy []archive.Folder

	currentID := folderID
	for currentID != "" && currentID != archive.FolderRootID {
		folder, exists := r.folders[currentID]
		if !exists {
			break
		}
		hierarchy = append([]archive.Folder{folder}, hierarchy...)
		currentID = folder.ParentID
	}

	if currentID == archive.FolderRootID {
		rootFolder, exists := r.folders[archive.FolderRootID]
		if exists {
			hierarchy = append([]archive.Folder{rootFolder}, hierarchy...)
		}
	}

	return hierarchy, nil
}

func NewFolderRepository() *FolderRepository {
	return &FolderRepository{
		folders: make(map[string]archive.Folder),
	}
}
