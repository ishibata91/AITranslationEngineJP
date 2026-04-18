---
name: review-ui-check
description: GitHub Copilot 側の UI check 知識 package。
---

# Review UI Check

## 目的

この skill は知識 package である。
`reviewer` agent が UI check を行う時に、主要導線、画面状態、console error、screenshot / snapshot の証跡を確認する判断基準を提供する。

## いつ参照するか

- frontend を含む実装の UI check を行う時
- Playwright MCP で browser surface を確認する時
- ui_evidence を reroute / pass 判断へ使う時

## 参照しない場合

- implementation review だけで足りる時
- design mock を作る時
- paid real AI API を呼ぶ危険がある時

## 原則

- Playwright MCP は `http://host.docker.internal:34115` を使う
- paid な real AI API を呼ばない
- design にない改善提案はしない
- UI 逸脱は reroute で返す

## DO / DON'T

DO:
- 主要導線、画面状態、console error を確認する
- 必要な screenshot または snapshot を evidence にする
- fake provider や test mode の安全性を確認する

DON'T:
- design review をしない
- UI check 中に修正しない
- active contract をこの skill に置かない

## Checklist

- [review-ui-check-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/review-ui-check/references/checklists/review-ui-check-checklist.md) を参照する。
