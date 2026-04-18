# Skill Template

skill は知識の集積地として扱う。
実行権限、agent contract、handoff 契約は agent の持ち物にする。

このテンプレートは draft であり、既存 skill へ機械適用しない。
Markdown 本文の file reference は `[ファイル名](/Users/iorishibata/Repositories/AITranslationEngineJP/<path>)` の形にする。
JSON、TOML、frontmatter、contract field の path は `/Users/iorishibata/Repositories/AITranslationEngineJP/` から始まるフルパス文字列で書く。
runtime や binding schema が相対 path を要求する場合だけ、schema に従って相対 path を使う。

## 概念

skill は actor ではない。
skill は、agent が参照する知識、判断基準、手順例、anti-pattern、checklist を持つ。

skill は次を持たない。

- 実行権限
- agent contract
- handoff contract
- write scope
- completion packet の責任

## 配置

```text
<absolute-skill-root>/
├── SKILL.md
├── agents/
│   └── openai.yaml
└── references/
    ├── checklists/
    │   └── <skill>-checklist.md
    ├── examples/
    │   └── <topic>.md
    └── patterns/
        └── <topic>.md
```

[checklists](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/<skill-name>/references/checklists/) は必須にする。
[examples](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/<skill-name>/references/examples/) と [patterns](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/<skill-name>/references/patterns/) は必要な時だけ置く。

`mode` や variant として分けたくなる知識は、原則として別 skill にする。
contract-level の責務、権限、handoff、output obligation が違う場合は、skill 側ではなく別 agent に切る。
差分が知識参照だけなら、agent は単一 contract のまま focused skill を読む。

[openai.yaml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/<skill-name>/agents/openai.yaml) は binding だけを持つ。
ここに権限、contract、default prompt は置かない。

## SKILL.md

```markdown
---
name: <actual-skill-name>
description: <この skill が提供する知識領域を 1 文で書く>
---

# <Human Readable Skill Name>

## 目的

<この skill がどの知識を提供するかを書く。>
<実行主体ではなく、agent が参照する知識であることを明記する。>

## いつ参照するか

- <knowledge trigger 1>
- <knowledge trigger 2>
- <knowledge trigger 3>

## 参照しない場合

- <non-trigger 1>
- <non-trigger 2>
- <non-trigger 3>

## 知識範囲

- <knowledge scope 1>
- <knowledge scope 2>
- <knowledge scope 3>

## 原則

- <principle 1>
- <principle 2>
- <principle 3>

## 標準パターン

1. <pattern step 1>
2. <pattern step 2>
3. <pattern step 3>

この手順は知識上の標準例である。
実行順、必須 input、完了条件は agent contract に従う。

## DO / DON'T

DO:
- <推奨する判断 1>
- <推奨する判断 2>
- <推奨する判断 3>

DON'T:
- <避ける判断 1>
- <避ける判断 2>
- <避ける判断 3>

## Checklist

- [<skill>-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/<skill-name>/references/checklists/<skill>-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。

## Examples

- [<topic>.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/<skill-name>/references/examples/<topic>.md)
- [<topic>.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/<skill-name>/references/patterns/<topic>.md)

## Agent が持つもの

- 実行権限
- agent 1:1 contract
- handoff 契約
- write scope
- stop / reroute 条件

## Maintenance

- 長い例は [examples](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/<skill-name>/references/examples/) に移す。
- 判断表は [patterns](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/<skill-name>/references/patterns/) に移す。
- checklist は 3〜6 項目の短い確認に保つ。
- 権限や契約を skill 本体へ戻さない。
- 知識差分が増えたら、内部 mode ではなく focused skill として分ける。
```

## agents/openai.yaml

```yaml
version: 1
agent: <actual-agent-name>
entry: SKILL.md
```

`entry: SKILL.md` は binding schema の値なので相対 path のままにする。
`default_prompt` は置かない。
権限と contract も置かない。

## references/checklists/<skill>-checklist.md

```markdown
# <Skill Name> Checklist

## Knowledge Check

- [ ] <確認観点 1>
- [ ] <確認観点 2>
- [ ] <確認観点 3>

## Common Pitfalls

- [ ] <見落とし 1> を避けた
- [ ] <見落とし 2> を避けた
- [ ] <見落とし 3> を避けた
```

checklist は skill の知識確認として必須にする。
agent が checklist を必須実行するかどうかは agent contract に書く。

## 作成チェック

- skill が知識領域を説明している。
- 権限、agent contract、output obligation を持っていない。
- [<skill>-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/<skill-name>/references/checklists/<skill>-checklist.md) がある。
- 長い例や判断表が [references](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/<skill-name>/references/) に分離されている。
- [openai.yaml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/<skill-name>/agents/openai.yaml) が binding だけになっている。
- 知識差分は内部 mode ではなく focused skill として分かれている。
