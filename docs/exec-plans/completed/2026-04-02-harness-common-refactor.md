# Impl Plan Template

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `scripts/harness/`
- task_id:
- task_catalog_ref:
- parent_phase:

## Request Summary

- `scripts/harness` の既存挙動を維持したまま、共通処理を `harness_common.py` に寄せて保守性と拡張性を上げる。

## Decision Basis

- `run.py`、`check_structure.py`、`check_design.py`、`check_execution.py` に PASS/FAIL 表示と実行制御の重複がある。
- suite 名、exit code、実行順は現状のまま維持したい。
- `scripts/harness` の範囲で閉じる refactor なので、workflow 契約や product docs の見直しは不要と判断する。

## Owned Scope

- `scripts/harness/run.py`
- `scripts/harness/harness_common.py`
- `scripts/harness/check_structure.py`
- `scripts/harness/check_design.py`
- `scripts/harness/check_execution.py`

## Out Of Scope

- `scripts/harness/README.md` の説明更新
- `.codex/` や `docs/` の source-of-truth 更新
- execution harness の環境依存 fail の仕様変更

## Dependencies / Blockers

- なし

## Parallel Safety Notes

- `scripts/harness/` の内部整理だけを扱う。

## UI

- なし

## Scenario

- なし

## Logic

- suite ごとの判定内容と CLI 契約は変えず、共通 helper を追加して実装重複だけを減らす。

## Implementation Plan

- `harness_common.py` に script 実行、section 表示、PASS/FAIL/SKIP 出力、summary 判定の helper を追加する。
- `run.py` を共通 helper 利用へ切り替える。
- `check_structure.py` を required path check と markdown link check に分割する。
- `check_design.py` と `check_execution.py` の report ロジックを共通 helper 経由へ寄せる。

## Acceptance Checks

- `python3 -m py_compile scripts/harness/run.py scripts/harness/harness_common.py scripts/harness/check_structure.py scripts/harness/check_design.py scripts/harness/check_execution.py`
- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite design`
- `python3 scripts/harness/run.py --suite all`

## Required Evidence

- suite 名と exit code 契約が維持されていること
- structure / design harness が pass すること
- full harness の結果と失敗理由

## 4humans Sync

- なし

## Outcome

- `harness_common.py` に script 実行、section 表示、PASS/FAIL/SKIP 出力、summary 判定の helper を追加した。
- `run.py` を共通 helper ベースに切り替え、suite 実行 loop から `subprocess.run(...)` 直書きを外した。
- `check_structure.py` を required path check と markdown link check に分割した。
- `check_design.py` と `check_execution.py` の report ロジックを共通 helper 経由へ寄せた。
- `python3 -m py_compile ...`、`python3 scripts/harness/run.py --suite structure`、`python3 scripts/harness/run.py --suite design` は pass した。
- `python3 scripts/harness/run.py --suite all` は既存 `npm run gate:execution` 内の `cargo clippy` が `index.crates.io` を解決できず fail した。
