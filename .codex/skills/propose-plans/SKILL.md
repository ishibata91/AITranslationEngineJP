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

## いつ参照するか

- 新規または継続 task の workflow state を整理する時
- distiller、designer、investigator が必要か最初に判定する時
- `docs/exec-plans/active/<task-id>/` と `plan.md` の扱いを判断する時
- human review と人間向け Copilot handoff の境界を確認する時
- Copilot 修正完了後に正本化が必要か判断する時

## 参照しない場合

- requirements、UI、scenario、implementation-scope、diagram の本文を作る時
- 再現、trace、実画面 observation を行う時
- docs 正本を更新する時
- product code または product test を直接変更する時

## 知識範囲

- task folder と `plan.md` の役割
- 最初に行う必要判定
- distiller、designer、investigator の独立 spawn 判断
- context を引き継がない handoff packet の作り方
- design bundle、human review、human Copilot handoff、Copilot 完了後正本化の順序
- work_history 用の Codex report 材料を整理するタイミング

## 原則

- いずれの工程も、必要かどうかを最初に判定する
- task に関連する事柄と必要資料の判断材料は、必要なら `distiller` で集める
- 資料作成は必要なら `designer` に渡す。diagram も designer の資料作成に含める
- 実画面観測は必要なら `investigator` に渡す
- spawned agent へ context を引き継がず、必要情報を packet に明示する
- Copilot handoff は Codex が直接渡さず、人間へ返す
- Copilot の修正完了が分かってから正本化へ進む
- closeout、停止、reroute 時は `codex-work-reporting` を参照し、最後に必ず報告材料を作る

## 標準パターン

1. memory、[README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/README.md)、[index.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md) を確認する。
2. active / completed task folder に同種 task がないか確認する。
3. distiller、designer、investigator が必要か最初に判定する。
4. 必要なら `distiller` を context 継承なしで spawn し、task 関連情報と必要資料の判断材料を集める。
5. 必要なら `designer` を context 継承なしで spawn し、requirements、UI、scenario、implementation-scope、diagram を含む資料を作る。
6. 必要なら `investigator` を context 継承なしで spawn し、実画面や観測対象を確認する。
7. 戻りを `plan.md` の workflow state に反映する。
8. design bundle 完了後に human review で停止する。
9. 承認後、人間が Copilot に渡せる handoff packet を返す。
10. Copilot の修正完了が分かった後、必要なら正本化へ進む。
11. closeout、停止、reroute 時は `codex-work-reporting` を参照し、`work_history` へ転記できる Codex report 材料を最後に必ず作る。

## DO / DON'T

DO:
- 論理名と actual skill / agent 名を同じ行に置く
- 最初に必要判定を明示する
- spawn packet に読む file、禁止事項、期待出力を書く
- human review gate と再開条件を見える化する

DON'T:
- spawned agent に会話文脈を暗黙継承させない
- diagram を別 flow として扱わない
- Codex から Copilot へ直接渡さない
- Copilot 修正完了前に正本化しない

## Checklist

- [propose-plans-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/propose-plans/references/checklists/propose-plans-checklist.md) を参照する。

## References

- workflow: [README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/README.md)
- docs index: [index.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md)
- agent spec: [propose_plans.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/propose_plans.agent.md)
- agent contract: [propose_plans.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/propose_plans/contracts/propose_plans.contract.json)
- report skill: [codex-work-reporting](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-work-reporting/SKILL.md)

## Maintenance

- `propose-plans` は workflow orchestration 知識だけを持つ。
- detailed design、investigation、diagram、docs 正本化の知識をこの skill に戻さない。
- workflow と actual skill / agent 名の対応を曖昧にしない。
- implementation lane の詳細は [.github/skills](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/) に置く。
