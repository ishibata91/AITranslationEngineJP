# Design: dictionary-article-exclusion-and-derived-extraction

`spec.md` を優先する。`plan.md` の各要求だけを扱う。

---

## R-1 冠詞を除外した照合形を追加する

### 現況の理解

`internal/store/prebuilt_dictionary.go` は事前作成済み辞書から、原語と各意味の訳候補を読み取る。`internal/engine/engine.go` の `referencesForSource` は、読み取った原語を `internal/core/dictionary.Dictionary.Extract` へ渡して本文の参考語を選ぶ。`Dictionary.Extract` は登録した原語と完全に同じ文字列だけを、語境界と最長一致で選ぶ。根拠は `internal/store/prebuilt_dictionary.go` の `References`、`internal/engine/engine.go` の `bodyReferences` と `referencesForSource`、`internal/core/dictionary/dictionary.go` の `NewDictionary` と `Extract` である。

| | 単位 |
| --- | --- |
| 要求が扱う対象 | 事前作成済み辞書の原語から作る冠詞を除外した照合形 |
| 既存の受け皿 | 参考語ごとの原語 |

### あるべき形

Engineは、事前作成済み辞書の各原語を保持した参考語から、先頭の英語冠詞を一つだけ除いた照合形を作る。照合形が本文へ一致した場合も、本文翻訳へ渡す参考語は元の原語と全ての訳候補を保持する。既存の原語と冠詞を除外した照合形は、同じ最長一致と語境界の規則で一つの候補集合として選ぶ。

### 変更点

- `internal/core/dictionary/dictionary.go` に、照合に使う文字列と本文翻訳へ渡す参考語を分離して最長一致を選ぶ機能を置く。
- `internal/engine/engine.go` の `referencesForSource` を変更する。事前作成済み辞書の原語から冠詞を除外した照合形を作り、照合結果を元の参考語へ戻す。
- `internal/core/dictionary/dictionary_test.go` と `internal/engine/engine_test.go` に、語境界、最長一致、元の原語の保持を確認する例を置く。

```mermaid
flowchart LR
  A[事前作成済み辞書の原語] --> B[参考語]
  B --> C[完全一致だけを選ぶ]
  C --> D[本文翻訳]
```

```mermaid
flowchart LR
  A[事前作成済み辞書の原語] --> B[参考語]
  B --> C[原語と冠詞を除外した照合形を選ぶ]
  C --> D[元の原語を持つ参考語]
  D --> E[本文翻訳]
```

---

## R-2 辞書単語から限定的な派生照合形を抽出する

### 現況の理解

事前作成済み辞書の原語から本文用の派生照合形を作る処理はない。`internal/core/termderive/termderive.go` はNPC名から短縮形などを作るが、事前作成済み辞書の複合語には接続していない。根拠は `internal/engine/engine.go` の `referencesForSource` と `internal/core/termderive/termderive.go` の `DeriveTerms` である。

| | 単位 |
| --- | --- |
| 要求が扱う対象 | 一部の複合語と本文照合用の派生語 |
| 既存の受け皿 | 参考語ごとの原語 |

### あるべき形

Engineは、採用された複合語の規則だけから派生照合形を作る。派生照合形が本文へ一致した場合も、本文翻訳へ渡す参考語は元の原語、訳候補、品詞、カテゴリ、由来を保持し、意味を渡さない。派生照合形は辞書項目、レビュー、採否フラグへ保存しない。

### 変更点

- `internal/core/dictionary/dictionary.go` に、派生照合形を元の参考語へ対応付けたまま最長一致へ加える機能を置く。
- `internal/engine/engine.go` の `referencesForSource` を変更する。確定した派生規則から候補を作り、同期翻訳とbatch翻訳が共通で使う参考語組立てへ渡す。
- `internal/core/dictionary/dictionary_test.go` と `internal/engine/engine_test.go` に、派生照合形から元の参考語が渡る例を置く。

---

## R-3 誤一致を抑止する

### 現況の理解

