package migrations

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
)

// V6_7_1 adds base_domain to custom tracking domains for optional branded subdomains.
func V6_7_1(db *sqlx.DB, fs stuffbin.FileSystem, ko *koanf.Koanf, lo *log.Logger) error {
	_, err := db.Exec(`
		ALTER TABLE custom_tracking_domains
			ADD COLUMN IF NOT EXISTS base_domain TEXT NOT NULL DEFAULT '';

		UPDATE custom_tracking_domains
		SET base_domain = CASE
			WHEN array_length(string_to_array(domain, '.'), 1) > 2
			THEN regexp_replace(domain, '^[^.]+\.', '')
			ELSE domain
		END
		WHERE base_domain = '';
	`)
	return err
}
