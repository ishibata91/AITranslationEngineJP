# 実装計画

- workflow: impl
- status: completed
- lane_owner: skill-modification
- scope: .codex/skills, .codex/workflow.md, package.json, scripts/dev
- task_id: 2026-04-12-reproduce-issues-wails-log-path
- task_catalog_ref: N/A
- parent_phase: N/A

## 要求要約

- `reproduce-issues` が Wails の Go 側ログを安定して読めるようにする。
- `npm run dev:wails:docker-mcp` を固定ログ path 付き起動へ変える。
- ログ file は起動のたびに削除し、溜まり続けないようにする。

## 判断根拠

- 現状の `dev:wails:docker-mcp` は `wails dev` を直接起動しており、skill が機械的に参照できるログ file path がない。
- `scripts/test/run-system-test.sh` には `wails dev` の stdout / stderr を file に退避する先例がある。
- fix lane の contract には `wails_log_targets` があり、既定 path を固定すると `reproduce-issues` と `logging-fixes` の handoff が安定する。

## 対象範囲

- `package.json`
- `scripts/dev/`
- `.codex/skills/reproduce-issues/`
- `.codex/skills/orchestrating-fixes/`
- `.codex/skills/logging-fixes/`
- `.codex/workflow.md`

## 対象外

- product code
- app 内 custom logger の新設
- `docs/` 正本更新

## 依存関係・ブロッカー

- なし

## 並行安全メモ

- 変更対象は repo root の起動 script と `.codex/` の workflow 文書に限定する。
- 既存の `test:system` 用 script は保持し、再現用の正本入口だけを `dev:wails:docker-mcp` に集約する。

## 機能要件

- `summary`: `dev:wails:docker-mcp` 起動時に `tmp/logs/wails-dev.log` を毎回削除してから当該起動の stdout / stderr を記録できるようにする。
- `in_scope`: 固定ログ path の追加、起動前削除、skill / contract / workflow の既定 path 明記。
- `non_functional_requirements`: ログ file は起動毎にリセットし、`reproduce-issues` が MCP read だけで確認できる。
- `out_of_scope`: Wails runtime logger の差し替え、別ログ保管方式、履歴保持。
- `open_questions`: なし。
- `required_reading`: `package.json`, `scripts/test/run-system-test.sh`, `.codex/skills/reproduce-issues/SKILL.md`, `.codex/workflow.md`

## UI モック

- N/A

## Scenario テスト一覧

- N/A

## 実装計画

- `parallel_task_groups`:
  - `group_id`: reproduce-issues-wails-log-path
  - `can_run_in_parallel_with`: none
  - `blocked_by`: none
  - `completion_signal`: 起動入口、skill、contract、workflow が `tmp/logs/wails-dev.log` を既定 path として参照する。
- `tasks`:
  - `task_id`: add-wails-dev-wrapper
  - `owned_scope`: `scripts/dev/`, `package.json`
  - `depends_on`: none
  - `parallel_group`: reproduce-issues-wails-log-path
  - `required_reading`: `package.json`, `scripts/test/run-system-test.sh`
  - `validation_commands`: `npm run dev:wails:docker-mcp`
  - `task_id`: sync-reproduce-issues-contracts
  - `owned_scope`: `.codex/skills/reproduce-issues/`, `.codex/skills/orchestrating-fixes/`, `.codex/skills/logging-fixes/`, `.codex/workflow.md`
  - `depends_on`: add-wails-dev-wrapper
  - `parallel_group`: reproduce-issues-wails-log-path
  - `required_reading`: `.codex/skills/reproduce-issues/SKILL.md`, `.codex/workflow.md`
  - `validation_commands`: `python3 scripts/harness/run.py --suite structure`

## 受け入れ確認

- `npm run dev:wails:docker-mcp` が `tmp/logs/wails-dev.log` を起動前に削除してから新しく出力する。
- `reproduce-issues` が既定の Wails ログ path と起動毎削除を明示している。
- `orchestrating-fixes` と `logging-fixes` の contract が `tmp/logs/wails-dev.log` を既定 path として渡せる。
- `.codex/workflow.md` が fix lane の Wails ログ既定 path を説明している。

## 必要な証跡

- `python3 scripts/harness/run.py --suite structure`
- `tmp/logs/wails-dev.log` の生成確認

## 機能要件 HITL 状態

- N/A

## 機能要件 承認記録

- N/A

## 詳細設計 HITL 状態

- N/A

## 詳細設計 承認記録

- N/A

## review 用差分図

- N/A

## 差分正本適用先

- `.codex/skills/`
- `.codex/workflow.md`
- `package.json`
- `scripts/dev/`

## Closeout Notes

- 完了後は plan を `docs/exec-plans/completed/` へ移す。

## 結果

- `scripts/dev/run-wails-docker-mcp.sh` を追加し、`tmp/logs/wails-dev.log` を起動前に削除してから当該起動の Wails stdout / stderr を書く入口へ置き換えた
- `package.json` の `dev:wails:docker-mcp` を新しい wrapper script 呼び出しへ変更した
- `reproduce-issues`、`orchestrating-fixes`、`logging-fixes`、`workflow.md` を更新し、既定の Wails ログ path を `tmp/logs/wails-dev.log` に固定した
- `python3 scripts/harness/run.py --suite structure` が通過した
- 検証で `tmp/logs/wails-dev.log` に事前に書いた `stale-marker` が消え、起動時削除を確認した
- 環境依存の既知事項として、検証時に `Port 5173 is already in use` が出た
