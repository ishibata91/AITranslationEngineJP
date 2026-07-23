# Design: strings-based-reference

この design は「どう実装しどう直すか」だけを持つ。scope 列挙・テスト設計は実装モジュールが扱う。

## 目的（不変）

参照訳（`reference_translation`）と固有名の確定訳語（`master_term`）の日本語 dest の供給を、xTranslator 英日 XML から Data フォルダの Strings（english / japanese）へ移し、XML 依存を廃止する。役割は不変で、供給源だけ替える。

## 論理的状態（定義）

AS-IS/TO-BE 表はこの ID で引く。同じ ID が両表で同じ論点を指す。I（照合キー）だけ現状維持で、他は変更点。

| ID | 論理的状態（何を決める論点か） |
| --- | --- |
| A | 原文が既訳と完全一致する叙述文・台詞へ流用する (英→日) 対を、どこから `reference_translation` に入れるか |
| B | 固有名の原語そのままの確定訳語 `(原語→訳語)` を、どこから `master_term` に入れるか（FULL 完全形） |
| C | フルネームから名のみ・短名を機械派生した確定訳語を、どこから `master_term` に入れるか（部分形） |
| D | 英語原文に対応する日本語既訳を、誰が突き合わせて (英→日) 対にするか |
| F | 抽出器が各 field を何言語ぶん解決するか（抽出の言語） |
| G | (英→日) 対をどの field から作るか |
| H | C# が書く受け皿テーブル `extracted_field` が持つ列とスコープ |
| I | 参照訳を原文へ突き合わせる照合キー |
| J | 実行時に版管理外の XML ディレクトリ `dictionaries/xTranslatorXMLs` を読む箇所 |
| K | 姓名分割（two）の派生を許してよいかの判定材料 |

## AS-IS

| ID | AS-IS 対応 | ソース（コードベースの根拠） |
| --- | --- | --- |
| A | xTranslator 英日 XML の全 rec:field（source・dest 非空）を取り込む | `engine/reference.go:28`・`core/termxml/reference.go:40-79`・`migration 0012` |
| B | C# が XML の `*:FULL`（`WOOP:FULL` 除外）を `(source,dest,category=rec)` で書く | `extractor/MasterTermXmlWriter.cs:15,87`・`Program.cs:82` |
| C | Go が XML の `NPC_:FULL`/`SHRT`・`INFO:NAM1` を termderive で派生 | `engine/engine.go:492-505`・`core/termxml/termxml.go:42-120` |
| D | 外部辞書 XML が対で持つ。抽出器は両言語を差分比較で触るが本文は捨てる | `xTranslatorXMLs`(gitignore)・`RecordDataIndex.cs:52-62`(先勝ち)・`PluginEnvironment.OwnsRecord` |
| F | 単一。本番は English 固定（`--language` 未指定） | `PluginEnvironment.cs:57`・`Program.cs:14`・`api/app.go` buildExtractorArgs |
| G | XML 群 `xTranslatorXMLs` が base 全体を対象と無関係に供給。取込は plugin で絞らない | `bootstrap.go:35`・`migration 0012`「plugin では絞らない」 |
| H | `(plugin, form_id, edid, rec, field, ordinal, source)`。dest 無し・英語 source のみ・対象 plugin 単位 | `migration 0006`・`ExtractedFieldSqliteWriter.cs:25-28` |
| I | `(rec, field, source)`。XML が FormID を持たないため | `engine/reference.go:18`・`migration 0012` |
| J | 3 経路が `dictionaries/xTranslatorXMLs` を読む | `api/app.go:487,494`・`Program.cs:78-82` |
| K | XML ファイル名の接頭（Skyrim/Dawnguard…） | `core/termxml/termxml.go:20-30` |

## TO-BE

| ID | TO-BE 対応 | ソース（変更先） |
| --- | --- | --- |
| A | `extracted_field` の dest 非空行から全 rec:field を組む | `engine/reference.go` `LoadReferenceTranslations(ctx)` 改・`store/ingest.go:16` `ListExtractedFields` 流用 |
| B | Go が `extracted_field` の 箱=固有名（`record_type_master`）行を `(source,dest,category=rec)` で書く（`DeriveMasterTerms` 手順1） | `engine/engine.go` `DeriveMasterTerms` 改・`record_type_master`(`migration 0006`) |
| C | Go が `extracted_field` の `NPC_` 対から termderive で派生（手順3、termderive 不変） | `engine/engine.go` 改・`core/termderive` 不変 |
| D | **C# 抽出器**が english+japanese Strings を両解決して対を作る（Mutagen 環境でのみ解決可） | `PluginExtractor.S` 改（両言語）・`RecordDataIndex.cs:206` `LoadStrings` 流用可 |
| F | english + japanese の 2 言語を解決 | `PluginEnvironment.Load` 改・`ExtractedFieldSqliteWriter` dest 書き込み |
| G | 抽出時に english と japanese の Strings が両方ある field だけを (英→日) 対にする（dest が入る）。base ゲームかは判定せず、base ゲームを別途抽出対象にもしない | `PluginExtractor` の 2 言語解決（F）・`ExtractedFieldSqliteWriter` は japanese がある時だけ dest を書く |
| H | dest 列（日本語本文）を追加。他は不変・対象 plugin 単位のまま | 新 `migration 0014`・`model/extracted_field.go`・`ExtractedFieldSqliteWriter` dest |
| I | `(rec, field, source)` を維持（form_id が使えるが対象横断再利用のため絞らない） | `engine/reference.go` 維持 |
| J | 廃止。3 経路・`termxml`・xmlDir・`xTranslatorXMLs` 配線を除去 | `MasterTermXmlWriter` 削除・`core/termxml` 削除・`api.New` termsXMLDir 除去 |
| K | その行に japanese dest がある（英日 Strings が揃う）かで判定。base ゲーム名の限定は廃止し、英日対のある行なら mod でも two を許す | `engine/engine.go` `DeriveMasterTerms` 改・termderive の base 判定を dest 有無由来へ・`baseGamePrefixes` 廃止 |

