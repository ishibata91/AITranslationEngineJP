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

- 作業 branch: `claude/translation-input-import-non-dialogue`
- 分岐元 branch: `master`
- 分岐元 commit: `e83cccad85ae294f26e98f81a6de437b0e65f3b9`

## 後続モジュール引き継ぎ

- task-id: `translation-input-import-non-dialogue`
- 依頼要約: 上記
- 想定 Y/N、artifact 集合、設計成果物、人間確認は design-module 入口で扱う。
- 仕様の出口は「13 種別 REC が抽出 JSON から DB へ正しく格納される」だが、画面導線・取り込み UI を変えるか、既存「翻訳入力取り込み」フローに乗せるかは設計差分で確定する。

## 前提・確認済み事項

- 既存翻訳データは reset 前提で互換 migration 不要（前 task の人間判断を継承）。
- DOOR:FULL、FLOR:FULL、FURN:FULL は 13 種別外であり取り込み対象でもない。
- xEdit 抽出側の出力 schema（key 名、type 表現、配列名）は実機 Lucien/Dawnguard JSON を一次根拠とする。

## preparation-module 出口

- 作業 branch は `claude/translation-input-import-non-dialogue` として作成済み。
- active plan folder は `docs/exec-plans/active/translation-input-import-non-dialogue/` として存在する。
- `plan.md` は依頼要約、分岐元 branch、分岐元 commit を記録済み。

## design-module 想定 Y/N 評価

- 仕様変更または仕様追加がある: Y。`docs/detail-specs/translation-input-intake.md` は xEdit 抽出 JSON を翻訳レコードと翻訳フィールドへ展開する仕様を持つが、`internal/service/translation_input_import_service.go` は `dialogue_groups` 以外を parse しない。
- 画面変更がある: N。依頼は既存「翻訳入力取り込み」フローの取り込み対象拡張であり、表示 layout、文言、style、表示構造の変更要求はない。
- 内部構造変更がある: Y。`translationInputDocument` と prepare 処理に、`items`、`magic`、`locations`、`cells`、`system`、`messages`、`load_screens`、`npcs`、`quests` 配下の stage/objective を取り込む構造追加が必要である。
- 画面の表示変更がある: N。Svelte 表示コンポーネント、props、style、story、fixture は変更しない想定である。
- frontend ロジック変更がある: N。state、API、Wails bridge、ルーティング、副作用、フォーム validation の変更要求はない。
- backend 変更がある: Y。`internal/service/translation_input_import_service.go` の JSON parse と DB 永続化準備が変更対象である。
- frontend と backend を接続する: N。既存 Wails DTO と既存 import usecase の入出力形状は変更しない想定である。
- 実装済み責務を独立に証明したい: Y。取り込み parse と 13 種別対象 REC の永続化は単体テストで独立に証明する必要がある。
- 実行時にしか確定しない値または原因分離が要る分岐がある: N。実機 JSON の shape は fixture と単体テストで固定でき、恒久観測ログ追加は不要である。

## design-module 出口

- 詳細仕様差分: `docs/exec-plans/active/translation-input-import-non-dialogue/detail-spec-diff.md`
- 設計差分図: `docs/exec-plans/active/translation-input-import-non-dialogue/design-diff-diagram.md`
- 画面設計差分: 省略。画面変更がないため不要である。
- 人間設計レビュー: 承認済み。会話上の進行指示を本 task の設計レビュー通過として扱う。
- 実装範囲: `docs/exec-plans/active/translation-input-import-non-dialogue/implementation-scope.md`
- テスト設計: `docs/exec-plans/active/translation-input-import-non-dialogue/test-design.csv`
- Storybook モジュール: 省略。画面の表示変更がないため不要である。

## implementation-module 入力

- 実装対象: backend 実装と単体テストである。
- backend 実装対象: `internal/service/translation_input_import_service.go`
- 単体テスト対象: `internal/service/translation_input_import_service_test.go`
- final validation: backend 変更のため `python3 scripts/harness/run.py --suite backend-local` を実行する。
- frontend、Wails、UI、Storybook、DB migration、`extractData.pas` は実装範囲外である。

## implementation-module 出口

- backend 実装: 完了。`internal/service/translation_input_import_service.go` で `dialogue_groups` 以外の xEdit 抽出 JSON 要素を decode と prepare の対象に追加した。
- 単体テスト: 完了。`internal/service/translation_input_import_service_test.go` で non-dialogue record 取り込み、`dialogue_groups` 非依存、空 importable record 拒否を追加した。
- 観測ログ追加: 省略。想定 Y/N 評価で「実行時にしか確定しない値または原因分離が要る分岐」が N のため不要である。
- frontend、Wails、UI、Storybook、DB migration、`extractData.pas`: 変更なし。

## 最終検証

- `go test ./internal/service`: 通過。
- `python3 scripts/harness/run.py --suite backend-local`: 通過。
- 実画面確認: 通過。`dictionaries/Lucien.esp_Export.json` を登録し、単語翻訳画面で `AI 翻訳対象語 206 件` と `1-50 / 206 件` を確認した。

## finalization-module 正本化判断

- 仕様変更対象: `docs/detail-specs/translation-input-intake.md` と `docs/detail-specs/term-translation-phase.md` に反映候補がある。
- 判断結果: 保留。docs 正本本文への反映は、人間が恒久仕様として明示承認した後に行う。
- 人間承認状態: docs 正本反映は未承認である。
- 詳細仕様正本反映: 未実行。人間承認なしの docs 正本本文変更は禁止である。

## finalization-module 作業 commit

- commit 対象: active plan 成果物、`internal/service/translation_input_import_service.go`、`internal/service/translation_input_import_service_test.go`、`CLAUDE.md`、`.claude/settings.json`
- commit hash: `6ac9a7c`
- 検証結果: `go test ./internal/service` 通過、`python3 scripts/harness/run.py --suite backend-local` 通過、実画面確認通過。
- 残留リスク: docs 正本反映は未承認のため未実行である。

## finalization-module マージ準備入力

- active plan folder: `docs/exec-plans/active/translation-input-import-non-dialogue/`
- source branch: `claude/translation-input-import-non-dialogue`
- target branch: `master`
- 作業 commit hash: `6ac9a7c`
- 最終検証結果: `go test ./internal/service` 通過、`python3 scripts/harness/run.py --suite backend-local` 通過、実画面確認通過。
- 残留リスク: docs 正本反映は未承認のため未実行である。
