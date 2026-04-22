---
name: review
description: GitHub Copilot 側の review 共通知識 package。UI check と implementation review の判断基準を提供する。
---

# Review

## 目的

`review` は知識 package である。
`reviewer` agent が、実装結果を `single_handoff_packet` と `lane_context_packet` に照合する時の共通判断を提供する。

実行権限、write scope、active contract、handoff は [reviewer.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/reviewer.agent.md) が持つ。

## いつ参照するか

- 実装差分を owned_scope と照合する時
- UI evidence を確認する時
- reroute すべき finding を返す時

## 参照しない場合

- design review を行う時
- 実装修正を同時に行う時
- 追加の再現や trace が先に必要な時

## 知識範囲

- single_handoff_packet と lane_context_packet に基づく review 判断
- actionable finding の返し方
- UI evidence、coverage 70%、repo-local Sonar issue gate、harness evidence の扱い
- focused skill の選び方

## 原則

- `single_handoff_packet` と `lane_context_packet` を判定の正本にする
- 好みや将来改善で判定しない
- design review をしない
- 修正と review を混ぜない

## Focused Skills

- [review-implementation](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/review-implementation/SKILL.md): implementation review
- [review-ui-check](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/review-ui-check/SKILL.md): UI check

## DO / DON'T

DO:
- findings は再現できる形で返す
- coverage 70%、repo-local Sonar issue gate、harness evidence を確認する
- 追加証跡が必要なら investigator へ戻せる理由を書く
- pass の場合も残リスクを明記する

DON'T:
- 新しい要件解釈を追加しない
- docs、`.codex`、`.github/skills`、`.github/agents` を変更しない
- mode 別 active contract を使わない

## 参照パターン

- [review-patterns.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/review/references/patterns/review-patterns.md) を参照する。
- 対象は confidence-based filtering、severity ordering、security / correctness / regression priority、silent failure detection である。
- output は `reviewer.contract.json` の fields に合わせて返す。

## Checklist

- [review-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/review/references/checklists/review-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。

## Agent が持つもの

- active contract: [reviewer.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/references/reviewer/contracts/reviewer.contract.json)
- permissions: [permissions.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/references/reviewer/permissions.json)

## Maintenance

- UI check と implementation review の知識差分は focused skill に置く。
- output obligation を skill 本体へ戻さない。
- 旧 mode contract は active 正本として扱わない。
