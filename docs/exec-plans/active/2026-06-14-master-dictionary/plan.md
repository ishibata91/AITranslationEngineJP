# Task Plan: 2026-06-14-master-dictionary

- `workflow`: work
- `status`: in_progress（着手。preparation-module 完了）
- `task_id`: 2026-06-14-master-dictionary
- `task_mode`: 重 task（preparation-module で確定）
- `request_summary`: Mod 横断マスター辞書（原語 → 確定訳語）のテーブル・登録・適用を実装する。固有名詞を抽出して辞書へ登録し、辞書解決の差込点を engine に新設して本文翻訳へ適用する。T2 は差込点を作らなかったため、差込点も本 task で作る。
- `goal`: 同一原語へ常に同一訳語を当て、Mod 横断で一貫性を保つ機構を実装する。
- `constraints`: ルール・プロンプト編集 UI（T4）は対象外。
- `close_conditions`: 翻訳実行画面から実 mod を流したとき、固有名の訳が辞書として作られ、叙述文・台詞の翻訳で本文中に出る固有名が辞書の確定訳語へ機械置換され、同一原語へ常に同一訳語が当たることを観測点で確認できる。
- `source_branch`: `master`
- `source_commit`: `18b8e5ab`
- `target_branch`: `master`
- `work_branch`: `claude/2026-06-14-master-dictionary`

## 趣旨訂正（2026-06-16、人間指摘）

当初の design・storybook は趣旨を取り違えていた。訂正後の趣旨を本節で正本とし、以降の節で古い記述と食い違う場合は本節を優先する。

- 正しい趣旨: 固有名の訳をあらかじめ辞書化しておき、叙述文・台詞を翻訳するときに、本文中に出る固有名を辞書の確定訳語へ当てて一貫させる。
- 適用方式（確定）: 本文への機械置換を主軸にする（プロンプト注入でなく置換）。本文翻訳の前に、原語（英語）→ 確定訳語（日本語）へ貪欲最長一致で置換し、AI は周りの英語だけを訳す。固有名は日本語で活用しないカタカナ語が多く置換と相性が良い。一貫性は固定要求なので決定的な置換が合う。誤置換が出た語だけ注入へ逃がす余地は残す。
- 取り違えていた点: (1) 固有名を独立した一覧（固有名辞書パネル）として見せる方向は誤り。観測は叙述文・台詞の結果で固有名訳が揃うこと（必要なら結果行に「置換した固有名」を併記）。(2) engine の適用を「固有名レコード自身の `proper_noun.dest` 充填」とした design は趣旨とずれる。適用は本文（叙述文・台詞）への置換。(3) 本文中の固有名揃え（concept-model の言及 e3/e4/e5）を「対象外」とした判断は誤りで、これが本 task の核。
- 巻き戻し: storybook で足した独立パネル（`TermResultRow`/`TermResultsPanel`、画面組み込み）は master へ戻して削除済み。

## 完了定義（preparation-module、趣旨訂正を反映）

システム上どこまで動かすかを 1 つに固定する。

- 動かす範囲: 翻訳実行画面から実 mod を流したとき、抽出した固有名詞がマスター辞書へ登録され、同一原語へ常に同一訳語が当たり、その一貫性を観測できる状態まで動かす。具体的に次の 4 つが実際に効く。
  1. extractor が名前付きレコード（武器・防具・NPC・地名など表示名を持つレコード）の固有名詞を `proper_noun` として中心 DB へ書く。現状の extractor は narration・台詞・話者属性だけを書くため、固有名詞の書き込みを新たに足す。
     - 前提（strings 制約）: base master（`Skyrim.esm` ほか）は localized plugin であり、固有名詞の表示名は plugin 本体でなく言語別 strings（`dictionaries/Data/strings/`）に持つ。extractor は `TargetLanguage` を指定して strings を解決して初めて原語を読む（`PluginEnvironment.cs`、`PluginExtractor.cs:43`）。strings が無い／entry が無いと原語が空になり、辞書登録が成立しない。よって固有名詞抽出は strings 解決を必須前提とする。
     - 確定訳語の供給は公式日本語版の既訳流用を主軸にする（preparation-module で確定。`system_requirements.md` §2 の未確定を本 task で既訳流用主軸に固定）。よって extractor は英語 strings（原語）と日本語 strings（既訳）の 2 言語を解決し、固有名詞ごとに原語と既訳を対応付ける。現状の extractor は 1 回の実行で 1 言語だけ読むため、2 言語を読んで対応付ける経路を新たに足す。既訳が無い固有名詞は辞書に登録しない（本文翻訳側では従来どおり AI が訳し、辞書解決はヒットしない）。原語と既訳の対応付けキー（FormID / EDID のどちらで突き合わせるか）は design-module で確定する。
  2. マスター辞書テーブル（Mod 横断、原語 → 確定訳語）を持ち、抽出した固有名詞を辞書へ登録する。テーブルの設計先（`er.md` 更新か辞書専用設計か）は design-module で決める。
  3. engine が原語をキーに辞書を引く辞書解決の差込点を engine 内に新設し、解決した確定訳語を本文翻訳へ適用する。T2 は差込点を作らなかったため、差込点の新設も本 task で行う。適用方式（プロンプト注入か機械置換か）は design-module で確定する。
  4. 同一原語が複数レコード・複数 Mod に現れても同一訳語が当たり、その一貫性を観測点で確認できる。
  差込点（interface 宣言・空テーブル・引数追加）を置くだけで、同一原語へ同一訳語が当たることを観測できない状態は「動く」と書かない。
