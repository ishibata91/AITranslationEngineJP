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
- facts、inferred、gap の分離
- single_handoff_packet と owned_scope の対応づけ
- focused skill の選び方

## 原則

- `single_handoff_packet` 1 件を唯一の source scope にする
- owned_scope に関係する code pointer を優先する
- 実装案を増やさず、実装に必要な制約だけを残す
- 設計不足は実装せず戻す

## Focused Skills

- [implementation-distill-implement](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-distill-implement/SKILL.md): 新機能や拡張の実装前 context
- [implementation-distill-fix](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-distill-fix/SKILL.md): fix scope、再現証跡、修正対象
- [implementation-distill-refactor](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-distill-refactor/SKILL.md): 不変条件、依存境界、preserved behavior

## DO / DON'T

DO:
- 実装者が最初に読む file と順番を残す
- validation entry を明示する
- gap と residual risk を分ける

DON'T:
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
