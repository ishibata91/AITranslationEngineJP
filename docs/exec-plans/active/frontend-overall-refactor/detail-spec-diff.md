# 詳細仕様差分: frontend-overall-refactor

- `skill`: detail-spec-design
- `status`: docs-canonicalization-decision-recorded
- `source_plan`: `./plan.md`
- `detail_spec_target`: `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/translation-job-setup.md`
- `screen_design_diff`: `N/A`
- `component_diagram`: `N/A`

## 詳細仕様差分

この refactor 内では docs 正本本文を変更しない。

`FSD-005` は人間判断で `実装が正` である。
そのため code 修正ではなく、docs 正本化候補として扱う。

## docs 正本化判断

- `decision`: docs 正本化が必要である。
- `reason`: 翻訳管理 shell の現行実装は、入力データ確認画面から翻訳ジョブを直接作成して `job-run` へ進む。既存 docs には `translation-job-setup` 相当の中間 step が残っている。
- `code_change`: 不要。`FSD-005` は code 修正対象外である。
- `canonical_docs_change`: 未実施。docs 正本文言の人間承認が未取得である。
- `handoff_target`: `updating-docs`
- `candidate_docs`: `docs/screen-design/screens/translation-management.md`, `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/translation-job-setup.md`, `docs/screen-design/screens/translation-job-setup.md`

## 未決

- `Q-001`: 仕様乖離整理で `実装が正` と判断された項目が出た場合、docs 正本化を別成果物として扱うか。

## 回答

- `Q-001`: 別成果物として扱う。
  - `回答者`: 人間
  - `回答日`: 2026-05-24
  - `反映先`: `docs正本化判断` と後続 `updating-docs`

## 根拠

- `source`: `docs/exec-plans/active/frontend-overall-refactor/plan.md`
- `review`: `reviewback.behavior.yaml`, `reviewback.contract.yaml`, `reviewback.trust-boundary.yaml`, `reviewback.state-invariant.yaml`, `reviewback.responsibility-boundary.yaml`
- `validation`: `python3 scripts/harness/run.py --suite structure`, `npm --prefix frontend run build-storybook`, `python3 scripts/harness/run.py --suite frontend-local`
