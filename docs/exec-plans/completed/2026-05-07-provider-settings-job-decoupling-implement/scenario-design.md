# Scenario Design: 2026-05-07-provider-settings-job-decoupling-implement

- `skill`: `scenario-design`
- `status`: `pending-human-review`
- `source_plan`: `./plan.md`
- `task_frame`: `./task-frame.md`
- `ui_source`: `required-after-scenario-review`
- `final_artifact_path`: `TBD after human review`
- `topic_abbrev`: `PSJD`

## 根拠

- `./task-frame.md`
- `./scenario-candidates.actor-goal.md`
- `./scenario-candidates.lifecycle.md`
- `./scenario-candidates.state-transition.md`
- `./scenario-candidates.failure.md`
- `./scenario-candidates.external-integration.md`
- `./scenario-candidates.operation-audit.md`
- `docs/detail-specs/ai-provider-settings-management.md`
- `docs/detail-specs/translation-job-setup.md`
- `docs/er.md`
- `docs/architecture.md`

## 固定要件

- `PROVIDER_SETTINGS` を provider ごとの共通設定正本にする。
- Job Setup は provider、model、execution mode、batch mode の選択値だけを扱う。
- Job 側 DB は secret store 情報、credential 参照実値、endpoint を所有しない。
- Ready job は実行開始前に最新 provider settings を再解決する。
- Running phase は開始時の再解決結果を使い、DB には非 secret の分類要約だけを残す。
- APIキー本体、secret、raw request、raw response、raw prompt、復号可能値は出力しない。

## DB 境界

### Job から外す値

- provider ごとの endpoint。
- secret store 参照実値と `credential_ref` 実値。
- APIキー本体、復号可能値、credential 管理。
- provider settings の更新履歴と Job 側 revision 所有。
- raw request、raw response、raw prompt。

### Job に残す値

- `TRANSLATION_JOB` の input 参照と job 状態集約に必要な値。
- `JOB_PHASE_RUN` の phase type、状態、provider、model、execution mode、batch mode。
- フェーズ対象、翻訳結果、出力成果物への参照。
- 実行時の非 secret 分類要約。
- 再実行時の同一 phase run 状態。

### provider settings 側に寄せる値

- provider ごとの endpoint。
- credential 参照状態。
- 接続確認状態と接続確認要約。
- 未設定化した provider settings row。
- secret store 内 APIキー本体への管理境界。

## 候補統合

候補ファイル上の見出し数は 42 件である。
`plan.md` は 43 候補と書くが、指定された 6 ファイルに存在する候補を母集団にした。

| 採用シナリオ | 統合候補 |
| --- | --- |
| `SCN-PSJD-001` | `actor-goal:CAND-PSJD-001`, `external-integration:CAND-PSJD-EI-001`, `external-integration:CAND-PSJD-EI-002` |
| `SCN-PSJD-002` | `actor-goal:CAND-PSJD-002`, `actor-goal:CAND-PSJD-003`, `lifecycle:CAND-PSJD-001`, `operation-audit:CAND-PSJD-OA-002` |
| `SCN-PSJD-003` | `actor-goal:CAND-PSJD-006`, `failure:CAND-PSJD-F01`, `failure:CAND-PSJD-F02`, `failure:CAND-PSJD-F04`, `external-integration:CAND-PSJD-EI-003`, `external-integration:CAND-PSJD-EI-004`, `external-integration:CAND-PSJD-EI-005` |
| `SCN-PSJD-004` | `actor-goal:CAND-PSJD-004`, `actor-goal:CAND-PSJD-005`, `lifecycle:CAND-PSJD-002`, `lifecycle:CAND-PSJD-006`, `state-transition:CAND-PSJD-ST-001`, `state-transition:CAND-PSJD-ST-002`, `external-integration:CAND-PSJD-EI-006`, `operation-audit:CAND-PSJD-OA-001`, `operation-audit:CAND-PSJD-OA-004`, `operation-audit:CAND-PSJD-OA-005` |
| `SCN-PSJD-005` | `lifecycle:CAND-PSJD-003`, `lifecycle:CAND-PSJD-004`, `state-transition:CAND-PSJD-ST-003`, `state-transition:CAND-PSJD-ST-004`, `failure:CAND-PSJD-F05`, `external-integration:CAND-PSJD-EI-007`, `operation-audit:CAND-PSJD-OA-003` |
| `SCN-PSJD-006` | `actor-goal:CAND-PSJD-007`, `lifecycle:CAND-PSJD-007`, `lifecycle:CAND-PSJD-008`, `state-transition:CAND-PSJD-ST-005`, `state-transition:CAND-PSJD-ST-006`, `state-transition:CAND-PSJD-ST-007`, `failure:CAND-PSJD-F03`, `failure:CAND-PSJD-F06`, `operation-audit:CAND-PSJD-OA-006`, `operation-audit:CAND-PSJD-OA-007` |

