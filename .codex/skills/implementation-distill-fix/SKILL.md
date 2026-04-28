---
name: implementation-distill-fix
description: Codex implementation lane 側の fix 向け context 圧縮知識 package。
---

# Implementation Distill Fix

## 目的

この skill は知識 package である。
`implementation_distiller` agent が fix handoff を整理する時に、症状、再現済み事実、修正対象、validation entry を分ける判断基準を提供する。

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
- 失敗を閉じるために必要な fix_ingredients を構造単位で残す
- 再現に似ているだけで修正に不要な context は distracting_context に分ける
- 修正対象は path、symbol、line number、変更種別で返す
- accepted fix scope、決定済み方針、禁止事項は implementation_implementer が再読不要な粒度で要約する
- 再現条件に関係しない整理を入れない

## DO / DON'T

DO:
- reproduction evidence を path / command と一緒に残す
- trace_or_analysis_result と accepted_fix_scope を対応づける
- fix_ingredients と distracting_context を分ける
- first_action と change_targets を path、symbol、line number 付きで返す
- requirements_policy_decisions に fix 方針と out of scope を残す
- residual risk と未解消ケースを分ける

DON'T:
- 原因断定を先取りしない
- 実 code を読まず handoff の文章を言い換えない
- 類似 context を required_reading に混ぜない
- fix 方針や決定事項を required_reading に丸投げしない
- product code / product test を変更しない
- active contract をこの skill に置かない

## Checklist

- [implementation-distill-fix-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-distill-fix/references/checklists/implementation-distill-fix-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。
