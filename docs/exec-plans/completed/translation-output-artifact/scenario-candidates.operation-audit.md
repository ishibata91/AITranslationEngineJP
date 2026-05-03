# Scenario Candidates: translation-output-artifact / operation-audit

- `generator`: `operation-audit`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TOA`
- `task_artifact_dir`: `docs/exec-plans/active/translation-output-artifact/`
- `candidate_artifact`: `docs/exec-plans/active/translation-output-artifact/scenario-candidates.operation-audit.md`
- `target_difference`: 完了ジョブの翻訳結果を確認し、xTranslator 互換の成果物として出力する。
- `candidate_count`: 9

## Generator Scope

- `viewpoint`: operation-audit
- `included_sources`:
  - `docs/exec-plans/active/translation-output-artifact/plan.md`
  - `tasks/usecases/translation-output-artifact.yaml`
  - `docs/spec.md`
  - `docs/er.md`
  - `docs/diagrams/er/combined-data-model-er.puml`
  - `docs/architecture.md`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md`
- `excluded_sources`:
  - プロダクトコード
  - プロダクトテスト
  - docs 正本本文の変更
  - `.codex` 配下の変更
  - 他観点のシナリオ候補生成
- `generation_notes`:
  - 採否、統合、最終 scenario-design は `designer` に残す。
  - 候補は運用確認、監査、履歴、再現性、保存禁止だけを扱う。
  - 監査保持期間、監査粒度、伏せ字範囲は候補として残し、AI だけで確定しない。

## Candidate Scenarios

### CAND-TOA-001 完了ジョブの成果物確認履歴を残す

- `source requirement`:
  - `docs/exec-plans/active/translation-output-artifact/plan.md:5-9`
  - `tasks/usecases/translation-output-artifact.yaml:20-27`
  - `docs/spec.md:103-121`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TOA-001`
- `actor`: 利用者、運用確認者
- `trigger`: Output Review を開き、完了済み翻訳ジョブと出力成果物を選択する。
- `audit event`: 完了ジョブの確認開始、対象ジョブ、対象成果物、確認結果。
- `saved summary`: job id、artifact id、job state、artifact status、result summary、確認時刻。
- `redaction rule`: result summary へ secret、API key、復号可能値を含めない。
- `expected outcome`: 完了済みジョブ、成果物状態、結果要約を後から確認できる。
- `observable point`: Output Review、gateway response、`TRANSLATION_ARTIFACT.status`。
- `related detail requirement type`: display / audit
- `adoption hint`: 完了ジョブ一覧と result summary の表示シナリオへ統合候補にする。
- `conflict hint`: ジョブが `Completed` 以外の場合の表示可否は state-transition 観点と競合しうる。

### CAND-TOA-002 xTranslator 出力行の再構成根拠を追跡する

- `source requirement`:
  - `docs/spec.md:65-67`
  - `docs/spec.md:222-235`
  - `docs/er.md:71-78`
  - `docs/diagrams/er/combined-data-model-er.puml:207-229`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TOA-002`
- `actor`: 利用者、障害調査者
- `trigger`: xTranslator 互換成果物を生成または確認する。
- `audit event`: 出力行の再構成、参照元 field、行数、生成結果。
- `saved summary`: artifact id、row count、`EDID` / `REC` / `FIELD` / `FORMID` の存在件数、status 分布。
- `redaction rule`: structured log には `Source` と `Dest` の全文一覧を重複保存しない。
- `expected outcome`: 各 `XTRANSLATOR_OUTPUT_ROW` が 1 つの `JOB_TRANSLATION_FIELD` に対応することを追跡できる。
- `observable point`: `TRANSLATION_ARTIFACT`、`XTRANSLATOR_OUTPUT_ROW`、row count、status distribution。
- `related detail requirement type`: persistence / audit
- `adoption hint`: xTranslator XML 生成の正常系シナリオへ統合候補にする。
- `conflict hint`: row-level の本文保持は ER 上は必要だが、監査ログへの全文複製は禁止候補にする。

### CAND-TOA-003 内部出力ステータスと xTranslator `Status` の写像を監査する

