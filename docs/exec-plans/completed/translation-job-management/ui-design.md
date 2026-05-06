# UI Design: translation-job-management

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
  - incomplete job list: 入力データ名、入力出自、job 状態、現在フェーズ、進捗、最終更新要約。
  - selected job summary: 選択中 job ID、入力出自、現在フェーズ、入力キャッシュ状態、AI 設定要約、再開不可理由。
  - operation availability: 停止可否、再開可否、削除可否、無効理由。
  - result and error summary: 停止要求結果、削除結果、一覧読み込み失敗、参照不能、進捗集約不能。
  - protected information: credential は参照状態だけを出し、API key 平文や provider 応答原文は出さない。
  - deferred boundary notes: API 設定の現在確認と外部通信の停止方式は、後続タスクの境界として表示文言に混ぜない。
- `primary_actions`:
  - 翻訳管理を開く。
  - 未完了 job を選択して Job Run の表示対象にする。
  - Running job の停止入口を開く。
  - Paused または RecoverableFailed job の再開入口を開く。
  - 非実行中 job の削除確認を開く。
- `button_enablement`:
  - 停止は Running job だけ有効にする。
  - 削除は Running job では無効にし、停止後再判定が必要な理由を出す。
  - 再開は Paused または RecoverableFailed のうち、入力キャッシュ欠落や terminal state（完了扱いまたは回復不能な状態）でない時だけ有効候補にする。
  - 参照不能、読み込み失敗、進捗集約不能では危険操作を無効にする。
  - Completed job 入口は本 task の導線に含めない。
  - 削除確認は、job 以下の DB 情報を削除し、入力データと抽出 JSON は残すことを表示する。
  - 再開入口は入力キャッシュ欠落、terminal state（完了扱いまたは回復不能な状態）、状態不整合を無効理由にする。
  - 再開入口は、API 設定が今も使えるかの確認結果までは管理画面で必須表示にしない。
  - 停止入口は Running job だけ有効にし、停止要求中と停止失敗を表示する。
  - 停止入口は、外部通信を実際にどう止めたかを完了条件として表示しない。
- `state_variants`:
  - loading、empty、error、disabled、progress、retry、success を持つ。
  - Ready、Running、Paused、RecoverableFailed、Failed、Canceled、Completed excluded を区別する。
  - cache missing、terminal state（完了扱いまたは回復不能な状態）、state projection inconsistent、stale selection、list load failure を区別する。
  - credential available、credential missing、credential inaccessible は値を出さず状態だけ表示する。
  - paid な real AI API を UI 確認の前提にしない。
- `post_implementation_review`:
  - 一覧と選択詳細で同じ job ID、状態、現在フェーズが読めるか確認する。
  - 停止、再開、削除の無効理由が専門知識なしで分かるか確認する。
  - 危険操作の削除が通常操作や再開入口と近すぎないか確認する。
  - 長い plugin 名、file path、phase 名、provider / model 名、理由文が overflow しないか確認する。
  - API key 平文、credential 値、provider 応答原文が UI、console、error summary に出ないか確認する。

## Interface Frame

- `purpose`: 未完了ジョブを俯瞰し、選択中 job の実行表示、停止、再開、削除可否、再開不可理由を判断できるようにする。
- `audience`: Skyrim Mod 翻訳作業者。
- `primary_workflow`: 翻訳管理を開き、未完了 job を比較し、選択中 job を Job Run に表示して次操作を判断する。
- `information_density`: 作業用画面として、一覧、選択詳細、操作可否を同時に追える密度にする。
- `visual_direction`: `docs/screen-design/design-system-ethereal-archive.md` の翻訳ジョブ表示の方向を保ちつつ、装飾より状態識別、無効理由、危険操作分離を優先する。
- `remembered_signal`: `docs/screen-design/code.html` には翻訳ジョブの過去 mock があるが、`data-design-status=placeholder` を含むため実装正本として流用しない。

## Structure Notes

