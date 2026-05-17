# Scenario Candidates: 2026-05-17-master-persona-componentization / state-transition

- `generator`: `state-transition`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `MPSC`

## Generator Scope

- `viewpoint`: state-transition
- `included_sources`: `./plan.md`, `docs/index.md`, `docs/architecture.md`, `docs/coding-guidelines-frontend.md`, `docs/coding-guidelines-tests.md`, `docs/UX-standard.md`, `docs/exec-plans/completed/2026-05-18-storybook-foundation/implementation-scope.md`, `docs/exec-plans/completed/2026-05-18-storybook-foundation/storybook-review.md`, `docs/spec.md`, `docs/er.md`, `docs/screen-design/screens/master-persona.md`, `docs/detail-specs/persona-generation-phase.md`, `frontend/src/ui/screens/master-persona/`
- `excluded_sources`: プロダクトコード変更、プロダクトテスト変更、docs 正本化、他の `scenario-candidates.*.md`
- `generation_notes`: マスターペルソナ生成画面の部品化と Storybook POC は、表示領域を props 境界へ移しても画面状態遷移を維持する必要がある。このファイルは state-transition 観点の候補だけを並べる。最終シナリオの採否、統合、却下、競合解消は `designer` が扱う。

## Candidate Scenarios

### CAND-MPSC-001 入力待ちから JSON 選択済みへ移る

- `source requirement`: `plan.md` は `MasterPersonaPage` を部品合成へ寄せ、表示部品へ画面 view model 由来の小さい props を渡すと定義している。`docs/screen-design/screens/master-persona.md` は生成準備パネルと入力 JSON パネルの状態を定義している。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-MPSC-001`
- `actor`: user
- `trigger`: 利用者が初期の入力待ち状態で JSON ファイルを選択する。
- `pre-transition state`: `selectedFileReference` は null、`selectedFileName` は `未選択`、`preview` は null、`runStatus.runState` は `入力待ち`、`canStartGeneration` は false である。
- `start condition`: ページ合成側からファイル選択 callback が渡され、生成準備部品が現在の view model 状態を受け取っている。
- `post-transition state`: `selectedFileReference` が設定され、`selectedFileName` が選択ファイル名へ変わり、preview 件数を表示でき、ファイル状態ラベルが `入力待ち` から `JSON 選択済み` へ変わる。
- `forbidden transition`: 生成準備部品は、親 view model とずれる独自ファイル状態を作らない。
- `expected outcome`: 入力 JSON パネルと生成準備パネルは、プロダクト画面と Storybook fixture で同じ選択ファイル状態を表示する。
- `observable point`: `master-persona-generation-setup-panel`, `master-persona-input-json-panel`, `対象件数`, `chooseJsonButton`, `resetJsonButton`, `executeGenerationButton`
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: この候補は、ファイル選択と preview 状態の props 境界を保護する。
- `conflict hint`: lifecycle 観点が preview 作成を別シナリオで扱う場合、`designer` はフロー全体を重複させず、この候補を状態句として統合する。

### CAND-MPSC-002 AI 設定更新中は選択操作を固定する

- `source requirement`: `docs/screen-design/screens/master-persona.md` は AI 設定更新中に AI 設定操作を無効にすると定義している。`frontend/src/ui/screens/master-persona/GenerationSetupPanel.svelte` は `isAISettingsRefreshing` とモデルカード view model 値を受け取っている。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-MPSC-002`
- `actor`: user
- `trigger`: 利用者が AI 設定カードからモデル一覧更新を開始する。
- `pre-transition state`: `isAISettingsRefreshing` は false であり、provider、model、execution の各操作はモデルカード状態に従って有効である。
- `start condition`: provider が選択され、更新 callback が渡されている。
- `post-transition state`: `isAISettingsRefreshing` が true になり、provider、model、execution の変更はページ handler で無視され、更新操作と保存関連操作は無効または更新中表示になる。
- `forbidden transition`: 利用者が provider を変更した後に、遅れて返ったモデル一覧応答が provider、model options、カード状態を変更しない。
- `expected outcome`: AI 設定カードは、入力 JSON 状態や実行状態を変えずに更新中状態を表示する。
- `observable point`: `master-persona-ai-settings-card`, `モデル一覧を更新`, model select の disabled 状態, provider select の disabled 状態, 状態文または警告文。
- `related detail requirement type`: `state_requirement`, `concurrency_requirement`, `consistency_requirement`
- `adoption hint`: この候補は、gateway 呼び出しなしで検査できる Storybook 状態として扱いやすい。
- `conflict hint`: external-integration 観点が provider 失敗を扱う可能性がある。この候補は UI 状態遷移と遅延応答の破棄だけを扱う。

