# Scenario Design Checklist

## Knowledge Check

- [ ] 必ず通す要件と risk を分けた
- [ ] 抽象要件を詳細要求タイプへ展開した
- [ ] 各詳細要求タイプを `explicit`、`derived`、`not_applicable`、`deferred`、`needs_human_decision` に分類した
- [ ] 仕様網羅を `scenario-design.requirement-coverage.json` に分離した
- [ ] `needs_human_decision` だけを `scenario-design.questions.md` へ集約した
- [ ] user journey と scenario matrix を分けた
- [ ] 受け入れテストを全 scenario case で先に固定した
- [ ] 各 scenario case に `実行テスト種別` と `実行段階` を書いた
- [ ] 開始条件、操作、期待結果、観測点を明示した
- [ ] `APIテスト` では受け入れ条件、public seam、入力開始点、主要 outcome、主要観測点、contract freeze を固定した
- [ ] `UI人間操作E2E` では開始操作、入力方法、主要操作列、主要観測点、UI-visible outcome、fake / stub 方針を固定した
- [ ] fake / fixture / validation command を確認した

## Common Pitfalls

- [ ] 実装方針を要件として固定しなかった
- [ ] 人間判断が必要な暗黙要求を AI 判断で固定しなかった
- [ ] paid な real AI API を前提にしなかった
- [ ] happy path だけにしなかった
- [ ] product test 実装詳細を書かなかった
- [ ] 裏側の直接呼び出しだけの検証を、UI 入口の `UI人間操作E2E` として扱わなかった
