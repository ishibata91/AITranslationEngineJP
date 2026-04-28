---
name: implementation-distill-refactor
description: Codex implementation lane 側の refactor 向け context 圧縮知識 package。
---

# Implementation Distill Refactor

## 目的

この skill は知識 package である。
`implementation_distiller` agent が refactor handoff を整理する時に、不変条件、依存境界、preserved behavior を抽出する判断基準を提供する。

## いつ参照するか

- refactor handoff を実装前 packet に圧縮する時
- 変更してはいけない振る舞いを整理する時
- affected package / component / tests をまとめる時

## 参照しない場合

- 新機能実装の facts を整理する時
- fix 再現条件を整理する時
- broad refactor を新たに提案する時

## 原則

- 不変条件、依存境界、変更候補を別々に圧縮する
- preserved behavior を守るための fix_ingredients を構造単位で残す
- refactor に似ているだけの周辺 context は distracting_context に分ける
- 似た責務の file は cluster としてまとめる
- 実装手順ではなく、守る制約を残す
- refactor 開始点は path、symbol、line number、変更種別で返す
- preserved behavior、決定済み方針、禁止事項は implementation_implementer が再読不要な粒度で要約する
- validation command と completion signal を明示する

## DO / DON'T

DO:
- preserved behavior を先に固定する
- 代表 path と差分だけを残す
- dependency direction を明示する
- fix_ingredients と distracting_context を分ける
- first_action と change_targets を path、symbol、line number 付きで返す
- requirements_policy_decisions に preserved behavior と out of scope を残す

DON'T:
- 追加の設計判断をしない
- 実 code を読まず handoff の文章を言い換えない
- 類似 context を required_reading に混ぜない
- refactor 方針や決定事項を required_reading に丸投げしない
- owned_scope 外の broad refactor を広げない
- active contract をこの skill に置かない

## Checklist

- [implementation-distill-refactor-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-distill-refactor/references/checklists/implementation-distill-refactor-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。
