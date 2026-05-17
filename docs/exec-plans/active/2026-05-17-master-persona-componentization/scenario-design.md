# Scenario Design: 2026-05-17-master-persona-componentization

- `skill`: scenario-design
- `status`: human-review-ready
- `source_plan`: `./plan.md`
- `ui_source`: `N/A`
- `final_artifact_path`: `N/A`
- `topic_abbrev`: `MPC`
- `candidate_sources`:
  - `./scenario-candidates.actor-goal.md`
  - `./scenario-candidates.lifecycle.md`
  - `./scenario-candidates.state-transition.md`
  - `./scenario-candidates.failure.md`
  - `./scenario-candidates.external-integration.md`
  - `./scenario-candidates.operation-audit.md`

## Fixed Requirements

- `must_pass_requirements`:
  - `MasterPersonaPage` は production controller 接続と部品合成へ寄せる。
  - 生成準備、進行状況、生成結果確認、操作モーダルは、画面 view model 由来の小さい props と callback で表示できる。
  - `AIModelSelectionCard` は既存共有 component を再利用し、手作りの代替カードを作らない。
  - Storybook story と fixture は backend、Wails runtime、generated `wailsjs`、Gateway、RuntimeEventAdapter、AI provider、secret store、DB、実 filesystem flow を要求しない。
  - Storybook fixture、story、review 記録は secret、API key、token、実 endpoint、ローカル絶対 path、実ユーザーデータ、raw request、raw response、raw prompt を含めない。
  - Storybook review URL、story ID、確認状態、未確認理由、再実行 command は task-local `storybook-review.md` に残す。
  - 実装後検証は `npm --prefix frontend run build-storybook` と `python3 scripts/harness/run.py --suite frontend-local` を入口にする。
- `non_goals`:
  - プロダクト backend、Wails binding、Gateway、RuntimeEventAdapter、AI provider、secret store、DB の実装変更。
  - docs 正本、`.codex`、既存 scenario candidate、`plan.md` の変更。
  - Storybook 上で実 AI 生成、実ファイル読み込み、DB 書き込み、provider network を再現すること。
  - UI 設計本文または画面設計書正本の作成。

## Scenario Candidate Coverage

正本: `./scenario-design.candidate-coverage.json`

6 種の候補成果物を統合した。
未解決競合と `needs_human_decision` は 0 件である。
候補で出ていた review 記録先の揺れは、完了済み Storybook 基盤の `storybook-review.md` 方針へ統合した。
候補で出ていた page story の扱いは、主要表示領域 story を必須、page 合成 story を任意として固定した。

## Detail Requirement Coverage

正本: `./scenario-design.requirement-coverage.json`

抽象要件は、部品境界、生成準備、実行状態、一覧詳細とモーダル、Storybook 外部境界、review 証跡と検証入口へ分けた。
`needs_human_decision` は 0 件である。
`deferred` は UI 設計または implementation-scope が扱う実装分割粒度だけに限定した。

## Human Decision Questionnaire

正本: `./scenario-design.questions.md`

none

## Scenario Matrix

