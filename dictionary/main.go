// Package main は事前作成辞書の生成 command と MCP server を提供する。
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

const defaultDictionaryPath = "dictionary/dictionary.sqlite3"
const defaultWordNetPath = "dictionary/reference/wnjpn.db"

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return errors.New("commandを指定する: import、classify、mcpのいずれか")
	}
	switch args[0] {
	case "import":
		return runImport(ctx, args[1:])
	case "classify":
		return runClassify(ctx, args[1:])
	case "mcp":
		return runMCPCommand(ctx, args[1:])
	default:
		return fmt.Errorf("未知のcommand %q: import、classify、mcpのいずれかを指定する", args[0])
	}
}

func runImport(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("import", flag.ContinueOnError)
	from := fs.String("from", "db/aitranslation.dev.sqlite3", "master_termを読む中心DB")
	dbPath := fs.String("db", defaultDictionaryPath, "生成する辞書DB")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("import commandの引数解析: %w", err)
	}
	s, err := openStore(*dbPath)
	if err != nil {
		return err
	}
	defer s.close() //nolint:errcheck
	result, err := importMasterTerms(ctx, *from, s)
	if err != nil {
		return err
	}
	return writeJSON(result)
}

func runClassify(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("classify", flag.ContinueOnError)
	dbPath := fs.String("db", defaultDictionaryPath, "分類する辞書DB")
	wordNetPath := fs.String("wordnet", defaultWordNetPath, "一般語判定に使う日本語WordNet")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("classify commandの引数解析: %w", err)
	}
	s, err := openStore(*dbPath)
	if err != nil {
		return err
	}
	defer s.close() //nolint:errcheck
	result, err := s.classifyGeneralDictionary(ctx, *wordNetPath)
	if err != nil {
		return err
	}
	return writeJSON(result)
}

func runMCPCommand(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("mcp", flag.ContinueOnError)
	dbPath := fs.String("db", defaultDictionaryPath, "MCPが操作する辞書DB")
	wordNetPath := fs.String("wordnet", defaultWordNetPath, "一般語判定に使う日本語WordNet")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("mcp commandの引数解析: %w", err)
	}
	s, err := openStore(*dbPath)
	if err != nil {
		return err
	}
	defer s.close() //nolint:errcheck
	return runMCP(ctx, s, *wordNetPath)
}

func writeJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("JSON出力: %w", err)
	}
	return nil
}
