# Scenario Candidates: 2026-05-17-master-persona-componentization / external-integration

- `generator`: `external-integration`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `MPC-EI`

## Generator Scope

- `viewpoint`: `external-integration`
- `included_sources`: `plan.md`, `docs/spec.md`, `docs/architecture.md`, `docs/coding-guidelines-frontend.md`, `docs/coding-guidelines-tests.md`, `docs/UX-standard.md`, `docs/screen-design/screens/master-persona.md`, `docs/detail-specs/ai-provider-settings-management.md`, `docs/er.md`, `docs/exec-plans/completed/2026-05-18-storybook-foundation/implementation-scope.md`, `docs/exec-plans/completed/2026-05-18-storybook-foundation/storybook-review.md`, `frontend/src/ui/screens/master-persona/`
- `excluded_sources`: プロダクトコード変更、プロダクトテスト変更、docs 正本化、他観点の scenario candidate、最終シナリオ採否、統合判断
- `generation_notes`: Storybook 基盤を前提に、マスターペルソナ生成画面の部品化が外部 provider、secret、Wails binding、file、network へ誤接続しない境界を候補化する。

## Candidate Scenarios

### CAND-MPC-EI-001 Storybook の表示部品が Wails binding に接続しない

- `source requirement`: `docs/architecture.md` は `View` と `UI Component` が generated binding と backend DTO を直接扱わないと定義する。Storybook 基盤の `implementation-scope.md` は story と fixture が generated `wailsjs`、Gateway、RuntimeEventAdapter、AI provider、secret store、DB を import しないと定義する。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-MPC-EI-001`
- `actor`: 実装者
- `trigger`: マスターペルソナ画面の `GenerationSetupPanel`、`RunStatusPanel`、`PersonaReviewPanel`、`PersonaActionModal`、または page 合成用 story を Storybook へ追加する。
- `expected outcome`: story と fixture は fixed props、view model fixture、callback stub だけで表示できる。story 表示時に Wails binding、Gateway、RuntimeEventAdapter、backend DTO、DB へ接続しない。
- `fake_or_stub`: callback は no-op または記録用 stub にする。controller を必要とする page story を作る場合は、Wails gateway ではなく画面 contract を満たす stub controller だけを使う。
- `observable point`: `npm --prefix frontend run build-storybook` が backend、Wails runtime、AI provider、secret store、DB なしで完了する。Storybook story source と fixture source に generated `wailsjs` と `@controller/wails` import がない。
- `related detail requirement type`: `testability_requirement`, `compatibility_requirement`, `security_requirement`
- `adoption hint`: 部品化の最小成立条件として採用しやすい。
- `conflict hint`: page story を要求する候補が controller の実接続まで求める場合、external-integration 観点では stub controller までに制限する。

### CAND-MPC-EI-002 AI サービス設定カードの story が secret 本体を持たない

- `source requirement`: `docs/spec.md` は APIKey を暗号化保存すると定義する。`docs/detail-specs/ai-provider-settings-management.md` は APIキー本体と raw payload を UI、DTO、要約、log、debug 出力に出さないと定義する。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-MPC-EI-002`
- `actor`: 実装者
- `trigger`: `GenerationSetupPanel` または `AIModelSelectionCard` の master-persona story fixture を作成する。
- `expected outcome`: fixture は provider、model、資格情報の状態分類だけを持つ。API key、token、secret、raw request、raw response、raw prompt を story、fixture、review 記録、スクリーンショット確認文へ含めない。
- `fake_or_stub`: 資格情報は `configured`、`missing`、`not_required` などの状態値だけで表現する。secret 本体の fake 値は置かない。
- `observable point`: story 画面には資格情報の存在状態または警告文だけが表示される。fixture source と `storybook-review.md` に secret 値、API key 例、token 例がない。
- `related detail requirement type`: `security_requirement`, `data_requirement`, `testability_requirement`
- `adoption hint`: AI 設定カードと Storybook review の安全条件として採用しやすい。
- `conflict hint`: visual review 用に「APIキー入力済み」を見せたい場合でも、表示は状態分類へ限定する。

