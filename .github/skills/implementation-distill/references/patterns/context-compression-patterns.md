# Context Compression Patterns

## 目的

`implementation-distiller` が handoff 1 件を実装可能な context packet へ圧縮するための判断パターンをまとめる。
agent contract の権限や output obligation は上書きしない。

## 採用する考え方

- entry point、execution flow、architecture layer、dependency を分けて読む。
- 周辺 code を読む前に、対象 handoff と owned_scope を固定する。
- facts、inferred、gap を混ぜず、実装者が使う順に圧縮する。
- similar implementation を探し、既存 pattern を優先する。
- planning 情報は file path、function、risk、validation entry を含める。

## 適用ルール

- Wails binding、frontend gateway、backend service / infra 境界を分けて記録する。
- `docs/architecture.md` と `docs/coding-guidelines.md` は必要な判断だけに圧縮する。
- source artifact の文章を写すのではなく、implementation handoff に必要な制約へ変換する。
- 不足情報は `gaps` に残し、実装案で埋めない。

## 赤旗

- `required_reading` が広すぎて実装者の最初の一手が分からない。
- inferred を fact として書いている。
- owned_scope 外の architecture tour が長い。
- validation entry がないまま implementer へ渡している。
