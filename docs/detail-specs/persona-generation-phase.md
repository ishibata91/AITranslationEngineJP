# 詳細仕様: NPC ペルソナ生成フェーズ

- `upper_scenario_id`: `persona-generation-phase`
- `status`: `approved`
- `source_plan`: `docs/exec-plans/completed/persona-generation-phase/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/plan.md`, `docs/exec-plans/active/2026-05-10-translation-job-state-machine-redesign/plan.md`
- `scenario_source`: `docs/exec-plans/completed/persona-generation-phase/scenario-design.md`, `docs/exec-plans/active/2026-05-10-translation-job-state-machine-redesign/scenario-design.md`
- `ui_source`: `docs/exec-plans/completed/persona-generation-phase/ui-design.md`
- `implementation_source`: `docs/exec-plans/completed/persona-generation-phase/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/final-validation.md`
- `review_source`: `docs/exec-plans/completed/persona-generation-phase/reviewback.*.yaml`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/review-summary.md`

## 要約

NPC ペルソナ生成フェーズは、単語翻訳フェーズ完了後に開始する。
このフェーズは、NPC の原文発話、NPC 属性メタデータ、会話文脈、共通ペルソナ参照から、ジョブ内で本文翻訳フェーズが参照する persona snapshot を作る。

Job Run は、フェーズの進行、生成対象、生成結果、persona snapshot 参照状態、本文翻訳フェーズの開始可否を表示する。
provider 失敗、入力不備、保存失敗、partial state は successful Completed として扱わない。

## 対象

- 対象利用者は、翻訳ジョブを進めるユーザーと、失敗理由を確認する運用確認者である。
- 開始条件は、単語翻訳フェーズが Completed であり、ジョブが terminal state ではなく、active phase run が存在しないことである。
- 完了状態は、生成対象の persona または snapshot 参照がそろい、本文翻訳フェーズの readiness が成立することである。
- 主要データは、NPC record、translation field reference、会話文脈、共通ペルソナ参照、`PERSONA`、`PERSONA_FIELD_EVIDENCE`、`PHASE_RUN_PERSONA`、persona snapshot summary である。

## 仕様

- 生成対象 summary は、NPC count、入力種類、対象件数、common persona hit / miss、対象外理由を含む。
- NPC ペルソナ生成フェーズ開始が許可された時だけ、NPC ペルソナ生成用の `JOB_PHASE_RUN` を作成する。
- 操作可否は `JOB_PHASE_RUN.state` と共通操作規則から決める。
- NPC ペルソナ生成フェーズ固有の `canRetry`、`canResume`、`canPause`、`canCancel` は持たない。
- 共通ペルソナ hit 時は新規 `PERSONA` を作らず、ジョブの persona snapshot 参照だけを固定する。
- persona 生成は 1 NPC を 1 request unit とし、NPC 属性と会話文脈を同じ request で扱う。
- provider、model、execution mode、batch mode は Job Setup の persona 専用設定を継承する。
- phase 開始と retry は、AIサービス設定から最新 endpoint と credential 参照状態を再解決する。
- job 側 runtime snapshot は provider、model、credential 状態分類、execution mode、batch mode だけを保存する。
- valid provider output は、ジョブ内ペルソナまたは persona snapshot 参照へ自動採用する。
- 生成対象 0 件は Completed とし、対象 0 件、provider 未実行、snapshot 空を result summary に出す。

失敗時は、provider failure、invalid response、input missing、save failure を成功として保存しない。
一部 NPC 失敗時は成功分を維持し、phase は RecoverableFailed として未処理 NPC だけ retry 対象にする。

再送、再開、リトライでは同じ `JOB_PHASE_RUN` を継続する。
成功済み persona と `PHASE_RUN_PERSONA` は重複作成しない。
terminal job では persona phase start、persona save、body readiness update を拒否し、既存 state を変更しない。

本文翻訳フェーズの readiness は、persona phase Completed かつ snapshot 参照成立の両方が true の時だけ成立する。
persona 未完了、失敗、snapshot 参照不能では本文翻訳フェーズの run を作成しない。

## UI 契約由来の恒久仕様

- Job Run は current phase として `NPC ペルソナ生成` を表示する。
- Job Run は phase state、progress、target count、generated count、failed count、skipped count を表示する。
- Job Run は persona snapshot ID または snapshot digest、snapshot 参照状態、missing count、body phase readiness を表示する。
- Job Run は provider、model、execution mode、batch mode、credential 状態分類、input count、output count、短い error kind を表示する。
- Job Run 再表示時は redacted phase result summary を復元できる形で表示する。
- 長い NPC 名、provider 名、model 名、error reason、snapshot digest は desktop と mobile で表示を破綻させない。

操作可否は phase state と readiness から決まる。
開始は、term phase Completed、非 terminal job、active phase run なしの場合だけ有効にする。
pause は Running の時だけ有効にする。
resume は Paused の時だけ有効にする。
retry は RecoverableFailed かつ retryable failure の時だけ有効にする。
cancel は Paused の時だけ有効にする。
body phase 開始は persona phase Completed と snapshot 参照成立が両方 true の時だけ有効にする。

phase state、retryable、body readiness は色だけで示さない。
disabled action は理由を読み取れる text または tooltip を持つ。
progress は数値と state label を併記する。

## 保護仕様

secret、API key 平文、credential 参照実値、secret store key、endpoint、provider raw request / response、raw prompt、原文発話全文、会話文脈全文は UI、error summary、structured log、final validation summary に出さない。
operation summary は DB に永続保存せず、必要な時に状態事実から導出する。
UI と導出 summary には ID、digest、件数、evidence ref、redacted phase result summary だけを出す。

debug log に prompt または request body を出す場合でも、secret と API key は出さない。
障害調査用の要約では、provider、model、execution mode、batch mode、credential 状態分類、input count、output count、prompt digest、error kind を確認できる。

## 受け入れ根拠

- `REQ-PGP-001`: `SCN-PGP-001`, `SCN-PGP-010`
- `REQ-PGP-002`: `SCN-PGP-002`, `SCN-PGP-003`
- `REQ-PGP-003`: `SCN-PGP-003`, `SCN-PGP-007`, `SCN-PGP-009`
- `REQ-PGP-004`: `SCN-PGP-004`, `SCN-PGP-008`
- `REQ-PGP-005`: `SCN-PGP-001`, `SCN-PGP-005`
- `REQ-PGP-006`: `SCN-PGP-004`, `SCN-PGP-006`
- `REQ-PGP-007`: `SCN-PGP-007`, `SCN-PGP-008`
- `REQ-PGP-008`: `SCN-PGP-005`, `SCN-PGP-009`

検証結果は pass である。
scenario gate、backend、frontend、system test、coverage を含む最終検証は通過済みである。
5 観点レビューはすべて `review_status: no_issue`、`must_fix_open: false`、`max_level: none` である。

## 対象外

- 本文翻訳フェーズの訳文生成。
- 共通ペルソナ構築フローの実行設計。
- ジョブ内ペルソナ flush の実行設計。
