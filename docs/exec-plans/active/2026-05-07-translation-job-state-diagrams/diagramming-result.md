# diagramming result

## 判断結果

- 判断結果: 図作成完了。
- 図成果物種別: `明示された補助図`。
- 図説明: 翻訳ジョブ状態を表す保存境界、状態投影境界、状態を変えうる操作入口だけを抽出した。
- 不足情報: なし。

## 根拠参照

- `docs/spec.md`: 翻訳ジョブの状態遷移と状態名。
- `docs/er.md`: `TRANSLATION_JOB`、`JOB_PHASE_RUN`、状態集約方針。
- `docs/detail-specs/translation-job-management.md`: 未完了一覧、操作可否、削除、停止、再開入口の仕様。
- `internal/repository/job_lifecycle_repository.go`: job、phase run、runtime snapshot の保存型。
- `internal/service/translation_job_management_service.go`: 状態投影、操作可否、削除・停止・再開入口。
- `frontend/src/application/gateway-contract/translation-job-management/translation-job-management-gateway-contract.ts`: UI 側の状態契約。

## 出力

- source path: `docs/exec-plans/active/2026-05-07-translation-job-state-diagrams/translation-job-state-er.puml`
- source path: `docs/exec-plans/active/2026-05-07-translation-job-state-diagrams/translation-job-state-components.puml`
- source path: `docs/exec-plans/active/2026-05-07-translation-job-state-diagrams/translation-job-state-sequence.puml`
- 描画結果 path: `docs/exec-plans/active/2026-05-07-translation-job-state-diagrams/translation-job-state-er.svg`
- 描画結果 path: `docs/exec-plans/active/2026-05-07-translation-job-state-diagrams/translation-job-state-components.svg`
- 描画結果 path: `docs/exec-plans/active/2026-05-07-translation-job-state-diagrams/translation-job-state-sequence.svg`

## 検証結果

- `plantuml -tsvg docs/exec-plans/active/2026-05-07-translation-job-state-diagrams/translation-job-state-er.puml docs/exec-plans/active/2026-05-07-translation-job-state-diagrams/translation-job-state-components.puml docs/exec-plans/active/2026-05-07-translation-job-state-diagrams/translation-job-state-sequence.puml`
- 結果: 通過。
