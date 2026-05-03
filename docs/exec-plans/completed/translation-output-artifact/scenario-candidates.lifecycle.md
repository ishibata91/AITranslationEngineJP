# Scenario Candidates: translation-output-artifact / lifecycle

- `generator`: `lifecycle`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TOA`
- `candidate_count`: `10`

## Generator Scope

- `viewpoint`: `lifecycle`
- `included_sources`:
  - `docs/exec-plans/active/translation-output-artifact/plan.md`
  - `tasks/usecases/translation-output-artifact.yaml`
  - `docs/spec.md`
  - `docs/er.md`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md`
  - `docs/exec-plans/completed/body-translation-phase/implementation-scope.md`
- `excluded_sources`:
  - product code
  - product test
  - docs 正本変更
  - `.codex` 変更
  - 他 generator の候補成果物
- `generation_notes`: lifecycle 段階、開始条件、期待結果、観測点だけを候補化する。採否、統合、最終 scenario-design は designer に残す。

## Candidate Scenarios

### CAND-TOA-001 完了済み job を Output Review で選択する

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:9`
  - `tasks/usecases/translation-output-artifact.yaml:20-26`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:386-401`
  - `docs/exec-plans/completed/body-translation-phase/implementation-scope.md:41`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TOA-001`
- `lifecycle stage`: 選択
- `actor`: ユーザー
- `start condition`: body phase Completed で job-level `Completed` になり、output readiness が成立している。
- `trigger`: Output Review を開き、completed job を選択する。
- `expected outcome`: 完了済み job が一覧され、選択した job の出力準備状態を確認できる。
- `observable point`: Output Review UI、completed job list、output readiness、拒否理由。
- `related detail requirement type`: workflow / display
- `adoption hint`: Output Review の入口 scenario として使える。
- `conflict hint`: completed job の一覧条件を job-level `Completed` だけにするか、output readiness true まで含めるかは designer で統合する。

### CAND-TOA-002 出力開始前に readiness を確認する

- `source requirement`:
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:386-401`
  - `docs/exec-plans/completed/body-translation-phase/implementation-scope.md:75`
  - `docs/exec-plans/completed/body-translation-phase/implementation-scope.md:322-323`
  - `tasks/usecases/translation-output-artifact.yaml:20-26`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TOA-002`
- `lifecycle stage`: 出力開始
- `actor`: ユーザー
- `start condition`: selected job の body phase が Completed で、field result と output status が整合している。
- `trigger`: output artifact の生成開始を試行する。
- `expected outcome`: readiness が true の job だけ生成開始でき、未完了、失敗中、status 不整合では開始不可理由を返す。
- `observable point`: readiness result、output status summary、開始不可理由、button enablement。
- `related detail requirement type`: workflow / state guard
- `adoption hint`: output artifact 生成開始の前置き scenario として使える。
- `conflict hint`: readiness の所有境界は body phase 由来か output artifact 由来かを designer で整理する。

### CAND-TOA-003 result summary を確認してから出力する

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:15-19`
  - `tasks/usecases/translation-output-artifact.yaml:23-26`
  - `docs/spec.md:43`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:293-300`
  - `docs/exec-plans/completed/body-translation-phase/implementation-scope.md:369-373`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TOA-003`
- `lifecycle stage`: 生成前確認
- `actor`: ユーザー
- `start condition`: completed job と output artifact 候補を選択済みである。
- `trigger`: result summary を確認する。
- `expected outcome`: 翻訳単位の件数、訳文、出力ステータス、保護要素検証結果、出力可能状態を生成前に確認できる。
- `observable point`: result summary、field result summary、output status summary、Output Review UI。
- `related detail requirement type`: display / workflow
- `adoption hint`: generation 前の確認 step として使える。
- `conflict hint`: result summary の表示単位を field 単位にするか集約単位にするかは UI 設計側と統合する。

### CAND-TOA-004 diff preview と再生成操作を生成前に確認する

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:17-18`
  - `tasks/usecases/translation-output-artifact.yaml:24-26`
  - `tasks/usecases/translation-output-artifact.yaml:31-35`
  - `docs/spec.md:43`
  - `docs/er.md:76-77`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TOA-004`
