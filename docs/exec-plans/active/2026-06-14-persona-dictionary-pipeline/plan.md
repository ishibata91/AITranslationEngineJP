# Task Plan: 2026-06-14-persona-dictionary-pipeline

- `workflow`: work
- `status`: in_progress（implementation-module 完了。finalization-module へ）
- `task_id`: 2026-06-14-persona-dictionary-pipeline
- `task_mode`: 重 task（preparation-module で確定）
- `request_summary`: 実行画面から実 mod を抽出し、台詞を話者属性からのペルソナ口調つきで本文翻訳し、本文翻訳の進捗と口調差を画面で確認できる縦切りを 1 本通す。固有名（辞書）は本 task から外し、マスター辞書 task（T3）へ移す。
- `goal`: 翻訳実行画面から `Innocence Lost - Quest Expansion.esp` を抽出して台詞を中心 DB へ入れ、台詞ごとに話者属性からペルソナ口調指示文を組んで本文翻訳へ注入し、本文翻訳の進捗をバーで見せ、翻訳結果と口調差を画面で確認できる状態を通す。
- `constraints`: 固有名解決とマスター辞書（テーブル・登録・適用）は本 task では扱わず、`2026-06-14-master-dictionary`（T3）で扱う。ルール・プロンプトの編集 UI（T4）は対象外。会話履歴解析と AI ペルソナ生成（`system_requirements.md` §3 で将来検討）は対象外。
- `close_conditions`: 翻訳実行画面から `Innocence Lost - Quest Expansion.esp` を実行したとき、(1) 抽出した台詞が期待件数（extractor の count と一致）で中心 DB に入り、(2) 本文翻訳の進捗がバーで表示され、(3) 話者属性から組んだペルソナ口調指示文が注入された翻訳結果と、口調の差が画面で確認できる。
- `source_branch`: `master`
- `source_commit`: `1b7ea778`
- `target_branch`: `master`
- `work_branch`: `claude/2026-06-14-persona-dictionary-pipeline`

## 完了定義（preparation-module、design-module で固有名を除外）

システム上どこまで動かすかを 1 つに固定する。

- 動かす範囲: 翻訳実行画面から実 mod `Innocence Lost - Quest Expansion.esp` を対象に、抽出 → ペルソナ口調つき本文翻訳 を一気通貫で動かす。具体的には次の 4 つが実際に効く。
  1. extractor が台詞（line）と話者属性（speaker、および master 側 NPC から解決した race / faction / voice_type）を中心 DB へ書く。今は narration（叙述文）だけを書くため、台詞と話者を新たに永続化する。書く値は識別子と事実（EDID / FormID）に限る。
  2. engine が話者属性から `system_requirements.md` §3 のとおり機械的・AI 不使用のテンプレートでペルソナ口調指示文（翻訳ディレクティブ）を組み、本文翻訳プロンプトへ注入する。話者属性から口調 traits を引く最小ルール 1 系統は engine 内に置く。
  3. engine が台詞を本文翻訳し、口調指示文を注入した結果を書き戻す。
  4. api が本文翻訳の進捗を runtime events で push し、frontend が本文翻訳の進捗バーを描く。
  差込点（interface 宣言・引数追加・空テーブル）を置くだけで、画面に抽出件数も進捗も訳文も口調差も出ない状態は「動く」と書かない。
- 観測点: 実 app の実画面と実データを主とする。`npm run dev:wails:run`（`http://localhost:34115`）で起動し、`Innocence Lost - Quest Expansion.esp` を選んで実行し、(a) 台詞が期待件数で抽出され（extractor の count モード出力と一致）、(b) 本文翻訳の進捗バーが進み、(c) 口調指示文が注入された翻訳結果と話者ごとの口調差が結果一覧に出ることを目視する。補助として、ペルソナ口調指示文の組み立ての単体テストと、extractor の台詞・話者書き込みの単体テストを置く。
- goal 整合: goal は実 mod を画面から流して本文翻訳と口調差を観測できることを要求する。本完了定義は、抽出件数・進捗・訳文・口調差のいずれも画面で観測できることを要求し、最小実装・空テーブル・単体テストだけで goal を満たしたとはしない。
- 矛盾検査: goal が要する手段（台詞・話者の抽出、口調注入、本文翻訳、進捗 UI）は、除外項目（固有名/マスター辞書＝T3、編集 UI＝T4、AI ペルソナ生成）のいずれも必要としない。口調ルールは機械テンプレート 1 系統で足りる。goal と「含まない」は矛盾しないため、本モジュールで停止しない。

