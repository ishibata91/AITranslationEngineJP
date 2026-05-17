# Scenario Design: 2026-05-18-storybook-foundation

- `skill`: scenario-design
- `status`: human-reviewed
- `human_review`: answered
- `source_plan`: `./plan.md`
- `ui_source`: `N/A`
- `screen_design_diff`: `N/A`
- `final_artifact_path`: `docs/scenario-tests/storybook-foundation.md`
- `topic_abbrev`: `SBF`
- `candidate_sources`:
  - `./scenario-candidates.actor-goal.md`
  - `./scenario-candidates.lifecycle.md`
  - `./scenario-candidates.state-transition.md`
  - `./scenario-candidates.failure.md`
  - `./scenario-candidates.external-integration.md`
  - `./scenario-candidates.operation-audit.md`

## Fixed Requirements

- `must_pass_requirements`:
  - Storybook は Master Persona 部品化より前に、最小基盤だけを追加する。
  - Storybook は `frontend/` を package root とし、`Svelte 5`、`TypeScript`、`Vite` の前提に従う。
  - Storybook dev 入口は、空または最小サンプル story を表示できる。
  - Storybook static build は、backend、Wails runtime、AI provider、secret store、DB に依存しない。
  - Storybook story は固定 props または view model fixture で表示し、gateway mock、backend DTO mock、実行フロー再現を含めない。
  - Storybook の import alias は既存 `frontend/vite.config.ts` の `@ui`、`@application`、`@controller` と矛盾しない。
  - 既存 lint は、プロダクトコードから Storybook への依存混入がないことを検査できる。
  - Storybook review URL、対象 story、確認状態、未確認理由は task-local 成果物へ残す。
  - 後続 Master Persona task は、最小 story と fixture 方針を起点に story を追加できる。
- `non_goals`:
  - Master Persona の部品化は含めない。
  - 画面再設計は含めない。
  - gateway mock、backend DTO mock、翻訳実行フロー再現は含めない。
  - backend 連携、AI provider 通信、secret 保存、DB 永続化、xTranslator 出力の成立証明は含めない。
  - プロダクトコード、プロダクトテスト、docs 正本、`.codex/` の変更指示は含めない。
  - `ui-design.md` と `screen-design-diff.*.md` はこの task では扱わない。

## Candidate Coverage

正本: `./scenario-design.candidate-coverage.json`

6 件の候補成果物は揃っている。
候補網羅 JSON では、生成 agent 名を付けた `generator:CAND-...` を一意な識別子として扱う。

`needs_human_decision` は 0 件である。
未解決の競合は 0 件である。
人間回答は `./scenario-design.questions.md` に記録済みである。

## Requirement Coverage

正本: `./scenario-design.requirement-coverage.json`

各抽象要件の詳細要求タイプは別 JSON に分離する。
この成果物は人間回答を反映済みである。
`scenario-design.requirement-coverage.json` は gate 通過用に回答済み状態へ更新済みである。

### `REQ-SBF-001` Storybook 最小基盤を起動できる

- `source_requirement`: `plan.md` の `goal` と `close_conditions`。Storybook dev 入口が動き、空または最小サンプル story で基盤確認できる。
- `requirement_kind`: frontend_tooling_foundation
- `needs_human_decision`: なし
- `fixed_decisions`: Storybook dev 入口は `frontend/` package root から起動する。表示対象は空または最小サンプル story に限定する。backend と Wails runtime は起動条件にしない。

### `REQ-SBF-002` Storybook static build を検証できる

- `source_requirement`: `plan.md` の `validation_commands`。`npm --prefix frontend run build-storybook` が検証入口に含まれる。
- `requirement_kind`: frontend_tooling_foundation
- `needs_human_decision`: なし
- `fixed_decisions`: Storybook build command 自体は必須検証入口である。Storybook build は Storybook 専用 gate に分ける。理由は、厳しく harness で制限する必要を感じないためである。

### `REQ-SBF-003` 後続 story が fixture 方針に沿って追加できる

- `source_requirement`: `plan.md` の `goal` と人間回答。fixture 配置方針を固定し、後続 Master Persona task が story を追加できる。
- `requirement_kind`: data_boundary
- `needs_human_decision`: なし
- `fixed_decisions`: fixture は fixed props または view model fixture とする。gateway mock、backend DTO mock、実行フロー再現は使わない。Storybook fixture は component 横の `frontend/src/ui/**/__fixtures__` に置く。

### `REQ-SBF-004` Storybook review URL を task-local に記録できる

