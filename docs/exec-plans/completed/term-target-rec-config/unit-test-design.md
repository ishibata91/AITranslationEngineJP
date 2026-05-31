# 単体テスト設計: term-target-rec-config

- `skill`: test-design
- `status`: ready-for-human-review
- `source_plan`: `./plan.md`
- `source_detail_spec_diff`: `./detail-spec-diff.md`
- `source_design_diff`: `./design-diff.md`
- `scenario_test_design`: `./test-design.csv`

## 目的

承認済み詳細仕様差分（`term-translation-phase-REQ-003` / `REQ-004` / `REQ-008`、`master-dictionary-REQ-004`）と設計差分図に従い、単語翻訳対象 REC 13 種別を共通 config `IsTermTarget` で固定し、`collectCandidates`、`processingTargetTermCountSQL` / `processingTargetTermListSQL`、`isAllowedImportREC` の各責務を独立して単体で証明する観点を固定する。シナリオ観点は `./test-design.csv` に分離する。

## 検証単位の分担

| 検証単位 | 主な責務 | 対応する仕様 |
| --- | --- | --- |
| 共通 config `IsTermTarget` 単体 | 13 種別集合の所属判定 | `term-translation-phase-REQ-008`、`master-dictionary-REQ-004` |
| `collectCandidates` 単体 | 翻訳入力から 13 種別 REC 候補を組み立て、candidate key を `RECORD:FIELD` 形式で生成 | `term-translation-phase-REQ-008`、`REQ-004` |
| `processingTargetTermCountSQL` / `processingTargetTermListSQL` の REC 絞り込み統合 | SQL 上の `record_type IN (...)` 絞り込みと `dictionary_scope = RECORD:FIELD` 比較 | `term-translation-phase-REQ-003`、`REQ-008` |
| `isAllowedImportREC` 単体 | XML 取り込みの REC 許可判定が共通 config を参照する | `master-dictionary-REQ-004` |

## 単体観点

### U-TTRC-001 `IsTermTarget` が 13 種別だけ true を返す

- 分類: 正常 / 境界
- 対象: `internal/recclassification` の `IsTermTarget(rec string) bool`
- 仕様根拠: `term-translation-phase-REQ-008`、`master-dictionary-REQ-004`
- 入力代表:
  - 13 種別の各 REC 文字列（`BOOK:FULL` / `NPC_:FULL` / `NPC_:SHRT` / `ARMO:FULL` / `WEAP:FULL` / `LCTN:FULL` / `CELL:FULL` / `CONT:FULL` / `MISC:FULL` / `ALCH:FULL` / `RACE:FULL` / `INGR:FULL` / `SHOU:FULL`）
  - 仕様変更で除外された 3 種別（`DOOR:FULL` / `FLOR:FULL` / `FURN:FULL`）
  - REC 形式違反（空文字、`NPC_`、`:FULL`、`NPC_:`、`npc_:full` 等）
- 期待値:
  - 13 種別の各入力は `true` を返す。
  - `DOOR:FULL` / `FLOR:FULL` / `FURN:FULL` を含む 13 種別外入力は `false` を返す。
  - REC 形式違反入力は `false` を返す（例外送出はしない）。
- 失敗診断: 13 種別集合の欠落、誤追加、`NPC_:FULL` と `NPC_:SHRT` の取り違えを assertion 名から識別できるようにする。
- 補足: table-driven tests とする。集合定義の重複や順序依存を入れない。

### U-TTRC-002 `collectCandidates` が 13 種別だけ候補化する

- 分類: 正常 / 境界
- 対象: `internal/service/term_translation_phase_service.go` の `collectCandidates`
- 仕様根拠: `term-translation-phase-REQ-008`
- 入力代表:
  - 13 種別 REC を持つ原語非空の翻訳レコード（複数 REC を混在）。
  - 13 種別外 REC を持つ原語非空の翻訳レコード（`DOOR:FULL` / `FLOR:FULL` / `FURN:FULL` を含む）。
  - 13 種別 REC で原語が空の翻訳レコード。
- 期待値:
  - 候補集合に 13 種別 REC かつ原語非空のレコードだけが含まれる。
  - 13 種別外 REC のレコードは候補集合から除外される。
  - 原語空のレコードは候補集合から除外される。
- 失敗診断: 失敗した assertion から「REC 絞り込み欠落」「原語空チェック欠落」のどちらかを識別できる名前にする。

### U-TTRC-003 `collectCandidates` の candidate key が `RECORD:FIELD` 形式で `NPC_:FULL` と `NPC_:SHRT` を分離する

- 分類: 境界
- 対象: `collectCandidates` が組み立てる candidate key（`RECORD:FIELD` 形式）
- 仕様根拠: `term-translation-phase-REQ-008` の `NPC_:FULL` と `NPC_:SHRT` を別 REC として扱う仕様
- 入力代表:
  - `NPC_:FULL` で原語 `X` を持つ翻訳レコード。
  - `NPC_:SHRT` で原語 `X` を持つ翻訳レコード（同一原語、別 REC）。
- 期待値:
  - 候補集合に `NPC_:FULL` 由来 candidate と `NPC_:SHRT` 由来 candidate が別キーで含まれる。
  - candidate key が `record.RecordType + ":" + field.SubrecordType` 形式である。
