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
- coverage 70%、Sonar gate、harness suite の closeout 判断
- subagent 戻り値だけから completion packet を作る判断
- design 不足を実装側で補わない reroute 判断

## 原則

- `implementation-scope` を唯一の実行正本にする
- 1 handoff を 1 RunSubagent 実行単位として扱う
- オーケストレーター自身は file read / search / edit / validation 実行をしない
- design 不足は実装せず戻す
- docs 正本化を implementation lane に混ぜない

## 標準パターン

1. `approval_record` と `implementation-scope` の handoff 見出し、owned_scope、depends_on、validation command だけを確認する。
2. `depends_on` を解消し、実行可能な handoff を 1 件だけ選ぶ。
3. context 整理、調査、実装、test、review の担当 agent を 1 つ選ぶ。
4. RunSubagent に active contract、handoff 1 件、禁止事項、期待 output を渡す。
5. validation、coverage、Sonar、harness evidence は subagent の戻り値だけから集約する。
6. subagent result に docs 変更がないこと、または reroute 理由があることを明記する。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は agent contract に従う。

## 実装パターン

### 通常パターン

新規実装または機能拡張の標準順である。
1 handoff の scope が広い場合は、実行せず `propose-plans` へ戻す。

1. `implementation-distiller` に渡し、implementation context packet を作る。
2. `tester` に渡し、handoff 資料のスコープ粒度で product test を追加または更新する。
3. `implementer` に渡し、同じ handoff 粒度で product code を実装する。
4. `reviewer` に渡し、implementation review または UI check を行う。
5. coverage、Sonar、harness の gate 結果を completion packet に集約する。

### 修正パターン

bug fix、regression、validation failure の標準順である。
原因不明のまま implementer へ渡さない。

1. `investigator` に渡し、再現条件、error output、log、UI state を集める。
2. 必要なら `implementation-distiller` に渡し、fix scope と required reading を圧縮する。
3. `tester` に渡し、再現を証明する failing test または regression test を handoff 粒度で追加する。
4. `implementer` に渡し、accepted fix scope だけを恒久修正する。
5. `reviewer` に渡し、再現が閉じたことと gate 未達がないことを確認する。

### Refactor パターン

外部 behavior を変えない整理の標準順である。
不変条件が未定義なら実行しない。

1. `implementation-distiller` に渡し、preserved behavior、dependency boundary、不変条件を固定する。
2. `tester` に渡し、不変条件を証明する既存 test または補強 test を handoff 粒度で整える。
3. `implementer` に渡し、owned_scope 内だけを refactor する。
4. `reviewer` に渡し、behavior drift、broad refactor、coverage / Sonar / harness gate を確認する。

### UI / Mixed パターン

frontend / backend 横断や UI evidence が必要な時の標準順である。
backend 完了前に frontend handoff を先行しない。

1. `implementation-distiller` に渡し、backend、frontend、Wails boundary、UI evidence point を分ける。
2. backend 側 handoff があれば、先に `tester` と `implementer` へ渡す。
3. frontend 側 handoff を `tester` と `implementer` へ渡し、visible state と test seam をそろえる。
4. 必要なら `investigator` に渡し、UI state、console、Wails binding evidence を集める。
5. `reviewer` に渡し、UI check と implementation review を scope 内で行う。

## DO / DON'T

DO:
- handoff 見出しを RunSubagent 単位にする
- depends_on を守る
- subagent 戻り値だけを集約する
- design 不足は `propose-plans` へ戻す

DON'T:
- RunSubagent 以外の tool を使う
- `implementation-scope` を書き換えない
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
