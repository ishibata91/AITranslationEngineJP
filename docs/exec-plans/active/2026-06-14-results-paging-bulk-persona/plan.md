# Task Plan: 2026-06-14-results-paging-bulk-persona

- `workflow`: work
- `status`: 着手中（preparation-module 完了。重 task のため design-module へ）
- `task_id`: 2026-06-14-results-paging-bulk-persona
- `task_mode`: 重 task（画面が動く＝結果一覧へページャ UI を追加するため。preparation-module の軽/重判定で確定）
- `request_summary`: 翻訳結果一覧が台詞数に比例した DB クエリ（N+1）を出さないよう話者属性を一括取得し、数万件でも一覧をページングで扱えるようにする。T2 の `/simplify` で記録した効率課題の恒久対応。
- `goal`: 結果一覧の取得 DB クエリ数を台詞数に依存させず（bulk または JOIN で一定本数）、結果一覧をページ単位で取得・表示できる状態にする。表示形（T2 で合意済みのコンパクト行＋口調チップ＋展開）は変えない。
- `constraints`: T2 で合意済みの表示（`TranslationProgress` / `TranslationResultRow` / `ResultsPanel` の props 形・style・story）を変えない。固有名（辞書）は対象外（T3）。AI 翻訳の本番実行は利用者 provider 依存のまま。
- `close_conditions`: (1) 結果取得（`ListResults` 相当）の DB クエリ数が台詞数 N に比例しない（話者属性を bulk/JOIN で取得し、`LoadLineSpeaker` の per-line 呼び出しを廃止）。(2) 結果一覧をページ単位（サーバ側ページング）で取得・表示でき、数万件規模でも一覧が破綻しない。(3) 実画面で複数ページを跨いで結果と口調差を確認できる。
- `source_branch`: `master`
- `source_commit`: `ed88b638`（master HEAD。T2 /simplify cleanup merge 後）
- `work_branch`: `claude/2026-06-14-results-paging-bulk-persona`
- `target_branch`: `master`

## 完了定義（preparation-module で固定）

task 後に観測できる振る舞いを 2 つ固定する。最小実装・空実装で goal を満たしたことにしない。

- 動かす範囲 1（N+1 廃止）: 結果一覧の取得（`ListResults` 相当）で、話者属性の取得 DB 呼び出し回数が台詞数 N に比例しない。`buildResults` が台詞ごとに `store.LoadLineSpeaker` を呼ぶ構造を、一括取得（`LoadLineSpeakers(lineIDs)` または `ListLines` への JOIN）1 経路へ置き換える。
  - 観測点: 単体テスト。話者取得呼び出し回数を数える fake store または query 計数で、台詞数 N を変えても一括取得が定数回（per-line 呼び出し 0 回）であることを assert する。
- 動かす範囲 2（ページング）: 結果一覧をサーバ側ページ単位で取得・表示でき、実画面で複数ページを跨いで結果と口調差を確認できる。数万件規模でも一覧表示が破綻しない。
  - 観測点: 実画面（`npm run dev:wails:run`、`http://localhost:34115`）。Innocence Lost - Quest Expansion.esp（121 台詞）を 1 ページの件数より小さい page size で流し、ページャ操作で複数ページに結果が分かれ、各ページで口調チップ・口調差が見えることを目視する。

- goal 整合: goal「DB クエリ数を台詞数に依存させない／ページ単位で取得・表示」を、上の 2 観測点が実際に動くことで満たす。close_conditions はこの 2 観測点で検証できる形にしてある。
- goal と「含まない」の矛盾: なし。固有名（T3）と口調ルール精緻化（T4）は goal に含まず、本 task の手段にも要らない。表示形（コンパクト行・口調チップ・展開）は T2 合意済みを維持し、ページャ追加は表示形の変更に当たらない。

## 軽 / 重判定（preparation-module で固定）

- 画面が動くか: Y。結果一覧へページャ UI（ページ送り表示部品、件数・現在ページ表示）を追加する。`storybook-module` で表示を承認する。
- `docs/architecture.md` 反映が要るか: N。N+1 廃止は store/engine/api 内部の一括取得化、ページングは既存 api→store 経路への offset/limit/total 追加で、層構成・依存方向・Bootstrap・Wails 境界の構造・強い制約を変えない。API 署名へのパラメータ追加は既存境界内の拡張で、構造記述（§8）の書き換えを要しない見込み。
- 判定: 片方（画面が動く）が Y のため重 task。経路は preparation-module → design-module（architecture 反映 N のため設計差分図は不要。実装範囲・テスト設計・人間設計レビューを固定）→ storybook-module（ページャ表示）→ implementation-module → finalization-module。

