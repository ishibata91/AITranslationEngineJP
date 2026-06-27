# Task Plan: 2026-06-23-record-type-translation-expansion

- `workflow`: work
- `status`: finalization-module 進行中（2026-06-27）。AI 翻訳の実走を実画面で確認済み（固有名フェーズ→本文フェーズの全段が hy-mt2-7b で動作、種別バッジ・口調・固有名注入を目視）。正本化判断は architecture.md・er.md とも反映不要と確定。次は作業 commit→local merge→completed 移動
- `task_id`: 2026-06-23-record-type-translation-expansion
- `request_summary`: 現在対象の会話・本以外、残りの概念モデルに登場するレコード種別も翻訳対象にする。固有名を本文より先に翻訳し、他レコード翻訳時に mod 用辞書として扱えるようにする。各レコード種別ごとに翻訳スタイルを指定できるようにする（クエスト→日記風、WEAP→武器、といった形）。
- `goal`: 概念モデルに登場するレコード種別を本（BOOK:DESC）・会話台詞（INFO:NAM1）以外へ広げて翻訳でき、固有名を本文翻訳より先に確定し、既訳は `master_term`（権威訳の横断永続辞書）・AI 訳は `proper_noun`（実行内）へ分けて本文翻訳へ注入し、レコード種別ごとに翻訳スタイル（プロンプト指示）を切り替えて訳せる。
- `source_branch`: `master`
- `source_commit`: `18de51c0`
- `work_branch`: `claude/2026-06-23-record-type-translation-expansion`
- `target_branch`: `master`

## 現状の起点（preparation-module 探索で確認、2026-06-23）

完了定義を矛盾なく固定するため、依頼 3 点に対する現状を記録する。

- 翻訳対象の絞り込み: 中心 DB に書かれて Go engine が訳すのは `narration`（BOOK:DESC のみ）と `line`（INFO:NAM1 のみ）の 2 種別だけ。`tools/extractor/Program.cs:57` の `narrationRecFields = { "BOOK:DESC" }` がハードコードで分類を決める。C# 抽出器は WEAP/ARMO/QUST/MESG など全種別を `ExtractionResult` に持つが、SQLite へは書かれず翻訳パイプラインに乗っていない。
- 固有名の辞書化: `master_term` テーブル（source/dest/category）が存在する。供給は xTranslator 既訳 XML からの FULL 流用（`tools/extractor/MasterTermXmlWriter.cs`）と人名部分形の派生（`internal/engine/termderive.go`）の 2 経路のみ。本文翻訳より先に固有名を確定して辞書化する「フェーズ順序」の機構はない。注入は `internal/engine/engine.go:144,158` の `dict.Apply` で本文へ機械置換される。
- 種別別スタイル: `prompt_template` テーブルは単一行（base_directive + persona_template）で、叙述文・台詞に共通。レコード種別ごとにスタイルを切り替える機構はない。`narration.style` カラムは存在するが翻訳プロンプトへ注入されていない（`internal/engine/prompt.go` に利用箇所なし）。

## 完了定義（preparation-module で固定、2026-06-23）

本 task は、現在対象の会話（INFO:NAM1）・本（BOOK:DESC）以外で、概念モデルに登場する翻訳対象の全レコード（REC:FIELD）を今回まとめて翻訳対象に加える。武器（WEAP）は代表例の 1 つで、対象を武器に絞らない（人間指示、2026-06-23）。差込点や空テーブルを置くだけで振る舞いが観測できない状態は「動く」と書かない。手段の細部（固有名の確定方式、スタイルの持ち方、新テーブル構造、レコード網羅の単位）は design-module で確定するが、完了定義の振る舞い自体は狭めない。

動かす範囲（task 後に検証者が観測できる振る舞い）:

