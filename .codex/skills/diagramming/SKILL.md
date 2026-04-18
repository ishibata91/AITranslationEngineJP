---
name: diagramming
description: Codex 側の図作成知識 package。D2、PlantUML、structure diff の source と review artifact の扱いを提供する。
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
- D2 または PlantUML source を作成、更新する時
- diagram source と rendered artifact の責務を分ける時

## 参照しない場合

- UI design の primary artifact が Figma で足りる時
- product code の構造を実装で変更する時
- docs 正本化だけが目的の時

## 知識範囲

- `structure-diff`、`d2`、`plantuml` の図種別
- source of truth と review artifact の分離
- render、validate、差分確認の扱い
- D2 や diagram library の事前確認順

## 原則

- 図の source of truth を先に明示する
- review 用 SVG や screenshot を正本にしない
- D2 や library の書き方は Context7 で確認する
- validation できない diagram は gap として返す

## 標準パターン

1. 図の目的、読者、正本 source を確認する。
2. diagram kind を選び、必要な source file を特定する。
3. D2 や library の書き方が関係する場合は Context7 を確認する。
4. source と review artifact を分けて作る。
5. render / validation 結果を残す。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は `designer` agent contract に従う。

## DO / DON'T

DO:
- source と review artifact を別物として扱う
- validation 結果と未確認 gap を残す
- 図が補う設計判断を明示する

DON'T:
- review 用 artifact を正本にしない
- product code 変更を diagramming に混ぜない
- source 不明のまま図を更新しない

## Checklist

- [diagramming-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/diagramming/references/checklists/diagramming-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は呼び出し元 agent contract が決める。

## References

- primary agent spec: [designer.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.agent.md)
- primary agent contract: [designer.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/designer/contracts/designer.contract.json)
- explicit helper agent: [diagrammer.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/diagrammer.agent.md)

## Maintenance

- 権限、write scope、output obligation を skill 本体へ戻さない。
- mode-specific active contract を skill 側に増やさない。
- 標準 flow では diagram を designer の資料作成 scope から外さない。
