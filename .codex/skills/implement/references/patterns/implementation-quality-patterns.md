# Implementation Quality Patterns

## 目的

`implementation_implementer` が 承認済み実装範囲 内で プロダクトコード を実装するための品質判断パターンをまとめる。
agent TOML の ツール権限 と skill の出力規約は上書きしない。

## 採用する考え方

- 読みやすさを優先し、clever code より明示的な構造を選ぶ。
- KISS、DRY、YAGNI を守り、必要になった時だけ抽象化する。
- エラー経路、空状態、境界 value を実装時に明示する。
- build / type エラー の解消は最小差分にする。
- 変更前に 単一引き継ぎ入力 の完了条件、実装対象、承認済み実装範囲、既存の naming、層、依存方向 を確認する。
- 変更前に `docs/lint-policy.md` のうち今回の touched 層 に効く rule を確認する。
- `APIテスト` 先行時だけ implementation_scenario_tester 出力 も確認する。
- 大きい関数、深いネスト、magic number、silent 代替処理 を赤旗として扱う。
- 振る舞いを変えない整理は、可読性が明確に上がる場合だけ行う。

## Backend 適用

- service / usecase / repository / infra adapter の責務を混ぜない。
- usecase から repository concrete や 実行定義 concrete を直接参照しない。
- usecase は 進行管理、port usage、business 不変条件 を担当し、SDK / DB / file system 詳細を持たない。
- outbound 依存 は consuming side の小さい interface / port で受ける。
- 検証、エラー mapping、トランザクション 境界 を実装の一部として扱う。
- N+1、unbounded query、timeout なしの外部呼び出しを避ける。
- internal エラー を user-facing response へ漏らさない。
- backend 変更では レーン内検証 結果 または未実行理由を返す。

## Frontend 適用

- 状態 update は明示的にし、stale closure と直接 mutation を避ける。
- 読み込み中、エラー、空、成功 の状態を implementation-scope に沿って扱う。
- Wails gateway を迂回せず、通信境界 を守る。
- generated `wailsjs` は `frontend/src/controller/wails/` に閉じ込め、View、ScreenController、Frontend UseCase から直接 import しない。
- component は表示、screen controller は状態遷移、gateway は Wails 紐づけ 境界に責務を分ける。
- UI check に必要な stable selector、表示状態、console 根拠 を残す。
- frontend 変更では レーン内検証 結果 または未実行理由を返す。

## 品質赤旗

- 承認済み実装範囲 外の cleanup、rename、format が混ざっている。
- 広域 refactor なしでは説明できない差分になっている。
- 検証 失敗 を握りつぶす 代替処理 がある。
- プロダクトテスト、検証データ、スナップショット、test helper を implementation_implementer が変更している。
- レーン内検証 の失敗または未実行理由がない。
- public API、Wails 紐づけ、DTO、storage schema の変更が 呼び出し箇所 と test に反映されていない。
- config、lint、test、coverage 設定を変更して 判定条件 を回避している。
- format、境界 rule、禁止 import のどれで止まるかを知らずに差分を広げている。

## 実装前確認

- 引き継ぎ 資料のスコープ粒度と 承認済み実装範囲 を確認する。
- 単一引き継ぎ入力 を確認する。
- `APIテスト` 先行時だけ implementation_scenario_tester 出力 を確認する。
- 実装対象 に対応する path、symbol、line number を確認する。
- 承認済み実装範囲 外の周辺 文脈 を実装対象から外す。
- 実装対象 の path、symbol、line number から着手する。
- 入口、呼び出し箇所、データ流れ、エラー経路、test 表面 を確認する。
- 既存の似た実装を探し、naming、constructor、DI、エラー return の形を合わせる。
- 追加する抽象化が既存 型 と一致するか確認する。
- 変更後に必要な レーン内検証 コマンド を先に確認する。

## 実装中の品質基準

- 1 関数は 1 つの責務に絞り、深いネストは early return や named helper で浅くする。
- 複雑な条件式は意味名を持つ変数または関数に分ける。
- duplicated logic は同じ 引き継ぎ対象範囲 内でだけ統合し、広域共通化へ広げない。
- dead code、commented-out code、stray debug 出力 を残さない。
- エラー は握りつぶさず、domain / application / adapter 境界で意味のある形へ変換する。
- TODO で未完成を隠さず、完了報告入力 の 残留リスク に残す。

## 境界チェック

- domain / usecase は framework、Wails、SDK、storage concrete を import しない。
- adapter は protocol / storage / external API 変換を担当し、business rule を持たない。
- composition / bootstrap は wiring を担当し、business rule を持たない。
- frontend は Wails 紐づけ を gateway 境界に閉じ込め、component から直接呼ばない。
- mixed 対象範囲 は API、DTO、Wails 紐づけ、gateway、adapter 契約 の接合点だけに限定する。
