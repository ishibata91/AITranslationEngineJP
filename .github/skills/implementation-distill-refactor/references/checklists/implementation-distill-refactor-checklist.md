# Implementation Distill Refactor Checklist

## Knowledge Check

- [ ] 不変条件、依存境界、変更候補を分けた
- [ ] preserved behavior を明示した
- [ ] affected package / component / tests を整理した

## Common Pitfalls

- [ ] 追加の設計判断をしなかった
- [ ] owned_scope 外の broad refactor を広げなかった
- [ ] product code / product test を変更しなかった
