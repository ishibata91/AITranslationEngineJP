# Codex Review Conductor Checklist

- [ ] diff、implementation-scope、implementation result、final validation result を確認した
- [ ] payload または validation 不足時は観点 group を spawn せず早期 return した
- [ ] review 可能な時だけ 4 観点 group を context 継承なしで並列 spawn した
- [ ] 各 group result に observed scope、violated invariant、root cause hypothesis、local patch assessment があることを確認した
- [ ] 各 group result の raw result を `reviewer_result_bundle` に落とさず残した
- [ ] `aggregation_trace` で各 aggregation field の参照元 group と source field を示した
- [ ] 採用しなかった観点 signal を `unselected_group_signals` または `residual_risks` に残した
- [ ] 欠落 field と低 confidence を `information_loss_notes` に残した
- [ ] primary failure mode、dominant invariant、minimum durable fix boundary を統合した
- [ ] `strict_pass` と `priority_override_pass` を分けた
- [ ] priority override した finding を `priority_overrides` と `residual_risks` に残した
- [ ] hard gate failure を平均 score で相殺しなかった
- [ ] `copilot_action` と remediation handoff を返し、Copilot の chosen strategy 判断に必要な材料を残した
