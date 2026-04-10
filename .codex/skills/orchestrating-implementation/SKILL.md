---
name: orchestrating-implementation
description: AITranslationEngineJp 専用。`workflow.md` の実装レーンを起点から完了まで順番に進め、必要な差し戻し先を決める orchestrator。
---

# Orchestrating Implementation

この skill は `workflow.md` の実装レーンをそのまま運用する入口です。
過去の packet 運用、独自 gate、追加 loop は持ち込まず、段階、成果物、差し戻し先だけを管理します。

## 使う場面

- 新機能実装
- 既存機能拡張
- 詳細設計から実装完了までを 1 本のレーンで進める通常実装

## Required Workflow

1. `docs/exec-plans/templates/impl-plan.md` を使って active plan を作成または更新する。
2. 要求整理として `phase-1-distill` へ handoff し、facts、constraints、gaps、required reading を集める。
3. 詳細設計として `phase-2-design` へ handoff し、active plan の `UI` / `Scenario` / `Logic` を固める。
5. 詳細設計 AI review として `phase-2.5-design-review` を 1 回だけ実行する。
6. `reroute` の時は第2段階へ戻す。人間確認が必要な論点が残る時は active plan に明示して第3段階で固定する。
7. 人間承認を active plan の `承認記録` と `HITL 状態` に固定する。未承認なら次へ進めない。
8. 実装計画として `phase-4-plan` へ handoff し、実装順、並列 task group、依存関係、担当範囲、検証コマンドを短い implementation brief に落とす。
9. 検証設計として `phase-5-test-implementation` へ handoff し、`Scenario` を tests / fixtures / acceptance checks / validation commands へ適用する。
10. 第4段階で固定した並列 task group と依存関係に従い、独立した scope は並列に、依存が残る scope は依存解消後に `phase-6-implement-frontend` または `phase-6-implement-backend` へ handoff する。
11. 実装で設計前提が崩れた時は第2段階へ戻す。実装修正だけで済む時は第6段階に留める。
12. UI確認として `phase-6.5-ui-check` へ handoff し、`chrome-devtools` で主要導線と画面状態を確認する。
13. UI確認が `reroute` の時は、第6段階へ戻すか、設計差分なら第2段階または第3段階へ戻す。
14. 単体テスト作成として `phase-7-unit-test` へ handoff し、責務と主要分岐を証明する unit test を補う。
15. 実装レビューとして `phase-8-review` を 1 回だけ実行する。
16. review が `reroute` の時は、第6段階へ戻すか、設計差分なら第2段階または第3段階へ戻す。
17. review が `pass` の後に `python3 scripts/harness/run.py --suite all` を実行する。
18. review 用差分図を使った時は承認済み差分を正本へ適用し、review 用差分図を削除する。
19. active plan を `completed/` へ移し、完了結果を記録する。

## 管理対象

- 要求整理
- 詳細設計
- 詳細設計 AI review
- 人間承認
- 実装計画
- 検証設計
- 実装と品質通過
- UI確認
- 単体テスト作成
- 実装レビュー
- 完了

## Rules

- `workflow.md` にない独自 phase や独自 gate を追加しない
- 詳細設計が固まる前に第4段階以降へ進めない
- 人間承認が必要な task では、第3段階通過前に実装を開始しない
- 第4段階で固定した並列 task group と依存関係を尊重し、独立した task は並列 handoff する
- `phase-6.5-ui-check` は UI 操作確認と証跡化だけを担当し、恒久修正は第6段階へ戻す
- task-local な設計は active exec-plan の `UI` / `Scenario` / `Logic` に閉じる
- review 用差分図は review の補助であり、正本にしない
- `changes/`、`context_board`、`tasks.md` を live 正本にしない
- 自身で product 実装や詳細なコード調査を抱え込まず、対応する phase skill へ handoff する
- 差し戻し先は `workflow.md` の戻り先に合わせて決める

## Reference Use

- downstream skill へ handoff する前に `references/orchestrating-implementation.to.<skill>.json` を参照し、渡す情報を揃える。
- downstream skill から受け取る時は、各 skill 側の `references/<skill>.to.orchestrating-implementation.json` を返却契約として扱う。
