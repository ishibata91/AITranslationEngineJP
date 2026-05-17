# Scenario Candidates: 2026-05-17-master-persona-componentization / failure

- `generator`: `failure`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `MPCF`

## Generator Scope

- `viewpoint`: 失敗観点
- `included_sources`: `plan.md`, `docs/index.md`, `docs/spec.md`, `docs/architecture.md`, `docs/er.md`, `docs/coding-guidelines-frontend.md`, `docs/coding-guidelines-tests.md`, `docs/UX-standard.md`, `docs/screen-design/screens/master-persona.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/scenario-tests/README.md`, `docs/exec-plans/completed/2026-05-18-storybook-foundation/implementation-scope.md`, `docs/exec-plans/completed/2026-05-18-storybook-foundation/storybook-review.md`, `frontend/src/ui/screens/master-persona/`
- `excluded_sources`: プロダクトコード変更、プロダクトテスト変更、docs 正本変更、他観点の候補、採否判断、統合判断、最終シナリオ表
- `generation_notes`: Storybook 基盤完了後のマスターペルソナ画面部品化について、失敗入力、参照不能、設定不整合、保存失敗、回復動作だけを候補化する。Storybook POC は fixed props または view model fixture で確認できる範囲に限定し、Gateway、Wails runtime、AI provider、secret store、DB の実行再現は候補に含めない。

## Candidate Scenarios

### CAND-MPCF-001 story が controller または Gateway を要求して render できない