却下候補は、Running phase が provider settings 更新へ途中追従する案である。
対象候補は `lifecycle:CAND-PSJD-005` である。
理由は、今回の固定前提が「endpoint と credential の所有は共通設定、phase には実行時の非 secret 要約だけ残す」案を優先しているためである。

保留候補はない。
UI 文言、レイアウト、操作列は `ui-design.md` で別成果物にする。

## 採用シナリオ

### SCN-PSJD-001 provider settings が共通設定として保存される

受け入れ条件:
provider settings row は provider ごとの endpoint、credential 参照状態、接続確認状態を扱う。
Job 系 DB は endpoint、secret store 参照実値、APIキー本体を作成しない。
未設定化では provider settings row を残し、endpoint と APIキー状態を未設定に戻す。

実行テスト種別: `APIテスト`
実行段階: `実装後`
公開接点 / API 境界: AIサービス設定の保存、未設定化、再読込。
入力開始点: provider、endpoint、APIキー保存状態、未設定化操作。
主要結果: provider settings 側の共通設定が更新され、secret 本体は secret store 境界に閉じる。
主要観測点: provider settings 保存要約、secret store fake、Job 系 DB の非保持。
公開接点確認: あり。

### SCN-PSJD-002 Job Setup が選択値だけで Ready job を作成する

受け入れ条件:
Job Setup は 3 つの翻訳段階で provider、model、execution mode、batch mode を選択する。
Ready job の永続値は選択値だけである。
作成後の設定要約は APIキー状態分類を表示できるが、Job 所有値として保存しない。

実行テスト種別: `UI人間操作E2E`
実行段階: `最終検証`
開始操作: Job Setup を開き、登録済み input を選択する。
入力方法: 画面で各翻訳段階の AIサービス、model、実行方法、一括処理を選ぶ。
主要操作列: input 選択、各段階の設定、job 作成、作成後要約確認。
主要観測点: 作成後要約、Job 系 DB、DTO、structured log。
UI-visible 結果: endpoint と secret store 参照値が表示されない。
fake / stub 方針: fake provider settings と fake secret store を使う。

### SCN-PSJD-003 Job Setup が不足状態と外部取得可否を分ける

受け入れ条件:
APIキーが必要な provider は、APIキー未設定なら model list API を呼ばない。
LM Studio は APIキーを要求せず、endpoint だけで model list API を扱う。
model 未選択、model list 取得失敗、APIキー未設定は別状態として扱う。

実行テスト種別: `UI人間操作E2E`
実行段階: `最終検証`
開始操作: Job Setup で provider を選ぶ。
入力方法: 設定済み、未設定、参照不能の provider settings を切り替える。
主要操作列: provider 選択、model list 更新、model 選択、job 作成可否確認。
主要観測点: model list 外部呼び出し回数、画面状態、作成不可理由。
UI-visible 結果: 不足理由が分類され、raw payload は出ない。
fake / stub 方針: fake transport と fake secret store を使い、有料の実 AI API は使わない。

### SCN-PSJD-004 Ready job 実行開始時に最新 provider settings を再解決する

受け入れ条件:
Ready job は実行開始前に最新 provider settings を再解決する。
Job 作成後に provider settings が変わった場合、実行開始は更新後の共通設定を使う。
再解決で未設定または参照不能を検出した場合、Running phase を開始しない。
失敗要約は分類と短い説明だけを残す。

実行テスト種別: `APIテスト`
実行段階: `実装後`
公開接点 / API 境界: Ready job 実行開始。
入力開始点: Ready job、更新後 provider settings、fake secret store 状態。
主要結果: 再解決済み分類が残り、Job 側 fallback は使われない。
主要観測点: 実行開始結果、phase 状態、fake provider execution 入力要約。
公開接点確認: あり。