- `lifecycle stage`: 生成前確認
- `actor`: ユーザー
- `start condition`: completed job の field result summary と output artifact 候補を取得済みである。
- `trigger`: diff preview を開く。
- `expected outcome`: 翻訳単位ごとの Source、Dest、Status 差分と、再生成操作の可否を確認できる。
- `observable point`: diff preview、field identity、再生成 button enablement、row preview count。
- `related detail requirement type`: display / re-output
- `adoption hint`: preview と再生成導線を同じ lifecycle stage の候補として扱える。
- `conflict hint`: diff preview が生成前 preview なのか、生成済み artifact との差分なのかは designer で確定する。

### CAND-TOA-005 xTranslator 互換 row を生成する

- `source requirement`:
  - `docs/spec.md:65-67`
  - `docs/spec.md:43`
  - `docs/er.md:71-77`
  - `tasks/usecases/translation-output-artifact.yaml:15-19`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:246-250`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TOA-005`
- `lifecycle stage`: 生成
- `actor`: システム
- `start condition`: output readiness が true で、生成対象の `JOB_TRANSLATION_FIELD` が確定している。
- `trigger`: ユーザーが xTranslator 互換成果物の出力を実行する。
- `expected outcome`: `EDID`、`REC`、`FIELD`、`FORMID`、`Source`、`Dest`、`Status` を持つ `XTRANSLATOR_OUTPUT_ROW` を生成できる。
- `observable point`: `TRANSLATION_ARTIFACT`、`XTRANSLATOR_OUTPUT_ROW`、row count、Status mapping summary。
- `related detail requirement type`: persistence / output format
- `adoption hint`: backend generation boundary の受け入れ候補として使える。
- `conflict hint`: row 順序、row identity、重複生成時の扱いは最終 scenario では未確定として扱う。

### CAND-TOA-006 cached status を xTranslator Status へ写像する

- `source requirement`:
  - `docs/spec.md:29-32`
  - `docs/spec.md:65-67`
  - `docs/er.md:76-77`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:246-250`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TOA-006`
- `lifecycle stage`: 生成
- `actor`: システム
- `start condition`: 出力対象 field に内部出力ステータス `cached` が含まれている。
- `trigger`: xTranslator 互換 row を生成する。
- `expected outcome`: 内部 `cached` は xTranslator の `Status=1` に写像され、辞書置換の内部観測情報は xTranslator `Status` とは別に保持される。
- `observable point`: row Status、internal observation summary、result summary。
- `related detail requirement type`: output format / status mapping
- `adoption hint`: status mapping の境界候補として CAND-TOA-005 に統合できる。
- `conflict hint`: `cached` 以外の内部出力ステータスと xTranslator `Status` の写像表は別観点または designer で補完が必要である。

### CAND-TOA-007 artifact 完了状態と出力履歴を確認する

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:16-19`
  - `tasks/usecases/translation-output-artifact.yaml:26`
  - `tasks/usecases/translation-output-artifact.yaml:31-35`
  - `docs/er.md:71-77`
  - `docs/spec.md:234-235`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TOA-007`
- `lifecycle stage`: 完了
- `actor`: ユーザー
- `start condition`: xTranslator 互換 row 生成が完了している。
- `trigger`: output artifact の状態を確認する。
- `expected outcome`: output artifact の完了状態、row count、result summary、再出力導線を確認できる。
- `observable point`: `TRANSLATION_ARTIFACT` state、artifact summary、Output Review UI、再出力導線。
- `related detail requirement type`: display / artifact lifecycle
- `adoption hint`: 生成後確認と終了後利用の候補として使える。
- `conflict hint`: artifact state の語彙、保存場所、履歴表示単位は根拠資料だけでは確定できない。

### CAND-TOA-008 生成済み artifact を再出力する

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:19`
  - `tasks/usecases/translation-output-artifact.yaml:24-26`
  - `tasks/usecases/translation-output-artifact.yaml:31-35`
  - `docs/spec.md:65-67`
  - `docs/er.md:73-77`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TOA-008`