| scenario_id | 受け入れシナリオ | 受け入れ条件 | 実行テスト種別 | 実行段階 | 主な候補 |
| --- | --- | --- | --- | --- | --- |
| `SCN-MPC-001` | controller 接続と表示部品の props 境界を分離する | `MasterPersonaPage` は production controller を作り、screen local component は `Store`、`Gateway`、generated binding を直接扱わない。各 component は props と callback だけで story 表示できる。 | `lower-level only` | `実装後` | `CAND-MPC-005`, `CAND-MPC-LC-001`, `CAND-MPCF-001`, `CAND-MPC-EI-001`, `CAND-MPCOA-004` |
| `SCN-MPC-002` | 生成準備、AI 設定、入力 JSON の状態を story で確認する | 未選択、JSON 選択済み、preview 成功、AI 設定不足、モデル一覧更新中、長い model 名を固定 fixture で表示できる。`ペルソナを作成` は作成開始可能状態だけ有効になる。 | `lower-level only` | `実装後` | `CAND-MPC-001`, `CAND-MPC-LC-002`, `CAND-MPSC-001`, `CAND-MPSC-002`, `CAND-MPCF-003`, `CAND-MPCF-004`, `CAND-MPC-EI-002`, `CAND-MPC-EI-004`, `CAND-MPC-EI-007` |
| `SCN-MPC-003` | 生成中と非生成中の進行状態を story で確認する | 生成中 fixture は進捗、処理済み件数、作成済み件数、既存スキップ数、現在の対象を表示する。生成中は一時停止と中止だけが有効になり、編集、削除、再生成は無効になる。 | `lower-level only` | `実装後` | `CAND-MPC-002`, `CAND-MPC-LC-003`, `CAND-MPC-LC-007`, `CAND-MPSC-003`, `CAND-MPSC-004`, `CAND-MPCF-005`, `CAND-MPCF-006`, `CAND-MPCOA-005` |
| `SCN-MPC-004` | 生成結果一覧と詳細の選択状態を story で確認する | 空一覧、一覧あり、行選択済み、検索またはページ変更後の未選択状態を固定 fixture で表示できる。古い選択が一覧から消えた場合、詳細は古いペルソナを表示しない。 | `lower-level only` | `実装後` | `CAND-MPC-003`, `CAND-MPC-LC-004`, `CAND-MPSC-005`, `CAND-MPCF-007` |
| `SCN-MPC-005` | 編集と削除のモーダル状態を story で確認する | 編集モーダルは選択中 identity、要約、話し方、本文、保存操作を表示する。削除モーダルは識別情報と危険操作を表示する。保存失敗または削除失敗 fixture ではモーダルを閉じず、入力値または対象識別情報を残す。 | `lower-level only` | `実装後` | `CAND-MPC-004`, `CAND-MPC-LC-005`, `CAND-MPC-LC-006`, `CAND-MPSC-006`, `CAND-MPSC-007`, `CAND-MPCF-008`, `CAND-MPCF-009` |
| `SCN-MPC-006` | Storybook fixture が外部接続と保存禁止情報を持たない | story と fixture は fixed props、view model fixture、callback stub だけで表示できる。button callback は AI 実行、Wails binding、provider network、DB 書き込み、実ファイル読み込みを開始しない。fixture は synthetic data だけを使う。 | `lower-level only` | `実装後` | `CAND-MPC-006`, `CAND-MPSC-008`, `CAND-MPCF-002`, `CAND-MPCF-010`, `CAND-MPCF-011`, `CAND-MPC-EI-003`, `CAND-MPC-EI-005`, `CAND-MPCOA-002`, `CAND-MPCOA-003`, `CAND-MPCOA-006` |
| `SCN-MPC-007` | Storybook review 証跡と検証結果を task-local に残す | `storybook-review.md` は story ID、review URL、確認状態、未確認理由、再実行 command、Storybook build 結果を持つ。review URL は Storybook localhost または iframe URL だけを指す。`frontend-local` と `build-storybook` は通過または未通過理由を持つ。 | `UI人間操作E2E` | `実装後` | `CAND-MPC-007`, `CAND-MPC-008`, `CAND-MPC-LC-008`, `CAND-MPCF-012`, `CAND-MPC-EI-006`, `CAND-MPCOA-001` |

## Acceptance Checks

### `SCN-MPC-001` controller 接続と表示部品の props 境界

- 受け入れ条件: screen local component は controller factory、Gateway、generated binding、RuntimeEventAdapter を import しない。
- 公開接点 / API 境界: Svelte component props と callback。
- 入力開始点: fixed props または view model fixture。
- 主要結果: component が production runtime なしで描画できる。
- 主要観測点: story source、fixture source、component import。
- 公開接点確認: あり。

### `SCN-MPC-002` 生成準備、AI 設定、入力 JSON

- 受け入れ条件: 作成開始可否、AI 設定状態、preview 件数、入力 JSON 状態が同じ fixture から一貫して表示される。
- 公開接点 / API 境界: `GenerationSetupPanel` props と `AIModelSelectionCard` props。
- 入力開始点: 未選択、JSON 選択済み、preview 成功、設定不足、更新中の fixture。
- 主要結果: 利用者は作成開始できるか判断できる。
- 主要観測点: `master-persona-generation-setup-panel`, `master-persona-ai-settings-card`, `master-persona-input-json-panel`, `executeGenerationButton`。
- 公開接点確認: あり。

### `SCN-MPC-003` 生成中と非生成中

