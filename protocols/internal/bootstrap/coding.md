# composition root

`internal/bootstrap/` は production の実装を生成して接続する。

- store、provider、外部辞書、engine、api の生成を集約する。
- 業務判断、翻訳規則、永続化処理を置かない。
- asset と環境依存の path は生成時に解決し、下位へ値として渡す。
