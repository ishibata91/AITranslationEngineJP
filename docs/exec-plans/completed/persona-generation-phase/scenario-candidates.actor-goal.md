# Scenario Candidates: persona-generation-phase / actor-goal

- `generator`: `actor-goal`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `PGP`
- `artifact_path`: `docs/exec-plans/active/persona-generation-phase/scenario-candidates.actor-goal.md`
- `candidate_count`: 6

## Generator Scope

- `viewpoint`: actor-goal
- `included_sources`: `./plan.md`, `tasks/index.yaml`, `tasks/usecases/persona-generation-phase.yaml`, `tasks/usecases/term-translation-phase.yaml`, `tasks/usecases/body-translation-phase.yaml`, `docs/spec.md`, `docs/er.md`, `docs/architecture.md`, `docs/exec-plans/completed/term-translation-phase/scenario-design.md`
- `excluded_sources`: product code, product tests, docs 正本化, 他 generator 成果物
- `generation_notes`: アクターの目的、開始操作、成功体験だけを候補化する。採否、統合、最終シナリオ表は `designer` に残す。

## Handoff Evidence

- `task artifact location`: `docs/exec-plans/active/persona-generation-phase/`
- `target delta`: 単語翻訳フェーズの後に、NPC の原文発話、属性メタデータ、会話文脈、共通ペルソナからジョブ内ペルソナを生成し、本文翻訳フェーズの入力として参照できるようにする。
- `viewpoint`: `actor-goal`
- `source boundary`: `task 枠` は completed。`scenario_candidates` は in_progress。
- `phase order evidence`: `tasks/index.yaml:4-10`, `tasks/usecases/term-translation-phase.yaml:21-27`, `tasks/usecases/persona-generation-phase.yaml:19-27`, `tasks/usecases/body-translation-phase.yaml:21-30`

## Candidate Scenarios

### CAND-PGP-001 Job Run から NPC ペルソナ生成フェーズを開始する

- `source requirement`: `tasks/usecases/term-translation-phase.yaml:21-27`, `tasks/usecases/persona-generation-phase.yaml:19-24`, `docs/spec.md:100-114`, `docs/spec.md:224-227`
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-PGP-001`
- `actor`: 翻訳ジョブを進めるユーザー
- `goal`: 単語翻訳フェーズ完了後に、本文翻訳前の NPC ペルソナ生成を開始する。
- `trigger`: Job Run を開き、NPC ペルソナ生成フェーズ開始を実行する。
- `expected outcome`: NPC ペルソナ生成フェーズが current phase になり、progress と開始結果を確認できる。
- `observable point`: Job Run UI、phase start result、`JOB_PHASE_RUN`
- `source requirement summary`: 単語翻訳フェーズ後、本文翻訳フェーズ前に NPC ペルソナ生成フェーズを実行する。
- `acceptance relevance`: `success_requirement`, `state_requirement`, `observability_requirement`
- `related detail requirement type`: `success_requirement`, `state_requirement`, `testability_requirement`
- `adoption hint`: 正常系の phase start 候補。開始可能状態の詳細は state-transition 側で固定する。
- `conflict hint`: Ready / Running / terminal job、既存 active phase run の扱いは他観点との統合対象にする。

### CAND-PGP-002 NPC ごとの生成対象を確認する

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:10-15`, `tasks/usecases/persona-generation-phase.yaml:21-24`, `docs/spec.md:21-25`, `docs/er.md:32-33`
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-PGP-002`
- `actor`: 生成対象を確認するユーザー
- `goal`: NPC 発話、NPC 属性メタデータ、会話文脈、共通ペルソナから、どの NPC が生成対象になるか判断する。
- `trigger`: Job Run の NPC ペルソナ生成フェーズで対象一覧または phase summary を開く。
- `expected outcome`: NPC ごとの生成対象、生成に使う入力種類、対象件数を確認できる。
- `observable point`: Job Run UI、target summary、`NPC_PROFILE` / `NPC_RECORD` 参照
- `source requirement summary`: NPC の原文発話、属性、会話文脈、共通ペルソナが入力であり、NPC ごとの生成対象を確認できる必要がある。
- `acceptance relevance`: `success_requirement`, `data_requirement`, `observability_requirement`
- `related detail requirement type`: `success_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: phase 開始前または実行中の対象確認候補。
- `conflict hint`: 共通ペルソナが存在する NPC を生成対象に含めるか、参照済みとして除外するかは designer の統合判断に残す。

### CAND-PGP-003 NPC の原文発話と属性からジョブ内ペルソナを生成する

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:9-18`, `docs/spec.md:22-24`, `docs/spec.md:212-216`, `docs/er.md:50-55`
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-PGP-003`
- `actor`: NPC ペルソナ生成フェーズの内部処理
- `goal`: NPC ごとの話し方、性格、属性の要約をジョブ内ペルソナとして作る。
- `trigger`: NPC ペルソナ生成フェーズが対象 NPC の生成 unit を実行する。
- `expected outcome`: 対象 job に紐づくジョブ内ペルソナが生成され、生成根拠の翻訳フィールドを追跡できる。
- `observable point`: `PERSONA`, `PERSONA_FIELD_EVIDENCE`, phase result
- `source requirement summary`: NPC 発話、種族、性別、属性から AI でペルソナを生成し、ジョブ内生成物として保持する。
- `acceptance relevance`: `success_requirement`, `data_requirement`, `consistency_requirement`
- `related detail requirement type`: `success_requirement`, `data_requirement`, `testability_requirement`
- `adoption hint`: 主要なデータ生成成功候補。
- `conflict hint`: AI provider request、応答不正、保存失敗の扱いは external-integration / failure 側へ残す。

