# 実装計画テンプレート

- workflow: impl
- status: in_progress
- lane_owner: orchestrating-implementation
- scope: dashboard-and-app-shell
- task_id: dashboard-and-app-shell
- task_catalog_ref: /Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/dashboard-and-app-shell.yaml
- parent_phase: implementation-lane

## 要求要約

- ダッシュボードとアプリシェルを追加し、アプリ全体の入口、グローバルナビゲーション、主要ページへの導線を成立させる。

## 判断根拠

<!-- Decision Basis -->

- `/Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/dashboard-and-app-shell.yaml` に completion criteria と manual check steps が定義されている。
- completion criteria には、ダッシュボード表示、主要ページ移動、マスター辞書、マスターペルソナ、翻訳準備、実行管理、出力管理への入口確認が含まれる。
- 実装レーン契約により、前段 HITL と後段 HITL の承認前に実装へ進めない。
- active plan には task-local artifact の path と要点だけを残し、本文は別 artifact に分離する。

## 対象範囲

- `tasks/usecases/dashboard-and-app-shell.yaml`
- `docs/exec-plans/active/2026-04-11-dashboard-and-app-shell.md`
- `docs/exec-plans/active/dashboard-and-app-shell.ui.html`
- `docs/exec-plans/active/dashboard-and-app-shell.scenario.md`
- dashboard / app shell 関連の frontend / backend 実装一式

## 対象外

- `docs/` 正本の恒久仕様変更
- 個別ページの詳細機能要件を超える追加業務要件

## 依存関係・ブロッカー

- 前段 HITL と後段 HITL の承認前は実装へ進めない。
- 主要ページの現状構成と導線の既存実装有無は phase-1 で確認が必要である。

## 並行安全メモ

- 詳細設計確定前は plan 本文と task-local artifact 以外へ変更を広げない。
- frontend / backend 実装の並列化は `実装計画` section で task group 固定後に判断する。

## 機能要件

- `summary`:
  - アプリ起動直後の既定表示はダッシュボードとし、初回到達で主要機能への導線を一覧できる状態を成立条件に含める。
  - `AppShell` はアプリ全体の共通枠としてグローバルナビゲーションと現在地表示を担い、ダッシュボードは主要ページへの独立入口を提示する入口画面に責務を限定する。
  - 主要ページの表示名と識別子は current mock のままで固定し、`翻訳準備` と `実行管理` は `翻訳管理` として統合して扱う。
  - 主要ページは未実装であるため、本タスクでは導線成立を優先し、未実装ページは共通シェル配下のプレースホルダー画面で受けられることを許容する。
- `in_scope`:
  - アプリ起動時に `AppShell` 配下の既定ページとしてダッシュボードを表示する。
  - グローバルナビゲーションから、ダッシュボード、マスター辞書、マスターペルソナ、翻訳管理、出力管理の各ページへ切り替えられる。
  - ダッシュボード上に、マスター辞書、マスターペルソナ、翻訳管理、出力管理への独立した入口を表示する。
  - ダッシュボードの各入口は、グローバルナビゲーションと独立に識別できるラベルを持ち、同じ遷移先へ到達できる。
  - 未実装の遷移先は、ページ名と未実装である旨を判別できるプレースホルダー画面で受け、導線切れを起こさない。
  - 画面切替後も共通シェルは保持され、主要ページ間の再移動を同一導線で継続できる。
- `out_of_scope`:
  - マスター辞書、マスターペルソナ、翻訳管理、出力管理それぞれの業務機能詳細や内部操作の確定。
  - 主要ページのデータ取得、Wails bind 追加、backend usecase 追加などの実データ連携要件の確定。
  - ジョブ一覧や進捗サマリなど、翻訳管理に統合する詳細情報カードの確定。
  - グローバルナビゲーションのレイアウト詳細、装飾、レスポンシブ表現などの UI 表現詳細の確定。
- `open_questions`:
  - なし
- `required_reading`:
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/dashboard-and-app-shell.yaml`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-11-dashboard-and-app-shell.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/views/AppShell.svelte`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/App.svelte`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/screen-design/code.html`

## UI モック

