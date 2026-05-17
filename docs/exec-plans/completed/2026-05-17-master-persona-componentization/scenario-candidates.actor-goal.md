# Scenario Candidates: 2026-05-17-master-persona-componentization / actor-goal

- `generator`: `actor-goal`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `MPC`

## Generator Scope

- `viewpoint`: マスターペルソナ生成画面を使う利用者、後続 UI 実装者、人間レビュアーの目的、開始操作、成功判定を候補化する。
- `included_sources`: `plan.md`, `docs/index.md`, `docs/architecture.md`, `docs/spec.md`, `docs/er.md`, `docs/coding-guidelines-frontend.md`, `docs/coding-guidelines-tests.md`, `docs/UX-standard.md`, `docs/screen-design/screens/master-persona.md`, `docs/detail-specs/ai-provider-settings-management.md`, `docs/scenario-tests/dashboard-and-app-shell.md`, `../completed/2026-05-18-storybook-foundation/implementation-scope.md`, `../completed/2026-05-18-storybook-foundation/storybook-review.md`, `frontend/src/ui/screens/master-persona/`
- `excluded_sources`: プロダクトコード変更指示、プロダクトテスト変更指示、docs 正本化、最終シナリオ採否、候補統合、競合解消、外部連携失敗、状態遷移網羅
- `generation_notes`: 候補は actor-goal 観点に限定する。Storybook 基盤は完了済みの前提として扱い、Master Persona 固有の部品境界と story 確認対象を候補化する。

## Candidate Scenarios

### CAND-MPC-001 利用者が生成準備を確認して作成開始を判断する

- `source requirement`: `plan.md` の `goal` と `close_conditions`。`docs/screen-design/screens/master-persona.md` の生成準備パネル、AI 設定カード、入力 JSON パネル。`docs/spec.md` の共通ペルソナ構築と目的に沿った AI 選択。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-MPC-001`
- `actor`: マスターペルソナを作成する利用者
- `trigger`: 利用者がマスターペルソナ画面を開き、AI サービス、モデル、実行方法、入力 JSON を確認する。
- `expected outcome`: 利用者は AI 設定状態、選択中 JSON、候補数、新規作成数、既存スキップ数を見て、作成開始できるか判断できる。
- `observable point`: 生成準備パネルに AI 設定、入力 JSON、対象件数、`ペルソナを作成` の有効状態が表示される。Storybook では固定 props または view model fixture で同じ判断材料が表示される。
- `related detail requirement type`: `success_requirement`, `data_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: 生成準備パネルと AI モデル選択カードの主正常系として採用候補にできる。
- `conflict hint`: `GenerationSetupPanel` が画面全体の view model を受け取る案と、小さい props 境界へ分ける案で実装境界が競合しうる。

### CAND-MPC-002 利用者が生成中の進行状況と停止操作を判断する

- `source requirement`: `docs/screen-design/screens/master-persona.md` の進行状況パネル。`docs/UX-standard.md` の状態表示、状態別操作、ローディング状態。`plan.md` の Storybook 上で主要状態を確認できる close condition。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-MPC-002`
- `actor`: マスターペルソナ作成を監視する利用者
- `trigger`: 利用者がペルソナ作成の開始後に進行状況パネルを見る。
- `expected outcome`: 利用者は生成中であること、処理済み件数、作成済み件数、既存スキップ件数、現在の対象、一時停止と中止の可否を判断できる。
- `observable point`: 進行状況パネルに進捗バー、件数、現在の対象、`一時停止`、`中止` が表示される。生成中 fixture では停止操作が有効であり、非生成中 fixture では停止操作が無効である。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `testability_requirement`
- `adoption hint`: 生成中状態の Storybook POC 対象として採用候補にできる。
- `conflict hint`: 生成中の停止操作そのものの状態遷移や失敗扱いは state-transition 観点または failure 観点と統合時に重複しうる。

### CAND-MPC-003 利用者が生成結果を検索して詳細を確認する

- `source requirement`: `docs/screen-design/screens/master-persona.md` の生成結果一覧パネルとペルソナ詳細パネル。`docs/spec.md` の翻訳補助メタデータを UI から観測可能にする要件。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-MPC-003`
- `actor`: 作成済みペルソナを確認する利用者
- `trigger`: 利用者が生成結果一覧で名前またはプラグイン名を検索し、ペルソナ行を選択する。
- `expected outcome`: 利用者は件数範囲、ページ、対象プラグイン、FormID、EditorID、声、話し方、ペルソナ本文を確認できる。
- `observable point`: 一覧パネルに検索欄、プラグイン選択、ページ操作、ペルソナ行が表示される。詳細パネルに選択中ペルソナの識別情報と本文が表示される。
- `related detail requirement type`: `success_requirement`, `data_requirement`, `testability_requirement`
- `adoption hint`: 一覧と詳細の表示部品を分ける根拠として採用候補にできる。
- `conflict hint`: 一覧と詳細を 1 つの story にまとめる案と、一覧パネル story と詳細パネル story に分ける案で review 粒度が競合しうる。

