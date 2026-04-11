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
3. 第1.5段階の機能要件固定として `phase-1.5-functional-requirements` へ handoff し、active plan の `機能要件` section に summary、in-scope、out-of-scope、open questions、required reading を固める。
4. 第1.6段階の UI モック作成として `phase-2-ui` へ handoff し、active plan の `UI モック` section と、主要なページの動きまである程度再現した page mock working copy を固める。
5. 第1.7段階の前段 HITL として `機能要件` と `UI モック` を人間承認へ回し、active plan の `機能要件 HITL 状態` と `機能要件 承認記録` に固定する。未承認なら次へ進めない。
6. 前段 HITL の review-back が返った時は、指摘を `機能要件固定 (phase-1.5-functional-requirements)` または `UI モック作成 (phase-2-ui)` へ切り分け、該当サブエージェントを再起動して active plan の参照情報を更新させる。
7. 以下第2段階は3つのエージェントを立ち上げ、並列で実行する。
8. 第2段階の Scenario テスト一覧作成として `phase-2-scenario` へ handoff し、active plan の `Scenario テスト一覧` section とシナリオ一覧 working copy を固める。
9. 第2段階の Logic 実装計画作成として `phase-2-logic` へ handoff し、active plan の `実装計画` を implementation brief として固める。
10. 第2段階の review 用構造差分図作成として `structure_diagrammer` で `diagramming-structure-diff` を `proposal_diff` mode で handoff し、正本のコンポーネント図があるかどうかを判断させ、ある時は更新対象を、ない時は新規作成対象を決めたうえで active plan 配下へ review 用差分 `.d2` / `.svg` を固める。
11. 詳細設計 AI review として `phase-2.5-design-review` を 1 回だけ実行する。
12. `reroute` の時は第2段階へ戻す。差し戻し理由が機能要件の取りこぼしや前段 UI 合意の崩れである時は第1.5段階または第1.6段階へ戻し、前段 HITL をやり直す。人間確認が必要な論点が残る時は active plan に明示して後段 HITL で固定する。
13. 後段の人間承認を active plan の `詳細設計 承認記録` と `詳細設計 HITL 状態` に固定する。未承認なら次へ進めない。
14. 後段 HITL の human review で review-back が返った時は、指摘を `機能要件固定 (phase-1.5-functional-requirements)`、`UI モック作成 (phase-2-ui)`、`Scenario テスト一覧作成 (phase-2-scenario)`、`Logic 実装計画作成 (phase-2-logic)`、`review 用構造差分図作成 (diagramming-structure-diff)` のどこへ返すか切り分け、該当サブエージェントを再起動して task-local artifact と active plan の参照情報を更新させる。AIレビューは必要なし。
15. 検証設計として `phase-5-test-implementation` へ handoff し、Scenario artifact を tests / fixtures / acceptance checks / validation commands へ適用する。
16. 第2段階で固定した並列 task group と依存関係に従い、独立した scope は並列に、依存が残る scope は依存解消後に `phase-6-implement-frontend` または `phase-6-implement-backend` へ handoff する。handoff には active exec-plan だけでなく、承認済み `UI モック` artifact path、`Scenario テスト一覧` artifact path、対象 `task_id`、承認済み `required_reading` を必ず含める。
17. 実装で設計前提が崩れた時は第2段階へ戻す。前段要件の崩れが見つかった時は第1.5段階または第1.6段階へ戻す。実装修正だけで済む時は第6段階に留める。
18. UI確認として `phase-6.5-ui-check` へ handoff し、`chrome-devtools` で主要導線、画面状態、承認済み HTML モック artifact に対する視覚構造差分を確認する。handoff には承認済み design bundle と review 用差分図 path を必ず含める。
19. UI確認が `reroute` の時は、第6段階へ戻すか、設計差分なら第1.5段階、第1.6段階、第2段階または後段 HITL へ戻す。
20. 単体テスト作成として `phase-7-unit-test` へ handoff し、承認済み `task_id` と design artifact に紐づく責務と主要分岐を証明する unit test を補う。
21. 実装レビューとして `phase-8-review` を 1 回だけ実行する。review には承認済み design bundle と review 用差分図 path を渡し、詳細設計との照合を省略しない。
22. review が `reroute` の時は、第6段階へ戻すか、設計差分なら第1.5段階、第1.6段階、第2段階または後段 HITL へ戻す。
23. review が `pass` の後に `python3 scripts/harness/run.py --suite all` を実行する。
24. review 用差分図を使った時は承認済み差分を `structure_diagrammer` で `diagramming-structure-diff` の `apply_to_source` mode へ handoff して正本へ適用し、review 用差分図を削除する。
25. 完了前に、第1.6段階で作った UI モック working copy を `docs/mocks/<page-id>/index.html` へ、第2段階の Scenario artifact working copy を `docs/scenario-tests/<topic-id>.md` へ移す。
26. active plan を `completed/` へ移し、最終正本 path と完了結果を記録する。

