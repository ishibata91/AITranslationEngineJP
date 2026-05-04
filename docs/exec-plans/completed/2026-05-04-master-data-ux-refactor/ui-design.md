# UI Design: 2026-05-04-master-data-ux-refactor

- `skill`: ui-design
- `status`: approved-for-frontend-implementation
- `source_plan`: `./plan.md`
- `ui_prototype`: `./prototype/index.svelte`
- `prototype_server_url`: `http://127.0.0.1:34118/prototype`
- `prototype_server_command`: `npm --prefix frontend run dev:prototype -- --task 2026-05-04-master-data-ux-refactor --port 34118`
- `human_review_server_required`: `yes`
- `human_review_designer_agent_required`: `yes`
- `human_feedback_route`: `designer_agent_direct`
- `designer_agent_close_after_review`: 人間 UI レビューは承認済み。`ux_refactor_lane` が `designer` agent を閉じ、frontend 実装へ進む。
- `ux_standard_source`: `docs/UX-standard.md`

## UI Contract

- `purpose`: 翻訳前の準備として、ベースゲームや大型 Mod の NPC ペルソナを事前に作成し、作成後に一覧と詳細で確認する。
- `audience`: Skyrim Mod 翻訳の作業者。AI 設定や生成状態の内部名を知らなくても操作できる利用者を前提にする。
- `hero_title`: `マスターペルソナ作成`
- `hero_description`: ベースゲームや大型 Mod の NPC を対象に、翻訳前の準備としてペルソナをまとめて作成する目的を短く説明する。
- `primary_action`: `ペルソナを作成`
- `secondary_actions`: `JSON を選択`, `モデル一覧を更新`, `一時停止`, `中止`, `編集`, `削除`
- `fixed_reuse`: モデル選択カードは既存 `AIModelSelectionCard.svelte` をそのまま使う。
- `prototype_reuse`: task-local UIプロトタイプは `AIModelSelectionCard.svelte` を `/@fs/.../frontend/src/ui/components/AIModelSelectionCard.svelte` から import する。

## Display Items

- `keep`:
  - 画面目的、JSON 選択状態、候補数、新規作成数、既存スキップ数を表示する。
  - AI サービス、モデル、処理方式、モデル一覧更新を既存モデル選択カードで表示する。
  - 生成状態、進捗、処理済み件数、作成済み件数、スキップ件数、現在の対象を表示する。
  - 作成後確認のため、一覧にはプラグイン名と NPC 名だけを表示する。
  - 一覧は 1000 から 10000 件規模を前提に、ページ範囲、総件数、ページ操作を表示する。
  - 詳細には識別情報、声、話し方、ペルソナ本文を表示する。
- `remove_or_deemphasize`:
  - `Gateway: ...` は利用者の判断材料ではないため画面から削除する。
  - `preview 状態: ...` のような内部状態名は削除し、`入力待ち`、`JSON 選択済み` などの日本語状態へ変換する。
  - プロンプトテンプレート説明は、利用者が変更できないため通常表示から削除する。
  - 一覧の `className`、`sourcePlugin`、`race`、`sex` は初期表示から削除し、必要なら詳細補助へ寄せる。
  - 一覧の `personaSummary` と `voiceType` は大量件数では行高を増やすため削除し、選択後の詳細だけで扱う。
  - 生成待ちのペルソナは一覧へ表示しない。生成済みまたは既存作成済みの確認対象だけを表示する。
  - `voice`、`class`、`source` の英語ラベルは、ユーザー向けには `声`、`クラス`、`元プラグイン` へ置き換える。
- `human_consultation_required_before_addition`:
  - 生成前に見積もり時間、料金目安、生成対象サンプルを追加するか。
  - 既存スキップの理由一覧を表示するか。
  - 作成後に「次に確認するペルソナ」を推奨表示するか。
  - プロンプト内容または生成方針の説明を利用者向けに表示するか。

## Main Operations

