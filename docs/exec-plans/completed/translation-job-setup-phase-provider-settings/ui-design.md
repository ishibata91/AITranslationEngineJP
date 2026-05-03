# UI Design: translation-job-setup-phase-provider-settings

- `skill`: ui-design
- `status`: pending-human-review
- `source_plan`: `./plan.md`
- `scenario_source`: `./scenario-design.md`
- `ui_prototype`: `./prototype.svelte`
- `prototype_server_url`: `http://127.0.0.1:34116/prototype`
- `prototype_server_command`: `npm --prefix frontend run dev:prototype -- --task translation-job-setup-phase-provider-settings --port 34116`
- `human_review_server_required`: `yes`

## UI Contract

- `display_items`:
  - 入力済み翻訳データの対象として、翻訳ジョブ作成画面に 3 つの翻訳段階を表示する。
  - 翻訳段階は、単語翻訳、NPC ペルソナ生成、本文翻訳とする。
  - 各翻訳段階には、AIサービス、使うモデル、APIキー状態、一括処理の設定、設定済みまたは未設定の状態を表示する。
  - 作成後の設定内容には、翻訳段階ごとの AIサービス、モデル、APIキー状態、一括処理の有無だけを表示する。
  - APIキーの文字列、秘密情報、外部サービスの生データ、内部ログ用の識別子は表示しない。
- `primary_actions`:
  - 利用者は翻訳段階ごとに AIサービスを選ぶ。
  - 利用者は APIキー未設定の AIサービスでは、同じ画面のモーダルで APIキーを登録する。
  - 利用者は AIサービスごとに取得したモデル一覧から使うモデルを選ぶ。
  - 利用者は Gemini と xAI の場合だけ、一括処理で実行するかをチェックボックスで切り替える。
  - 利用者はチェックボックス横の `?` の説明で、一括処理がバッチAPIを使い、安く済ませられることがある設定だと確認できる。
  - 利用者は 3 つの翻訳段階で APIキー不足とモデル未選択がない時だけ、翻訳ジョブを作成する。
- `button_enablement`:
  - APIキー登録は、APIキーが必要で未設定の AIサービスを選んだ場合だけ表示する。
  - モデル一覧更新は、APIキーが必要な AIサービスで APIキー未設定の場合は押せない。
  - モデル選択は、モデル一覧取得に成功した場合だけ表示する。
  - 翻訳段階ごとの設定確認ボタンは表示しない。
  - 翻訳ジョブ作成は、3 つの翻訳段階で APIキー不足とモデル未選択がない場合だけ押せる。
- `state_variants`:
  - モデル一覧未更新、更新中、取得済み、取得失敗、APIキー未設定で更新不可を分けて表示する。
  - 設定済み、APIキー未設定、モデル未選択、モデル一覧取得失敗、モデル一覧更新中を分けて表示する。
  - LM Studio は APIキー状態を不要として表示し、APIキー入力欄、APIキー未設定警告、APIキー選択欄を出さない。
  - モデル一覧取得失敗時と APIキー未設定時は、手動モデル入力欄を出さない。
  - APIキー登録モーダル保存後は、モデル一覧を再更新する必要があることを表示する。
- `post_implementation_review`:
  - desktop と mobile の両方で、入力確認、共通基盤、翻訳段階別設定、作成前確認、作成実行の順に読めることを確認する。
  - 各翻訳段階では、AIサービス、モデル、実行方法、現在状態が意味単位で分かれていることを確認する。
  - 翻訳ジョブ作成ボタンの無効理由が画面文言から読めることを確認する。

## Interface Frame

- `purpose`: 非エンジニアの日本語話者が、翻訳段階ごとに使う AIサービスとモデルを選び、作成前に不足を解消できるようにする。
- `audience`: Skyrim Mod 翻訳の作業者。外部 AI サービスの内部実装やログ名を知らない利用者を前提にする。
- `primary_workflow`: AIサービスを選ぶ、APIキーが未設定ならモーダルで登録する、OS 認証画面が出る可能性と理由を確認する、モデル一覧を更新する、モデルを選ぶ、一括処理を必要に応じて切り替える、不足がない状態で翻訳ジョブを作成する。
- `information_density`: 画面全体を UX 標準の判断順に再構成する。目的、前提、設定、確認、実行を分け、3 つの翻訳段階は 1 段階ずつ、AIサービス、モデル、実行方法、状態の小区画で読める密度にする。
- `visual_direction`: 既存 Job Setup 画面を土台にしつつ、英語の内部語を日本語の業務語へ置き換える。過剰な装飾は入れない。
- `remembered_signal`: 共通基盤は参照対象として確認し、画面本文では翻訳段階ごとの AI 設定に集中する。

