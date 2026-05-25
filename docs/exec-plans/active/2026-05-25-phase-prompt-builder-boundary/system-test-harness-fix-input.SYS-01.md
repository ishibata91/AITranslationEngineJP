# system-test harness fix 入力: SYS-01

- `task_id`: `2026-05-25-phase-prompt-builder-boundary`
- `handoff_id`: `SYS-01`
- `implementation_artifact`: system-test harness 修正
- `implementation_skill`: tests-scenario
- `source_scope`: `./implementation-scope.md`
- `human_review`: human requested harness pass and clarified `fakeAPI` is stale on 2026-05-25

## 目的

廃止済み `fakeAPI` query に依存していた translation-job-management の system-test を、実 Wails backend と deterministic な SQLite seed に置き換える。
`python3 scripts/harness/run.py --suite system-test` が未完了 job 0 件で失敗しない状態にする。

## 失敗内容

- command: `python3 scripts/harness/run.py --suite system-test`
- sandbox 内 result: Wails dev server が `http://127.0.0.1:34115` で ready にならず fail。
- sandbox 外 result: 10 件中 7 件 pass、translation-job-management 3 件 fail。
- fail detail: `/?fakeApi=1&fakeScenario=success#translation-management` で未完了 job card が 0 件になった。

## 変更許可範囲

- `tests/system/translation-job-management.spec.ts`
- `scripts/test/run-system-test.sh`
- system-test 専用 seed helper under `scripts/test/`

## 禁止範囲

- 廃止済み `fakeAPI`、`fakeScenario`、frontend fake API fixture の復活。
- production UI 文言、layout、style の変更。
- production 仕様の変更。
- `.codex/`。
- docs 正本本文。
- 他 agent の変更の取り消し。

## 完了条件

1. translation-job-management system-test は `/#translation-management` を開く。
2. Wails 実 backend が isolated system-test DB を読む。
3. seed は repository API 経由で running、ready、completed job を作る。
4. completed job は未完了一覧から除外される確認用に使う。
5. sandbox 外 `python3 scripts/harness/run.py --suite system-test` が 10 件 pass する。

## 期待出力

- `system-test-harness-fix-result.SYS-01.md`
- 変更ファイル一覧
- seed 方針
- 検証結果
- 残留リスク
