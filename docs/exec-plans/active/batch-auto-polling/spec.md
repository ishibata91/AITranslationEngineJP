# Spec: batch-auto-polling

`spec.md` はこの task の確定仕様として、要求ごとの仕様を持つ。`design.md` と食い違う場合は本 file を優先する。

---

## R-1 バッチ実行後は完了まで自動で状態確認する

- R-1-1（正常系）: OpenAI または xAI のバッチ実行を一度押すと、固有名段から本文段の完了まで10秒ごとに自動で進行状況を取得し、完了結果の取り込みと次のチャンク送信を続けること
  - 前提条件: 外部 batch が処理待ちまたは完了を返すこと
  - 確かめ方: 追加のボタン操作をせず、進行状況が固有名、本文、完了の順に変わり、結果一覧へ訳文が表示されることを見る
  - 対応する実テスト: `frontend/src/ui/screens/translation-run/TranslationRunContainer.test.ts` の「10秒ごとに複数チャンクと固有名から本文を完了まで進める」「状態確認の応答待ちには次の状態確認を開始しない」「自動状態確認中は開始時の接続情報を使い続ける」
- R-1-2（対象に入る側の境界）: 一つの進行段が複数チャンクに分かれる場合も、前のチャンクの完了後に次のチャンクを送り、同じ進行段の状態確認を続けること
  - 前提条件: 固有名段または本文段の送信対象が1000件を超えること
  - 確かめ方: 追加のボタン操作をせず、現在のチャンクの件数が更新され、最後のチャンク後に次の進行段または完了へ進むことを見る
  - 対応する実テスト: `frontend/src/ui/screens/translation-run/TranslationRunContainer.test.ts` の「10秒ごとに複数チャンクと固有名から本文を完了まで進める」
- R-1-3（対象に入らない側の境界）: 同期実行では10秒ごとの状態確認を開始しないこと
  - 前提条件: 配送方式が同期実行であること
  - 確かめ方: 実行を一度押すと従来どおり一回の同期実行が完了し、batch の進行状況を表示しないことを見る
  - 対応する実テスト: `frontend/src/ui/screens/translation-run/TranslationRunContainer.test.ts` の「同期実行ではbatchの状態確認を開始しない」

---

## R-2 画面を閉じた後は人の操作で既存のバッチ実行を再開する

- R-2-1（正常系）: 画面を閉じると、予約済みの次回状態確認を行わないこと
  - 前提条件: OpenAI または xAI の自動状態確認中に画面を閉じること
  - 確かめ方: 画面を閉じた後は外部 batch の状態確認が増えず、再表示しただけでは再開しないことを見る
  - 対応する実テスト: `frontend/src/ui/screens/translation-run/TranslationRunContainer.test.ts` の「画面終了後は予約済みの状態確認を行わない」
- R-2-2（対象に入る側の境界）: 画面を再表示してバッチ実行を押すと、保存済みの固有名段または本文段を新しく送信せずに自動状態確認を再開すること
  - 前提条件: 対象 plugin に完了していない batch の進行が保存されていること
  - 確かめ方: バッチ実行を押した後、新しい外部 batch を重複して送らず、保存済みの進行状況が表示されて完了まで進むことを見る
  - 対応する実テスト: `frontend/src/ui/screens/translation-run/TranslationRunContainer.test.ts` の「保存済みの進行を新規送信せず再開する」
- R-2-3（対象に入らない側の境界）: 保存済みの batch の進行が無い場合は、バッチ実行を押すと新しい外部 batch を送信して自動状態確認を開始すること
  - 前提条件: 対象 plugin に batch の進行が保存されていないこと
  - 確かめ方: バッチ実行を押した後、新しい外部 batch が一度だけ送信され、進行状況が表示されることを見る
  - 対応する実テスト: `frontend/src/ui/screens/translation-run/TranslationRunContainer.test.ts` の「進行なしでは一度だけ送信して直後に進行を取得する」
