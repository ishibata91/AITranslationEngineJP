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
- シナリオ候補生成共通 (`scenario-candidate-generation`): `skills/scenario-candidate-generation/SKILL.md`
- シナリオ候補生成 6 観点: `skills/scenario-actor-goal-generation/SKILL.md`、`skills/scenario-lifecycle-generation/SKILL.md`、`skills/scenario-state-transition-generation/SKILL.md`、`skills/scenario-failure-generation/SKILL.md`、`skills/scenario-external-integration-generation/SKILL.md`、`skills/scenario-operation-audit-generation/SKILL.md`
- 設計前調査: `skills/investigate/SKILL.md`
- UI 設計 (`ui-design`): `skills/ui-design/SKILL.md`
- シナリオ設計 (`scenario-design`): `skills/scenario-design/SKILL.md`
- 実装スコープ (`implementation-scope`): `skills/implementation-scope/SKILL.md`
- 実装時調査 (`implementation-investigate`): `skills/implementation-investigate/SKILL.md`
- プロダクトコード 実装 (`implement`): `skills/implement/SKILL.md`
- プロダクトテスト 実装 (`tests`): `skills/tests/SKILL.md`
- docs 正本化: `skills/updating-docs/SKILL.md`
- 作業流れ 契約変更: `skills/skill-modification/SKILL.md`
- run 全体レポート (`work_reporter`): `skills/codex-work-reporting/SKILL.md`
- 実装後 レビュー 観点: `skills/codex-review-behavior/SKILL.md`、`skills/codex-review-contract/SKILL.md`、`skills/codex-review-trust-boundary/SKILL.md`、`skills/codex-review-state-invariant/SKILL.md`

### 補助 skill

- 図作成補助: `skills/diagramming/SKILL.md`
- 実装 重点 skill: `skills/implement-backend/SKILL.md`、`skills/implement-frontend/SKILL.md`、`skills/implement-mixed/SKILL.md`、`skills/implement-fix-lane/SKILL.md`
- 実装時調査 重点 skill: `skills/implementation-investigate-reproduce/SKILL.md`、`skills/implementation-investigate-trace/SKILL.md`、`skills/implementation-investigate-observe/SKILL.md`、`skills/implementation-investigate-reobserve/SKILL.md`
- test 重点 skill: `skills/tests-unit/SKILL.md`、`skills/tests-scenario/SKILL.md`

## Agent / Skill Boundary

- live Codex agent は新規実装レーン 進行役 (`implement_lane`)、シナリオ候補生成 agent 6 体、設計成果物 agent (`designer`)、設計前調査 agent (`investigator`)、実装時調査 agent (`implementation_investigator`)、プロダクトコード 実装 agent (`implementation_implementer`)、プロダクトテスト 実装 agent (`implementation_tester`)、docs 更新 agent (`docs_updater`)、run レポート agent (`work_reporter`)、観点別 レビュー agent にする
- `implement_lane` は新規実装と機能拡張の task 内成果物 DAG、HITL、引き継ぎ、close 条件を管理する。全 close 条件には 作業レポート と ベンチマーク根拠 を必ず含める
- `scenario_actor_goal_generator`、`scenario_lifecycle_generator`、`scenario_state_transition_generator`、`scenario_failure_generator`、`scenario_external_integration_generator`、`scenario_operation_audit_generator` は、それぞれ 1 観点 だけを扱い、シナリオ 候補成果物 を作る
- `designer` は `implement_lane` が揃えた シナリオ 候補 成果物 を統合し、シナリオ を必須要件の固定点として作り、UI 変更がある時だけ `ui-design` を追加し、人間レビュー 後に `implementation-scope` を固定する
- シナリオ候補生成 agent 6 体、`designer`、`investigator`、`docs_updater` は 文脈 を引き継がず、引き継ぎ入力 だけで動く
- `implement_lane` は承認済み 実行成果物 を実行正本にし、`implementation_investigator`、`implementation_implementer`、`implementation_tester` を 文脈 継承なしで直接 起動 する。最終検証 後は観点別 レビュー agent を 文脈 継承なしで並列 起動 し、結果を 欠落なし集約 に統合する
- agent は代理人であり、職責、職能、ロール、ツール権限 の 担当者 として扱う。`agents/<agent>.toml` の中で「自分は何者か」と 書き込み許可 / 実行許可 を明示する
- skill は作業プロトコルであり、担当ロールが成果物を作る時の判断規約、成果物規約、完了規約、停止規約を持つ。手順、標準 型、参照タイミング一覧、知識範囲一覧は持たない
- Codex agent の人間可読な実行説明は対応する `skills/*/SKILL.md` に置き、紐づけ と ツール権限 は `agents/<agent>.toml` に置き、入力、出力、完了、停止の規約は対応する `skills/*/SKILL.md` に置く
- `.agent.md` は使わない

## 形式規約

- agent は人間の代わりに task を実行する担当ロールとして定義する
- agent は自分が何者か、職責、ツール権限、入力、出力、停止条件、戻し先を自分の 実行定義内に持つ
- skill は手順書ではなく作業プロトコルとして定義する
- skill は遵守すべき外部規約、判断規約、成果物規約、完了規約、停止規約を持つ
- skill には手順、網羅的な例外分岐、参照タイミング一覧、知識範囲一覧を置かない

## 責務境界

