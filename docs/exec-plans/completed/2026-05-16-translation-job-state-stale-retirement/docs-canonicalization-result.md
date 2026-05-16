# Docs Canonicalization Result

- `skill`: `updating-docs`
- `agent`: `docs_updater`
- `status`: `completed`
- `source_plan`: `docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement`
- `approval_record`: `human-design-review-request.md`
- `implementation_result`: `implementation-wave-result.md`
- `review_result`: `review-aggregation.md`

## 承認確認

- `human-design-review-request.md` は `approved` である。
- 人間判断 Q-001 は、`JobIOService` を stale として architecture 正本から外す選択である。
- 人間判断 Q-002 は、completed archive を変更しない選択である。
- 人間判断 Q-003 は、`cancelled` fixture spelling を `canceled` へそろえる選択である。
- `implementation-wave-result.md` は backend と tests の完了を記録している。
- `review-aggregation.md` は 5 観点レビュー通過と `implementation_action: close` を記録している。

## 正本化結果

- `docs/architecture.md`: `JobIOService` を構造主語と依存方向から外した。
- `docs/architecture.md`: job / phase run 状態の取得と保存を、既存 UseCase、Service、Repository 境界で扱う責務へ置換した。
- `docs/architecture.md`: `PolicyResult`、rule 名、policy 判定履歴を read model の永続値へ出さない制約を追加した。
- `docs/diagrams/backend/backend-architecture.puml`: backend 構造図から `JobIOService` ノードと依存線を外した。

## 変更しなかった正本

- `docs/spec.md`: `Ready` job に `JOB_PHASE_RUN` を事前作成しない記述が既にあるため変更しなかった。
- `docs/detail-specs/translation-job-management.md`: `Ready job` に `JOB_PHASE_RUN` を事前作成しない記述が既にあるため変更しなかった。
- `docs/detail-specs/term-translation-phase.md`: `Ready job` に `JOB_PHASE_RUN` を事前作成しない記述が既にあるため変更しなかった。
- `docs/detail-specs/persona-generation-phase.md`: phase 開始時だけ `JOB_PHASE_RUN` を作成する記述が既にあるため変更しなかった。
- `docs/detail-specs/body-translation-phase.md`: phase 開始時だけ `JOB_PHASE_RUN` を作成する記述が既にあるため変更しなかった。
- `docs/exec-plans/completed/**`: 人間判断 Q-002 に従い変更しなかった。

## 検索確認

- `rg -n "JobIOService|internal/jobio|jobio" docs/architecture.md docs/diagrams/backend/backend-architecture.puml docs/spec.md docs/detail-specs/translation-job-management.md docs/detail-specs/term-translation-phase.md docs/detail-specs/persona-generation-phase.md docs/detail-specs/body-translation-phase.md`: exit code `1`、出力なし。
- `rg -n "PolicyResult|rule 名|policy 判定履歴|read model|Ready job|JOB_PHASE_RUN を事前作成しない|canceled|cancelled" docs/architecture.md docs/spec.md docs/detail-specs/translation-job-management.md docs/detail-specs/term-translation-phase.md docs/detail-specs/persona-generation-phase.md docs/detail-specs/body-translation-phase.md`: exit code `0`、期待する既存仕様と追加制約だけを確認した。

## 検証

- `python3 scripts/harness/run.py --suite structure`: pass。
- `plantuml --check-syntax docs/diagrams/backend/backend-architecture.puml`: pass。
- `git diff --check -- docs/architecture.md docs/diagrams/backend/backend-architecture.puml docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/docs-canonicalization-result.md`: pass。

## 残留不足

- 現時点で docs 正本化の残留不足はない。