### CAND-MPC-004 利用者が選択中ペルソナを編集または削除する

- `source requirement`: `docs/screen-design/screens/master-persona.md` の編集モーダルと削除モーダル。`docs/UX-standard.md` の確認構造、危険操作、モーダル使用。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-MPC-004`
- `actor`: 作成済みペルソナを管理する利用者
- `trigger`: 利用者が詳細パネルで選択中ペルソナを確認し、`編集` または `削除` を押す。
- `expected outcome`: 利用者は編集モーダルで要約、話し方、本文を更新できる。利用者は削除モーダルで識別情報を確認して削除判断ができる。
- `observable point`: 編集モーダルに入力欄と `編集内容を保存` が表示される。削除モーダルに名前、FormID、EditorID、対象プラグイン、`削除する` が表示される。
- `related detail requirement type`: `success_requirement`, `data_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: モーダルを screen local component として story 化する候補にできる。
- `conflict hint`: 削除後の永続化、復元可否、監査保存は actor-goal 観点だけでは確定できず、他観点または人間判断と競合しうる。

### CAND-MPC-005 後続 UI 実装者が controller 接続と表示部品を分離して理解する

- `source requirement`: `plan.md` の `goal`。`docs/architecture.md` の `View` と `UI Component` 責務。`docs/coding-guidelines-frontend.md` の Svelte ファイル分割と component 分割基準。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-MPC-005`
- `actor`: 後続 UI 実装者
- `trigger`: 後続 UI 実装者が `MasterPersonaPage` と screen local component を読んで、変更対象の表示部品を探す。
- `expected outcome`: 後続 UI 実装者は `MasterPersonaPage` を controller 接続と部品合成の入口として読み、生成準備、進行状況、レビュー、操作モーダルの表示責務を個別部品として理解できる。
- `observable point`: screen local component は controller、Store、Gateway、generated binding を直接扱わず、props と callback で表示と操作を表す。Storybook では controller を作らずに主要部品を表示できる。
- `related detail requirement type`: `compatibility_requirement`, `testability_requirement`, `success_requirement`
- `adoption hint`: 部品化 refactor の責務境界シナリオとして採用候補にできる。
- `conflict hint`: 既存の画面専用部品を残す範囲と、新規に切り出す範囲は `plan.md` の未決事項であり、designer 統合時に人間判断候補になりうる。

### CAND-MPC-006 Story 作成者が固定 fixture で主要表示状態を追加する

- `source requirement`: `plan.md` の `goal` と `close_conditions`。完了済み `2026-05-18-storybook-foundation/implementation-scope.md` の fixture 配置方針と backend 境界禁止。完了済み `storybook-review.md` の review URL 記録方針。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-MPC-006`
- `actor`: Storybook story 作成者
- `trigger`: story 作成者がマスターペルソナ画面のパネル、カード、モーダルへ固定 props または view model fixture を渡す story を追加する。
- `expected outcome`: story 作成者は backend DTO mock、gateway mock、翻訳実行フロー再現を作らず、未設定、生成中、生成成功、生成失敗、編集中の主要表示状態を Storybook 上で確認できる。
- `observable point`: story source と fixture import path が component 横の `__fixtures__` から確認できる。Storybook 上で対象パネル、カード、モーダルの主要状態が表示される。
- `related detail requirement type`: `testability_requirement`, `data_requirement`, `compatibility_requirement`
- `adoption hint`: Storybook POC の中心候補として採用候補にできる。
- `conflict hint`: 主要表示状態をどこまで story に含めるかは `plan.md` の未決事項であり、failure 観点や state-transition 観点の候補と粒度が競合しうる。

