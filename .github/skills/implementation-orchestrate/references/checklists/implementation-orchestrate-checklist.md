# Implementation Orchestrate Checklist

## Knowledge Check

- [ ] 承認済み implementation-scope と approval record を確認した
- [ ] handoff 見出しを 1 実行単位として扱った
- [ ] depends_on と validation commands を確認した
- [ ] single_handoff_packet 1 件と lane_context_packet だけを subagent input にした
- [ ] distiller を tester / implementer より先に起動し、lane_context_packet を作った
- [ ] lane_context_packet に fix_ingredients、distracting_context、first_action、change_targets、requirements_policy_decisions、symbol / line number 付き related_code_pointers があることを確認した
- [ ] first_action が 1 completion_signal clause に固定され、partial や複数 clause でないことを確認した
- [ ] method / interface / field 追加候補が present / absent の code fact を持つことを確認した
- [ ] existing_patterns none と validation_entry が探索理由を持つことを確認した
- [ ] tester を implementer より先に起動した
- [ ] implementer へ lane_context_packet と tester output 以外の追加文脈を渡さなかった

## Common Pitfalls

- [ ] design 不足を実装側で補わなかった
- [ ] distiller に full implementation-scope、active work plan 全文、source artifacts、後続 handoff を渡さなかった
- [ ] handoff 文面の言い換えだけの distiller output を implementer に渡さなかった
- [ ] fix_ingredients がない distiller output を implementer に渡さなかった
- [ ] distracting_context を required_reading に混ぜた distiller output を implementer に渡さなかった
- [ ] first_action が partial または複数 clause の distiller output を implementer に渡さなかった
- [ ] 推測 method を fact にした distiller output を implementer に渡さなかった
- [ ] 要件、実装方針、決定事項を required_reading に丸投げした distiller output を implementer に渡さなかった
- [ ] implementer に full implementation-scope、active work plan 全文、source artifacts、後続 handoff を渡さなかった
- [ ] implementer に product test、fixture、snapshot、test helper の変更を依頼しなかった
- [ ] docs / workflow 文書変更を混ぜなかった
- [ ] completion packet に docs_changes を残した
