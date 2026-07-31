package migrations

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
)

func V6_3_0(db *sqlx.DB, fs stuffbin.FileSystem, ko *koanf.Koanf, lo *log.Logger) error {
	if _, err := db.Exec(`
		ALTER TABLE campaigns
		ADD COLUMN IF NOT EXISTS tracking_config JSONB NOT NULL DEFAULT '{}';
	`); err != nil {
		return err
	}

	return nil
}
