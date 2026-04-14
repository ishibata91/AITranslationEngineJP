---
name: sonar-gate
description: backend task の close 前に SonarQube / SonarCloud の open issue gate を確認し、HIGH/BLOCKER と reliability/security の残件有無を返す single-role skill。
---

# Sonar Gate

## Goal

- backend task の close 前に Sonar gate を明示確認する
- implementation-review と Sonar issue gate を混ぜずに扱う
- open issue 件数と quality gate 状態を close 判断へ返す

## Rules

- backend を含む task の close 前に使う
- `HIGH` / `BLOCKER` の open issue 件数を確認する
- open reliability issue 件数を確認する
- open security issue 件数を確認する
- quality gate 状態と security hotspot の有無を補足する
- issue の解消方針は返してよいが、実装修正そのものは行わない

## Output

- `decision`: `pass` | `reroute`
- `open_high_or_blocker_count`
- `open_reliability_count`
- `open_security_count`
- `quality_gate_status`
- `supporting_issues`
- `recheck`
- `closeout_notes`

## Reference Use

- handoff contract は `../orchestrate/references/contracts/orchestrate.to.sonar-gate.json` を正本とする
- 返却 contract は `references/contracts/sonar-gate.to.orchestrate.json` を正本とする
