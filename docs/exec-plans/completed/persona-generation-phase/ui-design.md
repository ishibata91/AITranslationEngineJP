# UI Design: persona-generation-phase

- `skill`: ui-design
- `status`: ready-for-human-review
- `source_plan`: `./plan.md`
- `scenario_source`: `./scenario-design.md`
- `related_screens`: `app-shell.md`, `job-run.md`

## UI Contract

- `display_items`:
  - Job Run は current phase として `NPC ペルソナ生成` を表示する。
  - phase state、progress、target count、generated count、failed count、skipped count を表示する。
  - 生成対象 summary は NPC count、common persona hit / miss、対象外理由を表示する。
  - phase result は persona snapshot ID または snapshot digest、snapshot 参照状態、missing count、body phase readiness を表示する。
  - AI 実行 summary は provider、model、execution mode、credential ref、input count、output count、短い error kind を表示する。
  - redacted phase result summary は Job Run 再表示時に復元できる形で表示する。
- `primary_actions`:
  - persona phase 開始。
  - persona phase の pause、resume、retry、cancel。
  - body phase readiness 確認。
  - body phase 開始は persona phase Completed かつ snapshot 参照成立後だけ表示または有効化する。
- `button_enablement`:
  - 開始は term phase Completed、非 terminal job、active phase run なしの場合だけ有効にする。
  - retry は retryable failure の時だけ有効にする。
  - pause は Running の時だけ有効にする。
  - resume は Paused または RecoverableFailed の時だけ有効にする。
  - cancel は cancel 可能な non-terminal job の時だけ有効にする。
  - body phase 開始は persona phase Completed と snapshot 参照成立が両方 true の時だけ有効にする。
  - terminal job では persona phase start、save 後続操作、body readiness update を無効にする。
- `state_variants`:
  - loading、not started、running、paused、recoverable failed、failed、completed、empty completed、blocked、snapshot missing。
  - partial success は成功分を維持し、phase は RecoverableFailed として未処理 NPC だけ retry する。
  - common persona hit は新規 `PERSONA` を作らず、job の `snapshot ref` として表示する。
  - 生成対象 0 件は empty completed とし、対象 0 件、provider 未実行、snapshot 空を表示する。
- `post_implementation_review`:
  - 長い NPC 名、provider 名、model 名、error reason、snapshot digest が desktop / mobile で overflow しない。
  - secret、API key 平文、raw prompt、raw response、原文発話全文、会話文脈全文が画面に出ない。
  - debug log に prompt / request body を出す場合でも、secret と API key は出ない。
  - refresh 後も phase result と snapshot 参照状態が崩れない。

## Interface Frame

- `purpose`: Job Run 上で persona phase の開始、進捗、生成結果、後続 phase readiness を判断できるようにする。
- `audience`: 翻訳ジョブを進めるユーザーと、失敗理由を確認する運用確認者。
- `primary_workflow`: Job Run を開く、persona phase を開始する、progress と result を確認する、body phase readiness を確認する。
- `information_density`: 作業用 dashboard として高密度にする。hero、marketing copy、装飾 card は不要にする。
- `visual_direction`: 既存 Job Run の phase summary と同じ静かな業務 UI に合わせる。
- `remembered_signal`: persona snapshot 参照状態、generated / failed count、body phase readiness を第一視認情報にする。

## Structure Notes

- `page_sections`:
  - phase header: current phase、phase state、progress、主要 action。
  - target summary: target count、common persona hit / miss、skipped / blocked reason。
  - generation result: generated count、failed count、snapshot reference status、missing count。
  - AI execution summary: provider、model、execution mode、credential ref、input / output count。
  - error summary: redacted error kind、retryable flag、latest error summary。
- `layout_constraints`:
  - phase header は横幅が狭い時に 2 行へ折り返す。
  - counters は固定幅または安定した min-width を持つ。
  - snapshot digest と error reason は折り返しまたは省略表示を使い、隣接項目を押し出さない。
- `responsive_constraints`:
  - mobile では action group を縦積みにする。
  - count summary は 2 列以下に落とし、数値と label の対応を崩さない。
  - 長い NPC 名や model 名は 1 要素内で折り返す。
- `accessibility_constraints`:
  - phase state、retryable、body readiness は色だけで示さない。
  - disabled action は理由を読み取れる text または tooltip を持つ。
  - progress は数値と state label を併記する。

## Interaction States

- `loading`:
  - phase summary と target summary の読み込み中を表示する。
  - action は二重実行を避けるため無効にする。
- `empty`:
  - persona phase 未開始の場合は開始可能条件または開始不可理由を表示する。
  - 生成対象 0 件は empty completed として、provider 未実行と snapshot 空を表示する。
- `error`:
  - provider failure、invalid response、input missing、save failure、snapshot missing を短い error kind で表示する。
  - secret、raw payload、原文全文は表示しない。
- `disabled`:
  - term phase 未完了、active phase あり、terminal job、snapshot missing、body readiness false の理由を表示する。
- `progress`:
  - target total、processed count、generated count、failed count を同時に表示する。
  - execution mode が batch の場合も request unit と target count の関係を崩さない。
- `retry`:
  - retryable failure の時だけ retry action を有効にする。
  - retry 後は同じ phase run ID の progress と latest error 更新を表示する。
  - 一部 NPC 失敗時は成功分を維持し、未処理 NPC だけ retry 対象として表示する。
- `success`:
  - Completed では generated count、snapshot reference status、body phase readiness を表示する。
  - body phase 開始 action は readiness true の時だけ有効にする。

## Post Implementation Review

- `desktop_review_points`:
  - phase header、summary、action group が 1280px 幅で折り返さず読める。
  - long provider / model / NPC 名で counters と button が押し出されない。
  - running から completed へ runtime event 後に表示が更新される。
- `mobile_review_points`:
  - 375px 幅で action が縦積みになり、text が button 外へはみ出さない。
  - summary counters の label と数値が対応したまま読める。
  - error summary が下部 content を重ねない。
- `overflow_risks`:
  - snapshot digest、provider model、error reason、NPC display name、credential ref。
  - batch execution summary の provider job id または batch item id。
  - common persona hit / miss の理由 text。
- `visual_polish_open_questions`:
  - state badge の色と label は既存 Job Run の phase state 表示へ合わせる。
  - count summary の密度は実装後に人間が実画面で確認する。
  - common persona hit の UI 表示名は snapshot ref として実装後に人間が実画面で確認する。

## Rules

- UI は実装前の mock ではなく、実装が満たす要件契約として書く。
- 実装前の見た目 artifact を新規必須にしない。
- 細かな visual polish は実装後に人間が実物を確認して直す。
- product component 名や owned scope は、implementation-scope で必要な時だけ扱う。
- implementation-scope の `owned_scope` や product code 対象 file は書かない。
