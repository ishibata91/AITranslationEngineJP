# 変更・判断履歴

正本（`requirements.md`、`system_requirements.md`、`architecture.md` など）には現在の状態だけを書く。
「なぜ変えたか」「何を落としたか」などの判断履歴は本ファイルに残し、正本へ混ぜない。
新しい entry を上に追加する。1 entry は date 見出しで区切る。

## 2026-06-18 マスター辞書 T3 増分: 人名の部分形の派生（名のみ・短名を辞書化）

### 変更

- `internal/engine/termderive.go`（新規）: 人名の部分形を 3 種で派生する純関数 `DeriveTerms` を追加。shrt（`NPC_:SHRT` の短縮別名通過）・byname（` the ` を含む名の前部と Dest 末尾カタカナ連）・two（base ゲームの空白 2 語姓名を中黒 2 語で整列）。安全フィルタ（用法比 lc/uc・称号・種族/ハウス語・縮約素体・純カタカナ・最小長 4・base 衝突 skip・two の base ゲーム限定）を同関数に持つ。副作用なし。
- `internal/engine/termusage.go`（新規）: 会話文の英語原文から各英単語の用法分布（lc=小文字始まり一般語用法 / uc=文頭以外の大文字始まり固有名用法）を作る純関数 `BuildUsage` を追加。
- `internal/engine/termderive_test.go`・`termusage_test.go`（新規）: 純粋ルールの全分岐単体テスト。新規ルール関数のカバレッジ 100%。
- `scripts/dict/derive-master-terms/main.go`（新規）: ビルド時コマンド。xTranslator 英日 XML を解析し、用法分布を作り、純粋ルールを呼び、base 衝突を除いて `master_term` へ `category="derive:<種別>"` で追記する。`db.Apply` で schema を ensure し、`INSERT OR IGNORE` で二重追記を防ぐ。
- `scripts/dict/derive-master-terms/main_test.go`（新規）: XML 解析の単体テストと、temp DB へ XML→派生→追記・base 衝突 skip の結合テスト。
- 実行時の置換器 `internal/engine/dictionary.go`・`engine.go`、テーブル `master_term` は無改造。`loadDictionary` が category を問わず全件を読むため派生行を自動で取り込む。

### 判断

- 派生規則は副作用の無い単一の純粋ルール（`DeriveTerms`）へ分離し、ユニットテストカバレッジ 100% を基準にした（人間が固定した基準）。XML 解析・DB 書込の I/O 配線はルールの外（ビルド時コマンド）へ出し、結合テストで見る。
- 置き場所は判定が属する言語（置換器が Go なので Go）に合わせ、ビルド時生成で `master_term` へ焼く（人間承認の構成）。永続した派生行は category 印で目視・差し戻しできる。
- base 辞書は現行の C# extractor のまま。派生は base 書き込み後に走る Go コマンドが追記する 2 段構成。同じ既訳をより単純に得られ、純粋ルールを Go に置けるため。
- 由来種別は既存 `category` 列に `derive:<種別>` として持たせ、スキーマは変えない。実行時の照合は source だけを見るため影響しない。
- two（姓名分割）は base ゲーム XML 限定にした。patch/mod の Source/Dest 対応ずれ（USSEP で観測）による誤訳を避けるため。
- 安全フィルタは観測した失敗の手書き除外でなく、用法分布と語形による構造判定にした。一般語（`Master`・`Blood`・`Mine`）・種族語（`Imperial`・`Nord`）・称号（`Lord`・`Captain`）・縮約（`Aren`）を捨て、固有名（`Grelod`・`Mercer`）を残す。`Imperial` は用法比だけでは残る側だが種族/ハウス語集合で捨てる（地の文の誤置換回避）。
- two の category 形容語フィルタ（creature ラベル混入抑制）は本実装に入れない。検証済み数値（破壊型過剰置換 0・被覆 99.9%・held-out 汎用性）はフィルタ無しで得たもので、追加は未検証のため。起動条件は「ハーネスで過剰置換が減り被覆が落ちないと実測したとき」。

### 残課題

