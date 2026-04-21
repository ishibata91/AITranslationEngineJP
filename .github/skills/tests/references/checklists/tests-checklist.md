# Tests Checklist

## Knowledge Check

- [ ] implemented scope と owned_scope を確認した
- [ ] tester_context_packet の test_ingredients、test_required_reading、test_validation_entry を確認した
- [ ] test_subscope がある場合はその範囲だけを証明した
- [ ] insufficient_context_criteria の structural gate に一致する場合だけ insufficient_context、needed_context、remaining_test_subscopes を返した
- [ ] not_insufficient_context に該当する局所確認や期待どおり fail する test を停止理由にしなかった
- [ ] deterministic setup にした
- [ ] focused skill の知識だけを追加で参照した

## Common Pitfalls

- [ ] 新しい要件解釈を足さなかった
- [ ] full lane_context_packet、fix_ingredients、change_targets、broad related_code_pointers を直接追わなかった
- [ ] insufficient_context を広い調査で埋めなかった
- [ ] criteria mismatch になる insufficient_context を返さなかった
- [ ] paid real AI API を呼ばなかった
- [ ] mode 別 active contract を使わなかった
