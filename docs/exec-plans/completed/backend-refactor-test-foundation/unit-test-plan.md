# 恒久ユニットテスト計画（不変ルールの純粋部括り出しと 100% 目標）

本書は完了定義⑤の成果物。backend 包括リファクタの前後で劣化させてはいけない不変ルールを列挙し、各ルールの「純粋部の括り出し方針」と「100% カバレッジ目標」を固定する。実装は本 task では行わず、本書を入力にリファクタ本体 task（純粋関数の別 package 化を含む）で進める。

## 用語と基準

- 不変ルール: 入力が同じなら出力が同じになる、翻訳品質や DB 整合に直結する決定的な規則。例は固有名の機械置換、プロンプト合成、取込の箱振り分け。
- 純粋単位: 副作用と IO を持たず、引数と戻り値だけで完結する関数。単体テストの対象はこの純粋単位に限る。
- IO 統合部: store・filesystem・provider に触れる手続き。単体テストの対象にせず、②の合成入力 golden harness（`internal/harness`）と api 統合テストで担保する。
- 100% 基準: `go test -cover` の手元確認で純粋単位の行・分岐を 100% にする。常設の coverage ゲート（CI 失敗条件）にはしない（Sonar 廃止に伴う方針）。
- 到達不能な防御分岐の扱い: 構造上の不変条件で到達しない防御分岐は、その不変条件をコメントで根拠づけたうえで除く。④の `Dictionary.Apply` で、`re` が `bySource` のキーだけから組まれる不変条件を根拠に `if !ok` を除いた先例に従う。

## ④で完了済み（純粋部の括り出しとテストが済んだルール）

| ルール | 純粋単位（関数） | 場所 | 現状 |
| --- | --- | --- | --- |
| 固有名の機械置換 | `NewDictionary`・`Apply` | `internal/engine/dictionary.go` | 100%（④で達成） |
| プロンプト合成 | `FillVariables`・`ComposePrompt`・`RenderPrompt` | `internal/engine/prompt.go` | 100%（④で達成） |

## 残る不変ルールの計画

各行の「現状」は本書作成時点の `go test -cover` 実測値。リファクタ後に同値以上を保つ。

### 1. 取込の箱振り分け（ingest 分類）

- 不変ルール: C# 抽出器が素朴吸い出しした `extracted_field` を、`record_type_master` の対応で叙述文・定型句・固有名・台詞・対象外へ振り分ける規則。
- 純粋単位: `Dispatch`・`recordMasterMap`（`internal/engine/ingest.go`）。現状いずれも 100%。
- IO 統合部: `engine.Ingest`（現状 0%）と `internal/store/ingest.go` の `INSERT OR IGNORE` 群。store に触れるため単体対象外。
- 括り出し方針: 分類規則は既に純粋関数へ括り出し済みで追加作業は不要。`Ingest` は分類済み結果を store へ投入する薄い手続きに保ち、振り分けの正否は harness golden の DB 最終状態で観測する。

### 2. 口調の性質生成と役割語（persona 性質/役割語）

- 不変ルール: 話者の本文特徴と声型から基底口調（性質文・段階・印）を決め、性別・年齢・基底口調で一人称と語尾の役割語を引く規則。
- 純粋単位と現状:
  - `internal/engine/tone_catalog.go`: `toneTraitOf`・`raceMarkerTrait`・`buildToneTraits`・`buildToneDirective`・`buildToneLabel`・`personaMetaOf` は 100%。`roleSpeechLine` は 66.7%。
  - `internal/engine/role_speech.go`: `ParseRoleSpeech`・`Lookup`・`matchScore`・`roleClassOfRace` は 100%。
  - `internal/engine/tone/`（`Classifier.Classify` ほか）: `classifier_test.go` で担保。
  - 本文特徴の抽出（`ExtractFeatures` ほか）: `linefeatures_test.go` で担保。
- IO 統合部: `PersonaGenerator.Generate`（86.4%）・`ensureLineAnalyses`（90.5%）。store に触れる集計・キャッシュ手続きのため単体対象外。
- 括り出し方針: 純粋部は概ね括り出し済み。残る穴は `roleSpeechLine`（66.7%）の未到達分岐のみ。役割語が引けない経路（該当テンプレートなし）の戻り値を分岐網羅して 100% にする。`Generate`・`ensureLineAnalyses` は本文ハッシュのキャッシュ判断と集計の手続きに保ち、口調出力の正否は harness golden の `persona_character`・`line_analysis` と送信プロンプトで観測する。

