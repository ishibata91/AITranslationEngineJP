---
name: codex-work-reporting
description: Codex 側の1ラン報告知識 package。work_history の Codex report とラン横断 finding を残す判断基準を提供する。
---

# Codex Work Reporting

## 目的

`codex-work-reporting` は知識 package である。
Codex workflow の完了、停止、reroute 時に、`work_history` へ残す Codex 側の報告材料を整理する。

この skill は実行主体ではない。
書き込み権限、完了条件、停止条件は参照元 agent の contract に従う。

## いつ参照するか

- Codex run の closeout、停止、reroute を記録する時
- Copilot completion report を受けて work_history に残す材料を整理する時
- `codex.md` とラン `README.md` の記入観点を確認する時

## 参照しない場合

- product code または product test を変更する時
- Copilot 側 implementation lane の事実を推測で補う時
- docs 正本化の承認や scope を代替する時

## 知識範囲

- `work_history/templates/run/codex.md` の記入観点
- `work_history/templates/run/README.md` の横断 finding
- 改善、時間、無駄、困りごとの分離
- Codex 固有の design、HITL、handoff、正本化判断の記録

## 原則

- `propose_plans` は最後に必ず Codex report 材料を作る。
- 置き場所は `work_history/runs/YYYY-MM-DD-<task-id>-run/codex.md` に固定する。
- 事実と判断材料を分ける。
- 分からない項目は `未確認`、`不明`、`なし` のいずれかで明示する。
- Copilot 側の実装事実は Copilot completion report からだけ転記する。
- `.codex/history` には触れず、`work_history/` を使う。
- レポートは次回の prompt、handoff、template 改善へ戻せる粒度にする。

## 標準パターン

1. `work_history/runs/YYYY-MM-DD-<task-id>-run/` があるか確認する。
2. なければ `work_history/templates/run/` を複製して run folder を作る。
3. Codex 側の報告は run folder の `codex.md` だけに書く。
4. 両 lane の比較は run folder の `README.md` だけに書く。
5. `docs/exec-plans/`、`.codex/history/`、handoff file には run report を置かない。
6. HITL、handoff packet、docs 正本化判断、spawn / 調査の必要判定を記録する。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は agent contract に従う。

## DO / DON'T

DO:
- `propose_plans` の最後に必ず Codex report 材料を作る
- `work_history/runs/YYYY-MM-DD-<task-id>-run/codex.md` を唯一の Codex report 置き場所にする
- Codex が実際に見た evidence と推測を分ける
- 人間が次に見るべき path や command を残す
- 重要エラーと未実行 validation を短く明示する

DON'T:
- Copilot の作業時間や実装内容を推測で埋めない
- docs 正本化や implementation-scope の代わりにしない
- `docs/exec-plans/`、`.codex/history/`、handoff file に run report を置かない
- `.codex/history` へ移行や参照ルールを追加しない
- 長い経緯説明や感想を増やさない

## Checklist

- [codex-work-reporting-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-work-reporting/references/checklists/codex-work-reporting-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。

## References

- run index template: [README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/README.md)
- Codex report template: [codex.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/codex.md)

## Agent が持つもの

- 実行権限
- agent 1:1 contract
- write scope
- stop / reroute 条件

## Maintenance

- 権限や contract を skill 本体へ戻さない。
- template 変更時は checklist の観点も同期する。
- Copilot lane 固有の実装観点は `copilot-work-reporting` 側へ置く。
