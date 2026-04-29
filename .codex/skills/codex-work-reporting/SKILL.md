---
name: codex-work-reporting
description: Codex 側の run-wide reporting 知識 package。Codex / Codex implementation lane benchmark score と completion evidence から work_history report と次回改善 finding を残す判断基準を提供する。
---

# Codex Work Reporting

## 目的

`codex-work-reporting` は知識 package である。
Codex workflow の完了、停止、reroute 時に、`work_history` へ残す run-wide report 材料を整理する。
Codex と Codex implementation lane の benchmark score、completion evidence、validation result を同じ run 単位で集約する。
completion evidence は明示 packet だけでなく、Codex / Codex implementation lane transcript から source_ref 付きで抽出した完了報告も含む。

この skill は実行主体ではない。
tool policy は参照元 agent TOML に従い、完了条件と停止条件は参照元 agent の contract に従う。

## いつ参照するか

- Codex run の closeout、停止、reroute を記録する時
- Codex implementation lane completion evidence を受けて、または Codex implementation lane transcript から抽出して `codex.md` へ転記する時
- `analysis/benchmark-score.json` から run-wide benchmark を作る時
- `README.md`、`codex.md`、`codex.md` の記入観点を確認する時

## 参照しない場合

- product code または product test を変更する時
- Codex implementation lane 側 implementation lane の事実を推測で補う時
- docs 正本化の承認や scope を代替する時
- 速度の数値閾値で close 可否を判定する時

## 知識範囲

- `work_history/templates/run/README.md` の run-wide 要約と benchmark block
- `work_history/templates/run/codex.md` の記入観点
- `work_history/templates/run/codex.md` の記入観点
- `analysis/benchmark-score.json` の session、metrics、scores、evidence_refs
- 改善、時間、無駄、困りごとの分離
- Codex 固有の design、HITL、handoff、正本化判断の記録
- Codex implementation lane 固有の completed_handoffs、touched_files、validation、residual の記録
- Codex implementation lane transcript / chat session file から completion packet、final report、validation result を source_ref 付きで抽出する判断

## 原則

- `work_reporter` は最後に必ず run-wide report を作る。
- 置き場所は `work_history/runs/YYYY-MM-DD-<task-id>-run/` に固定する。
- 一次データは home の Codex / Codex implementation lane transcript、`analysis/benchmark-score.json`、completion evidence とする。
- `README.md` は人間向け run-wide report と benchmark summary にする。
- `codex.md` と `codex.md` は `work_reporter` が evidence から生成する。
- 事実と判断材料を分ける。
- 分からない項目は `未確認`、`不明`、`なし` のいずれかで明示する。
- Codex implementation lane 側の実装事実は、明示 Codex implementation lane completion evidence または Codex implementation lane transcript / chat session file 内の完了報告からだけ転記する。
- Codex implementation lane transcript / chat session file を読む時は、`completed_handoffs`、`touched_files`、`test_results`、`ui_evidence`、`codex_review_result`、`reviewer_result_bundle`、`aggregation_trace`、`harness_gate_result`、`completion_evidence` の completion packet 欄を探し、report 項目へ source_ref を残す。
- transcript / chat session 内に実装完了を示す文があっても、対象 task と紐づく completion packet 欄、final report 欄、または validation result 欄を確認できない場合は推測扱いにする。
- benchmark score 欠落、source_ref 欠落、壊れた transcript JSONL は次回改善 finding として扱う。
- 速度指標は改善観測であり、初期 close 判定には使わない。
- `.codex/history` には触れず、`work_history/` を使う。
- レポートは次回の prompt、handoff、template 改善へ戻せる粒度にする。

## 標準パターン

1. `work_history/runs/YYYY-MM-DD-<task-id>-run/` があるか確認する。
2. なければ `work_history/templates/run/` を複製して run folder を作る。
3. 必要なら `scripts/work-history/score_transcripts.py` で Codex / Codex implementation lane transcript を score 化する。
4. `analysis/benchmark-score.json` を先に読む。
5. 明示 completion evidence が不足する lane は、`transcript_refs.json` または caller-provided transcript / chat session path から session file を読む。
6. Codex implementation lane transcript / chat session file では task id、completion packet、final validation、integrated review result、gate result を検索し、確認できた項目だけを `codex.md` へ転記する。
7. Codex completion evidence または Codex transcript から `codex.md` を作る。
8. Codex implementation lane completion evidence または Codex implementation lane transcript / chat session file から `codex.md` を作る。
9. 両 lane の比較と benchmark summary は `README.md` へ集約する。
10. `docs/exec-plans/`、`.codex/history/`、handoff file には run report を置かない。
11. benchmark score / evidence 欠落は close 判定ではなく、次回改善 finding にする。

