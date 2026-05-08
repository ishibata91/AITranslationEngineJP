# Scenario Design: observability-logging-foundation

## 目的

frontend と backend の観測ログ基盤を固定する。
対象は恒久ログであり、一時 trace の追加と除去ではない。

## 必須要件

- backend は `log/slog` を使い、構造化ログを出せる。
- backend の診断ログは `tmp/observability-logging-foundation/diagnostic-log.sqlite` に保存できる。
- frontend は `pino` を使い、repo 側の診断 logger 境界から browser console へ構造化ログを出せる。
- trace ID は全関数引数へ広げず、処理入口、logger scope、context で扱う。
- secret、API key、provider raw payload、過剰本文はログへ出さない。
- ループや大量処理では同種ログを増やさず、集約情報を優先する。

## 受け入れ条件

- frontend の恒久ログは `console.log` 直書きではなく、診断 logger 境界を通る。
- frontend の診断ログは agent-browser の console 証跡で取得できる。
- backend の恒久ログは `slog` の key-value 形式で検索できる。
- backend の診断ログは `sqlite3` CLI で検索できる。
- frontend と backend のログは、同じ操作を trace ID で追える。
- 観測ログ追加 agent が、実行時にしか確定しない値と実行後に消える分岐理由を残す判断をできる。
- 禁止情報が DTO、UI、structured log、debug log、runtime event へ露出しない。

## 非対象

- 外部監視 SaaS 連携は扱わない。
- OpenTelemetry 導入は扱わない。
- runtime event をログ正本へ置き換える変更は扱わない。
- frontend log を Wails command で逐次送信する変更は扱わない。
- frontend log を SQL へ永続保存する変更は扱わない。
- 既存業務処理の振る舞い変更は扱わない。

## 検証観点

- frontend lint が通る。
- backend test が通る。
- logger 境界の単体テストで、trace ID と redaction の扱いを確認できる。
- agent-browser で frontend console 証跡を取得できる。
- `sqlite3 tmp/observability-logging-foundation/diagnostic-log.sqlite` で backend 診断ログを検索できる。
- 大量処理のログが件数、分類、代表 ID、最初の失敗、最後の失敗へ集約されている。
