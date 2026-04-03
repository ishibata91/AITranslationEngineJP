# Harness Scripts

## Entry Points

- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite design`
- `python3 scripts/harness/run.py --suite execution`
- `python3 scripts/harness/run.py --suite all`

## Suites

- `structure`: 必須ファイル、必須ディレクトリ、Markdown リンクの検査
- `design`: 主要文書の契約語と最低限の内容確認
- `execution`: 将来の Rust / frontend 実装に対する標準 format / test / lint / build 入口

## Execution Notes

- repo root の `package.json` に `gate:execution` がある時、`execution` suite はその script を正規入口として実行する
- repo root の `gate:execution` は `format:check`、`lint`、`src-tauri/` の Rust gate、`scan:sonar`、`test`、`src-tauri/` の Rust test、`build` をこの順でまとめる
- `gate:execution` が未定義の repo では、`execution` suite は `package.json` の `format:check` / `lint` / `test` / `build` と `Cargo.toml` の標準 command を個別実行する
- repo root に `sonar-project.properties` がある時、Sonar step の正本は repo root の `scan:sonar` script とし、未定義の場合だけ `sonar-scanner` を直接実行する
- Sonar issue の取得と remediation loop は harness ではなく implementation lane の skill 契約で扱う
