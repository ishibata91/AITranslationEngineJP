# Scenario Design: 2026-05-10-translation-job-state-machine-redesign

- `skill`: `scenario-design`
- `status`: `ready-for-human-design-review`
- `source_plan`: `./plan.md`
- `ui_source`: `N/A`
- `final_artifact_path`: `docs/scenario-tests/translation-job-state-machine-redesign.md`
- `topic_abbrev`: `TJSM`
- `candidate_coverage`: `./scenario-design.candidate-coverage.json`
- `requirement_coverage`: `./scenario-design.requirement-coverage.json`
- `questions`: `./scenario-design.questions.md`

## 判断結果

シナリオ設計の未回答質問は解消済みである。
`scenario-design.requirement-gate.md` は pass である。

採用シナリオは、人間回答を反映して確定する。
本成果物では、候補統合後の受け入れ観点、入出力境界、質問票、ゲート結果を固定する。

## 根拠参照

- `./plan.md`: 状態遷移規則と `JobIOService` への責務分離、状態正本論点、禁止事項。
- `./scenario-candidates.*.md`: 6 観点の候補母集団。
- `docs/spec.md`: `Draft`、`Ready`、`Running`、`Paused`、`RecoverableFailed`、`Completed`、`Failed`、`Canceled` の job state。
- `docs/er.md`: ジョブ状態は `JOB_PHASE_RUN` 群から集約するというデータモデル方針。
- `docs/architecture.md`: 既存正本では `StateMachine` が状態遷移規則だけを保持し、`JobIOService` が job 状態の取得と保存だけを扱う。
- `docs/detail-specs/*.md`: phase 開始条件、retry、resume、cancel、terminal job、output readiness の詳細仕様。

## 必須要件

- `translationjobpolicy` は DB を読まない pure rule とする。
- `translationjobpolicy` は、共通操作規則を先に評価する。
- `translationjobpolicy` は、`start` の時だけ phase 別開始前提を評価する。
- `translationjobpolicy` の判断結果、rule 名、判定履歴は DB に永続化しない。
- `PolicyResult` は UseCase 内だけで消費する一時値とする。
- `retry`、`resume`、`pause`、`cancel` の可否は phase type で分けない。
- `JobIOService` は状態の取得と保存だけを扱い、遷移可否を判断しない。
- `JobIOService` は確定済みの job / phase run 状態と、仕様で保存対象にした安全な状態事実だけを保存する。
- Job Run 表示だけでは `Ready` job を `Running` へ暗黙遷移させない。
- active な `JOB_PHASE_RUN` がある job では、重複した phase run 作成を拒否する。
- phase の retry、resume、開始再送は、同じ `JOB_PHASE_RUN` を継続する。
- phase type で分ける対象は、開始前提データ、完了判定、呼び出す service method だけにする。
- terminal job では、phase run 作成、保存、readiness 更新、late response 後書きを拒否する。
- provider failure、invalid response、保存失敗、検証失敗は successful `Completed` として扱わない。
- secret、API key、credential 参照実値、endpoint、provider raw payload、prompt 全文、翻訳本文全文を UI、DTO、ログ、fake transport log に出さない。

## 非対象

- プロダクトコードとプロダクトテストの変更。
- docs 正本の変更。
- 実画面確認。
- UI 実装。
- job 削除後の履歴保持と監査表示の新規仕様化。
- paid real API を前提にしたシステムテスト。

## 状態語彙の衝突

### 衝突 1: job state の正本

`docs/spec.md` は job state を直接列挙している。
一方で、`docs/er.md` はジョブ状態を `JOB_PHASE_RUN` 群から集約すると定義している。

回答: 各フェーズ画面の操作可否は、現在フェーズの `JOB_PHASE_RUN.state` を正本にする。
回答: 大枠の一覧、導線、ジョブ全体の表示は、`TRANSLATION_JOB.state` を正本にする。
影響: 画面責務ごとに参照する状態を分けるため、ページ要求が状態値を直接変更しない構造が必要になる。
回答: 保存済み `TRANSLATION_JOB.state` と現在フェーズの `JOB_PHASE_RUN.state` が食い違う場合は、表示だけで状態を書き換えず、危険操作を無効化する。
扱い: `Q-TJSM-001` は回答済みとして扱う。

### 衝突 2: `Ready` と `pending`

