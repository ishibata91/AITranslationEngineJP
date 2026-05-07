# 人間設計レビュー入力

- `task_id`: `2026-05-07-provider-settings-job-decoupling-implement`
- `status`: `approved`
- `next_artifact_after_approval`: `実装範囲`

## レビュー対象

- `scenario-design.md`
- `ui-design.md`
- `diagramming-result.md`
- `design-diff-components.svg`
- `design-diff-sequence.svg`

## 承認してほしい判断

- Job 側 DB は endpoint、secret store 参照実値、`credential_ref` 実値を所有しない。
- Job Setup は provider、model、execution mode、batch mode の選択値だけを扱う。
- Ready job は実行開始前に最新 provider settings を再解決する。
- Running phase は開始時の非 secret 要約だけを保存する。
- provider settings revision と更新履歴は Job 側に持たせない。

## 次に進む条件

人間が `scenario-design.md`、`ui-design.md`、設計差分図を承認した場合だけ、`実装範囲` を作成する。
差し戻しがある場合は、差し戻し対象の成果物だけを再作成する。

## 承認記録

- `approved_at`: `2026-05-07`
- `human_response`: `ok`
- `approved_artifacts`: `scenario-design.md`, `ui-design.md`, `diagramming-result.md`, `design-diff-components.svg`, `design-diff-sequence.svg`

## 検証済み

- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/scenario-design.md --coverage docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/scenario-design.candidate-coverage.json --json`
- `plantuml -checkonly docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/design-diff-components.puml`
- `plantuml -checkonly docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/design-diff-sequence.puml`
