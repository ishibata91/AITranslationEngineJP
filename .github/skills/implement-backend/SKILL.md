---
name: implement-backend
description: GitHub Copilot 側の backend 実装知識 package。layer 責務、lane-local validation の判断基準を提供する。
---

# Implement Backend

## 目的

この skill は知識 package である。
`implementer` agent が backend owned_scope を実装する時に、usecase、service、repository、adapter の責務整合と dependency direction を守る判断基準を提供する。

## いつ参照するか

- backend package を変更する時
- lane-local validation と error path を実装する時

## 参照しない場合

- frontend だけの変更を実装する時
- UI check を行う時
- backend 境界を設計し直す時

## 原則

- layer 責務と依存方向を守る
- error path と validation を owned_scope 内で閉じる
- lane_context_packet を確認して product code だけを変更する
- scenario 先行時だけ tester output も確認する
- lane-local validation result または未実行理由を返す

## DO / DON'T

DO:
- usecase / service / repository / adapter の責務を確認する
- lane-local validation を実行した場合は結果を closeout に残す

DON'T:
- owned_scope 外の layer refactor を混ぜない
- product test、fixture、snapshot、test helper を変更しない
- docs や workflow 文書を変更しない
- active contract をこの skill に置かない

## Checklist

- [implement-backend-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implement-backend/references/checklists/implement-backend-checklist.md) を参照する。
