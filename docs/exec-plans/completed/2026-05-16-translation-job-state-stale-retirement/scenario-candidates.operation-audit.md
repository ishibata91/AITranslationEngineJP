# Scenario Candidates: 2026-05-16-translation-job-state-stale-retirement / operation-audit

- `generator`: `operation-audit`
- `source_plan`: `./implement-lane-task-frame.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJSSR-OA`

## Generator Scope

- `viewpoint`: operation-audit
- `included_sources`: `./implement-lane-task-frame.md`, `./plan.md`, `./state-knowledge-investigation.md`, `../observability-log-addition/scenario-design.md`, `../observability-log-addition/scenario-candidates.operation-audit.md`, `docs/observability-logging.md`, `docs/spec.md`, `docs/architecture.md`, `docs/er.md`, `docs/detail-specs/README.md`, `docs/scenario-tests/README.md`, `docs/screen-design/README.md`
- `excluded_sources`: product code 変更、product test 変更、docs 正本本文変更、`docs/exec-plans/completed/**`、UI 変更、観測ログ追加の実装判断、他 agent の候補生成、最終シナリオ採否
- `generation_notes`: 候補は、旧名参照、状態語彙、残留参照の扱いを後から確認できるかに限定する。`stale_selection`、`validation_stale`、`model_selection_stale` は削除対象にしない。

## Candidate Scenarios

### CAND-TJSSR-OA-001 active observability task-local の旧名参照を今回更新した事実を監査する

- `source requirement`: `./implement-lane-task-frame.md` は、`observability-log-addition` の旧名参照を今回の active task-local 更新に含めるかを最初に固定する判断としている。`./state-knowledge-investigation.md` は、`StateMachine` 旧名が active observability task-local に残ると古い責務境界を再注入する可能性があるとしている。
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TJSSR-OA-001`
- `actor`: 設計レビュー担当者、実装引き継ぎ担当者
- `trigger`: designer が、`observability-log-addition` の active 成果物を今回の stale 廃止で更新対象に含める判断を候補にする。
- `expected outcome`: 更新対象に含めた場合、active observability task-local の `StateMachine` 参照が `TranslationJobPolicy` 相当の説明へ置き換わったことを後から確認できる。更新した範囲と、更新していない completed archive の範囲を区別できる。
- `observable point`: `docs/exec-plans/active/observability-log-addition/scenario-design.md`、`scenario-candidates.*.md`、`design-diff.*.puml` に残る旧名参照の検索結果
- `saved summary`: 更新対象 path、旧名、置換後の責務名、更新しない path 分類、確認 command の要約
- `redaction rule`: docs path と固定名だけを残す。provider 応答、prompt、翻訳本文、secret、API key、endpoint 実値は残さない。
- `forbidden saved information`: product runtime payload、provider raw payload、全文本文、credential 参照実値
- `related detail requirement type`: `observability_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: 今回の task-local 更新をシナリオ設計に含める場合の候補にする。
- `conflict hint`: `docs/exec-plans/completed/**` まで更新する案と衝突する。観測ログ追加の実装範囲を広げる案とも衝突する。
- `open issue`: `JobIOService` 参照を同時に更新するかは、人間判断または designer の統合判断に残る。

### CAND-TJSSR-OA-002 active observability task-local の旧名参照を別 task に送った事実を監査する

- `source requirement`: `./implement-lane-task-frame.md` は、`observability-log-addition` の旧名参照を今回更新する場合と別 task に送る場合を判断対象にしている。`./state-knowledge-investigation.md` は、active observability task-local の旧名参照を影響ファイル候補に挙げている。
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TJSSR-OA-002`
- `actor`: 設計レビュー担当者、後続 task 起票者
- `trigger`: designer が、`observability-log-addition` の active 成果物更新を今回の stale 廃止から外す判断を候補にする。
- `expected outcome`: 別 task に送った場合、残る `StateMachine` / `JobIOService` 参照が既知の残留参照として記録される。後続 task の入口、残留理由、誤って product code の更新対象にしない条件を確認できる。
- `observable point`: 今回の task-local 成果物の未決事項、後続 task 送り理由、旧名参照検索結果
- `saved summary`: 残留 path、残留固定名、残留理由、戻し先、後続 task で再確認する条件
- `redaction rule`: 残留参照の path と固定名だけを残す。runtime log payload や利用者データは残さない。
- `forbidden saved information`: frontend console log 本文、backend JSON log 本文、provider raw payload、翻訳本文全文
- `related detail requirement type`: `observability_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: active observability task-local を今回変更しない場合の候補にする。
- `conflict hint`: 残留参照を未記録のまま close する案は、後続 task で旧責務境界を再利用するリスクと衝突する。
- `open issue`: 別 task の owner と起動条件は、この候補では確定しない。

