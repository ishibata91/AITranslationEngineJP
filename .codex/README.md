# .codex

このディレクトリは Codex 作業流れ の正本です。
Codex は設計 作業流れ、承認済み 対象範囲 からの実装、実装後 レビュー、docs 正本化を進めます。

プロダクト仕様と設計判断の正本は `docs/` です。
作業流れ、skill、agent、引き継ぎ 契約の正本は `.codex/` に置きます。
live 作業流れ の説明本文と判断基準の正本はこの `README.md` とします。
`.codex/workflow.md` は補助図であり、live 判断を上書きしません。

## Live Skills

### 主 skill

- 新規実装レーン (`implement-lane`): `skills/implement-lane/SKILL.md`
- 設計壁打ち: `skills/wall-discussion/SKILL.md`
- design bundle 進行: `skills/design-bundle/SKILL.md`
- シナリオ候補生成 6 観点: `skills/scenario-actor-goal-generation/SKILL.md`、`skills/scenario-lifecycle-generation/SKILL.md`、`skills/scenario-state-transition-generation/SKILL.md`、`skills/scenario-failure-generation/SKILL.md`、`skills/scenario-external-integration-generation/SKILL.md`、`skills/scenario-operation-audit-generation/SKILL.md`
- 設計前調査: `skills/investigate/SKILL.md`
- UI 設計 (`ui-design`): `skills/ui-design/SKILL.md`
- シナリオ設計 (`scenario-design`): `skills/scenario-design/SKILL.md`
- 実装スコープ (`implementation-scope`): `skills/implementation-scope/SKILL.md`
- 探索テスト計画 (`exploration-test-planning`): `skills/exploration-test-planning/SKILL.md`
- 探索テストレーン (`exploration-test-lane`): `skills/exploration-test-lane/SKILL.md`
- UX改善レーン (`ux-refactor-lane`): `skills/ux-refactor-lane/SKILL.md`
- 実装時調査 (`implementation-investigate`): `skills/implementation-investigate/SKILL.md`
- 修正レーン (`fix-lane`): `skills/fix-lane/SKILL.md`
- プロダクトコード 実装 重点 skill: `skills/implement-backend/SKILL.md`、`skills/implement-frontend/SKILL.md`、`skills/implement-integration/SKILL.md`
- シナリオテスト 実装 (`tests-scenario`): `skills/tests-scenario/SKILL.md`
- 単体テスト 実装 (`tests-unit`): `skills/tests-unit/SKILL.md`
- docs 正本化: `skills/updating-docs/SKILL.md`
- 作業流れ 契約保守 (`workflow-contract-maintenance`): `skills/workflow-contract-maintenance/SKILL.md`
- run 全体レポート (`work_reporter`): `skills/codex-work-reporting/SKILL.md`
- 実装後 レビュー 観点: `skills/codex-review-behavior/SKILL.md`、`skills/codex-review-contract/SKILL.md`、`skills/codex-review-trust-boundary/SKILL.md`、`skills/codex-review-state-invariant/SKILL.md`、`skills/codex-review-responsibility-boundary/SKILL.md`

### 補助 skill

- 図作成補助: `skills/diagramming/SKILL.md`
- 実装時調査 重点 skill: `skills/implementation-investigate-reproduce/SKILL.md`、`skills/implementation-investigate-trace/SKILL.md`、`skills/implementation-investigate-observe/SKILL.md`、`skills/implementation-investigate-reobserve/SKILL.md`

## Agent / Skill Boundary

