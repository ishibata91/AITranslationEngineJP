# E2E Test Design

この directory は、UI 人間操作 E2E のテスト観点表正本を置く。
テスト観点表は、ユースケース、対象画面、前提条件、手順、期待値、備考を CSV で固定する。

## Records

- [`test-design.csv`](./test-design.csv): UC ベース UI 人間操作 E2E テスト観点表

## Rules

- `test-design.csv` は `docs/e2e-test-guidelines.md` の CSV header に従う。
- 前提条件は、各テストが単独実行に必要な画面表示状態として書く。
- 手順と期待値は、selector、状態変化、表示内容を特定できる粒度で書く。
- repository、DB、fake、seed、helper、fixture、保存先、外部境界の準備手順は前提条件に書かない。
- task 内の検討資料は `docs/exec-plans/` に残し、正本 CSV には実装判断へ渡す観点だけを残す。
