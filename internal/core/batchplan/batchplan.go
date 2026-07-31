// Package batchplan は batch 翻訳のライフサイクル決定規則（IO を伴わない純粋関数）を持つ。
// 送信リクエストの custom_id 割り当て、進行段の遷移判定（二重送信を防ぐ冪等を含む）、
// 取得結果 1 件の行への書き戻し可否を決める。os・store・engine を import しない（core 規約）。
// provider は失敗種別の番兵（同期経路と共有）を参照するためだけに import する（core/prompt と同じ前例）。
// model は進行段の定数を参照するために import する。
//
// 本 package は「バッチ管理を単一モジュールへ寄せ、外から見て同期と batch を変えない」ための判断を集める。
// 決定規則だけを持ち、IO（DB・HTTP）は engine 側の薄いシェルへ寄せる。ユニットテスト 100% を基準にする。
package batchplan

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"aitranslationenginejp/internal/model"
	"aitranslationenginejp/internal/provider"
)

// EncodeCustomID は行種別と行 id を batch の custom_id「種別:id」へ組む。
// 叙述文・台詞は id がテーブルごと独立採番で衝突するため、種別接頭で名前空間化する（外部 batch の一意制約も満たす）。
func EncodeCustomID(kind string, id int64) string {
	return kind + ":" + strconv.FormatInt(id, 10)
}

// DecodeCustomID は custom_id「種別:id」を行種別と行 id へ分解する。形式不正・id 非数は error。
func DecodeCustomID(s string) (kind string, id int64, err error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", 0, fmt.Errorf("custom_id の形式が不正: %q", s)
	}
	id, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("custom_id の id 解析に失敗: %q: %w", s, err)
	}
	return parts[0], id, nil
}

// PlannedRequest は送信 1 件の計画。行種別・行 id と、engine が組んだ完成プロンプト。
// custom_id の割り当ては BuildBatchRequests が行う（種別と id から一意化する）。
type PlannedRequest struct {
	Kind   string
	RowID  int64
	Prompt provider.Prompt
}

// BuildBatchRequests は計画群を provider への送信要求（custom_id 付き）へ写す純粋関数。
func BuildBatchRequests(planned []PlannedRequest) []provider.BatchRequest {
	out := make([]provider.BatchRequest, len(planned))
	for i, p := range planned {
		out[i] = provider.BatchRequest{CustomID: EncodeCustomID(p.Kind, p.RowID), Prompt: p.Prompt}
	}
	return out
}

// ApplyKind は取得結果 1 件を行へどう適用するかの分類。
type ApplyKind int

// ApplyKind の分類値。確定と、種別ごとの未訳据え置き、および同期だけが run を止める Fatal。
const (
	ApplyConfirm                ApplyKind = iota // 訳文を確定として書き戻す
	ApplySkipTagLost                             // 実行時タグが AI 出力から欠落。未訳のまま据え置く
	ApplySkipStructuredParse                     // 構造化出力の空・スキーマ違反。未訳据え置き
	ApplySkipResponseUnreadable                  // 応答エンベロープの読み取り失敗。未訳据え置き
	ApplySkipServerTransient                     // サーバ一時失敗（429・5xx）。未訳据え置き
	ApplyFatal                                   // 通信断・4xx 認証など。同期は run 停止、batch は据え置き
)

// ApplyOutcome は適用分類と、確定時の訳文。
type ApplyOutcome struct {
	Kind ApplyKind
	Dest string
}

// DecideApply は取得結果 1 件（訳文・失敗・欠落タグ数）から行への書き戻し可否を決める純粋関数。
// 同期の本文フェーズと batch の反映が同じ関数を通し、「外から見て同期と batch が変わらない」を担保する。
//   - 成功（err==nil）でタグ欠落なし → 確定。
//   - 成功でタグ欠落あり → 未訳据え置き（ApplySkipTagLost）。
//   - skippable な失敗 3 種（構造化出力・応答エンベロープ・サーバ一時） → 種別ごとに未訳据え置き。
//   - それ以外の失敗（通信断・4xx など） → ApplyFatal。同期は run を止め、batch は据え置きにする（呼び出し側が決める）。
func DecideApply(translation string, err error, missingTags int) ApplyOutcome {
	switch {
	case err == nil:
		if missingTags > 0 {
			return ApplyOutcome{Kind: ApplySkipTagLost}
		}
		return ApplyOutcome{Kind: ApplyConfirm, Dest: translation}
	case errors.Is(err, provider.ErrStructuredParse):
		return ApplyOutcome{Kind: ApplySkipStructuredParse}
	case errors.Is(err, provider.ErrResponseUnreadable):
		return ApplyOutcome{Kind: ApplySkipResponseUnreadable}
	case errors.Is(err, provider.ErrServerTransient):
		return ApplyOutcome{Kind: ApplySkipServerTransient}
	default:
		return ApplyOutcome{Kind: ApplyFatal}
	}
}

