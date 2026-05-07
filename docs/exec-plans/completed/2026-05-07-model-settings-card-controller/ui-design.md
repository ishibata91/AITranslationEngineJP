# UI Design: 2026-05-07-model-settings-card-controller

- `skill`: ui-design
- `status`: approved
- `source_plan`: `./plan.md`
- `scenario_source`: `./scenario-design.md`
- `questions_source`: `./scenario-design.questions.md`
- `ux_standard_source`: `docs/UX-standard.md`

## UI Contract

- `display_items`:
  - 利用者向け provider list は `gemini`、`lm_studio`、`xai` だけを表示する。
  - model list 状態は、未更新、更新中、取得済み、取得済み 0 件、取得失敗、APIキー未設定で更新不可を分ける。
  - model 選択状態は、未選択、選択済み、未保存変更、保存中、保存済み、保存失敗を分ける。
  - APIキー状態は、設定済み、未設定、不要を分ける。APIキー本体、secret、raw payload は表示しない。
  - fake mode の内部状態、`fake` provider ID、fake transport の詳細は表示しない。
- `primary_actions`:
  - provider を選ぶ。
  - model list を更新する。
  - model を選ぶ。
  - model 設定を保存する。
  - APIキー未設定時は、共有カード内に AIサービス設定を開く導線を出さず、状態表示だけを確認する。
- `button_enablement`:
  - model list 更新は、APIキーが必要な provider で APIキー未設定の場合は無効にする。
  - model 選択は、現在 provider の model list 取得済み状態でだけ有効にする。
  - 保存は、現在 provider と model が整合し、参照側ごとの保存対象が決まっている場合だけ有効にする。
  - Job Setup の job 作成は、3 翻訳段階すべてで APIキー不足、model 未選択、更新中、取得失敗がない時だけ有効にする。
- `state_variants`:
  - 初期取得中、取得済み、取得失敗、未保存変更あり、保存中、保存済み、保存失敗。
  - model list 未更新、更新中、取得済み、取得済み 0 件、取得失敗、APIキー未設定で更新不可。
  - provider 変更直後、遅延応答破棄後、fake mode 取得結果表示。
- `post_implementation_review`:
  - デスクトップとモバイルで provider、model、状態、操作が近接していることを確認する。
  - 長い model 名、長いエラー文、長い provider 表示がカード幅からはみ出さないことを確認する。
  - fake mode で `fake` provider ID が表示されないことを確認する。
  - APIキー本体と raw payload が UI、要約、エラー文に出ないことを確認する。

## Interface Frame

- `purpose`: マスターペルソナと翻訳ジョブ設定で、provider と model の選択、model list 更新、保存状態を同じ操作規則で扱う。
- `audience`: Skyrim Mod 翻訳 job を作る利用者、共通ペルソナを構築する利用者、人間レビュー担当者。
- `primary_workflow`: provider を選び、model list を更新し、model を選び、必要に応じて保存し、Job Setup では job 作成条件へ反映する。
- `information_density`: カード内に provider、model、APIキー状態、model list 状態、保存状態、主要操作だけを置く。接続詳細と raw payload は出さない。
- `visual_direction`: 既存画面と既存カード構造を維持し、変更対象区画だけを状態制御対応へ差し替える。
- `remembered_signal`: 利用者は「どの AIサービスで、どの model を、保存済みとして使えるか」をカード内で判断できる。

## Structure Notes

- `page_sections`:
  - マスターペルソナ画面のモデル設定カード区画。
  - Job Setup の単語翻訳、NPC ペルソナ生成、本文翻訳の各カード区画。
  - Job Setup の作成前確認と作成後の設定内容。
- `layout_constraints`:
  - provider、model list 更新、model 選択、状態表示、保存操作は同じ判断単位として近くに置く。
  - APIキー未設定や取得失敗の文言は、対象 provider と更新ボタンの近くに置く。
  - 保存済み、未保存、保存失敗の表示は保存操作の近くに置く。
  - 画面固有の job 作成条件は Job Setup 側の確認区画へ残す。
- `responsive_constraints`:
  - 狭い幅では provider、更新操作、model 選択、保存操作の順に縦積みする。
  - 長い model 名とエラー文は折り返し、操作ボタンの幅を押し広げない。
  - 3 翻訳段階カードは同じ順序と同じ状態表示位置を維持する。
- `accessibility_constraints`:
  - エラー、警告、成功、選択状態は色だけに依存しない。
  - 無効ボタンには、近接する説明文で理由を示す。
  - 状態ラベルは短い名詞句にし、説明文は自然な日本語に分ける。

## UI Component Judgment

- `componentized_targets`:
  - `AIModelSelectionCard.svelte` は表示部品として維持する。
  - provider 選択、model list 状態、model 選択、保存状態の表示は共有カードの入力として扱う。
  - 画面固有の job 作成可否や作成後設定内容は Job Setup 側に残す。
- `placement`:
  - 共有表示部品は UI Component 層に置く。
  - 保存、取得、model list 更新、遅延応答破棄は ScreenController、Frontend UseCase、Store 側で扱う。
  - Wails binding と backend DTO の変換は Gateway 側に閉じる。
- `not_componentized`:
  - 保存単位、取得公開口、状態 namespace は UI Component の props だけで解決しない。
  - APIキー保存や secret 管理はモデル設定カードの表示部品へ入れない。