### CAND-MPC-EI-003 モデル一覧更新の story が provider network を呼ばない

- `source requirement`: `docs/detail-specs/ai-provider-settings-management.md` は実装後検証で fake transport DI と fake secret store を使い、有料の実 AI API を呼ばないと定義する。Storybook 基盤の `implementation-scope.md` は story と fixture が AI provider を import しないと定義する。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-MPC-EI-003`
- `actor`: 実装者
- `trigger`: モデル一覧更新中、更新成功、資格情報不足、provider 到達不能の表示 story を追加する。
- `expected outcome`: story は `isAISettingsRefreshing`、`modelListStatusText`、`modelOptions`、`credentialWarningText` などの props だけで状態を示す。provider network、real API、endpoint、secret store を呼ばない。
- `fake_or_stub`: `refreshAISettings` は resolved Promise を返す stub にする。失敗表示は callback 実行ではなく fixture の view model 状態で表す。
- `observable point`: Storybook 上で `モデル一覧を更新` を押しても外部通信が発生しない。build-storybook はネットワーク接続なしの環境で成立する。
- `related detail requirement type`: `testability_requirement`, `failure_handling_requirement`, `security_requirement`
- `adoption hint`: 更新中と失敗表示を Storybook で確認する場合に採用しやすい。
- `conflict hint`: 失敗シナリオ候補が provider failure を UI 操作で再現させる場合、external-integration 観点では fixture 状態で再現する。

### CAND-MPC-EI-004 入力 JSON パネルの story が実ファイルパスへ依存しない

- `source requirement`: `docs/screen-design/screens/master-persona.md` は入力 JSON の選択、選び直し、候補数、新規作成数、既存スキップ数を表示すると定義する。Storybook 基盤の `implementation-scope.md` は fixture と review 記録にローカル絶対 path と実ユーザーデータを含めないと定義する。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-MPC-EI-004`
- `actor`: 実装者
- `trigger`: `GenerationSetupPanel` の未選択、JSON 選択済み、preview 取得済み、作成開始不可の story を追加する。
- `expected outcome`: story は `selectedFileName`、`selectedFileReference`、`preview` の fixture 値だけで表示する。実ファイル選択、filesystem、ローカル絶対 path、実ユーザー由来 JSON を必須にしない。
- `fake_or_stub`: `handleJsonSelected`、`chooseJsonFile`、`resetJsonSelection`、`startGeneration` は callback stub にする。file reference は `"fixture://master-persona/input.json"` などの非実パス識別子にする。
- `observable point`: Storybook 表示と build は実 JSON ファイルなしで成立する。review 記録と fixture にローカル絶対 path がない。
- `related detail requirement type`: `boundary_requirement`, `testability_requirement`, `security_requirement`
- `adoption hint`: file input の表示状態を Storybook で扱う場合に採用しやすい。
- `conflict hint`: lifecycle 候補が実ファイル読み込み完了までを求める場合、external-integration 観点では Storybook の fixture 境界と実 filesystem 境界を分ける。

### CAND-MPC-EI-005 ペルソナ作成開始の story が AI 実行を開始しない

