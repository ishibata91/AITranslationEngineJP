# Scenario Candidates: 2026-05-08-translation-flow-navigation-overhaul / operation-audit

- `generator`: `operation-audit`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TFNO-OA`

## Generator Scope

- `viewpoint`: operation-audit
- `included_sources`: `plan.md`, `navigation-state-machine.puml`, `docs/spec.md`, `docs/architecture.md`, `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/translation-job-setup.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`, `docs/detail-specs/translation-output-artifact.md`
- `excluded_sources`: product code, product tests, docs canonical updates, `.codex`, tool permissions, implementation instructions
- `generation_notes`: 最終シナリオ表の確定、候補の採否、候補の統合判断は行わない。ローカルアプリに明示されていない保持期間、監査ログ形式、永続化形式は固定しない。

## Candidate Scenarios

### CAND-TFNO-OA-001 Job Setup 作成結果から単語翻訳ページへ進む材料を後追い確認する

- `source requirement`: `plan.md` lines 16, 67-76; `navigation-state-machine.puml` lines 27-31, 70-71; `translation-job-setup.md` lines 45-48; `term-translation-phase.md` lines 21-31
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TFNO-OA-001`
- `actor`: 利用者、運用確認者
- `trigger`: 利用者が `Job Setup` で job 作成を完了し、単語翻訳ページへ進む。
- `expected outcome`: 単語翻訳ページは、作成結果から受け取った job を対象にし、Ready job と初期 phase 状態を確認できる。旧 `Job Run` のセッション取得を使わず、対象 job が作成結果由来であることを後から確認できる。
- `audit event`: job 作成完了、初期 phase 状態作成、単語翻訳ページへの引き継ぎ。
- `saved summary`: job ID、入力出自要約、作成済み job の状態、単語翻訳開始可否、phase runtime 要約、credential 状態分類。
- `redaction rule`: API key 平文、credential 参照実値、secret store key、endpoint、モデル一覧鮮度 token、外部 provider raw data は保存要約と表示へ出さない。
- `observable point`: `Job Setup` の作成後 summary、単語翻訳ページの jobId と単語翻訳 summary、未完了 job 一覧で同じ job を識別できる情報。
- `related detail requirement type`: `observability_requirement`, `data_requirement`, `consistency_requirement`, `security_requirement`
- `adoption hint`: 新規 job 作成直後の移動候補として扱える。
- `conflict hint`: 単語翻訳フェーズ開始条件は Ready job かつ active phase run なしであるため、作成直後に実行開始まで暗黙に進める候補とは競合する。

### CAND-TFNO-OA-002 未完了 job 一覧から再開対象を固定した材料を後追い確認する

- `source requirement`: `plan.md` lines 17, 42-43, 69-76; `navigation-state-machine.puml` lines 32-35, 74-76; `translation-job-management.md` lines 13-22, 26-49, 70-78
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TFNO-OA-002`
- `actor`: 利用者、運用確認者
- `trigger`: 利用者が未完了 job 一覧で Ready、Paused、RecoverableFailed などの job を選ぶ。
- `expected outcome`: 選択した job ID と表示フェーズが固定され、フェーズページは選択済み job だけを対象にする。一覧表示または再開不可理由の表示だけでは job 状態を変えない。
- `audit event`: 未完了 job 一覧表示、再開対象 job 選択、表示フェーズ固定。
- `saved summary`: job ID、作成日時、入力出自、job state、現在フェーズ、進捗、操作可否、無効理由、再開不可理由カテゴリ、AI 設定要約。
- `redaction rule`: endpoint、credential 参照実値、secret store key、API key 本文、外部 provider 応答原文は一覧、エラー、履歴要約へ出さない。
- `observable point`: 未完了 job 一覧のカード、フェーズページの current phase、操作可否、再開不可理由。
- `related detail requirement type`: `observability_requirement`, `state_requirement`, `recovery_requirement`, `security_requirement`
- `adoption hint`: 途中再開の運用確認候補として扱える。
- `conflict hint`: 旧 `Job Run` のセッション取得で任意 job を表示する候補とは競合する。

### CAND-TFNO-OA-003 フェーズページ直移動防止と復帰先を後追い確認する

- `source requirement`: `plan.md` lines 56-65, 115-119; `navigation-state-machine.puml` lines 66-68, 130-147
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TFNO-OA-003`
- `actor`: 利用者、運用確認者
- `trigger`: 対象 job が確定していない状態で、単語翻訳ページ、NPC ペルソナ生成ページ、本文翻訳ページへ直接入ろうとする。
- `expected outcome`: フェーズページは曖昧な job を推測せず、未完了 job 一覧へ戻す。曖昧な対象 job の summary 取得、操作可否判定、実行操作は行わない。
- `audit event`: job 未確定状態でのフェーズページ入場拒否、未完了 job 一覧への復帰。
- `saved summary`: 対象 job 未確定という短い理由、復帰先が未完了 job 一覧であること、job 状態を変更していないこと。
- `redaction rule`: 直移動の試行元に個人情報、secret、入力本文、provider 情報が含まれる場合は保存要約へ出さない。
- `observable point`: 未完了 job 一覧への復帰、フェーズページ操作が表示されないこと、状態変更がないこと。
- `related detail requirement type`: `observability_requirement`, `consistency_requirement`, `state_requirement`, `security_requirement`
- `adoption hint`: navigation guard の運用確認候補として扱える。
- `conflict hint`: グローバルナビからフェーズページへ直接入れる候補、または route state から job を推測する候補とは競合する。