- `source requirement`:
  - `docs/spec.md:29-32`
  - `docs/spec.md:41-43`
  - `docs/spec.md:65-67`
  - `docs/er.md:23-24`
  - `docs/diagrams/er/combined-data-model-er.puml:156-165`
  - `docs/diagrams/er/combined-data-model-er.puml:218-229`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TOA-003`
- `actor`: 利用者、運用確認者
- `trigger`: 辞書置換または翻訳済み field を成果物へ出力する。
- `audit event`: 内部 `output_status` から xTranslator `Status` への写像。
- `saved summary`: input field count、cached count、translated count、status mapping summary。
- `redaction rule`: 辞書置換である事実は要約に残し、辞書 entry の過剰本文は監査ログへ複製しない。
- `expected outcome`: 内部 `cached` が xTranslator の `Status=1` へ写像されたことを確認できる。
- `observable point`: `JOB_TRANSLATION_FIELD.output_status`、`XTRANSLATOR_OUTPUT_ROW.status`、result summary。
- `related detail requirement type`: persistence / audit
- `adoption hint`: output status 変換の受け入れ条件へ統合候補にする。
- `conflict hint`: xTranslator `Status` と内部観測情報を同一概念に統合すると、辞書置換の監査情報が失われる。

### CAND-TOA-004 diff preview と再生成操作の履歴を残す

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:15-19`
  - `tasks/usecases/translation-output-artifact.yaml:23-26`
  - `tasks/usecases/translation-output-artifact.yaml:31-35`
  - `docs/spec.md:120-121`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TOA-004`
- `actor`: 利用者
- `trigger`: diff preview を確認し、translation unit 単位の再生成操作または再出力導線を使う。
- `audit event`: diff preview 確認、再生成開始、再出力要求、終了状態。
- `saved summary`: artifact id、対象 unit count、changed count、regeneration requested flag、re-output state。
- `redaction rule`: preview 本文の全文を監査ログへ保存せず、件数、field id、digest を保存対象にする。
- `expected outcome`: 再生成または再出力の理由と結果を後追いできる。
- `observable point`: diff preview、result summary、再出力状態、operation result。
- `related detail requirement type`: display / workflow / audit
- `adoption hint`: diff preview と再出力導線の表示系シナリオへ統合候補にする。
- `conflict hint`: translation unit 単位の再生成が本文翻訳フェーズの再実行を伴う場合、lifecycle 観点と競合しうる。

### CAND-TOA-005 同一完了ジョブからの再出力を再現可能にする

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:15-19`
  - `docs/er.md:71-78`
  - `docs/diagrams/er/combined-data-model-er.puml:207-229`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:380-402`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TOA-005`
- `actor`: 利用者、障害調査者
- `trigger`: 同じ完了済みジョブから成果物を再出力する。
- `audit event`: 再出力開始、参照した field summary、生成済み成果物との差分有無、終了状態。
- `saved summary`: job id、artifact id、source field digest、row digest、generated_at、status。
- `redaction rule`: digest には本文復元に使える値を含めない。
- `expected outcome`: 同じ入力 summary から同じ成果物を再現できるか確認できる。
- `observable point`: `TRANSLATION_ARTIFACT.generated_at`、row digest、artifact status、result summary。
- `related detail requirement type`: workflow / audit
- `adoption hint`: 再出力の正常系または回復系シナリオへ統合候補にする。
- `conflict hint`: ER は `TRANSLATION_ARTIFACT.translation_job_id` を `UNIQUE` としているため、履歴を複数世代で残すか、最新状態だけにするかは人間判断候補にする。

### CAND-TOA-006 成果物生成失敗を成功扱いにしない監査要約を残す

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:15-19`
  - `tasks/usecases/translation-output-artifact.yaml:23-26`
  - `docs/er.md:71-78`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:304-328`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TOA-006`
- `actor`: 利用者、運用確認者
- `trigger`: 成果物生成、XML 構成、保存、再出力のいずれかが失敗する。
- `audit event`: output artifact 生成失敗、error kind、対象行数、成功公開の拒否。
- `saved summary`: artifact id、attempted row count、failed stage、error kind、retryable flag、status。
- `redaction rule`: ファイル内容、外部 provider 応答原文、secret は error summary と log に含めない。
- `expected outcome`: 失敗した成果物が成功状態として表示または再利用されない。
- `observable point`: artifact status、error summary、result summary、structured log。
- `related detail requirement type`: failure / audit
- `adoption hint`: output artifact 生成失敗シナリオへ統合候補にする。
- `conflict hint`: failed artifact を保存するか、保存せず失敗要約だけを残すかは人間判断候補にする。

### CAND-TOA-007 入力ファイルの出自情報と成果物を結び付ける

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:10-19`
  - `docs/spec.md:11-15`
  - `docs/spec.md:205-207`
  - `docs/diagrams/er/combined-data-model-er.puml:26-35`
  - `docs/diagrams/er/combined-data-model-er.puml:144-154`
  - `docs/diagrams/er/combined-data-model-er.puml:207-216`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TOA-007`
- `actor`: 利用者、障害調査者
- `trigger`: 出力成果物の出自を確認する。
- `audit event`: 入力データ、翻訳ジョブ、出力成果物の関連確認。
- `saved summary`: source tool、target plugin、record count、job id、artifact id、artifact format、target game。
- `redaction rule`: source file path の保存粒度は伏せ字候補にし、全文 path を監査ログへ無条件に複製しない。
- `expected outcome`: 成果物がどの xEdit 抽出データと翻訳ジョブから作られたか後追いできる。
- `observable point`: `X_EDIT_EXTRACTED_DATA`、`TRANSLATION_JOB`、`TRANSLATION_ARTIFACT`。
- `related detail requirement type`: audit / persistence
- `adoption hint`: 出力成果物 detail または result summary の監査情報へ統合候補にする。
- `conflict hint`: 入力ファイル出自の可観測性と path redaction の範囲は人間判断候補にする。

### CAND-TOA-008 出力監査ログに保存禁止情報が出ないことを確認する

- `source requirement`:
  - `docs/spec.md:57-58`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:107-113`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:404-426`
  - `docs/architecture.md:92-98`
  - `docs/architecture.md:173-180`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TOA-008`
