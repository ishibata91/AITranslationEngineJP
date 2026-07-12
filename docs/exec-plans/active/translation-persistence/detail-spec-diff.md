# translation-persistence 詳細仕様差分

親要件・確定仕様・未決・回答を固定する。DDL・SQL は持たず、決定と固定名だけを残す。実 DDL は `db/` migration が正本。

## 用語

翻訳の永続化単位は、対象 plugin と 1 対 1 の「翻訳対象 plugin」とする。`plugin` と 1 対 1 で、やり直しは削除して作り直すため、独立した作業単位（Job）の名は付けず plugin を実体名にする。table は `target_plugin`（plugin ファイル名をキーにした登録表）とする。

## 親要件

`plan.md` を参照。翻訳を対象 plugin 単位で永続化し、起動時 DB 全消去（dev flush）を廃止して plugin 単位の削除でやり直す。後続 batch plan の土台にする。

## 確定仕様

### target_plugin テーブル

- `target_plugin` を足す。plugin ファイル名をキーにして plugin と 1 対 1。作成時刻（`created_at`）と選んだ plugin のフルパス（`source_path`）を持つ。
- `source_path` は一覧の行から実行・結果画面へ戻る導線に使う（結果表示と再実行に、元の Data フォルダのフルパスを復元する）。
- 翻訳開始時（`RunExtractAndTranslate` の先頭、抽出の前）に同 plugin の `target_plugin` を upsert する。開始時に作るので、翻訳が途中で失敗しても登録が残り、削除でやり直せる。upsert は衝突時に `source_path` だけ更新し、`created_at`（初回登録時刻）は保つ。
- 翻訳状態（未訳/訳済）を `target_plugin` に持たせない。状態は既存行の `status` に残す。`target_plugin` は識別と束ねに限定する。一覧の進捗（総数・訳済み数）は既存行から計算する。

### 束ねと削除（手続き削除、2026-07-12 確定）

削除方式は FK cascade でなく Go 側の明示 DELETE にする（人間確認済み）。FK cascade は対象スコープの約 11 表を作り直す必要があり、加えて Go 接続の FK 強制 ON の副作用（挿入順序依存の実行時失敗）を伴う。手続き削除は変更が Go に閉じ、既存表の作り直しと FK 強制を避けられる。

- `target_plugin` は追加のみ（既存表の作り直し・FK 追加はしない）。plugin 列（`narration.plugin` 等）と `info_plugin` は現状のまま、値は対象 plugin ファイル名で束ねる。
- 削除は `store.DeleteTargetPlugin` が 1 トランザクションで、対象 plugin の連関（`narration_mention`・`narration_described`・`line_mention`・`line_speaker`・`line_condition`）を親行 id の副問い合わせで先に消し、続いて本体（`narration`・`line`・`proper_noun`・`extracted_field`）と staging（`extracted_info_speaker`・`extracted_info_condition`、`info_plugin`）、最後に `target_plugin` 行を消す。
- 対象は plugin 列 / info_plugin 列が対象 plugin 値の行に限る。共有 entity（`speaker`/`race`/`faction`/`voice_type`）・横断辞書（`master_term`）・各キャッシュ・seed は削除文に含めない。

### C# スキーマ 1 回適用化（永続化の前提）

- C# 抽出器はこれまで毎抽出で `db/migrations/*.sql` を全連結して生実行していた。起動時 flush を廃止すると、`0009` の `proper_noun` 作り直し（`DROP`→`CREATE`）を毎抽出で再実行し、別 plugin の固有名訳を毎回消す。
- そこで C# 側に `SchemaMigrator` を足し、Go の `db.Apply` と同じく `PRAGMA user_version` で適用済み migration を飛ばす。適用済み（多くは Go 起動時に到達）なら C# は何もしない。破壊的 DDL を毎抽出で再実行しない。

### 削除で残すもの（plugin をまたぐ資産）

- `speaker`・`race`・`faction`・`voice_type` は entity の出自 plugin（多くはマスタ `Skyrim.esm`）でスコープされ、複数 target が共有する。削除で消さない。
- `master_term`（横断・権威訳）、`line_analysis`（本文ハッシュ横断キャッシュ）、`persona_character`（話者キーのキャッシュ）、global 設定・seed（`prompt_template`・`directive`・`record_type_master`・`tone_default`）も消さない。

### 起動時 flush 廃止

- `scripts/dev/run-wails.sh` の中心 DB 削除処理を外す。dev 起動で DB を維持する。やり直しは plugin 単位の削除 UI 操作へ移す。

## 未決 → 回答

| 未決 | 回答 | 根拠 |
|---|---|---|
| 永続化単位を独立した作業単位（Job）にするか、plugin を実体名にするか | plugin を実体名（`target_plugin`）にする | plugin と 1 対 1 で、やり直しは削除して作り直す方式のため、Job は履歴も状態も持たず名前を稼げない。実体は「翻訳した plugin」そのもの |
| entity（`speaker`/`race`/`faction`/`voice_type`）を削除で消すか | 消さない | 抽出器が entity の plugin を出自 mod で書く（`SpeakerSqliteWriter.cs:82,119`）。target でなく共有マスタなので、消すと別 plugin が壊れる |
| 束ねを新規 `job_id` 列にするか、既存 plugin 列で束ねるか | 既存 plugin 列で束ねる（新規列は足さない） | C# 抽出器が既に書く plugin をそのまま束ねキーにでき、C# へ id を渡す必要がない |
| 削除を FK cascade にするか手続き削除にするか | Go 側の明示 DELETE（手続き削除）にする | FK cascade は約 11 表の作り直しと FK 強制 ON の副作用を伴う。手続き削除は変更が Go に閉じ低リスク（2026-07-12 人間確定） |
| `target_plugin` に status/progress を持たせるか | 持たせない | 状態は既存行から導出する。batch の非同期状態は batch plan の内部へ閉じる |
| schema 変更・データ移行の扱い | schema は作り直し可、データ移行不要（人間確認済み） | dev 専用・greenfield。migration を書き直して FK を入れられる |

## scope 外

- 配送の振る舞いを抽象する interface（同期と batch が共有する port）は本 plan で作らない。後続 batch plan が provider の 2 つ目の port として足す。
- `target_plugin` に batch 固有の永続情報を持たせない。batch 固有情報は batch plan の内部テーブルが plugin へぶら下げて持つ。
