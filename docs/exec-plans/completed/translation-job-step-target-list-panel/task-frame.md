# task 枠: 翻訳ジョブステップ処理対象一覧表示パネル

## 依頼

- 翻訳ジョブステップ全体に処理対象一覧表示パネルを追加する。
- フロントエンドコンポーネントは共通化する。

## 作業場所

- 作業場所: `/Users/iorishibata/Repositories/AITranslationEngineJP`
- 作業 branch: `codex/translation-job-step-target-list-panel`
- 統合先 branch: `master`

## 対象

- 対象画面: 翻訳管理のジョブ実行画面。
- 対象段階: 単語翻訳、NPC ペルソナ生成、本文翻訳、翻訳結果確認。
- 対象情報: 選択中ジョブに対して、各段階で扱う処理対象の一覧。

## 今回扱う範囲

- 翻訳ジョブ実行画面へ処理対象一覧表示パネルを追加する。
- フロントエンドの表示部品を `job-run` 共通部品として設計する。
- 人間設計レビュー後に frontend 実装範囲を固定する。
