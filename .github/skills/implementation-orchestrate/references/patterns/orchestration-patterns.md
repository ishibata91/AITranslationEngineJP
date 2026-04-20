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
- tester は implementer より先に必ず起動する。
- implementer には lane_context_packet と tester output 以外の追加文脈を渡さない。
- 1 handoff に複数 owned_scope が混ざる場合は実行せず reroute する。
- validation は subagent が返した result だけを集約する。
- closeout では coverage、Sonar、harness の gate evidence または blocked reason を必ず返す。

## 実行順パターン

- 通常: distiller -> tester -> implementer -> reviewer。
- 修正: investigator -> distiller -> tester -> implementer -> reviewer。
- refactor: distiller -> tester -> implementer -> reviewer。
- UI / mixed: backend handoff を先行し、それぞれ distiller -> tester -> implementer -> reviewer で扱う。
- distiller: default path で使い、single_handoff_packet 1 件だけから lane_context_packet を作る。

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
- implementer が tester より先に起動している。
- implementer に full implementation-scope、active work plan 全文、source artifacts、後続 handoff を渡している。
- オーケストレーターが直接 file read / search / edit / validation 実行をしている。
- validation failure の原因が設計不足なのに product code で吸収しようとしている。
- coverage、Sonar、harness の未実行理由が completion packet にない。