- 機械置換は実 DB（base 24,554＋派生 517）で観測済み: `Grelod`→`グレロッド`（名のみ、原問題解消）、`Mercer Frey`→`メルセル・フレイ`（最長一致）、`Mercer`→`メルセル`（姓のみ）、一般語 `master` は無置換。AI 無しで `engine.NewDictionary`+`Apply` を実 DB に対し実行して確認。
- 実 app の end-to-end も観測済み（実 LLM、2026-06-19）: 「Innocence Lost - Quest Expansion.esp」の台詞 121 件を gemma-4-12b で翻訳し、`Grelod` 28/28→`グレロッド`、`Riften` 4/4→`リフテン`、崩れ 0。観測前は前回（派生なし）の訳で裸 `Grelod` が `グレロド`/`グレロッド` に揺れていたのが、派生辞書で全て `グレロッド` に揃った。plugin はパス直接入力欄から手入力（人間補助不要。先の「ネイティブダイアログのため人間依頼」は誤り）。
- 所見（モデル依存・新規発見）: 機械置換は AI 前段で確定訳語を注入するが、end-to-end の一貫性は AI が注入カタカナを保持するかに依存する。shisa-v2.1-qwen3-8b は注入済み `リフテン` を `リヴェン` 等へ書き換え（`Riften` 0/4・`Grelod` 24/28 のみ保持）、gemma-4-12b は完全保持。派生のバグではなく（base 名 `Riften` も崩れる）、弱い小型モデルが注入トークンを保持しない問題。翻訳モデルは注入トークンを保持できる能力のものを選ぶ。
- 派生はビルド時 1 段（`go run ./scripts/dict/derive-master-terms --sqlite db/aitranslation.dev.sqlite3`）。中心 DB は wipe されず writer が `INSERT OR IGNORE` のため、app の再抽出で派生行は消えず runtime 無改造で base+派生が効く。
- `docs/architecture.md` 反映の要否は finalization-module で判断する（本増分は engine 内の純粋ルール追加とビルド時コマンドで、Wails 境界・実行時の層構造は変えていない）。
- `go build ./...` は既存の壊れ（`scripts/test/seed-system-test-db/main.go` が存在しない `internal/repository` を import）で失敗する。本増分の範囲外。

## 2026-06-15 結果一覧の N+1 廃止＋keyset cursor ページング（T2 の効率課題の恒久対応）

### 変更

- store: 台詞ごとに話者を引く `LoadLineSpeaker`（1 台詞 2 クエリ）を、台詞 id 群を一括で引く `LoadLineSpeakers(lineIDs)→map`（IN 句・host parameter 上限回避の chunk）へ置換。keyset 範囲取得 `NarrationsAfter`／`LinesAfter`（`WHERE id > ? ORDER BY id LIMIT ?`）と `CountNarrations`／`CountLines`、共通 helper `query.go` を追加。未使用の `ListNarrations`／`ListLines` を削除。
- engine: `LineStore` を一括取得へ差し替え、`LineDirective`／`linePersona`（per-line）を `LinePersonas(lineIDs)→map[int64]Persona` へ置換。`Run` はループ前に話者を 1 度だけ一括取得し、ループ内の個別 DB 問い合わせを廃止。
- api: `ResultPage{Total,Results,NextCursor,HasMore}` と `ListResultsPage(cursor,limit)`、連結列のページ範囲を決める `pageRows`、cursor 解析（`""`／`n:<id>`／`l:<id>`）を追加。`buildResults`（全件）を廃し、ページ内台詞ぶんだけ口調を一括生成する形へ。`RunExtractAndTranslate` は全件 `Results` を返すのをやめ件数要約だけ返す。
- frontend: gateway に `listResultsPage(cursor,limit)`、container に keyset state（cursor 履歴・ページ index・total・nextCursor・hasMore・ページサイズ 50）と前へ/次へ。ページャ表示 `ResultsPager`（順次送り・端で無効化・現在ページ番号）を追加し、件数バッジを総件数へ（Storybook 人間レビュー承認）。`TranslationResultRow` の表示形（コンパクト行・口調チップ・展開）は不変。

### 判断

