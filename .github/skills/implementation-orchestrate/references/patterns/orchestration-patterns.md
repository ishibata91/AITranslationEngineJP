# Orchestration Patterns

## 目的

`implementation-orchestrate` が RunSubagent だけで handoff を分配するための判断パターンをまとめる。
agent contract の権限や output obligation は上書きしない。

## 適用ルール

- `implementation-scope` の handoff 見出しを RunSubagent 実行単位にする。
- オーケストレーター自身は RunSubagent 以外の tool を使わない。
- `depends_on` を守り、後続 handoff を先行しない。
- 1 handoff に複数 owned_scope が混ざる場合は実行せず reroute する。
- validation は subagent が返した result だけを集約する。
- closeout では coverage、Sonar、harness の gate evidence または blocked reason を必ず返す。

## 実行順パターン

- 通常: distiller -> tester -> implementer -> reviewer。
- 修正: investigator -> distiller -> tester -> implementer -> reviewer。
- refactor: distiller -> tester -> implementer -> reviewer。
- UI / mixed: backend handoff を先行し、接合点 evidence を集めてから reviewer へ渡す。

## 赤旗

- handoff が owned_scope、depends_on、validation command を持たない。
- オーケストレーターが直接 file read / search / edit / validation 実行をしている。
- validation failure の原因が設計不足なのに product code で吸収しようとしている。
- coverage、Sonar、harness の未実行理由が completion packet にない。
