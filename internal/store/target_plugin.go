package store

import (
	"context"
	"fmt"
	"time"

	"aitranslationenginejp/internal/model"
)

// UpsertTargetPlugin は翻訳した対象 plugin を登録する。翻訳開始時（抽出の前）に 1 度呼ぶ。
// 初回は plugin・source_path・created_at を書き、2 回目以降は source_path だけ更新して
// created_at（初回登録時刻）を保つ。plugin ファイル名がキーで plugin と 1 対 1。冪等。
func (s *Store) UpsertTargetPlugin(ctx context.Context, plugin, sourcePath string) error {
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO target_plugin (plugin, source_path, created_at) VALUES (?, ?, ?)
		 ON CONFLICT(plugin) DO UPDATE SET source_path = excluded.source_path`,
		plugin, sourcePath, createdAt)
	if err != nil {
		return fmt.Errorf("target_plugin の登録: %w", err)
	}
	return nil
}

// targetPluginListQuery は登録済み plugin を、束ねる翻訳対象の総数（total）と訳済み数（translated）付きで返す。
// total / translated は叙述文（narration）・台詞（line）・固有名（proper_noun）を plugin で数え合わせる。
// 訳済みは status != 0（未訳は status = 0）。新しい登録が先頭に来るよう created_at 降順で並べる。
const targetPluginListQuery = `
SELECT tp.plugin, tp.source_path, tp.created_at,
  (SELECT COUNT(*) FROM narration   n WHERE n.plugin = tp.plugin)
  + (SELECT COUNT(*) FROM line        l WHERE l.plugin = tp.plugin)
  + (SELECT COUNT(*) FROM proper_noun p WHERE p.plugin = tp.plugin) AS total,
  (SELECT COUNT(*) FROM narration   n WHERE n.plugin = tp.plugin AND n.status != 0)
  + (SELECT COUNT(*) FROM line        l WHERE l.plugin = tp.plugin AND l.status != 0)
  + (SELECT COUNT(*) FROM proper_noun p WHERE p.plugin = tp.plugin AND p.status != 0) AS translated
FROM target_plugin tp
ORDER BY tp.created_at DESC, tp.plugin`

// ListTargetPlugins は登録済みの対象 plugin を一覧する（新しい順）。各行に進捗（total / translated）を付ける。
func (s *Store) ListTargetPlugins(ctx context.Context) ([]model.TargetPlugin, error) {
	var rows []model.TargetPlugin
	if err := s.db.SelectContext(ctx, &rows, targetPluginListQuery); err != nil {
		return nil, fmt.Errorf("target_plugin の一覧取得: %w", err)
	}
	return rows, nil
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
