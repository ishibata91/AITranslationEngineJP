-- C#↔Go 契約の SQL schema 正本（repo-owned migration）の 13 本目。
-- gemini-xai-batch-translation: 非同期の大量翻訳（batch）の橋渡し情報を永続する土台を足す。
-- 足すテーブル:
--   batch_translation … 対象 plugin 1 つの batch 進行の本体（外部 batch ID・進行段・model）。plugin と 1 対 1。
--   batch_request     … 送信したリクエストと結果の対応（custom_id ↔ 翻訳対象行）。batch_translation の子。
--
-- 設計（docs/exec-plans/active/gemini-xai-batch-translation/）:
--   - batch は「送信」と「後日反映」の 2 時点へ割れる非同期配送。2 時点の橋渡しに外部 batch ID と
--     送信行の対応を永続する。同期経路（narration・line・proper_noun の status/dest 直接更新）は本テーブルを経由しない。
--   - 進行は固有名 batch → 本文 batch の 2 段を逐次でたどる（固有名を本文より先に確定して注入するため）。
--     stage は 'proper_noun'（固有名 batch 進行中）/ 'body'（本文 batch 進行中）/ 'done'（完了）。
--     proper_batch_id / body_batch_id は各段の外部 batch ID（未送信は空文字）。二重送信は外部 ID の有無で防ぐ。
--   - 同一 plugin に進行中は 1 進行だけ（UNIQUE(plugin)）。やり直し（再送信）は同じ行を reset して使い回す。
--   - custom_id は「種別:id」の複合キー（種別 n=叙述文 / l=台詞 / p=固有名）。叙述文と台詞は id が
--     テーブルごと独立採番で衝突するため、種別接頭で名前空間化する（結果ページングの cursor と同じ流儀）。
--     kind / row_id は custom_id の分解を列で持ち、結果適用時の行特定に使う。
--   - plugin 列は narration.plugin 等と同値の plugin ファイル名で束ねる。対象 plugin の削除（Go 側の
--     手続き DELETE）で batch_request（子）→ batch_translation（親）の順に消す（migration 0010 と同方式）。
-- 注記: C# 抽出器は本 migration も毎回 ensure するが、追加のみ（CREATE TABLE IF NOT EXISTS）で冪等。
--   接続情報（endpoint・API key）は永続しない。反映のたびに UI から渡す。論理 ER は docs/er.md。

-- batch_translation: 対象 plugin 1 つの batch 進行の本体。plugin と 1 対 1。
CREATE TABLE IF NOT EXISTS batch_translation (
    id INTEGER PRIMARY KEY,
    plugin TEXT NOT NULL UNIQUE,
    model TEXT NOT NULL,
    stage TEXT NOT NULL,
    proper_batch_id TEXT NOT NULL DEFAULT '',
    body_batch_id TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- batch_request: 送信したリクエストと結果の対応。custom_id は「種別:id」。batch_translation の子。
CREATE TABLE IF NOT EXISTS batch_request (
    id INTEGER PRIMARY KEY,
    batch_id INTEGER NOT NULL,
    external_batch_id TEXT NOT NULL,
    custom_id TEXT NOT NULL,
    kind TEXT NOT NULL,
    row_id INTEGER NOT NULL,
    UNIQUE (batch_id, custom_id)
);

CREATE INDEX IF NOT EXISTS idx_batch_request_external ON batch_request (external_batch_id);