- `source requirement`: `docs/spec.md` は NPC の発言と属性情報を元に AI にペルソナを生成させられることを要求する。Storybook 基盤の `implementation-scope.md` は backend、Wails runtime、AI provider、secret store、DB の接続実装を扱わないと定義する。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-MPC-EI-005`
- `actor`: 実装者
- `trigger`: `ペルソナを作成` が有効な状態、生成中状態、生成失敗通知の story を追加する。
- `expected outcome`: story は `canStartGeneration`、`runStatus`、`errorMessage`、`progressPercent` の fixture で状態を示す。button callback は AI 実行、Wails binding、provider adapter、DB 書き込みを開始しない。
- `fake_or_stub`: `startGeneration`、`interruptGeneration`、`cancelGeneration` は no-op または action logger 互換の stub にする。生成結果は fixture の `items` と `selectedEntry` で表す。
- `observable point`: Storybook 上の操作は表示確認に閉じる。外部 API 呼び出し、Wails binding 呼び出し、DB 更新を前提にしない。
- `related detail requirement type`: `testability_requirement`, `security_requirement`, `compatibility_requirement`
- `adoption hint`: 生成中、完了、失敗表示を Storybook で増やす場合に採用しやすい。
- `conflict hint`: 受け入れテスト候補が実生成を求める場合、Storybook POC の受け入れ条件からは分離する。

### CAND-MPC-EI-006 Storybook review 記録が review URL だけを外部境界として扱う

- `source requirement`: Storybook 基盤の `implementation-scope.md` は Storybook review URL と確認状態を task-local の `storybook-review.md` に記録すると定義する。完了済み `storybook-review.md` は Storybook localhost URL、story ID、確認状態、未確認理由、build 結果を記録している。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-MPC-EI-006`
- `actor`: 実装者
- `trigger`: マスターペルソナ部品の Storybook review URL と確認状態を記録する。
- `expected outcome`: review URL は Storybook の localhost URL または iframe URL だけを指す。fakeAPI URL、Wails runtime URL、backend API URL、real provider URL を review URL として扱わない。
- `fake_or_stub`: Storybook server は UI review の表示境界として扱う。backend や provider の代替 server は起動しない。
- `observable point`: review 記録に story ID、review URL、確認状態、未確認理由、再実行 command がある。URL、query、story ID に secret、API key、token、ローカル絶対 path、実ユーザーデータがない。
- `related detail requirement type`: `observability_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: 人間見た目レビューを Storybook へ寄せる場合に採用しやすい。
- `conflict hint`: UI 人間操作 E2E が Wails app URL を要求する場合、Storybook review と実アプリ確認の境界を分ける。

### CAND-MPC-EI-007 Storybook fixture が provider 表示対象を実サービス選択に固定しない

- `source requirement`: `docs/spec.md` は Gemini、xAI、LMStudio を利用可能な AI として定義する。`docs/detail-specs/ai-provider-settings-management.md` は利用者向け provider list が `gemini`、`lm_studio`、`xai` だけを扱い、fake provider を表示しないと定義する。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-MPC-EI-007`
- `actor`: 実装者
- `trigger`: AI サービス選択を含む master-persona story fixture を作成する。
- `expected outcome`: fixture は実サービス名の表示を確認できるが、real API 利用、実 endpoint、実 credential を前提にしない。fake provider を利用者向け選択肢として表示しない。
- `fake_or_stub`: provider option は `gemini`、`lm_studio`、`xai` の表示用 fixture にする。接続可否は credential 状態や model list 状態の fixture で表す。
- `observable point`: story 画面の provider 選択肢に fake provider が表示されない。fixture は実 endpoint と secret を持たない。
- `related detail requirement type`: `compatibility_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: 既存 `AIModelSelectionCard` 再利用時の provider 表示安全条件として採用しやすい。
- `conflict hint`: テスト容易性のために fake provider を UI に出す候補がある場合、外部連携観点では user-facing 表示から除外する。

## Open Notes

- `human decision candidate`: page story を作る場合、page component を直接 story 化するか、review-only wrapper を作るかは `plan.md` の未決事項である。external-integration 観点では、どちらでも Wails gateway と provider network へ接続しないことを候補条件にする。
- `human decision candidate`: Storybook review の確認状態を `ui-design.md` 内へ残すか、別 artifact の `storybook-review.md` へ残すかは `plan.md` の未決事項である。external-integration 観点では、記録先に secret、token、ローカル絶対 path、実ユーザーデータを含めないことを候補条件にする。
- `merge candidate`: `CAND-MPC-EI-001`、`CAND-MPC-EI-003`、`CAND-MPC-EI-005` は Storybook の外部接続遮断として統合される可能性がある。
- `merge candidate`: `CAND-MPC-EI-002`、`CAND-MPC-EI-006`、`CAND-MPC-EI-007` は secret と provider 表示の安全条件として統合される可能性がある。
- `rejection candidate`: 部品化対象に AI 設定カードや入力 JSON パネルを含めない場合、対応する個別 story 候補は採用されない可能性がある。
