# UX実装修正入力

## 対象成果物

- `frontend 実装`
- 使用 skill: `implement-frontend`

## 承認済み UI改善契約

- [ui-design.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/ux-dashboard-refactor-20260504/ui-design.md)
- 人間UIレビュー: 承認
- 承認条件: `準備中` の文言変更を追加する。

## 実装対象

- [frontend/src/ui/stores/shell-state.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/stores/shell-state.ts)

## 対象変更範囲

- `SHELL_ROUTE_CONTRACT` の入口カード状態値だけを変更する。
- `master-dictionary` の `state` を `準備中` から `確認可能` へ変更する。
- `master-persona` の `state` を `準備中` から `確認可能` へ変更する。
- `output-management` の `state` を `準備中` から `確認可能` へ変更する。

## 禁止変更範囲

- route id、route label、lead、description を変更しない。
- `dashboardEntryRoutes` の抽出条件を変更しない。
- `AppShell.svelte` の構造整理を行わない。
- backend、統合境界、Wails bind、DTO、docs 正本本文を変更しない。
- プロダクトテストを変更しない。

## 確認観点

- ダッシュボード入口カードから `準備中` が消える。
- `マスター辞書`、`マスターペルソナ`、`出力管理` の状態値に `確認可能` が出る。
- `AIサービス設定` と `翻訳管理` の状態値は維持する。
- 入口カード 5 件とグローバルナビゲーション 6 件は欠落しない。
- `ジョブ一覧` と `進捗サマリ` はダッシュボード本文に出ない。

## 検証コマンド

- `npm --prefix frontend run test -- AppShell`
- `python3 scripts/harness/run.py --suite frontend-local`

## 書き戻し先

- 実装結果: `docs/exec-plans/active/ux-dashboard-refactor-20260504/implementation-result.frontend.md`