- `select_json`: JSON 未選択では `ペルソナを作成` を無効にする。
- `generate_persona`: JSON 選択済みで、AI サービスとモデルが利用可能な場合だけ有効にする。
- `pause_or_cancel`: 生成中だけ `一時停止` と `中止` を有効にする。
- `review_result`: 生成中または完了後も一覧と詳細を確認できる。
- `paginate_persona_list`: 大量件数を前提に、一覧はページ単位で切り替える。
- `edit_or_delete`: 選択中ペルソナがある場合だけ `編集` と `削除` を有効にし、押下時はモーダルで確認する。

## Button Enablement

- `JSON を選択`: 常に有効にする。
- `ペルソナを作成`: JSON 選択済み、AI 設定利用可能、生成中ではない場合だけ有効にする。
- `モデル一覧を更新`: `AIModelSelectionCard.svelte` の既存 props で制御する。
- `一時停止` と `中止`: 生成中だけ有効にする。
- `前へ` と `次へ`: 先頭ページまたは最終ページでは無効にする。
- `編集` と `削除`: 選択中ペルソナがあり、生成処理で編集がロックされていない場合だけ有効にする。

## Wording Policy

- `provider` は画面では `AI サービス` と表示する。
- `model` は画面では `モデル` と表示する。
- `executionMethod` は画面では `処理方式` と表示する。
- `preview` は画面に出さない。JSON 内容確認機能がないため、専用ボタンを置かない。
- `runState` は画面では `入力待ち`、`JSON 選択済み`、`生成中`、`完了`、`エラー` へ変換する。
- `Gateway`、raw request、raw response、raw prompt は表示しない。

## State Variants

- `empty`: JSON 未選択。候補数は 0、主要 CTA は無効。
- `ready`: JSON 選択済み。主要 CTA を有効にする。
- `running`: 進捗、処理済み件数、現在の対象を表示し、生成開始の二重実行を防ぐ。
- `complete`: 完了文言を表示し、一覧と詳細で結果確認へ誘導する。
- `error`: エラー理由と次操作を表示し、同じ JSON 選択または再確認へ戻れるようにする。
- `locked`: 生成中に編集不可の場合は、編集と削除を無効化して理由を表示する。

## Structure Notes

- `page_sections`:
  - Hero のタイトルは `マスターペルソナ作成` だけにする。
  - Hero の説明文はタイトルと分け、ベースゲームや大型 Mod 向けの事前作成目的を短く示す。
  - 生成準備区画にモデル選択カードと JSON 選択を置く。
  - 進行状況区画に生成中の状態と中止操作を置く。
  - 生成結果区画に一覧と詳細を置く。
  - 一覧はページング前提の細い行にし、1 行はプラグイン名と NPC 名だけにする。
  - 編集と削除は詳細区画の補助操作にし、押下後はモーダルを表示する。
- `layout_constraints`:
  - desktop は生成準備を 2 カラム、生成結果を一覧と詳細の 2 カラムにする。
  - mobile はすべて 1 カラムへ落とし、JSON 操作、進行状況、一覧、詳細の順に読む。
  - 主要 CTA は JSON 選択区画の末尾に置き、危険操作から離す。
- `responsive_constraints`:
  - 390px 幅でモデルカード、JSON ファイル名、エラー文、ペルソナ本文が横にはみ出さない。
  - 一覧は mobile でカード化し、列見出しを表示しない。
  - 一覧行は 44px 程度の高さを保ち、ページング表示と組み合わせて大量件数でも密度を維持する。
  - 状態確認用ナビゲーションは折り返しても主要操作を隠さない。
- `accessibility_constraints`:
  - 状態は色だけで示さず、状態 pill と説明文で示す。
  - 主要操作後の結果は `role=status` で通知する。
  - 検索、プラグイン選択、モデル選択にはラベルを付ける。

## UI Component Decision

- `reuse_shared_component`:
  - `AIModelSelectionCard.svelte`: 複数画面で使う共有部品であり、変更せず利用する。
