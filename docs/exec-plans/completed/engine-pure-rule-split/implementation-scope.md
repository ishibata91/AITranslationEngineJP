# 実装範囲: engine-pure-rule-split

承認済み設計（plan.md「設計確定」）を入力にした、implementation-module へ渡す scope 境界・依存・検証単位。

## scope の境界（触る対象）

### 新設する集積ディレクトリ `internal/core/` と 9 純粋 package（`internal/core/<name>`）

`internal/core/` は親ディレクトリで package を持たず、副作用のない決定的な計算ロジック（functional core）だけを束ねる。
各 package は `internal/engine/<name>.go` の純粋部を移し、package 宣言を `<name>` に変え、対応する `<name>_test.go` も同 package へ移す。

1. `internal/core/dictionary` ← `dictionary.go`＋`dictionary_test.go`
2. `internal/core/prompt` ← `prompt.go`＋`prompt_test.go`
3. `internal/core/termderive` ← `termderive.go`＋`termderive_test.go`
4. `internal/core/termusage` ← `termusage.go`＋`termusage_test.go`（`core/termderive` を import）
5. `internal/core/termxml` ← `termxml.go` の純粋部（`parseTermXML`→`ParseTermXML`、`isBaseGame`→`IsBaseGame`、`xmlString`、`baseGamePrefixes`）＋`termxml_test.go` の純粋部。さらに純粋結合 `DeriveTermsFromFiles(files []XMLFile, baseSources map[string]bool) []termderive.DerivedTerm` を新設（`XMLFile{Name string; Data []byte}` を受け取り、各 `ParseTermXML`＋`termusage.BuildUsage`＋`termderive.DeriveTerms` を束ねる）。`core/termusage`・`core/termderive` を import。
6. `internal/core/rolespeech` ← `role_speech.go` の純粋部（`ParseRoleSpeech`・`RoleSpeechTable`・`RoleSpeechTemplate`・`roleSpeechRow`・`roleWildcard`・`Lookup`・`matchScore`・`roleClassOfRace`→`RoleClassOfRace`）＋`role_speech_test.go`。os 読み `LoadRoleSpeech` は移さず廃止する。
7. `internal/core/linefeatures` ← `linefeatures.go`＋`linefeatures_test.go`（`core/tone` を import）
8. `internal/core/personatone` ← `tone_catalog.go`＋`tone_catalog_test.go`。engine が呼ぶ `buildToneDirective`・`buildToneTraits`・`buildToneLabel`・`personaMetaOf` を公開（`Build*`・`PersonaMetaOf`）。内部 helper（`toneTraitOf`・`raceMarkerTrait`・`roleSpeechLine`・`toneTraits`）は非公開のまま。`rolespeech.RoleClassOfRace`・`RoleSpeechTable.Lookup`・`tone.CellName`・`model.LinePersonaInput` を使う。
9. `internal/core/tone` ← `internal/engine/tone/`（`classifier.go`・`voice_traits.go`・`classifier_test.go`）をディレクトリごと移す。package 宣言 `tone` は変えず、import path を `internal/engine/tone`→`internal/core/tone` に直す。`internal/lexicon` は IO 持ちのため移さない。

### `engine` に残し、呼び出しを新 package 参照へ書き換える

