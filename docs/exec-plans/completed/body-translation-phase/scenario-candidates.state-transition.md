# Scenario Candidates: body-translation-phase / state-transition

- `generator`: `state-transition`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `BTP`
- `candidate_count`: 12

## Generator Scope

- `viewpoint`: state-transition
- `included_sources`:
  - `./plan.md`
  - `tasks/index.yaml`
  - `tasks/usecases/body-translation-phase.yaml`
  - `tasks/usecases/term-translation-phase.yaml`
  - `tasks/usecases/persona-generation-phase.yaml`
  - `tasks/usecases/translation-output-artifact.yaml`
  - `docs/spec.md`
  - `docs/er.md`
  - `docs/architecture.md`
  - `docs/exec-plans/completed/term-translation-phase/scenario-design.md`
  - `docs/exec-plans/completed/persona-generation-phase/scenario-design.md`
- `excluded_sources`:
  - 最終シナリオ表の確定
  - 候補の採否
  - 他観点候補との統合判断
  - プロダクトコード変更
  - プロダクトテスト変更
  - docs 正本化
- `generation_notes`:
  - 状態、許可遷移、禁止遷移、冪等再実行だけを候補化する。
  - UI 操作列、外部 provider 詳細、失敗原因分類、監査保存項目の最終判断は扱わない。
  - body phase 完了後の job-level `Completed` 反映時点は人間判断候補として残す。

## Candidate Scenarios

### CAND-BTP-001 persona phase 完了後だけ body phase を開始する

- `source requirement`: `tasks/usecases/body-translation-phase.yaml` は precondition を NPC ペルソナ生成フェーズ完了とし、`docs/spec.md` は単語翻訳フェーズ、NPC ペルソナ生成フェーズ、本文翻訳フェーズの順序を求める。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-BTP-001`
- `actor`: Job Run から body phase を開始するユーザーまたは phase start boundary
- `trigger`: persona phase Completed、persona snapshot 参照可能、ジョブ内辞書参照可能、非 terminal job、active phase run なしの job で body phase start を要求する。
- `pre-state`: term phase Completed、persona phase Completed、body phase not started、job は非 terminal。
- `start condition`: `JOB_PHASE_RUN` 群から前段 phase 完了と active phase なしを確認できる。
- `post-state`: body phase の `JOB_PHASE_RUN` が Running になり、current phase と progress を確認できる。
- `expected outcome`: body phase が current phase として観測でき、対象翻訳フィールドの input summary が生成される。
- `observable point`: phase start result、`JOB_PHASE_RUN`、Job Run current phase、progress、body input summary。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `observability_requirement`
- `adoption hint`: PGP の body readiness 候補と接続し、開始条件の正本候補にできる。
- `conflict hint`: persona snapshot 参照不能時の扱いは PGP 側候補と衝突しうる。

### CAND-BTP-002 body phase 開始条件を満たさない job では phase run を作らない

- `source requirement`: `tasks/usecases/body-translation-phase.yaml` は persona phase 完了を precondition にし、completed PGP scenario は persona 未完了、失敗、snapshot 参照不能では body phase run を作らない候補を持つ。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-BTP-002`
- `actor`: phase start boundary
- `trigger`: persona 未完了、persona RecoverableFailed、snapshot missing、ジョブ内辞書参照不能、既存 active phase run あり、または terminal job で body phase start を要求する。
- `pre-state`: body phase start 不許可の job state または phase state。
- `start condition`: 開始拒否条件のいずれかが成立している。
- `post-state`: `JOB_PHASE_RUN`、`PHASE_RUN_TRANSLATION_FIELD`、`JOB_TRANSLATION_FIELD` は変更されない。
- `expected outcome`: 開始不可理由を確認でき、phase run 重複と terminal job 後書きを防ぐ。
- `observable point`: phase transition result、phase run 件数、拒否理由、state snapshot。
- `related detail requirement type`: `state_requirement`, `boundary_requirement`, `consistency_requirement`
- `adoption hint`: 禁止遷移候補として、開始許可候補と対にできる。
- `conflict hint`: failure 観点が persona missing を recoverable failure とするか blocked とするかで分類が衝突しうる。

### CAND-BTP-003 body phase 開始時に参照 snapshot を固定する

