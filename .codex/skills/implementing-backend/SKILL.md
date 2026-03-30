---
name: implementing-backend
description: backend 側の allowed scope を直接実装し、指定 checks を返す。
---

# Implementing Backend

## Rules

- 編集前に `docs/coding-guidelines.md` を読む
- active exec-plan と work brief を読んでから編集する
- backend owned scope だけを変更する
- 指定された checks を実行する
- `sonar-scanner` 実行後に `powershell -File .codex/skills/directing-implementation/scripts/get-open-sonar-issues.ps1 -Project ishibata91_AITranslationEngineJP -OwnedPaths <backend-owned-paths>` で `status == OPEN` の issue だけを見て、backend owned scope issue が出た場合はその issue が消えるまで修正を継続する
- plan の書き換えや lane 切り替えはしない

## Reference Use

- 着手前に `../directing-implementation/references/directing-implementation.to.implementing-backend.json` を参照して入力契約を確認する。
- `directing-implementation` へ返す時は `references/implementing-backend.to.directing-implementation.json` を返却契約として使う。
