# Scenario Candidates: 2026-05-08-translation-flow-navigation-overhaul / lifecycle

- `generator`: `lifecycle`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TFNO-LC`

## Generator Scope

- `viewpoint`: lifecycle
- `included_sources`: `./plan.md`, `./navigation-state-machine.puml`, `docs/spec.md`, `docs/architecture.md`, `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/translation-job-setup.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`, `docs/detail-specs/translation-output-artifact.md`
- `excluded_sources`: 他観点の候補成果物、未引き継ぎ会話文脈、プロダクトコード、プロダクトテスト、docs 正本変更、`.codex` 変更
- `generation_notes`: 作成、再開、フェーズ進行、完了後利用の時間順だけを候補化する。採否、統合、競合解消は `designer` に残す。

## Candidate Scenarios

### CAND-TFNO-LC-001 翻訳セクション入口で新規開始または再開入口を選ぶ

- `source requirement`: `./plan.md:34-43`, `./plan.md:56-65`, `./navigation-state-machine.puml:18-20`, `./navigation-state-machine.puml:66-68`, `docs/spec.md:132-133`
- `viewpoint`: lifecycle
- `lifecycle stage`: 翻訳セクション開始
- `start condition`: 利用者がグローバルナビまたは dashboard から翻訳セクションを開く。
- `candidate scenario id`: `CAND-TFNO-LC-001`
- `actor`: 翻訳 job を開始または再開したい利用者
- `trigger`: 翻訳セクション入口を表示する。
- `expected outcome`: 新規開始は入力データページへ進み、途中再開は未完了 job 一覧へ進む。フェーズページへ直接入ろうとした復元状態は未完了 job 一覧へ戻る。
- `observable point`: 入口画面で新規開始と途中再開の経路を区別できる。対象 job が未確定のフェーズ画面は表示されない。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `consistency_requirement`
- `adoption hint`: フロー入口と直リンク防止を 1 つの起点候補として扱える。
- `conflict hint`: state-transition 観点の直リンク禁止候補と重複する可能性がある。

### CAND-TFNO-LC-002 入力データを固定して Job Setup へ進む

- `source requirement`: `./navigation-state-machine.puml:22-30`, `./navigation-state-machine.puml:70-71`, `docs/spec.md:13-15`, `docs/detail-specs/translation-job-setup.md:19-22`, `docs/detail-specs/translation-job-setup.md:26-38`
- `viewpoint`: lifecycle
- `lifecycle stage`: job 作成前
- `start condition`: 登録済み入力データ、共通辞書、共通ペルソナ、AIサービス設定の参照状態を確認できる。
- `candidate scenario id`: `CAND-TFNO-LC-002`
- `actor`: 新規翻訳 job を作成したい利用者
- `trigger`: 入力データページで対象 input を固定して `Job Setup` へ進む。
- `expected outcome`: `Job Setup` は対象 input と 3 つの翻訳段階の設定を受け取り、作成前検証へ進める状態になる。
- `observable point`: 選択済み input、共通基盤、単語翻訳、NPC ペルソナ生成、本文翻訳の設定状態を確認できる。
- `related detail requirement type`: `success_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: 新規 job 作成前の入力固定段階として扱える。
- `conflict hint`: actor-goal 観点の job 作成候補と統合される可能性がある。

### CAND-TFNO-LC-003 Job Setup 完了直後に単語翻訳ページへ進む

- `source requirement`: `./plan.md:67-76`, `./navigation-state-machine.puml:27-31`, `./navigation-state-machine.puml:70-72`, `docs/detail-specs/translation-job-setup.md:37-47`, `docs/detail-specs/term-translation-phase.md:21-24`
- `viewpoint`: lifecycle
- `lifecycle stage`: job 作成完了直後
- `start condition`: `Job Setup` で 3 つの翻訳段階の APIキー不足と model 未選択がない。
- `candidate scenario id`: `CAND-TFNO-LC-003`
- `actor`: 新規 job を作成した利用者
- `trigger`: `Job Setup` の job 作成が完了する。
- `expected outcome`: 作成された Ready job と初期 phase 状態を受け取り、単語翻訳ページへ遷移する。旧 `Job Run` のセッション取得は要求しない。
- `observable point`: 単語翻訳ページに jobId、単語翻訳 summary、開始可否、設定要約が表示される。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `compatibility_requirement`
- `adoption hint`: 旧セッション取得廃止の主候補として扱える。
- `conflict hint`: 既存詳細仕様に残る `Job Run` 表示対象という表現と用語衝突する可能性がある。

### CAND-TFNO-LC-004 未完了 job 一覧から再開対象と表示フェーズを固定する

