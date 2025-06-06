package sqlite

import (
	"context"
	"embed"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unterlagen/features/common"
	"unterlagen/platform/configuration"

	"github.com/jmoiron/sqlx"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*
var migrations embed.FS

func Initialize(shutdown *common.Shutdown, scheduler *common.Scheduler, configuration configuration.Configuration) *sqlx.DB {
	var db *sqlx.DB
	var err error
	if configuration.Production {
		err = os.MkdirAll(configuration.Data.Directory, os.ModePerm)
		if err != nil {
			panic(err)
		}

		db, err = sqlx.Open("sqlite3", filepath.Join(configuration.Data.Directory, "unterlagen.db"))
		if err != nil {
			panic(err)
		}
	} else {
		db, err = sqlx.Open("sqlite3", ":memory:")
		if err != nil {
			panic(err)
		}
	}

	// Set connection pool settings
	db.SetMaxOpenConns(1)                // SQLite handles one writer at a time
	db.SetMaxIdleConns(1)                // Keep one connection alive
	db.SetConnMaxLifetime(time.Hour * 1) // Rotate connections hourly

	// Apply additional PRAGMA settings
	pragmas := []string{
		"PRAGMA journal_mode = WAL",
		"PRAGMA txlock = immediate",
		"PRAGMA busy_timeout = 30000",
		"PRAGMA foreign_keys = ON",
		"PRAGMA cache_size = -64000",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA temp_store = memory",
		"PRAGMA mmap_size = 268435456",
		"PRAGMA optimize",
		"PRAGMA analysis_limit = 1000",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			panic(err)
		}
	}

	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		panic(err)
	}
	if err := goose.Up(db.DB, "migrations"); err != nil {
		panic(err)
	}

	scheduler.Schedule(func(ctx context.Context) {
		ticker := time.NewTicker(24 * time.Hour) // Daily optimization
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if _, err := db.Exec("PRAGMA optimize"); err != nil {
					slog.Warn("failed to run periodic optimization", "error", err)
				}
			case <-ctx.Done():
				slog.Info("sqlite periodic optimization stopped")
				return
			}
		}
	})

	shutdown.AddCallback(func() {
		db.Exec("PRAGMA optimize")
		db.Close()
		slog.Info("closed database connection")
	})

	db.MapperFunc(func(s string) string {
		pattern := regexp.MustCompile("(\\p{Lu}+\\P{Lu}*)")
		s2 := pattern.ReplaceAllString(s, "${1}_")
		s2, _ = strings.CutSuffix(strings.ToLower(s2), "_")
		return s2
	})
	return db
}
