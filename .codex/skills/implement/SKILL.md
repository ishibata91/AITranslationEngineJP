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
- 単一引き継ぎ入力: implementation-scope から抽出済みの 引き継ぎ 1 件。
- 承認記録: 人間が承認した実装範囲の根拠参照。
- 実装対象: 変更してよいファイル、symbol、公開接点。
- 承認済み実装範囲: 実装してよいプロダクトコード範囲。
- 依存解消状態: 依存対象が完了済みかを示す状態。
- 任意入力: 実装小範囲、参照ヒント、レーン内検証コマンド、implementation_scenario_tester 出力。
- 文脈不足基準: 完了条件、公開接点 / API 境界、実装対象、承認済み実装範囲、検証コマンド の不足を判定する基準。
- 文脈不足対象: 実装対象がファイル、symbol、公開接点のいずれにも対応せず、承認済み実装範囲の拡張または非対象変更が必要になる状態。
- 文脈不足非対象: 単一引き継ぎ入力内の局所確認、既存型への通常追従、プロダクトコード内で閉じるレーン内検証失敗。
- 文脈不足返却: 理由、必要文脈、推奨分割軸、実装済み小範囲、残り実装小範囲。

## 外部参照規約

- エージェント実行定義とツール権限は [implementation_implementer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_implementer.toml) の 書き込み許可 / 実行許可 とする。
- エージェント実行定義: [implementation_implementer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_implementer.toml)
- コーディング規約: [coding-guidelines.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/coding-guidelines.md) とする。
- lint 規約: [lint-policy.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/lint-policy.md) とする。
- architecture 規約: 引き継ぎに architecture constraint がある場合だけ [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) を参照する。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-backend/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-frontend/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-mixed/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-fix-lane/SKILL.md

## 内部参照規約

### 拘束観点

- 承認済み実装範囲 を超えない実装判断
- 引き継ぎ 資料のスコープ粒度に合わせる判断
- コーディング規約 と既存 型 の確認
- lint 規約 と architecture constraint の局所確認
- 境界、エラー経路、test 表面 の実装品質判断
- 検証結果 と 残留リスク の返し方
- 重点 skill の選び方

- 参照 型 は [implementation-quality-patterns.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement/references/patterns/implementation-quality-patterns.md) とする。

## 判断規約

- `implementation-scope` と 承認済み実装範囲 を超えない
- 引き継ぎ 資料のスコープ粒度で実装する
- 単一引き継ぎ入力 と 実装対象 に合わせて プロダクトコード だけを変更する
- `APIテスト` 先行時だけ implementation_scenario_tester 出力 も確認する
- 実装小範囲 が渡された場合はその 小範囲 内だけを実装する
- 実装完了後、引き継ぎ を終える前に変更層に対応する局所検証を実行する
- 実装対象 に対応するコード経路を優先し、承認済み実装範囲 外へ寄り道しない
- 単一引き継ぎ入力 の完了条件、公開接点、検証コマンド から着手する
- 文脈不足基準 は 構造判定条件 とし、完了条件、公開接点、実装対象、承認済み実装範囲、検証コマンド の不足時に返す
- 実装対象 がファイル、symbol、公開接点に対応していない場合は 文脈不足 を返す
- 単一引き継ぎ入力 内の局所確認、既存 型 への通常追従、レーン内検証 失敗 は not_insufficient_context として扱う
- 既存 型、naming、層 に合わせる
- 広域構造整理を混ぜない
- シナリオテスト、単体テスト、検証データ、スナップショット、test helper は test 担当 agent が扱う
- docs 正本化をしない

- コーディング規約 と lint 規約 から、引き継ぎ に効く静的 check の責務を確認する
- 単一引き継ぎ入力 の完了条件、承認済み実装範囲、実装対象、関連 根拠参照、検証コマンド を確認する
- 引き継ぎ に architecture constraint がある場合は、その範囲だけ architecture 規約 を局所確認する
- 実装小範囲 があれば 完了条件 clause、公開接点、変更対象 / symbol、検証コマンド を確認する
- 文脈不足 を返す場合は reason、必要文脈、推奨分割軸、remaining_実装小範囲s を 構造判定条件 に対応づける
- 入口、呼び出し箇所、データ流れ、エラー経路、test 表面 を確認する
- 既存 型 に naming、constructor、DI、エラー return を合わせる
- 生成 import、層 依存、境界 rule、整形逸脱など、変更層で踏みやすい lint 観点を先に確認する
- レーン内検証 結果または未実行理由を返す
- backend 引き継ぎ は `python3 scripts/harness/run.py --suite backend-local`、frontend 引き継ぎ は `python3 scripts/harness/run.py --suite frontend-local` を使う
- mixed 引き継ぎ は変更層に応じて両方を実行する
- 変更ファイルは プロダクトコード だけにする
- active 規約 は agent 1:1。backend / frontend / mixed / fix-lane の差分は 重点 skill で扱い、出力 obligation はこの 規約 に固定する。implementation_implementer は承認済み 引き継ぎ 1 件の プロダクトコード 実装を扱い、`APIテスト` 先行時だけ implementation_scenario_tester 出力 を受け取る。プロダクトテスト / 検証データ / スナップショット / test helper は変更しない。

## 非対象規約

