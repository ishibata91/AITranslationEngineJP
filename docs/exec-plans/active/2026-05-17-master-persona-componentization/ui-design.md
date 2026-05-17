# UI Design: 2026-05-17-master-persona-componentization

- `skill`: ui-design
- `status`: human-review-ready
- `source_plan`: `./plan.md`
- `scenario_source`: `./scenario-design.md`
- `requirement_gate_source`: `./scenario-design.requirement-gate.md`
- `ux_standard_source`: `docs/UX-standard.md`
- `screen_design_source`: `docs/screen-design/screens/master-persona.md`
- `storybook_foundation_source`: `docs/exec-plans/completed/2026-05-18-storybook-foundation/storybook-review.md`
- `screen_design_diff`: `N/A`

## UI Contract

- `display_items`:
  - 画面状態領域は、通知、作成状態、状態値を維持する。
  - 生成準備は、AI 設定、入力 JSON、候補数、新規作成数、既存スキップ数を維持する。
  - 進行状況は、進捗、処理済み件数、作成済み件数、既存スキップ数、現在の対象を維持する。
  - 生成結果確認は、検索、プラグイン絞り込み、一覧、ページ操作、詳細を維持する。
  - 操作モーダルは、編集対象、編集入力、削除対象、確認操作を維持する。
- `primary_actions`:
  - `ペルソナを作成`: 生成準備の主要操作として扱う。
  - `モデル一覧を更新`: `AIModelSelectionCard` の既存操作として扱う。
  - `編集内容を保存`: 編集モーダルの主要操作として扱う。
  - `削除する`: 削除モーダルの危険操作として扱う。
- `button_enablement`:
  - `ペルソナを作成` は、入力 JSON、AI 設定、作成可能件数、実行状態が作成開始可能な時だけ有効にする。
  - `一時停止` と `中止` は、生成中だけ有効にする。
  - `編集` と `削除` は、選択中ペルソナがあり、変更可能状態の時だけ有効にする。
  - `前へ` と `次へ` は、現在ページと総ページから移動可能な時だけ有効にする。
- `state_variants`:
  - 未選択、JSON 選択済み、preview 成功、AI 設定不足、モデル一覧更新中、長い model 名。
  - 生成前、生成中、生成失敗、生成完了。
  - 空一覧、一覧あり、行選択済み、絞り込み後空、古い選択なし。
  - 編集モーダル、削除モーダル、保存失敗、削除失敗。
- `post_implementation_review`:
  - Storybook でパネル、カード、モーダルを順に確認する。
  - Storybook review URL、story ID、確認状態、未確認理由、再実行 command は `storybook-review.md` に残す。
  - `build-storybook` と `frontend-local` の結果または未実行理由を確認する。

## Interface Frame

- `purpose`: マスターペルソナ生成画面を、実画面と Storybook の両方で確認できる表示部品へ分ける。
- `audience`: Skyrim Mod 向け翻訳作業を行う利用者。
- `primary_workflow`: JSON を選ぶ。AI 設定を確認する。ペルソナを作成する。結果を検索して確認する。必要に応じて編集または削除する。
- `information_density`: 作業画面として中密度にする。生成準備、進行状況、結果確認、モーダルを別の判断単位に分ける。
- `visual_direction`: 既存の `master-persona` 画面の panel、card、modal、button、余白、配色を維持する。独自 page shell、独自 card、独自 grid、独自配色を作らない。
- `remembered_signal`: `AIModelSelectionCard` は既存共有 component を再利用する。手作りの代替カードを作らない。

## Storybook Review Contract

- `review_unit_order`:
  1. `AIModelSelectionCard`: Storybook 基盤で確認済みの共有 card を基準表示として確認する。
  2. `GenerationSetupPanel`: AI 設定 card と入力 JSON panel の組み合わせを確認する。
  3. `RunStatusPanel`: 生成中と非生成中の進行状態を確認する。
  4. `PersonaReviewPanel`: 一覧、詳細、空状態、選択状態を確認する。
  5. `PersonaActionModal`: 編集、削除、失敗時の対象保持を確認する。
  6. `MasterPersonaPage` 合成 story: 任意とする。主要表示領域 story の不足を page 合成 story で代替しない。