## Scope（含む / 含まない）

含む:
- 話者属性の一括取得: `internal/store` に `LoadLineSpeakers(ctx, lineIDs)` の bulk 取得（`IN (...)` で speaker/race/voice_type JOIN ＋ faction を 2 クエリ程度）を足す、または `ListLines` に話者属性を JOIN して返す形へ拡張。`internal/model` に必要なら話者属性つきの投影型を足す。
- N+1 廃止: `internal/api` の `buildResults` が台詞ごとに `engine.LineDirective`（→ `store.LoadLineSpeaker`）を呼ぶ構造を、一括取得した `map[int64]SpeakerIdentity` から口調指示文・口調要約を組む形へ置き換える。`internal/engine` に identity→persona をバッチ適用する経路を用意。
- サーバ側ページング: `api` の結果取得にページ範囲（offset/limit または keyset）と総件数を足し、`store` の `ListNarrations`/`ListLines` を範囲取得へ拡張。
- frontend のページング: 結果一覧のページ取得・表示（ページ番号送りまたは仮想スクロール）。表示コンポーネント（ページャ等）は storybook-module、state/取得は implementation-module。
- 翻訳実行中（`engine.Run`）の per-line 話者クエリも、翻訳前に話者属性を一括取得して再利用する（任意。翻訳は AI latency 支配のため優先度は低い）。

含まない:
- 固有名解決・マスター辞書（T3）。
- 口調ルールの精緻化、編集 UI（T4）。
- 表示形（コンパクト行・口調チップ・展開）の変更。

## 依存

- T2（`2026-06-14-persona-dictionary-pipeline`、完了・master へ merge `a9b4125e`）の台詞抽出・話者解決・結果一覧・進捗。
- 現状コード: `internal/api/app.go`（`buildResults`/`ListResults`/`RunExtractAndTranslate`）、`internal/engine/engine.go`（`Run`/`LineDirective`/`linePersona`）、`internal/store/line.go`（`ListLines`/`LoadLineSpeaker`）、`frontend/src/ui/screens/translation-run/`（結果一覧・container・gateway）。

## 設計メモ（T2 の /simplify レビュー由来）

- N+1 の所在:
  - `internal/api/app.go` の `buildResults` が `store.ListLines` で全台詞を取得後、台詞ごとに `engine.LineDirective(ctx, l.ID)` を呼ぶ。`LineDirective` は `engine.linePersona` 経由で `store.LoadLineSpeaker(lineID)` を実行し、1 台詞あたり 2 クエリ（speaker+race+voice JOIN、faction SELECT）。N 台詞で 2N クエリ。
  - `internal/engine/engine.go` の `Run` も台詞翻訳ループで台詞ごとに `linePersona` を呼び、同じ 2 クエリ／台詞。翻訳完了後に `buildResults` が再度同じ取得をするため、1 実行で話者取得が 2 回（翻訳時＋表示時）走る。
- 修正案: `store` に lineID 群の一括取得を足し、`api`（表示）と任意で `engine`（翻訳前）が `map[int64]SpeakerIdentity` を 1 度作って使い回す。`buildResults` のループ内 DB アクセスを廃止。
- ページング: サーバ側 LIMIT/OFFSET か keyset で結果をページ取得し総件数を返す。進捗 event（extract/translate）はページングと独立のため整合に注意（翻訳は全件、表示はページ単位）。

## 設計判断（design-module で未確定を解決。人間設計レビュー対象）

3 つの未確定を次のように確定する。

- 判断 1（ページング方式）= サーバ側 keyset（cursor）ページング ＋ 総件数。ページサイズ 50（人間設計レビューで確定。LIMIT/OFFSET から変更）。cursor は連結列上の現在位置を表す不透明文字列（`""`=先頭、`n:<id>`=叙述文の id まで読了、`l:<id>`=台詞の id まで読了）。前進は cursor を使い、後退は frontend が訪問済み cursor を履歴として持って戻す。理由: 大量件数で深い位置でも `WHERE id > ?` の index 走査で取得でき、OFFSET の読み飛ばしコストを避ける。並び順は id 昇順で安定。
- 判断 2（話者一括取得の形）= `LoadLineSpeakers(lineIDs) → map[int64]SpeakerIdentity`（map 返し）。理由: 所属勢力は話者に対し 1 対多のため `ListLines` への単純 JOIN は行が増殖し group_concat 等が要る。map 返しは所属勢力を `[]string` のまま保て、T2 の責務分担（store は識別子・事実、engine が口調へ解釈）を維持する。同じ map を表示（buildResults）と翻訳（engine.Run）が使い回す。
- 判断 3（混在ページング順序とキー）= 叙述文（id 昇順）を先頭に、台詞（id 昇順）を続けた連結列を 1 つのページ列とみなす。cursor が叙述文区間にある間は叙述文を `id > afterID` で読み、ページに満たなければ台詞を先頭から補充して台詞区間へ移る。次 cursor は読み終えた最後の行の区間と id で組む。frontend のキーは従来通りページ内 index（ページ単位でまるごと差し替えるため衝突しない。表示形を変えない）。

