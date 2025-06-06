package test

import (
	"log/slog"
	"os"
	"testing"
	"time"
	"unterlagen/features/administration"
	"unterlagen/features/archive"
	"unterlagen/features/common"
	"unterlagen/platform/configuration"
	"unterlagen/platform/database/memory"
	"unterlagen/platform/database/sqlite"
	"unterlagen/platform/messaging/synchronous"
	"unterlagen/platform/storage/filesystem"
	"unterlagen/platform/web"

	"github.com/playwright-community/playwright-go"
)

var (
	pw      *playwright.Playwright
	browser playwright.Browser
)

type TestEnvironment struct {
	server   *web.Server
	shutdown *common.Shutdown
}

func (t *TestEnvironment) StartServer() {
	t.server.Start()
	time.Sleep(100 * time.Millisecond)
}

func (t *TestEnvironment) StopServer() {
	t.shutdown.Execute()
}

func NewTestEnvironment() *TestEnvironment {
	// General configuration
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	configuration := configuration.Load()
	configuration.Production = false

	// Graceful shutdown
	shutdown := common.NewShutdown()
	scheduler := common.NewScheduler(shutdown)

	// Repositories
	db := sqlite.Initialize(shutdown, scheduler, configuration)
	userRepository := sqlite.NewUserRepository(db)
	settingsRepository := memory.NewSettingsRepository()
	documentRepository := sqlite.NewDocumentRepository(db)
	folderRepository := sqlite.NewFolderRepository(db)

	// Event
	userEvents := synchronous.NewUserMessages()
	documentEvents := synchronous.NewDocumentMessages()

	// Storage
	documentStorage := filesystem.NewDocumentStorage(configuration)
	documentPreviewStorage := filesystem.NewDocumentPreviewStorage(configuration)

	// Features
	administration := administration.New(settingsRepository, userRepository, userEvents)
	archive := archive.New(documentRepository, documentStorage, documentPreviewStorage, documentEvents, folderRepository, userEvents, scheduler)

	// Web
	server := web.NewServer(administration, archive, shutdown, configuration)

	return &TestEnvironment{
		server:   server,
		shutdown: shutdown,
	}
}

func TestMain(m *testing.M) {
	os.Setenv("UNTERLAGEN_SERVER_SESSION_KEY", "my-key")
	defer os.Unsetenv("UNTERLAGEN_SERVER_SESSION_KEY")

	var err error
	err = playwright.Install()
	if err != nil {
		panic(err)
	}

	pw, err = playwright.Run()
	if err != nil {
		panic(err)
	}
	defer pw.Stop()
	browser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
		//SlowMo:   playwright.Float(500),
	})
	if err != nil {
		panic(err)
	}
	defer browser.Close()

	code := m.Run()
	os.Exit(code)
}
