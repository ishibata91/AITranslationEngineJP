# Implementation: skip-already-japanese-translation

## 変更したfile

- `.go-arch-lint.yml`: 日本語文字判定の純粋規則をarchitecture lintへ登録した。
- `internal/core/japanesetext/japanesetext.go`: ひらがな、カタカナ、漢字のいずれかを含む判定を追加した。
- `internal/core/japanesetext/japanesetext_test.go`: 日本語文字と対象外文字の判定を検証した。
- `internal/engine/ingest.go`: 日本語を含む原文を原文と同じ訳文、翻訳済み状態で取込むようにした。
- `internal/store/ingest.go`: 取込時の訳文と翻訳状態を保存するようにした。
- `internal/engine/japanese_completion.go`: 保存済み未訳行を原文と同じ訳文、翻訳済み状態で完了する処理を追加した。
- `internal/engine/engine.go`: 同期翻訳前に日本語を含む保存済み未訳行を完了するようにした。
- `internal/engine/batch.go`: バッチ送信計画から日本語を含む保存済み未訳行を除外して完了するようにした。
- `internal/engine/ingest_test.go`: 取込時の翻訳済み保存を検証した。
- `internal/engine/engine_test.go`: 同期翻訳が外部翻訳を呼ばないことを検証した。
- `internal/engine/batch_integration_test.go`: バッチ翻訳が外部batchを作らないことを検証した。

## 仕様との対応

- R-1-1: 取込、同期翻訳、バッチ翻訳の各経路で、ひらがな、カタカナ、漢字を含む原文を原文と同じ訳文で翻訳済みにした。対応する実テストは `TestDispatchMarksJapaneseSourceAsTranslated`、`TestTranslateUntranslatedCompletesJapaneseRowsWithoutProviderCall`、`TestBatchSkipsJapaneseRows` である。
- R-1-2: 一文字を含む場合を含める文字判定を `TestContains` で検証した。
- R-1-3: 日本語文字を含まない原文を対象外にしない判定を `TestContains` と既存の `TestTranslateUntranslatedUsesOnlyPendingRowsWithoutPersonaRegeneration` で検証した。

## 検証結果

- `GOCACHE=/private/tmp/aitranslationenginejp-go-cache go test ./internal/engine ./internal/store ./internal/core/japanesetext`: 通過。
- `GOCACHE=/private/tmp/aitranslationenginejp-go-cache npm run verify:backend`: 通過。
- masterへのlocal merge後に `GOCACHE=/private/tmp/aitranslationenginejp-go-cache npm run test:backend`: 通過。

## 未確認事項と停止理由

- なし。

## 人間の指摘