- `source_requirement`: `plan.md` の `goal`。review URL 記録方針を最小単位で固定する。
- `requirement_kind`: operation_record
- `needs_human_decision`: なし
- `fixed_decisions`: review URL は Storybook 表示確認先だけを指す。fakeAPI URL、Wails runtime URL、backend API URL は Storybook review URL として扱わない。Storybook review URL と確認状態は専用 `storybook-review.md` に記録する。

### `REQ-SBF-005` Storybook は backend 境界を越えない

- `source_requirement`: `plan.md` の `constraints`。gateway mock、backend DTO mock、実行フロー再現を Storybook に入れない。
- `requirement_kind`: boundary
- `needs_human_decision`: なし
- `fixed_decisions`: Storybook は UI 表示確認に限定する。story と fixture は generated `wailsjs`、backend DTO、Gateway、RuntimeEventAdapter、AI provider、secret store、DB、filesystem business flow を import しない。

### `REQ-SBF-006` 検証結果と後続引き継ぎを追跡できる

- `source_requirement`: `plan.md` の `close_conditions` と `Routing Notes`。Storybook dev / build の入口、fixture 配置方針、review URL 記録方針を後続 task へ渡す。
- `requirement_kind`: operation_record
- `needs_human_decision`: なし
- `fixed_decisions`: task-local 成果物には command、結果、未実行理由、review URL、story ID、残留リスクを残す。Storybook 運用は、後続 task の plan と POC task が成功した後に、skill、agent、docs へ反映する。

### `REQ-SBF-007` 既存 lint で Storybook 依存混入を検査できる

- `source_requirement`: 人間追加 goal。既存 lint に、Storybook への依存がないことをチェックするものを追加する。
- `requirement_kind`: lint_boundary
- `needs_human_decision`: なし
- `fixed_decisions`: 既存 lint は、プロダクトコードが Storybook package、Storybook runtime、Storybook 専用 module へ依存していないことを検査する。Storybook 専用設定、story、fixture は検査対象の許可範囲として扱う。

## Human Decision Questionnaire

正本: `./scenario-design.questions.md`

未回答質問は 0 件である。
`Q-SBF-001` から `Q-SBF-004` までの人間回答を反映済みである。

## Risks

- Storybook build を Storybook 専用 gate に分けるため、通常 `frontend-local` とは別に実行記録が必要になる。
- Storybook 依存混入チェックを既存 lint に追加するため、Storybook 専用ファイルとプロダクトコードの対象範囲を誤ると、許可すべき story まで止める可能性がある。
- Storybook fixture を component 横に置くため、UI 階層の移動時は fixture も同じ境界で移動する必要がある。
- Storybook review URL と確認状態を専用 `storybook-review.md` に記録するため、implementation closeout で artifact 作成を忘れると人間レビュー入口が欠落する可能性がある。
- Storybook 運用の skill、agent、docs 反映は後続 task の plan と POC task 成功後に回るため、この task では docs 正本を変更しない。
- Storybook story に backend DTO mock や gateway mock を入れると、Storybook が backend 連携検証の代替に見える可能性がある。

## Rules

- ケース ID は `SCN-SBF-NNN` 形式にする。
- Markdown table は使わず、1 ケースごとの縦型ブロックで書く。
- 受け入れテストは全ケースで先に固定する。
- `実行テスト種別` は `APIテスト | UI人間操作E2E | lower-level only` に固定する。
- `実行段階` は `実装後 | 最終検証` に固定する。
- `期待結果` は観測可能な結果にする。
- `needs_human_decision` が残る場合はシナリオ完了にしない。
- 未解決競合が残る場合はシナリオ完了にしない。
- 有料の実 AI API を前提にしない。

## Scenario Matrix

このシナリオ表は人間回答を反映済みである。
implementation-scope 作成判断は `implement_lane` に戻す。

### SCN-SBF-001 Storybook dev 入口で最小 story を確認する

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: 後続 UI 実装者が Storybook dev 入口から最小 story を確認できる。
- `受け入れ条件`: Storybook dev script を実行すると Storybook が起動し、空または最小サンプル story を表示できる。
- `事前条件`: Storybook script、最小 config、最小 story が存在する。
- `public_seam_or_api_boundary`: Storybook dev server と browser のローカル表示境界。
- `入力開始点`: Storybook dev script。
- `主要 outcome`: Storybook URL で最小 story が見える。
- `開始操作`: Storybook dev script を実行する。
- `入力方法`: `frontend/` package root の npm script を使う。
- `主要操作列`: dev script を起動し、Storybook URL を開き、最小 story を選択する。
- `期待結果`:
  1. dev command が Storybook server を起動する。
  2. Storybook 画面で最小 story が選択できる。
  3. Wails runtime、backend Controller、AI provider は起動条件にならない。
