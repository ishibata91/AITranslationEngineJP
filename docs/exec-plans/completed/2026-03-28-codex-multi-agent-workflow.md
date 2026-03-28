# Codex Multi-Agent Workflow

- Date: 2026-03-28
- Status: Completed

## Goal

作業フローの正本を `.codex/` 配下へ移し、role と workflow を分離したマルチエージェント運用に切り替える。

## Scope

- `.codex/README.md` の追加
- `.codex/agents/*.md` に役割仕様を追加
- `.codex/skills/` に heavy / lightweight workflow skill を追加
- `docs/exec-plans/templates/` に plan template を追加
- 入口文書と harness を `.codex/` 前提へ更新

## Acceptance Criteria

- `.codex/README.md` から heavy / light 判定と参照順序が辿れる
- `.codex/agents/*.md` に責務境界が明記される
- `.codex/skills/*/SKILL.md` と `agents/openai.yaml` が存在する
- heavy / light 用の plan template が存在する
- harness が `.codex/` の存在と必須語を検査できる

## Outcome

- `.codex/` をチーム共有の workflow 正本として追加した
- agent role と workflow skill を分離した
- docs は仕様と設計の正本として維持し、作業方法だけ `.codex/` へ移した

## Follow-up

- プロダクト固有 skill を必要に応じて追加する
- design harness をより意味的なチェックへ拡張する
