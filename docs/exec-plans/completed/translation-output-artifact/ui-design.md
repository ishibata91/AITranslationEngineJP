# UI Design: translation-output-artifact

- `skill`: ui-design
- `status`: ready-human-review
- `source_plan`: `./plan.md`
- `scenario_source`: `./scenario-design.md`

## UI Contract

- `display_items`:
  - completed job list、selected job summary、input provenance summary。
  - output readiness、拒否理由、result summary、output status distribution。
  - translation unit diff preview、Source、Dest、Status、row reflection summary。
  - artifact status、row count、generated_at、file path、re-output state。
  - compatibility summary、redacted error summary、operation summary。
- `primary_actions`:
  - completed job を選択する。
  - diff preview を開く。
  - xTranslator XML を出力する。
  - 出力済み artifact を再出力する。
  - compatibility / error summary から対象 unit へ移動する。
- `button_enablement`:
  - 出力 action は output readiness true、row validation pass、出力先 path valid の時だけ有効にする。
  - 再出力 action は existing artifact あり、output readiness true、stale または再出力可能状態の時だけ有効にする。
  - diff preview action は selected job と field summary が参照可能な時だけ有効にする。
  - invalid job、`Canceled` job、field result 不整合、status 不整合では出力 action を無効にし、理由を隣接表示する。
- `state_variants`:
  - loading、empty、not-ready、ready、preview-ready、generating、success、failed、stale、re-output-needed。
  - target count 0、row count 0、readonly path、XML parse failure、compatibility warning を区別する。
  - provider / network / secret が不要であることは操作説明ではなく、error summary と監査要約の非露出条件で保証する。
- `post_implementation_review`:
  - desktop と mobile で job list、summary、diff preview、action rail が重ならないことを確認する。
  - 長い plugin 名、file path、Source、Dest、error reason が overflow しないことを確認する。
  - redacted 表示が伏せ字ではなく secret 再表示になっていないことを確認する。

## Interface Frame

- `purpose`: 完了済み翻訳 job から、xTranslator 互換 XML の出力前確認、出力、再出力、失敗確認を行う。
- `audience`: Skyrim Mod 翻訳成果物を確認して配布形式へ出したい利用者。
- `primary_workflow`: completed job 選択、result summary 確認、diff preview 確認、XML 出力、artifact 状態確認、必要なら再出力。
- `information_density`: 運用ツールとして密度を高める。装飾より、job、status、row、error を比較しやすくする。
- `visual_direction`: 既存 design system に従い、静かな作業画面にする。hero、marketing 表現、過剰な AI 風装飾は使わない。
- `remembered_signal`: 選択中 job、選択中 artifact、diff filter、compatibility filter は画面内 state として保持する。

## Structure Notes

- `page_sections`:
  - 左側または上部に completed job list を置く。
  - 中央に selected job の result summary と readiness を置く。
  - 下部または詳細領域に diff preview と row validation result を置く。
  - action rail は出力、再出力、summary 更新をまとめ、無効理由を近くに出す。
- `layout_constraints`:
  - job list、summary、diff preview、action rail をカードの入れ子にしない。
  - diff preview は列幅を固定し、Source と Dest の長文で action rail を押し出さない。
  - error summary は上部の状態帯に出し、row detail の本文と混ぜない。
- `responsive_constraints`:
  - mobile では job list、summary、diff preview を縦順にし、action は sticky footer ではなく section 内に置く。
  - file path、plugin 名、Source、Dest は折り返し、長い単語は省略表示と詳細展開で扱う。
  - row preview は横スクロールを前提にせず、主要列を優先表示する。
- `accessibility_constraints`:
  - 出力不可理由は色だけで伝えず、短い日本語文言と status icon を併用する。
  - action button は disabled の理由を tooltip または隣接文言で確認できる。
  - diff preview は追加、変更、欠損をテキスト label で区別する。

## Interaction States

- `loading`:
  - completed job list 読み込み中、selected job summary 読み込み中、diff preview 読み込み中を分ける。
  - 読み込み中に既存 summary を消さず、更新中 state を表示する。
- `empty`:
  - completed job がない場合は、出力可能な job がないことを表示する。
  - target count 0 の completed job は empty と扱わず、row count 0 の出力可能 summary として表示する。
- `error`:
  - readiness 取得失敗、row validation 失敗、XML 出力失敗、file write 失敗を別表示にする。
  - secret、API key、provider raw payload、本文全文を error summary に出さない。
- `disabled`:
  - 出力不可の action は disabled にし、理由を `not_completed`、`canceled`、`status_mismatch`、`missing_required_row_field` などの表示語彙へ写像する。
- `progress`:
  - XML 出力中は generating state とし、row count、stage、cancel ではなく待機表示を出す。
  - provider 実行の progress と混同しない。
- `retry`:
  - 保存失敗、XML 出力失敗、readonly path は再出力または出力先変更へ誘導する。
  - field 再翻訳の retry はこの画面の primary action にしない。
- `success`:
  - artifact status、row count、file path、generated_at、compatibility summary を表示する。
  - success 後も diff preview と re-output state を確認できる。

## Post Implementation Review

- `desktop_review_points`:
  - completed job list、summary、diff preview、action rail が 1280px 幅で重ならない。
  - long plugin name、long path、long Source / Dest、multi-line error reason が layout shift を起こさない。
  - success、failed、stale、row count 0 の state が一目で区別できる。
- `mobile_review_points`:
  - 390px 幅で job selection、summary、diff preview、actions の順序が保たれる。
  - action button text が折り返しても押下領域と隣接要素が重ならない。
  - diff preview の主要情報が横スクロールなしで確認できる。
- `overflow_risks`:
  - XML file path、plugin path、FormID list、Source、Dest、compatibility warning、rejection reason。
  - redacted placeholder が長くなり、summary の主要 count を押し出す可能性がある。
- `visual_polish_open_questions`:
  - compatibility warning を badge、inline warning、summary panel のどれで強調するか。
  - diff preview で Source と Dest を同時表示するか、選択 unit detail へ逃がすか。
  - artifact state を timeline で見せるか、latest summary だけにするか。

## Rules

- UI は実装前の mock ではなく、実装が満たす要件契約として書く。
- 実装前の見た目 artifact を新規必須にしない。
- 細かな visual polish は実装後に人間が実物を確認して直す。
- product component 名や owned scope は、implementation-scope で必要な時だけ扱う。
- implementation-scope の `owned_scope` や product code 対象 file は書かない。
