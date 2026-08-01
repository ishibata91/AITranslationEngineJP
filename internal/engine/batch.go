package engine

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"aitranslationenginejp/internal/core/batchplan"
	"aitranslationenginejp/internal/core/dictionary"
	"aitranslationenginejp/internal/core/prompt"
	"aitranslationenginejp/internal/core/runtimetag"
	"aitranslationenginejp/internal/model"
	"aitranslationenginejp/internal/provider"
)

// errNoBatchProvider は batch provider が未配線のときに送信・反映が返すエラー。
var errNoBatchProvider = errors.New("batch provider が未配線")

// errBatchProviderMismatch は保存済み進行と画面の提供元が違う場合に、外部 API と DB を触る前に返す。
var errBatchProviderMismatch = errors.New("batch provider が保存済み進行と一致しない")

// BatchStore は batch 進行の永続に使う中心データアクセス（使う分だけ宣言する）。
// 行の読み書き（未訳一覧・dest 更新）や辞書・テンプレートは Engine が持つ store を通すため、ここには batch 固有だけを置く。
type BatchStore interface {
	StartBatchProgression(ctx context.Context, plugin, provider, model string) (int64, error)
	RecordBatchExternalID(ctx context.Context, batchID int64, stage, externalID string) error
	AdvanceBatchStage(ctx context.Context, batchID int64, stage string) error
	GetBatchProgression(ctx context.Context, plugin string) (model.BatchTranslation, bool, error)
	InsertBatchRequests(ctx context.Context, rows []model.BatchRequest) (int, error)
	ListBatchRequests(ctx context.Context, externalBatchID string) ([]model.BatchRequest, error)
	MarkSyncRetryReady(ctx context.Context, plugin string) error
}

// BatchRunner は OpenAI と xAI の batch 翻訳のオーケストレーション（薄いシェル）。
// 純粋核 batchplan の判断を、Engine の読み書き部品・provider batch port・batch 永続へ束ねるだけにする。
// 送信（SubmitBatch）で固有名 batch を作り、状態確認（ProgressStatus・副作用なし）で進行を観測し、
// 取り込み（RefreshPlugin）で結果反映・本文 batch 送信を進める。
// 固有名 → 本文 の 2 段を逐次でたどり、外から見て同期経路と結果が変わらないようにする。
type BatchRunner struct {
	e       *Engine
	batches map[string]provider.BatchTranslator
	store   BatchStore
}

// NewBatchRunner は BatchRunner を生成する。bootstrap が唯一の生成点。
// provider に対応する実装が無い場合、送信・反映は明示エラーを返す。
func NewBatchRunner(e *Engine, batches map[string]provider.BatchTranslator, store BatchStore) *BatchRunner {
	return &BatchRunner{e: e, batches: batches, store: store}
}

func (r *BatchRunner) batchFor(providerName string, conn provider.Connection) (provider.BatchTranslator, error) {
	if providerName == provider.BatchProviderOpenAI && strings.TrimSpace(conn.APIKey) == "" {
		return nil, fmt.Errorf("OpenAI API キーが空")
	}
	batch := r.batches[providerName]
	if batch == nil {
		return nil, fmt.Errorf("%w: %s", errNoBatchProvider, providerName)
	}
	return batch, nil
}

func (r *BatchRunner) batchForProgress(ctx context.Context, providerName string, conn provider.Connection, prog model.BatchTranslation) (provider.BatchTranslator, error) {
	if providerName != prog.Provider {
		slog.WarnContext(ctx, "batch provider の不一致を拒否した",
			slog.String("event", "batch_provider_mismatch"),
			slog.String("where", "engine.batch"),
			slog.String("result", "rejected"),
			slog.String("id", prog.Plugin),
			slog.String("reason", "provider_mismatch"),
		)
		return nil, fmt.Errorf("%w: 保存済み=%s 選択中=%s", errBatchProviderMismatch, prog.Provider, providerName)
	}
	return r.batchFor(providerName, conn)
}

