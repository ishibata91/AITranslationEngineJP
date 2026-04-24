---
name: copilot-work-reporting
description: GitHub Copilot 側の1ラン報告知識 package。work_history の Copilot report に必要な implementation lane 事実を completion packet へ残す判断基準を提供する。
---

# Copilot Work Reporting

## 目的

`copilot-work-reporting` は知識 package である。
GitHub Copilot の implementation lane が、1ランの問題点、時間、無駄、困りごとを completion packet へ残すための判断基準を提供する。

この skill は実行主体ではない。
file mutation、RunSubagent、完了条件、停止条件は参照元 agent の contract に従う。

## いつ参照するか

- implementation-orchestrate が completion packet を作る時
- tester、implementer、investigator、final validation、Codex review の戻り値から報告材料を集約する時
- Codex が `work_history/templates/run/copilot.md` へ転記できる材料を返す時

## 参照しない場合

- implementation-scope を変更する時
- docs、`.codex`、`.github/skills`、`.github/agents` を変更する時
- オーケストレーター自身が直接実装、調査、test、review を行う時

## 知識範囲

- `work_history/templates/run/copilot.md` の記入観点
- implementation-scope の読み取り、実装分割、調査、test、final validation、Codex review、reroute の記録
- touched files、validation、完了報告不足、report skeleton の整理
- Codex が work_history に転記するための report-ready summary

## 原則

- implementation-orchestrate は最後に必ず `copilot_work_report` を作らせる。
- 置き場所は `work_history/runs/YYYY-MM-DD-<task-id>-run/copilot.md` に固定する。
- completion packet 直下の `copilot_work_report` は report-ready skeleton として返す。
- subagent 戻り値だけから報告材料を集約する。
- 推測で作業時間、変更 file、validation 結果を補わない。
- `docs_changes` は implementation-orchestrate contract に従い `none` を返す。
- file に直接書かず、completion packet に `copilot_work_report` として返す。
- 分からない項目は `未確認`、`不明`、`なし` のいずれかで明示する。

## 標準パターン

1. report path を `work_history/runs/YYYY-MM-DD-<task-id>-run/copilot.md` に固定する。
2. file には直接書かず、completion packet の `copilot_work_report.report_path` に固定 path を返す。
3. `copilot_work_report` に `report_path`、`status`、`改善すべきこと`、`時間がかかったこと`、`無駄だったこと`、`困ったこと`、`次に見るべき場所` を同じ順序で返す。
4. completed handoff、touched files、implemented scope、test results、final validation result、Codex review result を completion packet から集める。
5. implementation-scope の読み取り、実装分割、調査、validation、reroute を整理する。
6. `docs/exec-plans/`、`.codex/history/`、handoff file を report path にしない。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は agent contract に従う。

## DO / DON'T

DO:
- implementation-orchestrate の最後に必ず `copilot_work_report` を作らせる
- `copilot_work_report.report_path` を `work_history/runs/YYYY-MM-DD-<task-id>-run/copilot.md` に固定する
- report skeleton の field 順序を固定する
- tester、implementer、investigator、final validation、Codex review の戻り値を根拠にする
- reroute と blocked reason を report 材料として残す
- 未実行 validation は未実行理由と一緒に書く
- Codex が `copilot.md` へ転記できる短い粒度にする

DON'T:
- オーケストレーター自身で file read / search / edit を行わない
- implementation-scope の不足を実装判断で補わない
- `docs/exec-plans/`、`.codex/history/`、handoff file を report path にしない
- docs 正本化や workflow 変更を implementation lane に混ぜない
- 空欄や曖昧な成功報告だけで completion packet を返さない

## Checklist

- [copilot-work-reporting-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/copilot-work-reporting/references/checklists/copilot-work-reporting-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。

## References

- Copilot report template: [copilot.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/copilot.md)
- run index template: [README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/README.md)

## Agent が持つもの

- 実行権限
- agent 1:1 contract
- handoff 契約
- stop / reroute 条件

## Maintenance

- 権限や contract を skill 本体へ戻さない。
- completion packet の field 名を変える時は implementation-orchestrate contract と同期する。
- Codex lane 固有の設計、HITL、正本化判断は `codex-work-reporting` 側へ置く。
