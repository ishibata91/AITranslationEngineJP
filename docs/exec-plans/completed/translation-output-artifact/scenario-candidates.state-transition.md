# Scenario Candidates: translation-output-artifact / state-transition

- `generator`: `state-transition`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TOA-ST`
- `candidate_count`: 8

## Generator Scope

- `viewpoint`: `state-transition`
- `included_sources`:
  - `docs/exec-plans/active/translation-output-artifact/plan.md:5-17`
  - `tasks/usecases/translation-output-artifact.yaml:1-35`
  - `docs/spec.md:29-32`
  - `docs/spec.md:43-67`
  - `docs/spec.md:135-199`
  - `docs/spec.md:222-235`
  - `docs/er.md:20-23`
  - `docs/er.md:61-77`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:20-32`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:88-105`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:116-119`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:380-400`
- `excluded_sources`:
  - プロダクトコード、プロダクトテスト、docs 正本、`.codex` は変更対象にしない。
  - 採否、統合、最終シナリオ表の確定は `designer` に残す。
  - 画面固有の表示契約は、後続の設計成果物で確定する。
- `generation_notes`:
  - 出力成果物の状態名は候補ラベルであり、正式な状態値ではない。
  - `TRANSLATION_ARTIFACT` と `XTRANSLATOR_OUTPUT_ROW` は、出力状態の観測点として扱う。
  - 再出力と差分プレビューの状態遷移は、usecase 上の要求から候補化する。

## Candidate Scenarios

### CAND-TOA-ST-001 完了済みジョブだけが出力準備済みに遷移する

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:9-26`
  - `docs/spec.md:146-150`
  - `docs/spec.md:187-199`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:386-400`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TOA-ST-001`
- `actor`: Output Review を操作するユーザー。
- `trigger`: 完了済み job を選び、出力 readiness を確認する。
- `transition before state`: 本文翻訳フェーズが Completed で、job-level `Completed` が成立している。
- `start condition`: 訳文、出力ステータス、保護要素検証結果が field 単位で整合している。
- `transition after state`: 出力準備済み候補状態になり、出力開始が可能になる。
- `expected outcome`: result summary と output readiness が確認できる。
- `observable point`: readiness result、`JOB_TRANSLATION_FIELD`、phase result、output status summary。
- `related detail requirement type`: state gate、display、persistence。
- `adoption hint`: `designer` は、上流の readiness boundary と同じ開始条件に統合できる。
- `conflict hint`: job-level `Completed` と body phase Completed の反映を分ける設計は開始条件と衝突する。

### CAND-TOA-ST-002 未完了または不整合の job は出力成果物を作らない

- `source requirement`:
  - `docs/spec.md:162-183`
  - `docs/spec.md:187-199`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:31-32`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:397-400`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TOA-ST-002`
- `actor`: Output Review を操作するユーザー。
- `trigger`: 未完了、失敗中、取り消し済み、不整合ありの job で出力開始を試す。
- `transition before state`: job が `Draft`、`Ready`、`Running`、`Paused`、`RecoverableFailed`、`Failed`、`Canceled` のいずれか、または field result が不整合である。
- `start condition`: 出力開始条件を満たさない。
- `transition after state`: job と field result は変化せず、出力成果物は未生成のままになる。
- `expected outcome`: 出力開始は拒否され、拒否理由が確認できる。
- `observable point`: rejection reason、`TRANSLATION_ARTIFACT` 未作成、`XTRANSLATOR_OUTPUT_ROW` 未作成。
- `related detail requirement type`: forbidden transition、state invariant、display。
- `adoption hint`: `designer` は、禁止遷移として readiness false と開始拒否を統合できる。
- `conflict hint`: 部分成功結果を成果物出力へ使う設計は、`Canceled` 後の途中成功結果を使わない上流判断と衝突する。

### CAND-TOA-ST-003 出力準備済み job から xTranslator 互換成果物を生成する

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:15-19`
  - `docs/spec.md:14`
  - `docs/spec.md:43`
  - `docs/spec.md:65-67`
  - `docs/spec.md:222-235`
  - `docs/er.md:20-23`
  - `docs/er.md:71-77`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TOA-ST-003`
- `actor`: 出力成果物生成処理。
- `trigger`: 出力準備済み job で xTranslator 互換成果物の出力を実行する。
- `transition before state`: 出力準備済み job に、まだ現行の出力成果物がない。
- `start condition`: 入力ファイルの出自情報、訳文、出力ステータス、xTranslator 行に必要な識別情報を参照できる。
- `transition after state`: `TRANSLATION_ARTIFACT` が生成済み候補状態になり、`XTRANSLATOR_OUTPUT_ROW` が作成される。
- `expected outcome`: xTranslator 互換 XML の出力行を再構成できる。
- `observable point`: `TRANSLATION_ARTIFACT`、`XTRANSLATOR_OUTPUT_ROW`、EDID、REC、FIELD、FORMID、Source、Dest、Status。
- `related detail requirement type`: output format、persistence、state transition。
- `adoption hint`: `designer` は、成果物生成の成功状態と行生成の整合確認を同じ候補へ統合できる。
- `conflict hint`: 出力行の生成規則を詳細化しすぎると、failure または external-integration 観点の候補と重複する。

### CAND-TOA-ST-004 内部 cached は xTranslator 出力時に Status=1 へ遷移する

- `source requirement`:
  - `docs/spec.md:29-32`
  - `docs/spec.md:65-67`
  - `docs/er.md:76-77`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TOA-ST-004`
