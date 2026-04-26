---
name: propose-plans
description: Codex workflow orchestration 知識 package。必要判定、distill、designer、investigator、human review、人間向け Copilot handoff、正本化入口の進め方を提供する。
---

# Propose Plans

## 目的

`propose-plans` は workflow orchestration の知識 package である。
`propose_plans` agent が Codex workflow を進めるために、必要判定、task folder、agent spawn、human review、人間向け Copilot handoff、正本化入口の判断基準を提供する。

この skill はオーケストレーションの考え方だけを持つ。
資料作成、実画面観測、docs 正本化、product 実装は、それぞれの agent、人間、Copilot lane に渡す。

## 原則

- いずれの工程も、必要かどうかを最初に判定する
- task に関連する事柄と必要資料の判断材料は、必要なら `distiller` で集める
- 承認済み design bundle がない限り、資料作成は `designer` に渡す
- 実画面観測は必要なら `investigator` に渡す
- spawned agent へ context を引き継がず、必要情報を packet に明示する
- human review 前に `implementation-scope` を作らない
- scenario-design の `needs_human_decision` が残る場合は、design bundle review へ進めず質問票回答待ちにする
- Copilot handoff は Codex が直接渡さず、人間へ返す
- Copilot の修正完了が分かってから正本化へ進む
- closeout、停止、reroute 時は `work_reporter` に渡す benchmark score と completion evidence を整理し、最後に必ず報告材料を作る

## Runtime Boundary

- binding: [propose_plans.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/propose_plans.toml)
- permissions: [permissions.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/propose_plans/permissions.json)
- contract: [propose_plans.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/propose_plans/contracts/propose_plans.contract.json)
- allowed: task folder と `plan.md` の作成、更新、closeout、独立 agent spawn の packet 作成、人間向け handoff packet 作成、work_reporter への report packet 作成
- forbidden: product code、product test、docs 正本、詳細設計 artifact 本文の代筆、Copilot への直接 handoff
- write scope: `docs/exec-plans/active/` と `docs/exec-plans/completed/` の workflow state

## 標準パターン

1. memory、[README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/README.md)、[index.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md) を確認する。
2. active / completed task folder に同種 task がないか確認する。
3. distiller と investigator の要否を判定し、承認済み design bundle がない限り `designer` を使う。
4. `distiller` を context 継承なしで spawn し、`task_frame`、`canonical_evidence`、`code_evidence`、`effective_prior_decisions`、`observation_evidence` の available input を packet に明示して渡す。
5. `designer` を context 継承なしで spawn し、`scenario-design` を必須で作り、詳細要求タイプの未決検出と質問票出力を含める。UI 変更がある時だけ `ui-design` を追加する。
6. 必要なら `investigator` を context 継承なしで spawn し、実画面や観測対象を確認する。
7. 戻りを `plan.md` の workflow state に反映する。
8. `needs_human_decision` が残る場合は質問票を人間へ返して停止し、0 件になった design bundle だけ human review へ進める。
9. 承認後に `designer` を再度 context 継承なしで spawn し、`implementation-scope` を固定する。
10. 人間が Copilot に渡せる handoff packet を返す。
11. Copilot の修正完了が分かった後、必要なら正本化へ進む。
12. closeout、停止、reroute 時は `work_reporter` へ渡せる benchmark score と completion evidence を整理し、`work_history` へ転記できる report 材料を最後に必ず作る。

## Stop / Reroute

- distill、資料作成、実画面観測の必要判定ができない場合は停止する。
- human review が必要な判断を AI だけで確定しそうな場合は停止する。
- scenario-design の質問票に未回答項目がある場合は停止する。
- Copilot handoff packet を人間へ返す前に implementation-scope が不足している場合は停止する。
- Copilot の修正完了が分からない場合は正本化へ進まない。

## Handoff

- Codex spawn 先: `distiller`, `designer`, `investigator`, `docs_updater`, `work_reporter`
- human handoff 先: `human`
- Copilot handoff: `propose_plans` は直接渡さず、人間に packet を返す
- 渡す contract: [propose_plans.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/propose_plans/contracts/propose_plans.contract.json)
- 渡す scope: task id、現在状態、必要判定、読む file、禁止事項、期待出力、再開条件

## DO / DON'T

DO:
- 論理名と actual skill / agent 名を同じ行に置く
- 最初に必要判定を明示する
- `distiller` packet では入口 evidence の種類を明示する
- spawn packet に読む file、禁止事項、期待出力を書く
- `designer` を optional artifact writer として扱わない
- human review gate と再開条件を見える化する

DON'T:
- spawned agent に会話文脈を暗黙継承させない
- Codex から Copilot へ直接渡さない
- Copilot 修正完了前に正本化しない

## Checklist

- [propose-plans-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/propose-plans/references/checklists/propose-plans-checklist.md) を参照する。

## References

- workflow: [README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/README.md)
- docs index: [index.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md)
- binding: [propose_plans.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/propose_plans.toml)
- agent contract: [propose_plans.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/propose_plans/contracts/propose_plans.contract.json)
- report skill: [codex-work-reporting](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-work-reporting/SKILL.md)

## Maintenance

- `propose-plans` は workflow orchestration 知識だけを持つ。
- detailed design、investigation、docs 正本化の知識をこの skill に戻さない。
- workflow と actual skill / agent 名の対応を曖昧にしない。
- implementation lane の詳細は [.github/skills](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/) に置く。
