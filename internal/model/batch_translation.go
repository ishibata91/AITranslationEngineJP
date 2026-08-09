package model

// batch 進行段の定数。batch_translation.stage 列に対応する。
// 進行は固有名 batch → 本文 batch → 完了 の順にたどる。
const (
	BatchStageProperNoun = "proper_noun" // 固有名 batch を送信済みで反映待ち
	BatchStageBody       = "body"        // 固有名を反映し、本文 batch を送信済みで反映待ち
	BatchStageDone       = "done"        // 本文 batch まで反映し、進行を完了
)

// batch リクエストの行種別（custom_id の接頭・batch_request.kind 列）。
// 叙述文・台詞は id がテーブルごと独立採番で衝突するため、種別接頭で名前空間化する。
const (
	BatchKindNarration = "n" // 叙述文（narration）
	BatchKindLine      = "l" // 台詞（line）
	BatchKindProper    = "p" // 固有名（proper_noun）
)

// BatchTranslation は対象 plugin 1 つの batch 進行の本体（batch_translation テーブル）。plugin と 1 対 1。
// Stage は進行段（proper_noun / body / done）。ProperBatchID / BodyBatchID は各段で現在処理している外部 batch ID（未送信は空）。
// Provider は外部 batch ID を発行した提供元。Model は送信に使う翻訳モデル名。CreatedAt は進行開始時刻の文字列。
type BatchTranslation struct {
	ID            int64  `db:"id"`
	Plugin        string `db:"plugin"`
	Provider      string `db:"provider"`
	Model         string `db:"model"`
	Stage         string `db:"stage"`
	ProperBatchID string `db:"proper_batch_id"`
	BodyBatchID   string `db:"body_batch_id"`
	CreatedAt     string `db:"created_at"`
}

// BatchRequest は送信したリクエストと結果の対応 1 件（batch_request テーブル）。batch_translation の子。
// ExternalBatchID はどの外部 batch に載せたか。CustomID は「種別:id」。Kind / RowID は結果適用時の行特定用。
type BatchRequest struct {
	ID              int64  `db:"id"`
	BatchID         int64  `db:"batch_id"`
	ExternalBatchID string `db:"external_batch_id"`
	CustomID        string `db:"custom_id"`
	Kind            string `db:"kind"`
	RowID           int64  `db:"row_id"`
	ReferencesJSON  string `db:"references_json"`
	PromptHash      string `db:"prompt_hash"`
	SendState       string `db:"send_state"`
}
