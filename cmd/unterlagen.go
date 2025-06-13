package main

import (
	"log/slog"
	_ "net/http/pprof"
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
	jobScheduler := common.NewJobScheduler(shutdown)

	configuration := configuration.Load()

	// Database
	db := sqlite.Initialize(shutdown, jobScheduler, configuration)
	userRepository := sqlite.NewUserRepository(db)
	documentRepository := sqlite.NewDocumentRepository(db)
	folderRepository := sqlite.NewFolderRepository(db)
	taskRepository := sqlite.NewTaskRepository(db)
	settingsRepository := memory.NewSettingsRepository()

	// Messaging
	userMessages := synchronous.NewUserMessages()
	documentMessages := synchronous.NewDocumentMessages()

	// Storage
	documentStorage := filesystem.NewDocumentStorage(configuration)
	documentPreviewStorage := filesystem.NewDocumentPreviewStorage(configuration)

	// Features
	taskScheduler := common.NewTaskScheduler(shutdown, taskRepository)
	administration := administration.New(settingsRepository, userRepository, userMessages, taskRepository)
	archive := archive.New(documentRepository, documentStorage, documentPreviewStorage, documentMessages, folderRepository, userMessages, jobScheduler, taskScheduler, shutdown)

	// Web
	server := web.NewServer(administration, archive, shutdown, configuration)
	server.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	shutdown.Execute()
	slog.Info("unterlagen stopped. Bye!")
}
