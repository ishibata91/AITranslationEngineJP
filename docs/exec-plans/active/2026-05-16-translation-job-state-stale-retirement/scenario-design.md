# Scenario Design: 2026-05-16-translation-job-state-stale-retirement

- `skill`: `scenario-design`
- `status`: `stopped_for_human_decision`
- `source_plan`: `./implement-lane-task-frame.md`
- `ui_source`: `N/A`
- `final_artifact_path`: `N/A`
- `topic_abbrev`: `TJSR`
- `candidate_sources`:
  - `./scenario-candidates.actor-goal.md`
  - `./scenario-candidates.lifecycle.md`
  - `./scenario-candidates.state-transition.md`
  - `./scenario-candidates.failure.md`
  - `./scenario-candidates.external-integration.md`
  - `./scenario-candidates.operation-audit.md`

## Fixed Requirements

- `must_pass_requirements`:
  - `pending` は `TRANSLATION_JOB.state`、`JOB_PHASE_RUN.state`、read model の phase state、operation summary、Wails DTO、詳細仕様の state 一覧へ正本 state として出さない。
  - `Ready` job には `JOB_PHASE_RUN` を事前作成しない。phase start が許可された時だけ対象 phase run を作成し、観測可能な最初の phase state は `Running` にする。
  - `pause`、`resume`、`retry`、`cancel` の操作可否は phase type に依存しない共通操作規則で決める。
  - read model の `CanPause`、`CanResume`、`CanRetry`、`CanCancel` は、`TranslationJobPolicy` が扱う共通操作規則と同じ state 事実から導出する。
  - `TranslationJobPolicy` の rule 名、判定履歴、`PolicyResult` は DB、DTO、repository 永続契約、read model の永続値へ出さない。
  - 保存済み `TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` が食い違う場合、読み取りだけで状態を書き換えず、危険操作を無効にする。
  - terminal job では phase run 作成、保存、readiness 更新、late response 後書きを拒否する。
  - provider raw payload、prompt 全文、翻訳本文全文、credential 実値、API key、endpoint 実値は、今回の stale 廃止で新しく保存、表示、ログ出力しない。
  - `stale_selection`、`validation_stale`、`model_selection_stale` はドメイン仕様の理由分類として残す。
- `non_goals`:
  - UI、画面文言、layout、style を変更しない。
  - DB schema、Wails 公開 DTO、新しい永続 state を追加しない。
  - `docs/exec-plans/completed/**` を変更しない。
  - product code、product test、docs 正本本文の変更可否を、このシナリオ設計だけで承認しない。
  - `JobIOService` の廃止または実体化を、AI 判断だけで確定しない。
  - active `observability-log-addition` の task-local 更新範囲を、AI 判断だけで確定しない。
  - `cancelled` fixture spelling の今回実装範囲を、AI 判断だけで確定しない。

## Scenario Candidate Coverage

正本: `./scenario-design.candidate-coverage.json`

6 種の候補成果物は存在する。
`designer` は候補生成器を再起動していない。

候補統合の結果、`pending` の正本 state 昇格は採用しない。
理由は `docs/spec.md` が `TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` を限定し、`pending` を含めていないためである。

`TranslationJobPolicy` の共通操作規則は、利用者が観測する操作可否の同一仕様として採用する。
実装方式は、`TranslationJobPolicy` を UseCase だけが呼ぶという architecture 正本に従い、implementation-scope で人間レビュー後に決める。

`JobIOService`、active `observability-log-addition`、`cancelled` fixture spelling は、人間レビューが必要な範囲として質問票へ分離する。

## Detail Requirement Coverage

正本: `./scenario-design.requirement-coverage.json`

詳細要求タイプは JSON に分離した。
`needs_human_decision` が残るため、scenario-design は完了ではなく人間回答待ちで停止する。

## Scenario Matrix

### SCN-TJSR-001 正本 state と内部一時 state を分離する

