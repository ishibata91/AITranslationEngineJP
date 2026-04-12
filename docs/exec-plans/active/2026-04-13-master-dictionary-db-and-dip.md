# 実装計画

- workflow: impl
- status: planned
- lane_owner: orchestrating-implementation
- scope: master-dictionary-db-and-dip
- task_id: master-dictionary-db-and-dip
- task_catalog_ref: N/A
- parent_phase: implementation-lane

## 要求要約

- マスター辞書をメモリ保持ではなく DB 永続化へ移行する。
- マスター辞書まわりの上位層 test を、SQLite concrete 実装へ直結しない構成へ改める。
- controller / service / usecase / repository の境界に沿って、レイヤー間の DIP を入れ、unit test を mock / fake ベースで組める状態へ寄せる。

## 判断根拠

<!-- Decision Basis -->

- 直前の SQLite 化試行は破棄した。
- 破棄理由は、DB 化そのものよりも、default concrete 実装へ直結した test 構造が arch lint、driver 登録、coverage 要求と衝突したためである。
- 今回は `辞書DB化` と `アーキテクチャレイヤー間 DIP 化` を分離せず、同一 task として扱う。
- 高 coverage を unit test 中心で維持するには、上位層が concrete repository や driver 登録へ引きずられない構成が前提になる。

## 対象範囲

- `internal/controller/`
- `internal/usecase/`
- `internal/service/`
- `internal/repository/`
- master dictionary の DB 永続化導線
- master dictionary 関連 unit / integration test の責務整理

## 対象外

- master dictionary 以外の persistence 移行
- docs 正本の恒久仕様変更
- arch lint ルール自体の緩和

## 依存関係・ブロッカー

- `docs/tech-selection.md` の SQLite 方針と整合すること。
- `docs/architecture.md` の依存方向を壊さないこと。
- DB 化より先に、どのレイヤーまでを mock / fake 化するかの注入点整理が必要。

## 並行安全メモ

- まず依存注入点と default wiring を整理しないと、DB 化 task と test 再設計 task が相互に競合する。
- `controller` / `service` の unit test は concrete SQLite 実装から切り離す。
- SQLite 実行証明は repository integration test または top-level wiring test へ寄せる。

## 機能要件

- `summary`:
  - master dictionary を SQLite へ永続化する。
  - controller / service / usecase の unit test が DB や driver 登録なしで成立するようにする。
  - default wiring と SQLite 実行証明は責務を限定した integration test へ分離する。
- `in_scope`:
  - master dictionary の DB schema と migration。
  - repository 抽象と default factory の見直し。
  - 上位層への DIP 導入。
  - test 戦略の再配置。
- `non_functional_requirements`:
  - arch lint を例外追加なしで通す。
  - coverage 目標を unit / integration の責務分離を保ったまま満たせる構造にする。
  - driver 依存は composition root か専用 integration 境界へ閉じ込める。
- `out_of_scope`:
  - `_test.go` 全除外などの lint 緩和。
  - controller / service unit test からの DB 直接実行継続。
- `open_questions`:
  - repository の抽象境界を service 直下までに留めるか、usecase 入力境界まで広げるか。
  - SQLite 実行証明を `repository` integration に置くか top-level package に置くか。
- `required_reading`:
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/tech-selection.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/master-dictionary.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/scenario-tests/master-dictionary-management.md`

## UI モック

- `artifact_path`: `docs/mocks/master-dictionary/index.html`
- `final_artifact_path`: `docs/mocks/master-dictionary/index.html`
- `summary`:
  - UI 追加要件はない。既存 mock を参照する。

## Scenario テスト一覧

- `artifact_path`: `docs/scenario-tests/master-dictionary-management.md`
- `final_artifact_path`: `docs/scenario-tests/master-dictionary-management.md`
- `template_path`: `docs/exec-plans/templates/scenario-tests.md`
- `summary`:
  - 既存 scenario を参照し、削除後の page clamp、未インポート空状態、再起動後保持、起動時 migration を再証明対象とする。

## 実装計画

<!-- Implementation Plan -->

- `parallel_task_groups`:
  - `group_id`: `md-dip-design`
  - `can_run_in_parallel_with`: なし
  - `blocked_by`: なし
  - `completion_signal`: 上位層 unit test が concrete SQLite 実装へ依存しない注入点が確定している。
  - `group_id`: `md-sqlite-implementation`
  - `can_run_in_parallel_with`: なし
  - `blocked_by`: `md-dip-design`
  - `completion_signal`: DB 永続化と test 責務分離が同時に成立している。
- `tasks`:
  - `task_id`: `master-dictionary-dbization`
  - `owned_scope`: master dictionary SQLite persistence, migration, startup bootstrap, repository integration proof
  - `depends_on`: `master-dictionary-layer-dip`
  - `parallel_group`: `md-sqlite-implementation`
  - `required_reading`:
    - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/tech-selection.md`
    - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md`
  - `validation_commands`:
    - `python3 scripts/harness/run.py --suite backend-lint`
    - `python3 scripts/harness/run.py --suite execution`
    - `python3 scripts/harness/run.py --suite coverage`
  - `task_id`: `master-dictionary-layer-dip`
  - `owned_scope`: controller/service/usecase/repository injection boundaries, unit test seam design, default wiring separation
  - `depends_on`: なし
  - `parallel_group`: `md-dip-design`
  - `required_reading`:
    - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md`
    - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/master-dictionary.md`
  - `validation_commands`:
    - `python3 scripts/harness/run.py --suite backend-lint`
    - `python3 scripts/harness/run.py --suite coverage`

## 受け入れ確認

- controller / service の unit test が mock / fake で閉じる。
- SQLite 永続化は repository integration または専用 wiring test で証明される。
- driver 依存が `internal` の不要な層へ漏れない。
- arch lint / execution / coverage を満たす。

## 必要な証跡

<!-- Required Evidence -->

- レイヤー間 DIP 導入後の unit test 証跡
- SQLite 永続化と migration の integration 証跡
- arch lint / execution / coverage の結果

## 機能要件 HITL 状態

- pending

## 機能要件 承認記録

- 2026-04-13 human: 直前の SQLite 化試行はすべて discard し、follow-up として `辞書DB化` と `アーキテクチャレイヤー間 DIP 化` を task に積むよう指示。

## 詳細設計 HITL 状態

- pending

## 詳細設計 承認記録

- N/A

## review 用差分図

- N/A

## 差分正本適用先

- N/A

## Closeout Notes

- この plan は再着手用の placeholder であり、実装は未開始。

## 結果

<!-- Outcome -->

- planned
