---
name: codex-work-reporting
description: Codex 側の run 全体レポート作業プロトコル。Codex / Codex implementation レーン ベンチマーク値 と レビュー差し戻しレポート から work_history レポート と次回改善事項を残す判断基準を提供する。
---
# Codex Work Reporting

## 目的

`codex-work-reporting` は作業プロトコルである。
Codex 作業流れ の完了、停止、戻し時に、`work_history` へ残す run 全体レポート 材料を整理する。
Codex と Codex implementation レーン の ベンチマーク値、レビュー差し戻し、検証結果 を同じ run 単位で集約する。
問題点の抽出は、ベンチマークスクリプトの出力と レビュー差し戻しレポート を主材料にする。
work_reporter は benchmark script 結果とレビュー差し戻しレポートだけを前提にする。

この skill は実行主体ではない。
実行境界は参照元 agent TOML に従い、完了条件と停止条件は参照元 skill に従う。

## 対応ロール

- 呼び出し元は 終了処理、停止、戻しを扱う Codex agent とする。
- 返却先は人間と `work_history` レポート とする。
- 担当成果物は `codex-work-reporting` の出力規約で固定する。

## 入力規約

- run 対象: run 全体レポートを作る `work_history/runs/YYYY-MM-DD-<task-id>-run/`。

## 外部参照規約

- run index 雛形: [README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/README.md)
- Codex レポート 雛形: [codex.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/codex.md)
- 実行定義 agent: [work_reporter.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/work_reporter.toml)
- エージェント実行定義と実行境界は [work_reporter.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/work_reporter.toml) に従う。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-work-reporting/SKILL.md

## 内部参照規約

### 拘束観点

- `work_history/templates/run/README.md` の run 全体 要約とベンチマーク欄
- `work_history/templates/run/codex.md` の記入観点
- `work_history/templates/run/codex.md` の記入観点
- `analysis/benchmark-score.json` の session、metrics、scores
- `work_history/runs/<run>/review-reject-*.md` のレビュー差し戻し出力
- 改善、時間、無駄、困りごとの分離
- Codex 固有の設計、人間介入、引き継ぎ、正本化判断の記録
- Codex implementation レーン 固有の 完了済み引き継ぎ、変更ファイル、検証、残留 の記録

### Benchmark Script

hook は使わない。
benchmark script は時間と摩擦の機械指標を出す。
work_reporter は benchmark script の出力だけを読む。
script は改善案、原因推定、責務判断を行わない。

生成物:
- `run-title.txt`
- `analysis/benchmark-score.json`

フォルダ名は最初の user prompt を安全化して作る。
同名フォルダがある場合は統合する。

### Benchmark

benchmark は次回改善用の観測値である。
速度の閾値や ベンチマーク値 欠落を初期終了判定には使わない。

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
- `input_gaps`

`time_cost` は `active_duration_ms_total` から算出する。
`long_idle_gaps` は 根拠 として残すが、評価値 には加算しない。
`人間修正回数` は人間の実メッセージだけを数え、skill 入力内容、下位 agent 通知、承認 レビュー用入力内容 は除外する。

## 判断規約

- `work_reporter` は最後に必ず run 全体レポート を作る。
- 置き場所は `work_history/runs/YYYY-MM-DD-<task-id>-run/` に固定する。
- 問題点の抽出は `analysis/benchmark-score.json` と `review-reject-*.md` を主材料にする。
- `README.md` は人間向け run 全体レポート と benchmark summary にする。
- `codex.md` は `work_reporter` が 根拠 から生成する。
- 事実と判断材料を分ける。
- 分からない項目は `未確認`、`不明`、`なし` のいずれかで明示する。
- Codex implementation レーン 側の実装事実は、run 内レポート、benchmark script 結果、レビュー差し戻しレポートから確認できる範囲だけ転記する。
- ベンチマーク値 欠落、レビュー差し戻しレポート 欠落、benchmark script 入力不足 は次回改善 指摘 として扱う。
- 速度指標は改善観測であり、初期終了判定には使わない。
- `.codex/history` には触れず、`work_history/` を使う。
- レポートは次回の指示、引き継ぎ、雛形 改善へ戻せる粒度にする。

