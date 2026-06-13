// Package db は C#↔Go 契約の SQL schema 正本（migrations/*.sql）を embed して公開する。
// Go app は起動時にこれを適用し、同じファイルを C# extractor も冪等 ensure する。
package db

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Migration は 1 つの migration ファイル。Name はファイル名順で適用順を決める。
type Migration struct {
	Name string
	SQL  string
}

// Migrations は migrations/ の SQL をファイル名昇順で返す。
// fs.ReadDir はファイル名昇順でソート済みの entry を返す（io/fs 仕様）ため、その順序が適用順になる。
func Migrations() ([]Migration, error) {
	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("migration ディレクトリ読み込み: %w", err)
	}
	migrations := make([]Migration, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		content, err := migrationsFS.ReadFile("migrations/" + e.Name())
		if err != nil {
			return nil, fmt.Errorf("migration %s 読み込み: %w", e.Name(), err)
		}
		migrations = append(migrations, Migration{Name: e.Name(), SQL: string(content)})
	}
	return migrations, nil
}