- R-2-4（provider が異なる境界）: 選択した provider が保存済みの batch の進行と異なる場合は、保存済みの進行を変更せず、新しい外部 batch を送信しないこと
  - 前提条件: 対象 plugin に、選択した provider と異なる provider の batch の進行が保存されていること
  - 確かめ方: バッチ実行を押すと provider が異なる理由を表示し、新しい外部 batch の送信と自動状態確認を開始しないことを見る
  - 対応する実テスト: `frontend/src/ui/screens/translation-run/TranslationRunContainer.test.ts` の「provider不一致では新規送信を行わない」
- R-2-5（画面の対象が変わる境界）: 自動状態確認中に対象 plugin または provider を変更した場合は、変更前に開始した処理の応答で変更後の進行状況と結果一覧を書き換えないこと
  - 前提条件: 自動状態確認の応答を待っている間に対象 plugin または provider を変更すること
  - 確かめ方: 変更前の状態確認が後から完了しても、変更後の対象 plugin の進行状況と結果一覧が変更前の内容にならないことを見る
  - 対応する実テスト: `frontend/src/ui/screens/translation-run/TranslationRunContainer.test.ts` の「provider変更後は変更前の遅延応答を無視する」「provider変更後は変更前の遅延した結果一覧を無視する」「plugin変更後は変更前の遅延応答を無視する」

---

## R-3 状態確認ボタンを削除する

- R-3-1（正常系）: OpenAI と xAI の画面には状態確認ボタンと手動の取り込みボタンを表示せず、バッチ実行を開始または再開する主操作を一つだけ表示すること
  - 前提条件: 配送方式が OpenAI または xAI であること
  - 確かめ方: batch の操作欄に状態確認ボタンと取り込みボタンがなく、主操作が一つだけ表示されることを見る
  - 対応する実テスト: `frontend/src/ui/screens/translation-run/TranslationRunContainer.test.ts` の「batch画面に状態確認と手動取り込みを表示しない」
- R-3-2（対象に入る側の境界）: バッチ実行の開始処理と自動状態確認中は主操作に実行中と表示して無効にし、二つの処理を重ねて開始しないこと
  - 前提条件: バッチ実行の開始処理または自動状態確認中であること
  - 確かめ方: 主操作に実行中と表示され、押せないことを見る
  - 対応する実テスト: `frontend/src/ui/screens/translation-run/TranslationRunContainer.test.ts` の「開始処理中の多重押下を防ぐ」
- R-3-3（状態確認が止まる境界）: 外部 batch の状態確認が失敗した場合は自動状態確認を止めて理由を表示し、途中の進行では主操作にバッチ実行を再開と表示して有効にすること
  - 前提条件: 外部 batch が失敗を返すこと
  - 確かめ方: 画面に外部 batch ID と失敗理由が表示され、主操作にバッチ実行を再開と表示されて押せることを見る
  - 対応する実テスト: `frontend/src/ui/screens/translation-run/TranslationRunContainer.test.ts` の「状態確認エラーで停止して再開操作を有効にする」
- R-3-4（完了の境界）: 未訳が無い完了後は主操作に完了と表示して無効にし、未訳が残る完了後は未訳だけを再送信と表示して有効にすること
  - 前提条件: 固有名段と本文段が完了していること
  - 確かめ方: 未訳件数が0件なら完了と表示されて押せず、未訳が1件以上なら未訳だけを再送信と表示されて押せることを見る
  - 対応する実テスト: `frontend/src/ui/screens/translation-run/TranslationRunContainer.test.ts` の「完了済みの未訳件数に応じて終了または未訳再送信を選ぶ」、`frontend/src/ui/screens/translation-run/translation-run-presentation.test.ts` の「未訳が複数残る場合は正確な件数と再送信操作を表示する」
- R-3-5（対象に入らない側の境界）: 同期実行の実行ボタンと進行表示を変更しないこと
  - 前提条件: 配送方式が同期実行であること
  - 確かめ方: 同期実行では従来の実行ボタンと進行表示が表示されることを見る
  - 対応する実テスト: `frontend/src/ui/screens/translation-run/TranslationRunContainer.test.ts` の「同期実行ではbatchの状態確認を開始しない」
