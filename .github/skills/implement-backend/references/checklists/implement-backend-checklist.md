# Implement Backend Checklist

## Knowledge Check

- [ ] layer 責務と dependency direction を確認した
- [ ] backend lint の format、static、arch、module 観点を確認した
- [ ] validation と error path を owned_scope 内で確認した
- [ ] lane_context_packet と lane-local validation を確認した
- [ ] `APIテスト` 先行時だけ tester output を確認した

## Common Pitfalls

- [ ] owned_scope 外の layer refactor を混ぜなかった
- [ ] usecase、service、controller で concrete 実装を new しなかった
- [ ] product test、fixture、snapshot、test helper を変更しなかった
- [ ] docs / workflow 文書を変更しなかった