## 変える IF（どこを触るか）

C# 抽出器（英日ペアを作る本体。ここだけが Mutagen で英日を解決できる）。

- `PluginEnvironment.Load`（`PluginEnvironment.cs`）: english に加え japanese Strings も解決可能に読む（現状は `TargetLanguage=English` 単一）。
- `TranslationString`（`TranslationCounts.cs:6`、`(RecField, Id, EditorId, Text)`）: 日本語本文フィールドを足し、英語 `Text` と対で持つ。
- `PluginExtractor`（`PluginExtractor.cs`）: 各 field を english と japanese の両方で解決して `TranslationString` に入れる。
- `ExtractedFieldSqliteWriter.Write`（`ExtractedFieldSqliteWriter.cs:14`）: INSERT に dest 列を足し、日本語本文を書く。

DB schema。

- 新 migration: `extracted_field` に `dest TEXT NOT NULL DEFAULT ''` を ALTER で足す。

Go（XML 読みを DB→DB の組み替えへ替える）。

- `model.ExtractedField`・`extractedFieldColumns`（`store/ingest.go:12`）: dest 列を足す。読み口 `Store.ListExtractedFields`（`store/ingest.go:16`）はそのまま流用する。
- `Engine.LoadReferenceTranslations(ctx, xmlDir)` → `(ctx)`（`reference.go:28`）: `ListExtractedFields` の dest 非空行から `reference_translation` を組む。
- `Engine.DeriveMasterTerms(ctx, xmlDir)` → `(ctx)`（`engine.go:488`）: master_term への書き込みを 1 関数に畳む。順序は依存があるため固定する。
  1. `ListExtractedFields` のうち `record_type_master` で箱＝固有名の全 FULL（source・dest 非空）を `(source, dest, category=rec)` で `master_term` へ書く（旧 `MasterTermXmlWriter` の役目）。
  2. `ListMasterTerms` を baseSources として読む（手順 1 で書いた FULL を含む）。
  3. 同じ `ListExtractedFields` の `NPC_:FULL`/`SHRT`・`INFO:NAM1` から termderive の入力（`NamePair`・dialogues）を作り、baseSources と衝突する部分形を除外して派生し `master_term` へ追記する。
  - 手順 1 を派生より前に置くのは、派生が baseSources（既存 master_term の原語）との衝突を除外するため（`engine.go:493`）。FULL を先に入れないと除外が効かない。termderive・termusage は不変。two（姓名分割）の可否は行の japanese dest（英日対）の有無で判定し、base ゲーム名の判定は廃止する。
- `api.New(..., termsXMLDir, ...)`（`app.go:115`）: termsXMLDir 引数と `App.termsXMLDir` を除く。`prepareForTranslation`（`app.go:487,494`）は `DeriveMasterTerms(ctx)`・`LoadReferenceTranslations(ctx)` を呼ぶ。
- 配線から termsXMLDir を除く: `bootstrap.go:35,95`・`harness/run.go:28,45,87`・`goldcap main.go:115`。
- 削除: `internal/core/termxml`（`ReferencesFromFiles`・`ParseReferences`・`DeriveTermsFromFiles`・`ParseTermXML`・`XMLFile`・`IsBaseGame`）と engine の `readXMLDir`。

## 片側 Strings 欠け時の警告

english / japanese の片方しか無いと対を作れず、固有名は確定訳語を再利用できず全文 AI 翻訳になる。この状態を利用者へ知らせる画面警告を出す（固有名機能の観測可能な成果に含む）。判定材料と配線は plan の scope に含める。表示範囲は storybook-module で扱う。

## 決定が要る点

- **対の置き場**（推奨: `extracted_field.dest` を足して Go が組む）。英日解決だけを境界に置き、振り分け（`record_type_master`）と派生（termderive）は既存 Go ロジックを再利用する。対案は抽出器が `reference_translation`・`master_term` を直接書く形。
- **照合キー**（推奨: 現状の `(rec, field, source)` を維持）。form_id が使えるが、`reference_translation` は対象横断で同一原文を再利用する設計のため form_id では絞れない。