// SubmitBatch は対象 plugin の未訳固有名を固有名 batch として送信し、進行を開始する。
// 既訳ありの固有名（権威訳）は AI を呼ばず即確定し、既訳なしだけを batch へ載せる（同期の固有名フェーズと同じ選別）。
// 固有名 batch の対象が無い場合は固有名段を飛ばして本文 batch を送る。
// 進行中の batch がある plugin への再送信は拒否する（反映で完了させてから送り直す）。
func (r *BatchRunner) SubmitBatch(ctx context.Context, providerName string, conn provider.Connection, modelName, plugin string) error {
	batch, err := r.batchFor(providerName, conn)
	if err != nil {
		return err
	}
	if err := r.ensureNoActiveProgression(ctx, plugin); err != nil {
		return err
	}
	batchID, err := r.store.StartBatchProgression(ctx, plugin, providerName, modelName)
	if err != nil {
		return fmt.Errorf("batch 進行の開始: %w", err)
	}
	planned, err := r.planProperRequests(ctx, plugin)
	if err != nil {
		return err
	}
	if len(planned) == 0 {
		// 固有名の AI 訳が無い（全て既訳流用 or 固有名なし）。固有名 batch を挟まず本文 batch を送る。
		return r.submitBodyBatch(ctx, batch, conn, modelName, plugin, batchID)
	}
	return r.sendStage(ctx, batch, conn, modelName, batchID, model.BatchStageProperNoun, planned)
}

// ensureNoActiveProgression は外部へ送信済みで待機中の batch がある plugin への再送信を拒否する。
// 拒否は「現段の外部 batch ID が非空」の進行に限る（反映で回収・前進できるため）。
// 現段の外部 ID が空の半端な進行（送信 HTTP 失敗などで途中で終わった進行）は、再送信で reset して作り直せるよう拒否しない。
func (r *BatchRunner) ensureNoActiveProgression(ctx context.Context, plugin string) error {
	prog, ok, err := r.store.GetBatchProgression(ctx, plugin)
	if err != nil {
		return fmt.Errorf("batch 進行の確認: %w", err)
	}
	if ok && batchplan.BlocksResubmit(prog.Stage, prog.ProperBatchID, prog.BodyBatchID) {
		return fmt.Errorf("対象 plugin に送信済みで待機中の batch がある。反映で完了させてから再送信する: %s", plugin)
	}
	return nil
}

// planProperRequests は未訳固有名を供給源で選別する。既訳ありは即確定し、既訳なしを送信計画へ積んで返す。
func (r *BatchRunner) planProperRequests(ctx context.Context, plugin string) ([]batchplan.PlannedRequest, error) {
	propers, err := r.e.store.ListUntranslatedProperNouns(ctx, plugin)
	if err != nil {
		return nil, fmt.Errorf("未訳固有名の取得: %w", err)
	}
	authoritative, err := r.e.authoritativeTerms(ctx)
	if err != nil {
		return nil, err
	}
	tmpl, err := r.e.store.GetPromptTemplate(ctx)
	if err != nil {
		return nil, fmt.Errorf("プロンプトテンプレートの取得: %w", err)
	}
	instructionByKey, _, err := r.e.directiveLookups(ctx)
	if err != nil {
		return nil, err
	}
	properInstruction := instructionByKey[directiveProperNoun]

	var planned []batchplan.PlannedRequest
	for _, pn := range propers {
		sup := SelectSupply(pn.Source, authoritative)
		if sup.Kind == SupplyAuthoritative {
			// 既訳流用。AI を呼ばず権威訳を proper_noun へ確定として書く（同期と同じ）。
			if err := r.e.store.UpdateProperNounDest(ctx, pn.ID, sup.Dest, statusTranslated); err != nil {
				return nil, fmt.Errorf("固有名（既訳流用）の書き戻し: %w", err)
			}
			continue
		}
		planned = append(planned, batchplan.PlannedRequest{
			Kind:   model.BatchKindProper,
			RowID:  pn.ID,
			Prompt: prompt.ComposePrompt(tmpl.BaseDirective, properInstruction, pn.Source),
		})
	}
	return planned, nil
}

