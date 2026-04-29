# .codex

このディレクトリは Codex workflow の正本です。
Codex は設計 workflow、承認済み scope からの実装、実装後 review、docs 正本化を進めます。

プロダクト仕様と設計判断の正本は `docs/` です。
workflow、skill、agent、handoff 契約の正本は `.codex/` に置きます。
live workflow の説明本文と判断基準の正本はこの `README.md` とします。
`.codex/workflow.md` は補助図であり、live 判断を上書きしません。

## Live Skills

### Main Skills

- 新規実装レーン (`implement-lane`): `skills/implement-lane/SKILL.md`
- 設計壁打ち: `skills/wall-discussion/SKILL.md`
- 設計用文脈整理: `skills/distill/SKILL.md`
- design bundle 進行: `skills/design-bundle/SKILL.md`
- シナリオ候補生成共通 (`scenario-candidate-generation`): `skills/scenario-candidate-generation/SKILL.md`
- シナリオ候補生成 6 観点: `skills/scenario-actor-goal-generation/SKILL.md`、`skills/scenario-lifecycle-generation/SKILL.md`、`skills/scenario-state-transition-generation/SKILL.md`、`skills/scenario-failure-generation/SKILL.md`、`skills/scenario-external-integration-generation/SKILL.md`、`skills/scenario-operation-audit-generation/SKILL.md`
- 設計前調査: `skills/investigate/SKILL.md`
- UI 設計 (`ui-design`): `skills/ui-design/SKILL.md`
- シナリオ設計 (`scenario-design`): `skills/scenario-design/SKILL.md`
- 実装スコープ (`implementation-scope`): `skills/implementation-scope/SKILL.md`
- 実装前文脈整理 (`implementation-distill`): `skills/implementation-distill/SKILL.md`
- 実装時調査 (`implementation-investigate`): `skills/implementation-investigate/SKILL.md`
- product code 実装 (`implement`): `skills/implement/SKILL.md`
- product test 実装 (`tests`): `skills/tests/SKILL.md`
- docs 正本化: `skills/updating-docs/SKILL.md`
- workflow 契約変更: `skills/skill-modification/SKILL.md`
- run-wide report (`work_reporter`): `skills/codex-work-reporting/SKILL.md`
- 実装後 review 観点: `skills/codex-review-behavior/SKILL.md`、`skills/codex-review-規約/SKILL.md`、`skills/codex-review-trust-boundary/SKILL.md`、`skills/codex-review-state-invariant/SKILL.md`

### Support Skills

- 設計向け文脈整理: `skills/distill-design/SKILL.md`
- 調査向け文脈整理: `skills/distill-investigate/SKILL.md`
- 図作成補助: `skills/diagramming/SKILL.md`
- 実装 focused skill: `skills/implement-backend/SKILL.md`、`skills/implement-frontend/SKILL.md`、`skills/implement-mixed/SKILL.md`、`skills/implement-fix-lane/SKILL.md`
- 実装前文脈整理 focused skill: `skills/implementation-distill-implement/SKILL.md`、`skills/implementation-distill-fix/SKILL.md`、`skills/implementation-distill-refactor/SKILL.md`
- 実装時調査 focused skill: `skills/implementation-investigate-reproduce/SKILL.md`、`skills/implementation-investigate-trace/SKILL.md`、`skills/implementation-investigate-observe/SKILL.md`、`skills/implementation-investigate-reobserve/SKILL.md`
- test focused skill: `skills/tests-unit/SKILL.md`、`skills/tests-scenario/SKILL.md`

## Agent / Skill Boundary

