# Orchestration Patterns

## 目的

`implementation-orchestrate` が RunSubagent だけで handoff を分配するための判断パターンをまとめる。
agent contract の権限や output obligation は上書きしない。

## 適用ルール

- `implementation-scope` の handoff 見出しを RunSubagent 実行単位にする。
- オーケストレーター自身は RunSubagent 以外の tool を使わない。
- `depends_on` を守り、後続 handoff を先行しない。
- subagent に渡す source scope は `single_handoff_packet` 1 件と、その distill 結果に限定する。
- distiller は tester / implementer より先に必ず起動する。
- distiller output は fix_ingredients、distracting_context、first_action、change_targets、requirements_policy_decisions、symbol / line number 付き related_code_pointers を持つ。
- distiller output は tester_context_packet、test_ingredients、test_required_reading、test_validation_entry を持つ。
- tester には tester_context_packet と test_subscope だけを渡し、full lane_context_packet、fix_ingredients 全体、change_targets 全体を渡さない。
- distiller output の first_action は 1 completion_signal clause に固定し、partial や複数 clause なら implementer に渡さない。
- distiller output が推測 method を fact にしている場合は implementer に渡さない。
- existing_patterns と validation_entry が探索理由を持たない場合は implementer に渡さない。
- tester は implementer より先に必ず起動する。
- implementer には lane_context_packet と tester output 以外の追加文脈を渡さない。
- tester / implementer の無応答、timeout、空 output、required field 欠落、insufficient_context は autonomous narrowing trigger として扱う。
- insufficient_context は各 agent contract の insufficient_context_criteria に一致する場合だけ narrowing trigger にする。
- criteria mismatch は contract violation として completion packet に残し、narrowing trigger にしない。
- narrowing は同じ single_handoff_packet 内で行い、completion_signal を削らず remaining subscopes を残す。
- narrowing 軸は completion_signal clause、public seam / API boundary、test target file、change target / symbol、validation command のいずれか 1 つに限定する。
- 1 handoff に複数 owned_scope が混ざる場合は実行せず reroute する。
- validation は subagent が返した result だけを集約する。
- closeout では coverage、Sonar、harness の gate evidence または blocked reason を必ず返す。

## 実行順パターン

- 通常: distiller -> tester -> implementer -> reviewer。
- 修正: investigator -> distiller -> tester -> implementer -> reviewer。
- refactor: distiller -> tester -> implementer -> reviewer。
- UI / mixed: backend handoff を先行し、それぞれ distiller -> tester -> implementer -> reviewer で扱う。
- distiller: default path で使い、single_handoff_packet 1 件だけから lane_context_packet を作る。
- narrowing retry: tester / implementer が criteria に一致する insufficient_context の場合は、同一 handoff 内で sub-scope を狭めて最大 2 回 retry し、進まなければ blocked_after_narrowing を返す。

## Harness と Sonar gate の確認方法

- `python3 scripts/harness/run.py --suite all` は `check_structure`, `check_execution`, `check_system_test`, `check_coverage` を順に実行する。
- `check_coverage.py` が `sonar-scanner` を起動し、`api/issues/search` 経由で以下の Sonar gate 条件を確認する:
  - coverage >= 70%
  - Security issue = 0
  - Reliability issue = 0
  - Maintainability HIGH issue = 0
- harness と Sonar gate の確認は `mcp_mcp_docker_mcp-exec` ツールで実行する。
- `sonar-project.properties` が存在しない場合は sonar-scanner をスキップする。

## 赤旗

- handoff が owned_scope、depends_on、validation command を持たない。
- distiller を tester / implementer より先に起動していない。
- distiller に full implementation-scope、active work plan 全文、source artifacts、後続 handoff を渡している。
- distiller output が handoff 文面の言い換えだけである。
- distiller output に fix_ingredients がない。
- distiller output に tester_context_packet がない。
- tester に full lane_context_packet、fix_ingredients、change_targets を渡している。
- distiller output が distracting_context を required_reading から分離していない。
- distiller output が要件、実装方針、決定事項を required_reading に丸投げしている。
- distiller output の first_action が partial、複数 clause、または曖昧な advance 表現である。
- distiller output が存在確認していない method / interface / field を fact にしている。
- distiller output の existing_patterns none に探索範囲と impact がない。
- distiller output の validation_entry が broad command だけで cheap check の検討理由がない。
- first_action や symbol / line number 付き related_code_pointers がない lane_context_packet を implementer に渡している。
- implementer が tester より先に起動している。
- tester / implementer の insufficient_context を broad investigation で埋めようとしている。
- criteria mismatch の insufficient_context を narrowing trigger にしている。
- autonomous narrowing で completion_signal を削っている。
- autonomous narrowing を理由に docs、implementation-scope、active work plan を書き換えている。
- implementer に full implementation-scope、active work plan 全文、source artifacts、後続 handoff を渡している。
- オーケストレーターが直接 file read / search / edit / validation 実行をしている。
- validation failure の原因が設計不足なのに product code で吸収しようとしている。
- coverage、Sonar、harness の未実行理由が completion packet にない。
