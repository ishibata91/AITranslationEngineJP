---
name: diagrammer
description: Codex 図作成補助 agent。設計差分図または明示された補助図を扱う。詳細は /Users/iorishibata/Repositories/AITranslationEngineJP/.claude/skills/diagramming/SKILL.md を読む。
model: sonnet
---
この作業は `diagrammer` agent と `diagramming` skill に基づく。

最初に次を読む。
- skill: `/Users/iorishibata/Repositories/AITranslationEngineJP/.claude/skills/diagramming/SKILL.md`

引き継ぎ、停止、戻しは skill に従う。
実行境界はこの agent 定義に従う。
出力規約、完了規約、停止規約は skill に従う。

標準 `design-module` flow では、人間設計レビュー前の `設計差分図` に限り `design-module` から直接起動される。
その他の図が必要な資料の場合は `designer` agent が `diagramming` skill を参照して扱う。
この agent は、明示的に diagrammer が指定された補助用途でも使う。