## 実装範囲（design-module で固定）

scope の境界と依存だけ示す。詳細実装は implementation-module で Claude 本体が文脈を保って書く。

- N+1 廃止（per-line メソッドを一括メソッドへ置換し、N+1 を表現できない形にする）:
  - `internal/store/line.go`: `LoadLineSpeaker(lineID)`（per-line、2 クエリ）を削除し、`LoadLineSpeakers(ctx, lineIDs) → map[int64]model.SpeakerIdentity` を新設。IN 句で speaker/race/voice を 1 度、所属勢力を 1 度引く。SQLite の変数上限を避けるため IN 句は内部で分割（chunk）する。先頭話者の採り方（id 昇順先頭）は per-line と同基準。
  - `internal/engine/engine.go`: `LineDirective(lineID)`・`linePersona(lineID)`（per-line）を削除し、`LinePersonas(ctx, lineIDs) → map[int64]Persona{Directive,Label}` を新設。内部で `store.LoadLineSpeakers` を 1 度呼び、各 identity を `personaFromIdentity`→`buildPersonaDirective`/`buildPersonaLabel` で口調へ写す。`LineStore` interface の `LoadLineSpeaker` を `LoadLineSpeakers` へ差し替え。
  - `internal/engine/engine.go` の `Run`: 未訳台詞の id 群を集め、ループ前に `LinePersonas` を 1 度呼んで map を作り、ループ内は map 参照に変える（per-line DB 廃止）。進捗通知は変えない。
  - `internal/api/app.go` の `buildResults`: ページ内台詞の id 群で `engine.LinePersonas` を 1 度呼び、map から各行の口調指示・要約を引く（ループ内 `LineDirective` 廃止）。
- サーバ側ページング（keyset cursor）:
  - `internal/store`: `CountNarrations`/`CountLines`（総件数）、`NarrationsAfter(ctx, afterID, limit)`/`LinesAfter(ctx, afterID, limit)`（`WHERE id > ? ORDER BY id LIMIT ?`）を新設。次ページ有無の判定に limit+1 取得を使う。
  - `internal/api/app.go`: `ResultPage{Total int, Results []ResultView, NextCursor string, HasMore bool}` を新設。`ListResultsPage(cursor string, limit int) → ResultPage` を公開（起動時・ページ送り・実行後の取得を 1 経路に統一）。cursor 解析（`""`=先頭／`n:<id>`／`l:<id>`）と区間跨ぎの補充・次 cursor 組み立てを api orchestration で行う。叙述文区間で limit に満たなければ台詞先頭から補充して台詞区間へ移る。`RunExtractAndTranslate` は全件 `Results` を返すのをやめ、`RunResult{TranslatedCount}` の要約だけ返す（数万件をインラインで返さない）。
  - `bootstrap`: 新 store メソッドの配線確認（New の署名は変えない見込み。engine の Store interface 差し替えに追従）。
- frontend ロジック（implementation-module）:
  - `gateway/translation-gateway.ts`: `listResultsPage(cursor, limit) → {total, results, nextCursor, hasMore}` を追加、`listResults`（全件）を置換。`runExtractAndTranslate` の戻りを `{translatedCount}` へ。
  - `TranslationRunContainer.svelte`: keyset 状態（訪問済み cursor 履歴・現在ページ index・総件数・nextCursor・hasMore・ページサイズ定数 50）を持つ。`次へ`=nextCursor で取得し履歴へ積む、`前へ`=履歴を 1 つ戻して再取得。onMount と実行後は cursor `""` から page 0 を読む。ページャ操作 callback を View へ渡す。
- frontend 表示（storybook-module）:
  - 結果一覧へページャ表示部品（前へ・次へ＝端で無効化、現在ページ番号、総件数）を追加。件数バッジを「現在ページ件数」から「総件数」へ。cursor ページングは順次送りのため番号ジャンプは持たない。`TranslationResultRow` の表示形（コンパクト行・口調チップ・展開）は変えない。