- `reason`:
  - UI Component は backend DTO、generated binding、Store、Gateway を直接扱わないという architecture 正本に従う。
  - 共有カード制御は複数画面で使うが、Job Setup の job 作成条件は画面固有である。

## UX Standard Review

- `source`: `docs/UX-standard.md`
- `screen_structure_high_priority_results`:
  - 状態表示、状態別操作、表示条件、エラー状態、入力制約、変更差分を UI 契約へ反映した。
  - UI 文言は内部状態名を避け、利用者が次操作を判断できる日本語へ寄せる。
  - 検査可能性は scenario matrix と state variants で固定した。
- `screen_structure_applicable_results`:
  - 空状態は、取得済み 0 件を取得失敗と分ける文言にする。
  - 保存失敗後の変更差分は、未保存変更として残す文言にする。
  - APIキー未設定時は、共有カード内導線なしで状態表示だけにする。
- `layout_responsive_high_priority_results`:
  - 関連情報の近接、操作対象との対応、状態の目立ち方、要素サイズ統一、折り返し順、長文耐性を実装後確認観点にした。
  - デスクトップとモバイルで、provider、model list、model、保存、job 作成条件の読む順を維持する。
- `layout_responsive_applicable_results`:
  - 3 翻訳段階カードの同一構造を維持する。
  - 長い model 名、翻訳データ名、エラー文のはみ出しを実装後に確認する。
- `deferred_items`:
  - 幅別スクリーンショットと状態別スクリーンショットは、frontend 実装後の人間確認で扱う。

## Interaction States

- `loading`:
  - 初期取得中、model list 更新中、保存中を分けて表示する。
- `empty`:
  - model list 未更新と取得済み 0 件を分ける。取得済み 0 件では model 未選択として保存と job 作成を拒否する。
- `error`:
  - model list 取得失敗、provider settings 参照不能、保存失敗、endpoint 参照不能を短い日本語要約で表示する。
  - raw payload、secret、内部ログ用識別子は出さない。
- `disabled`:
  - APIキー未設定時の model list 更新、model 未選択時の保存、3 翻訳段階不足時の job 作成を無効にする。
- `progress`:
  - model list 更新中は現在 provider の更新として表示し、provider 変更後は古い応答を反映しない。
- `retry`:
  - model list 取得失敗後は再更新を可能にする。
  - 保存失敗後は未保存変更を残し、保存操作の近くで再試行できるようにする。
- `success`:
  - 取得済み、model 選択済み、保存済み、Job Setup 作成可能を分けて表示する。

## Wording Review

- `fixed_names_preserved`:
  - `gemini`、`lm_studio`、`xai`、`fake-model` は固定名として残す。
  - `fake` provider ID は利用者向け表示へ出さない。
- `business_japanese_terms`:
  - `model list` は画面表示では「モデル一覧」とする。
  - `credential missing` は「APIキーが未設定です」とする。
  - `dirty` は「未保存の変更があります」とする。
- `internal_state_names_hidden`:
  - fake mode、request token、gateway、DTO、Store、adapter は画面表示へ出さない。
- `next_action_wording`:
  - APIキー未設定時は「APIキーが未設定のため、モデル一覧を更新できません」とする。共有カード内に AIサービス設定を開く導線は出さない。
  - model list 取得失敗時は「モデル一覧を取得できませんでした。もう一度更新してください」とする。
- `allowed_english_labels`:
  - provider 固定 ID と model 固定名だけを残す。

## Post Implementation Review

- `desktop_review_points`:
  - Job Setup の 3 翻訳段階でカード構造、状態表示位置、操作位置が揃っている。
  - 作成前確認と作成後の設定内容に、段階ごとの AIサービス、model、APIキー状態、一括処理の有無だけが出る。
- `mobile_review_points`:
  - provider、model list 更新、model 選択、保存操作が縦積みでも意味順に読める。
  - 無効理由と対象ボタンの対応が崩れない。
- `overflow_risks`:
  - 長い model 名、長い翻訳データ名、長いエラー文。
  - provider 名と状態ラベルが横並びになった時の折り返し。
- `visual_polish_open_questions`:
  - なし。細かな visual polish は実装後の実物確認で扱う。

## Agent Browser Review

- `command_source`: `agent-browser`
- `checked_url`: `not-checked`
- `checked_viewports`: []
- `ux_standard_review`:
  - `source`: `docs/UX-standard.md`
  - `high_priority_results`: UI 契約上の確認に限定した。
  - `applicable_results`: 実画面確認は frontend 実装後に行う。
  - `deferred_items`: 幅別スクリーンショット、状態別スクリーンショット、console error。
- `wording_review`:
  - `review_timing`: `before_agent_browser_review`
  - `fixed_names_preserved`: provider ID と `fake-model` だけを固定名として残した。
  - `business_japanese_terms`: model list は「モデル一覧」として扱う。
  - `internal_state_names_hidden`: fake mode、request token、Store、Gateway は画面文言に出さない。
  - `next_action_wording`: APIキー未設定と model list 取得失敗の文言を固定した。
  - `allowed_english_labels`: provider ID と model 名だけ。
  - `plain_language_next_action_judgement`: APIキー未設定時は、既存バナー導線と重複しない状態表示だけにする。
- `console_errors`: not-checked
- `screenshot_or_snapshot_refs`: []
- `layout_breaks`: not-checked
- `ambiguous_interactions`: []
- `open_issues`: []
- `not_checked_reason`: frontend 実装前の UI 契約固定段階であり、実画面確認対象がまだないため。
