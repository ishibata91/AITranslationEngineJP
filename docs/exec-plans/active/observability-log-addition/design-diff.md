# 設計差分図説明

## 根拠参照

- `docs/exec-plans/active/observability-log-addition/scenario-design.md`
- `docs/exec-plans/active/observability-log-addition/scenario-design.requirement-coverage.json`
- `docs/exec-plans/active/observability-log-addition/scenario-design.candidate-coverage.json`
- `docs/observability-logging.md`
- `docs/architecture.md`
- `docs/exec-plans/active/observability-log-addition/plan.md`

## 図が補う判断

人間設計レビュー前に、観測ログを追加する境界を最小範囲で確認する。
同時に、既存接続を維持する箇所と、追加しない接続を分けて確認する。

対象は 4 境界に限定する。

- backend 状態遷移境界
- backend 外部境界
- frontend runtime event 境界
- 大量処理境界

## 追加予定

- backend 状態遷移境界に、`slog` の JSON log を追加する。
  追加内容: `event`、`where`、`result`、必要な `id`、`reason`、変更前状態、変更後状態。
- backend 外部境界に、`slog` の JSON log を追加する。
  追加内容: provider、secret、file、DB、Wails response 変換、成果物出力の失敗分類。
- frontend runtime event 境界に、`pino` の browser console log を追加する。
  追加内容: subscribed、accepted、dropped、skipped、detached の分類。
- 大量処理境界に、集約 log を追加する。
  追加内容: `count`、最初の失敗分類、最後の失敗分類。

## 削除予定

既存接続の削除予定はない。
差分図では、次の廃案または禁止接続を `removed` として明示した。

- frontend log を backend へ送る接続。
- backend log と frontend log を同じ file へ集約する接続。
- trace ID の一律追加。
- 全 command の start / finish log。
- DTO 全体 dump を含む重い payload。

## 変更しない接続先

- backend の JSON log 出力先は `stderr` のまま維持する。
- `dev:wails:agent-browser` 実行時の backend 保存先は `tmp/logs/wails-dev.log` のまま維持する。
- frontend log の出力先は browser console のまま維持する。
- Wails runtime event は push 通知専用のまま維持する。
- backend の command 経路は `Wails Bind Controller -> Backend UseCase / Service` を維持する。
- UI 表示、画面文言、layout、style は変更しない。

## 検証結果

- PlantUML source を 2 件作成した。
  - `docs/exec-plans/active/observability-log-addition/design-diff.component.puml`
  - `docs/exec-plans/active/observability-log-addition/design-diff.sequence.puml`
- PlantUML syntax check を実行し、2 件とも成功した。
  - `plantuml --check-syntax --no-error-image docs/exec-plans/active/observability-log-addition/design-diff.component.puml`
  - `plantuml --check-syntax --no-error-image docs/exec-plans/active/observability-log-addition/design-diff.sequence.puml`
- 描画結果は作成していない。
  理由: 本依頼の禁止範囲に「担当出力ファイル以外を変更しない」が含まれるため。

## 未決事項

- 最初に実装する slice の順序は未決である。
- 最初に対象にする frontend runtime event の画面範囲は未決である。
- 既存の `observability-logger-lightweight` を完了扱いにするか、本 task へ統合するかは未決である。
