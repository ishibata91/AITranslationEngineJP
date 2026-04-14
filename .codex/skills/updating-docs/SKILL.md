---
name: updating-docs
description: "`docs/` 配下の正本更新だけを扱う。human 承認済みの docs-only task を、human 直起動または orchestrate handoff で処理する。"
---

# Updating Docs

## Output

- `touched_docs_files`
- `updated_source_of_truth`
- `validation_results`
- `remaining_gaps`

## Rules

- `docs/` 配下の正本更新だけを扱う
- human 承認済みの `docs-only` task でだけ使う
- human が直接起動した時も、`orchestrate` から handoff された時も同じ contract で扱う
- 変更前に対象の正本文書と関連する tests / acceptance checks / validation commands を読む
- `.codex/` と product code は更新しない
- 未確定の仕様を agent が独断で補完しない
- workflow 契約変更は `skill-modification` へ戻す

## Reroute

- 更新対象が `docs/` ではなく `.codex/` である
- human の明示判断なしに恒久仕様を追加または変更しようとしている
- product code の更新を同時に要求している

## Detailed Guides

- `references/mode-guides/docs-only.md`

## Reference Use

- quick overview は `references/updating-docs.to.orchestrate.json` を使う
- mode 別 contract は `references/contracts/updating-docs.to.orchestrate.docs-only.json` を正本とする