- `review_target_state_order`:
  - 生成準備: 未選択、JSON 選択済み、preview 成功、AI 設定不足、モデル一覧更新中、長い model 名。
  - 進行状況: 生成前、生成中、生成失敗、生成完了。
  - 結果確認: 空一覧、一覧あり、行選択済み、絞り込み後空。
  - モーダル: 編集、削除、保存失敗、削除失敗。
- `fixture_boundary`:
  - fixture は fixed props、view model fixture、callback stub だけを使う。
  - fixture は外部接続、secret、API key、token、実 endpoint、ローカル絶対 path、実ユーザーデータ、raw request、raw response、raw prompt を持たない。
  - story と fixture は backend、Wails runtime、generated `wailsjs`、Gateway、RuntimeEventAdapter、AI provider、secret store、DB、実 filesystem flow を要求しない。

## UI Component Boundary

- `reuse_components`:
  - `AIModelSelectionCard`: 共有 component として再利用する。provider、model、実行方法、資格情報状態、モデル一覧更新をこの card へ渡す。
- `screen_local_components`:
  - `GenerationSetupPanel`: 生成準備の section frame、AI 設定 card、入力 JSON、preview 件数、作成開始操作を扱う。
  - `RunStatusPanel`: 生成状態、進捗、処理件数、現在の対象、一時停止、中止を扱う。
  - `PersonaReviewPanel`: 検索、プラグイン絞り込み、一覧、ページ操作、詳細、編集開始、削除開始を扱う。
  - `PersonaActionModal`: 編集モーダル、削除モーダル、保存失敗、削除失敗を扱う。
- `thin_page_conditions`:
  - `MasterPersonaPage` は controller の mount、dispose、subscribe、file input bridge、通知、状態領域、表示部品の合成だけを持つ。
  - `MasterPersonaPage` は production controller factory、Gateway、generated binding、RuntimeEventAdapter を screen local component へ渡さない。
  - story 対象 component は `viewModel` 全体ではなく、画面 view model から作った小さい props と callback を受け取る。
  - callback は DOM event または画面操作結果を親へ返すだけにする。
- `do_not_split`:
  - page shell は新規部品にしない。既存 App shell と画面 wrapper の責務を増やさないためである。
  - `AIModelSelectionCard` の内部 UI はこの task で分割しない。Storybook 基盤で共有 component として確認済みであり、再利用が今回の前提である。
  - list row や stat card を共有 component 化しない。現時点では master-persona 専用の表示規則であり、共有化すると画面固有条件を props へ増やす可能性がある。

## Structure Notes

- `page_sections`:
  - 画面状態領域。
  - 生成準備パネル。
  - 進行状況パネル。
  - 生成結果一覧と詳細パネル。
  - 編集モーダルと削除モーダル。
- `layout_constraints`:
  - 生成準備は、AI 設定 card と入力 JSON panel の 2 カラムを維持する。
  - 生成結果確認は、一覧と詳細の 2 カラムを維持する。
  - モーダルは、対象識別情報と操作を同じ dialog 内で確認できる配置を維持する。
- `responsive_constraints`:
  - 生成準備は狭い幅で 1 カラムにする。
  - 生成結果確認は狭い幅で一覧、詳細の順に 1 カラムにする。
  - モーダルは狭い幅で幅いっぱいに近づけ、入力欄と操作ボタンを縦に並べる。
  - 長い model 名、plugin 名、EditorID、persona 本文は container 内で折り返す。
- `accessibility_constraints`:
  - 画面状態領域は `aria-label="マスターペルソナ作成状態"` を維持する。
  - 進行状況は `aria-label="生成進捗"` を維持する。
  - 結果確認は `aria-label="生成結果の確認"` を維持する。
  - ページ操作は `aria-label="ペルソナ一覧のページ操作"` を維持する。
  - モーダルは dialog として対象識別情報、入力欄、確定操作を読み上げ順に並べる。
