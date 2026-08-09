---
name: design-workflow
description: メインエージェントが設計を作成し、画面変更時は Storybook 人間レビューを設計HITLで行い、design_reviewer を fresh で起動して設計承認まで進めるオーケストレーター。プロダクト変更の設計を作成して人間が承認する時に使う。
---

# Design Workflow

## 責務

設計作業の順序と設計HITLの位置を固定する。

メインエージェントは最初に `design-protocol` と `specification-protocol` を読む。
設計判断はメインエージェントが `design-protocol` に従って行う。
仕様作成はメインエージェントが `specification-protocol` に従って行う。
人間向けの説明が必要な場合はメインエージェントが `presentation` に従う。
画面の見た目が変わる場合は、メインエージェントが `storybook-module` に従い、設計HITLで見た目を固定する。
`design-protocol` をagentとして起動しない。
`plan_compactor` は、設計HITL前の exec-plan を整理するためにだけ起動する。`context_reviewer` は起動しない。

## 入力

- 人間の要求。
- 対象 repository。
- task-id。
- 統合先branch。指定がない場合は `master`。
- 確定済みの事実と制約。

作業成果物は `docs/exec-plans/active/<task-id>/` に置く。

## 設計を作成する

メインエージェントは作業branchを `codex/<task-id>` として準備する。
メインエージェントは `docs/exec-plans/templates/task-folder/plan.md`、`design.md`、`spec.md`、`pending.md`、`references.md` を雛形として使う。

次の順序で進める。

| 順序 | 担当 | context | 作業 | 成果物 |
| --- | --- | --- | --- | --- |
| 1 | `codebase-explorer` | fresh | 要求に関係する実装、呼び出し元、依存先、testの探索 | sourceの場所と探索結果 |
| 2 | メインエージェント | 現在のtask | 要求の整理、参照、設計、仕様、未決定事項の分離 | 作業branch、`plan.md`、`design.md`、`spec.md`、`pending.md`、`references.md` |
| 3 | `plan_compactor` | fresh | 規約違反の移動と重複の集約 | 整理済み exec-plan と未整理箇所 |
| 4 | メインエージェント | 現在のtask | 画面の見た目が変わる場合だけ `storybook-module` に従って表示を固定し、人間が確認できる story を作る | story、fixture、表示コンポーネント |
| 5 | `design_reviewer` | fresh | 要求、設計、仕様、実ソースの照合 | 検証結果 |

`codebase-explorer` へ要求、repository、確定済みの事実と制約、探索対象を渡す。
メインエージェントの会話文脈と設計案を `codebase-explorer` へ渡さない。

メインエージェントは探索結果を受け取った後に `references.md` へ source と外部資料の所在だけを記録する。メインエージェントは `design-protocol` と `specification-protocol` に従い、`plan.md`、`design.md`、`spec.md`、`pending.md` を作成する。
未決定事項とブロッカーは `pending.md` にだけ置く。解決した項目は結論を正本へ反映してから `pending.md` から削除する。
メインエージェントは、設計HITL前に `plan_compactor` へ task folder と `protocols/docs/exec-plans/coding.md` を渡す。`plan_compactor` は `plan-compaction` に従い、規約違反の移動と重複の集約だけを直接行う。`log.jsonl` は渡さない。
`plan_compactor` が整理できない箇所を返した場合は、メインエージェントが task の設計判断として扱う。`plan_compactor` に設計判断を再依頼しない。
人間向けの説明が必要な場合だけ `presentation` を読む。

画面変更の有無は `plan.md`、`design.md`、`spec.md` から判定する。画面変更がある場合は、`design_reviewer` を起動する前に `storybook-module` を読み、Storybook 人間レビューを完了する。

`pending.md` が空でない場合は、design-review と設計HITLを開始せずに停止する。

成果物と必要な Storybook 人間レビューが揃い、`pending.md` が空の場合に `design_reviewer` へ要求、`plan.md`、`design.md`、`spec.md`、`pending.md`、`references.md`、repository、語彙の正本、確定済みの事実と制約を渡す。`log.jsonl` は渡さない。
メインエージェントの会話文脈を `design_reviewer` へ渡さない。

workflowが起動するagentは `codebase-explorer`、`plan_compactor`、`design_reviewer` だけとする。
forkまたは親文脈を継承するagentを起動しない。

`storybook-module` が固定した story、fixture、表示コンポーネントは、承認済み設計の一部として `implementation-workflow` へ渡す。

## agentを維持する

起動した三つのagentを閉じない。
三つのagentの識別子を保持する。
追加の探索は同じ `codebase-explorer` を再開して依頼する。
再検証は同じ `design_reviewer` を再開して依頼する。

設計レビューの回数と3回目の扱いは `design-review` に従う。各回の判定と更新した正本を指す event を `log.jsonl` へ追記する。`log.jsonl` を再レビューの入力へ渡さない。

人間の指摘をメインエージェントが要約または言い換えてagentへ渡さない。
人間は `plan.md`、`design.md`、`spec.md`、`pending.md`、`references.md` の該当箇所を直接変更する。
メインエージェントは人間が記入した内容を削除または言い換えない。
メインエージェントは変更済み成果物を読み、必要な作業だけを続ける。

## HITL

`design_reviewer` が通過を返した後に設計HITLを置く。

Storybook 人間レビューで設計または仕様の変更が必要になった場合は、メインエージェントが `design.md` または `spec.md` を直してから Storybook 人間レビューを完了し、その後に `design_reviewer` を起動する。表示範囲だけの指摘は `storybook-module` の範囲で直す。

人間へ次を示して停止する。

- `plan.md`、`design.md`、`spec.md`、`pending.md`、`references.md` のpath。
- `design_reviewer` の判定と根拠。
- 開いたまま維持している三つのagent。
- 画面の見た目が変わる場合は、承認された story と画面状態。

人間が成果物を変更した場合は、最初に `pending.md` が空であることを確かめる。空でない場合は、設計HITLまたは `design_reviewer` の再開をせずに停止する。空の場合は、`design-review` が定める回数の上限内でメインエージェントが必要な作業を続け、同じ `design_reviewer` を再開する。

人間の設計HITL後に成果物が変更された場合は、`plan_compactor` を再開しない。

回数の上限に達した後に人間が成果物を変更した場合は、変更後の成果物を未検証として停止する。Codex 本体は、再レビューできない理由と、変更を別の設計作業として始める必要があることを人間へ返す。
追加の探索が必要な場合だけ同じ `codebase-explorer` を再開する。
設計レビューの書き直しで固定した見た目が変わる場合は、`storybook-module` の Storybook 人間レビューへ戻ってから、`design-review` が定める回数の範囲で同じ `design_reviewer` を再開する。
再検証が通過した後に設計HITLへ戻る。

人間が明示的に承認した時だけ完了する。

## 返す成果物

- `plan.md`、`design.md`、`spec.md`、`pending.md`、`references.md` のpath。
- `codebase-explorer`、`plan_compactor`、`design_reviewer` の識別子。
- 検証結果。
- 画面の見た目が変わる場合は、承認された story、fixture、表示コンポーネント、Storybook 人間レビューの承認状態。
- 設計HITLの承認状態。

## 停止条件

- 必須入力が不足する。
- `pending.md` が空でない。
- agentをfreshで起動できない。
- 同じagentを再開できない。
- `design_reviewer` の否決を成果物の変更で解消できない。
- `design-review` が定める回数の上限後も `design_reviewer` が `目的未固定` または否決を返す。
- 画面変更があるが、`storybook-module` の完了条件を満たせない。

停止時は不足項目、agentの状態、未解決の検証結果を返す。
