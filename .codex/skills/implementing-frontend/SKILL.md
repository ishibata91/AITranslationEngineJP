---
name: implementing-frontend
description: frontend 側の allowed scope を直接実装し、指定 checks を返す。
---

# Implementing Frontend

## Rules

- active exec-plan と work brief を読んでから編集する
- frontend owned scope だけを変更する
- 指定された checks を実行する
- `sonar-scanner` 実行後に SonarQube MCP で frontend owned scope の open issue が出た場合は、その issue が消えるまで修正を継続する
- plan の書き換えや lane 切り替えはしない

## Reference Use

- 着手前に `../directing-implementation/references/directing-implementation.to.implementing-frontend.json` を参照して入力契約を確認する。
- `directing-implementation` へ返す時は `references/implementing-frontend.to.directing-implementation.json` を返却契約として使う。
