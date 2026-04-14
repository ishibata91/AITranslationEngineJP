# Design: ui-mock

## Goal

- task-local UI mock を working copy として作る

## Rules

- path は `docs/exec-plans/active/<task-id>.ui.html` に固定する
- 最終適用先は `docs/mocks/<page-id>/index.html` として plan に記録する
- 正本ではなく working copy として扱う
- 主要導線、状態差分、確認したい edge case を含める
- 実装可能性より、review と handoff の明瞭さを優先する
- `canonicalization_targets` に mock 正本の反映先を追加する
