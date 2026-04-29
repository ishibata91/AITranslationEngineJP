---
name: codex-work-reporting
description: Codex 側の run-wide reporting 作業プロトコル。Codex / Codex implementation lane benchmark score と completion evidence から work_history report と次回改善 finding を残す判断基準を提供する。
---
# Codex Work Reporting

## 目的

`codex-work-reporting` は作業プロトコルである。
Codex workflow の完了、停止、reroute 時に、`work_history` へ残す run-wide report 材料を整理する。
Codex と Codex implementation lane の benchmark score、completion evidence、validation result を同じ run 単位で集約する。
completion evidence は明示 packet だけでなく、Codex / Codex implementation lane transcript から source_ref 付きで抽出した完了報告も含む。

この skill は実行主体ではない。
tool policy は参照元 agent TOML に従い、完了条件と停止条件は参照元 skill に従う。

## 対応ロール

- 呼び出し元は closeout、停止、reroute を扱う Codex agent とする。
- 返却先は人間と `work_history` report とする。
- owner artifact は `codex-work-reporting` の出力規約で固定する。

## 入力規約

- Codex run の closeout、停止、reroute を記録する時
- Codex implementation lane completion evidence を受けて、または Codex implementation lane transcript から抽出して `codex.md` へ転記する時
- `analysis/benchmark-score.json` から run-wide benchmark を作る時
- `README.md`、`codex.md`、`codex.md` の記入観点を確認する時
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。
- 必須入力: run_folder, benchmark_score, codex_evidence_source, implementation_evidence_source
- 任意入力: codex_completion_evidence, implementation_completion_evidence, codex_transcript_refs, implementation_transcript_refs, implementation_chat_session_refs, related_plan, related_handoff, validation_results, known_gaps
- 必須 artifact: /Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/README.md, /Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/codex.md, /Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/codex.md

## 外部参照規約

- run index template: [README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/README.md)
- Codex report template: [codex.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/codex.md)
- Codex implementation lane report template: [codex.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/codex.md)
- runtime agent: [work_reporter.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/work_reporter.toml)
- agent runtime と tool policy は [work_reporter.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/work_reporter.toml) の `allowed_write_paths` / `allowed_commands` とする。
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-work-reporting/SKILL.md

## 内部参照規約

### 拘束観点

- `work_history/templates/run/README.md` の run-wide 要約と benchmark block
- `work_history/templates/run/codex.md` の記入観点
- `work_history/templates/run/codex.md` の記入観点
- `analysis/benchmark-score.json` の session、metrics、scores、evidence_refs
- 改善、時間、無駄、困りごとの分離
- Codex 固有の design、HITL、handoff、正本化判断の記録
- Codex implementation lane 固有の completed_handoffs、touched_files、validation、residual の記録
- Codex implementation lane transcript / chat session file から completion packet、final report、validation result を source_ref 付きで抽出する判断

### Benchmark Score

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

### Benchmark

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

## 判断規約

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

