# engine-pure-rule-split

## 依頼要約

`internal/engine` に orchestration（副作用持ち）と同居している純粋不変ルールを、各々の純粋 package へ切り出す。
`engine` には Run・ingest・persona/proper_noun の orchestration だけ残し、純粋ルールは `engine` が一方向に import する独立 package にする。
先送り中だった「純粋関数の別 package 化」（backend-refactor-test-foundation の除外範囲に挙げたリファクタ本体）に当たる。

## 分岐元

- 統合先 branch: `master`
- 分岐元 commit: `ecf104ff`
- 作業 branch: `claude/engine-pure-rule-split`

## 対象の純粋ルール（現状 `internal/engine` 同居）

| 現ファイル | 切り出す純粋関数 | 役割 |
| --- | --- | --- |
| `dictionary.go` | `Apply` | 機械置換ルール |
| `prompt.go` | `ComposePrompt`・`FillVariables`・`RenderPrompt` | プロンプト組み立て |
| `termderive.go` | `DeriveTerms` | 固有名派生 |
| `termusage.go` | `BuildUsage` | 用例抽出 |
| `termxml.go` | `parseTermXML` | 固有名 XML 解析（純粋部） |
| `role_speech.go` | `ParseRoleSpeech`・`Lookup`・`matchScore`・`roleClassOfRace` | 役割語テーブル（解析・参照部） |
| `linefeatures.go` | `ExtractFeatures`・`SourceHash` ほか計算部 | 行特徴抽出 |
| `tone_catalog.go` | catalog 引き | 口調カタログ参照 |

`internal/engine/tone` も純粋（stdlib のみ）なので集積対象に加え、`internal/core/tone` へ移す。
`internal/lexicon` は `os.Open` でファイルを読む IO 持ちのため対象外（adapter として現状維持）。

`engine` に残す副作用持ち: `engine.Run`、`ingest`、`persona_generate`、`proper_noun`、および
`role_speech.LoadRoleSpeech`・`termxml.DeriveTermsFromXMLDir`（os 読みの薄い wrapper）。

## 完了定義

### 動かす範囲（task 後に観測できる振る舞い）

1. 純粋不変ルールが各純粋 package へ移り、`engine` は orchestration だけを持ち、純粋 package を一方向に import する。
   逆方向（純粋 package → engine）依存はない。
2. 翻訳 Run・固有名派生・役割語適用・プロンプト組み立ての出力が、分割前と一致する（非劣化）。
   package を移しただけで、出力を変えるロジック改修はしない。
3. 各純粋 package の不変ルールがその package のユニットだけで検証でき、不変ルール関数のカバレッジが 100% に保たれる。

### 観測点

- 依存方向: `npm run lint:backend:arch`（arch-lint）で新 component の一方向依存が緑。`npm run lint:backend:boundary` も緑。
- 非劣化: `npm run test:backend` の harness golden 比較が分割前と一致（送信プロンプト列＋DB 最終状態）。
- カバレッジ: 各純粋 package で `go test -cover` を手元実行し、不変ルール関数が 100%。
- 実 app: `npm run dev:wails:run` で実画面から `.esp` を抽出し翻訳を回し、機械置換と固有名注入が従来どおり効くことを目視。

### 含まない（除外範囲）

- 出力を変えるロジック改修（振る舞いは不変に保つ。非劣化が条件）。
- gocognit `//nolint` の解消・認知複雑度リファクタ（`backend-violation-cleanup` の担当）。
- os 読み wrapper（`LoadRoleSpeech`・`DeriveTermsFromXMLDir`）の再設計。純粋部を package へ出し、薄い読み込みは `engine` 側に残すに留める。
- 表示（svelte・story・fixture）の変更。

goal（純粋クラスの別 package 化）と 含まない は矛盾しない。
goal が要る手段（純粋関数の移動・呼び出し元の更新・arch-lint 反映）は 含まない で除外していない。

### close_conditions

