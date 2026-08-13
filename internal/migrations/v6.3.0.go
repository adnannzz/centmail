package migrations

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
)

// V6_3_0 performs the DB migrations for v6.3.0.
func V6_3_0(db *sqlx.DB, fs stuffbin.FileSystem, ko *koanf.Koanf, lo *log.Logger) error {
	// Saved, embeddable subscription forms, each tied to a set of lists.
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS forms (
			id           SERIAL PRIMARY KEY,
			uuid         UUID NOT NULL UNIQUE,
			name         TEXT NOT NULL,
			list_ids     INTEGER[] NOT NULL DEFAULT '{}',
			redirect_url TEXT NOT NULL DEFAULT '',
			created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		);
	`); err != nil {
		return err
	}

	return nil
}
