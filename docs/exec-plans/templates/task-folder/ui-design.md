# UI Design: <task-id>

- `skill`: ui-design
- `status`: draft
- `source_plan`: `./plan.md`
- `requirements_source`: `./requirements-design.md`

## Figma Artifact

- `figma_file_url`:
- `figma_file_key`:
- `figma_page`:
- `primary_node_id`:
- `review_frame_ids`:
- `screenshot_artifacts`:

## Interface Frame

- `purpose`:
- `audience`:
- `primary_workflow`:
- `information_density`:
- `visual_direction`:
- `remembered_signal`:

## Structure Notes

- `frames`:
- `components`:
- `variants`:
- `layout_constraints`:
- `tokens_or_variables`:

## Interaction States

- `loading`:
- `empty`:
- `error`:
- `disabled`:
- `progress`:
- `retry`:
- `success`:

## Review Evidence

- `figma_screenshot`:
- `design_context_summary`:
- `variable_defs_summary`:
- `known_gaps`:

## Canonicalization

- `final_mock_path`: `docs/mocks/<page-id>/figma.md`
- `canonicalization_targets`:

## Rules

- Figma file / node を UI artifact の主入力にする
- repo には Figma URL、node id、判断、状態差分、証跡だけを残す
- HTML / CSS mock を作らない
- Figma authoring tool がない session では、human-provided Figma file / node を要求して停止する
- implementation-scope の `owned_scope` や product code 対象 file は書かない
