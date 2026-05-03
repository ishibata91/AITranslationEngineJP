# UI Design: <task-id>

- `skill`: ui-design
- `status`: draft
- `source_plan`: `./plan.md`
- `scenario_source`: `./scenario-design.md`
- `ui_prototype`: `./prototype.svelte`
- `prototype_server_url`: `http://127.0.0.1:34116/prototype`
- `prototype_server_command`: `npm --prefix frontend run dev:prototype -- --task <task-id> --port 34116`
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

## UI Prototype Contract

- `prototype_kind`: `existing_screen_change | new_screen`
- `source_basis`:
- `prototype_path`: `./prototype.svelte`
- `required_before_human_review`: `yes`
- `required_for_frontend_handoff`: `yes`
- `framework_conversion`: UIプロトタイプの構造を frontend framework へ変換する
- `prototype_server_url`: `http://127.0.0.1:34116/prototype`
- `prototype_server_command`: `npm --prefix frontend run dev:prototype -- --task <task-id> --port 34116`
- `human_review_server_required`: `yes`
- `mock_data_root`: `./mock-data/`
- `mock_data_migration`: `forbidden`
- `sample_data_root`: `[data-ui-prototype-sample-data-root]`
- `sample_data_migration`: `forbidden`
- `production_reference_direction`: `product_code_must_not_reference_ui_prototype`
- `interaction_review`:
- `state_transition_review`:
- `wording_review`:
- `structure_to_preserve`:
- `allowed_changes_during_conversion`:
- `forbidden_changes_during_conversion`:

## Agent Browser Review

- `command_source`: `agent-browser`
- `served_url`: `http://127.0.0.1:34116/prototype`
- `server_command`: `npm --prefix frontend run dev:prototype -- --task <task-id> --port 34116`
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
  - `display_wording`:
  - `input_effort`:
  - `eye_movement`:
  - `responsive_continuity`:
- `wording_review`:
  - `review_timing`: `after_agent_browser_review`
  - `fixed_names_preserved`:
  - `business_japanese_terms`:
  - `internal_state_names_hidden`:
  - `next_action_wording`:
  - `allowed_english_labels`:
  - `plain_language_next_action_judgement`:
- `console_errors`:
- `screenshot_or_snapshot_refs`:
- `layout_breaks`:
- `ambiguous_interactions`:
- `open_issues`:
- `not_checked_reason`:

## Rules

- UI は `ui-design.md` で固定する
- UIプロトタイプは `docs/exec-plans/active/<task-id>/` 配下を正本にする
- UIプロトタイプを作る場合は確認サーバーの URL を `agent-browser` で開き、UX 観点から確認する
- `agent-browser` 確認後に、専門知識がなくても次に何をするか分かる表現水準かを表示文言レビューで確認する
- 固定名以外の画面表示文言は、日本語の業務語へ置き換える
- 内部状態名は画面に出さず、利用者の次操作を示す文へ変換する
- 英語ラベルは、利用者が設定画面で見る既存語だけに限定する
- 人間確認中は UIプロトタイプ確認サーバーを起動したままにする
- 既存画面変更では、既存画面または既存 UI 部品を土台にする
- 新規画面では、`docs/screen-design` の画面設計に従う
- frontend 実装がある task では、UIプロトタイプを task-local 確認用として扱う
- frontend 実装では UIプロトタイプの構造を framework へ変換し、主要区画、導線、状態表示を維持する
- 本番コードから UIプロトタイプを参照しない
- 実 API、永続化、本番 gateway 接続、業務ロジック完全再現は UIプロトタイプの対象外にする
- `mock-data/` 配下の値は状態表示確認用であり、frontend 実装へ移植しない
- `[data-ui-prototype-sample-data-root]` の範囲はモックデータ置き場であり、frontend 実装へ移植しない
- 細かな visual polish は実装後に人間が実物を確認して直す
- product component 名や owned scope は、implementation-scope で必要な時だけ扱う
- implementation-scope の `owned_scope` や product code 対象 file は書かない
