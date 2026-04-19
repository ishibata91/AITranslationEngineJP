# UI Design: <task-id>

- `skill`: ui-design
- `status`: draft
- `source_plan`: `./plan.md`
- `requirements_source`: `./requirements-design.md`

## HTML Mock Artifact

- `working_mock_path`: `./<task-id>.ui.html`
- `preview_url`:
- `desktop_screenshot_artifacts`:
- `mobile_screenshot_artifacts`:
- `final_mock_path`: `docs/mocks/<page-id>/index.html` または `N/A`

## Interface Frame

- `purpose`:
- `audience`:
- `primary_workflow`:
- `information_density`:
- `visual_direction`:
- `remembered_signal`:

## Structure Notes

- `page_sections`:
- `component_like_sections`:
- `state_variants`:
- `layout_constraints`:
- `visual_tokens`:

## Interaction States

- `loading`:
- `empty`:
- `error`:
- `disabled`:
- `progress`:
- `retry`:
- `success`:

## Review Evidence

- `desktop_screenshot`:
- `mobile_screenshot`:
- `interaction_notes`:
- `known_gaps`:

## Canonicalization

- `final_mock_path`: `docs/mocks/<page-id>/index.html` または `N/A`
- `canonicalization_targets`:

## Rules

- task-local HTML mock を UI artifact の主入力にする
- working mock は `docs/exec-plans/active/<task-id>/<task-id>.ui.html` に置く
- HTML / CSS / 必要最小限の素の JavaScript だけで主要導線と状態変化を再現する
- framework 記法や product component 名を mock の正本に持ち込まない
- implementation-scope の `owned_scope` や product code 対象 file は書かない
