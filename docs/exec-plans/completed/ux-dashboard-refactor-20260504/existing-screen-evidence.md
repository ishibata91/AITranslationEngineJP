# 既存画面根拠

## 対象

- `screen`: ダッシュボード
- `url`: `http://127.0.0.1:34115/#dashboard`
- `確認日時`: 2026-05-04
- `確認者`: `ux_refactor_lane`

## 根拠参照

- [docs/scenario-tests/dashboard-and-app-shell.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/scenario-tests/dashboard-and-app-shell.md)
- [docs/exec-plans/completed/2026-04-11-dashboard-and-app-shell.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-11-dashboard-and-app-shell.md)
- [frontend/src/ui/views/AppShell.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/views/AppShell.svelte)
- [frontend/src/ui/stores/shell-state.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/stores/shell-state.ts)

## 実物確認結果

- `agent-browser open http://127.0.0.1:34115/#dashboard`: pass
- `agent-browser snapshot`: pass
- `agent-browser screenshot tmp/agent-browser/ux-dashboard-refactor-before.png`: pass
- `agent-browser errors`: 出力は空行のエラー印のみで、具体的な console error text は取得できなかった。

## 維持対象

- 共通ヘッダに `AITranslationEngineJp` と `翻訳エンジン` を表示する。
- グローバルナビゲーションに `ダッシュボード`、`AIサービス設定`、`マスター辞書`、`マスターペルソナ`、`翻訳管理`、`出力管理` を表示する。
- ページ見出しとして `現在のページ` と `ダッシュボード` を表示する。
- リード文として `最初に移動したい作業を選び、共通ナビゲーションからいつでも別の主要ページへ切り替えられます。` を表示する。
- ダッシュボード本文に `主要ページ` と `作業を選ぶ` を表示する。
- 入口カードとして `AIサービス設定`、`マスター辞書`、`マスターペルソナ`、`翻訳管理`、`出力管理` を表示する。

## 入口カード維持対象

- `AIサービス設定`: 状態 `設定状態を確認`、説明 `エンドポイントと APIキー状態を AIサービスごとに確認します。`、操作 `開く`
- `マスター辞書`: 状態 `準備中` を `確認可能` へ変更する。説明 `用語と訳語の基盤データを確認します。`、操作 `開く`
- `マスターペルソナ`: 状態 `準備中` を `確認可能` へ変更する。説明 `翻訳に使うペルソナ設定を確認します。`、操作 `開く`
- `翻訳管理`: 状態 `Body Phase UI 追加`、説明 `入力確認、validation、ready job 作成、term phase、persona phase、body phase の実行状況をまとめて確認します。`、操作 `開く`
- `出力管理`: 状態 `準備中` を `確認可能` へ変更する。説明 `生成物と書き出し結果を確認します。`、操作 `開く`

## 禁止表示

- ダッシュボード本文に `ジョブ一覧` を表示しない。
- ダッシュボード本文に `進捗サマリ` を表示しない。
- ダッシュボード入口カードに `ダッシュボード` 自身を表示しない。

## リファクタ時の確認観点

- 既存の全表示項目が欠落しない。
- 既存の各入口カードは同じ route に遷移する。
- ダッシュボードの責務は主要ページ入口に限定する。
- 状態値と説明文は、既存契約から意味を広げない。
- デスクトップとモバイルで、ヘッダ、見出し、入口カード、操作 `開く` が確認できる。
