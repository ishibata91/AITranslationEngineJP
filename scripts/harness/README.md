# Harness Scripts

## Entry Points

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite execution`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Suites

- `structure`: 必須ファイル、必須ディレクトリ、Markdown リンクの検査
- `design`: 主要文書の契約語と最低限の内容確認
- `execution`: 将来の Rust / frontend 実装に対する標準 test / lint / build 入口

## Execution Notes

- `execution` suite は `package.json` の `lint` / `test` / `build` と `Cargo.toml` の標準 command に加えて、repo root に `sonar-project.properties` がある時は `lint` の後段で `sonar-scanner` を実行する
- Sonar issue の取得と remediation loop は harness ではなく implementation lane の skill 契約で扱う