- `source requirement`: `tasks/usecases/body-translation-phase.yaml` は確定訳語、ジョブ内辞書、ジョブ内ペルソナ、翻訳補助メタデータを input とする。TTP scenario は共通辞書 snapshot 固定を採用し、PGP scenario は persona snapshot 参照を body phase input とする。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-BTP-003`
- `actor`: body phase execution boundary
- `trigger`: body phase start が許可された直後に input snapshot を確定する。
- `pre-state`: body phase Running、input snapshot 未固定。
- `start condition`: 参照対象のジョブ内辞書、確定訳語、persona snapshot、翻訳補助メタデータが読み取り可能である。
- `post-state`: body phase の入力 snapshot digest または summary が固定され、同じ phase run の再開で同じ入力を参照できる。
- `expected outcome`: 翻訳中に辞書や persona が変わっても、同一 phase run の結果が同じ入力基準で再現できる。
- `observable point`: body input summary、snapshot digest、`PHASE_RUN_DICTIONARY_ENTRY`、`PHASE_RUN_PERSONA`、`JOB_PHASE_RUN`。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: 再開、リトライ、監査の候補と統合しやすい。
- `conflict hint`: input snapshot の永続化場所は ER 正本に明示が薄く、designer 側で保持場所を決める必要がある。

### CAND-BTP-004 翻訳対象フィールドを phase run に紐づける

- `source requirement`: `docs/er.md` は `PHASE_RUN_TRANSLATION_FIELD` を、フェーズが対象にしたジョブ内翻訳フィールドとして扱う。`tasks/usecases/body-translation-phase.yaml` は翻訳フィールドを input とする。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-BTP-004`
- `actor`: body phase target selection boundary
- `trigger`: body phase Running で翻訳対象フィールドを確定する。
- `pre-state`: body phase Running、対象 field link 未作成。
- `start condition`: 翻訳フィールドが job に属し、出力対象として扱える。
- `post-state`: 対象 `JOB_TRANSLATION_FIELD` と body phase run の関連が作成され、対象外フィールドは理由付きで summary に残る。
- `expected outcome`: body phase の対象件数、対象外件数、処理済み件数が phase progress と整合する。
- `observable point`: `PHASE_RUN_TRANSLATION_FIELD`、target summary、progress、対象外理由。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `observability_requirement`
- `adoption hint`: progress と訳文保存候補の起点にできる。
- `conflict hint`: 対象外判定の詳細は actor-goal / failure 観点の候補と重複しうる。

### CAND-BTP-005 フィールド翻訳成功で訳文、出力ステータス、保護要素検証結果を更新する

- `source requirement`: `tasks/usecases/body-translation-phase.yaml` は outputs を訳文、出力ステータス、保護要素検証結果とする。`docs/spec.md` は埋め込み要素を損なわない翻訳と、訳文、出力ステータスの lossless 保持を求める。`docs/er.md` は翻訳結果と出力ステータスを `JOB_TRANSLATION_FIELD` に保持すると定義する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-BTP-005`
- `actor`: body translation field execution boundary
- `trigger`: 翻訳結果が保護要素検証を通過し、保存可能な状態になる。
- `pre-state`: 対象 `JOB_TRANSLATION_FIELD` は untranslated または pending、body phase Running。
- `start condition`: 翻訳結果、出力ステータス、保護要素検証結果が同じ対象フィールドへ対応付いている。
- `post-state`: `JOB_TRANSLATION_FIELD` に訳文、出力ステータス、保護要素検証結果が保存され、対象 field は translated / output-ready 相当の状態になる。
- `expected outcome`: 訳文と出力ステータスと保護要素検証結果を同じ field 単位で確認できる。
- `observable point`: `JOB_TRANSLATION_FIELD`、`PHASE_RUN_TRANSLATION_FIELD`、phase result、Job Run field summary。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: body phase の主要正常遷移候補にできる。
- `conflict hint`: 出力ステータス値の語彙は正本で未固定のため、designer が status vocabulary を確認する必要がある。

### CAND-BTP-006 保護要素検証が失敗した field を成功扱いにしない

- `source requirement`: `docs/spec.md` は `<10gold>` などの埋め込み要素を損なわない翻訳を求め、`tasks/usecases/body-translation-phase.yaml` は保護要素検証結果を output にする。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-BTP-006`
- `actor`: protected element validation boundary
- `trigger`: 翻訳結果が保護要素検証で不一致になる。
- `pre-state`: 対象 field は pending、body phase Running。
- `start condition`: 原文の保護要素と訳文の保護要素検証結果が不一致である。
- `post-state`: 対象 field は output-ready にならず、body phase は RecoverableFailed または field-level retryable state を観測可能にする。
- `expected outcome`: 不正な訳文が `Completed` 相当や出力可能状態へ遷移しない。
- `observable point`: protected element validation result、field state、phase result、retryable flag、訳文保存有無。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `recovery_requirement`
- `adoption hint`: 保護要素を state invariant として固定できる。
- `conflict hint`: failure 観点が扱う provider response 不正、保存失敗、回復単位と重複しうる。