## 軽 / 重判定（preparation-module）

- 画面が動くか: Y。本文翻訳の進捗バーを翻訳実行画面に新設し、進捗 event を購読する。結果一覧に口調指示文を併記する。表示構造・props・story・fixture を変える。
- `docs/architecture.md` 反映が要るか: Y。§8「現在の状態と移行」は「extractor の SQLite writer は narration のみ」を前提に書かれているが、本 task で台詞・話者の書き込みを足すため現状記述が変わる。engine の進捗 event 経路（§5 runtime events）とペルソナ生成の実装も反映対象になりうる。最終判断は finalization-module。
- 結論: 両方 Y → 重 task。経路は preparation-module → design-module → storybook-module（進捗バーと口調表示）→ implementation-module → finalization-module。

## Scope（含む / 含まない）

含む:
- extractor（C#/Mutagen）に、台詞（line）と話者（speaker）および話者の race / faction / voice_type を中心 DB へ書く経路を足す。話者は esp に無く master 側 NPC にいるため、台詞 → 話者 NPC → 属性を master 連鎖から解決する。
- engine に、話者属性からのペルソナ口調指示文生成（機械テンプレート、AI 不使用）と、本文翻訳プロンプトへの注入を足す。話者属性 → 口調 traits の最小ルール 1 系統を engine 内に置く。
- api に、本文翻訳の進捗を runtime events で push する経路を足す。
- frontend に、本文翻訳の進捗バーと event 購読、結果一覧への口調指示文の併記を足す。
- 実 mod `Innocence Lost - Quest Expansion.esp` を対象にした実行確認。

含まない:
- 固有名解決とマスター辞書（テーブル・登録・適用、`proper_noun` 抽出、`line_mention` e5、name 関連 e8/e13/e14）。すべて `2026-06-14-master-dictionary`（T3）で扱う。
- ルール・プロンプトの編集 UI と実プロンプト参照（T4）。口調ルールは本 task では機械テンプレート 1 系統に留める。
- 会話履歴解析と AI ペルソナ生成（`system_requirements.md` §3 で将来検討）。

## 依存

- T1（`2026-06-14-extract-translate`）の narration 抽出 → 翻訳 → 画面表示の縦切りが master に存在すること。
- 翻訳実行画面（`frontend/src/ui/screens/translation-run/`）と `RunExtractAndTranslate`（`internal/api/app.go`）が存在すること。
- 実 mod と master 群が `dictionaries/Data/`（`Innocence Lost - Quest Expansion.esp`、`Skyrim.esm` ほか）に存在すること。

## 後続 task への影響（固有名の移送）

- 固有名解決は `2026-06-14-master-dictionary`（T3）へ移す。T3 の plan は依存に「T2 の辞書解決の差込点」と書くが、本 task は差込点を持たないため、その依存は無効になる。T3 着手時に T3 自身が辞書解決の差込点を engine へ作る。T3 の plan 修正は T3 着手時に行う。

## 現状の事実（preparation-module / design-module で確認）

- extractor は narration（BOOK:DESC）のみ SQLite へ書く（`tools/extractor/NarrationSqliteWriter.cs`）。台詞・話者は in-memory 抽出のみで未永続。INFO は話者条件（ANAM・CTDA）を持ち、`PluginEnvironment` の LinkCache で master 連鎖の NPC→race/voice/faction を解決できる。
- engine は narration を 1 件ずつ翻訳するだけで、ペルソナ・進捗 event を持たない（`internal/engine/engine.go`）。provider の system prompt は定数で directive 注入口が無い（`internal/provider/openai_compatible.go`）。
- 翻訳実行画面は idle / running / done / error の 4 状態のみ表示し、進捗バーと event 購読を持たない（`frontend/src/ui/screens/translation-run/`）。
- `Innocence Lost - Quest Expansion.esp` の実抽出 count: response lines（台詞 NAM1）121、info nodes 107、dialogues 61、books（叙述文）0、speakers（esp 内 NPC_/TACT）0。話者は base-game NPC で master 側にいる。

## Routing Notes

