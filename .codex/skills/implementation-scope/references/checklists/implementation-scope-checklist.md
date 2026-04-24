# Implementation Scope Checklist

## Knowledge Check

- [ ] human review approval を確認した
- [ ] scenario-design に `needs_human_decision` が残っていないことを確認した
- [ ] 承認済み詳細要求タイプを validation intent の根拠にした
- [ ] handoff を owned_scope、depends_on、validation で分けた
- [ ] 各 handoff が `1 e2e use case × 1 validation intent` に収まっている
- [ ] 各 validation command が `completion_signal` を直接検証している
- [ ] 各 validation command が `owned_scope` と解消済み `depends_on` だけで pass できる
- [ ] 各 handoff の想定 touched files と changed lines を見積もった
- [ ] `15 files` / `800 changed lines` 以下を normal として扱った
- [ ] `16-25 files` または `801-1500 changed lines` の caution handoff には、1 件にする理由を `notes` に書いた
- [ ] `26 files` 以上または `1501 changed lines` 以上の split required handoff を 1 件として渡していない
- [ ] `40 files` 以上または `2500 changed lines` 以上の hard stop handoff は propose-plans へ戻した
- [ ] import / generation / settings save / preview / create / update / delete / export のうち、別 use case になっている処理を同一 handoff に混ぜていない
- [ ] domain 名や画面名だけを根拠に、複数 use case を同一 handoff にまとめていない
- [ ] layer をまたぐ handoff は、e2e completion_signal で完了判定できる
- [ ] broad validation を途中 handoff に置く場合は、required downstream scope と理由を `notes` に書いた
- [ ] 人間が Copilot に渡す entry、禁止事項、期待完了報告を明示した

## Common Pitfalls

- [ ] human review 前に implementation-scope を作らなかった
- [ ] 人間判断が残る scenario-design から implementation-scope を作らなかった
- [ ] layer だけを根拠に micro handoff を量産しなかった
- [ ] file 数と changed lines の基準を超える handoff を根拠なしに残さなかった
- [ ] Codex から Copilot へ直接 handoff しなかった
- [ ] docs 正本化を Copilot handoff に混ぜなかった
- [ ] validation command なしで handoff しなかった
- [ ] 未実装の後続 handoff を必要とする validation command を途中 handoff に入れなかった
- [ ] final validation で見るべき broad command を lane-local validation として扱わなかった
