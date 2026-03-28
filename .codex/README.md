# .codex

このディレクトリは、AITranslationEngineJp のマルチエージェント作業フローの正本です。
プロダクト仕様と設計は `docs/` を正本とし、agent の役割、入口、handoff、実装フローは `.codex/` を正本とします。
標準原則は `flow light, gate heavy` とし、品質は review の段数ではなく gate 契約と evidence で担保します。

## 入口

- 標準入口: `skills/architect-direction/SKILL.md`
- 軽量入口: `skills/light-direction/SKILL.md`
- 実装 skill: `skills/light-work/SKILL.md`
- Gate skill: `skills/workflow-gate/SKILL.md`
- 補助 review skill: `skills/light-review/SKILL.md`
- 役割定義:
  - `agents/architect.toml`
  - `agents/research.toml`
  - `agents/coder.toml`

## 標準フロー

### Heavy flow

`User -> Architect -> Research -> Plan Stabilization Loop -> Coder -> Workflow Gate -> Architect accept`

- User の要求は Architect が受ける
- Architect は必要なら Research に調査を委任する
- Architect は `docs/exec-plans/templates/heavy-plan.md` で heavy plan を固める
- `Plan Stabilization Loop` では blocking unknown がなくなるまで `Research -> plan update` を反復する
- Coder は確定した plan を見て実装する
- `Workflow Gate` は plan 適合性、証跡不足、docs 同期漏れ、reroute 要否だけを判定する
- Architect が gate 結果を読んで最終 accept / reroute / docs handoff を決める
- 最終的な accept / reroute / docs handoff は Architect が決める

### Light flow

`User -> Architect -> Short plan -> Coder -> Workflow Gate -> Architect accept`

- 仕様判断と受け入れ条件が固定済み
- blocking unknown がない
- 単一責務
- 短い plan で実装判断が固定できる

軽量フローでは `light-direction -> light-work -> workflow-gate` を使い、必要なら Architect が `light-review` を補助 checklist として使います。

## 重いフローと軽いフローの使い分け

`heavy` を使う:

- blocking unknown があり、実装前に plan を固定できない
- 仕様解釈、acceptance checks、変更境界、docs 同期先、fallback 方針のいずれかが未確定
- 調査結果がないと safe default か reroute 条件を決められない

`light` を使う:

- 仕様判断と受け入れ条件が固定済み
- blocking unknown がない
- 単一責務
- 短い plan で判断が固定できる

## Unknown の扱い

- unknown は `blocking` と `non-blocking` に分類する
- `blocking unknown` は scope、acceptance checks、変更境界、docs 同期先、fallback 方針を固定できない unknown を指す
- `non-blocking unknown` は safe default や assumption で固定でき、実装開始を止めない unknown を指す
- Heavy flow では blocking unknown が残る間は Coder に handoff しない
- Light flow では blocking unknown が見つかった時点で heavy へ reroute する

## Gate の役割

- 品質は review の往復回数ではなく、plan / checks / evidence / harness の強さで担保する
- `workflow-gate` は approved plan、変更要約、checks 結果、docs 更新有無、残る non-blocking unknown を入力として受ける
- gate の出力は `decision`、`missing_evidence`、`contract_breaks`、`docs_sync`、`recheck` を含む
- gate は美しさや好みを判定しない
- 繰り返し見つかる review 指摘は gate や review に残さず、`docs/executable-specs.md` または harness へ昇格する

## 守ること

- User からの要求はまず Architect が受ける
- Architect は自分で実装せず、plan と handoff に責任を持つ
- Research は read-only で事実だけを集める
- Coder は plan 外の仕様判断を増やさない
- 最終 accept は Architect が持つ
- 作業計画の保存先は `docs/exec-plans/` を使う
- skill は最小構成を維持し、細かい lane 分割を増やしすぎない