- `required_reading`:
  - `docs/system_requirements.md`（§3 ペルソナ＝ルールベース、構造化属性＋翻訳ディレクティブ、永続境界）
  - `docs/concept-model.md`（話者・素材・性質の合成で口調が決まる関係）
  - `docs/er.md`（`line` / `speaker` / `race` / `faction` / `voice_type` / `line_speaker` / `speaker_faction`）
  - `docs/architecture.md`（engine の責務、§5 runtime events、§6 C#↔Go 境界、§8 現在の状態）

## design-module の成果物

- 設計差分図: `design-diff.md`（コンポーネント差分図、シーケンス図、確認観点）。
- 人間設計レビュー: 非 UI 部分を承認（責務分担、ペルソナ最小ルール、スキーマ再利用、進捗の backend 経路、口調差を観測する方針）。UI 表示（進捗バー、口調併記）は design-module では承認せず、storybook-module の Storybook 人間レビューループで設計・承認する。`summary.md` はレビュー終了につき削除。

## 実装範囲（design-module 承認後に固定。UI 表示は storybook-module）

非 UI（design-module 承認済み、implementation-module で実装）:
- extractor（C#）: `tools/extractor` に台詞・話者 writer を足す。`line` へ INFO:NAM1 を書き（ordinal に response 番号）、INFO の話者条件（ANAM・CTDA）から LinkCache で話者 NPC を解決し `speaker`/`race`/`faction`/`voice_type`/`line_speaker`/`speaker_faction` を識別子・事実で書く。`NarrationSqliteWriter` のパターン（schema ensure＋INSERT OR IGNORE＋件数返却）を踏襲し、`Program.cs` の `--sqlite` 経路で narration に加えて台詞・話者も書く。
- db: backup `cdd8798c` の `db/migrations/0002` を再利用（line/speaker/race/faction/voice_type/line_speaker/speaker_faction）。
- model/store（Go）: `internal/model/line.go`（`Line`、`SpeakerPersona`）と `internal/store/line.go`（未訳台詞＋話者属性の join 読み出し、dest 更新）を backup から再利用・調整。
- engine（Go）: 話者属性→口調 traits の最小ルールと `buildPersonaDirective`（backup `persona.go`）でペルソナ口調指示文を組み、`Run` に台詞翻訳と directive 注入、本文翻訳 phase の進捗 callback を足す。固有名（辞書）は配線しない（T3）。
- provider（Go）: `Translate` に directive 引数を足し system prompt の base 指示文後段へ注入。
- api（Go）: `RunExtractAndTranslate` を台詞抽出・台詞翻訳・進捗 event 配線へ拡張。進捗 event の契約: event 名＋payload `{ phase, done, total }`。結果 `RunResult` の各行に注入した口調指示文を載せる。
- bootstrap（Go）: line store を engine/api へ配線。

UI（storybook-module が story＋svelte コンポーネントで設計・承認。implementation-module が event/state を配線）:
- 表示（storybook-module）: 本文翻訳の進捗バー、結果行への口調指示文併記。fixture と story で各状態（抽出中／翻訳中 N/M／完了、口調あり／無し行）を固定し Storybook 人間レビューで承認。
- frontend ロジック（implementation-module）: Wails runtime events の購読（EventsOn）、gateway の event 経路、container の進捗 state と結果 state、`RunExtractAndTranslate` 呼び出し。

依存: T1 narration 経路、既存 run 画面、実 mod と master 群（`dictionaries/Data/`）。

## テスト設計（design-module）

単体テスト（書く）:
- ペルソナ口調指示文の組み立て: `buildPersonaDirective`（`SpeakerPersona`→指示文、属性 0 件で空文字）。
- 口調ルール: 声型/種族/勢力 識別子→口調 traits の純関数（境界、未知識別子）。
- provider の directive 注入: directive が system メッセージへ入る。
- extractor の台詞・話者書込: C# 単体テスト（temp sqlite、行数・FK・UNIQUE 冪等）。

単体テストで書かない（E2E・実画面に任せる）:
- store の join 読み出し（DB アクセス込み）。
- api の `RunExtractAndTranslate` 統合経路（extractor 子プロセス＋engine＋event）。
- 進捗 event 購読・進捗 state derive・画面表示（storybook の story と実画面で確認）。

