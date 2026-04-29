---
name: implement
description: Codex implementation レーン 側の プロダクトコード 実装の共通作業プロトコル。承認済み実装範囲 を実装する判断基準を提供する。
---
# Implement

## 目的

`implement` は作業プロトコルである。
`implementation_implementer` agent が、承認済み `implementation-scope` の 引き継ぎ 1 件を 承認済み実装範囲 内へ実装する時の共通判断を提供する。

ツール権限 は [implementation_implementer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_implementer.toml) が持ち、引き継ぎ は skill に従う。

## 対応ロール

- `implementation_implementer` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- 担当成果物は `implement` の出力規約で固定する。

## 入力規約

- 不足時の扱い: 入力に 根拠参照、担当者、承認状態が不足する場合は推測で補わない。
- 必須入力: 単一引き継ぎ入力, 承認記録, 実装対象, 承認済み実装範囲, depends_on_resolved
- 任意入力: 実装小範囲, 参照ヒント, レーン内検証コマンド, implementation_tester 出力
- 入力注記: {"単一引き継ぎ入力": "implementation-scope から抽出済みの 引き継ぎ 1 件だけ。implementation-scope 全文、進行中作業計画 全文、根拠成果物、後続 引き継ぎ は入力に含めない。", "実装小範囲": "implement_lane が 文脈 枯渇時に同一 引き継ぎ 内で狭めた implementation_implementer 用 小範囲。完了条件 clause、公開接点 / API 境界、変更対象 / symbol、検証コマンド のいずれか 1 軸で切られる。完了条件 を削るものではない。", "implementation_tester 出力": "`APIテスト` 先行 引き継ぎ で implementation_tester が先に返した プロダクトテスト 結果。通常、単体、原因未確定の 回帰 引き継ぎ では入力に含めない。", "参照ヒント": "implement-backend、implement-frontend、implement-mixed、implement-fix-lane の参照ヒント。共通規約と完了条件は変えない。implement-mixed は API / Wails / DTO / gateway など接合点 対象範囲 に限定する。"}
- 文脈不足基準: {"判定条件": "構造判定条件", "文脈不足として返す条件": ["単一引き継ぎ入力 に完了条件、公開接点 / API 境界、実装対象、承認済み実装範囲、検証コマンド のいずれかが欠けている", "実装対象 が file / symbol / 公開接点 のいずれにも対応していない", "実装に 承認済み実装範囲 拡張、プロダクトテスト / 検証データ / スナップショット / test helper 変更、docs / 作業流れ 変更、新規設計判断、広域 refactor が必要になる"], "文脈不足としない条件": ["単一引き継ぎ入力 内の局所確認だけで first edit に入れる", "既存 型 への通常追従で実装できる", "レーン内検証 失敗 を プロダクトコード 内の 対象範囲 で修正できる"], "該当時の必須返却": ["reason", "必要文脈", "推奨分割軸", "実装済み小範囲", "remaining_実装小範囲s"]}

## 外部参照規約

- エージェント実行定義とツール権限は [implementation_implementer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_implementer.toml) の 書き込み許可 / 実行許可 とする。
- エージェント実行定義: [implementation_implementer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_implementer.toml)
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-backend/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-frontend/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-mixed/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-fix-lane/SKILL.md

## 内部参照規約

### 拘束観点

- 承認済み実装範囲 を超えない実装判断
- 引き継ぎ 資料のスコープ粒度に合わせる判断
- coding guidelines と既存 型 の確認
- lint policy と architecture constraint の局所確認
- 境界、エラー経路、test 表面 の実装品質判断
- 検証結果 と 残留リスク の返し方
- 重点 skill の選び方

- 参照 型 は [implementation-quality-patterns.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement/references/patterns/implementation-quality-patterns.md) とする。

## 判断規約

