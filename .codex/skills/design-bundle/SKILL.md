---
name: design-bundle
description: Codex 側の 設計成果物 進行 skill。必須要件、UI、シナリオ、implementation-scope を task 内成果物 として固定するための 正本、進め方、引き継ぎ を提供する。
---
# Design Bundle

## 目的

`design-bundle` は作業プロトコルである。
`designer` agent と top-level Codex が、必須要件、UI、シナリオ、implementation-scope を task 内成果物 として固定する時の、人間可読な実行説明の正本として使う。

作業流れ の次 行動判断、task folder 進行管理、人間向け Codex implementation レーン 引き継ぎ の返却は `implement_lane` が担当する。
プロダクトコードとプロダクトテスト は変更しない。

## 対応ロール

- `designer` が使う。
- 呼び出し元は `implement_lane` または人間とする。
- 返却先は 人間レビュー または `implement_lane` とする。
- 担当成果物は `design-bundle` の出力規約で固定する。

## 入力規約

- 入力は 呼び出し元 から渡された task 内成果物、根拠参照、必要な承認状態を含む。
- 入力に 根拠参照、担当者、承認状態が不足する場合は推測で補わない。
- 必須入力: 呼び出し元, 引き継ぎ入力, user_instruction_or_task_summary, active_task_folder, 設計範囲, レーン担当
- 任意入力: 人間レビュー記録, 対象 skill, 既存設計成果物, シナリオ候補成果物, 既知不足
- 入力規約: 引き継ぎ入力 だけで作業できること。引き継いでいない会話文脈に依存しない。
- 必須 成果物: /Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md, 進行中 task folder

## 外部参照規約

- エージェント実行定義とツール権限は [designer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.toml) の 書き込み許可 / 実行許可 とする。
- secondary: 入力一式 に明示された関連 docs、関連 skill、human の現在指示
- エージェント実行定義: [designer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.toml)
- ツール権限: エージェント実行定義の 書き込み許可 / 実行許可 に従う
- ui: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/ui-design/SKILL.md)
- 候補 common: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-candidate-generation/SKILL.md)
- 候補 focused: [actor-goal](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-actor-goal-generation/SKILL.md)、[lifecycle](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-lifecycle-generation/SKILL.md)、[state-transition](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-state-transition-generation/SKILL.md)、[失敗](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-failure-generation/SKILL.md)、[external-integration](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-external-integration-generation/SKILL.md)、[operation-audit](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-operation-audit-generation/SKILL.md)
- シナリオ: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-design/SKILL.md)
- implementation 対象範囲: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-scope/SKILL.md)
- wall discussion: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/wall-discussion/SKILL.md)
- diagramming: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/diagramming/SKILL.md)
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/ui-design/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-design/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-scope/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/skill-modification/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/wall-discussion/SKILL.md

## 内部参照規約

### 実装対象範囲判定条件

implementation-scope を扱う時は、Codex 実装レーンの agent 起動 の token 量を事前計算しない。
代わりに [implementation-scope](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-scope/SKILL.md) の 引き継ぎ分割規約 と 規模判定条件 に従い、論理境界と規模の目安で分割する。

各 引き継ぎ は原則として `1 受け入れユースケース × 1 検証 intent` に収める。
Codex implementation レーンから 対象範囲 過大で 戻し された場合は、既存 承認 を維持せず `pending-human-review` に戻す。

### シナリオ完備判定条件

scenario-design は、抽象要件から直接 シナリオ を作って完了にしない。
`designer` は `implement_lane` が揃えた 6 種の `scenario-candidates.<viewpoint>.md` を読み、候補の重複、採用、統合、不採用、競合を固定してから シナリオ表 を作る。
`designer` は候補生成器を再 起動 しない。
[scenario-design](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-design/SKILL.md) の詳細要求タイプを使い、明示的ではない判断を先に検出する。
詳細要求タイプの仕様網羅は `scenario-design.requirement-coverage.json` に分ける。
シナリオ 候補の採否と競合は `scenario-design.candidate-coverage.json` に分ける。
人間向け質問票は `scenario-design.questions.md` に分ける。
`scenario-design.md` に長い JSON や質問票本文を埋め込まない。

