# 人間UIレビュー

## 状態

- `artifact`: `人間UIレビュー`
- `status`: `approved`
- `review_url`: `http://localhost:34115/?refresh=20260505#translation-management`
- `review_command`: `agent-browser open 'http://localhost:34115/?refresh=20260505#translation-management'`

## 確認済み証跡

- `agent-browser snapshot` で `翻訳入力データロード` を確認した。
- `agent-browser snapshot` で `データロード / Job Setup / Job Run` を確認した。
- `agent-browser snapshot` で `ロード準備`、`読み込み済みデータ`、`選択データの内容` を確認した。
- `rg -n "Input Review" frontend/src/ui/screens/translation-input frontend/src/ui/stores/shell-state.ts frontend/src/ui/views/AppShell.svelte` は該当なしである。

## レビュー観点

- 翻訳管理タブと画面見出しが、`Input Review` ではなくデータロードの概念で読める。
- 上部の状態帯で、接続状態、作業状態、次操作がすぐ判断できる。
- JSON 選択、登録、選び直しが主作業としてまとまって見える。
- 一覧と詳細が、読み込み済みデータの確認画面として自然に読める。
- キャッシュ再構築が主操作ではなく、選択済みデータの補助操作として分かる。

## 停止中の検証

- `python3 scripts/harness/run.py --suite frontend-local` は失敗した。
- 失敗原因は旧テストの `Input Review` 見出し期待である。
- 人間指示により、旧表示名へ依存する frontend 単体テスト調整を次成果物として扱う。
- 後続の `テスト修正証跡` で `python3 scripts/harness/run.py --suite frontend-local` は pass へ変わった。
