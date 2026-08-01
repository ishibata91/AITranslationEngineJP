# Task Plan: fix-batch-failure-handling

`plan.md` は branch 情報と、人間が見た事象、そこから起こした要求を持つ。
設計は `design.md`、確定仕様は `spec.md` が持つ。

## 事象

- 外部 batch 自体が `failed` になっても成功扱いになり、0 件のまま取り込みを始めようとする。
- 1 回に送った入力が queued prompt tokens の上限を超えた可能性がある。
- 1000 件ごとに分けても、複数の外部 batch を続けて送ると、先に送った分が queued prompt tokens を占有して後続分が失敗する可能性がある。

## 要求

- **R-1 batch を最大1000件ずつ順番に送る**: OpenAI または xAI への batch 送信は最大1000件に分け、現在の外部 batch が完了してから次の外部 batch を送り、画面には現在の外部 batch の総数、処理待ち、成功、失敗を表示する。
- **R-2 failed の取り込みを止めて失敗理由を表示する**: OpenAI の外部 batch が `failed` の場合は結果の取り込みと次の送信を開始せず、進行段と翻訳結果を変えず、外部 batch ID と OpenAI が返した失敗理由を画面に表示する。OpenAI の外部 batch が `completed` で一部だけが失敗した場合は、成功した結果を取り込み、失敗した分の翻訳結果を変えない。

## branch 情報

- `execution_branch`: `codex/fix-batch-failure-handling`
- `target_branch`: `master`
- `source_commit`: `a738b0b93dc4249afcd827d447a97453b91db3b4`

## やらないこと

- 入力 file の内容を取得または保存しない。
- 送信前に prompt tokens を算出して、OpenAI の model ごとに変わる上限を予測しない。
- `failed` になった外部 batch の自動再送信または自動取消は扱わない。
- `expired` または `cancelled` の既存動作は変更しない。