- `page_sections`:
  - header: 翻訳管理の現在地、未完了 job 件数、再読込。
  - incomplete job list: 未完了 job の一覧、検索または状態フィルタ、状態と進捗。
  - selected job detail: 入力出自、現在フェーズ、入力キャッシュ状態、再開不可理由、AI 設定要約。
  - operation panel: 停止、再開、削除、無効理由、確認結果。
  - feedback area: 読み込み失敗、stale selection、進捗集約不能、削除失敗、停止失敗。
- `layout_constraints`:
  - desktop は一覧を主列、選択詳細と操作 panel を補助列に置く。
  - mobile は一覧、選択詳細、操作 panel の順で縦積みにする。
  - 削除は通常の停止、再開、表示導線から離して配置する。
  - cards の入れ子を避け、一覧行と詳細 section を分ける。
  - status badge だけに依存せず、状態説明 text を併記する。
- `responsive_constraints`:
  - file path、plugin 名、provider / model 名、reason text は折り返しまたは省略表示を持つ。
  - 操作ボタンは mobile で横はみ出しせず、危険操作を独立行または確認領域へ分ける。
  - 進捗は数値と短い説明を併記し、バーだけに依存しない。
  - 最小幅でも job 名、状態、主要操作可否を失わない。
- `accessibility_constraints`:
  - 一覧選択、再読込、停止、再開、削除確認、再試行は keyboard 操作可能にする。
  - error と warning は色だけで伝えない。
  - progress、state、disabled reason は screen reader が読める text を持つ。
  - 危険操作の削除は確認前に影響範囲を読み上げ可能にする。

## UX Standard Review

- `source`: `docs/UX-standard.md`
- `screen_structure_high_priority_results`:
  - No.1 画面目的: 未完了 job 管理に固定する。Completed 履歴や成果物確認は混ぜない。
  - No.5 情報階層: 入力出自、状態、現在フェーズ、進捗、操作可否の順で判断できるようにする。
  - No.7 状態表示: Ready、Running、Paused、RecoverableFailed、Failed、Canceled を表示文言へ変換する。
  - No.8 状態別操作: 停止、再開、削除の有効条件と無効理由を state ごとに持つ。
  - No.19 危険操作: 削除は通常操作と分離し、入力データが残ることを確認できる。
- `screen_structure_applicable_results`:
  - No.12 空状態: 未完了 job がない時は、何も管理対象がないことと次に見る場所を表示する。
  - No.14 エラー状態: 読み込み失敗、参照不能、集約不能を空状態と分ける。
  - No.21 一覧構造: 初回は検索、状態フィルタ、再読込を優先し、一括操作は含めない。
  - No.22 詳細構造: 選択中 job の現在状態、操作可否、理由、設定要約を分ける。
  - No.29 UI文言: 内部状態名だけを出さず、日本語の業務語を併記する。
- `layout_responsive_high_priority_results`:
  - No.4 危険操作位置: 削除は停止、再開と近接しすぎない。
  - No.8 操作対象との対応: 行操作と選択詳細の操作対象 job ID を一致表示する。
  - No.25 ブレークポイント: desktop 2 カラム、mobile 1 カラムへ切り替える。
  - No.34 長文耐性: file path、plugin 名、provider / model 名、reason text を破綻させない。
  - No.38 誤操作防止: 停止、再開、削除の隣接と連打を避ける。
- `layout_responsive_applicable_results`:
  - No.13 ラベルと値: job 状態、現在フェーズ、入力出自の対応を崩さない。
  - No.15 状態の目立ち方: terminal state（完了扱いまたは回復不能な状態）、cache missing、projection inconsistent を通常状態より強く示す。
  - No.28 テーブル崩し: mobile では一覧をカード化し、主要値を縦に積む。
  - No.30 CTA配置変形: mobile では操作 panel を選択詳細の直下に置く。
  - No.37 タップ領域: 状態別操作のボタン間隔を確保する。