1. 叙述文の全種別を翻訳対象に加える: 各 DESC/DNAM 系（武器・防具・呪文の説明、`MGEF:DNAM`、`QUST:CNAM` ログ・`QUST:NNAM` 目標、`LSCR`、`PERK:DESC` など、概念モデルの叙述文に当たる全 REC:FIELD）を翻訳対象に加え、翻訳実行で訳文が生成される。各 REC:FIELD を文体（説明体・書物体・日記体・世界観断片）へ割り当て、文体ごとのプロンプトで訳す。
2. 固有名の全種別を先に確定して辞書化する: 各 FULL 系（武器・防具・呪文・任務・地名・本タイトル・人名・種族名・勢力名など、概念モデルの固有名に当たる全 REC:FIELD）を、本文翻訳より先に確定する。既訳は権威訳辞書（`master_term`）を流用し、既訳の無いものは AI で翻訳して `proper_noun`（中心 DB、実行内スコープ）へ書く（`master_term` へは昇格しない）。権威訳と AI 訳を合わせて本文（叙述文・台詞）へ機械置換で注入する。
3. 定型句を翻訳対象に加える: `RNAM`（起動動作）・`MESG:ITXT`（メッセージのボタン）・`WOOP:TNAM`（龍語の語義）を翻訳対象に加える。
4. 台詞の残りを翻訳対象に加える: `INFO:RNAM`（選択肢上書き）・`DIAL:FULL`（選択肢既定文）を翻訳対象に加える。
5. 文体別スタイルが効く: レコード種別（REC:FIELD）ごとに訳し方（叙述文の文体・固有名・定型句・台詞の口調）が選ばれ、叙述文は文体ごとにプロンプト指示が変わる。既存の本（書物体）・会話（口調ペルソナ）の訳し方は維持される。

無訳片（`WOOP:FULL`、龍語綴り）は翻訳禁止のため対象外とする。

観測点:

- 1〜4 → 実データ。抽出→翻訳実行後、会話・本以外の各 REC:FIELD の行が中心 DB に入り訳文が書き戻ることを確認する。固有名は本文翻訳より先に確定され、既訳の無いものが AI 翻訳されて `proper_noun` へ入り、`master_term` には権威訳だけが残ることを確認する。
- 5 → 実画面と単体テスト。実 app（`npm run dev:wails:run`、`http://localhost:34115`）で抽出→翻訳を実行し、複数の代表種別（武器・呪文・任務・地名など）で文体に応じた訳と固有名の一貫注入を目視確認する。種別（REC:FIELD）→訳し方の割り当てと、文体→プロンプト指示の選択は純粋 IO クラスへ分離して単体テスト（カバレッジ 100%）で確認する。
- 退行防止 → 実画面。既存の本（BOOK:DESC）・会話（INFO:NAM1）の翻訳が壊れないことを確認する。

含まない（除外範囲）:

- 固有名辞書の実行をまたぐ永続化。中心 DB は起動ごとに空にする意向があり永続化方針が未決のため、本 task は同一実行内のフェーズ順序（固有名→本文）に限る。AI 訳の固有名は `proper_noun`（実行内）に留め `master_term` へ昇格しない（方針 A）。永続化方針そのものは本 task で決めない。
- 叙述文・台詞の文中の固有名を自動検出して固有名へ結ぶ「言及」（概念モデル弱点 2、関連 e4/e5）の実装。固有名注入は辞書の機械置換で行い、言及の自動検出は別 task とする。

## close_conditions（観測点で検証する）

- `close-1`: 実データで、会話・本以外の概念モデル全レコード（叙述文・固有名・定型句・台詞の未対応 REC:FIELD）が抽出→翻訳で訳文を持つ。
- `close-2`: 実データで、固有名（全 FULL 系）が本文翻訳より先に確定され、既訳の無いものは AI 翻訳されて `proper_noun`（実行内）へ入り（`master_term` は権威訳専用で AI 訳を持たない）、権威訳と AI 訳が本文へ機械置換で注入される。
- `close-3`: 実画面で、複数の代表種別で文体別プロンプトが効き、固有名が本文へ一貫注入され、既存の本・会話翻訳の訳し方が維持される。訳し方の純粋判定 4 つ（レコード分類・文体割り当て・供給源選別・プロンプト合成）は engine 配下 1 サブパッケージへ集約し、単体テスト（カバレッジ 100%）で確認する。供給源選別は AI 訳が `master_term` を宛先に返さないこと（不変ルール）を含める。

## 軽 / 重判定（preparation-module で固定、2026-06-23）