シナリオ（実画面・実データ、完了定義の観測点）:
- `Innocence Lost - Quest Expansion.esp` を実行画面から流し、台詞が期待件数で抽出され、本文翻訳の進捗バーが進み、口調指示文つき訳文と話者ごとの口調差が結果一覧に出ることを目視する。

## design-module への引き継ぎ（実装順の素案）

- 重 task として design-module → storybook-module → implementation-module → finalization-module を通す。
- 実装順（1 本の観測可能成果として並べる。横の seam 層を別 task に切り出さない）: extractor に台詞・話者書き込み（LinkCache で話者解決）→ engine にペルソナ口調生成と本文翻訳の directive 注入 → api に本文翻訳の進捗 event → frontend に進捗バーと口調併記 → 実 mod 実行確認。
- 再利用候補: backup branch `backup/2026-06-14-persona-dictionary-pipeline-prev`（`cdd8798c`）の `db/migrations/0002`（line / speaker / race / faction / voice_type / line_speaker / speaker_faction）、`internal/model/line.go`、`internal/store/line.go`、`internal/engine/persona.go` の `buildPersonaDirective`。固有名側（`internal/engine/dictionary.go`、`buildGlossaryDirective`）は本 task では使わず T3 へ送る。extractor の台詞・話者書き込みと UI 進捗・口調併記は backup には無い。

## storybook-module の成果物（合意済み frontend 保護）

Storybook 人間レビュー承認済み。表示の正本は次の story と svelte コンポーネント。

- 承認済み画面・表示規則:
  - 本文翻訳の進捗バー（`TranslationProgress.svelte`）。抽出中は不定バー＋「台詞を抽出しています」、本文翻訳中は確定バー＋「本文を翻訳しています done / total（percent%）」。`phase==="running"` かつ進捗ありで表示。
  - 結果行はコンパクト 1 行（`TranslationResultRow.svelte`、`details/summary`）。畳んだ summary に状態・EDID・口調チップ（声質などの短い要約）・原文/訳文の抜粋を 1 行で出し、行クリックで原文・訳文・口調指示の全文を展開する。口調を持たない行は「口調なし」を控えめに出す。一覧のまま口調チップで口調差を観測できる。
  - 結果一覧（`ResultsPanel.svelte`）は行間 `gap-2` の密な並び。
- 反映先 frontend ファイル: `TranslationProgress.svelte`（追加）、`TranslationResultRow.svelte`（変更）、`TranslationRunScreen.svelte`（進捗バー配置を追加）、`ResultsPanel.svelte`（行間）、`translation-run-view.ts`（`RunProgress`・`ProgressStage`・`NarrationResultRow.directive`・`personaLabel`）、`translation-run-presentation.ts`（`PROGRESS_STAGE_LABEL`）、`translation-run.fixtures.ts`。
- 通常分類へ戻した story:
  - `UI Components/TranslationProgress`（抽出中・翻訳中 途中・翻訳中 ほぼ完了）
  - `UI Components/TranslationResultRow`（畳む 口調あり、展開 口調あり、畳む 口調あり 老女、畳む 口調なし）
  - `Screens/翻訳実行`（既存 4 状態＋実行中 抽出・実行中 本文翻訳・完了 口調差）
- 後続実装で表示を変えずに済む境界（変更禁止範囲）: 上記 svelte 表示コンポーネントの構造・props 形・style・story・fixture。implementation-module はこれらの表示を変えずに、props へ値を供給する。
- 検証: `npm --prefix frontend run check`（本 task 変更にエラーなし）、`npm --prefix frontend run build-storybook`（成功）、`python3 scripts/harness/run.py --suite frontend-local`（通過）。

## implementation-module へ渡す frontend ロジック（表示範囲外）

- Wails runtime events の購読（`EventsOn`）で本文翻訳の進捗（契約 `{ phase, done, total }`）を受け、画面の `progress` props へ橋渡しする。
- gateway に進捗 event の購読経路を足す。
- container（`TranslationRunContainer.svelte`）に進捗 state と結果 state を持つ。結果行の `directive`・`personaLabel` は backend の `RunResult` から供給する（api の契約に短い口調要約 `personaLabel` と全文 `directive` を含める）。
- `RunExtractAndTranslate` 呼び出しと、台詞抽出・台詞翻訳・進捗 event 配線（backend 側）。

## implementation-module の成果物