## Structure Notes

- `page_sections`:
  - 既存 `hero-card` は画面目的と短い説明文だけを置く。
  - 既存 `content-grid` に、入力データ、共通辞書、共通ペルソナ、翻訳段階別 AI 設定を判断順で置く。
  - 既存の単一 runtime select 区画だけを、単語翻訳、NPC ペルソナ生成、本文翻訳の設定区画へ差し替える。
  - 翻訳段階別 AI 設定は、各段階の中で AIサービス、モデル、実行方法、状態文を小区画へ分ける。
  - 作成前確認はカードにせず、ページ下部固定バーで未設定理由と `次へ` CTA だけを表示する。
  - 作成後は、既存 `summary-grid` と `detail-grid` で作成後の設定内容を読み取り専用で表示する。
- `layout_constraints`:
  - desktop でも既存 `content-grid` と `summary-grid` を 1 列にし、読み順とスクロール順を一致させる。
  - desktop ではモデルカード以外のヒーロー、入力、共通基盤を compact 表示にし、画面全体の縦量を抑える。
  - 作成前確認の CTA は `次へ` とし、アンカー移動に頼らずページ全体の下部固定バーで見失わない配置にする。
  - 翻訳段階別 AI 設定は、desktop では単語翻訳、NPC ペルソナ生成、本文翻訳を横に並べて比較できるようにし、各段階の内部は判断単位で分離する。
  - 共通辞書と共通ペルソナは別テーブルに分け、参照項目が増えても画面全体を押し下げないように高さ上限を付ける。
  - 横幅が狭い場合は既存 media query に従い、カード見出しと操作を縦に並べる。
  - 各翻訳段階では、AIサービス、APIキー状態、モデル一覧、モデル選択、一括処理、設定済み状態の順に置く。
- `responsive_constraints`:
  - mobile では横スクロールなしで主要操作を完了できることを実装後に確認する。
  - 作成後の設定内容だけは表形式のため、狭い画面では表内スクロールを許容する。
  - 長いモデル名は折り返しまたは表内スクロールで表示する。
- `accessibility_constraints`:
  - 各 select と checkbox には日本語ラベルを付ける。
  - 状態は色だけで伝えず、設定済み、未設定ありなどの文言を併記する。
  - ボタン無効状態には、近くに次に必要な操作を表示する。
  - APIキー登録モーダルは、背景ブラーと暗幕で背面情報を抑え、本文を読めるコントラストにする。

## 画面構造UX Review

- `screen_purpose`: pass。画面目的は、翻訳ジョブ作成前に 3 つの翻訳段階の AI 設定を確認する 1 点に固定した。
- `target_user`: pass。対象は、外部 AI サービスの内部名を知らない Skyrim Mod 翻訳作業者に固定した。
- `primary_action`: pass。主要操作は「翻訳ジョブを作成」に絞り、モデル一覧更新と APIキー登録は補助操作として扱う。
- `information_hierarchy`: pass。画面冒頭は目的と説明文だけに絞り、その後に入力データ、共通基盤、翻訳段階別 AI 設定を配置した。未設定理由と次操作はページ下部固定バーに集約した。
- `screen_responsibility`: pass。一覧、編集、確認、完了表示は 1 画面内に残るが、作成前設定という既存 Job Setup の責務範囲内に収めた。作成前確認カードは削除し、未設定理由と `次へ` だけを下部固定バーへ移した。
- `state_and_allowed_actions`: pass。APIキー未設定、モデル一覧未更新、取得済み、取得失敗、モデル未選択で可能操作と禁止操作を分けた。
- `display_and_blocking_conditions`: pass。APIキー登録、モーダル、モデル選択、モデル一覧更新、作成実行の表示条件と無効条件を UI Contract に固定した。
- `recovery_route`: pass。APIキー未設定時は APIキー登録ボタンを表示し、モーダル保存後にモデル一覧更新から再開する。
- `permission_difference`: not_applicable。ユーザー権限差分は今回の Job Setup UI 変更の対象外である。
- `input_constraints`: pass。手動モデル入力を出さず、モデル一覧から選ぶ制約を画面文言で示す。
- `input_grouping`: pass。入力項目は翻訳段階ごとにまとまり、各段階の中でも AIサービス、モデル、実行方法へ分けた。
- `async_empty_loading_error`: pass。モデル一覧未更新、更新中、取得失敗、APIキー未設定で次行動を分けた。
- `change_delta`: pass。設定変更後は model 未選択に戻し、作成不可理由に不足だけを表示する。
- `dangerous_action`: not_applicable。削除や権限剥奪などの危険操作はない。
- `navigation_after_completion`: pass。作成後は `summary-grid` に結果を表示し、設定画面へ戻る補助操作を置く。
- `maintainability`: pass。実装流用可否、禁止移植、確認 URL、確認コマンドを UI Prototype Contract と Agent Browser Review に分けて残した。