- `status`: `ready_after_human_review`
- `受け入れテスト`: `Ready` job の読み取り、phase start、既存 `pending` 観測時の安全側扱いを確認する。
- `実行者`: 翻訳 job を実行する利用者
- `開始条件`: `TRANSLATION_JOB.state` が `Ready` で、active な `JOB_PHASE_RUN` が存在しない。
- `操作`: Job Run または公開 API から対象 phase の start を要求する。
- `期待結果`: 読み取りだけでは state を変更しない。start 許可時だけ対象 phase run を作成し、観測可能な phase state は `Running` になる。
- `期待結果`: `pending` が既存データまたは内部途中値として観測されても、正本 phase state、DTO state、read model state、operation summary の成功値として返さない。
- `観測点`: job state、phase run 作成有無、phase state、operation availability、拒否理由カテゴリ、DB state の変更有無
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装後`
- `公開接点 / API 境界`: Wails backend の job / phase query と phase start command
- `入力開始点`: Ready job、phase start request、仕様外 `pending` を含む検証データ
- `主要 結果`: `pending` は正本 state へ昇格しない。危険操作は安全側へ倒れる。
- `主要観測点`: response DTO、repository snapshot、operation summary
- `公開接点確認`: あり

### SCN-TJSR-002 共通操作規則で phase 操作可否をそろえる

- `status`: `ready_after_human_review`
- `受け入れテスト`: `Running`、`Paused`、`RecoverableFailed`、terminal state の操作可否を phase type 横断で確認する。
- `実行者`: 翻訳 job を進める利用者
- `開始条件`: 対象 phase run が `Running`、`Paused`、`RecoverableFailed`、`Completed`、`Failed`、`Canceled` のいずれかで存在する。
- `操作`: `pause`、`resume`、`retry`、`cancel` の可否を query し、許可される操作は command として実行する。
- `期待結果`: `Running` では pause だけを許可する。`Paused` では resume と cancel を許可する。`RecoverableFailed` では retry を許可する。terminal job ではすべて拒否する。
- `期待結果`: read model の `CanPause`、`CanResume`、`CanRetry`、`CanCancel` は、共通操作規則と矛盾しない。
- `観測点`: action enablement、拒否理由カテゴリ、phase state、phase run id、operation result summary
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装後`
- `公開接点 / API 境界`: phase operation availability query と phase command
- `入力開始点`: phase state ごとの検証データ
- `主要 結果`: phase 固有の `canRetry`、`canResume`、`canPause`、`canCancel` は復活しない。
- `主要観測点`: response DTO、operation result、state snapshot
- `公開接点確認`: あり

### SCN-TJSR-003 状態不整合時に危険操作を拒否する

- `status`: `ready_after_human_review`
- `受け入れテスト`: job state と phase state が食い違う検証データで、読み取りと command の安全側動作を確認する。
- `実行者`: 未完了 job を確認する利用者
- `開始条件`: 保存済み `TRANSLATION_JOB.state` と現在 phase の `JOB_PHASE_RUN.state` が食い違っている。
- `操作`: 未完了一覧または Job Run を読み取り、続けて start、pause、resume、retry、cancel のいずれかを要求する。
- `期待結果`: 読み取りだけでは永続 state を変更しない。危険操作は無効になり、状態不整合を理由カテゴリとして返す。
- `観測点`: job 一覧 response、Job Run 表示対象可否、operation availability、拒否理由カテゴリ、DB state の変更有無
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装後`
- `公開接点 / API 境界`: job list query、Job Run query、phase command
- `入力開始点`: job state と phase state の不整合 fixture
- `主要 結果`: 表示確認だけで状態を修復しない。成功状態や空状態へ丸めない。
- `主要観測点`: response DTO、repository snapshot、reason category
- `公開接点確認`: あり

### SCN-TJSR-004 terminal guard と外部応答破棄を維持する

- `status`: `ready_after_human_review`
- `受け入れテスト`: terminal job に対する phase run 作成、保存、readiness 更新、late response 後書き拒否を確認する。
- `実行者`: phase 実行 usecase
- `開始条件`: job または phase が `Completed`、`Failed`、`Canceled` のいずれかである。
- `操作`: phase start、結果保存、readiness 更新、provider late response の処理を要求する。
- `期待結果`: 既存 state は変わらない。late response は保存されず、破棄事実だけが redacted summary として観測できる。
- `期待結果`: provider raw response、raw prompt、翻訳本文全文、credential 実値は保存、表示、ログへ増えない。
- `観測点`: terminal state、phase run 作成有無、result 保存有無、readiness 更新有無、redacted summary、raw payload 非露出
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装後`
- `公開接点 / API 境界`: phase command、provider adapter fake、operation summary
- `入力開始点`: terminal job、遅延 provider 応答 fake
- `主要 結果`: terminal guard が stale 廃止後も維持される。
- `主要観測点`: DB snapshot、operation summary、fake provider request count
- `公開接点確認`: あり

### SCN-TJSR-005 state 事実の保存境界を誤認しない