- live Codex agent は新規実装レーン conductor (`implement_lane`)、scenario candidate generator 6 体、design artifact agent (`designer`)、文脈圧縮 agent (`distiller`)、設計前調査 agent (`investigator`)、実装前文脈整理 agent (`implementation_distiller`)、実装時調査 agent (`implementation_investigator`)、product code 実装 agent (`implementation_implementer`)、product test 実装 agent (`implementation_tester`)、docs 更新 agent (`docs_updater`)、run report agent (`work_reporter`)、観点別 review agent にする
- `implement_lane` は新規実装と機能拡張の task-local artifact DAG、HITL、handoff、close 条件を管理する。全 close 条件には work report と benchmark evidence を必ず含める
- `scenario_actor_goal_generator`、`scenario_lifecycle_generator`、`scenario_state_transition_generator`、`scenario_failure_generator`、`scenario_external_integration_generator`、`scenario_operation_audit_generator` は、それぞれ 1 viewpoint だけを扱い、scenario candidate artifact を作る
- `designer` は `implement_lane` が揃えた scenario 候補 artifact を統合し、scenario を必須要件の固定点として作り、UI 変更がある時だけ `ui-design` を追加し、human review 後に `implementation-scope` を固定する
- `distiller`、scenario candidate generator 6 体、`designer`、`investigator`、`docs_updater` は context を引き継がず、handoff packet だけで動く
- `implement_lane` は承認済み execution artifact を実行正本にし、`implementation_distiller`、`implementation_investigator`、`implementation_implementer`、`implementation_tester` を context 継承なしで直接 spawn する。final validation 後は観点別 review agent を context 継承なしで並列 spawn し、結果を lossless aggregation に統合する
- agent は代理人であり、職責、職能、ロール、tool policy の owner として扱う。`agents/<agent>.toml` の中で「自分は何者か」と `allowed_write_paths` / `allowed_commands` を明示する
- skill は作業プロトコルであり、担当ロールが成果物を作る時の判断規約、成果物規約、完了規約、停止規約を持つ。手順、標準 pattern、参照タイミング一覧、知識範囲一覧は持たない
- Codex agent の人間可読な実行説明は対応する `skills/*/SKILL.md` に置き、binding と tool policy は `agents/<agent>.toml` に置き、入力、出力、完了、停止の規約は対応する `skills/*/SKILL.md` に置く
- `.agent.md` は使わない

## Format Policy

- agent は人間の代わりに task を実行する担当ロールとして定義する
- agent は自分が何者か、職責、tool policy、入力、出力、停止条件、reroute 先を自分の runtime 定義内に持つ
- skill は手順書ではなく作業プロトコルとして定義する
- skill は遵守すべき外部規約、判断規約、成果物規約、完了規約、停止規約を持つ
- skill には手順、網羅的な例外分岐、参照タイミング一覧、知識範囲一覧を置かない

## 責務境界

- `implement_lane` は新規実装レーンの進行役として artifact DAG、spawn packet、human review、human handoff、close 条件を扱う
- `implement_lane` は run の closeout、停止、reroute 時に `codex-work-reporting` を参照し、最後に必ず `work_history` 記録材料と benchmark evidence を作る
- scenario candidate generator 6 体は固定 viewpoint の scenario 候補だけを作り、採否、統合、最終 scenario matrix は扱わない
- `distiller` は task に関連する事柄と必要資料の判断材料を集める
- `designer` は scenario 候補を統合し、design bundle と implementation-scope の task-local artifact を作る
- `investigator` は必要な場合だけ実画面や観測対象を確認し、観測事実と risk を返す
- `implement_lane` は承認済み execution artifact DAG に従い、実装前整理、実装、test、final validation、観点別 review agent の並列 spawn、lossless aggregation、`implementation_action` 分岐を進める
- `implementation_distiller` は single handoff 1 件から implementation lane 用 context を作る
- `implementation_investigator` は承認済み owned_scope 内で実装時の証跡だけを扱う
- `implementation_implementer` は owned_scope 内の product code だけを変更する
- `implementation_tester` は owned_scope を証明する product test と必要最小限の test support だけを変更する
- `docs_updater` は実装と review の完了が分かった後、human 承認済み scope だけを正本化する
- `work_reporter` は Codex benchmark score と completion evidence から `work_history` の run-wide report を生成する。明示 completion evidence が不足する場合は Codex transcript または chat session file を source_ref 付き evidence として確認する
- `implement_lane` は全 implementation handoff と final validation 完了後、diff から取得した実コードを観点グループ別に score 化し、reviewer result bundle、aggregation trace、primary failure mode、dominant invariant、minimum durable fix boundary を completion evidence に残す
- 観点別 review agent は挙動正しさ、契約・互換性、権限・信頼境界、状態・データ不変条件のいずれか 1 つだけを扱い、修正範囲を命令せず修正判断に必要な情報を返す
- `implement_lane`、`designer`、`distiller`、`investigator`、`docs_updater`、`work_reporter`、review agent は product code と product test を変更しない
- product code は `implementation_implementer` だけが owned_scope 内で変更できる
- product test は `implementation_tester` だけが owned_scope 内で変更できる
- implementation lane は docs 正本、`.codex/` workflow 文書、agent runtime、tool policy を変更しない