- `artifact_path`: `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/mocks/dashboard-and-app-shell/index.html`
- `final_artifact_path`: `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/mocks/dashboard-and-app-shell/index.html`
- `summary`:
  - アプリ起動直後の既定表示を `dashboard` に固定した単一 HTML モックへ更新し、hash 遷移で共通シェル配下の主要ページ切替を再現した。
  - 上端グローバルナビゲーションとダッシュボード内の入口カードを `ダッシュボード`、`マスター辞書`、`マスターペルソナ`、`翻訳管理`、`出力管理` の 5 遷移先へ揃え、現在地表示を共通ヘッダへ残した。
  - HTML 本文から reviewer 向け設計注釈を除去し、ダッシュボード以外の主要ページは `ページ名`、`一文の状態説明`、`移動導線` だけを表示する最小プレースホルダーへ整理した。

## Scenario テスト一覧

- `artifact_path`: `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/scenario-tests/dashboard-and-app-shell.md`
- `final_artifact_path`: `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/scenario-tests/dashboard-and-app-shell.md`
- `template_path`: `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/templates/scenario-tests.md`
- `summary`:
  - 起動時のダッシュボード既定表示、グローバルナビゲーション 5 遷移、ダッシュボード入口カード 5 遷移を UI 観測可能なケースとして固定した。
  - 非ダッシュボード遷移先はプレースホルダーで受け、共通シェル維持と継続遷移可能性を主要例外系と責務境界ケースで固定した。
  - ダッシュボード本文にジョブ一覧と進捗サマリの表示領域が存在しないことを `SCN-DAS-001` の手順と期待結果で直接観測するよう更新した。

## 実装計画

<!-- Implementation Plan -->

- `ordered_scope`:
  - `1.` `frontend/src/ui/stores/shell-state.ts` に、既定ページ `dashboard`、主要ページ識別子、表示名、現在地表示、未実装ページ向け metadata をまとめた shell route contract を固定する。
  - `2.` `frontend/src/ui/views/AppShell.svelte` に、共通ナビゲーション、現在地表示、ダッシュボード入口群、未実装ページ共通プレースホルダーの描画責務を集約し、route contract だけを読む薄い view へ置き換える。
  - `3.` `frontend/src/ui/App.svelte` に、shell state 初期化と `AppShell` への受け渡しだけを残し、起動直後に `dashboard` が表示される bootstrap wiring を固定する。
- `parallel_task_groups`:
  - `group_id`: `shell-route-contract`
  - `can_run_in_parallel_with`: `none`
  - `blocked_by`: `none`
  - `completion_signal`: 5 主要ページの識別子、表示名、既定ページ、未実装ページ用 metadata が store 側で一元化され、`AppShell` が参照できる
  - `group_id`: `shell-view-assembly`
  - `can_run_in_parallel_with`: `app-bootstrap`
  - `blocked_by`: `shell-route-contract`
  - `completion_signal`: 共通ナビゲーション、現在地表示、ダッシュボード入口、プレースホルダー表示が route contract から描画される
  - `group_id`: `app-bootstrap`
  - `can_run_in_parallel_with`: `shell-view-assembly`
  - `blocked_by`: `shell-route-contract`
  - `completion_signal`: `App.svelte` が shell state を初期化し、起動直後に `dashboard` を表示する最小 wiring へ整理される
- `task_dependencies`:
  - `task_id`: `define-shell-route-contract`
  - `depends_on`: `none`
  - `enables`: `assemble-shell-view`, `wire-app-bootstrap`
  - `reason`: 主要ページの識別子と既定表示が未固定のままでは view と bootstrap の props / state 契約を決められない
  - `task_id`: `assemble-shell-view`
  - `depends_on`: `define-shell-route-contract`
  - `enables`: `wire-app-bootstrap`
  - `reason`: `AppShell` の受け取り形と表示責務を先に固定しないと root wiring の最小境界が定まらない
  - `task_id`: `wire-app-bootstrap`
  - `depends_on`: `define-shell-route-contract`, `assemble-shell-view`
  - `enables`: `manual-validation`
  - `reason`: 起動時の既定表示と画面切替導線は route contract と shell view の両方が揃って初めて確認できる
