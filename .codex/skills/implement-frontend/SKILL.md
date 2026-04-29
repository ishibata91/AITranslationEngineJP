---
name: implement-frontend
description: Codex implementation lane 側の frontend 実装作業プロトコル。画面導線、state、Wails bridge の判断基準を提供する。
---
# Implement Frontend

## 目的

この skill は作業プロトコルである。
`implementation_implementer` agent が frontend owned_scope を実装する時に、画面導線、state 反映、Wails bridge 呼び出しを守る判断基準を提供する。

## 対応ロール

- `implementation_implementer` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- owner artifact は `implement-frontend` の出力規約で固定する。

## 入力規約

- frontend component、state、gateway を変更する時
- console error を出さないことを確認する時
- UI check 前提の build / run 状態を整える時
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。

## 外部参照規約

- agent runtime と tool policy は [implementation_implementer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_implementer.toml) の `allowed_write_paths` / `allowed_commands` とする。
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

## 判断規約

- 画面導線と state 反映を implementation-scope に合わせる
- Wails bridge 呼び出しの境界を守る
- generated `wailsjs` は gateway 境界に閉じ込める
- affected UI の manual flow を確認できる状態にする
- UI check に必要な evidence を残す
- lane_context_packet を確認して product code だけを変更する
- `APIテスト` 先行時だけ implementation_tester output も確認する

- lane_context_packet、affected UI flow を確認する
- `APIテスト` 先行時だけ implementation_tester output を確認する
- console error の有無を closeout に残す
- UI state の初期値と更新条件を確認する
- [lint-policy.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/lint-policy.md) の frontend lint 内訳を確認し、`npm run lint` と `format:check` で拾われる観点を先に意識する
- [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) の frontend 境界に従い、View、ScreenController、Frontend UseCase、Gateway の責務を跨がない
- generated `wailsjs` と backend DTO の import は `frontend/src/controller/wails/` に閉じ込める

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力に tool policy、agent runtime、product code の変更義務を含めない。

## 完了規約

- 承認済み owned_scope 内の成果だけが返却されている。
- validation、未実行項目、residual risk が source_ref 付きで整理されている。
- 画面導線と state 反映を確認した。
- Wails bridge 境界を確認した。
- generated `wailsjs` を gateway 境界に閉じ込めた。
- affected UI flow と console error を確認した。
- frontend lint と format:check で拾われる境界違反を確認した。

## 停止規約

- backend だけの変更を実装する時
- design mock を作る時
- UI check だけを行う時
- design にない改善を足さない
- product test、fixture、snapshot、test helper を変更しない
- transport boundary を迂回しない
- View、ScreenController、Frontend UseCase から generated `wailsjs` を直接 import しない
- gateway 以外で backend DTO 変換をしない
- 停止時は不足項目、衝突箇所、reroute 先を返す。
- design にない改善を足さなかった場合は停止する。
- transport boundary を迂回しなかった場合は停止する。
- View、ScreenController、Frontend UseCase から generated `wailsjs` を直接 import しなかった場合は停止する。
- UI check に必要な evidence を残した。