### CAND-MPSC-003 生成開始可能から生成中へ移る

- `source requirement`: `docs/screen-design/screens/master-persona.md` は JSON、AI 設定、作成可能件数、実行状態が開始条件を満たす時だけ `ペルソナを作成` を有効にすると定義している。`frontend/src/application/presenter/master-persona/master-persona.presenter.ts` は `canStartGeneration` と `isRunActive` を導出している。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-MPSC-003`
- `actor`: user
- `trigger`: 利用者が `ペルソナを作成` を押す。
- `pre-transition state`: JSON は選択済み、AI 設定は完了済み、preview 状態は `生成可能`、`runStatus.runState` は `生成中` ではない。
- `start condition`: `canStartGeneration` は true である。
- `post-transition state`: `runStatus.runState` は `生成中` になり、`isRunActive` は true になり、`canStartGeneration` は false になり、進捗と現在対象の表示が更新される。
- `forbidden transition`: `isRunActive` が true の間、2 回目の開始操作を可能にしない。
- `expected outcome`: 作成ボタンは無効になり、進行状況パネルは実行中状態になり、実行中だけの操作が有効になる。
- `observable point`: `executeGenerationButton`, `master-persona-progress-panel`, `runProgressFill`, `interruptGenerationButton`, `cancelGenerationButton`
- `related detail requirement type`: `state_requirement`, `冪等性_requirement`, `consistency_requirement`
- `adoption hint`: この候補は、生成準備パネル props と進行状況パネル props の分離を保護する。
- `conflict hint`: failure 観点が開始失敗を扱う可能性がある。この候補は provider や backend 実行の振る舞いを決めない。

### CAND-MPSC-004 生成中から非生成中へ戻る

- `source requirement`: `docs/screen-design/screens/master-persona.md` は、生成中だけ `一時停止` と `中止` を有効にし、実行状態に応じて一覧と詳細の操作可否を変えると定義している。`frontend/src/application/usecase/master-persona/master-persona.usecase.ts` は実行中状態の終了後に一覧を再読み込みしている。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-MPSC-004`
- `actor`: user
- `trigger`: 利用者が `一時停止` または `中止` を押す。もしくは polling が生成終了を観測する。
- `pre-transition state`: `runStatus.runState` は `生成中`、`isRunActive` は true、変更操作は無効である。
- `start condition`: 中断、停止、実行状態取得の callback が非生成中の run status を返す。
- `post-transition state`: `isRunActive` は false になり、実行中だけの操作は無効になり、一覧と詳細は再読み込みされ、選択中ペルソナがある場合だけ変更操作が戻る。
- `forbidden transition`: 過去の story または部品ローカル状態が生成中だと見なすだけで、画面が編集と削除を無効のままにしない。
- `expected outcome`: 実行状態、進捗件数、一覧内容、選択詳細、編集と削除の可否が遷移後に整合する。
- `observable point`: `master-persona-progress-panel`, `master-persona-generation-result-list-panel`, `master-persona-persona-detail-panel`, `editButton`, `deleteButton`
- `related detail requirement type`: `state_requirement`, `recovery_requirement`, `consistency_requirement`
- `adoption hint`: この候補は、状態を持つ画面シナリオ、または実行中状態と実行後状態の Storybook fixture 対として扱える。
- `conflict hint`: lifecycle 観点が生成完了を広いフローとして扱う可能性がある。この候補は UI 状態遷移の句に留める。

