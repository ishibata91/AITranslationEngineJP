---
name: skill-modification
description: skill 自体の追加、整理、改名、権限調整を行う skill。`.codex/skills/` と関連する workflow docs の同期を扱う。
---

# Skill Modification

## Output

- touched skill files
- modified scope
- validation results
- follow-up docs sync

## Rules

- skill 変更前に既存 workflow と対象 skill の責務を確認する
- 変更は `.codex/skills/` と関連 workflow docs に限定する
- 新しい skill には `SKILL.md`、`agents/openai.yaml`、`references/permissions.json` をそろえる
- live workflow に合わない legacy artifact は持ち込まない
- 権限境界や handoff が曖昧なら停止して lane owner へ戻す
- スキルはロールとして扱う｡その役割に関係のない記述は極力排除すること｡

