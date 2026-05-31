# term-target-rec-config

## 依頼要約

単語翻訳フェーズの単語対象 REC と XML 辞書取り込み対象 REC を、Go 側の共通 config から参照する形へ整理する。
単語翻訳フェーズは 13 種別（`BOOK:FULL`、`NPC_:FULL`、`NPC_:SHRT`、`ARMO:FULL`、`WEAP:FULL`、`LCTN:FULL`、`CELL:FULL`、`CONT:FULL`、`MISC:FULL`、`ALCH:FULL`、`RACE:FULL`、`INGR:FULL`、`SHOU:FULL`）だけを対象にし、XML 辞書取り込み対象は既存 16 種別を維持する。
両者は別集合であり、共通 config から判定関数経由で参照できる構造に揃える。

参考:

- `internal/service/term_translation_phase_service.go`（`collectCandidates`）
- `internal/repository/processing_target_sqlite_repository.go`（`processingTargetTermCountSQL` / `processingTargetTermListSQL`）
- `internal/service/master_dictionary_service.go`（`allowedImportREC`）
- `docs/references/term-translation-target-record-candidates.md`
- `docs/detail-specs/master-dictionary.md`

確認済み前提:

- 既存の単語翻訳データは reset 前提で扱う。互換 migration は不要。
- 単語翻訳フェーズで `NPC_:FULL` と `NPC_:SHRT` は別候補として区別する必要がある。
- 共通辞書完全一致レコードは従来通り AI 対象から除外する。

## 分岐元

- 作業 branch: `codex/job-run-phase-fetch-redesign`（現 branch を継続使用、人間判断による）
- 分岐元 branch: `master`
- 分岐元 commit: `cb3bdbd44521b5034f68307afe1d68e779c751fa`
- task 開始時 HEAD: `b7c0df4de79421ddb0ea33a3620bec71d35a166d`

## 想定 Y/N 評価

| 想定 | Y/N | 根拠 |
| --- | --- | --- |
| 仕様変更または仕様追加がある | Y | 単語翻訳フェーズ候補が「全非空 source_text」から「13 種別 REC のみ」に絞られる。`internal/service/term_translation_phase_service.go:1667` `collectCandidates` の絞り込み欠如、`internal/repository/processing_target_sqlite_repository.go:247` `processingTargetTermCountSQL` の絞り込み欠如を変更する。 |
| 画面変更がある | N | 単語翻訳画面の表示構造、文言、layout は不変。処理対象一覧の件数と項目が減るのみで、画面側ロジックは変えない。 |
| 内部構造変更がある | Y | REC 分類を共通 config として新規に置き、`master_dictionary_service.go:28` `allowedImportREC` と service/repository の単語対象判定を共通 config 経由へ寄せる。`dictionary_scope` の意味も「`record_type` のみ」から「`RECORD:FIELD` 形式の REC」へ変える。 |
| 画面の表示変更がある | N | 同上、表示は変えない。`storybook-module` の起動条件にあたらない。 |
| frontend ロジック変更がある | N | state、API、Wails bridge、ルーティング、副作用、フォーム validation のいずれも触らない。 |
| backend 変更がある | Y | Go の service / repository / 共通 config の 4 箇所を変更する。 |
| frontend と backend を接続する | N | bridge IF も DTO も変えない。`internal/wails/bridge` の term-translation 入出力は不変。 |
| 実装済み責務を独立に証明したい | Y | REC 判定関数、collectCandidates の絞り込み、SQL の絞り込み、XML import の許可判定はそれぞれ独立した責務として単体テストで証明する。 |
| 実行時にしか確定しない値または原因分離が要る分岐がある | N | REC 値は静的な列挙集合で、観測ログ追加で原因分離する種類の分岐ではない。 |

省略 artifact:

- 画面設計差分: 表示変更なしのため省略。`storybook-module` も起動しない。

並列実行:

- `実装範囲` と `テスト設計` は `人間設計レビュー` 承認後に並列起動可。