- `source requirement`: `plan.md` の goal。`MasterPersonaPage` を薄くし、表示部品を story 化可能な props 境界へ寄せる。
- `viewpoint`: 参照不能、設定不整合
- `candidate scenario id`: `CAND-MPCF-001`
- `actor`: 実装者
- `failure start condition`: マスターペルソナ画面の story が `createController`、Gateway、Store、generated `wailsjs`、RuntimeEventAdapter のいずれかを必要とする。
- `rejected operation`: Storybook 上で表示部品を fixed props または view model fixture だけで render する操作。
- `expected error`: Storybook build または story render が失敗する。失敗しない場合でも、対象 story は props 境界違反として扱う。
- `expected outcome`: story は backend、Wails runtime、AI provider、secret store、DB なしで表示できる境界へ戻す。
- `observable point`: story import、fixture import、Storybook iframe の render 結果、`npm --prefix frontend run build-storybook` の結果。
- `related detail requirement type`: `failure_handling_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: 採否は `designer` が判断する。Storybook POC を props 境界の証明にする場合の基本候補である。
- `conflict hint`: page component 全体を story 化する候補と、panel/card/modal だけを story 化する候補の検証単位が競合する可能性がある。

### CAND-MPCF-002 story fixture が主要状態を表現できない

- `source requirement`: `plan.md` の未決事項。現行 view model だけで、未設定、生成中、生成成功、生成失敗、編集中をすべて表現できるか。
- `viewpoint`: 失敗入力、参照不能、人間判断候補
- `candidate scenario id`: `CAND-MPCF-002`
- `actor`: `designer`
- `failure start condition`: Storybook fixture が未設定、生成中、生成成功、生成失敗、編集中のいずれかを作れない。
- `rejected operation`: 主要状態を Storybook 上で確認済みとして扱う操作。
- `expected error`: 対象状態を未確認として記録する。fixture を推測で補って確認済みにしない。
- `expected outcome`: 不足状態は `scenario-design` または `ui-design` の人間判断候補へ残す。
- `observable point`: story 一覧、fixture file、review URL、確認状態、未確認理由。
- `related detail requirement type`: `state_requirement`, `testability_requirement`, `failure_handling_requirement`
- `adoption hint`: 採否は `designer` が判断する。Storybook POC の確認対象状態を固定する候補である。
- `conflict hint`: 最小 POC を 1 story だけにする候補と、主要状態を全部 story 化する候補が検証範囲で競合する可能性がある。

### CAND-MPCF-003 JSON 未選択または preview 失敗で生成開始を拒否する

- `source requirement`: `docs/screen-design/screens/master-persona.md` の入力 JSON パネル。`ペルソナを作成` は作成開始可能状態の場合だけ操作できる。
- `viewpoint`: 失敗入力、参照不能
- `candidate scenario id`: `CAND-MPCF-003`
- `actor`: 利用者
- `failure start condition`: JSON が未選択、preview が未取得、preview が失敗、または preview status が `生成可能` ではない。
- `rejected operation`: `ペルソナを作成` を実行する。
- `expected error`: JSON 選択または preview 失敗の理由を画面通知として表示する。生成開始要求は送信しない。
- `expected outcome`: 候補数、新規作成数、既存スキップ数は失敗状態と整合し、`ペルソナを作成` は操作不能のままになる。
- `observable point`: `master-persona-input-json-panel`、`executeGenerationButton` の disabled 状態、通知文、preview 件数、生成開始 callback の未発火。
- `related detail requirement type`: `failure_handling_requirement`, `boundary_requirement`, `state_requirement`, `testability_requirement`
- `adoption hint`: 採否は `designer` が判断する。入力 JSON と preview gate を UI 人間操作 E2E へ接続する候補である。
- `conflict hint`: JSON 選択直後に自動 preview する候補と、手動 preview を分ける候補が検証手順で競合する可能性がある。

### CAND-MPCF-004 AI 設定未完了または model list 参照不能で生成開始を拒否する

- `source requirement`: `plan.md` の goal、`docs/spec.md` の AI 実行基盤要件、`docs/screen-design/screens/master-persona.md` の AI 設定カード。
- `viewpoint`: 失敗入力、参照不能、設定不整合
- `candidate scenario id`: `CAND-MPCF-004`
- `actor`: 利用者
- `failure start condition`: AI サービス未選択、model 未選択、model list 未取得、model list 取得失敗、または AI 設定更新中である。
- `rejected operation`: `ペルソナを作成`、または未選択 model の保存済み扱いを実行する。
- `expected error`: 設定不足または model list 参照不能を短い日本語文で表示する。API key、secret、raw request、raw response は表示しない。
- `expected outcome`: model 選択、設定保存、生成開始の可否は view model と `AIModelSelectionCard` の props で一致する。
- `observable point`: `master-persona-ai-settings-card`、credential 状態、model list 状態、`executeGenerationButton` の disabled 状態、Storybook fixture の secret 非混入。
- `related detail requirement type`: `failure_handling_requirement`, `security_requirement`, `state_requirement`, `consistency_requirement`
- `adoption hint`: 採否は `designer` が判断する。共有 `AIModelSelectionCard` をマスターペルソナ画面へ再利用する失敗系候補である。
- `conflict hint`: LM Studio のように API key 不要な provider を API key 不足として扱うと、provider 別仕様と競合する。

### CAND-MPCF-005 生成中に編集、削除、再生成を拒否する

- `source requirement`: `docs/screen-design/screens/master-persona.md` の進行状況パネルとペルソナ詳細パネル。生成中は編集と削除を行えない。
- `viewpoint`: 設定不整合、状態整合性、回復動作
- `candidate scenario id`: `CAND-MPCF-005`
- `actor`: 利用者
- `failure start condition`: `runStatus.runState` が `生成中` で、一覧または詳細に選択中ペルソナがある。
- `rejected operation`: 編集開始、削除開始、生成の再開始を実行する。
- `expected error`: 生成中は編集、削除、再生成を行えないことを表示または disabled 理由として示す。
- `expected outcome`: `一時停止` と `中止` だけが進行中操作として有効になり、既存一覧と詳細は確認用途にとどまる。
- `observable point`: `master-persona-progress-panel`、`editButton`、`deleteButton`、`executeGenerationButton`、`interruptGenerationButton`、`cancelGenerationButton` の操作可否。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `failure_handling_requirement`
- `adoption hint`: 採否は `designer` が判断する。panel 分割後も `isRunActive` と `canMutate` の不変条件を崩さない候補である。
- `conflict hint`: 生成中も検索やページ移動を許す候補と、生成中の一覧操作を全部止める候補が操作範囲で競合する可能性がある。

### CAND-MPCF-006 生成状態取得失敗または停止操作失敗を成功扱いにしない

- `source requirement`: `docs/screen-design/screens/master-persona.md` の進行状況パネル。進行状態、処理件数、現在の対象、停止操作を表示する。
- `viewpoint`: 参照不能、回復動作
- `candidate scenario id`: `CAND-MPCF-006`
- `actor`: 利用者
- `failure start condition`: 生成状態取得、一時停止、中止のいずれかが失敗する。
- `rejected operation`: 失敗した取得または停止操作を成功として表示する。
- `expected error`: 生成状態取得失敗、一時停止失敗、中止失敗を画面通知に表示する。
- `expected outcome`: 直前に観測できていた進行状態と件数を勝手に完了扱いへ変えない。
- `observable point`: notice banner、`runStatus.message`、進捗バー、処理済み件数、現在の対象、停止操作 callback の結果。
- `related detail requirement type`: `failure_handling_requirement`, `recovery_requirement`, `observability_requirement`
- `adoption hint`: 採否は `designer` が判断する。進行状況 panel の失敗状態 story を追加する候補である。
- `conflict hint`: 失敗時に自動再読込する候補と、手動再実行だけを許す候補が回復動作で競合する可能性がある。

### CAND-MPCF-007 選択中ペルソナなしで編集または削除を拒否する

- `source requirement`: `docs/screen-design/screens/master-persona.md` のペルソナ詳細パネル、編集モーダル、削除モーダル。
- `viewpoint`: 失敗入力、設定不整合
- `candidate scenario id`: `CAND-MPCF-007`
- `actor`: 利用者
- `failure start condition`: `selectedEntry` が存在しない、または選択中 identity が一覧から消えている。
- `rejected operation`: 編集モーダルを保存する、削除モーダルで削除を確定する。
- `expected error`: 一覧からペルソナを選ぶ必要があることを表示する。削除または更新要求は送信しない。
- `expected outcome`: 詳細パネルは未選択表示になり、編集と削除は操作不能になる。
- `observable point`: `master-persona-persona-detail-panel`、`master-persona-edit-modal`、`master-persona-delete-modal`、`editButton`、`deleteButton`、保存または削除 callback の未発火。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `consistency_requirement`
- `adoption hint`: 採否は `designer` が判断する。modal を単体 story 化する時の null 選択 fixture 候補である。
- `conflict hint`: modal 単体 story で placeholder 表示を許す候補と、modal は選択済み状態だけを story 化する候補が競合する可能性がある。

### CAND-MPCF-008 編集入力不足または保存失敗を反映済み扱いにしない

- `source requirement`: `docs/screen-design/screens/master-persona.md` の編集モーダル。保存した内容が選択中ペルソナと一覧に反映される。
- `viewpoint`: 失敗入力、保存失敗、回復動作
- `candidate scenario id`: `CAND-MPCF-008`
- `actor`: 利用者
- `failure start condition`: ペルソナ要約またはペルソナ本文が空、または保存処理が失敗する。
- `rejected operation`: 編集内容を保存済みとして一覧と詳細へ反映する。
- `expected error`: 入力不足または保存失敗を画面通知に表示する。
- `expected outcome`: 編集モーダルは閉じず、入力値は再試行できる状態で残る。既存の選択中ペルソナは未保存内容で上書きされない。
- `observable point`: `master-persona-edit-modal`、notice banner、`saveEntryButton`、一覧行、詳細本文、保存 callback の結果。
- `related detail requirement type`: `failure_handling_requirement`, `data_requirement`, `recovery_requirement`
- `adoption hint`: 採否は `designer` が判断する。modal 部品の保存失敗 story を追加する候補である。
- `conflict hint`: 保存失敗時に modal を閉じる設計を採る場合、回復動作とユーザー入力保持の扱いが競合する。

### CAND-MPCF-009 削除失敗を一覧から削除済み扱いにしない

- `source requirement`: `docs/screen-design/screens/master-persona.md` の削除モーダル。削除したペルソナを一覧から外す。
- `viewpoint`: 保存失敗、回復動作
- `candidate scenario id`: `CAND-MPCF-009`
- `actor`: 利用者
- `failure start condition`: 削除処理が失敗する。
- `rejected operation`: 失敗した削除を削除済みとして一覧、件数、詳細へ反映する。
- `expected error`: 削除失敗を画面通知に表示する。
- `expected outcome`: 削除モーダルは閉じず、対象の識別情報を確認できる状態を維持する。件数範囲と選択中ペルソナは削除前と整合する。
- `observable point`: `master-persona-delete-modal`、notice banner、`confirmDeleteButton`、一覧件数、選択中詳細、削除 callback の結果。
- `related detail requirement type`: `failure_handling_requirement`, `data_requirement`, `recovery_requirement`
- `adoption hint`: 採否は `designer` が判断する。削除 modal の失敗状態 story を追加する候補である。
- `conflict hint`: 削除失敗時に modal を閉じるか、開いたままにするかは人間判断候補である。

### CAND-MPCF-010 長い状態文、model 名、NPC 名で Storybook 確認が破綻する

- `source requirement`: `docs/UX-standard.md` の長文耐性、`docs/detail-specs/persona-generation-phase.md` の長い NPC 名、provider 名、model 名、error reason の表示耐性。
- `viewpoint`: 失敗入力、競合候補
- `candidate scenario id`: `CAND-MPCF-010`
- `actor`: 人間レビュー担当者
- `failure start condition`: story fixture に長い model 名、長い NPC 名、長い error message、長い persona body のいずれかが入る。
- `rejected operation`: Storybook の表示確認を完了扱いにする。
- `expected error`: レイアウト崩れ、ボタン文字のはみ出し、カード内の重なり、モーダルの操作不能があれば review 未通過として記録する。
- `expected outcome`: 長文状態でも主要情報と主要操作が失われない。
- `observable point`: Storybook iframe screenshot、カード、一覧行、詳細本文、notice banner、modal、最小幅または標準幅の表示状態。
- `related detail requirement type`: `boundary_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: 採否は `designer` が判断する。Storybook POC を見た目レビュー入口にする場合の失敗系候補である。
- `conflict hint`: 長文 fixture を必須にすると最小 POC の scope と競合する可能性がある。

