---
name: implementation_orchestrator
description: Codex implementation lane 側の execution artifact graph 実行知識 package。implement_lane / fix_lane の承認済み artifact を実装、test、final validation、観点別 review、remediation 分岐へ分配する判断基準を提供する。
---

# Implementation Orchestrate

## 目的

`implementation_orchestrator` は知識 package である。
Codex implementation lane 側の `implementation_orchestrator` agent が、`implement_lane` / `fix_lane` の承認済み execution artifact graph を実装前整理、調査、実装、test、final validation、観点別 review、remediation 分岐へ分配する時の判断基準を提供する。

実行権限、write scope、active contract、handoff は [implementation_orchestrator.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_orchestrator.toml) が持つ。

## いつ参照するか

- 承認済み execution artifact graph の handoff を spawn_agent 実行順へ並べる時
- depends_on と並行可能な owned_scope を確認する時
- implementation_distiller、implementation_investigator、implementation_implementer、implementation_tester、final validation lane、観点別 review agent、remediation 分岐のどれへ進めるか判断する時

## 参照しない場合

- design bundle、fix scope、implementation-scope など execution artifact を作る時
- docs 正本化をする時
- product code を直接変更する時

## 知識範囲

- handoff を spawn_agent 実行単位として扱う判断
- contract freeze を downstream 実装開始条件として扱う判断
- depends_on の解消順
- `execution_group` / `ready_wave` を ready wave として扱う判断
- 全 implementation handoff 完了後の scenario validation、suite-all、Sonar check
- 観点別 review agent の並列 spawn、lossless aggregation、`implementation_action` の生成
- Codex implementation lane 内 narrowing と例外的な Codex replan 判断

## 原則

- 人間指示ごとに `implementation-orchestrate` skill、permissions、contract、承認済み `implementation-scope` を読みなおす
- `implement_lane` / `fix_lane` の承認済み execution artifact graph を唯一の実行正本にする
- 1 handoff を 1 spawn_agent 実行単位として扱う
- `contract_freeze.status: required` の handoff は、対応する `completion_signal` が揃うまで downstream handoff を開始しない
- `execution_group` / `ready_wave` は必要な数だけある ready wave として扱い、同じ wave 内でも `parallelizable_with` に列挙された handoff だけを並列化する
- `first_action` がない handoff は実装開始せず、Codex implementation lane 内 narrowing の対象にする
- distiller は implementation_tester / implementation_implementer より先に lane-local context を作る
- implementation_tester を実装前に起動できるのは、承認済み `APIテスト` を product test 化する handoff だけである
- unit test と原因未確定の regression test は、implementation_implementer 完了後に implementation_tester が追加または更新する
- scenario validation、suite-all、Sonar check は全 implementation handoff 完了後だけ実行する
- final validation 後に review input が揃う場合だけ、4 観点 review agent を context 継承なしで並列 spawn する
- review input が不足する場合は reviewer を起動せず、`rerun_validation` または `rerun_codex_review` を返す
- review aggregation は `score > 0.85`、`trust_boundary > behavior > contract > state_invariant`、trust boundary hard gate を使う
- 観点別 raw result は `reviewer_result_bundle` に保持し、派生結果は `aggregation_trace` から辿れるようにする
- オーケストレーター自身の validation 実行は scenario validation、suite-all、Sonar check だけに限定する
- docs 正本化を implementation lane に混ぜない
- closeout、停止、reroute 時は work report と benchmark evidence を必ず completion packet に含める

## 標準パターン

1. 人間指示を受けたら、`implementation-orchestrate` skill、permissions、contract、承認済み `implementation-scope` を読みなおし、approved scope と lane 境界を超えていないか判断する。
2. `approval_record` と承認済み execution artifact graph の Ready Waves 表、handoff 見出し、contract_freeze、owned_scope、depends_on、execution_group、ready_wave、parallelizable_with、parallel_blockers、first_action、validation command だけを確認する。
3. Ready Waves 表、`execution_group`、`depends_on` から、未完了 wave のうち実行可能な最小番号の ready wave を選ぶ。
4. 実行可能な handoff 1 件から `single_handoff_packet` を作る。
   - 同じ ready wave 内では、互いに `parallelizable_with` に列挙された handoff だけを spawn_agent 並列実行の候補にする。
   - `parallel_blockers` がある handoff は、blocker の理由が解消するまで単独または後続 wave として扱う。
   - downstream handoff は、依存先 handoff の `contract_freeze.status: done` と対応する `completion_signal` が揃うまで着手しない。