### CAND-BTP-007 全対象 field 完了で body phase を Completed にする

- `source requirement`: `tasks/usecases/body-translation-phase.yaml` は訳文、出力ステータス、保護要素検証結果を確認できることを completion criteria に含める。`tasks/usecases/translation-output-artifact.yaml` は本文翻訳フェーズ完了を precondition にする。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-BTP-007`
- `actor`: body phase completion boundary
- `trigger`: すべての対象 field が訳文、出力ステータス、保護要素検証結果を持ち、失敗中 field が残らない。
- `pre-state`: body phase Running、対象 field の一部または全部が処理済み。
- `start condition`: target count と successful field count が一致し、保護要素検証失敗と保存失敗が残っていない。
- `post-state`: body phase の `JOB_PHASE_RUN` が Completed になり、translation-output-artifact の入力にできる body phase result が成立する。
- `expected outcome`: Output Review 側が参照できる訳文、出力ステータス、保護要素検証結果の summary が成立する。
- `observable point`: `JOB_PHASE_RUN`、`JOB_TRANSLATION_FIELD`、phase result、completed field count、output readiness。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: 後続 output artifact の precondition 候補にできる。
- `conflict hint`: body phase Completed と job-level `Completed` を同時にするかは人間判断が必要である。

### CAND-BTP-008 body phase 完了前は output artifact へ進めない

- `source requirement`: `tasks/usecases/translation-output-artifact.yaml` は completed job、訳文、出力ステータスを input とし、precondition を本文翻訳フェーズ完了にする。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-BTP-008`
- `actor`: output artifact readiness boundary
- `trigger`: body phase 未開始、Running、Paused、RecoverableFailed、または保護要素検証失敗ありの job で output artifact readiness を確認する。
- `pre-state`: body phase Completed ではない、または output-ready field summary が成立していない。
- `start condition`: output artifact readiness の要求がある。
- `post-state`: output artifact の phase run または artifact は作成されず、拒否理由を確認できる。
- `expected outcome`: 未完了または不整合な訳文を xTranslator 互換成果物へ進めない。
- `observable point`: readiness result、artifact row count、拒否理由、body phase state。
- `related detail requirement type`: `state_requirement`, `compatibility_requirement`, `consistency_requirement`
- `adoption hint`: output artifact の入口候補と統合できる。
- `conflict hint`: output artifact 生成を別 task で扱うため、ここでは readiness 候補に留める。

### CAND-BTP-009 pause / resume は同じ body phase run を継続する

- `source requirement`: `tasks/usecases/body-translation-phase.yaml` は pause、resume、retry、cancel の可否確認を completion criteria に含める。`docs/spec.md` は翻訳ジョブの中断と再開を状態遷移として定義する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-BTP-009`
- `actor`: job control boundary
- `trigger`: body phase Running で pause、Paused で resume を要求する。
- `pre-state`: Running body phase または Paused body phase。
- `start condition`: job が terminal ではなく、対象 body phase run が存在する。
- `post-state`: pause では phase state が Paused になり、resume では同じ phase run が Running に戻る。
- `expected outcome`: phase run ID、input snapshot、処理済み field、progress が維持される。
- `observable point`: phase run ID、phase state、progress、processed field count、Job Run button enablement。
- `related detail requirement type`: `state_requirement`, `recovery_requirement`, `冪等性_requirement`
- `adoption hint`: lifecycle 候補と統合して user operation の state guard にできる。
- `conflict hint`: pause 時に provider 実行中 request をどう扱うかは external-integration / failure 観点と衝突しうる。

### CAND-BTP-010 RecoverableFailed の retry は未処理 field だけ進める

- `source requirement`: `docs/spec.md` は RecoverableFailed から Running への再開 / リトライを定義する。TTP と PGP の completed scenario は retry で既存成功分を維持し、未処理だけ進める候補を持つ。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-BTP-010`
- `actor`: phase retry boundary
- `trigger`: body phase RecoverableFailed で retry を要求する。
- `pre-state`: body phase RecoverableFailed、成功済み field と未処理または失敗 field が混在している。
- `start condition`: failure が retryable で、job が terminal ではない。
- `post-state`: 同じ `JOB_PHASE_RUN` が Running に戻り、成功済み `JOB_TRANSLATION_FIELD` は維持され、未処理または失敗 field だけが再実行対象になる。
- `expected outcome`: retry により成功済み訳文、出力ステータス、保護要素検証結果を重複作成または破損しない。
- `observable point`: phase run ID、field row count、成功済み field count、retry target count、latest error。
- `related detail requirement type`: `state_requirement`, `recovery_requirement`, `冪等性_requirement`
- `adoption hint`: 再実行時の二重保存防止候補にできる。
- `conflict hint`: failure 観点が field-level failure と phase-level RecoverableFailed の切り方を決める必要がある。

