---
name: merge-lane
description: active plan ごとの local merge、conflict 解消、merge 後検証、completed 移動、merge 結果 commit を固定する作業プロトコル。
---
# Merge Lane

## 目的

`merge-lane` は、各レーンが worktree 上で作成した local branch を統合先 branch へ取り込み、active plan を completed archive へ移す作業プロトコルである。
`merge_lane` がマージ準備確認、local merge、conflict 解消、merge 後検証、completed 移動、merge 結果 commit を管理する時に使う。

## 対応ロール

- `merge_lane` が使う。
- 呼び出し元は人間、`implement_lane`、`fix_lane`、`exploration_test_lane`、`light_change_lane` とする。
- 返却先は人間または呼び出し元レーンとする。
- 担当成果物は `マージ準備確認`、`local merge`、`conflict 解消`、`merge 後検証`、`completed 移動`、`merge 結果 commit` とする。

## 入力規約

- 呼び出し元: この skill を呼び出した人間またはレーン。
- 作業計画フォルダ: task 内成果物を置く `docs/exec-plans/active/<task-id>/`。
- マージ準備入力: 各レーンが作成した merge 前提の引き継ぎ情報。
- 作業worktree: source branch を checkout した worktree。
- source branch: 統合元の local branch。
- target branch: 統合先の local branch。
- commit hash: source branch の作業 commit。
- 検証結果: merge 前に各レーンが確認した検証結果。
- レビュー結果: merge 前に各レーンが確認したレビュー結果。
- 残留リスク: merge 前に各レーンが残した注意事項。

## 外部参照規約

- エージェント実行定義と実行境界は [merge_lane.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/merge_lane.toml) に従う。
- 作業流れ正本は [.codex/README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/README.md) とする。
- active plan folder 規約は [active/README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/README.md) に従う。
- completed plan folder 規約は [completed/README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/README.md) に従う。
- サンドボックス外で実行してよい command prefix は [.codex/rules/default.rules](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/rules/default.rules) に従う。
- 外部成果物が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

マージレーンの成果物DAGは次を必ず持つ。
各成果物は、`依存対象` の成果物が揃った時だけ着手できる。

| 成果物ID | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- |
| `マージ準備確認` | `merge_lane` | `マージ準備入力` | なし |
| `local merge` | `merge_lane` | `マージ準備確認` | なし |
| `conflict 解消` | `merge_lane` | `local merge` | なし |
| `merge 後検証` | `merge_lane` | `local merge`, `conflict 解消?` | なし |
| `completed 移動` | `merge_lane` | `merge 後検証` | なし |
| `merge 結果 commit` | `merge_lane` | `completed 移動` | なし |

## 判断規約

- `source branch` の既定名は `codex/<task-id>` とする。
- `target branch` の既定名は `master` とする。
- `マージ準備確認` は active plan folder、worktree path、source branch、target branch、commit hash、検証結果、レビュー結果、残留リスクを確認する。
- `source branch` が対象 worktree に checkout されていない場合は local merge へ進めない。
- `target branch` が人間指定なしで `master` 以外の場合は local merge へ進めない。
- `local merge` は remote repository を変更しない local command だけで行う。
- `conflict 解消` は active plan、source branch 差分、target branch 差分から判断できる範囲だけを扱う。
- conflict 解消が仕様判断、設計変更、レーン外の再実装を必要とする場合は停止する。
- `merge 後検証` はマージ準備入力の検証結果と conflict 解消範囲から必要な検証を選ぶ。
- `completed 移動` は `merge 後検証` の通過結果または未実行理由が記録された後だけ行う。
- `merge 結果 commit` は local merge 結果、completed 移動、検証記録を含める。
- `push`、tag push、remote branch delete など remote repository を変更する command は実行しない。

## 非対象規約

- 新規実装、恒久修正、探索テスト、軽量変更の成果物作成は扱わない。
- docs 正本本文の更新は扱わない。
- conflict 解消を超える仕様判断、設計変更、レーン外の再実装は扱わない。
- remote repository の変更は扱わない。
- destructive command は扱わない。

## 出力規約

- マージ準備確認: active plan folder、worktree path、source branch、target branch、commit hash、検証結果、レビュー結果、残留リスクの確認結果を返す。
- local merge: merge command、対象 branch、merge 結果、conflict 有無を返す。
- conflict 解消: conflict file、採用判断、根拠参照、解消結果を返す。
- merge 後検証: 実行 command、成否、証跡位置、未実行理由を返す。
- completed 移動: active plan folder から completed plan folder への移動結果を返す。
- merge 結果 commit: local commit の hash、対象 branch、commit 対象差分を返す。
- 停止返却: 不足項目、衝突箇所、固定できない判断、戻し先を返す。
- 禁止事項: 出力に remote repository の変更結果を含めない。

## 完了規約

- source branch が target branch へ local merge 済みである。
- conflict が残っていない。
- merge 後検証の通過結果または未実行理由が記録されている。
- 作業計画フォルダが `docs/exec-plans/completed/<task-id>/` へ移動済みである。
- merge 結果と completed 移動が local commit 済みである。
- remote repository を変更する command を実行していない。

## 停止規約

- マージ準備入力が不足する場合は停止する。
- source branch または target branch を特定できない場合は停止する。
- source branch が対象 worktree に checkout されていない場合は停止する。
- target branch が人間指定なしで `master` 以外の場合は停止する。
- レビュー必須問題が未解決の場合は停止する。
- conflict 解消が仕様判断、設計変更、レーン外の再実装を必要とする場合は停止する。
- merge 後検証が失敗し、conflict 解消範囲を超える修正が必要な場合は停止する。
- completed 移動を実行できない場合は停止する。
- local commit を作成できない場合は停止する。
- `push`、tag push、remote branch delete など remote repository を変更しそうな場合は停止する。
- `git reset --hard`、`git checkout --`、`git clean` など destructive command が必要な場合は停止する。
