---
name: reporting-risks
description: 実装後または bugfix 後の residual risk を evidence 付きで短くまとめる。
---

# Reporting Risks

## Output

- remaining risks
- why they remain
- recheck suggestion

## Rules

- diff や evidence のない推測を書かない
- plan や packet の別正本を作らない
- closeout に必要な最小情報だけを返す

## Reference Use

- 着手前に `../directing-fixes/references/directing-fixes.to.reporting-risks.json` を参照して入力契約を確認する。
- `directing-fixes` へ返す時は `references/reporting-risks.to.directing-fixes.json` を返却契約として使う。