- `deferred_items`:
  - 実装後 agent-browser で desktop と mobile の重なり、overflow、console error を確認する。
  - 同じ入力データから複数 job を作成できるため、実装時は job ID、作成日時、状態で見分けられる表示にする。
  - 削除確認では、job 以下の DB 情報が削除対象で、入力データと抽出 JSON が残ることを確認できる文言にする。
  - 停止中表示では、削除可否を再判定するまで待つ状態を明示する。
  - API 設定の現在確認は、再開実行側の後続タスクで扱う。
  - 外部通信の止め方、遅延応答の扱い、停止要求後の不整合防止は、翻訳実行側の後続タスクで扱う。
  - 停止機能が現状ない場合は、翻訳実行側の `tasks/` タスク化が必要である。

## Interaction States

- `loading`:
  - 未完了 job 一覧、選択中 job detail、操作結果の読み込みを分ける。
  - loading 中は停止、再開、削除を disabled にする。
- `empty`:
  - 未完了 job がないことを表示する。
  - 読み込み失敗と空状態を混同しない。
- `error`:
  - list load failure、stale selection、state projection inconsistent、stop failure、delete failure、cache missing を区別する。
  - secret 値、provider 応答原文、過剰な入力本文を error message に含めない。
- `disabled`:
  - Running では削除を disabled にする。
  - Running 以外では停止を disabled にする。
  - cache missing と terminal state（完了扱いまたは回復不能な状態）では再開を disabled にする。
  - disabled 理由は短い説明で表示する。
- `progress`:
  - 停止要求中、削除中、再読込中を区別する。
  - 停止要求中、停止失敗、停止完了を区別する。
- `retry`:
  - 読み込み失敗、参照不能、停止失敗、削除失敗では再読込または再試行を表示する。
  - 再試行前に対象 job ID が変わっていないか確認できるようにする。
- `success`:
  - job 選択成功、停止要求結果、削除結果、再読込成功を短く表示する。
  - 削除成功後は job が一覧から外れ、input data が残ることを表示する。

## Wording Review

- `fixed_names_preserved`:
  - 状態値として Ready、Running、Paused、RecoverableFailed、Failed、Canceled、Completed を設計内固定名として残す。
  - 画面では固定名だけにせず、日本語の状態説明を併記する。
- `business_japanese_terms`:
  - `Ready`: 実行前。
  - `Running`: 実行中。
  - `Paused`: 中断中。
  - `RecoverableFailed`: 再開可能な失敗。
  - `Failed`: 回復できない失敗。
  - `Canceled`: キャンセル済み。
  - `cache missing`: 入力キャッシュがありません。
  - `state projection inconsistent`: 進捗を確認できません。
- `internal_state_names_hidden`:
  - `JOB_PHASE_RUN` や repository 名は画面に出さない。
  - `credential missing` は `APIキーを確認してください` のように次操作へ変換する。
- `next_action_wording`:
  - 停止できません。実行中ではありません。
  - 削除できません。実行中のため、先に停止してください。
  - 再開できません。入力キャッシュを再構築してください。
  - 進捗を確認できません。再読込してください。
- `allowed_english_labels`:
  - provider、model、Batch API など、設定画面で既に見る固定名だけ許容する。

## Componentization Notes

- `componentization_targets`:
  - Job state badge: 状態値と日本語説明を表示する。
  - Job progress summary: 現在フェーズ、進捗、最終更新をまとめる。
  - Operation availability panel: 停止、再開、削除の有効条件と理由をまとめる。
  - Reason message: cache missing、terminal state（完了扱いまたは回復不能な状態）、projection failure の理由表示を統一する。
  - Protected setting summary: provider、model、credential 参照状態を secret なしで表示する。
- `placement`:
  - 画面専用部品を優先し、複数画面で再利用される状態 badge と reason message だけ共有候補にする。
  - 具体的な product component 名と file path は implementation-scope で必要な時だけ固定する。
- `do_not_split`:
  - 未完了 job 管理の page shell は独立部品化しない。
  - 一覧、選択詳細、操作 panel の全状態を 1 つの巨大コンポーネントに閉じ込めない。
- `rationale`:
  - 状態 badge と reason message は他画面でも意味が独立しやすい。
  - 操作 availability は job state と UI 操作の関係を閉じ込められる。

