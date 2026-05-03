# UI Design: <task-id>

- `skill`: ui-design
- `status`: draft
- `source_plan`: `./plan.md`
- `scenario_source`: `./scenario-design.md`
- `ui_prototype`: `./prototype.svelte`
- `prototype_server_url`: `http://127.0.0.1:34116/prototype`
- `prototype_server_command`: `npm --prefix frontend run dev:prototype -- --task <task-id> --port 34116`
- `human_review_server_required`: `yes`
- `human_review_designer_agent_required`: `yes | no`
- `human_feedback_route`: `designer_agent_direct | implement_lane`
- `designer_agent_close_after_review`:

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

## Screen Structure UX Check

- `screen_purpose`:
- `target_user`:
- `primary_action`:
- `secondary_actions`:
- `information_hierarchy`:
- `screen_responsibility`:
- `current_state_visibility`:
- `state_based_actions`:
- `display_conditions`:
- `disabled_conditions`:
- `permission_differences`:
- `input_constraints`:
- `input_grouping`:
- `empty_state`:
- `loading_state`:
- `error_state`:
- `change_diff_visibility`:
- `dangerous_action_separation`:
- `screen_transition`:
- `completion_route`:
- `ui_wording`:
- `reviewability`:
- `implementation_reuse_scope`:
- `evidence_storage`:

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
- `existing_screen_resource_refs`:
- `reused_screen_structure`:
- `changed_sections_only`:
- `new_visual_system_added`: `yes | no`
- `new_visual_system_reason`:
- `prototype_path`: `./prototype.svelte`
- `required_before_human_review`: `yes`
- `required_for_frontend_handoff`: `yes`
- `framework_conversion`: UIプロトタイプの構造を frontend framework へ変換する
- `prototype_server_url`: `http://127.0.0.1:34116/prototype`
- `prototype_server_command`: `npm --prefix frontend run dev:prototype -- --task <task-id> --port 34116`
- `human_review_server_required`: `yes`
- `human_review_designer_agent_required`: `yes | no`
- `human_feedback_route`: `designer_agent_direct | implement_lane`
- `designer_agent_close_after_review`:
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
- 画面構造UXチェック表の確認結果を `Screen Structure UX Check` に記録する
- `agent-browser` 確認後に、専門知識がなくても次に何をするか分かる表現水準かを表示文言レビューで確認する
- 固定名以外の画面表示文言は、日本語の業務語へ置き換える
- 内部状態名は画面に出さず、利用者の次操作を示す文へ変換する
- 英語ラベルは、利用者が設定画面で見る既存語だけに限定する
- 人間確認中は UIプロトタイプ確認サーバーを起動したままにする
- 人間確認中に UIプロトタイプ確認サーバーを `designer` agent が保持している場合は、人間レビュー終了まで `designer` agent を起動したままにする
- 人間確認中の UI 指摘は、起動中の `designer` agent が直接受け取り、`ui-design.md` と task-local UIプロトタイプへ反映する
- 既存画面変更では、既存画面または既存 UI 部品を土台にする
- 既存画面変更の UIプロトタイプは、対象画面の既存 Svelte、CSS、class、画面構造を再利用する
- 既存画面変更では、独自の page shell、card、grid、配色、余白体系を新規に作らない
- 既存画面変更では、変更対象区画だけを差し替え、変更しない区画は既存画面の構造と表示を維持する
- 既存画面リソースを再利用できない場合は、理由を `UI Prototype Contract` に記録し、完了扱いにしない
- `new_visual_system_added` が `yes` の場合は、既存画面変更として停止または差し戻しにする
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
