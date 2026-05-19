# Scenario Design: 2026-05-17-all-pages-componentization

- `skill`: scenario-design
- `status`: human-review-ready
- `human_review`: pending
- `source_plan`: `./plan.md`
- `ui_source`: 後続成果物。今回の起動入力では作成しない。
- `screen_design_diff`: UI 設計後に、画面設計書正本へ反映する差分がある場合だけ作成する。
- `final_artifact_path`: `docs/scenario-tests/all-pages-componentization.md`
- `topic_abbrev`: `APC`
- `candidate_sources`:
  - `./scenario-candidates.actor-goal.md`
  - `./scenario-candidates.lifecycle.md`
  - `./scenario-candidates.state-transition.md`
  - `./scenario-candidates.failure.md`
  - `./scenario-candidates.external-integration.md`
  - `./scenario-candidates.operation-audit.md`

## Fixed Requirements

- `must_pass_requirements`:
  - 全ページの部品化候補は、画面順、共有部品候補、画面専用部品候補、分けない候補へ分類する。
  - 画面順は、翻訳入力確認、ジョブ作成、翻訳ジョブ管理、ジョブ実行、単語翻訳段階、NPC ペルソナ生成段階、本文翻訳段階、翻訳完了、出力成果物、マスター辞書、Provider 設定、Dashboard、Master Persona の順を基準にする。
  - Master Persona は POC 済みの基準画面として扱い、全ページ review 対象には含めるが、同じ部品化判断を再設計しない。
  - ページ component は controller 接続、購読、dispose、通知、表示部品の合成へ寄せる。
  - パネル、カード、モーダルは小さい props と callback で表示できる境界にする。
  - 共有部品は、複数画面で使う表示規則や操作規則を集約する場合だけ `frontend/src/ui/components/` へ置く。
  - 画面固有の業務条件が増える候補は、screen local component に残す。
  - Storybook story は主要なパネル、カード、モーダルを必須対象にする。
  - page 合成 story は密度確認と配置確認に限定し、主要部品 story の不足を代替しない。
  - Storybook story と fixture は fixed props、view model fixture、callback stub だけを使う。
  - Storybook story と fixture は backend、Wails runtime、generated `wailsjs`、Gateway、RuntimeEventAdapter、AI provider、secret store、DB、実 filesystem flow を要求しない。
  - story、fixture、review 記録は secret、API key、token、実 endpoint、ローカル絶対 path、実ユーザーデータ、raw request、raw response、raw prompt、provider 応答原文を含めない。
  - 既存表示項目、状態文、操作、`aria-label` を削る場合は、この refactor の副作用として扱わず、UI 設計または screen design diff と人間承認へ戻す。
  - Storybook review URL、story ID、確認状態、未確認理由、再実行 command、`build-storybook` 結果は task-local `storybook-review.md` に残す。
  - 実装後検証は `npm --prefix frontend run build-storybook` と `python3 scripts/harness/run.py --suite frontend-local` を入口にする。
- `non_goals`:
  - プロダクト backend、Wails binding、Gateway、RuntimeEventAdapter、AI provider、secret store、DB の実装変更。
  - docs 正本、`.codex`、既存 scenario candidate、`plan.md` の変更。
  - Storybook 上で実 AI 生成、実ファイル読み込み、DB 書き込み、provider network、xTranslator 出力生成を再現すること。
  - UI 設計本文、画面設計書正本、implementation-scope の作成。
  - 画面内容の削除や導線変更を、部品化リファクタの範囲で確定すること。

## Scenario Candidate Coverage

正本: `./scenario-design.candidate-coverage.json`

6 種の候補成果物を統合した。
候補は 56 件である。
`needs_human_decision` は 0 件である。
未解決の競合は 0 件である。

候補で未決候補として出ていた項目は次の fixed decision で解消した。

- 部品化候補の画面順序は、現行 screen と翻訳フローの順に固定する。
- 共有部品へ上げる基準は `docs/architecture.md` の `UI Component` 判断表と `docs/coding-guidelines-frontend.md` の component 分割基準に固定する。
- 既存表示項目の削除は、今回の refactor では非対象にする。
- Storybook review 記録先は、完了済み Storybook 基盤と Master Persona POC に合わせて `storybook-review.md` に固定する。
- page 合成 story は任意の配置確認に限定する。

## Detail Requirement Coverage

正本: `./scenario-design.requirement-coverage.json`

