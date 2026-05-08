# Scenario Candidates: 2026-05-08-translation-flow-navigation-overhaul / external-integration

- `generator`: `external-integration`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TFNO-EI`

## Generator Scope

- `viewpoint`: external-integration
- `included_sources`: `plan.md`, `navigation-state-machine.puml`, `docs/spec.md`, `docs/architecture.md`, `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/translation-job-setup.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`, `docs/detail-specs/translation-output-artifact.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本の変更、`.codex` 変更、real API 前提の provider 実装方針
- `generation_notes`: AI provider、Wails gateway、runtime event、fake/stub、file 境界だけを候補化する。最終採否、統合、競合解消は designer に残す。

## Candidate Scenarios

### CAND-TFNO-EI-001 Job Setup 作成成功後に単語翻訳ページへ job を引き継ぐ

- `source requirement`: `plan.md:16`, `plan.md:67-76`, `navigation-state-machine.puml:27-31`, `navigation-state-machine.puml:70-72`, `translation-job-setup.md:37-48`
- `viewpoint`: external-integration
- `candidate scenario id`: `CAND-TFNO-EI-001`
- `external boundary`: Wails gateway、Job Setup の provider 設定参照、secret 境界
- `actor`: 翻訳 job を新規作成する利用者
- `trigger`: Job Setup で 3 つの翻訳段階の AI 設定が検証済みになり、job 作成を実行する
- `start condition`: 登録済み input と各翻訳段階の provider、model、execution mode、batch mode、API key 状態が作成可能である
- `expected outcome`: Ready job が作成され、作成結果の jobId と初期 phase 状態から単語翻訳ページへ移動する。公開 summary は provider 名、model 名、実行方法、API key 状態分類だけを含み、credential 参照実値、secret store key、endpoint、API key 本文を含まない。
- `fake_or_stub`: Job Setup gateway fake は、job 作成成功、stale モデル一覧拒否、API key 不足拒否を返せる。secret 本文を返す fake は使わない。
- `observable point`: 単語翻訳ページの jobId、current phase、redacted runtime summary、旧セッション取得表示がないこと
- `related detail requirement type`: `success_requirement`, `security_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: 新規作成直後の移動シナリオへ統合できる候補である。
- `conflict hint`: Job Setup の作成失敗時に単語翻訳ページへ遷移しない扱いは failure 観点と接続する。

### CAND-TFNO-EI-002 stale または未設定の provider 情報では job 作成後遷移を始めない

- `source requirement`: `translation-job-setup.md:39-45`, `translation-job-setup.md:61-70`, `plan.md:16`, `navigation-state-machine.puml:28-30`
- `viewpoint`: external-integration
- `candidate scenario id`: `CAND-TFNO-EI-002`
- `external boundary`: AI provider model list API、secret 参照状態、Job Setup gateway
- `actor`: Job Setup で AI 設定を選ぶ利用者
- `trigger`: モデル一覧が stale、model 未選択、または API key 不足の状態で job 作成を試みる
- `start condition`: Job Setup 画面で provider 設定が表示されているが、作成前検証を満たしていない
- `expected outcome`: job は作成されず、単語翻訳ページへ移動しない。UI は stale、model 未選択、API key 未設定、model list 取得失敗を区別して表示し、secret 本文と provider raw data を表示しない。
- `fake_or_stub`: model list gateway fake は、sourceToken stale、API key 未設定、取得失敗、取得済みを返せる。
- `observable point`: 作成ボタンの無効理由または作成拒否理由、移動しない画面状態、redacted error summary
- `related detail requirement type`: `failure_handling_requirement`, `security_requirement`, `boundary_requirement`, `testability_requirement`
- `adoption hint`: Job Setup の provider 境界候補として、正常作成候補と対になる。
- `conflict hint`: failure 観点の job 作成失敗候補と統合対象になりうる。

### CAND-TFNO-EI-003 直リンクまたは復元不整合では phase summary を読まず未完了 job 一覧へ戻す

- `source requirement`: `plan.md:56-65`, `plan.md:113-119`, `navigation-state-machine.puml:66-68`, `navigation-state-machine.puml:130-147`, `architecture.md:125-130`, `architecture.md:207-214`
- `viewpoint`: external-integration
- `candidate scenario id`: `CAND-TFNO-EI-003`
- `external boundary`: Wails gateway、route state、screen state 復元境界
- `actor`: 翻訳セクションを開く利用者
- `trigger`: 対象 jobId が未確定の状態で単語翻訳、NPC ペルソナ生成、本文翻訳のいずれかのページへ入る
- `start condition`: route state または復元 state に jobId がない、または jobId と表示 phase の対応が確定できない
- `expected outcome`: phase summary 取得、provider 設定再解決、runtime event 購読を始めず、未完了 job 一覧へ戻す。対象 job が曖昧なまま phase 操作を表示しない。
- `fake_or_stub`: phase summary gateway spy は、直リンク時に summary 取得が呼ばれないことを観測できる。未完了一覧 gateway fake は Completed 以外の job 一覧を返せる。
- `observable point`: 未完了 job 一覧への復帰、phase summary 取得呼び出しなし、直リンク防止理由
- `related detail requirement type`: `state_requirement`, `security_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: navigation guard の受け入れ候補として扱える。
- `conflict hint`: lifecycle 観点の画面復元候補と前提が重なる。

