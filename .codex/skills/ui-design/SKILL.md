---
name: ui-design
description: Codex 側の UI 設計作業プロトコル。UI 要件契約として主要操作、表示項目、状態差分、実装後確認観点を固定する基準を提供する。
---
# UI Design

## 目的

`ui-design` は作業プロトコルである。
`designer` agent が UI を実装前の見た目 成果物 ではなく UI 要件契約として扱うための、表示項目、操作、状態差分、実装後確認観点の見方を提供する。

実行境界、正本、引き継ぎ、stop / 戻し は [design-bundle](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md) を参照する。

## 対応ロール

- `designer` が使う。
- 呼び出し元は `implement_lane` または人間とする。
- 返却先は 人間レビュー または `implement_lane` とする。
- 担当成果物は `ui-design` の出力規約で固定する。

## 入力規約

- UI が関係する task で表示項目、操作、状態差分を固定する時
- 状態、variant、responsive、accessibility の差分を整理する時
- repo 側 `ui-design.md` に UI 契約と実装後確認観点を残す時
- 入力に 根拠参照、担当者、承認状態が不足する場合は推測で補わない。

## 外部参照規約

- エージェント実行定義とツール権限は [designer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.toml) の 書き込み許可 / 実行許可 とする。
- 雛形: [ui-design.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/templates/task-folder/ui-design.md)
- 実行定義 skill: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md)
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

### 拘束観点

- user-facing text、表示項目、主要 action、button enablement
- page section、状態 variant、layout constraints、accessibility
- 読み込み中、空、エラー、disabled、progress、再試行、成功
- desktop / mobile で破綻してはいけない条件と実装後確認観点

## 判断規約

- UI は見た目 成果物 ではなく実装が満たす契約として固定する
- 実装前の見た目 成果物 を新規必須にしない
- 細かな visual polish は実装後に人間が実物を確認して直す
- generic な AI 風 UI や過剰な装飾を要求しない

- UI 契約 と シナリオ の責務を分ける
- desktop と mobile の破綻条件を実装後確認観点として残す
- user-facing text は日本語を優先する

## 出力規約

- 出力は判断結果、根拠参照、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。

## 完了規約

- task 内成果物 が承認状態、根拠参照、未決事項を含んでいる。
- 人間レビュー が必要な判断を AI だけで完了扱いにしていない。
- 表示項目、主要 action、button enablement を確認した。
- 状態、variant、responsive、overflow リスク を実装後確認観点として確認した。
- `ui-design.md` は UI 要件契約と確認観点に限定した。

## 停止規約

- UI が不要で `plan.md` の `ui_design` が `N/A` の時
- プロダクト frontend code を実装する時
- docs 正本へ UI 仕様を反映するだけの時
- 実装前の見た目 成果物 を UI の必須 成果物 にしない
- プロダクトコード 実装へ踏み込まない
- 未承認で docs 正本化しない
- 停止時は不足項目、衝突箇所、戻し先を返す。
- 実装前の見た目 成果物 を必須 成果物 として復活させなかった場合は停止する。
- プロダクト frontend implementation に踏み込まなかった場合は停止する。
- human が実装後に確認すべき visual polish を隠さなかった場合は停止する。