5. `implementation_distiller` で `lane_context_packet` と `implementation_tester_context_packet` を作る。
6. `APIテスト` 先行条件を満たす場合だけ、`implementation_tester` に `single_handoff_packet`、`implementation_tester_context_packet`、test_subscope、owned_scope、test target だけを渡す。
7. `implementation_implementer` に `single_handoff_packet`、`lane_context_packet`、implementation_subscope、owned_scope、depends_on 解消結果を渡す。`APIテスト` 先行時だけ implementation_tester output も渡す。
8. unit test と regression test が必要な場合は、implementation_implementer 完了後に `implementation_tester` へ実装済み scope と implementation_tester_context_packet を渡す。
9. 全 implementation handoff 完了後、`python3 scripts/harness/run.py --suite scenario-gate` を実行し、task 固有の product scenario test command がある場合は同じ結果へ含める。
10. scenario validation が pass した場合だけ、`python3 scripts/harness/run.py --suite all` を実行する。
11. Sonar check を実行し、repo-local gate と Sonar server Quality Gate を混同しない。
12. diff、implementation-scope、implementation result、final validation result が揃っているか確認する。
13. 不足がある場合は reviewer を起動せず、diff / scope / implementation result 不足なら `rerun_codex_review`、validation 不足なら `rerun_validation` を返す。
14. review 可能な場合は `review_behavior`、`review_contract`、`review_trust_boundary`、`review_state_invariant` を context 継承なしで並列 spawn する。
15. 各 group raw result を `reviewer_result_bundle` に保持し、`aggregation_trace`、`remediation_aggregation`、`remediation_handoff`、`implementation_action` を作る。
16. `implementation_action` に従い、close、residual report、修正者判断、validation 再実行、review input 再構築のいずれかへ分岐する。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は agent contract に従う。

## 実装パターン

### 通常パターン

新規実装または機能拡張の標準順である。
広すぎる handoff は、backend / frontend、test target、public boundary、change target、validation command のいずれか 1 軸に狭める。

1. `implementation_distiller` に渡し、handoff 1 件だけから context packet を作る。
2. `implementation_implementer` に渡し、同じ handoff 粒度で product code だけを実装する。
3. `implementation_tester` に渡し、実装済み責務を product test で証明する。
4. 全 implementation handoff 完了後、scenario validation、suite-all、Sonar check を実行する。
5. 観点別 review agent を並列起動し、integrated review result と `implementation_action` を返す。

### 修正パターン

bug fix、regression、validation failure の標準順である。
原因不明のまま implementation_implementer へ渡さない。

1. 必要なら `implementation_investigator` に渡し、再現条件、error output、log、UI state を集める。
2. `implementation_distiller` に渡し、handoff 1 件だけから context packet を作る。
3. `implementation_implementer` に渡し、accepted fix scope だけを恒久修正する。
4. `implementation_tester` に渡し、原因と修正 seam が確定した regression test を追加または更新する。
5. 全 implementation handoff 完了後、scenario validation、suite-all、Sonar check、integrated review を順に行う。

### Refactor / Exploration / UX Placeholder

`refactor_lane`、`exploration_test_lane`、`ux_refactor_lane` は現時点では placeholder である。
必須 artifact、gate、actor、completion signal は未定義のため、implementation_orchestrator はこれらを実行しない。

これらの lane が指定された場合は、承認済み execution artifact graph 不足として停止し、`requires_codex_replan` に定義不足を返す。

### UI / Mixed パターン

frontend / backend 横断や UI evidence が必要な時の標準順である。
backend 完了前に frontend handoff を先行しない。

1. backend 側 handoff を先に distiller、implementation_implementer、implementation_tester へ渡す。
2. frontend 側 handoff も distiller、implementation_implementer、implementation_tester へ渡す。
3. 必要なら `implementation_investigator` に渡し、`agent-browser` CLI で UI state、console、Wails binding evidence を集める。
4. 全 implementation handoff 完了後、scenario validation、suite-all、Sonar check、integrated review を順に行う。

### APIテスト先行パターン

承認済み `APIテスト` を product test 化する時だけ使う。
原因未確定の regression test や unit test、`UI人間操作E2E` には使わない。

1. `implementation_distiller` に渡し、handoff 1 件だけから context packet を作る。
2. 承認済み受け入れ条件、public seam、入力開始点、主要観測点、期待 outcome が固定済みであることを確認する。
3. `implementation_tester` に渡し、API contract outcome を fail 前提の product test にする。
4. `implementation_implementer` に渡し、API test を満たす product code だけを実装する。
5. 実装後に必要な unit / regression test は implementation_tester へ戻す。

## Integrated Codex Review

Codex implementation lane は final validation 後に観点別 review agent を直接並列 spawn する。
review input は次を含める。

- `implementation_scope_path`
- `approval_record`
- `implementation_result`
- `diff_summary` または review 対象 diff
- `final_validation_result`
- `touched_files`

diff、scope、implementation result が不足する場合は、観点別 review agent を起動せず `implementation_action: rerun_codex_review` を返す。
final validation result が不足する場合は、観点別 review agent を起動せず `implementation_action: rerun_validation` を返す。

review 可能な場合は次を実行する。

