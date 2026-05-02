# Scenario Candidates: persona-generation-phase / external-integration

- `generator`: `external-integration`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `PGP`

## Generator Scope

- `viewpoint`: AI provider、credential、execution mode、fake transport、prompt/input/output 境界、network failure に限定する。
- `included_sources`: `./plan.md`, `tasks/index.yaml`, `tasks/usecases/persona-generation-phase.yaml`, `tasks/usecases/term-translation-phase.yaml`, `tasks/usecases/body-translation-phase.yaml`, `docs/spec.md`, `docs/er.md`, `docs/architecture.md`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md`
- `excluded_sources`: 他観点 candidate、最終 scenario 表、implementation-scope、product code、product test、docs 正本更新。
- `generation_notes`: 採否、統合、競合解消は行わない。paid real API を検証前提にしない。

## Candidate Scenarios

### CAND-PGP-001 ペルソナ生成用 AI provider を選択して実行できる

- `source requirement`: `docs/spec.md:49-58`, `docs/spec.md:98-101`, `docs/spec.md:130-131`, `docs/architecture.md:129-138`, `docs/er.md:25`
- `viewpoint`: provider 境界
- `candidate scenario id`: `CAND-PGP-001`
- `actor`: Job Run の NPC ペルソナ生成フェーズ実行者
- `trigger`: 単語翻訳フェーズ完了後、NPC ペルソナ生成フェーズを開始する。
- `expected outcome`: 選択済み provider、model、execution mode を使って NPC ペルソナ生成の AI 実行を開始し、`JOB_PHASE_RUN` にフェーズ別 AI 設定と最終 AI 実行情報を観測できる。
- `observable point`: Job Run の current phase、progress、provider / model / execution mode summary、`JOB_PHASE_RUN`。
- `fake_or_stub`: fake AI provider、provider selection fixture、phase run fixture。
- `related detail requirement type`: `success_requirement`, `data_requirement`, `testability_requirement`
- `acceptance relevance`: ペルソナ生成が本文翻訳フェーズの入力を作るため、provider 選択が phase 実行の開始条件になる。
- `adoption hint`: provider 呼び出し開始と phase run 記録の正常系候補として扱える。
- `conflict hint`: provider / model / execution mode を Job Setup から継承するか、phase 開始時に選ぶかは designer 側の人間判断候補になりうる。

### CAND-PGP-002 保存済み credential を参照し secret を露出しない

- `source requirement`: `docs/spec.md:57-58`, `docs/er.md:84`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:93-98`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:346-368`
- `viewpoint`: secret 境界
- `candidate scenario id`: `CAND-PGP-002`
- `actor`: Job Run の NPC ペルソナ生成フェーズ実行者
- `trigger`: API key 保存済み provider で NPC ペルソナ生成フェーズを実行する。
- `expected outcome`: API key の再入力なしで credential を参照でき、UI、error summary、structured log、fake transport log に API key、secret 本体、復号可能値が出ない。
- `observable point`: credential ref、fake secret store assertion、Job Run summary、structured log。
- `fake_or_stub`: fake secret store、fake transport、redaction assertion fixture。
- `related detail requirement type`: `security_requirement`, `observability_requirement`, `testability_requirement`
- `acceptance relevance`: AI provider 接続が必要なフェーズでも secret 非露出が成立することを受け入れ条件化できる。
- `adoption hint`: security / observability の外部連携候補として扱える。
- `conflict hint`: operation-audit 観点が保存したい監査情報と、secret / raw payload 保存禁止が衝突しうる。

### CAND-PGP-003 NPC 発話、属性、会話文脈、共通ペルソナを prompt 入力へ写像する

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:9-18`, `docs/spec.md:21-25`, `docs/spec.md:34-35`, `docs/spec.md:214-216`, `docs/spec.md:229`, `docs/er.md:50-55`
- `viewpoint`: prompt / input 境界
- `candidate scenario id`: `CAND-PGP-003`
- `actor`: NPC ペルソナ生成フェーズ実行者
- `trigger`: 生成対象 NPC の原文発話、NPC 属性メタデータ、会話文脈、共通ペルソナを持つジョブでフェーズを開始する。
- `expected outcome`: prompt input は NPC ごとの原文発話、属性、会話文脈、共通ペルソナ参照から構成され、生成根拠の翻訳フィールドを `PERSONA_FIELD_EVIDENCE` として追跡できる。
- `observable point`: fake transport request summary、prompt digest、input count、`PERSONA_FIELD_EVIDENCE`。
- `fake_or_stub`: prompt builder fixture、NPC record fixture、common persona fixture、fake transport request capture。
- `related detail requirement type`: `data_requirement`, `consistency_requirement`, `testability_requirement`
- `acceptance relevance`: 本文翻訳フェーズが参照する persona snapshot の品質は、AI へ渡す入力境界の欠落検知に依存する。
- `adoption hint`: prompt / input contract の候補として扱える。
- `conflict hint`: prompt 全文を保存して検証したい要求は、secret / raw payload / 過剰本文保存禁止と衝突しうる。

