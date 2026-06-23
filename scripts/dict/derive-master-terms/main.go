// Command derive-master-terms は人名の部分形（名のみ・短名）の確定訳語を master_term へ派生・追記する
// ビルド時コマンド。C# extractor が xTranslator 英日 XML から base 辞書（フルネーム :FULL）を master_term へ
// 書いた後に走らせる。XML 解析と派生の本体は engine.DeriveTermsFromXMLDir（翻訳 Run と共用）に閉じ、
// 本コマンドは DB 接続・base 原語の読み出し・派生行の追記だけを担う。派生行は category を "derive:<種別>" に
// して由来を残し、目視・差し戻しできる形にする。
//
// 実行順の前提: base 辞書（C# extractor の master_term 書き込み）が DB に入った後に走らせる。base が空のまま
// 走らせると、base 既出の単独名（例 Frost）との衝突回避が効かない。
//
//	go run ./scripts/dict/derive-master-terms --sqlite db/aitranslation.dev.sqlite3 \
//	    --terms-xml dictionaries/xTranslatorXMLs
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"aitranslationenginejp/db"
	"aitranslationenginejp/internal/engine"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite" // SQLite driver（pure-Go、CGO 不要）を "sqlite" 名で登録する。
)

func main() {
	dbPath := flag.String("sqlite", "", "中心 DB（SQLite）のパス。extractor が base 辞書を書いた後の DB を指す")
	xmlDir := flag.String("terms-xml", "dictionaries/xTranslatorXMLs", "xTranslator 英日 XML ディレクトリ")
	flag.Parse()
	if err := run(*dbPath, *xmlDir); err != nil {
		fmt.Fprintf(os.Stderr, "derive-master-terms: %v\n", err)
		os.Exit(1)
	}
}

// run は DB を開き、schema を適用し、base 原語集合を読み、XML から派生対を作って master_term へ追記する。
func run(dbPath, xmlDir string) error {
	if strings.TrimSpace(dbPath) == "" {
		return fmt.Errorf("--sqlite は必須")
	}
	conn, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("SQLite を開けない: %w", err)
	}
	defer conn.Close()
	if err := db.Apply(conn); err != nil {
		return fmt.Errorf("schema migration の適用: %w", err)
	}

	baseSources, err := loadBaseSources(conn)
	if err != nil {
		return err
	}
	rows, err := engine.DeriveTermsFromXMLDir(xmlDir, baseSources)
	if err != nil {
		return err
	}
	written, err := writeDerived(conn, rows)
	if err != nil {
		return err
	}
	fmt.Printf("[derive] master_term へ派生 %d 件を追記（候補 %d 件、base 既出を除く）\n", written, len(rows))
	return nil
}

// loadBaseSources は master_term の既存原語集合を読む。派生が base 既出の原語と衝突しないよう除外判定に使う。
func loadBaseSources(conn *sqlx.DB) (map[string]bool, error) {
	var sources []string
	if err := conn.Select(&sources, `SELECT DISTINCT source FROM master_term`); err != nil {
		return nil, fmt.Errorf("base 原語の取得: %w", err)
	}
	set := make(map[string]bool, len(sources))
	for _, s := range sources {
		set[s] = true
	}
	return set, nil
}

// writeDerived は派生対を master_term へ INSERT OR IGNORE で追記する。実際に追加した行数を返す。
// category は "derive:<種別>" にして由来を残す。UNIQUE(category, source) と OR IGNORE で二重追記を防ぐ。
func writeDerived(conn *sqlx.DB, rows []engine.DerivedTerm) (int, error) {
	tx, err := conn.Beginx()
	if err != nil {
		return 0, fmt.Errorf("トランザクション開始: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // commit 済みなら no-op。失敗時の後始末用。
	stmt, err := tx.Preparex(`INSERT OR IGNORE INTO master_term (source, dest, category) VALUES (?, ?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("INSERT 文の準備: %w", err)
	}
	defer stmt.Close()

	written := 0
	for _, r := range rows {
		res, err := stmt.Exec(r.Source, r.Dest, "derive:"+r.Kind)
		if err != nil {
			return 0, fmt.Errorf("派生 %q の書き込み: %w", r.Source, err)
		}
		n, _ := res.RowsAffected()
		written += int(n)
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("コミット: %w", err)
	}
	return written, nil
}
