---
name: scenario-external-integration-generation
description: Codex 側の external-integration シナリオ 候補生成 skill。外部 provider、secret、adapter、fake、network 境界 から シナリオ 候補を作る。
---
# 外部連携シナリオ候補生成

## 目的

`scenario-external-integration-generation` は作業プロトコルである。
`scenario_external_integration_generator` が external-integration 観点 の シナリオ 候補だけを作る時に使う。

共通規約と出力形は [scenario-candidate-generation](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-candidate-generation/SKILL.md) に従う。

## 対応ロール

- `scenario_external_integration_generator` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `designer` とする。
- 担当成果物は `scenario-external-integration-generation` の出力規約で固定する。

## 入力規約

- task 枠: 新規実装または機能拡張の対象範囲。
- 根拠要件: 候補生成の根拠にする要件または task 内成果物。
- 対象観点: `external-integration` 観点。
- 承認状態: 呼び出し元が渡す承認済みまたは未承認の状態。
- 不足時の扱い: 根拠参照、担当者、承認状態が不足する場合は推測で補わない。
- 呼び出し元: `implement_lane` を受け取る。
- 引き継ぎ入力: task 枠、根拠要件、作業計画フォルダ、承認状態を受け取る。
- 候補参照: 既存の候補参照パスがある場合だけ受け取る。
- 既知不足: 呼び出し元が把握している不足項目を受け取る。

## 外部参照規約

- エージェント実行定義とツール権限は [scenario_external_integration_generator.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/scenario_external_integration_generator.toml) の 書き込み許可 / 実行許可 とする。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。
- 共通規約: [scenario-candidate-generation](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-candidate-generation/SKILL.md) に従う。
- 統合先規約: [scenario-design](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-design/SKILL.md) を参照する。

## 内部参照規約

### 観点

- 外部 provider、secret store、adapter、ファイル境界 を起点にする
- fake / stub で検証できる境界を明示する
- secret 本体の露出、参照不能、provider 不一致を拾う
- lifecycle や 失敗 と失敗扱いが矛盾する場合は 競合候補 にする
- 有料の real API が必要に見える場合は 人間判断候補 にする

## 判断規約

- 判断は入力 成果物、外部参照規約、対象 agent の責務境界に従う。

## 非対象規約

- 有料の real API 前提や provider 実装方針の固定は扱わない。
- 最終シナリオ表の確定、候補の採否、統合判断は扱わない。
- プロダクト実装、未承認 docs 正本化、ツール権限、プロダクト仕様正本は扱わない。

## 出力規約

- 観点: `external-integration` 観点であることを返す。
- 成果物: `docs/exec-plans/active/<task-id>/scenario-candidates.external-integration.md` を返す。
- 候補: 外部境界、開始条件、期待結果、`fake_or_stub`、観測点を持つ候補を返す。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコード変更の指示を含めない。
- 候補数: 生成した 候補 シナリオ 数を返す。0 件なら不足理由を返す。
- 根拠網羅: 候補 ごとの 根拠要件、関連する詳細要求タイプ、観測点を返す。
- 競合候補: 他 観点 や最終 シナリオ 統合時に競合しうる前提、状態、結果、検証段階を返す。
- 人間判断 候補: AI が確定できない業務判断、状態遷移、外部連携、監査保存対象を返す。

## 完了規約

- 指定 観点 の 候補成果物 が出力規約の必須項目を満たしている。
- 採否や統合判断を行わず、designer が判断できる候補として返却されている。
- 必須 根拠: 根拠要件パス、task 内成果物パス、候補成果物パス、観点を返している。
- 完了判断材料: implement_lane が designer 起動入力に入れる 候補成果物パス、候補数、競合候補、人間判断候補 を判断できる。
- 残留リスク: AI が確定できない判断候補が返っている。

## 停止規約

- 停止時は不足項目、衝突箇所、戻し先を返す。
- 最終シナリオ表 の確定が求められている場合は停止する。
- 候補 採否または統合判断が求められている場合は停止する。
- プロダクト実装が求められている場合は停止する。
- 未承認 docs 正本化が求められている場合は停止する。
- 引き継ぎ入力 だけでは 根拠要件 を特定できない場合は停止する。
- 作業計画フォルダが不足している場合は停止する。
- 候補成果物 の書き先が 作業計画フォルダ 外である場合は停止する。
- 人間レビュー が必要な判断を AI だけで確定しそうな場合は停止する。