- 画面が動くか: Y。会話・本以外の全レコード種別の翻訳結果を結果一覧へ表示する。レコード種別の表示・絞り込みや文体別スタイルの指定・確認に、layout・文言・表示構造・svelte 表示コンポーネント・props・story・fixture のいずれかを変える可能性が高い。最終要否は design-module / storybook-module で確定するが、画面が動く前提で進む。
- `docs/architecture.md` 反映が要るか: Y（design-module で確定）。翻訳パイプラインへ取込段（抽出生テーブル→箱別テーブルへ振り分け）と固有名先行フェーズ（固有名→本文の段階順序）が加わる。C# 抽出器は素朴吸い出しに変え（意味判定を持たない）、訳し方の純粋判定 4 つ（レコード分類・文体割り当て・供給源選別・プロンプト合成）を engine 配下 1 サブパッケージへ集約する（横レイヤーは増やさない）。抽出生テーブルと `proper_noun` の新設、`narration.style` の使用開始が、パイプライン段階・依存方向・DB スキーマ・Wails 境界へ触れる。`master_term` は権威訳専用に保ち変更しない。最終要否は design-module 出口で確定する。
- 判定結果: 重 task（画面 Y だけで重 task 確定）。経路は `preparation-module` → `design-module` → `storybook-module`（画面が動くため）→ `implementation-module` → `finalization-module`。

## 後続 task（起動条件付き、2026-06-23）

本 task の対象外として、起動条件が満たされたとき別 task で扱う。

- 固有名辞書の実行をまたぐ永続化（中心 DB を作り直しても固有名訳を残す）。起動条件: 中心 DB の永続化方針が人間判断で確定した後（`project-db-wipe-on-launch-intent`）。`master_term` が中心 DB に同居しつつ永続意図を持つ衝突（起動消去から除外するか別 DB へ分離するか）も、この task で決める。
- 既訳 XML の取り込み機構（大型依存 Mod の固有名を権威訳として `master_term` へ取り込む）。起動条件: 大型依存 Mod の固有名一貫性が AI 訳＋実行内注入だけでは不足と判明した時、または既訳 XML 取り込みを人間が要求した時。
- 叙述文・台詞の文中の固有名を自動検出して固有名へ結ぶ「言及」（概念モデル弱点 2、関連 e4/e5）。起動条件: 辞書の機械置換注入だけでは固有名一貫性が不足と判明した後。

## Branch Status（preparation-module 分）

- `source_branch`: `master`
- `source_commit`: `18de51c0`
- `execution_branch`: `claude/2026-06-23-record-type-translation-expansion`
- `branch_ready`: Y
- `remote_operation`: `not-performed`

## Artifact Index（design-module で記入、2026-06-25）

- `implementation-scope.md`: 実装範囲・テスト設計（design-module 出口で固定）。
- `docs/changelog.md` 2026-06-25 entry: 人間設計レビュー承認と 3 確認点の判断履歴。
- `design-review.md`: 人間設計レビュー材料（一時）。承認後に削除済み（2026-06-25）。

## Routing Notes（design-module で記入、2026-06-25）

- required_reading: `docs/concept-model.md`、`docs/er.md`、`docs/architecture.md`、`db/migrations/0001_init.sql`（narration）・`0002_persona_dictionary.sql`（line）・`0003_master_term.sql`（master_term）・`0004_prompt_template.sql`、`internal/engine/engine.go`・`dictionary.go`・`tone/`、`tools/extractor/Program.cs`。
- canonicalization_targets: `db/migrations/`（0005 以降の新設テーブル＋seed）、`docs/er.md`（`record_type_master`・`style_template`・`extracted_field`・`proper_noun` の追記）、`docs/architecture.md` §8（取込段・固有名フェーズ・新設テーブル・Wails 境界）。架空構造の先行記述を避け、実在後に追記する。
- validation_commands: `go test ./internal/engine/...`、`go build ./...`、実画面 `npm run dev:wails:run`（`http://localhost:34115`）。

## HITL Status（design-module で記入、2026-06-25）

- detail_spec_hitl: 承認済み（2026-06-25、人間「先へ進めて」で人間設計レビュー通過）。
- storybook_review_loop: 承認済み（2026-06-25）。レビューで設計が「Base ＋ 種別ごとの指示文（directive）」の MECE モデルへ収束。
- frontend_human_review: 承認済み（2026-06-25、人間「おk」）。
- approval_record: 2026-06-25 人間設計レビュー承認＋Storybook 人間レビュー承認。設計の変遷（論理名追加→行編集 select 撤去→プロンプトテンプレートのタブ化→MECE モデル）は changelog 参照。

## 合意済み frontend 保護（storybook-module、2026-06-25）

