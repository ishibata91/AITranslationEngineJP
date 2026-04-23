# UI Design: <task-id>

- `skill`: ui-design
- `status`: draft
- `source_plan`: `./plan.md`
- `scenario_source`: `./scenario-design.md`

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

## Rules

- UI は実装前の mock ではなく、実装が満たす要件契約として書く
- 実装前の見た目 artifact を新規必須にしない
- 細かな visual polish は実装後に人間が実物を確認して直す
- product component 名や owned scope は、implementation-scope で必要な時だけ扱う
- implementation-scope の `owned_scope` や product code 対象 file は書かない
