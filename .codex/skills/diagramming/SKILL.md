---
name: diagramming
description: Codex 側の図作成知識 package。PlantUML と structure diff の source と review artifact の扱いを提供する。
---

# Diagramming

## 目的

`diagramming` は知識 package である。
`designer` agent が diagram を必要資料として作る時に、source of truth、review artifact、validation の見方を提供する。

標準 `propose_plans` flow では、diagram は `designer` の資料作成 scope に含める。
明示的に diagrammer が指定された補助用途では [diagrammer.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/diagrammer.agent.md) も参照できる。

## いつ参照するか

- design bundle に diagram が必要な時
- structure diff を review 可能な図にする時
- PlantUML source を作成、更新する時
- diagram source と rendered artifact の責務を分ける時

## 参照しない場合

- UI design の primary artifact が HTML mock で足りる時
- product code の構造を実装で変更する時
- docs 正本化だけが目的の時

## 知識範囲

- PlantUML を使った structure diff と設計補助図
- source of truth と review artifact の分離
- render、validate、差分確認の扱い
- 一時 PNG と AI 目視による可読性確認
- 図を分割すべき大きさの判断

## 原則

- 図の source of truth を先に明示する
- review 用 SVG、PNG、screenshot を正本にしない
- PlantUML や diagram library の書き方は Context7 で確認する
- 図の方向は常に上から下にする
- validation できない diagram は gap として返す
- 原則日本語で書くが、必要に応じて英語も使う

## 標準パターン

1. 図の目的、読者、正本 source を確認する。
2. diagram kind を選び、必要な source file を特定する。
3. PlantUML や library の書き方が関係する場合は Context7 を確認する。
4. PlantUML source と review artifact を分けて作る。
5. 一時 PNG を生成し、AI が画像で可読性を確認する。
6. render / validation 結果と未確認 gap を残す。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は `designer` agent contract に従う。

## DO / DON'T

DO:
- source と review artifact を別物として扱う
- validation 結果と未確認 gap を残す
- 図が補う設計判断を明示する
- 大きすぎる図は overview と detail に分ける

DON'T:
- review 用 artifact を正本にしない
- product code 変更を diagramming に混ぜない
- source 不明のまま図を更新しない
- PlantUML 以外の diagram source を新規作成しない

## Checklist

- [diagramming-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/diagramming/references/checklists/diagramming-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は呼び出し元 agent contract が決める。

## Templates

- [review-diff-style-legend.puml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/diagramming/references/templates/review-diff-style-legend.puml) を差分図の標準 style / legend として使う。
- [review-diff-template.puml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/diagramming/references/examples/review-diff-template.puml) を新規差分図の開始点にする。

## Readability Gate

- 一時 PNG は `/tmp` へ生成し、repo の正本にしない。
- AI が画像を確認し、線、文字、余白、線長、交差、legend の干渉を確認する。
- 可読性が悪い場合は、配置調整ではなく図の分割を先に検討する。

## Readability Patterns

- 正本図の用語、package 名、layer 名は、対象 docs の構造主語と節名に合わせる。
- architecture 図では [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) の構造主語、依存方向、節名を優先する。
- 正本図では可読性のために node を勝手に畳まず、まず edge、legend、note の粒度で調整する。
- package / group は意味単位または layer 単位で作り、上から下に主依存を読める順序に固定する。
- 長距離の wiring や bootstrap の concrete 生成線は、全 edge を描くより legend へ要約する。
- ER 図は table の意味単位で group を作り、関係線を最小限にして読む順序を安定させる。
- docs 正本図は neutral な構造図にし、差分図 style は `docs/exec-plans/` 配下の review artifact に限定する。
- `skinparam linetype ortho` のような角ばった線指定は、読みやすさが下がる場合があるため既定では使わない。

## Split Rule

- primary node が 12 個を超える時は分割する。
- attribute 付き class / table が 8 個を超える時は分割する。
- visible edge が 20 本を超える時は分割する。
- package / boundary が 4 個を超える時は分割する。
- package をまたぐ長距離 edge が複数ある時は overview と detail に分ける。

## References

- primary agent spec: [designer.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.agent.md)
- primary agent contract: [designer.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/designer/contracts/designer.contract.json)
- explicit helper agent: [diagrammer.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/diagrammer.agent.md)

## Maintenance

- 権限、write scope、output obligation を skill 本体へ戻さない。
- mode-specific active contract を skill 側に増やさない。
- 標準 flow では diagram を designer の資料作成 scope から外さない。
