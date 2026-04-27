# Orchestration Patterns

## 目的

`implementation-orchestrate` が Copilot 実装 lane を分配し、全 implementation handoff 完了後に final validation と人間実行用 Codex review request を作るための判断パターンをまとめる。
agent contract の権限や output obligation は上書きしない。

## 適用ルール

- `implementation-scope` の handoff 見出しを RunSubagent 実行単位にする。
- Ready Waves 表または `ready_wave` から、実行可能な最小番号の wave を選ぶ。
- `first_action` を含む `single_handoff_packet` だけを subagent に渡す。
- distiller は tester / implementer より先に必ず起動する。
- tester を implementer より先に起動できるのは、承認済み `APIテスト` を product test 化する場合だけである。
- unit test と原因未確定の regression test は、implementer 完了後に tester が追加または更新する。
- subagent に渡す source scope は `single_handoff_packet` 1 件と、その distill 結果に限定する。
- scenario validation、suite-all、Sonar check は全 implementation handoff 完了後だけ実行する。
- scenario validation が fail した場合は close せず、Copilot 側 blocker として返す。
- Codex review は Copilot が直接呼び出さず、final validation 後に人間実行用 request payload と `codex exec` command を返す。
- closeout では final validation と Codex review request の evidence または blocked reason を必ず返す。
- 人間から Codex review の戻り値が戻された場合だけ、`copilot_action` で受け取り、再解釈しない。

## 実行順パターン

- 通常: distiller -> implementer -> tester。
- 修正: investigator -> distiller -> implementer -> tester。
- refactor: distiller -> implementer -> tester。
- APIテスト先行: distiller -> tester -> implementer。
- mixed: backend handoff を先行し、各 handoff を通常順または APIテスト先行順で扱う。
- final validation: 全 implementation handoff 完了後に scenario validation -> suite-all -> Sonar check を実行する。
- Codex review request: final validation 後に人間実行用 payload と `codex exec` command を返す。

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

## Codex Review Request

人間が `codex exec` で Codex review conductor を呼び出せるよう、completion packet に次を渡す。

- `implementation_scope_path`
- `approval_record`
- `implementation_result`
- `diff_summary` または review 対象 diff
- `final_validation_result`
- `touched_files`
- `human_codex_exec_command`

## Codex Review Result

人間から戻された `codex_review_result.copilot_action` は次の分岐だけに使う。

- `close`: completion packet に review result を残して終了する。
- `report_residual`: priority override または confidence residual を residual risk に残して終了する。
- `fix`: reviewer result bundle、aggregation trace、remediation handoff を読み、chosen strategy、chosen scope、why_not_narrower、why_not_wider、planned changes、invariant tests、used review signals を決めてから修正し、final validation と Codex review request を再作成する。
- `rerun_validation`: 指定された不足 validation だけを再実行し、Codex review request を再作成する。
- `rerun_codex_review`: payload を補い、product code を変更せず Codex review request だけを再作成する。

## 赤旗

- final validation 前に scenario validation、suite-all、Sonar check を実行している。
- scenario validation failure を close 可能な residual risk として扱っている。
- `first_action` 欠落を広い調査で補っている。
- `parallelizable_with` がない handoff を同一 wave という理由で並列実行している。
- Codex review request payload に diff または scope path がない。
- `rerun_codex_review` で product code を変更している。
- `fix` で symptoms だけを潰す局所修正に閉じている。
- `fix` で why_not_narrower と why_not_wider なしに scope を選んでいる。
- repo-local Sonar issue gate と Sonar server Quality Gate を混同している。
- coverage、Sonar、harness、Codex review request の未実行理由が completion packet にない。
