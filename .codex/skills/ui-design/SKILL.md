---
name: ui-design
description: Codex 側の UI 設計知識 package。HTML mock を primary artifact として主要操作、状態差分、visual system 判断を固定する基準を提供する。
---

# UI Design

## 目的

`ui-design` は知識 package である。
`designer` agent が HTML mock を UI design の primary artifact として扱うための、構造、状態差分、visual system、evidence の見方を提供する。

実行権限、agent contract、handoff、stop / reroute は [designer.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.agent.md) が持つ。

## いつ参照するか

- UI が関係する task で HTML mock artifact を固定する時
- state、variant、responsive、accessibility の差分を整理する時
- repo 側 `ui-design.md` に HTML mock 参照と証跡を残す時

## 参照しない場合

- UI が不要で `plan.md` の `ui_design` が `N/A` の時
- product frontend code を実装する時
- docs 正本へ UI 仕様を反映するだけの時

## 知識範囲

- HTML mock file path、preview URL、screenshot artifacts
- page structure、component-like sections、state variant、visual tokens
- loading、empty、error、disabled、progress、retry、success
- desktop / mobile frame と overflow risk

## 原則

- task-local HTML mock を UI design の primary artifact にする
- `ui-design.md` は参照、判断、状態差分、証跡に限定する
- 既存 `docs/mocks/` と `docs/screen-design/` を優先する
- generic な AI 風 UI や過剰な装飾を避ける

## 標準パターン

1. interface の目的、利用者、主要 workflow を定義する。
2. HTML mock path、preview URL、review screenshot を固定する。
3. HTML / CSS / interaction state と screenshot evidence を確認する。
4. 状態差分を HTML mock の section、state control、または別表示として示す。
5. human review が必要な visual decision を open question に残す。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は `designer` agent contract に従う。

## DO / DON'T

DO:
- HTML mock と `ui-design.md` の責務を分ける
- desktop と mobile の破綻を evidence で確認する
- user-facing text は日本語を優先する

DON'T:
- product frontend code を UI 構造正本にしない
- product code 実装へ踏み込まない
- 未承認で docs 正本化しない

## Checklist

- [ui-design-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/ui-design/references/checklists/ui-design-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は `designer` agent contract が決める。

## References

- template: [ui-design.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/templates/task-folder/ui-design.md)
- mock canonicalization: [README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/mocks/README.md)
- agent spec: [designer.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.agent.md)
- agent contract: [designer.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/designer/contracts/designer.contract.json)

## Maintenance

- 権限、write scope、output obligation を skill 本体へ戻さない。
- HTML mock primary 方針を崩さない。
- 長い visual 判断表は references に分離する。
