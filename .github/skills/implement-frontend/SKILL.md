---
name: implement-frontend
description: GitHub Copilot 側の frontend 実装知識 package。画面導線、state、Wails bridge の判断基準を提供する。
---

# Implement Frontend

## 目的

この skill は知識 package である。
`implementer` agent が frontend owned_scope を実装する時に、画面導線、state 反映、Wails bridge 呼び出しを守る判断基準を提供する。

## いつ参照するか

- frontend component、state、gateway を変更する時
- console error を出さないことを確認する時
- UI check 前提の build / run 状態を整える時

## 参照しない場合

- backend だけの変更を実装する時
- design mock を作る時
- UI check だけを行う時

## 原則

- 画面導線と state 反映を implementation-scope に合わせる
- Wails bridge 呼び出しの境界を守る
- generated `wailsjs` は gateway 境界に閉じ込める
- affected UI の manual flow を確認できる状態にする
- UI check に必要な evidence を残す
- lane_context_packet を確認して product code だけを変更する
- `APIテスト` 先行時だけ tester output も確認する

## DO / DON'T

DO:
- lane_context_packet、affected UI flow を確認する
- `APIテスト` 先行時だけ tester output を確認する
- console error の有無を closeout に残す
- UI state の初期値と更新条件を確認する
- [lint-policy.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/lint-policy.md) の frontend lint 内訳を確認し、`npm run lint` と `format:check` で拾われる観点を先に意識する
- [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) の frontend 境界に従い、View、ScreenController、Frontend UseCase、Gateway の責務を跨がない
- generated `wailsjs` と backend DTO の import は `frontend/src/controller/wails/` に閉じ込める

DON'T:
- design にない改善を足さない
- product test、fixture、snapshot、test helper を変更しない
- transport boundary を迂回しない
- View、ScreenController、Frontend UseCase から generated `wailsjs` を直接 import しない
- gateway 以外で backend DTO 変換をしない
- active contract をこの skill に置かない

## Checklist

- [implement-frontend-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implement-frontend/references/checklists/implement-frontend-checklist.md) を参照する。
