# Wails gateway

`frontend/src/gateway/` は generated Wails API を画面から隠す。

- generated binding と backend model の import を閉じ込める。
- backend DTO を frontend の named type へ変換して返す。
- backend の状態を画面文言へ変換しない。
- API key と接続設定を永続化しない。
- generated file を編集しない。
- test は binding を置き換え、変換結果と error の伝播を検証する。
