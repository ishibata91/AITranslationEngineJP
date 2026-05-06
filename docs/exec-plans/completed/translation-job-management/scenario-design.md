# Scenario Design: translation-job-management

- `skill`: scenario-design
- `status`: draft-human-review-ready
- `source_plan`: `./plan.md`
- `ui_source`: `./ui-design.md`
- `final_artifact_path`: `docs/scenario-tests/translation-job-management.md`
- `topic_abbrev`: `TJM`
- `human_review_ready`: `true`
- `stop_reason`: `none`
- `review_note`: 人間レビュー前の draft であり、承認済みではない。
- `candidate_sources`:
  - `./scenario-candidates.actor-goal.md`
  - `./scenario-candidates.lifecycle.md`
  - `./scenario-candidates.state-transition.md`
  - `./scenario-candidates.failure.md`
  - `./scenario-candidates.external-integration.md`
  - `./scenario-candidates.operation-audit.md`

## Fixed Requirements

- `must_pass_requirements`:
  - Completed 以外の翻訳ジョブを一覧し、入力出自、状態、現在フェーズ、進捗を確認できる。
  - 一覧表示、選択、再読込だけでは job 状態を変更しない。
  - 選択した未完了 job を Job Run の表示対象にし、入力出自、現在フェーズ、進捗、操作可否を表示できる。
  - 同じ入力データから複数の翻訳ジョブを作成できる。
  - 既存 job がある入力を選んでも、新規 job 作成を禁止しない。
  - Running job は削除できず、停止操作または停止不可理由を表示する。
  - Running ではない未完了 job は、入力データと入力出自を残したまま削除できる。
  - Paused と RecoverableFailed は再開入口を表示し、再開不可の時は理由を表示する。
  - 入力キャッシュ欠落、terminal state（完了扱いまたは回復不能な状態）、状態不整合は、再開不可理由または安全側表示として観測できる。
  - API key 平文、credential 値、外部 provider 応答原文、過剰な入力本文や翻訳本文は一覧、詳細、履歴要約に出さない。
  - paid な real AI API を scenario validation の前提にしない。
- `non_goals`:
  - Completed job の成果物確認、完了履歴画面、再出力導線はこの task で確定しない。
  - 入力キャッシュ再構築の実行 UI はこの task で作らない。
  - 実際の phase 実行、翻訳生成、provider SDK 実装は含めない。
  - docs 正本、product code、product test、implementation-scope は扱わない。

## Scenario Candidate Coverage

正本: `./scenario-design.candidate-coverage.json`

6 件の candidate artifact は揃っている。
candidate id は generator 間で重複しているため、coverage JSON では `generator:source_candidate_id` を一意 key として扱う。

- `adopted`: 13 件。
- `merged`: 31 件。
- `rejected`: 9 件。
- `needs_human_decision`: 0 件。
- `conflicted`: 0 件。

未解決 conflict は 0 件である。
Q-TJM-001 から Q-TJM-004 は回答済み判断として固定した。
Q-TJM-003 は再開実行側へ送る。
Q-TJM-004 は翻訳実行側へ送る。
Q-TJM-004 は後続の `tasks/` タスク化が必要である。

## Detail Requirement Coverage

正本: `./scenario-design.requirement-coverage.json`

詳細要求タイプは sidecar JSON に分離する。
`needs_human_decision` は残っていない。
この scenario design は人間レビュー待ちの draft とする。

### `REQ-TJM-001` 未完了ジョブ一覧と状態集約を表示する

- `source_requirement`: Completed 以外の job を一覧し、job ごとの入力出自、状態、現在フェーズ、進捗、停止可否、再開可否、削除可否を確認できる。
- `requirement_kind`: display
- `needs_human_decision`: なし
- `fixed_decisions`: Completed は未完了一覧から除外する。Failed と Canceled は Completed ではないため、未完了管理内では terminal state（完了扱いまたは回復不能な状態）の再開不可理由を表示する。

### `REQ-TJM-002` 同じ入力データから複数 job を作成できる

