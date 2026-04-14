# Work Plan Template

- workflow: work
- status: completed
- lane_owner: skill-modification
- scope: .codex/skills, .codex/README.md, .codex/workflow.md, docs/exec-plans/templates, docs/exec-plans/active
- task_id: 2026-04-14-skill-role-rehydration
- task_catalog_ref: N/A
- parent_phase: N/A

## Request Summary

- role-based 統合は維持する
- 旧 skill directory は復活させない
- 旧 skill の運用知識、手順、contract 粒度は新 skill に残す

## Decision Basis

- 前回の圧縮は live 名称統一には成功したが、専門運用の本文が削られすぎた
- user は asset の継承を要求しており、旧 directory の復活は要求していない
- 本体と refs を分け、判断ルールは `SKILL.md`、詳細手順と mode 別 contract は `references/` に置く方針が確定している

## Task Mode

- `task_mode`: implement
- `goal`: role-based skill を維持したまま、旧 specialized skill の知識を新 skill 側へ再配置する
- `constraints`: 旧 skill directory は復活させない。`work-plan.md` を唯一の live template とする。旧名は対応表だけに残す
- `close_conditions`: 新 skill の本文と refs だけで旧運用知識の大半を追える。workflow docs が新名中心で読める。mode 別 contract が残る

## Facts

- new skill の `SKILL.md` を共通判断ルール中心に書き直した
- `references/mode-guides/` を新設し、旧 specialized skill の濃い手順を各 role 配下へ再配置した
- `references/contracts/` を新設し、mode 別 handoff / return contract を新 skill 名のまま保持した
- `.codex/README.md` と `.codex/workflow.md` は新名中心の本文を維持しつつ、末尾に旧名対応表を追加した

## Functional Requirements

- `summary`: 新 skill を knowledge-rich に戻す
- `in_scope`: `.codex/skills/*/SKILL.md`, `references/mode-guides`, mode 別 contract, `.codex/README.md`, `.codex/workflow.md`, active plan README
- `non_functional_requirements`: old names を live flow に戻さない。新旧対応は文書末尾の対応表だけにする
- `out_of_scope`: product code, docs 正本の恒久仕様変更, 旧 skill directory の復活, `impl-plan.md` と `fix-plan.md` の復元
- `open_questions`: なし
- `required_reading`: `docs/exec-plans/completed/2026-04-10-phase-2-skill-split.md`, `docs/exec-plans/completed/2026-04-12-fix-lane-reproduce-issues.md`, `docs/exec-plans/completed/2026-04-14-skill-role-compression.md`

## Artifacts

- `ui_artifact_path`: N/A
- `scenario_artifact_path`: N/A
- `review_diff_diagrams`: N/A
- `source_diagram_targets`: N/A

## Work Brief

- `implementation_target`: backend
- `accepted_scope`: skill docs, refs, contracts, workflow docs
- `parallel_task_groups`: none
- `tasks`: rehydrate skill bodies; add mode guides; split contracts by mode; add legacy-name mapping tables; sync plan docs
- `validation_commands`: `find .codex/skills -name '*.json' -print0 | xargs -0 -n1 python3 -m json.tool >/dev/null`, `python3 scripts/harness/run.py --suite structure`, `python3 scripts/harness/run.py --suite all`

## Investigation

- `reproduction_status`: N/A
- `trace_hypotheses`: N/A
- `observation_points`: N/A
- `residual_risks`: execution harness の backend failure は skill 再編起因ではなく既存 product failure のまま残る

## Acceptance Checks

- new skill 本体に共通判断ルールが残っている
- mode 別 detailed guide が `references/mode-guides/` にある
- mode 別 contract が new skill 名で読める
- `.codex/README.md` と `.codex/workflow.md` に旧名対応表がある

## Required Evidence

- JSON parse success
- structure harness result
- all harness result with failure attribution if any

## HITL Status

- `functional_or_design_hitl`: 不要
- `approval_record`: user approved consolidation-with-asset-retention on 2026-04-14

## Closeout Notes

- D2 記法は Context7 で `shape: class` と `shape: sql_table` を再確認した上で `diagramming` guide に反映した
- quick overview contract は維持し、mode 別 contract を `references/contracts/` に追加した

## Outcome

- `find .codex/skills -name '*.json' -print0 | xargs -0 -n1 python3 -m json.tool >/dev/null` が通過した
- `python3 scripts/harness/run.py --suite structure` が通過した
- `python3 scripts/harness/run.py --suite all` は前回と同じ backend build failure により失敗した
- 失敗点は `internal/usecase/master_dictionary_usecase.go` と `internal/service/master_dictionary_service_test.go` の未解消シンボル不整合である
- frontend lint / test と `sonar-scanner` は通過した