- ページング方式は keyset（cursor）を採用し LIMIT/OFFSET を不採用にした（人間設計レビューで確定、ページサイズ 50）。結果一覧は閲覧中に行が増減しない静的集合だが、数万件の深い位置を `WHERE id > ?` の index 走査で取れる keyset を選んだ。仮想スクロールは全件を frontend へ載せ payload を増やすため不採用。
- 話者一括取得は map 返し（`LoadLineSpeakers(lineIDs)→map`）にした。所属勢力が話者に対し 1 対多のため `ListLines` への JOIN は行が増殖し group_concat 等が要る。map なら所属勢力を `[]string` のまま保て、T2 の責務分担（store は識別子・事実、engine が口調へ解釈）を維持できる。同じ map を表示と翻訳で使い回す。
- N+1 は per-line メソッドを削除して一括メソッドへ置換し、「N+1 を表現できない」形にした（特殊対応の追加でなく機構の置き換え）。
- `RunExtractAndTranslate` は全件結果のインライン返却をやめた。数万件を実行応答に載せず、実行後・起動時・ページ送りを `ListResultsPage` の 1 経路に統一した。
- `docs/architecture.md` への反映は不要と判断（§5 Wails 境界・§7 ディレクトリ正本・§8 現在の状態の構造を変えていない。keyset と一括 persona は層内の内部精緻化、`ListResults`→`ListResultsPage` は既存 Bind 境界内）。

### 残課題

- 固有名解決・マスター辞書は T3、口調ルールの精緻化・編集 UI は T4（いずれも対象外）。
- 翻訳実行中（`engine.Run`）の話者一括取得は実装済みだが、翻訳は AI latency 支配のため最適化の主目的は表示経路（`buildResults`）にある。
- 本 task の実画面検証はローカル LLM stub で pipeline 全体を通した（121 台詞・3 ページ）。AI 翻訳の本番実行は利用者の OpenAI 互換 provider で行う。

## 2026-06-14 T2 ペルソナ口調 pipeline（台詞抽出→話者解決→口調注入翻訳→進捗・口調差を画面で観測）

### 変更

- extractor（C#）に台詞（INFO:NAM1）と話者属性（speaker / race / faction / voice_type）の SQLite 書込を追加。INFO の話者 FormKey を LinkCache で NPC へ解決し、種族・声型・所属勢力の EditorID を書く（`LineSpeakerSqliteWriter`）。
- engine（Go）に台詞翻訳とペルソナ口調生成を追加。話者の声型/種族/勢力 EditorID から口調 traits を引く最小ルール（`persona_rule.go`）と、口調指示文を組む `buildPersonaDirective`／チップ用 `buildPersonaLabel`。provider の `Translate` に directive 引数を足し system prompt の base 指示後段へ注入。`Run` に本文翻訳の進捗 callback。
- api に本文翻訳の進捗 runtime events（extract / translate）と、結果へ口調指示文・口調要約を載せる `ResultView`／`ListResults` を追加。
- frontend に本文翻訳の進捗バー（`TranslationProgress`）と結果行のコンパクト化（口調チップ＋展開）、進捗 event 購読を追加。db migration `0002`（line / speaker / race / faction / voice_type / line_speaker / speaker_faction）。
- `architecture.md` §8「現在の状態」を現状へ更新（extractor が台詞・話者も書く、engine が台詞翻訳・ペルソナ・進捗を持つ）。

### 判断

- task の完了定義を「縦切り（観測可能な成果）」へ置き直した。当初の「翻訳プロンプトへ差込点を 1 本通す（単体テストで確認）」は seam 層を完了と呼ぶ弱い条件で、実 mod で観測できる成果が無かった。実 mod `Innocence Lost - Quest Expansion.esp` を実行画面から流し、台詞抽出・口調差・進捗・翻訳を実画面で観測することを完了条件にした。
- 固有名（辞書）解決は本 task から外し、マスター辞書 task（T3）へ移送。T3 の依存「T2 の辞書解決の差込点」は無効化（T3 が自身で差込点を作る）。
- 責務分担: 事実の抽出は extractor（C#）、口調などの解釈は engine（Go）。extractor は識別子・事実（EditorID）だけ書き、口調 traits は engine の最小ルールで与える。ルールの永続化と編集 UI は T4（対象外）。
- 結果行 UI は数万件・ページングに耐えるコンパクト 1 行＋口調チップ（展開で全文）にした（Storybook 人間レビュー承認）。口調差は一覧のままチップで観測。
- architecture 構造（§1〜7）は変えていない。engine 責務（辞書解決・ペルソナ生成）と runtime events は §3・§5 に既記載で、追加スキーマは `er.md` に既定義。§8 の現状記述だけ人間承認のうえ更新。

### 残課題

- 固有名解決・マスター辞書（proper_noun 抽出、line_mention e5、name 関連 e8/e13/e14）は T3。ルール・プロンプト編集 UI は T4。
- ペルソナ口調ルールは engine 内の固定最小 1 系統。声型 EditorID の網羅と気質テキストの精緻化、結果一覧の仮想スクロール／ページングは後続で扱う。
- AI 翻訳の本番実行は利用者の OpenAI 互換 provider で行う。本 task の検証はローカル stub で pipeline 全体を通した。

