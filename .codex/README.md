# .codex

このディレクトリは、AITranslationEngineJp のマルチエージェント作業フローの正本です。
プロダクト仕様と設計は `docs/` を正本とし、agent の役割、入口、handoff、実装フローは `.codex/` を正本とします。

## 入口

- 標準入口: `skills/architect-direction/SKILL.md`
- 軽量入口: `skills/light-direction/SKILL.md`
- 実装 skill: `skills/light-work/SKILL.md`
- 軽量 review skill: `skills/light-review/SKILL.md`
- 役割定義:
  - `agents/architect.toml`
  - `agents/research.toml`
  - `agents/coder.toml`

## 標準フロー

### Heavy flow

`User -> Architect -> Research -> Architect plan -> Coder -> Architect review`

- User の要求は Architect が受ける
- Architect は必要なら Research に調査を委任する
- Architect は `docs/exec-plans/templates/heavy-plan.md` で heavy plan を固める
- Coder は確定した plan を見て実装する
- Architect が read-only で最終レビューする
- 最終的な accept / reroute / docs handoff は Architect が決める

### Light flow

`User -> Architect -> Short plan -> Coder -> Architect review`

- 仕様変更なし
- 低リスク
- 単一責務
- 短い plan で実装判断が固定できる

軽量フローでは `light-direction -> light-work` を使い、必要なら Architect が `light-review` を review checklist として使います。

## 重いフローと軽いフローの使い分け

`heavy` を使う:

- 仕様変更あり
- 複数レイヤーへ影響
- 複数ファイルや複数文書へ波及
- 高リスクまたは非可逆
- 調査結果がないと plan を固められない

`light` を使う:

- 仕様変更なし
- 低リスク
- 単一責務
- 短い plan で判断が固定できる

## 守ること

- User からの要求はまず Architect が受ける
- Architect は自分で実装せず、plan と handoff に責任を持つ
- Research は read-only で事実だけを集める
- Coder は plan 外の仕様判断を増やさない
- 最終 review は Architect が持つ
- 作業計画の保存先は `docs/exec-plans/` を使う
- skill は最小構成を維持し、細かい lane 分割を増やしすぎない
