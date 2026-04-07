# データモデル / ER 図仕様

関連文書: [`index.md`](./index.md), [`spec.md`](./spec.md), [`architecture.md`](./architecture.md)

本書は、入力データ、基盤マスター、翻訳ジョブを含む現行データモデルの正本とする。
未解消論点は末尾に明示し、解決までは本書に記載した構造を現行仕様として扱う。

- xEdit 抽出 JSON の出力カテゴリと項目を canonical input として整理する
- 抽出 JSON は正本、DB は `PLUGIN_EXPORT` 単位の実行キャッシュとして扱う
- 内部主キーはシーケンシャル PK、外部 FormID は `form_id` として別保持する
- `dialogue_groups -> responses` の階層をそのままエンティティ化する
- raw JSON 互換項目と、DB 正規化後に採用する canonical 項目を分けて記述する
- 辞書は単純化して `source_text` / `dest_text` のみ保持する

## 入力データ ER 図

### JSON 入力

関連ファイル: [`input-data-er.d2`](./diagrams/er/input-data-er.d2), [`input-data-er.svg`](./diagrams/er/input-data-er.svg)

![入力データ ER 図](./diagrams/er/input-data-er.svg)

### 基盤マスター

関連ファイル: [`foundation-master-er.d2`](./diagrams/er/foundation-master-er.d2), [`foundation-master-er.svg`](./diagrams/er/foundation-master-er.svg)

![基盤マスター ER 図](./diagrams/er/foundation-master-er.svg)

## 翻訳ジョブ ER 図

関連ファイル: [`translation-job-er.d2`](./diagrams/er/translation-job-er.d2), [`translation-job-er.svg`](./diagrams/er/translation-job-er.svg)

![翻訳ジョブ ER 図](./diagrams/er/translation-job-er.svg)

## 入力データ補足

- `PLUGIN_EXPORT` は JSON のルートにある `target_plugin` を表す親エンティティ
- `PLUGIN_EXPORT` は JSON 原本に対応する実行キャッシュ親であり、ジョブ実行中だけ DB に入力データを保持する
- `dialogue_groups` は `DIALOGUE_GROUP`、その `responses` は `DIALOGUE_RESPONSE` として分離する
- `quests`, `items`, `magic`, `locations`, `system`, `messages`, `load_screens` は、JSON のトップレベル配列ごとに独立エンティティ化する
- raw JSON には互換のため `cells` ルートも出るが、現行入力では空配列として扱う
- `DIALOGUE_GROUP.nam1` は抽出 JSON が条件付きで出す補助項目として保持する
- `QUEST_OBJECTIVE` と `QUEST_STAGE_LOG` は `objectives` / `stages` を分解したものとして扱う
- `ITEM.text` と `ITEM.type_hint` は item 系レコードの追加プロパティとして扱う
- `npcs` は配列ではなく ID をキーにしたオブジェクトなので、永続化時は `NPC.id` をキーとして正規化する想定とする
- `DIALOGUE_GROUP.quest_id` と `MESSAGE.quest_id` は抽出時に正規化し、表示用文字列とは別に `QUEST.id` を参照する FK として保持する
- `DIALOGUE_RESPONSE.previous_response_id` は抽出時に正規化し、前段応答を自己参照 FK で保持する
- raw JSON の `voicetype` は互換入力として扱い、canonical では `speaker_*` と `NPC.voice` を基準にする
- `CELL FULL` は `locations` 配列へ `type = "CELL FULL"` として入る前提で、独立した `CELL` エンティティは持たず `LOCATION` に集約する
- `MASTER_PERSONA` と `MASTER_DICTIONARY` は、仕様上の基盤データとして JSON 入力とは独立に持つ
- `TRANSLATION_UNIT` は import 時に各 translatable field から生成する canonical 翻訳単位であり、xTranslator 出力と標準配布形式出力の両方の基準にする
- Mermaid では参照先コメントを列に埋め込まず、関係線と `FK` 表記で外部キーを表現する

## 翻訳ジョブ補足

- `TRANSLATION_JOB` は 1 つ以上の `PLUGIN_EXPORT` を `JOB_PLUGIN_EXPORT` 経由で参照し、複数入力ファイルをまとめて 1 ジョブで扱える
- `TRANSLATION_UNIT` は `record_signature` / `field_name` / `form_id` / `editor_id` / `source_text` を保持し、xTranslator XML の `<String>` を lossless に再構成できる
- `JOB_TRANSLATION_UNIT` はジョブごとの翻訳進捗、リトライ回数、翻訳文、xTranslator `Status` 相当の `translation_status_code` を保持する
- `JOB_PHASE_TYPE` はテーブル化せず、`phase_code` を定数としてアプリケーション側で管理する前提にした
- `JOB_DIALOGUE_GROUP` は、ジョブがどの会話グループを対象にしているかを表す中間テーブル
- `JOB_DICTIONARY_ENTRY` は、ジョブが再利用する辞書項目を表す中間テーブル
- `JOB_PERSONA_ENTRY` は mod 追加 NPC を含むジョブ内ペルソナを保持し、`MASTER_PERSONA` は基盤データ、`JOB_PERSONA_ENTRY` は実行時生成データとして分離する
- `FOUNDATION_PHASE_RUN` は基盤構築フロー側の AI 実行単位を表し、`MASTER_PERSONA` を AI 実行基盤へ接続するための materialization 境界として扱う
- `AI_RUN` は provider 側の run / batch 識別子と時刻、失敗理由を保持し、中断 / 再開 / 進捗観測に使う
- `JOB_OUTPUT_ARTIFACT` は format ごとの出力ファイルと生成状態を保持し、UI 観測と再出力判断に使う
- 入力キャッシュ削除判定は `PLUGIN_EXPORT` 単位で行い、同一 `PLUGIN_EXPORT` を参照する未完了ジョブが `JOB_PLUGIN_EXPORT` 上に残っていない場合のみ削除する