- live Codex agent は新規実装レーン 進行役 (`implement_lane`)、修正レーン 進行役 (`fix_lane`)、探索テストレーン 進行役 (`exploration_test_lane`)、UX改善レーン 進行役 (`ux_refactor_lane`)、探索テスト計画 agent (`exploration_test_planner`)、シナリオ候補生成 agent 6 体、設計成果物 agent (`designer`)、調査 agent (`investigator`)、実装時調査 agent (`implementation_investigator`)、プロダクトコード 実装 agent (`implementation_implementer`)、シナリオテスト 実装 agent (`implementation_scenario_tester`)、単体テスト 実装 agent (`implementation_unit_tester`)、docs 更新 agent (`docs_updater`)、run レポート agent (`work_reporter`)、観点別 レビュー agent にする
- `implement_lane` は新規実装と機能拡張の task 内成果物 DAG、HITL、引き継ぎ、close 条件を管理する。全 close 条件には 作業レポート、作業観測根拠、作業計画 folder の `docs/exec-plans/completed/<task-id>/` への移動を必ず含める
- `fix_lane` は人間が確認した不具合、レビュー非通過、検証失敗の task 内成果物 DAG、起動入力、担当 agent 起動、停止、戻し、close 条件を管理する
- `exploration_test_lane` は探索テストの task 内成果物 DAG、起動入力、停止、戻し、close 条件を管理する。探索計画、探索証跡、実装、回帰確認の担当 agent を分ける
- `ux_refactor_lane` は frontend 実装修正を起点にした画面体験改善の task 内成果物 DAG、起動入力、人間UIレビュー、テスト修正、責務レビュー、close 条件を管理する
- `exploration_test_planner` は探索計画だけを作り、観測、ログ確認、画面確認、原因仮説の作成を扱わない
- `scenario_actor_goal_generator`、`scenario_lifecycle_generator`、`scenario_state_transition_generator`、`scenario_failure_generator`、`scenario_external_integration_generator`、`scenario_operation_audit_generator` は、それぞれ 1 観点 だけを扱い、シナリオ 候補成果物 を作る
- `designer` は `implement_lane` が揃えた シナリオ 候補 成果物 を統合し、シナリオ を必須要件の固定点として作る。UI 変更がある時は `ui-design` を独立成果物として作り、`ui-design.md` を揃える。人間レビュー 後に `implementation-scope` を固定する
- シナリオ候補生成 agent 6 体、`designer`、`exploration_test_planner`、`investigator`、`docs_updater` は 文脈 を引き継がず、引き継ぎ入力 だけで動く
- `implement_lane` は承認済み 実行成果物 を実行正本にし、`implementation_investigator`、`implementation_implementer`、`implementation_scenario_tester`、`implementation_unit_tester` を 文脈 継承なしで直接 起動 する。UI がある task では frontend 実装後に人間レビューを挟む。最終検証 後は観点別 レビュー agent を 文脈 継承なしで並列 起動 し、結果を 欠落なし集約 に統合する
- `ux_refactor_lane` は `task 枠` を実行正本にし、`implementation_implementer`、`implementation_scenario_tester`、`implementation_unit_tester`、`review_responsibility_boundary`、`work_reporter` を 文脈 継承なしで直接 起動 する。frontend 実装 agent は人間UIレビュー結果の記録まで維持し、その他の完了済みサブエージェントは完了結果を集約した後に閉じる
- agent は代理人であり、職責、職能、ロール、ツール権限 の 担当者 として扱う。`agents/<agent>.toml` の中で「自分は何者か」と 実行境界 を明示する
- skill は作業プロトコルであり、担当ロールが成果物を作る時の判断規約、成果物規約、完了規約、停止規約を持つ。手順、標準 型、参照タイミング一覧、知識範囲一覧は持たない
- Codex agent の人間可読な実行説明は対応する `skills/*/SKILL.md` に置き、紐づけ と `sandbox_mode` は `agents/<agent>.toml` に置き、入力、出力、完了、停止の規約は対応する `skills/*/SKILL.md` に置く
- サンドボックス外で実行してよい command prefix は `.codex/rules/default.rules` の Codex rules に置く
- `.agent.md` は使わない

## 形式規約

