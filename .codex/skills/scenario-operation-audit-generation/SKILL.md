---
name: scenario-operation-audit-generation
description: Codex 側の operation-audit scenario 候補生成 skill。運用確認、監査、ログ、履歴、再現性から scenario 候補を作る。
---
# Scenario Operation Audit Generation

## 目的

`scenario-operation-audit-generation` は knowledge package である。
`scenario_operation_audit_generator` が operation-audit viewpoint の scenario 候補だけを作る時に使う。

共通規約と出力形は [scenario-candidate-generation](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-candidate-generation/SKILL.md) に従う。

## 対応ロール

- `scenario_operation_audit_generator` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `designer` とする。
- owner artifact は `scenario-operation-audit-generation` の出力規約で固定する。

## 入力規約

- 入力は `task_frame`、source requirement、対象 viewpoint を含む。
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。
- 必須入力: caller, handoff_packet, active_task_folder, lane_owner
- 任意入力: candidate_source_paths, known_gaps
- 必須 artifact: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-operation-audit-generation/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/templates/task-folder/scenario-candidates.viewpoint.md, active task folder または caller-provided task context

## 外部参照規約

- エージェント実行定義とツール権限は [scenario_operation_audit_generator.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/scenario_operation_audit_generator.toml) の `allowed_write_paths` / `allowed_commands` とする。
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-operation-audit-generation/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-candidate-generation/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-design/SKILL.md

## 内部参照規約

### 観点

- 運用者が後から確認する情報を起点にする
- audit log、履歴、再現材料、error summary を分ける
- 保存すべきものと保存してはいけないものを分ける
- security / data requirement と衝突する保存対象は conflict hint にする
- 監査粒度が不明な場合は human decision candidate にする

## 判断規約

- 判断は入力 artifact、外部参照規約、対象 agent の責務境界に従う。
- 対象外の成果物、ツール権限、product 仕様正本はこの skill で変更しない。

## 出力規約

- `viewpoint`: `operation-audit`
- `artifact`: `scenario-candidates.operation-audit.md`
- `candidate`: audit event、stored summary、redaction rule、observable point を必ず持つ
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。
- 観点名: operation-audit 観点であることを返す。
- 候補 artifact: `docs/exec-plans/active/<task-id>/scenario-candidates.operation-audit.md` を返す。
- 候補数: 生成した candidate scenario 数を返す。0 件なら不足理由を返す。
- source coverage: candidate ごとの source requirement、関連 detail requirement type、観測点を返す。
- 競合 hint: 他 viewpoint や最終 scenario 統合時に競合しうる前提、状態、outcome、検証段階を返す。
- human decision 候補: AI が確定できない業務判断、状態遷移、外部連携、監査保存対象を返す。

## 完了規約

- 指定 viewpoint の candidate artifact が出力規約の必須項目を満たしている。
- 採否や統合判断を行わず、designer が判断できる候補として返却されている。
- 必須 evidence: source requirement path or task-local artifact path, candidate artifact path, viewpoint
- 完了判断材料: implement_lane が designer packet に入れる candidate artifact path、候補数、conflict hint、human decision candidate を判断できる。
- 残留リスク: AI が確定できない判断候補が返っている。

## 停止規約

- observability を実装ログ形式として固定しない
- secret や個人情報の保存を推測で許可しない
- 採用、不採用、統合を確定しない
- 停止時は不足項目、衝突箇所、戻し先を返す。
- 拒否条件: final scenario matrix の確定が求められている
- 拒否条件: candidate 採否または統合判断が求められている
- 拒否条件: product implementation request
- 拒否条件: unapproved docs canonicalization
- 停止条件: handoff_packet だけでは source requirement を特定できない
- 停止条件: active task folder が不足している
- 停止条件: candidate artifact の書き先が active task folder 外である
- 停止条件: human review が必要な判断を AI だけで確定しそうである
