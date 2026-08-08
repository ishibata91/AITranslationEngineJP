---
name: design-workflow
description: メインエージェントが designer と design_reviewer を fresh で起動し、設計HITLまでを進めるオーケストレーター。プロダクト変更の設計を作成して人間が承認する時に使う。
---

# Design Workflow

## 責務

設計を担当する agent と設計を検証する agent を順に起動し、設計HITLで停止する。

メインエージェントは設計、仕様、検証を行わない。
メインエージェントは成果物を書き換えない。

## 入力

- 人間の要求。
- 対象 repository。
- task-id。
- 統合先branch。指定がない場合は `master`。
- 確定済みの事実と制約。

作業成果物は `docs/exec-plans/active/<task-id>/` に置く。

## agentを起動する

| 順序 | agent | context | 担当 | 成果物 |
| --- | --- | --- | --- | --- |
| 1 | `designer` | fresh | 作業branchの準備、要求の整理、設計、仕様 | 作業branch、`plan.md`、`design.md`、`spec.md` |
| 2 | `design_reviewer` | fresh | 要求、設計、仕様、実ソースの照合 | 検証結果 |

`designer` へ要求、repository、task-id、統合先branch、確定済みの事実と制約、成果物のpathを渡す。
作業branchは `codex/<task-id>` として準備させる。
`designer` へ `docs/exec-plans/templates/task-folder/plan.md`、`design.md`、`spec.md` を雛形として渡す。
`designer` の完了後に `design_reviewer` へ要求、三つの成果物、repository、語彙の正本、確定済みの事実と制約を渡す。

workflowが直接起動する対象は表にある二つのfresh agentだけとする。
fork、親文脈を継承するagentを起動しない。
agentへメインエージェントの会話文脈を継承しない。
agentが必要とする情報は入力と成果物だけで渡す。

`designer` がコードベース探索のために起動した `codex-explorer` の識別子を受け取る。

## agentを維持する

起動した二つのagentと `designer` が起動した `codex-explorer` を閉じない。
三つのagentの識別子を保持する。
追加作業は新しいagentへ渡さず、同じagentを再開して依頼する。

人間の指摘をメインエージェントが要約または言い換えてagentへ渡さない。
人間は `plan.md`、`design.md`、`spec.md` の該当箇所を直接変更する。
人間が記入した内容をagentが削除または言い換えない。
再開したagentは変更済み成果物を読み、必要な作業だけを続ける。

## HITL

`design_reviewer` が通過を返した直後に設計HITLを置く。

人間へ次を示して停止する。

- `plan.md`、`design.md`、`spec.md` のpath。
- `design_reviewer` の判定と根拠。
- 開いたまま維持している三つのagent。

人間が成果物を変更した場合は、同じ `designer` を再開し、続けて同じ `design_reviewer` を再開する。
再検証が通過した後に設計HITLへ戻る。

人間が明示的に承認した時だけ完了する。

## 返す成果物

- 三つの設計成果物のpath。
- `designer`、`codex-explorer`、`design_reviewer` の識別子。
- 検証結果。
- 設計HITLの承認状態。

## 停止条件

- 必須入力が不足する。
- agentをfreshで起動できない。
- 同じagentを再開できない。
- `design_reviewer` の否決を成果物の変更で解消できない。

停止時は不足項目、agentの状態、未解決の検証結果を返す。
