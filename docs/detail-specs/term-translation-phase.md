# 詳細仕様: 単語翻訳フェーズ

- `upper_scenario_id`: `term-translation-phase`
- `status`: `approved`
- `source_plan`: `docs/exec-plans/completed/term-translation-phase/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/plan.md`
- `scenario_source`: `docs/exec-plans/completed/term-translation-phase/scenario-design.md`
- `ui_source`: `docs/exec-plans/completed/term-translation-phase/ui-design.md`
- `implementation_source`: `docs/exec-plans/completed/term-translation-phase/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/final-validation.md`
- `review_source`: `docs/exec-plans/completed/term-translation-phase/reviewback.*.yaml`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/review-summary.md`

## 要約

単語翻訳フェーズは、本文翻訳フェーズの前に用語と固有名詞の訳語を確定する。
確定した訳語は対象ジョブのジョブ内辞書へ反映し、後続フェーズの入力にする。

ユーザーは Job Run で単語翻訳フェーズを開始し、進行状況、結果、失敗理由、後続フェーズへ進める条件を確認する。
共通辞書に完全一致する語は provider request へ送らず、phase 開始時の snapshot に基づく置換対象として扱う。

## 対象

- 対象利用者は、翻訳ジョブを実行するユーザーである。
- 開始条件は、対象ジョブが Ready であり、active な単語翻訳 phase run が存在しないことである。
- 完了状態は、単語翻訳フェーズが Completed になり、ジョブ内辞書参照が成立することである。
- 主要データは、`JOB_PHASE_RUN`、`DICTIONARY_ENTRY`、`PHASE_RUN_DICTIONARY_ENTRY`、共通辞書 snapshot、provider / model / execution mode / batch mode の要約である。

## 仕様

Ready 以外のジョブ、terminal job、既存 active phase run があるジョブでは、単語翻訳フェーズを開始できない。
job は Running のまま維持し、単語翻訳フェーズの状態で完了、中断、回復可能失敗、再実行準備を区別する。
phase 開始と retry は、AIサービス設定から最新 endpoint と credential 参照状態を再解決する。
job 側 runtime snapshot は provider、model、credential 状態分類、execution mode、batch mode だけを保存する。

共通辞書は phase 開始時の snapshot で固定する。
共通辞書に完全一致する語は provider request へ含めず、共通辞書完全一致ではない語は確定済みとして保存しない。
共通辞書除外後に対象語が 0 件でも、phase result は Completed とし、provider 未実行を結果要約へ表示する。

共通辞書にない用語と固有名詞は provider へ送る。
provider 実行は 1 対象語 1 request unit を基本とし、Batch API を使う場合も batch item は 1 対象語単位にする。
provider の valid response は source term と translated term の対応を保持し、自動で確定訳語として扱う。

確定訳語は対象 job のジョブ内辞書として保存する。
`DICTIONARY_ENTRY.translation_job_id` と `PHASE_RUN_DICTIONARY_ENTRY` から、対象 job と phase run を追跡できる必要がある。
同一 job、同一 record type、同一 source term では重複 entry を作らない。
別 record type の同一 source term は別 entry として扱える。

再開、リトライ、開始再送では、同じ `JOB_PHASE_RUN` を使う。
既存 entry は維持し、未処理 term だけ provider request へ進める。
retryable failure では最新 error と progress を更新する。

provider 失敗、応答不正、保存失敗は成功扱いにしない。
invalid response、response 欠落、余分な応答、空訳語は対象語単位の failed / retryable として扱う。
保存途中失敗では partial dictionary state を成功扱いにしない。
別 provider への暗黙 fallback は行わない。

単語翻訳フェーズ未完了、失敗中、辞書参照不能の場合、後続 phase run を作成しない。
単語翻訳フェーズ完了後だけ、後続 phase の入力 summary が成立する。
terminal job への後書きは拒否する。

secret、API key 平文、復号可能な値、credential 参照実値、secret store key、endpoint、provider raw request / response、翻訳フィールド本文の全文は表示しない。
同じ情報は error summary、structured log、fake transport log、保存データにも出さない。
監査要約には provider、model、execution mode、batch mode、credential 状態分類、input count、output count、prompt version または digest、共通辞書 snapshot の digest または version を残す。

## UI 契約由来の恒久仕様

Job Run は current phase、phase state、progress、開始時刻、完了時刻、対象語件数、共通辞書 hit 件数、AI 実行対象語件数を表示する。
phase result は確定訳語件数、ジョブ内辞書反映件数、置換対象件数、未一致件数、provider / model / execution mode / batch mode / credential 状態分類の要約を表示する。

エラー時は error kind、短い理由、retryable flag、後続 phase 不可理由を表示する。
secret、API key 平文、provider raw request / response、翻訳フィールド本文の全文は表示しない。

開始操作は Ready job かつ active な単語翻訳 phase run がない時だけ有効にする。
後続 phase へ進む操作は、単語翻訳フェーズ完了と辞書参照成立後だけ有効にする。
リトライは retryable failure の時だけ有効にする。

Job Run は `idle_ready`、`running`、`empty_completed`、`completed`、`paused`、`recoverable_failed`、`blocked` の状態差分を示す。
loading 中は既存 summary を保持し、更新中であることを表示する。
phase state は色だけでなく text label で示す。
disabled button は理由を近接表示する。

desktop と mobile では、長い source term、provider 名、model 名、error reason が横にはみ出さない。
狭幅では summary を 1 列にし、件数と状態 label が折り返しても操作列を押し出さない。
table が必要な情報でも、狭幅では key-value list へ落とせる構造にする。

## 受け入れ根拠

- `SCN-TTP-001`: Ready job から単語翻訳フェーズを開始し、current phase と progress を確認する。
- `SCN-TTP-002`: 共通辞書の完全一致語を provider request から除外する。
- `SCN-TTP-003`: provider 応答を確定訳語候補として扱う。
- `SCN-TTP-004`: 確定訳語をジョブ内辞書へ反映する。
- `SCN-TTP-005`: Job Run で単語翻訳フェーズ結果を確認する。
- `SCN-TTP-006`: provider 失敗や応答不正で辞書を汚さない。
- `SCN-TTP-007`: 単語翻訳フェーズ未完了では後続フェーズへ進めない。
- `SCN-TTP-008`: phase 再送、再開、リトライで重複作成しない。
- `SCN-TTP-009`: 監査要約と redaction を確認する。

## 検証根拠

`term-translation-phase` の plan は `workflow_state: implementation-review-passed` である。
design bundle は human approved であり、`implementation-scope.md` は `human_review_status: approved` である。

最終検証では scenario gate、frontend、backend、全体検証が pass している。
Sonar は coverage 74.6%、line 75.8%、branch 64.4%、security 0、reliability 0、maintainability HIGH 0 である。

5 観点の reviewback はすべて `review_status: no_issue`、`must_fix_open: false`、`max_level: none` である。
trust-boundary reviewback は `hard_gate_result: passed` である。
`implementation_action` は `close` である。

## 対象外

- 本文翻訳フェーズの訳文生成。
- NPC ペルソナ生成フェーズの設計。
- 共通辞書管理 UI、xTranslator import、xTranslator export。
- task-local 確認用の一時成果物。
