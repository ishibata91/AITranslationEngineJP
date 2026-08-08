# 観測ログ仕様

関連文書: [`index.md`](./index.md), [`tech-selection.md`](./tech-selection.md), [`coding-guidelines-backend.md`](./coding-guidelines-backend.md), [`coding-guidelines-frontend.md`](./coding-guidelines-frontend.md)

本書は、backend と frontend の観測ログ構造を定義する。
観測ログは、実行後に消える状態、分岐理由、外部境界の失敗分類を後続調査で分離するために使う。

## 1. 出力先

- backend は `log/slog` の JSON log を `stderr` へ出す。
- `dev:wails:agent-browser` 起動時は、既存 script が backend の stdout / stderr を `tmp/logs/wails-dev.log` へ起動毎に削除してから保存する。
- frontend は `pino` の browser console 出力を使う。
- frontend log は Wails 経由で backend へ送らない。
- backend log と frontend log を同じ file へ集約しない。

## 2. 共通 payload

全ログで意識する共通 payload は次の 3 項目だけにする。

- `event`: 何が起きたか。
- `where`: どの境界か。
- `result`: 結果分類。

logger が自動で持つ timestamp、level、message、source は共通 payload に含めない。

## 3. 任意 payload

必要なログだけに次を追加する。

- `id`: 原因分離に必要な代表 ID。
- `count`: 大量処理の集約数。
- `reason`: 拒否、破棄、失敗分類。

個別 ID を共通必須項目にしない。
`jobId`、`phaseRunId`、`inputSourceId` は必要な場所だけで使う。
複数 ID を入れる場合でも、ログ追加箇所ごとに原因分離へ必要な最小数にする。

## 4. 追加対象

- 状態遷移の変更前、変更後、拒否理由が同じ場所で取れる箇所。
- provider、file、DB、Wails 境界の失敗分類が同じ場所で取れる箇所。
- 大量処理の件数、分類、最初の失敗、最後の失敗が集約済みの箇所。
- frontend runtime event の破棄理由が画面操作後に消える箇所。

## 5. 禁止事項

- trace ID を追加しない。
- 全 command の start / finish log を追加しない。
- loop 内で 1 件ごとのログを出さない。
- DTO 全体、secret、API key、provider raw payload、prompt 全文、翻訳本文全文、XML 全文を出さない。
- logger のために constructor 引数を広げない。
- context へ logger を埋め込まない。
- frontend から backend へログを送らない。

## 6. 例

backend:

```go
slog.InfoContext(ctx, "phase rejected",
	slog.String("event", "phase_rejected"),
	slog.String("where", "backend.service"),
	slog.String("result", "rejected"),
	slog.String("id", "job:12"),
	slog.String("reason", "invalid_state"),
)
```

frontend:

```ts
logger.warn("runtime event dropped", {
  event: "runtime_event_dropped",
  where: "frontend.runtime",
  result: "skipped",
  id: "job:12",
  reason: "stale_event"
})
```