## Task Type Lanes

- task run は task type ごとの lane として扱い、各 lane が自分の必須 artifact DAG を持つ
- live lane は `implement_lane` と `fix_lane` にする
- `implement_lane` は新規実装と機能拡張だけを扱う
- `fix_lane` は bug fix、regression、validation failure の恒久修正だけを扱う
- `refactor_lane`、`exploration_test_lane`、`ux_refactor_lane` は placeholder とし、必須 artifact、actor、next agent は未定義のままにする
- 各 lane は task-local artifact DAG を持ち、順序は phase 名ではなく `depends_on` と対象 skill の完了規約で固定する
- agent は lane そのものではなく、artifact を作る実行主体として扱う
- 全 lane の close 条件には work report と benchmark evidence を必須で含める


## Implement Lane Artifact DAG

新規実装レーンの成果物DAGは次を標準形にする。
順序は `depends_on` と対象 skill の完了規約で固定し、phase 名では固定しない。

| artifact_id | owner | depends_on | next_agent |
| --- | --- | --- | --- |
| `task_frame` | `implement_lane` | `[]` | none |
| `context_distill` | `distiller` | `task_frame` | `distiller` |
| `scenario_candidates` | scenario generators | `task_frame`, `context_distill?` | scenario candidate generators |
| `design_bundle` | `designer` | `scenario_candidates` | `designer` |
| `human_design_review` | human | `design_bundle` | human |
| `implementation_scope` | `designer` | `human_design_review` | `designer` |
| `implementation_handoff_packet` | `implement_lane` | `implementation_scope` | none |
| `implementation_execution` | `implement_lane` | `implementation_handoff_packet` | `implementation_distiller`, `implementation_investigator?`, `implementation_implementer`, `implementation_tester` |
| `final_validation` | `implement_lane` | `implementation_execution` | none |
| `pass_review_evidence` | `implement_lane` | `final_validation` | review agents |
| `canonicalization_decision` | `implement_lane` | `pass_review_evidence` | `docs_updater?` |
| `work_report_packet` | `implement_lane` / `work_reporter` | all completed or stopped artifacts | `work_reporter` |

## Exec Plan Folder

- 新規 task は `docs/exec-plans/active/<task-id>/` に folder として作る
- `plan.md` は索引、状態、HITL、validation、closeout だけを書く
- 各 skill の資料は同じ folder の skill 名つき file に分ける
- AI は最初に `plan.md` だけ読み、必要な資料だけ追加で読む
- 完了後は folder ごと `docs/exec-plans/completed/<task-id>/` へ移す

## Docs 正本化

- docs 正本化は実装と review の完了が分かった後に扱う
- docs 正本化は Codex 側だけで扱う
- human 承認済みの artifact だけ `docs_updater` が `updating-docs` を参照して正本へ反映する
- task-local UI 要件契約と scenario は task folder に置く
- UI の細かな visual polish は実装後に人間が実物を確認して直す
- `implementation-scope` は handoff 履歴であり docs 正本へ昇格しない

## 非 live 扱い

- 旧 `design` は `scenario-design`、`ui-design`、`implementation-scope` 中心の design bundle に再整理した
- 旧 flat file 形式の exec-plan は legacy とし、新規 task では使わない
- UI check 専用、log instrumentation agent は live から外した
- Codex 側 `distill` は design / investigate の入口整理だけを扱い、implement / fix / refactor の実装前整理は `.codex/skills/implementation-distill` が扱う
- Codex 側の人間可読な runtime 説明は skill へ集約し、`.codex/agents/*.agent.md` は持たない
- `.codex/workflow.md` は補助図として残し、live workflow の正本にはしない
- 旧 skill / agent の退避物は live workflow に残さない

## Work Plan

- 非自明な変更は `docs/exec-plans/active/<task-id>/` に置く
- 完了後は `docs/exec-plans/completed/<task-id>/` へ移す
- completed plan は履歴として残す
