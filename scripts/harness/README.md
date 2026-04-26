# Harness Scripts

## Entry Points

- `python3 scripts/harness/run.py --suite frontend-lint`
- `python3 scripts/harness/run.py --suite backend-lint`
- `python3 scripts/harness/run.py --suite frontend-local`
- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite frontend-test`
- `python3 scripts/harness/run.py --suite backend-test`
- `python3 scripts/harness/run.py --suite system-test`
- `python3 scripts/harness/run.py --suite scenario-gate`
- `python3 scripts/harness/run.py --suite coverage`
- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite execution`
- `python3 scripts/harness/run.py --suite all`

## Suites

- `frontend-lint`: repo root package の `lint:frontend` を入口にして frontend lint だけを実行する
- `backend-lint`: repo root package または backend workspace の `lint:backend` を入口にして backend lint を実行する
- `frontend-local`: `frontend-lint`、`frontend-test` をこの順で実行する実装中の軽量 gate
- `backend-local`: `backend-lint`、`backend-test` をこの順で実行する実装中の軽量 gate
- `frontend-test`: repo root package の `test:frontend` を入口にして frontend test を実行する
- `backend-test`: repo root package の `test:backend` を入口にして backend test を実行する
- `system-test`: repo root package の `test:system` を入口にして Playwright system test を実行する
- `scenario-gate`: active task の `scenario-design.md` にある candidate coverage と詳細要求 coverage を検査し、漏れ report と人間質問票を生成する
- `coverage`: repo root package の `test:frontend:coverage` と `test:backend:coverage` を入口にして Sonar 互換の project coverage を 70% 基準で検査し、Sonar 用 report path と集計値を `test-results/coverage-manifest.json` にまとめる
- `structure`: `docs/index.md` を repo の地図として扱い、リンク切れを検査する
- `execution`: `lint:backend`、`lint:frontend`、`test:backend`、`test:frontend`、Sonar をまとめて確認する入口

## Execution Notes

- implementation lane の implementer / tester は backend handoff で `backend-local`、frontend handoff で `frontend-local` を local validation に使う
- mixed handoff は touched layer に応じて `backend-local` と `frontend-local` の両方を実行する
- local validation は実装中の feedback 用であり、`execution`、`coverage`、`system-test`、`all`、Sonar check は final validation lane に残す
- implementer / tester は local validation の結果、または未実行理由を completion output に残す
- `execution` suite は repo root の `package.json` を唯一の入口として扱い、`lint:backend`、`lint:frontend`、`test:backend`、`test:frontend`、Sonar step をこの順で実行する
- `all` suite は `structure`、`scenario-gate`、`execution`、`system-test`、`coverage` をこの順で実行する
- `coverage` suite は単独でも実行できる独立 gate として維持しつつ、`all` からも実行する
- repo root に `sonar-project.properties` がある時、Sonar step の正本は repo root の `scan:sonar` script とし、未定義の場合だけ `sonar-scanner` を直接実行する
- Sonar issue の取得と remediation loop は harness ではなく implementation lane の skill 契約で扱う
