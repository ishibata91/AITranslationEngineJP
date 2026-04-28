---
name: implementation-investigate-trace
description: Codex implementation lane 側の実装中 trace 知識 package。
---

# Implementation Investigate Trace

## 目的

この skill は知識 package である。
`implementation_investigator` agent が実装中の原因候補、観測点、不足情報を整理する時の判断基準を提供する。

## いつ参照するか

- observed facts と hypotheses を分ける時
- 次の observation point を整理する時
- implement、tests、review、reroute の次 action を判断材料として返す時

## 参照しない場合

- 実装前再現だけを行う時
- 一時観測点の add / remove が主目的の時
- 恒久修正が主目的の時

## 原則

- 観測済み事実と仮説を混ぜない
- trace は owned_scope 内に限定する
- 不足情報を remaining_gaps に残す
- evidence のない結論を固定しない

## DO / DON'T

DO:
- hypotheses に根拠と未確認点を付ける
- observation_points を最小にする
- recommended_next_step を根拠付きで返す

DON'T:
- 恒久修正をしない
- product test を追加しない
- active contract をこの skill に置かない

## Checklist

- [implementation-investigate-trace-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-investigate-trace/references/checklists/implementation-investigate-trace-checklist.md) を参照する。
