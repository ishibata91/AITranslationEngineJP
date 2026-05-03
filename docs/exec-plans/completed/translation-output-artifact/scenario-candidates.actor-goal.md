# Scenario Candidates: translation-output-artifact / actor-goal

- `generator`: `actor-goal`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TOA`

## Generator Scope

- `viewpoint`: actor-goal
- `task_artifact_root`: `docs/exec-plans/active/translation-output-artifact/`
- `target_delta`: 完了ジョブの翻訳結果を確認し、xTranslator 互換の成果物として出力する。
- `candidate_artifact`: `./scenario-candidates.actor-goal.md`
- `generated_candidate_count`: 8
- `included_sources`:
  - `./plan.md`
  - `tasks/index.yaml`
  - `tasks/usecases/translation-output-artifact.yaml`
  - `docs/spec.md`
  - `docs/er.md`
  - `docs/architecture.md`
  - `docs/screen-design/README.md`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md`
- `excluded_sources`:
  - product code
  - product test
  - docs canonicalization
  - final scenario adoption decision
  - other scenario candidate generator output
- `generation_notes`: 実行者の目的、開始経路、成功結果、観測点だけを候補化する。状態遷移網羅、失敗詳細、外部ファイル I/O 詳細、監査保存詳細は他観点または designer の統合対象へ残す。

## Candidate Scenarios

### CAND-TOA-001 Output Review で completed job を選び出力対象を確認する

- `source requirement`: `tasks/usecases/translation-output-artifact.yaml:20-26` は本文翻訳フェーズ完了、completed job 一覧、result summary、diff preview、output artifact 状態を求める。`docs/exec-plans/completed/body-translation-phase/scenario-design.md:386-400` は body phase Completed と field result 整合時だけ output readiness が成立するとする。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TOA-001`
- `actor`: xTranslator 互換成果物を作りたいユーザー
- `trigger`: Output Review を開き、completed job を一覧から選ぶ。
- `expected outcome`: 出力可能な completed job が確認でき、選択した job の output readiness と拒否理由の有無を判断できる。
- `observable point`: Output Review UI、completed job list、readiness result、`JOB_TRANSLATION_FIELD`
- `related detail requirement type`: `success_requirement` / `state_requirement` / `display_requirement`
- `adoption hint`: designer は completed job の一覧表示と readiness 表示を正常系の入口候補として扱える。
- `conflict hint`: state-transition 観点の terminal job 除外、failure 観点の status 不整合拒否と統合が必要になる。

### CAND-TOA-002 result summary で出力前の翻訳結果を確認する

- `source requirement`: `tasks/usecases/translation-output-artifact.yaml:10-18` は訳文、出力ステータス、入力ファイルの出自情報、result summary を入出力に含める。`docs/spec.md:43` は最終出力に必要な識別子、原文、訳文、出力ステータスを lossless に保持することを求める。`docs/er.md:23-24` は翻訳結果と出力ステータスを `JOB_TRANSLATION_FIELD` に保持するとする。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TOA-002`
- `actor`: 翻訳結果を出力前に確認するユーザー
- `trigger`: Output Review で completed job を選択し、result summary を開く。
- `expected outcome`: 訳文件数、出力ステータス件数、入力ファイルの出自、出力可能性を summary として確認できる。
- `observable point`: result summary、input provenance summary、output status summary、`JOB_TRANSLATION_FIELD`
- `related detail requirement type`: `display_requirement` / `data_requirement` / `consistency_requirement`
- `adoption hint`: designer は result summary を diff preview と export 操作の前提表示候補として扱える。
- `conflict hint`: operation-audit 観点で保存する summary 項目、security 観点で表示してよい入力出自情報との統合が必要になる。

### CAND-TOA-003 translation unit 単位の diff preview を確認する

