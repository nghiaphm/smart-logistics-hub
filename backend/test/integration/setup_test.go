package integration

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	db, err := setupTestDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup test database: %v\n", err)
		os.Exit(1)
	}
	testDB = db
	code := m.Run()
	testDB.Close()
	os.Exit(code)
}

func setupTestDB() (*sql.DB, error) {
	host := getEnv("MARIADB_HOST", "localhost")
	port := getEnv("MARIADB_PORT", "3307")
	user := getEnv("MARIADB_USER", "root")
	password := getEnv("MARIADB_PASSWORD", "root")
	dbName := getEnv("MARIADB_DB_NAME", "smart_logistics")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true&charset=utf8mb4",
		user, password, host, port, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if err := applyMigrations(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func applyMigrations(db *sql.DB) error {
	upPath := filepath.Join("..", "..", "migrations", "000001_initial_schema.up.sql")
	downPath := filepath.Join("..", "..", "migrations", "000001_initial_schema.down.sql")

	downContent, err := os.ReadFile(downPath)
	if err != nil {
		return fmt.Errorf("read down migration file: %w", err)
	}
	if _, err := db.Exec(string(downContent)); err != nil {
		// Ignore "table does not exist" during reset; other errors are real failures.
		if !strings.Contains(err.Error(), "Error 1051") {
			return fmt.Errorf("reset migrations: %w", err)
		}
	}

	content, err := os.ReadFile(upPath)
	if err != nil {
		return fmt.Errorf("read migration file: %w", err)
	}
	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func truncateTables(t *testing.T) {
	t.Helper()
	_, err := testDB.Exec("SET FOREIGN_KEY_CHECKS=0")
	if err != nil {
		t.Fatalf("disable FK checks: %v", err)
	}
	tables := []string{
		"tracking_events",
		"order_items",
		"orders",
		"inventory",
		"drivers",
		"ai_events",
		"inbound_items",
		"inbounds",
		"billing",
		"trip_stops",
		"trips",
		"products",
		"warehouses",
	}
	for _, table := range tables {
		if _, err := testDB.Exec("TRUNCATE TABLE " + table); err != nil {
			t.Fatalf("truncate %s: %v", table, err)
		}
	}
	_, err = testDB.Exec("SET FOREIGN_KEY_CHECKS=1")
	if err != nil {
		t.Fatalf("enable FK checks: %v", err)
	}
}
