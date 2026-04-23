---
name: skill-modification
description: Codex 側の skill / agent 変更知識 package。skill を knowledge package、agent を実行主体として整理する基準を提供する。
---

# Skill Modification

## 目的

`skill-modification` は知識 package である。
`.codex` の skill と agent runtime を整理する時に、配置、path policy、skill 正本、agent TOML、agent-owned contract の扱いを判断するための知識を提供する。

実行境界、source of truth、handoff、stop / reroute は [design-bundle](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md) を参照する。

## いつ参照するか

- skill 自体を追加、整理、改名、分割する時
- skill、agent TOML、permissions、contract の配置を変える時
- workflow docs と skill / agent の責務を同期する時
- `.codex` を直接変更できず、`tmp/codex` staged apply が必要な時

## 参照しない場合

- product code または product test を変更する時
- docs 正本の product 仕様を変更する時
- 権限境界が不明で lane owner 判断が必要な時

## 知識範囲

- [skill-agent-concept.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/skill-modification/references/skill-agent-concept.md) の概念分担
- skill template と agent runtime template
- Markdown / JSON / TOML の path policy
- agent 1:1 contract と legacy references の扱い
- `tmp/codex` staged apply と人間実行 script の扱い

## 原則

- skill は人間可読な正本、agent は runtime binding と機械契約の owner として扱う
- permissions と contract は agent 側に置く
- handoff、stop / reroute、source of truth の人間可読説明は skill 側に置く
- contract は agent 1:1 にする
- mode / variant ごとの active contract file を増やさない
- live workflow にない legacy artifact を復活させない
- staged apply は反映元を破壊せず、削除差分を明示確認してから正本へ写す

## 標準パターン

1. `skill-agent-concept.md` と対象 agent の permissions を読む。
2. 既存 workflow と対象 skill / agent の責務を確認する。
3. role、source of truth、handoff、stop / reroute の人間可読説明を skill 側へ集約する。
4. agent 側には TOML binding、permissions、1:1 contract を置く。
5. checklist を skill references に置き、旧 `.agent.md` は削除する。
6. path policy と workflow 名の actual name 対応を確認する。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は `designer` agent contract に従う。

## Staged Apply Pattern

`.codex` へ直接書けない時は、[staged-apply-flow.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/skill-modification/references/patterns/staged-apply-flow.md) を使う。
反映済み file は `tmp/codex/files/<repo-relative-path>` に置き、人間が `scripts/codex/apply_tmp_codex.py` または VSCode task で正本へ上書き、追加、削除する。
通常 apply は基本的に人間が実行する。
Codex は staged file の作成と `--check-only` による final gate 確認までを担当する。
人間から明示指示がある場合だけ、Codex が通常 apply を試行できる。

反映 script は最終チェックとして次を行う。

- 反映元 file の hash を反映前後で比較し、反映元を破壊していないことを確認する
- 反映先と反映元の diff を表示し、削除行があれば停止する
- file 削除が必要な時は `tmp/codex/delete-paths.txt` に対象を列挙する
- 記載削除または file 削除が必要な時は `tmp/codex/deletion-rationale.md` に削除対象、理由、代替参照先を記録してから再実行する
- JSON / TOML / Markdown / PlantUML など、対象 file の最低限の構文確認を行う
- `.codex` へ直接反映しない段階では `--check-only` で同じ final gate だけを確認する
- 通常 apply は人間実行を基本とし、Codex は明示指示なしに通常 apply へ進まない
- 通常 apply が成功したら `tmp/codex` を全削除する

## DO / DON'T

DO:
- 論理名と actual skill / agent 名を同じ行に置く
- Markdown 本文の file reference はフルパスリンクにする
- 権限境界が曖昧なら停止する
- staged apply script は copy 前に diff と削除行を見せる

DON'T:
- skill 本体へ hard permissions や active contract を戻さない
- product 実装や docs product 仕様変更を混ぜない
- default_prompt を導入しない
- staged apply script で反映元 directory を削除しない

## Checklist

- [skill-modification-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/skill-modification/references/checklists/skill-modification-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は `designer` agent contract が決める。

## References

- concept: [skill-agent-concept.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/skill-modification/references/skill-agent-concept.md)
- skill template: [skill-template.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/skill-modification/references/skill-template.md)
- agent template: [agent-template.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/skill-modification/references/agent-template.md)
- Codex TOML template: [codex-agent-template.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/skill-modification/references/codex-agent-template.toml)
- staged apply pattern: [staged-apply-flow.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/skill-modification/references/patterns/staged-apply-flow.md)
- design runtime: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md)

## Maintenance

- template 方針と実適用の差分を放置しない。
- skill は knowledge package、agent は actor という分担を崩さない。
- workflow docs の同期が必要なら範囲を明示する。
- staged apply 手順は、反映元保全と削除妥当性確認を外さない。