- `source requirement`: `./plan.md:42-43`, `./navigation-state-machine.puml:32-35`, `./navigation-state-machine.puml:74-76`, `docs/detail-specs/translation-job-management.md:19-22`, `docs/detail-specs/translation-job-management.md:26-44`
- `viewpoint`: lifecycle
- `lifecycle stage`: 途中再開
- `start condition`: Ready、Running、Paused、RecoverableFailed、Failed、Canceled の job が存在し、翻訳管理を開ける。
- `candidate scenario id`: `CAND-TFNO-LC-004`
- `actor`: 未完了 job を再開したい利用者
- `trigger`: 未完了 job 一覧で対象 job を選択する。
- `expected outcome`: 選択した jobId と表示フェーズが固定され、単語翻訳、NPC ペルソナ生成、本文翻訳の該当ページへ進む。一覧表示だけでは job 状態を変えない。
- `observable point`: 一覧に job state、現在フェーズ、進捗、操作可否、再開不可理由が表示される。
- `related detail requirement type`: `alternative_success_requirement`, `state_requirement`, `recovery_requirement`
- `adoption hint`: 新規作成とは別の lifecycle 再開候補として扱える。
- `conflict hint`: 旧 `Job Run` への表示対象設定を、フェーズページ表示対象へ読み替える必要がある。

### CAND-TFNO-LC-005 単語翻訳ページから NPC ペルソナ生成ページへ進む

- `source requirement`: `./plan.md:45-55`, `./plan.md:121-130`, `./navigation-state-machine.puml:37-43`, `./navigation-state-machine.puml:78-80`, `docs/spec.md:129-130`, `docs/detail-specs/term-translation-phase.md:21-24`, `docs/detail-specs/term-translation-phase.md:55-75`
- `viewpoint`: lifecycle
- `lifecycle stage`: 単語翻訳完了後
- `start condition`: 単語翻訳フェーズが Completed になり、ジョブ内辞書参照が成立している。
- `candidate scenario id`: `CAND-TFNO-LC-005`
- `actor`: 翻訳 job を次フェーズへ進めたい利用者
- `trigger`: 単語翻訳ページの `次へ進む` を使う。
- `expected outcome`: NPC ペルソナ生成ページへ遷移する。未完了、失敗中、辞書参照不能の場合は進めず、理由を表示する。
- `observable point`: `sticky footer` が `次へ進む`、`一覧へ戻る`、進めない理由を表示する。実行、一時停止、再開、再試行はページ本文にある。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `consistency_requirement`
- `adoption hint`: フェーズページ間を footer の `次へ進む` に制限する候補として扱える。
- `conflict hint`: UI 観点の `sticky footer` 候補と統合される可能性がある。

### CAND-TFNO-LC-006 NPC ペルソナ生成ページから本文翻訳ページへ進む

- `source requirement`: `./plan.md:45-55`, `./navigation-state-machine.puml:44-50`, `./navigation-state-machine.puml:82-84`, `docs/spec.md:130-131`, `docs/detail-specs/persona-generation-phase.md:21-24`, `docs/detail-specs/persona-generation-phase.md:40-60`, `docs/detail-specs/body-translation-phase.md:20-23`
- `viewpoint`: lifecycle
- `lifecycle stage`: NPC ペルソナ生成完了後
- `start condition`: NPC ペルソナ生成フェーズが Completed であり、persona snapshot 参照が成立している。
- `candidate scenario id`: `CAND-TFNO-LC-006`
- `actor`: 翻訳 job を本文翻訳へ進めたい利用者
- `trigger`: NPC ペルソナ生成ページの `次へ進む` を使う。
- `expected outcome`: 本文翻訳ページへ遷移する。persona 未完了、失敗、snapshot 参照不能では本文翻訳フェーズを開始しない。
- `observable point`: persona snapshot 参照状態、body phase readiness、進めない理由を確認できる。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `consistency_requirement`
- `adoption hint`: フェーズ順序の中間候補として扱える。
- `conflict hint`: state-transition 観点の readiness 候補と統合される可能性がある。

### CAND-TFNO-LC-007 本文翻訳完了後に翻訳完了ページへ進む

