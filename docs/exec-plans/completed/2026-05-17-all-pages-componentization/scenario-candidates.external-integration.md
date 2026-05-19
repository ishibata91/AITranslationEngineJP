# Scenario Candidates: 2026-05-17-all-pages-componentization / external-integration

- `generator`: `external-integration`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `APC-EI`
- `candidate_count`: 9

## Generator Scope

- `viewpoint`: `external-integration`
- `included_sources`: `./plan.md`, `docs/architecture.md`, `docs/coding-guidelines-frontend.md`, `docs/UX-standard.md`, `docs/exec-plans/completed/2026-05-18-storybook-foundation/scenario-design.md`, `docs/exec-plans/completed/2026-05-17-master-persona-componentization/scenario-design.md`, `docs/exec-plans/completed/2026-05-17-master-persona-componentization/ui-design.md`, `frontend/src/ui/screens/README.md`, `frontend/src/ui/screens/`, `frontend/src/ui/components/`
- `excluded_sources`: プロダクトコード変更、テスト変更、docs 正本変更、最終シナリオ表、候補採否、候補統合、他 agent の成果物。
- `generation_notes`: 全ページ部品化と Storybook story 追加について、Wails gateway、generated binding、Store、ScreenController、Storybook fixture、file 表示値の境界だけを候補化する。Storybook は fixed props または view model fixture を使い、実 backend、Wails runtime、AI provider、secret store、DB、実 filesystem flow を要求しない前提で候補を書く。

## Candidate Scenarios

### CAND-APC-EI-001 page component が controller 接続だけを持ち、表示部品へ外部境界を渡さない

- `source requirement`: `plan.md` は、ページ component を controller 接続と部品合成へ寄せる。`docs/architecture.md` は、View が `ScreenController` へ DOM event を渡し、UI Component は backend DTO、generated binding、`Store`、`Gateway` を直接扱わないと定義する。
- `viewpoint`: external-integration / adapter 境界
- `candidate scenario id`: `CAND-APC-EI-001`
- `actor`: 実装者
- `external boundary`: Wails gateway、generated binding、`Store`、`ScreenController`
- `start condition`: 全ページから panel、card、modal を screen local component または shared component へ切り出す。
- `trigger`: 切り出した表示部品を story または親 page から描画する。
- `expected outcome`: 切り出した表示部品は小さい props と callback だけを受け取り、`Gateway`、generated `wailsjs`、`Store`、controller factory を import しない。
- `fake_or_stub`: fixed props、view model fixture、callback stub。
- `observable point`: story source、fixture source、component import、`npm --prefix frontend run build-storybook` の結果。
- `related detail requirement type`: `compatibility_requirement`, `testability_requirement`, `consistency_requirement`
- `adoption hint`: controller 接続と表示部品の props 境界を確認するシナリオへ統合できる。
- `conflict hint`: page story を必須にする場合、page component が controller factory を要求するため、主要部品 story の検証段階と衝突する可能性がある。

### CAND-APC-EI-002 Storybook story は Wails runtime と generated binding を要求しない

- `source requirement`: Storybook 基盤の固定要件は、Storybook static build が backend、Wails runtime、AI provider、secret store、DB に依存しないことである。frontend 規約は generated `wailsjs` の import を gateway 境界に閉じ込める。
- `viewpoint`: external-integration / generated binding 境界
- `candidate scenario id`: `CAND-APC-EI-002`
- `actor`: 実装者
- `external boundary`: generated `wailsjs`、`globalThis.go.wails`、Wails runtime
- `start condition`: 全ページの主要 panel、card、modal に Storybook story を追加する。
- `trigger`: `npm --prefix frontend run build-storybook` を実行する。
- `expected outcome`: story と fixture は generated `wailsjs`、`globalThis.go.wails`、Wails runtime を参照せず、Storybook static build が成立する。
- `fake_or_stub`: fixed props、view model fixture、callback stub。
- `observable point`: story import、fixture import、build-storybook の exit code、Wails binding 未接続 error の有無。
- `related detail requirement type`: `testability_requirement`, `compatibility_requirement`, `failure_handling_requirement`
- `adoption hint`: Storybook 外部境界の禁止 import 確認として採用できる。
- `conflict hint`: 失敗観点が「render 失敗」として扱う場合、外部連携観点では「Wails runtime 要求の禁止」として統合候補になる。

### CAND-APC-EI-003 Storybook fixture は secret と実 endpoint を持たない

