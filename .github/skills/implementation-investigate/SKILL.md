---
name: implementation-investigate
description: GitHub Copilot 側の実装時調査の共通知識 package。single_handoff_packet 1 件内で evidence first に調査する判断基準を提供する。
---

# Implementation Investigate

## 目的

`implementation-investigate` は知識 package である。
`investigator` agent が、`single_handoff_packet` 1 件と owned_scope 内で実装時の証拠を集める時の共通判断を提供する。

実行権限、write scope、active contract、handoff は [investigator.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/investigator.agent.md) が持つ。

## いつ参照するか

- 実装前再現、trace、再観測を行う時
- 一時観測点を add / remove する時
- evidence と仮説を分けて返す時

## 参照しない場合

- 恒久修正を行う時
- product test を追加する時
- design-time investigation を行う時

## 知識範囲

- evidence first の観測
- observed facts と hypotheses の分離
- temporary observation の cleanup
- focused skill の選び方

## 原則

- `single_handoff_packet` 1 件と owned_scope を超えない
- evidence のない結論を固定しない
- 一時観測点は返却前に除去する
- 恒久修正と product test 追加を混ぜない

## Focused Skills

- [implementation-investigate-reproduce](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-investigate-reproduce/SKILL.md): 実装前再現
- [implementation-investigate-trace](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-investigate-trace/SKILL.md): 実装中 trace
- [implementation-investigate-observe](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-investigate-observe/SKILL.md): 一時観測点
- [implementation-investigate-reobserve](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-investigate-reobserve/SKILL.md): 修正後再観測
## DO / DON'T

DO:
- 観測条件、command、結果を残す
- temporary changes と cleanup_status を返す
- recommended next step を根拠付きで返す

DON'T:
- 恒久修正を同時に行わない
- product test を追加しない
- mode 別 active contract を使わない

## 参照パターン

- [investigation-patterns.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-investigate/references/patterns/investigation-patterns.md) を参照する。
- 対象は execution path tracing、silent failure hunting、temporary observation、minimal error isolation である。
- validation は repo の command と agent contract の出力要件に従って扱う。

## Checklist

- [implementation-investigate-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-investigate/references/checklists/implementation-investigate-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。

## Agent が持つもの

- active contract: [investigator.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/references/investigator/contracts/investigator.contract.json)
- permissions: [permissions.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/references/investigator/permissions.json)

## Maintenance

- 調査種別の知識差分は focused skill に置く。
- output obligation を skill 本体へ戻さない。
