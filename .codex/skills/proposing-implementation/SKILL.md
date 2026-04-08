---
name: proposing-implementation
description: AITranslationEngineJp 専用。実装要求の proposal 入口。active exec-plan、implementation distill、task-local design、review 用構造差分図、HITL を管理する。
---

# Proposing Implementation

この skill は implementation proposal lane の入口です。

## 使う場面

- 新機能実装の方針確定
- 既存機能拡張の task-local design 固定
- execution 前に human review と LGTM が必要な通常実装

## Required Workflow

1. `docs/exec-plans/templates/impl-plan.md` を使って日本語の active plan を作成または更新する。
2. MCP memory bucket (`repo_conventions`, `recurring_pitfalls`) を読み、今回の task に関係する項目だけを MCP memory recall として整理する。MCP memory は repo 作法と再発失敗の recall に限定し、仕様や設計の正本代替には使わない。
3. `<ctx_loader>` を `distilling-implementation` でスポーンし、active plan と入口情報から最小限の repo 調査を行わせ、facts、constraints、gaps、closeout notes、required reading を整理する。
4. `<task_designer>` を `designing-implementation` でスポーンし、distill 結果を前提に active exec-plan の `UI` / `Scenario` / `Logic` を固める。
5. `<structure_diagrammer>` を `structure diagram diff skill (`diagramming-structure-diff`)` で起動し、active plan、差分意図、入力図、active exec-plan 配下の出力先を渡して、更新対象または new component detail 図を判断させた上で human review 用の構造差分 D2 / SVG を作成させる。
6. active plan の `HITL 状態`、`承認記録`、`review 用差分図`、`差分正本適用先` を更新し、human review と LGTM を待つ。
7. LGTM が確認できた時だけ implementation lane owner (`directing-implementation`) へ handoff する。

## 許可すること

- 各エージェントのスポーンは `fork_context=false` で呼ぶ。
- 各エージェントの契約パケットを読む。

## Rules

- active plan と重複確認に必要な最小限の入口情報だけを読み、詳細なコードベース調査は `distilling-implementation` へ委譲する
- MCP memory から読むのは bucket (`repo_conventions`, `recurring_pitfalls`) の recall だけに限定し、`docs/` 正本の代わりにしない
- `changes/`、`context_board`、`tasks.md` を live 正本にしない
- task-local な設計は active exec-plan の `UI` / `Scenario` / `Logic` に閉じる
- MCP memory recall は今回の task に効く項目だけを context summary に残し、無関係な項目は handoff に混ぜない
- human LGTM 前に `planning-implementation` 以降へ進めない
- review 用差分図は active exec-plan と同じフォルダに置き、source of truth にしない
- `diagrams/backend/` 正本は proposal 中に書き換えず、承認済み差分を execution close で適用する
- skill 権限が曖昧な場合は停止して適切な handoff を選ぶ

## Reference Use

- downstream skill へ handoff する前に `references/proposing-implementation.to.<skill>.json` を参照し、渡す情報を揃える。
- downstream skill から受け取る時は、各 skill 側の `references/<skill>.to.proposing-implementation.json` を返却契約として扱う。