- 観測点: 実 app の実画面と実データを主とする。前提として localized strings 資源（`dictionaries/Data/strings/`）が揃っていること（無いと固有名詞の原語が空になり観測が成立しない）。`npm run dev:wails:run`（`http://localhost:34115`）で実 mod を実行し、(a) 固有名詞が `proper_noun` として抽出され（strings から原語が解決される）、(b) 同一原語へ常に同一訳語が当たることを画面で目視する。実画面での一貫性の見せ方（結果一覧への固有名表示など）の具体は design-module・storybook-module で決める。補助として、辞書解決（原語 → 確定訳語の lookup と適用）の単体テストと、extractor の固有名詞書き込みの単体テストを置く。
- goal 整合: goal は同一原語へ常に同一訳語を当て Mod 横断で一貫性を保つ機構を要求する。本完了定義は、抽出・登録・適用・一貫性のいずれも実データまたは実画面で観測できることを要求し、空テーブル・差込点・単体テストだけで goal を満たしたとはしない。
- close_conditions 整合: 本 task の close_conditions は、固有名詞の抽出・辞書登録・同一原語への同一訳語適用を実 mod 実行で観測する形にしてある。観測点を欠く close_conditions は残さない。
- 矛盾検査: goal が要する手段（固有名詞抽出、マスター辞書、辞書解決の差込点、確定訳語の適用）は、除外項目（ルール・プロンプト編集 UI＝T4、本文・会話文中にだけ現れる語）のいずれも必要としない。goal と「含まない」は矛盾しないため、本モジュールで停止しない。

## 軽 / 重判定（preparation-module）

- 画面が動くか: Y。一貫性を実画面で観測するため、結果一覧へ固有名（`proper_noun`）の訳の単位を表示する変更が要る見込み。表示構造・fixture・story を変える。表示の具体と確定は design-module（適用方式の決定）と storybook-module に従う。
- `docs/architecture.md` 反映が要るか: Y。§8「現在の状態と移行」は extractor が narration・台詞を書き engine が両者を訳す前提で書かれているが、本 task で extractor の固有名詞書き込みと engine の辞書解決を足すため現状記述が変わる。`er.md` は master 辞書・proper_noun を「本書では設計しない」とするため、`proper_noun` テーブルと辞書テーブルの設計先（`er.md` 更新か辞書専用設計か）も design-module で決める。最終判断は finalization-module。
- 結論: 両方 Y → 重 task。経路は preparation-module → design-module → storybook-module → implementation-module → finalization-module。

## Scope（含む / 含まない）

含む:
- マスター辞書テーブルの設計と実装。`er.md` のスコープ外だったため、ここで ER に追加する（`er.md` 更新または辞書専用設計）。
- 揃える対象用語の特定（名前付きレコードの固有名詞の機械抽出。`system_requirements.md` §2）。`er.md` の `proper_noun`（訳の単位・重複排除 e1）と接続する。
- 確定訳語の既訳流用主軸での供給（preparation-module で確定）。英語 strings（原語）と日本語 strings（公式既訳）の 2 言語を解決し、既訳のある固有名詞を辞書へ登録する。
- 辞書の登録・適用（プロンプト注入または機械置換。適用方式は design-module で確定）。