## 2026-06-14 T1 後の architecture.md との構造差異を整合

### 変更

- keyring secret store を `internal/repository/` から `internal/store/secret/`（package `secret`）へ移動。`architecture.md` §3・§7 の「secret 子に置く」に合わせた。`.go-arch-lint.yml` の component を `repository` → `secret` に更新。
- `db.Apply` に schema version の読み込み時検査を追加。DB の `user_version` がアプリの想定 migration 数より新しければ適用せずエラーにする（`architecture.md` §6「Go は読み込み時に version を検査する」を実装）。
- `architecture.md` §4: 多態の port は `provider` 1 つのみと明記し、`store` 用の狭い interface（consumer 側・実装 1 つ・単体テスト用）は port ではない切り離しとして許容を追記。
- `architecture.md` §5: runtime の閉じ込め先を `bootstrap` と `api`（Bind 公開面）に明記。`api` が runtime を進捗 push とファイル選択ダイアログに使うことを許容し、下位層へは漏らさないと固定。
- `architecture.md` §7: `db/` に migration 適用（`db.Apply`）を追記。`store` が起動時に委譲する旨を記載。

### 判断

- T1 実装が `architecture.md` と食い違った 3 点を、コード修正と doc 改訂に振り分けて整合した。
  - keyring 場所: doc が明示指定（secret 子）。コードを doc に合わせて移動。
  - store の狭い interface: テスト容易性の利益が大きく、design-module のテスト設計（engine を mock 越しに試す）と整合する。doc を実態に合わせて改訂。
  - migration 適用の場所: ユーザー指示「migration とリポジトリは分けて」で `db` パッケージへ分離済み。doc に明記。
- Wails runtime の `api` 直接利用は、§2 図と §3 が「`api` が runtime events を push」と示し、§5 の閉じ込め先 adapter に Bind 公開面（`api`）が含まれるため、乖離ではないと確定。§5 を明文化して曖昧さを除いた。
- 残る差異は意図的な未実装（provider 3 系統・engine の重複排除/辞書/ペルソナ/XML・進捗 push）で、後続タスクで埋める。§8「現在の状態と移行」の陳腐化は別途更新する。

### 検証

- Go test 緑、backend lint（format/vet/static/arch/module）0 issues。store の version 検査込みで store test 緑。

## 2026-06-14 T1 最小縦切り（抽出 → 翻訳 → DB → 画面）を実装

### 変更

- backend（Go）を初実装。`internal/model`（Narration）、`internal/store`（sqlx ＋ modernc.org/sqlite、migration 適用）、`internal/provider`（Translator port ＋ OpenAI 互換実装）、`internal/engine`（未訳を翻訳し仮訳で書き戻す手続き）、`internal/api`（Wails Bind 公開面）、`internal/bootstrap`（composition root）、`main.go`（Wails entry）。
- `db/migrations/0001_init.sql`：narration テーブルの DDL（C#↔Go 契約 1 本）。`db/migrations.go`：embed して公開。
- `tools/extractor`（C#）に `NarrationSqliteWriter` と `--sqlite` モードを追加。BOOK:DESC を narration へ UPSERT。
- frontend を daisyUI で再構築。Tailwind v4 ＋ daisyUI v5 の独自テーマ `dovahkael`、汎用部品（Field/TextField/SelectField/FileSelectField/StatusBadge）、画面 `TranslationRunScreen` ＋ container、gateway、`main.ts`/`App.svelte`。
- lint 整備：`.go-arch-lint.yml` 新設（新層の依存方向）、`.golangci.yml` の static 違反解消、frontend eslint で生成 `wailsjs` を除外、`wails-boundary.test.mjs` で Wails 境界を検証。
- `.gitignore`：`db` 全体無視を `db/*.sqlite3*` に絞り、`db/` の source を追跡。

### 判断

- 叙述文 1 種は `BOOK:DESC`（書物本文）。装備 DESC への拡張は `TranslationCounts.Enumerate` のフィルタ追加だけで済む。
- provider 接続情報（endpoint/apiKey/model）は永続化せず画面から都度渡す。API キーなしの OpenAI 互換（LM Studio 等）に対応するため、キーが空のとき Authorization を付けない。base URL は `/v1` 配下へ正規化（`http://127.0.0.1:1234` でも届く）。
- 抽出は Go の api が C# extractor を `dotnet run` で子プロセス起動し、続けて engine を呼ぶ同期手続き。進捗 push は対象外。
- AI 翻訳は訳状態 3（仮訳）で書き戻す。
- 起動時に中心 DB の現状を読み込み、前回の結果を画面に出す。