## Post Implementation Review

- `desktop_review_points`:
  - 一覧、選択詳細、操作 panel が同時に読める。
  - Running の削除不可、Paused の再開、cache missing の再開不可が迷わず分かる。
  - 危険操作の削除が通常操作から分離されている。
  - API key 平文が画面、console、error summary に出ない。
- `mobile_review_points`:
  - job list が縦積みになっても job 名、状態、現在フェーズ、操作可否が読める。
  - 選択詳細と操作 panel の対応 job が分かる。
  - 長い理由文が画面外へはみ出さない。
  - 削除確認が誤タップしにくい。
- `overflow_risks`:
  - plugin 名、file path、provider / model 名、phase 名、再開不可理由、削除影響説明。
  - 複数の理由が同時にある場合の selected job detail。
  - job ID、input ID、provider capability（API が必要な実行方式に対応しているか）の要約。
- `visual_polish_open_questions`:
  - terminal state（完了扱いまたは回復不能な状態）の強調度。
  - 停止要求中の一時状態表示。
  - 削除確認の表示形式。
  - Completed job 入口が別画面になる場合の誘導文言。

## Agent Browser Review

- `command_source`: `agent-browser`
- `checked_url`: `not_checked`
- `checked_viewports`: `not_checked`
- `ux_standard_review`:
  - `source`: `docs/UX-standard.md`
  - `high_priority_results`: `UX Standard Review` に設計時確認として記録済み。
  - `applicable_results`: `UX Standard Review` に設計時確認として記録済み。
  - `deferred_items`: 実装後に desktop / mobile / console / screenshot を確認する。
- `wording_review`:
  - `review_timing`: `design_time_before_agent_browser`
  - `fixed_names_preserved`: 固定状態名は設計内に残し、画面文言では日本語説明を併記する。
  - `business_japanese_terms`: `Wording Review` に記録済み。
  - `internal_state_names_hidden`: repository 名、phase run 内部名、secret store 内部名は画面に出さない。
  - `next_action_wording`: 再読込、再構築、停止、削除確認など、次操作が分かる文にする。
  - `allowed_english_labels`: provider、model、Batch API など既存設定語だけ許容する。
  - `plain_language_next_action_judgement`: draft pass。実装後 agent-browser 確認で再判定する。
- `console_errors`: `not_checked`
- `screenshot_or_snapshot_refs`: `none`
- `layout_breaks`: `not_checked`
- `ambiguous_interactions`: `none`
- `open_issues`: human review 承認前である。実装後に削除確認、停止中表示、再開不可理由の文言を agent-browser で再確認する。
- `deferred_design_boundaries`: Q-TJM-003 は再開実行側へ送る。Q-TJM-004 は翻訳実行側へ送り、後続の `tasks/` タスク化を必要事項として扱う。
- `not_checked_reason`: 未実装画面の UI 要件契約作成段階であり、実画面確認は UI 設計根拠ではなく実装後確認観点に回す。

## Source Gaps

- `docs/screen-design/incomplete-job-list.md`: 存在しない。
- `docs/screen-design/app-shell.md`: 存在しない。
- `docs/screen-design/job-run.md`: 存在しない。
- `docs/screen-design/code.html`: 翻訳ジョブ mock はあるが placeholder を含むため正本ではない。
- `docs/scenario-tests/translation-job-management.md`: 未作成。human review 後の昇格候補である。
- `docs/detail-specs/translation-job-management.md`: 未作成。human 承認済み UI 要件だけが将来の昇格候補である。

## Rules

- UI は実装前の mock ではなく、実装が満たす要件契約として書く。
- 独自の page shell、card、grid、配色、余白体系を新規に作らない。
- `docs/screen-design/code.html` の placeholder 文言はそのまま流用しない。
- 細かな visual polish は実装後に人間が実物を確認して直す。
- product component 名や owned scope は、implementation-scope で必要な時だけ扱う。
- implementation-scope の `owned_scope` や product code 対象 file は書かない。
