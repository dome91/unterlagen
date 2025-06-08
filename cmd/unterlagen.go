package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"unterlagen/features/administration"
	"unterlagen/features/archive"
	"unterlagen/features/common"
	"unterlagen/platform/configuration"
	"unterlagen/platform/database/memory"
	"unterlagen/platform/database/sqlite"
	"unterlagen/platform/messaging/synchronous"
	"unterlagen/platform/storage/filesystem"
	"unterlagen/platform/web"
)

func main() {
	// Common functionality
	shutdown := common.NewShutdown()
	scheduler := common.NewScheduler(shutdown)

	configuration := configuration.Load()

	// Database
	db := sqlite.Initialize(shutdown, scheduler, configuration)
	userRepository := sqlite.NewUserRepository(db)
	documentRepository := sqlite.NewDocumentRepository(db)
	folderRepository := sqlite.NewFolderRepository(db)
	taskRepository := sqlite.NewTaskRepository(db)
	settingsRepository := memory.NewSettingsRepository()
	
	// Set task repository on scheduler to break cyclic dependency
	scheduler.SetTaskRepository(taskRepository)

	// Event
	userMessages := synchronous.NewUserMessages()
	documentMessages := synchronous.NewDocumentMessages()

	// Storage
	documentStorage := filesystem.NewDocumentStorage(configuration)
	documentPreviewStorage := filesystem.NewDocumentPreviewStorage(configuration)

	// Features
	administration := administration.New(settingsRepository, userRepository, userMessages, taskRepository)
	archive := archive.New(documentRepository, documentStorage, documentPreviewStorage, documentMessages, folderRepository, userMessages, scheduler)

	// Web
	server := web.NewServer(administration, archive, shutdown, configuration)
	server.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	shutdown.Execute()
	slog.Info("unterlagen stopped. Bye!")
}
