# Spec: batch-retry-untranslated-records

`spec.md` はこの task の確定仕様として、要求ごとの仕様を持つ。要求は `plan.md`、設計理由・変更手順・図は `design.md` が持つ。
`design.md` と `spec.md` が食い違う場合は `spec.md` を優先する。

---

## R-1 batch の完了画面で未訳と再送信を示す

- R-1-1（正常系）: OpenAI と xAI の batch の取り込み後に、対象 plugin に残る未訳件数を表示し、未訳だけを再送信する操作であることを画面に示すこと
    - 前提条件: 本文 batch の取り込み後に対象 plugin へ未訳が複数残る
    - 確かめ方: 画面に対象 plugin の未訳と同じ件数と「未訳だけを再送信」が表示される
    - 対応する実テスト: `internal/api/app_test.go` の `TestGetBatchProgressCountsUntranslatedAfterCompletion`、`frontend/src/ui/screens/translation-run/translation-run-presentation.test.ts` の「未訳が複数残る場合は正確な件数と再送信操作を表示する」
- R-1-2（対象に入る側の境界）: OpenAI と xAI の batch の取り込み後に、未訳が 1 件残ったことと未訳だけを再送信する操作であることを画面に示すこと
    - 前提条件: 本文 batch の取り込み後に対象 plugin へ未訳が 1 件だけ残る
    - 確かめ方: 画面に「1 件が未訳のまま残りました」と「未訳だけを再送信」が表示される
    - 対応する実テスト: `frontend/src/ui/screens/translation-run/translation-run-presentation.test.ts` の「未訳が1件の場合は1件の案内を表示する」
- R-1-3（対象に入らない側の境界）: OpenAI と xAI の batch の取り込み後に未訳件数と「未訳だけを再送信」を表示しないこと
    - 前提条件: 本文 batch の取り込み後に対象 plugin へ未訳が残らない
    - 確かめ方: 画面に本文の取り込み完了が表示され、未訳件数と「未訳だけを再送信」が表示されない
    - 対応する実テスト: `frontend/src/ui/screens/translation-run/translation-run-presentation.test.ts` の「未訳が0件の場合は未訳案内を表示しない」
- R-1-4（未訳件数を取得できない場合）: OpenAI と xAI の batch の取り込み後に未訳件数を取得できない場合、未訳 0 件として表示しないこと
    - 前提条件: 本文 batch の取り込み後に対象 plugin の未訳件数を取得できない
    - 確かめ方: 画面に本文の取り込みは完了し、未訳件数の更新に失敗したことが表示され、未訳 0 件を表す画面が表示されない
    - 対応する実テスト: `internal/api/app_test.go` の `TestGetBatchProgressFailsWhenUntranslatedCannotBeCounted`、`frontend/src/ui/screens/translation-run/TranslationRunContainer.test.ts` の「取り込み成功後の未訳件数取得失敗で以前の表示を保つ」

---

## R-2 batch の再送信では未訳だけを処理する

- R-2-1（正常系）: 登録済み plugin の OpenAI と xAI の batch の再送信が、既訳の収集、抽出、横断辞書の派生、取込、口調の集計を繰り返さず、中心 DB に残る未訳だけを処理し、横断辞書または既訳を適用できない未訳だけを送り、訳のある本文を変更しないこと
    - 前提条件: 翻訳に必要な準備の完了が中心 DB に保存された登録済み plugin に、訳のある本文と未訳がある
    - 確かめ方: 画面に保存済みの準備を使って未訳だけを処理したことが表示され、固有名と本文の各 batch へ送った件数が横断辞書または既訳を適用できない未訳の件数と一致し、取り込み後も再送信前から訳のある本文が変わらない
    - 対応する実テスト: `internal/api/app_test.go` の `TestSubmitBatchTranslationReusesPreparedTarget`、`internal/engine/batch_integration_test.go` の `TestBatchRetryUsesOnlyPendingRowsWithoutPersonaRegeneration`、`internal/harness/retry_untranslated_test.go` の `TestBatchRetriesOnlyUntranslatedRows`
- R-2-2（対象に入る側の境界）: 登録済み plugin の OpenAI と xAI の batch の再送信が、中心 DB に残る 1 件の未訳だけを処理し、訳のある本文を変更しないこと
    - 前提条件: 翻訳に必要な準備の完了が中心 DB に保存された登録済み plugin に、横断辞書または既訳を適用できない未訳が 1 件だけあり、別に訳のある本文がある
    - 確かめ方: 画面に保存済みの準備を使って未訳だけを処理したことと batch へ送った件数 1 件が表示され、取り込み後も再送信前から訳のある本文が変わらない
    - 対応する実テスト: `internal/engine/batch_integration_test.go` の `TestBatchRetryUsesOnlyPendingRowsWithoutPersonaRegeneration`、`internal/harness/retry_untranslated_test.go` の `TestBatchRetriesOnlyUntranslatedRows`