## Interaction States

- `loading`: モデル一覧更新中は、対象翻訳段階に「モデル一覧を更新しています。」と表示する。
- `empty`: モデル一覧未更新では「モデル一覧を更新してください。」と表示し、モデル手入力欄は出さない。
- `error`: モデル一覧取得失敗では「モデル一覧を取得できませんでした。時間をおいて再実行してください。」と表示し、手入力では進めないことを補足する。
- `disabled`: APIキーが必要な AIサービスで APIキー未設定の場合、モデル一覧更新と設定確認を無効にし、APIキー登録モーダルを開く補助操作を表示する。
- `progress`: 設定確認中は対象翻訳段階を「確認中」にする。
- `retry`: モデル一覧取得失敗時は、モデル一覧更新ボタンから再実行できる。
- `success`: 設定済みの翻訳段階は「設定済み」と表示し、作成前の設定済み数に反映する。
- `credential_modal`: APIキー登録モーダルは、APIキー入力、AIサービスへのリクエストに使う理由、OS 認証画面が出る可能性、OS の安全な保管場所へ保存する理由、保存後の次行動を表示する。
- `return_from_modal`: APIキー登録モーダル保存後は、対象翻訳段階のモデル一覧更新から再開する。

## Wording Rules

- `provider` は画面に出さず、利用者向けには「AIサービス」と表示する。
- `credential` は画面に出さず、利用者向けには「APIキー状態」と表示する。
- `getModels` は画面に出さず、利用者向けには「モデル一覧を更新」と表示する。
- 設定変更後の内部状態名は画面に出さず、利用者向けには「本文翻訳のモデルを選んでください」のような不足理由で表示する。
- `Ready job summary` は画面に出さず、利用者向けには「作成後の設定内容」と表示する。

## Post Implementation Review

- `desktop_review_points`:
  - 1440px 幅でも、画面目的、次操作、入力確認、共通基盤、フェーズ別設定の順で読めること。
  - モデルカード以外の周辺セクションが過大にならず、ヒーロー、入力、共通基盤が compact に収まること。
  - 未設定理由と `次へ` がページ下部の固定バーとして残ること。
  - 各翻訳段階の中で AIサービス、モデル、実行方法、状態が意味単位として分離して見えること。
  - モデル名が長い場合でも、select と作成後の設定内容が破綻しないこと。
  - Gemini と xAI だけに一括処理チェックボックスが出ること。
- `mobile_review_points`:
  - 390px 幅で各翻訳段階が縦に並び、操作順が読めること。
  - 共通辞書と共通ペルソナの別テーブルが高さ上限内で読め、画面全体の横スクロールを発生させないこと。
  - 作成前の未設定理由と `次へ` ボタンがページ下部で見失われないこと。
  - LM Studio の区画に APIキー入力欄、APIキー未設定警告、APIキー選択欄が出ないこと。
- `overflow_risks`:
  - 長いモデル名、長い翻訳データ名、エラー文が区画外へはみ出す可能性がある。
  - 共通辞書と共通ペルソナは利用対象が増える可能性があるため、別テーブル内スクロールで高さを抑える。
  - 作成後の設定内容は列数が多いため、mobile では表内スクロールを許容する。
