// Package main は事前作成辞書を管理する MCP server を提供する。
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
)

const defaultDictionaryPath = "db/dictionary.sqlite3"
const defaultWordNetPath = "dictionary/reference/wnjpn.db"

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) == 0 || args[0] != "mcp" {
		return fmt.Errorf("commandはmcpだけを指定できる")
	}
	return runMCPCommand(ctx, args[1:])
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
