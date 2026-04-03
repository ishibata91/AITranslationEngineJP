---
name: directing-implementation
description: AITranslationEngineJp 専用。実装要求の正式入口。必要なら `designing-implementation` に active exec-plan の `UI` / `Scenario` / `Logic` を固めさせ、そのまま実装、close まで進める。
---

# Directing Implementation

この skill は実装 lane の入口です。

## 使う場面

- 新機能実装
- 既存機能の拡張
- UI 変更
- 設計判断を少し含む通常実装

## Required Workflow

1. `docs/exec-plans/templates/impl-plan.md` を使って active plan を作成または更新する。
2. `<ctx_loader>` を `distilling-implementation` でスポーンし、active plan と入口情報から最小限の repo 調査を行わせ、facts、constraints、gaps、closeout notes、required reading を整理する。
3. `<task_designer>` を `designing-implementation` でスポーンし、distill 結果を前提に active plan の `UI` / `Scenario` / `Logic` を固める。
4. `<workplan_builder>` を `planning-implementation` でスポーンし、ordered scope、required reading、validation commands を短い brief にする。
5. `<test_architect>` を `architecting-tests` でスポーンし、failing tests、fixtures、acceptance checks、validation commands を先に固定し、必要な test / fixture を最小範囲で実装させる。
6. `<implementer>` を `implementing-frontend` または `implementing-backend` でスポーンして実装する。
7. 実装後は project root で `sonar-scanner` を実行し、`codexmcps` profile の `mcp/sonarqube` を使う `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP` で `project == ishibata91_AITranslationEngineJP` かつ `status == OPEN` issue を取得する。
8. owned scope の確認が必要な時は `--owned-paths` に touched path を渡し、gate 判定へ `CLOSED` issue を混ぜない。
9. owned scope の Sonar issue が残る間は lane に差し戻し、同じ active plan を更新して implementing skill を再実行する。
10. Sonar issue が解消した後に `<review_cycler>` を `reviewing-implementation` で **1** 回だけ実行する。
11. review が `reroute` を返したら lane に差し戻し、同じ active plan を更新して再実行する。
12. 差し戻しが修正されたら､再レビューはせず次へ進む｡
13. 必要な `4humans sync` を整理し、コードベース境界や実行フローが変わる時は `<diagrammer>` を `diagramming-d2` でスポーンして `4humans/class-diagrams/` または `4humans/sequence-diagrams/` の `.d2` / `.svg` を同一変更で更新してから plan を `completed/` へ移す。
14. タスクがアサインされている場合、タスクのstatusをdoneにする。

## 許可すること
- 各エージェントのスポーン
- 各エージェントの契約パケットを読む

## Rules

- active plan と重複確認に必要な最小限の入口情報だけを読み、詳細なコードベース調査は `distilling-implementation` へ委譲する
- `changes/`、`context_board`、`tasks.md` を live 正本にしない
- review が `pass` でも `4humans sync` と plan 完了前に close とみなさない
- active plan の `4humans Sync` には、必要な `4humans/...diagrams` 更新対象を明記する
- skill 権限が曖昧な場合は停止して適切な handoff を選ぶ

## Reference Use

- downstream skill へ handoff する前に `references/directing-implementation.to.<skill>.json` を参照し、渡す情報を揃える。
- downstream skill から受け取る時は、各 skill 側の `references/<skill>.to.directing-implementation.json` を返却契約として扱う。
