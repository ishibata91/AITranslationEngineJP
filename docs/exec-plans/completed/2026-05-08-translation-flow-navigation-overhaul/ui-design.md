# UI Design: 2026-05-08-translation-flow-navigation-overhaul

- `skill`: ui-design
- `status`: draft-human-review-ready
- `source_plan`: `./plan.md`
- `scenario_source`: `./scenario-design.md`
- `ux_standard_source`: `docs/UX-standard.md`
- `human_review_ready`: `true`
- `stop_reason`: `none`
- `review_note`: 人間レビュー前の draft であり、承認済みではない。

## UI Contract

- `display_items`:
  - 未完了 job 一覧ページ: 翻訳管理の初期ページ。新規翻訳を開始、Ready、Running、Paused、RecoverableFailed、Failed、Canceled の job、現在フェーズ、進捗、再開可否、再開不可理由、現在の制約説明。
  - 入力データページ: 登録済み入力データ、選択状態、登録状態、`Job Setup` へ進める条件。
  - `Job Setup` ページ: 選択 input、共通基盤、単語翻訳、NPC ペルソナ生成、本文翻訳の AI 設定、作成前検証、作成後 summary。
  - 単語翻訳ページ: job ID、phase state、progress、辞書参照状態、結果 summary、失敗理由、実行操作、次へ進めない理由。
  - NPC ペルソナ生成ページ: job ID、phase state、progress、snapshot 参照状態、body readiness、失敗理由、実行操作、次へ進めない理由。
  - 本文翻訳ページ: job ID、phase state、progress、field result summary、output readiness、失敗理由、実行操作、完了ページへ進めない理由。
  - 翻訳完了ページ: Completed job の原文、訳文、ページング、一覧へ戻る、出力管理へ移動。
  - 出力管理入口: Completed job 一覧、selected job summary、output readiness、拒否理由、`Output Review`。
  - 保護情報: credential は状態分類だけを表示し、API key 平文、secret、provider raw payload、過剰な本文全文は表示しない。
- `primary_actions`:
  - 新規翻訳を開始する。
  - 未完了 job 一覧で job を選んで途中から再開する。
  - input を選び、`Job Setup` へ進む。
  - job を作成し、単語翻訳ページへ進む。
  - 未完了 job 一覧から job を選ぶ。
  - 各フェーズ本文で開始、一時停止、再開、再試行、取消を実行する。
  - `sticky footer` で `次へ進む`、`一覧へ戻る`、`出力管理へ移動` を使う。
  - 出力管理側で Completed job を選ぶ。
- `button_enablement`:
  - job 作成は 3 フェーズ分の API key 状態と model 選択が作成可能な時だけ有効にする。
  - 未完了 job 一覧の再開は参照可能で phase progress を集約できる job だけ有効にする。
  - 単語翻訳から次へ進む操作は、単語翻訳 Completed とジョブ内辞書参照成立時だけ有効にする。
  - NPC ペルソナ生成から次へ進む操作は、persona phase Completed と snapshot 参照成立時だけ有効にする。
  - 本文翻訳から翻訳完了ページへ進む操作は、body phase Completed、field result 整合、output status 整合時だけ有効にする。
  - `Canceled`、`Failed`、RecoverableFailed、field result 不整合、status 不整合では出力管理導線を無効または非表示にし、理由を表示する。
  - 出力管理の XML 出力は、出力管理側で Completed job を選び、output readiness true、row validation pass、出力先 path valid の時だけ有効にする。
  - `sticky footer` は phase 実行、retry、cancel を起動しない。
- `state_variants`:
  - loading、empty、error、disabled、progress、retry、success を持つ。
  - job state は Ready、Running、Paused、RecoverableFailed、Failed、Canceled、Completed excluded を区別する。
  - phase state は not-ready、ready、running、paused、recoverable failed、completed、canceled、failed を区別する。
  - 参照不能、辞書参照不能、snapshot 参照不能、field result 不整合、output status 不整合、stale selection を区別する。
  - 直リンク防止状態は、未完了 job 一覧への復帰と短い理由で表す。
