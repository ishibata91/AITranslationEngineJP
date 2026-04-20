---
name: implementation-orchestrate
description: GitHub Copilot 側の実装入口知識 package。承認済み implementation-scope を実装 lane に分配する判断基準を提供する。
---

# Implementation Orchestrate

## 目的

`implementation-orchestrate` は知識 package である。
GitHub Copilot 側の `implementation-orchestrate` agent が、承認済み `implementation-scope` を実装前整理、調査、実装、test、review へ分配する時の判断基準を提供する。

実行権限、write scope、active contract、handoff は [implementation-orchestrate.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/implementation-orchestrate.agent.md) が持つ。

## いつ参照するか

- 承認済み implementation-scope の handoff を RunSubagent 実行順へ並べる時
- depends_on と並行可能な owned_scope を確認する時
- implementation-distiller、investigator、implementer、tester、reviewer のどれへ渡すか判断する時

## 参照しない場合

- design bundle や implementation-scope を作る時
- docs 正本化をする時
- product code を直接変更する時

## 知識範囲

- handoff を RunSubagent 実行単位として扱う判断
- depends_on の解消順
- lane-local result と final validation lane の closeout 判断
- subagent 戻り値だけから completion packet を作る判断
- design 不足を実装側で補わない reroute 判断

## 原則

- `implementation-scope` を唯一の実行正本にする
- 1 handoff を 1 RunSubagent 実行単位として扱う
- RunSubagent に渡す source scope は `single_handoff_packet` 1 件と、その distill 結果に限定する
- implementation-distiller は tester / implementer より先に lane-local context を作る
- tester は全 implementation handoff で implementer より先に起動する
- implementer へ full `implementation-scope`、active work plan 全文、source artifacts、後続 handoff を渡さない
- オーケストレーター自身は file read / search / edit / validation 実行をしない
- design 不足は実装せず戻す
- docs 正本化を implementation lane に混ぜない

## 標準パターン

1. `approval_record` と `implementation-scope` の handoff 見出し、owned_scope、depends_on、validation command だけを確認する。
2. `depends_on` を解消し、実行可能な handoff を 1 件だけ選ぶ。
3. 選んだ handoff から `single_handoff_packet` を作る。
4. `implementation-distiller` に `single_handoff_packet` だけを渡し、`lane_context_packet` を作る。
5. `lane_context_packet` に fix_ingredients、distracting_context、first_action、change_targets、requirements_policy_decisions、symbol / line number 付き related_code_pointers があることを確認する。
6. first_action が 1 clause に固定され、推測 method が fact 化されず、existing_patterns と validation_entry の探索理由があることを確認する。
7. 不足していれば tester / implementer へ渡さず reroute reason にする。
8. `tester` に `single_handoff_packet`、`lane_context_packet`、owned_scope、test target だけを渡す。
9. `implementer` に `single_handoff_packet`、`lane_context_packet`、owned_scope、depends_on 解消結果、tester output、禁止事項だけを渡す。
10. `reviewer` に lane-local の実装結果と test result だけを渡す。
11. validation、coverage、Sonar、harness evidence は subagent の戻り値だけから集約する。
12. subagent result に docs 変更がないこと、または reroute 理由があることを明記する。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は agent contract に従う。

## 実装パターン

### 通常パターン

新規実装または機能拡張の標準順である。
1 handoff の scope が広い場合は、実行せず `propose-plans` へ戻す。

1. `implementation-distiller` に渡し、handoff 1 件だけから `lane_context_packet` を作る。
2. `tester` に渡し、handoff 資料のスコープ粒度で product test を追加または更新する。
3. `implementer` に渡し、同じ handoff 粒度で product code だけを実装する。
4. `reviewer` に渡し、implementation review または UI check を行う。
5. lane-local result を completion packet に集約する。

### 修正パターン

bug fix、regression、validation failure の標準順である。
原因不明のまま implementer へ渡さない。

1. `investigator` に渡し、再現条件、error output、log、UI state を集める。
2. `implementation-distiller` に渡し、調査結果を足さず handoff 1 件だけから `lane_context_packet` を作る。
3. `tester` に渡し、再現を証明する failing test または regression test を handoff 粒度で追加する。
4. `implementer` に渡し、accepted fix scope だけを恒久修正する。
5. `reviewer` に渡し、再現が閉じたことと lane-local validation の未達がないことを確認する。

