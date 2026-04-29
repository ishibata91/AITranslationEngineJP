---
name: scenario-external-integration-generation
description: Codex 側の external-integration シナリオ 候補生成 skill。外部 provider、secret、adapter、fake、network 境界 から シナリオ 候補を作る。
---
# Scenario External Integration Generation

## 目的

`scenario-external-integration-generation` は knowledge package である。
`scenario_external_integration_generator` が external-integration 観点 の シナリオ 候補だけを作る時に使う。

共通規約と出力形は [scenario-candidate-generation](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-candidate-generation/SKILL.md) に従う。

## 対応ロール

- `scenario_external_integration_generator` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `designer` とする。
- 担当成果物は `scenario-external-integration-generation` の出力規約で固定する。

## 入力規約

- 入力は `task 枠`、根拠要件、対象 観点 を含む。
- 入力に 根拠参照、担当者、承認状態が不足する場合は推測で補わない。
- 必須入力: 呼び出し元, 引き継ぎ入力, active_task_folder, レーン担当
- 任意入力: candidate_source_paths, 既知不足
- 必須 成果物: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-external-integration-generation/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/templates/task-folder/scenario-candidates.viewpoint.md, 進行中 task folder または 呼び出し元提供 task 文脈

## 外部参照規約

- エージェント実行定義とツール権限は [scenario_external_integration_generator.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/scenario_external_integration_generator.toml) の 書き込み許可 / 実行許可 とする。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-external-integration-generation/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-candidate-generation/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-design/SKILL.md

## 内部参照規約

### 観点

- 外部 provider、secret store、adapter、file 境界 を起点にする
- fake / stub で検証できる境界を明示する
- secret 本体の露出、参照不能、provider mismatch を拾う
- lifecycle や 失敗 と失敗扱いが矛盾する場合は 競合候補 にする
- real paid API が必要に見える場合は 人間判断候補 にする

## 判断規約

- 判断は入力 成果物、外部参照規約、対象 agent の責務境界に従う。
- 対象外の成果物、ツール権限、プロダクト 仕様正本はこの skill で変更しない。

## 出力規約

- `観点`: `external-integration`
- `成果物`: `scenario-candidates.external-integration.md`
- `候補`: external 境界、trigger、expected 結果、fake_or_stub、observable point を必ず持つ
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。
- 観点名: external-integration 観点であることを返す。
- 候補 成果物: `docs/exec-plans/active/<task-id>/scenario-candidates.external-integration.md` を返す。
- 候補数: 生成した 候補 シナリオ 数を返す。0 件なら不足理由を返す。
- source coverage: 候補 ごとの 根拠要件、関連 detail requirement type、観測点を返す。
- 競合 hint: 他 観点 や最終 シナリオ 統合時に競合しうる前提、状態、結果、検証段階を返す。
- 人間判断 候補: AI が確定できない業務判断、状態遷移、外部連携、監査保存対象を返す。

## 完了規約

- 指定 観点 の 候補成果物 が出力規約の必須項目を満たしている。
- 採否や統合判断を行わず、designer が判断できる候補として返却されている。
- 必須 根拠: 根拠要件 path or task 内成果物 path, 候補成果物 path, 観点
- 完了判断材料: implement_lane が designer 入力一式 に入れる 候補成果物 path、候補数、競合候補、人間判断候補 を判断できる。
- 残留リスク: AI が確定できない判断候補が返っている。

## 停止規約

- real paid API を前提にしない
- provider 実装方針を シナリオ として固定しない
- 採用、不採用、統合を確定しない
- 停止時は不足項目、衝突箇所、戻し先を返す。
- 拒否条件: final シナリオ表 の確定が求められている
- 拒否条件: 候補 採否または統合判断が求められている
- 拒否条件: プロダクト実装 request
- 拒否条件: unapproved docs 正本化
- 停止条件: 引き継ぎ入力 だけでは 根拠要件 を特定できない
- 停止条件: 進行中 task folder が不足している
- 停止条件: 候補成果物 の書き先が 進行中 task folder 外である
- 停止条件: 人間レビュー が必要な判断を AI だけで確定しそうである