- R-2-3（対象に入らない側の境界）: 翻訳に必要な準備の完了が中心 DB に保存されていない plugin の batch の送信が、未訳だけの再送信として扱われないこと
    - 前提条件: 登録済みでない plugin、または翻訳に必要な準備の完了が中心 DB に保存されていない登録済み plugin を OpenAI または xAI の batch へ送る
    - 確かめ方: 画面に batch を送信したことが表示され、保存済みの準備を使ったことは表示されない
    - 対応する実テスト: `internal/engine/batch_integration_test.go` の `TestBatchMatchesSyncEndToEnd`、`frontend/src/ui/screens/translation-run/TranslationRunContainer.test.ts` の「初回のbatch送信を未訳だけの再送信として表示しない」
- R-2-4（横断辞書または既訳を適用できる場合）: 登録済み plugin の OpenAI と xAI の batch の再送信が、横断辞書または既訳を適用できる未訳を OpenAI または xAI の batch へ送らずに処理すること
    - 前提条件: 翻訳に必要な準備の完了が中心 DB に保存された登録済み plugin に、横断辞書または既訳を適用できる未訳がある
    - 確かめ方: 再送信後の画面で対象の未訳に訳が入り、「完了」と表示され、batch へ送った件数に対象の件数が含まれない
    - 対応する実テスト: `internal/engine/batch_integration_test.go` の `TestBatchRetryCompletesWithoutExternalBatchWhenReferencesResolveAllRows`、`frontend/src/ui/screens/translation-run/TranslationRunContainer.test.ts` の「送信成功後の画面更新失敗で以前の表示を保つ」

---

## R-3 結果一覧を未訳だけに絞る

- R-3-1（正常系）: 結果一覧の「未訳のみ」を選択した場合、対象 plugin の結果全体を未訳だけに絞り、先頭ページから表示すること
    - 前提条件: 対象 plugin に複数ページにまたがる訳のある本文と未訳がある
    - 確かめ方: 「未訳のみ」を選択すると先頭ページが表示され、未訳件数が対象 plugin の未訳件数と一致し、一覧の全ての行が未訳になる
    - 対応する実テスト: `internal/api/app_test.go` の `TestListResultsPageFiltersUntranslatedAcrossSections`、`internal/store/result_pagination_test.go` の `TestResultPagesFilterUntranslatedRows`
- R-3-2（対象に入る側の境界）: 結果一覧の「未訳のみ」を選択した場合、未訳がないことを表示すること
    - 前提条件: 対象 plugin に訳のある本文があり、未訳がない
    - 確かめ方: 「未訳のみ」を選択すると、「未訳はありません」と表示される
    - 対応する実テスト: `frontend/src/ui/screens/translation-run/ResultsPanel.test.ts` の「未訳のみの選択と未訳なしを表示する」
- R-3-3（対象に入らない側の境界）: 結果一覧の「未訳のみ」を選択していない場合、対象 plugin の結果を訳の有無で絞らないこと
    - 前提条件: 対象 plugin に訳のある本文と未訳がある
    - 確かめ方: 「未訳のみ」を選択していない画面に、訳のある本文と未訳の両方が表示され、件数が両方の合計と一致する
    - 対応する実テスト: `internal/api/app_test.go` の `TestListResultsPageFiltersUntranslatedAcrossSections`、`internal/store/result_pagination_test.go` の `TestResultPagesFilterUntranslatedRows`
- R-3-4（xTranslator への書き出し）: 結果一覧の「未訳のみ」を選択して未訳が 0 件になっても、訳のある本文があれば xTranslator への書き出しを選べる状態を保つこと
    - 前提条件: 対象 plugin に訳のある本文があり、未訳がない
    - 確かめ方: 「未訳のみ」を選択して「未訳はありません」と表示された画面に「xTranslator へ書き出し」が表示される
    - 対応する実テスト: `frontend/src/ui/screens/translation-run/ResultsPanel.test.ts` の「未訳のみの選択と未訳なしを表示する」
- R-3-5（取得に失敗した場合）: 結果一覧の「未訳のみ」への切り替えで取得に失敗した場合、選択前の表示を保ってエラーを表示すること
    - 前提条件: 「未訳のみ」を選択した先頭ページの取得に失敗する
    - 確かめ方: チェックボックス、結果、表示中のページが選択前のままで、画面にエラーが表示される
    - 対応する実テスト: `internal/api/app_test.go` の `TestListResultsPageFailsWhenUnfilteredTotalCannotBeCounted`、`frontend/src/ui/screens/translation-run/TranslationRunContainer.test.ts` の「未訳のみの取得失敗時に選択前の一覧とページを保つ」

---

前提条件に「なし」と書いた仕様が、状況によらず成立させる振る舞いになる。
「対応する実テスト」は設計段階では空にする。`implementation-module` が最終検証で埋め、対応する実テストを置けなかった仕様は停止理由または残課題として人間へ上げる。
