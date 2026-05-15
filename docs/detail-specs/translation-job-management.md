# 詳細仕様: 翻訳ジョブ管理

- `upper_scenario_id`: `translation-job-management`
- `status`: `approved`
- `source_plan`: `docs/exec-plans/completed/translation-job-management/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/plan.md`, `docs/exec-plans/active/2026-05-10-translation-job-state-machine-redesign/plan.md`
- `scenario_source`: `docs/exec-plans/completed/translation-job-management/scenario-design.md`, `docs/exec-plans/active/2026-05-10-translation-job-state-machine-redesign/scenario-design.md`
- `ui_source`: `docs/exec-plans/completed/translation-job-management/ui-design.md`
- `implementation_source`: `docs/exec-plans/completed/translation-job-management/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/final-validation.md`
- `review_source`: `docs/exec-plans/completed/translation-job-management/reviewback.*.yaml`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/review-summary.md`

## 要約

- 利用者は Completed 以外の翻訳 job を一覧し、次に扱う job を選べる。
- 利用者は選択した job を Job Run の表示対象にできる。
- 利用者は停止可否、削除可否、再開可否、リトライ可否、実行不可理由を確認できる。

## 対象

- 対象利用者は、作成済みの未完了翻訳 job を確認し、再開、リトライ、停止、削除の判断をしたい利用者である。
- 開始条件は、翻訳 job が作成済みであり、翻訳管理を開けることである。
- 完了状態は、未完了 job の一覧、状態、操作可否、理由カテゴリ、Job Run への導線を確認できることである。
- 主要データは `TRANSLATION_JOB`、入力出自、phase progress、AI 設定要約、credential 状態分類、reason category である。

## 仕様

- Ready、Running、Paused、RecoverableFailed、Failed、Canceled の job は未完了一覧の対象にする。
- Completed job は未完了一覧に表示しない。
- 未完了一覧、Job Run 導線、job-level terminal 判定は `TRANSLATION_JOB.state` を参照する。
- フェーズ画面の操作可否は、対象フェーズの `JOB_PHASE_RUN.state` を参照する。
- 一覧表示だけでは job 状態を変えない。
- Ready job には `JOB_PHASE_RUN` を事前作成しない。
- フェーズ開始が許可された時だけ、対象フェーズの `JOB_PHASE_RUN` を作成する。
- 同一入力データに対する job 一意制約は持たない。
- 同じ入力データに Completed、未完了、terminal の既存 job があっても新規 job を作成できる。
- 新規 job 作成時に、既存 job は削除または上書きしない。
- 同じ入力から作成された複数 job は job ID と作成日時で区別できる。
- 未完了一覧で選択した job は Job Run の表示対象になる。
- Job Run 表示だけでは Ready job を Running へ暗黙遷移させない。
- Ready job は再編集ではなく read-only の実行入口として見える。
- Running job は削除できない。削除拒否理由と停止入口を表示する。
- 停止要求中は削除できず、停止完了後に削除可否を再判定する。
- 非実行中 job を削除しても、input data と抽出 JSON 正本は残る。
- 削除成功後、対象 job は未完了一覧から外れる。
- Paused では再開入口、RecoverableFailed ではリトライ入口、現在フェーズ、進捗、実行不可理由を確認できる。
- 入力キャッシュ欠落、terminal state、状態不整合では、実行不可理由を理由カテゴリとして表示する。
- 実行不可理由の表示だけでは job 状態を変えない。
- 保存済み `TRANSLATION_JOB.state` と現在フェーズの `JOB_PHASE_RUN.state` が食い違う場合、表示だけで状態を書き換えず、危険操作を無効にする。
- 一覧読み込み失敗は空一覧にしない。
- 参照不能 job は Job Run の表示対象にしない。
- phase progress 集約不能は成功値として表示せず、危険操作を無効にする。
- provider、model、execution mode、batch mode、credential 状態分類は表示できる。
- phase start と retry は、AIサービス設定から最新 endpoint と credential 参照状態を再解決する。
- job 一覧と操作結果 summary は、endpoint、credential 参照実値、secret store key、API key 本文を表示しない。
- operation summary は DB に永続保存せず、必要な時に状態事実から導出する。
- API key 平文、credential 値、外部 provider 応答原文は UI、エラー、履歴要約に表示しない。
- Job Management は job 未作成 input を一覧へ混ぜない。
- Data Load の新規登録は input 作成だけであり、job は自動作成しない。

## 受け入れ根拠

- `SCN-TJM-001`: 未完了ジョブを一覧し操作可否を比較する。
- `SCN-TJM-002`: 同じ入力データから新しい job を作成する。
- `SCN-TJM-003`: 選択した job を Job Run の表示対象にする。
- `SCN-TJM-004`: Running job の削除を拒否し停止入口を表示する。
- `SCN-TJM-005`: Paused job の再開入口と RecoverableFailed job のリトライ入口を表示する。
- `SCN-TJM-006`: 非実行中 job を削除して入力データを保持する。
- `SCN-TJM-007`: 実行不可理由を表示し状態を変えない。
- `SCN-TJM-008`: 参照不能と集約不能を安全側に表示する。
- `SCN-TJM-009`: 保存済み AI 設定と secret 参照状態を平文なしで表示する。
- 人間設計レビューは 2026-05-06 に承認済みである。
- 最終検証は `python3 scripts/harness/run.py --suite all` 通過済みである。
- 5 観点 reviewback は再集約後にすべて `no_issue` である。

## UI 契約由来の恒久仕様

- 表示項目は未完了 job 一覧、job state、現在フェーズ、進捗、入力出自、操作可否、無効理由、AI 設定要約である。
- Completed job は未完了一覧に表示しない。
- 一覧カードは合意済み UI とし、過剰な詳細パネルや一覧の過剰表示項目を復活させない。
- stepper と Job Run 連携は、合意済み UI として維持する。
- 実行不可理由と操作可否は、一覧内の操作ボタン、無効理由、Job Run 表示対象で確認できる範囲に収める。
- 登録済みまたは警告ありの selected input だけ、Data Load から `Job Setup へ進む` 導線を表示する。
- 失敗または再構築が必要な selected input では、`Job Setup へ進む` 導線を表示しない。
- 読み込み失敗、stale selection、集約不能は空状態や成功状態と区別して表示する。
- secret は credential 状態分類だけを表示し、API key 本文を表示しない。

## 対象外

- Completed job の成果物確認。これは翻訳成果物出力で扱う。
- phase resume の実行本体。
- 実行中通信の止め方。
- 削除後の履歴保持と監査表示。
- 入力キャッシュ再構築の実行。