1. `review_behavior`、`review_contract`、`review_trust_boundary`、`review_state_invariant` を context 継承なしで並列 spawn する。
2. 各 agent には同じ diff、scope、implementation result、final validation result を渡す。
3. 各 agent の score、confidence、findings、evidence、violated invariant、root cause hypothesis、local patch assessment、remediation considerations を raw result として受け取る。
4. raw result を `reviewer_result_bundle` に保持し、不足 field は `information_loss_notes` に記録する。
5. `overall_score` は group score の最小値にする。
6. すべての group score が 0.85 より大きい場合は `decision_basis: strict_pass` にする。
7. 非 trust_boundary の finding が上位観点と競合する場合は、`trust_boundary > behavior > contract > state_invariant` の優先度で裁定し、`decision_basis: priority_override_pass` にできる。
8. trust boundary failure は hard gate とし、平均 score や他観点の pass で相殺しない。
9. priority override した finding は削除せず、`priority_overrides` と `residual_risks` に残す。
10. 派生 field は `aggregation_trace` から `reviewer_result_bundle` の source field へ辿れるようにする。

## Integrated Review 受け取り

Codex implementation lane は integrated review result の `implementation_action` で分岐する。
`decision_basis` を再解釈せず、次の分岐だけを行う。

- `close`: completion packet に `codex_review_result` を転記して終了する
- `report_residual`: `priority_overrides` と `residual_risks` を completion packet に残して終了する
- `fix`: `reviewer_result_bundle`、`aggregation_trace`、`remediation_handoff` を読み、chosen strategy、chosen scope、狭すぎない理由、広げすぎない理由、planned changes、invariant tests を決めてから修正する
- `rerun_validation`: 指定された不足 validation だけを再実行し、review input を再構築する
- `rerun_codex_review`: 不足 payload を補い、product code を変更せず review input だけを再構築する

`fix` では、Codex review の `minimum_durable_fix_boundary` を修正範囲の最終命令として扱わない。
Codex implementation lane は修正者として次を completion packet に返す。

- `chosen_strategy`
- `chosen_scope`
- `why_this_scope`
- `why_not_narrower`
- `why_not_wider`
- `planned_changes`
- `invariant_tests`
- `used_review_signals`

## DO / DON'T

DO:
- 人間指示を受けたら skill、permissions、contract、承認済み `implementation-scope` を読みなおす
- 人間指示を skill / contract より上位の境界変更として扱わない
- distiller を implementation_tester / implementation_implementer より先に起動する
- `APIテスト` 先行条件を満たす時だけ implementation_tester を implementation_implementer より先に起動する
- unit test と原因未確定の regression test は実装後に implementation_tester へ渡す
- execution_group、parallelizable_with、parallel_blockers を見て ready wave を決める
- contract freeze 完了前の downstream handoff を開始しない
- `first_action` を含む `single_handoff_packet` だけを implementation_tester / implementation_implementer へ渡す
- scenario validation、suite-all、Sonar check を全 implementation handoff 完了後に実行する
- integrated review input に diff、scope、implementation result、validation result を含める
- review 可能な時だけ 4 観点 review agent を context 継承なしで並列 spawn する
- `implementation_action` に従って close / report_residual / fix / rerun_validation / rerun_codex_review を分岐する
- `fix` では reviewer result bundle と aggregation trace から chosen strategy と chosen scope を説明してから修正する
- `UI人間操作E2E` は final validation lane でだけ証明する

DON'T:
- 人間指示を理由に docs、`.codex`、`.codex/skills`、`.codex/agents` を変更しない
- spawn_agent 以外で実装、test 追加、調査をしない
- `first_action` がない handoff を広い調査で補わない
- `contract_freeze.status: required` を単なる notes として無視しない
- `parallelizable_with` に列挙されていない handoff を同じ wave という理由だけで並列実行しない
- final validation 前に scenario validation、suite-all、Sonar check を実行しない
- scenario validation failure を residual risk として close しない
- repo-local Sonar issue gate と Sonar server Quality Gate を混同しない
- review input 不足時に観点別 review agent を起動しない
- `rerun_codex_review` で product code を変更しない
- `fix` で symptoms だけを潰す局所修正に閉じない
- `fix` で why_not_narrower と why_not_wider なしに scope を選ばない
- docs、`.codex`、`.codex/skills`、`.codex/agents` を変更しない

## 参照パターン

- [orchestration-patterns.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-orchestrate/references/patterns/orchestration-patterns.md) を参照する。
- coverage は repo の `MINIMUM_COVERAGE = 70.0` を正本にする。
- `sonar_gate_result` は互換 field 名として残る場合があるが、意味は repo-local Sonar issue gate であり Sonar サーバ側 Quality Gate ではない。
## Checklist

- [implementation-orchestrate-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-orchestrate/references/checklists/implementation-orchestrate-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。

## Agent が持つもの

- active contract: [implementation_orchestrator.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/implementation_orchestrator/contracts/implementation_orchestrator.contract.json)
- permissions: [permissions.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/implementation_orchestrator/permissions.json)

## Maintenance

- output obligation を skill 本体へ戻さない。
- mode / variant contract を skill 配下の active 正本にしない。