- `implement_lane` は新規実装レーンの進行役として 成果物 DAG、起動入力、人間レビュー、人間向け引き継ぎ、close 条件を扱う
- `implement_lane` は run の 終了処理、停止、戻し 時に `codex-work-reporting` を参照し、最後に必ず `work_history` 記録材料と ベンチマーク根拠 を作る
- シナリオ候補生成 agent 6 体は固定 観点 の シナリオ 候補だけを作り、採否、統合、最終 シナリオ表 は扱わない
- `designer` は シナリオ 候補を統合し、design bundle と implementation-scope の task 内成果物 を作る
- `investigator` は必要な場合だけ実画面や観測対象を確認し、観測事実と リスク を返す
- `implement_lane` は承認済み 実行成果物 DAG に従い、実装時調査、実装、test、最終検証、観点別 レビュー agent の並列 起動、欠落なし集約、`implementation_action` 分岐を進める
- `implementation_investigator` は承認済み実装範囲 内で実装時の証跡だけを扱う
- `implementation_implementer` は 承認済み実装範囲 内の プロダクトコード だけを変更する
- `implementation_tester` は 承認済み実装範囲 を証明する プロダクトテスト と必要最小限の test 補助 だけを変更する
- `docs_updater` は実装と レビュー の完了が分かった後、human 承認済み 対象範囲 だけを正本化する
- `work_reporter` は Codex ベンチマーク値 と 完了根拠 から `work_history` の run 全体レポート を生成する。明示 完了根拠 が不足する場合は Codex 会話ログ または chat session file を 根拠参照 付き 根拠 として確認する
- `implement_lane` は全 implementation 引き継ぎ と 最終検証 完了後、diff から取得した実コードを観点グループ別に 評価値 化し、レビュー結果一式、集約記録、主な失敗種別、主要不変条件、最小恒久修正境界 を 完了根拠 に残す
- 観点別 レビュー agent は挙動正しさ、契約・互換性、権限・信頼境界、状態・データ不変条件のいずれか 1 つだけを扱い、修正範囲を命令せず修正判断に必要な情報を返す
- `implement_lane`、`designer`、`investigator`、`docs_updater`、`work_reporter`、レビュー agent は プロダクトコード と プロダクトテスト を変更しない
- プロダクトコード は `implementation_implementer` だけが 承認済み実装範囲 内で変更できる
- プロダクトテスト は `implementation_tester` だけが 承認済み実装範囲 内で変更できる
- implementation レーン は docs 正本、`.codex/` 作業流れ 文書、agent 実行定義、ツール権限 を変更しない


## task 種別レーン

- task run は task type ごとの レーン として扱い、各 レーン が自分の必須 成果物 DAG を持つ
- live レーン は `implement_lane` と `fix_lane` にする
- `implement_lane` は新規実装と機能拡張だけを扱う
- `fix_lane` は bug fix、回帰、検証 失敗 の恒久修正だけを扱う
- `refactor_lane`、`exploration_test_lane`、`ux_refactor_lane` は placeholder とし、必須 成果物、実行者、next agent は未定義のままにする
- 各 レーン は task 内成果物 DAG を持ち、順序は phase 名ではなく `依存対象` と対象 skill の完了規約で固定する
- agent は レーン そのものではなく、成果物 を作る実行主体として扱う
- 全 レーン の close 条件には 作業レポート と ベンチマーク根拠 を必須で含める


## 実装レーン成果物DAG

新規実装レーンの成果物DAGは次を標準形にする。
順序は `依存対象` と対象 skill の完了規約で固定し、phase 名では固定しない。

| 成果物ID | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- |
| `task 枠` | `implement_lane` | `[]` | なし |
| `scenario_candidates` | シナリオ 生成 agent | `task 枠` | シナリオ候補 生成 agent |
| `設計成果物束` | `designer` | `scenario_candidates` | `designer` |
| `人間設計レビュー` | human | `設計成果物束` | human |
| `実装範囲` | `designer` | `人間設計レビュー` | `designer` |
| `実装引き継ぎ入力` | `implement_lane` | `実装範囲` | なし |
| `実装実行` | `implement_lane` | `実装引き継ぎ入力` | `implementation_investigator?`, `implementation_implementer`, `implementation_tester` |
| `最終検証` | `implement_lane` | `実装実行` | なし |
| `レビュー通過根拠` | `implement_lane` | `最終検証` | レビュー agents |
| `正本化判断` | `implement_lane` | `レビュー通過根拠` | `docs_updater?` |
| `作業レポート入力` | `implement_lane` / `work_reporter` | 全完了または停止済み 成果物 | `work_reporter` |

## 実行計画 folder

- 新規 task は `docs/exec-plans/active/<task-id>/` に folder として作る
- `plan.md` は索引、状態、HITL、検証、終了処理 だけを書く
- 各 skill の資料は同じ folder の skill 名つき file に分ける
- AI は最初に `plan.md` だけ読み、必要な資料だけ追加で読む
- 完了後は folder ごと `docs/exec-plans/completed/<task-id>/` へ移す

## Docs 正本化

- docs 正本化は実装と レビュー の完了が分かった後に扱う
- docs 正本化は Codex 側だけで扱う
- human 承認済みの 成果物 だけ `docs_updater` が `updating-docs` を参照して正本へ反映する
- task 内 UI 要件契約と シナリオ は task folder に置く
- UI の細かな visual polish は実装後に人間が実物を確認して直す
- `implementation-scope` は 引き継ぎ 履歴であり docs 正本へ昇格しない

## 非 live 扱い

- 旧 `design` は `scenario-design`、`ui-design`、`implementation-scope` 中心の design bundle に再整理した
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