- `screen_local_component_candidates`:
  - `GenerationSetupPanel`: 入力 JSON と生成開始条件を閉じ込める画面専用部品。
  - `RunStatusPanel`: 生成状態と中止操作を閉じ込める画面専用部品。
  - `PersonaReviewPanel`: 一覧と詳細の選択確認を閉じ込める画面専用部品。
  - `PersonaActionModal`: 編集フォームと削除確認を閉じ込める画面専用部品。
- `do_not_componentize`:
  - 画面全体の並びと状態確認タブは画面責務が大きいため、共有部品化しない。
  - mock data は確認用であり、部品として本番へ移植しない。
- `architecture_basis`: `docs/architecture.md` の `UI Component` に従い、画面専用部品は `frontend/src/ui/screens/master-persona/` 配下の候補として扱う。共有化は複数画面で使う場合だけ検討する。

## UX Standard Review

- `source`: `docs/UX-standard.md`
- `screen_structure_high_priority_results`:
  - `screen_purpose`: pass。目的をベースゲームや大型 Mod 向けの事前作成と作成後確認へ絞った。
  - `primary_action`: pass。主要 CTA は `ペルソナを作成` に固定した。
  - `information_hierarchy`: pass。入力、確認、生成、結果確認の判断順に並べた。
  - `state_and_allowed_actions`: pass。入力前、JSON 選択済み、生成中、完了、エラーで操作差分を分けた。
  - `ui_wording`: pass。内部状態名と英語ラベルを日本語の業務語へ変換する。
- `screen_structure_applicable_results`:
  - `empty_state`: pass。JSON 未選択時に次操作を表示する。
  - `error_state`: pass。JSON 読み取り失敗時はファイル形式確認へ誘導する。
  - `dangerous_action`: pass。削除は詳細区画の補助操作として主要 CTA から離す。
  - `list_structure`: pass。一覧は大量件数での視認性を優先し、プラグイン名と NPC 名へ絞る。
  - `large_list_density`: pass。一覧は 1000 から 10000 件を前提に、プラグイン名と NPC 名だけの細い行へ絞る。
- `layout_responsive_high_priority_results`:
  - `column_drop`: pass。980px 以下で 1 カラムへ切り替える。
  - `long_text`: pass。ペルソナ本文、ファイル名、エラー文は折り返す。
  - `tap_target`: pass。主要ボタンは 40px 以上の高さを保つ。
- `layout_responsive_applicable_results`:
  - `table_breakdown`: pass。一覧はテーブルではなくカード型にする。
  - `fixed_element_interference`: pass。固定 CTA は置かず、スクロール順を優先する。
- `deferred_items`:
  - 実装後に既存 AppShell 内で navigation 幅と画面余白を確認する。
  - 実データの長い NPC 名、長いペルソナ本文、モデル名で崩れを確認する。

## Interaction States

- `loading`: モデル一覧更新中は `AIModelSelectionCard.svelte` の更新表示を使う。
- `empty`: JSON 未選択ではファイル選択だけを促す。
- `error`: JSON 読み取り失敗と生成失敗は、原因と再実行導線を表示する。
- `disabled`: 無効ボタンは理由が近くの文言で分かるようにする。
- `progress`: 生成中は進捗、処理済み件数、現在の対象を表示する。
- `retry`: エラー後は JSON 選択から再開できる。
- `success`: 完了後は生成結果一覧を確認できる状態にする。

## Post Implementation Review

- `desktop_review_points`:
  - 1440px 幅で、目的、生成準備、進行状況、生成結果の順に読める。
  - モデル選択カードが既存部品の見た目と挙動を維持している。
  - 生成準備と生成結果の 2 カラムが横幅不足で破綻しない。
  - 削除操作が主要 CTA と混同されない。
- `mobile_review_points`:
  - 390px 幅で、主要 CTA と状態文が横にはみ出さない。
  - 一覧がカード型になり、詳細へ移る読み順が自然である。
  - モデル選択カードの select と更新ボタンが横スクロールを発生させない。