### CAND-TFNO-EI-004 未完了 job 一覧の選択結果だけで phase page を開始する

- `source requirement`: `plan.md:17`, `plan.md:42-43`, `plan.md:67-76`, `translation-job-management.md:13-22`, `translation-job-management.md:26-49`, `navigation-state-machine.puml:32-35`, `navigation-state-machine.puml:74-76`
- `viewpoint`: external-integration
- `candidate scenario id`: `CAND-TFNO-EI-004`
- `external boundary`: 未完了 job 一覧 gateway、phase summary gateway、secret redaction 境界
- `actor`: 未完了 job を再開したい利用者
- `trigger`: 未完了 job 一覧から Ready、Running、Paused、RecoverableFailed、Failed、Canceled のいずれかの job を選ぶ
- `start condition`: 未完了 job 一覧が読み込まれ、選択した job が参照可能である
- `expected outcome`: 選択結果の jobId と現在フェーズから対象 phase page を表示する。Ready job 表示だけで Running へ暗黙遷移しない。provider、model、execution mode、batch mode、credential 状態分類は表示できるが、secret 本文と外部 provider 応答原文を表示しない。
- `fake_or_stub`: 未完了一覧 gateway fake は、各 job state、参照不能 job、phase progress 集約不能を返せる。phase summary fake は jobId scope の結果だけを返す。
- `observable point`: 選択 jobId、表示 phase、操作可否、redacted AI 設定要約、旧セッション取得なし
- `related detail requirement type`: `alternative_success_requirement`, `state_requirement`, `security_requirement`, `compatibility_requirement`
- `adoption hint`: 途中再開の主要候補として扱える。
- `conflict hint`: 旧 Job Run 表示対象という既存語彙を新 phase page 語彙へ置き換える必要がある。

### CAND-TFNO-EI-005 sticky footer の次へ進むは readiness summary だけを使い provider 実行を起動しない

- `source requirement`: `plan.md:45-54`, `plan.md:121-130`, `navigation-state-machine.puml:37-57`, `navigation-state-machine.puml:78-88`, `term-translation-phase.md:55-72`, `persona-generation-phase.md:44-60`, `body-translation-phase.md:21-43`
- `viewpoint`: external-integration
- `candidate scenario id`: `CAND-TFNO-EI-005`
- `external boundary`: phase summary gateway、runtime event 表示更新境界、AI provider 非起動境界
- `actor`: phase page で次工程へ進もうとする利用者
- `trigger`: sticky footer の `次へ進む` を押す、または次へ進めない理由が更新される
- `start condition`: phase page に jobId があり、対象 phase の summary が読めている
- `expected outcome`: 単語翻訳ページは辞書参照成立後だけ NPC ペルソナ生成ページへ進む。NPC ペルソナ生成ページは snapshot 参照成立後だけ本文翻訳ページへ進む。本文翻訳ページは job Completed と output readiness を確認して翻訳完了ページへ進む。footer 操作だけでは provider request、phase start、retry、cancel を起動しない。
- `fake_or_stub`: phase summary fake は completed、blocked、recoverable failed、参照不能を返せる。AI provider fake は footer 操作では呼ばれないことを観測できる。
- `observable point`: footer の移動可否、近接表示されるブロック理由、provider fake 呼び出しなし、phase 実行操作が本文側に残ること
- `related detail requirement type`: `consistency_requirement`, `state_requirement`, `testability_requirement`, `failure_handling_requirement`
- `adoption hint`: フェーズ間移動と実行操作分離の候補として扱える。
- `conflict hint`: state-transition 観点の phase state 遷移候補と統合対象になる。

### CAND-TFNO-EI-006 runtime event は現在選択中 job だけを更新し、古い job の event で画面を飛ばさない

