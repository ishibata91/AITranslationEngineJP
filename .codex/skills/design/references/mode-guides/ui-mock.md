# Design: ui-mock

## Goal

- task-local UI mock を working copy として作る
- human review と実装者 handoff の両方で、主要操作と状態差分を誤解なく読めるようにする

## 必須の書き方

- 画面要素の有無だけで終えず、主要操作、状態差分、失敗時、再実行時を読める形にする
- UI 論点は必要に応じて `issue`、`background`、`options`、`recommendation`、`reasoning`、`open_risks` で補足する
- 状態遷移は mock 本体か補足のどちらかで明示する
- 固有名詞、既存 field 名、既存 contract 名、mode 名を除き、日本語優先で書く

## Rules

- path は `docs/exec-plans/active/<task-id>.ui.html` に固定する
- 最終適用先は `docs/mocks/<page-id>/index.html` として plan に記録する
- 正本ではなく working copy として扱う
- 実装可能性より、review と handoff の明瞭さを優先する
- `canonicalization_targets` に mock 正本の反映先を追加する

## Work Plan Mapping

- 主要 UI 要素と導線は `Functional Requirements` と `Acceptance Checks` の根拠にする
- 失敗時と再実行時の状態差分は `Acceptance Checks` と `Required Evidence` に接続できる形にする
- human review が必要な UI 判断は `open_questions` と `HITL Status` に接続する

## Avoid

- 状態差分なしの static mock
- UI、データ、コマンドを 1 文で混在させること
- human review 前に確定表現で narrow scope を書くこと
