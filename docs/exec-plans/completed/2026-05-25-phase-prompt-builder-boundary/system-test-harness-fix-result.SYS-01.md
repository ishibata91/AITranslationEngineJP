# system-test harness fix 結果: SYS-01

- `task_id`: `2026-05-25-phase-prompt-builder-boundary`
- `handoff_id`: `SYS-01`
- `implementation_skill`: `tests-scenario`
- `status`: completed

## 実装結果

- translation-job-management system-test の URL から廃止済み `fakeApi` / `fakeScenario` query を削除した。
- `scripts/test/run-system-test.sh` は `test-results/system-test/translation-job-management.sqlite3` を実 backend DB path として export する。
- system-test 起動前に seed helper が isolated SQLite DB を reset してから seed する。
- 既存 Wails dev server がある場合は、既定で停止して stale DB の混入を避ける。

## Seed 方針

- `running` job: `body_translation` の running phase run と runtime snapshot を持つ。
- `ready` job: phase run を持たず、runtime snapshot を持つ。
- `completed` job: 未完了一覧から除外される確認用として作る。
- 各 job は xEdit 入力と translation record を持つ。
- seed は frontend fake API fixture ではなく repository API だけを使う。

## 変更ファイル

- `tests/system/translation-job-management.spec.ts`
- `scripts/test/run-system-test.sh`
- `scripts/test/seed-system-test-db/main.go`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/system-test-harness-fix-input.SYS-01.md`
- `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/system-test-harness-fix-result.SYS-01.md`

## 検証結果

- `GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./scripts/test/seed-system-test-db`: pass
- sandbox 内 `python3 scripts/harness/run.py --suite system-test`: fail
  - reason: Wails dev server が `http://127.0.0.1:34115` で ready にならなかった。
- sandbox 外 `python3 scripts/harness/run.py --suite system-test`: pass
  - result: 10 tests passed
- sandbox 外 `python3 scripts/harness/run.py --suite all`: pass
  - structure: pass
  - execution: backend lint pass、frontend lint pass、backend test pass、frontend test pass
  - system-test: 10 tests passed
  - coverage: Sonar coverage `71.1%`, line `73.0%`, branch `57.1%`
  - Sonar issues: security `0`, reliability `0`, maintainability HIGH `0`

## 残留リスク

- sandbox 内 Wails dev server 起動は ready timeout のままである。
- sandbox 外 harness は通過しているため、残留リスクは sandbox 実行環境に限定される。
- `WAILS_SYSTEM_TEST_ALLOW_EXISTING_SERVER=1` を使う場合、既存 server が isolated DB を読んでいる保証はない。
