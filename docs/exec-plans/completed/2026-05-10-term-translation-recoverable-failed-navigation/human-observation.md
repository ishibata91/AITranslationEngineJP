# 人間観測記録

## 観測条件

- 観測日: 2026-05-10
- 画面: 単語翻訳フェーズ
- 対象 job: 画面上は `translation-job-15` と見える。
- 対象 phase: `term_translation`

## 人間が見た症状

- 単語翻訳フェーズの操作領域に `開始: active phase run already exists` が表示された。
- 単語翻訳フェーズの操作領域に `中断: phase is not running` が表示された。
- 単語翻訳フェーズの操作領域に `再開: phase is not paused` が表示された。
- 進行状況は `0 / 4,930 件 / AI 対象 4,930 件 / retry_waiting` と表示された。
- 失敗理由として `provider response was invalid` が表示された。
- `次へ進む` は表示されているが、ほか 2 件の確認が必要という補足が見える。

## 期待との差分

- `retry_waiting` であれば、再試行可能な導線が使える必要がある。
- 開始、中断、再開の拒否理由だけが並ぶ状態では、利用者が次に押すべき操作を判断しづらい。
- 未完了一覧から現在の翻訳段階へ戻れる必要がある可能性がある。

## 添付 UI 根拠

- 入力画像: 会話添付画像。
- 追加 screenshot: `tmp/agent-browser/2026-05-10-term-translation-investigate/current.png`
- 追加 screenshot: `tmp/agent-browser/2026-05-10-term-translation-investigate/job6-term-phase.png`

## 注意

- provider 応答不正の根本原因は、今回の観測だけでは固定しない。
- 今回の入口は、失敗後の状態投影と操作導線の不整合である。

