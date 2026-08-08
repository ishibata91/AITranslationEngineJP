# Go

- `gofmt` と `goimports` の出力を正とする。
- exported name は利用側が文脈なしで役割を読める名前にする。
- interface は利用側で定義し、必要な method だけを持たせる。
- `context.Context` は処理をまたいで伝播させ、構造体へ保存しない。
- error は無視せず、処理と対象が分かる文脈を加えて `%w` で包む。
- `panic` を通常の失敗経路に使わない。
- cleanup は `defer` で失敗経路にも適用する。
- 分岐を並べて検証する単体テストは table-driven test にする。
