# Task Plan: 2026-06-14-prompt-persona-customization

- `workflow`: work
- `status`: planned（未着手。T3 完了後に着手）
- `task_id`: 2026-06-14-prompt-persona-customization
- `task_mode`: 重 task 見込み（画面が動く）
- `request_summary`: 属性 → 翻訳指示ルールの編集、プロンプトテンプレートの編集、実際に投げるプロンプトの参照をできるようにする。
- `goal`: ユーザーがルールとプロンプトを翻訳前に編集し、実プロンプトを目視確認できる。
- `constraints`: 翻訳エンジン本体の手続き（T1〜T3）は変えない。編集と参照の面を足す。
- `close_conditions`: ルール編集 UI・プロンプトテンプレート編集・実プロンプト表示が動く。
- `source_branch`: `master`
- `target_branch`: `master`

## Scope（含む / 含まない）

含む:
- 属性 → 翻訳指示のルール集の編集 UI（翻訳前に設定・編集。`system_requirements.md` §3）。ルール集を永続化する。
- プロンプトテンプレートの編集。
- 実際に翻訳 AI へ投げるプロンプトの参照（目視確認）。
- 画面が動くため storybook-module 経由。

含まない:
- AI によるペルソナ生成・会話履歴解析（将来検討）。

## 依存

- T2（ペルソナ機構）と T3（マスター辞書）。

## Routing Notes

- `required_reading`:
  - `docs/system_requirements.md`（§3 ルール集の永続化と翻訳前編集）
  - `docs/concept-model.md`（素材の性質と口調の合成）
  - `docs/UX-standard.md`（UI 設計の正本）
  - `docs/references/storybook.md`（画面実装の作法）

## Outcome

- 未着手。
