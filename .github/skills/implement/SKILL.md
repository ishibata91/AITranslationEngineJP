---
name: implement
description: GitHub Copilot 側の product code 実装の共通知識 package。承認済み owned_scope を実装する判断基準を提供する。
---

# Implement

## 目的

`implement` は知識 package である。
`implementer` agent が、承認済み `implementation-scope` の handoff 1 件を owned_scope 内へ実装する時の共通判断を提供する。

実行権限、write scope、active contract、handoff は [implementer.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/implementer.agent.md) が持つ。

## いつ参照するか

- owned_scope 内の product code を実装する時
- lane_context_packet に基づいて product code を実装する時
- scenario 先行時の tester output を product code 実装へ反映する時
- lane-local validation の扱いを確認する時

## 参照しない場合

- 実装前 context 整理だけを行う時
- UI check や implementation review を行う時
- docs や workflow 文書を変更する時

## 知識範囲

- owned_scope を超えない実装判断
- handoff 資料のスコープ粒度に合わせる判断
- coding guidelines と既存 pattern の確認
- boundary、error path、test surface の実装品質判断
- validation result と residual risk の返し方
- focused skill の選び方

## 原則

- `implementation-scope` と owned_scope を超えない
- handoff 資料のスコープ粒度で実装する
- lane_context_packet に合わせて product code だけを変更する
- scenario 先行時だけ tester output も確認する
- implementation_subscope が渡された場合はその sub-scope 内だけを実装する
- 実装完了後、handoff を終える前に touched layer に対応する local validation を実行する
- fix_ingredients に対応する code path を優先し、distracting_context へ寄り道しない
- first_action と change_targets から着手する
- insufficient_context_criteria は structural gate とし、fix_ingredients、first_action、change_targets、requirements_policy_decisions、existing pattern、validation_entry の不足時に返す
- first_action が 1 clause に固定されていない、line / symbol / public seam が不明、closure chain がない場合は insufficient_context を返す
- listed required_reading 内の局所確認、既存 pattern への通常追従、lane-local validation failure は not_insufficient_context として扱う
- 既存 pattern、naming、layer に合わせる
- broad refactor を混ぜない
- product test、fixture、snapshot、test helper は tester が扱う
- docs 正本化をしない

## Focused Skills

- [implement-backend](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implement-backend/SKILL.md): backend layer と lane-local validation
- [implement-frontend](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implement-frontend/SKILL.md): UI state と Wails bridge
- [implement-mixed](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implement-mixed/SKILL.md): API / Wails / DTO / gateway など frontend と backend の接合点 scope
- [implement-fix-lane](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implement-fix-lane/SKILL.md): accepted fix scope の恒久修正

## DO / DON'T

DO:
- 実装前に [coding-guidelines.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/coding-guidelines.md) を読む
- lane_context_packet の fix_ingredients、distracting_context、first_action、change_targets、requirements_policy_decisions、related_code_pointers を確認する
- implementation_subscope があれば completion_signal clause、public seam、change target / symbol、validation command を確認する
- insufficient_context を返す場合は reason、needed_context、suggested_narrowing_axis、remaining_implementation_subscopes を structural gate に対応づける
- entry point、call site、data flow、error path、test surface を確認する
- 既存 pattern に naming、constructor、DI、error return を合わせる
- lane-local validation 結果または未実行理由を返す
- backend handoff は `python3 scripts/harness/run.py --suite backend-local`、frontend handoff は `python3 scripts/harness/run.py --suite frontend-local` を使う
- mixed handoff は touched layer に応じて両方を実行する
- touched files は product code だけにする

DON'T:
- 要件や設計を追加しない
- fix_ingredients がないまま実装を始めない
- insufficient_context を返さず広い調査で不足 context を埋めない
- criteria mismatch になる不安や通常の局所確認を insufficient_context にしない
- implementation_subscope 外へ実装を広げない
- distracting_context を実装対象に混ぜない
- first_action がないまま広い調査を始めない
- config、lint、test、coverage 設定を変更して gate を回避しない
- product test、fixture、snapshot、test helper を変更しない
- coverage、harness all、repo-local Sonar issue gate を implementer の必須 closeout にしない
- owned_scope 外の cleanup、rename、format を混ぜない
- docs、`.codex`、`.github/skills`、`.github/agents` を変更しない
- mode 別 active contract を使わない

## 参照パターン

- [implementation-quality-patterns.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implement/references/patterns/implementation-quality-patterns.md) を参照する。
- 対象は readability、KISS、DRY、YAGNI、error handling、backend / frontend boundary、minimal build fix である。
- Svelte、Wails gateway、Go backend の責務境界に沿って判断する。

## Checklist

- [implement-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implement/references/checklists/implement-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。

## Agent が持つもの

- active contract: [implementer.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/references/implementer/contracts/implementer.contract.json)
- permissions: [permissions.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/references/implementer/permissions.json)

## Maintenance

- backend / frontend / mixed / fix-lane の知識差分は focused skill に置く。
- output obligation を skill 本体へ戻さない。
