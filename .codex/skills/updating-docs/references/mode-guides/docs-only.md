# Updating Docs: docs-only

## Goal

- human 承認済みの `docs/` 正本だけを更新する

## Rules

- human 直起動または `orchestrate` handoff の時だけ使う
- `approval_record` がない場合は停止する
- `docs/` の正本更新だけを扱い、`.codex/` と product code は触らない
- 更新前に対象の正本文書、関連 tests、acceptance checks、validation commands を読む
- workflow 契約変更は `skill-modification` へ戻す

## Return

- touched_docs_files
- updated_source_of_truth
- validation_results
- remaining_gaps
