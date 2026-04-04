# Implementer Lint Suites

- workflow: impl
- status: completed
- lane_owner: directing-implementation
- scope: `scripts/harness/`, `package.json`, `.codex/skills/`, `.codex/README.md`, `.codex/workflow.md`, `scripts/harness/README.md`
- task_id:
- task_catalog_ref:
- parent_phase:

## Request Summary

- implementation lane の `implementer` が full harness や Sonar を回さず、frontend / backend lint だけを安定した suite 名で実行できるようにする。

## Decision Basis

- 現在の repo には `frontend-lint` / `backend-lint` の harness suite がなく、`implementer` が具体 command 群へ依存すると lint 内容の増減のたびに skill 契約の更新が必要になる。
- `directing-implementation` が final gate として `all` と Sonar を一元実行し、`implementer` には軽い validation だけを任せる方が lane 全体の速度と保守性に合う。

## Owned Scope

- `scripts/harness/run.py`
- `scripts/harness/check_frontend_lint.py`
- `scripts/harness/check_backend_lint.py`
- `scripts/harness/README.md`
- `package.json`
- `.codex/skills/directing-implementation/`
- `.codex/skills/implementing-frontend/`
- `.codex/skills/implementing-backend/`
- `.codex/agents/implementer.toml`
- `.codex/README.md`
- `.codex/workflow.md`

## Out Of Scope

- product code や bugfix
- `docs/` 正本の恒久仕様更新
- fix lane の validation 再設計

## Dependencies / Blockers

- `design` harness は着手前から `.codex/skills/directing-fixes/SKILL.md` の既存 `tasks.md` pattern 欠落で fail している。
- `.codex/agents/implementer.toml` と `.codex/agents/ctx_loader.toml` に既存 dirty change があるため、差分を上書きしない。

## Parallel Safety Notes

- `package.json` と `scripts/harness/run.py` は shared file なので、この変更で一貫して更新する。
- `.codex/agents/implementer.toml` は既存 dirty change を保持したまま必要最小限で追記する。

## UI

- なし。

## Scenario

- implementation lane owner (`directing-implementation`) は `frontend-lint` または `backend-lint` suite を `implementer` に渡す。
- `implementer` は assigned lint suite だけを実行して返却する。
- implementation lane owner (`directing-implementation`) は close 前にだけ `python3 scripts/harness/run.py --suite all` と Sonar gate を実行する。

## Logic

- `frontend-lint` は repo root frontend lint 群のみを担当する。
- `backend-lint` は `src-tauri/` の `cargo fmt --all --check` と `cargo clippy --all-targets --all-features -- -D warnings` のみを担当する。
- `lint` / `gate:execution` は新しい lint script 名に寄せて、lint 実装詳細を skill 契約から切り離す。

## Implementation Plan

- harness に `frontend-lint` と `backend-lint` suite を追加する。
- `package.json` に `lint:frontend` と `lint:backend` を追加し、既存 aggregate script を整理する。
- implementation lane の skill / handoff / workflow 文書を lint suite 前提に更新する。
- validation を実行し、plan を completed へ移す。

## Acceptance Checks

- `python3 scripts/harness/run.py --suite frontend-lint`
- `python3 scripts/harness/run.py --suite backend-lint`
- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite design`
- `python3 scripts/harness/run.py --suite all`

## Required Evidence

- 新 lint suite の pass/fail 出力
- structure / design / all の結果
- implementation lane 文書で `implementer` が lint suite だけを実行する証跡

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`
- 今回の変更では diagram 更新は不要想定

## Outcome

- `scripts/harness/run.py --suite frontend-lint` と `--suite backend-lint` を追加し、repo 側で lint 責務を stable suite 名に固定した。
- `package.json` に `lint:frontend` と `lint:backend` を追加し、`lint` と `gate:execution` を新 split に寄せた。
- `directing-implementation` と `implementing-*` の skill / handoff / workflow 文書を、implementer が lint suite のみ実行し、direction が full harness と Sonar を担当する契約へ更新した。
- `design` harness の既存 blocker だった `.codex/skills/directing-fixes/SKILL.md` の `tasks.md` 文言不足も workflow 整合として修正した。
- Validation passed:
  - `python3 scripts/harness/run.py --suite frontend-lint`
  - `python3 scripts/harness/run.py --suite backend-lint`
  - `python3 scripts/harness/run.py --suite structure`
  - `python3 scripts/harness/run.py --suite design`
  - `python3 scripts/harness/run.py --suite all`
- `4humans` は品質状態の変更を伴わない workflow / harness 更新のため、追加更新なし。
