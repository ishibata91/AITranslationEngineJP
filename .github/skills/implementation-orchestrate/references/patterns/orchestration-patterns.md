# Orchestration Patterns

## 目的

`implementation-orchestrate` が Copilot 実装 lane を分配し、全 implementation handoff 完了後に final validation と Codex review を行うための判断パターンをまとめる。
agent contract の権限や output obligation は上書きしない。

## 適用ルール

- `implementation-scope` の handoff 見出しを RunSubagent 実行単位にする。
- Ready Waves 表または `ready_wave` から、実行可能な最小番号の wave を選ぶ。
- `first_action` を含む `single_handoff_packet` だけを subagent に渡す。
- distiller は tester / implementer より先に必ず起動する。
- tester を implementer より先に起動できるのは、承認済み scenario artifact を product test 化する場合だけである。
- unit test と原因未確定の regression test は、implementer 完了後に tester が追加または更新する。
- subagent に渡す source scope は `single_handoff_packet` 1 件と、その distill 結果に限定する。
- scenario validation、suite-all、Sonar check は全 implementation handoff 完了後だけ実行する。
- scenario validation が fail した場合は close せず、Copilot 側 blocker として返す。
- Codex review は final validation 後に `codex exec` で呼び出す。
- closeout では final validation と Codex review の evidence または blocked reason を必ず返す。
- Codex review の戻り値は `copilot_action` で受け取り、再解釈しない。

## 実行順パターン

- 通常: distiller -> implementer -> tester。
- 修正: investigator -> distiller -> implementer -> tester。
- refactor: distiller -> implementer -> tester。
- scenario 先行: distiller -> tester -> implementer。
- mixed: backend handoff を先行し、各 handoff を通常順または scenario 先行順で扱う。
- final validation: 全 implementation handoff 完了後に scenario validation -> suite-all -> Sonar check を実行する。
- Codex review: final validation 後に `codex exec` で `review_conductor` を呼び出す。

## Final Validation

- scenario validation の default は `python3 scripts/harness/run.py --suite scenario-gate` を使う。
- task 固有の product scenario test command が `implementation-scope` にある場合は、scenario validation result に含める。
- suite-all は `python3 scripts/harness/run.py --suite all` を使う。
- suite-all は既に scenario-gate を含むが、completion packet では scenario validation result を別 field として抜き出す。
- Sonar check は repo root の `scan:sonar` script を優先する。
- repo-local gate は Sonar サーバ側 Quality Gate ではない。
- `npm run test:system` または harness all が Wails、sandbox、OS 権限で止まる場合は `FAIL_ENVIRONMENT` として扱う。
- `FAIL_ENVIRONMENT` は blocked reason、再実行環境、再実行コマンドを residual risk に残す。

## Codex Review Payload

`codex exec` で Codex review conductor を呼び出す時は、次を渡す。

- `implementation_scope_path`
- `approval_record`
- `implementation_result`
- `diff_summary` または review 対象 diff
- `final_validation_result`
- `touched_files`

## Codex Review Result

`codex_review_result.copilot_action` は次の分岐だけに使う。

- `close`: completion packet に review result を残して終了する。
- `report_residual`: priority override または confidence residual を residual risk に残して終了する。
- `fix`: `copilot_patch_scope` 内だけを修正し、final validation と Codex review を再実行する。
- `rerun_validation`: 指定された不足 validation だけを再実行し、Codex review を再実行する。
- `rerun_codex_review`: payload を補い、product code を変更せず Codex review だけを再実行する。

## 赤旗

- final validation 前に scenario validation、suite-all、Sonar check を実行している。
- scenario validation failure を close 可能な residual risk として扱っている。
- `first_action` 欠落を広い調査で補っている。
- `parallelizable_with` がない handoff を同一 wave という理由で並列実行している。
- Codex review payload に diff または scope path がない。
- `rerun_codex_review` で product code を変更している。
- `fix` で `copilot_patch_scope` 外を変更している。
- repo-local Sonar issue gate と Sonar server Quality Gate を混同している。
- coverage、Sonar、harness、Codex review の未実行理由が completion packet にない。
