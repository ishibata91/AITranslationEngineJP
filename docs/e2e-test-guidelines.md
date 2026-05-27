# E2E テスト規約

関連文書: [`index.md`](./index.md), [`coding-guidelines.md`](./coding-guidelines.md), [`coding-guidelines-tests.md`](./coding-guidelines-tests.md), [`scenario-tests/README.md`](./scenario-tests/README.md)

本書は、UI 人間操作 E2E の設計、実装、観点表の規約を定義する。
backend / frontend 全体のテスト品質規約は [`coding-guidelines-tests.md`](./coding-guidelines-tests.md) を正本にする。

## 1. 基本方針

- E2E テストは、利用者の画面操作を既定の開始点にする。
- 各テストは、他のテストの実行順に依存しない。
- 各テストは、単独実行に必要な前提条件を自分で満たす。
- 前提データ投入は、主シナリオの証明ではなく状態準備として扱う。
- 実外部 API、実 secret、実利用者データへ到達する E2E テストは書かない。

## 2. 前提条件

前提条件は、テスト実行前に成立している必要がある状態を明示する。
前提条件には、実行済みである必要があるデータ投入テスト、画面状態、データ種類を含める。

前提条件を作る操作は、可能な限り画面操作で行う。
画面操作だけでは主シナリオ以外の準備が重くなる場合は、DB seed、API helper、test helper を補助として使ってよい。

補助による前提準備は、利用者操作で証明する対象を置き換えない。
補助は、状態準備、外部境界の fake、決定的な test data 作成に限定する。

## 3. selector 方針

`data-testid` は、画面設計時点で大まかに固定する selector として扱う。
画面設計書に selector がある場合、E2E テストはその selector を使う。

画面文言だけに依存する selector は、利用者表示そのものを検証する場合に限定する。
CSS class や DOM 階層に依存する selector は、画面設計上の固定対象でない限り使わない。

## 4. Page Object

Page Object は、Playwright 専用の test helper 配下に置く。
Page Object は、画面操作の語彙をまとめる。

Page Object が持ってよい責務は次の通り。

- selector を使った画面操作
- 利用者操作を表す method
- 画面遷移後の待機
- file chooser などの Playwright 操作補助

Page Object が持たない責務は次の通り。

- 前提データ投入
- 期待値の計算
- assertion
- 業務判断
- 複数シナリオにまたがる条件分岐

## 5. テスト独立性

各 E2E テストは、共有状態、前回実行結果、他テストの副作用に依存しない。
各 E2E テストは、必要な状態とデータ種類を観点表またはテスト内の Arrange で明示する。

テスト間の依存が必要に見える場合は、関連する確認を別テストの完了条件ではなく前提条件として表現する。
前提条件は、同じテスト内の準備または決定的な helper で作る。

## 6. テスト観点表

テスト観点表は CSV とする。
CSV header は次に固定する。

```csv
ID,関連UC,対象画面,前提条件,手順,期待値,備考
```

各列の意味は次の通り。

| 列 | 意味 |
| --- | --- |
| `ID` | テスト観点を一意に識別する ID。 |
| `関連UC` | 関連するユースケース。 |
| `対象画面` | 操作または検証の主対象画面。 |
| `前提条件` | 実行前に終わらせている必要があるテスト、状態、データ種類。 |
| `手順` | selector レベルで指定した操作手順。 |
| `期待値` | selector レベルで指定した期待結果。 |
| `備考` | 補足、制約、未決事項。 |

`手順` と `期待値` には、画面設計時点で決まっている `data-testid` をできるだけ書く。
selector が未確定の場合は、画面設計側の未決事項として扱う。

## 7. CSV 記入例

```csv
ID,関連UC,対象画面,前提条件,手順,期待値,備考
E2E-LOGIN-001,UC-LOGIN,ログイン画面,有効な利用者データが存在する,"[data-testid=email] に email を入力; [data-testid=password] に password を入力; [data-testid=login] をクリック","[data-testid=dashboard] が表示される","認証方式の詳細は対象機能の仕様を参照する"
```
