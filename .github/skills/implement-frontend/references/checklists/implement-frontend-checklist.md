# Implement Frontend Checklist

## Knowledge Check

- [ ] 画面導線と state 反映を確認した
- [ ] Wails bridge 境界を確認した
- [ ] generated `wailsjs` を gateway 境界に閉じ込めた
- [ ] affected UI flow と console error を確認した
- [ ] frontend lint と format:check で拾われる境界違反を確認した

## Common Pitfalls

- [ ] design にない改善を足さなかった
- [ ] transport boundary を迂回しなかった
- [ ] View、ScreenController、Frontend UseCase から generated `wailsjs` を直接 import しなかった
- [ ] UI check に必要な evidence を残した