- 承認済み画面: プロンプトテンプレート（`Screens/プロンプトテンプレート`）。サブタブ [ベース][レコード別]。
  - ベース: Base 翻訳指示（全種別共通の前提）。
  - レコード別: 種別ごとの指示文（directive）を一律に並べる（説明体・書物体・日記体・世界観断片・口調〔`{traits}` 変数つき〕・固有名・定型句）。各 directive は 本文（編集可）＋変数（あれば）＋対象 REC:FIELD（読み取り専用）。
- 表示規則: プロンプト = Base ＋ その REC:FIELD の directive（変数を実行時に埋める）。REC:FIELD → directive の割り当ては固定。編集は Base と directive の指示文だけ。翻訳しない種別は載せない。
- 通常分類へ戻した story: `Screens/プロンプトテンプレート`、`UI Components/DirectiveEditor`、`UI Components/TranslationResultRow`（結果行に種別バッジ追加）。
- 反映先 frontend: `frontend/src/ui/screens/template-editor/`（`TemplateEditorScreen.svelte`・`TemplateBasePane.svelte`・`DirectiveEditor.svelte`・`directive-view.ts`・`directive-presentation.ts`・`template-editor.fixtures.ts`・`template-editor-view.ts`）、`frontend/src/ui/screens/translation-run/`（`TranslationResultRow.svelte`・`translation-run-view.ts`）。
- 変更禁止範囲（後続実装で表示を変えない境界）: 承認済み svelte 表示コンポーネントの構造・props 形・style。implementation-module は state・gateway・タブ状態・directive 保存の配線だけを足す（表示は変えない）。
- storybook-module 内で触れた frontend ロジック（最小）: `TemplateEditorContainer.svelte` から不要になった placeholders 受け渡しを撤去。新 props（activeTab・directives・assignments・onInstructionInput・onTabChange）への本配線は implementation-module。
- 削除: 独立画面のレコード種別マスター一式（`frontend/src/ui/screens/record-type-master/`）。口調/文体/変数なしの分割を directive へ畳んだため。

## implementation-module 実装記録（2026-06-25）

変更（1 行 1 ファイル群）:

- `db/migrations/0006_record_type_translation.sql`: directive（7 指示文 seed）・record_type_master（65 REC:FIELD→box・directive・論理名 seed）・extracted_field・proper_noun・extracted_info_speaker を新設。
- `internal/model/`: `directive.go`（Directive・RecordType）・`extracted_field.go`・`proper_noun.go` を追加。
- `internal/store/`: `directive.go`・`proper_noun.go`・`ingest.go`（extracted_field 読み・batch 投入・line_speaker 解決）・`seed_test.go`（seed 整合）を追加。
- `internal/engine/`: `ingest.go`（Dispatch 純粋・Ingest）・`proper_noun.go`（SelectSupply 純粋・固有名フェーズ）を追加。`engine.go` Run を MECE 化（固有名フェーズ→本文フェーズ、directive 引き当て、辞書 master_term∪proper_noun）。`prompt.go` に FillVariables 追加。`ingest_test.go`・`proper_noun_test.go`・`prompt_test.go` でテスト。
- `internal/api/app.go`: 取込段 Ingest 呼び出し、結果行 recordType、固有名結果区間、GetDirectiveEditing・SaveDirective。pageRows を 3 区間化。
- `tools/extractor/`: `ExtractedFieldSqliteWriter.cs`（全 REC:FIELD 素朴吸い出し）・`SpeakerSqliteWriter.cs`（話者素材＋INFO→speaker 橋渡し）へ置換。`Program.cs` 更新。旧 Narration/LineSpeaker writer 削除。
- `frontend/src/gateway/template-gateway.ts`・`translation-gateway.ts`: directive 編集 API と結果行 recordType を追加。
- `frontend/src/ui/screens/template-editor/TemplateEditorContainer.svelte`: 新 props（activeTab・directives・assignments・onInstructionInput・onTabChange・directive 保存）へ配線。表示コンポーネントは storybook-module 確定のまま変更なし。
- `frontend/wailsjs/`: Wails bindings 再生成（GetDirectiveEditing・SaveDirective・新 DTO）。

最終検証:

