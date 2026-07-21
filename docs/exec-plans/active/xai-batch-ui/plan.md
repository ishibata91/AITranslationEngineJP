# Task Plan: xai-batch-ui

`plan.md` は branch 情報と、この task でやること・やらないことの要点を持つ。
設計判断、判断履歴、検証結果、実装結果は持たない。設計は `design.md`、恒久的に残す判断は `docs/changelog.md` に書く。

## やること

- xAI batch 翻訳（非同期の大量翻訳）を frontend の画面から使えるようにする。backend は master 済みで、Wails 公開面 `SubmitBatchTranslation`・`RefreshBatchTranslations`・`GetXAIModels` の binding も生成済み。
- 翻訳実行画面に、xAI（batch）を選ぶ導線、batch 送信の操作、結果の反映の操作を足す。xAI 選択時はモデル一覧を `GetXAIModels` から取る。
- batch で反映した訳が、同期翻訳の結果と区別なく結果一覧へ出ることを実画面で確かめられる状態にする。

## branch 情報

- `execution_branch`: `claude/xai-batch-ui`
- `target_branch`: `master`
- `source_commit`: 87eb46d2

## やらないこと

- backend の翻訳ロジック・provider・永続の変更はしない（master 済みの公開面をそのまま呼ぶ）。
- 実 xAI batch API への課金する自動テストは持たない。実 API 疎通と構造化出力の実挙動確認は手動 e2e で行う。
- Gemini の batch 対応はしない（本 task は xAI のみ）。