- `implementation-scope` と 承認済み実装範囲 を超えない
- 引き継ぎ 資料のスコープ粒度で実装する
- 単一引き継ぎ入力 と 実装対象 に合わせて プロダクトコード だけを変更する
- `APIテスト` 先行時だけ implementation_tester 出力 も確認する
- 実装小範囲 が渡された場合はその 小範囲 内だけを実装する
- 実装完了後、引き継ぎ を終える前に touched 層 に対応する local 検証 を実行する
- 実装対象 に対応する code path を優先し、承認済み実装範囲 外へ寄り道しない
- 単一引き継ぎ入力 の完了条件、公開接点、検証コマンド から着手する
- 文脈不足基準 は 構造判定条件 とし、完了条件、公開接点、実装対象、承認済み実装範囲、検証コマンド の不足時に返す
- 実装対象 が file / symbol / 公開接点 に対応していない場合は 文脈不足 を返す
- 単一引き継ぎ入力 内の局所確認、既存 型 への通常追従、レーン内検証 失敗 は not_insufficient_context として扱う
- 既存 型、naming、層 に合わせる
- 広域 refactor を混ぜない
- プロダクトテスト、検証データ、スナップショット、test helper は implementation_tester が扱う
- docs 正本化をしない

- 実装前に [coding-guidelines.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/coding-guidelines.md) を読む
- 実装前に [lint-policy.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/lint-policy.md) を読み、引き継ぎ に効く静的 check の責務を確認する
- 単一引き継ぎ入力 の完了条件、承認済み実装範囲、実装対象、関連 根拠参照、検証コマンド を確認する
- 引き継ぎ に architecture constraint がある場合は、その範囲だけ [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) を局所確認する
- 実装小範囲 があれば 完了条件 clause、公開接点、変更対象 / symbol、検証コマンド を確認する
- 文脈不足 を返す場合は reason、必要文脈、推奨分割軸、remaining_実装小範囲s を 構造判定条件 に対応づける
- 入口、呼び出し箇所、データ流れ、エラー経路、test 表面 を確認する
- 既存 型 に naming、constructor、DI、エラー return を合わせる
- generated import、層 依存、境界 rule、format 逸脱など、touched 層 で踏みやすい lint 観点を先に確認する
- レーン内検証 結果または未実行理由を返す
- backend 引き継ぎ は `python3 scripts/harness/run.py --suite backend-local`、frontend 引き継ぎ は `python3 scripts/harness/run.py --suite frontend-local` を使う
- mixed 引き継ぎ は touched 層 に応じて両方を実行する
- touched files は プロダクトコード だけにする
- active 規約 は agent 1:1。backend / frontend / mixed / fix-lane の差分は 重点 skill で扱い、出力 obligation はこの 規約 に固定する。implementation_implementer は承認済み 引き継ぎ 1 件の プロダクトコード 実装を扱い、`APIテスト` 先行時だけ implementation_tester 出力 を受け取る。プロダクトテスト / 検証データ / スナップショット / test helper は変更しない。

## 非対象規約

- UI check、implementation レビュー、要件追加、設計追加は扱わない。
- docs、`.codex`、`.codex/skills`、`.codex/agents` は変更しない。
- プロダクトテスト、検証データ、スナップショット、test helper は変更しない。
- 承認済み実装範囲外の cleanup、rename、format、広域 refactor は扱わない。
- coverage、harness all、repo-local Sonar issue 判定条件は必須終了処理にしない。

## 出力規約

- 基本出力: 出力は判断結果、根拠参照、不足情報、次 agent が判断できる材料を含む。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。
- 返却先: implement_lane
- 実装成果物: 単一引き継ぎ入力 の 承認済み実装範囲 に対応する プロダクトコード だけを返す。プロダクトテスト、検証データ、スナップショット、test helper は含めない
- 引き継ぎ対応: 単一引き継ぎ入力 1 件と 実装小範囲 に対応づけ、複数 引き継ぎ を束ねない
- 実装済み完了条件: 実際に実装した 完了条件 clause、公開接点 / API 境界、変更対象 / symbol、検証コマンド を返す。実装小範囲 が入力された場合はそれに対応づける
- 未実装小範囲: 同じ 引き継ぎ 内で未実装の 小範囲 を返す。完了条件は削らず、未処理分を明示する
- レーン内検証結果: 実装完了後、引き継ぎ を終える前に touched 層 に対応する local 検証 結果を返す。backend は `python3 scripts/harness/run.py --suite backend-local`、frontend は `python3 scripts/harness/run.py --suite frontend-local`、mixed は touched 層 に応じて両方を実行する。未実行なら 阻害理由 を返す。coverage、Sonar、harness all は implementation_implementer の必須 終了処理 にしない
- 実装根拠: 入口、呼び出し箇所、データ流れ、エラー経路、test 表面、既存 型 への整合を簡潔に返す。mixed の場合は接合点 契約 を明記する
- 文脈不足判定: 文脈不足基準 の 構造判定条件 に一致する場合だけ true とし、reason、必要文脈、推奨分割軸、実装済み小範囲、remaining_実装小範囲s を返す。自力で広く調査して埋めない。判定基準 に一致しない不安、通常の局所確認、レーン内検証 失敗 だけでは true にしない。問題がなければ false または なし
- 文脈不足該当条件: 文脈不足 true 時は 文脈不足基準 のどの 構造判定条件 に一致したかを返す。false 時は なし または未使用にする
- 不足文脈: 文脈不足 時に不足している完了条件、実装対象、公開接点、承認済み実装範囲、existing 型、検証コマンド を列挙する
- 推奨分割軸: 文脈不足 時に orchestrator が次に狭めるべき軸を 完了条件 clause、公開接点 / API 境界、変更対象 / symbol、検証コマンド のいずれかで返す
- 阻害理由: 未実行 検証、対象範囲 超過、設計不足、test / 検証データ 変更が必要になった場合の 阻害理由 を分ける

