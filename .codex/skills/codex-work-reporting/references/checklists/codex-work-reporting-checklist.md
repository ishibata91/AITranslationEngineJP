# Codex Work Reporting Checklist

## Knowledge Check

- [ ] `work_reporter` が run-wide report を作った
- [ ] `work_history/templates/run/README.md` の必須項目を確認した
- [ ] `work_history/templates/run/codex.md` の必須項目を確認した
- [ ] `work_history/templates/run/copilot.md` の必須項目を確認した
- [ ] `telemetry.jsonl` を run-wide benchmark の一次データとして扱った
- [ ] 改善、時間、無駄、困りごとを分けた
- [ ] HITL、handoff、docs 正本化判断を記録対象にした
- [ ] Copilot 事実を completion evidence からだけ扱った

## Common Pitfalls

- [ ] `.codex/history` へ記録先を戻さなかった
- [ ] 推測で Copilot 側の実装事実を補わなかった
- [ ] レポートを docs 正本や implementation-scope の代替にしなかった
- [ ] Markdown report を telemetry の一次データにしなかった
- [ ] 速度指標を初期 close 判定に使わなかった
