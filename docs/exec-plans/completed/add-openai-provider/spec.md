# Spec: add-openai-provider

`spec.md` はこの task の確定仕様として、要求ごとの仕様を持つ。要求は `plan.md`、設計理由・変更手順・図は `design.md` が持つ。
`design.md` と `spec.md` が食い違う場合は `spec.md` を優先する。

---

## R-1 OpenAI Batch API で翻訳できるようにする

- R-1-1（正常系）: 翻訳実行画面で OpenAI（batch）を選び、OpenAI API キーと翻訳対象 plugin を指定すると、`gpt-5.6-luna` を使って固有名と本文を順に送信、状態確認、取り込みできること
    - 前提条件: OpenAI API キーで `gpt-5.6-luna` と OpenAI Batch API を利用できる。
    - 確かめ方: 翻訳実行画面に OpenAI（batch）、OpenAI の公式 endpoint、`gpt-5.6-luna` が表示され、固有名と本文の進行を状態確認でき、取り込み後の訳された本文を結果一覧で確認する。
    - 対応する実テスト: `TestOpenAIBatch送信でLunaのChatCompletionsJSONLを使う`、`TestOpenAIBatch状態確認は終端状態だけDoneにする`、`TestBatchMatchesSyncEndToEnd`
- R-1-2（対象に入る側の境界）: OpenAI batch で成功した文と失敗した文が混在した場合、成功した訳された本文を結果一覧で確認でき、失敗した文は未訳のまま再送信できること
    - 前提条件: OpenAI batch の状態確認で取り込みできる状態になり、成功した文と失敗した文がある。
    - 確かめ方: 取り込み後の結果一覧で成功した訳された本文と未訳の文を確認し、同じ翻訳対象 plugin の未訳の文を再送信できることを確認する。
    - 対応する実テスト: `TestOpenAIBatch成功と失敗の結果を同時に取得する`、`TestBatchLeavesUntranslatedOnFailureLikeSync`
- R-1-3（対象に入る側のモデル一覧の境界）: OpenAI のモデル一覧に `gpt-5.6-luna` と他のモデルがある場合、`gpt-5.6-luna` が初期選択され、他のモデルも選べること
    - 前提条件: OpenAI API キーで取得したモデル一覧に `gpt-5.6-luna` と他のモデルがある。
    - 確かめ方: 翻訳実行画面のモデル一覧で `gpt-5.6-luna` が選択され、他のモデルも選べることを確認する。
    - 対応する実テスト: `OpenAIのモデル一覧ではLunaを先頭にして他モデルも残す`
- R-1-4（対象に入らない側の境界）: OpenAI（batch）以外を選んだ場合、OpenAI の公式 endpoint と `gpt-5.6-luna` を自動で使わず、OpenAI 互換 API または xAI batch で翻訳できること
    - 前提条件: 翻訳実行画面で OpenAI 互換 API または xAI batch を選ぶ。
    - 確かめ方: 翻訳実行画面に OpenAI 互換 API または xAI batch の endpoint とモデル一覧が表示され、OpenAI 互換 API では翻訳を実行でき、xAI batch では送信、状態確認、取り込みできることを確認する。
    - 対応する実テスト: `Test既存Batch進行のProviderはXAIになる`、`同期とxAIはOpenAIのendpointとLunaを使わない`
- R-1-5（進行中の提供元の境界）: 進行中の OpenAI batch を xAI batch として、または進行中の xAI batch を OpenAI（batch）として状態確認または取り込みできないこと
    - 前提条件: 翻訳対象 plugin に OpenAI batch または xAI batch の進行があり、翻訳実行画面で進行中と異なる選択肢を選ぶ。
    - 確かめ方: 翻訳実行画面で進行中の OpenAI batch または xAI batch と選択肢が一致しないことを確認し、進行中の固有名、本文、結果一覧が変更されないことを確認する。
    - 対応する実テスト: `Test進行中のOpenAIBatchをXAIとして状態確認または取り込みできない`
- R-1-6（OpenAI API キーの境界）: OpenAI API キーが空の場合、OpenAI batch の送信、状態確認、取り込みができないこと
    - 前提条件: 翻訳実行画面で OpenAI（batch）を選び、OpenAI API キーを空にする。
    - 確かめ方: 翻訳実行画面で送信、状態確認、取り込みを開始できないことを確認する。
    - 対応する実テスト: `TestOpenAIAPIキーが空ならBatch操作を開始しない`、`OpenAI APIキーが空ならBatch操作を許可しない`
- R-1-7（失敗した文だけの境界）: OpenAI batch に成功した訳された本文がなく失敗した文だけの場合も、失敗した文を未訳のまま再送信できること
    - 前提条件: OpenAI batch の状態確認で取り込みできる状態になり、失敗した文だけがある。
    - 確かめ方: 取り込み後の結果一覧で未訳の文を確認し、同じ翻訳対象 plugin の未訳の文を再送信できることを確認する。
    - 対応する実テスト: `TestOpenAIBatch結果Fileなしは空結果を返す`、`TestOpenAIBatchが全件失敗しても未訳を再送信できる`
- R-1-8（選択肢切替の境界）: OpenAI（batch）、OpenAI 互換 API、xAI batch の一つから別の選択肢へ切り替えた場合、切替前のモデル一覧と進行表示を持ち越さないこと
    - 前提条件: 翻訳実行画面で一つの選択肢のモデル一覧と進行表示を確認した後、別の選択肢へ切り替える。
    - 確かめ方: 翻訳実行画面で切替前のモデル一覧と進行表示が消え、切替後の選択肢の endpoint とモデル一覧が表示されることを確認する。
    - 対応する実テスト: `提供元を切り替えると選択先のendpointとモデルだけを持つ`

---

「対応する実テスト」は設計段階では空にする。`implementation-module` が最終検証で埋める。
