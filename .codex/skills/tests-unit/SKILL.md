---
name: tests-unit
description: Codex 実装系レーン側の 単体テスト 実装作業プロトコル。
---
# Tests Unit

## 目的

この skill は作業プロトコルである。
`implementation_unit_tester` agent が実装済み責務 または 軽量変更レーンの `task 枠` の 公開振る舞い、分岐、エラー経路 を 単体テスト で証明する時の判断基準を提供する。

## 対応ロール

- `implementation_unit_tester` が使う。
- 呼び出し元は `implement_lane`、`light_change_lane`、`refactor_lane` のいずれかとする。
- 返却先は呼び出し元レーンとする。
- 担当成果物は `tests-unit` の出力規約で固定する。

## 入力規約

- 単一引き継ぎ入力: `implementation-scope` から切り出された tests-unit 用 引き継ぎ 1 件、または 軽量変更レーンの `テスト修正証跡` 用 引き継ぎ 1 件。
- 実行中タスク成果物場所: テスト成果、検証結果、停止理由を書き戻す作業計画フォルダまたは run 成果物フォルダ。
- 仕様根拠: 承認済み `detail-spec-diff.md` または関連する `docs/detail-specs/<detail-spec-id>.md`。
- 対象テスト範囲: 変更してよい 単体テスト と必要最小限の テスト補助 の path。
- 実装済み対象: 実装種別別 agent が変更済みのファイル、公開接点、symbol。
- 証明対象: 公開振る舞い、分岐、エラー経路 のいずれを証明するかを示す対象。
- 検証コマンド: 実行を許可された backend-local または frontend-local の harness command。
- 網羅率検証コマンド: `python3 scripts/harness/run.py --suite coverage` で実行する harness command。

## 外部参照規約

- エージェント実行定義と実行境界は [implementation_unit_tester.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_unit_tester.toml) に従う。
- 詳細仕様正本: [detail-specs](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/README.md) とする。
- テストコーディング規約: [coding-guidelines-tests.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/coding-guidelines-tests.md) とする。
- lint 規約: [lint-policy.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/lint-policy.md) とする。
- architecture 規約: 引き継ぎに architecture constraint がある場合だけ [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) を参照する。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

- `公開振る舞い`: 公開接点から観測できる最小の振る舞いを証明する。
- `分岐`: 条件ごとの結果を 1 テスト 1 分岐で証明する。
- `エラー経路`: 入力不備、依存先失敗、不整合などの戻り方を証明する。
- `単体テスト`: 実装済み範囲に対応する局所テストとして扱う。

## 判断規約

- 各テストは 1 つの 公開振る舞い、分岐、エラー経路 のどれか 1 つを証明する
- 期待結果は仕様根拠、承認済み実装範囲、実装済み対象から導く
- setup は決定的にする
- テスト本体に条件分岐を入れない
- implementation_task_ids の外まで広げない
- 原因未確定の 回帰テスト は実装前に書かない
- `harness all` と repo-local Sonar issue 判定条件 は 最終検証 レーン へ送る
- 網羅率検証は `python3 scripts/harness/run.py --suite coverage` を実行し、全体網羅率が 70.0% を上回ることを確認する

- Arrange / Act / Assert を空行または短いコメントで判別できる状態にする
- テスト本体には意味的に何の振る舞いを証明するテストかを短いコメントで書く
- 分岐 ごとに テストケース を分ける
- clock、random、ID、repository 応答順序を固定する
- テストコーディング規約の良いテストの品質観点に従う

## 非対象規約

- シナリオ成果物の結果、統合 flow、新しい要件解釈は扱わない。
- テストのためだけの広いプロダクトコード変更は扱わない。

## 出力規約

