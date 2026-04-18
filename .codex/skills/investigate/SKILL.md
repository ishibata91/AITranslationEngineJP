---
name: investigate
description: Codex 側の設計前調査知識 package。再現、UI 証跡、trace、risk-report を evidence first で扱う判断基準を提供する。
---

# Investigate

## 目的

`investigate` は知識 package である。
`investigator` agent が設計前に必要な証拠を集めるための、観測事実、UI 証跡、仮説、remaining gap の分け方を提供する。

UI check 専用 skill / agent は置かない。
設計前の UI evidence は `investigator` が `investigate` の一部として扱う。

## いつ参照するか

- 設計前に再現可否を確認する時
- UI evidence、console、画面状態を設計判断の証跡として確認する時
- trace の観測点と不足情報を整理する時
- design continuation の risk を短く返す時

## 参照しない場合

- implementation-scope 承認後の再現や再観測を扱う時
- 恒久修正や product test 追加が必要な時
- implementation review が主目的の時

## 知識範囲

- `reproduce`、`ui-evidence`、`trace`、`risk-report` の観点
- observed fact、UI evidence、hypothesis の分離
- evidence path と再現条件の残し方
- 設計を止める residual risk の表現

## 原則

- evidence のない結論を書かない
- 観測事実と仮説を混ぜない
- UI evidence は画面状態、console、screenshot、操作条件を分けて残す
- 実装 lane の調査は Copilot 側へ戻す

## 標準パターン

1. 調査目的と設計判断への影響を確認する。
2. 既知 facts、再現条件、UI check scope、未観測情報を分ける。
3. 最小の観測を行い、根拠 path を残す。
4. hypothesis は evidence level を明示する。
5. designer が次判断できる形で gaps と risks を返す。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は `investigator` agent contract に従う。

## DO / DON'T

DO:
- observed、UI evidence、inferred を分ける
- 証跡 path と再現条件を優先する
- 設計継続可否に効く gap を残す

DON'T:
- 恒久修正を始めない
- implementation-time investigation を扱わない
- owned_scope や対象 file を確定しない

## Checklist

- [investigate-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/investigate/references/checklists/investigate-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は `investigator` agent contract が決める。

## References

- agent spec: [investigator.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/investigator.agent.md)
- agent contract: [investigator.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/investigator/contracts/investigator.contract.json)

## Maintenance

- 権限、write scope、output obligation を skill 本体へ戻さない。
- 実装時調査は [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-investigate/SKILL.md) へ分ける。
- UI check 専用 skill / agent を戻さない。
