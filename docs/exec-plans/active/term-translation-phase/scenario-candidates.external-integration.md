# Scenario Candidates: term-translation-phase / external-integration

- `generator`: `external-integration`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TTP`

## Generator Scope

- `viewpoint`: external-integration
- `included_sources`: `./plan.md`, `tasks/usecases/term-translation-phase.yaml`, `tasks/index.yaml`, `docs/spec.md`, `docs/er.md`, `docs/architecture.md`, `docs/exec-plans/completed/translation-job-setup/plan.md`, `docs/exec-plans/completed/translation-job-setup/scenario-design.md`, `docs/exec-plans/completed/translation-job-setup/implementation-scope.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本化、他 generator の成果物、最終 scenario matrix
- `generation_notes`: 本文翻訳フェーズ前に単語翻訳フェーズを追加し、用語や固有名詞の訳語をジョブ内辞書へ反映する変更だけを対象にする。採否、統合、競合解消は designer に残す。

## Candidate Scenarios

### CAND-TTP-001 共通辞書の完全一致語を provider 呼び出しから除外する

- `source requirement`: `docs/spec.md` の「共通辞書として登録済み、または用語翻訳で翻訳済みの対象は翻訳対象とせず、置き換えること」「完全一致のみを構築済みとして扱うこと」、`tasks/usecases/term-translation-phase.yaml` の「共通辞書に存在する語を翻訳対象から除外できる」
- `viewpoint`: external-integration
- `candidate scenario id`: `CAND-TTP-001`
- `external boundary`: 共通辞書参照 boundary、AIProvider request boundary
- `actor`: 翻訳ジョブ実行者
- `trigger`: 単語翻訳フェーズ開始時に、翻訳対象語の中に共通辞書と完全一致する語が含まれる。
- `start condition`: Ready 相当の翻訳ジョブ、翻訳対象語、共通辞書参照が存在する。
- `expected outcome`: 完全一致した語は provider request へ含めず、辞書置換対象として観測できる。未一致語だけが AIProvider へ渡される。
- `fake_or_stub`: 共通辞書 fixture、翻訳対象語 fixture、fake AI transport
- `observable point`: provider request payload、phase result、ジョブ内辞書または置換対象判定結果
- `related detail requirement type`: `external_integration`, `dictionary_reuse`, `provider_request_filtering`
- `adoption hint`: 共通辞書の完全一致除外と provider 呼び出し payload の分離を 1 つの acceptance 候補にできる。
- `conflict hint`: state-transition / failure 観点の「対象語なしで phase を完了扱いにするか」と競合する可能性がある。

### CAND-TTP-002 用語翻訳 provider 応答をジョブ内辞書へ変換する

- `source requirement`: `docs/spec.md` の「事前に単語翻訳フェーズを実行し、会話文やクエストの本文翻訳フェーズに辞書として再利用できること」、`docs/er.md` の `DICTIONARY_ENTRY` と `PHASE_RUN_DICTIONARY_ENTRY`
- `viewpoint`: external-integration
- `candidate scenario id`: `CAND-TTP-002`
- `external boundary`: AIProvider response adapter、ジョブ内辞書 persistence boundary
- `actor`: 翻訳ジョブ実行者
- `trigger`: 未翻訳の用語や固有名詞に対して provider 応答が返る。
- `start condition`: provider 応答に原語、訳語、対象語 ID または対応付け可能な key が含まれる。
- `expected outcome`: provider 応答は内部契約へ写像され、ジョブ内 `DICTIONARY_ENTRY` と phase run の関連として保存候補になる。本文翻訳フェーズは確定訳語を参照できる。
- `fake_or_stub`: fixed provider response、job dictionary repository fake、temp DB
- `observable point`: adapter output、ジョブ内辞書 entry、phase run と辞書 entry の関連、本文翻訳フェーズ入力の参照可否
- `related detail requirement type`: `external_integration`, `adapter_mapping`, `job_dictionary_persistence`
- `adoption hint`: provider 固有応答を内部辞書 entry に変換する contract freeze 候補にできる。
- `conflict hint`: lifecycle 観点の「確定前の候補訳語をいつ確定扱いにするか」と競合する可能性がある。

### CAND-TTP-003 secret を露出せず用語翻訳 provider を呼び出す