含まない（再掲）: 固有名解決・マスター辞書（T3）、口調ルール精緻化・編集 UI（T4）、表示形そのものの変更。

## テスト設計（design-module で固定）

単体テストで書く対象（純ロジック・ビジネスロジック・過去基準の再発防止）:

- keyset cursor 送り（api、fake store で orchestration を確認）: cursor `""` から page を順に取り、(1) 叙述文だけで埋まるページ→次 cursor が `n:<id>`、(2) 叙述文を読み切り台詞へ補充して跨ぐページ→次 cursor が `l:<id>`、(3) 台詞区間のページ→次 cursor が `l:<id>`、(4) 末尾ページ→`HasMore=false`、(5) `Total` が叙述文＋台詞の合計、を assert。fake store なので DB アクセスは含まない（純 orchestration）。
- N+1 廃止の観測（engine、完了定義の観測点）: 話者取得呼び出し回数を数える fake store で `LinePersonas` を呼び、台詞数 N（例: 2 件と 50 件）を変えても `LoadLineSpeakers` 呼び出しが定数（1 回）で、per-line 呼び出しが 0 であることを assert。
- `LinePersonas`（engine）: fake store が返す map から、解決できた話者に口調指示・要約が付き、map に無い台詞は口調なし（空）になること。既存 `TestLineDirective` を batch 版へ書き換える。
- 既存 engine テストの追従: `fakeStore` を `LoadLineSpeakers(lineIDs)→map` 実装へ変更。`Run` 系テスト（ペルソナ注入・話者なし・進捗）は挙動不変で通ること。

単体テストで書かない対象（実画面・実データで観測。完了定義の観測点）:

- `LoadLineSpeakers`・`*Range`・`Count*` の SQLite クエリ（DB アクセス込みの統合経路）。先頭話者の採取・所属勢力の付与が従来と一致することは、実 mod（Innocence Lost、Grelod=気難しい老女・Aventus=ノルドの子供 等）の口調差が実画面で従来通り見えることで観測する。
- ページャ操作・Wails 経路・複数ページ跨ぎ表示は実画面（localhost:34115）で観測する。

## 合意済み frontend 保護（storybook-module で固定）

storybook-module で人間レビュー承認した画面表示。implementation-module はこの表示を変えずに配線する。

- 承認済み画面: 結果一覧へ keyset ページャを追加。`ResultsPager.svelte`（新規、順次送り「← 前へ ｜ ページ N ｜ 次へ →」、端で無効化）、`ResultsPanel.svelte`（件数バッジを総件数へ、結果ありのときページャを下に配置）。
- 表示規則: ページャは presentation のみで state を持たず、表示値（`pageNumber`・`canPrev`・`canNext`）と操作 callback（`onPrev`・`onNext`）を props で受ける。`ResultsPanel` の paging props は任意（既定は単一ページ＝前後無効）。総件数は `ResultsPanel` の `total` props（未指定なら現在ページ件数）。
- 変更禁止範囲: `TranslationResultRow.svelte` の表示形（コンパクト 1 行・口調チップ・展開）。番号ジャンプ UI は持たない（cursor は順次送り）。
- 反映先 frontend ファイル: `ResultsPager.svelte`（追加）、`ResultsPanel.svelte`・`TranslationRunScreen.svelte`（変更）、`translation-run-view.ts`（`ResultsPaging` 型）、`ResultsPager.stories.ts`・`ResultsPanel.stories.ts`・`TranslationRunScreen.stories.ts`・`translation-run.fixtures.ts`。
- 後続実装で表示を変えずに済む境界（props 形）: `ResultsPanel` は `total?`・`pageNumber?`・`canPrev?`・`canNext?`・`onPrev?`・`onNext?` を受ける。`TranslationRunScreen` は `paging?: ResultsPaging`・`onPagePrev?`・`onPageNext?` を受け `ResultsPanel` へ写す。container はこの props を満たす state・取得を配線する。

## Routing Notes

- `required_reading`:
  - `docs/exec-plans/completed/2026-06-14-persona-dictionary-pipeline/plan.md`（T2 の実装範囲・成果物）
  - `docs/architecture.md`（engine/store/api の責務、Wails 境界）
  - `docs/coding-guidelines-backend.md` / `docs/coding-guidelines-frontend.md`

## Outcome

