# Skill Modification Checklist

## Knowledge Check

- [ ] `skill-agent-concept.md` と permissions を読んだ
- [ ] skill は knowledge package、agent は actor として分けた
- [ ] agent-owned permissions と 1:1 contract を確認した
- [ ] staged apply が必要な時は反映元保全と削除差分を確認した
- [ ] 廃止対象の file / directory を削除した

## Common Pitfalls

- [ ] skill 本体へ permissions や output obligation を戻さなかった
- [ ] mode / variant ごとの active contract file を増やさなかった
- [ ] 廃止対象の pointer、stub、legacy 説明を残さなかった
- [ ] product 実装や docs product 仕様変更を混ぜなかった
- [ ] 反映元を破壊する script や削除妥当性のない上書きを作らなかった
