# UI Design: translation-job-setup

- `skill`: ui-design
- `status`: approved
- `source_plan`: `./plan.md`
- `scenario_source`: `./scenario-design.md`

## UI Contract

- `display_items`:
  - input selector: 入力データ名、出自、登録日時、翻訳レコード件数、既存 job 状態。
  - foundation selector: 共通辞書、共通ペルソナ、参照状態、参照不能理由。
  - AI runtime selector: provider、model、credential 参照状態、実行方式、capability 表示。
  - validation summary: pass / fail / warning、失効状態、失敗理由、最終検証時刻または検証断面。
  - create result: 作成済み job ID、`Ready` 状態、入力出自、実行設定要約。
- `primary_actions`:
  - Job Setup を開く。
  - 入力データ、共通辞書、共通ペルソナ、AI runtime、実行方式を選ぶ。
  - validation を実行または再実行する。
  - validation pass 後に create job を実行する。
  - failure reason から設定を修正する。
- `button_enablement`:
  - create job は validation pass かつ未失効の時だけ有効にする。
  - validation は入力データと最低限の設定が選ばれている時に有効にする。
  - 設定変更後は create job を無効にし、再 validation が必要であることを表示する。
  - 同一入力に既存 job がある場合は、状態に関係なく create job を無効にする。
  - Ready job は再表示だけ許可し、入力、基盤参照、AI runtime、実行方式の再編集導線は表示しない。
- `state_variants`:
  - loading、empty、progress、success、error、disabled、retry、dirty-validation を持つ。
  - input missing、cache missing、foundation ref missing、credential missing、unsupported provider / mode、partial create failure を区別する。
  - cache missing は Job Setup 内で再構築せず、Input Review の再構築導線へ戻す。
  - credential 解決、provider capability、ネットワーク到達性の失敗は blocking validation failure にする。
  - API key 平文、secret 本体、復号可能な値を表示しない。
  - paid な real AI API を UI 確認の前提にしない。
- `post_implementation_review`:
  - 入力、基盤参照、AI runtime、validation 結果、create 可否が 1 画面で追えるか確認する。
  - desktop と mobile で selector、summary、error reason が重ならないか確認する。
  - 長い plugin 名、file path、provider / model 名、failure reason が overflow しないか確認する。
  - keyboard だけで選択、validation、create、retry に到達できるか確認する。

## Interface Frame

- `purpose`: 翻訳ジョブ作成前に、入力データ、共通基盤、AI 実行設定、validation 結果を確認し、`Ready` job を作成する。
- `audience`: Skyrim Mod 翻訳作業者。
- `primary_workflow`: Job Setup で対象入力と実行設定を選び、validation pass を確認して create job を実行する。
- `information_density`: 作業用画面として、選択欄と validation summary を同時に追える密度にする。
- `visual_direction`: task-local UI 契約として、装飾より状態識別、長い値の可読性、操作可否の明確さを優先する。
- `remembered_signal`: `docs/screen-design/job-setup.md` は存在しない。今回の UI 契約は task-local の一時正本である。

## Structure Notes

- `page_sections`:
  - header: Job Setup の現在状態と create job action。
  - input panel: 入力データ選択、出自、既存 job 状態。
  - foundation panel: 共通辞書、共通ペルソナ、参照状態。
  - runtime panel: provider、model、credential 参照、実行方式。
  - validation panel: pass / fail / warning、dirty、failure reason、retry。
  - result panel: 作成後 job 要約と job detail への導線。
- `layout_constraints`:
  - desktop は input / foundation / runtime と validation summary を近接配置する。
  - mobile は section を縦積みにし、create job action を横にはみ出さない。
  - cards の入れ子を避け、section と repeated item を分ける。
  - selector と validation summary は固定幅前提にせず、長い値を折り返す。
- `responsive_constraints`:
  - file path、plugin 名、provider / model 名、failure reason は折り返しまたは省略表示を持つ。
  - create job と validation action は mobile で誤タップしにくい間隔を持つ。
  - status badge だけに依存せず、状態説明 text を併記する。
- `accessibility_constraints`:
  - selector、validation、create、retry は keyboard 操作可能にする。
  - error と warning は色だけで伝えない。
  - progress と validation result は screen reader が読める text を持つ。
  - validation failure 後の focus は failure reason または修正対象へ移せる。

## Interaction States

- `loading`:
  - 入力データ、共通基盤、AI runtime 設定の読込中を分けて表示する。
  - create job は disabled にする。
- `empty`:
  - 入力データがない場合は Job Setup を開始できない理由を表示する。
  - 共通辞書または共通ペルソナがない場合は選択欄を空状態にし、管理 UI 作成は本 scope に含めない。
- `error`:
  - validation failure、参照不能、credential 不備、provider / mode 不整合、create 失敗を区別する。
  - secret 値を error message に含めない。
- `disabled`:
  - validation pass がない、または失効している場合は create job を disabled にする。
  - disabled 理由を短く表示する。
- `progress`:
  - validation 中、create 中を区別する。
  - 実行中は重複 create を止める。
- `retry`:
  - validation failure は設定修正後に再 validation できる。
  - create 失敗は partial state が残らない前提で再試行可否を表示する。
- `success`:
  - 作成済み job、`Ready` 状態、入力出自、実行設定要約、validation 結果を表示する。
  - 作成後の Ready job は read-only 要約として表示し、再編集 action は出さない。
  - 翻訳開始、phase 実行、成果物出力は success state の必須 action にしない。

## Post Implementation Review

- `desktop_review_points`:
  - 入力、基盤、runtime、validation summary、create result が同時に読める。
  - 設定変更時に validation が失効したことが明確に見える。
  - API key 平文が画面、console、error summary に出ない。
- `mobile_review_points`:
  - selector を縦に積んでも現在の選択と validation 状態を見失わない。
  - action が折り返されても押し間違えにくい。
  - 長い failure reason が画面外にはみ出さない。
- `overflow_risks`:
  - input file path、plugin 名、provider / model 名、foundation 名、failure reason。
  - 複数 warning と blocking failure が同時に出る場合の validation summary。
  - 作成済み job ID や input ID の横長表示。
- `visual_polish_open_questions`:
  - validation warning の強調方法。
  - job 作成後の result panel を残す時間と位置。
  - app-shell navigation 上の表示名と位置。

## Source Gaps

- `docs/screen-design/job-setup.md`: 存在しない。
- `docs/scenario-tests/translation-job-setup.md`: 未作成。task-local scenario 承認後に昇格候補。
- `docs/detail-specs/job-setup.md`: 未作成。human 承認済み UI 要件だけが将来の昇格候補。

## Rules

- UI は実装前の mock ではなく、実装が満たす要件契約として書く。
- 実装前の見た目 artifact を新規必須にしない。
- 細かな visual polish は実装後に人間が実物を確認して直す。
- product component 名や owned scope は、implementation-scope で必要な時だけ扱う。
- implementation-scope の `owned_scope` や product code 対象 file は書かない。
