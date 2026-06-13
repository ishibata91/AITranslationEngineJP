---
name: finalization-module
description: "正本化判断、正本反映、作業 commit、local merge、merge 後検証、completed 移動、merge 結果 commit を固定する出口モジュール。正本反映対象は docs/architecture.md に限定。TRIGGER when: 実装モジュールの最終検証が通過し、正本化判断と作業 branch の統合先 branch への取り込みを固定する必要がある。SKIP when: 最終検証通過前、または作業 branch が固定されていない。"
---
# Finalization Module

## 目的

`finalization-module` は、実装モジュールの最終検証通過後に、仕様正本化、作業 commit、作業 branch の統合先 branch への local merge、completed 移動、merge 結果 commit を固定するモジュール skill である。
人間承認済みの恒久仕様だけを正本へ反映し、remote repository の変更は行わない。
conflict 発生時だけ `conflict_resolver` agent を起動して conflict 解消を任せる。

正本反映対象は次に限定する:

- `docs/architecture.md`（層構成、依存方向、強い制約、Wails 境界の正本）

画面の正本は Storybook（story と svelte コンポーネント）であり、frontend source として作業 commit に含める。docs 正本反映の対象にはしない。
`docs/detail-specs/`、`docs/usecases/`、`docs/screen-design/` は廃止済みのため正本反映対象に含まれない。

## 呼び出し関係

- 呼び出し元: 人間、または上位モジュール skill。
- 返却先: 呼び出し元。
- モジュールが呼ぶ下位 skill: `conflict-resolver`（merge 時に conflict が発生した時のみ）。正本反映は本 SKILL.md 内の手順で Claude 本体が直接行う。
- モジュールが呼ぶ下位 agent: `conflict_resolver`（merge 時に conflict が発生した時のみ）。
- `詳細仕様正本反映`、`local merge`、`merge 後検証`、`completed 移動`、`merge 結果 commit` は Claude 本体が直接実行する。

## 入口条件

- 実装モジュールの `最終検証` が通過している、または成立条件不成立で停止理由が固定されている。
- 想定 Y/N（仕様変更または仕様追加がある）が `design-module` または `investigation-module` の `想定 Y/N 評価` で固定されている。「人間承認済みの恒久仕様がある」は `正本化判断` 後に確定する。
- 作業 branch が `claude/<task-id>` として存在する。

## 出口条件

- `作業 commit` の commit hash が `plan.md` に固定されている。
- 仕様変更または仕様追加がある場合は、`正本化判断` の結果（人間承認済みの恒久仕様がある／ない／後続課題に切り出す）が固定されている。
- `local merge` が統合先 branch（既定 `master`）へ完了している、または成立条件不成立で停止理由が固定されている。
- conflict が発生した場合は `conflict_resolver` が解消し、`merge 後検証` が通過している。
- active plan folder が `docs/exec-plans/completed/<task-id>/` へ移動済みである。
- `merge 結果 commit` の commit hash が記録されている。

## 担当 artifact

| 成果物ID | 担当 | 依存対象 | 起動先 |
| --- | --- | --- | --- |
| `正本化判断` | Claude 本体 | `最終検証` | なし |
| `正本反映` | Claude 本体 | `正本化判断` | なし |
| `作業 commit` | Claude 本体 | 全完了または停止済み成果物 | なし |
| `local merge` | Claude 本体 | `作業 commit` | なし |
| `conflict 解消` | `conflict_resolver` | `local merge`（conflict 発生時のみ） | `conflict_resolver` |
| `merge 後検証` | Claude 本体 | `local merge`, `conflict 解消?` | なし |
| `completed 移動` | Claude 本体 | `merge 後検証` | なし |
| `merge 結果 commit` | Claude 本体 | `completed 移動` | なし |

## decision table

複数想定が同時に Y なら、各行で「要」になった artifact を全部作る。

| 想定 | 正本化判断 | 正本反映 |
| --- | --- | --- |
| `docs/architecture.md` への反映が要る | 要 | - |
| 人間承認済みの恒久仕様がある | - | 要 |

「人間承認済みの恒久仕様がある」は、`正本化判断` 後に人間が恒久仕様として承認した場合だけ Y になる。

## 各 artifact の詳細

### 正本化判断

- Claude 本体が固定する。
- 次を `plan.md` に記録する。
    - 反映対象（`docs/architecture.md`）、影響範囲、対象 docs パス候補。
    - 恒久仕様として承認するか、後続課題に切り出すか、廃案にするかの判断。
    - 人間承認状態（承認済み、保留、差し戻し）。
