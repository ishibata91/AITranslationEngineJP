---
name: implementation-distill-implement
description: Codex implementation lane 側の新規実装・拡張向け context 圧縮知識 package。
---

# Implementation Distill Implement

## 目的

この skill は知識 package である。
`implementation_distiller` agent が新規実装や拡張の handoff を整理する時に、implementation facts、constraints、validation entry を抽出する判断基準を提供する。

## いつ参照するか

- single_handoff_packet 1 件を実装可能な facts へ落とす時
- handoff、owned_scope、validation entry を明示する時
- 変更対象 package / component / test surface を整理する時

## 参照しない場合

- fix の再現症状を整理する時
- refactor の不変条件を整理する時
- product code / product test を変更する時

## 原則

- single_handoff_packet、owned_scope、validation entry を先に固定する
- 既存境界と依存方向を、handoff に効く最小単位の architecture constraint として実装前 context に残す
- lint-sensitive な規約は `docs/lint-policy.md` と `docs/architecture.md` から、今回の handoff に効く禁止依存、許可境界、format / lint / type check 観点だけを抽出して残す
- patch 生成に必要な fix_ingredients を構造単位で残す
- 類似していても実装不要な context は distracting_context に分ける
- 新規 method、interface、field が必要そうな時は、既存定義の present / absent を code fact として確認する
- first_action は 1 completion_signal clause に限定し、同じ clause の最小 closure chain を示す
- existing_patterns と cheap validation entry を実装開始前 context に残す
- 実装者が最初に触る file、symbol、line number、変更種別を残す
- 要件、実装方針、決定事項は implementation_implementer が再読不要な粒度で要約する
- docs や design artifact は必要な判断だけに圧縮する

## DO / DON'T

DO:
- path catalog から必要 file だけ summary / full に上げる
- owned_scope に直接関係する code pointer を優先する
- implementation_implementer が architecture 正本を全文再読しなくてよいよう、今回の handoff に効く境界だけを `requirements_policy_decisions` へ圧縮する
- frontend handoff では generated `wailsjs`、gateway、screen controller、usecase の境界違反を起こす禁止 import と許可経路を残す
- backend handoff では usecase、service、repository、adapter の依存方向と禁止 import を残す
- mixed handoff では API、Wails binding、DTO、gateway、adapter contract の接合点だけを change target と対応づけて残す
- fix_ingredients と distracting_context を分ける
- first_action と change_targets を path、symbol、line number 付きで返す
- 推測ではなく present / absent の code fact で追加対象を説明する
- existing_patterns がない場合は searched scope と impact を返す
- validation command は最初の cheap check を優先する
- requirements_policy_decisions に implementation_implementer impact を残す
- local validation 前に踏みやすい lint 規約があれば、どの command で拾われるかまで要約する
- validation command と completion signal を残す

DON'T:
- 要件や設計を追加しない
- 実 code を読まず handoff の文章を言い換えない
- 類似 context を required_reading に混ぜない
- `docs/architecture.md` や `docs/lint-policy.md` の全文読解を implementation_implementer へ丸投げしない
- 推測だけで method / interface / field 追加を決めない
- first_action に partial、複数 clause、曖昧な advance を書かない
- 要件、実装方針、決定事項を required_reading に丸投げしない
- owned_scope 外を広く探索しない
- active contract をこの skill に置かない

## Checklist

- [implementation-distill-implement-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-distill-implement/references/checklists/implementation-distill-implement-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。