benchmark score は必要なら次の helper を使う。

```bash
python3 scripts/work-history/score_transcripts.py \
  --codex-transcript /Users/<user>/.codex/sessions/YYYY/MM/DD/<session>.jsonl \
  --implementation-transcript "/Users/<user>/Library/Application Support/Code/User/workspaceStorage/<workspace>/codex/session-or-transcript/<session>.jsonl" \
  --output-root work_history/runs
```

この手順は知識上の標準例である。
実行順、必須 input、完了条件は agent contract に従う。

## Benchmark Score

hook は使わない。
Codex / Codex implementation lane の home transcript を正本にし、スコア script で時間と摩擦の機械指標を出す。

Codex transcript は `session_meta`、`event_msg`、`response_item` を読む。
Codex implementation lane transcript は `session.start`、`user.message`、`assistant.message`、`tool.execution_*` を読む。
VS Code `chatSessions/*.jsonl` は `requests` 配列、`message.text`、`response`、`toolCallRounds` を読む。
script は時刻、turn 数、tool 数、非 0 終了、長い idle、再実行、user correction を数える。
script は改善案、原因推定、責務判断は行わない。

生成物:
- `run-title.txt`
- `transcript_refs.json`
- `analysis/benchmark-score.json`

folder 名は最初の user prompt を安全化して作る。
同名 folder がある場合は merge し、`transcript_refs.json` に session を追記する。

## Benchmark

benchmark は次回改善用の観測値である。
速度の閾値や benchmark score 欠落を初期 close 判定には使わない。

集計対象:
- `session_count`
- `metrics.duration_ms_total`
- `metrics.active_duration_ms_total`
- `metrics.user_turns`
- `metrics.assistant_turns`
- `metrics.tool_calls`
- `metrics.subagent_calls`
- `metrics.nonzero_tool_results`
- `metrics.user_corrections`
- `metrics.long_idle_gaps`
- `metrics.repeated_tool_commands`
- `scores.time_cost`
- `scores.interaction_cost`
- `scores.tool_churn`
- `scores.rework_cost`
- `evidence_refs`
- `transcript_gaps`

`time_cost` は `active_duration_ms_total` から算出する。
`long_idle_gaps` は evidence として残すが、score には加算しない。
`user_corrections` は人間の実メッセージだけを数え、skill payload、subagent notification、approval reviewer 用 payload は除外する。

## DO / DON'T

DO:
- `work_reporter` で run-wide report を作る
- `work_history/runs/YYYY-MM-DD-<task-id>-run/` を唯一の report 置き場所にする
- `analysis/benchmark-score.json` を agent が最初に読む材料として扱う
- 必要な時だけ `source_ref` から transcript 原文へ戻る
- Codex が実際に見た evidence と推測を分ける
- Codex implementation lane facts は Codex implementation lane completion evidence または Codex implementation lane transcript / chat session file の source_ref 付き抽出からだけ扱う
- 人間が次に見るべき path や command を残す
- 重要エラーと未実行 validation を短く明示する

DON'T:
- Codex implementation lane の作業時間や実装内容を推測で埋めない
- Codex implementation lane transcript / chat session file を読めるのに読まず、completion packet 未提示だけで close blocker にしない
- docs 正本化や implementation-scope の代わりにしない
- `docs/exec-plans/`、`.codex/history/`、handoff file に run report を置かない
- `.codex/history` へ移行や参照ルールを追加しない
- Markdown report を benchmark score の一次データにしない
- 速度指標を初期 close 判定に使わない
- 長い経緯説明や感想を増やさない

## Checklist

- [codex-work-reporting-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-work-reporting/references/checklists/codex-work-reporting-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。

## References

- run index template: [README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/README.md)
- Codex report template: [codex.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/codex.md)
- Codex implementation lane report template: [codex.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/codex.md)
- runtime agent: [work_reporter.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/work_reporter.toml)
- benchmark scorer: [score_transcripts.py](/Users/iorishibata/Repositories/AITranslationEngineJP/scripts/work-history/score_transcripts.py)

## Agent が持つもの

- tool policy
- agent 1:1 contract
- tool policy
- stop / reroute 条件

## Maintenance

- tool policy や contract を skill 本体へ戻さない。
- template 変更時は checklist の観点も同期する。
- Codex implementation lane 固有の実装事実は completion evidence から受ける。
- completion evidence が明示入力にない場合は、Codex implementation lane transcript / chat session file から task id と completion packet 欄を確認し、source_ref 付きで抽出する。
- Codex implementation lane 側に report 作成責務を戻さない。