非互換変更メモ:

- 既存の `dictionary_scope` は `record_type` だけを保持しているため、`NPC_:FULL` と `NPC_:SHRT` を分離するには `dictionary_scope` の表現を `RECORD:FIELD` 形式に変える必要がある。既存翻訳データは reset 前提で互換 migration を作らない（人間判断 2026-05-31）。

## 実装完了記録

- wave-1 B1: `internal/recclassification/term_target.go` 新規追加（`IsTermTarget`, `TermTargets`, `TermTargetRECList`）
- wave-2 B2: `internal/service/term_translation_phase_service.go` `collectCandidates` を REC 単位で絞り込み（既存テスト fixture も `RECORD:FIELD` 形式へ追従修正）
- wave-2 B3: `internal/repository/processing_target_sqlite_repository.go` 2 SQL を関数化、`recclassification.TermTargetRECList()` 由来の IN 句で絞り込み。`.go-arch-lint.yml` に `recclassification` component を追加
- wave-2 B4: `internal/service/master_dictionary_service.go` `allowedImportREC` 削除、`isAllowedImportREC` を `recclassification.IsTermTarget` 呼び出しへ、`categoryFromREC` から DOOR/FLOR/FURN 分岐削除
- wave-3 UT1: `internal/recclassification/term_target_test.go` 13 種別判定・防御コピーの単体検証
- wave-3 UT2: `internal/service/term_translation_phase_service_test.go` に `collectCandidates` の REC 絞り込み・key 形式・NPC_:FULL/SHRT 分離の 4 テスト追加
- wave-3 UT3: `internal/repository/processing_target_sqlite_repository_test.go` に SQL 絞り込み・`dictionary_scope` REC 比較・NPC_ 分離の 3 テスト追加
- wave-3 UT4: `internal/service/master_dictionary_service_test.go` 新規追加（許可 REC、category 分岐）
- wave-3 ST1: `internal/service/term_translation_phase_scenario_test.go` 新規追加（処理対象一覧と実行対象の REC 集合一致 2 シナリオ）
- arch lint 調整: `service.mayDependOn` に `infra_sqlite` を追加（scenario test で同一 SQLite を service と repository 両方から読むため）

## 最終検証

- `python3 scripts/harness/run.py --suite backend-local`: 通過（2026-05-31）
  - backend lint, backend test ともに pass

## 正本化判断

### 仕様変更対象

- `term-translation-phase-REQ-003`: 完全一致条件に「原語と REC が一致」を追加
- `term-translation-phase-REQ-004`: 一意キーを「同一 REC、同一原語」へ変更
- `term-translation-phase-REQ-008`: 単語翻訳対象 13 種別を新規固定
- `master-dictionary-REQ-004`: 16 種別 → 13 種別へ縮小、両集合の同一性と判定関数共有を明示

### 影響範囲（対象 docs）

- `docs/detail-specs/master-dictionary.md`（許可 REC リスト、判定関数共有）
- `docs/detail-specs/term-translation-phase.md`（13 種別、REC 一致、NPC_:FULL/SHRT 分離）
- `docs/references/term-translation-target-record-candidates.md`（13 種別と除外 3 種別の説明）

### 判断結果

- 恒久仕様として承認する（2026-05-31 人間レビュー承認済み）。
- `import 拡張`（dialogue_groups 以外を parse する変更）は別 task として切り出す（2026-05-31 人間判断）。
- DOOR:FULL / FLOR:FULL / FURN:FULL は仕様から外れたが、既存翻訳データは reset 前提で互換配慮不要（2026-05-31 人間判断）。

### 人間承認状態

- 承認済み（design-module の人間設計レビューで 1 回差し戻し後、最終承認）。

## 後続モジュール引き継ぎ

- task-id: `term-target-rec-config`
- 依頼要約: 上記
- 想定 Y/N、artifact 集合、設計成果物、人間確認は design-module 入口で扱う。