## 管理対象

- 要求整理
- 機能要件固定
- 前段 HITL
- 詳細設計
- 詳細設計 AI review
- 後段 HITL
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
- 人間承認が必要な task では、前段 HITL と後段 HITL の両方を通過する前に実装を開始しない
- 前段 HITL の review-back は orchestrator 自身で抱え込まず、`phase-1.5-functional-requirements` または `phase-2-ui` の該当サブエージェントへ返して修正させる
- 後段 HITL の review-back は orchestrator 自身で抱え込まず、`phase-1.5-functional-requirements`、`phase-2-ui`、`phase-2-scenario`、`phase-2-logic`、`diagramming-structure-diff` の該当サブエージェントへ返して修正させる
- 後段 HITL の review-back を反映した後は `phase-2.5-design-review` を通し直し、後段 HITL へ戻す
- 完了時は task-local working copy を `docs/mocks/` と `docs/scenario-tests/` の最終正本へ移してから plan close を行う
- 第2段階の `phase-2-logic` で固定した並列 task group と依存関係を尊重し、独立した task は並列 handoff する
- downstream skill の起動は、`fork_context: false` を明示したサブエージェント呼び出しに限定し、active plan と handoff contract にある必要最小限の情報だけを渡す
- phase-5 以降へ渡す execution handoff では、承認済み `ui_artifact_path`、`scenario_artifact_path`、`implementation_task_ids`、`implementation_required_reading` を design bundle として必須にする
- `phase-6.5-ui-check` と `phase-8-review` へ渡す時は、review 用差分図がない task でも `review_diff_diagrams` を空配列で明示する
- `phase-6.5-ui-check` は UI 操作確認、承認済み HTML モック artifact に対する視覚構造差分の確認、証跡化だけを担当し、恒久修正は第6段階へ戻す
- 前段の成果物は active exec-plan の `機能要件` / `UI モック` と、対応する別 artifact に閉じる
- 第2段階の成果物は active exec-plan の `Scenario テスト一覧` / `実装計画` / `review 用差分図` / `差分正本適用先` と、対応する別 artifact に閉じる
- review 用差分図は review の補助であり、正本にしない
- `changes/`、`context_board`、`tasks.md` を live 正本にしない
- 自身で product 実装や詳細なコード調査を抱え込まず、対応する phase skill へ handoff する
- 差し戻し先は `workflow.md` の戻り先に合わせて決める
 
## Handoff Agents

- `ctx_loader` `phase-1-distill`
- `task_designer` `phase-1.5-functional-requirements`
- `task_designer` `phase-2-ui`
- `test_architect` `phase-2-scenario`
- `workplan_builder` `phase-2-logic`
- `review_cycler` `phase-2.5-design-review`
- `implementer` `phase-5-test-implementation`
- `implementer` `phase-6-implement-frontend`
- `implementer` `phase-6-implement-backend`
- `ui_checker` `phase-6.5-ui-check`
- `test_architect` `phase-7-unit-test`
- `review_cycler` `phase-8-review`
- `structure_diagrammer` `diagramming-structure-diff`

## Reference Use

- downstream skill へ handoff する前に `references/orchestrating-implementation.to.<skill>.json` を参照し、渡す情報を揃える。
- downstream skill から受け取る時は、各 skill 側の `references/<skill>.to.orchestrating-implementation.json` を返却契約として扱う。
