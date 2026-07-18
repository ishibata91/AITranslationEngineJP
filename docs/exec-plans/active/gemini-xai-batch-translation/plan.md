# Plan: gemini-xai-batch-translation

`task_id`: `gemini-xai-batch-translation`

`状態`: 未着手（`translation-persistence` の完了後に着手する）

`依存`: `translation-persistence`（翻訳を対象 plugin 単位で永続化する土台。この上へ乗せる）

## 私がやりたいこと

Gemini・xAI の batch API（非同期の大量翻訳）に対応する。

アプリを閉じても、送信した batch の情報を失わない。後から状態を確認し、結果を対象行へ書き戻せる。

## 決まった仕様

- batch を「非同期の配送方式」として対象 plugin 単位の翻訳永続化の上に乗せる
  - `translation-persistence` が作る対象 plugin 単位の翻訳永続化を土台にする。
  - 同期の翻訳本体と既存行アクセスは変えない。
- batch 固有の永続情報は batch 側の内部に閉じる
  - 外部 batch ID、送信行と結果の対応を持つ。
  - 同期はこの情報を経由しない。
- batch は `provider` の 2 つ目の port として足す
  - 同期の共通実装 `openai_compatible` は変更しない。
  - vendor 固有 API のため、Gemini 用・xAI 用を別に実装する。
- 状態確認は常駐プロセスもバックグラウンドポーリングも作らない
  - 起動時または画面操作の時点だけで行う。
- 接続情報（`provider.Connection`）は永続化せず、都度 UI から渡す。
- 結果の書き戻しは既存経路（`dest` 更新）を再利用する。