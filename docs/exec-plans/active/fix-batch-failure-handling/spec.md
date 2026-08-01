# Spec: fix-batch-failure-handling

`spec.md` はこの task の確定仕様として、要求ごとの仕様を持つ。要求は `plan.md`、設計理由と変更箇所は `design.md` が持つ。

---

## R-1 batch を最大1000件ずつ順番に送る

- R-1-1（正常系）: OpenAI または xAI へ1000件を超える batch を送る場合は、最大1000件を一つだけ送り、現在の外部 batch が完了してから次の最大1000件を送ること
    - 前提条件: batch へ送る対象が1001件以上ある
    - 確かめ方: OpenAI または xAI で、現在の外部 batch が完了する前に次の外部 batch が送られず、各外部 batch の件数が1000件以下であることを見る
    - 対応する実テスト: `internal/engine/batch_integration_test.go` の `TestBatchは1001件を最大1000件ずつ順番に送る`
- R-1-2（対象に入る側の境界）: OpenAI または xAI へちょうど1000件の batch を送る場合は、1000件を一つの外部 batch として送ること
    - 前提条件: batch へ送る対象が1000件ある
    - 確かめ方: OpenAI または xAI で、1000件の外部 batch が一つ作成されることを見る
    - 対応する実テスト: `internal/engine/batch_integration_test.go` の `TestBatchは1000件を一つの外部Batchとして送る`
- R-1-3（対象に入らない側の境界）: batch へ送る対象が無い場合は、OpenAI または xAI の外部 batch を作らないこと
    - 前提条件: 既訳の反映後に batch へ送る対象が0件である
    - 確かめ方: 画面に翻訳の完了が表示され、OpenAI または xAI に外部 batch が増えないことを見る
    - 対応する実テスト: `internal/engine/batch_integration_test.go` の `TestBatchは対象0件なら外部Batchを作らない`
- R-1-4（表示）: 画面の総数、処理待ち、成功、失敗は、現在処理している最大1000件の外部 batch の件数を表示すること
    - 前提条件: 1001件以上の対象を分けて送り、現在の外部 batch の状態を確認する
    - 確かめ方: 画面の件数が1000件以下であり、OpenAI または xAI で現在処理している外部 batch の件数と一致することを見る
    - 対応する実テスト: `internal/engine/batch_integration_test.go` の `TestBatchは1001件を最大1000件ずつ順番に送る`

---

## R-2 failed の取り込みを止めて失敗理由を表示する

- R-2-1（正常系）: OpenAI の外部 batch が `failed` の場合は、結果の取り込み、次の送信、進行段の更新を行わず、外部 batch ID と OpenAI が返した失敗理由を画面に表示すること
    - 前提条件: OpenAI の外部 batch が `status=failed` であり、`errors.data` に `code` と `message` がある
    - 確かめ方: 画面に外部 batch ID と OpenAI が返した `code` と `message` が表示され、翻訳結果と進行段が変わらないことを見る
    - 対応する実テスト: `internal/provider/openai_batch_test.go` の `TestOpenAIBatch状態確認はFailed理由を返す`、`internal/engine/batch_integration_test.go` の `TestBatch状態確認失敗は取り込みと進行を止める`
- R-2-2（対象に入る側の境界）: OpenAI の外部 batch が `failed` で失敗理由が空の場合は、結果を0件として取り込まず、外部 batch ID と `failed` を画面に表示すること
    - 前提条件: OpenAI の外部 batch が `status=failed` であり、`errors.data` に失敗理由が無い
    - 確かめ方: 画面に外部 batch ID と `failed` が表示され、翻訳結果と進行段が変わらないことを見る
    - 対応する実テスト: `internal/provider/openai_batch_test.go` の `TestOpenAIBatch結果取得は理由なしFailedを空結果にしない`、`internal/engine/batch_integration_test.go` の `TestBatch状態確認失敗は取り込みと進行を止める`
- R-2-3（対象に入らない側の境界）: OpenAI の外部 batch が `completed` で一部だけが失敗した場合は、成功した結果を取り込み、失敗した分の翻訳結果を変えないこと
    - 前提条件: OpenAI の外部 batch が `status=completed` であり、成功した分と失敗した分がある
    - 確かめ方: 画面で成功した翻訳が反映され、失敗した分の翻訳結果が変わらないことを見る
    - 対応する実テスト: `internal/provider/openai_batch_test.go` の `TestOpenAIBatch成功と失敗の結果を同時に取得する`、`internal/engine/batch_integration_test.go` の `TestBatchLeavesUntranslatedOnFailureLikeSync`

---

「対応する実テスト」は設計段階では空にする。`implementation-module` が最終検証で埋める。
