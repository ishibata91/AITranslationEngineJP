---
name: implementing-frontend
description: frontend 側の allowed scope を直接実装し、指定 checks を返す。
---

# Implementing Frontend

## Rules

- 編集前に `docs/coding-guidelines.md` を読む
- 画面構成、導線、情報配置に触れる変更では `docs/screen-design/` と関連する `docs/screen-design/wireframes/` を読む
- active exec-plan と work brief を読んでから編集する
- frontend owned scope だけを変更する
- 指定された checks を実行する
- `sonar-scanner` 実行後に `powershell -File .codex/skills/directing-implementation/scripts/get-open-sonar-issues.ps1 -Project ishibata91_AITranslationEngineJP -OwnedPaths <frontend-owned-paths>` で `status == OPEN` の issue だけを見て、frontend owned scope issue が出た場合はその issue が消えるまで修正を継続する
- plan の書き換えや lane 切り替えはしない

## Reference Use

- 着手前に `../directing-implementation/references/directing-implementation.to.implementing-frontend.json` を参照して入力契約を確認する。
- 画面変更時は `docs/index.md` を入口に、関連する `docs/screen-design/` と `docs/screen-design/wireframes/` を参照する。
- `directing-implementation` へ返す時は `references/implementing-frontend.to.directing-implementation.json` を返却契約として使う。