- arch-lint・boundary 走査が緑で、依存方向が engine→純粋 package の一方向であることを検査できる。
- harness golden が分割前と一致し、非劣化を検証できる。
- 各純粋 package のユニットで不変ルール 100% を `go test -cover` で確認できる。
- 実画面で翻訳 Run が従来どおり動くことを目視で確認できる。

## 想定 Y/N

- 仕様変更または仕様追加: N（振る舞いは不変。package 構成だけ変える）。
- `docs/architecture.md` 反映: Y（§4 の component と依存方向を変える。finalization で反映）。

## 設計確定（人間設計レビュー承認済み）

配置・os 読み・口調組立の 3 論点を人間設計レビューで確定した。

- 配置: 全 8 純粋 package を集積ディレクトリ `internal/core/<name>` 配下に置く（engine の外）。`internal/core/` は親ディレクトリで package を持たず、副作用のない決定的な計算ロジック（functional core）だけを束ねる。判定基準は「純粋な計算は core、IO は engine」。
- os 読み: 純粋関数は object（`io.Reader`・`[]byte`・スライス）だけ受け取り、ファイル open は呼び出し側へ寄せる。これで純粋部は全て core へ移り、engine に純粋ロジックを残さない。
    - role-speech: `ParseRoleSpeech`（既に `io.Reader` 受け取り）を core へ移す。`Engine.New` は `*rolespeech.RoleSpeechTable` を注入で受け取るため、engine 自身はファイルを読まない。`engine.LoadRoleSpeech`（os.Open 包み）は廃し、composition root（`bootstrap`）と `cmd/goldcap` がファイルを開いて `rolespeech.ParseRoleSpeech` を呼ぶ。
    - 固有名 XML: 純粋結合 `DeriveTermsFromFiles([]XMLFile, baseSources)`（parse×N＋`BuildUsage`＋`DeriveTerms`）を `core/termxml` へ置く。`engine.DeriveMasterTerms` は実行時の os 読み（glob＋ReadFile）で `[]XMLFile` を作り、この core 関数を呼ぶだけにする（runtime の os 読みは engine 残置）。
- 口調指示組立: 現 `tone_catalog.go` を独立 package `personatone` として公開し切り出す。

確定した 8 package と依存（一方向、逆依存・循環なし）:

| 新 package | 公開シンボル（移す） | 内部依存 |
| --- | --- | --- |
| `internal/core/dictionary` | `Apply`・`NewDictionary`・`DictionaryTerm` | 標準ライブラリのみ |
| `internal/core/prompt` | `ComposePrompt`・`FillVariables`・`RenderPrompt` | `provider` |
| `internal/core/termderive` | `DeriveTerms`・`NamePair`・`Usage`・`DerivedTerm`・`DeriveConfig`・`DefaultDeriveConfig` | 標準ライブラリのみ |
| `internal/core/termusage` | `BuildUsage` | `core/termderive` |
| `internal/core/termxml` | `ParseTermXML`・`IsBaseGame`（現 `parseTermXML`/`isBaseGame` を公開） | `core/termderive` |
| `internal/core/rolespeech` | `ParseRoleSpeech`・`RoleSpeechTable`・`RoleSpeechTemplate`・`Lookup`・`RoleClassOfRace`（現 `roleClassOfRace` を公開）・`matchScore`(非公開) | 標準ライブラリのみ |
| `internal/core/linefeatures` | `ExtractFeatures`・`SourceHash`・`EmotionLexicon` | `core/tone` |
| `internal/core/personatone` | `BuildToneDirective`・`BuildToneTraits`・`BuildToneLabel`・`PersonaMetaOf`（現 `build*`/`personaMetaOf` を公開） | `core/tone`・`model`・`core/rolespeech` |
| `internal/core/tone` | `Classifier`・`NewClassifier`・`CellName`・`Features`・`Persona`（現 `internal/engine/tone` を移設、公開名そのまま） | 標準ライブラリのみ |