抽象要件は、部品化候補、薄い page 境界、画面別 story、状態維持、既存表示維持、外部境界、保存禁止情報、review 証跡、検証再実行へ分けた。
`needs_human_decision` は 0 件である。
`deferred` は、UI 設計または implementation-scope が扱う成果物粒度に限定した。

## Human Decision Questionnaire

正本: `./scenario-design.questions.md`

none

## Scenario Matrix

### SCN-APC-001 全ページの部品化候補を分類する

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: 実装者と reviewer が、画面ごとの分割対象、共有候補、画面専用候補、分けない候補を追跡できる。
- `受け入れ条件`: UI 設計または component candidate 証跡に、対象画面、候補名、配置先、story 対象、分ける理由、分けない理由が残る。
- `事前条件`: Storybook 基盤と Master Persona POC の成果物を参照できる。
- `public_seam_or_api_boundary`: UI Component の props と callback 境界。
- `入力開始点`: `frontend/src/ui/screens/` と `frontend/src/ui/components/` の現行構成。
- `主要 outcome`: 全ページ一括 task の部品境界が、実装前に追跡できる粒度になる。
- `期待結果`:
  1. 画面順は fixed requirements の順序で記録される。
  2. 共有部品は複数画面で使う UI 規則に限られる。
  3. 画面固有条件が多い候補は screen local component に残る。
  4. 親画面状態を大量に読む候補、業務フロー全体を持つ候補、props が条件分岐の塊になる候補は分けない理由を持つ。
- `観測点`: component candidate 証跡、component path、story path、fixture path、分けない理由。
- `公開接点確認`: あり。UI Component の props 境界を公開接点として扱う。

### SCN-APC-002 ページ component を薄い合成役へ寄せる

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: ページ component が production 接続を持ち、表示部品が runtime 依存を直接持たない。
- `受け入れ条件`: ページ component は controller 接続、購読、dispose、通知、表示部品合成へ寄り、切り出した部品は props と callback だけで描画できる。
- `事前条件`: 対象画面の部品化候補が UI 設計で固定されている。
- `public_seam_or_api_boundary`: Svelte component props、callback、story args。
- `入力開始点`: 対象 page component と screen local component。
- `主要 outcome`: Storybook は主要表示部品を production runtime なしで描画できる。
- `期待結果`:
  1. 表示部品は `Store`、Gateway、generated `wailsjs`、RuntimeEventAdapter、controller factory を import しない。
  2. story fixture は backend DTO ではなく、Presenter 後の view model または component props を使う。
  3. 操作 callback は DOM event または画面操作結果だけを親へ返す。
  4. 操作可否と無効理由は同じ fixture 内で整合する。
- `観測点`: component import、props 型、callback 型、story source、fixture source、`build-storybook` 結果。
- `公開接点確認`: あり。

### SCN-APC-003 主要パネル、カード、モーダルを Storybook で確認する

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: 人間レビュアーが全ページの主要表示部品を Storybook 上で順に確認できる。
- `受け入れ条件`: 各画面の主要パネル、カード、モーダルに story があり、通常、空、読み込み、失敗、長文、選択状態の代表 fixture を確認できる。
- `事前条件`: story と fixture が追加され、Storybook dev server または static preview を開ける。
- `public_seam_or_api_boundary`: Storybook review URL と story ID。
- `入力開始点`: `storybook-review.md` に記録された review URL。
- `主要 outcome`: 主要表示部品の表示崩れ、長文耐性、空状態、エラー状態を確認できる。
- `開始操作`: Storybook review URL を開く。
- `入力方法`: task-local に記録された story ID と review URL を使う。
- `主要操作列`: Storybook を起動する。対象 story を開く。状態 variant を確認する。確認状態または未確認理由を記録する。
- `期待結果`:
  1. 翻訳入力、ジョブ作成、ジョブ管理、出力、辞書、Provider 設定の主要表示領域 story が確認対象になる。
  2. Master Persona は POC 済み story を基準にして、全ページ review 対象一覧へ含める。
  3. page 合成 story は配置確認用であり、主要部品 story の不足を代替しない。
- `観測点`: Storybook canvas、story ID、fixture state、`storybook-review.md` の確認状態。
- `UI-visible outcome`: 主要パネル、カード、モーダルの代表状態が表示される。
- `fake_or_stub`: fixed props、view model fixture、callback stub。

### SCN-APC-004 job と phase の状態表示、操作可否、navigation を維持する

