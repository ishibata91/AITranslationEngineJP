# Implementation Orchestrate Checklist

## Knowledge Check

- [ ] 承認済み implementation-scope と approval record を確認した
- [ ] Ready Waves 表または `ready_wave` から実行可能 wave を選んだ
- [ ] `single_handoff_packet` に `first_action` が含まれていることを確認した
- [ ] distiller を implementation_tester / implementation_implementer より先に起動した
- [ ] `APIテスト` 先行条件を満たす場合だけ implementation_tester を implementation_implementer より先に起動した
- [ ] unit test と原因未確定の regression test を実装後の implementation_tester に回した
- [ ] scenario validation、suite-all、Sonar check を全 implementation handoff 完了後に実行した
- [ ] scenario validation が fail した場合は close せず blocker として返した
- [ ] `UI人間操作E2E` を final validation lane で扱った
- [ ] review input が揃う場合だけ 4 観点 review agent を context 継承なしで並列 spawn した
- [ ] review input 不足時は reviewer を起動せず `rerun_validation` または `rerun_codex_review` を返した
- [ ] `implementation_action` に従って close / report_residual / fix / rerun_validation / rerun_codex_review を分岐した
- [ ] `fix` では reviewer result bundle、aggregation trace、remediation handoff から chosen strategy、chosen scope、why_not_narrower、why_not_wider、used_review_signals を返した
## Common Pitfalls

- [ ] final validation 前に scenario validation、suite-all、Sonar check を実行しなかった
- [ ] `first_action` 欠落を広い調査で補わなかった
- [ ] `parallelizable_with` がない handoff を並列実行しなかった
- [ ] repo-local Sonar issue gate と Sonar server Quality Gate を混同しなかった
- [ ] integrated review input に diff、scope、implementation result、validation result を含めた
- [ ] `rerun_codex_review` で product code を変更しなかった
- [ ] `fix` で symptoms だけを潰す局所修正に閉じなかった
- [ ] docs / workflow 文書変更を implementation lane に混ぜなかった
