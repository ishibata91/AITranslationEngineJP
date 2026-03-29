---
name: updating-docs
description: `docs/` 配下の正本更新だけを扱う。human が直接起動した時だけ使い、恒久仕様や設計判断を repo 正本へ反映する。
---

# Updating Docs

## Output

- touched_docs_files
- updated_source_of_truth
- validation_results
- remaining_gaps

## Rules

- `docs/` 配下の正本更新だけを扱う
- human が直接起動した時だけ使う
- 変更前に対象の正本文書と関連する tests / acceptance checks / validation commands を読む
- `.codex/`、`4humans/`、product code は更新しない
- 未確定の仕様を agent が独断で補完しない
- workflow 契約変更は `skill-modification` へ戻す

## Reroute

- 更新対象が `docs/` ではなく `.codex/` である
- human の明示判断なしに恒久仕様を追加または変更しようとしている
- product code や `4humans/` の更新を同時に要求している