- agent は人間の代わりに task を実行する担当ロールとして定義する
- agent は自分が何者か、職責、実行境界、入力、出力、停止条件、戻し先を自分の 実行定義内に持つ
- skill は手順書ではなく作業プロトコルとして定義する
- skill は遵守すべき外部規約、判断規約、成果物規約、完了規約、停止規約を持つ
- skill には手順、網羅的な例外分岐、参照タイミング一覧、知識範囲一覧を置かない

## 責務境界

- `implement_lane` は新規実装レーンの進行役として 成果物 DAG、起動入力、人間レビュー、人間向け引き継ぎ、close 条件、作業計画 folder の完了移動を扱う
- `implement_lane` は run の 終了処理、停止、戻し 時に `codex-work-reporting` を参照し、最後に必ず `work_history` 記録材料と 作業観測根拠 を作る
- `fix_lane` は人間観測、レビュー非通過、検証失敗、修正前調査を読み、修正実行入力、レビュー通過根拠を管理する。調査、実装、テスト、レビュー、作業レポート本文は担当 agent を起動して委任し、プロダクトコードとプロダクトテストは変更しない
- シナリオ候補生成 agent 6 体は固定 観点 の シナリオ 候補だけを作り、採否、統合、最終 シナリオ表 は扱わない
- `designer` は シナリオ 候補を統合し、シナリオ設計、UI 設計、implementation-scope の task 内成果物 を作る。UI 設計は design bundle 本体へ含めず、独立成果物として扱う
- `exploration_test_lane` は探索計画と探索証跡を読み、バグ一覧、ログ、影響ファイルを集約する。プロダクトコードとプロダクトテストは変更しない
- `ux_refactor_lane` は人間依頼、既存画面根拠、変更許可範囲、禁止範囲を `task 枠` に固定し、frontend 実装、人間UIレビュー、テスト修正、責務レビューを管理する。プロダクトコードとプロダクトテストは変更しない
- `exploration_test_planner` は探索テストの観測対象、探索観点、テストデータ方針、停止条件だけを固定する
- `investigator` は設計前調査、探索テスト証跡、修正前調査のために実画面や観測対象を確認し、観測事実、UI 証跡、ログ、未確認事項を返す。探索テストレーンでは探索証跡だけを担当し、探索範囲を広げる判断をしない。修正レーンでは修正前調査だけを担当し、修正実行入力を作らない
- `implement_lane` は承認済み 実行成果物 DAG に従い、実装時調査、実装、テスト、最終検証、検証証跡を渡した観点別 レビュー agent の並列 起動、`reviewback.<観点>.yaml` 群の欠落なし集約、`implementation_action` 分岐を進める
- `implement_lane` は観点別 レビュー結果を `reviewback.<観点>.yaml` の `must_fix_open` と `max_level` から集約し、behavior、security、responsibility_boundary、その他 の優先度で上位観点の失敗または停止を下位観点の通過で相殺しない
- `implementation_investigator` は承認済み実装範囲 内で実装時の証跡だけを扱う
- `implementation_implementer` は 承認済み実装範囲 内の プロダクトコード だけを変更する
- `implementation_implementer` は backend、frontend、integration のいずれか 1 つの実装 skill を主契約として選ぶ。共通親 skill は置かない
- `implementation_scenario_tester` は 承認済みシナリオ と 承認済み実装範囲 を証明する シナリオテスト と必要最小限の テスト補助 だけを変更する
- `implementation_unit_tester` は 実装済み責務 と 承認済み実装範囲 を証明する 単体テスト と必要最小限の テスト補助 だけを変更する
- `docs_updater` は実装と レビュー の完了が分かった後、human 承認済み 対象範囲 だけを正本化する
- `work_reporter` は 完了根拠、`transcript_refs.json`、レビュー最終状態 YAML、改善ログ、検証結果 から `work_history` の run 全体レポート を生成する。明示 完了根拠 が不足する場合は Codex 会話ログ または chat session file を 根拠参照 付き 根拠 として確認する
- `implement_lane` は全 implementation 引き継ぎ と 最終検証 完了後、diff から取得した実コードを観点グループ別に 評価し、`reviewback.<観点>.yaml`、集約記録、主な失敗種別、主要不変条件、最小恒久修正境界 を 完了根拠 に残す
- `implement_lane` は run 中に見つけた構造問題、作業流れ問題、権限問題、実行問題、人間フィードバック、レビュー由来の改善示唆を `work_history/runs/<run>/workflow-improvement-log.jsonl` へ逐次追記する
- `ux_refactor_lane` は人間UIレビュー 完了後、UI 構造または表示文言の変更で落ちるテストを `implementation_scenario_tester` または `implementation_unit_tester` へ渡し、その後に `review_responsibility_boundary` を必須レビューとして起動する。責務外の実装範囲、secret、外部入力、保存先、ログ出力先を触る差分が必要な場合は停止して人間へ返す
- 観点別 レビュー agent は挙動正しさ、契約・互換性、権限・信頼境界、状態・データ不変条件、責務境界のいずれか 1 つだけを扱い、`reviewback.<観点>.yaml` を作成、追記、解決更新、削除する
- 観点別 レビュー agent は広い ハーネス 再実行を担当せず、`implement_lane` から渡された検証証跡をレビュー入力として扱う
- 観点別 レビュー agent は 失敗 または 停止 の場合も `reviewback.<観点>.yaml` に結果、根拠、未解決指摘を記録する
- `reviewback.<観点>.yaml` はゲート判断用レビュー成果物とし、work_history 側に観点別の非通過 YAML は作らない
- `workflow-improvement-log.jsonl` は作業流れ改善用の run 内観測ログとし、ゲート判断には使わない
- `implement_lane`、`fix_lane`、`exploration_test_lane`、`ux_refactor_lane`、`exploration_test_planner`、`designer`、`investigator`、`docs_updater`、`work_reporter`、レビュー agent は プロダクトコード と プロダクトテスト を変更しない
- プロダクトコード は `implementation_implementer` だけが 承認済み実装範囲 内で変更できる
- シナリオテスト は `implementation_scenario_tester` だけが 承認済み実装範囲 内で変更できる
- 単体テスト は `implementation_unit_tester` だけが 承認済み実装範囲 内で変更できる
- implementation レーン は docs 正本、`.codex/` 作業流れ 文書、agent 実行定義、ツール権限 を変更しない


