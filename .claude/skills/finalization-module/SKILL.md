---
name: finalization-module
description: "正本化判断、詳細仕様正本反映、作業 commit、マージ準備入力を固定する出口モジュール。TRIGGER when: 実装モジュールの最終検証が通過し、正本化判断、作業 commit、マージ準備入力を固定する必要がある。SKIP when: 最終検証通過前、または作業 branch が固定されていない。"
---
# Finalization Module

## 目的

`finalization-module` は、実装モジュールの最終検証通過後に、仕様正本化、作業 commit、マージ準備入力を固定するモジュール skill である。
人間承認済みの恒久仕様だけを正本へ反映し、local merge と remote 変更は行わない。

## 呼び出し関係

- 呼び出し元: 人間、または上位モジュール skill。
- 返却先: 呼び出し元。
- モジュールが呼ぶ下位 skill: `updating-docs`。
- モジュールが呼ぶ下位 agent: `docs_updater`。

## 入口条件

- 実装モジュールの `最終検証` が通過している、または成立条件不成立で停止理由が固定されている。
- 想定 Y/N（仕様変更または仕様追加がある）が `design-module` または `investigation-module` の `想定 Y/N 評価` で固定されている。「人間承認済みの恒久仕様がある」は `正本化判断` 後に確定する。
- 作業 branch が `claude/<task-id>` として存在する。

## 出口条件

- `作業 commit` の commit hash が `plan.md` に固定されている。
- `マージ準備入力` が `merge_lane` へ渡せる形で揃っている。
- 仕様変更または仕様追加がある場合は、`正本化判断` の結果（人間承認済みの恒久仕様がある／ない／後続課題に切り出す）が固定されている。

## 担当 artifact

| 成果物ID | 担当 | 依存対象 | 起動先 |
| --- | --- | --- | --- |
| `正本化判断` | 呼び出し元 agent | `最終検証` | `docs_updater?` |
| `詳細仕様正本反映` | `docs_updater` | `正本化判断` | `docs_updater` |
| `作業 commit` | 呼び出し元 agent | 全完了または停止済み成果物 | なし |
| `マージ準備入力` | 呼び出し元 agent | `作業 commit` | `merge_lane` |

## decision table

複数想定が同時に Y なら、各行で「要」になった artifact を全部作る。

| 想定 | 正本化判断 | 詳細仕様正本反映 |
| --- | --- | --- |
| 仕様変更または仕様追加がある | 要 | - |
| 人間承認済みの恒久仕様がある | - | 要 |

「人間承認済みの恒久仕様がある」は、`正本化判断` 後に人間が恒久仕様として承認した場合だけ Y になる。設計モジュールの `想定 Y/N 評価` で N でも、`正本化判断` 後に Y へ変わる可能性があり、その場合は `plan.md` の評価結果を更新する。

## 各 artifact の詳細

### 正本化判断

- 呼び出し元 agent が固定する。
- 次を `plan.md` に記録する。
    - 仕様変更または仕様追加の対象（影響範囲、対象 docs パス候補）。
    - 恒久仕様として承認するか、後続課題に切り出すか、廃案にするかの判断。
    - 人間承認状態（承認済み、保留、差し戻し）。
- 人間承認なしに「恒久仕様として承認」を確定しない。
- 承認済みなら `詳細仕様正本反映` の起動入力（対象 docs パス、反映する変更要点、根拠 active plan）を固定する。

### 詳細仕様正本反映

- `docs_updater` を Task ツールで起動して作らせる。
- `docs_updater` は `updating-docs` skill に従い、`docs/detail-specs/` 配下の正本へ承認済み変更だけを反映する。
- 呼び出し元 agent は本文を代筆しない。

### 作業 commit

- 呼び出し元 agent が固定する。
- 次を行う。
    - 全完了または停止済み成果物を確認する。
    - 作業 branch（`claude/<task-id>`）上で local commit を作る。
    - commit hash、変更ファイル一覧、検証結果、残留リスクを `plan.md` に記録する。
- remote repository を変更しない（push、tag push、remote branch delete は行わない）。
- `docs/exec-plans/completed/` への移動は行わない。

### マージ準備入力

- 呼び出し元 agent が固定する。
- `merge_lane` へ渡せる形で次を `plan.md` に記録する。
    - active plan folder のパス。
    - source branch（作業 branch）名。
    - target branch（既定 `master`）名。
    - 作業 commit hash。
    - 最終検証結果。
    - 残留リスク。

## 不変条件

- 人間承認なしの docs 正本本文を変更しない。
- `正本化判断` を経ずに `詳細仕様正本反映` へ進めない。
- 仕様変更または仕様追加があるのに `正本化判断` を固定できない場合は停止する。恒久仕様承認があるのに `詳細仕様正本反映` を固定できない場合も停止する。
- local merge、`docs/exec-plans/completed/` への移動、remote repository の変更は行わない。
- 呼び出し元 agent は `docs_updater` の担当 artifact 本文を代筆しない。

## 返す成果物

- 正本化判断: 仕様変更対象、判断結果、人間承認状態。
- 詳細仕様正本反映: 反映 docs パス、変更要点、`docs_updater` 返却。
- 作業 commit: commit hash、変更ファイル一覧、検証結果、残留リスク。
- マージ準備入力: active plan folder、source branch、target branch、commit hash、最終検証結果、残留リスク。
- 停止判定: 停止理由、不足項目、戻し先。

## 作業を止める条件

- `最終検証` が通過していない、かつ停止理由も固定されていない。
- 仕様変更または仕様追加があるのに `正本化判断` を固定できない。
- 恒久仕様承認があるのに `詳細仕様正本反映` を固定できない。
- 作業 branch が `claude/<task-id>` として存在しない、または local commit を作れない。
- `merge_lane` へ渡す `マージ準備入力` の必須項目（active plan folder、source/target branch、commit hash、最終検証結果）を埋められない。
- 停止時は不足項目、衝突箇所、戻し先を返す。
