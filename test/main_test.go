package test

import (
	"log/slog"
	"os"
	"testing"
	"time"
	"unterlagen/features/administration"
	"unterlagen/features/archive"
	"unterlagen/features/common"
	"unterlagen/features/search"
	"unterlagen/platform/configuration"
	"unterlagen/platform/database/memory"
	"unterlagen/platform/database/sqlite"
	"unterlagen/platform/llm"
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
	slog.SetDefault(slog.New(slog.DiscardHandler))

	// Common functionality
	shutdown := common.NewShutdown()
	jobScheduler := common.NewJobScheduler(shutdown)

	configuration := configuration.Load()
	configuration.Production = false

	// Database
	db := sqlite.Initialize(shutdown, jobScheduler, configuration)
	userRepository := sqlite.NewUserRepository(db)
	documentRepository := sqlite.NewDocumentRepository(db)
	folderRepository := sqlite.NewFolderRepository(db)
	taskRepository := sqlite.NewTaskRepository(db)
	settingsRepository := memory.NewSettingsRepository()
	searchRepository := sqlite.NewSearchRepository(db)

	// Messaging
	userMessages := synchronous.NewUserMessages()
	documentMessages := synchronous.NewDocumentMessages()

	// Storage
	documentStorage := filesystem.NewDocumentStorage(configuration)
	documentPreviewStorage := filesystem.NewDocumentPreviewStorage(configuration)

	// LLM
	documentSummarizer := llm.GetSummarizer(configuration)

	// Features
	taskScheduler := common.NewTaskScheduler(shutdown, taskRepository, common.TaskSchedulerModeSynchronous)
	administration := administration.New(settingsRepository, userRepository, userMessages, taskRepository)
	archive := archive.New(documentRepository, documentStorage, documentPreviewStorage, documentMessages, documentSummarizer, folderRepository, userMessages, jobScheduler, taskScheduler, shutdown)
	search := search.New(searchRepository, documentMessages, taskScheduler)

	// Web
	server := web.NewServer(administration, archive, search, shutdown, configuration)

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
	// Configure browser launch options based on GITHUB_ACTIONS environment variable
	launchOptions := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
		SlowMo:   playwright.Float(500),
	}

	// If running in GitHub Actions, use headless mode with no slowdown
	if os.Getenv("GITHUB_ACTIONS") != "" {
		launchOptions.Headless = playwright.Bool(true)
		launchOptions.SlowMo = playwright.Float(0)
	}

	browser, err = pw.Chromium.Launch(launchOptions)
	if err != nil {
		panic(err)
	}
	defer browser.Close()

	code := m.Run()
	os.Exit(code)
}
