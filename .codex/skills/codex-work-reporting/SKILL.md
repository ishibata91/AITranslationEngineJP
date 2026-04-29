---
name: codex-work-reporting
description: Codex 側の run 全体 reporting 作業プロトコル。Codex / Codex implementation レーン ベンチマーク値 と 完了根拠 から work_history レポート と次回改善事項を残す判断基準を提供する。
---
# Codex Work Reporting

## 目的

`codex-work-reporting` は作業プロトコルである。
Codex 作業流れ の完了、停止、戻し時に、`work_history` へ残す run 全体レポート 材料を整理する。
Codex と Codex implementation レーン の ベンチマーク値、完了根拠、検証結果 を同じ run 単位で集約する。
完了根拠 は明示 入力一式 だけでなく、Codex / Codex implementation レーン 会話ログ から 根拠参照 付きで抽出した完了報告も含む。

この skill は実行主体ではない。
ツール権限 は参照元 agent TOML に従い、完了条件と停止条件は参照元 skill に従う。

## 対応ロール

- 呼び出し元は 終了処理、停止、戻しを扱う Codex agent とする。
- 返却先は人間と `work_history` レポート とする。
- 担当成果物は `codex-work-reporting` の出力規約で固定する。

## 入力規約

- 不足時の扱い: 入力に 根拠参照、担当者、承認状態が不足する場合は推測で補わない。
- 必須入力: 実行記録 folder, 基準評価, Codex 根拠元, 実装根拠元
- 任意入力: Codex 完了根拠, 実装完了根拠, Codex 会話参照, 実装会話参照, 実装チャット参照, 関連計画, 関連引き継ぎ, 検証結果, 既知不足
- 必須 成果物: /Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/README.md, /Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/codex.md, /Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/codex.md

## 外部参照規約

- run index 雛形: [README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/README.md)
- Codex レポート 雛形: [codex.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/codex.md)
- Codex implementation レーン レポート 雛形: [codex.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/codex.md)
- 実行定義 agent: [work_reporter.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/work_reporter.toml)
- エージェント実行定義とツール権限は [work_reporter.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/work_reporter.toml) の 書き込み許可 / 実行許可 とする。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-work-reporting/SKILL.md

## 内部参照規約

### 拘束観点

- `work_history/templates/run/README.md` の run 全体 要約と benchmark block
- `work_history/templates/run/codex.md` の記入観点
- `work_history/templates/run/codex.md` の記入観点
- `analysis/benchmark-score.json` の session、metrics、scores、根拠参照
- 改善、時間、無駄、困りごとの分離
- Codex 固有の design、HITL、引き継ぎ、正本化判断の記録
- Codex implementation レーン 固有の 完了済み引き継ぎ、変更ファイル、検証、残留 の記録
- Codex implementation レーン 会話ログ / chat session file から 完了報告入力、final レポート、検証結果 を 根拠参照 付きで抽出する判断

### Benchmark Score

hook は使わない。
Codex / Codex implementation レーン の home 会話ログ を正本にし、スコア script で時間と摩擦の機械指標を出す。

Codex 会話ログ は `session_meta`、`event_msg`、`response_item` を読む。
Codex implementation レーン 会話ログ は `session.start`、`user.message`、`assistant.message`、`tool.execution_*` を読む。
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
速度の閾値や ベンチマーク値 欠落を初期 close 判定には使わない。

集計対象:
- `session_count`
- `metrics.duration_ms_total`
- `metrics.active_duration_ms_total`
- `metrics.user_turns`
- `metrics.assistant_turns`
- `metrics.tool_calls`
- `metrics.subagent_calls`
- `metrics.nonzero_tool_results`
- `metrics.人間修正回数`
- `metrics.long_idle_gaps`
- `metrics.repeated_tool_commands`
- `scores.time_cost`
- `scores.interaction_cost`
- `scores.tool_churn`
- `scores.rework_cost`
- `根拠参照`
- `transcript_gaps`

`time_cost` は `active_duration_ms_total` から算出する。
`long_idle_gaps` は 根拠 として残すが、評価値 には加算しない。
`人間修正回数` は人間の実メッセージだけを数え、skill 入力内容、下位 agent 通知、承認 レビュー用入力内容 は除外する。

## 判断規約

- `work_reporter` は最後に必ず run 全体レポート を作る。
- 置き場所は `work_history/runs/YYYY-MM-DD-<task-id>-run/` に固定する。
- 一次データは home の Codex / Codex implementation レーン 会話ログ、`analysis/benchmark-score.json`、完了根拠 とする。
- `README.md` は人間向け run 全体レポート と benchmark summary にする。
- `codex.md` と `codex.md` は `work_reporter` が 根拠 から生成する。
- 事実と判断材料を分ける。
- 分からない項目は `未確認`、`不明`、`なし` のいずれかで明示する。
- Codex implementation レーン 側の実装事実は、明示 Codex implementation レーン 完了根拠 または Codex implementation レーン 会話ログ / chat session file 内の完了報告からだけ転記する。
- Codex implementation レーン 会話ログ / chat session file を読む時は、`完了済み引き継ぎ`、`変更ファイル`、`test 結果`、`UI 根拠`、`Codex レビュー結果`、`レビュー結果一式`、`集約記録`、`harness 判定結果`、`完了根拠` の 完了報告入力 欄を探し、レポート 項目へ 根拠参照 を残す。
- 会話ログ / chat session 内に実装完了を示す文があっても、対象 task と紐づく 完了報告入力 欄、final レポート 欄、または 検証結果 欄を確認できない場合は推測扱いにする。
- ベンチマーク値 欠落、根拠参照 欠落、壊れた 会話ログ JSONL は次回改善 指摘 として扱う。
- 速度指標は改善観測であり、初期 close 判定には使わない。
- `.codex/history` には触れず、`work_history/` を使う。
- レポートは次回の prompt、引き継ぎ、雛形 改善へ戻せる粒度にする。

