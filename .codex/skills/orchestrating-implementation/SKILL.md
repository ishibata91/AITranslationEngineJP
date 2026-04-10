---
name: orchestrating-implementation
description: AITranslationEngineJp 専用。`workflow.md` の実装レーンを起点から完了まで順番に進め、必要な差し戻し先を決める orchestrator。
---

# Orchestrating Implementation

この skill は `workflow.md` の実装レーンをそのまま運用する入口です。
過去の packet 運用、独自 gate、追加 loop は持ち込まず、段階、成果物、差し戻し先だけを管理します。
この orchestrator 自身は product 実装や詳細調査を担当せず、各 phase skill を `fork_context: false` で起動するサブエージェントへ必要最小限の handoff 情報だけを渡す配線役として振る舞います。

## 使う場面

- 新機能実装
- 既存機能拡張
- 詳細設計から実装完了までを 1 本のレーンで進める通常実装

## Required Workflow

1. `docs/exec-plans/templates/impl-plan.md` を使って active plan を作成または更新する。
2. 要求整理として `phase-1-distill` へ handoff し、facts、constraints、gaps、required reading を集める。
3. 以下第2段階は3つのエージェントを立ち上げ、並列で実行する。
3. 第2段階の UI モック作成として `phase-2-ui` へ handoff し、active plan の `UI モック` section と、主要なページの動きまである程度再現した page mock working copy を固める。
4. 第2段階の Scenario テスト一覧作成として `phase-2-scenario` へ handoff し、active plan の `Scenario テスト一覧` section とシナリオ一覧 working copy を固める。
5. 第2段階の Logic 実装計画作成として `phase-2-logic` へ handoff し、active plan の `実装計画` を implementation brief として固める。
6. 詳細設計 AI review として `phase-2.5-design-review` を 1 回だけ実行する。
7. `reroute` の時は第2段階へ戻す。人間確認が必要な論点が残る時は active plan に明示して第3段階で固定する。
8. 人間承認を active plan の `承認記録` と `HITL 状態` に固定する。未承認なら次へ進めない。
9. phase3 の human review で review-back が返った時は、指摘を `UI モック作成 (phase-2-ui)`、`Scenario テスト一覧作成 (phase-2-scenario)`、`Logic 実装計画作成 (phase-2-logic)` のどこへ返すか切り分け、該当サブエージェントを再起動して task-local artifact と active plan の参照情報を更新させる。
10. review-back の反映後は `phase-2.5-design-review` を再実行し、`pass` を確認してから第3段階へ戻す。
11. 検証設計として `phase-5-test-implementation` へ handoff し、Scenario artifact を tests / fixtures / acceptance checks / validation commands へ適用する。
12. 第2段階で固定した並列 task group と依存関係に従い、独立した scope は並列に、依存が残る scope は依存解消後に `phase-6-implement-frontend` または `phase-6-implement-backend` へ handoff する。
13. 実装で設計前提が崩れた時は第2段階へ戻す。実装修正だけで済む時は第6段階に留める。
14. UI確認として `phase-6.5-ui-check` へ handoff し、`chrome-devtools` で主要導線と画面状態を確認する。
15. UI確認が `reroute` の時は、第6段階へ戻すか、設計差分なら第2段階または第3段階へ戻す。
16. 単体テスト作成として `phase-7-unit-test` へ handoff し、責務と主要分岐を証明する unit test を補う。
17. 実装レビューとして `phase-8-review` を 1 回だけ実行する。
18. review が `reroute` の時は、第6段階へ戻すか、設計差分なら第2段階または第3段階へ戻す。
19. review が `pass` の後に `python3 scripts/harness/run.py --suite all` を実行する。
20. review 用差分図を使った時は承認済み差分を正本へ適用し、review 用差分図を削除する。
21. 完了前に、第2段階で作った UI モック working copy を `docs/mocks/<page-id>/index.html` へ、Scenario artifact working copy を `docs/scenario-tests/<topic-id>.md` へ移す。
22. active plan を `completed/` へ移し、最終正本 path と完了結果を記録する。

## 管理対象

- 要求整理
- 詳細設計
- 詳細設計 AI review
- 人間承認
- 検証設計
- 実装と品質通過
- UI確認
- 単体テスト作成
- 実装レビュー
- 完了

## Rules

- `workflow.md` にない独自 phase や独自 gate を追加しない
- 詳細設計 artifact と `実装計画` が固まる前に第5段階以降へ進めない
- 不要になったサブエージェントは逐次で閉じること。
- 人間承認が必要な task では、第3段階通過前に実装を開始しない
- phase3 の review-back は orchestrator 自身で抱え込まず、`phase-2-ui`、`phase-2-scenario`、`phase-2-logic` の該当サブエージェントへ返して修正させる
- phase3 の review-back を反映した後は `phase-2.5-design-review` を通し直し、第3段階へ戻す
- 完了時は task-local working copy を `docs/mocks/` と `docs/scenario-tests/` の最終正本へ移してから plan close を行う
- 第2段階の `phase-2-logic` で固定した並列 task group と依存関係を尊重し、独立した task は並列 handoff する
- downstream skill の起動は、`fork_context: false` を明示したサブエージェント呼び出しに限定し、active plan と handoff contract にある必要最小限の情報だけを渡す
- `phase-6.5-ui-check` は UI 操作確認と証跡化だけを担当し、恒久修正は第6段階へ戻す
- 第2段階の成果物は active exec-plan の `UI モック` / `Scenario テスト一覧` / `実装計画` と、対応する別 artifact に閉じる
- review 用差分図は review の補助であり、正本にしない
- `changes/`、`context_board`、`tasks.md` を live 正本にしない
- 自身で product 実装や詳細なコード調査を抱え込まず、対応する phase skill へ handoff する
- 差し戻し先は `workflow.md` の戻り先に合わせて決める
 
## Handoff Agents

- `ctx_loader` `phase-1-distill`
- `task_designer` `phase-2-ui`
- `test_architect` `phase-2-scenario`
- `workplan_builder` `phase-2-logic`
- `review_cycler` `phase-2.5-design-review`
- `test_architect` `phase-5-test-implementation`
- `implementer` `phase-6-implement-frontend`
- `implementer` `phase-6-implement-backend`
- `ui_checker` `phase-6.5-ui-check`
- `test_architect` `phase-7-unit-test`
- `review_cycler` `phase-8-review`
- `structure_diagrammer` `diagramming-structure-diff`

## Reference Use

- downstream skill へ handoff する前に `references/orchestrating-implementation.to.<skill>.json` を参照し、渡す情報を揃える。
- downstream skill から受け取る時は、各 skill 側の `references/<skill>.to.orchestrating-implementation.json` を返却契約として扱う。
