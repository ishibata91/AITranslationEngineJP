---
name: codex-work-reporting
description: Codex 側の run-wide reporting 知識 package。Codex / Copilot telemetry と completion evidence から work_history report と次回改善 finding を残す判断基準を提供する。
---

# Codex Work Reporting

## 目的

`codex-work-reporting` は知識 package である。
Codex workflow の完了、停止、reroute 時に、`work_history` へ残す run-wide report 材料を整理する。
Codex と Copilot の telemetry、completion evidence、validation result を同じ run 単位で集約する。

この skill は実行主体ではない。
書き込み権限、完了条件、停止条件は参照元 agent の contract に従う。

## いつ参照するか

- Codex run の closeout、停止、reroute を記録する時
- Copilot completion evidence を受けて `copilot.md` へ転記する時
- `telemetry.jsonl` から run-wide benchmark を作る時
- `README.md`、`codex.md`、`copilot.md` の記入観点を確認する時

## 参照しない場合

- product code または product test を変更する時
- Copilot 側 implementation lane の事実を推測で補う時
- docs 正本化の承認や scope を代替する時
- 速度の数値閾値で close 可否を判定する時

## 知識範囲

- `work_history/templates/run/README.md` の run-wide 要約と benchmark block
- `work_history/templates/run/codex.md` の記入観点
- `work_history/templates/run/copilot.md` の記入観点
- `telemetry.jsonl` の assistant response event
- 改善、時間、無駄、困りごとの分離
- Codex 固有の design、HITL、handoff、正本化判断の記録
- Copilot 固有の completed_handoffs、touched_files、validation、residual の記録

## 原則

- `work_reporter` は最後に必ず run-wide report を作る。
- 置き場所は `work_history/runs/YYYY-MM-DD-<task-id>-run/` に固定する。
- 一次データは `telemetry.jsonl` と completion evidence とする。
- `README.md` は人間向け run-wide report と benchmark summary にする。
- `codex.md` と `copilot.md` は `work_reporter` が evidence から生成する。
- 事実と判断材料を分ける。
- 分からない項目は `未確認`、`不明`、`なし` のいずれかで明示する。
- Copilot 側の実装事実は Copilot completion evidence からだけ転記する。
- telemetry 欠落、event field 欠落、壊れた JSONL は次回改善 finding として扱う。
- 速度指標は改善観測であり、初期 close 判定には使わない。
- `.codex/history` には触れず、`work_history/` を使う。
- レポートは次回の prompt、handoff、template 改善へ戻せる粒度にする。

## 標準パターン

1. `work_history/runs/YYYY-MM-DD-<task-id>-run/` があるか確認する。
2. なければ `work_history/templates/run/` を複製して run folder を作る。
3. `telemetry.jsonl` を読み、`runtime: codex | copilot` ごとの response 数と elapsed を集計する。
4. Codex completion evidence から `codex.md` を作る。
5. Copilot completion evidence から `copilot.md` を作る。
6. 両 lane の比較と benchmark summary は `README.md` へ集約する。
7. `docs/exec-plans/`、`.codex/history/`、handoff file には run report を置かない。
8. telemetry / evidence 欠落は close 判定ではなく、次回改善 finding にする。

telemetry 集計は必要なら次の helper を使う。

```bash
python3 scripts/work-history/aggregate_telemetry.py work_history/runs/<run>/telemetry.jsonl
```

この手順は知識上の標準例である。
実行順、必須 input、完了条件は agent contract に従う。

## Telemetry Event

初期 event は assistant response 単位で扱う。
event は OpenTelemetry 風の trace / event / aggregate に寄せるが、外部 service への送信は必須にしない。

必須 field:
- `event_type`: `assistant_response`
- `run_id`
- `runtime`: `codex` または `copilot`
- `elapsed_ms_from_run_start`
- `phase`
- `status`
- `mechanical_summary`

任意 field:
- `turn_index`
- `related_paths`
- `commands`
- `blocked_reason`

## Benchmark

benchmark は次回改善用の観測値である。
速度の閾値や event 欠落を初期 close 判定には使わない。

集計対象:
- `response_count_by_runtime`
- `elapsed_ms_total`
- `phase_elapsed_ms`
- `runtime_elapsed_ms`
- `blocked_elapsed_ms`
- `reroute_count`
- `validation_elapsed_ms`

## DO / DON'T

DO:
- `work_reporter` で run-wide report を作る
- `work_history/runs/YYYY-MM-DD-<task-id>-run/` を唯一の report 置き場所にする
- `telemetry.jsonl` を機械集計の一次データとして扱う
- Codex が実際に見た evidence と推測を分ける
- Copilot facts は Copilot completion evidence からだけ扱う
- 人間が次に見るべき path や command を残す
- 重要エラーと未実行 validation を短く明示する

DON'T:
- Copilot の作業時間や実装内容を推測で埋めない
- docs 正本化や implementation-scope の代わりにしない
- `docs/exec-plans/`、`.codex/history/`、handoff file に run report を置かない
- `.codex/history` へ移行や参照ルールを追加しない
- Markdown report を telemetry の一次データにしない
- 速度指標を初期 close 判定に使わない
- 長い経緯説明や感想を増やさない

## Checklist

- [codex-work-reporting-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-work-reporting/references/checklists/codex-work-reporting-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。

## References

- run index template: [README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/README.md)
- Codex report template: [codex.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/codex.md)
- Copilot report template: [copilot.md](/Users/iorishibata/Repositories/AITranslationEngineJP/work_history/templates/run/copilot.md)
- runtime agent: [work_reporter.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/work_reporter.toml)
- telemetry aggregator: [aggregate_telemetry.py](/Users/iorishibata/Repositories/AITranslationEngineJP/scripts/work-history/aggregate_telemetry.py)

## Agent が持つもの

- 実行権限
- agent 1:1 contract
- write scope
- stop / reroute 条件

## Maintenance

- 権限や contract を skill 本体へ戻さない。
- template 変更時は checklist の観点も同期する。
- Copilot lane 固有の実装事実は completion evidence から受ける。
- Copilot 側に report 作成責務を戻さない。