`Dictionary.Extract` は語境界と長さ降順を守るが、事前作成済み辞書の候補にはstoplistを使わない。`internal/core/dictionary/stoplist.go` のstoplistは既存の翻訳前機械置換にだけ適用される。根拠は `internal/core/dictionary/dictionary.go` の `Extract`、`internal/core/dictionary/stoplist.go`、`internal/engine/engine.go` の `translationVocabulary` と `referencesForSource` である。

### あるべき形

冠詞を除外した照合形と派生照合形は、一般語、短すぎる語、複数の元の原語を持つ語を本文翻訳の候補に加えない。既存の原語は、事前作成済み辞書の読取り条件を変えずに、現在の語境界と最長一致の規則で扱う。

### 変更点

- `internal/core/dictionary/dictionary.go` に、派生照合形ごとの採用可否を判定する機能を置く。
- `internal/core/dictionary/dictionary_test.go` に、一般語、短すぎる語、複数の元の原語を持つ語が候補に入らない例を置く。

---

## R-4 参考語化した本文翻訳へ接続する

### 現況の理解

`model.PrebuiltDictionaryReference` は意味を持つ。`model.TranslationReference`、`prompt.BodyReference`、参考語snapshot、`api.TermView`、frontendの結果行は意味を持たない。完了済み `replace-extraction-to-prebuilt-dictionary` の設計と仕様は、本文翻訳、snapshot、結果表示から意味を除く契約を確定している。根拠は `internal/model/translation_reference.go`、`internal/core/prompt`、`internal/api/app.go` の `TermView` と `setResultSnapshot`、`frontend/src/gateway/translation-gateway.ts`、`docs/exec-plans/completed/replace-extraction-to-prebuilt-dictionary/design.md` と `spec.md` である。

### あるべき形

一致した冠詞を除外した照合形または派生照合形は、元の原語、訳候補、品詞、カテゴリ、由来を本文翻訳の参考語として渡す。意味は本文翻訳、送信時の参考語snapshot、結果表示へ渡さない。同期翻訳、batch翻訳、送信時の参考語snapshot、結果表示は同じ項目を保持する。

### 変更点

- `internal/engine/engine.go` の `referencesForSource` を変更する。照合形から元の参考語へ戻す。
- `internal/model/translation_reference.go`、`internal/core/prompt`、snapshot、`internal/api/app.go`、`frontend/src/gateway/translation-gateway.ts` は、意味を持たない既存の共有契約を維持する。
- `internal/engine/engine_test.go`、`internal/engine/batch_integration_test.go`、`internal/api/app_test.go` に、照合形から作る参考語が意味を送信、保存、表示しない例を置く。

---

## R-5 同期翻訳とbatch翻訳を一致させる

### 現況の理解

同期翻訳は `Engine.TranslateUntranslated` から `referencesForSource` と `composeBodyPrompt` を使う。batch翻訳は `internal/engine/batch.go` の `planBodyRequests` から同じ `composeBodyPrompt` と `referencesForSource` を使い、参考語とprompt hashを保存する。結果表示は `internal/api/app.go` の `setResultSnapshot` が保存済み参考語から再構成する。根拠は `internal/engine/engine.go`、`internal/engine/batch.go`、`internal/api/app.go` である。

### あるべき形

同期翻訳とbatch翻訳は、一つの参考語組立てから冠詞を除外した照合形、派生照合形、既存の原語を選ぶ。二つの経路は同じ最長一致規則、参考語の項目、保存内容、結果表示を使う。

### 変更点

- `internal/engine/engine.go` の共通の `referencesForSource` を変更する。
- `internal/engine/batch.go` は共通の参考語組立てを直接実装せず、既存の呼出しを維持する。
- `internal/engine/engine_test.go`、`internal/engine/batch_integration_test.go`、`internal/api/app_test.go` に、同期翻訳、batch翻訳、保存済み参考語、結果表示の一致を確認する例を置く。

---

## 検討が必要なこと

- R-2 とR-3の派生規則、最低文字数、一般語の判定、複数の元の原語を持つ派生照合形の扱いは、`plan.md` が未確定と明記している。人間の決定なしに一つの規則へ絞れない。
- なし。
