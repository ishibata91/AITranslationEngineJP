package main

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type searchInput struct {
	Query              string `json:"query,omitempty" jsonschema:"原語、訳語、意味の検索文字列。空なら状態だけで絞る"`
	Category           string `json:"category,omitempty" jsonschema:"Skyrimのcategoryの完全一致"`
	GeneralMatchStatus string `json:"general_match_status,omitempty" jsonschema:"一般辞書との一致状態"`
	InclusionDecision  string `json:"inclusion_decision,omitempty" jsonschema:"undecided、include、excludeのいずれか"`
	ReviewStage        string `json:"review_stage,omitempty" jsonschema:"unreviewed、ai_reviewed、human_reviewedのいずれか"`
	Limit              int    `json:"limit,omitempty" jsonschema:"返す最大件数。既定50、最大200"`
}

type idInput struct {
	ID int64 `json:"id" jsonschema:"dictionary_senseのid"`
}

type senseAddInput struct {
	Source       string `json:"source" jsonschema:"英語表記"`
	Dest         string `json:"dest" jsonschema:"公式日本語訳"`
	PartOfSpeech string `json:"part_of_speech,omitempty" jsonschema:"noun、verb、adjective、adverb、other、unknownのいずれか"`
	Meaning      string `json:"meaning,omitempty" jsonschema:"同じ表記の別の意味と区別する説明"`
}

type senseUpdateInput struct {
	ID           int64  `json:"id" jsonschema:"dictionary_senseのid"`
	Revision     int64  `json:"revision" jsonschema:"dictionary_getで取得したrevision"`
	Dest         string `json:"dest" jsonschema:"更新後の公式日本語訳"`
	PartOfSpeech string `json:"part_of_speech" jsonschema:"更新後の品詞"`
	Meaning      string `json:"meaning" jsonschema:"更新後の意味を区別する説明"`
	ChangedBy    string `json:"changed_by" jsonschema:"変更したAIまたは人間の識別子"`
	Reason       string `json:"reason" jsonschema:"変更理由"`
}

type occurrenceAssignInput struct {
	OccurrenceID int64  `json:"occurrence_id" jsonschema:"dictionary_occurrenceのid"`
	SenseID      int64  `json:"sense_id" jsonschema:"割り当て先dictionary_senseのid"`
	ChangedBy    string `json:"changed_by" jsonschema:"変更したAIまたは人間の識別子"`
	Reason       string `json:"reason" jsonschema:"割り当て理由"`
}

type matchUpdateInput struct {
	MatchID   int64  `json:"match_id" jsonschema:"general_dictionary_matchのid"`
	Status    string `json:"status" jsonschema:"same_mean_and_translationまたはdifferent_meaning_or_translation"`
	ChangedBy string `json:"changed_by" jsonschema:"判定したAIまたは人間の識別子"`
	Reason    string `json:"reason" jsonschema:"意味と訳の判定理由"`
}

type matchQueueInput struct {
	Status            string `json:"status,omitempty" jsonschema:"既定はsame_mean_candidate。確定済み状態も指定できる"`
	InclusionDecision string `json:"inclusion_decision,omitempty" jsonschema:"undecided、include、excludeで収録判断を絞る"`
	ReviewStage       string `json:"review_stage,omitempty" jsonschema:"unreviewed、ai_reviewed、human_reviewedでレビュー段階を絞る"`
	AfterSenseID      int64  `json:"after_sense_id,omitempty" jsonschema:"前回のnext_sense_id。最初は0"`
	Limit             int    `json:"limit,omitempty" jsonschema:"返す最大件数。既定50、最大200"`
}

type reviewAddInput struct {
	SenseID           int64  `json:"sense_id" jsonschema:"dictionary_senseのid"`
	Revision          int64  `json:"revision" jsonschema:"dictionary_getで取得したrevision"`
	ReviewerKind      string `json:"reviewer_kind" jsonschema:"aiまたはhuman"`
	ReviewerReference string `json:"reviewer_reference" jsonschema:"判断したAIまたは人間の識別子"`
	Decision          string `json:"decision" jsonschema:"include、exclude、needs_humanのいずれか"`
	Reason            string `json:"reason" jsonschema:"収録判断の理由"`
}

type historyInput struct {
	TargetTable string `json:"target_table" jsonschema:"変更対象のtable名"`
	TargetID    int64  `json:"target_id" jsonschema:"変更対象のid"`
}

type emptyInput struct{}

func runMCP(ctx context.Context, s *store, wordNetPath string) error {
	server := newMCPServer(s, wordNetPath)
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		return fmt.Errorf("MCP server: %w", err)
	}
	return nil
}