`engine` に残す orchestration: `Run`・`Ingest`・`LinePersonas`・`persona_generate`・`proper_noun`・`LoadDictionary`（store 読み→`dictionary.NewDictionary`）・`DeriveMasterTerms`（store 読み＋XML dir の glob・ReadFile→`core/termxml.DeriveTermsFromFiles`）。
`engine.LoadRoleSpeech` と `engine.DeriveTermsFromXMLDir` は廃止する（前者は composition root へ、後者は純粋結合を core へ移し engine 側は os 読みだけ残す）。

利用元の import 変化:
- `api`: `engine.ComposePrompt`/`RenderPrompt` → `prompt.*`、`engine.DictionaryTerm` → `dictionary.DictionaryTerm`。`dictionary`・`prompt` を import。
- `harness`: `engine.ParseRoleSpeech` → `rolespeech.ParseRoleSpeech`。
- `bootstrap`・`cmd/goldcap`: `engine.LoadRoleSpeech(path)` → 自前で `os.Open` し `rolespeech.ParseRoleSpeech(f)` を呼ぶ（asset 読みは composition root と cmd の責務）。

`tone` は `internal/engine/tone` から `internal/core/tone` へ移し、利用元（engine の `persona_generate`・`core/linefeatures`・`core/personatone`）の import path を直す。
移設後は core 内で完結し、core package が engine 配下を跨いで import する状態を残さない。
`internal/lexicon` は IO 持ちのため移さない。

人間承認状態: 承認済み（3 論点とも推奨案または指定案で確定）。

## 最終検証（機械）

- ビルド: `go build ./...` 緑、`go vet ./...` 緑。
- 単体: `go test ./ ./internal/...` 全 package 緑。core 9 package すべて緑。
- カバレッジ（`go test -cover ./internal/core/...`）: dictionary・prompt・rolespeech・termderive・termusage・tone は 100%。termxml は 95.0%（`IsBaseGame`・`DeriveTermsFromFiles` のユニットを追加。残 5% は `DecodeElement` のエラー経路で誘発困難）。linefeatures 87.2%・personatone 88.9% は移設前からの既存ギャップ（prose 経路・`roleSpeechLine` 分岐）で、コード不変のため維持。
- 依存方向: `npm run lint:backend:arch` 緑（engine→core 一方向、逆依存・循環なし）。`npm run lint:backend:boundary` 緑。
- 非劣化: `npm run test:backend` の harness golden（送信プロンプト列＋DB 最終状態）が分割前と一致。
- 静的解析: `npm run lint:backend:static` は 9 件（errcheck 1・errorlint 1・gosec 2・staticcheck 1・wrapcheck 4）。すべて移設前からの baseline で、本 task で新規追加した違反は 0（revive の stutter・package コメント欠落・公開関数コメントは本 task 内で全解消）。残 9 件は `backend-violation-cleanup` 管轄。
- `npm run verify:backend` exit 0。

## 最終検証（実機目視）

`npm run dev:wails:run` で実 app を起動し、`Innocence Lost - Quest Expansion.esp` を実画面から抽出→翻訳（実 LLM = LM Studio `hy-mt2-7b`、LAN `http://192.168.0.226:1234`）。197 件を翻訳し、移設コードの全 runtime 経路が従来どおり効くことを目視した。

- 起動成功: bootstrap の asset 読み（`os.Open`＋`rolespeech.ParseRoleSpeech`）が効き、Wails binding 接続（console "Connected to backend"）。
- 抽出→固有名派生→ingest: 197 件を取り込み（`readXMLDir`→`core/termxml.DeriveTermsFromFiles`→`termusage`/`termderive`、ingest）。
- 機械置換・固有名注入（`core/dictionary.Apply`、派生辞書）: `Grelod`→`グレロッド`、`Grelod the Kind`→`親切者のグレロッド`、`Riften`→`リフテン`、`Aventus Aretino`→`アベンタス・アレティノ`、`Honorhall Orphanage`→`オナーホール孤児院` が全行で一貫。
- 口調指示（`core/personatone`→`LinePersonas`、`core/linefeatures`→`tone.Classifier`）: Grelod=「口調: ぞんざい」、AventusAretino=「口調: 平明」、叙述文=「口調なし」と話者別。
- 実プロンプト再構成（api の `core/prompt.ComposePrompt`/`RenderPrompt`＋`dict.Apply`）: Grelod 台詞の system に「- 口調: ぶっきらぼうで乱暴な口調」＋「- 人称と言い回し: 一人称は「わたし」。年配の女性らしい…」（役割語は `assets/role-speech.tsv` を bootstrap が `os.Open`→`rolespeech.ParseRoleSpeech` で読んだ結果）が注入されている。