// submitBodyBatch は確定した固有名で本文の機械置換辞書を組み、未訳の叙述文・台詞を本文 batch として送る。
// 既存訳と完全一致する本文は AI を呼ばず即確定し、batch へ載せない（同期の本文フェーズと同じ）。
// 本文の未訳が無い場合は本文 batch を挟まず進行を完了にする。送信後は進行段を本文へ進める。
func (r *BatchRunner) submitBodyBatch(ctx context.Context, batch provider.BatchTranslator, conn provider.Connection, modelName, plugin string, batchID int64) error {
	planned, err := r.planBodyRequests(ctx, plugin)
	if err != nil {
		return err
	}
	if err := r.store.MarkSyncRetryReady(ctx, plugin); err != nil {
		return err
	}
	if len(planned) == 0 {
		// 本文の未訳が無い（全て既訳流用 or 空）。本文 batch を挟まず完了にする。
		if err := r.store.AdvanceBatchStage(ctx, batchID, model.BatchStageDone); err != nil {
			return fmt.Errorf("進行段の更新（完了）: %w", err)
		}
		return nil
	}
	if err := r.sendStage(ctx, batch, conn, modelName, batchID, model.BatchStageBody, planned); err != nil {
		return err
	}
	if err := r.store.AdvanceBatchStage(ctx, batchID, model.BatchStageBody); err != nil {
		return fmt.Errorf("進行段の更新（本文）: %w", err)
	}
	return nil
}

// planBodyRequests は確定固有名で辞書を組み、未訳の叙述文・台詞を送信計画へ積んで返す。既存訳一致は即確定し載せない。
func (r *BatchRunner) planBodyRequests(ctx context.Context, plugin string) ([]batchplan.PlannedRequest, error) {
	// 確定した NPC 名から人名の部分形を派生し、対象 plugin の proper_noun へ足す（同期 Run と同じ位置＝辞書を組む前）。
	// 固有名 batch を反映した後に本文を組むこの経路でも、部分形が機械置換辞書へ載るようにする。
	if _, err := r.e.deriveRunProperNouns(ctx, plugin); err != nil {
		return nil, err
	}
	// 確定した固有名で本文の機械置換辞書を組む（master_term ∪ proper_noun）。固有名反映後に呼ぶ。
	dict, err := r.e.LoadDictionary(ctx)
	if err != nil {
		return nil, err
	}
	refIndex, err := r.e.referenceIndex(ctx)
	if err != nil {
		return nil, err
	}
	tmpl, err := r.e.store.GetPromptTemplate(ctx)
	if err != nil {
		return nil, fmt.Errorf("プロンプトテンプレートの取得: %w", err)
	}
	instructionByKey, keyByRF, err := r.e.directiveLookups(ctx)
	if err != nil {
		return nil, err
	}
	narrations, err := r.e.store.ListUntranslatedNarrations(ctx, plugin)
	if err != nil {
		return nil, fmt.Errorf("未訳叙述文の取得: %w", err)
	}
	lines, err := r.e.store.ListUntranslatedLines(ctx, plugin)
	if err != nil {
		return nil, fmt.Errorf("未訳台詞の取得: %w", err)
	}
	// 口調ペルソナを最新化し、台詞ごとの口調指示を一括で組む（同期 Run と同じ）。
	if _, genErr := r.e.GeneratePersonas(ctx); genErr != nil {
		return nil, genErr
	}
	personas, err := r.e.LinePersonas(ctx, lines, instructionByKey[directiveTone], ToneDefaults{
		Generic: tmpl.GenericToneText,
		PC:      tmpl.PcToneText,
		PcSex:   tmpl.PcSex,
	})
	if err != nil {
		return nil, err
	}

	nPlanned, err := r.planNarrations(ctx, narrations, refIndex, dict, tmpl.BaseDirective, instructionByKey, keyByRF)
	if err != nil {
		return nil, err
	}
	lPlanned, err := r.planLines(ctx, lines, refIndex, dict, tmpl.BaseDirective, personas)
	if err != nil {
		return nil, err
	}
	return append(nPlanned, lPlanned...), nil
}

// planNarrations は未訳叙述文を送信計画へ積む。既存訳一致は AI を呼ばず即確定し、載せない（同期と同じ）。
func (r *BatchRunner) planNarrations(ctx context.Context, narrations []model.Narration, refIndex map[referenceKey]string,
	dict *dictionary.Dictionary, base string, instructionByKey map[string]string, keyByRF map[RecordKey]string) ([]batchplan.PlannedRequest, error) {
	var planned []batchplan.PlannedRequest
	for _, n := range narrations {
		if dest, ok := refIndex[referenceKey{Rec: n.Rec, Field: n.Field, Source: n.Source}]; ok {
			if err := r.e.store.UpdateNarrationDest(ctx, n.ID, dest, statusTranslated); err != nil {
				return nil, fmt.Errorf("叙述文の既訳流用: %w", err)
			}
			continue
		}
		instruction := instructionByKey[keyByRF[RecordKey{Rec: n.Rec, Field: n.Field}]]
		p, _ := r.e.composeBodyPrompt(base, instruction, n.Source, dict)
		planned = append(planned, batchplan.PlannedRequest{Kind: model.BatchKindNarration, RowID: n.ID, Prompt: p})
	}
	return planned, nil
}

