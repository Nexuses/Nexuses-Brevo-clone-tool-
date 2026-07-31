package migrations

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
)

func V6_4_0(db *sqlx.DB, fs stuffbin.FileSystem, ko *koanf.Koanf, lo *log.Logger) error {
	if _, err := db.Exec(`ALTER TYPE content_type ADD VALUE IF NOT EXISTS 'maily';`); err != nil {
		return err
	}
	if _, err := db.Exec(`ALTER TYPE template_type ADD VALUE IF NOT EXISTS 'campaign_maily';`); err != nil {
		return err
	}

	mailySrc, err := fs.Get("/static/email-templates/default-maily.json")
	if err != nil {
		return err
	}
	if _, err := db.Exec(`
		INSERT INTO templates (name, type, subject, body, body_source)
		SELECT $1, $2, $3, $4, $5
		WHERE NOT EXISTS (SELECT 1 FROM templates WHERE type = 'campaign_maily' LIMIT 1)`,
		"Sample Maily template", "campaign_maily", "", `<p>Hello {{ .Subscriber.Name }}</p>`, mailySrc.ReadBytes()); err != nil {
		return err
	}

	return nil
}