## task 種別レーン

- task run は task type ごとの レーン として扱い、各 レーン が自分の必須 成果物 DAG を持つ
- live レーン は `implement_lane`、`fix_lane`、`exploration_test_lane`、`ux_refactor_lane` にする
- `implement_lane` は新規実装と機能拡張だけを扱う
- `fix_lane` は人間が確認した不具合、レビュー非通過、検証失敗の恒久修正だけを扱う
- `exploration_test_lane` は探索計画、テストデータ、探索証跡、バグ一覧、ログ、影響ファイル、実装証跡、回帰テスト証跡を扱う
- `ux_refactor_lane` は frontend 実装修正を起点にした画面体験改善だけを扱う
- `refactor_lane` は placeholder とし、必須 成果物、実行者、next agent は未定義のままにする
- 各 レーン は task 内成果物 DAG を持ち、順序は phase 名ではなく `依存対象` と対象 skill の完了規約で固定する
- agent は レーン そのものではなく、成果物 を作る実行主体として扱う
- 全 レーン の close 条件には 作業レポート、作業観測根拠、作業計画 folder の `docs/exec-plans/completed/<task-id>/` への移動を必須で含める


## 実装レーン成果物DAG

新規実装レーンの成果物DAGは次を標準形にする。
順序は `依存対象` と対象 skill の完了規約で固定し、phase 名では固定しない。