### CAND-PGP-004 共通ペルソナを参照しつつジョブ内ペルソナと分離する

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:15-18`, `tasks/usecases/persona-generation-phase.yaml:24-25`, `docs/spec.md:24-25`, `docs/spec.md:34-35`, `docs/spec.md:243-247`, `docs/er.md:50-55`
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-PGP-004`
- `actor`: 共通ペルソナを再利用したいユーザー
- `goal`: ジョブをまたぐ共通ペルソナと、翻訳ジョブ内で使うペルソナを混同せずに利用する。
- `trigger`: 共通ペルソナを持つ NPC を含むジョブで NPC ペルソナ生成フェーズを開始する。
- `expected outcome`: 共通ペルソナの参照状態と、ジョブ内ペルソナまたは persona snapshot の分離状態を確認できる。
- `observable point`: Job Run UI、`PERSONA.persona_scope`, `PERSONA.translation_job_id`, persona snapshot summary
- `source requirement summary`: 共通ペルソナは再利用可能であり、mod 翻訳中の NPC ペルソナは共通ペルソナと分離して保持する。
- `acceptance relevance`: `data_requirement`, `compatibility_requirement`, `consistency_requirement`
- `related detail requirement type`: `data_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: 共通データ再利用と job-scoped data の境界候補。
- `conflict hint`: ER は同一 NPC profile に共通とジョブ内を同時保持しない前提を持つため、生成、skip、snapshot 参照のどれを成功扱いにするかは人間判断候補にする。

### CAND-PGP-005 本文翻訳フェーズが参照する persona snapshot を用意する

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:16-25`, `tasks/usecases/body-translation-phase.yaml:9-16`, `tasks/usecases/body-translation-phase.yaml:21-30`, `docs/spec.md:102-115`, `docs/spec.md:130-131`
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-PGP-005`
- `actor`: 本文翻訳フェーズへ進めたいユーザー
- `goal`: 生成済みペルソナを本文翻訳フェーズの入力として参照できる状態にする。
- `trigger`: NPC ペルソナ生成フェーズ完了後に、phase result または本文翻訳フェーズ開始前の入力 summary を確認する。
- `expected outcome`: 本文翻訳フェーズ入力として参照できる persona snapshot と対象 NPC 数を確認できる。
- `observable point`: phase result、body phase input summary、`PHASE_RUN_PERSONA`
- `source requirement summary`: 本文翻訳フェーズはジョブ内ペルソナと翻訳補助メタデータを参照して本文を翻訳する。
- `acceptance relevance`: `success_requirement`, `state_requirement`, `compatibility_requirement`
- `related detail requirement type`: `success_requirement`, `state_requirement`, `consistency_requirement`
- `adoption hint`: 後続フェーズ連携の成功候補。
- `conflict hint`: snapshot の固定時点、欠落時の後続 phase block、再生成時の参照更新は state-transition / lifecycle 側と統合する。

### CAND-PGP-006 Job Run で生成結果と参照状態を確認する

- `source requirement`: `tasks/usecases/persona-generation-phase.yaml:21-34`, `docs/spec.md:35`, `docs/spec.md:53-55`, `docs/spec.md:133-134`, `docs/architecture.md:60-66`, `docs/architecture.md:173-180`
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-PGP-006`
- `actor`: phase result を確認するユーザー
- `goal`: NPC ペルソナ生成フェーズの結果、progress、persona snapshot 参照状態を UI で確認する。
- `trigger`: NPC ペルソナ生成フェーズの実行中または完了後に Job Run を開く。
- `expected outcome`: current phase、progress、phase result、生成済み件数、参照可能状態が表示される。
- `observable point`: Job Run UI、view model、gateway response
- `source requirement summary`: 翻訳補助メタデータ、辞書、共通基盤データは実行前後とも UI から観測可能であり、進捗も UI から観測する。
- `acceptance relevance`: `observability_requirement`, `success_requirement`, `testability_requirement`
- `related detail requirement type`: `success_requirement`, `observability_requirement`, `testability_requirement`
- `adoption hint`: UI 人間操作 E2E の表示確認候補。
- `conflict hint`: UI の表示項目、有効条件、詳細な画面状態差分は `ui-design` と designer 統合に残す。

## Open Notes

- `human decision candidate`: 共通ペルソナが存在する NPC で、ジョブ内ペルソナを新規生成するか、共通ペルソナを snapshot 参照するか、生成対象から除外するか。
- `human decision candidate`: AI provider のペルソナ生成結果を自動採用するか、人間確認を挟むか。
- `human decision candidate`: persona snapshot の固定時点を phase 開始時、phase 完了時、本文翻訳フェーズ開始時のどれにするか。
- `merge candidate`: `CAND-PGP-001` と `CAND-PGP-006` は Job Run 起点の UI 成功体験として統合される可能性がある。
- `merge candidate`: `CAND-PGP-003` と `CAND-PGP-005` は生成と後続参照のデータ整合シナリオとして統合される可能性がある。
- `rejection candidate`: なし。採否は `designer` の `scenario-design.candidate-coverage.json` に残す。
