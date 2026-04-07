---
name: proposing-implementation
description: AITranslationEngineJp 専用。実装要求の proposal 入口。active exec-plan、implementation distill、task-local design、review 用差分図、HITL を管理する。
---

# Proposing Implementation

この skill は implementation proposal lane の入口です。

## 使う場面

- 新機能実装の方針確定
- 既存機能拡張の task-local design 固定
- execution 前に human review と LGTM が必要な通常実装

## Required Workflow

1. `docs/exec-plans/templates/impl-plan.md` を使って日本語の active plan を作成または更新する。
2. `<ctx_loader>` を `distilling-implementation` でスポーンし、active plan と入口情報から最小限の repo 調査を行わせ、facts、constraints、gaps、closeout notes、required reading を整理する。
3. `<task_designer>` を `designing-implementation` でスポーンし、distill 結果を前提に active exec-plan の `UI` / `Scenario` / `Logic` を固める。
4. review 用差分図が必要な時は `<diagrammer>` を起動し、active plan、差分意図、入力図、出力先を渡して `diagramming-d2` で差分 D2 / SVG を作成させる。
5. active plan の `HITL 状態`、`承認記録`、`review 用差分図` を更新し、human review と LGTM を待つ。
6. LGTM が確認できた時だけ implementation lane owner (`directing-implementation`) へ handoff する。

## 許可すること

- 各エージェントのスポーンは `fork_context=false` で呼ぶ。
- 各エージェントの契約パケットを読む。

## Rules

- active plan と重複確認に必要な最小限の入口情報だけを読み、詳細なコードベース調査は `distilling-implementation` へ委譲する
- `changes/`、`context_board`、`tasks.md` を live 正本にしない
- task-local な設計は active exec-plan の `UI` / `Scenario` / `Logic` に閉じる
- human LGTM 前に `planning-implementation` 以降へ進めない
- review 用差分図は active exec-plan と同じフォルダに置き、source of truth にしない
- skill 権限が曖昧な場合は停止して適切な handoff を選ぶ

## Reference Use

- downstream skill へ handoff する前に `references/proposing-implementation.to.<skill>.json` を参照し、渡す情報を揃える。
- downstream skill から受け取る時は、各 skill 側の `references/<skill>.to.proposing-implementation.json` を返却契約として扱う。
