---
name: skill-modification
description: skill 自体の追加、整理、改名、権限調整を行う skill。`.codex/skills/` と関連する workflow docs の同期を扱う。
---

# Skill Modification

## Output

- touched skill files
- modified scope
- validation results
- follow-up workflow sync

## Rules

- skill 変更前に既存 workflow と対象 skill の責務を確認する
- 変更は `.codex/skills/` と関連 workflow docs に限定する
- 新しい skill には `SKILL.md`、`agents/openai.yaml`、`references/permissions.json` をそろえる
- live workflow に合わない legacy artifact は持ち込まない
- 権限境界や handoff が曖昧なら停止して lane owner へ戻す
- スキルはロールとして扱う｡その役割に関係のない記述は極力排除すること｡
- workflow 記述では、論理レベルの名前と実際の skill / agent 名をできるだけ同じ行に置くこと
- 推奨表記は `論理名 (`actual-skill-name`)` または `論理名 (`actual-agent-name`)` とする
- 人間 review の理解速度と、actual name による検索可能性を同時に満たす記述を優先すること
