package migrations

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
)

// V6_7_0 adds custom tracking domains and app.tracking_url.
func V6_7_0(db *sqlx.DB, fs stuffbin.FileSystem, ko *koanf.Koanf, lo *log.Logger) error {
	_, err := db.Exec(`
		INSERT INTO settings (key, value, updated_at)
			VALUES ('app.tracking_url', '""', NOW())
			ON CONFLICT (key) DO NOTHING;

		CREATE TABLE IF NOT EXISTS custom_tracking_domains (
			id               SERIAL PRIMARY KEY,
			user_id          INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
			domain           TEXT NOT NULL UNIQUE,
			status           TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'verified', 'failed')),
			dns_record_type  TEXT NOT NULL DEFAULT 'CNAME',
			dns_record_name  TEXT NOT NULL DEFAULT '',
			dns_record_value TEXT NOT NULL DEFAULT '',
			verified_at      TIMESTAMP WITH TIME ZONE NULL,
			last_error       TEXT NOT NULL DEFAULT '',
			created_at       TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at       TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_ctd_user_id ON custom_tracking_domains(user_id);
		CREATE INDEX IF NOT EXISTS idx_ctd_status ON custom_tracking_domains(status);

		ALTER TABLE campaigns ADD COLUMN IF NOT EXISTS tracking_domain_id INTEGER NULL
			REFERENCES custom_tracking_domains(id) ON DELETE SET NULL ON UPDATE CASCADE;
	`)
	return err
}