### 3. 固有名の部分形派生（termderive 派生）

- 不変ルール: 確定済みの人名（FULL）から名のみ・短名・二つ名前部を派生し、本文機械置換へ載せる単位を決める規則。地雷語（一般語化する短名）を除く判定を含む。
- 純粋単位: `DeriveTerms`・`safePair`・`landmine`・`hasByname`・`isKanaToken`・`isKatakana`・`hasKanji`・`trailingKanaRun`・`wordSet`・`DefaultDeriveConfig`（`internal/engine/termderive.go`）。現状すべて 100%。
- 括り出し方針: 既に純粋で 100%。追加作業は不要。リファクタで同値を保つことだけ確認する。

### 4. xTranslator XML の整形（termxml 整形）

- 不変ルール: xTranslator 英日 XML から `NPC_:FULL` の原日対を取り出し、base ゲーム plugin を接頭辞で判定して姓名分割（two）の可否を決める規則。
- 純粋単位と現状:
  - `parseTermXML`（XML バイト列 → 原日対。92.0%）。
  - `isBaseGame`（ファイル名接頭辞判定。0.0%。純粋だが未テスト）。
- IO 統合部: `DeriveTermsFromXMLDir`（0.0%。ディレクトリ走査と読み込み）。filesystem に触れるため単体対象外。
- 括り出し方針: `parseTermXML` は XML バイト列を引数に取る純粋関数のため、未到達分岐（不正タグ・空 Content・REC 不一致など）を分岐網羅して 100% にする。`isBaseGame` は plugin 名を引数に取る純粋関数のため、base ゲーム名（Skyrim・Dawnguard ほか）と mod 名の両方で 100% にする。`DeriveTermsFromXMLDir` は走査と読み込みに保ち、整形の正否は `parseTermXML` の単体と harness golden の `master_term` で観測する。

### 5. 結果取得の cursor と DTO 再構成（api の cursor/DTO 再構成）

- 不変ルール: 叙述文→台詞→固有名の連結列を cursor で区間送りする規則と、結果行へ機械置換内訳・元レコード種別バッジ・口調メタを再構成して載せる写像。
- 純粋単位と現状（`internal/api/app.go`）:
  - `parseCursor`・`makeCursor`（cursor の符号化・復号。いずれも 100%）。
  - `narrationResultView`（叙述文行 → DTO。100%）。
  - `directiveViews`・`assignmentViews`・`recordTypeView`・`termViews`（いずれも 0.0%。純粋な写像だが未テスト）。
- IO 統合部: `ListResultsPage`・`buildResultsPage`・`pageRows`（80.9%）。store に触れるため単体対象外。
- 括り出し方針: `directiveViews`・`assignmentViews`・`recordTypeView`・`termViews` は引数のスライス・マップだけから DTO を組む純粋関数のため、テーブル駆動で 100% にする。`recordTypeView` は箱の有無（nil 返し）の両分岐、`termViews` は空・重複・複数語の各経路を含める。`buildResultsPage`・`pageRows` は store 取得とページングの手続きに保ち、再構成の正否は ②の api 統合 golden（送信プロンプト user の置換済み原文、結果 DTO）で観測する。

## まとめ（リファクタ本体 task への入力）

- 追加作業が要る純粋単位: `roleSpeechLine`（2）、`parseTermXML`・`isBaseGame`（4）、`directiveViews`・`assignmentViews`・`recordTypeView`・`termViews`（5）。いずれも分岐網羅の単体テスト追加で 100% にする。
- 括り出し済みで同値維持だけ要る純粋単位: ingest 分類（1）、termderive 派生（3）、tone_catalog と role_speech の大半（2）。
- IO 統合部は単体対象にせず、②の合成入力 golden harness と api 統合テストで担保する。
- 「純粋関数の別 package 化」（完了定義⑦の対象外・リファクタ本体 task）は本書の純粋単位の一覧を入力にする。別 package へ移す際も `go test -cover` 手元確認の 100% を保つ。