- 受け入れ条件: `isRunActive` と `canMutate` の組み合わせが、進行状況、停止操作、編集削除可否と矛盾しない。
- 公開接点 / API 境界: `RunStatusPanel` props と `PersonaReviewPanel` props。
- 入力開始点: 生成中、生成失敗、生成完了の fixture。
- 主要結果: 生成中は再生成、編集、削除ができない。
- 主要観測点: `master-persona-progress-panel`, `interruptGenerationButton`, `cancelGenerationButton`, `editButton`, `deleteButton`。
- 公開接点確認: あり。

### `SCN-MPC-004` 生成結果一覧と詳細

- 受け入れ条件: 一覧、検索、ページ、選択中詳細、空状態が同じ props から整合して表示される。
- 公開接点 / API 境界: `PersonaReviewPanel` props と callback。
- 入力開始点: 空一覧、一覧あり、選択済み、絞り込み後空の fixture。
- 主要結果: 利用者は対象ペルソナを検索し、詳細を確認できる。
- 主要観測点: `master-persona-generation-result-list-panel`, `master-persona-persona-detail-panel`, `prevPageButton`, `nextPageButton`。
- 公開接点確認: あり。

### `SCN-MPC-005` 編集と削除のモーダル

- 受け入れ条件: 編集、削除、保存失敗、削除失敗の各状態で、対象識別情報と操作可否が失われない。
- 公開接点 / API 境界: `PersonaActionModal` props と callback。
- 入力開始点: `modalState`, `selectedEntry`, `editForm` の fixture。
- 主要結果: 利用者は編集または削除の対象を確認して操作できる。
- 主要観測点: `master-persona-edit-modal`, `master-persona-delete-modal`, `saveEntryButton`, `confirmDeleteButton`。
- 公開接点確認: あり。

### `SCN-MPC-006` Storybook fixture 外部境界

- 受け入れ条件: Storybook story は外部実行を開始せず、保存禁止情報を含まない。
- 公開接点 / API 境界: story、fixture、callback stub。
- 入力開始点: component 横の `__fixtures__`。
- 主要結果: Storybook build は backend、Wails runtime、AI provider、secret store、DB、実 filesystem なしで成立する。
- 主要観測点: story import、fixture import、Storybook build 結果、保存禁止情報の不在。
- 公開接点確認: あり。

### `SCN-MPC-007` Storybook review 証跡

- 開始操作: 人間レビュアーまたは実装後確認者が Storybook review URL を開く。
- 入力方法: `storybook-review.md` に記録された URL と story ID を使う。
- 主要操作列: Storybook を起動する。対象 story を開く。表示、story ID、未確認理由、build 結果を確認する。
- 主要観測点: Storybook URL、iframe URL、story ID、確認状態、未確認理由、`build-storybook` 結果。
- UI-visible 結果: 主要パネル、カード、モーダルの story が表示される。
- fake / stub 方針: fixed props、view model fixture、callback stub だけを使う。実 AI API、secret、DB、Wails runtime は使わない。

## Validation Entry

- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/2026-05-17-master-persona-componentization/scenario-design.md --report-out docs/exec-plans/active/2026-05-17-master-persona-componentization/scenario-design.requirement-gate.md --questionnaire-out docs/exec-plans/active/2026-05-17-master-persona-componentization/scenario-design.questions.md`
- `npm --prefix frontend run build-storybook`
- `python3 scripts/harness/run.py --suite frontend-local`

## Risks

- UI 設計本文はこの起動では作成しないため、visual review の詳細 story 順序は次の `ui-design.md` または `storybook-review.md` で固定する必要がある。
- page 合成 story は任意である。主要表示領域 story の不足を page story だけで代替してはならない。
- 現行 component は `viewModel` 全体を受け取っている箇所がある。実装では scenario の小さい props 境界を満たすため、props 型の再分割が必要になる可能性がある。

## Next Artifacts

- `ui-design.md`: Storybook review 対象、見た目確認順、確認状態の扱いを固定する。
- `screen-design-diff.master-persona.md`: 画面設計書正本へ反映する差分がある場合だけ作成する。
- `implementation-scope.md`: 人間レビュー後に、部品化、story fixture、検証、review 証跡の実装 handoff を分割する。
