# ER 図ドラフト

`docs/spec.md` と `extractData.pas` をもとにした、概念レベルの ER 図たたき台。

- `extractData.pas` の抽出ロジックを正として、出力カテゴリと項目を整理
- 抽出 JSON は正本、DB は `PLUGIN_EXPORT` 単位の実行キャッシュとして扱う
- 内部主キーはシーケンシャル PK、外部 FormID は `form_id` として別保持する
- `dialogue_groups -> responses` の階層をそのままエンティティ化
- 辞書は単純化して `source_text` / `dest_text` のみ保持

## 入力データ ER 図

### JSON 入力

```mermaid
erDiagram
    PLUGIN_EXPORT ||--o{ DIALOGUE_GROUP : contains
    PLUGIN_EXPORT ||--o{ QUEST : contains
    PLUGIN_EXPORT ||--o{ ITEM : contains
    PLUGIN_EXPORT ||--o{ MAGIC : contains
    PLUGIN_EXPORT ||--o{ LOCATION : contains
    PLUGIN_EXPORT ||--o{ SYSTEM_RECORD : contains
    PLUGIN_EXPORT ||--o{ MESSAGE : contains
    PLUGIN_EXPORT ||--o{ LOAD_SCREEN : contains
    PLUGIN_EXPORT ||--o{ NPC : contains

    DIALOGUE_GROUP ||--o{ DIALOGUE_RESPONSE : has
    QUEST ||--o{ QUEST_OBJECTIVE : has
    QUEST ||--o{ QUEST_STAGE_LOG : has
    DIALOGUE_GROUP }o--o| QUEST : references
    DIALOGUE_RESPONSE }o--o| NPC : speaks
    MESSAGE }o--o| QUEST : references

    PLUGIN_EXPORT {
        bigint id PK
        string target_plugin
        string source_json_path
        datetime imported_at
    }

    DIALOGUE_GROUP {
        bigint id PK
        bigint plugin_export_id FK
        string form_id
        string editor_id
        string type
        string player_text
        string source
        string quest_id
        boolean is_services_branch
        string services_type
        string nam1
    }

    DIALOGUE_RESPONSE {
        bigint id PK
        bigint dialogue_group_id FK
        string form_id
        string editor_id
        string type
        integer response_order
        string source
        string text
        string prompt
        string topic_text
        string menu_display_text
        string speaker_id
        string voicetype
        integer index
        string previous_id
    }

    QUEST {
        bigint id PK
        bigint plugin_export_id FK
        string form_id
        string editor_id
        string source
        string name
        string type
    }

    QUEST_OBJECTIVE {
        bigint id PK
        bigint quest_id FK
        string objective_index
        string type
        string parent_id
        string parent_editor_id
        string text
    }

    QUEST_STAGE_LOG {
        bigint id PK
        bigint quest_id FK
        integer stage_index
        integer log_index
        string type
        string parent_id
        string parent_editor_id
        string text
    }

    ITEM {
        bigint id PK
        bigint plugin_export_id FK
        string form_id
        string editor_id
        string type
        string source
        string name
        string description
        string text
        string type_hint
    }

    MAGIC {
        bigint id PK
        bigint plugin_export_id FK
        string form_id
        string editor_id
        string type
        string source
        string name
        string description
    }

    LOCATION {
        bigint id PK
        bigint plugin_export_id FK
        string form_id
        string editor_id
        string type
        string source
        string name
        string parent_id
    }

    SYSTEM_RECORD {
        bigint id PK
        bigint plugin_export_id FK
        string form_id
        string editor_id
        string type
        string source
        string name
        string description
    }

    MESSAGE {
        bigint id PK
        bigint plugin_export_id FK
        string form_id
        string editor_id
        string type
        string source
        string text
        string title
        string quest_id
    }

    LOAD_SCREEN {
        bigint id PK
        bigint plugin_export_id FK
        string form_id
        string editor_id
        string type
        string source
        string text
    }

    NPC {
        bigint id PK
        bigint plugin_export_id FK
        string form_id
        string editor_id
        string type
        string source
        string name
        string race
        string voice
        string sex
        string class_name
    }

    DICTIONARY_ENTRY {
        bigint id PK
        string source_text
        string dest_text
    }
```

### 基盤マスター

```mermaid
erDiagram
    MASTER_PERSONA ||--o{ MASTER_PERSONA_ENTRY : contains
    MASTER_DICTIONARY ||--o{ MASTER_DICTIONARY_ENTRY : contains

    MASTER_PERSONA {
        bigint id PK
        string persona_name
        string source_type
        datetime built_at
    }

    MASTER_PERSONA_ENTRY {
        bigint id PK
        bigint master_persona_id FK
        string npc_form_id
        string npc_name
        string race
        string sex
        string voice
        string persona_text
    }

    MASTER_DICTIONARY {
        bigint id PK
        string dictionary_name
        string source_type
        datetime built_at
    }

    MASTER_DICTIONARY_ENTRY {
        bigint id PK
        bigint master_dictionary_id FK
        string source_text
        string dest_text
    }
```

## 翻訳ジョブ ER 図

