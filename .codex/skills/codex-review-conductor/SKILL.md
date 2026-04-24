---
name: codex-review-conductor
description: Codex 実装後 review conductor 知識 package。観点グループ別 review agent を並列 spawn し、score と hard gate を集約する判断基準を提供する。
---

# Codex Review Conductor

## 目的

`codex-review-conductor` は知識 package である。
Copilot 実装完了後に `codex exec` で呼ばれた `review_conductor` が、観点グループ別 review agent を並列 spawn し、結果を集約する判断基準を提供する。

## いつ参照するか

- Copilot completion packet から Codex review を行う時
- diff から取得した実コードを観点別に score 化する時
- hard gate と overall decision を集約する時

## 知識範囲

- 観点グループ別 review agent の spawn
- score threshold `score > 0.85`
- hard gate の非相殺
- blocking findings と confidence notes の集約

## 原則

- review agent は context を引き継がず並列 spawn する
- 各 agent には同じ diff、scope、validation result を渡す
- overall pass は全観点 score が 0.85 より大きい場合だけ返す
- 権限・信頼境界は hard gate とし、平均 score で相殺しない
- Codex review は実装修正、docs 正本化、design review を行わない

## 標準パターン

1. Copilot payload に diff、implementation-scope、implementation result、final validation result があることを確認する。
2. `review_behavior`、`review_contract`、`review_trust_boundary`、`review_state_invariant` を context 継承なしで並列 spawn する。
3. 各 agent の score、confidence、findings、evidence を受け取る。
4. `overall_score` は group score の最小値にする。
5. すべての group score が 0.85 より大きい場合だけ `pass` にする。
6. blocking findings、confidence notes、next actions を分けて返す。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は agent contract に従う。

## DO / DON'T

DO:
- diff から取得した実コードを review 対象にする
- score と confidence を分ける
- hard gate failure を blocking finding にする
- blocked reason には不足 payload と再実行条件を書く

DON'T:
- 観点別 review を conductor が代替しない
- low score を平均で隠さない
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