### 検証

- TDD：provider（/v1 正規化・auth・getModels・翻訳）、store（migration・未訳取得・dest 更新）、engine（仮訳書き戻し・provider エラー伝播）、api（status ラベル・DTO・extractor 引数）、C# NarrationSqliteWriter（BOOK:DESC 書き込み・冪等）を失敗テスト先行で実装。
- Go test 緑、backend lint（format/vet/static/arch/module）緑、frontend lint（eslint/tsc/knip/boundaries）緑、C# 17 テスト緑、build-storybook 緑。
- 実 app（`dev:wails:run`、localhost:34115）で end-to-end を目視確認。OpenAI 互換モック（`127.0.0.1:1234`）に対し、getModels でモデル選択、Dawnguard.esm から 65 件抽出 → 翻訳 → SQLite → 見開き対訳表示まで動作。LM Studio を同 endpoint に立てれば同経路で実訳になる。

### 残課題

- ファイル選択ダイアログ（Wails OpenFileDialog）は実装済みだが、ネイティブダイアログのため自動 UI テストは未。
- 大量レコードの同期翻訳は進捗表示が無く待ち時間が長い（進捗 push は後続）。
- 書物本文の HTML 様タグ（`<font ...>`）の扱いは未整理。
- フォントは Google Fonts CDN。デスクトップ app 用の self-host は後続。
- greenfield 未配線の `diagnostic`・`shell-state`・`pino` は knip ignore で保持（将来配線で解消）。

## 2026-06-14 ER 設計（抽出入力）の正本 er.md を新設

### 変更

- `docs/er.md`: 新設。`concept-model.md` の箱（抽出入力）を `SQLite` の物理テーブルへ写す ER 設計。テーブル定義・関係・concept-model 対応・既知の論点を記述。
- `docs/index.md`: `er.md` を Read Order・Directory Contract・Choose The Right Record に登録。

### 判断

- スコープは抽出入力（`concept-model.md` の 10 箱と関連 e1〜e14）に限定。マスター辞書・ペルソナルール・翻訳ジョブ/結果キャッシュ・schema version 管理は対象外（あとで別途追加）。
- 概念モデルから外れない。テーブルは `concept-model.md` の箱と 1 対 1。箱を統合せず、属性（`人称`・`口調`・`背景`・`性質` を含む）も落とさない。
- 実現方式を ER に持ち込まない。重複排除のタイミング、属性の充填時期、永続化の有無は `concept-model.md` L7 のとおり実現方式の責務とし、ER は構造だけを固定する。
- 正規化は根拠を明示する。多対多と可変多重度（e4/e5/e6/e7/e8/e10/e14）は連関テーブル（第1正規形）、1 対多（e1/e2/e3/e9/e11/e12/e13）は FK 1 本。訳の単位の分離は更新異常の除去。レコード識別の分解は第1正規形と xTranslator 出力要件。form_id と edid の同居は出力自己完結のための意図的冗長。
- レコード識別キーは xTranslator String 行の `(plugin, form_id, rec, field, ordinal)`。`status` は xTranslator `Status`（0-4）を踏襲（`references/xtranslator_ref.md`）。
- 実 SQL DDL の正本は `db/` migration（`architecture.md` §7）。`er.md` は論理 ER 設計に限定し、DDL を二重に持たない。

### 経緯

- 初版で `配置`・`叙述文`・`台詞`・`無訳片` を `extracted_string` 1 テーブルに、`固有名`・`定型句` を `translation_unit` 1 テーブルに統合し、重複排除責務を `engine` に寄せる実現方式を持ち込んだ。これは概念モデルから外れる逸脱で、人間指摘により 10 箱 1 対 1 へ作り直した。

### 残課題

- 初版。ファイルレビューで確定する。
- 言及（e4/e5）の検出方式、純汎用台詞の話者群の口調決定は実現方式で決める（`concept-model.md` 弱点）。
- 対象外テーブル（辞書・ルール・ジョブ・schema version）は別設計。

## 2026-06-14 tech-selection.md の責務外記述を除去（採用技術へ純化）

### 変更

