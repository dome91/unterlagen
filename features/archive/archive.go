package archive

import (
	"unterlagen/features/administration"
	"unterlagen/features/common"
)

type Archive struct {
	*documents
	*folders
}

func (a *Archive) Synchronize(owner string) error {
	return a.rescheduleAllDocumentTasks(owner)
}

func New(
	documentRepository DocumentRepository,
	documentStorage DocumentStorage,
	documentPreviewStorage DocumentPreviewStorage,
	documentMessages DocumentMessages,
	folderRepository FolderRepository,
	userMessages administration.UserMessages,
	jobScheduler *common.JobScheduler,
	taskScheduler *common.TaskScheduler,
	shutdown *common.Shutdown,
) *Archive {
	return &Archive{
		documents: newDocuments(
			documentRepository,
			documentStorage,
			documentPreviewStorage,
			documentMessages,
			jobScheduler,
			taskScheduler,
			shutdown,
		),
		folders: newFolders(
			folderRepository,
			userMessages,
		),
	}
}
