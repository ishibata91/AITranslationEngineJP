---
name: tests
description: Codex implementation レーン 側の プロダクトテスト 共通作業プロトコル。承認済み実装範囲 を test で証明する判断基準を提供する。
---
# Tests

## 目的

`tests` は作業プロトコルである。
`implementation_tester` agent が、単一引き継ぎ入力 と 承認済み実装範囲 を プロダクトテスト で証明する時の共通判断を提供する。

ツール権限 は [implementation_tester.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_tester.toml) が持ち、引き継ぎ は skill に従う。

## 対応ロール

- `implementation_tester` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- 担当成果物は `tests` の出力規約で固定する。

## 入力規約

- 承認済み 引き継ぎ または実装済み 対象範囲 を プロダクトテスト で証明する時
- シナリオ 成果物 または 単体 responsibility を test に落とす時
- fake provider や DI seam で paid real AI API を避ける時
- 入力に 根拠参照、担当者、承認状態が不足する場合は推測で補わない。
- 必須入力: 単一引き継ぎ入力, 承認記録, 承認済み実装範囲, テスト対象, 検証コマンド
- 任意入力: test 小範囲, 参照ヒント, シナリオ成果物, 実装済み範囲, 再現根拠, 既存テスト参照
- 入力注記: {"単一引き継ぎ入力": "implementation-scope から抽出済みの 引き継ぎ 1 件だけ。implementation-scope 全文、進行中作業計画 全文、根拠成果物、後続 引き継ぎ は入力に含めない。", "test 小範囲": "implement_lane が 文脈 枯渇時に同一 引き継ぎ 内で狭めた implementation_tester 用 小範囲。完了条件 clause、公開接点 / API 境界、test 対象 file、検証コマンド のいずれか 1 軸で切られる。", "参照ヒント": "tests-scenario または tests-unit の参照ヒント。共通規約と完了条件は変えない。"}
- 文脈不足基準: {"判定条件": "構造判定条件", "文脈不足として返す条件": ["単一引き継ぎ入力 に 証明対象の挙動、公開接点、test 対象、検証主眼、fixture/helper 方針、対象限定検証 のいずれかが欠けている", "`UI人間操作E2E` で開始操作、検証対象の入口、入力模倣方針のいずれかが欠けている", "`APIテスト` で 公開接点、要求 / 応答契約、入力開始点、主要観測点のいずれかが欠けている", "test 小範囲 が 完了条件 clause、公開接点 / API 境界、test 対象 file、検証コマンド のいずれにも対応していない", "test 作成に 承認済み実装範囲 外探索、プロダクトコード 変更、paid API 呼び出しが必要になる"], "文脈不足としない条件": ["承認済み シナリオ を元に期待どおり fail する プロダクトテスト を追加できる", "局所的な import 修正または既存 test file 内の軽微な確認だけで進められる", "単一引き継ぎ入力 の 列挙済み file / symbol 内で テスト接点 を確認できる"], "該当時の必須返却": ["reason", "必要文脈", "テスト済み小範囲", "残り test 小範囲"]}

## 外部参照規約

- エージェント実行定義とツール権限は [implementation_tester.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_tester.toml) の 書き込み許可 / 実行許可 とする。
- エージェント実行定義: [implementation_tester.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_tester.toml)
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/tests/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/tests-scenario/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/tests-unit/SKILL.md

## 内部参照規約

### 拘束観点

- test の責務分割
- 単一引き継ぎ入力 の読み順
- 文脈不足 の返し方
- deterministic setup
- 最終検証 レーン へ defer する coverage / harness all の扱い
- Arrange / Act / Assert の読みやすさ
- 重点 skill の選び方

- 参照 型 は [test-patterns.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/tests/references/patterns/test-patterns.md) とする。

## 判断規約

- 各 test は 1 つの振る舞いを扱う
- 単一引き継ぎ入力 の完了条件、公開接点、test 対象、検証コマンド の順に読む
- 列挙済み file / symbol 外を探索して 文脈 不足を埋めない
- 文脈不足基準 は 構造判定条件 とし、証明対象の挙動、公開接点、test 対象、検証主眼、fixture/helper 方針、対象限定検証 の不足時に返す
- `UI人間操作E2E` で開始操作、検証対象の入口、入力模倣方針が不足する場合は 文脈不足 を返す
- `APIテスト` で 公開接点、要求 / 応答契約、入力開始点、主要観測点が不足する場合は 文脈不足 を返す
- test 小範囲 が 完了条件 clause、公開接点、test 対象 file、検証コマンド のいずれにも対応しない場合は 文脈不足 を返す
- 承認済み シナリオ を元に期待どおり fail する test、局所的 import 修正、既存 test file 内の軽微な確認は not_insufficient_context として扱う
- 原因未確定の 回帰 test は実装前に書かない
- setup は決定的にする
- test 追加または更新後、引き継ぎ を終える前に touched 層 に対応する local 検証 を実行する
- paid real AI API を呼ばない
- 新しい要件解釈を足さない

- 単一引き継ぎ入力 の 完了条件 clause、証明対象の挙動、公開接点、assertion_focus に沿って test を作る
- `UI人間操作E2E` では、承認済みシナリオの開始操作と入力模倣方針に沿って試験を作る
- `APIテスト` では、承認済み受け入れ条件、公開接点、要求 / 応答契約、入力開始点、主要観測点に沿って試験を作る
- test 小範囲 が渡された場合はその 小範囲 だけを証明し、残りを 残り test 小範囲 に残す
- 文脈不足 を返す場合は reason、必要文脈、残り test 小範囲 を 構造判定条件 に対応づける
- 検証データ の入力値、cロック、実行定義 応答、seed を固定する
- test body に条件分岐を入れない
- coverage、harness all、repo-local Sonar issue 判定条件 は 最終検証 レーン へ defer する
- backend 引き継ぎ は `python3 scripts/harness/run.py --suite backend-local`、frontend 引き継ぎ は `python3 scripts/harness/run.py --suite frontend-local` を使う
- mixed 引き継ぎ は touched 層 に応じて両方を実行する
- 検証結果 と 残り 不足 を返す
- active 規約 は agent 1:1。シナリオ / 単体 の差分は 重点 skill で扱い、出力 obligation はこの 規約 に固定する。