- `分類`: 状態遷移
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: 翻訳ジョブ状態と phase 実行状態を混同せず、部品化後も操作可否が維持される。
- `受け入れ条件`: job card、phase status、action card、progress、readiness、navigation footer は、状態ラベル、進捗、失敗情報、次段階 readiness、無効理由を story fixture で表現できる。
- `事前条件`: job と phase の view model が component props へ整形されている。
- `public_seam_or_api_boundary`: phase panel props、navigation footer props、action callback。
- `入力開始点`: Ready、Running、Paused、RecoverableFailed、Completed、Failed、Canceled の代表 fixture。
- `主要 outcome`: 利用者は開始、停止、再開、retry、削除、次段階移動の可否を誤読しない。
- `期待結果`:
  1. job state と phase run state の表示責務が混ざらない。
  2. `Running`、`Paused`、`RecoverableFailed` の操作可否が既存仕様と一致する。
  3. 直リンク相当の job 未選択時は、未完了一覧へ戻る案内を維持する。
  4. readiness 未達では次段階へ進めない理由を表示する。
- `観測点`: state label、disabled 属性、無効理由、progress、readiness props、navigation footer story。
- `公開接点確認`: あり。

### SCN-APC-005 一覧、詳細、モーダル、file 表示の整合を維持する

- `分類`: 回帰防止
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: 部品化によって、選択状態、詳細表示、モーダル対象、file 表示、既存表示項目が欠落しない。
- `受け入れ条件`: 一覧と詳細を分けても古い選択を表示せず、モーダル失敗時は対象識別情報と入力値を維持し、file 参照は synthetic 表示値として確認できる。
- `事前条件`: 対象画面の状態別表示が fixture 化されている。
- `public_seam_or_api_boundary`: list props、detail props、modal props、file summary props。
- `入力開始点`: 空一覧、一覧あり、選択済み、絞り込み後空、保存失敗、削除失敗、file 未選択、file 選択済みの fixture。
- `主要 outcome`: 既存表示項目、状態文、操作、`aria-label` が refactor 後も維持される。
- `期待結果`:
  1. 一覧から消えた対象を詳細 component が表示しない。
  2. モーダル失敗時は dialog を閉じず、対象識別情報と入力値を残す。
  3. file story は実 file picker、実 file read、実 output generation を開始しない。
  4. 画面設計書にある表示項目、状態文、操作、`aria-label` を削る場合は human review へ戻す。
- `観測点`: selected id props、modal state props、error props、fixture state、screen design 対応、Storybook 表示。
- `公開接点確認`: あり。

### SCN-APC-006 Storybook は外部 runtime と実行系へ遷移しない

- `分類`: 禁止遷移
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: Storybook が backend 連携検証や実行フロー再現の代替にならない。
- `受け入れ条件`: story と fixture は generated `wailsjs`、Wails runtime、Gateway、RuntimeEventAdapter、AI provider、secret store、DB、実 filesystem flow を import または起動しない。
- `事前条件`: story、fixture、component が追加されている。
- `public_seam_or_api_boundary`: story source、fixture source、Storybook static build。
- `入力開始点`: `npm --prefix frontend run build-storybook`
- `主要 outcome`: Storybook static build は backend と Wails runtime なしで成立する。
- `期待結果`:
  1. story callback は no-op または story-local な表示確認だけに留まる。
  2. page story を作る場合も production gateway を暗黙生成しない。
  3. runtime event と polling adapter は Storybook に持ち込まない。
  4. Wails binding 未接続 error で story が止まらない。
- `観測点`: import 境界、callback stub、Storybook iframe render、`build-storybook` exit code。
- `公開接点確認`: あり。

### SCN-APC-007 Storybook fixture と review 証跡に保存禁止情報を含めない

- `分類`: セキュリティ
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: story、fixture、review 記録から secret と実データが漏れない。
- `受け入れ条件`: story、fixture、Storybook URL、screenshot 名、検証記録は保存禁止情報を含まない。
- `事前条件`: provider 設定、AI 設定、file 参照、進捗、review 証跡を持つ story が追加されている。
- `public_seam_or_api_boundary`: fixture source、story source、`storybook-review.md`。
- `入力開始点`: story、fixture、review 証跡。
- `主要 outcome`: Storybook review は synthetic data だけで再現できる。
- `期待結果`:
  1. fixture は secret、API key、token、実 endpoint、ローカル絶対 path、実ユーザーデータを持たない。
  2. fixture は raw request、raw response、raw prompt、provider 応答原文を持たない。
  3. review URL は Storybook localhost または iframe URL だけを指す。
  4. fakeAPI URL、Wails runtime URL、backend API URL を review URL として扱わない。