- `source requirement`: Storybook 基盤と Master Persona の固定要件は、fixture、story、review 記録が secret、API key、token、実 endpoint、ローカル絶対 path、実ユーザーデータ、raw request、raw response、raw prompt を含まないことである。`docs/spec.md` は API key を保存でき、暗号化して保存する要件を持つ。
- `viewpoint`: external-integration / secret 境界
- `candidate scenario id`: `CAND-APC-EI-003`
- `actor`: 実装者
- `external boundary`: secret store、AI provider 設定、provider endpoint
- `start condition`: Provider 設定、AI モデル選択、翻訳フェーズ、Master Persona などの story fixture を追加する。
- `trigger`: fixture を story から import して表示する。
- `expected outcome`: fixture は synthetic data だけを持ち、secret、API key、token、実 endpoint、raw provider payload を含まない。
- `fake_or_stub`: synthetic provider 名、synthetic model 名、秘匿値を含まない状態 fixture。
- `observable point`: fixture source、story source、Storybook review URL、screenshot path、保存禁止情報の不在。
- `related detail requirement type`: `security_requirement`, `data_requirement`, `testability_requirement`
- `adoption hint`: Storybook fixture 外部境界の必須受け入れ条件へ統合できる。
- `conflict hint`: Provider 設定画面の story で「保存済み secret 状態」を表示する場合、表示文言だけで表現し、secret 値を fixture に置かない制約と合わせる必要がある。

### CAND-APC-EI-004 shared component は Store、Gateway、screen-specific DTO を持たない

- `source requirement`: `docs/architecture.md` は、複数画面で使う部品だけ `frontend/src/ui/components/` に置き、UI Component は `Store`、`Gateway`、generated binding を直接扱わないと定義する。`plan.md` は共有部品へ上げる候補を事前調査で分ける。
- `viewpoint`: external-integration / adapter 境界
- `candidate scenario id`: `CAND-APC-EI-004`
- `actor`: 実装者
- `external boundary`: shared component と application/controller 層の境界
- `start condition`: 複数ページで使う card、button、footer、status 表示を `frontend/src/ui/components/` へ上げる。
- `trigger`: shared component の story を fixed props で描画する。
- `expected outcome`: shared component は `Store`、`Gateway`、generated binding、screen-specific DTO を import せず、表示規則と callback だけを公開する。
- `fake_or_stub`: shared component 用 props fixture、callback stub。
- `observable point`: `frontend/src/ui/components/` の import、component props 型、shared component story、build-storybook の結果。
- `related detail requirement type`: `compatibility_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: 共有部品化の外部境界確認として採用できる。
- `conflict hint`: 共有化候補を増やしすぎると、画面固有条件を props 分岐へ押し込む候補と衝突する可能性がある。

### CAND-APC-EI-005 page story または App story が production gateway を暗黙生成しない

- `source requirement`: `docs/coding-guidelines-frontend.md` は、production gateway、controller factory、外部 adapter の生成を composition root に置き、View component が production wiring を作らないと定義する。現行 `App.svelte` は一部 fallback wiring を持つため、page story 化では注入境界の扱いが候補になる。
- `viewpoint`: external-integration / Wails gateway 境界
- `candidate scenario id`: `CAND-APC-EI-005`
- `actor`: 実装者
- `external boundary`: `Frontend Bootstrap`、production gateway、controller factory
- `start condition`: page 合成 story、App shell story、または page-level story を追加する。
- `trigger`: Storybook で page story を開く。
- `expected outcome`: story は production gateway を暗黙生成せず、fixture controller または null controller 状態で表示できる。Storybook 表示は Wails binding 未接続 error で止まらない。
- `fake_or_stub`: fixture controller、null controller、fixed view model fixture。
- `observable point`: page story props、App story props、controller factory 注入、Wails binding 未接続 error の有無。
- `related detail requirement type`: `testability_requirement`, `failure_handling_requirement`, `compatibility_requirement`
- `adoption hint`: page story を任意にする場合でも、作る時の外部境界制約として残せる。
- `conflict hint`: page component を production 接続に寄せる方針と、Storybook で page を単体描画する方針が衝突する可能性がある。主要部品 story を優先するかは designer 判断に残す。

### CAND-APC-EI-006 file 表示を fixture 化し、実 filesystem flow を起動しない

- `source requirement`: Storybook 基盤と Master Persona の固定要件は、story と fixture が実 filesystem flow を要求しないことである。全ページ対象には翻訳入力、辞書取り込み、成果物出力など file 境界を持つ画面が含まれる。
- `viewpoint`: external-integration / ファイル境界
- `candidate scenario id`: `CAND-APC-EI-006`
- `actor`: 実装者
- `external boundary`: file picker、import file path、output artifact path、filesystem adapter
- `start condition`: 翻訳入力、辞書、成果物出力など file 参照を表示する画面部品へ story を追加する。
- `trigger`: Storybook で file 参照を含む story を開き、操作 callback を発火可能な状態にする。
- `expected outcome`: story は synthetic file reference を表示し、実 file picker、実 file read、実 output generation を開始しない。ローカル絶対 path は fixture と review URL に含まれない。
- `fake_or_stub`: synthetic file name、synthetic import summary、callback stub。
- `observable point`: fixture source、story callback、Storybook URL、ローカル絶対 path の不在、build-storybook の結果。
- `related detail requirement type`: `security_requirement`, `data_requirement`, `testability_requirement`
- `adoption hint`: file boundary を持つページの Storybook fixture 条件として採用できる。
- `conflict hint`: actor-goal 観点が「利用者がファイルを選ぶ」操作を候補化する場合、Storybook では操作開始ではなく表示状態の fixture として扱う必要がある。

### CAND-APC-EI-007 runtime event と polling adapter を Storybook に持ち込まない

- `source requirement`: `docs/architecture.md` は `RuntimeEventAdapter` を Wails event 購読と screen local handler への写像に限定する。Storybook 基盤は `RuntimeEventAdapter` を story と fixture へ持ち込まない方針である。
- `viewpoint`: external-integration / network 境界
- `candidate scenario id`: `CAND-APC-EI-007`
- `actor`: 実装者
- `external boundary`: Wails runtime event、runtime polling adapter、push 通知
- `start condition`: 進捗、完了、失敗、通知を表示する phase panel、管理画面、Master Dictionary、Master Persona の story を追加する。
- `trigger`: Storybook で実行中、失敗、完了の fixture を表示する。
- `expected outcome`: story は runtime event 購読や polling を開始せず、進捗と通知は fixed view model fixture として表示する。
- `fake_or_stub`: progress fixture、notification fixture、callback stub。
- `observable point`: story source、component import、runtime adapter import の不在、Storybook canvas の進捗表示。
- `related detail requirement type`: `testability_requirement`, `compatibility_requirement`, `observability_requirement`
- `adoption hint`: 進捗表示部品の Storybook 外部境界として採用できる。
- `conflict hint`: lifecycle 観点が実行中の進捗更新を扱う場合、Storybook では実 push ではなく固定状態表示に限定する前提を合わせる必要がある。

### CAND-APC-EI-008 backend DTO ではなく view model を story 入力にする

- `source requirement`: `docs/architecture.md` は、View が backend DTO や generated binding を直接扱わず、view model だけを前提に描画すると定義する。`docs/coding-guidelines-frontend.md` は、Wails bridge の戻り値を使用前に絞り込むことを定義する。
- `viewpoint`: external-integration / adapter 境界
- `candidate scenario id`: `CAND-APC-EI-008`
- `actor`: 実装者
- `external boundary`: backend DTO、GatewayDTO、Presenter view model
- `start condition`: 既存 page の表示領域を story 対象 component へ分ける。
- `trigger`: story fixture を component props として渡す。
- `expected outcome`: story fixture は backend response DTO や GatewayDTO ではなく、Presenter 後の view model または component props を使う。DTO 変換失敗や backend shape 変更を Storybook fixture が肩代わりしない。
- `fake_or_stub`: view model fixture、component props fixture。
- `observable point`: fixture type import、story source、component props 型、GatewayDTO import の不在。
- `related detail requirement type`: `consistency_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: 全ページの props 境界と DTO 境界を確認するシナリオへ統合できる。
- `conflict hint`: 既存 component が gateway-contract の domain type を直接 props に使う場合、採否時に「許容する型」と「DTO 依存」の線引きが必要になる。

