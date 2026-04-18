# データモデル / ER 図仕様

関連文書: [`index.md`](./index.md), [`spec.md`](./spec.md), [`architecture.md`](./architecture.md), [`diagrams/conceptual/combined_perspective.puml`](./diagrams/conceptual/combined_perspective.puml)

本書は、入力データ、基盤データ、翻訳ジョブを含むデータモデルの正本とする。
概念モデルを DB 化するため、ER は `combined_perspective.puml` の主語と責務境界を基準にする。

## ER 図

正本: [`combined-data-model-er.d2`](./diagrams/er/combined-data-model-er.d2)

review artifact: [`combined-data-model-er.svg`](./diagrams/er/combined-data-model-er.svg)

ER 図の正本は 1 枚だけとする。
旧分割 ER の `input-data-er.d2`、`foundation-master-er.d2`、`translation-job-er.d2` は廃止済みの pointer として扱う。

## 全体方針

- `TRANSLATION_JOB` は 1 つの `X_EDIT_EXTRACTED_DATA` だけを参照する
- ジョブ状態は `JOB_PHASE_RUN` 群から集約する
- 翻訳結果と出力ステータスは `JOB_TRANSLATION_FIELD` に保持する
- `PERSONA` と `DICTIONARY_ENTRY` はマスター / ジョブ内を同じテーブルで扱う
- フェーズ別 AI 設定、指示構成、最終 AI 実行情報は `JOB_PHASE_RUN` に保持する

## 入力データ

`X_EDIT_EXTRACTED_DATA` は xEdit 抽出 JSON の取込単位を表す。
`TRANSLATION_RECORD` は FormID、EditorID、RecordType を持つ Skyrim/xTranslator 上の識別単位を表す。

`NPC_RECORD` は `TRANSLATION_RECORD` の派生テーブルとして扱う。
Book、Item、Magic、Quest、Dialogue、Message、Location などは専用派生を作らず、標準の `TRANSLATION_RECORD` と `TRANSLATION_FIELD` で扱う。

## 翻訳フィールド

`TRANSLATION_FIELD` は翻訳対象になる原文フィールドを表す。
訳文と出力ステータスはジョブごとに異なるため、ここではなく `JOB_TRANSLATION_FIELD` に保持する。

`TRANSLATION_FIELD_DEFINITION` は `RecordType + SubrecordType` ごとの説明テーブルである。
AI 向け説明、翻訳対象フラグ、順序あり、順序スコープ、参照要件を保持する。

`TRANSLATION_FIELD_RECORD_REFERENCE` は DB 上の多対多中間テーブルである。
概念モデル上の汎用参照箱ではなく、翻訳フィールドが発話者 NPC、親クエスト、会話トピックなどの別レコードを参照するための実装上の関係として扱う。

## ペルソナと辞書

`PERSONA` は NPC と 1:1 で紐づく。
マスターペルソナはジョブ内ペルソナ生成の重複スキップに使うため、同一 NPC にマスターとジョブ内を同時保持しない。

マスター / ジョブ内、適用範囲、作成経路は `persona_kind`、`persona_scope`、`persona_source` で区別する。
`PERSONA_FIELD_EVIDENCE` はペルソナ生成根拠の翻訳フィールドを保持する。
`DICTIONARY_ENTRY` はマスター辞書とジョブ内辞書を同じテーブルで扱い、`dictionary_kind`、`dictionary_scope`、`dictionary_source` で区別する。

## ジョブとフェーズ

`JOB_PHASE_RUN` は翻訳ジョブ内のフェーズ実行だけを表す。
マスターペルソナ生成やマスター辞書構築などの基盤構築は `JOB_PHASE_RUN` に含めず、`PERSONA` / `DICTIONARY_ENTRY` の作成経路で表す。

フェーズ再実行は同じ `JOB_PHASE_RUN` の状態を戻す扱いにする。
Attempt 履歴テーブルは持たない。

`PHASE_RUN_TRANSLATION_FIELD`、`PHASE_RUN_PERSONA`、`PHASE_RUN_DICTIONARY_ENTRY` は、フェーズが対象にしたジョブ内翻訳フィールド、ペルソナ、辞書項目を表す。

## 出力

`TRANSLATION_ARTIFACT` はジョブが生成する成果物を表す。
標準出力形式は xTranslator 互換 XML とする。

`XTRANSLATOR_OUTPUT_ROW` は xTranslator 互換 XML の各出力行を表す。
1 つの `JOB_TRANSLATION_FIELD` に対応し、EDID、REC、FIELD、FORMID、Source、Dest、Status を保持する。

## Migration 化時の注意

- D2 は概念 ER の正本であり、SQLite migration では `NOT NULL`、`UNIQUE`、index、cascade 方針を別途固定する
- `kind`、`scope`、`source`、`state`、`phase_type` は初期 migration では文字列列として扱い、定数はアプリケーション側で管理する
- `credential_ref` は暗号化済み API key そのものではなく、secret store への参照だけを保持する