- `docs/tech-selection.md`: §2・§3・§4・§7 から、採用技術でない記述（データ配置、C#↔Go 設計、責務・依存方向、別プロセス構成、抽出の挙動・制約、観測ログ出力先）を除去した。`SQLite`・`sqlx`・SQL migration・`pino`・`log/slog`・`C#/.NET`＋`Mutagen`・`.NET 8` などの採用技術そのものは残した。

### 判断

- `tech-selection.md` の責務は `index.md` L34 と `core-beliefs.md` §2 で「採用技術と品質基盤＝実装技術の選択」と定義されている。データ配置・内部境界・依存方向は `architecture.md`、観測ログ出力先は `observability-logging.md` の責務で、`core-beliefs.md` §3 が「同じ責務を複数文書で別定義している状態」を除去対象としている。
- §4 永続化・§7 抽出基盤の架構記述は、同日の「アーキテクチャ再構築」entry で `tech-selection.md` に追加したが、`architecture.md` と重複していた。責務分離のため `architecture.md` 側へ一本化し、`tech-selection.md` からは除去した。
- 除去した記述はすべて他の正本に既出で、情報損失はない。対応は次のとおり。
  - DB が持つ内容（抽出入力・マスター辞書・ペルソナルール・翻訳ジョブ）= `architecture.md` §3。
  - C# extractor 直書き・中間形式なし、SQL schema が C#↔Go の唯一の契約、migration 適用責務と extractor の冪等 ensure = `architecture.md` §1/§6。
  - 上位層へ driver 固有 API を漏らさない = `architecture.md` §4。
  - 別プロセス構成・構造体モデル準拠・Data folder の明示パス指定・macOS 実行可・抽出結果の SQLite 書き込み = `architecture.md` §1/§6。
  - backend/frontend 観測ログの出力先（`stderr`／browser console）= `observability-logging.md` §1。

### 残課題

- なし。採用技術の選択内容は変えていない。

## 2026-06-14 アーキテクチャ再構築（データ中心＋手続き、Go 維持＋SQLite 境界）

### 変更

- `docs/architecture.md`: 旧層構成（`UseCase` / `Service` / per-entity `Repository` / `Presenter` ＋ 厚い手動 DI）前提を破棄し、データ中心かつ手続き中心の骨格へ全面書き換え。Mermaid コンポーネント図、各箱の責務、依存方向、Wails 境界、C#↔Go 境界（SQLite 契約）、ディレクトリ正本を記述。
- `docs/tech-selection.md`: §4 永続化を SQLite 中心へ書き換え（抽出 sink を SQLite 正本に、JSON 中間形式を廃止、SQL schema を C#↔Go 契約に）。§7 として翻訳対象抽出基盤（C#/.NET ＋ Mutagen）を新設。§6 に抽出ツールの `xUnit` 検証を追記。公式参照に Mutagen を追加。

### 判断

- 概念モデルが示す実体（中心は Skyrim データ、翻訳は一本の手続き）に対し、旧層構成は過剰と判断。層を薄くして間接化を削る。
- engine の runtime は Go を維持する。Wails / Svelte / 既存 harness を温存するため。
- C#↔Go の受け渡し境界は SQLite とする（案1 ＝ C# extractor が SQLite へ直接書く）。境界専用の JSON 中間形式は持たない。理由: 旧 JSON 境界は xEdit Pascal が翻訳ロジックを載せられない制約に由来する。Mutagen は通常の .NET ライブラリ呼び出しで、tech-selection が既に「入力データを SQLite に持つ」方針のため、境界を SQLite に寄せれば境界専用形式を作らずに済む。
- 抽象（port）は AI provider の境界 1 つだけに置く。実装が 4 系統（Gemini / xAI / OpenAI 互換 / Claude）に分かれる唯一の箇所のため。
- 手動 DI は composition root 1 箇所へ集約する。
- 抽出の検証は C# テスト（`CountParityTests` / `ModelInvariantTests`）へ移管済み。Python の `validate_extraction.py`・`compare_counts.py` は重複のため削除済み。
- Mutagen は macOS でも動くことを公式 docs で確認した（`GameEnvironment.Typical.Builder().WithTargetDataFolder()` で registry 自動検出を回避する）。「macOS では Mutagen を動かせない」という当初の想定は誤りで撤回した。
- engine 内部のパッケージは approach A（`engine` / `model` / `store` / `provider`）を採用する。
- `.NET` 統合（Go を捨て Mutagen と engine を 1 プロセスへ統合する案）は今回不採用。Go 維持で進める。

