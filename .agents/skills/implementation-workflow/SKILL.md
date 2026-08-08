---
name: implementation-workflow
description: メインエージェントが画面変更時の Storybook 人間レビューと承認済み設計の実装を行い、implementation_reviewer を fresh で起動して実装HITLまで進めるオーケストレーター。承認済み設計を実装して人間が確認する時に使う。
---

# Implementation Workflow

## 責務

実装作業の順序と実装HITLの位置を固定する。

画面の見た目が変わる場合は、メインエージェントが最初に `storybook-module` を読み、Storybook 上で見た目を固定する。
Storybook 人間レビューの承認後に、メインエージェントが `implementation-protocol` を読んで残りの実装と検証を行う。
画面の見た目が変わらない場合は、メインエージェントが `storybook-module` を読まずに `implementation-protocol` を読む。
`implementation-protocol` をagentとして起動しない。

## 入力

- 人間が承認済みの `plan.md`、`design.md`、`spec.md`。
- 対象 repository。
- 調査結果または追加制約がある場合は該当成果物。
- `implementation.md` の出力先。

## 実装する

次の順序で進める。

| 順序 | 担当 | context | 作業 | 成果物 |
| --- | --- | --- | --- | --- |
| 1 | メインエージェント | 現在のtask | 画面の見た目が変わる場合だけ `storybook-module` に従って見た目を固定 | 承認された story、fixture、表示コンポーネント |
| 2 | メインエージェント | 現在のtask | `implementation-protocol` に従って実装、テスト、最終検証 | product code、test、`spec.md`、`implementation.md` |
| 3 | `implementation_reviewer` | fresh | 承認済み設計と実装結果の照合 | 検証結果 |

画面変更の有無は `plan.md`、`design.md`、`spec.md` から判定する。
`storybook-module` が固定した見た目は、後続実装の入力として保つ。

メインエージェントは `docs/exec-plans/templates/task-folder/implementation.md` を雛形として使う。

実装完了後に `implementation_reviewer` へ承認済み成果物、repository、変更差分、テスト結果、`implementation.md` を渡す。
画面の見た目が変わる場合は、承認された story と画面状態も渡す。
メインエージェントの会話文脈を `implementation_reviewer` へ渡さない。

workflowが起動するagentは `implementation_reviewer` だけとする。
forkまたは親文脈を継承するagentを起動しない。

## agentを維持する

起動した `implementation_reviewer` を閉じない。
`implementation_reviewer` の識別子を保持する。
再検証は同じ `implementation_reviewer` を再開して依頼する。

人間は `implementation.md` の `人間の指摘`、または変更対象の成果物を直接変更する。
メインエージェントは人間が記入した内容を削除または言い換えない。
メインエージェントは変更済み成果物を読み、必要な作業だけを続ける。

## HITL

画面の見た目が変わる場合は、`storybook-module` が定める人間レビューを実装前に完了する。
`implementation_reviewer` が通過を返した直後に実装HITLを置く。

人間へ次を示して停止する。

- 変更差分。
- `implementation.md` のpath。
- テスト結果。
- `implementation_reviewer` の判定と根拠。
- 開いたまま維持している `implementation_reviewer`。
- 画面の見た目が変わる場合は、承認された story と画面状態。

人間が成果物を変更した場合は、メインエージェントが `implementation-protocol` に従って必要な作業を続け、同じ `implementation_reviewer` を再開する。
固定した見た目が変わる場合は、先に `storybook-module` の人間レビューへ戻る。
再検証が通過した後に実装HITLへ戻る。

人間が `design.md` または `spec.md` を変更した場合は設計承認を無効にし、承認を行った設計作業へ戻す。
人間が明示的に承認した時だけ完了する。

## 返す成果物

- 変更差分と `implementation.md` のpath。
- `implementation_reviewer` の識別子。
- テスト結果と検証結果。
- 画面の見た目が変わる場合は、Storybook 人間レビューの承認状態。
- 実装HITLの承認状態。

## 停止条件

- 設計成果物に対する人間承認がない。
- 必須成果物が不足する。
- 画面変更があるが、`storybook-module` の完了条件を満たせない。
- `implementation-protocol` が設計変更の必要を返す。
- `implementation_reviewer` をfreshで起動できない。
- 同じ `implementation_reviewer` を再開できない。
- 検証の否決を承認済み設計の範囲で解消できない。

停止時は不足項目、agentの状態、未解決の検証結果を返す。
