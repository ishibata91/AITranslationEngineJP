---
name: design-workflow
description: メインエージェントが設計を作成し、codex-explorer と design_reviewer を fresh で起動して設計HITLまで進めるオーケストレーター。プロダクト変更の設計を作成して人間が承認する時に使う。
---

# Design Workflow

## 責務

設計作業の順序と設計HITLの位置を固定する。

メインエージェントは最初に `design-protocol` と `specification-protocol` を読む。
設計判断はメインエージェントが `design-protocol` に従って行う。
仕様作成はメインエージェントが `specification-protocol` に従って行う。
人間向けの説明が必要な場合はメインエージェントが `presentation` に従う。
`design-protocol` をagentとして起動しない。

## 入力

- 人間の要求。
- 対象 repository。
- task-id。
- 統合先branch。指定がない場合は `master`。
- 確定済みの事実と制約。

作業成果物は `docs/exec-plans/active/<task-id>/` に置く。

## 設計を作成する

メインエージェントは作業branchを `codex/<task-id>` として準備する。
メインエージェントは `docs/exec-plans/templates/task-folder/plan.md`、`design.md`、`spec.md` を雛形として使う。

次の順序で進める。

| 順序 | 担当 | context | 作業 | 成果物 |
| --- | --- | --- | --- | --- |
| 1 | `codex-explorer` | fresh | 要求に関係する実装、呼び出し元、依存先、testの探索 | sourceの場所と探索結果 |
| 2 | メインエージェント | 現在のtask | 要求の整理、設計、仕様 | 作業branch、`plan.md`、`design.md`、`spec.md` |
| 3 | `design_reviewer` | fresh | 要求、設計、仕様、実ソースの照合 | 検証結果 |

`codex-explorer` へ要求、repository、確定済みの事実と制約、探索対象を渡す。
メインエージェントの会話文脈と設計案を `codex-explorer` へ渡さない。

メインエージェントは探索結果を受け取った後に `design-protocol` と `specification-protocol` に従い、三つの成果物を作成する。
人間向けの説明が必要な場合だけ `presentation` を読む。

成果物の作成後に `design_reviewer` へ要求、三つの成果物、repository、語彙の正本、確定済みの事実と制約を渡す。
メインエージェントの会話文脈を `design_reviewer` へ渡さない。

workflowが起動するagentは `codex-explorer` と `design_reviewer` だけとする。
forkまたは親文脈を継承するagentを起動しない。

## agentを維持する

起動した二つのagentを閉じない。
二つのagentの識別子を保持する。
追加の探索は同じ `codex-explorer` を再開して依頼する。
再検証は同じ `design_reviewer` を再開して依頼する。

人間の指摘をメインエージェントが要約または言い換えてagentへ渡さない。
人間は `plan.md`、`design.md`、`spec.md` の該当箇所を直接変更する。
メインエージェントは人間が記入した内容を削除または言い換えない。
メインエージェントは変更済み成果物を読み、必要な作業だけを続ける。

## HITL

`design_reviewer` が通過を返した直後に設計HITLを置く。

人間へ次を示して停止する。

- `plan.md`、`design.md`、`spec.md` のpath。
- `design_reviewer` の判定と根拠。
- 開いたまま維持している二つのagent。

人間が成果物を変更した場合は、メインエージェントが必要な作業を続け、同じ `design_reviewer` を再開する。
追加の探索が必要な場合だけ同じ `codex-explorer` を再開する。
再検証が通過した後に設計HITLへ戻る。

人間が明示的に承認した時だけ完了する。

## 返す成果物

- 三つの設計成果物のpath。
- `codex-explorer` と `design_reviewer` の識別子。
- 検証結果。
- 設計HITLの承認状態。

## 停止条件

- 必須入力が不足する。
- agentをfreshで起動できない。
- 同じagentを再開できない。
- `design_reviewer` の否決を成果物の変更で解消できない。

停止時は不足項目、agentの状態、未解決の検証結果を返す。