- `sticky_footer_responsibility`:
  - フェーズページ間の移動判断だけを持つ。
  - `次へ進む`、`一覧へ戻る`、`出力管理へ移動` を画面状態に応じて表示する。
  - 次へ進めない理由を既存 sticky footer のエラー表示方式で近接表示する。
  - 実行、一時停止、再開、再試行、取消、XML 出力、preview、再出力は持たない。
- `direct_link_recovery`:
  - 対象 job 未確定でフェーズページへ入った場合は、phase summary 取得を始めず未完了 job 一覧へ戻す。
  - 参照不能 job はフェーズページへ進めず、未完了 job 一覧に理由を表示する。
  - 復帰先は入力データページや `Job Setup` ではない。
- `post_implementation_review`:
  - 実画面で翻訳管理入口、各フェーズページ、翻訳完了ページ、出力管理を確認する。
  - desktop と mobile で footer、一覧、phase summary、長文理由が重ならないことを確認する。
  - 旧 `Job Run` のセッション取得、フェーズ直リンク、前工程戻り導線が残っていないことを確認する。
  - `Canceled` と `Failed` が翻訳完了ページに入らないことを確認する。
  - 出力管理へ移動後に job が自動選択されず、completed job list から選ぶことを確認する。
  - console error と secret 平文表示がないことを確認する。

## Interface Frame

- `purpose`: 翻訳 job の作成、再開、フェーズ実行、結果確認、成果物出力への移動を、状態整合性を崩さない導線へ整理する。
- `audience`: Skyrim Mod 翻訳作業者。
- `primary_workflow`: 翻訳管理を開くと未完了 job 一覧ページを初期表示し、そのページから新規作成または未完了 job 選択を行い、フェーズページを順に進め、Completed job の結果を確認して出力管理へ移動する。
- `information_density`: 作業用画面として、現在 job、現在フェーズ、操作可否、進めない理由を同時に追える密度にする。
- `visual_direction`: 既存のアプリシェル、翻訳管理、出力管理の構造を土台にする。独自の page shell、card、grid、配色、余白体系は作らない。
- `remembered_signal`: 既存実画面では翻訳管理の `実行` タブに旧 `Job Run` と 3 フェーズが縦に並ぶ。新設計では 旧 `Job Run` 大箱を残さず、各フェーズを翻訳セクション直下のページとして扱う。

## Structure Notes

- `page_sections`:
  - app shell: グローバルナビは `翻訳管理` と `出力管理` だけを主要入口にする。
  - incomplete job list page: 翻訳管理の初期ページ。新規開始導線、Completed 以外の job 一覧、再開可否。
  - input data page: 登録済み入力データの選択と `Job Setup` への条件。
  - job setup page: 3 フェーズの AI 設定と job 作成。
  - phase pages: 単語翻訳、NPC ペルソナ生成、本文翻訳を同列ページとして表示。
  - translation complete page: Completed job の結果確認。
  - output management: Completed job 一覧と Output Review。
- `layout_constraints`:
  - 各フェーズページは上部に job summary と phase state を置く。
  - フェーズ実行操作はページ本文の操作領域に置く。
  - `sticky footer` は移動操作と進めない理由に限定する。
  - フェーズページから入力データページや `Job Setup` へ戻る button、tab、breadcrumb を表示しない。
  - 翻訳完了ページと Output Review の操作責務を混ぜない。
  - cards の入れ子を避け、一覧、summary、操作、失敗情報を区画として分ける。
- `responsive_constraints`:
  - desktop は summary、進捗、操作を 2 カラムまたは明確な縦区画で並べる。
  - mobile は job summary、phase result、操作、失敗情報、footer の順に縦積みにする。
  - sticky footer は mobile で本文と重ならない余白を持つ。
  - 長い job ID、plugin 名、file path、source text、translated text、model 名、error reason は折り返しまたは省略と詳細展開を持つ。
  - 一覧は mobile でカード化し、状態、現在フェーズ、操作可否、理由を失わない。
- `accessibility_constraints`:
  - 状態は色だけでなく label と icon を併用する。
  - disabled button の理由は tooltip または隣接文言で読める。
  - sticky footer のエラーは screen reader が読める text を持つ。
  - グローバルナビ、一覧選択、footer 操作、フェーズ本文操作は keyboard 操作可能にする。
  - 危険操作の cancel は通常移動と離して置く。

## UX Standard Review

