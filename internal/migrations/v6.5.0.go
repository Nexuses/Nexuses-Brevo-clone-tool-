package migrations

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
)

func V6_5_0(db *sqlx.DB, fs stuffbin.FileSystem, ko *koanf.Koanf, lo *log.Logger) error {
	_, err := db.Exec(`
		INSERT INTO settings (key, value) VALUES(
			'cloudflare_email',
			'{"enabled":false,"account_id":"","api_token":"","max_conns":10,"timeout":"10s","max_msg_retries":2}'
		) ON CONFLICT DO NOTHING;
	`)
	return err
}