- `engine.go`: `NewDictionary`→`dictionary.NewDictionary`、`ComposePrompt`→`prompt.ComposePrompt`、`DeriveTermsFromXMLDir` 内で `termxml.ParseTermXML`/`termusage.BuildUsage`/`termderive.DeriveTerms`、`buildToneDirective` ほか→`personatone.Build*`/`PersonaMetaOf`。
- `proper_noun.go`: `ComposePrompt`→`prompt.ComposePrompt`。
- `persona_generate.go`: `ExtractFeatures`→`linefeatures.ExtractFeatures`、`SourceHash`→`linefeatures.SourceHash`、`EmotionLexicon` 型→`linefeatures.EmotionLexicon`、`tone.NewClassifier` ほか tone 参照の import path を `core/tone` へ。
- `role_speech.go`: ファイルごと削除する。`Engine.New` は `*rolespeech.RoleSpeechTable` を注入で受け取り、engine 自身はファイルを読まない。`engine.LoadRoleSpeech` は廃止。`Engine.roleSpeech` フィールドと `LinePersonas` の参照型を `*rolespeech.RoleSpeechTable` にする。
- `termxml.go`: 純粋部を core へ移したあと、engine 側には `DeriveMasterTerms` 内の os 読み（`filepath.Glob`＋`os.ReadFile` ループで `[]termxml.XMLFile` を作る）だけ残し、`termxml.DeriveTermsFromFiles` を呼ぶ。`engine.DeriveTermsFromXMLDir` は廃止。
- `engine.go`: `Engine.lexicon` の型を `linefeatures.EmotionLexicon`、`Engine.roleSpeech` を `*rolespeech.RoleSpeechTable` に。`DeriveMasterTerms` を上記の os 読み＋core 呼びに書き換え。
- `dictionary` を返す `LoadDictionary`（store 読み）は engine 残置で、戻り値型を `*dictionary.Dictionary` にする。

### 利用元の import 書き換え

- `internal/api/app.go`: `engine.RenderPrompt`/`ComposePrompt`→`prompt.*`、`engine.DictionaryTerm`→`dictionary.DictionaryTerm`、`a.engine.LoadDictionary` の戻り型は `*dictionary.Dictionary`。
- `internal/harness/run.go`: `engine.ParseRoleSpeech`→`rolespeech.ParseRoleSpeech`。
- `internal/bootstrap/bootstrap.go`・`cmd/goldcap/main.go`: `engine.LoadRoleSpeech(path)`→自前で `os.Open(path)`（defer close）し `rolespeech.ParseRoleSpeech(f)` を呼ぶ。

### 機械検査の更新

- `.go-arch-lint.yml`: 新 component 9 個（`internal/core/<name>` を各 component に対応。既存 `tone` component の path を `internal/engine/tone`→`internal/core/tone` に直す）を登録し、依存方向を明示（下記「依存」）。
- 境界走査（`run-boundary-scan.sh`）: 純粋 package は os/driver/runtime を触らないため許可リスト変更なし。走査が緑であることだけ確認。

## 依存（arch-lint に固定する一方向グラフ）

```
api          → engine, model, provider, dictionary, prompt
engine       → model, provider, tone, dictionary, prompt, termderive,
               termxml, rolespeech, linefeatures, personatone, lexicon, store
bootstrap    → api, engine, provider, store, lexicon, rolespeech
goldcap      → api, engine, harness, lexicon, rolespeech
harness      → api, engine, provider, store, rolespeech
prompt       → provider
termusage    → termderive
termxml      → termderive, termusage
linefeatures → tone
personatone  → tone, model, rolespeech
dictionary   → （内部依存なし）
termderive   → （内部依存なし）
tone         → （内部依存なし）
rolespeech   → （内部依存なし）
```

逆依存・循環はない（core の純粋 package は engine を import しない。core 内は termusage/termxml→termderive、linefeatures/personatone→tone、personatone→rolespeech の一方向のみ）。

## 検証単位

- 単体（package 単位）: 移設した各 `_test.go` が新 package で緑。`go test -cover` で不変ルール関数 100%。
- 依存方向: `npm run lint:backend:arch`・`npm run lint:backend:boundary` 緑。
- 非劣化（統合）: `npm run test:backend`（harness golden = 送信プロンプト列＋DB 最終状態）が分割前と一致。
- 全体: `npm run verify:backend` exit 0。
- 実 app: `npm run dev:wails:run` で実画面から `.esp` 抽出→翻訳し、機械置換・固有名注入・口調指示が従来どおり効くことを目視。

## 含まない

- 出力を変えるロジック改修（非劣化が条件）。
- gocognit `//nolint` の解消（`backend-violation-cleanup` の担当。移設で nolint コメントはそのまま運ぶ）。
- os 読み wrapper の再設計。
- 表示（svelte・story・fixture）の変更。
