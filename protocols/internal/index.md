# internal

Go backend の実装を置く。

- `bootstrap/` は production の依存を組み立てる。
- `api/` は Wails Bind の公開面を持つ。
- `engine/` は翻訳手続きを進行する。
- `core/` は副作用のない決定規則を持つ。
- `store/` は SQLite と keyring へ接続する。
- `provider/` は AI provider の port と実装を持つ。
- `model/` は backend 内で共有するデータ構造を持つ。
- `lexicon/` は外部辞書を読む adapter を持つ。
- `harness/` は翻訳手続きの決定的な統合検証を組み立てる。
