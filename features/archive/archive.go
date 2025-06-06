package archive

import (
	"unterlagen/features/administration"
	"unterlagen/features/common"
)

type Archive struct {
	*documents
	*folders
}

func New(
	documentRepository DocumentRepository,
	documentStorage DocumentStorage,
	documentPreviewStorage DocumentPreviewStorage,
	documentMessages DocumentMessages,
	folderRepository FolderRepository,
	userMessages administration.UserMessages,
	scheduler *common.Scheduler,
) *Archive {
	return &Archive{
		documents: newDocuments(
			documentRepository,
			documentStorage,
			documentPreviewStorage,
			documentMessages,
			scheduler,
		),
		folders: newFolders(
			folderRepository,
			userMessages,
		),
	}
}
