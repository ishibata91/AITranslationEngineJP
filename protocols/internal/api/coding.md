# Wails Bind

`internal/api/` は frontend と backend の transport boundary を担当する。

- request と response の変換、native dialog、runtime event、子プロセス起動を置く。
- CRUD は store へ、翻訳手続きは engine へ委譲する。
- 純粋な翻訳規則と SQL を置かない。
- runtime event は push 通知に使い、query と command の主経路にしない。
- response は状態、種別、条件を返し、画面固有の操作可否を返さない。