### SCN-PSJD-005 Running phase が開始時の非 secret 要約だけを保持する

受け入れ条件:
Running phase は phase 開始時の provider settings 再解決結果を実行に使う。
provider settings が実行中に更新または未設定化されても、実行中 phase は途中で参照結果を混在させない。
phase runtime snapshot は endpoint 原文、credential 参照実値、APIキー本体を保存しない。
保存できる値は provider、model、execution mode、batch mode、credential 状態分類、接続確認状態、再解決結果分類、再解決時刻である。

実行テスト種別: `APIテスト`
実行段階: `実装後`
公開接点 / API 境界: phase 実行開始、provider execution adapter。
入力開始点: Running phase、実行中の provider settings 更新または未設定化。
主要結果: 実行中 phase の設定由来が非 secret 要約で説明できる。
主要観測点: phase runtime snapshot、fake provider execution 入力要約、structured log。
公開接点確認: あり。

### SCN-PSJD-006 完了、再実行、監査要約が secret を残さない

受け入れ条件:
Completed phase は provider settings 更新だけで再評価されない。
Failed phase の再実行は同じ `JOB_PHASE_RUN` を戻し、再実行開始時に最新 provider settings を再解決する。
Attempt 履歴テーブルは作らない。
job detail、phase run detail、保存要約、structured log は secret と raw payload を出さない。

実行テスト種別: `APIテスト`
実行段階: `実装後`
公開接点 / API 境界: phase 再実行、job detail、phase run detail。
入力開始点: Completed phase、Failed phase、更新後 provider settings。
主要結果: 完了済み結果は維持され、再実行は最新共通設定を参照する。
主要観測点: `JOB_PHASE_RUN` 状態、出力成果物、実行要約、直近失敗分類。
公開接点確認: あり。

## 受け入れテスト

| シナリオ | 必須受け入れ条件 | 検証入口 |
| --- | --- | --- |
| `SCN-PSJD-001` | provider settings 側だけが endpoint と credential 状態を扱う | `python3 scripts/harness/run.py --suite backend-local` |
| `SCN-PSJD-002` | Ready job は選択値だけを永続化する | `python3 scripts/harness/run.py --suite frontend-local` |
| `SCN-PSJD-003` | 不足状態と外部取得抑止を画面から確認できる | `python3 scripts/harness/run.py --suite frontend-local` |
| `SCN-PSJD-004` | 実行開始前に最新 provider settings を再解決する | `python3 scripts/harness/run.py --suite backend-local` |
| `SCN-PSJD-005` | Running phase は非 secret 要約だけを保存する | `python3 scripts/harness/run.py --suite backend-local` |
| `SCN-PSJD-006` | 完了、再実行、監査要約に secret と raw payload が出ない | `python3 scripts/harness/run.py --suite backend-local` |

## UI 設計

UI 設計は必要である。
理由は、Job Setup が扱う値から credential 参照と endpoint を外すため、既存の Job Setup 表示契約と保存契約を独立して確認する必要がある。
UI 設計は `ui-design.md` に分け、この成果物には混ぜない。

## 未決事項

人間判断が必要な未決事項はない。
Running phase 追従、provider settings revision、`credential_ref` 完全削除の候補は、今回の固定前提に従って採否を固定した。

AI 採用判断:
provider settings revision は Job 側に保持しない。
理由は、provider settings の更新履歴を保存しない既存仕様と、Job 側へ AIサービス設定値を永続所有させない前提があるためである。

AI 採用判断:
`credential_ref` 実値は Job 側から外し、phase には credential 状態分類だけを残す。
理由は、secret store 情報を Job 側 DB が所有しない前提と、実行時の非 secret 要約だけを残す前提があるためである。

## 正本化対象

人間レビュー後に、必要なら次の正本更新を別レーンで扱う。

- `docs/detail-specs/ai-provider-settings-management.md`
- `docs/detail-specs/translation-job-setup.md`
- `docs/er.md`
- `docs/diagrams/er/combined-data-model-er.puml`

## 完了判定

シナリオ設計は作成済みである。
状態は人間レビュー待ちである。
implementation-scope は作成していない。