### CAND-APC-EI-009 Storybook review URL は Storybook 表示先だけを指す

- `source requirement`: Storybook 基盤と Master Persona の固定要件は、Storybook review URL、story ID、確認状態、未確認理由、再実行 command を task-local `storybook-review.md` に残すことである。review URL は Storybook 表示確認先であり、fakeAPI URL、Wails runtime URL、backend API URL は扱わない。
- `viewpoint`: external-integration / network 境界
- `candidate scenario id`: `CAND-APC-EI-009`
- `actor`: 実装後確認者
- `external boundary`: Storybook dev server、Storybook static preview、local browser URL
- `start condition`: 全ページの主要 story 実装後に Storybook review 証跡を残す。
- `trigger`: Storybook review URL と story ID を task-local 成果物へ記録する。
- `expected outcome`: review URL は Storybook localhost または iframe URL だけを指し、Wails runtime URL、backend API URL、fake provider URL、secret、API key、token、ローカル絶対 path を含まない。
- `fake_or_stub`: Storybook story の fixed props、view model fixture。
- `observable point`: `storybook-review.md` の URL、story ID、確認状態、未確認理由、build-storybook の結果。
- `related detail requirement type`: `observability_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: Storybook review 証跡シナリオへ統合できる。
- `conflict hint`: operation-audit 観点の証跡候補と重なる可能性がある。外部連携観点では URL が外部実行先を指さないことだけを扱う。

## Open Notes

- `human decision candidate`: page 合成 story を全ページで必須にするか、主要 panel、card、modal の story を必須にして page story は任意にするかは、designer が最終シナリオ統合時に判断する。
- `human decision candidate`: 既存 component が gateway-contract の domain type を props に使う場合、DTO 依存として禁止する範囲と、画面 view model の型として許容する範囲を designer が線引きする必要がある。
- `merge candidate`: `CAND-APC-EI-001`、`CAND-APC-EI-002`、`CAND-APC-EI-008` は、controller 接続、generated binding、DTO 境界をまとめた props 境界シナリオへ統合できる。
- `merge candidate`: `CAND-APC-EI-003`、`CAND-APC-EI-006`、`CAND-APC-EI-009` は、Storybook fixture と review 証跡の保存禁止情報シナリオへ統合できる。
- `rejection candidate`: 実 AI API、secret store、DB、backend Controller、Wails binding の実接続成立を証明する候補は、この task の external-integration 候補から除外する。理由は、全ページ部品化と Storybook story 追加の成果物境界を越えるためである。