### CAND-TJSSR-OA-003 `JobIOService` の扱い別に監査対象の境界を分ける

- `source requirement`: `docs/architecture.md` は `JobIOService` を job と phase run の状態取得と保存だけを扱う構造主語としている。`./state-knowledge-investigation.md` は、`JobIOService` が architecture 正本と lint component に残る一方で実体 package は `doc.go` だけだとしている。
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TJSSR-OA-003`
- `actor`: 設計レビュー担当者、実装範囲作成者
- `trigger`: designer が、`JobIOService` を今回 architecture 正本から外す候補、または別 task で実体化する候補を比較する。
- `expected outcome`: `JobIOService` を外す場合、状態事実の取得と保存の監査対象が usecase、service、repository のどの境界へ移るかを確認できる。別 task で実体化する場合、active observability task-local の `JobIOService` 参照が残留または補足対象であることを確認できる。
- `observable point`: `docs/architecture.md` の `JobIOService` 節、active observability task-local の `JobIOService` 参照、今回 task-local の人間判断候補
- `saved summary`: 判断分岐、対象境界、更新対象 path、残留理由、正本化判断の必要有無
- `redaction rule`: 状態名、境界名、path だけを残す。operation summary、policy 判定履歴、provider raw payload、secret は残さない。
- `forbidden saved information`: policy rule 名の永続化、policy 判定履歴、operation summary、credential 参照実値
- `related detail requirement type`: `observability_requirement`, `data_requirement`, `consistency_requirement`, `security_requirement`
- `adoption hint`: `JobIOService` の扱いが今回の設計判断に残る場合に候補にする。
- `conflict hint`: `JobIOService` 廃止をこの agent が確定する案は、docs 正本化判断と衝突する。`JobIOService` を残したまま実体 package の欠落を監査対象から外す案とも衝突する。
- `open issue`: `JobIOService` を廃止するか実体化するかは人間判断に残る。

### CAND-TJSSR-OA-004 `pending` と正本 state の差分を後から追跡できるようにする

- `source requirement`: `docs/spec.md` は `TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` に `pending` を含めていない。`./state-knowledge-investigation.md` は、3 phase service に `pending` が内部 state として残ると観測している。
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TJSSR-OA-004`
- `actor`: 設計レビュー担当者、状態不整合調査者
- `trigger`: designer が、`pending` を canonical state に昇格するか、内部一時 state として隔離するかを設計入力へ入れる。
- `expected outcome`: `pending` の扱いを決めた後、正本 state との差分、DB 永続 state として作られる経路、read model 用の派生 state との違いを後から確認できる。
- `observable point`: `docs/spec.md` の state 定義、`./state-knowledge-investigation.md` の repository write path 未追跡事項、phase service の state 派生知識
- `saved summary`: state 名、正本 state との差分、内部一時 state かどうか、未追跡 write path、確認対象境界
- `redaction rule`: state 名と境界名だけを残す。利用者入力、翻訳本文、provider 応答原文は残さない。
- `forbidden saved information`: phase 入力本文、provider raw response、prompt 全文、DTO 全体
- `related detail requirement type`: `observability_requirement`, `state_requirement`, `consistency_requirement`, `data_requirement`
- `adoption hint`: `pending` の扱いを今回の stale 廃止で設計判断に含める場合に候補にする。
- `conflict hint`: `pending` を正本 state として扱う案は、現行 `docs/spec.md` と衝突する可能性がある。read model 用の表示 state と永続 state を同じ監査対象に混ぜる案とも衝突する。
- `open issue`: `pending` が DB 永続 state として実際にどの経路で作られるかは追加追跡が必要である。

### CAND-TJSSR-OA-005 ドメイン仕様の stale reason が削除対象に混ざっていないことを監査する