- `source_requirement`: 人間判断により、同じファイルまたは同じ入力データから複数の翻訳ジョブを作成できることにする。
- `requirement_kind`: workflow
- `needs_human_decision`: なし
- `fixed_decisions`: 同じ入力データから新規 job 作成を許可する。Completed job、未完了 job、terminal job が同じ入力データに紐づいていても、新しい job を作成できる。実装では同一入力に対する job 一意制約を削除する。
- `human_decision`: 2026-05-06 に「同じファイルのジョブを何個でも作れる」方針へ変更した。

### `REQ-TJM-003` 状態別に停止、再開、削除の可否を制御する

- `source_requirement`: Running job は削除ではなく停止操作を実行でき、停止または中断後に削除可否を再判定する。
- `requirement_kind`: operation
- `needs_human_decision`: なし
- `fixed_decisions`: Running でない job の停止操作は拒否する。表示と選択だけでは状態を変えない。Running job では削除を拒否し、停止要求中、停止失敗、Paused への収束後の削除可否再判定を表示する。

### `REQ-TJM-004` 非実行中 job を削除して入力データを残す

- `source_requirement`: 実行中ではない job は入力データを残したまま削除できる。入力ファイルの出自を失わずに保持する。
- `requirement_kind`: persistence
- `needs_human_decision`: なし
- `fixed_decisions`: Running は削除不可である。非実行中 job の削除成功後は、DB に保持している job 以下の情報をすべて削除し、input data と抽出 JSON 正本は残す。ローカルアプリに監査はなく、削除後の履歴保持は不要である。

### `REQ-TJM-005` 再開入口と再開不可理由を表示する

- `source_requirement`: Paused または RecoverableFailed の job では再開入口を確認できる。再開不可の job では入力キャッシュ欠落や terminal state（完了扱いまたは回復不能な状態）などの理由を確認できる。
- `requirement_kind`: external_integration
- `needs_human_decision`: なし
- `fixed_decisions`: 入力キャッシュ欠落、terminal state（完了扱いまたは回復不能な状態）、状態不整合は再開不可理由として表示する。Job Management では保存済み AI 設定要約だけを表示し、再開直前の外部設定確認は再開を実行する側で扱う。Job Management 内で入力キャッシュ再構築は実行しない。

### `REQ-TJM-006` 保存禁止情報を一覧、詳細、履歴に出さない

- `source_requirement`: APIKey は暗号化して保存する。translation-job-setup は API key 平文を表示または保存要約に出さない。
- `requirement_kind`: security
- `needs_human_decision`: なし
- `fixed_decisions`: credential は値ではなく参照状態として表示する。provider 応答原文や過剰な入力本文は監査要約へ出さない。

### `REQ-TJM-007` 外部 API 費用なしで管理境界を検証する

- `source_requirement`: paid な real AI API を scenario validation の前提にしない。
- `requirement_kind`: non_functional
- `needs_human_decision`: なし
- `fixed_decisions`: user-facing provider list は real provider のまま保ち、外部 request または SDK transport だけを fake に差し替える。

## Human Decision Questionnaire

正本: `./scenario-design.questions.md`

未回答質問はない。

- `Q-TJM-001`: 回答済み。同じ入力データから複数 job 作成を許可する。
- `Q-TJM-002`: 回答済み。非実行中 job の削除では job 以下の DB 情報を削除し、入力データと抽出 JSON を残す。
- `Q-TJM-003`: 回答済み。管理画面では保存済み AI 設定要約と再開不可理由を表示し、API 設定が今も使えるかの直前確認は再開実行側の後続タスクで扱う。
- `Q-TJM-004`: 回答済み。管理画面では停止要求中、停止失敗、Paused への収束、削除可否再判定を表示し、外部通信の止め方と停止要求後の不整合防止は翻訳実行側の後続タスクで扱う。

## Risks