含まない:
- ルール・プロンプトの編集 UI（T4）。
- 本文・会話文中にだけ現れる語の拾い上げ（`system_requirements.md` §2 で未確定）。
- 既訳の無い固有名詞への確定訳語の AI 生成（本 task は既訳流用主軸とし、AI による訳語生成は採らない。本文翻訳側の AI 訳は従来どおり動くが、辞書登録の供給源にはしない）。

## 依存

- T2（`2026-06-14-persona-dictionary-pipeline`）の narration・台詞の抽出 → 翻訳 → 結果一覧表示の縦切りが master に存在すること。翻訳実行画面（`frontend/src/ui/screens/translation-run/`）と engine の `Run`（`internal/engine/engine.go`）を起点にする。
- `er.md` の `proper_noun`（固有名の訳の単位）。ただし `er.md` は line 17 で master 辞書・proper_noun を「本書では設計しない」と除外するため、proper_noun テーブルと辞書テーブルの設計先は design-module で決める。
- localized strings 資源（`dictionaries/Data/strings/`）。固有名詞の原語は base master の言語別 strings にしか無いため、英語 strings（原語）が揃っていることを前提にする。既訳流用を採るなら日本語 strings（`*_japanese.strings`/`.ilstrings`/`.dlstrings`）も前提になる。現状この資源は揃っている。

### 依存の訂正（preparation-module）

- 当初 plan は依存に「T2 の辞書解決の差込点」と書いていた。T2 は差込点を作らず完了した（`internal/engine/engine.go` に辞書解決は無い。T2 plan 末尾「後続 task への影響」も同旨）。
- よって本 task は engine の辞書解決の差込点を持たない前提で着手し、差込点の新設も本 task の動かす範囲に含める。

## 確定（preparation-module）

- 訳語の供給方式: 公式日本語版 strings の既訳流用を主軸にする。AI による訳語生成は本 task では採らない。既訳が無い固有名詞は辞書に登録しない。

## 確定（design-module、趣旨訂正後）

- マスター辞書 `master_term`（原語 → 確定訳語）を新規正本として中心 SQLite（`db/aitranslation.dev.sqlite3`）に追加する（案A 同居、人間承認済み）。起動時 wipe は未実装で本 task 対象外。
- 辞書の構築（実装で確定）: extractor が xTranslator 英日辞書 XML（`dictionaries/Data/` 兄弟の `dictionaries/xTranslatorXMLs/`、公式日本語版既訳）を解析し、固有名（`REC:FULL`、龍語綴り `WOOP:FULL` を除く）の `<Source>`（原語）→ `<Dest>`（確定訳語）を `master_term` へ登録する。当初設計の「英語/日本語 strings を 2 言語 Mutagen 読みで突き合わせる」案は、同じ既訳をより単純・高速に得られ vanilla 固有名（Skyrim.esm 等）も covers する XML 解析へ置き換えた。種別（REC 接頭）も保持するが、本文置換の照合キーは原語（文字列）。実 6 XML から 24,554 件を構築できることを確認済み。
- 適用方式（確定・趣旨訂正後）: engine の辞書解決の差込点で、叙述文・台詞の原文（英語）に対し辞書の原語を貪欲最長一致で照合し、確定訳語（日本語）へ機械置換してから AI 翻訳する。照合は最長一致優先・語境界・大文字始まりを手掛かりにし、誤置換を減らす。同綴り異義の誤置換は概念モデル弱点1と同じ性質で許容する。
- 表示（確定・趣旨訂正後）: 独立した固有名辞書パネルは持たない。叙述文・台詞の結果行に、その本文で置換した固有名（原語 → 確定訳語）を併記する（口調指示の併記と同じ要領）。
- proper_noun／placement の訳の単位としての出力は本 task の対象外にする（本文への置換に絞る）。
- 設計成果物: `design-diff.md`（設計差分図、趣旨訂正後に作り直す）、`implementation-scope.md`（実装範囲・テスト設計、趣旨訂正後に作り直す）。

## Routing Notes

- `required_reading`:
  - `docs/system_requirements.md`（§2 一貫性＝Mod 横断マスター辞書）
  - `docs/concept-model.md`（固有名・重複排除 e1）
  - `docs/er.md`（`proper_noun`/`set_phrase`/`placement`）
  - `docs/architecture.md`（store / engine の責務）

## Outcome

