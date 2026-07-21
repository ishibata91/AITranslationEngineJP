package batchplan

import (
	"errors"
	"fmt"
	"testing"

	"aitranslationenginejp/internal/model"
	"aitranslationenginejp/internal/provider"
)

func TestEncodeCustomID(t *testing.T) {
	if got := EncodeCustomID(model.BatchKindNarration, 10); got != "n:10" {
		t.Fatalf("叙述文 custom_id = %q, want n:10", got)
	}
	if got := EncodeCustomID(model.BatchKindLine, 20); got != "l:20" {
		t.Fatalf("台詞 custom_id = %q, want l:20", got)
	}
	if got := EncodeCustomID(model.BatchKindProper, 30); got != "p:30" {
		t.Fatalf("固有名 custom_id = %q, want p:30", got)
	}
}

func TestDecodeCustomID(t *testing.T) {
	kind, id, err := DecodeCustomID("l:42")
	if err != nil {
		t.Fatalf("正常な custom_id で error: %v", err)
	}
	if kind != "l" || id != 42 {
		t.Fatalf("decode = (%q, %d), want (l, 42)", kind, id)
	}
}

func TestDecodeCustomID_エラー(t *testing.T) {
	cases := map[string]string{
		"コロン無し":  "n10",
		"種別が空":   ":10",
		"id が空":  "n:",
		"id が非数": "n:abc",
	}
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			if _, _, err := DecodeCustomID(in); err == nil {
				t.Fatalf("%q で error を期待したが nil", in)
			}
		})
	}
}

func TestBuildBatchRequests(t *testing.T) {
	if got := BuildBatchRequests(nil); len(got) != 0 {
		t.Fatalf("空計画で len=%d, want 0", len(got))
	}
	planned := []PlannedRequest{
		{Kind: model.BatchKindNarration, RowID: 1, Prompt: provider.Prompt{System: "s1", User: "u1"}},
		{Kind: model.BatchKindLine, RowID: 2, Prompt: provider.Prompt{System: "s2", User: "u2"}},
	}
	got := BuildBatchRequests(planned)
	if len(got) != 2 {
		t.Fatalf("len=%d, want 2", len(got))
	}
	if got[0].CustomID != "n:1" || got[0].Prompt.User != "u1" {
		t.Fatalf("req[0] = %+v", got[0])
	}
	if got[1].CustomID != "l:2" || got[1].Prompt.System != "s2" {
		t.Fatalf("req[1] = %+v", got[1])
	}
}

func TestDecideApply(t *testing.T) {
	cases := []struct {
		name        string
		translation string
		err         error
		missing     int
		wantKind    ApplyKind
		wantDest    string
	}{
		{"成功・タグ欠落なし", "訳文", nil, 0, ApplyConfirm, "訳文"},
		{"成功・タグ欠落あり", "訳文", nil, 1, ApplySkipTagLost, ""},
		{"構造化出力の失敗", "", fmt.Errorf("wrap: %w", provider.ErrStructuredParse), 0, ApplySkipStructuredParse, ""},
		{"応答エンベロープ失敗", "", fmt.Errorf("wrap: %w", provider.ErrResponseUnreadable), 0, ApplySkipResponseUnreadable, ""},
		{"サーバ一時失敗", "", fmt.Errorf("wrap: %w", provider.ErrServerTransient), 0, ApplySkipServerTransient, ""},
		{"その他の失敗は Fatal", "", errors.New("通信断"), 0, ApplyFatal, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := DecideApply(c.translation, c.err, c.missing)
			if got.Kind != c.wantKind {
				t.Fatalf("Kind = %d, want %d", got.Kind, c.wantKind)
			}
			if got.Dest != c.wantDest {
				t.Fatalf("Dest = %q, want %q", got.Dest, c.wantDest)
			}
		})
	}
}

