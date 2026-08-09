package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"my-web-app.com/smart-logistic-hub/internal/infrastructure/config"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	dsn := cfg.MariaDBURI
	if dsn == "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			cfg.MariaDBUser,
			cfg.MariaDBPassword,
			cfg.MariaDBHost,
			cfg.MariaDBPort,
			cfg.MariaDBName,
		)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}

func Close(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}
