# translation-input-import-non-dialogue

## 依頼要約

xEdit 抽出 JSON（例: `dictionaries/Lucien.esp_Export.json`）の `dialogue_groups` 以外の配列（`items`、`cells`、`locations`、`magic`、`messages`、`quests`、`responses`、`stages`、`system`、`objectives`、`load_screens`）を翻訳入力取り込みで parse し、単語翻訳フェーズ対象 13 種別 REC（BOOK:FULL、NPC_:FULL、NPC_:SHRT、ARMO:FULL、WEAP:FULL、LCTN:FULL、CELL:FULL、CONT:FULL、MISC:FULL、ALCH:FULL、RACE:FULL、INGR:FULL、SHOU:FULL）のレコードを DB に格納できるようにする。

## 背景（観察事実）

- 完了 task `term-target-rec-config`（2026-05-31 merge）で単語翻訳フェーズ対象 REC を 13 種別へ統一済み。
- 実機で Lucien.esp_Export.json をロードすると単語翻訳フェーズの処理対象が 0 件になる。
- DB の TRANSLATION_FIELD 内訳は INFO/NAM1 = 4563 行、DIAL/FULL = 488 行のみ。BOOK/NPC_/ARMO/CELL/LCTN/CONT/MISC/WEAP 等は格納されていない。
- 原因: `internal/service/translation_input_import_service.go:312` の `translationInputDocument` が `DialogueGroups []translationInputDialogueGroup` だけを `json:"dialogue_groups"` で受け取り、他配列を parse しない。
- Lucien.esp_Export.json のトップ配列: `cells`、`dialogue_groups`、`items`、`load_screens`、`locations`、`magic`、`messages`、`objectives`、`quests`、`responses`、`stages`、`system`。
- xEdit 抽出側（`extractData.pas`）は前 task で NPC SHRT 抽出と汎用 `ExtractNamedRecord` 対応を追加済み。

## 親要件（推定）

- 単語翻訳フェーズの処理対象が、xEdit 抽出 JSON 内の 13 種別 REC レコードを正しく拾えること。
- master dictionary の XML 取り込み対象 13 種別と整合すること（両集合は同一の `recclassification.IsTermTarget` を共有）。

## 関連ファイル

- `internal/service/translation_input_import_service.go`（取り込み parse の正本）
- `internal/recclassification/term_target.go`（13 種別共通 config、`IsTermTarget`）
- `dictionaries/Lucien.esp_Export.json`（実機検証 input）
- `dictionaries/Dawnguard.esm_Export.json`（参考 input）
- `extractData.pas`（xEdit 抽出側、出力フォーマットの一次定義）

## 分岐元

- 作業 branch: 未確定（preparation-module 入口で確定する）
- 分岐元 branch: `master`
- 分岐元 commit: `ef2a9f3411990de421d588189d916a49003deb27`

## 後続モジュール引き継ぎ

- task-id: `translation-input-import-non-dialogue`
- 依頼要約: 上記
- 想定 Y/N、artifact 集合、設計成果物、人間確認は design-module 入口で扱う。
- 仕様の出口は「13 種別 REC が抽出 JSON から DB へ正しく格納される」だが、画面導線・取り込み UI を変えるか、既存「翻訳入力取り込み」フローに乗せるかは設計差分で確定する。

## 前提・確認済み事項

- 既存翻訳データは reset 前提で互換 migration 不要（前 task の人間判断を継承）。
- DOOR:FULL、FLOR:FULL、FURN:FULL は 13 種別外であり取り込み対象でもない。
- xEdit 抽出側の出力 schema（key 名、type 表現、配列名）は実機 Lucien/Dawnguard JSON を一次根拠とする。