### 実装（変更ファイル）

backend（Go）:
- 追加 `db/migrations/0002_persona_dictionary.sql` — line/speaker/race/faction/voice_type/line_speaker/speaker_faction（backup から再利用）。
- 追加 `internal/model/line.go` — `Line`、`SpeakerIdentity`（事実=EDID）、`SpeakerPersona`（口調 traits）。
- 追加 `internal/store/line.go` — 未訳台詞・全台詞の取得、dest 更新、`LoadLineSpeaker`（話者の race/voice/faction EDID を join 読み出し）。
- 追加 `internal/engine/persona.go` — `buildPersonaDirective`（口調指示文）、`buildPersonaLabel`（口調チップ、一般声は種族へ後退）。
- 追加 `internal/engine/persona_rule.go` — 識別子→口調 traits の最小ルール（声型は完全一致＋命名推定、種族は子供推定、勢力は最小表）。
- 変更 `internal/engine/engine.go` — `Store`（Narration＋Line）、`Run` に台詞翻訳・進捗 callback、`LineDirective`（表示用）。
- 変更 `internal/provider/openai_compatible.go` — `Translate` に `directive` 引数、`systemPrompt` で base 指示の後段へ注入。
- 変更 `internal/api/app.go` — `RunExtractAndTranslate` を台詞抽出・台詞翻訳・進捗 runtime events 配線へ拡張、`ResultView`（directive/personaLabel）、`ListResults`、`buildResults`。
- 不変 `internal/bootstrap/bootstrap.go` — 既存の `engine.New(s, p)` が新署名と一致し変更不要（固有名は配線しない）。

extractor（C#）:
- 追加 `tools/extractor/LineSpeakerSqliteWriter.cs` — 台詞（INFO:NAM1）と話者属性を書く。INFO の話者 FormKey を LinkCache で NPC へ解決し race/voice_type/faction の EditorID を書く。
- 変更 `tools/extractor/Program.cs` — `--sqlite` 経路で narration に加えて台詞・話者を書く。

frontend ロジック（implementation-module 範囲）:
- 変更 `frontend/src/gateway/translation-gateway.ts` — `ListResults`、結果に directive/personaLabel、進捗 event 購読 `onRunProgress`。
- 変更 `frontend/src/ui/screens/translation-run/TranslationRunContainer.svelte` — 進捗 state、event 購読、結果配線。
- 変更 `frontend/src/ui/screens/translation-run/ResultsPanel.svelte` — `{#each}` key を一意な index へ（実データの重複台詞で起きる `each_key_duplicate` crash の修正、表示は不変＝UX 中立のため implementation-module で対応）。
- 再生成 `frontend/wailsjs/go/api/App.*`、`frontend/wailsjs/go/models.ts` — Go API 変更を反映（`wails generate module`）。

### テスト

- `internal/engine/engine_test.go` — narration 翻訳、provider error、台詞のペルソナ注入、話者なし、進捗 callback、`LineDirective`。
- `internal/engine/persona_test.go` — `buildPersonaDirective` / `buildPersonaLabel`。
- `internal/engine/persona_rule_test.go` — 声型の完全一致と命名推定、識別子→traits、未知の畳み込み。
- `internal/provider/openai_compatible_test.go` — directive の system メッセージ注入、空 directive は base のみ。
- `internal/api/app_test.go` — `narrationResultView`（口調なし）、statusLabel、extractor 引数。
- `tools/extractor.Tests/LineSpeakerSqliteWriterTests.cs` — 台詞書込（列・ordinal）と冪等性。

### 観測ログ

- 追加なし。ペルソナ口調指示文の組み立てと識別子→traits は決定的で単体テストで証明済み。extractor/engine のエラーは `fmt.Errorf` で文脈付き。実行時にしか確定しない値や原因分離が要る分岐は無い。

### 最終検証