- `観測点`: dev command の起動状態、Storybook URL、browser 上の最小 story 表示。
- `UI-visible outcome`: 最小 story が Storybook canvas または docs view に表示される。
- `fake_or_stub`: 固定 props または最小表示用 fixture。gateway mock と backend DTO mock は使わない。

### SCN-SBF-002 Storybook static build を検証する

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: Storybook の静的 build が backend なしで成立する。
- `受け入れ条件`: `npm --prefix frontend run build-storybook` が成功し、Storybook static output が生成される。
- `事前条件`: Storybook script、最小 config、最小 story が存在する。
- `public_seam_or_api_boundary`: npm script と Storybook static build output。
- `入力開始点`: `npm --prefix frontend run build-storybook`
- `主要 outcome`: Storybook static build が成功する。
- `期待結果`:
  1. build command が成功する。
  2. Storybook の生成物は frontend tooling の成果物として作られる。
  3. Wails build、backend 実行、AI provider 通信、secret store、DB は不要である。
  4. 既存 app build の `dist` と Storybook output を混同しない。
- `観測点`: command exit code、build output directory、backend 接続を前提にしない static asset 構成。
- `公開接点確認`: あり。npm script を公開接点として扱う。
- `決定`: `Q-SBF-001`。Storybook build は Storybook 専用 gate に分ける。理由は、厳しく harness で制限する必要を感じないためである。

### SCN-SBF-003 後続 task が story を追加できる

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: 後続 Master Persona task が Storybook 基盤を起点に story を追加できる。
- `受け入れ条件`: 後続 task が、最小 story と fixture 方針を見て、backend DTO mock や gateway mock を作らずに story を追加できる。
- `事前条件`: Storybook dev と build の入口が確認済みである。
- `public_seam_or_api_boundary`: UI Component props と story fixture の境界。
- `入力開始点`: 最小 story と fixture 方針。
- `主要 outcome`: story 追加先、fixture 参照先、props 入力境界が分かる。
- `期待結果`:
  1. story は UI Component または screen local component に固定 props を渡す。
  2. story は generated `wailsjs`、backend DTO、Gateway、RuntimeEventAdapter を import しない。
  3. fixture は secret、API key、token、外部 provider 応答原文、実ユーザーデータを含まない。
  4. 後続 Master Persona 固有の部品境界はこの task で確定しない。
- `観測点`: story source、fixture import path、禁止 import の有無、task-local fixture 方針。
- `公開接点確認`: あり。story source と fixture root を後続 task の入力境界として扱う。
- `決定`: `Q-SBF-002`。Storybook fixture は component 横の `frontend/src/ui/**/__fixtures__` に置く。理由は、component 単位の見た目検証用であり、業務データを Storybook で扱わないためである。

### SCN-SBF-004 Storybook review URL と確認状態を記録する

- `分類`: 運用記録
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `最終検証`
- `観点`: 人間レビュー担当者が同じ Storybook 表示へ到達できる。
- `受け入れ条件`: task-local 成果物に Storybook URL、story ID、確認状態、未確認理由、起動 command が残る。
- `事前条件`: Storybook dev server が起動しているか、Storybook static preview URL がある。
- `public_seam_or_api_boundary`: local browser と Storybook server または Storybook static preview の確認境界。
- `入力開始点`: Storybook review URL。
- `主要 outcome`: 人間レビュー対象の URL と確認状態が追跡できる。
- `開始操作`: Storybook URL を開く。
- `入力方法`: task-local に記録された URL と story ID を使う。
- `主要操作列`: Storybook URL を開き、対象 story を確認し、確認状態または未確認理由を task-local 成果物へ残す。
- `期待結果`:
  1. review URL は Storybook 表示確認先を指す。
  2. fakeAPI URL、Wails runtime URL、backend API URL を Storybook review URL として扱わない。
  3. URL query に secret、API key、token、ローカル絶対 path、実ユーザーデータを含めない。
- `観測点`: task-local 成果物の URL、story ID、確認状態、未確認理由。
- `UI-visible outcome`: 記録された Storybook URL で対象 story が表示される。
- `fake_or_stub`: Storybook story の固定 props。
- `決定`: `Q-SBF-003`。Storybook review URL と確認状態は、専用 `storybook-review.md` に記録する。