- `source requirement`: `./plan.md:78-88`, `./navigation-state-machine.puml:51-61`, `./navigation-state-machine.puml:86-92`, `docs/detail-specs/body-translation-phase.md:13-23`, `docs/detail-specs/body-translation-phase.md:39-43`
- `viewpoint`: lifecycle
- `lifecycle stage`: 本文翻訳完了後
- `start condition`: 本文翻訳フェーズが Completed であり、field result 整合と output status 整合が成立している。
- `candidate scenario id`: `CAND-TFNO-LC-007`
- `actor`: 翻訳結果を確認したい利用者
- `trigger`: 本文翻訳ページで job 完了条件が成立する。
- `expected outcome`: 翻訳完了ページへ進み、原文と訳文をページング表示できる。翻訳完了ページでは出力処理そのものを扱わない。
- `observable point`: 原文、訳文、output readiness、出力管理への移動ボタンを確認できる。
- `related detail requirement type`: `success_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: 翻訳管理内の完了確認候補として扱える。
- `conflict hint`: `navigation-state-machine.puml:87` は `Canceled` と `Failed` でも翻訳完了ページへ遷移すると読める。成功完了ページへ含めるかは未解決候補として残す。

### CAND-TFNO-LC-008 翻訳完了ページから出力管理へ移動する

- `source requirement`: `./plan.md:78-98`, `./plan.md:132-139`, `./plan.md:146-148`, `./navigation-state-machine.puml:58-62`, `./navigation-state-machine.puml:90-92`, `docs/detail-specs/translation-output-artifact.md:19-28`
- `viewpoint`: lifecycle
- `lifecycle stage`: 完了後利用
- `start condition`: 翻訳完了ページで Completed job の結果を確認できる。
- `candidate scenario id`: `CAND-TFNO-LC-008`
- `actor`: 翻訳成果物を出力したい利用者
- `trigger`: 翻訳完了ページの出力管理への移動ボタンを使う。
- `expected outcome`: 翻訳管理から成果物出力セクションへ移動する。翻訳管理は XML 出力、preview、再出力、互換性確認を開始しない。
- `observable point`: 出力管理側で completed job 一覧または selected job summary を確認できる。
- `related detail requirement type`: `alternative_success_requirement`, `consistency_requirement`, `data_requirement`
- `adoption hint`: 翻訳完了後の別セクション移動候補として扱える。
- `conflict hint`: 出力対象 job を自動選択するか、出力管理側で選ばせるかは未決である。

### CAND-TFNO-LC-009 Completed job 一覧から Output Review へ進む

- `source requirement`: `./plan.md:89-98`, `./navigation-state-machine.puml:95-118`, `docs/spec.md:61-67`, `docs/detail-specs/translation-output-artifact.md:19-28`, `docs/detail-specs/translation-output-artifact.md:61-68`
- `viewpoint`: lifecycle
- `lifecycle stage`: 成果物出力セクション開始
- `start condition`: body phase が Completed で、job-level 状態も Completed である翻訳 job が存在する。
- `candidate scenario id`: `CAND-TFNO-LC-009`
- `actor`: 完了済み job から成果物を出力したい利用者
- `trigger`: 成果物出力セクションで Completed job を選択する。
- `expected outcome`: 出力対象 jobId が固定され、Output Review ページへ進む。未完了、失敗中、Canceled、不整合 job は出力候補にしない。
- `observable point`: completed job list、selected job summary、output readiness、拒否理由、result summary を確認できる。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `data_requirement`
- `adoption hint`: 翻訳管理から分離された成果物出力 lifecycle 候補として扱える。
- `conflict hint`: operation-audit 観点の artifact 出力履歴候補と統合される可能性がある。

### CAND-TFNO-LC-010 フェーズページから入力データページまたは Job Setup へ戻らない

- `source requirement`: `./plan.md:34-43`, `./navigation-state-machine.puml:124-129`, `docs/detail-specs/translation-job-management.md:33-35`, `docs/detail-specs/translation-job-setup.md:29-36`
- `viewpoint`: lifecycle
- `lifecycle stage`: 作成済み job の前工程戻り禁止
- `start condition`: 作成済み job のフェーズページを表示している。
- `candidate scenario id`: `CAND-TFNO-LC-010`
- `actor`: 作成済み job を進行中に見直したい利用者
- `trigger`: フェーズページから前工程へ戻ろうとする。
- `expected outcome`: 入力データページまたは `Job Setup` へ戻る導線は表示されない。入力や設定を変えたい場合は既存 job を巻き戻さず、新しい job 作成の経路を使う。
- `observable point`: フェーズページの移動導線は `次へ進む`、`一覧へ戻る`、`出力管理へ移動` の範囲に収まる。
- `related detail requirement type`: `consistency_requirement`, `state_requirement`, `compatibility_requirement`
- `adoption hint`: lifecycle の巻き戻し禁止候補として扱える。
- `conflict hint`: failure 観点の作り直し候補と統合される可能性がある。

## Open Notes

- `human decision candidate`: `CAND-TFNO-LC-008` は、出力管理へ移動した後に出力対象 job を自動選択するか、出力管理側で選ばせるかが未決である。
- `human decision candidate`: `CAND-TFNO-LC-007` は、`Canceled` または `Failed` の job を翻訳完了ページへ出すか、未完了一覧または別の終端表示へ戻すかが未決である。
- `merge candidate`: `CAND-TFNO-LC-003` と `CAND-TFNO-LC-004` は、旧 `Job Run` のセッション取得廃止を扱う候補と統合される可能性がある。
- `merge candidate`: `CAND-TFNO-LC-005`、`CAND-TFNO-LC-006`、`CAND-TFNO-LC-007` は、フェーズ順序と readiness の状態遷移候補と統合される可能性がある。
- `rejection candidate`: なし。採否判断は `designer` に残す。
