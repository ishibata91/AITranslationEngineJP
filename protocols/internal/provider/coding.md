# AI provider

`internal/provider/` は完成 prompt を外部 AI API へ送る。

- 同期処理は `Translator`、非同期 batch は `BatchTranslator` の契約を実装する。
- prompt の文面、辞書、口調、batch の段遷移を決めない。
- provider 固有の request、response、status、error を共通契約へ変換する。
- HTTP response body と接続を成功時と失敗時の両方で閉じる。
- test server か fake transport で request と失敗分類を検証し、実 API を単体テストから呼ばない。
