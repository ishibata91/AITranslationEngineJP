# Implementation Investigate Observe Checklist

## Knowledge Check

- [ ] temporary_changes に path と目的を残した
- [ ] 観測点を返却前に除去した
- [ ] cleanup_status を必ず返した

## Common Pitfalls

- [ ] 観測点を恒久修正として残さなかった
- [ ] owned_scope 外を変更しなかった
- [ ] cleanup 不能な時に続行しなかった
