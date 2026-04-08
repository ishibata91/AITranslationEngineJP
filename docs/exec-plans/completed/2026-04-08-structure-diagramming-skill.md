# 実装計画テンプレート

- workflow: impl
- status: completed
- lane_owner: skill-modification
- scope: `.codex/skills/` と関連 workflow docs の同期
- task_id:
- task_catalog_ref:
- parent_phase:

## 要求要約

- `structure_diagrammer` が使う dedicated skill を追加し、active exec-plan の task-local design から更新対象の既存 backend 図を特定する。
- 対応する図がない時は、新規 component detail 図の作成を判断する。
- proposal では active exec-plan 配下に review 用の構造差分 `.d2` / `.svg` を生成する契約にする。

## 判断根拠

<!-- Decision Basis -->

- `structure_diagrammer` が現在 `diagramming-d2` を直接使っており、対象特定、新規作成判断、差分生成の role が skill 契約として分離されていない。
- proposal lane と close lane の両方で同じ structure 図 role を使うため、agent 専用 skill と handoff 契約を追加した方が責務境界が明確になる。

## 対象範囲

- `.codex/skills/` 配下の新 skill 追加
- `proposing-implementation` / `directing-implementation` の handoff 契約更新
- `.codex/README.md`、`.codex/workflow.md`、`.codex/workflow_activity_diagram.puml` の workflow 同期

## 対象外

- `docs/` 正本の恒久仕様更新
- product code の内容変更
- 既存 baseline の design harness failure 解消

## 依存関係・ブロッカー

- `structure_diagrammer` agent 契約
- `diagramming-d2` の D2 validate / render 規約

## 並行安全メモ

- `proposing-implementation` と `directing-implementation` は shared workflow 契約なので、同じ変更で同期する。
- 既存の dirty worktree にある `.DS_Store` と未コミット workflow 変更は巻き戻さない。

## UI

- N/A

## Scenario

- proposal lane で task-local design 固定後に structure 用 skill が active exec-plan と既存図を読み、更新対象の特定または新規 detail 図の要否判断を行う。
- proposal lane では source of truth を直接書き換えず、review 用差分図だけを active exec-plan 配下へ出力する。
- execution close では同じ skill 契約で承認済み差分を `diagrams/backend/` 正本へ適用する。

## Logic

- 新 skill には `SKILL.md`、`agents/openai.yaml`、`references/permissions.json` を追加する。
- `proposing-implementation -> new skill` と `directing-implementation -> new skill` の handoff JSON を追加する。
- live workflow 文書の `diagramming-d2` 直接参照を、structure 専用 skill 経由へ置き換える。

## 実装計画

<!-- Implementation Plan -->

- 新 skill 名、責務、入出力、停止条件を定義する。
- `structure_diagrammer` agent 契約を new skill 前提へ更新する。
- proposal / close lane の skill 契約と workflow 文書を同期する。
- harness を実行し、baseline failure と今回差分を切り分ける。

## 受け入れ確認

- `structure_diagrammer` の proposal step が new skill 名で説明されている。
- new skill が active exec-plan の `UI` / `Scenario` / `Logic` と既存 `diagrams/backend/` を入力に持つ。
- review 用差分図の出力先が active exec-plan 配下で明記されている。

## 必要な証跡

<!-- Required Evidence -->

- structure harness 結果
- design harness 結果
- final harness 結果

## HITL 状態

- N/A

## 承認記録

- N/A

## review 用差分図

- N/A

## 差分正本適用先

- N/A

## Closeout Notes

- N/A

## 結果

<!-- Outcome -->

- `diagramming-structure-diff` skill を追加し、`structure_diagrammer` が proposal_diff と apply_to_source の両方で使う dedicated skill 契約を導入した。
- `proposing-implementation` と `directing-implementation` に new handoff JSON を追加し、古い `proposing-implementation.to.diagramming-d2.json` は削除した。
- `.codex/README.md`、`.codex/workflow.md`、`.codex/workflow_activity_diagram.puml` を new skill 前提へ同期した。
- `python3 scripts/harness/run.py --suite structure` は通過した。
- `python3 scripts/harness/run.py --suite design` と `python3 scripts/harness/run.py --suite all` は既存 baseline failure のまま失敗した。失敗内容は `docs/core-beliefs.md`、`docs/index.md` の pattern 不足で、今回差分起因ではない。
