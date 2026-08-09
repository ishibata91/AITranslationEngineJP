# Go backend の境界

- production の concrete dependency は `bootstrap/` だけで生成する。
- Wails runtime は `api/` と `bootstrap/` の外へ渡さない。
- SQLite driver は product code では `store/` に閉じ込める。
- 多態の interface は AI provider の境界に置く。
- 単一実装の interface は、利用側のテスト境界として必要な method だけを宣言する場合に限る。
- backend は状態、種別、条件を返す。
- 画面の操作可否を表す boolean を返さない。
- error は失敗した処理と対象を加えて上位へ返す。

# 観測

- backend の観測ログは `log/slog` で出す。
- 状態遷移、外部境界の失敗、大量処理の集約結果だけを記録する。
- `event`、`where`、`result` を基本属性とし、原因分離に必要な場合だけ `id`、`count`、`reason` を加える。
- loop の要素ごとのログ、secret、API key、prompt 全文、翻訳本文全文、外部 payload 全体を出さない。

# 検証

- test は公開された振る舞いか責務境界から観測できる結果を検証する。
- clock、random、外部応答、並び順を固定し、filesystem と database は一時資源へ閉じる。
- backend の変更後は `npm run verify:backend` を実行する。
