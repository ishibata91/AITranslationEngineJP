# Orchestration Patterns

## 目的

`implementation-orchestrate` が Copilot 実装 lane を分配し、全 implementation handoff 完了後に final validation と Codex review を行うための判断パターンをまとめる。
agent contract の権限や output obligation は上書きしない。

## 適用ルール

- `implementation-scope` の handoff 見出しを RunSubagent 実行単位にする。
- distiller は tester / implementer より先に必ず起動する。
- tester は implementer より先に必ず起動する。
- subagent に渡す source scope は `single_handoff_packet` 1 件と、その distill 結果に限定する。
- suite-all と Sonar check は全 implementation handoff 完了後だけ実行する。
- Codex review は final validation 後に `codex exec` で呼び出す。
- closeout では final validation と Codex review の evidence または blocked reason を必ず返す。

## 実行順パターン

- 通常: distiller -> tester -> implementer。
- 修正: investigator -> distiller -> tester -> implementer。
- refactor: distiller -> tester -> implementer。
- mixed: backend handoff を先行し、各 handoff を distiller -> tester -> implementer で扱う。
- final validation: 全 implementation handoff 完了後に suite-all -> Sonar check を実行する。
- Codex review: final validation 後に `codex exec` で `review_conductor` を呼び出す。

## Final Validation

- suite-all は `python3 scripts/harness/run.py --suite all` を使う。
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

## 赤旗

- final validation 前に suite-all または Sonar check を実行している。
- Codex review payload に diff または scope path がない。
- repo-local Sonar issue gate と Sonar server Quality Gate を混同している。
- coverage、Sonar、harness、Codex review の未実行理由が completion packet にない。
