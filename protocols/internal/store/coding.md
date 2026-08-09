# SQLite access

`internal/store/` は中心 SQLite への読書きと keyring secret store を担当する。

- SQL、transaction、row と model の変換を置く。
- schema の変更は `db/migrations/` へ置き、起動時の適用へ委譲する。
- 翻訳手順と画面判断を置かない。
- 複数 table を一つの操作として更新する場合は transaction で原子性を保つ。
- 再実行される書込は一意性と冪等性を test する。
- test は一時 database を使い、保存値と再実行結果を確認する。