- `source`: `docs/UX-standard.md`
- `screen_structure_high_priority_results`:
  - No.1 画面目的: 各ページの主目的を未完了一覧、入力、設定、フェーズ実行、完了確認、出力管理に分ける。
  - No.5 情報階層: job、phase、readiness、次操作、理由の順で判断できるようにする。
  - No.6 画面責務: 旧 `Job Run` の大箱を分解し、翻訳完了ページと Output Review を混ぜない。
  - No.8 状態別操作: job state、phase state、readiness ごとの有効条件を定義する。
  - No.10 禁止条件: 直リンク、前工程戻り、未選択 job 出力、Completed 以外の出力を禁止条件として明示する。
  - No.14 エラー状態: 参照不能、辞書参照不能、snapshot 参照不能、不整合を空状態と分ける。
  - No.24 画面遷移: 業務フロー順の `次へ進む` と一覧再開に限定する。
  - No.29 UI文言: 内部状態名だけを出さず、利用者が次に何をするか分かる日本語へ変換する。
  - No.32 検査可能性: UI 状態、操作、禁止条件を scenario と対応づける。
- `screen_structure_applicable_results`:
  - No.12 空状態: 未完了 job なし、Completed job なし、summary 未取得を分ける。
  - No.19 危険操作: cancel は phase 本文の危険操作として分離する。
  - No.20 完了後導線: 翻訳完了ページは出力管理への移動だけを示す。
  - No.21 一覧構造: 未完了 job 一覧と Completed job 一覧を別セクションに分ける。
  - No.25 現在地表示: ページタイトルと現在フェーズで現在地を示す。
- `layout_responsive_high_priority_results`:
  - No.5 関連情報の近接: 次へ進めない理由は footer の `次へ進む` 近くに表示する。
  - No.8 操作対象との対応: フェーズ本文操作と対象 job ID を同一区画で確認できるようにする。
  - No.25 ブレークポイント: desktop と mobile で構造を変える。
  - No.30 CTA配置変形: mobile では footer 操作を下部で維持し、本文に必要な余白を確保する。
  - No.33 固定要素の干渉: sticky footer が本文やエラー文と重ならないようにする。
  - No.34 長文耐性: source、Dest、error reason、model 名、file path を破綻させない。
  - No.38 誤操作防止: 実行操作と移動操作を離す。
- `layout_responsive_applicable_results`:
  - No.13 ラベルと値: job ID、phase state、readiness の対応を崩さない。
  - No.28 テーブル崩し: field result は mobile で key-value list または要約カードにする。
  - No.37 タップ領域: footer button と phase action の間隔を確保する。
  - No.45 エラー表示: footer のエラー表示で layout jump が大きくならないようにする。
  - No.61 証跡: 実装後に desktop と mobile のスクリーンショットを保存する。
- `deferred_items`:
  - 実装後に各状態 fixture で `agent-browser` screenshot を保存する。
  - 実装後に mobile viewport の sticky footer 干渉を確認する。
  - 実装後に Canceled / Failed の表示が翻訳完了ページへ入らないことを確認する。
  - 実装後に 出力管理へ移動後の自動選択なし表示を確認する。

## Interaction States

- `loading`:
  - 未完了 job 一覧、各 phase summary、completed result paging、Completed job 一覧を分ける。
  - loading 中は危険操作と次へ進むを disabled にする。
- `empty`:
  - 未完了 job なし、Completed job なし、summary 未取得、field result なしを分ける。
  - 空状態では次に見る場所を短く示す。
- `error`:
  - job 作成失敗、一覧読み込み失敗、参照不能、phase progress 集約不能、辞書参照不能、snapshot 参照不能、field result 不整合、output status 不整合を区別する。
  - secret、provider raw payload、過剰本文は error に含めない。
- `disabled`:
  - disabled 理由を button 近くまたは footer 近くに出す。
  - `次へ進む` は前 phase readiness 不成立で disabled にする。
  - 出力 action は Completed job 未選択、readiness false、path invalid で disabled にする。
- `progress`:
  - phase running、pause requested、resume requested、retrying、output generating を区別する。
  - progress は数値と状態文を併記する。
