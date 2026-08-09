package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"aitranslationenginejp/internal/model"
)

// UpsertTargetPlugin は翻訳した対象 plugin を登録する。翻訳開始時（抽出の前）に 1 度呼ぶ。
// 初回は plugin・source_path・created_at を書き、2 回目以降は source_path を更新して
// sync_retry_ready を解除する。created_at（初回登録時刻）は保つ。plugin ファイル名がキーで plugin と 1 対 1。
func (s *Store) UpsertTargetPlugin(ctx context.Context, plugin, sourcePath string) error {
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO target_plugin (plugin, source_path, created_at, sync_retry_ready) VALUES (?, ?, ?, 0)
		 ON CONFLICT(plugin) DO UPDATE SET source_path = excluded.source_path, sync_retry_ready = 0`,
		plugin, sourcePath, createdAt)
	if err != nil {
		return fmt.Errorf("target_plugin の登録: %w", err)
	}
	return nil
}

// IsSyncRetryReady は保存済みの準備結果を使って未訳だけを同期翻訳できるかを返す。
// 未登録の plugin は準備未完了として false を返す。
func (s *Store) IsSyncRetryReady(ctx context.Context, plugin string) (bool, error) {
	var ready int
	err := s.db.GetContext(ctx, &ready, `SELECT sync_retry_ready FROM target_plugin WHERE plugin = ?`, plugin)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("同期再実行準備状態の取得: %w", err)
	}
	return ready == 1, nil
}

// MarkSyncRetryReady は既訳収集、抽出、辞書派生、取込、口調集計が完了した状態を保存する。
func (s *Store) MarkSyncRetryReady(ctx context.Context, plugin string) error {
	result, err := s.db.ExecContext(ctx, `UPDATE target_plugin SET sync_retry_ready = 1 WHERE plugin = ?`, plugin)
	if err != nil {
		return fmt.Errorf("同期再実行準備状態の保存: %w", err)
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("同期再実行準備状態の更新件数取得: %w", err)
	}
	if changed != 1 {
		return fmt.Errorf("同期再実行準備状態の保存: 対象 plugin が未登録: %s", plugin)
	}
	return nil
}

// translationTargetTables は 1 つの対象 plugin について翻訳対象を数える表。
// 叙述文（narration）・台詞（line）・固有名（proper_noun）を数え合わせる。alias は状態列と絞りの前置。
// extra は表ごとの追加の絞りで、固有名は機械派生した人名の部分形を翻訳対象から外す条件を持つ。
// 数え合わせる表と外す条件をここ 1 箇所で持ち、一覧の進捗と対象 plugin 1 件の未訳件数が別々の数え方にならないようにする。
var translationTargetTables = []struct{ table, alias, extra string }{
	{table: "narration", alias: "n"},
	{table: "line", alias: "l"},
	{table: "proper_noun", alias: "p", extra: translationTargetProperNounAsP},
}

// translationTargetCountExpr は翻訳対象の件数を数える SQL 式（3 表の COUNT の和）を組む。
// pluginRef は plugin 値の参照で、一覧は相関副問い合わせの列（tp.plugin）、対象 plugin 1 件の集計は placeholder（?）を渡す。
// statusCond は状態の絞りで、空なら総数、`status != 0` なら訳済み、`status = 0` なら未訳を数える（前置の alias は本関数が付ける）。
func translationTargetCountExpr(pluginRef, statusCond string) string {
	terms := make([]string, 0, len(translationTargetTables))
	for _, t := range translationTargetTables {
		conds := []string{t.alias + ".plugin = " + pluginRef}
		if statusCond != "" {
			conds = append(conds, t.alias+"."+statusCond)
		}
		if t.extra != "" {
			conds = append(conds, t.extra)
		}
		terms = append(terms,
			"(SELECT COUNT(*) FROM "+t.table+" "+t.alias+" WHERE "+strings.Join(conds, " AND ")+")")
	}
	return strings.Join(terms, " + ")
}

// targetPluginListQuery は登録済み plugin を、束ねる翻訳対象の総数（total）と訳済み数（translated）付きで返す。
// 訳済みは status != 0（未訳は status = 0）。新しい登録が先頭に来るよう created_at 降順で並べる。
var targetPluginListQuery = `
SELECT tp.plugin, tp.source_path, tp.created_at,
  ` + translationTargetCountExpr("tp.plugin", "") + ` AS total,
  ` + translationTargetCountExpr("tp.plugin", "status != 0") + ` AS translated
FROM target_plugin tp
ORDER BY tp.created_at DESC, tp.plugin`

// targetPluginUntranslatedQuery は対象 plugin 1 件の未訳件数（status = 0）を返す。
// 3 表それぞれの副問い合わせが plugin 値の placeholder を 1 つ持つため、同じ plugin 値を 3 回渡す。
var targetPluginUntranslatedQuery = `SELECT ` + translationTargetCountExpr("?", "status = 0")

// ListTargetPlugins は登録済みの対象 plugin を一覧する（新しい順）。各行に進捗（total / translated）を付ける。
func (s *Store) ListTargetPlugins(ctx context.Context) ([]model.TargetPlugin, error) {
	var rows []model.TargetPlugin
	if err := s.db.SelectContext(ctx, &rows, targetPluginListQuery); err != nil {
		return nil, fmt.Errorf("target_plugin の一覧取得: %w", err)
	}
	return rows, nil
}

// CountUntranslated は対象 plugin に残る未訳（status = 0）の件数を返す。
// 翻訳の実行は未訳の全件を対象に取るため、実行後のこの件数がその実行で未訳のまま残した件数になる。
// 数え方は一覧の進捗（total / translated）と同じ表・同じ絞りを共有する。
func (s *Store) CountUntranslated(ctx context.Context, plugin string) (int, error) {
	return s.count(ctx, targetPluginUntranslatedQuery, plugin, plugin, plugin)
}

// targetPluginDeleteStmts は対象 plugin の翻訳成果を消す削除文の並び。
// 連関（親行 id の副問い合わせ）を先に消し、続いて親（対象スコープ本体）・staging・登録本体を消す。
// 共有資産（speaker / race / faction / voice_type / master_term / 各キャッシュ・seed）は消さない。
// %[1]s は narration の plugin 列、%[2]s は line の plugin 列など、すべて対象 plugin 値で束ねる（info_plugin も同値）。
var targetPluginDeleteStmts = []string{
	// narration の連関（本文言及・説明対象）を、対象 plugin の narration に紐づく分だけ消す。
	`DELETE FROM narration_mention   WHERE narration_id IN (SELECT id FROM narration WHERE plugin = ?)`,
	`DELETE FROM narration_described WHERE narration_id IN (SELECT id FROM narration WHERE plugin = ?)`,
	// line の連関（本文言及・話者連関・条件由来の性別）を、対象 plugin の line に紐づく分だけ消す。
	`DELETE FROM line_mention   WHERE line_id IN (SELECT id FROM line WHERE plugin = ?)`,
	`DELETE FROM line_speaker   WHERE line_id IN (SELECT id FROM line WHERE plugin = ?)`,
	`DELETE FROM line_condition WHERE line_id IN (SELECT id FROM line WHERE plugin = ?)`,
	// 対象スコープの本体（叙述文・台詞・固有名・素朴吸い出し）。
	`DELETE FROM narration       WHERE plugin = ?`,
	`DELETE FROM line            WHERE plugin = ?`,
	`DELETE FROM proper_noun     WHERE plugin = ?`,
	`DELETE FROM extracted_field WHERE plugin = ?`,
	// staging（INFO→speaker / INFO→条件由来の性別）。info_plugin は対象 plugin 値。
	`DELETE FROM extracted_info_speaker   WHERE info_plugin = ?`,
	`DELETE FROM extracted_info_condition WHERE info_plugin = ?`,
	// 本文送信時の参考語snapshot。batch requestと進行本体より先に消す。
	`DELETE FROM translation_reference_snapshot WHERE plugin = ?`,
	// batch 進行（gemini-xai-batch-translation）。子（送信行対応）→ 親（進行本体）の順で消す。
	// batch_request は plugin 列を持たないため、対象 plugin の進行本体に紐づく分を副問い合わせで消す。
	`DELETE FROM batch_request     WHERE batch_id IN (SELECT id FROM batch_translation WHERE plugin = ?)`,
	`DELETE FROM batch_translation WHERE plugin = ?`,
	// 登録本体。
	`DELETE FROM target_plugin WHERE plugin = ?`,
}

// DeleteTargetPlugin は対象 plugin の翻訳成果を 1 トランザクションで消す。
// FK cascade は使わず、対象スコープの行と連関を明示 DELETE で順に消す（削除方式の確定判断）。
// 共有 entity・横断辞書・各キャッシュ・seed は残すため削除対象に含めない。
func (s *Store) DeleteTargetPlugin(ctx context.Context, plugin string) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("削除のトランザクション開始: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // commit 済みなら no-op。失敗時の後始末用。
	for _, stmt := range targetPluginDeleteStmts {
		if _, err := tx.ExecContext(ctx, stmt, plugin); err != nil {
			return fmt.Errorf("target_plugin の削除: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("削除のコミット: %w", err)
	}
	return nil
}
