---
name: implementation-distill
description: GitHub Copilot 側の実装前文脈整理の共通知識 package。single_handoff_packet 1 件から実装に必要な facts を圧縮する判断基準を提供する。
---

# Implementation Distill

## 目的

`implementation-distill` は知識 package である。
`implementation-distiller` agent が、`single_handoff_packet` 1 件から lane_context_packet を作る時の共通判断を提供する。

実行権限、write scope、active contract、handoff は [implementation-distiller.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/implementation-distiller.agent.md) が持つ。

## いつ参照するか

- 実装前に facts、constraints、gaps を圧縮する時
- required reading と code pointer を owned_scope に絞る時
- validation entry を実装前に明示する時

## 参照しない場合

- product code または product test を変更する時
- review や UI check を行う時
- design 不足を実装側で補う時

## 知識範囲

- path catalog から必要 file だけを summary / full に上げる圧縮
- fix_ingredients と distracting_context の分離
- first_action、change_targets、related_code_pointers の具体化
- 要件、実装方針、決定事項の要約
- facts、inferred、gap の分離
- single_handoff_packet と owned_scope の対応づけ
- focused skill の選び方

## 原則

- `single_handoff_packet` 1 件を唯一の source scope にする
- owned_scope に関係する code pointer を優先する
- patch 生成に必要な fix ingredients を file / function / block の構造単位で残す
- 類似していても修正に不要な context は distracting_context として分ける
- repository method、interface、field の追加が必要そうに見えても、実 code で present / absent を確認するまで fact にしない
- first_action は 1 completion_signal clause に限定し、partial や複数 clause にしない
- 1 edit で clause が閉じない場合は、同じ clause の最小 closure chain を上流 symbol から leaf まで残す
- existing_patterns が none なら、探索範囲と実装判断への影響を残す
- validation entry は最初に試せる cheap check を優先する
- 実 code を読んでから first_action を返す
- 実装開始点は path、symbol、line number、変更種別で返す
- 要件、実装方針、決定事項は distiller が要約し、implementer の再読を原則不要にする
- 実装案を増やさず、実装に必要な制約だけを残す
- 設計不足は実装せず戻す

## Focused Skills

- [implementation-distill-implement](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-distill-implement/SKILL.md): 新機能や拡張の実装前 context
- [implementation-distill-fix](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-distill-fix/SKILL.md): fix scope、再現証跡、修正対象
- [implementation-distill-refactor](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-distill-refactor/SKILL.md): 不変条件、依存境界、preserved behavior

## DO / DON'T

DO:
- fix_ingredients に path、symbol/type/function、line number、why_needed_for_patch を残す
- distracting_context に why_excluded と risk_if_read を残す
- 実装者が最初に触る file、symbol、line number、変更種別を残す
- requirements_policy_decisions に要件、実装方針、決定事項、out of scope、禁止事項を要約する
- repository method、interface、field の有無を present / absent の code fact として残す
- first_action が閉じる clause を 1 つだけ明示する
- existing_patterns がない場合は searched scope と impact を添えて none とする
- required_reading は読む目的と symbol を添えて順序づける
- 要件や決定事項の原文は、要約では判断できない時だけ required_reading に残す
- related_code_pointers は path、symbol/type/function、line number、読み取った事実を残す
- validation entry を明示する
- gap と residual risk を分ける

DON'T:
- fix_ingredients を特定せず first_action だけを返さない
- 類似 context を required_reading に混ぜない
- 存在確認していない repository method、interface、field を追加前提にしない
- first_action に partial、複数 clause、曖昧な advance を書かない
- existing_patterns の none を探索範囲なしで返さない
- cheap validation を検討せず広い command だけを返さない
- 実 code を読まず handoff の文章を言い換えない
- 要件、実装方針、決定事項を要約せず required_reading に丸投げしない
- required_reading をファイル名だけの列挙にしない
- implementer に「どこから調べるか」を委ねない
- product code / product test を変更しない
- 要件や設計を追加しない
- full implementation-scope、active work plan 全文、source artifacts、後続 handoff を要求しない
- mode 別 active contract を使わない

## 参照パターン

- [context-compression-patterns.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-distill/references/patterns/context-compression-patterns.md) を参照する。
- 対象は entry point discovery、execution flow、architecture layer mapping、dependency documentation である。
- Wails + Go + Svelte 境界に沿って facts、inferred、gaps を分ける。

## Checklist

- [implementation-distill-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-distill/references/checklists/implementation-distill-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。

## Agent が持つもの

- active contract: [implementation-distiller.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/references/implementation-distiller/contracts/implementation-distiller.contract.json)
- permissions: [permissions.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/references/implementation-distiller/permissions.json)

## Maintenance

- implement / fix / refactor の知識差分は focused skill に置く。
- output obligation を skill 本体へ戻さない。
- 旧 mode guide は active 正本として扱わない。