## 非対象規約

- プロダクトコード恒久修正、design / シナリオ成果物の新規作成、レビューだけの作業は扱わない。
- docs、`.codex`、`.codex/skills`、`.codex/agents` は変更しない。
- paid real AI API は呼ばない。
- 新しい要件解釈や mode 別 個別 JSON 規約は追加しない。
- UI 入口の `UI人間操作E2E` を裏側の直接呼び出しだけで代替しない。

## 出力規約

- 出力は判断結果、根拠参照、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。
- 返却先: implement_lane
- 引き継ぎ 資料のスコープ粒度に対応する プロダクトテスト と必要最小限の 検証データ / helper だけを返す
- 単一引き継ぎ入力 1 件と test 小範囲 の観点と対応づけ、複数 引き継ぎ を束ねない
- 実際に test で証明した 完了条件 clause、公開接点 / API 境界、test 対象 file、検証コマンド を返す。test 小範囲 が入力された場合はそれに対応づける
- 同じ 引き継ぎ 内で未証明の implementation_tester 小範囲 を返す。完了条件は削らず、未処理分を明示する
- test 追加または更新後、引き継ぎ を終える前に touched 層 に対応する local 検証 結果を返す。backend は `python3 scripts/harness/run.py --suite backend-local`、frontend は `python3 scripts/harness/run.py --suite frontend-local`、mixed は touched 層 に応じて両方を実行する。未実行なら 阻害理由 を返す
- 文脈不足基準 の 構造判定条件 に一致する場合だけ true とし、reason、必要文脈、テスト済み小範囲、残り test 小範囲 を返す。自力で広く調査して埋めない。判定基準 に一致しない不安や通常の局所確認では true にしない。問題がなければ false または なし
- 文脈不足 true 時は 文脈不足基準 のどの 構造判定条件 に一致したかを返す。false 時は なし または未使用にする
- 文脈不足 時に不足している 公開接点、要求 / 応答契約、existing test 対象、fixture/helper、検証主眼、検証コマンド、`UI人間操作E2E` の開始操作、検証対象の入口、入力模倣方針、`APIテスト` の入力開始点、主要観測点を列挙する
- 最終検証 レーン に defer する。実行した場合だけ結果を返し、未実行なら 最終検証延期 と理由を返す
- 最終検証 レーン に defer する。実行した場合だけ結果を返し、未実行なら 最終検証延期 と理由を返す
- 未証明の振る舞い、未実行 検証、阻害理由 を分ける

## 完了規約

- 承認済み実装範囲 内の成果だけが返却されている。
- 検証、未実行項目、残留リスク が 根拠参照 付きで整理されている。
- implemented 対象範囲 と 承認済み実装範囲 を確認した。
- 単一引き継ぎ入力 の完了条件、公開接点、test 対象、検証コマンド を確認した。
- `UI人間操作E2E` では開始操作、検証対象の入口、入力模倣方針を確認した。
- `APIテスト` では 公開接点、要求 / 応答契約、入力開始点、主要観測点を確認した。
- test 小範囲 がある場合はその範囲だけを証明した。
- 文脈不足基準 の 構造判定条件 に一致する場合だけ 文脈不足、必要文脈、残り test 小範囲 を返した。
- not_insufficient_context に該当する局所確認や承認済み シナリオ 由来の fail test を停止理由にしなかった。
- 原因未確定の 回帰 test を実装前に書かなかった。
- deterministic setup にした。
- 重点 skill の知識だけを追加で参照した。
- 必須 根拠: 単一引き継ぎ入力 id, 引き継ぎ対象範囲 granularity, 承認済み実装範囲, テスト対象, テスト済み小範囲 or 文脈不足 reason, 変更層の局所 検証結果 or 阻害理由
- 完了判断材料: implement_lane が implementation_implementer、レビュー、戻し の次 action を 変更層の局所 検証 結果込みで判断できる
- 残留リスク: remaining_gaps

## 停止規約

- プロダクトコードの恒久修正が主目的の時
- design や シナリオ 成果物 を新規に作る時
- レビュー だけを行う時
- 文脈不足を返さず広く調査する必要がある場合は停止する。
- 判定基準 mismatch になる不安や通常の局所確認を文脈不足扱いにする必要がある場合は停止する。
- プロダクトコード を広く直さない
- 停止時は不足項目、衝突箇所、戻し先を返す。
- 拒否条件: 不足 単一引き継ぎ入力
- 拒否条件: 不足 テスト対象
- 拒否条件: 不足 承認済み実装範囲
- 拒否条件: 不足 必須 文脈 for UI人間操作E2E or APIテスト
- 拒否条件: paid real AI API リスク
- 停止条件: 設計や ownership の整理が先に必要である
- 停止条件: プロダクトコードの広い変更が必要になる
- 停止条件: paid real AI API を呼ぶ危険がある
- 規約違反条件: 文脈不足 を返さず広く調査する
- 規約違反条件: 判定基準 mismatch: 文脈不足基準 外の理由で 文脈不足 を返す
- 規約違反条件: 原因未確定の 回帰 test を実装前に書く
- 規約違反条件: 残り test 小範囲 を隠して完了扱いにする
- 規約違反条件: test 追加または更新後に touched 層 の local 検証 結果または未実行理由を返さない