- `visual_polish_open_questions`:
  - 実装後の実画面で、既存 Job Setup 画面のカード密度と余白へ合わせる必要がある。
  - 実装後の実画面で、既存の戻る導線と作成後表示の切り替え位置を確認する必要がある。

## UI Prototype Contract

- `prototype_kind`: `existing_screen_change`
- `source_basis`: `./scenario-design.md` と `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`
- `prototype_path`: `./prototype.svelte`
- `required_before_human_review`: `yes`
- `required_for_frontend_handoff`: `yes`
- `existing_screen_resource_reference`: `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`
- `reused_screen_structure`: `job-setup-shell`, `job-setup-card`, `hero-card`, `content-grid`, `summary-grid`, `section-head`, `detail-grid`, `status-pill`, `field-block`, `button-primary`, `button-secondary` を再利用した。
- `changed_section_only`: 既存の `provider / model / execution mode` 単一選択区画を、3 翻訳段階別の AIサービス、モデル、APIキー状態、一括処理、設定確認へ差し替えた。
- `api_key_registration`: APIキー未設定の AIサービスだけ、同じ Job Setup 画面内のモーダルで APIキーを登録する。別画面への遷移は作らない。
- `os_credential_warning`: APIキー保存時に OS 認証画面が出る可能性と、OS の安全な保管場所へ保存する理由をモーダル内に表示する。
- `preserved_sections`: 入力データ、共通辞書、共通ペルソナ、作成前確認、作成後表示の基本構造を既存 Job Setup 画面へ寄せた。
- `new_visual_system`: `none`
- `new_page_shell_card_grid_palette_spacing`: `none`
- `framework_conversion`: 既存 Job Setup の Svelte、CSS、class、画面構造を土台に、差し替え対象区画だけを frontend 実装へ変換する。
- `prototype_server_url`: `http://127.0.0.1:34116/prototype`
- `prototype_server_command`: `npm --prefix frontend run dev:prototype -- --task translation-job-setup-phase-provider-settings --port 34116`
- `human_review_server_required`: `yes`
- `mock_data_root`: `N/A`
- `mock_data_migration`: `forbidden`
- `sample_data_root`: `[data-ui-prototype-sample-data-root]`
- `sample_data_migration`: `forbidden`
- `production_reference_direction`: `product_code_must_not_reference_ui_prototype`
- `sample_data_notice`: `prototype.svelte` 内のサンプル値は `data-ui-prototype-sample-data-root` 範囲の状態表示確認用データであり、product code、fixture、default state、test data へ移植してはいけない。
  - `interaction_review`: AIサービス変更、APIキー登録モーダル、OS 認証警告、モデル一覧更新、モデル選択、一括処理切り替え、翻訳ジョブ作成を prototype で確認する。
  - `state_transition_review`: APIキー未設定、APIキー登録後の再更新待ち、LM Studio、モデル一覧取得失敗、設定変更後の model 未選択、作成後の設定内容を prototype で確認する。
- `wording_review`: 内部語を画面に出さず、日本語の次操作文へ置き換える。
- `structure_to_preserve`: 既存 Job Setup の shell、card、grid、summary、field、tag、button の構造を維持しつつ、画面全体を UX 標準の判断順へ再構成する。3 つの翻訳段階は desktop では横比較できる設定ブロックとして並べ、mobile では縦順へ戻す。
- `allowed_changes_during_conversion`: 実 view model 名への変換、既存 controller 操作への接続、既存 Job Setup CSS への細部調整。
- `forbidden_changes_during_conversion`: master-persona の AI 設定への依存復活、手動モデル入力欄の追加、LM Studio への APIキー警告追加、Gemini と xAI 以外への一括処理チェックボックス追加。

## Agent Browser Review

