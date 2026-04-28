---
name: implement-backend
description: Codex implementation lane 側の backend 実装知識 package。layer 責務、lane-local validation の判断基準を提供する。
---

# Implement Backend

## 目的

この skill は知識 package である。
`implementation_implementer` agent が backend owned_scope を実装する時に、usecase、service、repository、adapter の責務整合と dependency direction を守る判断基準を提供する。

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
- `APIテスト` 先行時だけ implementation_tester output も確認する
- lane-local validation result または未実行理由を返す
- `lint:backend` の format、vet、static、arch、module で落ちる境界違反を事前に避ける

## DO / DON'T

DO:
- usecase / service / repository / adapter の責務を確認する
- [lint-policy.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/lint-policy.md) の backend lint 内訳を確認する
- [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) の backend 依存方向に従い、usecase、service、repository、adapter concrete の境界を跨がない
- usecase から repository concrete、runtime concrete、driver API を直接参照しない
- lane-local validation を実行した場合は結果を closeout に残す

DON'T:
- owned_scope 外の layer refactor を混ぜない
- controller、usecase、service で concrete 実装を new しない
- service core から filesystem、Wails runtime、DB driver の concrete API を直接呼ばない
- product test、fixture、snapshot、test helper を変更しない
- docs や workflow 文書を変更しない
- active contract をこの skill に置かない

## Checklist

- [implement-backend-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-backend/references/checklists/implement-backend-checklist.md) を参照する。