- `retry`:
  - RecoverableFailed の retry は phase 本文で扱う。
  - footer は retry を持たない。
  - 読み込み失敗や参照不能は再読込または一覧へ戻るを表示する。
- `success`:
  - job 作成成功、phase Completed、translation complete、output success を分ける。
  - phase Completed は次へ進む readiness を更新する。

## Wording Review

- `fixed_names_preserved`:
  - 固定名として `Job Setup`、`Output Review`、Ready、Running、Paused、RecoverableFailed、Failed、Canceled、Completed を設計内で残す。
  - 画面では固定名だけにせず、日本語説明を併記する。
- `business_japanese_terms`:
  - `Job Run`: 旧実行画面。新 UI では画面名として使わない。
  - `Ready`: 実行前。
  - `Running`: 実行中。
  - `Paused`: 中断中。
  - `RecoverableFailed`: 再開可能な失敗。
  - `Failed`: 回復できない失敗。
  - `Canceled`: キャンセル済み。
  - `output readiness`: 出力準備。
  - `field result`: 翻訳結果。
  - `snapshot`: 参照スナップショット。
- `internal_state_names_hidden`:
  - `JOB_PHASE_RUN`、repository 名、secret store key、route state は画面に出さない。
  - `credential missing` は `APIキーを確認してください` のように次操作へ変換する。
  - `state projection inconsistent` は `進捗を確認できません` と表示する。
- `next_action_wording`:
  - 次へ進めません。単語翻訳の辞書参照を確認してください。
  - 次へ進めません。ペルソナ参照を確認してください。
  - 出力できません。本文翻訳が完了していません。
  - 出力管理で Completed job を選んでください。
  - job を選び直してください。
- `allowed_english_labels`:
  - provider、model、Batch API、Output Review など、既存設定画面または詳細仕様の固定名だけ許容する。
- `plain_language_next_action_judgement`:
  - 利用者が「次に押すもの」と「押せない理由」を専門用語なしで読める水準にする。

## Componentization Notes

- `componentization_targets`:
  - Phase state badge: phase state と日本語説明を表示する。
  - Phase progress summary: progress、target count、result count を表示する。
  - Readiness reason message: 次へ進めない理由を表示する。
  - Sticky navigation footer: 移動操作とブロック理由だけを扱う。
  - Job selection summary: 未完了 job 一覧とフェーズページの対象 job を示す。
  - Completed result pager: 原文と訳文のページング表示を扱う。
  - Protected setting summary: provider、model、credential 状態分類を secret なしで表示する。
- `placement`:
  - 画面専用部品を優先する。
  - 状態 badge、reason message、protected setting summary は共有候補にできる。
  - product component 名と file path は implementation-scope で必要な時だけ固定する。
- `do_not_split`:
  - 翻訳管理全体の page shell は新規独自化しない。
  - 業務フロー全体の進行状態を 1 つの共有部品へ閉じ込めない。
  - `sticky footer` に phase 実行操作を持たせない。
- `rationale`:
  - 状態 badge と reason message は複数ページで意味が独立する。
  - phase page の大きな layout は各ページ固有の判断が多く、無理に共有化しない。

## Post Implementation Review

- `desktop_review_points`:
  - 未完了 job 一覧ページから新規開始と途中再開が分かれる。
  - 各フェーズページで job ID、phase state、進捗、本文操作、footer 操作が重ならない。
  - 旧 `Job Run` 名とセッション取得ボタンが表示されない。
  - 翻訳完了ページで XML 出力や preview が表示されない。
  - 出力管理で Completed job 一覧から選ぶ流れが分かる。
  - secret 平文が画面、console、error summary に出ない。
- `mobile_review_points`:
  - sticky footer が本文、error、keyboard と干渉しない。
  - job list、phase summary、操作、失敗情報が自然な順に縦積みされる。
  - 長い原文、訳文、model 名、file path、error reason がはみ出さない。
  - footer button が横はみ出しせず、無効理由が読める。
- `overflow_risks`:
  - 原文、訳文、EditorID、FormID、plugin 名、file path、provider / model 名、snapshot digest、error reason。
  - 未完了 job 一覧で複数理由が同時に出る場合。
  - sticky footer に長いブロック理由が出る場合。
