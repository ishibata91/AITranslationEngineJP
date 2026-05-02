# UI Design: term-translation-phase

- `skill`: ui-design
- `status`: ready-for-human-review
- `source_plan`: `./plan.md`
- `scenario_source`: `./scenario-design.md`

## UI Contract

- `display_items`:
  - Job Run に current phase、phase state、progress、開始時刻、完了時刻、対象語件数、共通辞書 hit 件数、AI 実行対象語件数を表示する。
  - phase result に確定訳語件数、ジョブ内辞書反映件数、置換対象件数、未一致件数、provider / model / execution mode の要約を表示する。
  - エラー時は error kind、短い理由、retryable flag、後続 phase 不可理由を表示する。
  - secret、API key 平文、provider raw request / response、翻訳フィールド本文の全文は表示しない。
- `primary_actions`:
  - 単語翻訳フェーズ開始。
  - 中断、再開、リトライ。
  - phase result 更新。
  - 後続 phase へ進む操作。
- `button_enablement`:
  - 開始は Ready job かつ active な単語翻訳 phase run がない時だけ有効にする。
  - 後続 phase へ進む操作は単語翻訳フェーズ完了と辞書参照成立後だけ有効にする。
  - リトライは retryable failure の時だけ有効にする。
  - job は Running のまま、phase state で完了、中断、回復可能失敗、再実行準備を区別する。
- `state_variants`:
  - `idle_ready`: Ready job で開始可能。
  - `running`: current phase と progress を表示する。
  - `empty_completed`: 対象語 0 件。Completed として扱い、provider 未実行を result summary に表示する。
  - `completed`: 確定訳語とジョブ内辞書 summary を表示する。
  - `paused`: 再開可能であることを表示する。
  - `recoverable_failed`: retryable flag と短い理由を表示する。
  - `blocked`: Ready 未成立、terminal job、term phase 未完了、provider setting 不整合。
- `post_implementation_review`:
  - desktop / mobile で長い語、provider 名、model 名、error reason がはみ出さない。
  - secret と raw response が UI、console、error summary に出ない。
  - progress、phase result、button state が同じ phase run を指す。

## Interface Frame

- `purpose`: 本文翻訳前に用語訳語の確定状況とジョブ内辞書反映を確認し、次へ進めるか判断する。
- `audience`: 翻訳ジョブを実行するユーザー。
- `primary_workflow`: Job Run を開き、単語翻訳フェーズを開始し、progress と phase result を確認し、必要なら再試行する。
- `information_density`: 実行状態、結果 summary、エラー理由を同一画面で比較できる密度にする。
- `visual_direction`: 既存 Job Run の運用画面として、装飾を増やさず状態と操作を優先する。
- `remembered_signal`: 単語翻訳フェーズが本文翻訳前に辞書を準備する段階であることを current phase と result summary で示す。

## Structure Notes

- `page_sections`:
  - phase header: current phase、state、progress、操作。
  - execution summary: 対象語件数、共通辞書 hit、AI 実行対象、provider summary、共通辞書 snapshot 要約。
  - result summary: 確定訳語、ジョブ内辞書反映、置換対象、後続 phase 可否。
  - error summary: error kind、短い理由、retryable flag、修正対象。
- `layout_constraints`:
  - phase header と primary actions は画面上部に固定して、result と error は下に積む。
  - table が必要な場合でも、狭幅では key-value list へ落とせる構造にする。
  - button text は短くし、状態説明は button の外に出す。
- `responsive_constraints`:
  - mobile 幅では summary を 1 列にし、件数と状態 label が折り返しても操作列を押し出さない。
  - 長い source term、provider 名、model 名、error reason は wrapping または truncation と tooltip 相当の確認導線を持つ。
  - progress 表示は幅に依存して文字が潰れない。
- `accessibility_constraints`:
  - phase state は色だけでなく text label で示す。
  - disabled button は理由を近接表示する。
  - error summary は retryable / blocking を text で区別する。

## Interaction States

- `loading`: phase result refresh 中は既存 summary を保持し、更新中 indicator を出す。
- `empty`: 対象語 0 件は空結果ではなく Completed の phase result として表示する。provider 未実行を result summary に出す。
- `error`: secret と raw response を含まない short reason、error kind、retryable flag を表示する。
- `disabled`: Ready 未成立、terminal job、active phase run あり、term phase 未完了、provider setting 不整合では操作を無効にし、理由を表示する。
- `progress`: progress percent、処理済み件数、AI 実行対象件数、現在 step を表示する。
- `retry`: retryable failure では retry action を表示し、同じ phase run を使うことを result に示す。
- `success`: 確定訳語件数、ジョブ内辞書反映件数、後続 phase 可否を表示する。

## Post Implementation Review

- `desktop_review_points`:
  - phase header、progress、primary actions、result summary が重ならない。
  - 確定訳語一覧または summary が長い語でも横にはみ出さない。
  - error summary に secret、raw response、本文全量がない。
- `mobile_review_points`:
  - current phase、progress、次操作が最初の縦スクロール範囲で確認できる。
  - summary card を入れ子にせず、縦に読みやすく積む。
  - button label が折り返しても隣接要素と重ならない。
- `overflow_risks`:
  - Skyrim の長い固有名詞。
  - provider / model 名。
  - stale dictionary や response mismatch の error reason。
  - file path または plugin 名が result summary に出る場合。
- `visual_polish_open_questions`:
  - Completed かつ provider 未実行の result label。
  - job Running と phase state label の組み合わせ。

## Rules

- UI は実装前の mock ではなく、実装が満たす要件契約として書く。
- 実装前の見た目 artifact を新規必須にしない。
- 細かな visual polish は実装後に人間が実物を確認して直す。
- product component 名や owned scope は implementation-scope で必要な時だけ扱う。
- implementation-scope の `owned_scope` や product code 対象 file は書かない。
