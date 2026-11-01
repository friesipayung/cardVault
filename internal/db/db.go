package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// DB wraps *sql.DB and the on-disk file path.
type DB struct {
	*sql.DB
	Path string
}

// Open opens (or creates) the SQLite database at path, enabling WAL + foreign keys.
func Open(path string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	sqlDB, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	// SQLite pragmas for a single-user desktop/mobile app.
	for _, p := range []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA foreign_keys=ON;",
		"PRAGMA busy_timeout=5000;",
	} {
		if _, err := sqlDB.Exec(p); err != nil {
			sqlDB.Close()
			return nil, fmt.Errorf("pragma %s: %w", p, err)
		}
	}
	return &DB{DB: sqlDB, Path: path}, nil
}

// Migrate runs the schema migrations. Idempotent.
func (d *DB) Migrate() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS providers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			logo_url TEXT,
			barcode_format TEXT,
			country TEXT,
			notes TEXT,
			created_at DATETIME NOT NULL DEFAULT (datetime('now'))
		);`,
		`CREATE TABLE IF NOT EXISTS cards (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			provider_id INTEGER,
			provider_name TEXT NOT NULL,
			card_number TEXT NOT NULL,
			card_holder_name TEXT,
			barcode_type TEXT,
			barcode_value TEXT,
			nfc_uid TEXT,
			nfc_data BLOB,
			notes TEXT,
			created_at DATETIME NOT NULL DEFAULT (datetime('now')),
			updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
			last_used_at DATETIME,
			FOREIGN KEY (provider_id) REFERENCES providers(id) ON DELETE SET NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_cards_provider ON cards(provider_id);`,
		`CREATE INDEX IF NOT EXISTS idx_cards_number ON cards(card_number);`,
		`CREATE TABLE IF NOT EXISTS nfc_captures (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			card_id INTEGER,
			uid TEXT NOT NULL,
			atqa TEXT,
			sak TEXT,
			ndef TEXT,
			raw_dump BLOB,
			captured_at DATETIME NOT NULL DEFAULT (datetime('now')),
			FOREIGN KEY (card_id) REFERENCES cards(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		);`,
	}
	for _, s := range stmts {
		if _, err := d.Exec(s); err != nil {
			return err
		}
	}
	return d.seedProviders()
}

// seedProviders inserts a small preset list of loyalty providers if empty.
func (d *DB) seedProviders() error {
	var n int
	if err := d.QueryRow("SELECT COUNT(*) FROM providers").Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	presets := []struct {
		name, format, country string
	}{
		{"Starbucks", "QR", "Global"},
		{"IKEA Family", "Barcode", "Global"},
		{"Tesco Clubcard", "Barcode", "UK"},
		{"Carrefour", "Barcode", "EU"},
		{"Alfamart", "Barcode", "ID"},
		{"Indomaret", "Barcode", "ID"},
		{"GrabRewards", "QR", "SEA"},
		{"GoPay", "QR", "ID"},
	}
	for _, p := range presets {
		_, _ = d.Exec(
			"INSERT OR IGNORE INTO providers(name, barcode_format, country) VALUES(?,?,?)",
			p.name, p.format, p.country)
	}
	return nil
}

// DefaultPath returns the default DB file location.
func DefaultPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cardvault", "cardvault.db")
}
