---
name: light-direction
description: AITranslationEngineJp 専用。仕様変更のない低リスク修正に対して、Architect が short plan を作って coder 実装へ直接 handoff したいときに使う。
---

# Light Direction

この skill は、軽量フローの入口です。
Architect が短い plan を固め、Coder に直接実装させ、自身の review で閉じます。

## 使う条件

- 仕様変更なし
- 低リスク
- 単一責務
- 短い plan で判断を固定できる

## Required Workflow

1. 依頼が軽量条件を満たすか確認する。
2. 以下を含む short plan を固める。
   - Request Summary
   - Why Light Flow Applies
   - Short Plan
   - Checks
   - Record Updates
3. `coder` に `light-work` を使わせて実装させる。
4. 実装後に Architect 自身が `light-review` を review checklist として使う。
5. review 結果で accept / reroute を決める。

## 禁止

- 仕様変更を含む依頼を軽量扱いしない
- plan を省略しない
- review を飛ばさない
