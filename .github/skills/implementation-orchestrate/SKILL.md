---
name: implementation-orchestrate
description: GitHub Copilot 側の実装入口知識 package。承認済み implementation-scope を実装、test、final validation、Codex review 呼び出しへ分配する判断基準を提供する。
---

# Implementation Orchestrate

## 目的

`implementation-orchestrate` は知識 package である。
GitHub Copilot 側の `implementation-orchestrate` agent が、承認済み `implementation-scope` を実装前整理、調査、実装、test、final validation、Codex review 呼び出しへ分配する時の判断基準を提供する。

実行権限、write scope、active contract、handoff は [implementation-orchestrate.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/implementation-orchestrate.agent.md) が持つ。

## いつ参照するか

- 承認済み implementation-scope の handoff を RunSubagent 実行順へ並べる時
- depends_on と並行可能な owned_scope を確認する時
- implementation-distiller、investigator、implementer、tester、final validation lane、Codex review conductor のどれへ渡すか判断する時

## 参照しない場合

- design bundle や implementation-scope を作る時
- docs 正本化をする時
- product code を直接変更する時

## 知識範囲

- handoff を RunSubagent 実行単位として扱う判断
- depends_on の解消順
- 全 implementation handoff 完了後の suite-all と Sonar check
- `codex exec` による Codex review conductor 呼び出し
- Copilot 内 narrowing と例外的な Codex replan 判断

## 原則

- `implementation-scope` を唯一の実行正本にする
- 1 handoff を 1 RunSubagent 実行単位として扱う
- distiller は tester / implementer より先に lane-local context を作る
- tester は全 implementation handoff で implementer より先に起動する
- suite-all と Sonar check は全 implementation handoff 完了後だけ実行する
- Codex review は final validation 後に `codex exec` で呼び出す
- オーケストレーター自身の validation 実行は suite-all と Sonar check だけに限定する
- docs 正本化を implementation lane に混ぜない

## 標準パターン

1. `approval_record` と `implementation-scope` の handoff 見出し、owned_scope、depends_on、validation command だけを確認する。
2. 実行可能な handoff 1 件から `single_handoff_packet` を作る。
3. `implementation-distiller` で `lane_context_packet` と `tester_context_packet` を作る。
4. `tester` に `single_handoff_packet`、`tester_context_packet`、test_subscope、owned_scope、test target だけを渡す。
5. `implementer` に `single_handoff_packet`、`lane_context_packet`、implementation_subscope、owned_scope、depends_on 解消結果、tester output だけを渡す。
6. 全 implementation handoff 完了後、`python3 scripts/harness/run.py --suite all` を実行する。
7. Sonar check を実行し、repo-local gate と Sonar server Quality Gate を混同しない。
8. `codex exec` で `review_conductor` を呼び出し、diff、scope、implementation result、validation result を渡す。
9. `copilot-work-reporting` を参照し、completion packet に `codex_review_result` と `copilot_work_report` を必ず含める。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は agent contract に従う。

## 実装パターン

### 通常パターン

新規実装または機能拡張の標準順である。
広すぎる handoff は、backend / frontend、test target、public boundary、change target、validation command のいずれか 1 軸に狭める。

1. `implementation-distiller` に渡し、handoff 1 件だけから context packet を作る。
2. `tester` に渡し、handoff 粒度で product test を追加または更新する。
3. `implementer` に渡し、同じ handoff 粒度で product code だけを実装する。
4. 全 implementation handoff 完了後、suite-all と Sonar check を実行する。
5. `codex exec` で Codex review conductor を呼び出す。

### 修正パターン

bug fix、regression、validation failure の標準順である。
原因不明のまま implementer へ渡さない。

1. 必要なら `investigator` に渡し、再現条件、error output、log、UI state を集める。
2. `implementation-distiller` に渡し、handoff 1 件だけから context packet を作る。
3. `tester` に渡し、failing test または regression test を handoff 粒度で追加する。
4. `implementer` に渡し、accepted fix scope だけを恒久修正する。
5. 全 implementation handoff 完了後、suite-all、Sonar check、Codex review を順に行う。

### Refactor パターン

外部 behavior を変えない整理の標準順である。
不変条件が未定義なら実行しない。

1. `implementation-distiller` に渡し、不変条件を context packet に固定する。
2. `tester` に渡し、不変条件を証明する test を handoff 粒度で整える。
3. `implementer` に渡し、owned_scope 内だけを refactor する。
4. 全 implementation handoff 完了後、suite-all、Sonar check、Codex review を順に行う。

### UI / Mixed パターン

frontend / backend 横断や UI evidence が必要な時の標準順である。
backend 完了前に frontend handoff を先行しない。

1. backend 側 handoff を先に distiller、tester、implementer へ渡す。
2. frontend 側 handoff も distiller、tester、implementer へ渡す。
3. 必要なら `investigator` に渡し、UI state、console、Wails binding evidence を集める。
4. 全 implementation handoff 完了後、suite-all、Sonar check、Codex review を順に行う。

## Codex Review 呼び出し

`codex exec` は Copilot の実装完了後に呼び出す。
Copilot の課金形態では待機が request cost にならないため、実装完了後の同期 review として扱う。

渡す payload は次を含める。

- `implementation_scope_path`
- `approval_record`
- `implementation_result`
- `diff_summary` または review 対象 diff
- `final_validation_result`
- `touched_files`

## DO / DON'T

DO:
- distiller を tester / implementer より先に起動する
- tester を implementer より先に起動する
- suite-all と Sonar check を全 implementation handoff 完了後に実行する
- `codex exec` の review payload に diff と validation result を含める
- 最後に必ず `copilot-work-reporting` で completion packet の報告材料を作る

DON'T:
- RunSubagent 以外で実装、test 追加、調査をしない
- final validation 前に suite-all や Sonar check を実行しない
- repo-local Sonar issue gate と Sonar server Quality Gate を混同しない
- docs、`.codex`、`.github/skills`、`.github/agents` を変更しない

## 参照パターン

- [orchestration-patterns.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-orchestrate/references/patterns/orchestration-patterns.md) を参照する。
- coverage は repo の `MINIMUM_COVERAGE = 70.0` を正本にする。
- `sonar_gate_result` は互換 field 名として残る場合があるが、意味は repo-local Sonar issue gate であり Sonar サーバ側 Quality Gate ではない。
- report: [copilot-work-reporting](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/copilot-work-reporting/SKILL.md) を参照する。

## Checklist

- [implementation-orchestrate-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-orchestrate/references/checklists/implementation-orchestrate-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。

## Agent が持つもの

- active contract: [implementation-orchestrate.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/references/implementation-orchestrate/contracts/implementation-orchestrate.contract.json)
- permissions: [permissions.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/references/implementation-orchestrate/permissions.json)

## Maintenance

- output obligation を skill 本体へ戻さない。
- mode / variant contract を skill 配下の active 正本にしない。
