---
name: updating-docs
description: Codex 側の docs 正本化作業プロトコル。implementation completion後に、human 承認済み docs-only artifact を正本へ反映する判断基準を提供する。
---
# Updating Docs

## 目的

`updating-docs` は作業プロトコルである。
`docs_updater` agent が implementation completion後に human 承認済み artifact を docs 正本へ反映するための、source of truth、承認確認、validation の見方を提供する。

人間可読な実行境界、handoff、stop / reroute はこの skill を正本にする。

## 対応ロール

- `docs_updater` が使う。
- 返却先は caller または次 agent とする。
- owner artifact は `updating-docs` の出力規約で固定する。

## 入力規約

- Codex implementation lane の修正完了が分かっている時
- human 承認済み docs-only artifact を docs 正本へ移す時
- canonicalization target と validation を整理する時
- task-local artifact と docs source of truth の対応を確認する時
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。
- 必須入力: caller, implementation_completion_report, approval_record, approved_artifact, canonicalization_targets
- 任意入力: validation_commands, source_docs
- 必須 artifact: Codex implementation completion report, approved docs-only artifact, /Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md

## 外部参照規約

- agent runtime と tool policy は [docs_updater.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/docs_updater.toml) の `allowed_write_paths` / `allowed_commands` とする。
- binding: [docs_updater.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/docs_updater.toml)
- agent runtime: [docs_updater.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/docs_updater.toml)
- forbidden: product code、product test、workflow / skill / agent runtime の変更
- tool policy: agent runtime の `allowed_write_paths` / `allowed_commands` に従う
- docs index: [index.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md)
- binding: [docs_updater.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/docs_updater.toml)
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/updating-docs/SKILL.md

## 内部参照規約

### 拘束観点

- Codex implementation completion report の確認
- docs source of truth の選び方
- human approval record の確認
- approved artifact と canonical target の対応
- validation と remaining gaps の記録

## 判断規約

- implementation completion後にだけ正本化へ進む
- human 承認済み artifact だけを反映する
- docs-only scope を超えない
- implementation-scope を docs 正本へ自動昇格しない
- 未確定仕様を独断で補完しない

- Codex implementation completion report を根拠として残す
- approval record を根拠として残す
- source of truth と task-local artifact を分ける
- validation 結果を残す

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力に tool policy、agent runtime、product code の変更義務を含めない。

### Handoff

- handoff 先: `implement_lane`
- 渡す scope: docs 更新結果、validation、remaining gaps
- 必須出力: touched_docs_files, updated_source_of_truth, validation_results, remaining_gaps

## 完了規約

- 出力規約を満たし、次の actor が再解釈なしで判断できる。
- 不足情報または停止理由がある場合は明示されている。
- Codex implementation completion report を確認した。
- human approval record を確認した。
- approved artifact と canonical target を対応づけた。
- validation 結果と remaining gaps を記録した。
- 必須 evidence: Codex implementation completion report, approval record, source artifact path, validation result
- completion signal: implementation completion後の docs 正本が approved artifact と同期している
- residual risk key: remaining_gaps

## 停止規約

- Codex implementation lane の修正完了が未確認の時
- workflow / skill / agent runtime や skill / agent を変更する時
- product code や product test の変更が必要な時
- human approval が不足している時
- Codex implementation lane の修正完了が分からない場合は停止する。
- approval がない場合は停止する。
- workflow 変更なら `implement_lane` へ戻す。
- product 実装が必要なら `implement_lane` へ戻す。
- implementation completion前に正本化しない
- 未承認 draft を正本化しない
- workflow 変更を docs 更新に混ぜない
- product implementation を同時に進めない
- 停止時は不足項目、衝突箇所、reroute 先を返す。
- implementation completion前に正本化しなかった場合は停止する。
- 未承認 draft を正本化しなかった場合は停止する。
- implementation-scope を自動昇格しなかった場合は停止する。
- workflow 変更や product 実装を混ぜなかった場合は停止する。
- 拒否条件: Codex implementation completion missing
- 拒否条件: approval missing
- 拒否条件: product implementation required
- 拒否条件: workflow / skill / agent runtime change
- 停止条件: Codex implementation lane の修正完了が分からない
- 停止条件: docs-only scope ではない
- 停止条件: human approval がない
