# 詳細仕様: AIサービス設定管理

- `upper_scenario_id`: `ai-provider-settings-management`
- `status`: `approved`
- `source_plan`: `docs/exec-plans/completed/ai-provider-settings-management/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/plan.md`
- `scenario_source`: `docs/exec-plans/completed/ai-provider-settings-management/scenario-design.md`
- `ui_source`: `docs/exec-plans/completed/ai-provider-settings-management/ui-design.md`
- `implementation_source`: `docs/exec-plans/completed/ai-provider-settings-management/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/final-validation.md`
- `review_source`: `docs/exec-plans/completed/ai-provider-settings-management/reviewback.*.yaml`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/review-summary.md`

## 要約

- 利用者は AI サービスごとに endpoint と APIキー状態を保存し、Job Setup、翻訳フェーズ、master-persona から参照できる。
- AIサービス設定は、model、処理方法、Batch API 切り替え、利用 provider の選択を保存しない。
- APIキー本体と raw payload は、UI、DTO、要約、log、debug 出力に出さない。

## 対象

- 対象利用者は、翻訳実行前に Gemini、xAI、LM Studio の接続設定を管理したい利用者である。
- 開始条件は、app-shell から `AIサービス設定` を開けることである。
- 完了状態は、provider ごとの endpoint、credential 参照状態、直近確認状態を確認できることである。
- 主要データは provider settings row、endpoint、credential 参照状態、secret store 内の APIキー本体、接続確認要約である。

## 仕様

- 利用者向け provider list は `gemini`、`lm_studio`、`xai` だけを扱い、fake provider を表示しない。
- AIサービス設定は provider ごとに endpoint と credential 参照状態を持つ。
- APIキー本体は secret store に保存する。DB は APIキー平文と復号可能値を保持しない。
- APIキー、raw request、raw response、raw prompt は UI、DTO、error summary、structured log、fake transport log、保存要約へ出さない。
- endpoint 変更後は接続確認状態を未確定へ戻し、古い確認結果を現在設定へ混入させない。
- 未設定へ戻す操作では provider settings row を残し、endpoint と APIキー状態を未設定へ戻し、secret 本体を削除する。
- endpoint はローカル運用の画面と保存要約で表示できる。secret は伏せ字または存在状態だけを表示する。
- provider settings の更新履歴は保存しない。
- Job Setup と master-persona は provider settings を参照し、個別の secret や endpoint を fallback にしない。
- Ready job の実行開始と retry は、AIサービス設定から最新 endpoint と credential 参照状態を再解決する。
- Running phase の job 側 runtime snapshot は provider、model、credential 状態分類、execution mode、batch mode だけを保存する。
- Running phase は provider adapter へ渡す endpoint と secret を AIサービス設定から解決し、job 側 summary や UI へ出さない。
- 保存結果と接続確認結果は raw payload ではなく分類と要約で観測する。
- 実装後検証は fake transport DI と fake secret store を使い、有料の実 AI API を呼ばない。

## 受け入れ根拠

- `SCN-AIPSM-001`: app-shell から AIサービス設定へ移動する。
- `SCN-AIPSM-002`: provider 単位で endpoint と APIキー状態を保存する。
- `SCN-AIPSM-003`: endpoint 変更後に接続確認状態を再評価する。
- `SCN-AIPSM-004`: 各参照側で provider、model、処理方法を選ぶ。
- `SCN-AIPSM-005`: 保存済み provider settings を再読込と再起動後に復元する。
- `SCN-AIPSM-006`: Job Setup と master-persona が provider settings を参照する。
- `SCN-AIPSM-007`: 失敗時も secret と raw payload を露出しない。
- `SCN-AIPSM-008`: fake transport DI で provider settings を検証する。
- `SCN-AIPSM-009`: provider settings の更新、未設定化、直近要約を扱う。
- `Q-AIPSM-001` から `Q-AIPSM-006` は人間回答済みである。
- 5 観点 reviewback は plan 上で `no_issue` として記録されている。

## UI 契約由来の恒久仕様

- 表示項目は AIサービス一覧、provider 別 endpoint、APIキー状態、接続確認状態、保存結果、接続確認結果である。
- AIサービス設定画面には model と Batch API 切り替えを出さない。
- 主要操作は AIサービス設定を開く、endpoint を編集する、APIキーを保存する、接続確認を行う、未設定へ戻すことである。
- endpoint または APIキーを変更した直後は、再確認待ちを表示する。
- secret は伏せ字または存在状態だけを表示し、APIキー本文は表示しない。
- 接続確認、保存失敗、endpoint 参照不能、provider 不正応答は分類と短い要約で表示する。
- error は色だけで伝えず、短い日本語文言を併記する。

## 対象外

- model、処理方法、Batch API 切り替え、利用 provider の保存。
- provider settings の更新履歴保存。
- provider SDK 実装方式、migration 番号、repository owner。
- 有料の実 AI API を使う必須検証。