### SCN-SBF-005 禁止された backend mock と実行フロー再現を入れない

- `分類`: 禁止遷移
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: Storybook 基盤が backend 連携検証の代替にならない。
- `受け入れ条件`: story と fixture に gateway mock、backend DTO mock、generated `wailsjs`、翻訳実行フロー再現が入らない。
- `事前条件`: Storybook story または fixture が追加されている。
- `public_seam_or_api_boundary`: story source と import 境界。
- `入力開始点`: Storybook story source と fixture source。
- `主要 outcome`: Storybook は UI 表示確認に閉じる。
- `期待結果`:
  1. story と fixture は Gateway、generated `wailsjs`、backend DTO、RuntimeEventAdapter を import しない。
  2. story と fixture は AI provider、secret store、DB、filesystem business flow を再現しない。
  3. scope 逸脱が見つかった場合は Storybook 基盤の成果物に含めない。
- `観測点`: import 境界、story props、fixture の種類、backend DTO 参照の有無。
- `公開接点確認`: なし。source 境界確認として扱う。
- `fake_or_stub`: props fixture、view model fixture。gateway mock と backend DTO mock は使わない。

### SCN-SBF-006 検証結果と後続引き継ぎを残す

- `分類`: 運用記録
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `最終検証`
- `観点`: Storybook 基盤 task の完了証跡を後続 Master Persona task が再利用できる。
- `受け入れ条件`: task-local 成果物に script 名、config path、最小 story path、build 検証結果、review URL 記録形式、残留リスクが残る。
- `事前条件`: Storybook dev、Storybook build、fixture 方針、review URL 記録方針の確認が終わっている。
- `public_seam_or_api_boundary`: task-local closeout 証跡。
- `入力開始点`: 実装結果 artifact または closeout 証跡。
- `主要 outcome`: 後続 task が追加 story の検証入口を特定できる。
- `期待結果`:
  1. 実行した command、通過または失敗、未実行理由、失敗時の要約が残る。
  2. command 出力全文、依存 cache の絶対 path、secret、token、長いローカル path は保存しない。
  3. 後続 task が Storybook scripts、config、最小 story、review URL 記録形式を参照できる。
  4. docs 正本化が未承認の場合、docs 正本は変更しない。
- `観測点`: task-local 成果物の command 結果、review URL 記録形式、残留リスク。
- `公開接点確認`: あり。task-local closeout 証跡を後続 task の入力として扱う。
- `決定`: `Q-SBF-004`。Storybook 運用は、後続 task の plan と POC task が成功した後に、skill、agent、docs へ反映する。

### SCN-SBF-007 既存 lint で Storybook 依存混入を検査する

- `分類`: 依存境界
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: プロダクトコードが Storybook 実行環境へ依存しない。
- `受け入れ条件`: 既存 lint を実行すると、Storybook 専用設定、story、fixture 以外から Storybook package、Storybook runtime、Storybook 専用 module への依存混入を検出できる。
- `事前条件`: Storybook package、設定、最小 story、fixture が追加済みである。
- `public_seam_or_api_boundary`: frontend lint の依存境界チェック。
- `入力開始点`: `npm --prefix frontend run lint`
- `主要 outcome`: 既存 lint が Storybook 依存混入を検査する。
- `期待結果`:
  1. Storybook 専用設定、story、fixture は Storybook 依存を持てる。
  2. プロダクトコードは Storybook package、Storybook runtime、Storybook 専用 module を import しない。
  3. 依存混入がある場合は lint が失敗する。
- `観測点`: lint command の exit code、依存境界チェックの対象範囲、許可対象 file pattern、禁止 import pattern。
- `公開接点確認`: あり。既存 lint を公開検査入口として扱う。

## Validation Entry

- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/2026-05-18-storybook-foundation/scenario-design.md --coverage docs/exec-plans/active/2026-05-18-storybook-foundation/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/active/2026-05-18-storybook-foundation/scenario-design.candidate-coverage.json --json`
- `python3 scripts/harness/run.py --suite scenario-gate`
- `python3 scripts/harness/run.py --suite frontend-local`
- `npm --prefix frontend run lint`
- `npm --prefix frontend run build-storybook`

## Handoff State

- `status`: human-reviewed
- `reason`: `Q-SBF-001` から `Q-SBF-004` までの人間回答を反映済みである。
- `return_to`: `implement_lane`
