---
name: implementation-workflow
description: メインエージェントが implementer と implementation_reviewer を fresh で起動し、実装HITLまでを進めるオーケストレーター。承認済み設計を実装して人間が確認する時に使う。
---

# Implementation Workflow

## 責務

実装を担当する agent と実装を検証する agent を順に起動し、実装HITLで停止する。

メインエージェントは実装、テスト、検証を行わない。
メインエージェントは成果物を書き換えない。

## 入力

- 人間が承認済みの `plan.md`、`design.md`、`spec.md`。
- 対象 repository。
- 調査結果または追加制約がある場合は該当成果物。
- `implementation.md` の出力先。

## agentを起動する

| 順序 | agent | context | 担当 | 成果物 |
| --- | --- | --- | --- | --- |
| 1 | `implementer` | fresh | 実装、テスト、最終検証 | product code、test、`spec.md`、`implementation.md` |
| 2 | `implementation_reviewer` | fresh | 承認済み設計と実装結果の照合 | 検証結果 |

`implementer` へ承認済み成果物、repository、追加制約、`implementation.md` のpathを渡す。
`implementer` へ `docs/exec-plans/templates/task-folder/implementation.md` を雛形として渡す。
`implementer` の完了後に `implementation_reviewer` へ同じ入力、変更差分、テスト結果、`implementation.md` を渡す。

workflowが直接起動する対象は表にある二つのfresh agentだけとする。
fork、親文脈を継承するagent、下位agentを起動しない。
agentへメインエージェントの会話文脈を継承しない。
agentが必要とする情報は入力と成果物だけで渡す。

## agentを維持する

起動した二つのagentを閉じない。
agentの識別子を保持する。
追加作業は新しいagentへ渡さず、同じagentを再開して依頼する。

人間の指摘をメインエージェントが要約または言い換えてagentへ渡さない。
人間は `implementation.md` の `人間の指摘`、または変更対象の成果物を直接変更する。
人間が記入した内容をagentが削除または言い換えない。
再開したagentは変更済み成果物を読み、必要な作業だけを続ける。

## HITL

`implementation_reviewer` が通過を返した直後に実装HITLを置く。

人間へ次を示して停止する。

- 変更差分。
- `implementation.md` のpath。
- テスト結果。
- `implementation_reviewer` の判定と根拠。
- 開いたまま維持しているagent。

人間が成果物を変更した場合は、同じ `implementer` を再開し、続けて同じ `implementation_reviewer` を再開する。
再検証が通過した後に実装HITLへ戻る。

人間が `design.md` または `spec.md` を変更した場合は設計承認を無効にし、承認を行った `design-workflow` または `fix-workflow` へ戻す。
人間が明示的に承認した時だけ完了する。

## 返す成果物

- 変更差分と `implementation.md` のpath。
- `implementer` と `implementation_reviewer` の識別子。
- テスト結果と検証結果。
- 実装HITLの承認状態。

## 停止条件

- 設計成果物に対する人間承認がない。
- 必須成果物が不足する。
- agentをfreshで起動できない。
- 同じagentを再開できない。
- 検証の否決を承認済み設計の範囲で解消できない。

停止時は不足項目、agentの状態、未解決の検証結果を返す。
