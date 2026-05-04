# Exploration Test Data: exploration-normal-flow-20260503

- `skill`: exploration-test-lane
- `status`: complete
- `source_plan`: `./plan.md`
- `source_exploration_plan`: `./exploration-test-plan.md`
- `owner_agent`: exploration_test_lane

## Data Sets

- `inputs`:
  - `./normal-flow-lucien-mini.json`
  - `./normal-flow-lucien-reviewfix-mini.json`
  - `./normal-flow-lucien-complete-mini.json`
  - `target_plugin`: `Lucien.esp`
  - `dialogue_groups`: 1件
  - `records`: `DIAL FULL` 1件、`INFO NAM1` 1件
- `initial_state`:
  - アプリ起動直後の `ダッシュボード` から開始する。
  - 翻訳ジョブ未作成状態を前提にする。
  - 入力 import で重複 hash が発生した場合は、同一 DB の既存入力が残っている状態として記録し、通常フロー観測を停止する。
- `fixtures`:
  - task 内 fixture: `./normal-flow-lucien-mini.json`
  - task 内 fixture: `./normal-flow-lucien-reviewfix-mini.json`
  - task 内 fixture: `./normal-flow-lucien-complete-mini.json`
  - 入力形式根拠: `internal/service/translation_input_import_service_test.go`
  - 画面導線根拠: `frontend/src/ui/stores/shell-state.ts`
- `accounts_or_permissions`:
  - 追加アカウントは使わない。
  - API key や secret は使わない。
- `external_conditions`:
  - Wails dev server を `npm run dev:wails:agent-browser` で起動する。
  - UI から観測できる状態表示だけを探索証跡に残す。
  - 実 AI provider の品質比較は行わない。

## Mapping To Plan

- `observation_target_coverage`:
  - 区間1: `ダッシュボード` 既定表示を確認する。
  - 区間2: `翻訳管理` の `Input Review` で `normal-flow-lucien-mini.json` を使う。
  - 区間3: `Job Setup` で import 済み入力から ready job 作成を確認する。
  - 区間4: `Job Run` で term phase、persona phase、body phase の順序と状態接続を確認する。
  - 区間5: `出力管理` で output readiness と成果物確認導線を確認する。
- `viewpoint_coverage`:
  - 途中停止: 次区間へ進めない箇所を観測する。
  - 前提未充足操作: 前区間完了前に後続区間を実行できるかを観測する。
  - 状態接続: `Draft -> Ready -> Running -> Completed` の UI 表示を観測する。
  - 再入場: 同一区間の再表示で続行不能が起きるかを観測する。
- `stop_condition_coverage`:
  - fixture import 不可、dev server 起動不可、通常導線不明、出力確認範囲不明を停止条件にする。

## Reuse And Cleanup

- `reuse_policy`:
  - `./normal-flow-lucien-mini.json`、`./normal-flow-lucien-reviewfix-mini.json`、`./normal-flow-lucien-complete-mini.json` は task 内探索データとして再利用する。
  - repo 正本 docs やプロダクトテストへ昇格しない。
- `cleanup_policy`:
  - runtime が生成した一時 DB、ログ、スクリーンショットは `tmp/` または `test-results/` に置く。
  - 探索終了時に task artifact から参照できる証跡だけ残す。
- `forbidden_data`:
  - 実 API key、実 secret、外部 paid API 実行、巨大入力、異常入力、複数入力を使わない。

## Output

- `decision`: complete
- `evidence_refs`:
  - `./exploration-test-plan.md`
  - `./normal-flow-lucien-mini.json`
  - `./normal-flow-lucien-reviewfix-mini.json`
  - `./normal-flow-lucien-complete-mini.json`
  - `internal/service/translation_input_import_service_test.go`
  - `frontend/src/ui/stores/shell-state.ts`
- `missing_info`:
  - `翻訳管理` と `出力管理` の正常系 scenario 正本は未整備である。
  - 実 AI provider なしで phase 完了まで到達できるかは探索証跡で確認する。
- `next_artifact`: `./exploration-test-evidence.md`
