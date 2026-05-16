# Merge Prep Input

- `skill`: `implement-lane`
- `status`: `ready_for_merge_lane`
- `target_lane`: `merge-lane`
- `branch`: `codex/translation-job-state-stale-retirement`
- `active_plan`: `docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/`

## 完了状態

実装レーンとして必要な設計、実装、docs 正本化、レビュー、最終検証、作業履歴を完了した。
local merge、active plan の completed 移動、merge 結果 commit は merge lane の担当として残す。

## 主要成果物

- `branch-preparation.md`: branch 判断。
- `scenario-design.md`: stale 廃止の scenario 設計。
- `design-diff.md`: 設計差分。
- `implementation-scope.md`: 実装範囲。
- `implementation-handoff.md`: backend 実装引き継ぎ。
- `implementation-wave-result.md`: backend 実装結果。
- `docs-canonicalization-result.md`: docs 正本化結果。
- `canonicalization-decision.md`: docs 正本化判断。
- `review-aggregation.md`: review gate 集約。
- `final-validation.md`: 最終検証。
- `work-report-input.md`: work_reporter 入力。

## 実装結果

- `pending` は canonical state へ昇格しない。
- Ready job 作成時に未開始 `JOB_PHASE_RUN` を事前作成しない。
- term、persona、body phase は phase 開始時に必要な phase run を作成する。
- `commonPhaseActionAvailability` に集約した操作可否判断を phase service read model で使う。
- `JobIOService` は stale として architecture 正本と backend 図から外した。
- active observability task-local の旧名参照は更新済み。
- `cancelled` fixture spelling は `canceled` に統一済み。
- `stale_selection`、`validation_stale`、`model_selection_stale` は削除対象にしていない。

## Review Gate

| 観点 | ファイル | 状態 | 未解決 | 最大重大度 |
| --- | --- | --- | --- | --- |
| 挙動正しさ | `reviewback.behavior.yaml` | `no_issue` | `false` | `none` |
| 契約互換性 | `reviewback.contract.yaml` | `no_issue` | `false` | `none` |
| 責務境界 | `reviewback.responsibility-boundary.yaml` | `no_issue` | `false` | `none` |
| 状態・データ不変条件 | `reviewback.state-invariant.yaml` | `no_issue` | `false` | `none` |
| 権限・信頼境界 | `reviewback.trust-boundary.yaml` | `no_issue` | `false` | `none` |

## 検証

- `python3 scripts/harness/run.py --suite backend-local`: pass
- `python3 scripts/harness/run.py --suite backend-lint`: pass
- `python3 scripts/harness/run.py --suite structure`: pass
- `python3 scripts/harness/run.py --suite coverage`: pass
- `go test ./internal/apitest ./internal/integrationtest`: pass
- `go test ./internal/usecase ./internal/service`: pass
- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/scenario-design.md --json`: pass
- `git diff --check`: pass

## 作業履歴

- `work_history/runs/2026-05-16-translation-job-state-stale-retirement-run/README.md`
- `work_history/runs/2026-05-16-translation-job-state-stale-retirement-run/codex.md`
- `work_history/runs/2026-05-16-translation-job-state-stale-retirement-run/workflow-improvement-log.jsonl`
- `work_history/runs/2026-05-16-translation-job-state-stale-retirement-run/transcript_refs.json`

## 残留リスク

- 外部 export 済み transcript は未利用である。
- system test 件数は `final-validation.md` では未確認である。
- `.codex/environments/environment.toml` は自動生成 file と判断し、commit 対象から外す。

## merge lane への依頼

- active plan を completed へ移動する。
- local merge を実施する。
- merge 後検証を実施する。
- merge 結果 commit を作成する。
