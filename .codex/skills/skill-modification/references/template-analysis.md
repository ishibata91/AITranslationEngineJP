# Skill And Agent Template Analysis

## 結論

この draft では、agent を実行主体、skill を知識の集積地として定義する。
権限と contract は agent の持ち物にする。

contract は agent に対して 1:1 にする。
mode / variant ごとの active contract file は作らない。
contract-level の差分が必要なら別 agent に切り、知識差分だけなら focused skill として分ける。

skill は agent の判断を助けるが、agent の権限や output obligation を上書きしない。

## everything-claude-code から見える設計

`everything-claude-code` では root agent が強い実行単位になっている。
agent は persona、tools、model、workflow、判断基準、output format を持つ。

skill は TDD、API design、security review のような再利用知識を持つ。
原則、具体例、GOOD / BAD、checklist、anti-pattern が中心である。

この構造は、agent を単体で呼ぶカタログ型運用に合理的である。

## 採用する概念

- agent は actor である。
- skill は knowledge package である。
- permissions は actor に属する。
- contract は actor に属し、agent 1:1 にする。
- contract-level の差分は別 agent に切る。
- checklist は skill に置けるが、実行義務は agent contract が決める。
- handoff は agent 間契約として扱う。

## Agent が定義するもの

- 実行主体としての role
- tool 権限
- user-invocable か subagent 専用か
- source of truth
- allowed / forbidden action
- allowed / forbidden write scope
- agent 1:1 contract
- handoff contract
- stop / reroute 条件

agent は「誰が何をしてよいか」と「何を返すか」を定義する。

## Skill が定義するもの

- 知識領域
- いつ参照するか
- 原則
- 標準パターン
- DO / DON'T
- examples
- anti-pattern
- checklist

skill は「その仕事をどう考えるか」と「何を良い判断と見るか」を定義する。
mode や variant として分けたくなる知識は focused skill に分ける。

## 変更したテンプレート方針

- `skill-template.md` から permissions、agent contract、write scope、completion responsibility を外した。
- `agent-template.md` に permissions、agent 1:1 contract、handoff contract を移した。
- contract-level の差分を selector_requirements で持たず、別 agent に切る方針にした。
- `references/checklists/` は skill 側に残した。
- checklist の実行義務は agent contract 側が決めることにした。
- `default_prompt` は引き続き不採用にした。

## 現行 repo との衝突

現行には旧方式の skill-side `permissions.json` や contract slice が残っている。
これは互換用 legacy として残し、新しい正本にはしない。

この draft を広く採用する場合は、次を別 task で同期する必要がある。

- [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/skill-modification/SKILL.md)
- [README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/README.md)
- [.github/agents](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/) と [.github/skills](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/) の配置規約
- 既存 skill 配下の `references/permissions.json` の扱い
- 既存 mode / variant contract slice の扱い

今回は template draft と `distiller` trial の再定義だけで、既存 workflow 全体には適用しない。

## 判断理由

権限と契約を skill に置くと、知識と実行責任が混ざる。
その結果、agent が何に責任を持つのか読みにくくなる。

contract-level の差分を selector として単一 contract に詰めると、実質的には複数 agent の責務が混ざる。
別 agent に切る方が、入力、出力、停止条件、handoff の所在が明確になる。

agent が権限と契約を持つ方が、実行主体、責任、handoff、停止条件の所在が揃う。
skill は知識として複数 agent から再利用しやすくなる。

## 次の判断ポイント

- agent-owned references の配置を [.github/agents/references](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/references/) に固定するか。
- Codex 側にも [.codex/agents](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/) を作るか。
- 既存 skill の `permissions.json` を legacy として残すか移すか。
- checklist を skill 必須にするか、agent contract が参照する時だけ作るか。
- 現行 live workflow にいつ適用するか。

## 今回の結論

conceptual template は作り直した。
実適用は `distiller` trial の範囲に留める。

今後の基本線は次である。

- agent: 実行主体、権限、agent 1:1 contract、handoff
- skill: 知識、判断基準、例、checklist
- contract: agent-owned かつ agent 1:1
- permissions: agent-owned
- contract 差分: 別 agent
