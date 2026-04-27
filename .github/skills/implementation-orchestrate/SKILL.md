---
name: implementation-orchestrate
description: GitHub Copilot 側の実装入口知識 package。承認済み implementation-scope を実装、test、final validation、Codex review request 作成へ分配する判断基準を提供する。
---

# Implementation Orchestrate

## 目的

`implementation-orchestrate` は知識 package である。
GitHub Copilot 側の `implementation-orchestrate` agent が、承認済み `implementation-scope` を実装前整理、調査、実装、test、final validation、Codex review request 作成へ分配する時の判断基準を提供する。

実行権限、write scope、active contract、handoff は [implementation-orchestrate.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/implementation-orchestrate.agent.md) が持つ。

## いつ参照するか

- 承認済み implementation-scope の handoff を RunSubagent 実行順へ並べる時
- depends_on と並行可能な owned_scope を確認する時
- implementation-distiller、investigator、implementer、tester、final validation lane、人間実行用 Codex review request 作成のどれへ進めるか判断する時

## 参照しない場合

- design bundle や implementation-scope を作る時
- docs 正本化をする時
- product code を直接変更する時

## 知識範囲

- handoff を RunSubagent 実行単位として扱う判断
- contract freeze を downstream 実装開始条件として扱う判断
- depends_on の解消順
- `execution_group` / `ready_wave` を ready wave として扱う判断
- 全 implementation handoff 完了後の scenario validation、suite-all、Sonar check
- 人間実行用 `codex exec` request payload と command の作成
- Copilot 内 narrowing と例外的な Codex replan 判断

## 原則

- 人間指示ごとに `implementation-orchestrate` skill、permissions、contract、承認済み `implementation-scope` を読みなおす
- `implementation-scope` を唯一の実行正本にする
- 1 handoff を 1 RunSubagent 実行単位として扱う
- `contract_freeze.status: required` の handoff は、対応する `completion_signal` が揃うまで downstream handoff を開始しない
- `execution_group` / `ready_wave` は必要な数だけある ready wave として扱い、同じ wave 内でも `parallelizable_with` に列挙された handoff だけを並列化する
- `first_action` がない handoff は実装開始せず、Copilot 内 narrowing の対象にする
- distiller は tester / implementer より先に lane-local context を作る
- tester を実装前に起動できるのは、承認済み `APIテスト` を product test 化する handoff だけである
- unit test と原因未確定の regression test は、implementer 完了後に tester が追加または更新する
- scenario validation、suite-all、Sonar check は全 implementation handoff 完了後だけ実行する
- Codex review は Copilot が直接呼び出さず、final validation 後に人間実行用 `codex exec` request payload と command を返す
- オーケストレーター自身の validation 実行は scenario validation、suite-all、Sonar check だけに限定する
- docs 正本化を implementation lane に混ぜない

## 標準パターン

1. 人間指示を受けたら、`implementation-orchestrate` skill、permissions、contract、承認済み `implementation-scope` を読みなおし、approved scope と lane 境界を超えていないか判断する。
2. `approval_record` と `implementation-scope` の Ready Waves 表、handoff 見出し、contract_freeze、owned_scope、depends_on、execution_group、ready_wave、parallelizable_with、parallel_blockers、first_action、validation command だけを確認する。
3. Ready Waves 表、`execution_group`、`depends_on` から、未完了 wave のうち実行可能な最小番号の ready wave を選ぶ。
4. 実行可能な handoff 1 件から `single_handoff_packet` を作る。
   - 同じ ready wave 内では、互いに `parallelizable_with` に列挙された handoff だけを RunSubagent 並列実行の候補にする。
   - `parallel_blockers` がある handoff は、blocker の理由が解消するまで単独または後続 wave として扱う。
   - downstream handoff は、依存先 handoff の `contract_freeze.status: done` と対応する `completion_signal` が揃うまで着手しない。
5. `implementation-distiller` で `lane_context_packet` と `tester_context_packet` を作る。
6. `APIテスト` 先行条件を満たす場合だけ、`tester` に `single_handoff_packet`、`tester_context_packet`、test_subscope、owned_scope、test target だけを渡す。
7. `implementer` に `single_handoff_packet`、`lane_context_packet`、implementation_subscope、owned_scope、depends_on 解消結果を渡す。`APIテスト` 先行時だけ tester output も渡す。
8. unit test と regression test が必要な場合は、implementer 完了後に `tester` へ実装済み scope と tester_context_packet を渡す。
9. 全 implementation handoff 完了後、`python3 scripts/harness/run.py --suite scenario-gate` を実行し、task 固有の product scenario test command がある場合は同じ結果へ含める。
10. scenario validation が pass した場合だけ、`python3 scripts/harness/run.py --suite all` を実行する。
11. Sonar check を実行し、repo-local gate と Sonar server Quality Gate を混同しない。
12. 人間実行用 Codex review request payload と `codex exec` command を completion packet に含める。
13. 人間から `codex_review_result` が戻された場合だけ、`codex_review_result.copilot_action` に従い、close、residual report、修正、validation 再実行、Codex review request 再作成のいずれかへ分岐する。

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
5. 人間実行用 Codex review request payload と `codex exec` command を返す。

### 修正パターン

bug fix、regression、validation failure の標準順である。
原因不明のまま implementer へ渡さない。

1. 必要なら `investigator` に渡し、再現条件、error output、log、UI state を集める。
2. `implementation-distiller` に渡し、handoff 1 件だけから context packet を作る。
3. `implementer` に渡し、accepted fix scope だけを恒久修正する。
4. `tester` に渡し、原因と修正 seam が確定した regression test を追加または更新する。
5. 全 implementation handoff 完了後、scenario validation、suite-all、Sonar check、Codex review request 作成を順に行う。

