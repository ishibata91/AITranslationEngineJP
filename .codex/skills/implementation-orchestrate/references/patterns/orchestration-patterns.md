# Orchestration Patterns

## 目的

`implementation_orchestrator` が Codex implementation lane 実装 lane を分配し、全 implementation handoff 完了後に final validation と integrated review を行うための判断パターンをまとめる。
agent contract の権限や output obligation は上書きしない。

## 適用ルール

- `implementation-scope` の handoff 見出しを spawn_agent 実行単位にする。
- Ready Waves 表または `ready_wave` から、実行可能な最小番号の wave を選ぶ。
- `first_action` を含む `single_handoff_packet` だけを subagent に渡す。
- distiller は implementation_tester / implementation_implementer より先に必ず起動する。
- implementation_tester を implementation_implementer より先に起動できるのは、承認済み `APIテスト` を product test 化する場合だけである。
- unit test と原因未確定の regression test は、implementation_implementer 完了後に implementation_tester が追加または更新する。
- subagent に渡す source scope は `single_handoff_packet` 1 件と、その distill 結果に限定する。
- scenario validation、suite-all、Sonar check は全 implementation handoff 完了後だけ実行する。
- scenario validation が fail した場合は close せず、Codex implementation lane blocker として返す。
- final validation 後に review input が揃う場合だけ、4 観点 review agent を context 継承なしで並列 spawn する。
- review input 不足時は reviewer を起動せず、`rerun_validation` または `rerun_codex_review` を返す。
- closeout では final validation と integrated review の evidence または blocked reason を必ず返す。
- integrated review の `implementation_action` は再解釈せず、そのまま分岐する。

## 実行順パターン

- 通常: distiller -> implementation_implementer -> implementation_tester。
- 修正: implementation_investigator -> distiller -> implementation_implementer -> implementation_tester。
- refactor: distiller -> implementation_implementer -> implementation_tester。
- APIテスト先行: distiller -> implementation_tester -> implementation_implementer。
- mixed: backend handoff を先行し、各 handoff を通常順または APIテスト先行順で扱う。
- final validation: 全 implementation handoff 完了後に scenario validation -> suite-all -> Sonar check を実行する。
- integrated review: final validation 後に 4 観点 review agent を並列 spawn し、lossless aggregation と `implementation_action` を返す。

## Final Validation

- scenario validation の default は `python3 scripts/harness/run.py --suite scenario-gate` を使う。
- scenario validation result は `APIテスト群 + UI人間操作E2E 群` の集約結果として扱う。
- task 固有の product scenario test command が `implementation-scope` にある場合は、scenario validation result に含める。
- suite-all は `python3 scripts/harness/run.py --suite all` を使う。
- suite-all は既に scenario-gate を含むが、completion packet では scenario validation result を別 field として抜き出す。
- Sonar check は repo root の `scan:sonar` script を優先する。
- repo-local gate は Sonar サーバ側 Quality Gate ではない。
- `npm run test:system` または harness all が Wails、sandbox、OS 権限で止まる場合は `FAIL_ENVIRONMENT` として扱う。
- `FAIL_ENVIRONMENT` は blocked reason、再実行環境、再実行コマンドを residual risk に残す。

## Integrated Review

観点別 review agent には同じ input を渡す。

- `implementation_scope_path`
- `approval_record`
- `implementation_result`
- `diff_summary` または review 対象 diff
- `final_validation_result`
- `touched_files`

diff、scope、implementation result が不足する場合は `implementation_action: rerun_codex_review` を返す。
final validation result が不足する場合は `implementation_action: rerun_validation` を返す。
review 可能な場合だけ `review_behavior`、`review_contract`、`review_trust_boundary`、`review_state_invariant` を context 継承なしで並列 spawn する。

aggregation は次を守る。

- all group score `> 0.85` の場合だけ `strict_pass` にする。
- trust boundary failure は hard gate とし、他観点の pass で相殺しない。
- conflict は `trust_boundary > behavior > contract > state_invariant` で裁定する。
- 観点別 raw result は `reviewer_result_bundle` に保持する。
- 派生結果は `aggregation_trace` から raw result へ辿れるようにする。

## Integrated Review Result

`codex_review_result.implementation_action` は次の分岐だけに使う。

- `close`: completion packet に review result を残して終了する。
- `report_residual`: priority override または confidence residual を residual risk に残して終了する。
- `fix`: reviewer result bundle、aggregation trace、remediation handoff を読み、chosen strategy、chosen scope、why_not_narrower、why_not_wider、planned changes、invariant tests、used review signals を決めてから修正し、final validation と integrated review を再実行する。
- `rerun_validation`: 指定された不足 validation だけを再実行し、review input を再構築する。
- `rerun_codex_review`: payload を補い、product code を変更せず review input だけを再構築する。

## 赤旗

- final validation 前に scenario validation、suite-all、Sonar check を実行している。
- scenario validation failure を close 可能な residual risk として扱っている。
- `first_action` 欠落を広い調査で補っている。
- `parallelizable_with` がない handoff を同一 wave という理由で並列実行している。
- integrated review input に diff または scope path がない。
- `rerun_codex_review` で product code を変更している。
- `fix` で symptoms だけを潰す局所修正に閉じている。
- `fix` で why_not_narrower と why_not_wider なしに scope を選んでいる。
- repo-local Sonar issue gate と Sonar server Quality Gate を混同している。
- coverage、Sonar、harness、integrated review の未実行理由が completion packet にない。
