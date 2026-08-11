# Spec: skip-already-japanese-translation

## R-1 既に日本語になっている文章を翻訳対象にしない

- R-1-1（正常系）: ひらがな、カタカナ、漢字のいずれかを含む翻訳対象の原文を、原文と同じ訳文で翻訳を実行せず翻訳済みとして保存すること
  - 前提条件: 翻訳対象の原文を翻訳対象として保存する前、または保存済みの翻訳対象の原文を翻訳する前
  - 確かめ方: 外部の翻訳処理へ送信する内容、保存された訳文、翻訳済みの状態を確認する
  - 対応する実テスト: `internal/engine/ingest_test.go` の `TestDispatchMarksJapaneseSourceAsTranslated`、`internal/engine/engine_test.go` の `TestTranslateUntranslatedCompletesJapaneseRowsWithoutProviderCall`、`internal/engine/batch_integration_test.go` の `TestBatchSkipsJapaneseRows`
- R-1-2（対象に入る側の境界）: ひらがな、カタカナ、漢字のいずれかを一文字だけ含む翻訳対象の原文を、原文と同じ訳文で翻訳を実行せず翻訳済みとして保存すること
  - 前提条件: 翻訳対象の原文が他の文字と一文字のひらがな、カタカナ、または漢字を含む
  - 確かめ方: 外部の翻訳処理へ送信する内容、保存された訳文、翻訳済みの状態を確認する
  - 対応する実テスト: `internal/core/japanesetext/japanesetext_test.go` の `TestContains`
- R-1-3（対象に入らない側の境界）: ひらがな、カタカナ、漢字を含まず既訳と完全一致しない翻訳対象の原文を、外部の翻訳処理へ送信し、返された訳文を保存すること
  - 前提条件: 翻訳対象の原文がひらがな、カタカナ、漢字を含まず、既訳と完全一致しない
  - 確かめ方: 外部の翻訳処理へ送信する内容と、返された訳文として保存された内容を確認する
  - 対応する実テスト: `internal/core/japanesetext/japanesetext_test.go` の `TestContains`、`internal/engine/engine_test.go` の `TestTranslateUntranslatedUsesOnlyPendingRowsWithoutPersonaRegeneration`