- `tasks`:
  - `task_id`: `define-shell-route-contract`
  - `owned_scope`: `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/stores/shell-state.ts`
  - `depends_on`: `none`
  - `parallel_group`: `shell-route-contract`
  - `required_reading`: `/Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/dashboard-and-app-shell.yaml`, `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-11-dashboard-and-app-shell.md`, `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/dashboard-and-app-shell.ui.html`, `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md`
  - `validation_commands`: `python3 scripts/harness/run.py --suite structure`, `npm --prefix frontend run check`
  - `task_id`: `assemble-shell-view`
  - `owned_scope`: `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/views/AppShell.svelte`
  - `depends_on`: `define-shell-route-contract`
  - `parallel_group`: `shell-view-assembly`
  - `required_reading`: `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-11-dashboard-and-app-shell.md`, `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/dashboard-and-app-shell.ui.html`, `/Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/dashboard-and-app-shell.yaml`, `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md`
  - `validation_commands`: `npm --prefix frontend run check`, `npm --prefix frontend run lint`
  - `task_id`: `wire-app-bootstrap`
  - `owned_scope`: `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/App.svelte`
  - `depends_on`: `define-shell-route-contract`, `assemble-shell-view`
  - `parallel_group`: `app-bootstrap`
  - `required_reading`: `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/App.svelte`, `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/views/AppShell.svelte`, `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/stores/shell-state.ts`, `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md`
  - `validation_commands`: `npm --prefix frontend run check`, `npm --prefix frontend run build`
- `owned_scope`:
  - frontend-only。`App.svelte`、`AppShell.svelte`、`shell-state.ts` の shell bootstrap と view 切替に限定し、backend bind / usecase / data integration には触れない。
- `required_reading`:
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-11-dashboard-and-app-shell.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/dashboard-and-app-shell.ui.html`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/tasks/usecases/dashboard-and-app-shell.yaml`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/App.svelte`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/views/AppShell.svelte`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/stores/shell-state.ts`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md`
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite structure`
  - `npm --prefix frontend run check`
  - `npm --prefix frontend run lint`
  - `npm --prefix frontend run build`

## 受け入れ確認

- アプリ起動時にダッシュボードを表示できる。
- グローバルナビゲーションから主要ページへ移動できる。
- ダッシュボードから主要機能の入口を辿れる。
- マスター辞書ページへの独立導線を確認できる。
- マスターペルソナページへの独立導線を確認できる。
- 翻訳管理ページへの独立導線を確認できる。
- 出力管理ページへの独立導線を確認できる。

## 必要な証跡

<!-- Required Evidence -->

- phase-1 以降の artifact path と要約
- 前段 HITL と後段 HITL の承認記録
- 実装レビュー結果
- `python3 scripts/harness/run.py --suite all` の最終結果

## 機能要件 HITL 状態

- approved

## 機能要件 承認記録

- 2026-04-11 human review: 表示名と識別子は current mock のままで承認。`翻訳準備` と `実行管理` は `翻訳管理` へ統合。ダッシュボードへジョブ一覧と進捗サマリは置かない。

## 詳細設計 HITL 状態

- approved

## 詳細設計 承認記録

- 2026-04-11 phase-2.5-design-review: `pass`。task-local design artifact 間の整合確認を通過。human の後段 HITL 承認待ち。
- 2026-04-11 human review: phase3 として後段 HITL を完了。第4段階は workflow 上で `phase-2-logic` へ吸収済みのため欠番とし、第5段階 `phase-5-test-implementation` から再開する。

## review 用差分図

- closeout で `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/diagrams/components/frontend/dashboard-and-app-shell.d2` と `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/diagrams/components/frontend/dashboard-and-app-shell.svg` へ正本適用済み。review copy は active plan 配下から退避済み。

## 差分正本適用先

- `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/diagrams/components/frontend/dashboard-and-app-shell.d2`
- `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/diagrams/components/frontend/dashboard-and-app-shell.svg`

## Closeout Notes

- review 用に active exec-plan 配下へ置いた差分 D2 / SVG copy は、`diagrams/backend/` または `diagrams/frontend/` 正本適用後に削除し、completed plan へ持ち越さない。
- 第1.6段階で作った UI モック working copy は、完了前に `docs/mocks/dashboard-and-app-shell/index.html` へ移す。
- 第2段階で作った Scenario artifact working copy は、完了前に `docs/scenario-tests/dashboard-and-app-shell.md` へ移す。

## 結果

<!-- Outcome -->

- in_progress
- 2026-04-12 human review: 承認済みモック `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/mocks/dashboard-and-app-shell/index.html` と実装 `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/views/AppShell.svelte` の視覚構造差分を理由に、第6段階 frontend 実装へ差し戻し。
- 差し戻し論点: header 3 カラム構造、hero 情報ブロック、dashboard entry-card 構造、placeholder 構造、mobile nav affordance、route metadata がモック準拠でない。
- 次回再開位置: `phase-6-implement-frontend`
