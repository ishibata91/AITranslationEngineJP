# Harness Scripts

## Entry Points

- `python3 scripts/harness/run.py --suite frontend-lint`
- `python3 scripts/harness/run.py --suite backend-lint`
- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite execution`
- `python3 scripts/harness/run.py --suite all`

## Suites

- `frontend-lint`: repo root package の `lint:frontend` を入口にして frontend lint だけを実行する
- `backend-lint`: repo root package または backend workspace の `lint:backend` を入口にして backend lint を実行する
- `structure`: `docs/index.md` を repo の地図として扱い、リンク切れを検査する
- `execution`: `Wails + Go + frontend` 実装に対する標準 format / test / lint / build 入口

## Execution Notes

- implementation lane の implementer は `frontend-lint` または `backend-lint` を local validation に使い、direction は review が `pass` になった後で `all` を final harness として実行する
- repo root の `package.json` に `gate:execution` がある時、`execution` suite はその script を正規入口として実行する
- repo root の `gate:execution` は `format:check`、`lint`、`scan:sonar`、`test`、backend test、`build` をこの順でまとめる
- `gate:execution` が未定義の repo では、`execution` suite は `package.json` の `format:check` / `lint` / `test` / `build` と `go test ./...` を個別実行する
- repo root に `sonar-project.properties` がある時、Sonar step の正本は repo root の `scan:sonar` script とし、未定義の場合だけ `sonar-scanner` を直接実行する
- Sonar issue の取得と remediation loop は harness ではなく implementation lane の skill 契約で扱う
