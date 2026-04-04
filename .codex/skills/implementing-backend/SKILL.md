---
name: implementing-backend
description: backend 側の allowed scope を直接実装し、指定 checks を返す。
---

# Implementing Backend

## Rules

- 編集前に `docs/coding-guidelines.md` を読む
- active exec-plan と work brief を読んでから編集する
- backend owned scope だけを変更する
- implementation lane owner (`directing-implementation`) から渡された `python3 scripts/harness/run.py --suite backend-lint` だけを local validation として実行する
- plan の書き換えや lane 切り替えはしない

## Reference Use

- 着手前に `../directing-implementation/references/directing-implementation.to.implementing-backend.json` を参照して入力契約を確認する。
- `directing-implementation` へ返す時は `references/implementing-backend.to.directing-implementation.json` を返却契約として使う。