- `actor`: 出力行生成処理。
- `trigger`: 内部出力ステータス `cached` を持つ翻訳 field を出力行へ変換する。
- `transition before state`: `JOB_TRANSLATION_FIELD` が内部ステータス `cached` を保持している。
- `start condition`: xTranslator 互換形式への出力が開始されている。
- `transition after state`: `XTRANSLATOR_OUTPUT_ROW.Status` は `1` として観測できる。
- `expected outcome`: 辞書置換である内部観測情報は、xTranslator の `Status` とは別に保持される。
- `observable point`: `JOB_TRANSLATION_FIELD` の内部出力ステータス、`XTRANSLATOR_OUTPUT_ROW.Status`、辞書置換の内部観測情報。
- `related detail requirement type`: status mapping、output format、persistence。
- `adoption hint`: `designer` は、出力ステータス語彙と xTranslator `Status` の写像確認として扱える。
- `conflict hint`: 内部 `cached` を xTranslator `Status` と同一概念として扱う設計は、内部観測情報の保持要件と衝突する。

### CAND-TOA-ST-005 本文翻訳対象 0 件の完了 job も出力可能状態になる

- `source requirement`:
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:29-30`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:116-119`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:397-400`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TOA-ST-005`
- `actor`: Output Review を操作するユーザー。
- `trigger`: 本文翻訳対象 0 件で正常完了した job の出力 readiness を確認する。
- `transition before state`: 本文翻訳対象が 0 件で、body phase と job-level `Completed` が成立している。
- `start condition`: provider 未実行、target count、skipped count、output readiness への影響を確認できる。
- `transition after state`: 出力準備済み候補状態になり、単語翻訳結果を成果物出力へ渡せる。
- `expected outcome`: target count 0 を理由に出力不可へ遷移しない。
- `observable point`: target count、skipped count、readiness result、result summary。
- `related detail requirement type`: edge state、state gate、display。
- `adoption hint`: `designer` は、空本文 job を正常系の境界値として採否判断できる。
- `conflict hint`: 出力行が 0 件になる扱いと、単語翻訳結果だけを出す扱いは人間判断が必要である。

### CAND-TOA-ST-006 翻訳 unit 再生成後は成果物が再出力必要状態になる

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:17-19`
  - `tasks/usecases/translation-output-artifact.yaml:23-26`
  - `docs/spec.md:43`
  - `docs/er.md:73-77`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TOA-ST-006`
