# Implementation Quality Patterns

## 目的

`implementer` が owned_scope 内で product code を実装するための品質判断パターンをまとめる。
agent contract の権限や output obligation は上書きしない。

## 採用する考え方

- 読みやすさを優先し、clever code より明示的な構造を選ぶ。
- KISS、DRY、YAGNI を守り、必要になった時だけ抽象化する。
- error path、empty state、boundary value を実装時に明示する。
- build / type error の解消は最小差分にする。
- 変更前に lane_context_packet の fix_ingredients、distracting_context、first_action、change_targets、既存の naming、layer、dependency direction を確認する。
- scenario 先行時だけ tester output も確認する。
- 大きい関数、深いネスト、magic number、silent fallback を赤旗として扱う。
- 振る舞いを変えない整理は、可読性が明確に上がる場合だけ行う。

## Backend 適用

- service / usecase / repository / infra adapter の責務を混ぜない。
- usecase は orchestration、port usage、business invariant を担当し、SDK / DB / file system 詳細を持たない。
- outbound dependency は consuming side の小さい interface / port で受ける。
- validation、error mapping、transaction boundary を実装の一部として扱う。
- N+1、unbounded query、timeout なしの外部呼び出しを避ける。
- internal error を user-facing response へ漏らさない。
- backend 変更では lane-local validation result または未実行理由を返す。

## Frontend 適用

- state update は明示的にし、stale closure と直接 mutation を避ける。
- loading、error、empty、success の状態を implementation-scope に沿って扱う。
- Wails gateway を迂回せず、transport boundary を守る。
- component は表示、screen controller は状態遷移、gateway は Wails binding 境界に責務を分ける。
- UI check に必要な stable selector、visible state、console evidence を残す。
- frontend 変更では lane-local validation result または未実行理由を返す。

## 品質赤旗

- owned_scope 外の cleanup、rename、format が混ざっている。
- broad refactor なしでは説明できない差分になっている。
- validation failure を握りつぶす fallback がある。
- product test、fixture、snapshot、test helper を implementer が変更している。
- lane-local validation の失敗または未実行理由がない。
- public API、Wails binding、DTO、storage schema の変更が call site と test に反映されていない。
- config、lint、test、coverage 設定を変更して gate を回避している。

## 実装前確認

- handoff 資料のスコープ粒度と owned_scope を確認する。
- lane_context_packet を確認する。
- scenario 先行時だけ tester output を確認する。
- fix_ingredients に対応する path、symbol、line number を確認する。
- distracting_context に挙がった周辺 context を実装対象から外す。
- first_action の path、symbol、line number から着手する。
- 入口、call site、data flow、error path、test surface を確認する。
- 既存の似た実装を探し、naming、constructor、DI、error return の形を合わせる。
- 追加する抽象化が既存 pattern と一致するか確認する。
- 変更後に必要な lane-local validation command を先に確認する。

## 実装中の品質基準

- 1 関数は 1 つの責務に絞り、深いネストは early return や named helper で浅くする。
- 複雑な条件式は意味名を持つ変数または関数に分ける。
- duplicated logic は同じ handoff scope 内でだけ統合し、広域共通化へ広げない。
- dead code、commented-out code、stray debug output を残さない。
- error は握りつぶさず、domain / application / adapter 境界で意味のある形へ変換する。
- TODO で未完成を隠さず、completion packet の residual risk に残す。

## 境界チェック

- domain / usecase は framework、Wails、SDK、storage concrete を import しない。
- adapter は protocol / storage / external API 変換を担当し、business rule を持たない。
- composition / bootstrap は wiring を担当し、business rule を持たない。
- frontend は Wails binding を gateway 境界に閉じ込め、component から直接呼ばない。
- mixed scope は API、DTO、Wails binding、gateway、adapter contract の接合点だけに限定する。