// planLines は未訳台詞を送信計画へ積む。既存訳一致は AI を呼ばず即確定し、載せない（同期と同じ）。
func (r *BatchRunner) planLines(ctx context.Context, lines []model.Line, refIndex map[referenceKey]string,
	dict *dictionary.Dictionary, base string, personas map[int64]Persona) ([]batchplan.PlannedRequest, error) {
	var planned []batchplan.PlannedRequest
	for _, l := range lines {
		if dest, ok := refIndex[referenceKey{Rec: l.Rec, Field: l.Field, Source: l.Source}]; ok {
			if err := r.e.store.UpdateLineDest(ctx, l.ID, dest, statusTranslated); err != nil {
				return nil, fmt.Errorf("台詞の既訳流用: %w", err)
			}
			continue
		}
		p, _ := r.e.composeBodyPrompt(base, personas[l.ID].Directive, l.Source, dict)
		planned = append(planned, batchplan.PlannedRequest{Kind: model.BatchKindLine, RowID: l.ID, Prompt: p})
	}
	return planned, nil
}

// sendStage は計画を外部 batch として送り、外部 ID を記録し、送信行対応を永続する（固有名段・本文段で共通）。
func (r *BatchRunner) sendStage(ctx context.Context, batch provider.BatchTranslator, conn provider.Connection, modelName string, batchID int64, stage string, planned []batchplan.PlannedRequest) error {
	externalID, err := batch.SubmitBatch(ctx, conn, modelName, batchplan.BuildBatchRequests(planned))
	if err != nil {
		return fmt.Errorf("%s batch の送信: %w", stage, err)
	}
	if err := r.store.RecordBatchExternalID(ctx, batchID, stage, externalID); err != nil {
		return fmt.Errorf("外部 batch ID の記録: %w", err)
	}
	return r.persistRequests(ctx, batchID, externalID, planned)
}

// ProgressStatus は対象 plugin の進行状況を副作用なしで返す（状態確認）。
// 現段の外部 batch を PollBatch するだけで、dest 更新・batch 送信・段更新・DB 書き込みを一切しない。
// 現段の外部 ID が空の半端な進行と完了段は PollBatch もしない。進行が無い plugin は ok=false を返す。
func (r *BatchRunner) ProgressStatus(ctx context.Context, providerName string, conn provider.Connection, plugin string) (batchplan.BatchProgress, bool, error) {
	prog, ok, err := r.store.GetBatchProgression(ctx, plugin)
	if err != nil {
		return batchplan.BatchProgress{}, false, fmt.Errorf("batch 進行の確認: %w", err)
	}
	if !ok {
		return batchplan.BatchProgress{}, false, nil
	}
	batch, err := r.batchForProgress(ctx, providerName, conn, prog)
	if err != nil {
		return batchplan.BatchProgress{}, false, err
	}
	externalID := prog.ProperBatchID
	if prog.Stage == model.BatchStageBody {
		externalID = prog.BodyBatchID
	}
	hasCurrent := prog.Stage != model.BatchStageDone && externalID != ""
	var status provider.BatchStatus
	if hasCurrent {
		status, err = batch.PollBatch(ctx, conn, externalID)
		if err != nil {
			return batchplan.BatchProgress{}, false, fmt.Errorf("batch 状態確認: %w", err)
		}
	}
	return batchplan.BuildProgress(prog.Stage, hasCurrent, status), true, nil
}

// RefreshPlugin は対象 plugin の進行 1 件だけを反映する（前進）。
// 現段が完了していれば結果を取り込み、固有名段なら本文 batch を送る（中身は refreshOne）。
// 起動時・画面操作の時点だけ呼ぶ（常駐ポーリングはしない）。進行が無い plugin は何もしない。接続情報は都度渡す。
func (r *BatchRunner) RefreshPlugin(ctx context.Context, providerName string, conn provider.Connection, plugin string) error {
	prog, ok, err := r.store.GetBatchProgression(ctx, plugin)
	if err != nil {
		return fmt.Errorf("batch 進行の確認: %w", err)
	}
	if !ok {
		return nil
	}
	batch, err := r.batchForProgress(ctx, providerName, conn, prog)
	if err != nil {
		return err
	}
	return r.refreshOne(ctx, batch, conn, prog)
}

