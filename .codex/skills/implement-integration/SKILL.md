---
name: implement-integration
description: Codex implementation レーン 側の API / Wails / DTO / gateway など frontend と backend の統合境界実装作業プロトコル。
---
# Implement Integration

## 目的

この skill は作業プロトコルである。
`implementation_implementer` agent が 対象範囲 freeze 済みの API、Wails 紐づけ、DTO、gateway、adapter 契約 など frontend と backend の統合境界 承認済み実装範囲 を実装する時の判断基準を提供する。

integration は広い frontend / backend 同時変更の許可ではない。
片側だけで閉じる UI 実装や backend 実装は、それぞれ `implement-frontend` または `implement-backend` を使う。

## 対応ロール

- `implementation_implementer` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- 担当成果物は `implement-integration` の出力規約で固定する。

## 入力規約

- 不足時の扱い: 入力に 根拠参照、担当者、承認状態が不足する場合は推測で補わない。
- 単一引き継ぎ入力: implementation-scope から抽出済みの 引き継ぎ 1 件。
- 承認記録: 人間が承認した実装範囲の根拠参照。
- 実装対象: 変更してよい統合境界のファイル、symbol、公開接点。
- 承認済み実装範囲: 実装してよい統合境界のプロダクトコード範囲。
- 依存解消状態: 依存対象が完了済みかを示す状態。
- 任意入力: 実装小範囲、レーン内検証コマンド、implementation_scenario_tester 出力。

## 外部参照規約

- エージェント実行定義とツール権限は [implementation_implementer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_implementer.toml) の 書き込み許可 / 実行許可 とする。
- frontend コーディング規約: frontend 変更がある場合は [coding-guidelines-frontend.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/coding-guidelines-frontend.md) とする。
- backend コーディング規約: backend 変更がある場合は [coding-guidelines-backend.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/coding-guidelines-backend.md) とする。
- lint 規約: [lint-policy.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/lint-policy.md) とする。
- architecture 規約: [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) の統合境界だけを参照する。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

## 判断規約

- implementation-scope の 承認済み実装範囲 を守る
- integration の対象を API、Wails 紐づけ、DTO、gateway、adapter 契約 の統合境界だけに限定する
- 片側だけで閉じない理由を 対象範囲 成果物 で確認する
- 単一引き継ぎ入力 と 承認済み実装範囲 を確認して プロダクトコード だけを変更する
- `APIテスト` 先行時だけ implementation_scenario_tester 出力 も確認する
- 検証 は frontend、backend、統合境界 契約 の証跡を分ける

- API / Wails / DTO / gateway / adapter 契約 のどれを統合境界として変更したか 終了処理 に残す
- 両側の touched files を 引き継ぎ と対応づける
- frontend / backend / 統合境界 契約 の レーン内検証 根拠 を分ける
- レーン内検証 コマンド の不足を 残留リスク にする

## 非対象規約

- frontend または backend の片側だけで閉じる変更は扱わない。
- integration を広い frontend / backend 同時変更の口実にしない。
- 承認済み統合境界外の API / Wails / DTO / gateway / adapter 契約変更は扱わない。
- プロダクトテスト、検証データ、スナップショット、test helper は変更しない。
- docs や作業流れ文書は変更しない。
- coverage、harness all、repo-local Sonar issue 判定条件は必須終了処理にしない。

## 出力規約

- 判断結果: 統合境界プロダクトコード実装の完了、未完了、停止の判定を返す。
- 根拠参照: 実装の根拠にした入力、変更箇所、検証結果を返す。
- 不足情報: 実装を完了できない不足項目を返す。
- 次判断材料: `implement_lane` が次を判断できる材料を返す。
- 実装成果物: 単一引き継ぎ入力 の 承認済み実装範囲 に対応する統合境界プロダクトコードだけを返す。
- レーン内検証結果: backend-local と frontend-local の結果または未実行理由を変更層別に返す。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコード変更の指示を含めない。

## 完了規約

- 承認済み実装範囲 内の成果だけが返却されている。
- 検証、未実行項目、残留リスク が 根拠参照 付きで整理されている。
- 単一引き継ぎ入力、承認記録、実装対象、承認済み実装範囲を確認した。
- API / Wails / DTO / gateway / adapter 契約 の統合境界 対象範囲 が承認済みであることを確認した。
- 両側の touched files を 引き継ぎ と対応づけた。
- 単一引き継ぎ入力 と レーン内検証 根拠 を分けた。
- backend 側の変更がある場合は `python3 scripts/harness/run.py --suite backend-local` を実行し、結果または未実行理由を返した。
- frontend 側の変更がある場合は `python3 scripts/harness/run.py --suite frontend-local` を実行し、結果または未実行理由を返した。
- backend と frontend の両方を含む場合は両方の局所ハーネスを実行し、結果または未実行理由を返した。
- `APIテスト` 先行時だけ implementation_scenario_tester 出力 を確認した。

## 停止規約

- frontend または backend の片側だけで閉じる時
- API / Wails / DTO / gateway / adapter 契約 の統合境界変更がない時
- 横断範囲が未承認の時
- 追加設計で横断 対象範囲 を広げる時
- API 統合境界を変えずに UI と backend を同時に触らない
- 単一引き継ぎ入力、実装対象、承認記録、承認済み実装範囲が不足する場合は停止する。
- プロダクトテスト、検証データ、スナップショット、test helper の変更が必要になる場合は停止する。
- 承認済み実装範囲外へ実装を広げる必要がある場合は停止する。
- 停止時は不足項目、衝突箇所、戻し先を返す。
