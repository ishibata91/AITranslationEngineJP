# Implementation Orchestrate Checklist

## Knowledge Check

- [ ] 承認済み implementation-scope と approval record を確認した
- [ ] distiller を tester / implementer より先に起動した
- [ ] tester を implementer より先に起動した
- [ ] suite-all と Sonar check を全 implementation handoff 完了後に実行した
- [ ] `codex exec` で Codex review conductor を呼び出した
- [ ] completion packet に final validation、Codex review、copilot_work_report を含めた

## Common Pitfalls

- [ ] final validation 前に suite-all または Sonar check を実行しなかった
- [ ] repo-local Sonar issue gate と Sonar server Quality Gate を混同しなかった
- [ ] Codex review payload に diff、scope、validation result を含めた
- [ ] docs / workflow 文書変更を implementation lane に混ぜなかった
