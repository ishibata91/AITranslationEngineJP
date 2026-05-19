# Scenario Candidates: 2026-05-17-all-pages-componentization / state-transition

- `generator`: `state-transition`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `APC-ST`
- `candidate_count`: 12

## Generator Scope

- `viewpoint`: `state-transition`
- `included_sources`: `plan.md`, `docs/architecture.md`, `docs/coding-guidelines-frontend.md`, `docs/UX-standard.md`, Storybook 基盤の完了済み `scenario-design.md`, Master Persona 部品化の完了済み `scenario-design.md` と `ui-design.md`, `frontend/src/ui/screens/`, `frontend/src/ui/components/`
- `excluded_sources`: プロダクトコード変更、テスト変更、docs 正本変更、他 generator の採否、最終シナリオ統合
- `generation_notes`: 全ページ部品化で変わる状態境界、fixture 状態、story 状態、props 境界、既存状態表示維持だけを候補化する。

## Candidate Scenarios

### CAND-APC-ST-001 ページ状態を表示部品 props へ遷移する

- `source requirement`: `plan.md` の goal と close_conditions。ページ component を薄くし、主要表示領域が props 境界と Storybook story を持つ。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-APC-ST-001`
- `actor`: frontend implementation lane
- `transition before`: ページ component が controller 接続、view model 購読、表示分岐、主要表示領域を同じ Svelte file に持つ。
- `trigger`: パネル、カード、モーダル単位で screen local component を切り出す。
- `transition after`: ページ component は controller 接続と部品合成を持ち、表示部品は小さい props と callback だけで描画できる。
- `expected outcome`: 表示部品は `Store`、`Gateway`、generated binding、`RuntimeEventAdapter` を直接扱わない。
- `observable point`: component import、props 型、story source、fixture import。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: 全ページ共通の境界シナリオ候補として扱える。
- `conflict hint`: 画面専用の大きなレイアウトを無理に分ける候補と衝突する可能性がある。

### CAND-APC-ST-002 読み込み、空、失敗、通常表示を fixture 状態へ遷移する

- `source requirement`: `docs/UX-standard.md` は状態表示、空状態、ローディング状態、エラー状態を高優先または適用観点にしている。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-APC-ST-002`
- `actor`: frontend implementation lane
- `transition before`: 実画面の view model だけが読み込み中、空、失敗、通常表示を表す。
- `trigger`: 各主要部品へ固定 props または view model fixture を追加する。
- `transition after`: Storybook story は読み込み、空、失敗、通常表示を runtime なしで表示できる。
- `expected outcome`: 既存画面の状態文言、状態ラベル、空表示、エラー表示が部品化後も消えない。
- `observable point`: Storybook story variant、fixture object、状態表示 selector、既存 Svelte の状態分岐。
- `related detail requirement type`: `state_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: 入力レビュー、ジョブ管理、出力管理、辞書、各翻訳段階へ横断適用できる。
- `conflict hint`: failure 観点の異常系候補と重複する可能性がある。

### CAND-APC-ST-003 story callback は状態を実行系へ遷移させない

- `source requirement`: Storybook 基盤の完了済み `scenario-design.md` は story と fixture を fixed props または view model fixture に限定し、gateway mock と backend DTO mock を含めない。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-APC-ST-003`
- `actor`: frontend implementation lane
- `transition before`: 実画面の操作 callback は controller を経由して backend、Wails runtime、DB、filesystem flow へ進みうる。
- `trigger`: Storybook story で操作 callback を stub 化する。
- `transition after`: story 上の操作は表示確認用の no-op または story-local な状態確認に留まり、実行系状態へ遷移しない。
- `expected outcome`: story は AI provider、secret store、DB、実 filesystem、generated `wailsjs` を起動条件にしない。
- `observable point`: callback stub、禁止 import、`build-storybook` 結果。
- `related detail requirement type`: `state_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: Storybook 外部境界の共通シナリオ候補として扱える。
- `conflict hint`: external-integration 観点の境界候補と重複する可能性がある。

### CAND-APC-ST-004 操作可否状態を部品化後も維持する

- `source requirement`: `docs/architecture.md` は状態変更を event として View へ返し、View から `ScreenController` へ渡すと定義している。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-APC-ST-004`
- `actor`: 利用者
- `transition before`: 実画面は view model の `disabled`、有効条件、理由表示で操作可否を表す。
- `trigger`: 操作ボタン、sticky footer、action card、行操作を部品化する。
- `transition after`: 部品は props の操作可否と理由を表示し、許可状態だけ callback を親へ返す。
- `expected outcome`: 作成、開始、停止、再開、retry、削除、ページ移動、出力の有効状態と無効理由が維持される。
- `observable point`: button disabled、理由表示、callback 呼び出し境界、Storybook の許可状態と禁止状態。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: 操作部品を持つ全画面の共通候補として扱える。
- `conflict hint`: actor-goal 観点の操作シナリオと検証段階が重複する可能性がある。