- `visual_polish_open_questions`:
  - phase state badge の色と icon の強調度。
  - sticky footer の error 表示密度。
  - 翻訳完了ページの原文と訳文の比較表示密度。
  - 出力管理へ移動後の selected job 未選択表示。

## Agent Browser Review

- `command_source`: `agent-browser`
- `app_start_command`: `npm run dev:wails:agent-browser`
- `checked_url`:
  - `http://localhost:34115/#dashboard`
  - `http://localhost:34115/#translation-management`
  - `http://localhost:34115/#output-management`
- `checked_viewports`:
  - desktop default viewport。mobile viewport は未確認。
- `ux_standard_review`:
  - `source`: `docs/UX-standard.md`
  - `high_priority_results`: `UX Standard Review` に設計時確認として記録済み。
  - `applicable_results`: `UX Standard Review` に設計時確認として記録済み。
  - `deferred_items`: mobile viewport、状態別 fixture、実装後スクリーンショット。
- `existing_screen_observations`:
  - dashboard は `翻訳管理` と `出力管理` をグローバルナビに持つ。
  - 翻訳管理は `ジョブ管理`、`データロード`、`セットアップ`、`実行` の tab を持つ。
  - 既存の `実行` tab には `Job Run` と単語翻訳、NPC ペルソナ生成、本文翻訳の情報が同じ大きな表示内に並ぶ。
  - 既存の `Job Run` には job id 入力と `summary 取得` がある。
  - 出力管理は `Output Review`、completed job list、selected job summary、XML 出力、再出力を持つ別ページである。
- `wording_review`:
  - `review_timing`: `after_agent_browser_review`
  - `fixed_names_preserved`: `Job Setup`、`Output Review`、provider、model などの固定名は必要箇所だけ残す。
  - `business_japanese_terms`: 既存画面には `Job Run`、`phase control`、`execution summary` など英語語彙が多い。新 UI では見出しを日本語の業務語へ寄せる。
  - `internal_state_names_hidden`: 新 UI では route state、JOB_PHASE_RUN、credential ref を利用者向け文言へ置き換える。
  - `next_action_wording`: `job id を入力して summary を取得してください` は、新 UI では一覧選択または作成結果から job が固定されるため不要にする。
  - `allowed_english_labels`: provider、model、Batch API、Output Review。
  - `plain_language_next_action_judgement`: 既存画面の `Job Run` と `summary 取得` は、新方針では利用者の次操作を曖昧にするため廃止対象である。
- `console_errors`: none。`agent-browser errors` は出力なし。
- `screenshot_or_snapshot_refs`:
  - `tmp/agent-browser/translation-flow-navigation-translation.png`
  - `tmp/agent-browser/translation-flow-navigation-output.png`
- `layout_breaks`:
  - desktop default では設計対象の旧画面確認に使える表示が得られた。
  - mobile viewport は未確認であり、実装後確認へ送る。
- `ambiguous_interactions`:
  - 既存 `Job Run` の job id 入力と `summary 取得` は、未完了 job 一覧から対象 job を選ぶ新方針と競合する。
  - 既存 `実行` tab の step 表示は、フェーズページ直移動に見えやすい。
- `open_issues`:
  - 実装後に mobile viewport と複数状態 fixture の agent-browser 確認が必要である。
- `not_checked_reason`:
  - mobile viewport と状態別 fixture は、設計根拠ではなく実装後確認観点として扱うため未確認にした。

## Rules

- UI は `ui-design.md` の UI 要件契約で固定する。
- UI 確認は実装後に agent-browser で行う。
- UX 標準の確認結果を `UX Standard Review` に記録する。
- 固定名以外の画面表示文言は、日本語の業務語へ置き換える。
- 内部状態名は画面に出さず、利用者の次操作を示す文へ変換する。
- 英語ラベルは、利用者が設定画面で見る既存語だけに限定する。
- 既存画面変更では、既存画面または既存 UI 部品を土台にする。
- 既存画面変更では、独自の page shell、card、grid、配色、余白体系を新規に作らない。
- 既存画面変更では、変更対象区画だけを差し替え、変更しない区画は既存画面の構造と表示を維持する。
- 細かな visual polish は実装後に人間が実物を確認して直す。
- implementation-scope の `owned_scope` や product code 対象 file は書かない。
