---
name: codex-review-conductor
description: Codex 実装後 review conductor 知識 package。観点グループ別 review agent を並列 spawn し、score と hard gate を集約する判断基準を提供する。
---

# Codex Review Conductor

## 目的

`codex-review-conductor` は知識 package である。
Copilot 実装完了後に人間が `codex exec` で起動した `review_conductor` が、観点グループ別 review agent を並列 spawn し、結果を集約する判断基準を提供する。

## いつ参照するか

- Copilot completion packet の review request payload から Codex review を行う時
- diff から取得した実コードを観点別に score 化する時
- hard gate と overall decision を集約する時

## 知識範囲

- 観点グループ別 review agent の spawn
- score threshold `score > 0.85`
- conflict priority `trust_boundary > behavior > contract > state_invariant`
- hard gate の非相殺
- blocking findings と confidence notes の集約
- Copilot が再解釈なしで受け取れる `copilot_action`

## 原則

- review agent は context を引き継がず並列 spawn する
- 各 agent には同じ diff、scope、validation result を渡す
- overall pass は `strict_pass` または `priority_override_pass` の場合だけ返す
- 権限・信頼境界は hard gate とし、平均 score で相殺しない
- 観点間 finding が競合する時は `trust_boundary > behavior > contract > state_invariant` の順で裁定する
- 優先度で退けた finding は削除せず、`priority_overrides` と `residual_risks` に残す
- Codex review は実装修正、docs 正本化、design review を行わない

## 標準パターン

1. Copilot payload に diff、implementation-scope、implementation result、final validation result があることを確認する。
2. diff、scope、implementation result、final validation result が不足している場合は、観点別 review を起動せず `copilot_action` を返す。
3. `review_behavior`、`review_contract`、`review_trust_boundary`、`review_state_invariant` を context 継承なしで並列 spawn する。
4. 各 agent の score、confidence、findings、evidence を受け取る。
5. `overall_score` は group score の最小値にする。
6. すべての group score が 0.85 より大きい場合は `decision_basis: strict_pass` にする。
7. 非 trust_boundary の finding が上位観点と競合する場合は、優先度で上位観点を採用し `decision_basis: priority_override_pass` にできる。
8. blocking findings、priority overrides、confidence notes、residual risks、copilot action を分けて返す。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は agent contract に従う。

## Copilot 返答

`review_conductor` は、人間が Copilot へ戻せる形で次 action を返す。
Copilot は戻された `copilot_action` を再解釈せず分岐する。
`copilot_action` は次のいずれかにする。

- `close`: `strict_pass` で residual がない場合に返す
- `report_residual`: `priority_override_pass` または低 confidence residual が残る場合に返す
- `fix`: Copilot が `copilot_patch_scope` 内で修正すべき場合に返す
- `rerun_validation`: validation 不足で観点別 review を起動しない場合に返す
- `rerun_codex_review`: diff、scope、implementation result など payload 不足で観点別 review を起動しない場合に返す

`rerun_validation` と `rerun_codex_review` は早期 return である。
この場合、観点別 review agent は spawn しない。

## DO / DON'T

DO:
- diff から取得した実コードを review 対象にする
- score と confidence を分ける
- hard gate failure を blocking finding にする
- blocked reason には不足 payload と再実行条件を書く
- `copilot_action` と `decision_basis` を必ず分ける

DON'T:
- 観点別 review を conductor が代替しない
- low score を平均で隠さない
- priority override した finding を消さない
- 実装修正や docs 正本化をしない

## Checklist

- [codex-review-conductor-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-review-conductor/references/checklists/codex-review-conductor-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。

## Agent が持つもの

- active contract: [review_conductor.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/review_conductor/contracts/review_conductor.contract.json)
- permissions: [permissions.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/review_conductor/permissions.json)

## Maintenance

- 観点グループを増減する時は conductor contract と `.codex/README.md` を同期する。
- threshold を変える時は全 review skill の pass 条件も同期する。
