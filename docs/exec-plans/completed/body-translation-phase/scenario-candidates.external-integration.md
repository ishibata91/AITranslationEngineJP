# Scenario Candidates: body-translation-phase / external-integration

- `generator`: `external-integration`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `BTP`

## Generator Scope

- `viewpoint`: AI provider、credential、execution mode、fake transport、prompt/input/output adapter、network boundary、外部互換出力へ渡る status 境界に限定する。
- `included_sources`:
  - `./plan.md`
  - `../../../../tasks/index.yaml`
  - `../../../../tasks/usecases/body-translation-phase.yaml`
  - `../../../../tasks/usecases/term-translation-phase.yaml`
  - `../../../../tasks/usecases/persona-generation-phase.yaml`
  - `../../../../tasks/usecases/translation-output-artifact.yaml`
  - `../../../spec.md`
  - `../../../er.md`
  - `../../../architecture.md`
  - `../term-translation-phase/scenario-design.md`
  - `../persona-generation-phase/scenario-design.md`
- `excluded_sources`: 採否決定、最終シナリオ表、実装範囲、プロダクトコード、プロダクトテスト、docs 正本化、real paid API 前提の検証。
- `generation_notes`: `designer` が統合判断できる候補だけを書く。外部 provider の実装方式、provider 選択、失敗時の最終状態は固定しない。

## Candidate Scenarios

### CAND-BTP-001 本文翻訳用 AI provider を fake transport で実行できる

- `source requirement`: `tasks/usecases/body-translation-phase.yaml` の「翻訳レコード種別に応じた翻訳指示を構成できる」、`docs/spec.md` の AI 実行基盤、`docs/er.md` の `JOB_PHASE_RUN` にフェーズ別 AI 設定と最終 AI 実行情報を保持する要件。
- `viewpoint`: provider 境界
- `candidate scenario id`: `CAND-BTP-001`
- `external boundary`: AIProvider boundary、phase execution setting boundary、fake transport boundary。
- `actor`: 翻訳ジョブ実行者
- `trigger`: persona phase Completed の job で本文翻訳フェーズを開始する。
- `start condition`: provider、model、execution mode、credential ref、本文翻訳対象 field 件数を phase 実行設定として解決できる。
- `expected outcome`: paid real API を呼ばずに fake transport で本文翻訳 provider 呼び出しを観測できる。phase result には provider、model、execution mode、request unit count、output count の要約が残る。
- `observable point`: fake transport log、provider execution summary、Job Run の current phase / progress、`JOB_PHASE_RUN`。
- `fake_or_stub`: fake AI provider、provider selection fixture、phase run fixture、temp DB。
- `related detail requirement type`: `external_integration_requirement`, `testability_requirement`, `state_requirement`
- `adoption hint`: 本文翻訳フェーズの provider 正常実行候補として扱える。
- `conflict hint`: provider / model / execution mode を Job Setup から継承するか、本文翻訳 phase 開始時に再選択するかは人間判断候補になりうる。

### CAND-BTP-002 保存済み credential を参照し secret を露出しない

- `source requirement`: `docs/spec.md` の「各フェーズのAPI選択、APIKeyは再入力不要で保存ができること」「APIKeyは暗号化して保存すること」、`docs/er.md` の `credential_ref` は secret store 参照だけを保持する要件。
- `viewpoint`: secret 境界
- `candidate scenario id`: `CAND-BTP-002`
- `external boundary`: secret store boundary、AIProvider request boundary、runtime log boundary、Job Run summary boundary。
- `actor`: 翻訳ジョブ実行者
- `trigger`: 保存済み API key を使う provider 設定で本文翻訳フェーズを実行する。
- `start condition`: credential ref は解決できるが、API key 平文は UI と phase result に渡さない。
- `expected outcome`: API key、secret 本体、復号可能値は UI、error summary、structured log、fake transport log、provider request 観測値へ出ない。credential ref と参照状態だけを観測できる。
- `observable point`: fake secret store assertion、redaction assertion、Job Run summary、structured log。
- `fake_or_stub`: fake secret store、fake transport、redaction assertion fixture。
- `related detail requirement type`: `security_requirement`, `external_integration_requirement`, `observability_requirement`
- `adoption hint`: body phase の AI 実行でも secret 非露出を受け入れ条件にできる。
- `conflict hint`: operation-audit 観点が raw prompt や raw response を保存したい場合、secret / 過剰本文保存禁止と衝突しうる。

### CAND-BTP-003 翻訳フィールド本文と補助入力を provider request へ写像する