### CAND-APC-ST-005 選択状態と詳細表示の整合を維持する

- `source requirement`: 既存 screen は一覧、選択行、詳細、空詳細を同一画面で表示する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-APC-ST-005`
- `actor`: 利用者
- `transition before`: 一覧未選択、行選択済み、絞り込み後空、古い選択が画面 view model 内で整合している。
- `trigger`: 一覧 component と詳細 component を分け、Storybook fixture を作る。
- `transition after`: 一覧状態が変わっても、詳細 component は存在しない古い対象を表示しない。
- `expected outcome`: 入力レビュー、辞書、未完了ジョブ一覧、出力候補、Master Persona 結果確認で選択と詳細が矛盾しない。
- `observable point`: selected id props、items props、detail props、空詳細 fixture、絞り込み後空 fixture。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: 一覧と詳細を持つ画面の横断候補として扱える。
- `conflict hint`: lifecycle 観点の一覧操作候補と統合される可能性がある。

### CAND-APC-ST-006 モーダル状態を対象保持つきで維持する

- `source requirement`: Master Persona 完了済み `ui-design.md` は編集、削除、保存失敗、削除失敗のモーダル状態を Storybook で確認すると固定している。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-APC-ST-006`
- `actor`: 利用者
- `transition before`: モーダルは closed、create、edit、delete、busy、failure の状態を画面状態から表示する。
- `trigger`: モーダルを screen local component として切り出し、fixture を作る。
- `transition after`: 失敗状態ではモーダルが閉じず、対象識別情報と入力値が維持される。
- `expected outcome`: Master Persona、マスター辞書、未完了ジョブ削除、API キー入力 panel の確認状態が欠落しない。
- `observable point`: modal state props、selected target props、form props、error props、dialog 表示、Storybook variant。
- `related detail requirement type`: `state_requirement`, `recovery_requirement`, `compatibility_requirement`
- `adoption hint`: modal と modal 相当 panel の状態候補として扱える。
- `conflict hint`: failure 観点の保存失敗、削除失敗候補と重複する可能性がある。

### CAND-APC-ST-007 ファイル選択と import 実行状態を維持する

- `source requirement`: 入力レビューとマスター辞書は、ファイル未選択、選択済み、実行中、結果表示、失敗表示を持つ。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-APC-ST-007`
- `actor`: 利用者
- `transition before`: JSON または XML の選択状態は hidden file input と画面 view model にまたがっている。
- `trigger`: import 準備 panel、選択ファイル表示、実行結果表示を部品化する。
- `transition after`: Storybook fixture は未選択、選択済み、実行中、完了、失敗を固定 props で表示できる。
- `expected outcome`: 選択ファイル名、保存場所、hash、進捗、結果件数、選び直し可否が維持される。
- `observable point`: `hasStagedFile`, `isImporting`, progress props、result props、button disabled、fixture variant。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `compatibility_requirement`
- `adoption hint`: 入力レビューと辞書 import の共通候補として扱える。
- `conflict hint`: external-integration 観点の filesystem 境界候補と重複する可能性がある。

### CAND-APC-ST-008 job と phase の実行状態表示を部品化後も維持する

- `source requirement`: `docs/spec.md` と `docs/er.md` は `TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` を分け、一覧とフェーズ画面の操作可否へ使う。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-APC-ST-008`
- `actor`: 利用者
- `transition before`: job 一覧、`JobRunPage`、単語翻訳、NPC ペルソナ生成、本文翻訳は Ready、Running、Paused、RecoverableFailed、Completed、Failed、Canceled などの状態を表示する。
- `trigger`: job card、phase status、action card、progress、readiness、navigation footer を部品化する。
- `transition after`: 状態値ごとの表示、操作可否、進捗、失敗情報、次段階 readiness が Storybook fixture で再現できる。
- `expected outcome`: job state と phase run state を混同せず、状態ラベルと操作可否が既存画面と同じ意味を保つ。
- `observable point`: state label、state pill、action disabled、blocked reason、progress、readiness props、phase story variant。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: 翻訳実行系画面の中心候補として扱える。
- `conflict hint`: state-transition 以外の lifecycle 候補と統合時に前提状態が衝突する可能性がある。

### CAND-APC-ST-009 phase 間 navigation 状態を footer へ遷移する

