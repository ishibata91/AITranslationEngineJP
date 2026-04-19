# Diagramming Checklist

## Knowledge Check

- [ ] 図の source of truth と review artifact を分けた
- [ ] diagram kind と読者を明示した
- [ ] PlantUML や library が関係する場合は Context7 を確認した
- [ ] 一時 PNG を生成し、AI が画像で可読性を確認した
- [ ] 図の分割条件に当たらないことを確認した
- [ ] 正本 docs の用語、構造主語、layer 名と diagram の package 名を揃えた
- [ ] 正本図では node を勝手に畳まず、edge、legend、note の粒度で調整した

## Common Pitfalls

- [ ] review 用 SVG / PNG を正本にしなかった
- [ ] validation なしで完了扱いにしなかった
- [ ] product code 変更を diagramming に混ぜなかった
- [ ] 線、文字、legend、note の重なりを放置しなかった
- [ ] docs 正本図を差分図 style にしなかった
- [ ] 可読性向上を理由に正本 node を消さなかった
- [ ] 角ばった線指定で読みやすさを下げなかった
