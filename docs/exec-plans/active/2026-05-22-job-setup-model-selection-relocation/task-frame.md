# task 枠: ジョブセットアップ画面廃止と AI モデル選択移動

## 目的

ジョブセットアップ画面を廃止する。
AI モデル選択は、単語翻訳、NPC ペルソナ生成、本文翻訳の各段階画面で行えるようにする。

## 入力

- 人間依頼: ジョブセットアップ画面は実質 AI モデル選択だけなので廃止し、各ページへ AI モデル選択を移動する。
- [docs/detail-specs/translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-setup.md)
- [docs/detail-specs/term-translation-phase.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/term-translation-phase.md)
- [docs/detail-specs/persona-generation-phase.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/persona-generation-phase.md)
- [docs/detail-specs/body-translation-phase.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/body-translation-phase.md)
- [docs/screen-design/screens/translation-management.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/screen-design/screens/translation-management.md)
- [docs/screen-design/screens/translation-input-review.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/screen-design/screens/translation-input-review.md)
- [docs/screen-design/screens/translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/screen-design/screens/translation-job-setup.md)
- [docs/screen-design/screens/term-translation-phase.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/screen-design/screens/term-translation-phase.md)
- [docs/screen-design/screens/persona-generation-phase.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/screen-design/screens/persona-generation-phase.md)
- [docs/screen-design/screens/body-translation-phase.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/screen-design/screens/body-translation-phase.md)

## 完了条件

- 翻訳管理の段階表示からジョブセットアップ画面が消える。
- 入力データ確認から単語翻訳画面へ進む導線が成立する。
- 単語翻訳、NPC ペルソナ生成、本文翻訳の各画面で、開始前に AI サービス、モデル、実行方式を選択できる。
- 各段階の開始時と再試行時に、その段階の AI 設定だけを利用する。
- 各翻訳段階の画面は、開始判断、実行状態、結果判断に必要な要素へ整理されている。
- 認証不足は、対象段階の開始可否として判断できる。
- 秘密値本体は UI、DTO、read model、log、error summary に出さない。
- Storybook で変更した部品と表示状態を人間レビューできる。

## 設計前提

- 入力データの登録と確認は、翻訳入力レビュー画面が担当する。
- 翻訳ジョブの作成は、入力データ確認後に成立する。
- 翻訳ジョブ作成時に 3 段階すべての AI モデル選択を求めない。
- AI モデル選択は、翻訳段階を開始する直前の判断として扱う。
- 翻訳段階画面は、常時表示する診断情報を減らし、詳細情報は必要な状態だけで表示する。
- 既存の AI サービス設定管理は、接続先と認証状態の正本として残す。

## 非対象

- AI サービス設定画面の廃止。
- 共通辞書、共通ペルソナ、入力データ登録の責務変更。
- 本文翻訳完了後の出力管理導線の変更。
- docs 正本本文の直接更新。
- 人間設計レビュー前のプロダクトコード変更。
