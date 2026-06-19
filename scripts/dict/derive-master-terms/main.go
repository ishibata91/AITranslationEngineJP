// Command derive-master-terms は人名の部分形（名のみ・短名）の確定訳語を master_term へ派生・追記する
// ビルド時コマンド。C# extractor が xTranslator 英日 XML から base 辞書（フルネーム :FULL）を master_term へ
// 書いた後に走らせる。本コマンドは同じ XML 群から NPC の人名対（NPC_:FULL / NPC_:SHRT）と会話文（INFO:NAM1）を
// 読み、本文の用法分布を作り、純粋ルール engine.DeriveTerms で安全な部分形だけを派生し、base 衝突を除いて
// master_term へ追記する。派生行は category を "derive:<種別>" にして由来を残し、目視・差し戻しできる形にする。
//
// 実行順の前提: base 辞書（C# extractor の master_term 書き込み）が DB に入った後に走らせる。base が空のまま
// 走らせると、base 既出の単独名（例 Frost）との衝突回避が効かない。
//
//	go run ./scripts/dict/derive-master-terms --sqlite db/aitranslation.dev.sqlite3 \
//	    --terms-xml dictionaries/xTranslatorXMLs
package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
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
	rows, err := deriveFromXMLDir(xmlDir, baseSources)
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

// baseGamePrefixes は base ゲーム（Bethesda 公式本体・DLC）の XML ファイル名接頭。姓名分割（two）を
// base ゲーム限定にするための判定に使う。第三者 mod は Source/Dest 対応が信頼できないため two を派生しない。
var baseGamePrefixes = []string{"Skyrim", "Dawnguard", "Dragonborn", "HearthFires", "Update"}

// isBaseGame は XML ファイル名が base ゲーム由来かを返す。
func isBaseGame(filename string) bool {
	for _, p := range baseGamePrefixes {
		if strings.HasPrefix(filename, p) {
			return true
		}
	}
	return false
}

// deriveFromXMLDir はディレクトリ内の全 XML を解析し、用法分布を作り、純粋ルールで派生対を返す。
func deriveFromXMLDir(xmlDir string, baseSources map[string]bool) ([]engine.DerivedTerm, error) {
	entries, err := filepath.Glob(filepath.Join(xmlDir, "*.xml"))
	if err != nil {
		return nil, fmt.Errorf("XML の列挙: %w", err)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("XML が無い: %s", xmlDir)
	}
	sort.Strings(entries)

	var fulls, shrts []engine.NamePair
	var dialogues []string
	for _, path := range entries {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("%s の読み込み: %w", path, err)
		}
		f, s, d, err := parseTermXML(data, isBaseGame(filepath.Base(path)))
		if err != nil {
			return nil, fmt.Errorf("%s の解析: %w", path, err)
		}
		fulls = append(fulls, f...)
		shrts = append(shrts, s...)
		dialogues = append(dialogues, d...)
	}
	usage := engine.BuildUsage(dialogues)
	return engine.DeriveTerms(fulls, shrts, usage, baseSources, engine.DefaultDeriveConfig()), nil
}

// xmlString は xTranslator XML の 1 レコード（<String>）。REC（種別）と原語・確定訳語を読む。
type xmlString struct {
	REC    string `xml:"REC"`
	Source string `xml:"Source"`
	Dest   string `xml:"Dest"`
}

// parseTermXML は 1 つの XML から NPC のフルネーム対・短名対と会話文の英語原文を取り出す。
// baseGame は供給元が base ゲームかを示し、フルネーム対へ印を付ける（two の base ゲーム限定判定用）。
// encoding/xml はエンティティ（&lt; など）を復号する。先頭 UTF-8 BOM は除いてから解析する。
func parseTermXML(data []byte, baseGame bool) (fulls, shrts []engine.NamePair, dialogues []string, err error) {
	data = bytes.TrimPrefix(data, []byte("\uFEFF"))
	dec := xml.NewDecoder(bytes.NewReader(data))
	for {
		tok, tokErr := dec.Token()
		if tokErr == io.EOF {
			break
		}
		if tokErr != nil {
			return nil, nil, nil, tokErr
		}
		start, ok := tok.(xml.StartElement)
		if !ok || start.Name.Local != "String" {
			continue
		}
		var s xmlString
		if decErr := dec.DecodeElement(&s, &start); decErr != nil {
			return nil, nil, nil, decErr
		}
		rec := strings.ToUpper(strings.TrimSpace(s.REC))
		src := strings.TrimSpace(s.Source)
		dst := strings.TrimSpace(s.Dest)
		switch rec {
		case "NPC_:FULL":
			if src != "" && dst != "" {
				fulls = append(fulls, engine.NamePair{Source: src, Dest: dst, BaseGame: baseGame})
			}
		case "NPC_:SHRT":
			if src != "" && dst != "" {
				shrts = append(shrts, engine.NamePair{Source: src, Dest: dst, BaseGame: baseGame})
			}
		case "INFO:NAM1":
			if src != "" {
				dialogues = append(dialogues, src)
			}
		}
	}
	return fulls, shrts, dialogues, nil
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
