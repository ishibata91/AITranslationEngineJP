---
name: review-implementation
description: GitHub Copilot 側の implementation review 知識 package。
---

# Review Implementation

## 目的

この skill は知識 package である。
`reviewer` agent が implementation review を行う時に、差分が owned_scope と一致するか、test と validation が十分かを確認する判断基準を提供する。

## いつ参照するか

- 差分が owned_scope に収まっているか確認する時
- 必要な product test と validation command を確認する時
- backend を含む場合に Sonar gate を確認する時

## 参照しない場合

- UI check が主目的の時
- design review が必要な時
- 修正を同時に行う時

## 原則

- implementation-scope と review target diff を照合する
- 好みや将来改善で reroute しない
- finding は再現できる形で返す
- 修正は行わない

## DO / DON'T

DO:
- scope 外 diff、missing test、failed validation を分ける
- Sonar gate の該当可否を明示する
- pass の場合も未実行 validation を残す

DON'T:
- design review をしない
- 新しい要件解釈を追加しない
- active contract をこの skill に置かない

## Checklist

- [review-implementation-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/review-implementation/references/checklists/review-implementation-checklist.md) を参照する。
