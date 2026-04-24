# Implement Checklist

## Knowledge Check

- [ ] owned_scope と implementation target を確認した
- [ ] single_handoff_packet と lane_context_packet を確認した
- [ ] scenario 先行時だけ tester output を確認した
- [ ] implementation_subscope がある場合はその範囲だけを実装した
- [ ] fix_ingredients と distracting_context を確認した
- [ ] insufficient_context_criteria の structural gate に一致する場合だけ insufficient_context、needed_context、suggested_narrowing_axis を返した
- [ ] not_insufficient_context に該当する局所確認、既存 pattern 追従、lane-local validation failure を停止理由にしなかった
- [ ] first_action と change_targets から着手した
- [ ] coding guidelines と lane-local validation commands を確認した
- [ ] focused skill の知識だけを追加で参照した
- [ ] touched files が product code だけであることを確認した

## Common Pitfalls

- [ ] broad refactor を混ぜなかった
- [ ] insufficient_context を広い調査で埋めなかった
- [ ] criteria mismatch になる insufficient_context を返さなかった
- [ ] implementation_subscope 外へ実装を広げなかった
- [ ] distracting_context を実装対象に混ぜなかった
- [ ] first_action 不足を広い調査で埋めなかった
- [ ] product test、fixture、snapshot、test helper を変更しなかった
- [ ] docs / workflow 文書を変更しなかった
- [ ] mode 別 active contract を使わなかった