`translation-job-management.md` は `Ready` job を read-only の実行入口として扱う。
過去の不整合候補と今回の候補は、`Ready` job と `pending` phase run の混在を危険状態として扱っている。

回答: `Ready` job には `JOB_PHASE_RUN` を事前作成しない。
影響: 開始要求が許可された時だけ、対象フェーズの `JOB_PHASE_RUN` を作る。
扱い: `Q-TJSM-002` は回答済みとして扱う。

### 衝突 3: `RecoverableFailed -> Ready`

`docs/spec.md` は `RecoverableFailed --> Ready : 再実行準備` を持つ。
一方で、`docs/er.md` と phase 詳細仕様は、再実行を同じ `JOB_PHASE_RUN` の状態を戻す扱いにしている。

回答: `RecoverableFailed -> Ready` は廃止する。
回答: retry と resume は同じ `JOB_PHASE_RUN` を継続する。
影響: `Ready` へ戻して新規開始に見せる経路は作らない。
扱い: `Q-TJSM-003` は回答済みとして扱う。

### 衝突 4: cancel と resume の粒度

`docs/spec.md` は `Ready` と `Paused` から `Canceled` へ進む経路を持つ。
`body-translation-phase.md` は cancel を body phase `Paused` の時だけ有効にする。
`persona-generation-phase.md` は `Paused` または `RecoverableFailed` から resume 可能とするが、body phase は `Paused` からだけ resume 可能とする。

回答: `Ready` cancel は job-level 操作として残す。
回答: phase 開始後の cancel は、`Paused` の対象フェーズからだけ許可する。
回答: resume は `Paused` だけに許可し、retry は `RecoverableFailed` だけに許可する。
回答: `retry`、`resume`、`pause`、`cancel` の可否は phase type で分けない。
影響: job-level 操作と phase-level 操作を分けるため、Running からの直接 cancel は作らない。
影響: persona phase だけ resume 条件を変える分岐は作らない。
扱い: `Q-TJSM-004` と `Q-TJSM-005` は回答済みとして扱う。

### 衝突 5: phase 別 ruleset

各 phase 詳細仕様は、開始、再開、retry、cancel の条件を phase ごとに持つ。
一方で、操作可否の大半は `JOB_PHASE_RUN.state` と要求イベントだけで決まる。

回答: phase ごとの `canRetry`、`canResume`、`canPause`、`canCancel` は作らない。
回答: `translationjobpolicy` は、共通操作規則と phase 別開始前提を分ける。
回答: phase 別に残す対象は、開始前提データ、完了判定、呼び出す service method だけにする。
影響: ルール表は phase type を主キーにしない。`start` だけが target phase の prerequisite を参照する。
扱い: `Q-TJSM-008` は回答済みとして扱う。

## translationjobpolicy 入出力

### 入力候補

- job ID。
- job state の候補値。
- `JOB_PHASE_RUN` 群の集約結果。
- current phase type。
- 操作種別。
- 開始対象 phase type。
- phase 別開始前提の判定結果。
- terminal 判定。
- active phase run 有無。
- retryable 判定。
- 同一再送判定。

### 出力候補

- `allowed` または `rejected`。
- 共通操作規則による許可または拒否。
- `start` の時だけ phase 別開始前提による許可または拒否。
- 状態不変または遷移後 job state 候補。
- 状態不変または遷移後 phase state 候補。
- 継続する `JOB_PHASE_RUN` id。
- 作成してよい phase run type。
- 呼び出してよい service method の種類。
- 保存してよい状態事実の候補。
- 危険操作の有効可否。
- UI 表示向け reason category。
- 導出してよい redacted summary 種別。

### 出力境界

`translationjobpolicy` は操作可否、拒否理由、状態遷移結果を返す。
`translationjobpolicy` の出力は UseCase 内の一時判断であり、DB に保存しない。
`translationjobpolicy` は phase service を呼ばない。
`translationjobpolicy` は phase type 別の retry / resume 可否を持たない。
`translationjobpolicy` は DB 取得、保存、provider 呼び出し、runtime event 発火を扱わない。
`translationjobpolicy` は UseCase だけが呼び出す。

### 共通操作規則

- `Running` phase run だけを `pause` できる。
- `Paused` phase run だけを `resume` できる。
- `RecoverableFailed` phase run だけを `retry` できる。
- terminal job では状態を変える操作を拒否する。
- active phase run がある時は、新しい phase run を作らない。