- `work_reporter` で run-wide report を作る
- `work_history/runs/YYYY-MM-DD-<task-id>-run/` を唯一の report 置き場所にする
- `analysis/benchmark-score.json` を agent が最初に読む材料として扱う
- 必要な時だけ `source_ref` から transcript 原文へ戻る
- Codex が実際に見た evidence と推測を分ける
- Codex implementation lane facts は Codex implementation lane completion evidence または Codex implementation lane transcript / chat session file の source_ref 付き抽出からだけ扱う
- 人間が次に見るべき path や command を残す
- 重要エラーと未実行 validation を短く明示する

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力に tool policy、agent runtime、product code の変更義務を含めない。
- 必須出力: run_report_paths, benchmark_summary, codex_report_summary, implementation_report_summary, cross_role_findings, benchmark_quality_findings, next_improvements, residual_gaps
- 出力 field 要件: {"run_report_paths": "README.md、codex.md、codex.md、analysis/benchmark-score.json、transcript_refs.json の path を返す", "benchmark_summary": "session_count、metrics、scores、evidence_refs を分かる範囲で返す。不足は blocker ではなく improvement finding にする", "codex_report_summary": "Codex completion evidence または Codex transcript source_ref から確認できる結果、未完了、重要エラー、検証不足、次に見るべき場所を返す", "implementation_report_summary": "Codex implementation completion evidence または Codex implementation transcript / chat session source_ref から確認できる completed_handoffs、touched_files、validation_result、integrated_review_result、residual_risks、次に見るべき場所を返す。推測で補わない", "cross_role_findings": "改善すべきこと、時間がかかったこと、無駄だったこと、困ったこと、検証で不足したことを run-wide に返す", "benchmark_quality_findings": "benchmark score / evidence / report 欠落、壊れた JSONL、runtime 欠落、source_ref 欠落、読めない transcript を次回改善 finding として返す", "next_improvements": "prompt、handoff、template、benchmark scoring の次回改善を返す", "residual_gaps": "未確認、不明、なしを区別して返す"}

## 完了規約

- 出力規約を満たし、次の actor が再解釈なしで判断できる。
- 不足情報または停止理由がある場合は明示されている。
- `work_reporter` が run-wide report を作った。
- `work_history/templates/run/README.md` の必須項目を確認した。
- `work_history/templates/run/codex.md` の必須項目を確認した。
- `analysis/benchmark-score.json` を run-wide benchmark の入力として扱った。
- 改善、時間、無駄、困りごとを分けた。
- HITL、handoff、docs 正本化判断を記録対象にした。
- implementation factを completion evidence または Codex implementation transcript / chat session file の source_ref 付き抽出からだけ扱った。
- 明示 completion evidence が不足する場合は Codex implementation transcript / chat session file を確認した。
- 必須 evidence: benchmark score json or benchmark score missing reason, Codex completion evidence or Codex transcript source_ref, Codex implementation completion evidence or Codex implementation transcript / chat session source_ref, report template paths, validation result when available
- completion signal: work_history/runs/<run>/README.md、codex.md、codex.md が benchmark score と evidence から生成され、次回改善 finding が明示されている
- residual risk key: residual_gaps

## 停止規約

- product code または product test を変更する時
- Codex implementation lane 側 implementation lane の事実を推測で補う時
- docs 正本化の承認や scope を代替する時
- 速度の数値閾値で close 可否を判定する時
- Codex implementation lane の作業時間や実装内容を推測で埋めない
- Codex implementation lane transcript / chat session file を読めるのに読まず、completion packet 未提示だけで close blocker にしない
- docs 正本化や implementation-scope の代わりにしない
- `docs/exec-plans/`、`.codex/history/`、handoff file に run report を置かない
- `.codex/history` へ移行や参照ルールを追加しない
- Markdown report を benchmark score の一次データにしない
- 速度指標を初期 close 判定に使わない
- 長い経緯説明や感想を増やさない
- 停止時は不足項目、衝突箇所、reroute 先を返す。
- `.codex/history` へ記録先を戻さなかった場合は停止する。
- 推測で Codex implementation laneの実装事実を補わなかった場合は停止する。
- レポートを docs 正本や implementation-scope の代替にしなかった場合は停止する。
- Markdown report を benchmark score の一次データにしなかった場合は停止する。
- 速度指標を初期 close 判定に使わなかった場合は停止する。
- 拒否条件: run_folder missing
- 拒否条件: Codex and implementation evidence sources both missing
- 拒否条件: report write target outside work_history/runs
- 停止条件: Codex implementation completion evidence is missing and Codex implementation transcript / chat session source cannot be determined
- 停止条件: implementation facts cannot be distinguished from inference after checking available completion evidence and transcript / chat session source_ref
- 停止条件: required report path cannot be determined
- 停止条件: report generation would require product or docs mutation
