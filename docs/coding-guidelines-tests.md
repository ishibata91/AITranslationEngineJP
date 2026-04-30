# テスト コーディング規約

関連文書: [`coding-guidelines.md`](./coding-guidelines.md), [`scenario-tests/README.md`](./scenario-tests/README.md), [`lint-policy.md`](./lint-policy.md)

本書は、backend / frontend のテスト実装規約を定義する。
受け入れ条件、単体分岐、統合境界、UI 人間操作 E2E の証明方法を対象にする。
プロダクト実装規約は frontend / backend の各文書を正本にする。

## 1. テスト分類

- 受け入れテストは実装判断より先に固定し、scenario-design の分類を正本にする
- 単体テストは公開振る舞い、分岐、エラー経路を小さい単位で証明する
- API / public seam test は UI 人間操作 E2E と混ぜず、入口と観測点を分ける
- UI 人間操作 E2E は人間操作の主要 flow と表示結果を証明する
- 検証データ、clock、random、ID、repository 応答順序は固定する

## 2. テスト構造

- Arrange、Act、Assert が読める構造にする
- テスト名は期待する振る舞いを説明し、実装手順名だけにしない
- 成功経路と失敗経路は別の test case に分ける
- mock は境界の代替に限定し、実装詳細の呼び出し順へ過度に依存しない
- 失敗した時に原因の層が分かる assertion にする

## 3. Go テスト

- 標準の `go test` を前提にし、分岐が増える場合は table-driven tests を使う
- race 条件が関係する変更では `-race` 相当の検証を検討し、未実行なら理由を残す
- repository、filesystem、外部プロセス境界は fake または test helper で明示する
- production seed や実外部環境へ依存しない

## 4. TypeScript / Svelte テスト

- Vitest を frontend の executable spec として扱う
- component の外から観測できる表示、状態、callback を優先して検証する
- generated `wailsjs` や backend DTO の直接 import を広げる test helper を作らない
- Playwright が必要な flow は、UI 人間操作 E2E として scenario 側の根拠に対応づける

## 5. 完了条件

- backend 側のテスト変更は `python3 scripts/harness/run.py --suite backend-local` の結果または未実行理由を残す
- frontend 側のテスト変更は `python3 scripts/harness/run.py --suite frontend-local` の結果または未実行理由を残す
- backend と frontend の両方を含む場合は両方の局所ハーネス結果または未実行理由を残す
- coverage 目標は個別 task の承認済み条件がある場合だけ完了条件にする

## 6. 参照元

- 輸入元: `../everything-claude-code/rules/common/testing.md`
- 輸入元: `../everything-claude-code/rules/golang/testing.md`
- 輸入元: `../everything-claude-code/rules/typescript/testing.md`