- `actor`: Output Review を操作するユーザー。
- `trigger`: 既存の出力成果物がある job で、translation unit 単位の再生成結果を確認する。
- `transition before state`: 生成済み出力成果物があり、元になった `JOB_TRANSLATION_FIELD` の内容と出力行が対応している。
- `start condition`: translation unit の訳文または出力ステータスが、生成済み成果物の元データから変化している。
- `transition after state`: 出力成果物は再出力必要候補状態になり、diff preview と再出力導線が表示される。
- `expected outcome`: 古い成果物を現行成果物として扱わない。
- `observable point`: diff preview、再出力状態、`JOB_TRANSLATION_FIELD`、`XTRANSLATOR_OUTPUT_ROW`。
- `related detail requirement type`: stale state、diff state、re-output workflow。
- `adoption hint`: `designer` は、再生成操作が output artifact task の範囲かどうかを判断して統合できる。
- `conflict hint`: translation unit 再生成を別 task に逃がす設計では、この候補は表示専用の stale 観測に縮小される。

### CAND-TOA-ST-007 同じ入力の再出力で重複成果物または重複行を作らない

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:19`
  - `tasks/usecases/translation-output-artifact.yaml:25-26`
  - `docs/er.md:66-77`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TOA-ST-007`
- `actor`: 出力成果物生成処理。
- `trigger`: 同じ Completed job と同じ field result で再出力を実行する。
- `transition before state`: 生成済み出力成果物が存在し、元データに変化がない。
- `start condition`: 同じ入力、同じ出力形式、同じ field 集合で再出力要求が来る。
- `transition after state`: 現行成果物は 1 つの生成済み候補状態に収束し、同一 field の出力行は重複しない。
- `expected outcome`: result summary と出力行件数が再実行前後で破損しない。
- `observable point`: `TRANSLATION_ARTIFACT`、`XTRANSLATOR_OUTPUT_ROW`、row count、result summary。
- `related detail requirement type`: idempotent transition、persistence、re-output state。
- `adoption hint`: `designer` は、同一成果物を置換するか新しい revision を作るかを人間判断へ回せる。
- `conflict hint`: artifact revision 履歴を残す設計では、現行成果物 1 つへの収束条件と履歴保存条件を分ける必要がある。

### CAND-TOA-ST-008 出力生成失敗を生成済み成果物として扱わない

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:16-19`
  - `tasks/usecases/translation-output-artifact.yaml:23-26`
  - `docs/spec.md:53-67`
  - `docs/er.md:71-77`
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TOA-ST-008`
- `actor`: 出力成果物生成処理。
- `trigger`: 出力生成中に成果物保存または行生成が完了できない。
- `transition before state`: 出力準備済み job から生成中候補状態へ入っている。
- `start condition`: xTranslator 互換成果物または出力行の生成完了条件を満たさない。
- `transition after state`: 生成済み候補状態へは遷移せず、失敗または再試行可能候補状態として観測できる。
- `expected outcome`: 不完全な成果物を completed output artifact として表示しない。
- `observable point`: output artifact state、失敗理由、再出力導線、`TRANSLATION_ARTIFACT`、`XTRANSLATOR_OUTPUT_ROW`。
- `related detail requirement type`: forbidden transition、failure state、re-output workflow。
- `adoption hint`: `designer` は、failure 観点候補と統合する時に状態不変条件だけを残せる。
- `conflict hint`: failure 観点が詳細な保存失敗パターンを持つ場合、この候補は状態遷移の禁止条件へ縮小する。

## Open Notes

- `human decision candidate`:
  - 出力成果物の正式な状態値を決める必要がある。
  - 再出力時に既存成果物を置換するか、revision 履歴を持つかを決める必要がある。
  - translation unit 再生成を output artifact task の操作範囲に含めるかを決める必要がある。
  - stale 判定の根拠を field digest、artifact revision、updated_at のどれにするかを決める必要がある。
  - 出力生成失敗後の retry 条件と表示状態を決める必要がある。
- `merge candidate`:
  - CAND-TOA-ST-002 は failure 観点の拒否理由候補と統合できる。
  - CAND-TOA-ST-003 と CAND-TOA-ST-004 は output row 生成規則の設計で統合できる。
  - CAND-TOA-ST-006 と CAND-TOA-ST-007 は lifecycle または operation-audit 観点の成果物履歴候補と統合できる。
- `rejection candidate`:
  - 状態変化を伴わない UI 操作列だけの候補は、この観点では採用しない候補に回す。
  - xTranslator XML の細かい serialization 規則だけの候補は、この観点では採用しない候補に回す。
  - provider 実行や本文翻訳 retry の候補は、上流 task または別観点へ回す。