- `actor`: 利用者、運用確認者
- `trigger`: 成果物出力の success、failure、re-output を実行する。
- `audit event`: structured log、debug log、runtime event summary、UI error summary の redaction 確認。
- `saved summary`: artifact id、operation kind、row count、digest、error kind、credential reference presence。
- `redaction rule`: secret、API key 平文、復号可能値、外部 provider 応答原文、過剰な本文全文を保存しない。
- `expected outcome`: 障害調査に必要な要約は確認できるが、保存禁止情報は表示やログへ出ない。
- `observable point`: UI error summary、structured log、debug log、runtime event payload。
- `related detail requirement type`: security / audit
- `adoption hint`: secret 非露出の横断シナリオへ統合候補にする。
- `conflict hint`: output task 自体は provider を直接呼ばない可能性があるため、credential reference の扱いは body phase 由来情報の表示範囲として統合判断する。

### CAND-TOA-009 出力不可理由を監査できる形で返す

- `source requirement`:
  - `docs/spec.md:135-199`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:17-33`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:380-402`
  - `tasks/usecases/translation-output-artifact.yaml:20-27`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TOA-009`
- `actor`: 利用者、運用確認者
- `trigger`: 未完了、失敗、取り消し、本文結果不整合の job から成果物出力を試みる。
- `audit event`: 出力開始拒否、拒否理由、参照した job state、field summary。
- `saved summary`: job state、body phase state、ready flag、rejection reason、output status summary。
- `redaction rule`: 拒否理由には本文全文や secret を含めず、状態、件数、error kind で説明する。
- `expected outcome`: 出力不可理由を後から確認でき、途中成功結果が成果物出力に使われない。
- `observable point`: readiness result、rejection reason、Job Run summary、Output Review。
- `related detail requirement type`: workflow / audit
- `adoption hint`: output readiness の禁止遷移または失敗系シナリオへ統合候補にする。
- `conflict hint`: 出力不可の状態条件は state-transition 観点、failure 観点の候補と統合判断が必要になる。

## Open Notes

- `human decision candidate`:
  - `TRANSLATION_ARTIFACT.translation_job_id` が `UNIQUE` であるため、再出力履歴を複数世代で残すか、最新状態と要約だけ残すか決める必要がある。
  - source file path、`Source`、`Dest`、diff preview の伏せ字範囲を決める必要がある。
  - 成果物生成失敗時に failed artifact record を残すか、失敗要約だけ残すか決める必要がある。
  - structured log、debug log、runtime event summary の保持期間と保持粒度を決める必要がある。
- `merge candidate`:
  - CAND-TOA-001 は actor-goal の Output Review 確認シナリオと統合候補にする。
  - CAND-TOA-004 と CAND-TOA-005 は lifecycle の再出力シナリオと統合候補にする。
  - CAND-TOA-006 と CAND-TOA-009 は failure / state-transition の拒否理由シナリオと統合候補にする。
  - CAND-TOA-008 は external-integration または trust-boundary 系の secret 非露出シナリオと統合候補にする。
- `rejection candidate`:
  - row-level の `Source` と `Dest` 全文を監査ログへ複製する候補は保存禁止として却下候補にする。
  - 採否はこの成果物で確定しない。