- `source requirement`: `docs/spec.md` の「各フェーズのAPI選択、APIKeyは再入力不要で保存ができること」「APIKeyは暗号化して保存すること」、`translation-job-setup` の secret 非露出決定
- `viewpoint`: external-integration
- `candidate scenario id`: `CAND-TTP-003`
- `external boundary`: secret store boundary、AIProvider request boundary、runtime log boundary
- `actor`: 翻訳ジョブ実行者
- `trigger`: 単語翻訳フェーズが保存済み credential 参照を使って provider request を作る。
- `start condition`: provider、model、credential 参照、実行方式が phase 実行設定として解決できる。
- `expected outcome`: API key 平文、secret 本体、復号可能な値は UI、phase result、error summary、structured log、provider request 観測値へ出ない。保存済み credential 参照状態だけを表示または記録する。
- `fake_or_stub`: fake secret store、fake transport、redaction assertion fixture
- `observable point`: phase setting summary、request log、error summary、structured log、UI 表示
- `related detail requirement type`: `external_integration`, `secret_boundary`, `redaction`
- `adoption hint`: Job Setup の secret 非露出規約を phase 実行へ継承する候補にできる。
- `conflict hint`: operation-audit 観点の「どの実行情報を監査保存するか」と競合する可能性がある。

### CAND-TTP-004 provider capability と実行方式不整合を phase 開始前に検出する

- `source requirement`: `docs/spec.md` の「LMStudio、Gemini、xAIを翻訳AIとして利用できること」「Gemini、xAIはBatchAPIが利用できること」「各フェーズではいずれのモデルでも選択できる」、`translation-job-setup` の provider capability blocking 決定
- `viewpoint`: external-integration
- `candidate scenario id`: `CAND-TTP-004`
- `external boundary`: provider capability adapter、phase start validation boundary
- `actor`: 翻訳ジョブ実行者
- `trigger`: 単語翻訳フェーズ開始時に、選択 provider と実行方式の組み合わせを検証する。
- `start condition`: provider、model、実行方式、翻訳対象語件数が phase 実行設定に含まれる。
- `expected outcome`: provider / mode 不整合は provider request 前に blocking として観測できる。fake provider は user-facing provider list に出ない。
- `fake_or_stub`: provider capability fixture、fake transport
- `observable point`: validation result、blocking failure category、provider request 未実行証跡、phase state または phase result
- `related detail requirement type`: `external_integration`, `provider_capability`, `execution_mode_validation`
- `adoption hint`: phase 開始前 contract として provider capability check を採用候補にできる。
- `conflict hint`: failure 観点の「blocking failure を RecoverableFailed にするか、開始拒否にするか」と競合する可能性がある。

### CAND-TTP-005 network timeout と provider 到達不能を paid API なしで検証する

- `source requirement`: `docs/spec.md` の「翻訳ジョブの中断、再開、失敗回復が継続的に行えること」「翻訳ジョブ、APIの実行進捗を確認できること」、`translation-job-setup` の「paid な real AI API を scenario validation の前提にしない」
- `viewpoint`: external-integration
- `candidate scenario id`: `CAND-TTP-005`
- `external boundary`: network boundary、AIProvider transport boundary
- `actor`: 翻訳ジョブ実行者
- `trigger`: 単語翻訳 provider request が timeout、到達不能、HTTP 相当の通信失敗を返す。
- `start condition`: provider request が fake transport で失敗 injection できる。
- `expected outcome`: paid な real API を呼ばずに、timeout、到達不能、応答不正を観測できる。失敗理由は secret を含まず、再試行可否の判断材料を返す。
- `fake_or_stub`: fake transport、timeout fixture、unreachable fixture、invalid response fixture
- `observable point`: transport log、phase result、error kind、progress 表示、external request 未実行証跡
- `related detail requirement type`: `external_integration`, `network_failure`, `testability`
- `adoption hint`: real provider list を保ったまま transport だけ fake にする lower-level acceptance 候補にできる。
- `conflict hint`: failure / state-transition 観点の失敗状態、retry、resume の分類と競合する可能性がある。

### CAND-TTP-006 Batch API 応答を対象語へ欠落なく対応付ける

- `source requirement`: `docs/spec.md` の「Gemini、xAIはBatchAPIが利用できること」「失敗しても特に別プロバイダフォールバックは必要ない」、`tasks/usecases/term-translation-phase.yaml` の「確定訳語を本文翻訳フェーズの入力として参照できる」
- `viewpoint`: external-integration
- `candidate scenario id`: `CAND-TTP-006`
- `external boundary`: Batch API adapter、provider response correlation boundary
- `actor`: 翻訳ジョブ実行者
- `trigger`: Gemini または xAI 相当の batch 実行方式で、複数の対象語に対する応答を取得する。
- `start condition`: batch request に対象語 ID または対応付け key が含まれる。
- `expected outcome`: batch 応答の順序差、部分欠落、余剰項目を検出し、対象語と訳語の対応を失わない。別 provider への自動フォールバックは前提にしない。
- `fake_or_stub`: batch response fixture、partial response fixture、correlation key fixture
- `observable point`: batch request payload、adapter correlation result、missing item error、ジョブ内辞書保存候補
- `related detail requirement type`: `external_integration`, `batch_api_adapter`, `response_correlation`
- `adoption hint`: batch 対応 provider の adapter 契約候補にできる。
- `conflict hint`: failure 観点の「部分欠落時に全体失敗か部分成功か」と競合する可能性がある。