### CAND-TFNO-OA-004 前工程へ戻れない理由と次へ進めない理由を後追い確認する

- `source requirement`: `plan.md` lines 34-54, 121-130; `navigation-state-machine.puml` lines 37-56, 78-88, 124-129; `term-translation-phase.md` lines 55-78; `persona-generation-phase.md` lines 44-64; `body-translation-phase.md` lines 64-80
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TFNO-OA-004`
- `actor`: 利用者、運用確認者
- `trigger`: 利用者がフェーズページで `次へ進む`、`一覧へ戻る`、または禁止された前工程戻りに関わる導線を確認する。
- `expected outcome`: フェーズページから `Job Setup` と入力データページへ戻る導線は出ない。`次へ進む` が無効な場合は、未完了、参照不能、失敗中などの理由を近接表示または既存 footer のエラー表示方式で確認できる。
- `audit event`: フェーズページでの移動可否判定、禁止導線の非表示、次フェーズ不可理由の表示。
- `saved summary`: 現在フェーズ、phase state、readiness、次へ進めない理由、後続フェーズに必要な参照状態。
- `redaction rule`: 原文全文、訳文全文、provider raw request / response、raw prompt、secret、API key 平文は理由表示、保存要約、エラー要約へ出さない。
- `observable point`: `sticky footer` の移動導線、次へ進めない理由、フェーズページ本文の実行操作。
- `related detail requirement type`: `observability_requirement`, `state_requirement`, `consistency_requirement`, `security_requirement`
- `adoption hint`: 移動制限と phase readiness の運用確認候補として扱える。
- `conflict hint`: state-transition 観点の禁止遷移候補と重なるため、designer 側で統合余地がある。

### CAND-TFNO-OA-005 分解後の各フェーズページで redacted result summary を再表示できることを確認する

- `source requirement`: `plan.md` lines 23-32, 121-130; `term-translation-phase.md` lines 59-69; `persona-generation-phase.md` lines 49-53, 66-72; `body-translation-phase.md` lines 64-75
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TFNO-OA-005`
- `actor`: 利用者、運用確認者
- `trigger`: 利用者が単語翻訳、NPC ペルソナ生成、本文翻訳の各ページを再表示する。
- `expected outcome`: 旧 `Job Run` 分解後も、各フェーズページで phase state、進捗、結果要約、失敗要約、後続可否を再確認できる。表示される要約は redacted summary に限定される。
- `audit event`: フェーズページ再表示、phase result summary 復元、失敗要約確認。
- `saved summary`: provider、model、execution mode、batch mode、credential 状態分類、input count、output count、prompt digest、snapshot digest、error kind、retryable flag、影響件数。
- `redaction rule`: secret、API key 平文、credential 参照実値、secret store key、endpoint、provider raw request / response、raw prompt、原文発話全文、会話文脈全文、翻訳フィールド本文全文は保存要約へ出さない。
- `observable point`: 各フェーズページの current phase、progress、result summary、redacted error summary、後続フェーズ可否。
- `related detail requirement type`: `observability_requirement`, `data_requirement`, `recovery_requirement`, `security_requirement`
- `adoption hint`: 旧 `Job Run` の表示責務をページ分解後も失わない候補として扱える。
- `conflict hint`: external-integration 観点の provider 監査要約候補と重なる可能性がある。