### CAND-MPCF-011 story または review 記録に secret、絶対 path、実ユーザーデータが混入する

- `source requirement`: `docs/exec-plans/completed/2026-05-18-storybook-foundation/implementation-scope.md` の `secret_boundary` と `storybook-review.md` の safety check。
- `viewpoint`: 設定不整合、保存失敗、競合候補
- `candidate scenario id`: `CAND-MPCF-011`
- `actor`: 実装者
- `failure start condition`: story ID、fixture、review URL、確認記録、screenshot 名、コマンド記録に secret、API key、token、ローカル絶対 path、実ユーザーデータが含まれる。
- `rejected operation`: Storybook review 証跡として保存する。
- `expected error`: 証跡保存前に安全条件違反として扱い、該当値を含む成果物を確認済みにしない。
- `expected outcome`: story と review 記録は synthetic data だけを使い、秘密値と実ユーザーデータを含まない。
- `observable point`: story fixture、Storybook URL、story ID、storybook review artifact、screenshot path、command 記録。
- `related detail requirement type`: `security_requirement`, `observability_requirement`, `testability_requirement`
- `adoption hint`: 採否は `designer` が判断する。Storybook 基盤の安全条件をマスターペルソナ POC へ継承する候補である。
- `conflict hint`: 実在 JSON サンプルを fixture に使う候補と、synthetic fixture だけを使う候補が安全条件で競合する可能性がある。