- `implementation_risks`:
  - 同一入力の job 一意制約を削除する必要があるため、translation-job-setup の過去設計と実装制約に差分が出る。
  - 停止要求中、停止失敗、Paused への収束後の削除可否再判定を混同すると、Running job の削除保護が弱くなる。
  - 再開直前の外部設定確認は本 task の受け入れ条件ではないため、後続の再開実行側で確認境界を落とさないようにする必要がある。
  - 外部通信の停止方式は本 task の受け入れ条件ではないため、後続の翻訳実行側で停止方式、遅延応答破棄、停止要求後の不整合防止を固定する必要がある。
  - 停止機能が現状ない場合、翻訳実行側の `tasks/` タスク化が必要である。
- `test_data_risks`:
  - Running、Paused、RecoverableFailed、Failed、Canceled、Completed を分けた fixture が必要である。
  - 入力キャッシュ欠落、stale selection、phase progress 集約不能、削除失敗、停止失敗を別 fixture にする必要がある。
  - fake transport と fake secret store を使い、paid API と secret 平文に依存しない検証が必要である。

## Rules

- ケース ID は `SCN-TJM-NNN` 形式にする。
- Markdown table は使わず、1 ケースごとの縦型ブロックで書く。
- 受け入れテストは全ケースで先に固定する。
- `実行テスト種別` は `APIテスト | UI人間操作E2E | lower-level only` に固定する。
- `実行段階` は `実装後 | final validation` に固定する。
- `期待結果` は観測可能な結果にする。
- human review 承認前であるため、下の scenario matrix は未承認 draft とする。

## Scenario Matrix

回答済み判断と後続送り境界を反映した draft である。
人間レビュー承認後に implementation-scope の根拠として扱う。

### SCN-TJM-001 未完了ジョブを一覧し操作可否を比較する

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: Completed 以外の job を一覧し、状態、現在フェーズ、進捗、操作可否を確認する。
- `受け入れ条件`: Ready、Running、Paused、RecoverableFailed、Failed、Canceled が一覧対象になる。Completed は表示されない。
- `事前条件`: Completed 以外の job と Completed job を含む fixture がある。
- `public_seam_or_api_boundary`: 未完了ジョブ一覧取得 boundary。詳細 API 名は human review 後の implementation-scope で固定する。
- `入力開始点`: 翻訳管理 UI。
- `主要 outcome`: 利用者が次に扱う job を選べる。
- `開始操作`: 翻訳管理を開く。
- `入力方法`: 既存 job fixture。
- `主要操作列`: 一覧表示、状態確認、操作可否確認。
- `手順`:
  1. 翻訳管理を開く。
  2. 未完了ジョブ一覧を確認する。
  3. 状態、現在フェーズ、進捗、停止可否、再開可否、削除可否を確認する。
- `期待結果`:
  1. Completed 以外の job が表示される。
  2. Completed job は未完了一覧に表示されない。
  3. 一覧表示だけでは job 状態が変わらない。
- `観測点`: 未完了ジョブ一覧、job state、phase progress、input provenance。
- `UI-visible outcome`: 入力出自、状態、現在フェーズ、進捗、操作可否が表示される。
- `fake_or_stub`: fixed job fixture、temp DB。
- `責務境界メモ`: Completed job は未完了一覧には表示しない。同じ入力データから新規 job を作成できるため、Completed job の入口はこの task の blocker ではない。