### CAND-MPSC-005 一覧選択から詳細選択へ移る

- `source requirement`: `docs/screen-design/screens/master-persona.md` は一覧選択、ページ操作、詳細表示を定義している。`frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte` は一覧、選択、ページ、詳細の値を view model から受け取っている。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-MPSC-005`
- `actor`: user
- `trigger`: 利用者がペルソナ行を選択する。もしくは検索条件、プラグイン条件、ページを変更する。
- `pre-transition state`: `items` は空または存在し、`selectedIdentityKey` は null または現在詳細を指し、`selectedEntry` は null または存在する。
- `start condition`: 選択可能な行が存在する。もしくは filter または page の操作が変更される。
- `post-transition state`: `selectedIdentityKey` と `selectedEntry` は新しい選択行を反映する。絞り込み後に対象行がない場合は null になる。
- `forbidden transition`: 一覧状態が選択 identity を含まない時、詳細パネルは古い selected entry を表示しない。
- `expected outcome`: 一覧選択、件数範囲、詳細見出し、識別情報、編集と削除の可否が一緒に更新される。
- `observable point`: `master-persona-generation-result-list-panel`, 選択行の `aria-pressed`, `master-persona-persona-detail-panel`, `prevPageButton`, `nextPageButton`
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `boundary_requirement`
- `adoption hint`: この候補は、空一覧、一覧あり、行選択済み、絞り込み後の空状態の Storybook fixture に向いている。
- `conflict hint`: actor-goal 観点が検索目的を扱う可能性がある。この候補は一覧と詳細の状態整合だけを保護する。

### CAND-MPSC-006 詳細選択から編集モーダルへ移る

- `source requirement`: `docs/screen-design/screens/master-persona.md` は、選択中ペルソナがあり変更可能な場合だけ編集できると定義している。`frontend/src/application/usecase/master-persona/master-persona.usecase.ts` は `modalState` を `edit` にし、`selectedEntry` から `editForm` を作る。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-MPSC-006`
- `actor`: user
- `trigger`: 利用者が `編集` を押し、その後に保存またはキャンセルする。
- `pre-transition state`: `selectedEntry` が存在し、`canMutate` は true、`modalState` は null である。
- `start condition`: 編集操作が有効である。
- `post-transition state`: `modalState` は `edit` になり、編集フォームは選択中ペルソナを反映し、保存成功時は選択と一覧を更新して `modalState` が null に戻る。
- `forbidden transition`: 保存失敗時にモーダルを閉じず、編集フォーム状態を観測できる状態に保つ。
- `expected outcome`: 編集モーダルは有効な選択からだけ開き、キャンセル、閉じる、保存成功のいずれかで閉じる。
- `observable point`: `editButton`, `master-persona-edit-modal`, `editPersonaSummaryInput`, `editSpeechStyleInput`, `editPersonaBodyInput`, `saveEntryButton`
- `related detail requirement type`: `state_requirement`, `failure_handling_requirement`, `consistency_requirement`
- `adoption hint`: この候補は、モーダル props 境界と編集フォーム状態を保護する。
- `conflict hint`: failure 観点が validation と保存失敗の詳細を追加する可能性がある。この候補はモーダル開閉遷移だけに留める。

### CAND-MPSC-007 詳細選択から削除モーダルへ移る

