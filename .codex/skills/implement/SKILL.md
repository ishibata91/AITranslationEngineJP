---
name: implement
description: implementation_target と owned scope に従い、frontend / backend / mixed の実装と品質通過を担当する role skill。
---

# Implement

## Goal

- `implementation_target` と `owned_scope` に従って必要なコードと関連 test を変更する
- implementation lane と fix lane の両方を同じ role で処理する
- frontend / backend 専用 skill は作らず、task ごとの scope で分岐する

## Required Inputs

- active work plan
- `task_mode`: `implement` | `fix` | `refactor`
- `implementation_target`: `frontend` | `backend` | `mixed`
- `owned_scope`
- `implement` / `refactor` では承認済み requirements、存在する `ui-mock` / `scenario`、`implementation-brief`、`implementation-scope`
- `implementation-brief` は背景、制約、判断理由の参照元として扱う
- `implementation-scope` は AI handoff 専用資料として扱い、実行単位、依存、validation の参照元とする
- `fix` では accepted scope と trace / reproduce evidence

## Common Rules

- 編集前に `docs/coding-guidelines.md` を読む
- `implementation-scope` を超える broad refactor や仕様追加を行わない
- ownership が曖昧なら停止して orchestrate へ返す
- plan の書き換えや lane 切り替えはしない
- 関連 test の更新は同一変更で行う
- 関連 test を追加または更新する時は、各 test method を Arrange / Act / Assert で判別できる構造にする
- 関連 test を追加または更新する時は、各 test method を 1 つの振る舞いと 1 つの検証対象に絞る。複数の検証対象が必要なら test case を分割する
- 関連 test の assertion bundle は、1 つの output object、1 つの state object、1 つの event payload を確認する場合だけに限る
- 関連 test の setup は決定的にし、test body に条件分岐を持ち込まない
- closeout 前に `python3 scripts/harness/run.py --suite all` を実行する
- backend を含む task は closeout 前に Sonar MCP で open `HIGH` / `BLOCKER`、open reliability、open security を確認し、すべて 0 件にする
- backend を含む task は `sonar_gate_result` に件数、対象 project、補足 issue を含めて返す
- backend を含む task は Sonar 件数ゲートを validation の一部として返す
- 役割を再確定せず、呼び出し元で確定した `implementation_target` と `task_mode` を前提に進める

## Target Notes

- `frontend`: 画面導線、state、UI event、Wails bridge 周辺を主対象にする
- `backend`: usecase、service、repository、adapter、validation を主対象にする
- 実装 task があるものは `implementation_target` に関係なく、scope freeze 済みの `owned_scope` を処理する
- `fix` では narrow permanent fix を優先し、ついでの整理を入れない

## Output

- `touched_files`
- `implemented_scope`
- `validation_results`
- `sonar_gate_result`
- `closeout_notes`
- `residual_risks`

## Detailed Guides

- `references/mode-guides/frontend.md`
- `references/mode-guides/backend.md`
- `references/mode-guides/mixed.md`
- `references/mode-guides/fix-lane.md`

## Reference Use

- quick overview は `../orchestrate/references/orchestrate.to.implement.json` を使う
- mode 別 contract は `../orchestrate/references/contracts/orchestrate.to.implement.<mode>.json` を正本とする
- 返却 contract は `references/contracts/implement.to.orchestrate.<mode>.json` を正本とする
