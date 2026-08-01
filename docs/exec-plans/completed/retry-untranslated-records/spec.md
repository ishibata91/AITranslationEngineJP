# Spec: retry-untranslated-records

`spec.md` はこの task の確定仕様として、要求ごとの仕様を持つ。要求は `plan.md`、設計理由・変更手順・図は `design.md` が持つ。
`design.md` と `spec.md` が食い違う場合は `spec.md` を優先する。

---

## R-1 未訳件数を表示する

- R-1-1（正常系）: 同期翻訳の完了時に、未訳のまま残った件数を画面へ表示すること
    - 前提条件: 対象 plugin に未訳が複数残る
    - 確かめ方: 画面の完了案内に、対象 plugin に残る未訳と同じ件数が表示される
    - 対応する実テスト: `frontend/src/ui/screens/translation-run/translation-run-presentation.test.ts` の `未訳が複数残る場合は正確な件数と再実行対象を表示する`、`internal/store/target_plugin_test.go` の `TestCountUntranslated`
- R-1-2（対象に入る側の境界）: 同期翻訳の完了時に、未訳が 1 件残ったことを画面へ表示すること
    - 前提条件: 対象 plugin に未訳が 1 件だけ残る
    - 確かめ方: 画面の完了案内に「1 件が未訳のまま残りました」と表示される
    - 対応する実テスト: `frontend/src/ui/screens/translation-run/translation-run-presentation.test.ts` の `未訳が1件残る場合は1件の案内を表示する`、`internal/store/target_plugin_test.go` の `TestCountUntranslated`
- R-1-3（対象に入らない側の境界）: 同期翻訳の完了時に、未訳件数の案内を画面へ表示しないこと
    - 前提条件: 対象 plugin に未訳が残らない
    - 確かめ方: 画面の完了案内に未訳件数が表示されない
    - 対応する実テスト: `frontend/src/ui/screens/translation-run/translation-run-presentation.test.ts` の `未訳が0件の場合は案内を表示しない`、`internal/store/target_plugin_test.go` の `TestCountUntranslated`

---

## R-2 未訳だけを再実行する

- R-2-1（正常系）: 登録済み plugin の再実行が、既訳の収集、抽出、横断辞書の派生、取込、口調の集計を繰り返さず、中心 DB に残る未訳だけを翻訳し、訳のある本文を変更しないこと
    - 前提条件: 既訳の収集、抽出、横断辞書の派生、取込、口調の集計の完了が中心 DB に保存された登録済み plugin に、訳のある本文と未訳がある
    - 確かめ方: 既訳の収集、抽出、横断辞書の派生、取込、口調の集計が表示されず、翻訳進捗の総数が再実行前の未訳件数と一致し、訳のある本文が変わらない
    - 対応する実テスト: `internal/harness/retry_untranslated_test.go` の `TestRunExtractAndTranslateRetriesOnlyUntranslatedRows`、`internal/engine/engine_test.go` の `TestTranslateUntranslatedUsesOnlyPendingRowsWithoutPersonaRegeneration`
- R-2-2（対象に入る側の境界）: 登録済み plugin の再実行が、既訳の収集、抽出、横断辞書の派生、取込、口調の集計を繰り返さず、中心 DB に残る未訳だけを翻訳すること
    - 前提条件: 既訳の収集、抽出、横断辞書の派生、取込、口調の集計の完了が中心 DB に保存され、登録済み plugin の全ての本文が未訳である
    - 確かめ方: 既訳の収集、抽出、横断辞書の派生、取込、口調の集計が表示されず、翻訳進捗の総数が再実行前の未訳件数と一致する
    - 対応する実テスト: `internal/engine/engine_test.go` の `TestTranslateUntranslatedReportsAllPendingRows`
- R-2-3（対象に入らない側の境界）: 既訳の収集、抽出、横断辞書の派生、取込、口調の集計の完了が中心 DB に保存されていない plugin の実行が、未訳だけの再実行として扱われないこと
    - 前提条件: 登録済みでない plugin、または既訳の収集、抽出、横断辞書の派生、取込、口調の集計の完了が中心 DB に保存されていない登録済み plugin を実行する
    - 確かめ方: 既訳の収集、plugin の抽出、横断辞書の派生、取込、口調の集計を行ってから翻訳を始める
    - 対応する実テスト: `internal/harness/retry_untranslated_test.go` の `TestRunExtractAndTranslateRetriesOnlyUntranslatedRows`、`internal/store/target_plugin_test.go` の `TestSyncRetryReadinessLifecycle`

---

前提条件に「なし」と書いた仕様が、状況によらず成立させる振る舞いになる。
「対応する実テスト」は設計段階では空にする。`implementation-module` が最終検証で埋める。
