# Scenario Candidates: 2026-05-16-translation-job-state-stale-retirement / failure

- `generator`: `failure`
- `source_plan`: `./implement-lane-task-frame.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJSSR`

## Generator Scope

- `viewpoint`: 失敗観点。正本 state から外れた state 名、状態不整合、拒否理由の重複、旧名参照、state spelling の検索漏れを候補化する。
- `included_sources`: `./implement-lane-task-frame.md`, `./state-knowledge-investigation.md`, `../../../spec.md`, `../../../architecture.md`, `../../../detail-specs/translation-job-management.md`, `../../../detail-specs/term-translation-phase.md`, `../../../detail-specs/persona-generation-phase.md`, `../../../detail-specs/body-translation-phase.md`
- `excluded_sources`: product code、product test、docs 正本本文、`docs/exec-plans/completed/**`、UI 変更、`stale_selection`、`validation_stale`、`model_selection_stale` の削除判断。
- `generation_notes`: 採否、統合、最終シナリオ表は `designer` に残す。候補は scenario-design が受け入れ条件へ変換しやすい単位に分ける。

## Candidate Scenarios

### CAND-TJSSR-001 `pending` が phase 実行状態として漏出する

- `source requirement`: `docs/spec.md` の `JOB_PHASE_RUN.state` は `Running`、`Paused`、`RecoverableFailed`、`Completed`、`Failed`、`Canceled` に限定される。`Ready` job には `JOB_PHASE_RUN` を事前作成しない。`state-knowledge-investigation.md` は 3 phase の service に `pending` が残ると記録している。
- `viewpoint`: 設定不整合、失敗入力、回復動作。
- `candidate scenario id`: `CAND-TJSSR-001`
- `actor`: 翻訳 job を実行または再表示する利用者。
- `trigger`: 既存の `JOB_PHASE_RUN.state` または read model 派生値に `pending` が残った状態で、単語翻訳フェーズ、NPC ペルソナ生成フェーズ、本文翻訳フェーズの開始または再表示を行う。
- `rejected operation`: `pending` を正本 phase state として表示すること。`pending` の `JOB_PHASE_RUN` を active run として操作可否判定に使うこと。Ready job に phase run を事前作成すること。
- `expected error`: 仕様外 state として扱い、成功状態または操作可能状態へ丸めない。危険操作は無効にし、状態を表示だけで書き換えない。エラー理由は実行不可理由カテゴリとして観測できる。
- `expected outcome`: 正本 state に存在しない `pending` が利用者向け state、DTO state、永続 state の成功値として漏出しない。開始許可時に作られる phase run は `Running` から始まる。
- `observable point`: `TRANSLATION_JOB.state`、`JOB_PHASE_RUN.state`、phase read model、操作可否、実行不可理由カテゴリ、phase run 作成有無。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `consistency_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: 3 phase 共通の fixture または repository seed で `pending` を入れ、API レベルで state と操作可否を確認する受け入れ条件へ変換しやすい。
- `conflict hint`: `pending` を内部一時 state として隔離するか、正本 state へ昇格するかは未決である。昇格する場合は `docs/spec.md` の state 正本と衝突する。

### CAND-TJSSR-002 job state と phase state が食い違う時に危険操作を拒否する

- `source requirement`: `docs/detail-specs/translation-job-management.md` は、保存済み `TRANSLATION_JOB.state` と現在フェーズの `JOB_PHASE_RUN.state` が食い違う場合、表示だけで状態を書き換えず、危険操作を無効にすると定義している。`docs/spec.md` は job state と phase state の用途を分けている。
- `viewpoint`: 設定不整合、回復動作、競合候補。
- `candidate scenario id`: `CAND-TJSSR-002`
- `actor`: 未完了 job 一覧または Job Run を開く利用者。
- `trigger`: `TRANSLATION_JOB.state` が `Ready` なのに active な `JOB_PHASE_RUN` が存在する、または `TRANSLATION_JOB.state` が terminal なのに non-terminal の `JOB_PHASE_RUN` が存在する。
- `rejected operation`: phase start、pause、resume、retry、cancel、late response 後書き、readiness update、削除などの危険操作。
- `expected error`: 状態不整合を理由カテゴリとして返す。空一覧、成功状態、暗黙の state 修復にはしない。
- `expected outcome`: 一覧表示と Job Run 表示は状態を保存し直さない。操作可否は安全側へ倒れ、参照不能 job は Job Run の表示対象にならない。
- `observable point`: job 一覧 response、Job Run 対象選択可否、操作可否、拒否理由カテゴリ、DB の `TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` の変更有無。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `consistency_requirement`, `recovery_requirement`, `testability_requirement`
- `adoption hint`: job state と phase state の不整合 fixture を 2 種類置き、query では状態が変わらず、command では拒否されることを受け入れ条件にできる。
- `conflict hint`: 正本化時に自動修復を認める仕様へ変える場合は、現在の「表示だけでは状態を書き換えない」仕様と競合する。

### CAND-TJSSR-003 操作可否の拒否理由が policy と service で重複または不一致になる

- `source requirement`: `docs/spec.md` は共通操作規則を `JOB_PHASE_RUN.state` に基づける。`docs/architecture.md` は `TranslationJobPolicy` が共通操作規則を評価し、`Backend UseCase` だけが呼ぶと定義している。`state-knowledge-investigation.md` は service helper が同じ state 規則を文字列で再実装していると記録している。
- `viewpoint`: 設定不整合、保存失敗、競合候補。
- `candidate scenario id`: `CAND-TJSSR-003`
- `actor`: pause、resume、retry、cancel の操作結果を確認する利用者。
- `trigger`: 同じ phase state に対して、`TranslationJobPolicy` 由来の拒否理由と service helper 由来の拒否理由が同時に評価される。
- `rejected operation`: `Running` 以外の pause、`Paused` 以外の resume、`RecoverableFailed` 以外の retry、`Paused` 以外の cancel。
- `expected error`: 1 つの操作に対して、利用者が読める拒否理由カテゴリは 1 つに収束する。policy と service の二重判定により、相反する可否や複数の理由カテゴリを返さない。
- `expected outcome`: phase 種別に依存しない共通操作規則が維持される。phase 固有の `canRetry`、`canResume`、`canPause`、`canCancel` は復活しない。
- `observable point`: API response の operation availability、reason category、phase service の result summary、structured log に出る reason category。
- `related detail requirement type`: `failure_handling_requirement`, `consistency_requirement`, `state_requirement`, `observability_requirement`, `compatibility_requirement`
- `adoption hint`: 各 phase state と各操作の組み合わせを table-driven にし、可否と拒否理由カテゴリが 1 件だけ返ることを受け入れ条件にできる。
- `conflict hint`: `TranslationJobPolicy` を UseCase 専用に固定する architecture 制約と、read model 集約で同じ規則を使いたい要求が衝突しうる。

### CAND-TJSSR-004 旧名参照が active task-local から再注入される

- `source requirement`: `implement-lane-task-frame.md` は `StateMachine` 旧名が product code から外れたが active observability task-local に残っていると記録している。`state-knowledge-investigation.md` は `StateMachine` と `JobIOService` の旧名参照が再注入リスクになると記録している。`docs/architecture.md` は現在の構造主語として `TranslationJobPolicy` と `JobIOService` を定義しているが、実体 package は `doc.go` だけであると調査に記録されている。
- `viewpoint`: 参照不能、設定不整合、回帰。
- `candidate scenario id`: `CAND-TJSSR-004`
- `actor`: stale 廃止後に設計または実装を進める実装者。
- `trigger`: active task-local の `StateMachine` 旧名、または実体の薄い `JobIOService` 参照を根拠にして、新しい実装範囲や観測ログ設計を作る。
- `rejected operation`: product code、arch lint、active task-local 成果物へ `StateMachine` 旧境界を再追加すること。実体が未確定な `JobIOService` を保存責務以上の状態遷移判断として扱うこと。
- `expected error`: 旧名参照を stale 参照として検出し、採用前に設計判断へ戻す。旧境界を根拠にした実装範囲をそのまま承認しない。
- `expected outcome`: stale 廃止後の成果物は、状態遷移判断を `TranslationJobPolicy`、状態事実保存を確定済み境界へ分けて扱う。旧名参照が残る場合は競合候補または人間判断候補として visible になる。
- `observable point`: active task-local 成果物の検索結果、`.go-arch-lint.yml` の component 名、architecture 図または実装範囲の構造主語、review 指摘の有無。
- `related detail requirement type`: `compatibility_requirement`, `consistency_requirement`, `testability_requirement`, `observability_requirement`
- `adoption hint`: scenario-design では検索観測を受け入れ条件にし、再注入の検出対象と許容対象を分けられる。
- `conflict hint`: `JobIOService` を architecture 正本から外すか、別 task で実体化するかは人間判断が必要である。

### CAND-TJSSR-005 `cancelled` spelling の検索漏れで旧 state が残る

- `source requirement`: `docs/spec.md` と phase detail-spec は正本 state を `Canceled` または `canceled` 表記で扱う。`state-knowledge-investigation.md` は `PersonaGenerationPhaseContractStub` の cancel fixture が `cancelled` を返すと記録している。
- `viewpoint`: 失敗入力、設定不整合、回帰。
- `candidate scenario id`: `CAND-TJSSR-005`
- `actor`: cancel 済み phase または job の状態を確認する利用者。
- `trigger`: state 語彙の棚卸しで `canceled` だけを検索し、`cancelled` を返す fixture、stub、test data、DTO 変換が残る。
- `rejected operation`: `cancelled` を正本 state として受理すること。`cancelled` を terminal 判定、操作可否、表示 state の成功値として流すこと。
- `expected error`: `cancelled` は旧 spelling として検出される。正本 state へ合わせるか、仕様外 state として拒否される。検索観測では `canceled` と `cancelled` の両方を対象にする。
- `expected outcome`: terminal 判定、cancel 後の操作不可、late response 後書き拒否は `Canceled` 正本 spelling だけで成立する。British spelling の残存で検索漏れや誤判定を起こさない。
- `observable point`: fixture response、contract stub、DTO state、operation availability、terminal guard、repository seed、grep 結果。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: `canceled` と `cancelled` の両 spelling を検索対象にし、`cancelled` が正本 state として残らないことを受け入れ条件にできる。
- `conflict hint`: 利用者向け表示文言で自然言語として `cancelled` を使う仕様が将来入る場合は、state value と表示文言を分ける必要がある。

## Open Notes

- `human decision candidate`: `pending` を canonical state へ昇格するか、内部一時 state として隔離するかは人間判断が必要である。
- `human decision candidate`: `JobIOService` を architecture 正本から外すか、別 task で実体化するかは人間判断が必要である。
- `human decision candidate`: active observability task-local の旧名参照を今回の stale 廃止に含めるかは人間判断が必要である。
- `merge candidate`: `CAND-TJSSR-001` と `CAND-TJSSR-002` は、仕様外 state と状態不整合の安全側表示シナリオとして統合できる可能性がある。
- `merge candidate`: `CAND-TJSSR-003` は state-transition 観点の操作可否候補と統合できる可能性がある。
- `rejection candidate`: UI 表示変更を前提にする候補は今回の failure 候補から除外する。
- `rejection candidate`: `stale_selection`、`validation_stale`、`model_selection_stale` の削除や名称変更を求める候補は今回の failure 候補から除外する。
