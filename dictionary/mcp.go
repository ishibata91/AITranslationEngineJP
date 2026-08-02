package main

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type searchInput struct {
	Query    string `json:"query" jsonschema:"検索する英語の原語または日本語の訳語"`
	Category string `json:"category,omitempty" jsonschema:"categoryの完全一致。空なら全category"`
	Limit    int    `json:"limit,omitempty" jsonschema:"返す最大件数。既定50、最大200"`
}

type idInput struct {
	ID int64 `json:"id" jsonschema:"辞書項目のid"`
}

type addInput struct {
	Source   string `json:"source" jsonschema:"英語の原語"`
	Dest     string `json:"dest" jsonschema:"日本語の訳語"`
	Category string `json:"category,omitempty" jsonschema:"Skyrimのrecord種別または空文字"`
}

type updateInput struct {
	ID       int64  `json:"id" jsonschema:"辞書項目のid"`
	Revision int64  `json:"revision" jsonschema:"dictionary_getで取得したrevision"`
	Source   string `json:"source" jsonschema:"更新後の英語の原語"`
	Dest     string `json:"dest" jsonschema:"更新後の日本語の訳語"`
	Category string `json:"category,omitempty" jsonschema:"更新後のcategory"`
}

type emptyInput struct{}

func runMCP(ctx context.Context, s *store) error {
	server := newMCPServer(s)
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		return fmt.Errorf("MCP server: %w", err)
	}
	return nil
}

func newMCPServer(s *store) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "skyrim-translation-dictionary", Version: "0.1.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "dictionary_search", Description: "原語、訳語、categoryで辞書項目を検索する。完全一致、前方一致、部分一致の順で返す。"},
		func(ctx context.Context, _ *mcp.CallToolRequest, in searchInput) (*mcp.CallToolResult, searchResult, error) {
			out, err := s.search(ctx, in.Query, in.Category, in.Limit)
			return nil, out, err
		})
	mcp.AddTool(server, &mcp.Tool{Name: "dictionary_get", Description: "idで辞書項目と出どころを1件取得する。"},
		func(ctx context.Context, _ *mcp.CallToolRequest, in idInput) (*mcp.CallToolResult, entry, error) {
			out, err := s.get(ctx, in.ID)
			return nil, out, err
		})
	mcp.AddTool(server, &mcp.Tool{Name: "dictionary_add", Description: "辞書項目を1件追加する。同じ原語に異なる訳語またはcategoryを追加できる。"},
		func(ctx context.Context, _ *mcp.CallToolRequest, in addInput) (*mcp.CallToolResult, entry, error) {
			out, err := s.add(ctx, in.Source, in.Dest, in.Category)
			return nil, out, err
		})
	mcp.AddTool(server, &mcp.Tool{Name: "dictionary_update", Description: "idと取得時のrevisionを使い、辞書項目を更新する。取得後に別の更新があれば失敗する。"},
		func(ctx context.Context, _ *mcp.CallToolRequest, in updateInput) (*mcp.CallToolResult, entry, error) {
			out, err := s.update(ctx, in.ID, in.Revision, in.Source, in.Dest, in.Category)
			return nil, out, err
		})
	mcp.AddTool(server, &mcp.Tool{Name: "dictionary_status", Description: "辞書項目数、出どころ別件数、検索用indexの状態を返す。"},
		func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, status, error) {
			out, err := s.status(ctx)
			return nil, out, err
		})
	return server
}
