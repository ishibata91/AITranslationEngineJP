---
name: implementation-distill-implement
description: GitHub Copilot 側の新規実装・拡張向け context 圧縮知識 package。
---

# Implementation Distill Implement

## 目的

この skill は知識 package である。
`implementation-distiller` agent が新規実装や拡張の handoff を整理する時に、implementation facts、constraints、validation entry を抽出する判断基準を提供する。

## いつ参照するか

- 承認済み implementation-scope を実装可能な facts へ落とす時
- handoff、owned_scope、validation entry を明示する時
- 変更対象 package / component / test surface を整理する時

## 参照しない場合

- fix の再現症状を整理する時
- refactor の不変条件を整理する時
- product code / product test を変更する時

## 原則

- source artifacts と implementation-scope の該当 handoff を先に固定する
- 既存境界と依存方向を実装前 context に残す
- 実装者が最初に読む file と順番を残す
- docs や design artifact は必要な判断だけに圧縮する

## DO / DON'T

DO:
- path catalog から必要 file だけ summary / full に上げる
- owned_scope に直接関係する code pointer を優先する
- validation command と completion signal を残す

DON'T:
- 要件や設計を追加しない
- owned_scope 外を広く探索しない
- active contract をこの skill に置かない

## Checklist

- [implementation-distill-implement-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-distill-implement/references/checklists/implementation-distill-implement-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。
