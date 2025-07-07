package migrator

import (
	"errors"
	"github.com/stlesnik/loyalty-system-project/internal/utils"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const migrationsDir = "./internal/repository/migrations"

func Run(dsn string) {
	m, err := migrate.New("file://"+migrationsDir, dsn)
	if err != nil {
		utils.Log.Errorw("migrate.New failed", "error", err)
		os.Exit(1)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		utils.Log.Errorw("migrations failed", "error", err)
		os.Exit(1)
	}
	utils.Log.Infow("migrations applied successfully")
}
