---
name: implementation-workflow
description: メインエージェントが必要なコードベース探索を codebase-explorer へ fresh で委譲し、設計側で承認済みの表示資産を保ちながら実装し、implementation_reviewer を fresh で起動して実装HITLまで進めるオーケストレーター。承認済み設計を実装して人間が確認する時に使う。
---

# Implementation Workflow

## 責務

実装作業の順序と実装HITLの位置を固定する。

画面の見た目が変わる場合は、`design-workflow` が `storybook-module` を使って承認した story、fixture、表示コンポーネントを入力として受け取る。
メインエージェントは `storybook-module` を起動しない。
メインエージェントは `implementation-protocol` を読んで実装と検証を行う。
`implementation-protocol` をagentとして起動しない。

## 入力

- 人間が承認済みの `plan.md`、`design.md`、`spec.md`、`pending.md`、`references.md`。
- 対象 repository。
- 調査結果または追加制約がある場合は該当成果物。
- `implementation.md` の出力先。

## 実装する

次の順序で進める。

| 順序 | 担当 | context | 作業 | 成果物 |
| --- | --- | --- | --- | --- |
| 1 | `codebase-explorer` | fresh | 実装に必要な探索情報が不足する場合だけ、変更対象、呼び出し元、依存先、関連testを探索 | sourceの場所と探索結果 |
| 2 | メインエージェント | 現在のtask | `implementation-protocol` に従って実装、テスト、最終検証 | product code、test、`spec.md`、`implementation.md` |
| 3 | `implementation_reviewer` | fresh | 承認済み設計と実装結果の照合 | 検証結果 |

`pending.md` が空でない場合は実装を開始しない。

承認済み成果物と既存の探索結果だけで変更対象、呼び出し元、依存先、関連testを特定できない場合は、実装前に `codebase-explorer` を起動する。
`codebase-explorer` へ承認済み成果物、`references.md`、repository、追加制約、探索で明らかにする問いを渡す。`log.jsonl` は渡さない。
メインエージェントの会話文脈を `codebase-explorer` へ渡さない。
メインエージェントは探索結果を受け取ってから実装を続ける。

画面変更の有無は `plan.md`、`design.md`、`spec.md` から判定する。
画面変更がある場合は、設計側で承認された story、fixture、表示コンポーネントが入力にあることを確かめる。固定した見た目は後続実装で保つ。

メインエージェントは `docs/exec-plans/templates/task-folder/implementation.md` を雛形として使う。

実装完了後に `implementation_reviewer` へ承認済み成果物、repository、変更差分、テスト結果、`implementation.md` を渡す。
`codebase-explorer` を起動した場合は探索結果も渡す。
画面の見た目が変わる場合は、設計側で承認された story と画面状態も渡す。
メインエージェントの会話文脈を `implementation_reviewer` へ渡さない。

コードベース探索は `codebase-explorer` へ委譲する。
実装結果の検証は `implementation_reviewer` へ委譲する。
実装を委譲するagent、fork、親文脈を継承するagentは起動しない。

## agentを維持する

起動した `implementation_reviewer` を閉じない。
`implementation_reviewer` の識別子を保持する。
再検証は同じ `implementation_reviewer` を再開して依頼する。

`codebase-explorer` を起動した場合は閉じず、識別子を保持する。
追加のコードベース探索は同じ `codebase-explorer` を再開して依頼する。

人間は `implementation.md` の `人間の指摘`、または変更対象の成果物を直接変更する。
メインエージェントは人間が記入した内容を削除または言い換えない。
メインエージェントは変更済み成果物を読み、必要な作業だけを続ける。

## HITL

`implementation_reviewer` が通過を返した直後に実装HITLを置く。

人間へ次を示して停止する。

- 変更差分。
- `implementation.md` のpath。
- テスト結果。
- `implementation_reviewer` の判定と根拠。
- 開いたまま維持している `implementation_reviewer`。
- 起動した場合は、開いたまま維持している `codebase-explorer`。
- 画面の見た目が変わる場合は、承認された story と画面状態。

人間が成果物を変更した場合は、最初に `pending.md` が空であることを確かめる。空でない場合は、実装、実装HITL、`implementation_reviewer` の再開をせずに停止する。空の場合は、メインエージェントが `implementation-protocol` に従って必要な作業を続け、同じ `implementation_reviewer` を再開する。
固定した見た目を変える必要がある場合は、実装を止めて `design-workflow` の Storybook 人間レビューへ戻す。
再検証が通過した後に実装HITLへ戻る。

人間が `design.md` または `spec.md` を変更した場合は設計承認を無効にし、承認を行った設計作業へ戻す。
人間が明示的に承認した時だけ完了する。

## 返す成果物

- 変更差分と `implementation.md` のpath。
- `implementation_reviewer` の識別子。
- 起動した場合は、`codebase-explorer` の識別子と探索結果。
- テスト結果と検証結果。
- 画面の見た目が変わる場合は、設計側の Storybook 人間レビューの承認状態。
- 実装HITLの承認状態。

## 停止条件

- 設計成果物に対する人間承認がない。
- 必須成果物が不足する。
- `pending.md` が空でない。
- 画面変更があるが、設計側で承認された story、fixture、表示コンポーネントがない。
- `implementation-protocol` が設計変更の必要を返す。
- 必要なコードベース探索を `codebase-explorer` で完了できない。
- `implementation_reviewer` をfreshで起動できない。
- 同じ `implementation_reviewer` を再開できない。
- 検証の否決を承認済み設計の範囲で解消できない。

停止時は不足項目、agentの状態、未解決の検証結果を返す。
