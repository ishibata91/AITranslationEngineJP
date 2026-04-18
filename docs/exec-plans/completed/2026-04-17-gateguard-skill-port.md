# Work Plan Template

- workflow: work
- status: completed
- lane_owner: skill-modification
- scope: .codex/skills/gateguard, .codex/README.md, .codex/workflow.md
- task_id: 2026-04-17-gateguard-skill-port
- task_catalog_ref: N/A
- parent_phase: N/A

## Request Summary

- `../everything-claude-code` の `gateguard` skill を、この repo 向けに日本語化して移植する。
- まずは hook 実装ではなく、編集前の事実確認 gate として価値を見極められる形にする。

## Decision Basis

- `.codex/README.md` は live workflow の入口を `orchestrate` に限定している。
- `orchestrate/SKILL.md` は別 skill を増やさない原則を持つため、`gateguard` は orchestrate handoff 対象にしない。
- `skill-modification` は `.codex/skills/` 配下の skill 追加と関連 workflow docs の同期を許可している。
- 移植元 `gateguard` は PreToolUse hook 前提だが、この repo ではまず direct-use の read-only guard として試す。

## Task Mode

- `task_mode`: implement
- `goal`: repo 固有制約を反映した日本語版 `gateguard` skill を追加する
- `constraints`: product code は変更しない。hook / package installer は持ち込まない。ファイル操作は MCP 経由に限定する
- `close_conditions`: skill 本体、agent 指定、permissions が揃い、workflow docs が direct-use 補助 gate としての位置づけを説明する

## Facts

- 移植元は `/Users/iorishibata/Repositories/everything-claude-code/skills/gateguard/SKILL.md` である。
- 既存 `.codex/skills/gateguard` は存在しなかった。
- 実装前の `python3 scripts/harness/run.py --suite structure` は通過した。
- `skill-creator` は MCP の許可ディレクトリ外で読めなかったため、repo-local `skill-modification` 契約を正本にした。

## Functional Requirements

- `summary`: 編集、ファイル作成、破壊的コマンドの前に、具体的な事実確認を要求する skill を追加する
- `in_scope`: `.codex/skills/gateguard/SKILL.md`, `agents/openai.yaml`, `references/permissions.json`, `.codex/README.md`, `.codex/workflow.md`
- `non_functional_requirements`: 日本語優先。self-review ではなく証拠収集を要求する。orchestrate の downstream にはしない
- `out_of_scope`: PreToolUse hook 実装、`.gateguard.yml` 導入、外部 package install、product code 修正
- `open_questions`: 実運用で hook 化するかは pilot 後に判断する
- `required_reading`: `.codex/README.md`, `.codex/workflow.md`, `.codex/skills/skill-modification/SKILL.md`, 移植元 `gateguard/SKILL.md`

## Artifacts

- `ui_artifact_path`: N/A
- `final_mock_path`: N/A
- `scenario_artifact_path`: N/A
- `final_scenario_path`: N/A
- `implementation_scope_artifact_path`: N/A
- `review_diff_diagrams`: N/A
- `source_diagram_targets`: N/A
- `canonicalization_targets`: N/A

## Work Brief

- `implementation_target`: workflow skill
- `accepted_scope`: direct-use 補助 gate として `gateguard` を追加し、live workflow への位置づけを明記する
- `parallel_task_groups`: none
- `tasks`: skill 本体追加、権限契約追加、agent 指定追加、README / workflow 追記、検証、plan close
- `validation_commands`: `python3 scripts/harness/run.py --suite structure`, `python3 scripts/harness/run.py --suite all`

## Investigation

- `reproduction_status`: N/A
- `trace_hypotheses`: N/A
- `observation_points`: skill 使用時に、編集前確認で実装方針が変わった回数と、false positive による手戻りを記録する
- `residual_risks`: direct-use では実際の tool call を強制停止できない

## Acceptance Checks

- `gateguard` が日本語で読める
- edit / new file / destructive command の gate が分かれている
- repo 固有の MCP file 操作制約と正本確認が含まれている
- `gateguard` が orchestrate handoff ではなく direct-use 補助 gate として説明されている

## Required Evidence

- 追加ファイル一覧
- structure harness 結果
- full harness 結果または失敗理由

## HITL Status

- `functional_or_design_hitl`: not-required
- `approval_record`: user requested Japanese gateguard port on 2026-04-17

## Closeout Notes

- `canonicalized_artifacts`: N/A

## Outcome

- `.codex/skills/gateguard/` に日本語版 direct-use 補助 gate を追加した。
- `.codex/README.md` と `.codex/workflow.md` に、orchestrate handoff ではない補助 gate として位置づけを追記した。
- `find .codex/skills -name '*.json' -print0 | xargs -0 -n1 python3 -m json.tool >/dev/null` が通過した。
- `python3 scripts/harness/run.py --suite structure` が通過した。
- `python3 scripts/harness/run.py --suite all` が通過した。