- `source requirement`: `JobRunPage` は単語翻訳、NPC ペルソナ生成、本文翻訳、翻訳完了の表示を選択中 job と readiness で切り替える。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-APC-ST-009`
- `actor`: 利用者
- `transition before`: `JobRunPage` が current phase page と footer の許可状態を同じ file 内で扱う。
- `trigger`: phase navigation footer を状態 props と callback で表示する。
- `transition after`: 次段階へ進める状態、進めない状態、未選択 job 状態、完了後の出力管理遷移状態を fixture で確認できる。
- `expected outcome`: 直リンク相当の job 未選択時は未完了一覧へ戻す案内を維持し、readiness 未達では次段階へ進めない。
- `observable point`: selected job props、current phase props、footer reasons、primary disabled、Storybook variant。
- `related detail requirement type`: `state_requirement`, `recovery_requirement`, `compatibility_requirement`
- `adoption hint`: `JobRunPage` の page story または footer story の候補として扱える。
- `conflict hint`: lifecycle 観点の画面遷移候補と重複する可能性がある。

### CAND-APC-ST-010 AI モデル選択状態を共有 component 境界で維持する

- `source requirement`: `docs/architecture.md` は複数画面で使う部品だけ `frontend/src/ui/components/` に置くと定義し、Master Persona 完了済み `ui-design.md` は `AIModelSelectionCard` 再利用を固定している。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-APC-ST-010`
- `actor`: 利用者
- `transition before`: AI 設定状態は各画面の phase card または設定詳細に分散している。
- `trigger`: `AIModelSelectionCard` または同等の共有 component へ props を渡す。
- `transition after`: provider、model、credential、モデル一覧更新、Batch API、処理方式、警告、長い model 名の状態を共有 component story で確認できる。
- `expected outcome`: 共有 component は画面固有条件を増やしすぎず、画面固有状態は props へ整形されて渡される。
- `observable point`: shared component props、screen local mapping、fixture variant、長い model 名 story。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: Provider Settings、Job Setup、Master Persona の AI 設定部品候補として扱える。
- `conflict hint`: 共通化基準が未確定のため、人間判断候補になりうる。

### CAND-APC-ST-011 出力候補と artifact 状態を維持する

- `source requirement`: 出力管理画面は completed job 一覧、選択中 job、readiness、出力操作、最新結果、diff preview を表示する。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-APC-ST-011`
- `actor`: 利用者
- `transition before`: 出力管理は job 未選択、job 選択済み、ready、not ready、submitting、last command あり、diff ありまたはなしを画面 view model で扱う。
- `trigger`: 出力候補一覧、選択 job summary、出力操作、diff preview を部品化する。
- `transition after`: Storybook fixture は出力不可、出力可能、再出力可能、失敗理由あり、diff empty、diff rows ありを表示できる。
- `expected outcome`: 出力可能状態と拒否理由、最新結果、stale 状態、xTranslator status 表示が維持される。
- `observable point`: viewState props、review props、canGenerate、canRegenerate、lastCommand、diffPreview、Storybook variant。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `compatibility_requirement`
- `adoption hint`: 出力管理画面固有の候補として扱える。
- `conflict hint`: output artifact の詳細仕様が未成熟な場合は人間判断候補へ回る可能性がある。

### CAND-APC-ST-012 story 追加状態と review 対象状態を task-local に遷移する

- `source requirement`: Storybook 基盤の完了済み `scenario-design.md` は review URL、story ID、確認状態、未確認理由を task-local 成果物へ残すと固定している。
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-APC-ST-012`
- `actor`: frontend implementation lane
- `transition before`: 対象 component に story がない、または確認状態が未記録である。
- `trigger`: 主要部品へ story と fixture を追加し、Storybook review 対象を記録する。
- `transition after`: 対象 story は未確認、確認済み、確認不能のいずれかの状態を task-local review 記録で追跡できる。
- `expected outcome`: 主要表示領域の story 不足を page story だけで代替しない。
- `observable point`: story file、story ID、fixture path、`storybook-review.md` の確認状態、`build-storybook` 結果。
- `related detail requirement type`: `state_requirement`, `observability_requirement`, `testability_requirement`
- `adoption hint`: operation-audit 観点と統合して review 証跡シナリオへできる可能性がある。
- `conflict hint`: review 記録そのものは operation-audit 観点と重複する。

## Open Notes

- `human decision candidate`: 全ページの部品化候補順序と、共有 component へ上げる具体基準は `plan.md` で未決である。
- `human decision candidate`: 既存表示項目を削る場合の承認粒度は `plan.md` で未決である。state-transition 候補では、削除を前提にしない。
- `human decision candidate`: page 合成 story をどの画面で任意または必須にするかは、Master Persona 以外ではまだ確定していない。
- `merge candidate`: `CAND-APC-ST-003` と external-integration 観点の Storybook 境界候補は統合候補である。
- `merge candidate`: `CAND-APC-ST-012` と operation-audit 観点の review 証跡候補は統合候補である。
- `rejection candidate`: 実行フロー再現、gateway mock、backend DTO mock、実ファイル読み込み、実 AI provider 通信を Storybook story の状態遷移として扱う候補は除外する。