### CAND-MPCF-012 Storybook review URL または story ID が記録されない

- `source requirement`: `plan.md` の close_conditions と `2026-05-18-storybook-foundation` の Storybook review 記録方針。
- `viewpoint`: 保存失敗、回復動作
- `candidate scenario id`: `CAND-MPCF-012`
- `actor`: 人間レビュー担当者
- `failure start condition`: Master Persona の panel、card、modal story は存在するが、review URL、iframe URL、story ID、確認状態、未確認理由が task-local 成果物へ記録されない。
- `rejected operation`: 人間レビュー担当者が Storybook 上の確認対象を追跡する操作。
- `expected error`: 見た目レビュー未準備として扱い、確認済み状態へ進めない。
- `expected outcome`: review URL、story ID、確認状態、未確認理由、再実行 command を task-local 成果物へ残す。
- `observable point`: `storybook-review.md` または `ui-design.md#storybook-review`、Storybook index、対象 story ID、確認状態。
- `related detail requirement type`: `observability_requirement`, `testability_requirement`, `recovery_requirement`
- `adoption hint`: 採否は `designer` が判断する。review 記録先が未決であるため、候補として残す。
- `conflict hint`: review 記録を `ui-design.md` に置くか別 artifact に置くかは未決であり、人間判断候補である。

## Open Notes

- `human decision candidate`: 主要状態の story 数を、未設定、生成中、生成成功、生成失敗、編集中の全件にするか、最小 POC に絞るかは未確定である。
- `human decision candidate`: page component の story を作るか、panel、card、modal だけを story 化するかは未確定である。
- `human decision candidate`: 削除失敗時に modal を閉じるか、開いたまま再試行可能にするかは未確定である。
- `human decision candidate`: Storybook review 記録先を `ui-design.md` と別 artifact のどちらにするかは未確定である。
- `merge candidate`: `CAND-MPCF-001` と `CAND-MPCF-002` は Storybook props 境界と fixture 網羅の候補として統合される可能性がある。
- `merge candidate`: `CAND-MPCF-003` と `CAND-MPCF-004` は生成開始 gate の失敗候補として統合される可能性がある。
- `merge candidate`: `CAND-MPCF-008` と `CAND-MPCF-009` は modal 保存失敗と回復動作の候補として統合される可能性がある。
- `rejection candidate`: 正常系の裏返しだけで、失敗開始条件、拒否される操作、期待エラー、観測点を持たない候補は採否判断前に除外対象になりうる。

## Completion Summary

- `viewpoint`: 失敗観点
- `candidate_count`: 12
- `artifact_path`: `docs/exec-plans/active/2026-05-17-master-persona-componentization/scenario-candidates.failure.md`
- `task_artifact_root`: `docs/exec-plans/active/2026-05-17-master-persona-componentization/`
- `target_diff`: Storybook 基盤を前提にしたマスターペルソナ生成画面の部品化と Storybook POC。
- `remaining_risk`: 採否、統合、競合解消、質問票化は `designer` が行う。