- `screen_design_diff_refs`:
  - `N/A`
  - 理由: `docs/screen-design/screens/master-persona.md` は既存の表示項目、操作、状態、主要 selector をすでに定義している。今回の UI 変更は、利用者に見える画面内容の追加、削除、文言変更ではなく、部品境界と Storybook review 単位の固定である。

## UX Standard Review

- `source`: `docs/UX-standard.md`
- `screen_structure_high_priority_results`:
  - 画面目的: pass。目的はマスターペルソナ作成、確認、編集、削除で固定済みである。
  - 主要CTA: pass。生成準備の主要操作は `ペルソナを作成` である。
  - 情報階層: pass。生成準備、進行状況、結果確認、モーダルを判断単位で分ける。
  - 画面責務: caution。一覧、詳細、編集、削除を 1 画面に持つため、Storybook では panel と modal ごとに分けて確認する。
  - 状態別操作: pass。生成中は停止操作だけを有効にし、編集と削除を無効にする。
  - エラー状態: caution。保存失敗と削除失敗は、モーダルを閉じずに対象識別情報を残す fixture が必要である。
  - UI文言: pass。固定名以外は日本語の業務語を使う。
  - 検査可能性: pass。story 単位、fixture、review 順を固定する。
- `screen_structure_applicable_results`:
  - 空状態: 空一覧と絞り込み後空を story で確認する。
  - ローディング状態: モデル一覧更新中と生成中を story で確認する。
  - 危険操作: 削除は削除モーダルに分離する。
  - 証跡: review URL、story ID、確認状態は `storybook-review.md` に残す。
- `layout_responsive_high_priority_results`:
  - 視線開始位置: pass。画面状態、生成準備、主要操作の順に置く。
  - 主要CTA位置: pass。`ペルソナを作成` は入力 JSON panel の操作列に置く。
  - 危険操作位置: pass。削除は詳細の補助操作と削除モーダル内の危険操作に分ける。
  - カラム構造: pass。生成準備と結果確認の 2 カラムを狭い幅で 1 カラムにする。
  - 長文耐性: caution。長い model 名、plugin 名、persona 本文、エラー文言を story で確認する。
  - モーダル変形: caution。狭い幅で入力欄と操作列が縦に並ぶことを story または実装後確認で見る。
- `layout_responsive_applicable_results`:
  - 一覧: デスクトップは一覧と詳細の 2 カラム、狭い幅は一覧から詳細の順にする。
  - フォーム: 編集モーダルの入力欄は label と入力欄の対応を維持する。
  - 幅別証跡: 実装後に標準デスクトップ、狭い幅、長文 fixture を確認する。
- `deferred_items`:
  - 実装後の pixel 単位の余白、border、font size の調整は人間見た目レビューで扱う。
  - Storybook での幅別 screenshot 保存は実装後の `storybook-review.md` で扱う。

## Interaction States

- `loading`:
  - モデル一覧更新中は AI 設定操作と保存操作を一時的に無効にし、更新中状態を card 内で示す。
- `empty`:
  - JSON 未選択時は入力待ちとして表示する。
  - 一覧が空の場合は、作成後に確認できることを表示する。
  - 絞り込み後空の場合は、検索条件に合うペルソナがないことを表示する。
- `error`:
  - AI 設定不足は、作成開始できない理由として表示する。
  - 生成失敗は、進行状況 panel の状態として表示する。
  - 保存失敗と削除失敗は、モーダルを閉じずに対象と入力値を維持する。
- `disabled`:
  - 作成開始不可、生成中の編集削除不可、ページ移動不可、モデル未選択状態を明示する。
- `progress`:
  - 生成中は進捗バー、処理済み件数、作成済み件数、既存スキップ数、現在の対象を表示する。
- `retry`:
  - モデル一覧更新は再実行できる操作として表示する。
  - 保存失敗と削除失敗は、同じモーダル内で再実行できる状態を残す。