### SCN-TJM-002 同じ入力データから新しい job を作成する

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装後`
- `観点`: 同じ入力データに既存 job があっても、新しい翻訳ジョブを作成できる。
- `受け入れ条件`: 同一入力データに対する job 一意制約がなく、Completed、未完了、terminal の既存 job があっても新規 job を作成できる。
- `事前条件`: 同じ入力データに紐づく既存 job fixture がある。
- `public_seam_or_api_boundary`: job create boundary。詳細 API 名は human review 後の implementation-scope で固定する。
- `入力開始点`: Job Setup または入力データ選択 UI。
- `主要 outcome`: 同じ入力データに複数 job が紐づく。
- `開始操作`: 既存 job がある入力データを選ぶ。
- `入力方法`: UI または API で同じ入力データを指定する。
- `主要操作列`: 入力選択、新規 job 作成、job 一覧または永続化結果確認。
- `期待結果`:
  1. 新しい job が作成される。
  2. 既存 job は削除または上書きされない。
  3. job ごとの入力出自は同じ入力データを指す。
- `観測点`: job count、input data ID、job ID、作成日時。
- `UI-visible outcome`: 同じ入力から作成された複数 job を区別できる。
- `fake_or_stub`: same input multi job fixture、temp DB。
- `責務境界メモ`: 実装では同一入力に対する job 一意制約を削除する。

### SCN-TJM-003 選択した job を Job Run の表示対象にする

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: 一覧で選択した job を Job Run の表示対象にし、表示だけで実行状態へ遷移させない。
- `受け入れ条件`: Ready job を選択しても Running へ暗黙遷移しない。Job Run は選択 job の入力出自、現在フェーズ、進捗を表示する。
- `事前条件`: Ready、Paused、RecoverableFailed の job fixture がある。
- `public_seam_or_api_boundary`: 選択中 job 表示 boundary。詳細 API 名は implementation-scope で固定する。
- `入力開始点`: 未完了ジョブ一覧。
- `主要 outcome`: Job Run の表示対象が選択 job と一致する。
- `開始操作`: 未完了ジョブ一覧で job を選択する。
- `入力方法`: UI の一覧行選択。
- `主要操作列`: job 選択、Job Run 表示、状態確認。
- `手順`:
  1. 未完了ジョブ一覧で Ready job を選択する。
  2. Job Run 表示へ進む。
  3. 入力出自、現在フェーズ、進捗、再開または開始入口を確認する。
- `期待結果`:
  1. Job Run の表示対象 job ID は選択した job と一致する。
  2. 表示だけでは job 状態が変わらない。
  3. Ready job は再編集ではなく read-only の実行入口として見える。
- `観測点`: Job Run 表示、selected job ID、job state。
- `UI-visible outcome`: 選択中 job の状態と現在フェーズが表示される。
- `fake_or_stub`: fixed job fixture。
- `責務境界メモ`: 実行開始そのものは phase 実行 task の責務である。

### SCN-TJM-004 Running job の削除を拒否し停止入口を表示する

- `分類`: 主要失敗系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装後`
- `観点`: Running job は削除できず、停止操作へ寄せる。
- `受け入れ条件`: Running では削除が拒否される。停止要求中、停止失敗、Paused への収束後の削除可否再判定を確認できる。
- `事前条件`: Running job と running phase fixture がある。
- `public_seam_or_api_boundary`: stop / delete guard boundary。詳細 API 名は implementation-scope で固定する。
- `入力開始点`: Running job fixture。
- `主要 outcome`: job と input data が壊れず、削除不可理由が観測できる。
- `開始操作`: Running job に削除または停止を実行する。
- `入力方法`: UI または公開接点。
- `主要操作列`: 削除要求、拒否理由確認、停止入口確認。
- `手順`:
  1. Running job に削除を試みる。
  2. 削除拒否理由を確認する。
  3. 停止入口と停止後再判定の表示を確認する。
- `期待結果`:
  1. 削除は成功しない。
  2. input data は残る。
  3. 停止要求中は削除できず、停止完了後に削除可否が再判定される。
- `観測点`: API response、job state、input data reference、UI operation state。
- `UI-visible outcome`: 削除不可理由と停止入口が表示される。
- `fake_or_stub`: running phase fixture、fake transport。
- `責務境界メモ`: 実行中通信の止め方は翻訳実行側で扱う。

### SCN-TJM-005 Paused または RecoverableFailed job の再開入口を表示する

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: 中断後または回復可能失敗後の job で再開入口と現在フェーズを表示する。
- `受け入れ条件`: Paused と RecoverableFailed で再開入口、現在フェーズ、進捗、再開不可理由を確認できる。
- `事前条件`: Paused と RecoverableFailed の job fixture がある。
- `public_seam_or_api_boundary`: resume entry boundary。詳細 API 名は implementation-scope で固定する。
- `入力開始点`: 未完了ジョブ一覧または Job Run。
- `主要 outcome`: 利用者が再開できるか、なぜ再開できないかを判断できる。
- `開始操作`: Paused または RecoverableFailed job を選択する。
- `入力方法`: UI の一覧行選択。
- `主要操作列`: job 選択、再開入口確認、理由確認。
- `手順`:
  1. Paused job を選択する。
  2. RecoverableFailed job を選択する。
  3. 再開入口と理由表示を確認する。
