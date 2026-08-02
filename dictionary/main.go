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

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return errors.New("commandを指定する: import または mcp")
	}
	switch args[0] {
	case "import":
		fs := flag.NewFlagSet("import", flag.ContinueOnError)
		from := fs.String("from", "db/aitranslation.dev.sqlite3", "master_termを読む中心DB")
		dbPath := fs.String("db", defaultDictionaryPath, "生成する辞書DB")
		if err := fs.Parse(args[1:]); err != nil {
			return err
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
	case "mcp":
		fs := flag.NewFlagSet("mcp", flag.ContinueOnError)
		dbPath := fs.String("db", defaultDictionaryPath, "MCPが操作する辞書DB")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		s, err := openStore(*dbPath)
		if err != nil {
			return err
		}
		defer s.close() //nolint:errcheck
		return runMCP(ctx, s)
	default:
		return fmt.Errorf("未知のcommand %q: import または mcpを指定する", args[0])
	}
}

func writeJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