- `source requirement`: `docs/screen-design/screens/master-persona.md` は、選択中ペルソナがあり変更可能な場合だけ削除できると定義している。`frontend/src/application/usecase/master-persona/master-persona.usecase.ts` は失敗時に `modalState` を維持し、成功時に null へ戻す。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-MPSC-007`
- `actor`: user
- `trigger`: 利用者が `削除` を押し、その後に確定またはキャンセルする。
- `pre-transition state`: `selectedEntry` が存在し、`canMutate` は true、`modalState` は null である。
- `start condition`: 削除操作が有効である。
- `post-transition state`: `modalState` は `delete` になり、削除確認は選択中ペルソナの識別情報を表示し、削除成功時は削除対象を一覧から外して `modalState` が null に戻る。
- `forbidden transition`: 削除失敗時にモーダルを閉じず、選択中ペルソナの識別情報を再試行またはキャンセルのために表示し続ける。
- `expected outcome`: 削除確認状態は選択中エントリに対してだけ表示され、空選択では表示されない。
- `observable point`: `deleteButton`, `master-persona-delete-modal`, `confirmDeleteButton`, 選択中ペルソナ名, FormID, EditorID, 対象プラグイン
- `related detail requirement type`: `state_requirement`, `failure_handling_requirement`, `consistency_requirement`
- `adoption hint`: この候補は、編集モーダルと fixture データを共有しつつ、モーダル状態を分けると扱いやすい。
- `conflict hint`: operation-audit 観点が削除証跡を扱う可能性がある。この候補は監査保持を決めない。

### CAND-MPSC-008 Storybook fixture 状態は controller なしで維持する

- `source requirement`: 完了済み Storybook 基盤の scope は、story と fixture が fixed props または view model fixture を使い、backend DTO、Gateway、RuntimeEventAdapter、AI provider、secret store、DB、filesystem business flow を import しないと定義している。今回の `plan.md` はパネル、カード、モーダルの状態を Storybook review 対象にしている。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-MPSC-008`
- `actor`: implementer
- `trigger`: 実装者が生成準備、進行状況、ペルソナ確認、操作モーダルの story を追加する。
- `pre-transition state`: Storybook は callback props 付きの固定 fixture 状態を描画する。
- `start condition`: Storybook dev server または static build から story を開く。
- `post-transition state`: 描画部品は、Storybook controls またはローカル fixture args が変わらない限り、fixture が定義した状態を維持する。
- `forbidden transition`: story 表示時に `MasterPersonaPage` の production controller を mount せず、Wails 呼び出し、Gateway 呼び出し、provider model 更新、file read、実画面状態の変更を行わない。
- `expected outcome`: Storybook は props-only fixture で、待機中、更新中、生成中、詳細選択済み、編集モーダル、削除モーダルを表示できる。
- `observable point`: Storybook story ID、部品の表示 text、disabled/enabled 状態、story と fixture import に product runtime 呼び出しがないこと。
- `related detail requirement type`: `testability_requirement`, `state_requirement`, `compatibility_requirement`
- `adoption hint`: この候補は、Storybook POC をプロダクト実行面ではなく状態検査面として保護する。
- `conflict hint`: external-integration 観点が provider 隔離を扱う可能性がある。この候補は Storybook で見える状態境界だけを扱う。

## Open Notes

- `human decision candidate`: 現行 view model が、未設定、更新中、JSON 選択済み、生成可能、生成中、完了、失敗、編集中、削除中の全 review 状態を表現できるかは未確定である。`plan.md` も未決事項にしているため、UI design または scenario design が補わない場合は `designer` が質問として残す。
- `human decision candidate`: `screen-design-diff.master-persona.md` は `plan.md` から参照されているが、候補生成時点の active task folder には存在しなかった。`designer` は、この不足をシナリオ確定 blocker にするか、別の UI design artifact gap にするかを判断する。
- `human decision candidate`: Storybook review 記録の配置は `plan.md` で未解決である。候補は `ui-design.md#storybook-review` と `storybook-review.md` のどちらを owner にするかを決めない。
- `merge candidate`: lifecycle 観点が生成開始と完了をすでに扱う場合、`CAND-MPSC-003` と `CAND-MPSC-004` は単一の実行状態シナリオへ統合できる。
- `merge candidate`: `CAND-MPSC-006` と `CAND-MPSC-007` は、編集と削除の振る舞い句を分けたまま、モーダル状態 fixture matrix を共有できる。
- `rejection candidate`: 実 provider 実行、secret 解決、backend 永続化、Wails runtime を必要とする詳細は、この観点から却下し、external-integration または lifecycle の候補へ移す。
