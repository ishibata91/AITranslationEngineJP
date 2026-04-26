---
name: design-bundle
description: Codex 側の design artifact 進行 skill。必須要件、UI、scenario、implementation-scope を task-local artifact として固定するための source of truth、進め方、handoff を提供する。
---

# Design Bundle

## 目的

`design-bundle` は知識 package である。
`designer` agent と top-level Codex が、必須要件、UI、scenario、implementation-scope を task-local artifact として固定する時の、人間可読な実行説明の正本として使う。

workflow の次 action 判断、task folder orchestration、人間向け Copilot handoff の返却は `propose_plans` が担当する。
product code と product test は変更しない。

## 参照 skill

- `scenario-design`: 必須要件、受け入れテスト観点、システムテスト分類、validation を参照する。
- `ui-design`: UI 要件契約と実装後確認観点を参照する。
- `implementation-scope`: human review 後の人間向け Copilot handoff 粒度を参照する。
- `wall-discussion`: read-only 壁打ちの質問設計を参照する。
- `diagramming`: diagram を必要資料として扱う時に参照する。
- `skill-modification`: skill / agent の境界整理を参照する。

## Source Of Truth

- primary: `propose_plans` から渡された handoff packet、active task folder、[README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/README.md)、[index.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md)
- secondary: packet に明示された関連 docs、関連 skill、human の現在指示
- forbidden source: 未承認の design review、旧 flat plan、Copilot の独自再設計、引き継いでいない会話文脈

## Runtime Boundary

- permissions: [permissions.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/designer/permissions.json)
- contract: [designer.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/designer/contracts/designer.contract.json)
- allowed: task-local design artifact を作成、更新、整理する
- forbidden: product code、product test、未承認 docs 正本を変更しない
- write scope: `docs/exec-plans/active/`、`.codex/` の workflow 範囲

## Implementation Scope Gate

implementation-scope を扱う時は、Copilot 側 RunSubagent の token 量を事前計算しない。
代わりに [implementation-scope](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-scope/SKILL.md) の Handoff Split Rule と Size Gate に従い、論理境界と規模の目安で分割する。

各 handoff は原則として `1 受け入れユースケース × 1 validation intent` に収める。
Copilot 側から scope 過大で reroute された場合は、既存 approval を維持せず `pending-human-review` に戻す。

## Scenario Completeness Gate

scenario-design は、抽象要件から直接 scenario を作って完了にしない。
[scenario-design](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-design/SKILL.md) の詳細要求タイプを使い、明示的ではない判断を先に検出する。
詳細要求タイプの仕様網羅は `scenario-design.requirement-coverage.json` に分ける。
人間向け質問票は `scenario-design.questions.md` に分ける。
`scenario-design.md` に長い JSON や質問票本文を埋め込まない。

design bundle を human review へ進める条件は次の通り。

- 必要な詳細要求タイプが `explicit`、`derived`、`not_applicable`、`deferred` のいずれかに分類されている
- `not_applicable` と `deferred` には理由がある
- `needs_human_decision` が 0 件である
- 人間判断が必要な項目がある場合は、scenario 完了ではなく `scenario-design.questions.md` 出力で停止している

## 標準パターン

1. `propose_plans` から渡された handoff packet を確認する。
2. handoff packet にない暗黙の会話文脈へ依存しない。
3. 必要な design skill と checklist を読む。
4. `scenario-design` を必須にし、詳細要求タイプと明示性 gate を通す。
5. UI 変更がある時だけ `ui-design` を扱い、`implementation-scope` は human review 後にだけ扱う。
6. task-local artifact と source of truth を分ける。
7. human review が必要な地点で停止する。
8. 作成、更新、未決事項、検証結果を `propose_plans` へ返す。

## Stop / Reroute

- scenario-design に `needs_human_decision` が残る場合は、質問票を返して human 回答待ちにする。
- workflow sequencing や task folder orchestration が主目的なら `propose_plans` へ戻す。
- 文脈圧縮が必要なら `propose_plans` へ戻す。
- 実画面 observation が必要なら `investigator` を使う前提で `propose_plans` へ戻す。
- docs 正本化が必要なら human 承認後に `docs_updater` を使う前提で `propose_plans` へ戻す。
- product 実装が必要なら `propose_plans` へ戻し、人間向け Copilot handoff の扱いを判断させる。

## Handoff

- handoff 先: `propose_plans`
- 渡す contract: [designer.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/designer/contracts/designer.contract.json)
- 渡す scope: design artifact、human review 状態、open questions

## References

- ui: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/ui-design/SKILL.md)
- scenario: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-design/SKILL.md)
- implementation scope: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-scope/SKILL.md)
- wall discussion: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/wall-discussion/SKILL.md)
- diagramming: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/diagramming/SKILL.md)

## Maintenance

- `designer` の人間可読な実行説明は、この skill を正本にする。
- hard permission と output obligation は JSON に残す。
- `.agent.md` を再導入しない。