### phase 別開始前提

- 単語翻訳 phase は、入力データと辞書生成対象を参照できる時だけ開始できる。
- NPC ペルソナ生成 phase は、単語翻訳 phase の完了結果を参照できる時だけ開始できる。
- 本文翻訳 phase は、persona snapshot と翻訳対象 field を参照できる時だけ開始できる。
- phase 別開始前提は、`start` 以外の操作可否へ使わない。

## JobIOService 境界

### 扱うこと

- job と phase run の取得。
- phase progress 集約に必要な永続化データの読み出し。
- UseCase が確定した状態変更の保存。
- `JOB_PHASE_RUN` の作成または継続対象の保存。
- 仕様で保存対象にした runtime snapshot の安全な値の保存。
- 保存失敗を成功遷移へ変換しないこと。

### 扱わないこと

- 遷移可否の判断。
- policy の判断結果、rule 名、判定履歴の保存。
- operation summary の DB 永続保存。
- terminal guard の判断。
- provider response validation。
- secret 実値の取得または露出。
- UI 表示文言の決定。
- Wails runtime event の採用可否。

## 候補統合後の受け入れ観点

### SCN-TJSM-001 Ready 作成と表示 no-op

- 受け入れ条件: job 作成後、対象 job は `Ready` として観測できる。Job Run 表示だけでは `Running` へ変わらない。
- 主な候補: actor-goal `001`、lifecycle `001`、state-transition `001`、`002`。
- 実行テスト種別: `APIテスト`。
- 実行段階: `実装後`。
- 主要観測点: job summary、phase run 有無、操作可否、secret 非露出。
- 状態: 確定。

### SCN-TJSM-002 phase 順序と開始 guard

- 受け入れ条件: 単語翻訳、NPC ペルソナ生成、本文翻訳は、前段完了、参照成立、active phase run 不在、非 terminal job の時だけ開始できる。
- 主な候補: actor-goal `002`、`003`、`004`、lifecycle `002`、`003`、`004`、state-transition `003`、`004`、`005`、failure `010`。
- 実行テスト種別: `APIテスト`。
- 実行段階: `実装後`。
- 主要観測点: 共通操作規則の allowed / rejected、phase 別開始前提の reason category、phase run 件数。
- 状態: 確定。

### SCN-TJSM-003 pause、resume、retry、cancel

- 受け入れ条件: pause、resume、retry、cancel は phase type に依存せず、`JOB_PHASE_RUN.state` と job terminal 判定で許可または拒否する。
- 主な候補: actor-goal `005`、`006`、lifecycle `005`、`006`、`007`、`008`、state-transition `006`、`007`、`008`、`011`、`014`。
- 実行テスト種別: `APIテスト`。
- 実行段階: `実装後`。
- 主要観測点: 同じ `JOB_PHASE_RUN` id、共通操作規則の reason category、retryable flag、重複作成なし、cancel 無効理由。
- 状態: 確定。

### SCN-TJSM-004 完了、対象 0 件、output readiness

- 受け入れ条件: term と persona の対象 0 件は phase `Completed` として扱うが、job 全体を `Completed` にしない。body phase `Completed` と field result 整合、output status 整合が成立した時だけ output readiness を true にする。
- 主な候補: actor-goal `009`、lifecycle `009`、state-transition `010`、`016`、external-integration `002`、failure `011`。
- 実行テスト種別: `APIテスト`。
- 実行段階: `実装後`。
- 主要観測点: provider 呼び出し回数、phase state、job-level `Completed`、output readiness、artifact action disabled reason。
- 状態: 確定。

### SCN-TJSM-005 削除安全性と input 保持

- 受け入れ条件: Running job と停止要求中 job の削除は拒否する。非実行中 job を削除しても input data と抽出 JSON 正本は保持する。
- 主な候補: actor-goal `008`、lifecycle `010`、state-transition `012`、failure `008`、operation-audit `004`、`005`。
- 実行テスト種別: `APIテスト`。
- 実行段階: `実装後`。
- 主要観測点: 削除拒否理由、停止入口、削除後一覧、input data 残存。
- 状態: 確定。

### SCN-TJSM-006 状態不整合と集約不能の安全側表示