### 残課題

- engine、store、provider、api、bootstrap、frontend 配線の実体は未実装。再構築は `docs/exec-plans/active/` の active plan で進める。
- SQL schema（中心データの具体テーブル）は未設計。concept-model の箱を SQLite テーブルへ写す設計を plan で詰める。
- C# extractor の SQLite writer は未実装（現状は in-memory `ExtractionResult` の件数検証まで）。

## 2026-06-13 業務要件2・3 のシステム要件を一部確定（対応策の発散と絞り込み）

### 変更

- `docs/system_requirements.md` 業務要件2「単語の一貫性」: `TBD` から、確定部分（Mod 横断マスター辞書、用語特定はレコード固有名詞の機械抽出）と未確定部分（訳語供給、漏れ語対応、適用方式、検証）を分けて記述。
- `docs/system_requirements.md` 業務要件3「NPC の口調」: 「属性と会話履歴から AI でペルソナ生成」から「属性からルールベースで生成」へ方針変更。表現（構造化属性＋翻訳ディレクティブ）、変換（テンプレート・機械的）、永続境界（ルール集を永続・翻訳前設定／個々ペルソナ非永続）、適用（プロンプト注入）を記述。

### 判断

- 一貫性のスコープは Mod 横断（永続マスター辞書）を選択。ジョブ内・Mod 内は不採用。
- 用語特定は、誤検出が無く既訳ヒット率が高いレコード固有名詞の機械抽出のみを採用。AI 抽出・頻度抽出・辞書マッチによる漏れ語対応（第2層）は保留。
- 頻度抽出は単言語では対訳を出せず、対訳コーパスがある場合の統計アラインメントでのみ対訳を出せると整理。用語特定（どの語を揃えるか）と訳語供給（訳語をどう得るか）は別軸と確認。
- 業務要件3 は機械的抽出（属性 → ルール）をまず採用し、AI 生成と会話履歴解析は保留。
- ペルソナ表現は構造化属性（内部表現）＋翻訳ディレクティブ（適用形）の 2 段とし、変換はテンプレートで機械的に行う。
- 機械的抽出ではペルソナを NPC レコードから都度導出できるため、個々ペルソナは永続化しない。永続資産は「属性 → 翻訳指示のルール集」とし、翻訳前にユーザーが設定可能とする。過去構想の `master-persona`（個々ペルソナを永続・編集）とは性質が異なる。
- VoiceType（声タイプ。Skyrim が音声収録のため NPC を声・性格でグループ化した分類）は属性の中で口調と相関が高く、ルールの主軸候補とする。

### 残課題

- 業務要件2: 訳語供給方式（既訳流用のみか AI 併用か）、本文・会話中の漏れ語対応、辞書の適用方式、一貫性検証が未確定。
- 業務要件3: 使用属性の選定、属性の分類、ルール合成の衝突優先順位を概念モデルで整理する。
- ペルソナ属性・Skyrim 構造の概念モデルの置き場が未定（`index.md` Read Order の `skyrim-structure-model.md` は実体なし）。

## 2026-06-13 screen-design 廃止、画面の正本を Storybook へ。tech-selection に Storybook / Tailwind / daisyUI

### 変更

- `docs/screen-design/` を削除（README.md、design-system-ethereal-archive.md、code.html、screens/ 配下すべて）。
- `docs/tech-selection.md`: フロントエンドに Tailwind CSS、daisyUI、Storybook を追加。CSS framework 不採用の行を削除。公式参照 3 件を追加。
- `docs/index.md`: Read Order・Directory Contract・Choose The Right Record から screen-design を削除。画面・表示の設計は Storybook（`frontend/`）と明記。
- `.claude/skills/`: design-module、storybook-module、finalization-module、coding-protocol、fix-decision、investigation-module、implementation-module、diagramming の 8 件を「画面の正本 = Storybook」へ作り替え。
- `docs/exec-plans/active/README.md`、`templates/work-plan.md`、`templates/task-folder/{README.md, plan.md, detail-spec-diff.md}`: `screen-design-diff` と「Storybook 後画面設計差分整合」前提を Storybook 正本へ付け替え。
- memory（`feedback_boundary_responsibility_separation`、`feedback_storybook_module_trigger`、`feedback_implementer_no_agent_split`、`MEMORY.md`）: 画面設計参照を Storybook へ更新。

### 判断