- 失敗診断: candidate key が `record_type` のみで集約されている退行を検出できる名前にする。

### U-TTRC-004 `processingTargetTermCountSQL` と `processingTargetTermListSQL` が 13 種別だけを集計対象にする

- 分類: 正常 / 境界
- 対象: `internal/repository/processing_target_sqlite_repository.go` の `processingTargetTermCountSQL` / `processingTargetTermListSQL`
- 仕様根拠: `term-translation-phase-REQ-003` / `REQ-008`
- 区分: SQLite を fake/in-memory として使う統合単体（repository 層単独）
- 入力代表（seed データ）:
  - 13 種別 REC のレコード（各 REC 1 件以上）。
  - 13 種別外 REC のレコード（`DOOR:FULL` / `FLOR:FULL` / `FURN:FULL`）。
  - `NPC_:FULL` と `NPC_:SHRT` で同一原語のレコード。
  - 共通辞書 `master_dictionary` に `NPC_:FULL` の原語 `X` を `dictionary_scope = 'NPC_:FULL'` で保持。
- 期待値:
  - `processingTargetTermCountSQL` が 13 種別 REC のレコードだけを件数集計する。
  - `processingTargetTermListSQL` が 13 種別 REC のレコードだけを返す。
  - 13 種別外 REC（`DOOR:FULL` / `FLOR:FULL` / `FURN:FULL`）は件数にも一覧にも含まれない。
  - `NPC_:FULL` と `NPC_:SHRT` は別 REC として件数と一覧に別行で現れる。
  - `master_dictionary` 完全一致（同一原語かつ同一 `dictionary_scope = RECORD:FIELD`）のレコードは AI 対象から除外される。
- 失敗診断: 「REC 絞り込み欠落」「dictionary_scope が RECORD:FIELD でなく record_type で比較されている」「`NOT EXISTS` 節の REC 比較欠落」のいずれを退行させたか識別できる名前にする。
- 補足: 件数 SQL と一覧 SQL は同じ seed で同じ REC 集合の数を返すことを比較する観点を 1 件含める（処理対象一覧 ↔ 実行対象 整合の repository 側証明）。

### U-TTRC-005 `isAllowedImportREC` が `IsTermTarget` 経由で 13 種別だけ許可する

- 分類: 正常 / 境界
- 対象: `internal/service/master_dictionary_service.go` の `isAllowedImportREC`
- 仕様根拠: `master-dictionary-REQ-004`
- 入力代表:
  - 13 種別の各 REC。
  - 旧仕様で許可されていた `DOOR:FULL` / `FLOR:FULL` / `FURN:FULL`。
  - REC 形式違反（空文字など）。
- 期待値:
  - 13 種別の入力は許可される。
  - `DOOR:FULL` / `FLOR:FULL` / `FURN:FULL` は許可されない。
  - REC 形式違反は許可されない。
- 失敗診断: `allowedImportREC` 固定 map の残存退行、または `IsTermTarget` 参照漏れを識別できる名前にする。
- 補足: `IsTermTarget` 単体（U-TTRC-001）と重複する集合判定は最小限にし、本観点では「`isAllowedImportREC` が `IsTermTarget` の結果に従う責務境界」だけを証明する。

## 分類別の網羅状況

| 分類 | 単体側で扱う観点 |
| --- | --- |
| 正常 | U-TTRC-001（13 種別 true）、U-TTRC-002（候補化）、U-TTRC-004（件数/一覧の正常）、U-TTRC-005（13 種別許可） |
| 代替 | 該当なし。REC 集合判定に「中断、取り消し、検索結果なし、差分なし」の代替経路は存在しない |
| 例外 | 該当なし。REC 値は静的な列挙集合で例外経路を構成しない |
| 境界 | U-TTRC-001（13 種別外、形式違反）、U-TTRC-002（13 種別外、原語空）、U-TTRC-003（NPC_:FULL と NPC_:SHRT の分離）、U-TTRC-004（13 種別外除外、共通辞書 hit 除外、件数=一覧 整合）、U-TTRC-005（13 種別外不許可） |

## 単体側で扱わない観点

- 画面表示、状態遷移、開始操作の UI 確認はシナリオ観点（`./test-design.csv`）に固定する。
- AIサービス境界の応答検証は本 task 範囲外（既存仕様）であり、本単体観点では扱わない。
- 既存 `dictionary_scope` を `record_type` から `RECORD:FIELD` へ変更する migration の検証は、既存翻訳データ reset 前提（人間判断 2026-05-31）に従い扱わない。

## 根拠

- `source`:
  - `./plan.md`
  - `./detail-spec-diff.md`
  - `./design-diff.md`
  - `docs/detail-specs/term-translation-phase.md`
  - `docs/detail-specs/master-dictionary.md`
  - `docs/usecases/uc-translation-management.md`（`翻訳段階を開始する`、`処理対象を確認する`）
  - `docs/usecases/uc-master-dictionary.md`（`XMLから辞書エントリを取り込む`）
  - `docs/coding-guidelines-tests.md`
- `review`: 人間設計レビュー（2026-05-31）承認済み確定事項に従い 13 種別と `IsTermTarget` 共有方針で固定。`status: ready-for-human-review`。
- `validation`: 観点固定のみ。実装は `implementation-module` 起動後に行う。
