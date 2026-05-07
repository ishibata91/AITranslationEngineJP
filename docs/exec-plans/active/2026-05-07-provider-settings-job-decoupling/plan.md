# provider settings job decoupling

- `task_id`: `2026-05-07-provider-settings-job-decoupling`
- `lane`: `light-change-lane`
- `status`: `stopped`
- `current_artifact`: `軽量変更計画`
- `decision`: `設計戻し`
- `return_to`: `design-bundle`

## task 枠

- 人間依頼: クレデンシャル管理を DB レベルで Job から分離する。
- 人間依頼: AIサービス設定で設定できる値を Job から分離する。
- 人間依頼: secret store 情報と endpoint を全 job provider 共通設定にする。
- 変更禁止範囲: 軽量変更レーンではプロダクトコード、プロダクトテスト、docs 正本本文を変更しない。
- 確認したい結果: Job 系 DB、Job Setup、翻訳フェーズ実行が provider settings を共通設定として参照する境界を設計し直せること。

## 成果物状態

- `task 枠`: 完了。
- `軽量変更計画`: 完了。判定は `設計戻し`。
- `設計差分図`: 未着手。軽量変更として進めないため作成しない。
- `実装証跡`: 未着手。軽量変更として進めないため起動しない。
- `レビュー通過根拠`: 未着手。実装差分がないため起動しない。

## 停止理由

- DB schema から Job 所有の credential / endpoint snapshot を外す判断が必要である。
- Job Setup の保存契約と翻訳フェーズ開始時の provider settings 再解決契約が変わる。
- 既存詳細仕様は Running phase が開始時 snapshot を使うと定義している。
- 新しい永続仕様と公開契約の確定が必要なため、軽量変更では扱わない。

## 次に渡す入力

- 戻し先: `design-bundle`
- 先に作る成果物: シナリオ設計、設計差分図、人間設計レビュー、implementation-scope。
- 重点: `PROVIDER_SETTINGS` を共通設定正本にし、Job 側に残す snapshot の有無と範囲を決める。