- `status`: `blocked_by_Q-001`
- `受け入れテスト`: `JobIOService` の扱いに応じて、状態事実の取得と保存の境界を確認する。
- `実行者`: 状態不整合を調査する運用確認者
- `開始条件`: `JobIOService` が architecture 正本に残り、実体 package が `doc.go` だけである。
- `操作`: job state、phase run state、進捗、失敗 reason category の取得または保存経路を追う。
- `期待結果`: 運用確認者は、状態事実の保存境界と状態遷移判断の境界を混同しない。
- `観測点`: architecture 正本、arch-lint component、active task-local の境界名、実体 package の有無
- `実行テスト種別`: `lower-level only`
- `実行段階`: `最終検証`
- `未決`: `Q-001` の回答が必要である。

### SCN-TJSR-006 active task-local の旧名参照を再注入させない

- `status`: `blocked_by_Q-002`
- `受け入れテスト`: active `observability-log-addition` の `StateMachine` / `JobIOService` 旧名参照を、今回更新するか残留管理するか確認する。
- `実行者`: active observability task を再開する運用確認者
- `開始条件`: active `observability-log-addition` に `StateMachine` または `JobIOService` 参照が残っている。
- `操作`: active task-local の scenario、候補、設計差分図を参照する。
- `期待結果`: 旧名参照を現在の状態境界として誤認しない。残す場合は既知の残留参照として記録される。
- `観測点`: active task-local の検索結果、残留理由、更新対象 path、更新しない path 分類
- `実行テスト種別`: `lower-level only`
- `実行段階`: `最終検証`
- `未決`: `Q-002` の回答が必要である。

### SCN-TJSR-007 ドメイン仕様の stale reason を保持する

- `status`: `ready_after_human_review`
- `受け入れテスト`: 状態関連 stale 廃止が、利用者向けまたは API 向け stale reason を削除しないことを確認する。
- `実行者`: 未完了 job 一覧と Job Run を確認する利用者
- `開始条件`: `stale_selection`、`validation_stale`、`model_selection_stale` を返す既存の理由分類がある。
- `操作`: job 一覧、Job Run、関連 API response の reason category を確認する。
- `期待結果`: `stale_selection`、`validation_stale`、`model_selection_stale` は保持される。状態関連 stale 廃止の削除対象に混ざらない。
- `観測点`: reason category、operation availability、API response、状態変更なし
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装後`
- `公開接点 / API 境界`: job list query、Job Run query、operation result
- `入力開始点`: stale reason を含む検証データ
- `主要 結果`: stale reason は空状態や成功状態と区別される。
- `主要観測点`: response DTO、reason category、state snapshot
- `公開接点確認`: あり

### SCN-TJSR-008 cancel state spelling を検索漏れさせない

- `status`: `blocked_by_Q-003`
- `受け入れテスト`: 正本 spelling の `Canceled` / `canceled` と、残留 spelling の `cancelled` を検索し、正本 state と fixture 差分を区別する。
- `実行者`: cancel 済み phase または job の状態を確認する運用確認者
- `開始条件`: `PersonaGenerationPhaseContractStub` に `cancelled` fixture spelling が残っている。
- `操作`: cancel 結果、fixture 応答、state 検索結果を確認する。
- `期待結果`: `cancelled` は正本 state として扱わない。今回含める場合は正本 spelling へそろえる。別 task に送る場合は残留参照として記録する。
- `観測点`: fixture response、contract stub、operation availability、terminal guard、検索結果
- `実行テスト種別`: `lower-level only`
- `実行段階`: `最終検証`
- `未決`: `Q-003` の回答が必要である。

## Human Review State

- `scenario-design`: 人間回答待ち。
- `ui-design`: 該当なし。UI 変更は想定しない。
- `implementation-scope`: 停止中。人間設計レビュー後にだけ作成する。

## Verification

- `scenario candidate presence`: 6 種の候補成果物を確認済み。
- `scenario candidate generator`: 起動していない。
- `manual browser check`: 未実行。UI 変更がないため不要。
- `requirement gate`: `needs_human_decision` が残るため fail になる想定で実行する。

## Open Questions

質問票正本: `./scenario-design.questions.md`

- `Q-001`: `JobIOService` を architecture 正本から外すか、別 task で実体化するか。
- `Q-002`: active `observability-log-addition` の旧名参照を今回の task-local 更新へ含めるか。
- `Q-003`: `cancelled` fixture spelling を今回の stale 廃止へ含めるか。

## Return To Implement Lane

`scenario-design` は人間回答待ちで停止する。
`implement_lane` は、人間設計レビューまたは質問回答を挟んでから、設計差分図、人間設計レビュー、implementation-scope へ進む判断を行う。