- `work_reporter` で run 全体レポート を作る
- `work_history/runs/YYYY-MM-DD-<task-id>-run/` を唯一の レポート 置き場所にする
- `analysis/benchmark-score.json` を agent が最初に読む材料として扱う
- 必要な時だけ `根拠参照` から 会話ログ 原文へ戻る
- Codex が実際に見た 根拠 と推測を分ける
- Codex implementation レーン facts は Codex implementation レーン 完了根拠 または Codex implementation レーン 会話ログ / chat session file の 根拠参照 付き抽出からだけ扱う
- 人間が次に見るべき path や コマンド を残す
- 重要エラーと未実行 検証 を短く明示する

## 非対象規約

- プロダクトコード、プロダクトテスト、docs 正本化は扱わない。
- docs 正本化の承認、対象範囲、implementation-scope を代替しない。
- `docs/exec-plans/`、`.codex/history/`、引き継ぎ file を run レポート置き場にしない。
- Markdown レポートをベンチマーク値の一次データにしない。
- 速度指標を初期 close 判定に使わない。

## 出力規約

- 基本出力: 出力は判断結果、根拠参照、不足情報、次 agent が判断できる材料を含む。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。
- レポート path: README.md、codex.md、codex.md、analysis/benchmark-score.json、transcript_refs.json の path を返す。
- benchmark summary: session count、metrics、scores、根拠 refs を分かる範囲で返す。不足は 阻害要因 ではなく次回改善事項にする。
- Codex レポート summary: Codex 完了根拠 または Codex 会話ログ 根拠参照 から確認できる結果、未完了、重要エラー、検証不足、次に見るべき場所を返す。
- Codex implementation レーン レポート summary: Codex implementation 完了根拠 または Codex implementation 会話ログ / chat session 根拠参照 から確認できる完了 引き継ぎ、変更ファイル、検証結果、統合 レビュー 結果、残留リスク、次に見るべき場所を返す。
- run 全体 指摘: 改善すべきこと、時間がかかったこと、無駄だったこと、困ったこと、検証で不足したことを返す。
- benchmark 品質 指摘: ベンチマーク値、根拠、レポート、実行定義、根拠参照、会話ログ の欠落または破損を次回改善事項として返す。
- 次回改善: prompt、引き継ぎ、雛形、benchmark scoring の改善を返す。
- 残留 不足: 未確認、不明、なしを区別して返す。

## 完了規約

- 出力規約を満たし、次の 実行者 が再解釈なしで判断できる。
- 不足情報または停止理由がある場合は明示されている。
- `work_reporter` が run 全体レポート を作った。
- `work_history/templates/run/README.md` の必須項目を確認した。
- `work_history/templates/run/codex.md` の必須項目を確認した。
- `analysis/benchmark-score.json` を run 全体ベンチマーク の入力として扱った。
- 改善、時間、無駄、困りごとを分けた。
- HITL、引き継ぎ、docs 正本化判断を記録対象にした。
- implementation factを 完了根拠 または Codex implementation 会話ログ / chat session file の 根拠参照 付き抽出からだけ扱った。
- 明示 完了根拠 が不足する場合は Codex implementation 会話ログ / chat session file を確認した。
- 必須根拠として、ベンチマーク値 json または不足理由、Codex 完了根拠 または Codex 会話ログ 根拠参照、Codex implementation 完了根拠 または Codex implementation 会話ログ / chat session 根拠参照、レポート 雛形 paths、利用可能な 検証結果 がある。
- 完了判断材料として、work_history/runs/<run>/README.md、codex.md、codex.md が ベンチマーク値 と 根拠 から生成され、次回改善事項が明示されている。
- 残留リスクとして、未確認または不明な 不足 が返っている。

## 停止規約

- プロダクトコードまたはプロダクトテスト を変更する時
- Codex implementation レーン 側 implementation レーン の事実を推測で補う時
- docs 正本化の承認や 対象範囲 を代替する時
- 速度の数値閾値で close 可否を判定する時
- 停止時は不足項目、衝突箇所、戻し先を返す。
- 拒否条件: 実行記録 folder 不足
- 拒否条件: Codex and implementation 根拠 sources both 不足
- 拒否条件: レポート write 対象 outside work_history/runs
- 停止条件: Codex implementation 完了根拠 is 不足 and Codex implementation 会話ログ / chat session source cannot be determined
- 停止条件: implementation facts cannot be distinguished from inference after checking available 完了根拠 and 会話ログ / chat session 根拠参照
- 停止条件: 必須 レポート path cannot be determined
- 停止条件: レポート generation would require プロダクト or docs mutation