- 受け入れ条件: 保存済み `TRANSLATION_JOB.state` と対象フェーズの `JOB_PHASE_RUN.state` が食い違う場合、表示だけで状態を書き換えず、危険操作を無効化する。
- 主な候補: actor-goal `010`、state-transition `015`、failure `001`、`002`、`009`、operation-audit `002`、`006`。
- 実行テスト種別: `APIテスト`。
- 実行段階: `実装後`。
- 主要観測点: state inconsistent reason、aggregate unavailable reason、状態不変、操作 disabled。
- 状態: 確定。

### SCN-TJSM-007 provider 失敗、応答不正、保存失敗

- 受け入れ条件: credential 未設定では開始しない。invalid response、correlation error、保存失敗は successful `Completed` にしない。correlation error は `RecoverableFailed` として扱う。
- 主な候補: failure `004`、`005`、`006`、`012`、external-integration `001`、`003`、`007`。
- 実行テスト種別: `APIテスト`。
- 実行段階: `実装後`。
- 主要観測点: error kind、retryable flag、failed count、成功済み結果の維持、secret 非露出。
- 状態: 確定。

### SCN-TJSM-008 terminal guard と late response rejection

- 受け入れ条件: terminal job では phase run 作成、保存、readiness 更新、late response 後書きを拒否する。frontend stale event は状態正本にしない。
- 主な候補: state-transition `013`、failure `007`、external-integration `005`、`006`、operation-audit `008`。
- 実行テスト種別: `APIテスト`。
- 実行段階: `実装後`。
- 主要観測点: late response rejected、保存行追加なし、readiness 変更なし、runtime event dropped。
- 状態: 確定。

### SCN-TJSM-009 redacted observability

- 受け入れ条件: operation summary は DB に永続保存しない。状態変更、拒否、削除、再開不可、provider 境界、runtime event 破棄は、必要な時にロジックで導出できる。
- 主な候補: operation-audit `001`、`003`、`007`、external-integration `008`。
- 実行テスト種別: `lower-level only`。
- 実行段階: `実装後`。
- 主要観測点: `event`、`where`、`result`、job ID、phase run ID、reason category、禁止値の不在。
- 状態: 確定。

## UI 設計へ進む条件

UI 設計へ進む条件は、受け入れ観点の確定後に判断する。
未完了一覧、Job Run、Output Review の表示値、操作可否、無効理由、reason category が変わる場合は `ui-design.md` を必須成果物にする。

UI 実装へは進めない。
UI の実画面確認が必要になった場合も、`ui-design` の独立成果物として扱う。

## 人間回答済み事項

- `Q-TJSM-001`: 各フェーズ画面は `JOB_PHASE_RUN.state`、大枠画面は `TRANSLATION_JOB.state` を正本にする。状態不整合時は表示だけで状態を書き換えず、危険操作を無効化する。
- `Q-TJSM-002`: `Ready` job には `JOB_PHASE_RUN` を事前作成しない。
- `Q-TJSM-003`: `RecoverableFailed -> Ready` は廃止し、同じ `JOB_PHASE_RUN` を継続する。
- `Q-TJSM-004`: `Ready` cancel は job-level、phase 開始後 cancel は `Paused` phase からだけ許可する。
- `Q-TJSM-005`: resume は `Paused` だけ、retry は `RecoverableFailed` だけに許可する。
- `Q-TJSM-006`: credential 未設定は開始拒否、correlation error は `RecoverableFailed` とする。
- `Q-TJSM-007`: operation summary は DB に永続保存せず、必要な時にロジックで導出する。
- `Q-TJSM-008`: retry、resume、pause、cancel の可否は phase type で分けず、phase type は開始前提データ、完了判定、service method の選択だけに使う。
- `Q-TJSM-009`: policy の判断結果、rule 名、判定履歴、`PolicyResult` は DB に永続化しない。JobIOService は確定済み状態事実だけを保存する。

## 完了条件

- 人間が `scenario-design.questions.md` に回答済みである。
- `scenario-design.requirement-coverage.json` の人間判断待ちは解消済みである。
- `scenario-design.candidate-coverage.json` の候補衝突は解消済みである。
- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/2026-05-10-translation-job-state-machine-redesign/scenario-design.md --report-out docs/exec-plans/active/2026-05-10-translation-job-state-machine-redesign/scenario-design.requirement-gate.md` は pass である。
