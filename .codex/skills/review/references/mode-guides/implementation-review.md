# Review: implementation-review

## Focus

- `implementation-scope` または accepted fix scope と diff の整合
- missing tests、validation gap、scope overrun の優先確認

## Rules

- `implementation-scope` は AI handoff 正本として扱う
- 好みや将来改善で判定しない
- fix lane では再現条件との整合も確認する
- `implementation-brief` の背景や理由は参照してよいが、判定の正本は `implementation-scope` と accepted fix scope に置く
