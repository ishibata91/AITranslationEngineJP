# 詳細仕様差分: term-target-rec-config

- `skill`: detail-spec-design
- `status`: ready-for-human-review
- `source_plan`: `./plan.md`
- `detail_spec_target`: `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/master-dictionary.md`
- `screen_design_diff`: `N/A`
- `component_diagram`: `N/A`

## 詳細仕様差分

### `term-translation-phase-REQ-003` 共通辞書一致で AI 翻訳対象を決める

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/term-translation-phase.md`

親要件:
共通辞書に完全一致する語は、フェーズ開始時の辞書参照に基づく置換対象として扱う。

仕様:
- 共通辞書はフェーズ開始時の辞書参照で固定する。
- 共通辞書の完全一致は、原語と REC が一致する場合だけ成立する。
- REC は、レコード種別とフィールド名の組として表現し、`BOOK:FULL`、`NPC_:FULL`、`NPC_:SHRT` のように `RECORD:FIELD` 形式で識別する。
- 共通辞書に完全一致する語は辞書置換対象にする。
- 共通辞書の完全一致条件外の語は AI 翻訳対象にする。
- 共通辞書の対象範囲を判定した後に AI 翻訳対象語が 0 件でも、単語翻訳フェーズは `Completed` として扱う。
- AI 翻訳対象語が 0 件の場合は、AIサービス未実行を結果として判断できる。

未決:
- なし

回答:
- なし

### `term-translation-phase-REQ-004` AIサービス応答をジョブ内辞書へ反映する

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/term-translation-phase.md`

親要件:
共通辞書対象外の用語と固有名詞は AI 翻訳対象とし、確定訳語を対象の翻訳ジョブ内辞書へ保存する。

仕様:
- 共通辞書対象外の用語と固有名詞は AIサービスへ送る。
- AIサービスへの実行単位は 1 対象語を基本とする。
- 一括処理を使う場合も、1 項目は 1 対象語に対応する。
- AIサービスへ渡す生成指示は、対象語、原文言語、訳文言語、応答対応に使う識別子を同じ実行単位に固定する。
- AIサービス応答は、対象語ごとに、原語と訳語の対応として検査する。
- 有効な応答は、要求した対象語と同じ原語を持ち、空ではない訳語を持つ応答である。
- AIサービスの有効な応答は、原語と訳語の対応を保持し、自動で確定訳語として扱う。
- 確定訳語は対象の翻訳ジョブ内辞書として保存する。
- 確定訳語は、対象の翻訳ジョブと単語翻訳の実行単位を追跡できる状態で保持する。
- 同一翻訳ジョブ、同一 REC、同一原語では一意の辞書項目として扱う。
- 別 REC の同一原語は別辞書項目として扱える。
- 翻訳ジョブ内辞書に保存する対応単位は、共通辞書 hit 判定の対応単位と同じ REC とする。

未決:
- なし

回答:
- なし

### `term-translation-phase-REQ-008` 単語翻訳対象 REC を 13 種別に固定する

- `変更種別`: 追加
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/term-translation-phase.md`

親要件:
利用者は、単語翻訳フェーズが処理対象とする REC の範囲を判断できる。

仕様:
- 単語翻訳対象 REC は、`BOOK:FULL`、`NPC_:FULL`、`NPC_:SHRT`、`ARMO:FULL`、`WEAP:FULL`、`LCTN:FULL`、`CELL:FULL`、`CONT:FULL`、`MISC:FULL`、`ALCH:FULL`、`RACE:FULL`、`INGR:FULL`、`SHOU:FULL` の 13 種別とする。
- 翻訳入力に出現する翻訳レコードまたは翻訳フィールドのうち、原語が空でなく、REC が 13 種別のいずれかに該当する場合だけ、単語翻訳フェーズの候補にする。
- 13 種別に該当しない REC は、単語翻訳フェーズの候補外として扱い、後続の本文翻訳フェーズの対象範囲に委ねる。
- `NPC_:FULL` と `NPC_:SHRT` は別 REC として扱い、同一原語でも別候補として識別できる。
- 単語翻訳対象 REC 集合は、XML 辞書取り込み対象 REC 集合と同一の 13 種別とする。
- 単語翻訳フェーズと XML 辞書取り込みは、同一の単語翻訳対象 REC 判定（仮称 `IsTermTarget`）を共有して対象判定を行う。

未決:
- なし

回答:
- なし

### `master-dictionary-REQ-004` XML から辞書データを取り込める

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/master-dictionary.md`

親要件:
利用者は XML ファイルを選択し、許可された REC だけから辞書データを取り込める。

仕様:
- 利用者は XML ファイルを選択して辞書データを取り込める。
- XML 取り込みは選択中ファイルを利用者が識別できる状態で開始する。
- XML 取り込みは xTranslator 形式の辞書 XML から単語を抽出できる。
- XML 取り込み対象 REC は、`BOOK:FULL`、`NPC_:FULL`、`NPC_:SHRT`、`ARMO:FULL`、`WEAP:FULL`、`LCTN:FULL`、`CELL:FULL`、`CONT:FULL`、`MISC:FULL`、`ALCH:FULL`、`RACE:FULL`、`INGR:FULL`、`SHOU:FULL` の 13 種別とする。
- REC は、レコード種別とフィールド名の組として `RECORD:FIELD` 形式で識別する。
- 13 種別の外の REC は、XML 取り込みの対象外として扱う。
- XML 辞書取り込み対象 REC 集合は、単語翻訳フェーズの対象 REC 集合と同一の 13 種別とする。
- XML 辞書取り込みと単語翻訳フェーズは、同一の単語翻訳対象 REC 判定（仮称 `IsTermTarget`）を共有して対象判定を行う。
- 既存に保存されている `DOOR:FULL`、`FLOR:FULL`、`FURN:FULL` を含む 13 種別外の master_dictionary レコードは、既存翻訳データ reset の対象として扱い、互換 migration を作らない。

未決:
- なし

回答:
- なし

## 根拠

- `source`:
  - `docs/exec-plans/active/term-target-rec-config/plan.md`
  - `docs/references/term-translation-target-record-candidates.md`
  - `docs/detail-specs/term-translation-phase.md`
  - `docs/detail-specs/master-dictionary.md`
  - `internal/service/term_translation_phase_service.go`（現行 `collectCandidates` に REC 絞り込みが存在しないことの根拠）
  - `internal/service/master_dictionary_service.go`（現行 `allowedImportREC` が 16 種別を保持することの根拠）
  - `internal/repository/processing_target_sqlite_repository.go`（現行 `processingTargetTermCountSQL` / `processingTargetTermListSQL` に REC 絞り込みが存在しないことの根拠）
  - 人間設計レビュー差し戻し記録（2026-05-31）: 両集合を 13 種別に統一し、判定関数を 1 つ（仮称 `IsTermTarget`）で共有する方針。
- `review`: 人間設計レビュー差し戻し対応後の再提出。`status: ready-for-human-review`。
- `validation`: 仕様差分のみのため検証未実行。`implementation-scope` 固定後に単体テストで責務単位を検証する予定。
