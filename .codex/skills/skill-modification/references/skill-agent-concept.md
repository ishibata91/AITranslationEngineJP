# Skill / Agent Concept

この文書は、AITranslationEngineJP における skill と agent の概念分担を説明する。
テンプレートを読む前に、この文書で前提をそろえる。

この文書は draft である。
既存 workflow へ適用するには、別 task で [README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/README.md)、[.github/agents](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/)、[.github/skills](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/)、既存 permissions の扱いを同期する。

## 基本方針

agent は実行主体として扱う。
skill は知識の集積地として扱う。

権限、agent contract、handoff、stop / reroute は agent が持つ。
skill は agent の判断を助けるが、agent の契約を上書きしない。

## Agent の概念

agent は actor である。
つまり、実際に task を受け取り、tool を使い、成果を返す主体である。

agent は次を定義する。

- 何を実行する agent なのか
- どの tool を使えるのか
- どの source of truth を読むのか
- どこまで書き換えてよいのか
- 何を入力として受け取り、何を出力として返すのか
- どこで停止し、どこへ reroute するのか

agent は責任の境界を持つ。
そのため permissions と agent contract は agent 側に置く。

## Contract の概念

contract は agent に対して 1:1 で置く。
mode / variant / task kind ごとの active contract file は作らない。

責務、権限、handoff、output obligation が分かれるなら、単一 contract に selector を足さず別 agent に切る。
差分が知識参照だけなら、agent は単一 contract のまま focused skill を読む。

この分け方により、agent が読む契約の正本が一か所に固定される。
同時に、実行責務の違いを contract 内の疑似分岐として隠さずに済む。

## Skill の概念

skill は knowledge package である。
つまり、agent が参照する知識、判断基準、標準パターン、例、anti-pattern を持つ。

skill は次を定義する。

- どの知識領域を扱うのか
- いつ参照するのか
- 何を良い判断と見るのか
- どの手順や pattern が有効なのか
- どの失敗を避けるべきか
- どの checklist で見落としを防ぐのか

skill は実行主体ではない。
そのため write scope、handoff contract、completion packet の責任を持たない。

## 分担

| 項目 | Agent | Skill |
| --- | --- | --- |
| 概念 | 実行主体 | 知識の集積地 |
| 主な責任 | 実行、権限、契約、handoff | 判断基準、例、pattern、checklist |
| permissions | 持つ | 持たない |
| contract | agent 1:1 で持つ | 持たない |
| tool 権限 | 持つ | 持たない |
| source of truth | 実行時に固定する | 読み方や判断軸を説明する |
| checklist | 実行義務を決める | 知識確認として提供する |
| handoff | agent 間契約として持つ | 持たない |

この分担により、agent は「何をしてよいか」を明確にする。
skill は「どう考えるとよいか」を再利用可能にする。

## なぜこの分担がよいか

権限と契約は、実行する主体に置く方が自然である。
書き換え、検証、handoff、停止判断の責任を持つのは agent だからである。

contract-level の差分を selector として単一 contract に詰めると、実質的には複数 agent の責務が混ざる。
差分が実行責務なら別 agent に切る方が、責任境界と読み順が安定する。

skill に権限を置くと、知識と実行責任が混ざる。
その結果、同じ skill を複数 agent が参照した時に、どの権限が正しいのか分かりにくくなる。

agent-owned contract にすると、責任境界が読みやすい。
agent は自分の input、output、write scope、reroute 条件を一か所で確認できる。

skill を知識に寄せると、再利用しやすい。
同じ API 設計 skill、調査 skill、review skill を、違う agent が別 contract で参照できる。

## everything-claude-code からの学び

`everything-claude-code` では、agent が強い実行単位として作られている。
agent は persona、tools、workflow、判断基準、output format を持つ。

一方で skill は、TDD、API design、security review のような知識 package として機能する。
原則、GOOD / BAD、checklist、具体例、anti-pattern が中心である。

この設計は、agent を直接呼ぶ catalog 型運用では合理的である。
AITranslationEngineJP では、そこから次の方針を採る。

- agent は実行責任を持つ
- skill は知識責任を持つ
- 契約と権限は agent に寄せる
- contract は agent 1:1 にする
- 具体例と判断基準は skill に寄せる

## 配置の考え方

agent-owned の情報は agent 側 references に置く。
例は次の形を基本にする。

```text
<agent-root>/
├── <agent-name>.agent.md
└── references/
    └── <agent-name>/
        ├── permissions.json
        └── contracts/
            └── <agent-name>.contract.json
```

skill-owned の情報は skill 側 references に置く。
例は次の形を基本にする。

```text
<skill-root>/
├── SKILL.md
├── agents/
│   └── openai.yaml
└── references/
    ├── checklists/
    │   └── <skill>-checklist.md
    ├── examples/
    └── patterns/
```

[openai.yaml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/<skill-name>/agents/openai.yaml) は binding だけを持つ。
`default_prompt`、permissions、contract は置かない。

## 設計上の注意

agent と skill に同じ項目を重複して書かない。
重複させると、どちらが正本か分からなくなる。

skill に output obligation を戻さない。
skill の checklist は知識確認であり、出力義務は agent contract が決める。

agent に長い知識集を持たせない。
長い判断表、例、anti-pattern は skill の references に分離する。

既存 repo には、skill 配下に `permissions.json` を置く旧方針が残っている。
この concept を live workflow に採用する場合は、既存方針を別 task で同期する。

## まとめ

agent は「誰が、何を、どこまで実行してよいか」を定義する。
skill は「その仕事をどう考え、何を良い判断と見るか」を定義する。

contract は agent ごとに 1 ファイルへ固定する。
contract-level の差分があるなら別 agent に切り、知識差分だけなら focused skill として扱う。
