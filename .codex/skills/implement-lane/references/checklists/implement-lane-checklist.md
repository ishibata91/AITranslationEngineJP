# Implement Lane Checklist

## Knowledge Check

- [ ] memory、`.codex/README.md`、`docs/index.md` を確認した
- [ ] active / completed task folder に同種 task がないか確認した
- [ ] 新規実装レーンの成果物DAGを `depends_on`、`gate`、`completion_signal` で固定した
- [ ] `pass_review_evidence` が成果物DAGに含まれている
- [ ] spawn packet が context 継承なしで足りる内容になっている
- [ ] human review と人間向け Codex implementation lane handoff を分けた
- [ ] closeout、停止、reroute のいずれでも work report と benchmark evidence を用意した

## Common Pitfalls

- [ ] spawned agent に会話文脈を暗黙継承させなかった
- [ ] Codex から Codex implementation lane へ直接 handoff しなかった
- [ ] implementation completion前に正本化へ進まなかった
- [ ] fix、refactor、探索テスト、UX 改善探索をこの skill で詳細化しなかった