// refreshOne は 1 進行を反映する。現段の外部 batch を状態確認し、純粋核の判定で次行動を決める。
func (r *BatchRunner) refreshOne(ctx context.Context, batch provider.BatchTranslator, conn provider.Connection, prog model.BatchTranslation) error {
	externalID := prog.ProperBatchID
	if prog.Stage == model.BatchStageBody {
		externalID = prog.BodyBatchID
	}
	if externalID == "" {
		return nil // 現段の外部 ID が空（異常系）。触らない。
	}
	status, err := batch.PollBatch(ctx, conn, externalID)
	if err != nil {
		return fmt.Errorf("batch 状態確認: %w", err)
	}
	switch batchplan.DecideRefreshStep(prog.Stage, prog.ProperBatchID, prog.BodyBatchID, status.Done) {
	case batchplan.StepApplyProperThenSubmitBody:
		if err := r.applyResults(ctx, batch, conn, prog.Plugin, prog.ProperBatchID); err != nil {
			return err
		}
		return r.submitBodyBatch(ctx, batch, conn, prog.Model, prog.Plugin, prog.ID)
	case batchplan.StepApplyProperThenAdvance:
		// 本文 batch は送信済み（クラッシュ復帰）。固有名を反映し、段だけ進める（本文を再送しない）。
		if err := r.applyResults(ctx, batch, conn, prog.Plugin, prog.ProperBatchID); err != nil {
			return err
		}
		if err := r.store.AdvanceBatchStage(ctx, prog.ID, model.BatchStageBody); err != nil {
			return fmt.Errorf("進行段の更新（本文）: %w", err)
		}
		return nil
	case batchplan.StepApplyBodyThenComplete:
		if err := r.applyResults(ctx, batch, conn, prog.Plugin, prog.BodyBatchID); err != nil {
			return err
		}
		if err := r.store.AdvanceBatchStage(ctx, prog.ID, model.BatchStageDone); err != nil {
			return fmt.Errorf("進行段の更新（完了）: %w", err)
		}
		return nil
	default:
		return nil // StepWait / StepNothing。
	}
}

// applyResults は外部 batch の全結果を取得し、custom_id で対応する行へ書き戻す。
// 書き戻しの可否は同期と共有する純粋規則（DecideApply）で決める。据え置き（タグ欠落・失敗）は未訳のまま残す（再送信で回収）。
// batch は個別失敗で全体を止めない（同期の abort とは違い、Fatal も据え置きにする）。
func (r *BatchRunner) applyResults(ctx context.Context, batch provider.BatchTranslator, conn provider.Connection, plugin, externalID string) error {
	results, err := batch.FetchResults(ctx, conn, externalID)
	if err != nil {
		return fmt.Errorf("batch 結果取得: %w", err)
	}
	rows, err := r.pendingRows(ctx, plugin)
	if err != nil {
		return err
	}

	var skipped, bad int
	for _, res := range results {
		kind, id, decErr := batchplan.DecodeCustomID(res.CustomID)
		if decErr != nil {
			bad++
			continue
		}
		applied, err := r.applyOneResult(ctx, res, kind, id, rows)
		if err != nil {
			return err
		}
		if !applied {
			skipped++
		}
	}
	logBatchApply(ctx, externalID, len(results), skipped, bad)
	return nil
}

// pendingRows は反映対象の未訳行を種別ごとに id→原文 で引けるようにまとめる（タグ照合の再計算に使う）。
type pendingRows struct {
	proper    map[int64]string
	narration map[int64]string
	line      map[int64]string
}