- `source requirement`: `tasks/usecases/body-translation-phase.yaml` の入力一式、`docs/spec.md` の翻訳指示、辞書再利用、ペルソナ提供、埋め込み要素保持要件、`docs/architecture.md` の `AIProvider` adapter 境界。
- `viewpoint`: adapter 境界
- `candidate scenario id`: `CAND-BTP-003`
- `external boundary`: prompt builder boundary、AIProvider request adapter boundary、job dictionary / persona snapshot reference boundary。
- `actor`: 翻訳ジョブ実行者
- `trigger`: 翻訳フィールド本文、確定訳語、ジョブ内辞書、ジョブ内ペルソナ、翻訳補助メタデータが揃った状態で本文翻訳 request を作る。
- `start condition`: 翻訳 field ID、record type、field type、原文、確定訳語参照、persona snapshot 参照、保護要素一覧を入力 summary として解決できる。
- `expected outcome`: provider request は本文翻訳に必要な入力参照を欠落なく持つ。fake transport では raw prompt 全文ではなく request summary、input count、prompt digest、保護要素 digest を観測できる。
- `observable point`: fake transport request summary、prompt digest、input reference summary、protection element digest。
- `fake_or_stub`: prompt builder fixture、job dictionary fixture、persona snapshot fixture、translation metadata fixture、fake transport request capture。
- `related detail requirement type`: `external_integration_requirement`, `consistency_requirement`, `security_requirement`
- `adoption hint`: body translation prompt/input adapter の contract freeze 候補にできる。
- `conflict hint`: prompt 全文を保存して調査したい要求は、raw prompt / 過剰本文保存禁止と衝突しうる。

### CAND-BTP-004 provider 応答を訳文、出力ステータス、保護要素検証結果へ写像する