- `overflow_risks`:
  - 長いモデル名、長い NPC 名、長いペルソナ本文、長い JSON ファイル名。
  - 大量件数の一覧表示では、ページングなしで全件を描画しない。
  - 生成失敗時の長いエラー文。
- `visual_polish_open_questions`:
  - 実 AppShell へ組み込んだ後、既存余白と panel density に合わせる必要がある。

## UI Prototype Contract

- `prototype_kind`: `existing_screen_change`
- `source_basis`: `./plan.md`, `frontend/src/ui/screens/master-persona/MasterPersonaPage.svelte`, `frontend/src/ui/components/AIModelSelectionCard.svelte`, `frontend/src/application/gateway-contract/master-persona/master-persona-gateway-contract.ts`, `docs/architecture.md`, `docs/UX-standard.md`
- `existing_screen_resource_refs`: `MasterPersonaPage.svelte`, `AIModelSelectionCard.svelte`
- `reused_screen_structure`: 既存画面の生成準備、進行状況、一覧、詳細の業務区画を維持する。表示順と文言は画面目的に合わせて再構成する。
- `changed_sections_only`: task-local prototype は生成画面全体の UI 改善契約確認用であり、product code は変更しない。
- `new_visual_system_added`: `no`
- `new_visual_system_reason`: 既存の dark panel、status pill、button、grid の範囲で構成する。
- `prototype_path`: `./prototype/index.svelte`
- `prototype_screen_list`: `入力前`, `選択済み`, `生成中`, `完了`, `エラー`, `編集モーダル`, `削除モーダル`
- `prototype_route_default`: `選択済み`
- `prototype_route_navigation`: `in_prototype_only`
- `required_before_human_review`: `yes`
- `required_for_frontend_handoff`: `yes`
- `framework_conversion`: UIプロトタイプの構造を frontend framework へ変換する。
- `prototype_server_url`: `http://127.0.0.1:34118/prototype`
- `prototype_server_command`: `npm --prefix frontend run dev:prototype -- --task 2026-05-04-master-data-ux-refactor --port 34118`
- `human_review_server_required`: `yes`
- `human_review_designer_agent_required`: `yes`
- `human_feedback_route`: `designer_agent_direct`
- `designer_agent_close_after_review`: 人間レビュー終了後に反映と確認を行ってから停止する。
- `ux_standard_source`: `docs/UX-standard.md`
- `mock_data_root`: `./mock-data/`
- `mock_data_file`: `./mock-data/master-persona-prototype.json`
- `mock_data_migration`: `forbidden`
- `mock_data_separation`: UIプロトタイプ本体へ一覧データやページング件数を直接書かず、`mock-data/` から import する。
- `sample_data_root`: `[data-ui-prototype-sample-data-root]`
- `sample_data_migration`: `forbidden`
- `production_reference_direction`: `product_code_must_not_reference_ui_prototype`
- `interaction_review`: 状態確認タブ、JSON 選択、生成開始、一時停止、中止、完了状態、一覧選択、ページ操作、編集モーダル、削除モーダルを確認する。
- `state_transition_review`: 入力前、選択済み、生成中、完了、エラー、編集モーダル、削除モーダルを確認する。
- `prototype_modal_review_tabs`: `編集モーダル` と `削除モーダル` は、agent-browser と人間レビューで画面外クリックに依存せず確認できる状態確認タブとして用意する。
- `wording_review`: 内部語を画面へ出さず、目的ベースの日本語文言へ変換する。
- `structure_to_preserve`: モデル選択カード、JSON 入力、プレビュー数値、進行状況、一覧、詳細、編集、削除を維持する。
- `allowed_changes_during_conversion`: view model 名への接続、既存 controller event への接続、画面専用部品への分割。
- `forbidden_changes_during_conversion`: product code から prototype 参照、mock data 移植、`AIModelSelectionCard.svelte` 変更、モデル選択カード自作、公開契約変更。

## Agent Browser Review