- preparation-module 完了。branch `claude/2026-06-14-master-dictionary`（分岐元 `master` `18b8e5ab`）を作成し、`完了定義` と `軽 / 重判定`（重 task）を固定した。
- design-module（初回）完了後、storybook-module の人間レビューで趣旨の取り違えが判明（独立した固有名辞書パネルは趣旨と違う）。適用は本文への機械置換、観測は叙述文・台詞での固有名訳の一貫性、と訂正。
- storybook で足した独立パネルは master へ戻して削除。趣旨訂正に合わせ design-diff.md・implementation-scope.md を作り直し、機械置換主軸（貪欲最長一致）で実装した。
- 実装済み（単体・コンポーネント検証通過）:
  - `internal/engine/dictionary.go`＋test: 貪欲最長一致の固有名置換（純関数）。
  - `internal/engine/engine.go`: `Run` で辞書を組み叙述文・台詞の原文を置換してから AI 翻訳する差込点。`DictStore` 追加。
  - `db/migrations/0003_master_term.sql`、`internal/model/master_term.go`、`internal/store/master_term.go`: 永続辞書テーブルと読み出し。
  - `tools/extractor/MasterTermXmlWriter.cs`＋`Program.cs`: xTranslator 英日 XML から `master_term` を構築（24,554 件確認）。
  - frontend: `TranslationResultRow.svelte` に結果行への「置換した固有名」併記（人間レビュー承認、説明文は削除）。
- 残り: 実 app での end-to-end 観測（実 mod を流し、叙述文・台詞で固有名が一貫置換されるのを目視）。plugin 選択はネイティブダイアログのため人間に依頼、翻訳は LAN の LM Studio。

## 増分: 人名の部分形の派生（2026-06-18、implementation-module）

機械置換の base 辞書はフルネーム（`NPC_:FULL`）しか持たず、台詞・叙述文で人名が「名のみ」「短名」で出ると辞書に当たらず AI 任せで訳が揺れた（例: 辞書に `Grelod the Kind` はあるが台詞の `Grelod` 単独は当たらない）。本増分は部分形の確定訳語を機械派生し、安全フィルタで一般語衝突を防いで `master_term` へビルド時に焼く。設計は `design-person-name-derivation.md` を正本とする。

- 実装（純粋ルール、置換器と同じ engine 配下）:
  - `internal/engine/termderive.go`: 人名の部分形を 3 種（shrt / byname / two）で派生し、安全フィルタ（用法比・称号・種族・縮約・純カタカナ・最小長・base 衝突・two の base ゲーム限定）で採否する純関数 `DeriveTerms` と各ヘルパ。副作用なし。
  - `internal/engine/termusage.go`: 会話文から用法分布（lc/uc）を作る純関数 `BuildUsage`。
- 実装（I/O 配線、ビルド時コマンド）:
  - `scripts/dict/derive-master-terms/main.go`: 同じ XML 群を解析し、用法分布を作り、純粋ルールを呼び、base 衝突を除いて `master_term` へ `category="derive:<種別>"` で追記する。実行時の置換器 `dictionary.go`・`engine.go` は無改造（`loadDictionary` が category を問わず全件を読むため派生行を自動で取り込む）。
- テスト:
  - `internal/engine/termderive_test.go`・`termusage_test.go`: 純粋ルールの全分岐を突く単体テスト。`go tool cover -func` で新規ルール関数のカバレッジ 100% を確認（task の承認済み基準）。
  - `scripts/dict/derive-master-terms/main_test.go`: XML 解析の単体テストと、temp DB へ XML→派生→`master_term` 追記・base 衝突 skip を確かめる結合テスト。
- 最終検証（backend、go 直実行。harness に backend suite が無いため）:
  - `go test ./internal/engine/ ./scripts/dict/derive-master-terms/` 通過。新規ルール関数カバレッジ 100%。`go vet`・`gofmt` 通過。
  - 実 XML 6 本で派生コマンドを実行し、検証済みロジックの再現を確認: 派生 582 件（byname 19 / shrt 275 / two 288、base 空時）。採用 `Grelod→グレロッド`(byname)・`Mercer→メルセル`/`Frey→フレイ`(two)・`Maro→マロ`(shrt)。棄却 `Master`・`Blood`・`Mine`・`Snow`・`Aren`・`Imperial`・`Nord`・`Lord`・`Captain`・`Guard`。
  - 注: `go build ./...` は既存の壊れ（`scripts/test/seed-system-test-db/main.go` が存在しない `internal/repository` を import、commit f2fbb14a 由来）で失敗する。本増分の範囲外で未修正。