- 画面・表示の設計判断の置き場を Storybook の story と svelte コンポーネントへ移す（ユーザー選択「Storybook を正本にする」）。
- `画面設計差分`（`screen-design-diff.<screen-id>.md`）doc を廃止。design-module は画面表示 doc を作らず、画面表示の設計は storybook-module が story とコンポーネントで直接行う。
- storybook-module の「Storybook 後画面設計差分整合」成果物を廃止。実装範囲を越える画面変更が要る場合は design-module へ差し戻して `実装範囲` を見直す経路に統一。
- finalization-module の docs 正本反映対象を `docs/architecture.md` のみに限定。画面の正本（Storybook）は frontend source として作業 commit に含める。
- fix-decision / investigation-module の画面再現確認の selector 正本を「実装済み画面の `data-testid` またはセレクタ」へ変更。
- AI プロバイダ指定「OpenAI API(llama, lmstudio)」は OpenAI 互換 API でローカル実行（llama、LM Studio）と OpenAI 本体を含む意味と解釈。

### 残課題

- exec-plans templates には screen-design と無関係の旧表記（Codex 名称、`docs/detail-specs/`、`agent-browser`）が残る。本変更の対象外。整理は別途判断する。
- `index.md` Read Order の `skyrim-structure-model.md`、`core-beliefs.md` の `er.md` は実体が無い既存リンク（本変更前から）。

## 2026-06-13 要件文書を業務要件とシステム要件へ分割

### 変更

- `docs/spec.md` を `docs/requirements.md` へ rename（git mv、業務要件の内容は不変）。
- `docs/system_requirements.md`: 新規作成。業務要件 1〜4 に対応するシステム要件を記述。1 = AI 利用（Gemini / xAI / OpenAI 互換 API（OpenAI、llama、LM Studio）/ Claude）、2 = TBD、3 = NPC 属性と会話履歴からペルソナ生成、4 = 機能要件＝業務要件。
- `docs/index.md`: Read Order、Directory Contract、Choose The Right Record を `requirements.md` / `system_requirements.md` へ更新。
- `docs/core-beliefs.md`: 関連文書リンクと記録方針を「業務要件 = `requirements.md`、システム要件 = `system_requirements.md`」へ更新。
- `docs/architecture.md` / `docs/tech-selection.md`: 関連文書リンクを `requirements.md` / `system_requirements.md` へ更新。
- `docs/screen-design/README.md`: `spec.md` 参照を `requirements.md` へ更新。

### 判断

- 業務要件（何をしたいか）とシステム要件（どう達成するか）を別文書に分ける。
- 単語の一貫性のシステム要件は TBD として明示的に保留（ユーザー判断）。
- AI プロバイダ指定「OpenAI API(llama, lmstudio)」は、OpenAI 互換 API でローカル実行（llama、LM Studio）と OpenAI 本体を含む意味と解釈した。

### 残課題

- `.claude/skills/coding-protocol/SKILL.md` の 2 箇所が `docs/spec.md` を参照（system 要件の参照行、docs 正本一覧）。auto mode で skill 編集が拒否されたため未修正。ユーザー承認後に `requirements.md` / `system_requirements.md` へ張り替える。
- `requirements.md` は用語集を廃止済みのため、`screen-design/README.md` の「用語」参照は形式的に古い。画面設計の用語運用を決める時に整理する。

## 2026-06-13 spec.md を業務要件専用へ書き換え

### 変更

- `docs/spec.md`: 恒久要件・用語集・状態機械を全削除し、4 つの業務要件（翻訳したい / 単語の一貫性 / NPC の口調 / xTranslator 形式出力）へ全面書き換え。各要件に目的を併記し、成功条件は不記載。

### 判断

- `spec.md` は業務要件（何をしたいか）だけにする。システム要件（どう実現するか）は別文書で扱う。
- 入力取得手段（xEdit 抽出など）はシステム要件側へ回す。業務要件側は対象を「Skyrim Mod のテキスト」とだけ書く。
- xTranslator 形式出力は、ツール固有だがユーザーの明示要望のため業務要件として残す。
- 成功条件は記載しない（ユーザー判断）。目的は記載する。

### 残課題

- システム要件の置き場が未定。入力取得手段、AI 基盤、ジョブ運用は置き場を決めてから書き起こす。
- `index.md`（`spec.md` を「恒久要件と用語集」と記述）と `core-beliefs.md`（「永続要件は `spec.md` に記録する」と記述）の文言が、業務要件専用化により古くなった。システム要件の文書を決める時に併せて直す。