- `source requirement`: `tasks/usecases/body-translation-phase.yaml` の出力、`docs/spec.md` の訳文と出力ステータス lossless 保持、`docs/er.md` の `JOB_TRANSLATION_FIELD` に翻訳結果と出力ステータスを保持する要件。
- `viewpoint`: adapter 境界
- `candidate scenario id`: `CAND-BTP-004`
- `external boundary`: AIProvider response adapter boundary、protection validation boundary、job translation field persistence boundary。
- `actor`: 翻訳ジョブ実行者
- `trigger`: fake provider が valid な本文翻訳応答を返す。
- `start condition`: provider 応答には field correlation key、訳文、status 判定材料、保護要素検証に必要な出力本文が含まれる。
- `expected outcome`: provider 応答は対象 `JOB_TRANSLATION_FIELD` の訳文、出力ステータス、保護要素検証結果へ写像される。FormID、EditorID、record type、field type、source text との対応を失わない。
- `observable point`: response adapter output、`JOB_TRANSLATION_FIELD`、protection validation result、phase result。
- `fake_or_stub`: fixed provider response fixture、response adapter fixture、protection validator fixture、temp DB。
- `related detail requirement type`: `external_integration_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: provider output adapter と persistence 境界の正常系候補にできる。
- `conflict hint`: output status の細分化は translation-output-artifact 観点や failure 観点と統合時に整理が必要である。

### CAND-BTP-005 provider 応答が保護要素を破壊した場合に成功保存しない

- `source requirement`: `docs/spec.md` の `<10gold>` など埋め込み要素保持要件、`tasks/usecases/body-translation-phase.yaml` の保護要素検証結果を確認する要件。
- `viewpoint`: adapter 境界 / network 応答不正
- `candidate scenario id`: `CAND-BTP-005`
- `external boundary`: AIProvider response adapter boundary、protection validation boundary、phase result boundary。
- `actor`: 翻訳ジョブ実行者
- `trigger`: fake provider が保護要素欠落、改変、順序不一致、対応不能な応答を返す。
- `start condition`: invalid provider response または protection failure response を fake transport で返せる。
- `expected outcome`: 保護要素を破壊した訳文は successful translated field として保存されない。phase result には field 単位の protection failure、retryable flag、error summary が secret なしで残る。
- `observable point`: protection validation result、field status、phase result、row count、error summary。
- `fake_or_stub`: invalid response fixture、protection failure fixture、temp DB、failure injection。
- `related detail requirement type`: `failure_handling_requirement`, `external_integration_requirement`, `data_requirement`
- `adoption hint`: provider 応答不正と保護要素 validation を 1 つの外部 adapter 失敗候補にできる。
- `conflict hint`: failure 観点の recoverable failure と統合される可能性がある。field 単位失敗か phase 全体失敗かは designer の判断対象である。

### CAND-BTP-006 Batch API 応答を翻訳フィールドへ欠落なく対応付ける

- `source requirement`: `docs/spec.md` の Gemini / xAI Batch API 利用要件、翻訳単位の lossless 保持要件、`tasks/usecases/body-translation-phase.yaml` の訳文と出力ステータス確認要件。
- `viewpoint`: provider 境界 / adapter 境界
- `candidate scenario id`: `CAND-BTP-006`
- `external boundary`: Batch API adapter boundary、provider response correlation boundary、progress boundary。
- `actor`: 翻訳ジョブ実行者
- `trigger`: Gemini または xAI 相当の batch execution mode で複数翻訳フィールドの応答を取得する。
- `start condition`: batch item に field correlation key、request unit ID、field ID を含められる。
- `expected outcome`: batch 応答の順序差、部分欠落、余剰項目を検出し、翻訳 field と訳文の対応を失わない。別 provider への暗黙 fallback は前提にしない。
- `observable point`: batch response adapter output、correlation error、phase progress、provider execution summary。
- `fake_or_stub`: batch request fixture、batch response fixture、partial response fixture、correlation key fixture。
- `related detail requirement type`: `external_integration_requirement`, `consistency_requirement`, `failure_handling_requirement`
- `adoption hint`: batch 対応 provider の adapter 契約候補にできる。
- `conflict hint`: batch 部分欠落を field 単位の retry 待ちにするか phase 全体の recoverable failure にするかは failure / state-transition 観点と競合しうる。

### CAND-BTP-007 network timeout と provider 到達不能を paid API なしで検証する

- `source requirement`: `docs/spec.md` の翻訳ジョブ中断、再開、失敗回復、API 実行進捗確認要件、`docs/architecture.md` の `AIProvider` adapter 境界。
- `viewpoint`: network 境界
- `candidate scenario id`: `CAND-BTP-007`
- `external boundary`: network boundary、AIProvider transport boundary、phase failure boundary。
- `actor`: 翻訳ジョブ実行者
- `trigger`: 本文翻訳 provider request が timeout、到達不能、rate limit、HTTP 相当の通信失敗を返す。
- `start condition`: fake transport で network failure を injection できる。
- `expected outcome`: paid real API を呼ばずに network failure を観測できる。phase result は error kind、retryable flag、progress、provider request 未完了数を secret なしで返す。
- `observable point`: fake transport log、phase result、error kind、progress、Job Run summary。
- `fake_or_stub`: fake transport、timeout fixture、unreachable fixture、rate limit fixture。
- `related detail requirement type`: `external_integration_requirement`, `recovery_requirement`, `testability_requirement`
- `adoption hint`: real provider list を保ったまま transport だけ fake にする lower-level acceptance 候補にできる。
- `conflict hint`: lifecycle / failure 観点の pause、resume、retry、cancel と状態遷移の最終扱いが競合しうる。

### CAND-BTP-008 確定訳語とジョブ内辞書を provider request で再翻訳対象にしない

- `source requirement`: `docs/spec.md` の共通辞書または用語翻訳済み対象を翻訳対象にせず置き換える要件、`term-translation-phase` の確定訳語とジョブ内辞書を本文翻訳フェーズの入力として参照する設計、`tasks/usecases/body-translation-phase.yaml` の入力一式。
- `viewpoint`: adapter 境界
- `candidate scenario id`: `CAND-BTP-008`
- `external boundary`: job dictionary reuse boundary、AIProvider request filtering boundary、output status boundary。
- `actor`: 翻訳ジョブ実行者
- `trigger`: 翻訳フィールド本文内に確定訳語またはジョブ内辞書 hit が含まれる。
- `start condition`: job dictionary snapshot と record type context を phase 開始時に固定できる。
- `expected outcome`: 確定訳語や辞書 hit は provider に自由翻訳させる対象ではなく、固定訳語制約または置換対象として扱われる。provider request summary と phase result で辞書適用件数を観測できる。
- `observable point`: provider request summary、dictionary snapshot digest、applied term count、field output status。
- `fake_or_stub`: job dictionary fixture、confirmed term fixture、prompt builder fixture、fake transport request capture。
- `related detail requirement type`: `external_integration_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: term phase の成果物を body phase provider request へ接続する候補にできる。
- `conflict hint`: 辞書 hit を provider request から完全除外するか、固定訳語制約として含めるかは人間判断候補になりうる。