- 対象外（明示）: two の category 形容語フィルタ。検証済み数値はフィルタ無しで得たもので、追加は未検証のため入れない。起動条件は「ハーネスで過剰置換が減り被覆が落ちないと実測したとき」。
- 運用（訂正）: plugin 選択はネイティブダイアログに加えてパス直接入力欄がある（`FileSelectField.svelte` の `onPathInput`、`TranslationRunScreen.svelte:36,95` で `onPluginPathInput` を配線）。よって end-to-end の実画面確認に人間補助は要らない。先の「ネイティブダイアログのため人間に依頼」は誤り。
- runtime 経路の確認: app の `RunExtractAndTranslate`（`internal/api/app.go:148`）は「C# extractor 子プロセス → 即 `engine.Run`」を 1 呼び出しで行い、間に派生段は無い。中心 DB は wipe されず（`project-db-wipe-on-launch-intent` の wipe は未実装）、全 writer が `INSERT OR IGNORE` のため再抽出は base を入れ直すだけで派生行を消さない。よって runtime 無改造のまま、事前に派生コマンドを 1 度流せば以後の翻訳は base+派生で動く。
- 実施済みの観測（AI 無し、実 DB）:
  - dev DB `db/aitranslation.dev.sqlite3`（base 24,554 件）へ派生コマンドを実行し、派生 517 件を追記（byname 17 / shrt 226 / two 274。base 衝突で空 base 時 582 から 65 件 skip）。
  - 実 DB の base+派生辞書で機械置換を確認。`Grelod runs the orphanage in Riften.` → `グレロッド runs the orphanage in リフテン.`（名のみ Grelod の原問題が解消）、`Mercer Frey` → `メルセル・フレイ`（base 最長一致）、`Mercer` 単独 → `メルセル`（派生 two）、`The master told me to wait.` は無置換（一般語 master の過剰置換なし）。
- end-to-end の実画面観測（完了、2026-06-19）: `npm run dev:wails:run` で起動し、plugin パスを手入力（`dictionaries/Data/Innocence Lost - Quest Expansion.esp`）、endpoint `http://192.168.0.226:1234`、実 LLM（LM Studio）で翻訳を実行した。検証材料は同 mod の台詞 121 件で、未訳台詞中に裸 `Grelod` が出る。観測前は前回（派生辞書なし）の訳が残り、裸 `Grelod` が `グレロド`（っ無し）と `グレロッド` に揺れていた。全 121 件を未訳へ戻し、派生辞書がある状態で翻訳し直した。
  - gemma-4-12b-it 実行: `Grelod` 28 行が 28 行とも `グレロッド`、`Riften` 4 行が 4 行とも `リフテン`。崩れ 0。close 条件「同一原語へ常に同一訳語」を満たす。
  - 所見（モデル依存）: 機械置換は AI 前段で確定訳語を注入する（`engine.Run` の `dict.Apply`）。一貫性は AI が注入カタカナを保持するかに依存する。shisa-v2.1-qwen3-8b 実行では `Grelod` 24/28・`Riften` 0/4 しか保持せず、注入済み `リフテン` を `リヴェン` 等へ書き換えた。同一の置換後 source を両モデルへ直接投げて確認した（gemma=保持、shisa=崩れ）。これは派生のバグではなく（base 名 `Riften` も同様に崩れる）、弱い小型モデルが注入トークンを保持しない問題。翻訳モデルは注入トークンを保持できる能力のものを選ぶ。
  - 後段修復（AI 翻訳後に注入訳語を再照合・修復）は不採用（人間判断、2026-06-19）。強いモデルまたはクラウドモデルで解く問題のため、本 task でも将来でも作らない。
  - 「実行後にどれが置き換えられたか」を見る表示は T4 へ回す（人間指示、2026-06-19）。現状 frontend の器（`TranslationResultRow` の `terms`）はあるが engine が `dict.Apply` の used を捨て backend 供給が無い。供給経路の実装は T4 plan（`2026-06-14-prompt-persona-customization`）の追加要望へ記録した。口調が不十分な点も同 T4 へ記録した。

## finalization-module（2026-06-19）

### 正本化判断

- `docs/architecture.md` 反映: 不要と判断（人間確認のうえ）。§1〜7（層構成・依存方向・強い制約・Wails 境界）は不変。`DictStore` は consumer 側の狭い interface（§4 で許容）、派生コマンドはビルドツールで実行時の層ではない。§8（現在の状態）は記述スナップショットで、構造不変のため churn しない（結果一覧ページング task と同じ判断）。当初 plan の「§8 反映 Y」は事前見込みで、実装が層を変えなかったため不要に確定。
- 人間承認済みの恒久仕様の正本反映: なし。
