# Copilot Work Reporting Checklist

## Knowledge Check

- [ ] implementation-orchestrate の最後に必ず `copilot_work_report` を作らせた
- [ ] `work_history/templates/run/copilot.md` の必須項目を確認した
- [ ] subagent 戻り値だけから report 材料を集約した
- [ ] 改善、時間、無駄、困りごとを分けた
- [ ] validation 未実行理由と reroute reason を残した

## Common Pitfalls

- [ ] オーケストレーター自身で file mutation しなかった
- [ ] implementation-scope の不足を推測で補わなかった
- [ ] docs 正本化や workflow 変更を implementation lane に混ぜなかった