// pendingRows は対象 plugin の未訳行（固有名・叙述文・台詞）を id→原文 の map にして返す。
func (r *BatchRunner) pendingRows(ctx context.Context, plugin string) (pendingRows, error) {
	propers, err := r.e.store.ListUntranslatedProperNouns(ctx, plugin)
	if err != nil {
		return pendingRows{}, fmt.Errorf("未訳固有名の取得: %w", err)
	}
	narrations, err := r.e.store.ListUntranslatedNarrations(ctx, plugin)
	if err != nil {
		return pendingRows{}, fmt.Errorf("未訳叙述文の取得: %w", err)
	}
	lines, err := r.e.store.ListUntranslatedLines(ctx, plugin)
	if err != nil {
		return pendingRows{}, fmt.Errorf("未訳台詞の取得: %w", err)
	}
	rows := pendingRows{
		proper:    make(map[int64]string, len(propers)),
		narration: make(map[int64]string, len(narrations)),
		line:      make(map[int64]string, len(lines)),
	}
	for _, x := range propers {
		rows.proper[x.ID] = x.Source
	}
	for _, x := range narrations {
		rows.narration[x.ID] = x.Source
	}
	for _, x := range lines {
		rows.line[x.ID] = x.Source
	}
	return rows, nil
}

// applyOneResult は 1 結果を種別に応じた行へ適用する。適用したら applied=true。
// 対象行が未訳一覧に無い（既に反映済み・別実行で確定済み）場合は applied=false で飛ばす（冪等）。
func (r *BatchRunner) applyOneResult(ctx context.Context, res provider.BatchResult, kind string, id int64, rows pendingRows) (bool, error) {
	switch kind {
	case model.BatchKindProper:
		source, ok := rows.proper[id]
		if !ok {
			return false, nil
		}
		return r.applyOne(res, source, func(dest string) error {
			return r.e.store.UpdateProperNounDest(ctx, id, dest, statusProvisional)
		})
	case model.BatchKindNarration:
		source, ok := rows.narration[id]
		if !ok {
			return false, nil
		}
		return r.applyOne(res, source, func(dest string) error {
			return r.e.store.UpdateNarrationDest(ctx, id, dest, statusProvisional)
		})
	case model.BatchKindLine:
		source, ok := rows.line[id]
		if !ok {
			return false, nil
		}
		return r.applyOne(res, source, func(dest string) error {
			return r.e.store.UpdateLineDest(ctx, id, dest, statusProvisional)
		})
	default:
		return false, nil // 未知の種別は飛ばす。
	}
}

// applyOne は 1 結果を DecideApply で判定し、確定なら update で書き戻す。据え置き（タグ欠落・失敗）は書かない。
// タグ照合は原文から実行時タグを引き直して行う（送信時のタグを永続せず再計算する）。
func (r *BatchRunner) applyOne(res provider.BatchResult, source string, update func(dest string) error) (bool, error) {
	_, tags := runtimetag.Mask(source)
	missing := runtimetag.CountMissing(res.Translation, tags)
	outcome := batchplan.DecideApply(res.Translation, res.Err, missing)
	if outcome.Kind != batchplan.ApplyConfirm {
		return false, nil // 据え置き（未訳のまま。再送信で回収）。
	}
	if err := update(outcome.Dest); err != nil {
		return false, fmt.Errorf("batch 結果の書き戻し: %w", err)
	}
	return true, nil
}

// persistRequests は送信した計画を送信行対応（custom_id ↔ 行）として永続する。反映時の照合に使う。
func (r *BatchRunner) persistRequests(ctx context.Context, batchID int64, externalID string, planned []batchplan.PlannedRequest) error {
	rows := make([]model.BatchRequest, len(planned))
	for i, p := range planned {
		rows[i] = model.BatchRequest{
			BatchID:         batchID,
			ExternalBatchID: externalID,
			CustomID:        batchplan.EncodeCustomID(p.Kind, p.RowID),
			Kind:            p.Kind,
			RowID:           p.RowID,
		}
	}
	if _, err := r.store.InsertBatchRequests(ctx, rows); err != nil {
		return fmt.Errorf("送信行対応の記録: %w", err)
	}
	return nil
}

// logBatchApply は反映 1 回の集計（総数・据え置き・壊れた custom_id）を観測ログへ出す。
// 据え置きや壊れた対応は再送信での回収・原因分離が要るため、件数を集約して 1 度だけ出す（loop 内で 1 件ずつ出さない）。
func logBatchApply(ctx context.Context, externalID string, total, skipped, bad int) {
	slog.InfoContext(ctx, "batch 結果を反映した",
		slog.String("event", "batch_apply"),
		slog.String("external_batch_id", externalID),
		slog.Int("total", total),
		slog.Int("skipped", skipped),
		slog.Int("bad_custom_id", bad),
	)
}