- `source requirement`: `tasks/usecases/translation-output-artifact.yaml:22-26` は translation unit 単位の差分と再生成操作、output artifact 状態と再出力導線を completion criteria にする。`docs/spec.md:65-67` は xTranslator 互換形式の出力と XML を標準配布形式にすることを求める。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TOA-003`
- `actor`: 出力前の差分を確認したいユーザー
- `trigger`: completed job と output artifact を選択し、translation unit の diff preview を開く。
- `expected outcome`: 1 translation unit ごとに原文、訳文、出力ステータス、出力対象 row への反映内容を確認できる。
- `observable point`: diff preview、selected translation unit、output row preview、output status summary
- `related detail requirement type`: `success_requirement` / `display_requirement` / `consistency_requirement`
- `adoption hint`: designer は diff preview を user-facing 成功確認の中心候補として扱える。
- `conflict hint`: 再生成操作が field result の再翻訳なのか、artifact row の再構成なのかは human decision candidate に残す。

### CAND-TOA-004 xTranslator 互換 XML を出力して成果物状態を確認する

- `source requirement`: `tasks/usecases/translation-output-artifact.yaml:15-19` は xTranslator 互換成果物、result summary、diff preview、再出力状態を outputs にする。`docs/spec.md:65-67` は xTranslator 互換形式と `.xml` を出力要件にする。`docs/er.md:71-77` は `TRANSLATION_ARTIFACT` と `XTRANSLATOR_OUTPUT_ROW` を出力モデルにする。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TOA-004`
- `actor`: 翻訳成果物を配布形式で出力するユーザー
- `trigger`: Output Review で出力可能な completed job を選び、xTranslator 互換成果物の出力を実行する。
- `expected outcome`: xTranslator 互換 XML の生成結果と成果物状態を確認できる。
- `observable point`: export result、`TRANSLATION_ARTIFACT`、`XTRANSLATOR_OUTPUT_ROW`、Output Review UI
- `related detail requirement type`: `success_requirement` / `data_requirement` / `compatibility_requirement`
- `adoption hint`: designer は実ファイル生成の詳細ではなく、ユーザーが生成結果を確認できる成功体験として扱える。
- `conflict hint`: external-integration 観点の XML adapter 境界、lifecycle 観点の artifact 状態更新と統合が必要になる。

### CAND-TOA-005 xTranslator row の必須列を再構成できることを確認する

- `source requirement`: `docs/spec.md:65-67` は `EDID`、`REC`、`FIELD`、`FORMID`、`Source`、`Dest`、`Status` の再構成を求める。`docs/er.md:76-77` は `XTRANSLATOR_OUTPUT_ROW` が 1 つの `JOB_TRANSLATION_FIELD` に対応し、同じ列を保持するとする。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TOA-005`
- `actor`: xTranslator 互換性を確認するユーザー
- `trigger`: Output Review で output row preview または export summary を確認する。
- `expected outcome`: 各出力行に必要な xTranslator 互換列が揃い、訳文と出力ステータスが対象 field に対応していることを確認できる。
- `observable point`: output row preview、`XTRANSLATOR_OUTPUT_ROW`、`JOB_TRANSLATION_FIELD`、result summary
- `related detail requirement type`: `compatibility_requirement` / `data_requirement` / `consistency_requirement`
- `adoption hint`: designer は row 必須列の確認を export 成功条件の詳細候補として扱える。
- `conflict hint`: row 生成順序、欠損列の扱い、preview と永続化 row の一致条件は state-transition / failure / external-integration 観点と統合が必要になる。

### CAND-TOA-006 cached の xTranslator Status 写像を含む成果物を確認する

- `source requirement`: `docs/spec.md:29-32` は完全一致の辞書置換を内部出力ステータス `cached` とし、xTranslator 互換形式では `Status=1` へ写像するとする。`docs/spec.md:65-67` は xTranslator row の `Status` 再構成を求める。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TOA-006`
- `actor`: 辞書置換を含む翻訳成果物を確認するユーザー
- `trigger`: cached 相当の field を含む completed job で output row preview または export summary を開く。
- `expected outcome`: 内部の `cached` が xTranslator row の `Status=1` として出力され、辞書置換であることは別の内部観測情報として確認できる。
- `observable point`: output row preview、`XTRANSLATOR_OUTPUT_ROW.Status`、internal output status summary、dictionary replacement summary
- `related detail requirement type`: `compatibility_requirement` / `observability_requirement` / `data_requirement`
- `adoption hint`: designer は cached 写像を xTranslator 互換性の代表的な正常系候補として扱える。
- `conflict hint`: 内部観測情報を UI に出す範囲と保存対象は operation-audit 観点と統合が必要になる。

