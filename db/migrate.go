package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Apply は未適用の migration だけを conn へ順に適用し、schema version を user_version に刻む。
// 既適用分（user_version 未満）は再実行しない。これにより将来 IF NOT EXISTS でない DDL を足しても二重適用しない。
// migration の適用は本パッケージ（db）の責務とし、repository（internal/store）からは分離する。
func Apply(conn *sqlx.DB) error {
	migrations, err := Migrations()
	if err != nil {
		return err
	}
	var current int
	if err := conn.Get(&current, "PRAGMA user_version"); err != nil {
		return fmt.Errorf("schema version の読み取り: %w", err)
	}
	if current >= len(migrations) {
		return nil
	}
	for _, m := range migrations[current:] {
		if _, err := conn.Exec(m.SQL); err != nil {
			return fmt.Errorf("migration %s の適用: %w", m.Name, err)
		}
	}
	if _, err := conn.Exec(fmt.Sprintf("PRAGMA user_version = %d", len(migrations))); err != nil {
		return fmt.Errorf("schema version の記録: %w", err)
	}
	return nil
}
