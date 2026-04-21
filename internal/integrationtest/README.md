# internal/integrationtest

**このディレクトリは integration test 専用です。production code を置いてはいけません。**

## 目的

複数の internal package (repository, infra/sqlite/dbinit など) をまたいで  
SQLite に対して end-to-end で動作を検証する integration test を格納します。

## ルール

- `_test.go` ファイルだけを置く。production code は置かない。
- package 宣言は `package integrationtest` とする。
- 複数コンポーネントにまたがる依存 (infra_sqlite + repository) は  
  この場所に限り `.go-arch-lint.yml` で許可されている。

## 関連

- arch lint 設定: `/.go-arch-lint.yml` の `integrationtest` コンポーネント
- 格納テストシナリオ: SCN-SMR-002〜005