### CAND-BTP-011 開始再送や同一 phase 再実行で二重 phase run を作らない

- `source requirement`: `docs/er.md` はフェーズ再実行を同じ `JOB_PHASE_RUN` の状態を戻す扱いにし、Attempt 履歴テーブルを持たないと定義する。TTP / PGP の completed scenario も same phase run reuse を候補にしている。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-BTP-011`
- `actor`: phase start / resume / retry boundary
- `trigger`: active body phase run がある job で start を再送する、または同じ body phase を再実行する。
- `pre-state`: body phase の `JOB_PHASE_RUN` が既に存在する。
- `start condition`: 同一 job、同一 phase type の再要求である。
- `post-state`: 新しい `JOB_PHASE_RUN` は作成されず、既存 phase run の状態を返すか、許可された状態戻しだけを行う。
- `expected outcome`: `PHASE_RUN_TRANSLATION_FIELD` と `JOB_TRANSLATION_FIELD` の重複更新を防ぐ。
- `observable point`: phase run ID、phase run row count、field link row count、progress、state transition result。
- `related detail requirement type`: `state_requirement`, `冪等性_requirement`, `consistency_requirement`
- `adoption hint`: API 再送、画面二重クリック、retry の共通 invariant にできる。
- `conflict hint`: 完全な再翻訳を許可する場合の trigger と区別する必要がある。

### CAND-BTP-012 cancel または terminal job では body translation の後書きを拒否する

- `source requirement`: `docs/spec.md` は Ready / Paused から Canceled、Failed / Canceled / Completed の terminal 終了を状態遷移として示す。PGP scenario は terminal job への後書きを拒否する候補を持つ。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-BTP-012`
- `actor`: job control boundary、body translation persistence boundary
- `trigger`: body phase Running / Paused で cancel する、または Completed / Failed / Canceled job へ field result 保存を試行する。
- `pre-state`: Running / Paused body phase、または terminal job。
- `start condition`: cancel 要求、または terminal job に対する body phase write 要求がある。
- `post-state`: cancel では job または phase が terminal 相当の canceled state へ遷移する。terminal job への訳文、出力ステータス、保護要素検証結果の後書きは拒否される。
- `expected outcome`: 終了済み job の翻訳結果が後から変わらない。
- `observable point`: state snapshot、row count、拒否理由、Job Run terminal summary。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: terminal guard の state invariant 候補にできる。
- `conflict hint`: cancel 時に partial translated field を保持するか破棄するかは人間判断候補になりうる。

## Open Notes

- `human decision candidate`:
  - `HD-BTP-001`: body phase Completed と同時に job-level `Completed` にするか、output artifact task まで job を Running または別状態で維持するか。
  - `HD-BTP-002`: body phase の対象 field が 0 件の場合、provider 未実行の Completed とするか、開始不可または blocked とするか。
  - `HD-BTP-003`: 保護要素検証失敗時の出力ステータス語彙をどうするか。例として retryable failure、validation failed、untranslated が考えられるが、候補段階では確定しない。
  - `HD-BTP-004`: cancel 時の成功済み field result を保持するか、output artifact readiness から必ず除外するか。
- `merge candidate`:
  - `CAND-BTP-001` と `CAND-BTP-002` は phase start / forbidden transition として統合できる。
  - `CAND-BTP-005` と `CAND-BTP-006` は field result state invariant として統合できる。
  - `CAND-BTP-009`、`CAND-BTP-010`、`CAND-BTP-011` は same phase run reuse として統合できる。
  - `CAND-BTP-007` と `CAND-BTP-008` は output artifact readiness 境界として統合できる。
- `rejection candidate`:
  - `CAND-BTP-003` は input snapshot の保持場所が design で不要と判断された場合、監査候補へ移す余地がある。
  - `CAND-BTP-008` は output artifact 側 candidate へ完全移管する場合、body phase 候補から外せる。
- `conflict candidate`:
  - `CONFLICT-BTP-001`: job-level `Completed` の遷移時点が `docs/spec.md` の「翻訳完了」と output artifact task の「completed job」入力で曖昧である。
  - `CONFLICT-BTP-002`: protected element validation failure を field-level retryable state にするか、phase-level RecoverableFailed にするかが未確定である。
  - `CONFLICT-BTP-003`: cancel 後に partial result を保存済み内部状態として残す場合、output artifact readiness と衝突しうる。
