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
6. `<implementer>` を `implementing-frontend` または `implementing-backend` でスポーンし、assigned lint suite だけを実行しながら実装する。
7. implementing skill の返却後に implementation lane owner (`directing-implementation`) が project root で `sonar-scanner` を実行し、 `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP` で `project == ishibata91_AITranslationEngineJP` かつ `status == OPEN` issue を取得する。
8. owned scope の確認が必要な時は `--owned-paths` に touched path を渡し、gate 判定へ `CLOSED` issue を混ぜない。
9. owned scope の Sonar issue が残る間は lane owner が同じ active plan を更新し、implementing skill を再スポーンして修正を差し戻す。
10. Sonar issue gate が解消した後に `<review_cycler>` を `reviewing-implementation` で **1** 回だけ実行する。
11. review が `reroute` を返したら lane owner が同じ active plan を更新し、implementing skill を再スポーンして修正する。
12. review が `pass` を返した後に implementation lane owner (`directing-implementation`) が project root で `python3 scripts/harness/run.py --suite all` を実行する。
13. `all` が失敗した時は lane owner が同じ active plan を更新し、implementing skill を再スポーンして修正する。
14. 差し戻しが修正されたら､再レビューはせず次へ進む｡
15. 必要な `4humans sync` を整理し、実装の変更または追加があった時は `<diagrammer>` を `diagramming-d2` でスポーンして `4humans/diagrams/processes/` の relevant `.d2` / `.svg` を更新し、構造の変更または追加があった時は `4humans/diagrams/structures/` の relevant `.d2` / `.svg` を同一変更で更新してから plan を `completed/` へ移す。new detail `.d2` を追加する時は `4humans/diagrams/overview-manifest.json` と manifest で紐づく overview `.d2` / `.svg` も同一変更で更新する。
16. タスクがアサインされている場合、タスクのstatusをdoneにする。

## 許可すること
- 各エージェントのスポーン spawn_agentのfork_contextはfalseで呼び出すこと。
- 各エージェントの契約パケットを読む

## Rules

- active plan と重複確認に必要な最小限の入口情報だけを読み、詳細なコードベース調査は `distilling-implementation` へ委譲する
- `changes/`、`context_board`、`tasks.md` を live 正本にしない
- implementing skill には raw lint command 群ではなく suite 名で渡し、lint の中身は repo 側の harness と package script を正本にする
- implementation lane owner (`directing-implementation`) が Sonar issue remediation loop、single-pass review、review pass 後の final harness を一元管理する
- review が `pass` でも `4humans sync` と plan 完了前に close とみなさない
- active plan の `4humans Sync` には、必要な `4humans/diagrams/processes/` と `4humans/diagrams/structures/` の更新対象を明記し、new detail `.d2` を追加する時は `4humans/diagrams/overview-manifest.json` と対応 overview `.d2` / `.svg` も必ず列挙する
- skill 権限が曖昧な場合は停止して適切な handoff を選ぶ

## Reference Use

- downstream skill へ handoff する前に `references/directing-implementation.to.<skill>.json` を参照し、渡す情報を揃える。
- downstream skill から受け取る時は、各 skill 側の `references/<skill>.to.directing-implementation.json` を返却契約として扱う。
