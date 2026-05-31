# Task Plan: fix-lucien-target-list-empty

- `workflow`: fix
- `status`: in-progress
- `lane_owner`: fix_lane
- `task_id`: fix-lucien-target-list-empty
- `task_mode`: fix
- `request_summary`: `dictionaries/Lucien.esp_Export.json` をデータロードで読み込むと、単語翻訳画面へ遷移でき、進捗パネルは4900件と表示するが、処理対象一覧パネルは0件と表示する。
- `goal`: 処理対象一覧パネルに、読み込んだ翻訳対象が表示され、進捗パネルの母数と整合する状態へ修正する。
- `constraints`: 差分最小化を目的に責務境界を壊さない。進捗パネルのカウント削減やフィルタ機能削除は禁止。
- `close_conditions`: 成果物DAGの成果物が揃い、harness が通過する。
- `worktree_path`: N/A（カレント作業ツリーで実施）
- `execution_branch`: `codex/fix-lucien-target-list-empty`
- `source_branch`: `codex/fix-lucien-target-list-empty`
- `target_branch`: `master`

## Artifact Index

- `観測記録`: `./observation.md`
- `修正方針判断`: `./fix-decision.md`（fix_decider が作成）
- `UC 差分候補`: `./uc-diff.md`（test_designer が作成）
- `E2E テスト観点差分`: `./e2e-diff.md`（test_designer が作成）
- `修正実行入力`: `./fix-input.md`（人間レビュー承認後に fix_lane が作成）
- `シナリオテスト追加証跡`: `./scenario-test.md`
- `実装修正証跡`: `./implementation.md`
- `単体テスト追加証跡`: `./unit-test.md`
- `実装後ブラウザ確認`: `./browser-confirmation.md`

## Routing Notes

- `required_reading`: `docs/index.md`, `docs/usecases/README.md`, `docs/screen-design/README.md`, `docs/observability-logging.md`
- `validation_commands`: `python3 scripts/harness/run.py --suite all`

## Branch Status

- `worktree_checkout`: current tree
- `branch_ready`: true
- `execution_branch`: `codex/fix-lucien-target-list-empty`
- `commit_hash`:
- `remote_operation`: `not-performed`

## Preparation

- `wails_process`: `npm run dev:wails:agent-browser`（PATH に `$HOME/go/bin` を前置して起動）。
- `wails_connect_target`: `http://localhost:34115`（Wails dev bridge）。
- `wails_note`: 非対話 shell では `wails` バイナリ（`/Users/iorishibata/go/bin/wails`）が PATH 外のため、PATH 前置が必要。

## HITL Status

- `human_fix_review`: `required-after-fix-decision-and-test-diffs`

## Outcome

- 判定: 未解決のまま close（2026-05-30）。fix-lane（局所修正）で7ラウンド原因究明したが、症状（単語翻訳画面の初回表示で処理対象一覧が0件）を解消できなかった。
- 理由: 原因が取得タイミング・bridge IPC 飽和（画面経由で約15秒遅延）・reactive 反映の複数層に絡み、局所パッチの積み上げでは収束しなかった。設計レベルの作り直しが必要と判断した。
- 局所修正は全 revert 済み（プロダクトコード変更なし）。
- 後継: `docs/exec-plans/active/job-run-phase-fetch-redesign/`。term/persona/body 3段階を貫く取得・表示フローの共通設計で作り直す。設計は専用レーン整備後に行う。
- 残置資産（後継へ引き継ぎ）: 不足テスト `tests/system/fix-lucien-target-list-empty.spec.ts`、E2E mock 追加 `tests/system/support/scenario-wails-mocks.ts`。
- 調査証跡: 本フォルダの `fix-decision.md`（追補1〜5）、`attempts/observe`・`attempts/obs-round2`・`attempts/obs-r4`〜`obs-r6`、`observation.md`、`uc-diff.md`、`e2e-diff.md`、`test-design.csv`、`data-testid-gaps.md`。
