---
name: conflict-resolver
description: "`finalization-module` 内で `conflict_resolver` agent が使う conflict 解消作業プロトコル。"
---
# Conflict Resolver

## 目的

`conflict-resolver` は、`finalization-module` の `local merge` で発生した conflict だけを解消する作業プロトコルである。
`conflict_resolver` agent が conflict file の採用判断、根拠参照、解消結果を固定する時に使う。
local merge 実行、merge 後検証、completed 移動、merge 結果 commit は本 skill が扱わず、呼び出し元の `finalization-module` が直接担当する。

## 対応ロール

- `conflict_resolver` agent が使う。
- 呼び出し元は `finalization-module`、または同等の入力を渡す上位 agent とする。
- 返却先は呼び出し元とする。
- 担当成果物は `conflict 解消` とする。

## 呼び出し元から渡される情報

- 呼び出し元: この skill を呼び出した agent。
- 作業計画フォルダ: `docs/exec-plans/active/<task-id>/`。
- source branch: 作業 branch 名（既定 `claude/<task-id>`）。
- target branch: 統合先 branch 名（既定 `master`）。
- conflict file 一覧: `local merge` で conflict と判定された file の path。
- 根拠参照: `plan.md` の検証結果、人間設計レビュー記録、active plan 内成果物。

## 作業前に読む正本

- エージェント実行定義と実行境界は [conflict_resolver.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.claude/agents/conflict_resolver.md) に従う。
- active plan の判断資料は `docs/exec-plans/active/<task-id>/plan.md` とする。
- 許可済みコマンドは [settings.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.claude/settings.json) の permissions に従う。
- 外部成果物が不足または衝突する場合は停止し、衝突箇所を返す。

## 担当ロールが判断してよい範囲

- conflict 解消は active plan、source branch 差分、target branch 差分から判断できる範囲だけを扱う。
- 採用判断の根拠を `plan.md` の検証結果または人間設計レビュー記録に対応づける。
- 自動マージできない箇所を最小範囲で書き換え、conflict marker を残さない。
- 解消が仕様判断、設計変更、レーン外の再実装を必要とする場合は停止する。
- `push`、tag push、remote branch delete など remote repository を変更する command は実行しない。
- `git reset --hard`、`git checkout --`、`git clean` など destructive command は実行しない。
- `git merge --continue` の実行は呼び出し元（`finalization-module`）の担当とし、本 skill では実行しない。

## skill が扱わない対象

- `local merge` 自体の実行は扱わない。
- `merge 後検証` の実行は扱わない。
- `completed 移動` は扱わない。
- `merge 結果 commit` は扱わない。
- 新規実装、恒久修正、探索テストは扱わない。
- docs 正本本文の更新は扱わない。
- conflict 解消を超える仕様判断、設計変更、レーン外の再実装は扱わない。
- remote repository の変更は扱わない。

## 返す成果物

- 判断結果: conflict 解消の完了、未完了、停止の判定を返す。
- conflict file: 解消対象の file path 一覧を返す。
- 採用判断: 各 conflict 箇所の採用元（source / target / 統合）と根拠を返す。
- 根拠参照: 採用判断の根拠にした `plan.md` 記録、active plan 内成果物、人間設計レビュー記録を返す。
- 解消結果: 変更ファイル一覧、conflict marker を残していないことの確認結果を返す。
- 不足情報: 解消を完了できない不足項目を返す。
- 禁止事項: 出力に `local merge` 実行、`merge 後検証` 実行、`completed 移動`、`merge 結果 commit`、remote repository の変更を含めない。

## 作業を完了できる条件

- 渡された conflict file 一覧すべてに対し、conflict marker が残っていない。
- 各 conflict 箇所の採用判断と根拠が返っている。
- 変更ファイル一覧が返っている。
- 解消が仕様判断、設計変更、レーン外の再実装を必要としない範囲で完了している。

## 作業を止める条件

- conflict file 一覧、source branch、target branch、作業計画フォルダのいずれかが不足する場合は停止する。
- 採用判断の根拠が `plan.md` または active plan 内成果物から得られない場合は停止する。
- conflict 解消が仕様判断、設計変更、レーン外の再実装を必要とする場合は停止する。
- レビュー必須問題が未解決の場合は停止する。
- `push`、tag push、remote branch delete など remote repository を変更しそうな場合は停止する。
- `git reset --hard`、`git checkout --`、`git clean` など destructive command が必要な場合は停止する。
- 停止時は不足項目、衝突箇所、戻し先を返す。