| 成果物ID | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- |
| `task 枠` | `implement_lane` | `[]` | なし |
| `scenario_candidates` | シナリオ 生成 agent | `task 枠` | シナリオ候補 生成 agent |
| `シナリオ設計` | `designer` | `scenario_candidates` | `designer` |
| `UI設計` | `designer` | `シナリオ設計` | `designer` |
| `人間設計レビュー` | human | `シナリオ設計`, `UI設計?` | human |
| `実装範囲` | `designer` | `人間設計レビュー` | `designer` |
| `実装引き継ぎ入力` | `implement_lane` | `実装範囲` | なし |
| `frontend 実装` | `implementation_implementer` / `implement-frontend` | `実装引き継ぎ入力` | `implementation_implementer` |
| `frontend 実装後人間レビュー` | human | `frontend 実装` | human |
| `backend 実装` | `implementation_implementer` / `implement-backend` | `実装引き継ぎ入力`, `frontend 実装後人間レビュー?` | `implementation_implementer` |
| `統合境界実装` | `implementation_implementer` / `implement-integration` | `backend 実装`, `frontend 実装後人間レビュー?` | `implementation_implementer` |
| `シナリオテスト` | `implementation_scenario_tester` | `backend 実装?`, `frontend 実装後人間レビュー?`, `統合境界実装?` | `implementation_scenario_tester` |
| `単体テスト` | `implementation_unit_tester` | `backend 実装?`, `frontend 実装後人間レビュー?`, `統合境界実装?` | `implementation_unit_tester` |
| `最終検証` | `implement_lane` | `backend 実装?`, `frontend 実装後人間レビュー?`, `統合境界実装?`, `シナリオテスト?`, `単体テスト?` | なし |
| `レビュー通過根拠` | `implement_lane` | `最終検証` | レビュー agents |
| `正本化判断` | `implement_lane` | `レビュー通過根拠` | `docs_updater?` |
| `詳細仕様正本反映` | `docs_updater` | `正本化判断` | `docs_updater?` |
| `作業レポート入力` | `implement_lane` / `work_reporter` | 全完了または停止済み 成果物 | `work_reporter` |
| `作業計画完了移動` | `implement_lane` | `作業レポート入力` | なし |

## 修正レーン成果物DAG

修正レーンの成果物DAGは次を標準形にする。
順序は `依存対象` と対象 skill の完了規約で固定し、phase 名では固定しない。

| 成果物ID | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- |
| `人間観測記録` | `fix_lane` | `task 枠` | なし |
| `修正前調査` | `investigator` | `人間観測記録` | `investigator` |
| `修正実行入力` | `fix_lane` | `人間観測記録`, `修正前調査` | なし |
| `実装証跡` | `implementation_implementer` / `implement-backend` または `implement-frontend` または `implement-integration` | `修正実行入力` | `implementation_implementer` |
| `回帰テスト証跡` | `implementation_scenario_tester` または `implementation_unit_tester` | `実装証跡` | `implementation_scenario_tester` または `implementation_unit_tester` |
| `レビュー通過根拠` | `fix_lane` | `人間観測記録`, `修正前調査`, `修正実行入力`, `実装証跡`, `回帰テスト証跡?` | `review_behavior`, `review_contract`, `review_trust_boundary`, `review_state_invariant`, `review_responsibility_boundary` |
| `作業レポート入力` | `fix_lane` / `work_reporter` | 全完了または停止済み 成果物, `レビュー通過根拠?` | `work_reporter` |
| `作業計画完了移動` | `fix_lane` | `作業レポート入力` | なし |

## 探索テストレーン成果物DAG

探索テストレーンの成果物DAGは次を標準形にする。
順序は `依存対象` と対象 skill の完了規約で固定し、phase 名では固定しない。