- backend（go ツール直実行）: `go build ./...` 通過。`go test ./internal/... ./db/...` 全通過。純粋判定（Dispatch・SelectSupply・FillVariables・ComposePrompt）カバレッジ 100%。lint は format・vet・module 通過。static・arch の残りはベースライン（本 task 前から存在）の既存債務のみで、本 task の新規コードは指摘ゼロ（shadow を 1 件むしろ削減）。
- frontend: `npm --prefix frontend run check` 通過（残る 1 error は node_modules の `@storybook/svelte` 型宣言で既存・本変更外）。`frontend-local` suite（eslint・tsc・knip・boundaries・vitest）通過。
- C#: `dotnet build` 通過、`dotnet test` 20 件通過（ExtractedField・Speaker writer の新テスト含む）。
- 実画面（`npm run dev:wails:run`・`http://localhost:34115`）:
  - directive 編集: レコード別タブが seed の 7 指示文・65 REC:FIELD 割り当て・論理名・口調の `{traits}` 変数を表示。指示文を編集→未保存→保存→リロードで永続を確認（SaveDirective→GetDirectiveEditing 往復）。
  - 結果行種別バッジ: 既存データで「台詞 ・ INFO:NAM1」など recordType バッジを目視確認。
  - 取込段（実データ）: 抽出器 CLI を inigo.esp に実行し extracted_field 8817 件・全 4 箱の REC:FIELD（INFO:NAM1・DIAL:FULL・QUST:CNAM・LSCR:DESC・各種 FULL/DESC など）を確認。Ingest を実コードパスで走らせ、振り分け（叙述文/定型句112・固有名146〔重複排除〕・台詞8545）と line_speaker のstaging 解決（5244）を確認。
- 未観測（環境制約）: 固有名フェーズ・本文フェーズの AI 翻訳実走と、新 REC:FIELD 行に訳文が入った結果一覧表示は、LM Studio（`192.168.0.226:1234`）が本セッションで不通のため end-to-end では未観測。判定ロジックは単体テスト（SelectSupply・プロンプト合成 100%、Run を fakeStore で）で、消費データは実データ取込段で担保済み。finalization-module で AI 復帰後に実走確認を行う。

## finalization-module 記録（2026-06-27）

### AI 翻訳実走の実画面確認（implementation-module の未観測項を解消）

LM Studio（`192.168.0.226:1234`、モデル `hy-mt2-7b`）復帰後に実 app（`npm run dev:wails:run`・`http://localhost:34115`）で end-to-end を観測した。対象 plugin は `Innocence Lost - Quest Expansion.esp`（小型）。中心 DB は辞書（`master_term` 25071）・seed（`directive` 7・`record_type_master` 65）を残し翻訳系テーブルだけ空にして実行した。

- 取込段: `extracted_field` 199 件 →（振り分け）叙述文/定型句 29・固有名 17・台詞 151。合計 197（残り 2 は無訳片/未知で除外）。
- 固有名フェーズ（本文より先行）: 17/17 を翻訳。供給源選別が機能（辞書一致は `master_term` 由来を採用、不一致は AI 訳を `proper_noun` へ）。例: `Honorhall Orphanage`→`オナーホール孤児院`、`Innocence Lost`→`失われた無垢`、`Skyrim`→`スカイリム`。
- ペルソナ生成: 話者 7 件分を生成（固有名フェーズと本文フェーズの間に走る）。
- 本文フェーズ 叙述文: 新規 REC:FIELD（`QUST:CNAM`・`QUST:NNAM`・`MGEF:DNAM`・`MESG:DESC` 等）を翻訳。固有名注入が一貫（`Aventus Aretino`→`アベンタス・アレティノ` が全文一致、`Riften`→`リフテン`、`Dark Brotherhood`→`闇の一党`）。注入カタカナを hy-mt2-7b が保持。
- 本文フェーズ 台詞: 口調が話者差で出る（少年 Aventus=「ぼく」、女院長 Grelod=「〜のよ／〜でしょう」）。固有名も台詞内で一貫。
- 結果一覧 UI: 「197 件」、種別バッジ（箱 ・ REC:FIELD）が新規種別を表示、台詞へ口調バッジ（平明・ぞんざい）、本文行へ固有名注入数バッジを表示。

`docs/changelog.md` 2026-06-27 entry に観測要点を追記。`project-injected-token-fidelity` memory に hy-mt2-7b の保持観測を追記。

### 検証（作業 commit 前、追跡コードは本セッション未変更）

