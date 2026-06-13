# Task Plan: 2026-06-14-minimal-extract-translate

- `workflow`: work
- `status`: planned
- `task_id`: 2026-06-14-minimal-extract-translate
- `task_mode`: 重 task（画面が動く）
- `request_summary`: 最小縦切り。1 つの plugin を抽出 → SQLite → engine が読む → OpenAI 互換 provider で叙述文 1 種を翻訳 → dest を SQLite へ戻し → 最小実行画面に結果を表示するまでを 1 本通す。
- `goal`: 抽出 → 翻訳 → DB → 画面表示の end-to-end を 1 本通し、GUI から翻訳を走らせて結果を見られる最小の価値を出す。
- `constraints`: ペルソナ・マスター辞書・重複排除の作り込み・provider 4 系統・全レコード種別・xTranslator XML 出力・設定/監視画面・進捗 push の作り込みは対象外。
- `close_conditions`: GUI で plugin を指定し実行すると、翻訳された叙述文の原文・訳文が画面に出る。SQLite の `narration.dest` に訳文が入る。
- `execution_branch`: `claude/2026-06-14-minimal-extract-translate`（作成済み）
- `source_branch`: `master`
- `branch_point`: `06503b66641bd5793feae700040cf442e3d73e2a`
- `target_branch`: `master`

## 軽 / 重判定

- 画面が動くか: Y（最小実行画面を作る）。
- `docs/architecture.md` 反映が要るか: N 想定（既存骨格の初実装で、層構成・依存方向は `architecture.md` に固定済み）。design-module 入口で境界の具体が要るか再評価する。
- 判定: 重 task。フロー = `preparation → design → storybook →（画面表示）→ implementation → finalization`。

## Scope（含む / 含まない）

含む:
- `tools/extractor`（C#/Mutagen）に SQLite writer を追加し、1 plugin から叙述文 1 種（候補: 武器・防具の DESC、または BOOK 本文）を抽出して `er.md` の `narration` テーブルへ書く。
- `db/` に初期 SQL migration（最小: `narration`、必要なら識別カラムだけ）。
- `internal/store`（sqlx 薄アクセス）: `narration` の読み書き。
- `internal/engine`: `narration` を読み、provider で翻訳し `dest` を書く最小手続き。重複排除なし、文体無視で素直に投げる。
- `internal/provider`: AI クライアント interface ＋ OpenAI 互換 1 実装（`net/http`）。
- `internal/api` ＋ `internal/bootstrap`: Bind 公開（抽出+翻訳実行 command、結果取得 query）、composition root。
- `internal/model`: `narration` 等のデータ構造。
- `frontend/src/ui`: 最小実行画面（plugin パス入力、provider 接続設定 endpoint/key/model、実行ボタン、結果一覧）。Store。
- `frontend/src/gateway`: Wails bindings ラッパ。

含まない:
- 固有名の重複排除（訳の単位化）、`proper_noun`/`set_phrase`/`placement` の本格運用。
- ペルソナ（`speaker`/`race`/`faction`/`voice_type`、口調指示）。
- マスター辞書。
- provider の Gemini / xAI / Claude 実装。
- xTranslator XML 出力。
- 設定画面・監視画面・進捗 push の作り込み。

## Routing Notes

- `required_reading`:
  - `docs/concept-model.md`（概念正本。叙述文の扱い）
  - `docs/er.md`（テーブル定義。今回は `narration` 中心）
  - `docs/architecture.md`（§3 各責務、§6 C#↔Go 境界、§7 ディレクトリ正本）
  - `docs/system_requirements.md`（§1 AI 翻訳、対応 provider）
  - `docs/tech-selection.md`（採用技術）
  - `docs/references/xtranslator_ref.md`（識別キー・status 値域）
  - `docs/coding-guidelines-backend.md` / `-frontend.md` / `-tests.md`
- `canonicalization_targets`: `docs/architecture.md`（必要時のみ。原則 N）
- `validation_commands`: backend `go test ./...` / frontend `Vitest` / `svelte-check`（実装時に確定）

## 着手前メモ（compact 跨ぎ）

- これは 4 タスク計画の T1。依存関係と T2〜T4 の輪郭は各 plan.md 参照。
- T1 の最小縦切りの確定事項: 叙述文 1 種、provider は OpenAI 互換、画面込み（重 task）。
- 次アクション: master の未コミット docs 変更（ER 設計・掃除）を整理し、`preparation-module` で branch `claude/2026-06-14-minimal-extract-translate` を切って着手。実装本体は Claude 本体が 1 文脈で縦通し（`coding-protocol` / `implementation-module`）。

## Branch Status

- `branch_ready`: yes（`claude/2026-06-14-minimal-extract-translate`、分岐元 `06503b66`）
- `remote_operation`: `not-performed`

## HITL Status

- `detail_spec_hitl`: `required-after-design-bundle`
- `frontend_human_review`: `required-after-storybook-review-loop-evidence`

## Outcome

- 未着手。
