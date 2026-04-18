---
name: implementation-distill-fix
description: GitHub Copilot 側の fix 向け context 圧縮知識 package。
---

# Implementation Distill Fix

## 目的

この skill は知識 package である。
`implementation-distiller` agent が fix handoff を整理する時に、症状、再現済み事実、修正対象、validation entry を分ける判断基準を提供する。

## いつ参照するか

- 再現済み症状と trace 結果を整理する時
- `accepted_fix_scope` を実装前 packet に圧縮する時
- 未解消ケースと residual risk を残す時

## 参照しない場合

- 新規実装の context を整理する時
- refactor の不変条件を整理する時
- 原因を evidence なしで断定する時

## 原則

- 症状、再現済み事実、仮説、未確認事項を混ぜない
- 長い log や stack trace は要点と path / command に圧縮する
- 修正対象と validation entry だけを実装者向けに残す
- 再現条件に関係しない整理を入れない

## DO / DON'T

DO:
- reproduction evidence を path / command と一緒に残す
- trace_or_analysis_result と accepted_fix_scope を対応づける
- residual risk と未解消ケースを分ける

DON'T:
- 原因断定を先取りしない
- product code / product test を変更しない
- active contract をこの skill に置かない

## Checklist

- [implementation-distill-fix-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-distill-fix/references/checklists/implementation-distill-fix-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。