### CAND-TOA-007 入力ファイルの出自情報を見て正しい成果物対象を確認する

- `source requirement`: `tasks/usecases/translation-output-artifact.yaml:10-14` は入力として完了済み翻訳ジョブ、訳文、出力ステータス、入力ファイルの出自情報を挙げる。`docs/spec.md:13-15` は複数入力を独立した翻訳ジョブとして管理し、1 job が 1 xEdit 抽出データを対象にして出自を保持することを求める。`docs/er.md:20` は `TRANSLATION_JOB` が 1 つの `X_EDIT_EXTRACTED_DATA` だけを参照するとする。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TOA-007`
- `actor`: 複数 mod の翻訳成果物を取り違えたくないユーザー
- `trigger`: Output Review で completed job と output artifact を選ぶ。
- `expected outcome`: 選択した completed job と出力成果物が、対象の入力データ由来であることを確認できる。
- `observable point`: input provenance summary、selected job summary、`TRANSLATION_JOB`、`X_EDIT_EXTRACTED_DATA`、`TRANSLATION_ARTIFACT`
- `related detail requirement type`: `data_requirement` / `consistency_requirement` / `display_requirement`
- `adoption hint`: designer は入力出自確認を export 誤操作防止の user-facing 候補として扱える。
- `conflict hint`: ファイルパスや抽出元情報の表示粒度は security / operation-audit 観点と統合が必要になる。

### CAND-TOA-008 出力済み artifact の状態を見て再出力へ進む

- `source requirement`: `tasks/usecases/translation-output-artifact.yaml:15-19` は再出力状態を outputs に含める。`tasks/usecases/translation-output-artifact.yaml:22-35` は output artifact の状態、再出力導線、manual check の出力状態確認を求める。`docs/er.md:73-77` は job が生成する成果物と xTranslator 出力 row を出力モデルにする。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TOA-008`
- `actor`: 修正後または確認後に成果物を再出力したいユーザー
- `trigger`: 既存の output artifact がある completed job を Output Review で開く。
- `expected outcome`: 現在の output artifact 状態を確認し、再出力できるかを判断できる。
- `observable point`: output artifact status、re-output action enablement、`TRANSLATION_ARTIFACT`、result summary
- `related detail requirement type`: `alternative_success_requirement` / `display_requirement` / `state_requirement`
- `adoption hint`: designer は再出力導線を正常系の代替成功候補として扱える。
- `conflict hint`: 既存 artifact を上書きするか、新しい artifact として履歴化するかは lifecycle / operation-audit 観点または human decision candidate に残す。

## Open Notes

- `human decision candidate`:
  - `tasks/usecases/translation-output-artifact.yaml:6-8` は `app-shell.md` と `output-review.md` を関連画面に挙げるが、`docs/screen-design/README.md:11-16` は新規 task の UI 判断を task-local `ui-design.md` に置くとしている。Output Review の詳細 UI 契約は designer 側で固定する必要がある。
  - 再生成操作が field result の再翻訳なのか、artifact row の再構成なのかを人間判断に残す。
  - 再出力時に既存 artifact を上書きするのか、新しい artifact として履歴化するのかを人間判断に残す。
  - 入力ファイルの出自情報として、フルパス、ファイル名、plugin 名、digest のどこまでを UI 表示するかを人間判断に残す。
- `merge candidate`:
  - `CAND-TOA-001` と `CAND-TOA-002` は Output Review の completed job 選択と result summary 表示として統合される可能性がある。
  - `CAND-TOA-004` と `CAND-TOA-005` は xTranslator 互換 export 成功条件として統合される可能性がある。
  - `CAND-TOA-003` と `CAND-TOA-008` は diff preview から再出力導線へ進む操作として統合される可能性がある。
- `rejection candidate`:
  - `CAND-TOA-006` は compatibility 観点または operation-audit 観点で十分に扱う場合、actor-goal 側では xTranslator row 確認候補へ縮小できる。
  - `CAND-TOA-007` は UI 表示対象が designer 判断で除外される場合、operation-audit 側の観測候補へ移せる。
