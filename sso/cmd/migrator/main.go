package migrator

import (
	"errors"
	"flag"
	"fmt"

	"github.com/HironixRotifer/grpc-authorization-service/internal/storage"
	"github.com/golang-migrate/migrate/v4"
)

func main() {
	var storagePath, migrationPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationPath, "migrations-path", "", "path to migration")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of the migration")
	flag.Parse()

	if storagePath == "" {
		panic("storage path is required")
	}
	if migrationPath == "" {
		panic("migration path is required")
	}

	m, err := migrate.New(
		"file://"+migrationPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationPath),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, storage.ErrNoChange) {
			fmt.Println("no migration applied")

			return
		}
		panic(err)
	}

	fmt.Println("migrations applied successfully")
}
