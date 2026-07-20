# Plan: gemini-xai-batch-translation

`task_id`: `gemini-xai-batch-translation`

## branch 情報

- 作業 branch: `claude/gemini-xai-batch-translation`
- 統合先 branch: `master`
- 分岐元 commit: `65316aeb`

## やることの要点

- 非同期の大量翻訳（batch）に対応する。今回は xAI の batch API だけを対象にする。
- アプリを閉じても、送信した batch の情報を失わない。後から状態を確認し、結果を対象行へ書き戻せる。
- batch を、対象 plugin 単位の翻訳永続化（`translation-persistence`）の上に乗せる。同期の翻訳本体と既存行アクセスは変えない。
- batch 固有の永続情報（外部 batch ID、送信行と結果の対応）は batch 側の内部に閉じる。
- batch は `provider` の 2 つ目の port として足す。同期の共通実装 `openai_compatible` は変えない。
- 状態確認は常駐プロセスもバックグラウンドポーリングも作らない。起動時または画面操作の時点だけで行う。
- 接続情報（`provider.Connection`）は永続化せず、都度 UI から渡す。
- 結果の書き戻しは既存経路（`dest` 更新）を再利用する。
- 依頼の検証条件: 最後の e2e は手動で行う。batch API と batch 管理は仕様を徹底調査し徹底テストする。batch 管理はできる限り単一モジュールでカバレッジ 100% にし、外から見て同期と batch の差が出ないようにする。

## やらないことの要点

- Gemini の batch API は今回対象にしない（本 plan の後続で足す）。
- 同期翻訳経路（`provider.Translator.Translate` と `openai_compatible`）の振る舞いは変えない。
- 常駐プロセス・自動バックグラウンドポーリングは作らない。
- 接続情報の永続化はしない。
- 課金が発生する実 xAI batch API への自動テストはしない（実 API への検証は手動 e2e だけ）。

判断履歴は本 plan に残さない。恒久的に残す判断は `docs/changelog.md` へ書く。
