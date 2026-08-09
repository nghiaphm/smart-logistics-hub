package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	migratemysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"my-web-app.com/smart-logistic-hub/internal/infrastructure/config"
)

func main() {
	cmd := "up"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	cfg := config.LoadConfig()

	m, err := newMigrator(cfg)
	if err != nil {
		log.Fatalf("init migrator: %v", err)
	}

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("migrate up: %v", err)
		}
		log.Println("migrations applied")
	case "down":
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("migrate down: %v", err)
		}
		log.Println("migrations reverted")
	case "version":
		v, dirty, err := m.Version()
		if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
			log.Fatalf("migrate version: %v", err)
		}
		fmt.Printf("version=%d dirty=%v\n", v, dirty)
	default:
		log.Fatalf("unknown command %q (want: up, down, version)", cmd)
	}
}

func newMigrator(cfg *config.Config) (*migrate.Migrate, error) {
	db, err := sql.Open("mysql", databaseDSN(cfg))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	driver, err := migratemysql.WithInstance(db, &migratemysql.Config{
		DatabaseName: cfg.MariaDBName,
	})
	if err != nil {
		return nil, err
	}

	return migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)
}

func databaseDSN(cfg *config.Config) string {
	if cfg.MariaDBURI != "" {
		return cfg.MariaDBURI
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		cfg.MariaDBUser,
		cfg.MariaDBPassword,
		cfg.MariaDBHost,
		cfg.MariaDBPort,
		cfg.MariaDBName,
	)
}
