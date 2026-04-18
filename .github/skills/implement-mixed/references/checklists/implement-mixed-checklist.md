# Implement Mixed Checklist

## Knowledge Check

- [ ] API / Wails / DTO / gateway / adapter contract の接合点 scope が承認済みであることを確認した
- [ ] 両側の touched files を handoff と対応づけた
- [ ] frontend、backend、接合点 contract の validation evidence を分けた

## Common Pitfalls

- [ ] mixed を広い frontend / backend 同時変更の口実にしなかった
- [ ] API 接合点を変えずに UI と backend を同時に触らなかった
- [ ] Sonar gate と UI evidence の該当可否を返した
- [ ] docs / workflow 文書を変更しなかった