### CAND-BTP-009 AI 実行の監査要約を残し raw prompt / raw response / 本文全文を保存しない

- `source requirement`: `docs/spec.md` の API 進捗確認、翻訳補助データの UI 観測、APIKey 暗号化要件、`term-translation-phase` と `persona-generation-phase` の redaction 設計。
- `viewpoint`: secret 境界 / adapter 境界
- `candidate scenario id`: `CAND-BTP-009`
- `external boundary`: runtime log boundary、audit summary boundary、AIProvider request / response boundary。
- `actor`: 翻訳ジョブ実行者
- `trigger`: provider 呼び出しを伴う success、invalid response、network failure のいずれかを実行する。
- `start condition`: provider、model、execution mode、credential ref、input count、output count、prompt digest、dictionary snapshot digest、persona snapshot digest を summary として生成できる。
- `expected outcome`: 障害調査に必要な provider 実行要約を確認できる。raw prompt、raw response、翻訳フィールド本文全文、secret は UI、DB summary、structured log に出ない。
- `observable point`: structured log、Job Run summary、fake transport log、redaction assertion。
- `fake_or_stub`: fake transport、log capture、redaction assertion fixture、fake secret store。
- `related detail requirement type`: `security_requirement`, `observability_requirement`, `external_integration_requirement`
- `adoption hint`: operation-audit 観点との統合候補にできる。
- `conflict hint`: prompt 調整用 debug log の扱いは、本文全文保存禁止と競合しうる。

### CAND-BTP-010 出力成果物へ渡る status と xTranslator 互換 status を混同しない

- `source requirement`: `docs/spec.md` の内部 `cached` と xTranslator `Status=1` 写像、訳文と出力ステータス lossless 保持、`translation-output-artifact.yaml` の訳文と出力ステータス入力要件。
- `viewpoint`: ファイル境界 / adapter 境界
- `candidate scenario id`: `CAND-BTP-010`
- `external boundary`: xTranslator compatibility adapter boundary、body phase output status boundary、translation-output-artifact handoff boundary。
- `actor`: 翻訳ジョブ実行者
- `trigger`: 本文翻訳フェーズが provider 翻訳、辞書置換、保護要素検証失敗、未処理 field を含む phase result を作る。
- `start condition`: `JOB_TRANSLATION_FIELD` に訳文、出力ステータス、保護要素検証結果が field 単位で保持される。
- `expected outcome`: body phase 内部 status は後続 output artifact で xTranslator status へ写像可能な形で保持される。provider translated、cached / dictionary replacement、failed / protection failed を混同しない。
- `observable point`: `JOB_TRANSLATION_FIELD`、phase result、output handoff summary、xTranslator output mapping fake。
- `fake_or_stub`: job translation field fixture、output mapping fake、temp DB。
- `related detail requirement type`: `external_integration_requirement`, `compatibility_requirement`, `data_requirement`
- `adoption hint`: body phase と translation-output-artifact の境界候補にできる。
- `conflict hint`: translation-output-artifact 観点が最終 status mapping を確定するため、本候補は handoff 境界までに留める。

## Open Notes

- `human decision candidate`: 本文翻訳フェーズの provider / model / execution mode を Job Setup から継承するか、phase 開始時に再選択させるかは指定資料だけでは確定しない。
- `human decision candidate`: 確定訳語とジョブ内辞書 hit を provider request から完全除外するか、固定訳語制約として provider request に含めるかは指定資料だけでは確定しない。
- `human decision candidate`: provider 応答で一部 field だけ失敗した場合、field 単位 retry、phase 全体 RecoverableFailed、partial success のどれを UI に出すかは final scenario 統合時の判断が必要である。
- `human decision candidate`: prompt 調整用 debug log に本文抜粋を許容するか、digest / ID だけに限定するかは redaction 方針として人間判断が必要である。
- `merge candidate`: `CAND-BTP-002` と `CAND-BTP-009` は secret redaction と audit summary の統合候補である。
- `merge candidate`: `CAND-BTP-003` と `CAND-BTP-004` は prompt/input adapter と response adapter の連続候補である。
- `merge candidate`: `CAND-BTP-005` と `CAND-BTP-007` は failure / state-transition 観点の recoverable failure 候補と統合される可能性がある。
- `merge candidate`: `CAND-BTP-010` は translation-output-artifact の scenario と統合される可能性がある。
- `rejection candidate`: paid real API 呼び出しを前提にする検証案は fake transport 方針と衝突するため designer 側の不採用候補になりうる。
