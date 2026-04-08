---
name: directing-implementation
description: AITranslationEngineJp 専用。承認済み implementation proposal を execution へ進める lane owner。`planning-implementation` 以降の実装、gate、close を管理する。
---

# Directing Implementation

この skill は implementation execution lane の入口です。

## 使う場面

- 新機能実装
- 既存機能の拡張
- UI 変更
- 設計判断を少し含む通常実装

## Required Workflow

1. `docs/exec-plans/templates/impl-plan.md` に基づく承認済み active exec-plan、context summary、review 用差分図、差分正本適用先、LGTM 記録を `proposing-implementation` から受け取る。
2. active plan の `HITL 状態` が承認済みであり、`承認記録`、`review 用差分図`、`差分正本適用先`、`UI`、`Scenario`、`Logic` が proposal 側で固定済みであることを確認する。task-local design skill (`designing-implementation`) はこの execution lane では再実行しない。
3. proposal handoff に `mcp_memory_recall` がある時は内容を確認し、実装判断に効くものだけを planning へ引き継ぐ。
4. `<workplan_builder>` を `planning-implementation` でスポーンし、ordered scope、required reading、validation commands、relevant な repo guardrail を短い brief にする。
5. `<test_architect>` を `architecting-tests` でスポーンし、failing tests、fixtures、acceptance checks、validation commands を先に固定し、必要な test / fixture を最小範囲で実装させる。
6. `<implementer>` を `implementing-frontend` または `implementing-backend` でスポーンし、assigned lint suite だけを実行しながら実装する。
7. implementing skill の返却後に implementation lane owner (`directing-implementation`) が project root で `sonar-scanner` を実行し、Sonar MCP の `search_sonar_issues_in_projects` で `project == ishibata91_AITranslationEngineJP` かつ `status == OPEN` issue を取得する。
8. owned scope の確認が必要な時は `--owned-paths` に touched path を渡し、gate 判定へ `CLOSED` issue を混ぜない。
9. owned scope の Sonar issue が残る間は lane owner が同じ active plan を更新し、implementing skill を再スポーンして修正を差し戻す。
10. Sonar issue gate が解消した後に `<review_cycler>` を `reviewing-implementation` で **1** 回だけ実行する。
11. review が `reroute` を返したら lane owner が同じ active plan を更新し、implementing skill を再スポーンして修正する。
12. review が `pass` を返した後に implementation lane owner (`directing-implementation`) が project root で `python3 scripts/harness/run.py --suite all` を実行する。
13. `all` が失敗した時は lane owner が同じ active plan を更新し、implementing skill を再スポーンして修正する。
14. 差し戻しが修正されたら､再レビューはせず次へ進む｡
15. close 前に completed work から task-local ではない知識だけを distill し、MCP memory bucket (`repo_conventions` または `recurring_pitfalls`) へ追加する。未承認設計、今回だけの TODO、一時回避策は MCP memory に残さない。
16. 必要な `4humans sync` を整理し、backend 構造の変更または追加があった時は `<structure_diagrammer>` を `structure diagram diff skill (`diagramming-structure-diff`)` でスポーンして承認済み review 用差分図を `差分正本適用先` に適用する。その後、処理変更がある時は `diagramming-d2` で `4humans/diagrams/processes/` の relevant `.d2` / `.svg` を更新し、構造変更がある時は `diagramming-d2` で `4humans/diagrams/structures/` の relevant `.d2` / `.svg` を同一変更で更新してから review 用差分図を削除し、plan を `completed/` へ移す。new detail `.d2` を追加する時は `4humans/diagrams/overview-manifest.json` と manifest で紐づく overview `.d2` / `.svg` も同一変更で更新する。
17. タスクがアサインされている場合、タスクのstatusをdoneにする。

## 許可すること
- 各エージェントのスポーン spawn_agentのfork_contextはfalseで呼び出すこと。
- 各エージェントの契約パケットを読む

## Rules

- proposal 前提が不足している時は開始せず `proposing-implementation` へ戻す
- `changes/`、`context_board`、`tasks.md` を live 正本にしない
- implementing skill には raw lint command 群ではなく suite 名で渡し、lint の中身は repo 側の harness と package script を正本にする
- implementation lane owner (`directing-implementation`) が Sonar issue remediation loop、single-pass review、review pass 後の final harness を一元管理する
- MCP memory は bucket (`repo_conventions`, `recurring_pitfalls`) の recall / distill に限定し、仕様や architecture の正本代替にしない
- close 前の MCP memory distill では repo-wide に再利用できる項目だけを残し、task-local な判断や one-off incident を bucket に昇格させない
- review が `pass` でも `4humans sync` と plan 完了前に close とみなさない
- `diagrams/backend/` 正本は proposal 中に書き換えず、承認済み差分を close 時に適用する
- review 用差分図は承認時点の一時 artifact として扱い、正本適用と `4humans` 正本同期後は削除する
- active plan の `4humans Sync` には、必要な `4humans/diagrams/processes/` と `4humans/diagrams/structures/` の更新対象を明記し、new detail `.d2` を追加する時は `4humans/diagrams/overview-manifest.json` と対応 overview `.d2` / `.svg` も必ず列挙する
- skill 権限が曖昧な場合は停止して適切な handoff を選ぶ

## Reference Use

- downstream skill へ handoff する前に `references/directing-implementation.to.<skill>.json` を参照し、渡す情報を揃える。
- `proposing-implementation` から受け取る時は `../proposing-implementation/references/proposing-implementation.to.directing-implementation.json` を入口契約として扱う。
- downstream skill から受け取る時は、各 skill 側の `references/<skill>.to.directing-implementation.json` を返却契約として扱う。