- `lifecycle stage`: 再出力
- `actor`: ユーザー
- `start condition`: 生成済み output artifact が存在し、元 job が output readiness を維持している。
- `trigger`: 再出力操作を実行する。
- `expected outcome`: 再出力状態を確認でき、xTranslator 互換成果物を再生成できる。
- `observable point`: 再出力状態、artifact summary、row count、diff preview。
- `related detail requirement type`: re-output / artifact lifecycle
- `adoption hint`: re-output scenario の lifecycle 候補として使える。
- `conflict hint`: 再出力が既存 artifact の上書きか、新しい artifact version の作成かは人間判断候補として残す。

### CAND-TOA-009 本文翻訳対象 0 件でも成果物出力へ進む

- `source requirement`:
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:114-119`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:397-400`
  - `docs/exec-plans/completed/body-translation-phase/implementation-scope.md:36`
  - `docs/exec-plans/completed/body-translation-phase/implementation-scope.md:157`
  - `tasks/usecases/translation-output-artifact.yaml:20-26`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TOA-009`
- `lifecycle stage`: 特殊完了から出力
- `actor`: ユーザー
- `start condition`: 本文翻訳対象 0 件の job が body phase と job-level `Completed` になり、output readiness が成立している。
- `trigger`: output artifact の生成を開始する。
- `expected outcome`: provider 未実行の Completed job でも、出力可能な summary を確認し、成果物出力へ進める。
- `observable point`: output readiness、target count、skipped count、result summary、output artifact state。
- `related detail requirement type`: boundary / workflow
- `adoption hint`: 0 件境界の lifecycle candidate として使える。
- `conflict hint`: `XTRANSLATOR_OUTPUT_ROW` が `JOB_TRANSLATION_FIELD` に対応するため、単語だけの plugin で row source をどう扱うかは designer で解消する。

### CAND-TOA-010 Canceled または不整合 job から artifact を作らない

- `source requirement`:
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:30-32`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:360-377`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:397-400`
  - `docs/exec-plans/completed/body-translation-phase/implementation-scope.md:40`
  - `docs/exec-plans/completed/body-translation-phase/implementation-scope.md:319-323`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-TOA-010`
- `lifecycle stage`: 終了後の出力禁止
- `actor`: ユーザー
- `start condition`: job が `Canceled`、未完了、失敗中、または field result と output status が不整合である。
- `trigger`: output artifact の生成開始を試行する。
- `expected outcome`: artifact 生成は開始されず、途中成功結果は output readiness に使われない。
- `observable point`: output readiness false、開始不可理由、artifact 未生成、field status summary。
- `related detail requirement type`: state guard / artifact lifecycle
- `adoption hint`: lifecycle の終端保護 scenario として使える。
- `conflict hint`: state-transition 観点の禁止遷移候補と重複しうるため、designer で統合する。

## Open Notes

- `human decision candidate`:
  - 再出力は既存 artifact の上書きか、新しい artifact version の追加かを決める必要がある。
  - artifact state の語彙と保存場所を決める必要がある。
  - diff preview は生成前 preview か、生成済み artifact との差分表示かを決める必要がある。
  - `cached` 以外の内部出力ステータスと xTranslator `Status` の写像表を決める必要がある。
  - 本文翻訳対象 0 件または単語だけの plugin で、`XTRANSLATOR_OUTPUT_ROW` の row source をどう扱うかを決める必要がある。
- `merge candidate`:
  - CAND-TOA-001 と CAND-TOA-002 は開始条件確認として統合される可能性がある。
  - CAND-TOA-003 と CAND-TOA-004 は生成前確認 UI として統合される可能性がある。
  - CAND-TOA-005 と CAND-TOA-006 は xTranslator row 生成として統合される可能性がある。
  - CAND-TOA-007 と CAND-TOA-008 は artifact 完了後利用として統合される可能性がある。
- `rejection candidate`:
  - なし。採否は designer が判断する。
- `conflict candidate`:
  - body phase 由来の output readiness と output artifact 自体の lifecycle state の所有境界が競合しうる。
  - state-transition 観点の禁止遷移候補と CAND-TOA-010 が重複しうる。
  - display 観点の result summary / diff preview 候補と CAND-TOA-003 / CAND-TOA-004 が重複しうる。
