# UI Design: body-translation-phase

- `skill`: ui-design
- `status`: ready-human-review
- `source_plan`: `./plan.md`
- `scenario_source`: `./scenario-design.md`

## UI Contract

- `display_items`:
  - current phase、phase state、progress、対象 field 件数、処理済み件数、未処理件数。
  - `Job Setup` で設定した本文翻訳用 provider / model / execution mode の要約、credential 参照状態、request unit count、output count。
  - 辞書適用件数、persona 参照件数、翻訳補助メタデータ summary、prompt digest。
  - field result summary、訳文、出力ステータス、保護要素検証結果、後続 output readiness。
  - failure state、error kind、retryable flag、影響 field 件数、redacted error summary。
- `primary_actions`:
  - 本文翻訳フェーズ開始、pause、resume、retry、cancel。
  - field result 表示切替、保護要素検証結果の詳細表示、後続出力 readiness 確認。
- `button_enablement`:
  - `start`: persona phase Completed、辞書と persona snapshot 参照成立、非 terminal job、active phase run なし。
  - `pause`: body phase Running。
  - `resume`: body phase Paused。
  - `retry`: body phase RecoverableFailed かつ retryable failure あり。
  - `cancel`: body phase Paused の時だけ有効。Running から直接 cancel しない。
  - `output readiness`: body phase Completed かつ field result 整合時だけ有効。
- `state_variants`:
  - not-ready、ready、starting、running、paused、recoverable failed、validation failed、empty completed、completed、canceled、failed。
  - provider skipped、provider running、provider partial failure、save failure、late response rejected。
- `post_implementation_review`:
  - long source / translated text が field result 領域からはみ出さない。
  - error summary が secret、API key 平文、復号可能値を出さない。
  - desktop と mobile で progress、field summary、action bar が重ならない。

## Interface Frame

- `purpose`: 翻訳ジョブの本文翻訳フェーズを開始、監視、回復し、結果を後続出力へ渡せるか確認する。
- `audience`: Skyrim Mod 翻訳を進めるユーザー、失敗時に再実行判断をするユーザー、運用確認者。
- `primary_workflow`: Job Run を開き、本文翻訳フェーズを開始し、progress と field result を確認し、失敗時は retry または pause 後の cancel を判断する。
- `information_density`: 作業画面として高密度にする。phase summary、field summary、actions、error summary を同一 Job Run で縦に追える構成にする。
- `visual_direction`: `The Ethereal Archive` の glass panel、amber active state、serif typography を守る。過剰な装飾より、状態と件数の判読性を優先する。
- `remembered_signal`: 本文翻訳フェーズは、辞書、persona、AI provider、保護要素検証が集約される最終翻訳作業として見える必要がある。

## Structure Notes

- `page_sections`:
  - phase header: current phase、state、progress、primary actions。
  - input summary: target count、dictionary digest、persona digest、metadata digest、provider setting。
  - execution summary: provider execution summary、request unit count、success / failure / skipped count。
  - field result list: field identity、source excerpt boundary、translated text、output status、protection validation result。
  - recovery panel: error kind、retryable flag、affected field count、next action。
  - output readiness panel: completed field count、status consistency、readiness result。
- `layout_constraints`:
  - action bar は phase header 内に置き、field result list のスクロールで位置が揺れない。
  - field result は固定列幅に依存せず、狭い幅では縦積みにする。
  - provider summary と secret 非露出 summary は secret を表示する領域にしない。
- `responsive_constraints`:
  - mobile 幅では phase header、actions、summary、field result を 1 column にする。
  - desktop 幅では summary と field result を隣接表示できるが、カード内カードにはしない。
  - 長い EditorID、FormID、error kind は折り返し、ボタン幅を押し広げない。
- `accessibility_constraints`:
  - phase state と validation result は色だけでなく label と icon で示す。
  - retryable / non-retryable は button enablement と説明 label の両方で示す。
  - progress は数値と状態文を併記する。

## Interaction States

- `loading`: phase start、input snapshot 構成、provider 実行、field result 保存で個別 loading を表示する。
- `empty`: 対象 field 0 件または空 source は completed として表示し、単語だけの plugin でも成果物出力へ進めることを示す。
- `error`: provider failure、invalid response、protection validation failure、save failure、late response rejected を error kind と retryable flag で表示する。
- `disabled`: precondition 不足、terminal job、active phase run、output readiness 不成立では操作を無効化し、理由を短く表示する。
- `progress`: target count、processed count、success count、failure count、skipped count を出す。
- `retry`: retry 実行時は同じ phase run を継続する前提で、retry target count と成功済み field の扱いを表示する。
- `success`: completed state では訳文、出力ステータス、保護要素検証 summary、job-level Completed、output readiness を表示する。

## Post Implementation Review

- `desktop_review_points`:
  - phase header の action 群が 1280px 幅で折り返しても summary を覆わない。
  - field result list の訳文と validation badge が同じ行で判読できる。
  - recovery panel と output readiness panel が状態ごとに不要表示を残さない。
- `mobile_review_points`:
  - 390px 幅で button label、long FormID、EditorID、error kind が横 overflow しない。
  - primary actions が縦積みになっても誤操作しにくい間隔を保つ。
  - field result の source / translated text が前後 section を押しつぶさない。
- `overflow_risks`:
  - 長い訳文、長い raw error kind、長い provider model 名、複数 validation error。
  - 辞書 / persona / metadata digest の横並び表示。
  - output status が未固定の間の仮ラベル増加。
- `visual_polish_open_questions`:
  - validation failed の severity 表示色。
  - provider skipped と dictionary applied の visual distinction。
  - output readiness panel を phase header に寄せるか result summary に寄せるか。

## Human Decision Resolution

- `Q-BTP-001`: provider setting は `Job Setup` の本文翻訳用設定を表示する。開始時の再選択 UI は作らない。
- `Q-BTP-003`: partial failure は成功済み field result を保持し、phase 全体を recoverable failed として表示する。
- `Q-BTP-004`: validation failed field の失敗訳文は表示しない。再試行可否と失敗件数を表示する。
- `Q-BTP-005`: empty state は completed として表示する。
- `Q-BTP-006`: 原文と訳文のローカル表示は許容する。secret、API key 平文、復号可能値の非露出だけを必須にする。
- `Q-BTP-007`: body phase Completed で job-level Completed と output readiness を表示する。
- `Q-BTP-008`: cancel は Paused の時だけ有効にする。Canceled 後の途中成功結果は output readiness に使わない。
- `Q-BTP-009`: 結果確認から戻る導線とフィールド単体編集 UI は本 task では表示しない。
- `Q-BTP-010`: output status badge は後続の成果物出力に必要な語彙だけを表示する。

## Rules

- UI は実装前の mock ではなく、実装が満たす要件契約として書く。
- 実装前の見た目 artifact を新規必須にしない。
- 細かな visual polish は実装後に人間が実物を確認して直す。
- product component 名や owned scope は、implementation-scope で必要な時だけ扱う。
- implementation-scope の `owned_scope` や product code 対象 file は書かない。
