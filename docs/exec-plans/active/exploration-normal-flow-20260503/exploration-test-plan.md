# Exploration Test Plan: exploration-normal-flow-20260503

- `skill`: exploration-test-planning
- `status`: complete
- `source_plan`: `./plan.md`
- `owner_agent`: exploration_test_planner
- `return_to`: exploration_test_lane

## Source Inputs

- `request_summary`: 通常フローを一貫して実行する探索テストを行う。
- `existing_artifacts`: `./plan.md`
- `human_constraints`: 観測、ログ確認、画面確認、原因仮説を計画段階で行わない。プロダクトコード、プロダクトテスト、docs 正本本文、`.codex/` を変更しない。
- `excluded_targets`: 中断、再開、失敗回復、キャンセル、回帰確認、恒久修正、バグ一覧作成。

## Observation Scope

- `targets`:
  - 区間1: アプリ起動からダッシュボード既定表示まで。
  - 区間2: ダッシュボードから `翻訳管理` へ入り、`Input Review` で入力登録と cache 再構築確認まで。
  - 区間3: `Job Setup` で validation と ready job 作成確認まで。
  - 区間4: `Job Run` で term phase、persona phase、body phase を正常順で進め、output readiness 到達まで。
  - 区間5: `出力管理` で生成成果物と書き出し結果の確認まで。
- `entry_points`:
  - アプリ起動直後の `ダッシュボード`。
  - `ダッシュボード` からの `翻訳管理` 導線。
  - `翻訳管理` 内の `Input Review`、`Job Setup`、`Job Run` 切替。
  - `出力管理` への主要導線。
- `state_or_log_targets`:
  - ジョブ状態は正常系の `Draft -> Ready -> Running -> Completed` に限定する。
  - 各区間で確認対象とする状態は、入力受付可否、validation 通過、ready job 作成済み、phase 進行中、phase 完了、output readiness、成果物出力済みとする。
  - ログ本文は計画対象に含めない。探索証跡では、必要時に UI から観測できる状態表示だけを扱う。
- `out_of_scope`:
  - `マスター辞書`、`マスターペルソナ` の個別機能詳細。
  - `Paused`、`RecoverableFailed`、`Failed`、`Canceled` への遷移。
  - 複数ジョブ並行管理。
  - AI プロバイダ切替比較。
  - 出力内容の翻訳品質評価。

## Exploration Viewpoints

- `failure_viewpoints`:
  - 正常入力でも次区間へ遷移できず通常フローが途中停止しないか。
  - 前区間の完了前に後続区間へ進めてしまわないか。
  - phase 完了後に output readiness または成果物確認へ接続できない断絶がないか。
- `state_transition_viewpoints`:
  - `ダッシュボード -> 翻訳管理 -> Input Review -> Job Setup -> Job Run -> 出力管理` の導線が一貫しているか。
  - `Draft -> Ready -> Running -> Completed` の正常系遷移を UI 上の状態変化で追えるか。
  - `Job Run` 内で term phase -> persona phase -> body phase の順序前提が崩れないか。
- `recovery_viewpoints`:
  - 通常フローでは回復操作を対象外とする。
  - 同一区間の再表示や再入場で正常系の続行不能が起きないかだけを確認対象に含める。
- `permission_or_trust_viewpoints`:
  - 通常フローに必要な入力、job 作成、phase 実行、出力確認の各操作が単一路線で到達できるか。
  - 前提未充足のまま実行操作を許していないか。
- `log_viewpoints`:
  - ログ本文は対象外とする。
  - 探索証跡では、UI 上で利用者が観測できる進捗、ready 状態、完了状態、成果物有無の表示だけを扱う。

## Test Data Policy

- `required_inputs`:
  - 1件の正常な xEdit 抽出入力データを用意する。
  - 通常フローを最後まで通すため、単一ジョブで完結する最小入力を優先する。
  - 本文翻訳フェーズまで進めるため、少なくとも用語対象と本文対象を含む入力が望ましい。
- `required_state`:
  - アプリ起動直後の既定状態から開始する。
  - 正常系のジョブ未作成状態から開始する。
  - 共通辞書と共通ペルソナは必須前提にしない。未使用でも通常フローが完走できる構成を優先する。
- `reusable_fixtures`:
  - 既存の repo 内 fixture または sample input があれば 1件だけ再利用する。
  - `Input Review`、`Job Setup`、`Job Run`、`出力管理` を同一入力で連結できる fixture を優先する。
- `data_constraints`:
  - 具体値はテストデータ成果物で固定する。
  - 複数入力、異常入力、境界値入力、巨大入力は通常フロー計画から外す。
  - 翻訳品質比較用の期待訳は不要とする。通常フローでは完走性と状態接続を優先する。

## Stop Conditions

- `sufficient_evidence`:
  - 区間1から区間5までの各区間で、入口、主要操作、次区間への接続条件、正常完了状態を観測できる計画になっている。
  - `Draft -> Ready -> Running -> Completed` の正常系が探索証跡で追える粒度に分割されている。
  - テストデータ成果物が単一正常入力を具体化できる。
- `not_reproducible`:
  - 正常入力の定義に必要な fixture または sample input を repo 内から特定できない場合。
  - `翻訳管理` または `出力管理` の正常導線が現行成果物から一意に読めず、探索証跡が再解釈を要する場合。
- `environment_blocker`:
  - アプリ起動前提、入力投入前提、phase 実行前提のいずれかが欠け、通常フローの着手条件を満たせない場合。
- `needs_human_decision`:
  - 通常フローに含める入力種別を、人間判断なしに 1 種へ固定できない場合。
  - `出力管理` で確認すべき成果物の最小範囲を、人間判断なしに 1 種へ固定できない場合。

## Output

- `decision`: complete
- `evidence_refs`:
  - `./plan.md`
  - `docs/spec.md`
  - `docs/scenario-tests/dashboard-and-app-shell.md`
  - `frontend/src/ui/stores/shell-state.ts`
  - `tmp/code-map/index.json`
- `missing_info`:
  - `翻訳管理` の正常系シナリオ正本が `docs/scenario-tests/` に未整備である。
  - `出力管理` の正常系シナリオ正本が `docs/scenario-tests/` に未整備である。
  - 通常フローで使う推奨 fixture 名または sample input path が未指定である。
- `next_artifact`: `./exploration-test-data.md`
