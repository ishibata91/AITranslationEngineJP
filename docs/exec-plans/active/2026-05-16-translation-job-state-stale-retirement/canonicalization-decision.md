# Canonicalization Decision

- `skill`: `implement-lane`
- `status`: `completed`
- `decision`: `docs_canonicalized`
- `source`: `implementation-scope.md`, `review-aggregation.md`, `docs-canonicalization-result.md`

## Decision

docs 正本化は完了済みとする。
`details-spec` 追加反映は不要とする。

## Reason

- `JobIOService` stale 廃止は architecture 正本へ反映済みである。
- `docs/spec.md` と `docs/detail-specs/*.md` は、`Ready` job に `JOB_PHASE_RUN` を事前作成しない仕様を既に持っている。
- completed archive は人間判断 Q-002 に従い変更しない。
- 5 観点レビューは `review-aggregation.md` で通過済みである。

## Canonicalized Files

- `docs/architecture.md`
- `docs/diagrams/backend/backend-architecture.puml`

## Not Changed

- `docs/spec.md`
- `docs/detail-specs/translation-job-management.md`
- `docs/detail-specs/term-translation-phase.md`
- `docs/detail-specs/persona-generation-phase.md`
- `docs/detail-specs/body-translation-phase.md`
- `docs/exec-plans/completed/**`

## Verification

- `python3 scripts/harness/run.py --suite structure`: pass。
- `plantuml --check-syntax docs/diagrams/backend/backend-architecture.puml`: pass。
- `rg -n "JobIOService|internal/jobio|jobio" docs/architecture.md docs/diagrams/backend/backend-architecture.puml internal .go-arch-lint.yml --glob '!**/*_test.go'`: exit code `1`、出力なし。

## Remaining Risk

残留不足はない。
