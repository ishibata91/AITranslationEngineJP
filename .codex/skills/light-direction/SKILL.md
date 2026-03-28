---
name: light-direction
description: AITranslationEngineJp 専用。仕様判断と受け入れ条件が固定済みで、blocking unknown がない修正に対して、Architect が short plan を作って coder 実装へ直接 handoff したいときに使う。
---

# Light Direction

この skill は、軽量フローの入口です。
Architect が短い plan を固め、Coder に直接実装させ、`workflow-gate` で閉じます。

## 使う条件

- 仕様判断と受け入れ条件が固定済み
- blocking unknown がない
- 単一責務
- 短い plan で判断を固定できる

## Required Workflow

1. 依頼が軽量条件を満たすか確認する。
2. blocking unknown が見つかったら、実装開始前に heavy へ reroute する。
3. 以下を含む short plan を固める。
   - Request Summary
   - Decision Basis
   - Why Light Flow Applies
   - Short Plan
   - Checks
   - Required Evidence
   - Reroute Trigger
   - Docs Sync
   - Record Updates
4. `coder` に `light-work` を使わせて実装させる。
5. 実装後に Architect 自身が `workflow-gate` を標準 gate として使う。
6. `light-review` は gate では判定できない設計論点が残る時だけ補助 checklist として使う。
7. gate 結果で accept / reroute を決める。

## 禁止

- 仕様変更を含む依頼を軽量扱いしない
- blocking unknown がある依頼を軽量扱いしない
- plan を省略しない
- workflow-gate を飛ばさない