### CAND-MPC-007 人間レビュアーが Storybook で部品の見た目を確認する

- `source requirement`: `plan.md` の `constraints`、`close_conditions`、`HITL Status`。完了済み `storybook-review.md` の review URL と確認状態の記録方法。`docs/UX-standard.md` の証跡、状態別スクリーンショット、検査可能性。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-MPC-007`
- `actor`: 人間レビュアー
- `trigger`: 人間レビュアーが task-local 成果物に記録された Storybook review URL を開く。
- `expected outcome`: 人間レビュアーはマスターペルソナ画面の主要パネル、カード、モーダルを Storybook 上で確認し、実装後の見た目レビュー対象を判断できる。
- `observable point`: task-local 成果物に review URL、story ID、確認状態、未確認理由が残る。Storybook 表示では対象 story が secret、API key、token、実ユーザーデータなしで表示される。
- `related detail requirement type`: `testability_requirement`, `observability_requirement`, `compatibility_requirement`
- `adoption hint`: frontend human review へつなぐ運用シナリオとして採用候補にできる。
- `conflict hint`: `plan.md` は review 記録先を `ui-design.md` または `storybook-review.md` の未決事項として残している。一方、Storybook 基盤成果物は専用 `storybook-review.md` に固定しているため、designer 統合時に扱いを確認する必要がある。

### CAND-MPC-008 実装検証者が Storybook と既存 frontend gate で回帰を確認する

- `source requirement`: `plan.md` の `validation_commands` と `close_conditions`。`docs/coding-guidelines-tests.md` の frontend local gate。完了済み Storybook 基盤の build gate。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-MPC-008`
- `actor`: 実装検証者
- `trigger`: 実装検証者がマスターペルソナ部品化と story 追加後に Storybook build と frontend local gate を確認する。
- `expected outcome`: 実装検証者は Storybook の表示確認対象と、既存アプリ側の frontend 回帰がどちらも成立しているか判断できる。
- `observable point`: Storybook build の結果、frontend local gate の結果、未通過理由が task-local 成果物へ残る。Storybook story は backend、Wails runtime、AI provider、secret store、DB を起動条件にしない。
- `related detail requirement type`: `testability_requirement`, `compatibility_requirement`
- `adoption hint`: POC 完了条件と既存 frontend 回帰確認をつなぐ候補にできる。
- `conflict hint`: 検証段階を `実装後` に寄せるか `最終検証` に寄せるかは、最終シナリオ統合時に他観点と調整が必要である。

## Open Notes

- `human decision candidate`: 既存部品を残す範囲と、新規部品として切り出す範囲。
- `human decision candidate`: `MasterPersonaPage` の page story を本体 component で作るか、review-only wrapper で作るか。
- `human decision candidate`: `GenerationSetupPanel`、`RunStatusPanel`、`PersonaReviewPanel`、`PersonaActionModal` が受け取る props を画面全体の view model のままにするか、小さい props 境界へ分けるか。
- `human decision candidate`: Storybook の必須表示状態を、未設定、生成中、生成成功、生成失敗、編集中の全件にするか、POC として代表状態へ絞るか。
- `merge candidate`: `CAND-MPC-001` と `CAND-MPC-006` は、生成準備パネルの固定 fixture story として統合できる可能性がある。
- `merge candidate`: `CAND-MPC-003` と `CAND-MPC-004` は、レビュー領域と操作モーダルの人間確認シナリオとして統合できる可能性がある。
- `merge candidate`: `CAND-MPC-007` と `CAND-MPC-008` は、Storybook POC の review 証跡と検証証跡として統合できる可能性がある。
- `rejection candidate`: Storybook に backend DTO mock、gateway mock、generated `wailsjs`、AI provider、secret store、DB、filesystem business flow を入れる候補は、完了済み Storybook 基盤の境界により不採用候補とする。