- UI check、implementation レビュー、要件追加、設計追加は扱わない。
- docs、`.codex`、`.codex/skills`、`.codex/agents` は変更しない。
- プロダクトテスト、検証データ、スナップショット、test helper は変更しない。
- 承認済み実装範囲外の掃除、改名、整形、広域構造整理は扱わない。
- coverage、harness all、repo-local Sonar issue 判定条件は必須終了処理にしない。

## 出力規約

- 判断結果: プロダクトコード実装の完了、未完了、停止の判定を返す。
- 根拠参照: 実装の根拠にした入力、変更箇所、検証結果を返す。
- 不足情報: 実装を完了できない不足項目を返す。
- 次判断材料: `implement_lane` が次を判断できる材料を返す。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコード変更の指示を含めない。
- 返却先: implement_lane
- 実装成果物: 単一引き継ぎ入力 の 承認済み実装範囲 に対応する プロダクトコード だけを返す。プロダクトテスト、検証データ、スナップショット、test helper は含めない
- 引き継ぎ対応: 単一引き継ぎ入力 1 件と 実装小範囲 に対応づけ、複数 引き継ぎ を束ねない
- 実装済み完了条件: 実際に実装した 完了条件 clause、公開接点 / API 境界、変更対象 / symbol、検証コマンド を返す。実装小範囲 が入力された場合はそれに対応づける
- 未実装小範囲: 同じ 引き継ぎ 内で未実装の 小範囲 を返す。完了条件は削らず、未処理分を明示する
- レーン内検証結果: 実装完了後、引き継ぎ を終える前に変更層に対応する局所検証結果を返す。backend は `python3 scripts/harness/run.py --suite backend-local`、frontend は `python3 scripts/harness/run.py --suite frontend-local`、mixed は変更層に応じて両方を実行する。未実行なら 阻害理由 を返す。coverage、Sonar、harness all は implementation_implementer の必須 終了処理 にしない
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
- `APIテスト` 先行時だけ implementation_scenario_tester 出力 を確認した。
- 実装小範囲 がある場合はその範囲だけを実装した。
- 文脈不足基準 の 構造判定条件 に一致する場合だけ 文脈不足、必要文脈、推奨分割軸 を返した。
- not_insufficient_context に該当する局所確認、既存 型 追従、レーン内検証 失敗 を停止理由にしなかった。
- 実装対象 と 公開接点 から着手した。
- コーディング規約、lint 規約、レーン内検証コマンドを確認した。
- backend 変更を含む場合は `python3 scripts/harness/run.py --suite backend-local` を実行し、結果または未実行理由を返した。
- frontend 変更を含む場合は `python3 scripts/harness/run.py --suite frontend-local` を実行し、結果または未実行理由を返した。
- backend と frontend の両方を含む場合は両方の局所ハーネスを実行し、結果または未実行理由を返した。
- 引き継ぎ にある architecture constraint を局所確認した。
- 重点 skill の知識だけを追加で参照した。
- 変更ファイルが プロダクトコード だけであることを確認した。
- 必須 根拠: 単一引き継ぎ入力 id, 実装対象, 承認済み実装範囲, 承認記録, implementation_scenario_tester 出力 APIテスト先行実装テストがある場合, 実装済み小範囲 or 文脈不足 reason, 入口, 呼び出し箇所, データ流れ または境界, エラー経路, test 表面, 変更層の局所 検証結果 or 阻害理由
- 完了判断材料: implement_lane が レビュー へ進める プロダクトコード 実装結果と 変更層の局所 検証 結果が返っている
- 残留リスク: residual_risks

## 停止規約

- UI 確認や implementation レビュー を行う必要がある場合は停止する。
- docs や 作業流れ 文書を変更する必要がある場合は停止する。
- 文脈不足 を返さず広い調査で不足 文脈 を埋める必要がある場合は停止する。
- 判定基準 mismatch になる不安や通常の局所確認を文脈不足扱いにする必要がある場合は停止する。
- 実装小範囲 外へ実装を広げる必要がある場合は停止する。
- 実装対象 がないまま広い調査を始める必要がある場合は停止する。
- lint 未確認のまま実装して局所検証で初めて境界違反を知る進め方になる場合は停止する。
- config、lint、test、coverage 設定変更で判定条件を回避する必要がある場合は停止する。
- 停止時は不足項目、衝突箇所、戻し先を返す。
- 単一引き継ぎ入力が不足する場合は停止する。
- 実装対象が不足する場合は停止する。
- 承認記録が不足する場合は停止する。
- APIテスト先行の実装前引き継ぎで implementation_scenario_tester 出力が不足する場合は停止する。
- 承認済み実装範囲が不明な場合は停止する。
- 設計判断が不足している場合は停止する。
- docs または 作業流れ の変更が必要になる場合は停止する。
- 広域構造整理なしでは実装できない場合は停止する。
- プロダクトテスト、検証データ、スナップショット、test helper の変更が必要になる場合は停止する。
- touched_test_files を返す必要がある場合は停止する。
- プロダクトテスト、検証データ、スナップショット、test helper を変更する必要がある場合は停止する。
- implementation-scope 全文 または後続 引き継ぎ を入力として要求する必要がある場合は停止する。
- 文脈不足を返さず広く調査して不足文脈を埋める必要がある場合は停止する。
- 文脈不足基準外の理由で文脈不足を返す必要がある場合は停止する。
- remaining_実装小範囲s を隠して完了扱いにする必要がある場合は停止する。
- 実装完了後に変更層の局所検証結果または未実行理由を返せない場合は停止する。
