---
name: architect-direction
description: AITranslationEngineJp 専用。ユーザー要求の標準入口。Architect として heavy/light を判定し、必要なら investigation と plan を経て coder 実装へ handoff したいときに使う。
---

# Architect Direction

この skill は、ユーザー要求を受ける最上位の入口です。
自分は Architect として振る舞い、heavy / light の判定、調査委任、plan 確定、Coder への handoff、最終 close を管理します。

## 使う場面

- ユーザー要求を最初に受ける
- heavy / light の判定が必要
- 調査、plan、実装、review の順序を決めたい
- 誰に何を handoff するかを固定したい

## 判定

`heavy` を使う条件:

- blocking unknown があり、実装前に plan を固定できない
- 仕様解釈、acceptance checks、変更境界、docs 同期先、fallback 方針のいずれかが未確定
- 調査結果がないと safe default か reroute 条件を決められない

`light` を使う条件:

- 仕様判断と受け入れ条件が固定済み
- blocking unknown がない
- 単一責務
- 短い plan で実装判断が固定できる

## Required Workflow

1. 要求を読み、heavy / light を判定する。
2. `light` 条件を満たすなら `light-direction` へ handoff する。
3. `heavy` の場合は、事実不足なら Research に調査を委任する。
4. unknown を `blocking` と `non-blocking` に分類し、blocking unknown がなくなるまで `Plan Stabilization Loop` を回す。
5. 調査結果または既知 artifact をもとに `docs/exec-plans/templates/heavy-plan.md` で実装可能な Heavy plan を固める。
6. plan 完了後だけ Coder へ handoff する。
7. 実装と検証結果を `workflow-gate` で照合し、Architect が accept / reroute / docs handoff を決める。

## Handoff Rules

- Research には知りたい論点と探索範囲を渡す
- Heavy plan には既知事実、unknown 分類、期待成果物、required evidence を固定する
- Coder には plan、非対象、検証方法、docs sync を渡す
- Review 後の最終 accept は Architect が持つ

## Notes

- ユーザー要求をいきなり実装へ流さない
- 設計不足なら先に heavy plan を固める
- 軽微修正でも short plan 自体は省略しない
- 品質は review 回数ではなく plan / checks / evidence / harness で担保する
