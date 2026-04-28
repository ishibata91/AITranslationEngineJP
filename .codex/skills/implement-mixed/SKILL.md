---
name: implement-mixed
description: Codex implementation lane 側の API / Wails / DTO / gateway など frontend と backend の接合点実装知識 package。
---

# Implement Mixed

## 目的

この skill は知識 package である。
`implementation_implementer` agent が scope freeze 済みの API、Wails binding、DTO、gateway、adapter contract など frontend と backend の接合点 owned_scope を実装する時の判断基準を提供する。

mixed は広い frontend / backend 同時変更の許可ではない。
片側だけで閉じる UI 実装や backend 実装は、それぞれ `implement-frontend` または `implement-backend` を使う。

## いつ参照するか

- implementation-scope が API、Wails binding、DTO、gateway、adapter contract の接合点変更を明示している時
- 片側だけでは contract 整合を証明できない owned_scope を扱う時
- validation を frontend / backend の接合点 evidence として返す時

## 参照しない場合

- frontend または backend の片側だけで閉じる時
- API / Wails / DTO / gateway / adapter contract の接合点を変更しない時
- 横断範囲が未承認の時
- 追加設計で横断 scope を広げる時

## 原則

- implementation-scope の owned_scope を守る
- mixed の対象を API、Wails binding、DTO、gateway、adapter contract の接合点だけに限定する
- 片側だけで閉じない理由を scope artifact で確認する
- lane_context_packet を確認して product code だけを変更する
- `APIテスト` 先行時だけ implementation_tester output も確認する
- validation は frontend、backend、接合点 contract の証跡を分ける

## DO / DON'T

DO:
- API / Wails / DTO / gateway / adapter contract のどれを接合点として変更したか closeout に残す
- 両側の touched files を handoff と対応づける
- frontend / backend / 接合点 contract の lane-local validation evidence を分ける
- lane-local validation command の不足を residual risk にする

DON'T:
- mixed を広い frontend / backend 同時変更の口実にしない
- 片側の都合で scope を広げない
- API 接合点を変えずに UI と backend を同時に触らない
- product test、fixture、snapshot、test helper を変更しない
- docs や workflow 文書を変更しない
- active contract をこの skill に置かない

## Checklist

- [implement-mixed-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-mixed/references/checklists/implement-mixed-checklist.md) を参照する。
