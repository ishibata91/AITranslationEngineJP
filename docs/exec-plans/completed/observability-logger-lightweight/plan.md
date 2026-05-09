# 観測ロガー軽量化 plan

## 状態

- task-id: `observability-logger-lightweight`
- 目的: 観測ログ基盤を軽くし、既存コードへ長い logging wrapper を増やさずに原因候補を分離できるようにする。
- 決定: trace id は使わない。
- 決定: command ごとの start/finish ログは既定にしない。
- 決定: frontend は `pino` で browser console へ出す。
- 決定: backend は `slog` を入口にする。

## 理由

- trace id は生成、伝播、紐付けの実装を増やす。
- command ごとの start/finish は差分量に対して観測価値が低い。
- 既存 ID、状態値、command 名、phase run ID、結果分類で多くの原因候補は分離できる。
- 観測ログ基盤は業務コードより薄くあるべき。

## 軽量化規約

- trace id を追加しない。
- context へ logger を埋め込まない。
- constructor へ logger 引数を広げない。
- 各 method を logging wrapper で包まない。
- DTO 全体、全文入力、secret、provider raw payload を出さない。
- loop 内では出さず、件数と分類だけを出す。

## backend 方針

- `slog.Default()` または小さい package helper だけを使う。
- 保存先の切り替えは app 起動時だけで扱う。
- 業務コードは `slog.InfoContext` 1 回から 2 回で済む箇所だけ対象にする。
- SQLite handler が長くなる場合は採用しない。
- backend の保存先は最初は stderr または既存 dev server log へ寄せる。
- SQLite が必要なら、後で handler だけを差し替える。

## frontend 方針

- repo 内の薄い `pino` wrapper だけを使う。
- Wails 経由でログ保存しない。
- gateway method 全体を包まない。
- runtime event の破棄、例外、画面から消える選択 ID だけを出す。
- 画面文言、layout、style は変更しない。

## 共通 payload

- `event`: 何が起きたか。
- `where`: どの境界か。
- `result`: 結果分類。
- `id`: 必要な場合だけ代表 ID を入れる。
- `count`: 必要な場合だけ集約数を入れる。
- `reason`: 必要な場合だけ拒否、破棄、失敗分類を入れる。

## 追加対象

- 状態遷移: 変更前、変更後、拒否理由が同じ場所で取れる箇所。
- 外部境界: provider、file、DB、Wails の失敗分類が同じ場所で取れる箇所。
- 大量処理: 件数、分類、最初の失敗、最後の失敗が集約済みの箇所。
- runtime event: 受信後に破棄された理由が画面操作後に消える箇所。

## 追加しない対象

- 成功した read-only query。
- 全 command の開始終了。
- phase 内の candidate 1 件ごとの進捗。
- provider request / response の本文。
- UI 操作の全クリック履歴。

## 実装単位

### P0 backend logger entry

- 目的: `slog` の出力先だけを app 起動時に固定する。
- 変更候補: `main.go` または bootstrap の小さい初期化箇所。
- 完了条件: 業務コード側に logger 引数が増えない。
- 停止条件: SQLite handler のために長い独自実装が必要になる。

### P1 backend representative logs

- 目的: 代表的な状態遷移と失敗分類だけを出す。
- 変更候補: job lifecycle service。
- 変更候補: phase execution service。
- 変更候補: import / artifact service。
- 上限: 1 file あたり 1 箇所から 3 箇所。
- 完了条件: ログ追加が業務処理の読みやすさを壊さない。

### P2 frontend representative logs

- 目的: browser console で runtime event 破棄と Wails 例外だけを確認できるようにする。
- 変更候補: runtime event adapter。
- 変更候補: 共通 gateway invoker が既にある場合だけ gateway。
- 上限: 1 file あたり 1 箇所から 3 箇所。
- 完了条件: gateway method ごとの wrapper を追加しない。

### P3 verification

- backend: `python3 scripts/harness/run.py --suite backend-local`
- frontend: `python3 scripts/harness/run.py --suite frontend-local`
- 確認: 出力ログに secret、全文入力、巨大 payload がない。
- 確認: 差分量が観測価値に対して過剰でない。

## agent 起動入力

- 完成済み実装成果物: 軽量 logger entry と対象 slice の既存コード。
- 変更ファイル: 代表ログを入れる file だけ。
- 禁止: trace id、全 command wrapper、constructor 引数拡散、DTO dump。
- 成果: 追加ログ、追加しない理由、禁止ログ確認、検証結果。