- 人間承認なしに「恒久仕様として承認」を確定しない。
- 人間承認を依頼する直前に、active plan folder に `summary.md` を一時作成し、承認終了後に削除する。固定セクションは「概要」と「図」の 2 つ。
- 承認済みなら `正本反映` の起動入力（対象 docs パス、反映する変更要点、根拠 active plan）を固定する。

### 正本反映

- Claude 本体が人間承認済みの恒久仕様だけを正本へ反映する。反映時は、変更前 / 変更後 / 根拠 active plan を `plan.md` に記録する。
- 反映対象は `docs/architecture.md` に限定する。
- 人間承認なしの本文変更は行わない。

### 作業 commit

- Claude 本体が固定する。
- 次を行う。
    - 全完了または停止済み成果物を確認する。
    - 作業 branch（`claude/<task-id>`）上で local commit を作る。
    - commit hash、変更ファイル一覧、検証結果、残留リスクを `plan.md` に記録する。
- remote repository を変更しない（push、tag push、remote branch delete は行わない）。

### local merge

- Claude 本体が固定する。
- 統合先 branch（既定 `master`）へ作業 branch を `git merge --no-ff` で取り込む。
- `plan.md` の source branch、target branch、commit hash が対応しているかを確認してから実行する。
- conflict が発生した場合は本体側で commit せず、`conflict_resolver` agent を Task ツールで起動して `conflict 解消` を任せる。
- remote repository を変更する command（push、tag push、remote branch delete）は実行しない。

### conflict 解消

- conflict が発生した場合だけ `conflict_resolver` を Task ツールで起動して作らせる。
- 下位 skill: `conflict-resolver`。
- 起動入力に含める内容: conflict file 一覧、source branch、target branch、active plan folder、`plan.md` の根拠参照。
- 解消結果（採用判断、根拠、変更ファイル）を `plan.md` に追記する。
- 解消が仕様判断、設計変更、レーン外の再実装を必要とする場合は停止し、呼び出し元へ戻す。

### merge 後検証

- Claude 本体が固定する。
- `plan.md` の検証結果と conflict 解消範囲から必要な検証 suite を選んで実行する。
- 通過結果または未実行理由を `plan.md` に記録する。

### completed 移動

- Claude 本体が固定する。
- `merge 後検証` の通過結果または未実行理由が記録された後だけ、active plan folder を `docs/exec-plans/completed/<task-id>/` へ移動する。

### merge 結果 commit

- Claude 本体が固定する。
- local merge 結果、conflict 解消結果、completed 移動、検証記録を含めて local commit を作る。
- remote repository を変更しない。

## 不変条件

- 人間承認なしの docs 正本本文を変更しない。
- `正本化判断` を経ずに `正本反映` へ進めない。
- 正本反映の対象は `docs/architecture.md` に限定する。
- `docs/architecture.md` への反映が要るのに `正本化判断` を固定できない場合は停止する。恒久仕様承認があるのに `正本反映` を固定できない場合も停止する。
- remote repository の変更は行わない（push、tag push、remote branch delete）。
- `git reset --hard`、`git checkout --`、`git clean` など destructive command を実行しない。
- 作業 commit を作らずに `local merge` へ進めない。
- conflict 発生時は `conflict_resolver` を起動し、本体は conflict file を直接書き換えない。

## 返す成果物

- 正本化判断: 反映対象、判断結果、人間承認状態。
- 正本反映: 反映 docs パス、変更要点。
- 作業 commit: commit hash、変更ファイル一覧、検証結果、残留リスク。
- local merge: merge command、対象 branch、merge 結果、conflict 有無。
- conflict 解消: conflict file、採用判断、根拠参照、解消結果（conflict 発生時のみ）。
- merge 後検証: 実行 command、成否、証跡位置、未実行理由。
- completed 移動: active plan folder から completed plan folder への移動結果。
- merge 結果 commit: local commit の hash、対象 branch、commit 対象差分。
- 停止判定: 停止理由、不足項目、戻し先。

## 作業を止める条件

- `最終検証` が通過していない、かつ停止理由も固定されていない。
- `docs/architecture.md` への反映が要るのに `正本化判断` を固定できない。
- 恒久仕様承認があるのに `正本反映` を固定できない。
- 作業 branch が `claude/<task-id>` として存在しない、または local commit を作れない。
- source branch と作業 commit hash が対応しない、または target branch が人間指定なしで `master` 以外。
- conflict 解消が仕様判断、設計変更、レーン外の再実装を必要とする（`conflict_resolver` の停止理由を受けて呼び出し元へ戻す）。
- `merge 後検証` が失敗し、conflict 解消範囲を超える修正が必要。
- `completed 移動` または `merge 結果 commit` を実行できない。
- 停止時は不足項目、衝突箇所、戻し先を返す。