- `work_reporter` で run 全体レポート を作る
- `work_history/runs/YYYY-MM-DD-<task-id>-run/` を唯一の レポート 置き場所にする
- `analysis/benchmark-score.json` を agent が最初に読む材料として扱う
- `review-reject-*.md` を レビュー差し戻しの問題抽出材料として扱う
- Codex が実際に見た 根拠 と推測を分ける
- Codex implementation レーンの事実は run 内レポート、benchmark script 結果、レビュー差し戻しレポートからだけ扱う
- 人間が次に見るべきパスや コマンド を残す
- 重要エラーと未実行 検証 を短く明示する

## 非対象規約

- プロダクトコード、プロダクトテスト、docs 正本化は扱わない。
- docs 正本化の承認、対象範囲、implementation-scope を代替しない。
- `docs/exec-plans/`、`.codex/history/`、引き継ぎファイルを run レポート置き場にしない。
- Markdown レポートをベンチマーク値の一次データにしない。
- 速度指標を初期終了判定に使わない。

## 出力規約

- 判断結果: run 全体レポートの完了、未完了、停止の判定を返す。
- 根拠参照: レポート生成に使った benchmark、レビュー差し戻し、検証結果を返す。
- 不足情報: レポート生成を完了できない不足項目を返す。
- 次判断材料: 人間または `implement_lane` が次を判断できる材料を返す。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコード変更の指示を含めない。
- レポートパス: README.md、codex.md、analysis/benchmark-score.json、review-reject-*.md のパスまたは未作成確認を返す。
- benchmark summary: session count、metrics、scores を分かる範囲で返す。不足は 阻害要因 ではなく次回改善事項にする。
- Codex レポート summary: run 内レポート、benchmark script 結果、レビュー差し戻しレポートから確認できる結果、未完了、重要エラー、検証不足、次に見るべき場所を返す。
- Codex implementation レーン レポート summary: run 内レポート、benchmark script 結果、レビュー差し戻しレポートから確認できる完了 引き継ぎ、変更ファイル、検証結果、統合 レビュー 結果、残留リスク、次に見るべき場所を返す。
- run 全体 指摘: 改善すべきこと、時間がかかったこと、無駄だったこと、困ったこと、検証で不足したことを返す。
- benchmark 品質 指摘: ベンチマーク値、レポート、実行定義、script 入力、script 出力 の欠落または破損を次回改善事項として返す。
- 次回改善: 指示、引き継ぎ、雛形、ベンチマーク採点の改善を返す。
- 残留 不足: 未確認、不明、なしを区別して返す。

## 完了規約

- 出力規約を満たし、次の 実行者 が再解釈なしで判断できる。
- 不足情報または停止理由がある場合は明示されている。
- `work_reporter` が run 全体レポート を作った。
- `work_history/templates/run/README.md` の必須項目を確認した。
- `work_history/templates/run/codex.md` の必須項目を確認した。
- `analysis/benchmark-score.json` を run 全体ベンチマーク の入力として扱った。
- `review-reject-*.md` を確認し、レビュー差し戻し出力から問題点を抽出した。
- 問題点はベンチマークスクリプト出力とレビュー差し戻しレポートから抽出した。
- 改善、時間、無駄、困りごとを分けた。
- 人間介入、引き継ぎ、docs 正本化判断を記録対象にした。
- implementation レーンの事実を run 内レポート、benchmark script 結果、レビュー差し戻しレポートからだけ扱った。
- 必須根拠として、ベンチマーク値 json または不足理由、レビュー差し戻しレポートまたは未作成確認、レポート 雛形 paths、利用可能な 検証結果 がある。
- 完了判断材料として、work_history/runs/<run>/README.md と codex.md が ベンチマーク値 と 根拠 から生成され、次回改善事項が明示されている。
- 残留リスクとして、未確認または不明な 不足 が返っている。

## 停止規約

- プロダクトコードまたはプロダクトテスト を変更する時
- Codex implementation レーン 側 implementation レーン の事実を推測で補う時
- docs 正本化の承認や 対象範囲 を代替する時
- 速度の数値閾値で終了可否を判定する時
- 停止時は不足項目、衝突箇所、戻し先を返す。
- run 対象が不足する場合は停止する。
- benchmark とレビュー差し戻しレポートの有無を確認できない場合は停止する。
- レポート書き込み先が `work_history/runs` 外になる場合は停止する。
- implementation レーンの事実と推測を区別できない場合は停止する。
- 必須レポートパスを特定できない場合は停止する。
- レポート生成にプロダクトまたは docs 正本の変更が必要な場合は停止する。
