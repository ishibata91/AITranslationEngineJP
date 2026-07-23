# Design: strings-based-reference

この design は「どう実装しどう直すか」だけを持つ。scope 列挙・テスト設計は実装モジュールが扱う。

## 目的（不変）

参照訳（`reference_translation`）と固有名権威訳（`master_term`）の日本語 dest の供給を、xTranslator 英日 XML から Data フォルダの Strings（english / japanese）へ移し、XML 依存を廃止する。役割は不変で、供給源だけ替える。

## AS-IS

### 日本語 dest は今すべて XML 由来（3 経路・同一ディレクトリ）

`reference_translation` と `master_term` の日本語 dest は、gitignore 対象の xTranslator 英日 XML から得ている。消費は 3 経路あり、すべて同じ XML ディレクトリ（`dictionaries/xTranslatorXMLs`）を読む。

| 経路 | 場所 | 生成物 | 採る対象 |
| --- | --- | --- | --- |
| C# `MasterTermXmlWriter` | `Program.cs:82` | `master_term`（FULL 完全形） | `*:FULL`（`WOOP:FULL` 除外）、英 source → 日 dest |
| Go `LoadReferenceTranslations` | `reference.go:28` | `reference_translation` | 全 rec:field で source・dest 非空の `<String>` |
| Go `DeriveMasterTerms` | `engine.go:488` | `master_term`（部分形） | `NPC_:FULL`/`SHRT`・`INFO:NAM1` → termderive で派生 |

### 入力読み込みがエンジン中核へ食い込んでいる

XML の形が、翻訳エンジンの API と照合ロジックまで規定している。ただの入力読み込みが境界に留まらず中核へ達している。

- engine の public メソッドが XML ディレクトリを引数に取る（`DeriveMasterTerms(ctx, xmlDir)`・`LoadReferenceTranslations(ctx, xmlDir)`）。翻訳エンジンが入力の置き場を知っている。
- engine 中核が形式パーサ `internal/core/termxml` を import している。
- 照合キーが `(rec, field, source)` なのは「XML が FormID を持たない」ため（migration 0012・`reference.go:18`）。しかし `extracted_field` は form_id を持つ。XML の制約に合わせて、使えるはずの強いキーを捨てている。
- base ゲーム判定が XML ファイル名の接頭（`Skyrim`/`Dawnguard`…、`termxml.go:20`）。ファイルの置き方が、姓名分割を派生してよいかという意味判断を決めている。
- `record_type_master` が全 REC:FIELD の箱（叙述文/固有名/定型句/台詞）を既に持つのに、XML 経路が `NPC_:FULL`/`INFO:NAM1` を再パースして同じ振り分けを二重に持っている。

### 抽出器は対を作る材料をほぼ持つが日本語を捨てている

- `extracted_field.source` は english Strings を Mutagen で解決した本文。`plugin, form_id, edid, rec, field, ordinal` 付き（対象 plugin 分）。
- japanese Strings は `RecordDataIndex` が差分正規化（`OwnsRecord` の master 一致判定）でだけ触り、本文はどの列にも書かず捨てる。本番起動は `--language` 未指定のため japanese を読み込みすらしない。
- `extracted_field` に dest 列が無い。**日本語本文を DB へ通す口が無い**。これが移行の最大ギャップ。

## TO-BE

### 軸: 入力読み込みを抽出器（入力の境界）へ集約する

Strings の英日解決は Mutagen を持つ C# 抽出器でしか作れない。よって対の生成を抽出境界に集約し、engine は DB だけを読む。engine から XML 形式・termxml・xmlDir を除く。

- **抽出器が japanese Strings も解決し、`extracted_field` に dest 列（日本語本文）を書く**。source（英語）は現状のまま。新テーブルは作らず、列追加 1 つ（ALTER）。
- **`reference_translation` と `master_term` を `extracted_field` の (source, dest) 対から組む**。REC:FIELD の振り分けは `record_type_master`（箱）を引く。XML 再パースによる二重振り分けを廃止する。
- **人名部分形の派生（termderive）は純粋・無 IO なのでそのまま**。XML を読む IO は呼び出し側（engine `DeriveMasterTerms` が `readXMLDir` で読み `store` へ書く／`termxml` が `[]byte` を解析）にある。この IO 側を、XML でなく `extracted_field` の `NPC_` 対を読む形へ替える。base ゲーム判定はファイル名接頭（`termxml`）でなく抽出時の master 連鎖から得た印を使う。
- **廃止**: C# `MasterTermXmlWriter`、Go の XML 読み（`internal/core/termxml`・xmlDir 引数）、`dictionaries/xTranslatorXMLs` 配線。

### 変える IF（どこを触るか）

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
  - 手順 1 を派生より前に置くのは、派生が baseSources（既存 master_term の原語）との衝突を除外するため（`engine.go:493`）。FULL を先に入れないと除外が効かない。termderive・termusage は不変。base ゲーム判定は行の plugin 名で行う。
- `api.New(..., termsXMLDir, ...)`（`app.go:115`）: termsXMLDir 引数と `App.termsXMLDir` を除く。`prepareForTranslation`（`app.go:487,494`）は `DeriveMasterTerms(ctx)`・`LoadReferenceTranslations(ctx)` を呼ぶ。
- 配線から termsXMLDir を除く: `bootstrap.go:35,95`・`harness/run.go:28,45,87`・`goldcap main.go:115`。
- 削除: `internal/core/termxml`（`ReferencesFromFiles`・`ParseReferences`・`DeriveTermsFromFiles`・`ParseTermXML`・`XMLFile`・`IsBaseGame`）と engine の `readXMLDir`。

### 供給範囲の帰結（XML との差）

XML は base ゲーム既訳を対象横断で一括供給していた。Strings 化後、base ゲームの日本語 dest は「base master を日本語 Strings 付きで抽出した時」に入る。`reference_translation`・`master_term` は対象横断・永続テーブルなので、抽出を重ねて蓄積する。bundled base 辞書は持たない。

### 片側 Strings 欠け時の警告

english / japanese の片方しか無いと対を作れず、固有名は権威訳を再利用できず全文 AI 翻訳になる。この状態を利用者へ知らせる画面警告を出す（固有名機能の観測可能な成果に含む）。判定材料と配線は plan の scope に含める。

## 決定が要る点

- **対の置き場**（推奨: `extracted_field.dest` を足して Go が組む）。英日解決だけを境界に置き、振り分け（`record_type_master`）と派生（termderive）は既存 Go ロジックを再利用する。対案は抽出器が `reference_translation`・`master_term` を直接書く形。
- **照合キー**（推奨: 現状の `(rec, field, source)` を維持）。form_id が使えるようになるが、`reference_translation` は対象横断で同一原文を再利用する設計のため form_id では絞れない。
