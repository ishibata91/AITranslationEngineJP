---
name: design-bundle
description: Codex 側の design artifact 進行 skill。必須要件、UI、scenario、implementation-scope を task-local artifact として固定するための source of truth、進め方、handoff を提供する。
---
# Design Bundle

## 目的

`design-bundle` は作業プロトコルである。
`designer` agent と top-level Codex が、必須要件、UI、scenario、implementation-scope を task-local artifact として固定する時の、人間可読な実行説明の正本として使う。

workflow の次 action 判断、task folder orchestration、人間向け Codex implementation lane handoff の返却は `implement_lane` が担当する。
プロダクトコードとプロダクトテスト は変更しない。

## 対応ロール

- `designer` が使う。
- 呼び出し元は `implement_lane` または人間とする。
- 返却先は human review または `implement_lane` とする。
- owner artifact は `design-bundle` の出力規約で固定する。

## 入力規約

- 入力は caller から渡された task-local artifact、source_ref、必要な承認状態を含む。
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。
- 必須入力: caller, handoff_packet, user_instruction_or_task_summary, active_task_folder, design_scope, lane_owner
- 任意入力: human_review_record, target_skill, existing_design_artifacts, scenario_candidate_artifacts, known_gaps
- input_policy: handoff_packet だけで作業できること。引き継いでいない会話文脈に依存しない。
- 必須 artifact: /Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md, active task folder

## 外部参照規約

- エージェント実行定義とツール権限は [designer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.toml) の `allowed_write_paths` / `allowed_commands` とする。
- secondary: packet に明示された関連 docs、関連 skill、human の現在指示
- エージェント実行定義: [designer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.toml)
- ツール権限: エージェント実行定義の `allowed_write_paths` / `allowed_commands` に従う
- ui: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/ui-design/SKILL.md)
- candidate common: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-candidate-generation/SKILL.md)
- candidate focused: [actor-goal](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-actor-goal-generation/SKILL.md)、[lifecycle](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-lifecycle-generation/SKILL.md)、[state-transition](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-state-transition-generation/SKILL.md)、[failure](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-failure-generation/SKILL.md)、[external-integration](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-external-integration-generation/SKILL.md)、[operation-audit](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-operation-audit-generation/SKILL.md)
- scenario: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-design/SKILL.md)
- implementation scope: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-scope/SKILL.md)
- wall discussion: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/wall-discussion/SKILL.md)
- diagramming: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/diagramming/SKILL.md)
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/ui-design/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-design/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-scope/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/skill-modification/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/wall-discussion/SKILL.md

## 内部参照規約

### Implementation Scope Gate

implementation-scope を扱う時は、Codex implementation lane spawn_agent の token 量を事前計算しない。
代わりに [implementation-scope](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-scope/SKILL.md) の Handoff Split Rule と Size Gate に従い、論理境界と規模の目安で分割する。

各 handoff は原則として `1 受け入れユースケース × 1 validation intent` に収める。
Codex implementation laneから scope 過大で reroute された場合は、既存 approval を維持せず `pending-human-review` に戻す。

### Scenario Completeness Gate

scenario-design は、抽象要件から直接 scenario を作って完了にしない。
`designer` は `implement_lane` が揃えた 6 種の `scenario-candidates.<viewpoint>.md` を読み、候補の重複、採用、統合、不採用、競合を固定してから scenario matrix を作る。
`designer` は候補生成器を再 spawn しない。
[scenario-design](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-design/SKILL.md) の詳細要求タイプを使い、明示的ではない判断を先に検出する。
詳細要求タイプの仕様網羅は `scenario-design.requirement-coverage.json` に分ける。
scenario 候補の採否と競合は `scenario-design.candidate-coverage.json` に分ける。
人間向け質問票は `scenario-design.questions.md` に分ける。
`scenario-design.md` に長い JSON や質問票本文を埋め込まない。

design bundle を human review へ進める条件は次の通り。

- 必要な詳細要求タイプが `explicit`、`derived`、`not_applicable`、`deferred` のいずれかに分類されている
- 6 種の scenario candidate artifact が task folder に存在する
- `scenario-design.candidate-coverage.json` で全 candidate の採用、統合、不採用、競合、要人間判断が分類されている
- `not_applicable` と `deferred` には理由がある
- `needs_human_decision` が 0 件である
- 未解決 conflict が 0 件である
- 人間判断が必要な項目がある場合は、scenario 完了ではなく `scenario-design.questions.md` 出力で停止している

## 判断規約

- 判断は入力 artifact、外部参照規約、対象 agent の責務境界に従う。
- 対象外の成果物、ツール権限、product 仕様正本はこの skill で変更しない。

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。

### Handoff

- handoff 先: `implement_lane`
- 渡す scope: design artifact、human review 状態、open questions
- 返却先: implement_lane
- 対象成果物: 扱った scenario、scenario candidate integration、UI、implementation-scope、skill-modification の状態を返す。
- 変更成果物: 作成または更新した task-local artifact path を返す。
- human review 状態: human review が必要な判断、承認待ち、承認済みの状態を返す。
- 確認結果: 実行した確認と未実行理由を返す。
- handoff または停止理由: `implement_lane` へ戻す理由または停止理由を返す。
- 未決事項: 設計継続に必要な未決事項を返す。

## 完了規約

- task-local artifact が承認状態、source_ref、未決事項を含んでいる。
- human review が必要な判断を AI だけで完了扱いにしていない。
- 必須根拠として、source artifact path、必要な human approval record、実行した validation result がある。
- 完了判断材料として、`implement_lane` が次の workflow action、human review、人間向け Codex implementation lane handoff を判断できる情報が返っている。
- 残留リスクとして、設計継続に必要な未決事項が返っている。

## 停止規約

- scenario-design に `needs_human_decision` または未解決 conflict が残る場合は、質問票を返して human 回答待ちにする。
- scenario candidate artifact が不足する場合は、`implement_lane` に戻し、候補生成器の不足を解消してから再開する。
- workflow sequencing や task folder orchestration が主目的なら `implement_lane` へ戻す。
- 文脈圧縮が必要なら `implement_lane` へ戻す。
- 実画面 observation が必要なら `investigator` を使う前提で `implement_lane` へ戻す。
- docs 正本化が必要なら human 承認後に `docs_updater` を使う前提で `implement_lane` へ戻す。
- product 実装が必要なら `implement_lane` へ戻し、人間向け Codex implementation lane handoff の扱いを判断させる。
- 停止時は不足項目、衝突箇所、戻し先を返す。
- 拒否条件: workflow orchestration request
- 拒否条件: missing handoff packet
- 拒否条件: product implementation request
- 拒否条件: unapproved docs canonicalization
- 拒否条件: implementation-owned implementation-time work
- 停止条件: human review が必要な判断を AI だけで確定しそうである
- 停止条件: scenario-design に必要な scenario_candidate_artifacts が不足している
- 停止条件: scenario candidate coverage に未解決 conflict が残っている
- 停止条件: active task folder が不足している
- 停止条件: design_scope が不明である
- 停止条件: handoff_packet だけでは作業できない
- 停止条件: product 実装が必要である