design bundle を 人間レビュー へ進める条件は次の通り。

- 必要な詳細要求タイプが `explicit`、`derived`、`not_applicable`、`deferred` のいずれかに分類されている
- 6 種の シナリオ 候補成果物 が task folder に存在する
- `scenario-design.candidate-coverage.json` で全 候補 の採用、統合、不採用、競合、要人間判断が分類されている
- `not_applicable` と `deferred` には理由がある
- `needs_human_decision` が 0 件である
- 未解決 conflict が 0 件である
- 人間判断が必要な項目がある場合は、シナリオ 完了ではなく `scenario-design.questions.md` 出力で停止している

## 判断規約

- 判断は入力 成果物、外部参照規約、対象 agent の責務境界に従う。

## 非対象規約

- 作業流れ順序決定、task folder 進行管理、作業前確認は扱わない。
- 実画面 observation、docs 正本化、プロダクト実装は扱わない。
- ツール権限、agent 実行定義、プロダクト仕様正本は変更しない。

## 出力規約

- 出力は判断結果、根拠参照、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。

### Handoff

- 引き継ぎ先: `implement_lane`
- 渡す対象範囲: 設計成果物、人間レビュー 状態、未回答質問
- 返却先: implement_lane
- 対象成果物: 扱った シナリオ、シナリオ候補 統合、UI、implementation-scope、skill-modification の状態を返す。
- 変更成果物: 作成または更新した task 内成果物 path を返す。
- 人間レビュー 状態: 人間レビュー が必要な判断、承認待ち、承認済みの状態を返す。
- 確認結果: 実行した確認と未実行理由を返す。
- 引き継ぎ または停止理由: `implement_lane` へ戻す理由または停止理由を返す。
- 未決事項: 設計継続に必要な未決事項を返す。

## 完了規約

- task 内成果物 が承認状態、根拠参照、未決事項を含んでいる。
- 人間レビュー が必要な判断を AI だけで完了扱いにしていない。
- 必須根拠として、根拠成果物 path、必要な 人間承認 記録、実行した 検証結果 がある。
- 完了判断材料として、`implement_lane` が次の 作業流れ action、人間レビュー、人間向け Codex implementation レーン 引き継ぎ を判断できる情報が返っている。
- 残留リスクとして、設計継続に必要な未決事項が返っている。

## 停止規約

- scenario-design に `needs_human_decision` または未解決 conflict が残る場合は、質問票を返して human 回答待ちにする。
- シナリオ 候補成果物 が不足する場合は、`implement_lane` に戻し、候補生成器の不足を解消してから再開する。
- 作業流れ 順序決定 や task folder 進行管理 が主目的なら `implement_lane` へ戻す。
- 作業前の影響範囲、実行計画、検証方法の確認が不足する場合は `implement_lane` へ戻す。
- 実画面 observation が必要なら `investigator` を使う前提で `implement_lane` へ戻す。
- docs 正本化が必要なら human 承認後に `docs_updater` を使う前提で `implement_lane` へ戻す。
- プロダクト 実装が必要なら `implement_lane` へ戻し、人間向け Codex implementation レーン 引き継ぎ の扱いを判断させる。
- 停止時は不足項目、衝突箇所、戻し先を返す。
- 拒否条件: 作業流れの進行管理要求
- 拒否条件: 不足 引き継ぎ入力
- 拒否条件: プロダクト実装 request
- 拒否条件: unapproved docs 正本化
- 拒否条件: implementation-owned implementation-time work
- 停止条件: 人間レビュー が必要な判断を AI だけで確定しそうである
- 停止条件: scenario-design に必要な シナリオ候補成果物 が不足している
- 停止条件: シナリオ候補 coverage に未解決 conflict が残っている
- 停止条件: 進行中 task folder が不足している
- 停止条件: 設計範囲 が不明である
- 停止条件: 引き継ぎ入力 だけでは作業できない
- 停止条件: プロダクト 実装が必要である