- 判断結果: 単体テストを実装したか、文脈不足で停止したかを返す。
- 根拠参照: 単一引き継ぎ入力、仕様根拠、実装済み対象、変更ファイルを返す。
- 不足情報: 不足した入力項目、衝突した根拠、戻し先を返す。
- テスト成果物: 実装済み範囲に対応する 単体テスト と必要最小限の 検証データ / 補助 だけを返す。
- 証明済み完了条件: テストで証明した 公開振る舞い、分岐、エラー経路、テスト対象ファイル、検証コマンドを返す。
- 網羅率検証結果: `python3 scripts/harness/run.py --suite coverage` の結果と全体網羅率値を返す。
- 未証明小範囲: 同じ 引き継ぎ 内で未証明の 公開振る舞い、分岐、エラー経路を返す。
- 影響範囲修正: 今回のテスト変更が直接壊したテスト補助、fixture、検証経路、担当単体テスト成果物を修正した場合に、対象、理由、変更結果を返す。
- レーン内検証結果: テスト追加または更新後、変更層 に対応する 局所検証 の失敗時は承認済み実装範囲、軽量変更レーンの `task 枠`、今回のテスト変更が直接壊した担当単体テスト成果物の影響範囲修正で直して再実行し、通過結果または未実行理由を返す。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコード変更の指示を含めない。

## 完了規約

- 承認済み実装範囲、軽量変更レーンの `task 枠`、今回のテスト変更が直接壊した担当単体テスト成果物の影響範囲修正 の成果だけが返却されている。
- 仕様根拠を読み、証明対象の期待結果と対応づけた。
- 検証、未実行項目、残留リスク が 根拠参照 付きで整理されている。
- 1 テストで 1 公開振る舞い / 分岐 / エラー経路 だけを証明した。
- setup の clock、random、ID、repository 応答順序を固定した。
- テスト本体に意味的に何の振る舞いを証明するテストかを示すコメントがある。
- テストコーディング規約の良いテストの品質観点に反するテスト品質問題が残っていない。
- implementation_task_ids の外へ広げなかった。
- 変更対象が 単体テスト と必要最小限の テスト補助だけである。
- backend 側の単体テストを変更した場合は `python3 scripts/harness/run.py --suite backend-local` を実行し、失敗した場合は承認済み実装範囲、軽量変更レーンの `task 枠`、今回のテスト変更が直接壊した担当単体テスト成果物の影響範囲修正 でその場で直して再実行し、通過結果または未実行理由を返した。
- frontend 側の単体テストを変更した場合は `python3 scripts/harness/run.py --suite frontend-local` を実行し、失敗した場合は承認済み実装範囲、軽量変更レーンの `task 枠`、今回のテスト変更が直接壊した担当単体テスト成果物の影響範囲修正 でその場で直して再実行し、通過結果または未実行理由を返した。
- backend と frontend の両方を含む場合は両方の局所ハーネスを実行し、失敗した場合は承認済み実装範囲、軽量変更レーンの `task 枠`、今回のテスト変更が直接壊した担当単体テスト成果物の影響範囲修正 でその場で直して再実行し、通過結果または未実行理由を返した。
- `python3 scripts/harness/run.py --suite coverage` を実行し、全体網羅率が 70.0% を上回る結果または未実行理由を返した。
- レーン内検証 の失敗時は承認済み実装範囲、軽量変更レーンの `task 枠`、今回のテスト変更が直接壊した担当単体テスト成果物の影響範囲修正で直して再実行し、通過結果または未実行理由を返した。

## 停止規約

- シナリオ 成果物 の 結果 を テストにする時
- 局所ハーネスの失敗原因が今回のテスト変更が直接壊したテスト補助、fixture、検証経路、担当単体テスト成果物に閉じない時
- `python3 scripts/harness/run.py --suite coverage` の全体網羅率が 70.0% 以下で、承認済み実装範囲、軽量変更レーンの `task 枠`、今回のテスト変更が直接壊した担当単体テスト成果物だけでは改善できない時
- プロダクトコード変更、UI 変更、secret / trust boundary 変更、API / DTO / DB / schema の意味拡張、docs 正本化、`.codex` 作業流れの変更が必要な時
- テストのためだけに広い プロダクトコード 変更が必要な時
- 統合 flow を証明する時
- 証明対象、対象テスト範囲、実装済み対象 のいずれかが不足している時
- 仕様根拠が不足している時
- 停止時は不足項目、衝突箇所、戻し先を返す。
- テスト本体に条件分岐が必要になる場合は停止する。
