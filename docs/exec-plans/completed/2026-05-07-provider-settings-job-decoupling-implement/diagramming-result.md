# Diagramming Result: 2026-05-07-provider-settings-job-decoupling-implement

- `skill`: `diagramming`
- `status`: `completed`
- `artifact_type`: `設計差分図`
- `source_plan`: `./plan.md`

## 判断結果

設計差分図を作成した。
図の対象は、Job から provider settings / credential / endpoint 所有を外す予定差分だけである。

## source path

- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/design-diff-components.puml`
- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/design-diff-sequence.puml`

## 描画結果 path

- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/design-diff-components.svg`
- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/design-diff-sequence.svg`

## 根拠参照

- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/plan.md`
- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/task-frame.md`
- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/scenario-design.md`
- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/ui-design.md`
- `docs/er.md`
- `docs/diagrams/er/combined-data-model-er.puml`
- `docs/architecture.md`

## 図説明

### コンポーネント図

- Job Setup の選択値境界と provider settings 共通設定境界を分けた。
- `JOB_PHASE_RUN` と `TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT` から外す credential / endpoint 所有を削除予定として明示した。
- Job 側に残してよい provider / model / execution mode / batch mode と状態分類の境界を残した。

### シーケンス図

- Ready job 実行開始時に最新 provider settings を再解決する流れを追加した。
- Running phase は開始時の非 secret 要約だけを保存する流れを追加した。
- provider settings 未設定または参照不能では Running を開始せず、短い分類要約だけ残す流れを示した。

## 検証結果

- `plantuml -checkonly docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/design-diff-components.puml` は成功した。
- `plantuml -checkonly docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/design-diff-sequence.puml` は成功した。
- `plantuml -tsvg docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/design-diff-components.puml` で SVG を生成した。
- `plantuml -tsvg docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/design-diff-sequence.puml` で SVG を生成した。

## 未決事項

- 追加の人間判断が必要な未決事項は、指定入力の範囲では見つからなかった。
- 監査表示に出せる状態分類の最終文言は、図では固定せず分類境界だけを示した。