### CAND-TTP-007 provider 応答不正時にジョブ内辞書を汚さない

- `source requirement`: `docs/spec.md` の「用語や固有名詞の訳語を確定してジョブ内辞書へ反映する」、`docs/er.md` の `DICTIONARY_ENTRY` と `JOB_PHASE_RUN`、`translation-job-setup` の partial state を残さない方針
- `viewpoint`: external-integration
- `candidate scenario id`: `CAND-TTP-007`
- `external boundary`: AIProvider response adapter、persistence transaction boundary
- `actor`: 翻訳ジョブ実行者
- `trigger`: provider 応答が JSON 形式不正、必須 key 欠落、対象語との対応不能、空訳語を含む。
- `start condition`: invalid provider response を fake transport で返せる。
- `expected outcome`: 不正応答は辞書 entry として確定保存されない。phase result は provider 応答不正を secret なしで観測できる。
- `fake_or_stub`: invalid response fixture、temp DB、failure injection
- `observable point`: adapter error、辞書 row count、phase result、error kind
- `related detail requirement type`: `external_integration`, `adapter_validation`, `persistence_atomicity`
- `adoption hint`: AI 応答 adapter と辞書保存の境界条件として採用候補にできる。
- `conflict hint`: state-transition 観点の「phase を failed にするか未確定候補だけ残すか」と競合する可能性がある。

### CAND-TTP-008 xTranslator 由来共通辞書と用語翻訳結果の出力ステータスを混同しない

- `source requirement`: `docs/spec.md` の「xTranslator互換形式へ出力する場合、内部の cached は xTranslator の Status=1 に写像すること」「辞書置換であることは、xTranslator の Status とは別の内部観測情報として保持できること」、`docs/er.md` の `XTRANSLATOR_TRANSLATION_XML` と `DICTIONARY_ENTRY`
- `viewpoint`: external-integration
- `candidate scenario id`: `CAND-TTP-008`
- `external boundary`: xTranslator compatibility adapter、dictionary source boundary
- `actor`: 翻訳ジョブ実行者
- `trigger`: 共通辞書由来の置換語と単語翻訳フェーズ由来のジョブ内辞書 entry が同じ本文翻訳入力に混在する。
- `start condition`: xTranslator 由来共通辞書 entry と用語翻訳由来ジョブ内辞書 entry が存在する。
- `expected outcome`: 内部の辞書置換情報、`cached` 相当の出力ステータス、辞書 source は混同されない。xTranslator 互換出力への写像は後続出力フェーズへ渡せる。
- `fake_or_stub`: xTranslator import fixture、job dictionary fixture、output mapping fake
- `observable point`: dictionary source、replacement observation、output status mapping candidate、本文翻訳フェーズ入力
- `related detail requirement type`: `external_integration`, `xtranslator_compatibility`, `dictionary_source_mapping`
- `adoption hint`: 外部形式互換の source / status 混同防止候補にできる。
- `conflict hint`: operation-audit 観点の観測情報保存、lifecycle 観点の本文翻訳フェーズ入力生成と競合する可能性がある。

## Open Notes

- `human decision candidate`: 単語翻訳フェーズの AI 実行設定を Job Setup の設定から継承するか、phase 開始時に provider / model / execution mode を再選択させるかは、指定資料だけでは確定しない。
- `human decision candidate`: Batch API の部分欠落時に、全体失敗、部分成功、再試行待ちのどれをユーザーに提示するかは、failure / state-transition 統合時に判断が必要である。
- `merge candidate`: `CAND-TTP-004` と `CAND-TTP-005` は provider validation / network failure として failure 観点候補と統合される可能性がある。
- `merge candidate`: `CAND-TTP-001` と `CAND-TTP-008` は dictionary reuse / xTranslator compatibility として operation-audit 観点候補と統合される可能性がある。
- `rejection candidate`: 単語翻訳フェーズ内で xTranslator XML ファイルを直接読み直す候補は、指定資料では phase 入力に含まれないため採用しない候補として designer 判断へ残す。
