# 詳細仕様: 翻訳ジョブ設定

- `upper_scenario_id`: `translation-job-setup`
- `status`: `approved`
- `source_plan`: `docs/exec-plans/completed/translation-job-setup/plan.md`, `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/plan.md`, `docs/exec-plans/completed/2026-05-06-job-setup-input-cards/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/plan.md`
- `scenario_source`: `docs/exec-plans/completed/translation-job-setup/scenario-design.md`, `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/scenario-design.md`
- `ui_source`: `docs/exec-plans/completed/translation-job-setup/ui-design.md`, `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/ui-design.md`
- `implementation_source`: `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/plan.md`, `docs/exec-plans/completed/2026-05-06-job-setup-input-cards/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/final-validation.md`
- `review_source`: `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/reviewback.*.yaml`, `docs/exec-plans/completed/2026-05-06-job-setup-input-cards/reviewback.*.yaml`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/review-summary.md`

## 要約

- 利用者は登録済み入力データを選び、共通基盤と 3 つの翻訳段階の AI 設定を確認して翻訳 job を作成する。
- Job Setup は master-persona の AI 設定を既定値または保存元として扱わない。
- Job Setup では job 未作成 input を選択または削除できる。

## 対象

- 対象利用者は、翻訳入力データから翻訳 job を作成したい利用者である。
- 開始条件は、登録済み入力データ、共通辞書、共通ペルソナ、AIサービス設定の参照状態を確認できることである。
- 完了状態は、選択 input と各翻訳段階の AI 設定を固定し、Ready job を作成できることである。
- 主要データは入力データ、既存 job summary、共通辞書、共通ペルソナ、phase runtime settings、provider settings 参照状態、モデル一覧鮮度 token である。

## 仕様

- 入力データ選択はプルダウンではなくカード一覧にする。
- 入力カードは input 名、出自、登録日時、レコード数を表示する。
- カード選択で Job Setup の対象 input を切り替える。
- 既存 job が参照している input は Job Setup の入力候補に表示しない。
- 既存 job summary は候補除外と独立した応答 field として扱う。
- 同じ input の `existingJob` は参考表示だけに使い、job 作成の block 理由にはしない。
- 削除ボタンは job 未作成 input だけを削除する。
- input 削除は `X_EDIT_EXTRACTED_DATA` の親行削除を正本にし、関連行は DB の cascade に任せる。
- input 削除中は対象カードを操作不能にし、カード内に `削除中...` を表示する。
- 削除成功後は options 全再読込ではなく、画面状態から対象候補だけを除去する。
- Job Management の job 削除は input を残す。Job Setup の input 削除と混ぜない。
- Job Setup は単語翻訳、NPC ペルソナ生成、本文翻訳の 3 つの翻訳段階を扱う。
- 各翻訳段階は provider、model、execution mode、batch mode、APIキー状態を扱う。
- model 候補は provider ごとの model list API から取得する。
- provider model list の `sourceToken` はモデル一覧の非 secret 鮮度 token である。`sourceToken` は credential 参照実値、secret store key、endpoint、secret を含めない。
- 作成前検証と作成処理は、モデル一覧取得時の非 secret 鮮度 token を使う。古いモデル一覧由来の選択は stale として拒否する。
- API key が設定済みの場合だけ、API key が必要な provider の model list 外部取得を試みる。
- LM Studio は API key を要求しないため、API key 入力、API key 未設定 warning、credential select に出さない。
- batch mode は暗黙推定にしない。対象 provider は Gemini と xAI だけに限定し、checkbox または select で明示する。
- 3 つの翻訳段階で APIキー不足と model 未選択がない時だけ、翻訳 job を作成できる。
- Job Setup は credential 参照実値、secret store key、endpoint、APIキー本文を公開 DTO、UI、作成後 summary に出さない。
- 作成後の設定内容には、翻訳段階ごとの AIサービス、model、APIキー状態、実行方法、一括処理の有無だけを表示する。
- APIキー文字列、secret、外部サービスの raw data、内部ログ用識別子、モデル一覧鮮度 token は表示しない。

## 受け入れ根拠

- `translation-job-setup-phase-provider-settings` は 2026-05-04 に人間設計レビュー承認済みである。
- `translation-job-setup-phase-provider-settings` の最終検証は `python3 scripts/harness/run.py --suite all` 通過済みである。
- `2026-05-06-job-setup-input-cards` の最終検証は frontend-local、backend-local、all、system test、coverage、Sonar issue gate を通過済みである。
- `2026-05-06-job-setup-input-cards` の 5 観点 reviewback はすべて `must_fix_open: false` である。
- `reviewback.contract.yaml` の `existingJob` 契約指摘は解決済みである。
- `reviewback.state_invariant.yaml` の末尾候補削除時の選択状態指摘は解決済みである。

## UI 契約由来の恒久仕様

- 表示項目は入力カード一覧、共通辞書、共通ペルソナ、単語翻訳設定、NPC ペルソナ生成設定、本文翻訳設定、作成後の設定内容である。
- 各翻訳段階は AIサービス、model、実行方法、APIキー状態、一括処理、設定済みまたは未設定の状態を表示する。
- 主要操作は input 選択、job 未作成 input 削除、AIサービス選択、APIキー登録、model list 更新、model 選択、一括処理切り替え、job 作成である。
- APIキー登録は、APIキーが必要で未設定の AIサービスを選んだ場合だけ表示する。
- model list 更新は、APIキーが必要な AIサービスで APIキー未設定の場合は押せない。
- model 選択は、model list 取得に成功した場合だけ表示する。
- モデル一覧未更新、更新中、取得済み、取得失敗、APIキー未設定で更新不可を分けて表示する。
- 設定済み、APIキー未設定、model 未選択、model list 取得失敗、model list 更新中を分けて表示する。
- APIキー登録モーダル保存後は、model list を再更新する必要があることを表示する。
- stale なモデル一覧に基づく選択は作成できないことを表示する。
- desktop と mobile の両方で、入力確認、共通基盤、翻訳段階別設定、作成前確認、作成実行の順に読めることを維持する。
- 長い model 名、長い翻訳データ名、エラー文は表示領域からはみ出さない。

## 対象外

- migration reset、DB 再作成、cutover 手順。
- translation phase 実行本体。
- provider SDK 実装方式。
- task-local UI prototype と一時的な実装運用情報。
