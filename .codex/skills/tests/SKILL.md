---
name: tests
description: Scenario 実装と unit test 拡張を mode 分岐で扱い、必要な test / fixture / helper を最小差分で追加する role skill。
---

# Tests

## Goal

- 実装済み scope を証明する test を追加する
- Scenario 系と unit 系を mode で切り替える
- product code の広い変更は行わず、証明不足を埋める

## Modes

- `scenario-implementation`: Scenario または fix 再現条件を E2E / fixture / acceptance checks に落とす
- `unit`: 実装済み責務と主要分岐の unit test を補う

## Common Rules

- 実装コードを広く直さない
- test / fixture / helper 以外の product code は必要最小限だけを触る
- 新しい要件解釈を足さない
- `scenario-implementation` では Scenario artifact または fix 再現条件をそのまま証明対象にする
- `unit` では `implementation_task_ids` にない責務へ広げない
- Wails runtime event を使う非同期完了は completion event を主要観測点とする
- review へ handoff する前に touched test files と残 gap を返す
- 役割を再確定せず、呼び出し元で確定した `test_mode` をそのまま実行する

## Output

- `touched_test_files`
- `implemented_test_scope`
- `validation_results`
- `remaining_gaps`

## Detailed Guides

- `references/mode-guides/scenario-implementation.md`
- `references/mode-guides/unit.md`

## Reference Use

- quick overview は `../orchestrate/references/orchestrate.to.tests.json` を使う
- mode 別 contract は `../orchestrate/references/contracts/orchestrate.to.tests.<mode>.json` を正本とする
- 返却 contract は `references/contracts/tests.to.orchestrate.<mode>.json` を正本とする