- `観測点`: fixture source、story source、review URL、story ID、検証結果要約。
- `公開接点確認`: あり。

### SCN-APC-008 Storybook review と UX / 人間 review の証跡を残す

- `分類`: 運用記録
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `最終検証`
- `観点`: 人間レビュアーと `implement_lane` が、全ページの見た目確認対象と承認状態を追跡できる。
- `受け入れ条件`: `storybook-review.md` は story ID、review URL、確認状態、未確認理由、再実行 command、`build-storybook` 結果を持つ。
- `事前条件`: Storybook dev server または static preview を開ける。
- `public_seam_or_api_boundary`: Storybook URL と task-local review 記録。
- `入力開始点`: `storybook-review.md`
- `主要 outcome`: review 対象、review 結果、指摘、未確認理由、承認状態、残留リスクが後から追跡できる。
- `開始操作`: Storybook review URL を開く。
- `入力方法`: task-local に記録された URL と story ID を使う。
- `主要操作列`: 対象 story を開く。表示、長文、狭い幅、状態差分を確認する。結果を task-local に残す。
- `期待結果`:
  1. story 欠落は未確認理由として残り、確認済み扱いにしない。
  2. UX review と frontend human review は混同せず、それぞれの承認状態を残す。
  3. review 記録は command 出力全文や不要な長文ログを保存しない。
- `観測点`: `storybook-review.md`、UX review 成果物、frontend human review 証跡、確認状態。
- `UI-visible outcome`: review URL で対象 story を再表示できる。
- `fake_or_stub`: fixed props、view model fixture、callback stub。

### SCN-APC-009 Storybook build、frontend-local、指摘対応を再実行できる

- `分類`: 完了判定
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `最終検証`
- `観点`: 全ページ部品化後の検証結果と再実行条件を追跡できる。
- `受け入れ条件`: `npm --prefix frontend run build-storybook` と `python3 scripts/harness/run.py --suite frontend-local` が通過する。または未通過理由、影響画面、再実行 command が task-local に残る。
- `事前条件`: 全対象ページの部品化、主要 story、review 記録が完了している。
- `public_seam_or_api_boundary`: npm script、harness suite、task-local 検証記録。
- `入力開始点`: `build-storybook` と `frontend-local`。
- `主要 outcome`: 実装後 gate の通過状態、未通過理由、再実行条件が判断できる。
- `期待結果`:
  1. Storybook static build は backend、Wails runtime、AI provider、secret store、DB なしで成立する。
  2. `frontend-local` は props 型、lint、既存画面テストの回帰を検出する。
  3. review 指摘後は対象 story と関連検証を再実行し、更新後の確認状態を残す。
  4. docs 正本化が必要な差分は、このシナリオで確定せず `implement_lane` へ戻す。
- `観測点`: command、exit code、短い失敗要約、対象 story、対象 file、再実行 command。
- `公開接点確認`: あり。

## Validation Entry

- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/2026-05-17-all-pages-componentization/scenario-design.md --report-out docs/exec-plans/active/2026-05-17-all-pages-componentization/scenario-design.requirement-gate.md --questionnaire-out docs/exec-plans/active/2026-05-17-all-pages-componentization/scenario-design.questions.md`
- `npm --prefix frontend run build-storybook`
- `python3 scripts/harness/run.py --suite frontend-local`

## Risks

- UI 設計本文はこの起動では作成しないため、画面ごとの最終 component candidate と Storybook review 順は `ui-design.md` で固定する必要がある。
- 全ページ一括 task のため、主要表示領域 story の不足が page story で隠れる可能性がある。
- 共有部品化を広く取りすぎると、画面固有条件が props 分岐として shared component に溜まる可能性がある。
- 現行 component が page view model 全体を受け取っている場合、小さい props 境界へ再分割する実装調査が必要になる可能性がある。
- 既存表示項目の削除が必要になった場合、この scenario-design だけでは進められず、人間レビューと screen design diff が必要になる。

## Next Artifacts

- `ui-design.md`: 全ページの component candidate、Storybook review 順、状態 variant、実装後見た目確認観点を固定する。
- `screen-design-diff.<screen-id>.md`: 画面設計書正本へ反映する差分がある場合だけ作成する。
- `implementation-scope.md`: 人間レビュー後に、部品化、story fixture、review 証跡、検証単位へ handoff を分割する。
