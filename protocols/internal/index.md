# backend

Go backend の実装は `internal/` 配下に置く。

- `bootstrap/` は composition root とする。
- `api/` は Wails Bind で公開する面とする。
- `engine/` は翻訳手続きの pipeline を orchestration する。
- `engine/` は純粋ルールを `internal/core/*` から import して束ねる。
- `core/` は副作用のない純粋不変ルールを集積する functional core とする。
- `core/` は 1 ルール 1 package とする。
- `core/` の対象 package は `dictionary`、`prompt`、`termderive`、`termusage`、`rolespeech`、`linefeatures`、`personatone`、`tone`、`batchplan` とする。
- `core/` は `os`、`provider`、`store`、`engine` を import しない。
- `store/` は SQLite へのアクセスを担当し、`sqlx` を使用する。
- keyring secret store は `store/secret/` に置く。
- `provider/` は AI クライアントの port と実装を担当する。
- `provider/` の多態 port は同期 `Translator` と非同期 `BatchTranslator` の 2 つとする。
- `model/` は概念的なデータモデルを持つ。
