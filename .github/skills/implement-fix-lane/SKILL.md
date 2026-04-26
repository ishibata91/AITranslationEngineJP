---
name: implement-fix-lane
description: GitHub Copilot 側の fix lane 恒久修正知識 package。
---

# Implement Fix Lane

## 目的

この skill は知識 package である。
`implementer` agent が `accepted_fix_scope` の恒久修正を行う時に、再現条件と矛盾しない変更へ限定する判断基準を提供する。

## いつ参照するか

- `task_mode: fix` の owned_scope を実装する時
- reproduction evidence または trace result に基づき修正する時
- residual risk と未解消ケースを closeout に残す時

## 参照しない場合

- 新機能や refactor の実装を行う時
- 再現条件が不足している時
- 原因が未確認なのに恒久修正する時

## 原則

- accepted_fix_scope を超えない
- 再現条件に関係しない整理を入れない
- trace_or_analysis_result と矛盾しない変更に限る
- lane_context_packet を確認して product code だけを変更する
- `APIテスト` 先行時だけ tester output も確認する
- 未解消ケースを closeout に残す

## DO / DON'T

DO:
- 修正前後で同じ条件の validation を比較する
- residual risk を明示する
- fix scope と touched files を対応づける

DON'T:
- unrelated cleanup を混ぜない
- 原因断定を evidence なしに広げない
- product test、fixture、snapshot、test helper を変更しない
- active contract をこの skill に置かない

## Checklist

- [implement-fix-lane-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implement-fix-lane/references/checklists/implement-fix-lane-checklist.md) を参照する。
