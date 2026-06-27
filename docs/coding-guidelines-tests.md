# テスト コーディング規約

関連文書: [`coding-guidelines.md`](./coding-guidelines.md), [`e2e-test-guidelines.md`](./e2e-test-guidelines.md), [`scenario-tests/README.md`](./scenario-tests/README.md), [`lint-policy.md`](./lint-policy.md)

本書は、backend / frontend のテスト実装規約を定義する。
受け入れ条件、単体分岐、統合境界、UI 人間操作 E2E の証明方法を対象にする。
プロダクト実装規約は frontend / backend の各文書を正本にする。

## 1. 良いテストの品質観点

| 観点 | 規約 |
| --- | --- |
| 仕様根拠 | テストは仕様、承認済み差分、または人間が正とした実装判断を根拠にする |
| 証明対象 | テストは 1 つの公開振る舞い、シナリオ結果、分岐、またはエラー経路だけを証明する |
| 意味コメント | テスト本体には、意味的に何の振る舞いを証明するテストかを短いコメントで明示する |
| 観測点 | テストは利用者、公開接点、または責務境界から観測できる結果を検証する |
| 期待値 | 期待値は実装都合ではなく、仕様根拠から読める結果にする |
| 失敗診断 | assertion が失敗した時に、壊れた責務、入力、期待値のどれかを特定できるようにする |
| 決定性 | clock、random、ID、並び順、外部応答を固定し、同じ入力で同じ結果にする |
| 独立性 | テストは他のテストの実行順、共有状態、前回実行結果に依存しない |
| 前提明示 | setup は必要な前提だけを作り、暗黙の repository 状態や既定値に依存しない |
| 入力代表性 | 入力データは正常、境界、失敗のどの条件を代表するか分かる名前にする |
| 境界値 | 境界値を扱うテストは、境界の意味と期待する内外差を明示する |
| 副作用隔離 | filesystem、network、database、外部プロセスの副作用は fake、temporary resource、または test helper へ閉じる |
| fixture 最小性 | fixture は証明対象に必要な最小量にし、無関係な field や巨大な copy を避ける |
| mock 境界 | mock は外部境界または遅い依存の代替に使い、内部実装の呼び出し順の固定に使わない |
| 可読性 | Arrange、Act、Assert の区切りが読み取れる構造にする |
| 保守性 | 仕様名、画面名、公開接点名が変わった時に、修正すべきテスト範囲が追える名前にする |
| 実行速度 | 単体テストは短時間で反復できる粒度にし、遅い検証はシナリオテストまたは局所ハーネスへ分ける |
| 回帰意図 | 回帰テストは観測済みの失敗、修正対象、再発条件を根拠にする |
| 網羅率 | coverage は目的ではなく補助指標として扱い、仕様を証明しない行数稼ぎをしない |

## 2. テスト分類

- 受け入れテストは実装判断より先に固定し、詳細仕様差分の親要件と仕様を根拠にする
- 単体テストは公開振る舞い、分岐、エラー経路を小さい単位で証明する
- API / public seam test は UI 人間操作 E2E と混ぜず、入口と観測点を分ける
- UI 人間操作 E2E は人間操作の主要 flow と表示結果を証明する
- UI 人間操作 E2E の selector、Page Object、観点表 CSV は [`e2e-test-guidelines.md`](./e2e-test-guidelines.md) に従う
- 検証データ、clock、random、ID、repository 応答順序は固定する

## 3. テスト構造

- Arrange、Act、Assert が読める構造にする
- テスト名は期待する振る舞いを説明し、実装手順名だけにしない
- テスト本体の先頭または Arrange 前に、意味的に何の振る舞いを証明するテストかを短いコメントで書く
- 成功経路と失敗経路は別の test case に分ける
- mock は境界の代替に限定し、実装詳細の呼び出し順へ過度に依存しない
- 失敗した時に原因の層が分かる assertion にする

## 4. Go テスト

- 標準の `go test` を前提にし、分岐が増える場合は table-driven tests を使う
- race 条件が関係する変更では `-race` 相当の検証を検討し、未実行なら理由を残す
- repository、filesystem、外部プロセス境界は fake または test helper で明示する
- production seed や実外部環境へ依存しない

## 5. TypeScript / Svelte テスト

- Vitest を frontend の executable spec として扱う
- component の外から観測できる表示、状態、callback を優先して検証する
- generated `wailsjs` や backend DTO の直接 import を広げる test helper を作らない
- Playwright が必要な flow は、UI 人間操作 E2E として scenario 側の根拠に対応づける

## 6. 完了条件

- backend 側のテスト変更は `npm run verify:backend`（`go test` ＋ arch-lint ＋ 境界違反走査を束ねる）の結果または未実行理由を残す
- frontend 側のテスト変更は `npm run test:frontend` の結果または未実行理由を残す
- backend と frontend の両方を含む場合は両方の結果または未実行理由を残す
- coverage 目標は個別 task の承認済み条件がある場合だけ完了条件にする。100% を求める純粋ルールは `go test -cover` の手元確認で見る（常設 coverage ゲートは持たない）

## 7. 参照元

- 輸入元: `../everything-claude-code/rules/common/testing.md`
- 輸入元: `../everything-claude-code/rules/golang/testing.md`
- 輸入元: `../everything-claude-code/rules/typescript/testing.md`