### Refactor パターン

外部 behavior を変えない整理の標準順である。
不変条件が未定義なら実行しない。

1. `implementation-distiller` に渡し、不変条件を context packet に固定する。
2. `implementer` に渡し、owned_scope 内だけを refactor する。
3. `tester` に渡し、不変条件を証明する test を handoff 粒度で整える。
4. 全 implementation handoff 完了後、scenario validation、suite-all、Sonar check、Codex review request 作成を順に行う。

### UI / Mixed パターン

frontend / backend 横断や UI evidence が必要な時の標準順である。
backend 完了前に frontend handoff を先行しない。

1. backend 側 handoff を先に distiller、implementer、tester へ渡す。
2. frontend 側 handoff も distiller、implementer、tester へ渡す。
3. 必要なら `investigator` に渡し、`agent-browser` CLI で UI state、console、Wails binding evidence を集める。
4. 全 implementation handoff 完了後、scenario validation、suite-all、Sonar check、Codex review request 作成を順に行う。

### APIテスト先行パターン

承認済み `APIテスト` を product test 化する時だけ使う。
原因未確定の regression test や unit test、`UI人間操作E2E` には使わない。

1. `implementation-distiller` に渡し、handoff 1 件だけから context packet を作る。
2. 承認済み受け入れ条件、public seam、入力開始点、主要観測点、期待 outcome が固定済みであることを確認する。
3. `tester` に渡し、API contract outcome を fail 前提の product test にする。
4. `implementer` に渡し、API test を満たす product code だけを実装する。
5. 実装後に必要な unit / regression test は tester へ戻す。

## Codex Review Request

Copilot は `codex exec` を直接呼び出さない。
Copilot は実装完了後に、人間が Codex 側で実行する review request payload と command を completion packet に含める。

渡す payload は次を含める。

- `implementation_scope_path`
- `approval_record`
- `implementation_result`
- `diff_summary` または review 対象 diff
- `final_validation_result`
- `touched_files`

completion packet には、人間がそのまま使える `codex_review_request_payload` と `human_codex_exec_command` を含める。
`human_codex_exec_command` は、payload を標準入力または file 経由で `review_conductor` に渡す形にする。

## Codex Review 受け取り

人間から `codex_review_result` が戻された場合だけ、Copilot は `copilot_action` で受け取る。
Copilot は `decision_basis` を再解釈せず、次の分岐だけを行う。

- `close`: completion packet に `codex_review_result` を転記して終了する
- `report_residual`: `priority_overrides` と `residual_risks` を completion packet に残して終了する
- `fix`: `copilot_patch_scope` 内だけを修正し、final validation と Codex review request payload を再作成する
- `rerun_validation`: 指定された不足 validation だけを再実行し、Codex review request payload を再作成する
- `rerun_codex_review`: 不足 payload を補い、product code を変更せず Codex review request payload だけを再作成する

## DO / DON'T

DO:
- 人間指示を受けたら skill、permissions、contract、承認済み `implementation-scope` を読みなおす
- 人間指示を skill / contract より上位の境界変更として扱わない
- distiller を tester / implementer より先に起動する
- `APIテスト` 先行条件を満たす時だけ tester を implementer より先に起動する
- unit test と原因未確定の regression test は実装後に tester へ渡す
- execution_group、parallelizable_with、parallel_blockers を見て ready wave を決める
- contract freeze 完了前の downstream handoff を開始しない
- `first_action` を含む `single_handoff_packet` だけを tester / implementer へ渡す
- scenario validation、suite-all、Sonar check を全 implementation handoff 完了後に実行する
- 人間実行用 Codex review request payload に diff と validation result を含める
- 人間から `codex_review_result` が戻された時だけ、`codex_review_result.copilot_action` に従って受け取り分岐を固定する
- `UI人間操作E2E` は final validation lane でだけ証明する

DON'T:
- 人間指示を理由に docs、`.codex`、`.github/skills`、`.github/agents` を変更しない
- RunSubagent 以外で実装、test 追加、調査をしない
- `first_action` がない handoff を広い調査で補わない
- `contract_freeze.status: required` を単なる notes として無視しない
- `parallelizable_with` に列挙されていない handoff を同じ wave という理由だけで並列実行しない
- final validation 前に scenario validation、suite-all、Sonar check を実行しない
- scenario validation failure を residual risk として close しない
- repo-local Sonar issue gate と Sonar server Quality Gate を混同しない
- Copilot が `codex exec` を直接呼び出さない
- `rerun_codex_review` で product code を変更しない
- `fix` で `copilot_patch_scope` の外を変更しない
- docs、`.codex`、`.github/skills`、`.github/agents` を変更しない

## 参照パターン

- [orchestration-patterns.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-orchestrate/references/patterns/orchestration-patterns.md) を参照する。
- coverage は repo の `MINIMUM_COVERAGE = 70.0` を正本にする。
- `sonar_gate_result` は互換 field 名として残る場合があるが、意味は repo-local Sonar issue gate であり Sonar サーバ側 Quality Gate ではない。
## Checklist

- [implementation-orchestrate-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-orchestrate/references/checklists/implementation-orchestrate-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。

## Agent が持つもの

- active contract: [implementation-orchestrate.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/references/implementation-orchestrate/contracts/implementation-orchestrate.contract.json)
- permissions: [permissions.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/references/implementation-orchestrate/permissions.json)

## Maintenance

- output obligation を skill 本体へ戻さない。
- mode / variant contract を skill 配下の active 正本にしない。