- preparation-module 完了（2026-06-14）。作業 branch `claude/2026-06-14-results-paging-bulk-persona` を master `ed88b638` から作成。完了定義（N+1 廃止＝話者取得呼び出しが台詞数非依存／ページング＝実画面で複数ページ確認）と軽/重判定（重 task。画面 Y、architecture N）を固定。
- design-module 完了（2026-06-14）。設計差分図は不要（architecture 反映 N）。人間設計レビュー承認（修正付き）: ページング方式を LIMIT/OFFSET から keyset（cursor）へ変更、ページサイズ 50 で確定。実装範囲・テスト設計を keyset 前提で固定。画面が動く（ページャ）ため次は storybook-module。
- storybook-module 完了（2026-06-15）。ページャ表示（順次送り・端で無効化・件数バッジを総件数へ）を Storybook 人間レビューで承認。`ResultsPager` 追加、`ResultsPanel`／`TranslationRunScreen`／fixture／story を更新。全 story を通常分類へ復帰、一時 review 資源を削除。`build-storybook`・`frontend-local` 通過、console エラーなし。`合意済み frontend 保護` を固定。次は implementation-module。
- implementation-module 完了（2026-06-15）。
  - backend: `store` に `LoadLineSpeakers(lineIDs)→map`（IN 句・chunk）／`NarrationsAfter`・`LinesAfter`・`CountNarrations`・`CountLines` を新設、per-line `LoadLineSpeaker`・`ListNarrations`・`ListLines` を削除（共通 helper `query.go`）。`engine` の `LineStore` を一括取得へ差し替え、`LineDirective`・`linePersona`(per-line) を `LinePersonas(lineIDs)→map[int64]Persona` へ置換、`Run` はループ前に一括取得。`api` に `ResultPage{Total,Results,NextCursor,HasMore}`・`ListResultsPage(cursor,limit)`・`pageRows`・cursor helper を新設、`buildResults` を keyset ページング＋一括 persona へ、`RunExtractAndTranslate` を `RunResult{TranslatedCount}` 要約だけへ。`bootstrap` は concrete の method 追加で追従（署名変更なし）。
  - frontend: `gateway` に `listResultsPage(cursor,limit)→{total,results,nextCursor,hasMore}`、`runExtractAndTranslate` を `{translatedCount}` へ。`container` に keyset state（cursor 履歴・ページ index・total・nextCursor・hasMore・PAGE_SIZE 50）と前へ/次へ、onMount/実行後は先頭ページ取得。wails binding 再生成。
  - テスト: engine `TestLinePersonas`（batch）・`TestLinePersonasBulkLoadsOnce`（N=2,50 で `LoadLineSpeakers` 1 回＝台詞数非依存）・Run 経路の一括 1 回観測、api `TestPageRows*`（連結列の順次送り・区間跨ぎ・末尾・空）。fakeStore を一括版へ更新。
  - 最終検証: backend `lint:backend`（format/vet/static 0 issues/arch OK/module）通過、`test:backend` 通過。frontend `svelte-check`（node_modules 既存エラーのみ、自作 0）・`frontend-local`・`build-storybook` 通過。
  - 完了定義 観測: (1) N+1 廃止＝単体テストで話者取得が台詞数非依存（1 回）を確認。(2) 実画面（localhost:34115）で Innocence Lost（121 台詞）を local LLM stub＋page size 50 で流し、3 ページ（50/50/21）を 次へ/前へ で跨ぎ、総件数バッジ 121・口調差（種族:ノルドの子供／声質:気難しい老女の声／幼い少年の声／幼い少女の声／若々しい女性の声）・端の無効化（先頭=前へ無効、末尾=次へ無効）を確認。console エラーなし。dev DB の stub 訳は gitignore 対象で commit されない。
  - 仕様変更・人間承認候補: なし（公開境界の意味拡張なし、docs 正本反映は finalization で判断）。次は finalization-module。
- finalization-module（2026-06-15）。
  - 正本化判断: `docs/architecture.md` への反映は不要。判断結果＝後続課題にも切り出さない（廃案でもなく、反映自体が不要）。根拠: §5 Wails 境界（Gateway→bindings→api Bind、進捗は runtime events）・§7 ディレクトリ正本（store/engine/api/gateway の責務）・§8 現在の状態（extractor が line/speaker を書き、engine が台詞翻訳＋ペルソナ＋進捗）は本 task で変えていない。keyset ページングと一括 persona は層内の内部精緻化、`ListResults`→`ListResultsPage` は既存 Bind 境界内のメソッド変更で、architecture の抽象度の記述を書き換えない。人間承認状態＝不要（恒久仕様追加なし）。
