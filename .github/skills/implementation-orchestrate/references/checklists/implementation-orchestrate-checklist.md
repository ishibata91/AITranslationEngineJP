# Implementation Orchestrate Checklist

## Knowledge Check

- [ ] 承認済み implementation-scope と approval record を確認した
- [ ] handoff 見出しを 1 実行単位として扱った
- [ ] depends_on と validation commands を確認した
- [ ] single_handoff_packet 1 件と lane_context_packet だけを subagent input にした
- [ ] distiller を tester / implementer より先に起動し、lane_context_packet を作った
- [ ] tester を implementer より先に起動した
- [ ] implementer へ lane_context_packet と tester output 以外の追加文脈を渡さなかった

## Common Pitfalls

- [ ] design 不足を実装側で補わなかった
- [ ] distiller に full implementation-scope、active work plan 全文、source artifacts、後続 handoff を渡さなかった
- [ ] implementer に full implementation-scope、active work plan 全文、source artifacts、後続 handoff を渡さなかった
- [ ] implementer に product test、fixture、snapshot、test helper の変更を依頼しなかった
- [ ] docs / workflow 文書変更を混ぜなかった
- [ ] completion packet に docs_changes を残した