func TestDecideRefreshStep(t *testing.T) {
	const (
		proper = model.BatchStageProperNoun
		body   = model.BatchStageBody
		done   = model.BatchStageDone
	)
	cases := []struct {
		name          string
		stage         string
		properBatchID string
		bodyBatchID   string
		terminal      bool
		want          RefreshStep
	}{
		{"固有名段・未送信", proper, "", "", true, StepNothing},
		{"固有名段・未終端", proper, "px", "", false, StepWait},
		{"固有名段・終端・本文未送信", proper, "px", "", true, StepApplyProperThenSubmitBody},
		{"固有名段・終端・本文送信済み（冪等）", proper, "px", "bx", true, StepApplyProperThenAdvance},
		{"本文段・未送信", body, "px", "", true, StepNothing},
		{"本文段・未終端", body, "px", "bx", false, StepWait},
		{"本文段・終端", body, "px", "bx", true, StepApplyBodyThenComplete},
		{"完了段", done, "px", "bx", true, StepNothing},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := DecideRefreshStep(c.stage, c.properBatchID, c.bodyBatchID, c.terminal)
			if got != c.want {
				t.Fatalf("step = %d, want %d", got, c.want)
			}
		})
	}
}

func TestBlocksResubmit(t *testing.T) {
	const (
		proper = model.BatchStageProperNoun
		body   = model.BatchStageBody
		done   = model.BatchStageDone
	)
	cases := []struct {
		name          string
		stage         string
		properBatchID string
		bodyBatchID   string
		want          bool
	}{
		{"固有名段・送信済み（拒否）", proper, "px", "", true},
		{"固有名段・外部ID空（半端・許可）", proper, "", "", false},
		{"本文段・送信済み（拒否）", body, "px", "bx", true},
		{"本文段・外部ID空（半端・許可）", body, "px", "", false},
		{"完了段（許可）", done, "px", "bx", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := BlocksResubmit(c.stage, c.properBatchID, c.bodyBatchID)
			if got != c.want {
				t.Fatalf("BlocksResubmit = %v, want %v", got, c.want)
			}
		})
	}
}

func TestBuildProgress(t *testing.T) {
	const (
		proper = model.BatchStageProperNoun
		body   = model.BatchStageBody
		done   = model.BatchStageDone
	)
	// 現段が処理待ちを残す状態と、終端（処理待ち 0）の状態。
	pending := provider.BatchStatus{Total: 113, Pending: 2, Succeeded: 111, Failed: 0, Done: false}
	terminal := provider.BatchStatus{Total: 10, Pending: 0, Succeeded: 10, Failed: 0, Done: true}

	cases := []struct {
		name       string
		stage      string
		hasCurrent bool
		status     provider.BatchStatus
		want       BatchProgress
	}{
		{
			"固有名段・処理中は件数を写し取り込み不可",
			proper, true, pending,
			BatchProgress{Stage: proper, Total: 113, Pending: 2, Succeeded: 111, Failed: 0, CanApply: false},
		},
		{
			"固有名段・終端は取り込み可",
			proper, true, terminal,
			BatchProgress{Stage: proper, Total: 10, Pending: 0, Succeeded: 10, Failed: 0, CanApply: true},
		},
		{
			"本文段・終端は取り込み可",
			body, true, terminal,
			BatchProgress{Stage: body, Total: 10, Pending: 0, Succeeded: 10, Failed: 0, CanApply: true},
		},
		{
			"完了段は状態確認せず件数 0・取り込み不可",
			done, false, provider.BatchStatus{},
			BatchProgress{Stage: done, CanApply: false},
		},
		{
			"現段の外部 ID 空（半端）は件数 0・取り込み不可",
			proper, false, provider.BatchStatus{},
			BatchProgress{Stage: proper, CanApply: false},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := BuildProgress(c.stage, c.hasCurrent, c.status)
			if got != c.want {
				t.Fatalf("BuildProgress = %+v, want %+v", got, c.want)
			}
		})
	}
}