### CAND-TFNO-OA-006 旧 Job Run セッション取得経路が廃止されたことを後追い確認する

- `source requirement`: `plan.md` lines 67-76, 121-130; `navigation-state-machine.puml` lines 134-138; `translation-job-management.md` lines 33-35, 43-45
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TFNO-OA-006`
- `actor`: 利用者、運用確認者
- `trigger`: 利用者がフェーズページを開く、または未完了 job 一覧から job を選ぶ。
- `expected outcome`: フェーズページは job を探す画面にならず、`Job Setup` 作成結果または未完了 job 一覧の選択結果だけを対象 job の入口にする。参照不能 job は表示対象にしない。
- `audit event`: 対象 job の入口確認、セッション取得操作の不在、参照不能 job の表示拒否。
- `saved summary`: job の入口種別、job ID、表示フェーズ、参照可否、拒否理由カテゴリ。
- `redaction rule`: セッション取得廃止の確認で、endpoint、secret、credential 値、provider raw payload、入力本文全文を保存対象にしない。
- `observable point`: フェーズページにセッション取得操作がないこと、未完了 job 一覧の選択結果、参照不能 job の拒否表示。
- `related detail requirement type`: `observability_requirement`, `compatibility_requirement`, `consistency_requirement`, `security_requirement`
- `adoption hint`: 旧 `Job Run` 由来の回帰防止候補として扱える。
- `conflict hint`: CAND-TFNO-OA-001 と CAND-TFNO-OA-002 の入口確認に統合される可能性がある。

### CAND-TFNO-OA-007 翻訳完了ページから出力管理へ移った後の対象 job 材料を確認する

- `source requirement`: `plan.md` lines 78-98, 132-148; `navigation-state-machine.puml` lines 58-61, 90-92, 95-118, 140-143; `body-translation-phase.md` lines 13-16, 42-45, 73-75; `translation-output-artifact.md` lines 19-43, 63-74
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TFNO-OA-007`
- `actor`: 利用者、運用確認者
- `trigger`: 本文翻訳フェーズ完了後、利用者が翻訳完了ページから出力管理へ移動する。
- `expected outcome`: 翻訳完了ページは結果確認と出力管理への案内だけを扱う。成果物出力は別セクションの completed job 一覧と Output Review で対象 job、output readiness、拒否理由、preview、operation summary を確認してから行う。
- `audit event`: job Completed 確認、翻訳完了ページ表示、出力管理への移動、completed job 一覧での出力対象 job 確認。
- `saved summary`: job ID、job-level Completed、output readiness、field result 整合、output status 整合、selected job summary、input provenance summary、operation summary。
- `redaction rule`: provider raw payload、secret、API key、復号可能値、過剰な本文全文は UI、DTO、summary、structured log、debug log、runtime event へ出さない。
- `observable point`: 翻訳完了ページのページング表示、出力管理への移動ボタン、completed job 一覧、Output Review の拒否理由と preview。
- `related detail requirement type`: `observability_requirement`, `data_requirement`, `consistency_requirement`, `security_requirement`
- `adoption hint`: 翻訳管理と成果物出力の責務分離を確認する候補として扱える。
- `conflict hint`: 出力管理へ移動した後に出力対象 job を自動選択するかは未決であり、actor-goal 観点や lifecycle 観点の候補と競合する可能性がある。

## Open Notes

- `human decision candidate`: 直リンク防止や禁止移動の拒否理由を、一時表示だけにするか、後から見返せる履歴要約に含めるかは未決である。
- `human decision candidate`: 出力管理へ移動した後、翻訳完了した job を自動選択するか、completed job 一覧で利用者に選ばせるかは未決である。
- `human decision candidate`: `Job Setup` 作成結果由来か未完了 job 一覧選択由来かを、利用者向け表示に出すか、内部要約だけにするかは未決である。
- `merge candidate`: CAND-TFNO-OA-001、CAND-TFNO-OA-002、CAND-TFNO-OA-006 は、対象 job の入口確認として統合される可能性がある。
- `merge candidate`: CAND-TFNO-OA-004 は、state-transition 観点の禁止遷移候補と統合される可能性がある。
- `rejection candidate`: 監査ログ形式、保持期間、永続化テーブルを固定する候補は、この観点成果物では不採用候補として designer へ残す。