func newMCPServer(s *store, wordNetPath string) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "skyrim-translation-dictionary", Version: "0.3.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "dictionary_search", Description: "原語、訳語、意味、Skyrimの分類、一般辞書との一致状態、収録判断、レビュー段階で意味候補を検索する。"},
		func(ctx context.Context, _ *mcp.CallToolRequest, in searchInput) (*mcp.CallToolResult, searchResult, error) {
			out, err := s.search(ctx, searchFilter(in))
			return nil, out, err
		})
	mcp.AddTool(server, &mcp.Tool{Name: "dictionary_get", Description: "idで意味、使用箇所、一般辞書との照合、レビューを取得する。"},
		func(ctx context.Context, _ *mcp.CallToolRequest, in idInput) (*mcp.CallToolResult, sense, error) {
			out, err := s.get(ctx, in.ID)
			return nil, out, err
		})
	mcp.AddTool(server, &mcp.Tool{Name: "dictionary_sense_add", Description: "英語表記へ意味と公式訳を追加する。"},
		func(ctx context.Context, _ *mcp.CallToolRequest, in senseAddInput) (*mcp.CallToolResult, sense, error) {
			out, err := s.add(ctx, in.Source, in.Dest, in.PartOfSpeech, in.Meaning)
			return nil, out, err
		})
	mcp.AddTool(server, &mcp.Tool{Name: "dictionary_sense_update", Description: "意味と訳語をrevision付きで更新し、変更理由を保存する。"},
		func(ctx context.Context, _ *mcp.CallToolRequest, in senseUpdateInput) (*mcp.CallToolResult, sense, error) {
			out, err := s.update(ctx, senseUpdate(in))
			return nil, out, err
		})
	mcp.AddTool(server, &mcp.Tool{Name: "dictionary_occurrence_assign", Description: "Skyrimの使用箇所を同じ英語表記の意味へ割り当て、理由を保存する。"},
		func(ctx context.Context, _ *mcp.CallToolRequest, in occurrenceAssignInput) (*mcp.CallToolResult, sense, error) {
			out, err := s.assignOccurrence(ctx, in.OccurrenceID, in.SenseID, in.ChangedBy, in.Reason)
			return nil, out, err
		})
	mcp.AddTool(server, &mcp.Tool{Name: "dictionary_general_match_update", Description: "一般辞書の意味候補がSkyrim側と同じ意味と訳かを確定し、理由を保存する。"},
		func(ctx context.Context, _ *mcp.CallToolRequest, in matchUpdateInput) (*mcp.CallToolResult, sense, error) {
			out, err := s.updateGeneralMatch(ctx, in.MatchID, in.Status, in.Reason, in.ChangedBy)
			return nil, out, err
		})
	mcp.AddTool(server, &mcp.Tool{Name: "dictionary_general_match_queue", Description: "一般辞書の意味候補を照合tableから順番に取得し、Skyrimの分類と辞書定義をまとめて確認する。"},
		func(ctx context.Context, _ *mcp.CallToolRequest, in matchQueueInput) (*mcp.CallToolResult, matchQueueResult, error) {
			out, err := s.generalMatchQueue(ctx, matchQueueFilter(in))
			return nil, out, err
		})
	mcp.AddTool(server, &mcp.Tool{Name: "dictionary_review_add", Description: "意味候補の収録判断と理由をレビューとして保存する。同じ意味と訳の確定がない一般語は除外できない。"},
		func(ctx context.Context, _ *mcp.CallToolRequest, in reviewAddInput) (*mcp.CallToolResult, sense, error) {
			out, err := s.addReview(ctx, reviewInput(in))
			return nil, out, err
		})
	mcp.AddTool(server, &mcp.Tool{Name: "dictionary_classify", Description: "日本語WordNetで全意味候補を照合し、候補と判定理由を保存する。収録判断は変更しない。"},
		func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, classifyResult, error) {
			out, err := s.classifyGeneralDictionary(ctx, wordNetPath)
			return nil, out, err
		})
	mcp.AddTool(server, &mcp.Tool{Name: "dictionary_history", Description: "指定したtableとidの変更履歴を取得する。"},
		func(ctx context.Context, _ *mcp.CallToolRequest, in historyInput) (*mcp.CallToolResult, changeHistory, error) {
			out, err := s.history(ctx, in.TargetTable, in.TargetID)
			return nil, out, err
		})
	mcp.AddTool(server, &mcp.Tool{Name: "dictionary_status", Description: "英語表記、意味、使用箇所、一般辞書判定、収録判断、レビュー段階の件数を返す。"},
		func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, status, error) {
			out, err := s.status(ctx)
			return nil, out, err
		})
	return server
}