- `command_source`: `agent-browser`
- `served_url`: `http://127.0.0.1:34118/prototype`
- `server_command`: `npm --prefix frontend run dev:prototype -- --task 2026-05-04-master-data-ux-refactor --port 34118`
- `standard_port_note`: `34116` と `34117` は使用中または古い確認状態だったため、人間レビュー用サーバーは `34118` で起動した。
- `server_status_during_human_review`: `running`
- `mock_data_refs`: `./mock-data/master-persona-prototype.json`, `[data-ui-prototype-sample-data-root]`
- `used_only_for_display_state_review`: `yes`
- `migration_to_product_code`: `forbidden`
- `migration_to_fixture_or_test_data`: `forbidden`
- `checked_viewports`: `desktop current agent-browser viewport`, `full-page screenshot`
- `ux_standard_review`:
  - `source`: `docs/UX-standard.md`
  - `high_priority_results`: pass。画面目的、主要 CTA、状態別操作、入力制約、エラー回復、検査可能性を snapshot で確認した。
  - `applicable_results`: pass。モデル選択カード、JSON 選択区画、進行状況、一覧、詳細が確認できた。
  - `deferred_items`: 実 AppShell 組み込み後の余白、既存 navigation 幅、実データ長文。
- `interaction_review`:
  - `edit_modal`: pass。状態確認タブの `編集モーダル` 押下で編集モーダルが表示される。
  - `edit_modal_font`: pass。編集モーダルの見出しは 24px、入力欄は 16px / line-height 25.6px で、周辺パネルと同じ本文基準に合わせた。
  - `delete_modal`: pass。状態確認タブの `削除モーダル` 押下で削除確認モーダルが表示される。
  - `persona_list_paging`: pass。`1-50 / 4280 件`、`1 / 86 ページ`、`前へ`、`次へ` を確認し、`次へ` 押下後に `51-100 / 4280 件` へ変わることを確認した。
  - `persona_list_row_density`: pass。一覧行はプラグイン名と NPC 名だけを表示し、行高 44px で確認した。
  - `generated_waiting_rows`: pass。生成待ちの行はモックデータと一覧表示から外した。
  - `component_action_wiring`: pass。主要操作は子コンポーネントから親へ渡す関数 props で扱い、プロトタイプ内の window event 中継へ依存しない。
- `wording_review`:
  - `review_timing`: `after_agent_browser_review`
  - `fixed_names_preserved`: `AIModelSelectionCard.svelte`, `Gemini`, `LM Studio`, `xAI`, `JSON`, `FormID`, `EditorID`
  - `business_japanese_terms`: `AI サービス`, `モデル`, `処理方式`, `JSON を選択`, `ペルソナを作成`
  - `internal_state_names_hidden`: pass。`Gateway`、`runState`、`preview status` などの内部名は主要表示に出していない。
  - `next_action_wording`: pass。JSON 選択、ペルソナ作成、完了後確認、編集、削除を画面目的ベースで表示した。
  - `allowed_english_labels`: 固有名、ファイル形式、識別子だけ許可する。
  - `plain_language_next_action_judgement`: pass。専門知識なしで次操作を読める日本語水準として確認した。
- `console_errors`: `none`
- `screenshot_or_snapshot_refs`: `tmp/agent-browser/master-persona-ux-prototype-desktop.png`, `tmp/agent-browser/master-persona-ux-prototype-full.png`, `tmp/agent-browser/master-persona-ux-prototype-final.png`, `agent-browser snapshot -i --compact --depth 7`
- `layout_breaks`: desktop snapshot では主要操作の消失なし。mobile は agent-browser CLI の viewport 指定がないため実画像未確認。
- `ambiguous_interactions`: 表示項目追加候補は人間判断待ち。
- `open_issues`: 人間相談事項を回答後、UI 改善契約を承認済みにする。
- `not_checked_reason`: mobile viewport の実画像は未確認。理由は agent-browser CLI に viewport 指定がないため。
