---
name: implement-mixed
description: Codex implementation レーン 側の API / Wails / DTO / gateway など frontend と backend の接合点実装作業プロトコル。
---
# Implement Mixed

## 目的

この skill は作業プロトコルである。
`implementation_implementer` agent が 対象範囲 freeze 済みの API、Wails 紐づけ、DTO、gateway、adapter 契約 など frontend と backend の接合点 承認済み実装範囲 を実装する時の判断基準を提供する。

mixed は広い frontend / backend 同時変更の許可ではない。
片側だけで閉じる UI 実装や backend 実装は、それぞれ `implement-frontend` または `implement-backend` を使う。

## 対応ロール

- `implementation_implementer` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- 担当成果物は `implement-mixed` の出力規約で固定する。

## 入力規約

- 不足時の扱い: 入力に 根拠参照、担当者、承認状態が不足する場合は推測で補わない。

## 外部参照規約

- エージェント実行定義とツール権限は [implementation_implementer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_implementer.toml) の 書き込み許可 / 実行許可 とする。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

## 判断規約

- implementation-scope の 承認済み実装範囲 を守る
- mixed の対象を API、Wails 紐づけ、DTO、gateway、adapter 契約 の接合点だけに限定する
- 片側だけで閉じない理由を 対象範囲 成果物 で確認する
- 単一引き継ぎ入力 と 承認済み実装範囲 を確認して プロダクトコード だけを変更する
- `APIテスト` 先行時だけ implementation_scenario_tester 出力 も確認する
- 検証 は frontend、backend、接合点 契約 の証跡を分ける

- API / Wails / DTO / gateway / adapter 契約 のどれを接合点として変更したか 終了処理 に残す
- 両側の touched files を 引き継ぎ と対応づける
- frontend / backend / 接合点 契約 の レーン内検証 根拠 を分ける
- レーン内検証 コマンド の不足を 残留リスク にする

## 非対象規約

- frontend または backend の片側だけで閉じる変更は扱わない。
- mixed を広い frontend / backend 同時変更の口実にしない。
- 承認済み接合点外の API / Wails / DTO / gateway / adapter 契約変更は扱わない。
- プロダクトテスト、検証データ、スナップショット、test helper は変更しない。
- docs や作業流れ文書は変更しない。

## 出力規約

- 基本出力: 出力は判断結果、根拠参照、不足情報、次 agent が判断できる材料を含む。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。

## 完了規約

- 承認済み実装範囲 内の成果だけが返却されている。
- 検証、未実行項目、残留リスク が 根拠参照 付きで整理されている。
- API / Wails / DTO / gateway / adapter 契約 の接合点 対象範囲 が承認済みであることを確認した。
- 両側の touched files を 引き継ぎ と対応づけた。
- 単一引き継ぎ入力 と レーン内検証 根拠 を分けた。
- `APIテスト` 先行時だけ implementation_scenario_tester 出力 を確認した。

## 停止規約

- frontend または backend の片側だけで閉じる時
- API / Wails / DTO / gateway / adapter 契約 の接合点変更がない時
- 横断範囲が未承認の時
- 追加設計で横断 対象範囲 を広げる時
- API 接合点を変えずに UI と backend を同時に触らない
- 停止時は不足項目、衝突箇所、戻し先を返す。