非劣化 golden（機械）と実機目視の両方で、package 分割後も振る舞いが分割前と一致することを確認した。`完了定義` の観測点（arch-lint・boundary・harness golden・カバレッジ・実画面）をすべて満たした。

## テスト設計

振る舞い不変の移設のため、新規テストは原則追加しない。既存テストを対応 package へ移し、安全網で非劣化を担保する。

- 純粋 class 単体: 各 `_test.go`（`dictionary_test.go`・`prompt_test.go`・`termderive_test.go`・`termusage_test.go`・`termxml_test.go` 純粋部・`role_speech_test.go`・`linefeatures_test.go`・`tone_catalog_test.go`）を新 package へ移し、package 宣言を合わせる。非公開→公開にしたシンボル（`ParseTermXML`・`IsBaseGame`・`RoleClassOfRace`・`Build*`・`PersonaMetaOf`）の参照を公開名へ直す。
- カバレッジ: 各純粋 package で `go test -cover` を手元実行し、不変ルール関数が 100% を保つ（恒久ゲートは持たない。`unit-test-plan.md` の方針に従う）。
- 非劣化の安全網: 既存 harness（`harness_test.go`）の golden 比較（送信プロンプト列＋DB 最終状態）が分割前と一致することを移設の正しさの根拠にする。
- 単体で書かない: os 読み wrapper（`LoadRoleSpeech`・`DeriveTermsFromXMLDir`、engine 残置）は統合（harness）に委ね、単体テストを新設しない。

## finalization

### 正本化判断

- 反映対象: `docs/architecture.md`（§3 backend 責務、§4.1 機械検査の component map、§7 ディレクトリ正本）。
- 影響範囲: 純粋ルールの層（core）新設と依存方向の変化。§5・§6・§8 は不変。
- 判断: 恒久仕様として承認（構造的変化のため反映が要る）。人間承認状態: 承認済み。

### 作業 commit

- commit hash: `5ae0536f`（branch `claude/engine-pure-rule-split`）。
- 変更: core 9 package 新設（移設）、engine/api/bootstrap/harness/goldcap の参照更新、`.go-arch-lint.yml`、`architecture.md` §3・§4.1・§7、active plan、`backend-violation-cleanup` plan の location 更新。
- 検証: `verify:backend` exit 0、実機目視 OK。残留リスク: なし（非劣化確認済み）。

### 正本反映

- §3: `engine` 記述を「プロンプト組み立て等の純粋ルールは core が持ち engine が束ねる」に修正。`core`（functional core、1 ルール 1 package、100% カバレッジ基準）の責務を 1 項追加。
- §4.1: component→ディレクトリ対応と `mayDependOn` を core 9 component と更新後の依存（engine→core/\*、api→dictionary・prompt、bootstrap/harness/goldcap→rolespeech ほか、core 内一方向、tone は core/tone）に差し替え。
- §7: `internal/core/`（純粋ルールの集積、1 ルール 1 package）を 1 項追加。`engine` 記述を orchestration として整えた。
- 根拠: 本 active plan（設計確定・実装範囲）。一時材料 `summary.md`・`design-review.md` は承認後に削除。

## 軽 / 重判定

- 画面が動くか: N（svelte 表示・props・story・fixture・表示構造を変えない）。
- `docs/architecture.md` 反映が要るか: Y（§4 の package 構成と依存方向が変わる。新 component を arch-lint へ登録する）。
- 判定: 片方 Y のため **重 task**。経路は `preparation` → `design` →（画面不動のため storybook bypass）→ `implementation` → `finalization`。