- `success`:
  - preview 成功時は、候補数、新規作成数、既存スキップ数を表示する。
  - 生成完了時は、一覧と詳細で結果確認へ進める状態を表示する。

## Post Implementation Review

- `desktop_review_points`:
  - 生成準備 2 カラムで、`AIModelSelectionCard` と入力 JSON panel の高さと余白が破綻しない。
  - 生成結果確認 2 カラムで、一覧、詳細、ページ操作、詳細操作が同時に読める。
  - 編集モーダルと削除モーダルで、対象識別情報と確定操作が同じ視界に入る。
- `mobile_review_points`:
  - 生成準備は AI 設定、入力 JSON の順で 1 カラムになる。
  - 結果確認は一覧、詳細の順で 1 カラムになる。
  - 操作ボタンは横幅不足で文字がはみ出さず、縦並びになる。
- `overflow_risks`:
  - 長い model 名。
  - 長い plugin 名。
  - 長い EditorID。
  - 長い persona 本文。
  - 長いエラー文言。
- `visual_polish_open_questions`:
  - `GenerationSetupPanel` と `AIModelSelectionCard` の二重枠が重く見えないか。
  - `PersonaReviewPanel` の一覧行密度が Storybook 上で狭すぎないか。
  - 編集モーダルの本文 textarea の高さが、標準デスクトップと狭い幅の両方で十分か。

## Agent Browser Review

- `command_source`: `agent-browser`
- `checked_url`: `N/A`
- `checked_viewports`: `N/A`
- `ux_standard_review`:
  - `source`: `docs/UX-standard.md`
  - `high_priority_results`: `UX Standard Review` に記録済み。
  - `applicable_results`: `UX Standard Review` に記録済み。
  - `deferred_items`: 実装後の Storybook review で確認する。
- `wording_review`:
  - `review_timing`: `before_agent_browser_review`
  - `fixed_names_preserved`: `MasterPersonaPage`、`AIModelSelectionCard`、Storybook story ID は固定名として扱う。
  - `business_japanese_terms`: 画面表示文言は既存の日本語業務語を維持する。
  - `internal_state_names_hidden`: story 名や fixture 内部名を画面表示文言として出さない。
  - `next_action_wording`: `ペルソナを作成`、`モデル一覧を更新`、`編集内容を保存`、`削除する` を維持する。
  - `allowed_english_labels`: `FormID`、`EditorID`、`JSON`、`AI`、`Storybook` は既存画面または task 固定名として許容する。
  - `plain_language_next_action_judgement`: pass。利用者が次に行う操作名として読める。
- `console_errors`: `N/A`
- `screenshot_or_snapshot_refs`: `N/A`
- `layout_breaks`: `N/A`
- `ambiguous_interactions`: `N/A`
- `open_issues`:
  - 実装後の Storybook で、長文、狭い幅、失敗状態の見た目確認が必要である。
- `not_checked_reason`: この段階ではプロダクトコード変更前であり、確認対象 story がまだ存在しない。既存画面の表示項目は `docs/screen-design/screens/master-persona.md` と現行 Svelte で確認できるため、アプリ起動と `agent-browser` 確認は必須ではない。

## Human Review Status

- `required`: yes
- `reason`: 部品境界、Storybook review 順、screen design diff 不要判断は、frontend implementation lane へ渡す前に人間レビューが必要である。
- `approval_target`:
  - `./ui-design.md`

## Rules

- UI は `ui-design.md` の UI 要件契約で固定する。
- 既存画面と既存 UI 部品を土台にする。
- `AIModelSelectionCard` を再利用する。
- 独自 page shell、card、grid、配色、余白体系を新規に作らない。
- 画面設計差分には、画面設計書正本へ適用できる恒久的な画面内容だけを書く。
- この task では恒久的な画面内容の差分がないため、`screen-design-diff.master-persona.md` は作成しない。
- Storybook fixture は外部接続、secret、実ユーザーデータを持たない。
