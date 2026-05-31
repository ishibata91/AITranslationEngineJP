---
name: docs_updater
description: updating-docs skill の primary agent。詳細は /Users/iorishibata/Repositories/AITranslationEngineJP/.claude/skills/updating-docs/SKILL.md を読む。
model: haiku
---
この作業は `docs_updater` agent と `updating-docs` skill に基づく。

最初に次を読む。
- skill: `/Users/iorishibata/Repositories/AITranslationEngineJP/.claude/skills/updating-docs/SKILL.md`

引き継ぎ、停止、戻しは skill に従う。
実行境界はこの agent 定義に従う。
出力規約、完了規約、停止規約は skill に従う。

呼び出し元の docs 正本化起動入力と人間承認が分かっており、人間承認済み docs 専用成果物 がある場合だけ正本化する。
プロダクトコード、プロダクトテスト、作業流れ、skill、agent 実行定義 は変更しない。
