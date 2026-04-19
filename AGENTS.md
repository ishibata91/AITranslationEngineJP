# AGENTS.md

会話と作業は日本語を基本にする。
英語の key、既存名、command は必要な時だけ使う。
tmp/code-map/index.jsonにコード地図がある

## 会話ルール

- 見出しなしの長文を避け、2〜4 個の短い見出しに分ける
- 見出しは `##` を使う
- 箇条書きは 3〜6 件に抑える
- 1 段落は 3 文以内にする
- 変更報告は 1 行 1 ファイルにする
- 長い出力の末尾には `SUMMARY` を付ける

## 目的

AITranslationEngineJp は Skyrim Mod 向け翻訳エンジンです。
この repo は agent-first で進めます。

## 参照マップ

最初に使う skill の `references/permissions.json` を読む。
その後、必要な文書だけ読む。

- workflow 正本: `.codex/README.md`
- Codex 入口: `.codex/skills/propose-plans/SKILL.md`
- Copilot 実装入口: `.github/skills/implementation-orchestrate/SKILL.md`
- 仕様入口: `docs/index.md`
- 長期原則: `docs/core-beliefs.md`
- 恒久要件: `docs/spec.md`
- architecture: `docs/architecture.md`
- 作業計画: `docs/exec-plans/`

## 強い制約

- Codex は設計、計画、handoff、docs 正本化を担当する
- GitHub Copilot は承認済み `implementation-scope` から実装する
- AI design review は行わず、人間が design bundle を review する
- Copilot は `docs/`、`.codex/`、`.github/skills`、`.github/agents` を変更しない
- docs 正本化は Codex の `updating-docs` だけが扱う

## 実装前に確認すること

1. 使う skill の `references/permissions.json` を読む
2. `docs/index.md` から関係する文書だけ読む
3. active / completed plan に同種 task がないか確認する
4. 編集前に gateguard の事実確認を行う

## 実装後にやること

1. 必要な follow-up を plan か issue に記録する
2. docs 正本更新は human 承認済みの時だけ行う
3. 完了した plan を `docs/exec-plans/completed/` へ移す

## 補足

- library の書き方は MCP_DOCKER 経由で Context7 を確認する
- wails は `npm run dev:wails:docker-mcp` で起動する
- Playwright MCP は `http://host.docker.internal:34115` に接続する
- Sonar project は `ishibata91_AITranslationEngineJP`