- `command_source`: `agent-browser`
- `served_url`: `http://127.0.0.1:34116/prototype`
- `server_command`: `npm --prefix frontend run dev:prototype -- --task translation-job-setup-phase-provider-settings --port 34116`
- `server_status_during_human_review`: `running on 34116`
- `server_command_result`: 指定 command で `http://127.0.0.1:34116/prototype` を起動できた。
- `actual_review_server`: 補助ポートは使っていない。
- `mock_data_refs`: `[data-ui-prototype-sample-data-root]`
- `used_only_for_display_state_review`: `yes`
- `migration_to_product_code`: `forbidden`
- `migration_to_fixture_or_test_data`: `forbidden`
- `checked_viewports`: `desktop 1440x1000 after modal backdrop blur`, `mobile 390x844 after modal backdrop blur`
- `screenshot_refs`: `/private/tmp/tjspps-modal-blur-desktop.png`, `/private/tmp/tjspps-modal-blur-mobile.png`
- `existing_resource_reuse_check`: `job-setup-shell`, `job-setup-card`, `hero-card`, `content-grid`, `summary-grid`, `section-head`, `detail-grid`, `status-pill`, `field-block` の利用を DOM と source で確認した。
- `new_visual_system_check`: 独自の `prototype-shell`, `phase-grid`, `phase-panel`, `page-header`, `create-section`, `summary-section` が残っていないことを DOM と source で確認した。
- `api_key_modal_check`: xAI 未設定状態で「APIキーを登録」ボタンを押すと、同じ画面内に `role=dialog` のモーダルが開くことを確認した。
- `os_auth_warning_check`: モーダル内に「OS認証」と「OS の安全な保管場所へ保存するためです。」が表示されることを確認した。
- `api_key_secret_review`: サンプル APIキー文字列が `document.body.innerText` に出ないことを確認した。
- `ux_review_points`:
  - `goal_completion`: 3 翻訳段階の設定、未設定 APIキーによる停止、未設定なしの翻訳ジョブ作成、作成後の設定内容表示を確認した。
  - `information_priority`: 画面冒頭に目的と説明文、中央に各翻訳段階、下部に固定の作成前確認が並ぶことを snapshot で確認した。
  - `ux_standard_overhaul`: 画面目的、入力確認、共通基盤、フェーズ別設定、下部固定の作成前確認と次操作の順に再構成した。各フェーズ内は AIサービス、モデル、実行方法、状態文へ分割した。
  - `human_review_layout_update`: 共通基盤の説明文から不要な AI 設定注記を削除した。3 つの翻訳段階は desktop で横比較できる配置にし、共通辞書と共通ペルソナは高さ上限付きの別テーブルにした。
  - `phase_density_update`: 3 つの翻訳段階カードは同じ高さに揃えた。モデル以外の小区画は余白と説明文を減らし、モデル選択区画を主な判断対象として見えるようにした。
  - `separate_table_check`: desktop 1440x1000 で共通辞書と共通ペルソナが 2 つの別テーブルになり、3 つの翻訳段階が同じ高さで 3 列に並ぶことを確認した。mobile 390x844 では共通基盤と翻訳段階が 1 列になり、横スクロールがないことを確認した。
  - `desktop_compact_shell_check`: desktop 1440x1000 でモデルカード自体の高さを維持し、周辺セクションの余白、見出し、表、確認区画だけを compact にした。mobile 390x844 では 1 列配置と横スクロールなしを確認した。
  - `validation_global_next_update`: 作成前確認カードを削除した。未設定理由と「次へ」はページ全体の下部固定バーへ移し、画面説明文は「AIサービスを設定します」に寄せた。
  - `hero_simplification_update`: 最上部のヒーローから設定状況、次操作、作成までの流れを削除し、画面説明だけに絞った。
  - `modal_backdrop_blur_update`: APIキー登録モーダルの背景にブラーと濃い暗幕を追加し、背面の透けで本文が読みにくくならないようにした。
  - `api_key_request_wording_update`: APIキー説明は翻訳ジョブ作成だけに限定せず、AIサービスへリクエストを送るために使う鍵として説明した。
  - `operation_order`: AIサービス、APIキー状態、モデル一覧、モデル選択、一括処理、設定済み状態の順で読めることを確認した。
  - `state_comprehension`: 設定済み、未設定あり、APIキー未設定、モデル未選択、モデル一覧取得失敗を文言で区別できることを確認した。
  - `recovery_path`: xAI の APIキー未設定では、モデル一覧更新が無効になり、APIキー登録モーダルを開けることを確認した。
  - `display_wording`: 禁止内部語が `document.body.innerText` に出ないことを確認した。
  - `input_effort`: APIキー未設定時とモデル一覧取得失敗時に手動モデル入力欄がなく、モデル一覧から選ぶ構造を確認した。
  - `eye_movement`: desktop と mobile の snapshot で主要操作が翻訳段階ごとにまとまっていることを確認した。
  - `responsive_continuity`: mobile 390x844 で主要操作が縦順に残ることを snapshot で確認した。
