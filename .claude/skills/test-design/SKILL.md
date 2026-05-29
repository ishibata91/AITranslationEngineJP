---
name: test-design
description: Codex 側のテスト設計作業プロトコル。作業計画フォルダ内の入力成果物から、指定されたテスト設計成果物を固定する。
---
# Test Design

## 目的

`test-design` は作業プロトコルである。
`test_designer` agent が作業計画フォルダ内の入力成果物を読み、指定されたテスト設計成果物を active plan 内に固定する時の判断基準を提供する。

## 対応ロール

- `test_designer` が使う。
- 呼び出し元は作業計画フォルダを渡す agent とする。
- 返却先は呼び出し元とする。
- 担当成果物は呼び出し元が指定したテスト設計成果物とする。

## 呼び出し元から渡される情報

- 作業計画フォルダ: 入力成果物の参照元とテスト設計成果物の書き戻し先にする `docs/exec-plans/active/<task-id>/`。

## 作業前に読む正本

- エージェント実行定義と実行境界は [test_designer.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.claude/agents/test_designer.md) に従う。
- E2E テスト規約は [e2e-test-guidelines.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/e2e-test-guidelines.md) とする。
- テストコーディング規約は [coding-guidelines-tests.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/coding-guidelines-tests.md) とする。
- 画面設計書正本は [screen-design](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/screen-design/README.md) とする。
- 詳細仕様正本は [detail-specs](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/README.md) とする。
- ユースケース正本は [usecases](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/usecases/README.md) とする。
- E2E テスト観点正本は [test-design.csv](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/e2e-test-design/test-design.csv) とする。
- 不足テストテンプレートは [missing-tests.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.claude/skills/test-design/assets/missing-tests.md) とする。
- 不足UCテンプレートは [missing-usecases.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.claude/skills/test-design/assets/missing-usecases.md) とする。
- 不足セレクタテンプレートは [missing-selectors.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.claude/skills/test-design/assets/missing-selectors.md) とする。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。

## skill 内の拘束条件

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
| `前提条件` | テスト開始時に利用者が画面上で確認できる表示状態。 |
| `手順` | selector、操作対象、操作種別、入力値を指定した操作手順。 |
| `期待値` | selector、状態変化、表示内容、後続導線を指定した期待結果。 |
| `備考` | 補足、制約、未決事項。 |

テスト分類は次の表に従う。
この表は、関連ユースケースから拾う観点を拘束する。

| 分類 | 意味 |
| --- | --- |
| 正常 | 利用者が目的を達成する標準経路。 |
| 代替 | 利用者が目的を中断、取り消し、検索結果なし、差分なしなどで終える経路。 |
| 例外 | 入力不正、資格情報不足、形式不正、実行条件不足などで目的を達成できない経路。 |
| 境界 | 無効状態、重複実行、件数 0、状態遷移直前直後など、仕様の端を確認する経路。 |

## 担当ロールが判断してよい範囲

- 作業計画フォルダ内の入力成果物と、作業前に読む正本だけを根拠にする。
- テスト設計成果物の行と列は、E2E テスト規約に従う。
- 関連ユースケースごとに、正常、代替、例外、境界の分類から必要な観点を拾う。
- 正常系だけでテスト観点表を閉じない。
- 実装手順、テスト実装、検証コマンドはテスト観点表に含めない。

## skill が扱わない対象

- プロダクトコード、プロダクトテスト、docs 正本本文は変更しない。
- テスト実装は扱わない。
- ユースケース正本と E2E テスト観点正本そのものを変更しない。

## 返す成果物

- 概要: 作成または更新した不足テスト、不足UC、不足セレクタの要約を返す。
- 不足テスト: active plan 内に作成または更新した不足テスト成果物 path を返す。
- 不足UC: active plan 内に作成または更新した不足UC成果物 path を返す。
- 不足セレクタ: active plan 内に作成または更新した不足セレクタ成果物 path を返す。

## 作業を完了できる条件

- 呼び出し元が指定したテスト設計成果物が active plan 内に存在する。
- CSV のテスト設計成果物は E2E テスト規約の固定 CSV header を持つ。
- 各行が E2E テスト規約に従っている。
- 関連ユースケースから必要な正常、代替、例外、境界の観点が抽出されている。
- `data-testid-gaps.md` が存在する場合は、対象画面、対象要素、必要 selector、関連テスト ID が書かれている。
- プロダクトコード、プロダクトテスト、docs 正本本文を変更していない。

## 作業を止める条件

- 作業計画フォルダが不足する場合は停止する。
- 呼び出し元が指定したテスト設計成果物が不明な場合は停止する。
- 作業計画フォルダ内の入力成果物が不足する場合は停止する。
- E2E テスト規約を満たす観点を根拠から判断できない場合は停止する。
- プロダクトコード、プロダクトテスト、docs 正本本文の変更が必要な場合は停止する。
- 停止時は不足項目、衝突箇所、戻し先を返す。
