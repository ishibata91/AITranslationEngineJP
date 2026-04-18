---
name: implementation-investigate-observe
description: GitHub Copilot 側の一時観測点知識 package。
---

# Implementation Investigate Observe

## 目的

この skill は知識 package である。
`investigator` agent が owned_scope 内に一時観測点を add / remove する時の判断基準を提供する。

## いつ参照するか

- 一時 log、probe、assertion などを使って観測する時
- temporary_changes と cleanup_status を返す時
- 観測点を返却前に除去する時

## 参照しない場合

- 恒久修正を行う時
- product test を追加する時
- cleanup 不能な観測変更が必要な時

## 原則

- 一時観測点は owned_scope 内に限る
- 観測目的を明確にする
- 返却前に必ず除去する
- cleanup_status を必ず返す

## DO / DON'T

DO:
- temporary_changes に path と目的を残す
- cleanup の validation を行う
- 除去不能なら stop する

DON'T:
- 観測点を恒久修正として残さない
- owned_scope 外を変更しない
- active contract をこの skill に置かない

## Checklist

- [implementation-investigate-observe-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-investigate-observe/references/checklists/implementation-investigate-observe-checklist.md) を参照する。
