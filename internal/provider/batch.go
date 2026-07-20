package provider

import "context"

// BatchRequest は非同期 batch 翻訳の 1 件の要求。
// CustomID は結果対応の一意キー（engine 側の「種別:id」複合キー）。Prompt は engine が組んだ完成プロンプト。
type BatchRequest struct {
	CustomID string
	Prompt   Prompt
}

// BatchStatus は外部 batch の状態。件数と終端判定だけを port の外へ渡す（vendor 固有の状態文字列は隠す）。
// Done は未処理が無くなり結果を取り出せる終端に達したか。件数は状態表示・観測に使う。
type BatchStatus struct {
	Total     int
	Pending   int
	Succeeded int
	Failed    int
	Cancelled int
	Done      bool
}

// BatchResult は外部 batch の 1 件の結果。CustomID で送信要求と対応づく。
// 成功時は Translation に訳文が入り Err は nil。失敗時は Translation は空で、Err に失敗種別（番兵）が入る。
// 失敗種別は同期経路と同じ番兵（ErrStructuredParse / ErrResponseUnreadable / ErrServerTransient）で表し、
// 呼び出し側（engine の反映）は同期と同じ思想で該当行を未訳のまま据え置く（再送信で回収）。
type BatchResult struct {
	CustomID    string
	Translation string
	Err         error
}

// BatchTranslator は非同期 batch 翻訳の port。同期 Translator と並ぶ 2 つ目の port。
// 送信で外部 batch を作り、状態確認で終端を判定し、結果取得で custom_id ごとの訳文または失敗を得る。
// 接続情報（Connection）は永続化せず、各操作へ都度渡す。
type BatchTranslator interface {
	// SubmitBatch は要求群を 1 つの外部 batch として送り、外部 batch ID を返す。
	SubmitBatch(ctx context.Context, conn Connection, model string, requests []BatchRequest) (externalBatchID string, err error)
	// PollBatch は外部 batch の状態を返す（未処理/成功/失敗の件数と終端判定）。
	PollBatch(ctx context.Context, conn Connection, externalBatchID string) (BatchStatus, error)
	// FetchResults は外部 batch の全結果を返す（ページングを最後までたどる）。custom_id ごとの訳文または失敗種別。
	FetchResults(ctx context.Context, conn Connection, externalBatchID string) ([]BatchResult, error)
}
