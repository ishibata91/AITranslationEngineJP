# 正本化判断

- `task_id`: `2026-05-10-translation-job-state-machine-redesign`
- `status`: `completed`
- `decided_at`: `2026-05-14`

## 判断

追加の docs 正本化は不要とする。
human 承認済みの恒久仕様は、既存の docs 正本差分に反映済みである。

## 反映済み正本

- `docs/spec.md`: `TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` の分離、共通操作規則、terminal guard、同じ `JOB_PHASE_RUN` 継続を反映済み。
- `docs/architecture.md`: `TranslationJobPolicy`、`JobIOService`、通知 module の責務境界を反映済み。
- `docs/er.md`: Ready job では `JOB_PHASE_RUN` を事前作成せず、フェーズ開始時だけ作成する方針を反映済み。
- `docs/detail-specs/term-translation-phase.md`: resume / retry / cancel の共通操作規則、同じ `JOB_PHASE_RUN` 継続、terminal 後書き拒否を反映済み。
- `docs/detail-specs/persona-generation-phase.md`: common operation rule、同じ `JOB_PHASE_RUN` 継続、terminal guard を反映済み。
- `docs/detail-specs/body-translation-phase.md`: retry / resume / cancel、同じ `JOB_PHASE_RUN` 継続、terminal guard を反映済み。
- `docs/detail-specs/translation-job-management.md`: job state と phase run state の参照分離、危険操作無効化を反映済み。

## 追加正本化しない理由

レビュー修正で追加した `UpdateJobPhaseRunWhenState` の expected state 条件は、実装上の状態不変条件である。
恒久仕様としては、共通操作規則と同じ `JOB_PHASE_RUN` 継続で表現済みである。

terminal job の read model 操作可否 false は、terminal guard と危険操作無効化の実装表現である。
恒久仕様としては、terminal job で状態変更操作を拒否する規則に含める。

## 根拠

- `review-aggregation.md`: 5 観点レビュー通過。
- `final-validation.md`: backend-local と coverage 通過。
- `reviewback.behavior.yaml`: terminal job 表示可否と `RecoverableFailed` 操作可否の通過確認。
- `reviewback.state-invariant.yaml`: expected state 条件付き更新の通過確認。

## 詳細仕様正本反映

`details-specs` への追加反映は不要である。
既存の docs 正本差分を今回の正本反映結果として扱う。