- `source requirement`: `architecture.md:103-115`, `architecture.md:125-130`, `architecture.md:148-153`, `architecture.md:207-214`, `body-translation-phase.md:39-44`, `body-translation-phase.md:73-75`
- `viewpoint`: external-integration
- `candidate scenario id`: `CAND-TFNO-EI-006`
- `external boundary`: RuntimeEventAdapter、Wails runtime event、screen local handler
- `actor`: phase 実行中または別 job を選び直した利用者
- `trigger`: backend から phase progress、完了、失敗、late response rejected の runtime event が届く
- `start condition`: 画面に現在選択中 jobId と表示 phase があり、runtime event が jobId または phase run scope を持つ
- `expected outcome`: 現在選択中 job に対応する event だけが summary と footer 理由を更新する。別 job または古い phase run の event は、画面遷移や provider 再実行を起こさない。通常の query / command は runtime event へ置き換えない。
- `fake_or_stub`: runtime event adapter fake は、現在 job、別 job、古い phase run、late response rejected を送れる。
- `observable point`: 現在 job の summary 更新、別 job event の無視、late response rejected 表示、query / command 境界維持
- `related detail requirement type`: `concurrency_requirement`, `state_requirement`, `observability_requirement`, `testability_requirement`
- `adoption hint`: runtime event 境界の候補として扱える。
- `conflict hint`: operation-audit 観点の event 記録候補と統合対象になる可能性がある。

### CAND-TFNO-EI-007 翻訳完了ページから出力管理へ移動しても成果物出力を直接開始しない

- `source requirement`: `plan.md:78-98`, `plan.md:132-148`, `navigation-state-machine.puml:58-61`, `navigation-state-machine.puml:90-118`, `translation-output-artifact.md:13-28`
- `viewpoint`: external-integration
- `candidate scenario id`: `CAND-TFNO-EI-007`
- `external boundary`: 翻訳管理 gateway、出力管理 gateway、Output Review 境界
- `actor`: 翻訳結果を確認して出力管理へ移動する利用者
- `trigger`: 翻訳完了ページで出力管理への移動ボタンを押す
- `start condition`: 本文翻訳ページから翻訳完了ページへ到達している
- `expected outcome`: 翻訳管理側では XML 出力、preview、再出力、互換性確認を開始しない。出力管理側は Completed job 一覧を入口にして、出力対象 jobId の固定を扱う。
- `fake_or_stub`: 出力管理 gateway fake は、completed job list、selected job summary、output readiness を返せる。翻訳管理 gateway spy は XML 出力 command が呼ばれないことを観測できる。
- `observable point`: 出力管理セクション表示、Completed job 一覧、翻訳管理側の出力 command 呼び出しなし
- `related detail requirement type`: `consistency_requirement`, `data_requirement`, `testability_requirement`, `alternative_success_requirement`
- `adoption hint`: 翻訳管理と成果物出力の責務分離候補として扱える。
- `conflict hint`: 出力対象 job を自動選択するかどうかは plan の未決事項である。

### CAND-TFNO-EI-008 成果物出力は AI provider、network、secret store を必須経路にしない

- `source requirement`: `spec.md:65-67`, `translation-output-artifact.md:24-43`, `translation-output-artifact.md:47-56`, `navigation-state-machine.puml:102-118`
- `viewpoint`: external-integration
- `candidate scenario id`: `CAND-TFNO-EI-008`
- `external boundary`: XML file boundary、Output Review gateway、AI provider 非依存境界
- `actor`: Completed job から xTranslator 互換 XML を出力する利用者
- `trigger`: Output Review で preview、XML 出力、再出力を行う
- `start condition`: body phase Completed、job-level Completed、field result 整合、output status 整合を満たす job が選択されている
- `expected outcome`: 出力処理は provider、network、secret store を必須経路にしない。row validation、XML serialization、file write、artifact 保存の失敗は成功 artifact として扱わない。UI、DTO、summary、structured log、runtime event へ secret、API key、provider raw payload、過剰な本文全文を出さない。
- `fake_or_stub`: XML adapter stub は parse 成功、parse 失敗、readonly path、file write 失敗を返せる。AI provider fake は出力処理で呼ばれないことを観測できる。
- `observable point`: output readiness、row count、file path、failed stage、provider fake 呼び出しなし、redacted error summary
- `related detail requirement type`: `boundary_requirement`, `failure_handling_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: 成果物出力セクションの外部境界候補として扱える。
- `conflict hint`: output-artifact 既存シナリオとの重複候補であり、designer が最終統合する。

## Open Notes

- `candidate_count`: 8
- `human decision candidate`: 出力管理へ移動した後、出力対象 job を自動選択するか、出力管理側で選ばせるか。根拠は `plan.md:86-87` と `plan.md:146-148`。
- `merge candidate`: `CAND-TFNO-EI-001` と `CAND-TFNO-EI-002` は Job Setup の provider 境界として統合される可能性がある。
- `merge candidate`: `CAND-TFNO-EI-003` と `CAND-TFNO-EI-004` は jobId 確定経路の候補として統合される可能性がある。
- `conflict candidate`: `CAND-TFNO-EI-007` は BodyPhasePage から `Canceled` / `Failed` でも翻訳完了ページへ進む図の表現と、出力候補を Completed job だけにする既存仕様の間で確認が必要である。
- `rejection candidate`: real API を使う provider 接続確認だけの候補は、有料 real API 前提を固定するため除外した。
