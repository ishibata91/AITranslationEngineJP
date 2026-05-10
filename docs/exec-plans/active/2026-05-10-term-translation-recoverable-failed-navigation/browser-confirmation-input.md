# 実装後ブラウザ確認起動入力

## 呼び出し元

- `fix_lane`。

## 確認 URL

- `http://localhost:34115`。

## 起動状態

- Wails agent-browser 用 dev server を起動する。
- 起動 command: `npm run dev:wails:agent-browser`。
- backend は fake provider または dev provider の既存状態を使う。
- paid real API 到達を必要とする操作は実行しない。

## 操作経路

- `agent-browser open http://localhost:34115` でアプリを開く。
- 未完了一覧を表示する。
- `recoverable_failed` または `pending` の current phase を持つ未完了 job が表示される場合は、「現在の翻訳段階へ進む」を確認する。
- 確認対象 job が見つかる場合は、「現在の翻訳段階へ進む」を押して phase page へ移動できることを確認する。
- 確認対象 job が見つからない場合は、未完了一覧の snapshot、console、backend log を取得し、未確認理由として記録する。

## 操作期待値

- 未完了一覧が表示できる。
- `recoverable_failed` の current phase を持つ job は、`progress_percent` の値に関係なく現在の翻訳段階へ進める導線を持つ。
- `pending` の current phase を持つ job は、`progress_percent` の値に関係なく現在の翻訳段階へ進める導線を持つ。
- `phase_progress_aggregation_failed` は、現在の翻訳段階へ進む導線の block 理由として表示されない。
- 画面 console に今回の操作で発生した runtime error がない。

## 禁止操作

- 翻訳 phase の開始、再試行、再開、中断、取り消しを実行しない。
- paid real API に到達する操作を実行しない。
- provider 設定、credential、secret を変更しない。
- job の削除、作成、状態変更を実行しない。
- provider raw response、prompt、翻訳本文全文を証跡へ出さない。

## 安全条件

- 確認は閲覧と navigation に限定する。
- fake provider または dev provider の既存状態確認だけを許可する。
- 外部送信を伴う操作はしない。
- 破壊的操作はしない。

## 証跡出力先

- `tmp/agent-browser/2026-05-10-term-translation-recoverable-failed-navigation/`。
- `tmp/logs/wails-dev.log`。

## 期待する成果物

- snapshot path。
- errors path。
- 必要な screenshot path。
- 操作確認結果。
- console または network 異常。
- 未確認理由。
- 戻し先。