- backend: `go build ./...` 通過。`go test ./internal/... ./db/...` 全通過。
- frontend: `npm --prefix frontend run check` exit 0（1 error は node_modules `@storybook/svelte`、既存）。harness `frontend-local` 全通過。
- C#: `dotnet test tools/extractor.Tests/` 20/20 通過。
- backend lint: 12 件は全て既存ファイル（`engine.go:293` の InsertDerivedTerms〔DeriveMasterTerms 経路、本 task 非変更〕・`termderive.go`・`termxml.go`・`nrc.go`）。本 task の新規/変更ファイルは指摘ゼロ。

### 正本化判断

- 反映対象（finalization-module は `docs/architecture.md` に限定）。判断: **反映不要**。
- 根拠（`docs/architecture.md`）: 層（新規 box/package なし、追加は engine・store・model の既存層内）・依存方向（engine の新規 consumer interface は §4 が許す既存パターンの踏襲）・Wails 境界（Bind 追加は §5 の機構を変えない）のいずれも不変。`feedback-architecture-reflection-structural-only` に従い §8 を churn せず人間承認も求めない。パイプライン段階追加（取込段・固有名フェーズ）と C#→素朴吸い出し/Go→振り分けの責務移動は changelog に記録する（正本ではない）。
- 根拠（`docs/er.md`、finalization の正本反映対象外だが整合を確認）: er.md は概念モデル 10 箱の論理 ER で、スコープ対象外に「マスター辞書・実現方式・schema version」を明記。新表は config（`directive`・`record_type_master`）と実現方式 buffer/staging（`extracted_field`・`extracted_info_speaker`）で概念箱でない。`proper_noun`/`narration`/`line` は既存概念箱。er.md §レコード識別は既に `QUST:CNAM`・`MESG:ITXT` の序数まで織り込み済みで本 task と整合。よって追記不要。
- 人間承認状態: 構造不変につき承認不要（memory 準拠）。設計時 Routing Notes の canonicalization_targets（er.md・architecture.md §8）は「実在後に追記」の暫定で、実装後に er.md の実スコープと構造不変判定により反映不要へ確定。
- 追補（merge 後、人間指示「他の乖離を直してくれ」を受領）: 構造は不変だが、`docs/architecture.md` §8 の「extractor は narration/line を書く」という現在状態の記述が実装と乖離していた。人間指示を承認として `architecture.md` を反映へ切り替えた。§3（engine 責務に取込段・固有名フェーズを反映）・§8 現状記述（extractor は `extracted_field` を書き engine の取込段が振り分ける）を修正し、§8 へ本 task の移行 entry（全 REC:FIELD 化・取込段・固有名フェーズ・directive・口調供給の prompt_template→口調 directive 移行・新 Bind）を追加。`er.md` は概念モデルの段でテーブルそのものではないと人間が確認し、据え置き。`docs/architecture.md`・本 plan・changelog の反映は同一 commit に含める（hash は git log 参照）。

### 作業 commit / local merge / completed 移動 / merge 結果 commit

- 作業 commit: `f05fddc9`（work branch `claude/2026-06-23-record-type-translation-expansion` 上、44 ファイル、2723 挿入・498 削除）。ノイズ（`.DS_Store`・`.claude/skills/presentation/SKILL.md`・`cmd|db|docs|scripts|tools/CLAUDE.md`・`docs/.obsidian/`）は本 task 外として除外し未コミットのまま残した。
- local merge: `git merge --no-ff` で master へ取り込み。merge commit `2af24c53`。merge-base が master 先端（`18de51c0`）と一致し divergence なし、conflict なし。
- merge 後検証: `go build ./...` 通過、`go test ./internal/... ./db/...` 全通過（cached＝merge 後ツリーは検証済み f05fddc9 と同一）。
- completed 移動: `git mv docs/exec-plans/active/2026-06-23-record-type-translation-expansion docs/exec-plans/completed/`。
- merge 結果 commit: 後述の hash（completed 移動と本記録を含めて master 上で commit）。
- remote 操作: 行わない（push・tag・remote delete なし）。
- 残留リスク: 固有名フェーズの大量翻訳は小型モデルでは時間がかかる（hy-mt2-7b で 1 件あたり数〜十数秒）。本 task の範囲外。後続 task（固有名辞書の永続化・既訳 XML 取り込み・言及自動検出）は本 plan の「後続 task」節に起動条件付きで記録済み。