### CAND-PGP-004 provider 応答をジョブ内ペルソナと persona snapshot へ写像する

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:16-25`, `tasks/usecases/body-translation-phase.yaml:9-16`, `docs/spec.md:23-24`, `docs/spec.md:34`, `docs/er.md:50-56`, `docs/er.md:69`
- `viewpoint`: output adapter 境界
- `candidate scenario id`: `CAND-PGP-004`
- `actor`: NPC ペルソナ生成フェーズ実行者
- `trigger`: fake provider が有効なペルソナ応答を返す。
- `expected outcome`: provider 応答はジョブ内ペルソナと翻訳時参照用 persona snapshot へ変換され、共通ペルソナと分離して保持される。
- `observable point`: adapter output、`PERSONA`、`PHASE_RUN_PERSONA`、persona snapshot summary、Job Run phase result。
- `fake_or_stub`: fixed provider response fixture、persona response adapter fixture、temp DB。
- `related detail requirement type`: `success_requirement`, `data_requirement`, `consistency_requirement`
- `acceptance relevance`: 生成済みペルソナを本文翻訳フェーズの入力として参照できることを確認する候補になる。
- `adoption hint`: valid response の保存と後続入力化の正常系候補として扱える。
- `conflict hint`: 共通ペルソナが存在する NPC をジョブ内生成対象に含めるかは state-transition / lifecycle 観点と競合しうる。

### CAND-PGP-005 単発実行と Batch API の実行単位を観測できる

- `source requirement`: `docs/spec.md:51-52`, `docs/spec.md:55-57`, `docs/spec.md:98-101`, `docs/er.md:25`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:65-70`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:389-399`
- `viewpoint`: execution mode 境界
- `candidate scenario id`: `CAND-PGP-005`
- `actor`: Job Run の NPC ペルソナ生成フェーズ実行者
- `trigger`: Gemini または xAI の Batch API 対応 execution mode、または単発 execution mode でフェーズを開始する。
- `expected outcome`: execution mode ごとに request unit、input count、output count、provider job id または batch item id の要約を観測でき、失敗時に別 provider へ暗黙 fallback しない。
- `observable point`: fake transport log、phase progress、provider execution summary、structured log。
- `fake_or_stub`: single request fixture、batch request fixture、batch item response fixture。
- `related detail requirement type`: `boundary_requirement`, `observability_requirement`, `failure_handling_requirement`
- `acceptance relevance`: 大量 NPC の progress と provider 実行進捗を一致させるための候補になる。
- `adoption hint`: execution mode 差分の contract 候補として扱える。
- `conflict hint`: persona 生成の request unit を NPC 単位にするか、会話単位にするかは未決であり、designer の質問票候補になりうる。

### CAND-PGP-006 provider failure、timeout、invalid response を回復可能失敗として扱う

- `source requirement`: `docs/spec.md:52-54`, `docs/spec.md:133`, `docs/spec.md:171-199`, `tasks/usecases/body-translation-phase.yaml:27-28`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:265-293`
- `viewpoint`: network 境界
- `candidate scenario id`: `CAND-PGP-006`
- `actor`: Job Run の NPC ペルソナ生成フェーズ実行者
- `trigger`: provider 接続失敗、timeout、rate limit、invalid response、response 欠落を fake transport で発生させる。
- `expected outcome`: 失敗は成功扱いにならず、retryable flag、error kind、phase state、後続本文翻訳フェーズの開始不可理由を観測できる。
- `observable point`: phase result、error kind、retryable flag、progress、Job Run error summary。
- `fake_or_stub`: fake transport failure injection、timeout fixture、invalid response fixture、temp DB。
- `related detail requirement type`: `failure_handling_requirement`, `recovery_requirement`, `state_requirement`
- `acceptance relevance`: recoverable failure の情報確認と後続フェーズ制御を受け入れ条件化できる。
- `adoption hint`: failure 観点との統合候補として扱える。
- `conflict hint`: RecoverableFailed と Failed の分類、retry / resume の扱いは state-transition / failure 観点が最終整理する必要がある。

### CAND-PGP-007 AI 実行の監査要約を残し raw prompt / raw response を保存しない

- `source requirement`: `docs/spec.md:35`, `docs/spec.md:53-54`, `docs/spec.md:133`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:93-98`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md:346-372`, `docs/architecture.md:126-138`
- `viewpoint`: adapter / observability 境界
- `candidate scenario id`: `CAND-PGP-007`
- `actor`: Job Run の NPC ペルソナ生成フェーズ実行者
- `trigger`: provider 呼び出しを伴う成功、invalid response、network failure のいずれかを実行する。
- `expected outcome`: provider、model、execution mode、credential ref、input count、output count、prompt digest、error kind は観測でき、raw prompt、raw response、過剰な原文本文、secret は保存または表示されない。
- `observable point`: structured log、Job Run summary、fake transport log、redaction assertion。
- `fake_or_stub`: fake transport、log capture、redaction assertion fixture。
- `related detail requirement type`: `observability_requirement`, `security_requirement`, `compatibility_requirement`
- `acceptance relevance`: AI 実行進捗と失敗理由を観測しつつ、保存禁止情報を漏らさない受け入れ条件になる。
- `adoption hint`: operation-audit 観点との統合候補として扱える。
- `conflict hint`: prompt digest だけで再現性が不足する場合、operation-audit 観点の保存要求と security requirement が衝突しうる。

## Open Notes

- `human decision candidate`: ペルソナ生成フェーズの provider / model / execution mode を Job Setup から継承するか、phase 開始時に選択するかは未決候補である。
- `human decision candidate`: persona 生成の request unit を NPC 単位、会話単位、field evidence 単位のどれにするかは未決候補である。
- `merge candidate`: `CAND-PGP-002` と `CAND-PGP-007` は secret redaction と監査要約の統合候補である。
- `merge candidate`: `CAND-PGP-003` と `CAND-PGP-004` は prompt/input と output adapter の境界として連続する候補である。
- `rejection candidate`: paid real API 呼び出しを前提にした検証案は、fake transport 方針と衝突するため designer 側の不採用候補になりうる。