### Refactor パターン

外部 behavior を変えない整理の標準順である。
不変条件が未定義なら実行しない。

1. `implementation-distiller` に渡し、handoff 1 件だけから不変条件の `lane_context_packet` を作る。
2. `tester` に渡し、不変条件を証明する既存 test または補強 test を handoff 粒度で整える。
3. `implementer` に渡し、owned_scope 内だけを refactor する。
4. `reviewer` に渡し、behavior drift と broad refactor がないことを確認する。

### UI / Mixed パターン

frontend / backend 横断や UI evidence が必要な時の標準順である。
backend 完了前に frontend handoff を先行しない。

1. backend 側 handoff があれば、先に `implementation-distiller`、`tester`、`implementer` へ渡す。
2. frontend 側 handoff も `implementation-distiller`、`tester`、`implementer` へ渡し、visible state と test seam をそろえる。
3. 必要なら `investigator` に渡し、UI state、console、Wails binding evidence を集める。
4. `reviewer` に渡し、UI check と implementation review を scope 内で行う。

### Distill パターン

`implementation-distiller` は default path で必ず使う。
ただし distiller に渡す input は `single_handoff_packet` 1 件だけに限定する。
distiller は full implementation-scope、active work plan 全文、source artifacts、後続 handoff を読まず、tester / implementer が使う `lane_context_packet` だけを返す。
`lane_context_packet` は fix_ingredients、distracting_context、first_action、change_targets、requirements_policy_decisions、required_reading、symbol / line number 付き related_code_pointers を含める。
handoff 1 件だけでは `lane_context_packet` を作れない設計不足は `propose-plans` へ戻す。

## DO / DON'T

DO:
- handoff 見出しを RunSubagent 単位にする
- depends_on を守る
- distiller に `single_handoff_packet` 1 件だけを渡す
- distiller output が patch 生成に必要な fix_ingredients を持つことを確認する
- distiller output が distracting_context を required_reading から分離していることを確認する
- distiller output が具体的な first_action と code pointer を持つことを確認する
- distiller output の first_action が 1 clause に固定されていることを確認する
- distiller output が推測 method を fact にしていないことを確認する
- distiller output の existing_patterns と validation_entry が探索理由を持つことを確認する
- distiller output が要件、実装方針、決定事項を要約していることを確認する
- tester を implementer より先に起動する
- implementer へ `single_handoff_packet`、`lane_context_packet`、tester output だけを渡す
- subagent 戻り値だけを集約する
- design 不足は `propose-plans` へ戻す

DON'T:
- RunSubagent 以外の tool を使う
- `implementation-scope` を書き換えない
- RunSubagent に full `implementation-scope`、active work plan 全文、source artifacts、後続 handoff を渡さない
- distiller に full `implementation-scope` や source artifacts を渡さない
- handoff 文面の言い換えだけの distiller output を implementer に渡さない
- fix_ingredients がない distiller output を implementer に渡さない
- first_action が partial または複数 clause の distiller output を implementer に渡さない
- 推測 method を fact にした distiller output を implementer に渡さない
- 類似 context を required_reading に混ぜた distiller output を implementer に渡さない
- 要件、実装方針、決定事項を required_reading に丸投げした distiller output を implementer に渡さない
- implementer に product test の新規作成、更新、fixture 調整を依頼しない
- file read / search / edit / validation 実行をしない
- docs、`.codex`、`.github/skills`、`.github/agents` を変更しない
- design review を行わない

## 参照パターン

- [orchestration-patterns.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-orchestrate/references/patterns/orchestration-patterns.md) を参照する。
- 対象は handoff 分割、depends_on 解消、validation failure の最小 reroute、closeout gate の判断である。
- coverage は repo の `MINIMUM_COVERAGE = 70.0` を正本にする。

## Checklist

- [implementation-orchestrate-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-orchestrate/references/checklists/implementation-orchestrate-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。

## Agent が持つもの

- active contract: [implementation-orchestrate.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/references/implementation-orchestrate/contracts/implementation-orchestrate.contract.json)
- permissions: [permissions.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/references/implementation-orchestrate/permissions.json)

## Maintenance

- output obligation を skill 本体へ戻さない。
- mode / variant contract を skill 配下の active 正本にしない。
- 実装責務が分かれる場合は agent 側の handoff で分ける。
