# Scenario Design Checklist

## Knowledge Check

- [ ] 必ず通す要件と risk を分けた
- [ ] 抽象要件を詳細要求タイプへ展開した
- [ ] 各詳細要求タイプを `explicit`、`derived`、`not_applicable`、`deferred`、`needs_human_decision` に分類した
- [ ] 仕様網羅を `scenario-design.requirement-coverage.json` に分離した
- [ ] `needs_human_decision` だけを `scenario-design.questions.md` へ集約した
- [ ] user journey と scenario matrix を分けた
- [ ] 開始条件、操作、期待結果、観測点を明示した
- [ ] UI が入口のシナリオは、ユーザー操作と入力から得られる値を開始条件に含めた
- [ ] fake / fixture / validation command を確認した

## Common Pitfalls

- [ ] 実装方針を要件として固定しなかった
- [ ] 人間判断が必要な暗黙要求を AI 判断で固定しなかった
- [ ] paid な real AI API を前提にしなかった
- [ ] happy path だけにしなかった
- [ ] product test 実装詳細を書かなかった
- [ ] 裏側の直接呼び出しだけの検証を、UI 入口の端から端までのシナリオとして扱わなかった