| 成果物ID | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- |
| `探索計画` | `exploration_test_planner` | `task 枠` | `exploration_test_planner` |
| `テストデータ` | `exploration_test_lane` | `探索計画` | なし |
| `探索証跡` | `investigator` | `探索計画`, `テストデータ` | `investigator` |
| `バグ一覧とログ、影響ファイル` | `exploration_test_lane` | `探索証跡` | なし |
| `実装証跡` | `implementation_implementer` | `バグ一覧とログ、影響ファイル` | `implementation_implementer` |
| `回帰テスト証跡` | `implementation_scenario_tester` または `implementation_unit_tester` | `実装証跡` | `implementation_scenario_tester` または `implementation_unit_tester` |
| `レビュー通過根拠` | `exploration_test_lane` | `探索計画`, `探索証跡`, `バグ一覧とログ、影響ファイル`, `実装証跡?`, `回帰テスト証跡?` | `review_behavior`, `review_contract`, `review_trust_boundary`, `review_state_invariant`, `review_responsibility_boundary` |
| `作業レポート入力` | `exploration_test_lane` / `work_reporter` | 全完了または停止済み 成果物, `レビュー通過根拠?` | `work_reporter` |

## UX改善レーン成果物DAG

UX改善レーンの成果物DAGは次を標準形にする。
順序は `依存対象` と対象 skill の完了規約で固定し、phase 名では固定しない。

| 成果物ID | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- |
| `task 枠` | `ux_refactor_lane` | `[]` | なし |
| `frontend 実装` | `implementation_implementer` / `implement-frontend` | `task 枠` | `implementation_implementer` |
| `人間UIレビュー` | human | `frontend 実装` | human |
| `テスト修正証跡` | `implementation_scenario_tester` または `implementation_unit_tester` | `frontend 実装`, `人間UIレビュー` | `implementation_scenario_tester` または `implementation_unit_tester` |
| `レビュー通過根拠` | `ux_refactor_lane` | `frontend 実装`, `人間UIレビュー`, `テスト修正証跡?` | `review_responsibility_boundary` |
| `作業レポート入力` | `ux_refactor_lane` / `work_reporter` | 全完了または停止済み 成果物, `レビュー通過根拠?` | `work_reporter` |
| `作業計画完了移動` | `ux_refactor_lane` | `作業レポート入力` | なし |

## 実行計画 folder

- 新規 task は `docs/exec-plans/active/<task-id>/` に folder として作る
- `plan.md` は索引、状態、HITL、検証、終了処理 だけを書く
- 各 skill の資料は同じ folder の skill 名つき file に分ける
- AI は最初に `plan.md` だけ読み、必要な資料だけ追加で読む
- close 時は folder ごと `docs/exec-plans/completed/<task-id>/` へ移す

## Docs 正本化

- docs 正本化は実装と レビュー の完了が分かった後に扱う
- docs 正本化は Codex 側だけで扱う
- human 承認済みの 成果物 だけ `docs_updater` が `updating-docs` を参照して正本へ反映する
- task 内 UI 要件契約、agent-browser 確認結果、シナリオ は task folder に置く
- UI の確認は、承認済み UI 要件契約と実画面確認結果で扱う
- UI の細かな visual polish は実装後の実物確認で差分を扱う
- `implementation-scope` は 引き継ぎ 履歴であり docs 正本へ昇格しない
- `detail-specs` は 上位シナリオ 単位の詳細仕様正本 とし、`scenario-design`、`ui-design`、実装結果、レビュー結果から human 承認済みの恒久仕様だけを製本する

## 非 live 扱い

- 旧 `design` は `scenario-design`、独立した `ui-design`、`implementation-scope` に再整理した
- 旧 flat file 形式の exec-plan は legacy とし、新規 task では使わない
- UI check 専用、log instrumentation agent は live から外した
- 作業前の影響範囲、実行計画、検証方法の確認は `AGENTS.md` の入口規約に集約する
- Codex 側の人間可読な 実行定義 説明は skill へ集約し、`.codex/agents/*.agent.md` は持たない
- `.codex/workflow.md` は補助図として残し、live 作業流れ の正本にはしない
- 旧 skill / agent の退避物は live 作業流れ に残さない

## 作業計画

- 非自明な変更は `docs/exec-plans/active/<task-id>/` に置く
- 完了後は `docs/exec-plans/completed/<task-id>/` へ移す
- completed plan は履歴として残す
