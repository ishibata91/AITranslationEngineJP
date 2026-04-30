---
name: implement-fix-lane
description: Codex implementation レーン 側の fix レーン 恒久修正作業プロトコル。
---
# Implement Fix Lane

## 目的

この skill は作業プロトコルである。
`implementation_implementer` agent が `承認済み修正範囲` の恒久修正を行う時に、再現条件と矛盾しない変更へ限定する判断基準を提供する。

## 対応ロール

- `implementation_implementer` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- 担当成果物は `implement-fix-lane` の出力規約で固定する。

## 入力規約

- 単一引き継ぎ入力: `implementation-scope` から切り出された fix 実装用 引き継ぎ 1 件。
- 実行中タスク成果物場所: 修正結果、検証結果、停止理由を書き戻す作業計画フォルダまたは run 成果物フォルダ。
- 修正対象: 変更してよいファイル、symbol、公開接点。
- 対象修正範囲: 実装してよい恒久修正のプロダクトコード範囲。
- 依存完了情報: 着手前に完了している必要がある依存対象の完了結果。
- 検証コマンド: 実行を許可された backend-local または frontend-local の harness command。

## 外部参照規約

- エージェント実行定義とツール権限は [implementation_implementer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_implementer.toml) の 書き込み許可 / 実行許可 とする。
- frontend コーディング規約: frontend 変更がある場合は [coding-guidelines-frontend.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/coding-guidelines-frontend.md) とする。
- backend コーディング規約: backend 変更がある場合は [coding-guidelines-backend.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/coding-guidelines-backend.md) とする。
- lint 規約: [lint-policy.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/lint-policy.md) とする。
- architecture 規約: 引き継ぎに architecture constraint がある場合だけ [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) を参照する。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

## 判断規約

- 承認済み修正範囲 を超えない
- 再現条件に関係しない整理を入れない
- trace_or_analysis_result と矛盾しない変更に限る
- 単一引き継ぎ入力 と 承認済み修正範囲 を確認して プロダクトコード だけを変更する
- `APIテスト` 先行時だけ implementation_scenario_tester 出力 も確認する
- 未解消ケースを 終了処理 に残す

- 修正前後で同じ条件の 検証 を比較する
- 残留リスク を明示する
- fix 対象範囲 と touched files を対応づける

## 非対象規約

- 新機能、refactor、unrelated cleanup は扱わない。
- 根拠なしの原因断定や、再現条件に関係しない整理は扱わない。
- プロダクトテスト、検証データ、スナップショット、test helper は変更しない。
- docs や作業流れ文書は変更しない。
- coverage、harness all、repo-local Sonar issue 判定条件は必須終了処理にしない。

## 出力規約

- 判断結果: fix-lane プロダクトコード実装の完了、未完了、停止の判定を返す。
- 根拠参照: 実装の根拠にした入力、変更箇所、検証結果を返す。
- 不足情報: 実装を完了できない不足項目を返す。
- 次判断材料: `implement_lane` が次を判断できる材料を返す。
- 実装成果物: 単一引き継ぎ入力 の 承認済み修正範囲 に対応するプロダクトコードだけを返す。
- レーン内検証結果: 変更層に対応する局所検証の失敗時はその場で直して再実行し、通過結果または未実行理由を返す。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコード変更の指示を含めない。

## 完了規約

- 承認済み実装範囲 内の成果だけが返却されている。
- 検証、未実行項目、残留リスク が 根拠参照 付きで整理されている。
- 単一引き継ぎ入力、修正対象、対象修正範囲、依存完了情報、検証コマンドを確認した。
- 対象修正範囲 と 再現条件 を確認した。
- trace_or_analysis_result と矛盾しない変更に限定した。
- backend 変更を含む場合は `python3 scripts/harness/run.py --suite backend-local` を実行し、失敗した場合は承認済み修正範囲 内でその場で直して再実行し、通過結果または未実行理由を返した。
- frontend 変更を含む場合は `python3 scripts/harness/run.py --suite frontend-local` を実行し、失敗した場合は承認済み修正範囲 内でその場で直して再実行し、通過結果または未実行理由を返した。
- backend と frontend の両方を含む場合は両方の局所ハーネスを実行し、失敗した場合は承認済み修正範囲 内でその場で直して再実行し、通過結果または未実行理由を返した。
- 残留リスク と未解消ケースを 終了処理 に残した。

## 停止規約

- 新機能や refactor の実装を行う時
- 再現条件が不足している時
- 原因が未確認なのに恒久修正する時
- 単一引き継ぎ入力、修正対象、対象修正範囲、依存完了情報、検証コマンドが不足する場合は停止する。
- プロダクトテスト、検証データ、スナップショット、test helper の変更が必要になる場合は停止する。
- `python3 scripts/harness/run.py --suite backend-local` または `python3 scripts/harness/run.py --suite frontend-local` の失敗原因が承認済み修正範囲 外にある場合は停止する。
- 承認済み修正範囲外へ実装を広げる必要がある場合は停止する。
- 停止時は不足項目、衝突箇所、戻し先を返す。
- task_mode が fix であることを確認した。