- Go（backend、harness に backend suite が無いため go ツール直実行）: `gofmt -l` 出力なし、`go vet ./internal/... ./db/...` 通過、`go build ./internal/... ./db/...` 通過、`go test ./internal/... ./db/...` 通過（api/engine/provider/store ok）。
- C#: `dotnet test tools/extractor.Tests`（19 件、台詞 writer 2 件含む）通過。
- frontend: `svelte-check`（本 task 変更にエラーなし）、`python3 scripts/harness/run.py --suite frontend-local`（lint＋test）通過。
- structure: fail 2 件。いずれも本 task 以前から存在する `docs/.DS_Store`（macOS 生成物）と `docs/mutagen-migration-plan.md`（index 未登録の既存 doc）で、本 task の変更集合と無関係（本 task の exec-plan docs はすべて PASS）。
- 実画面・実データ（完了定義の観測点、`npm run dev:wails:run` / `http://localhost:34115`）:
  - 抽出: `Innocence Lost - Quest Expansion.esp` を流し、台詞 121 件（extractor の count モードと一致）を中心 DB へ書き込み。話者は 7 件（Grelod=FemaleOldGrumpy、Aventus=NordRaceChild、Constance=FemaleYoungEager、子供 4 名）解決、121 件中 109 件に話者が紐づく。migration は user_version 1→2 適用。
  - ペルソナ口調・口調差: 結果一覧の各台詞に口調チップが出る（Grelod「声質: 気難しい老女の声」、Aventus「種族: ノルドの子供（…）」、子供「声質: 幼い少年/少女の声」、Constance「声質: 若々しい女性の声」、衛兵「口調なし」）。話者ごとの口調差を画面で観測。展開で口調指示文の全文を確認。
  - 進捗・翻訳: ローカル OpenAI 互換 stub で実行し、進捗 event は extract→translate（done 0→121、total 121）の 2 phase が発火、全 121 件に訳文（dest）が入り `仮訳` で表示。directive 注入は provider 単体テストで確認済み。
  - AI 翻訳の本番実行は利用者の OpenAI 互換 provider（画面の endpoint/API キー）で行う。本検証は stub で pipeline 全体（抽出→口調注入→翻訳→進捗→画面表示）を通した。

### finalization-module への引き継ぎ

- 仕様変更・追加の有無: 無し。scope（ペルソナ口調、固有名は T3）内で完結。
- 人間承認候補: 無し。
- `docs/architecture.md` 反映: §8「現在の状態と移行」は extractor の SQLite writer が narration のみと記すが、本 task で台詞・話者も書くため現状記述が変わる。engine が台詞翻訳・ペルソナ生成・進捗 event（§5 runtime events）を持つことも反映対象。最終判断は finalization-module。構造（コンポーネント・依存方向・port）は変えていない。
- 後続 task: 固有名解決は T3（`2026-06-14-master-dictionary`）。T3 の依存「T2 の辞書解決の差込点」は無効（T2 は差込点を持たない）。

## finalization-module の成果物

### 正本化判断

- 反映対象: `docs/architecture.md`。影響範囲: §8「現在の状態と移行」の現状記述のみ。
- 判断: 構造（§1〜7：コンポーネント・各責務・依存方向・手動 DI・Wails 境界・C#↔Go 境界・ディレクトリ正本）は不変。engine 責務（辞書解決・ペルソナ生成）は §3、runtime events は §5、extractor の SQLite 書込は §6 に既記載で、追加スキーマは `er.md` に既定義のため、新たな恒久仕様は無い。§8 の現状記述だけが事実（extractor が台詞・話者も書く、engine が台詞翻訳・ペルソナ・進捗を持つ）と食い違うため最小更新する。
- 人間承認状態: 承認済み（§8 を最小修正、構造は変えない）。

### 正本反映

- 反映 docs: `docs/architecture.md` §8。
- 変更前: 「extractor は in-memory の ExtractionResult を作り件数検証するところまで実装済み。SQLite writer は未実装。」
- 変更後: 「extractor は ExtractionResult を作り、叙述文（narration）と台詞（line）・話者属性（speaker / race / faction / voice_type）を中心 DB へ書く。engine は叙述文と台詞を AI 翻訳し、台詞は話者属性からのペルソナ口調指示を注入する。本文翻訳の進捗は runtime events で frontend へ push する。」
- 根拠 active plan: 本 plan の実装範囲・implementation-module 成果物。判断履歴は `docs/changelog.md` の T2 entry。

### 作業 commit / local merge / merge 後検証 / completed 移動 / merge 結果 commit

- 後続の手順で記録する。

## Outcome

- implementation-module 完了、finalization-module で §8 を人間承認のうえ反映。次は作業 commit → master へ local merge → merge 後検証 → completed 移動 → merge 結果 commit。