- `期待結果`:
  1. 再開入口が表示される。
  2. 入力キャッシュ欠落または terminal state（完了扱いまたは回復不能な状態）では再開不可理由が表示される。
  3. 保存済み AI 設定要約は確認できるが、再開直前の外部設定確認は管理画面の受け入れ条件にしない。
- `観測点`: Job Run 表示、resume availability、reason category。
- `UI-visible outcome`: 再開可否、現在フェーズ、理由カテゴリが表示される。
- `fake_or_stub`: paused job fixture、recoverable failed fixture、cache missing fixture。
- `責務境界メモ`: phase resume の実行本体は本 task で確定しない。

### SCN-TJM-006 非実行中 job を削除して入力データを保持する

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装後`
- `観点`: Running ではない未完了 job を削除し、input data と抽出 JSON 正本を残す。
- `受け入れ条件`: 削除成功後、対象 job は未完了一覧から外れる。input data と入力出自は残る。
- `事前条件`: Ready、Paused、RecoverableFailed、Failed、Canceled の job fixture がある。
- `public_seam_or_api_boundary`: non-running job delete boundary。詳細 API 名は implementation-scope で固定する。
- `入力開始点`: 非実行中 job fixture。
- `主要 outcome`: job を消しても入力データを再利用できる。
- `開始操作`: 非実行中 job に削除を実行する。
- `入力方法`: UI または公開接点。
- `主要操作列`: 削除確認、削除実行、一覧再取得、input data 確認。
- `手順`:
  1. 非実行中 job を選ぶ。
  2. 削除を実行する。
  3. 未完了一覧と input data を確認する。
- `期待結果`:
  1. job は未完了一覧から外れる。
  2. input data と入力出自は残る。
  3. job 以下の DB 情報は削除され、input data と抽出 JSON は残る。
- `観測点`: delete response、job list、input data list、input provenance。
- `UI-visible outcome`: 削除結果と入力保持が確認できる。
- `fake_or_stub`: non-running job fixture、temp DB。
- `責務境界メモ`: 削除後の履歴保持と監査表示は本 task に含めない。

### SCN-TJM-007 再開不可理由を表示し状態を変えない

- `分類`: 主要失敗系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: 入力キャッシュ欠落、terminal state（完了扱いまたは回復不能な状態）、状態不整合の job で再開不可理由を表示する。
- `受け入れ条件`: 再開不可理由が理由カテゴリとして表示され、表示だけでは job 状態が変わらない。
- `事前条件`: cache missing、Failed、Canceled、state projection inconsistent fixture がある。
- `public_seam_or_api_boundary`: resume guard and reason projection boundary。
- `入力開始点`: 未完了ジョブ一覧または Job Run。
- `主要 outcome`: 利用者が再開できない理由を確認できる。
- `開始操作`: 再開不可 job を選択する。
- `入力方法`: UI の一覧行選択。
- `主要操作列`: job 選択、理由表示、再開入口の無効状態確認。
- `手順`:
  1. 入力キャッシュ欠落 job を選択する。
  2. Failed または Canceled job を選択する。
  3. 再開不可理由を確認する。
- `期待結果`:
  1. 入力キャッシュ欠落が理由として表示される。
  2. terminal state（完了扱いまたは回復不能な状態）が理由として表示される。
  3. job と input data reference は変わらない。
- `観測点`: reason category、job state、input cache state。
- `UI-visible outcome`: 再開不可理由と次に確認すべき場所が表示される。
- `fake_or_stub`: cache missing fixture、terminal job fixture。
- `責務境界メモ`: 入力キャッシュ再構築の実行は別導線で扱う。

### SCN-TJM-008 参照不能と集約不能を安全側に表示する

- `分類`: 主要失敗系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装後`
- `観点`: 一覧選択後の参照不能、一覧読み込み失敗、phase progress 集約不能を空状態や成功値と混同しない。
- `受け入れ条件`: 参照不能 job は Job Run の表示対象にしない。読み込み失敗は空一覧にしない。集約不能進捗は成功値として表示しない。
- `事前条件`: stale selection、list load failure、progress projection failure fixture がある。
- `public_seam_or_api_boundary`: job projection boundary。
- `入力開始点`: job list read / selected job read boundary。
- `主要 outcome`: 別 job の誤表示や誤った成功表示を防ぐ。
- `開始操作`: 一覧取得、job 選択、詳細取得。
- `入力方法`: fixture と公開接点。
- `主要操作列`: 読み込み失敗、stale selection、集約不能を確認する。
- `手順`:
  1. 一覧読み込み失敗を発生させる。
  2. 選択後に対象 job を参照不能にする。
  3. phase progress 集約不能の job を表示する。