```mermaid
erDiagram
    PLUGIN_EXPORT ||--o{ TRANSLATION_JOB : owns
    TRANSLATION_JOB ||--o{ JOB_DIALOGUE_GROUP : targets
    JOB_DIALOGUE_GROUP }o--|| DIALOGUE_GROUP : includes

    TRANSLATION_JOB ||--o{ JOB_RECORD : targets

    TRANSLATION_JOB ||--o{ JOB_PHASE_RUN : executes
    JOB_PHASE_RUN }o--o| AI_RUN : performed_by
    AI_RUN }o--|| AI_PROVIDER : uses

    TRANSLATION_JOB ||--o{ TRANSLATION_INSTRUCTION : composes

    TRANSLATION_JOB ||--o{ JOB_DICTIONARY_ENTRY : uses
    JOB_DICTIONARY_ENTRY }o--|| DICTIONARY_ENTRY : references
    TRANSLATION_JOB }o--o| MASTER_PERSONA : references
    TRANSLATION_JOB }o--o| MASTER_DICTIONARY : references

    PLUGIN_EXPORT {
        bigint id PK
        string target_plugin
        string source_json_path
        datetime imported_at
    }

    DICTIONARY_ENTRY {
        bigint id PK
        string source_text
        string dest_text
    }

    MASTER_PERSONA {
        bigint id PK
        string persona_name
        string source_type
        datetime built_at
    }

    MASTER_DICTIONARY {
        bigint id PK
        string dictionary_name
        string source_type
        datetime built_at
    }

    DIALOGUE_GROUP {
        bigint id PK
        bigint plugin_export_id FK
        string form_id
        string editor_id
        string type
        string player_text
        string source
        string quest_id
        boolean is_services_branch
        string services_type
    }

    TRANSLATION_JOB {
        bigint id PK
        bigint plugin_export_id FK
        bigint master_persona_id FK
        bigint master_dictionary_id FK
        string job_name
        string status
        string current_phase
        datetime started_at
        datetime finished_at
    }

    JOB_DIALOGUE_GROUP {
        bigint id PK
        bigint job_id FK
        bigint dialogue_group_id FK
    }

    JOB_RECORD {
        bigint id PK
        bigint job_id FK
        string target_category
        bigint target_id
        string status
        integer retry_count
        string translated_text
    }

    JOB_PHASE_RUN {
        bigint id PK
        bigint job_id FK
        string phase_code
        bigint ai_run_id FK
        string status
        datetime started_at
        datetime finished_at
    }

    AI_PROVIDER {
        bigint id PK
        string provider_name
        string provider_type
        boolean supports_batch
    }

    AI_RUN {
        bigint id PK
        bigint ai_provider_id FK
        string execution_mode
        string request_payload_hash
        string response_payload_hash
        string status
    }

    TRANSLATION_INSTRUCTION {
        bigint id PK
        bigint job_id FK
        string phase_code
        string target_category
        bigint target_id
        text instruction_text
    }

    JOB_DICTIONARY_ENTRY {
        bigint id PK
        bigint job_id FK
        bigint dictionary_entry_id FK
    }
```

## 入力データ補足

- `PLUGIN_EXPORT` は JSON のルートにある `target_plugin` を表す親エンティティ
- `PLUGIN_EXPORT` は JSON 原本に対応する実行キャッシュ親であり、ジョブ実行中だけ DB に入力データを保持する
- `dialogue_groups` は `DIALOGUE_GROUP`、その `responses` は `DIALOGUE_RESPONSE` として分離
- `quests`, `items`, `magic`, `locations`, `system`, `messages`, `load_screens` は、JSON のトップレベル配列ごとに独立エンティティ化
- `DIALOGUE_GROUP.nam1` は、`extractData.pas` の `ExtractDialogue` が条件付きで出力する補助項目
- `QUEST_OBJECTIVE` と `QUEST_STAGE_LOG` は、`extractData.pas` の `ExtractQuest` が出力する `objectives` / `stages` を分解したもの
- `ITEM.text` と `ITEM.type_hint` は、`extractData.pas` の `ExtractItem` が出力する追加プロパティ
- `npcs` は配列ではなく ID をキーにしたオブジェクトなので、永続化時は `NPC.id` をキーとして正規化する想定
- `CELL FULL` は `extractData.pas` では `cells` 配列ではなく `locations` 配列へ `type = "CELL FULL"` として入る
- `cells` 配列は `extractData.pas` 上は出力枠があるが、現行ロジックでは実質未使用
- `MASTER_PERSONA` と `MASTER_DICTIONARY` は、仕様上の基盤データとして JSON 入力とは独立に持つ
- 外部 FormID は `form_id` として別保持し、DB の関連は内部シーケンシャル PK で張る
- Mermaid では参照先コメントを列に埋め込まず、関係線と `FK` 表記で外部キーを表現する
- `DIALOGUE_GROUP.quest_id` と `MESSAGE.quest_id` は、現状の JSON では表示用文字列を含む参照なので、厳密 FK にする前にパース仕様を決める必要がある
- `DIALOGUE_RESPONSE.previous_id` も文字列参照なので、必要なら後で自己参照 FK に変換する

## 翻訳ジョブ補足

- `JOB_RECORD` はカテゴリ横断で使えるよう、`target_category` と `target_id` で翻訳対象を指す単純形にしている
- `JOB_PHASE_TYPE` はテーブル化せず、`phase_code` を定数としてアプリケーション側で管理する前提にした
- Mermaid を壊しやすいポリモーフィック関連は図から省略し、`JOB_RECORD` の属性で表現している
- `JOB_DIALOGUE_GROUP` は、ジョブがどの会話グループを対象にしているかを表す中間テーブル
- `JOB_DICTIONARY_ENTRY` は、ジョブが再利用する辞書項目を表す中間テーブル
- `TRANSLATION_JOB` は `PLUGIN_EXPORT` を必須参照し、完了後は同一 `PLUGIN_EXPORT` に未完了ジョブが残っていない場合のみ入力キャッシュを削除する

## 次に詰める候補

1. `quest_id` と `previous_id` を文字列のまま持つか、抽出時に正規化するか
2. `JOB_RECORD` をポリモーフィック参照のままにするか、カテゴリ別に分けるか
3. `NPC.voice` と `DIALOGUE_RESPONSE.voicetype` の関係を統一するか
4. `cells` が使われるサンプルを追加で見て、エンティティ化するか
