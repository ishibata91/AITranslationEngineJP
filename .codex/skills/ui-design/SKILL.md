---
name: ui-design
description: Codex 側の UI 設計知識 package。UI 要件契約として主要操作、表示項目、状態差分、実装後確認観点を固定する基準を提供する。
---

# UI Design

## 目的

`ui-design` は知識 package である。
`designer` agent が UI を実装前の見た目 artifact ではなく UI 要件契約として扱うための、表示項目、操作、状態差分、実装後確認観点の見方を提供する。

実行境界、source of truth、handoff、stop / reroute は [design-bundle](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md) を参照する。

## いつ参照するか

- UI が関係する task で表示項目、操作、状態差分を固定する時
- state、variant、responsive、accessibility の差分を整理する時
- repo 側 `ui-design.md` に UI 契約と実装後確認観点を残す時

## 参照しない場合

- UI が不要で `plan.md` の `ui_design` が `N/A` の時
- product frontend code を実装する時
- docs 正本へ UI 仕様を反映するだけの時

## 知識範囲

- user-facing text、表示項目、主要 action、button enablement
- page section、state variant、layout constraints、accessibility
- loading、empty、error、disabled、progress、retry、success
- desktop / mobile で破綻してはいけない条件と実装後確認観点

## 原則

- UI は見た目 artifact ではなく実装が満たす契約として固定する
- 実装前の見た目 artifact を新規必須にしない
- 細かな visual polish は実装後に人間が実物を確認して直す
- generic な AI 風 UI や過剰な装飾を要求しない

## 標準パターン

1. interface の目的、利用者、主要 workflow を定義する。
2. 表示項目、主要 action、button enablement を固定する。
3. loading、empty、error、disabled、success などの状態差分を固定する。
4. responsive、overflow、accessibility の実装後確認観点を残す。
5. human が実装後に確認して直す visual polish を open question に残す。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は `designer` agent contract に従う。

## DO / DON'T

DO:
- UI contract と scenario の責務を分ける
- desktop と mobile の破綻条件を実装後確認観点として残す
- user-facing text は日本語を優先する

DON'T:
- 実装前の見た目 artifact を UI の必須 artifact にしない
- product code 実装へ踏み込まない
- 未承認で docs 正本化しない

## Checklist

- [ui-design-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/ui-design/references/checklists/ui-design-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は `designer` agent contract が決める。

## References

- template: [ui-design.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/templates/task-folder/ui-design.md)
- runtime skill: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md)
- agent contract: [designer.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/designer/contracts/designer.contract.json)

## Maintenance

- tool policy と output obligation を skill 本体へ戻さない。
- UI 要件契約を primary とし、見た目 artifact 必須へ戻さない。
- 長い判断表は references に分離する。
