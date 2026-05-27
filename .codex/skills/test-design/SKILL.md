---
name: test-design
description: Codex 側のテスト設計作業プロトコル。承認済み設計から active plan 内にテスト観点表を固定する。
---
# Test Design

## 目的

`test-design` は作業プロトコルである。
`test_designer` agent が承認済み設計を読み、active plan 内にテスト観点表を固定する時の判断基準を提供する。

## 対応ロール

- `test_designer` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- 担当成果物は `テスト設計` とする。

## 入力規約

- 作業計画フォルダ: テスト観点表を書き戻す `docs/exec-plans/active/<task-id>/`。
- 承認済み詳細仕様差分: テスト観点の仕様根拠にする `detail-spec-diff.md`。
- 承認済み画面設計差分: UI 人間操作 E2E の selector と対象画面の根拠にする `screen-design-diff.<screen-id>.md`。
- 関連ユースケース: テスト観点の `関連UC` に対応する docs 正本または task 内成果物。
- 承認記録: 詳細仕様差分、画面設計差分、関連ユースケースをテスト設計根拠として扱ってよい人間承認。
- 非必須入力: 既存テスト観点表。

## 外部参照規約

- エージェント実行定義と実行境界は [test_designer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/test_designer.toml) に従う。
- E2E テスト規約は [e2e-test-guidelines.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/e2e-test-guidelines.md) とする。
- テストコーディング規約は [coding-guidelines-tests.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/coding-guidelines-tests.md) とする。
- 画面設計書正本は [screen-design](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/screen-design/README.md) とする。
- 詳細仕様正本は [detail-specs](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/README.md) とする。
- 新規実装レーン入口は [implement-lane](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-lane/SKILL.md) とする。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

テスト観点表は active plan 内の `test-design.csv` に置く。
CSV header は次に固定する。

```csv
ID,関連UC,対象画面,前提条件,手順,期待値,備考
```

各列の意味は次の表に従う。

| 列 | 意味 |
| --- | --- |
| `ID` | テスト観点を一意に識別する ID。 |
| `関連UC` | 関連するユースケース。 |
| `対象画面` | 操作または検証の主対象画面。 |
| `前提条件` | 実行前に終わらせている必要があるテスト、状態、データ種類。 |
| `手順` | selector レベルで指定した操作手順。 |
| `期待値` | selector レベルで指定した期待結果。 |
| `備考` | 補足、制約、未決事項。 |

## 判断規約

- 承認済み詳細仕様差分、承認済み画面設計差分、関連ユースケースだけをテスト観点の根拠にする。
- UI 人間操作 E2E は画面操作を既定の開始点にする。
- `data-testid` は画面設計時点で固定する selector として扱う。
- `手順` は selector レベルで書く。
- `期待値` は selector レベルで書く。
- 前提条件は実行順依存ではなく、各テストが単独実行に必要な状態として書く。
- 前提データ投入は主シナリオの証明ではなく状態準備として扱う。
- 1 行は 1 つのテスト観点だけを表す。
- 未確定 selector は独断で補完せず、`備考` に未決として書く。
- 実装手順、テスト実装、検証コマンドはテスト観点表に含めない。

## 非対象規約

- プロダクトコード、プロダクトテスト、docs 正本本文は変更しない。
- `implementation-scope.md` は作らない。
- シナリオテスト実装と単体テスト実装は扱わない。
- 未承認の詳細仕様差分、画面設計差分、関連ユースケースは根拠にしない。

## 出力規約

- 判断結果: テスト設計を完了したか、文脈不足で停止したかを返す。
- 根拠参照: テスト観点表の根拠にした詳細仕様差分、画面設計差分、関連ユースケース、承認記録を返す。
- テスト観点表: active plan 内に作成または更新した `test-design.csv` を返す。
- 不足情報: 不足した入力項目、衝突した根拠、戻し先を返す。
- 次判断材料: `implement_lane` が後続の実装範囲、シナリオテスト、単体テストの扱いを判断できる材料を返す。
- 禁止事項: 出力にプロダクトコード、プロダクトテスト、docs 正本本文の変更を含めない。

## 完了規約

- `test-design.csv` が active plan 内に存在する。
- `test-design.csv` が固定 CSV header を持つ。
- 各行の `関連UC` が関連ユースケースに対応している。
- 各行の `対象画面` が画面設計差分または画面設計書正本に対応している。
- 各行の `手順` と `期待値` が selector レベルで書かれている。
- 各行の `前提条件` が単独実行に必要な状態として書かれている。
- 未確定 selector がある場合は `備考` に未決として書かれている。
- プロダクトコード、プロダクトテスト、docs 正本本文を変更していない。

## 停止規約

- 作業計画フォルダが不足する場合は停止する。
- 承認済み詳細仕様差分が不足する場合は停止する。
- 関連ユースケースが不足する場合は停止する。
- UI 人間操作 E2E が必要なのに画面設計差分または画面設計書正本が不足する場合は停止する。
- 承認記録が不足する場合は停止する。
- 根拠間に衝突がある場合は停止する。
- プロダクトコード、プロダクトテスト、docs 正本本文の変更が必要な場合は停止する。
- 停止時は不足項目、衝突箇所、戻し先を返す。
