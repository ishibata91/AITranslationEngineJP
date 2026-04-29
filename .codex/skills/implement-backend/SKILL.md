---
name: implement-backend
description: Codex implementation レーン 側の backend 実装作業プロトコル。層 責務、レーン内検証 の判断基準を提供する。
---
# Implement Backend

## 目的

この skill は作業プロトコルである。
`implementation_implementer` agent が backend 承認済み実装範囲 を実装する時に、usecase、service、repository、adapter の責務整合と 依存方向 を守る判断基準を提供する。

## 対応ロール

- `implementation_implementer` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- 担当成果物は `implement-backend` の出力規約で固定する。

## 入力規約

- backend package を変更する時
- レーン内検証 と エラー経路 を実装する時
- 入力に 根拠参照、担当者、承認状態が不足する場合は推測で補わない。

## 外部参照規約

- エージェント実行定義とツール権限は [implementation_implementer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_implementer.toml) の 書き込み許可 / 実行許可 とする。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

## 判断規約

- 層 責務と依存方向を守る
- エラー経路 と 検証 を 承認済み実装範囲 内で閉じる
- 単一引き継ぎ入力 と 承認済み実装範囲 を確認して プロダクトコード だけを変更する
- `APIテスト` 先行時だけ implementation_tester 出力 も確認する
- レーン内検証 結果 または未実行理由を返す
- `lint:backend` の format、vet、static、arch、module で落ちる境界違反を事前に避ける

- usecase / service / repository / adapter の責務を確認する
- [lint-policy.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/lint-policy.md) の backend lint 内訳を確認する
- [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) の backend 依存方向に従い、usecase、service、repository、adapter concrete の境界を跨がない
- usecase から repository concrete、実行定義 concrete、driver API を直接参照しない
- レーン内検証 を実行した場合は結果を 終了処理 に残す

## 出力規約

- 出力は判断結果、根拠参照、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。

## 完了規約

- 承認済み実装範囲 内の成果だけが返却されている。
- 検証、未実行項目、残留リスク が 根拠参照 付きで整理されている。
- 層 責務と 依存方向 を確認した。
- backend lint の format、static、arch、module 観点を確認した。
- 検証 と エラー経路 を 承認済み実装範囲 内で確認した。
- 単一引き継ぎ入力 と レーン内検証 を確認した。
- `APIテスト` 先行時だけ implementation_tester 出力 を確認した。

## 停止規約

- frontend だけの変更を実装する時
- UI check を行う時
- backend 境界を設計し直す時
- 承認済み実装範囲 外の 層 refactor を混ぜない
- controller、usecase、service で concrete 実装を new しない
- service core から filesystem、Wails 実行定義、DB driver の concrete API を直接呼ばない
- プロダクトテスト、検証データ、スナップショット、test helper を変更しない
- docs や 作業流れ 文書を変更しない
- 停止時は不足項目、衝突箇所、戻し先を返す。
- 承認済み実装範囲 外の 層 refactor を混ぜなかった場合は停止する。
- usecase、service、controller で concrete 実装を new しなかった場合は停止する。
- プロダクトテスト、検証データ、スナップショット、test helper を変更しなかった場合は停止する。
- docs / 作業流れ 文書を変更しなかった場合は停止する。
