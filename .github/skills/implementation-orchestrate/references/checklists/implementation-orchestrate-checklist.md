# Implementation Orchestrate Checklist

## Knowledge Check

- [ ] 承認済み implementation-scope と approval record を確認した
- [ ] handoff 見出しを 1 実行単位として扱った
- [ ] depends_on と validation commands を確認した
- [ ] single_handoff_packet 1 件、tester_context_packet、lane_context_packet だけを subagent input にした
- [ ] distiller を tester / implementer より先に起動し、lane_context_packet を作った
- [ ] distiller output に tester_context_packet、test_ingredients、test_required_reading、test_validation_entry があることを確認した
- [ ] lane_context_packet に fix_ingredients、distracting_context、first_action、change_targets、requirements_policy_decisions、symbol / line number 付き related_code_pointers があることを確認した
- [ ] first_action が 1 completion_signal clause に固定され、partial や複数 clause でないことを確認した
- [ ] method / interface / field 追加候補が present / absent の code fact を持つことを確認した
- [ ] existing_patterns none と validation_entry が探索理由を持つことを確認した
- [ ] tester を implementer より先に起動した
- [ ] tester へ tester_context_packet と test_subscope だけを渡し、implementer 用 full context を渡さなかった
- [ ] tester / implementer の insufficient_context は同一 handoff 内の narrowing trigger として扱った
- [ ] insufficient_context の reason が insufficient_context_criteria に一致することを確認した
- [ ] criteria mismatch の insufficient_context は contract violation として completion packet に残した
- [ ] narrowing した場合は remaining_test_subscopes / remaining_implementation_subscopes を残した
- [ ] implementer へ lane_context_packet と tester output 以外の追加文脈を渡さなかった

## Common Pitfalls

- [ ] 通常の不足や validation failure を Codex return 前提にせず Copilot 内 narrowing で扱った
- [ ] hard stop 条件に該当する場合だけ requires_codex_replan を true にした
- [ ] distiller に full implementation-scope、active work plan 全文、source artifacts、後続 handoff を渡さなかった
- [ ] handoff 文面の言い換えだけの distiller output を implementer に渡さなかった
- [ ] fix_ingredients がない distiller output を implementer に渡さなかった
- [ ] tester_context_packet がない distiller output を tester に渡さなかった
- [ ] tester へ full lane_context_packet、fix_ingredients、change_targets を渡さなかった
- [ ] distracting_context を required_reading に混ぜた distiller output を implementer に渡さなかった
- [ ] first_action が partial または複数 clause の distiller output を implementer に渡さなかった
- [ ] 推測 method を fact にした distiller output を implementer に渡さなかった
- [ ] 要件、実装方針、決定事項を required_reading に丸投げした distiller output を implementer に渡さなかった
- [ ] implementer に full implementation-scope、active work plan 全文、source artifacts、後続 handoff を渡さなかった
- [ ] implementer に product test、fixture、snapshot、test helper の変更を依頼しなかった
- [ ] autonomous narrowing で completion_signal を削らなかった
- [ ] autonomous narrowing を理由に docs、implementation-scope、active work plan を変更しなかった
- [ ] criteria mismatch を narrowing trigger にしなかった
- [ ] docs / workflow 文書変更を混ぜなかった
- [ ] completion packet に docs_changes を残した
