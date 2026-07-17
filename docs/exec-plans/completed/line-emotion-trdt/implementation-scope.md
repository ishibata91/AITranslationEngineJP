# 実装範囲: line-emotion-trdt

`design-module` の人間設計レビュー承認済み。scope の境界・依存・検証単位を固定する。詳細仕様は列挙せず、Claude 本体が文脈を保って実装する。

## 確定した設計選択（人間レビュー回答）

- **適用範囲**: 全台詞（名指し話者①・汎用台詞②・PC 発話③）に感情行を載せる。
- **感情行トーン**: 助言調。既存 `emotionAdvice`（上書きでなく助言）と一貫させる。
- **保存形**: `line.emotion_type` 列で持つ。
- **二重回避（②③の整理）**: TRDT 感情型が非 Neutral の台詞は TRDT 種別の感情行を出し、経路②③が現状使う本文推定の強度助言（`emotionAdvice`）は出さない。TRDT が Neutral または無しの台詞は、経路②③の本文推定助言を現状どおり維持する。

## scope 境界（触る層と粒度）

- **schema（`db/migrations` 新規）**: `extracted_info_emotion`（INFO 応答単位の感情。`info_plugin`・`info_form_id`・`ordinal`・`emotion_type`）を追加し、`line` へ `emotion_type` 列を追加する。C# 側 `SchemaMigrator.cs` も同 schema を冪等 ensure する（C#↔Go 契約の一致）。
- **抽出器（`tools/extractor`）**: `Model.cs` の `ResponseLine` へ感情型を足し、`PluginExtractor.cs` で `r.Emotion` を読む。`InfoConditionSqliteWriter.cs` に倣った新 writer で `extracted_info_emotion` へ書く。ordinal は本文と同じ出現順採番で `line` と結合可能にする。
- **取込（`internal/engine` 取込段・`internal/store`）**: `extracted_info_emotion` を `line.emotion_type` へ写す。
- **core（`internal/core/personatone`）**: 感情型8種→日本語感情語の写像（純粋ルール）。`BuildToneTraits`（経路①）と `BuildFreeToneTraits`（経路②③）へ感情行の加算と二重回避を組む。感情型の定数は `tone` へ置く（既存の口調定数と同居）。
- **model・結線（`internal/model`・`internal/engine`・`internal/store`）**: `Line` と `LinePersonaInput` へ感情型を足し、`LoadLinePersonas` と `freeTonePersona` が `line.emotion_type` を渡す。

## 依存順序

schema → 抽出器（C#）と取込（Go）→ core（`personatone` 写像・加算）→ engine 結線 → 実データ・実画面検証。

## 感情型8種→日本語感情語（写像の確定値）

Neutral は加算なし。Anger=怒り、Disgust=嫌悪、Fear=恐れ、Sad=悲しみ、Happy=喜び、Surprise=驚き、Puzzled=戸惑い。

## テスト設計（検証単位）

- **Go core 単体（`personatone`）**: 感情型→日本語感情語の写像、経路①②③で非 Neutral に感情行を加算、Neutral で加算しない、二重回避（TRDT 非 Neutral のとき本文推定助言を出さない）、基底口調・役割語・種族訛りが不変。
- **C#（`extractor.Tests`）**: 抽出器が INFO 応答の感情型を `extracted_info_emotion` へ書く。
- **取込（Go）**: `extracted_info_emotion` → `line.emotion_type` の写しが応答単位で対応する（純粋写像部分を単体、DB 込みは実データ観測へ）。
- **実データ**: 実 plugin（`Skyrim.esm` 等）抽出で `line.emotion_type` に感情型が入る。
- **実画面**: 翻訳実行で同一話者の異なる台詞に別々の感情行が載る（`http://localhost:34115`）。

## architecture.md 反映（finalization で実施）

- §6（C#↔Go の SQLite 契約）に `extracted_info_emotion` を1つ加える。
- §8（抽出器の責務）に「INFO 応答の感情型を書く」を追記する。
- 層構成・依存方向・Wails 境界は不変。
