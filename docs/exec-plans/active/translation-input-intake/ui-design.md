# UI Design: translation-input-intake

- `skill`: ui-design
- `status`: approved
- `source_plan`: `./plan.md`
- `scenario_source`: `./scenario-design.md`

## UI Contract

- `display_items`:
  - input file 一覧: 表示名、登録状態、登録日時、出自情報、再構築可否。
  - 入力データ概要: 翻訳レコード件数、翻訳フィールド件数、カテゴリ別件数。
  - sample field: RecordType、SubrecordType、FormID、EditorID、原文。
  - error summary: 失敗種別、対象 file、再試行可否。
- `primary_actions`:
  - xEdit 抽出 JSON を登録する。
  - 登録済み入力データを選択する。
  - 取り込み結果を再読込する。
  - 入力キャッシュを再構築する。
  - 失敗した登録または再構築を再試行する。
- `button_enablement`:
  - 登録 action は基盤データ管理が成立し、登録中でない時だけ有効にする。
  - 再構築 action は抽出 JSON 正本の出自情報があり、cache 欠落または再構築可能状態の時だけ有効にする。
  - retry action は失敗状態と retry 可能な error の時だけ有効にする。
  - job 作成、翻訳開始、出力生成の action はこの UI 契約に含めない。
- `state_variants`:
  - loading、empty、progress、success、error、disabled、retry を持つ。
  - duplicate input、invalid JSON、source file missing、cache missing を error variant として区別する。
  - response の warnings、categories、sample fields が null の場合は空配列として扱う。
  - 初回 import では browser file input から bare filename だけが渡る可能性を考慮し、backend へ file content または解決可能な source handle を渡す。
  - app-shell 導線は `dashboard-and-app-shell` 側の責務とし、Input Review はページ内で完結する。
- `post_implementation_review`:
  - input file 一覧、件数、カテゴリ、sample field が 1 画面で追えるか確認する。
  - desktop と mobile 幅で一覧と詳細が重ならないか確認する。
  - file name、path、error message が長い時に overflow しないか確認する。
  - keyboard だけで登録、一覧選択、retry に到達できるか確認する。

## Interface Frame

- `purpose`: xEdit 抽出 JSON の取り込み結果を翻訳開始前に確認する。
- `audience`: Skyrim Mod 翻訳作業者。
- `primary_workflow`: Input Review で JSON を登録し、一覧、件数、カテゴリ、sample field を確認する。
- `information_density`: 作業用画面として、一覧と詳細を同時に確認できる密度にする。
- `visual_direction`: `docs/screen-design/design-system-ethereal-archive.md` の visual direction に従い、装飾より可読性と状態識別を優先する。
- `remembered_signal`: `docs/screen-design/` に `app-shell.md` と `input-review.md` は存在しない。今回の UI 契約は task-local の一時正本である。

## Structure Notes

- `page_sections`:
  - header: Input Review の現在状態と登録 action。
  - input list: 登録済み input file 一覧。
  - selected input summary: 件数、カテゴリ、再構築可否。
  - sample fields: 代表 field の確認。
  - error / activity panel: import、rebuild、retry の結果。
- `layout_constraints`:
  - desktop は list と detail を横並びにできる。
  - mobile は list の下に detail を積む。
  - cards の入れ子を避け、section と repeated item を分ける。
  - 固定高さの一覧領域は overflow 表示を持つ。
- `responsive_constraints`:
  - file path、FormID、EditorID、error message は折り返しまたは省略表示を持つ。
  - 主要 action は mobile で横にはみ出さない。
  - 件数とカテゴリは小幅でも label と値が対応する。
- `accessibility_constraints`:
  - file registration、rebuild、retry は keyboard 操作可能にする。
  - error は色だけでなく text で伝える。
  - progress は screen reader が読める text を持つ。
  - focus 移動は登録後に結果 summary または error summary へ移る。

## Interaction States

- `loading`:
  - input file 一覧の読込中を表示する。
  - 登録 action は disabled にする。
- `empty`:
  - 登録済み input がないことを表示する。
  - xEdit 抽出 JSON の登録 action を有効にする。
- `error`:
  - invalid JSON、non-xEdit JSON、source file missing、cache rebuild failed を区別する。
  - 初回 import request が bare filename だけで content も source handle もない場合は invalid request として扱う。
  - source file missing は cache rebuild 時に保存済み正本が見つからない場合に限定する。
- `disabled`:
  - 基盤データ管理が未成立の場合は登録 action を disabled にする。
  - disabled 理由を画面上に短く表示する。
- `progress`:
  - import 中、parse 中、cache rebuild 中を区別する。
  - 同一 input に対する重複 action を止める。
- `retry`:
  - retry 可能な失敗だけ retry action を表示する。
  - retry 不能な失敗では source file や format の修正を促す。
- `success`:
  - 登録済み input、翻訳レコード件数、翻訳フィールド件数、カテゴリ別件数、sample field を表示する。
  - job 作成や翻訳開始を success state の必須 action にしない。

## Error States

- `invalid_json`: JSON として読めない。登録前拒否か失敗状態作成かは Q-TII-004 の回答待ち。
- `unsupported_extract_shape`: JSON だが xEdit 抽出形式として必要な構造がない。
- `missing_required_field`: 必須 field 欠落。拒否粒度は Q-TII-004 の回答待ち。
- `duplicate_input`: 同一抽出 JSON の再登録。扱いは Q-TII-003 の回答待ち。
- `source_file_missing`: cache 再構築に必要な抽出 JSON 正本が見つからない。初回 import の bare filename では使わない。
- `unknown_field_definition`: 未定義 RecordType + SubrecordType。扱いは Q-TII-006 の回答待ち。

## Post Implementation Review

- `desktop_review_points`:
  - 一覧、概要、sample field が同時に読める。
  - 長い file name と path が layout を壊さない。
  - error と retry の位置が登録 action と近い。
- `mobile_review_points`:
  - input list から detail へ移動しても文脈を失わない。
  - action が折り返されても押し間違えにくい。
  - sample field の横長値が画面外にはみ出さない。
  - browser file input 経由の実ファイル登録で source file missing にならない。
  - null 配列 response で console error が出ない。
- `overflow_risks`:
  - file path、plugin 名、EditorID、error message、category label。
  - record / field 件数が多い時の list scroll。
  - sample field の原文が長い時の高さ増加。
- `visual_polish_open_questions`:
  - app-shell navigation 上の表示名と位置は dashboard-and-app-shell 側で扱う。
  - カテゴリ color や badge の使い方。
  - empty state の説明文の粒度。

## Source Gaps

- `docs/screen-design/app-shell.md`: 存在しない。
- `docs/screen-design/input-review.md`: 存在しない。
- `docs/scenario-tests/translation-input-intake.md`: 未作成。task-local scenario 承認後に昇格候補。
- `docs/detail-specs/input-review.md`: 未作成。human 承認済み UI 要件だけが将来の昇格候補。

## Rules

- UI は実装前の mock ではなく、実装が満たす要件契約として書く。
- 実装前の見た目 artifact を新規必須にしない。
- 細かな visual polish は実装後に人間が実物を確認して直す。
- product component 名や owned scope は、implementation-scope で必要な時だけ扱う。
- implementation-scope の `owned_scope` や product code 対象 file は書かない。
