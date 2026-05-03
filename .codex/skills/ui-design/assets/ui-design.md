# UI Design: <task-id>

- `skill`: ui-design
- `status`: draft
- `source_plan`: `./plan.md`
- `scenario_source`: `./scenario-design.md`
- `html_mock`: `./ui-mock.html`
- `mock_server_url`: `http://127.0.0.1:34116/ui-mock.html`
- `mock_server_command`: `npm run dev:ui-mock -- --task <task-id> --port 34116`
- `human_review_server_required`: `yes`

## UI Contract

- `display_items`:
- `primary_actions`:
- `button_enablement`:
- `state_variants`:
- `post_implementation_review`:

## Interface Frame

- `purpose`:
- `audience`:
- `primary_workflow`:
- `information_density`:
- `visual_direction`:
- `remembered_signal`:

## Structure Notes

- `page_sections`:
- `layout_constraints`:
- `responsive_constraints`:
- `accessibility_constraints`:

## Interaction States

- `loading`:
- `empty`:
- `error`:
- `disabled`:
- `progress`:
- `retry`:
- `success`:

## Post Implementation Review

- `desktop_review_points`:
- `mobile_review_points`:
- `overflow_risks`:
- `visual_polish_open_questions`:

## HTML Mock Contract

- `mock_path`: `./ui-mock.html`
- `required_before_human_review`: `yes`
- `required_for_frontend_handoff`: `yes`
- `framework_conversion`: HTML の基本構造を frontend framework へ変換する
- `mock_server_url`: `http://127.0.0.1:34116/ui-mock.html`
- `mock_server_command`: `npm run dev:ui-mock -- --task <task-id> --port 34116`
- `human_review_server_required`: `yes`
- `mock_data_root`: `./mock-data/`
- `mock_data_migration`: `forbidden`
- `sample_data_root`: `[data-ui-mock-sample-data-root]`
- `sample_data_migration`: `forbidden`
- `structure_to_preserve`:
- `allowed_changes_during_conversion`:
- `forbidden_changes_during_conversion`:

## Agent Browser Review

- `command_source`: `agent-browser`
- `served_url`: `http://127.0.0.1:34116/ui-mock.html`
- `server_command`: `npm run dev:ui-mock -- --task <task-id> --port 34116`
- `server_status_during_human_review`:
- `mock_data_refs`:
- `used_only_for_display_state_review`:
- `migration_to_product_code`: `forbidden`
- `migration_to_fixture_or_test_data`: `forbidden`
- `checked_viewports`:
- `ux_review_points`:
  - `goal_completion`:
  - `information_priority`:
  - `operation_order`:
  - `state_comprehension`:
  - `recovery_path`:
  - `input_effort`:
  - `eye_movement`:
  - `responsive_continuity`:
- `console_errors`:
- `screenshot_or_snapshot_refs`:
- `layout_breaks`:
- `ambiguous_interactions`:
- `open_issues`:
- `not_checked_reason`:

## Rules

- UI は `ui-design.md` で固定する
- HTML モックを作る場合は確認サーバーの URL を `agent-browser` で開き、UX 観点から確認する
- 人間確認中は HTML モック確認サーバーを起動したままにする
- frontend 実装がある task では、HTML モックを task-local 確認用として扱う
- frontend 実装では HTML の基本構造を framework へ変換し、主要区画、導線、状態表示を維持する
- `mock-data/` 配下の値は状態表示確認用であり、frontend 実装へ移植しない
- `[data-ui-mock-sample-data-root]` の範囲はモックデータ置き場であり、frontend 実装へ移植しない
- 細かな visual polish は実装後に人間が実物を確認して直す
- product component 名や owned scope は、implementation-scope で必要な時だけ扱う
- implementation-scope の `owned_scope` や product code 対象 file は書かない
