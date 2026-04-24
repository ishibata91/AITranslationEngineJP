# Implementation Orchestrate Checklist

## Knowledge Check

- [ ] 承認済み implementation-scope と approval record を確認した
- [ ] Ready Waves 表または `ready_wave` から実行可能 wave を選んだ
- [ ] `single_handoff_packet` に `first_action` が含まれていることを確認した
- [ ] distiller を tester / implementer より先に起動した
- [ ] scenario 先行条件を満たす場合だけ tester を implementer より先に起動した
- [ ] unit test と原因未確定の regression test を実装後の tester に回した
- [ ] scenario validation、suite-all、Sonar check を全 implementation handoff 完了後に実行した
- [ ] scenario validation が fail した場合は close せず blocker として返した
- [ ] `codex exec` で Codex review conductor を呼び出した
- [ ] completion packet に final validation、Codex review、copilot_work_report を含めた

## Common Pitfalls

- [ ] final validation 前に scenario validation、suite-all、Sonar check を実行しなかった
- [ ] `first_action` 欠落を広い調査で補わなかった
- [ ] `parallelizable_with` がない handoff を並列実行しなかった
- [ ] repo-local Sonar issue gate と Sonar server Quality Gate を混同しなかった
- [ ] Codex review payload に diff、scope、validation result を含めた
- [ ] docs / workflow 文書変更を implementation lane に混ぜなかった
