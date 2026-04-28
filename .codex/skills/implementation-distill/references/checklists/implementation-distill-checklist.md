# Implementation Distill Checklist

## Knowledge Check

- [ ] facts、inferred、gap を分けた
- [ ] single_handoff_packet 1 件だけを source scope にした
- [ ] owned_scope の実 code を読んだ
- [ ] fix_ingredients に path、symbol/type/function、line number、why_needed_for_patch を入れた
- [ ] distracting_context に why_excluded と risk_if_read を入れた
- [ ] repository method、interface、field の present / absent を実 code fact として確認した
- [ ] first_action に path、symbol、line number、変更種別、1 つの clause_closed を入れた
- [ ] first_action が 1 edit で clause を閉じない場合は、同じ clause の最小 closure chain を change_targets に入れた
- [ ] change_targets と related_code_pointers に symbol/type/function と line number を入れた
- [ ] existing_patterns または none の探索範囲と影響を入れた
- [ ] validation_entry は cheap check を優先した
- [ ] requirements_policy_decisions に要件、実装方針、決定事項を要約した
- [ ] required reading と code pointer を owned_scope に絞った
- [ ] focused skill の知識だけを追加で参照した

## Common Pitfalls

- [ ] product code / product test を変更しなかった
- [ ] 類似 context を required_reading に混ぜなかった
- [ ] 推測だけで repository method、interface、field の追加を前提にしなかった
- [ ] first_action に partial、複数 clause、曖昧な advance を書かなかった
- [ ] existing_patterns の none を探索範囲なしで返さなかった
- [ ] cheap validation を検討せず広い command だけを返さなかった
- [ ] handoff 文章の言い換えだけで返さなかった
- [ ] 要件、実装方針、決定事項を required_reading に丸投げしなかった
- [ ] required_reading をファイル名だけの列挙にしなかった
- [ ] 要件や設計を追加しなかった
- [ ] full implementation-scope、active work plan 全文、source artifacts、後続 handoff を要求しなかった
- [ ] mode 別 active contract を使わなかった