## 完了規約

- 承認済み実装範囲 内の成果だけが返却されている。
- 検証、未実行項目、残留リスク が 根拠参照 付きで整理されている。
- 承認済み実装範囲 と implementation 対象 を確認した。
- 単一引き継ぎ入力 を確認した。
- `APIテスト` 先行時だけ implementation_tester 出力 を確認した。
- 実装小範囲 がある場合はその範囲だけを実装した。
- 文脈不足基準 の 構造判定条件 に一致する場合だけ 文脈不足、必要文脈、推奨分割軸 を返した。
- not_insufficient_context に該当する局所確認、既存 型 追従、レーン内検証 失敗 を停止理由にしなかった。
- 実装対象 と 公開接点 から着手した。
- coding guidelines、lint policy、レーン内検証 commands を確認した。
- 引き継ぎ にある architecture constraint を局所確認した。
- 重点 skill の知識だけを追加で参照した。
- touched files が プロダクトコード だけであることを確認した。
- 必須 根拠: 単一引き継ぎ入力 id, 実装対象, 承認済み実装範囲, 承認記録, implementation_tester 出力 APIテスト先行実装テストがある場合, 実装済み小範囲 or 文脈不足 reason, 入口, 呼び出し箇所, データ流れ または境界, エラー経路, test 表面, 変更層の局所 検証結果 or 阻害理由
- 完了判断材料: implement_lane が レビュー へ進める プロダクトコード 実装結果と 変更層の局所 検証 結果が返っている
- 残留リスク: residual_risks

## 停止規約

- UI check や implementation レビュー を行う時
- docs や 作業流れ 文書を変更する時
- 文脈不足 を返さず広い調査で不足 文脈 を埋めない
- 判定基準 mismatch になる不安や通常の局所確認を文脈不足扱いにする必要がある場合は停止する。
- 実装小範囲 外へ実装を広げない
- 実装対象 がないまま広い調査を始めない
- lint 未確認のまま実装して local 検証で初めて境界違反を知る進め方になる場合は停止する。
- config、lint、test、coverage 設定変更で判定条件を回避する必要がある場合は停止する。
- 停止時は不足項目、衝突箇所、戻し先を返す。
- 拒否条件: 不足 単一引き継ぎ入力
- 拒否条件: 不足 実装対象
- 拒否条件: 不足 承認記録
- 拒否条件: 不足 implementation_tester 出力 for API test pre-implementation 引き継ぎ
- 拒否条件: unclear 承認済み実装範囲
- 停止条件: 設計判断が不足している
- 停止条件: docs または 作業流れ の変更が必要になる
- 停止条件: 広域 refactor なしでは実装できない
- 停止条件: プロダクトテスト、検証データ、スナップショット、test helper の変更が必要になる
- 規約違反条件: touched_test_files を返す
- 規約違反条件: プロダクトテスト、検証データ、スナップショット、test helper を変更する
- 規約違反条件: implementation-scope 全文 または後続 引き継ぎ を入力として要求する
- 規約違反条件: 文脈不足 を返さず広く調査して不足 文脈 を埋める
- 規約違反条件: 判定基準 mismatch: 文脈不足基準 外の理由で 文脈不足 を返す
- 規約違反条件: remaining_実装小範囲s を隠して完了扱いにする
- 規約違反条件: 実装完了後に touched 層 の local 検証 結果または未実行理由を返さない
