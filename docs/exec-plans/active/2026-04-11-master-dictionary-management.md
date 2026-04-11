# 実装計画テンプレート

- workflow: impl
- status: planned
- lane_owner: orchestrating-implementation
- scope: master-dictionary-management
- task_id: master-dictionary-management
- task_catalog_ref: /Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/master-dictionary-management.yaml
- parent_phase: implementation-lane

## 要求要約

- マスター辞書ページで、辞書の取り込み、参照、作成、更新、削除導線を成立させる。

## 判断根拠

<!-- Decision Basis -->

- `/Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/master-dictionary-management.yaml` に completion criteria と manual check steps が定義されている。
- `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md` では、マスター辞書が基盤データとして UI から観測可能であることを要求している。
- user review-back により、マスターペルソナは `extractData.pas` 由来 JSON 前提で別 task へ切り分け、今回はマスター辞書に集中する。
- user review-back により、マスター辞書はダッシュボード配下ではなく独立ページとして扱う。

## 対象範囲

- `tasks/usecases/master-dictionary-management.yaml`
- `docs/exec-plans/active/2026-04-11-master-dictionary-management.md`
- `docs/exec-plans/active/master-dictionary-management.ui.html`
- `docs/exec-plans/active/master-dictionary-management.scenario.md`
- master dictionary 関連の frontend / backend 実装一式

## 対象外

- マスター辞書以外の画面詳細実装
- `docs/` 正本の恒久仕様変更

## 依存関係・ブロッカー

- 前段 HITL と後段 HITL の承認前は実装へ進めない。
- 個別 screen-design 正本が未整備のため、task-local UI mock で不足分を補う必要がある。

## 並行安全メモ

- 詳細設計確定前は plan 本文と task-local artifact 以外へ変更を広げない。
- frontend / backend 実装の並列化は `実装計画` section で task group 固定後に判断する。

## 機能要件

- `summary`:
  - 今回の task はマスター辞書だけを対象にし、マスターペルソナは別 task へ切り分ける。
  - マスター辞書ページは独立ページとして扱い、一覧参照、エントリ検索、選択中エントリの詳細確認、作成、更新、削除、`XMLから取り込み` を同一 task の機能要件に含める。
  - 一覧サマリ表示、xTranslator 写像表示、実装方針が見える UI 表現は要求しない。
- `in_scope`:
  - マスター辞書ページへ独立して到達できること。
  - マスター辞書一覧を参照できること。
  - マスター辞書一覧でエントリを検索できること。
  - 一覧で選択した辞書エントリの詳細情報を参照できること。
  - マスター辞書の新規作成導線を提供すること。
  - マスター辞書の既存エントリ更新導線を提供すること。
  - マスター辞書の既存エントリ削除導線を提供すること。
  - `XMLから取り込み` 導線を提供すること。
  - UI 上の文言、ラベル、説明には実装方針が見える表現を持ち込まないこと。
- `out_of_scope`:
  - マスターペルソナ画面および `extractData.pas` 由来 JSON を扱う導線。
  - 基盤データ以外の翻訳ジョブ画面、設定画面、翻訳成果物画面の詳細導線追加。
  - `docs/` 正本の恒久仕様変更、および usecase 完了条件を超える新規業務要件の追加。
  - 一覧サマリ表示、xTranslator 写像表示、利用状況表示、要件説明の露出。
  - UI モック、Scenario テスト一覧、実装計画、review 用差分図の確定。
- `open_questions`:
  - `XMLから取り込み` の対象ファイル選択をページ内で扱うか、別 dialog として扱うか。
  - 一覧、詳細、作成、更新を同一ページ内で扱うか、create / edit を別画面や modal に出すか。
  - 辞書 detail に表示する項目を、原文、訳語、由来、最終更新のどこまでに絞るか。
- `required_reading`:
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-11-master-dictionary-management.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/master-dictionary-management.yaml`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/screen-design/code.html`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/README.md`

## UI モック

- `artifact_path`: `docs/exec-plans/active/master-dictionary-management.ui.html`
- `final_artifact_path`: `docs/mocks/master-dictionary/index.html`
- `summary`:
  - マスター辞書は独立ページとして扱い、ダッシュボード配下の子ページの見え方は持ち込まない。
  - 一覧と詳細を横並びで確認できる構成にし、左ペインは一覧 summary を置かず、エントリ検索から目的の項目へ直接辿れるようにする。
  - 上段の細い操作バーに `XMLから取り込み` / `新規作成` / `更新` / `削除` をまとめる。
  - detail には基本情報だけを置き、メモ欄、利用状況、要件説明、xTranslator 写像のような補助情報は出さない。

## Scenario テスト一覧

- `artifact_path`: `docs/exec-plans/active/master-dictionary-management.scenario.md`
- `final_artifact_path`: `docs/scenario-tests/master-dictionary-management.md`
- `template_path`: `docs/exec-plans/templates/scenario-tests.md`
- `summary`:

## 実装計画

<!-- Implementation Plan -->

- `parallel_task_groups`:
  - `group_id`:
  - `can_run_in_parallel_with`:
  - `blocked_by`:
  - `completion_signal`:
- `tasks`:
  - `task_id`:
  - `owned_scope`:
  - `depends_on`:
  - `parallel_group`:
  - `required_reading`:
  - `validation_commands`:

## 受け入れ確認

- マスター辞書一覧を参照できる。
- マスター辞書を作成できる。
- 選択中の基盤エントリの詳細を確認できる。
- 選択中の基盤エントリを更新できる。
- 選択中の基盤エントリを削除できる。
- `XMLから取り込み` 導線を確認できる。

## 必要な証跡

<!-- Required Evidence -->

- phase-1 以降の artifact path と要約
- 前段 HITL と後段 HITL の承認記録
- 実装レビュー結果
- `python3 scripts/harness/run.py --suite all` の最終結果

## 機能要件 HITL 状態

- review_ready

## 機能要件 承認記録

- 2026-04-11 orchestrating-implementation: 機能要件と UI モックを前段 HITL へ回付。承認待ち。

## 詳細設計 HITL 状態

- pending

## 詳細設計 承認記録

- pending

## review 用差分図

- pending

## 差分正本適用先

- pending

## Closeout Notes

- review 用に active exec-plan 配下へ置いた差分 D2 / SVG copy は、`diagrams/backend/` または `diagrams/frontend/` 正本適用後に削除し、completed plan へ持ち越さない。
- 第1.6段階で作った UI モック working copy は、完了前に `docs/mocks/master-dictionary/index.html` へ移す。
- 第2段階で作った Scenario artifact working copy は、完了前に `docs/scenario-tests/master-dictionary-management.md` へ移す。

## 結果

<!-- Outcome -->

- in_progress