- `source requirement`: `./implement-lane-task-frame.md` と `./plan.md` は、`stale_selection`、`validation_stale`、`model_selection_stale` を削除対象にしないと明記している。`docs/observability-logging.md` は stale event の破棄理由を観測対象にしている。
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TJSSR-OA-005`
- `actor`: 設計レビュー担当者、回帰確認者
- `trigger`: stale 廃止の設計差分または実装範囲が、`stale` という語を含む候補を削除対象へ入れる。
- `expected outcome`: 旧設計名としての stale と、利用者向けまたは API 向け理由分類としての stale reason を区別して確認できる。ドメイン仕様の stale reason が保持対象として残る。
- `observable point`: `./implement-lane-task-frame.md` の禁止範囲、`./state-knowledge-investigation.md` の stale 語彙分離、active observability task-local の stale reason 候補
- `saved summary`: stale 語彙の分類、削除対象か保持対象か、保持理由、確認 command の要約
- `redaction rule`: reason 名と path だけを残す。runtime event payload 全体、画面状態全体、利用者入力は残さない。
- `forbidden saved information`: runtime event payload 全体、frontend store 全体、翻訳本文全文、credential 参照実値
- `related detail requirement type`: `observability_requirement`, `compatibility_requirement`, `security_requirement`
- `adoption hint`: stale 廃止の回帰確認をシナリオ設計へ入れる場合に候補にする。
- `conflict hint`: `stale_selection`、`validation_stale`、`model_selection_stale` を削除対象へ含める案と衝突する。
- `open issue`: runtime event の stale 判定をどの task で再確認するかは、この候補では確定しない。

### CAND-TJSSR-OA-006 状態 spelling 差分の検索漏れを監査する

- `source requirement`: `docs/spec.md` は terminal state を `Canceled` としている。`./state-knowledge-investigation.md` は、`PersonaGenerationPhaseContractStub` に `cancelled` fixture spelling が残ると観測している。
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-TJSSR-OA-006`
- `actor`: 設計レビュー担当者、回帰確認者
- `trigger`: designer が、`cancelled` fixture spelling を今回の stale 廃止に含めるか別 task に送るかを設計入力へ入れる。
- `expected outcome`: `canceled` と `cancelled` の差分が、仕様 state の違いなのか fixture spelling の違いなのかを後から確認できる。検索条件に片方だけを使って state 知識を見落とすリスクを分離できる。
- `observable point`: `docs/spec.md` の `Canceled` 定義、`./state-knowledge-investigation.md` の spelling 観測、fixture spelling の検索結果
- `saved summary`: 正本 spelling、差分 spelling、差分が残る path、今回対象か別 task 送りか
- `redaction rule`: state spelling と path だけを残す。翻訳本文、provider payload、利用者データは残さない。
- `forbidden saved information`: 入力 file 内容、翻訳対象 field、provider raw payload、secret
- `related detail requirement type`: `observability_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: spelling 差分を今回の stale 廃止で扱うかを検討する場合に候補にする。
- `conflict hint`: spelling 差分を state 仕様変更として扱う案は、現行 `docs/spec.md` と衝突する可能性がある。
- `open issue`: fixture spelling 修正を今回の実装範囲へ含めるかは、この候補では確定しない。

## Open Notes

- `human decision candidate`: `observability-log-addition` の active 成果物を今回更新するか、別 task に送るかは人間判断または designer の統合判断に残る。
- `human decision candidate`: `JobIOService` を architecture 正本から外すか、別 task で実体化するかは人間判断に残る。
- `human decision candidate`: `pending` を canonical state へ昇格するか、内部一時 state として隔離するかは人間判断に残る。
- `merge candidate`: `StateMachine` 旧名参照の更新候補と `JobIOService` 分岐候補は、active observability task-local 更新シナリオへ統合候補になる。
- `merge candidate`: stale reason 保持候補と frontend runtime event stale 候補は、観測ログ task の既存候補と統合候補になる。
- `rejection candidate`: product code、product test、docs 正本本文、completed archive、UI 変更、観測ログ追加実装を前提にする候補は不採用候補である。
- `conflict candidate`: 保存対象が secret、API key、credential 参照実値、provider raw payload、prompt 全文、翻訳本文全文を含む場合は `security_requirement` / `data_requirement` と衝突する。
- `candidate count`: 6
