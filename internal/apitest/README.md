# internal/apitest

**このディレクトリは API test 専用です。production code を置いてはいけません。**

## 目的

公開 API、Wails controller、bootstrap された backend graph を入口にして、
受け入れ条件を system-level で検証する API test を格納します。

`internal/integrationtest` は SQLite repository integration test 専用です。
API test は `internal/integrationtest` に置きません。

## ルール

- `_test.go` ファイルだけを置く。production code は置かない。
- package 宣言は `package apitest` とする。
- 開始点は公開 API、Wails controller、または bootstrap 済み controller にする。
- service や repository の直接呼び出しだけで完結する試験は置かない。

## 関連

- arch lint 設定: `/.go-arch-lint.yml` の `apitest` コンポーネント
- 既存の結合テスト置き場: `internal/integrationtest`
