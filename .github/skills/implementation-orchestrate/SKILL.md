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
- `execution_group` / `ready_wave` を ready wave として扱う判断
- 全 implementation handoff 完了後の scenario validation、suite-all、Sonar check
- `codex exec` による Codex review conductor 呼び出し
- Copilot 内 narrowing と例外的な Codex replan 判断

## 原則

- `implementation-scope` を唯一の実行正本にする
- 1 handoff を 1 RunSubagent 実行単位として扱う
- `execution_group` / `ready_wave` は必要な数だけある ready wave として扱い、同じ wave 内でも `parallelizable_with` に列挙された handoff だけを並列化する
- `first_action` がない handoff は実装開始せず、Copilot 内 narrowing の対象にする
- distiller は tester / implementer より先に lane-local context を作る
- tester を実装前に起動できるのは、承認済み scenario artifact を product test 化する handoff だけである
- unit test と原因未確定の regression test は、implementer 完了後に tester が追加または更新する
- scenario validation、suite-all、Sonar check は全 implementation handoff 完了後だけ実行する
- Codex review は final validation 後に `codex exec` で呼び出す
- オーケストレーター自身の validation 実行は scenario validation、suite-all、Sonar check だけに限定する
- docs 正本化を implementation lane に混ぜない

## 標準パターン

1. `approval_record` と `implementation-scope` の Ready Waves 表、handoff 見出し、owned_scope、depends_on、execution_group、ready_wave、parallelizable_with、parallel_blockers、first_action、validation command だけを確認する。
2. Ready Waves 表、`execution_group`、`depends_on` から、未完了 wave のうち実行可能な最小番号の ready wave を選ぶ。
3. 実行可能な handoff 1 件から `single_handoff_packet` を作る。
   - 同じ ready wave 内では、互いに `parallelizable_with` に列挙された handoff だけを RunSubagent 並列実行の候補にする。
   - `parallel_blockers` がある handoff は、blocker の理由が解消するまで単独または後続 wave として扱う。
4. `implementation-distiller` で `lane_context_packet` と `tester_context_packet` を作る。
5. scenario 先行条件を満たす場合だけ、`tester` に `single_handoff_packet`、`tester_context_packet`、test_subscope、owned_scope、test target だけを渡す。
6. `implementer` に `single_handoff_packet`、`lane_context_packet`、implementation_subscope、owned_scope、depends_on 解消結果を渡す。scenario 先行時だけ tester output も渡す。
7. unit test と regression test が必要な場合は、implementer 完了後に `tester` へ実装済み scope と tester_context_packet を渡す。
8. 全 implementation handoff 完了後、`python3 scripts/harness/run.py --suite scenario-gate` を実行し、task 固有の product scenario test command がある場合は同じ結果へ含める。
9. scenario validation が pass した場合だけ、`python3 scripts/harness/run.py --suite all` を実行する。
10. Sonar check を実行し、repo-local gate と Sonar server Quality Gate を混同しない。
11. `codex exec` で `review_conductor` を呼び出し、diff、scope、implementation result、validation result を渡す。
12. `copilot-work-reporting` を参照し、completion packet に `codex_review_result` と `copilot_work_report` を必ず含める。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は agent contract に従う。

## 実装パターン

### 通常パターン

新規実装または機能拡張の標準順である。
広すぎる handoff は、backend / frontend、test target、public boundary、change target、validation command のいずれか 1 軸に狭める。

1. `implementation-distiller` に渡し、handoff 1 件だけから context packet を作る。
2. `implementer` に渡し、同じ handoff 粒度で product code だけを実装する。
3. `tester` に渡し、実装済み責務を product test で証明する。
4. 全 implementation handoff 完了後、scenario validation、suite-all、Sonar check を実行する。
5. `codex exec` で Codex review conductor を呼び出す。

### 修正パターン

bug fix、regression、validation failure の標準順である。
原因不明のまま implementer へ渡さない。

1. 必要なら `investigator` に渡し、再現条件、error output、log、UI state を集める。
2. `implementation-distiller` に渡し、handoff 1 件だけから context packet を作る。
3. `implementer` に渡し、accepted fix scope だけを恒久修正する。
4. `tester` に渡し、原因と修正 seam が確定した regression test を追加または更新する。
5. 全 implementation handoff 完了後、scenario validation、suite-all、Sonar check、Codex review を順に行う。

### Refactor パターン

外部 behavior を変えない整理の標準順である。
不変条件が未定義なら実行しない。

1. `implementation-distiller` に渡し、不変条件を context packet に固定する。
2. `implementer` に渡し、owned_scope 内だけを refactor する。
3. `tester` に渡し、不変条件を証明する test を handoff 粒度で整える。
4. 全 implementation handoff 完了後、scenario validation、suite-all、Sonar check、Codex review を順に行う。

### UI / Mixed パターン

frontend / backend 横断や UI evidence が必要な時の標準順である。
backend 完了前に frontend handoff を先行しない。

1. backend 側 handoff を先に distiller、implementer、tester へ渡す。
2. frontend 側 handoff も distiller、implementer、tester へ渡す。
3. 必要なら `investigator` に渡し、`agent-browser` CLI で UI state、console、Wails binding evidence を集める。
4. 全 implementation handoff 完了後、scenario validation、suite-all、Sonar check、Codex review を順に行う。

### Scenario 先行パターン

承認済み scenario artifact を product test 化する時だけ使う。
原因未確定の regression test や unit test には使わない。

1. `implementation-distiller` に渡し、handoff 1 件だけから context packet を作る。
2. 承認済み scenario、public seam、観測点、期待 outcome が固定済みであることを確認する。
3. `tester` に渡し、scenario outcome を fail 前提の product test にする。
4. `implementer` に渡し、scenario test を満たす product code だけを実装する。
5. 実装後に必要な unit / regression test は tester へ戻す。

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
- scenario 先行条件を満たす時だけ tester を implementer より先に起動する
- unit test と原因未確定の regression test は実装後に tester へ渡す
- execution_group、parallelizable_with、parallel_blockers を見て ready wave を決める
- `first_action` を含む `single_handoff_packet` だけを tester / implementer へ渡す
- scenario validation、suite-all、Sonar check を全 implementation handoff 完了後に実行する
- `codex exec` の review payload に diff と validation result を含める
- 最後に必ず `copilot-work-reporting` で completion packet の報告材料を作る

DON'T:
- RunSubagent 以外で実装、test 追加、調査をしない
- `first_action` がない handoff を広い調査で補わない
- `parallelizable_with` に列挙されていない handoff を同じ wave という理由だけで並列実行しない
- final validation 前に scenario validation、suite-all、Sonar check を実行しない
- scenario validation failure を residual risk として close しない
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
