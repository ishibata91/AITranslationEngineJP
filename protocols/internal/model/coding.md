# Backend model

`internal/model/` は engine と store が共有するデータ構造を持つ。

- transport、SQL driver、Wails runtime、provider 固有型へ依存しない。
- 画面表示用の label と操作可否を持たない。
- 計算規則と I/O を置かない。