- `期待結果`:
  1. 読み込み失敗は空一覧ではなく error として表示される。
  2. 参照不能 job は Job Run 表示対象にならない。
  3. 集約不能進捗は安全側表示になり、危険操作は無効になる。
- `観測点`: API response、UI error state、selected job ID。
- `UI-visible outcome`: 再試行、再読込、理由カテゴリが表示される。
- `fake_or_stub`: stale selection fixture、projection failure fixture。
- `責務境界メモ`: state projection の内部実装は implementation-scope で固定する。

### SCN-TJM-009 保存済み AI 設定と secret 参照状態を平文なしで表示する

- `分類`: セキュリティ
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: provider、model、execution mode、credential 参照状態を表示し、secret 本体は出さない。
- `受け入れ条件`: API key 平文、credential 値、外部 provider 応答原文は UI、エラー、履歴要約に表示されない。
- `事前条件`: 保存済み AI 設定、fake secret store、redaction fixture がある。
- `public_seam_or_api_boundary`: AI setting summary boundary。
- `入力開始点`: Job Run または選択中 job detail。
- `主要 outcome`: 利用者は実行設定の参照状態だけを確認できる。
- `開始操作`: 保存済み AI 設定を持つ job を開く。
- `入力方法`: UI の一覧行選択。
- `主要操作列`: job 選択、AI 設定要約確認、エラー表示確認。
- `手順`:
  1. 保存済み AI 設定を持つ job を開く。
  2. provider、model、execution mode、credential 参照状態を見る。
  3. error summary と履歴要約に secret が出ないことを確認する。
- `期待結果`:
  1. provider と model は表示される。
  2. credential は参照状態だけが表示される。
  3. API key 平文は表示されない。
- `観測点`: UI text、console/error summary、structured log。
- `UI-visible outcome`: 保存済み credential があるかだけが表示される。
- `fake_or_stub`: fake secret store、redaction assertion。
- `責務境界メモ`: API key 保存 UI は本 task に含めない。

## Acceptance Checks

- `REQ-TJM-001`: `SCN-TJM-001`, `SCN-TJM-008`
- `REQ-TJM-002`: `SCN-TJM-002`
- `REQ-TJM-003`: `SCN-TJM-004`
- `REQ-TJM-004`: `SCN-TJM-006`
- `REQ-TJM-005`: `SCN-TJM-005`, `SCN-TJM-007`
- `REQ-TJM-006`: `SCN-TJM-009`
- `REQ-TJM-007`: 全 scenario の fake / stub 方針

## Validation Commands

- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/translation-job-management/scenario-design.md --coverage docs/exec-plans/active/translation-job-management/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/active/translation-job-management/scenario-design.candidate-coverage.json --json`
- `python3 scripts/harness/run.py --suite scenario-gate`

## Open Questions

- 未回答の人間質問はない。
- Q-TJM-003 と Q-TJM-004 は `./scenario-design.questions.md` の「後続タスクへ送る回答済み事項」に残す。
