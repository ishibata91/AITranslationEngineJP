# 正本化判断

- `task_id`: `2026-05-07-provider-settings-job-decoupling-implement`
- `decided_at`: `2026-05-07T21:58:26+0900`
- `decision`: `詳細仕様正本反映が必要`

## 理由

- Job Setup が credential 参照実値と endpoint を job 作成 payload、summary、公開 DTO へ持たない恒久仕様を追加した。
- AIサービス設定の endpoint と secret store 参照は、全 job provider 共通の参照元になった。
- Ready job は実行開始時に provider settings を再解決し、Running phase は開始時の非 secret 要約だけを job 側へ保存する仕様になった。
- モデル一覧の鮮度判定は、credential 参照実値を含まない token で行う仕様になった。

## 反映対象

- `docs/detail-specs/translation-job-setup.md`
- `docs/detail-specs/ai-provider-settings-management.md`
- `docs/detail-specs/translation-job-management.md`
- `docs/detail-specs/persona-generation-phase.md`
- `docs/detail-specs/body-translation-phase.md`
- `docs/detail-specs/term-translation-phase.md`

## 反映する恒久仕様

- Job Setup の各翻訳段階は provider、model、execution mode、batch mode、APIキー状態を扱う。
- Job Setup は credential 参照実値、secret store key、endpoint、APIキー本文を公開 DTO、UI、作成後 summary に出さない。
- Job Setup は provider model list の `sourceToken` に credential 参照実値、secret store key、endpoint、secret を含めない。
- 作成前検証と作成処理は、モデル一覧取得時の非 secret 鮮度 token を使い、古いモデル一覧由来の選択を stale として拒否する。
- Ready job の実行開始と retry は、AIサービス設定から最新 endpoint と credential 参照状態を再解決する。
- Running phase の job 側 runtime snapshot は provider、model、credential 状態分類、execution mode、batch mode だけを保存する。

## 反映しないもの

- 有料の実 AI API 呼び出しを必須検証にしない。
- provider SDK 実装方式、migration 番号、repository owner は詳細仕様正本へ入れない。
- task-local の review URL、fakeAPI 実行手順、検証ログ時刻は詳細仕様正本へ入れない。