- `structure_ux_review_points`:
  - `primary_action_strength`: 作成ボタンは `button-primary`、モデル一覧更新は `button-secondary` として強弱を分けた。
  - `responsibility_split`: 既存 Job Setup の入力、共通基盤、確認、作成結果は維持し、AI 設定だけを差し替えた。
  - `allowed_and_blocked_actions`: APIキー未設定ではモデル一覧更新が無効になり、APIキー登録だけが可能で、作成ボタンも無効のままになることを確認した。
  - `credential_modal`: APIキー保存後に登録済み状態へ変わり、モデル一覧更新が押せることを確認した。
  - `async_states`: 未更新、取得済み、取得失敗、APIキー未設定、設定済みを表示文言で確認した。
  - `completion_path`: 3 段階の不足解消後に作成ボタンが有効になり、作成後の設定内容へ切り替わることを確認した。
- `wording_review`:
  - `review_timing`: `after_agent_browser_review`
  - `fixed_names_preserved`: Gemini、xAI、OpenAI、LM Studio のサービス名だけを固定名として残した。
  - `business_japanese_terms`: `provider` は AIサービス、`credential` は APIキー状態、`getModels` はモデル一覧を更新へ置き換えた。
  - `internal_state_names_hidden`: `provider`, `credential`, `getModels`, `dirty-validation`, `Ready job summary` は画面本文へ出ていない。
  - `next_action_wording`: APIキー未設定時は「APIキーを登録」だけを表示し、OS 認証画面が出る可能性と保存理由はモーダル内だけで説明する。
  - `batch_help_wording`: 上側の状態表示から一括処理の重複文言を削除し、チェックボックス横の `?` だけに説明を残した。ホバーまたはフォーカス時に画面内ツールチップで「バッチAPIを使うと安く済ませられることがあります。」を読めるようにした。
  - `api_key_modal_actions`: APIキー登録モーダルの閉じる操作は「キャンセル」だけに統一した。
  - `allowed_english_labels`: AIサービスの固有名だけを英語表記として許可する。
  - `plain_language_next_action_judgement`: 非エンジニアの日本語話者が次の操作を読める水準として pass。
- `console_errors`: `none`
- `screenshot_or_snapshot_refs`: `/private/tmp/tjspps-34116-desktop.png`, `/private/tmp/tjspps-34116-mobile.png`, `/private/tmp/tjspps-34116-created.png`, `/private/tmp/tjspps-batch-tooltip-modal-actions.png`, `/private/tmp/tjspps-ux-overhaul-desktop.png`, `/private/tmp/tjspps-ux-overhaul-mobile.png`, `/private/tmp/tjspps-ux-overhaul-created.png`, `agent-browser snapshot -i --compact --depth 5`
- `layout_breaks`: UX-standard overhaul 後も desktop と mobile で横スクロールなし。accessibility snapshot 上の重なりはなし。
- `interaction_detail`: agent-browser の snapshot / errors / screenshot に加え、APIキー登録、モデル一覧更新、作成ボタン有効化、作成後表示を確認した。
- `failure_state_review`: 本文翻訳で Gemini に切り替え、モデル一覧更新失敗時に「手入力欄はありません。」が表示されないこと、text input が存在しないこと、作成不可理由に「本文翻訳のモデルを選んでください」が出ることを確認した。
- `open_issues`: `none`
- `not_checked_reason`: `none`

## Rules

- UI は `ui-design.md` で固定する。
- UIプロトタイプは task-local 確認用であり、docs 正本へ昇格しない。
- UIプロトタイプは `docs/exec-plans/active/translation-job-setup-phase-provider-settings/prototype.svelte` を正本配置にする。
- 本番コードから UIプロトタイプを参照しない。
- UIプロトタイプ内のサンプル値は `data-ui-prototype-sample-data-root` 範囲だけで扱う。
- UIプロトタイプ内のサンプル値は product code、fixture、default state、test data へ移植してはいけない。
- 旧 `ui-mock.html` 方式へ戻さない。
- 旧 `npm run dev:ui-mock` を使わない。