// BatchProgress は進行 1 件の表示用スナップショット（副作用なしの状態確認で組む）。
// Stage は進行段（model.BatchStage*）。件数は現段 batch の状態確認由来。CanApply は取り込める完了段があるか。
type BatchProgress struct {
	Stage     string
	Total     int
	Pending   int
	Succeeded int
	Failed    int
	CanApply  bool
}

// BuildProgress は進行段と現段 batch の状態から表示用スナップショットを組む純粋関数。
// hasCurrentBatch=false（完了段、または現段の外部 ID が空の半端な進行）は件数 0・CanApply=false にする。
// それ以外は件数を status から写し、CanApply は現段が終端（処理待ち 0）かで決める。
// done 段は取り込む対象が無いので CanApply=false（早期に返す）。
func BuildProgress(stage string, hasCurrentBatch bool, status provider.BatchStatus) BatchProgress {
	if stage == model.BatchStageDone || !hasCurrentBatch {
		return BatchProgress{Stage: stage}
	}
	return BatchProgress{
		Stage:     stage,
		Total:     status.Total,
		Pending:   status.Pending,
		Succeeded: status.Succeeded,
		Failed:    status.Failed,
		CanApply:  status.Done,
	}
}

// RefreshStep は反映操作で、進行に対して次に取る行動。
type RefreshStep int

// RefreshStep の行動値。待機・何もしない、または固有名/本文の反映と段送り。
const (
	StepWait                      RefreshStep = iota // 現段の batch が未終端。何もしない
	StepApplyProperThenSubmitBody                    // 固有名段が終端。固有名を反映し、本文 batch を送る
	StepApplyProperThenAdvance                       // 固有名段が終端かつ本文 batch は送信済み。固有名を反映し段だけ進める
	StepApplyBodyThenComplete                        // 本文段が終端。本文を反映し完了にする
	StepNothing                                      // 完了済み、または現段の外部 batch が未送信。何もしない
)

// BlocksResubmit は進行中の batch への再送信を拒否すべきかを決める純粋関数。
// 拒否するのは、現在の進行段の外部 batch ID が非空（外部へ送信済みで、反映で回収・前進できる）場合だけ。
// 現段の外部 ID が空の進行（送信 HTTP 失敗などで途中で終わった半端な進行）は、再送信で作り直せるよう拒否しない。
// 完了段（done）その他も拒否しない。これで送信失敗の詰まりを削除なしに復旧できる。
func BlocksResubmit(stage, properBatchID, bodyBatchID string) bool {
	switch stage {
	case model.BatchStageProperNoun:
		return properBatchID != ""
	case model.BatchStageBody:
		return bodyBatchID != ""
	default:
		return false
	}
}

// DecideRefreshStep は進行の現段と、現段の外部 batch が終端かで、反映の次行動を決める純粋関数。
// 二重送信を防ぐ冪等: 固有名段の終端でも body 外部 ID が既にあれば本文を再送せず段だけ進める
// （固有名 dest 確定と本文 batch 送信の間でクラッシュし、本文送信済みのまま再反映した場合に効く）。
func DecideRefreshStep(stage, properBatchID, bodyBatchID string, terminal bool) RefreshStep {
	switch stage {
	case model.BatchStageProperNoun:
		switch {
		case properBatchID == "":
			return StepNothing // 固有名 batch が未送信（異常系）。触らない。
		case !terminal:
			return StepWait
		case bodyBatchID != "":
			return StepApplyProperThenAdvance // 本文は送信済み（クラッシュ復帰）。再送しない。
		default:
			return StepApplyProperThenSubmitBody
		}
	case model.BatchStageBody:
		switch {
		case bodyBatchID == "":
			return StepNothing // 本文 batch が未送信（異常系）。触らない。
		case !terminal:
			return StepWait
		default:
			return StepApplyBodyThenComplete
		}
	default:
		return StepNothing // 完了段（done）その他。
	}
}
