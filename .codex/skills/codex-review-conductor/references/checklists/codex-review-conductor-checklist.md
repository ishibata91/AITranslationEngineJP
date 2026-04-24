# Codex Review Conductor Checklist

- [ ] diff、implementation-scope、implementation result、final validation result を確認した
- [ ] payload または validation 不足時は観点 group を spawn せず早期 return した
- [ ] review 可能な時だけ 4 観点 group を context 継承なしで並列 spawn した
- [ ] `strict_pass` と `priority_override_pass` を分けた
- [ ] priority override した finding を `priority_overrides` と `residual_risks` に残した
- [ ] hard gate failure を平均 score で相殺しなかった
- [ ] `copilot_action` を返し、Copilot の受け取り分岐を再解釈不要にした
