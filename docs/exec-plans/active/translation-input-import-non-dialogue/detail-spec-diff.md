# 詳細仕様差分: translation-input-import-non-dialogue

- `skill`: detail-spec-design
- `status`: ready-for-human-review
- `source_plan`: `./plan.md`
- `detail_spec_target`: `docs/detail-specs/translation-input-intake.md`, `docs/detail-specs/term-translation-phase.md`
- `screen_design_diff`: `N/A`
- `component_diagram`: `./design-diff-diagram.md`

## 詳細仕様差分

### `translation-input-intake-REQ-002` xEdit 抽出 JSON の翻訳レコードと翻訳フィールドを展開できる

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/translation-input-intake.md`

親要件:
利用者は登録済み入力データの規模、分類、代表例、失敗理由を判断できる。

現行仕様:
- 登録した入力データから翻訳レコードと翻訳フィールドを展開し、件数とカテゴリを判断できる。
- 未定義 RecordType と SubrecordType の組み合わせは警告として扱い、非翻訳対象として観測可能に保持する。
- 代表例は RecordType、SubrecordType、FormID、EditorID、原文を根拠として示せる。

追加後の仕様:
- xEdit 抽出 JSON の登録結果は、`dialogue_groups` の本体と応答に加えて、実機 JSON が持つ会話以外の翻訳レコードも展開対象にする。
- 会話以外の翻訳レコードは、`items`、`magic`、`locations`、`cells`、`system`、`messages`、`load_screens`、`npcs`、`quests` に含まれる要素として扱う。
- `quests` に含まれる `stages` と `objectives` は、親のクエストと区別できる翻訳フィールドとして扱う。
- 実機 JSON 根拠では、最上位の `responses`、`stages`、`objectives` は確認できない。応答は `dialogue_groups` の子要素として扱い、ステージと目的は `quests` の子要素として扱う。
- 各翻訳フィールドの REC は、JSON 内の `type` が示すレコード種別とフィールド名から、`RECORD:FIELD` 形式で判断できる状態にする。
- `type` が `BOOK FULL` の場合は `BOOK:FULL` として扱い、`type` が `NPC_ FULL` の場合は `NPC_:FULL` として扱う。
- 翻訳フィールドの原文は、対象要素が持つ本文を根拠にする。
- 原文が空の翻訳フィールドは、単語翻訳フェーズの候補にしない状態で登録結果を判断できる。
- 会話以外の要素であっても、RecordType と SubrecordType の組み合わせが未定義の場合は既存の警告規則で扱う。
- `target_plugin` があり、取り込み可能な翻訳レコードが 1 件以上ある JSON は、`dialogue_groups` だけを必須とせず登録結果を作れる。

未決:
- なし

回答:
- なし

### `term-translation-phase-REQ-008` 単語翻訳対象 REC を 13 種別に固定する

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/term-translation-phase.md`

親要件:
利用者は、単語翻訳フェーズが処理対象とする REC の範囲を判断できる。

現行仕様:
- 単語翻訳対象 REC は、`BOOK:FULL`、`NPC_:FULL`、`NPC_:SHRT`、`ARMO:FULL`、`WEAP:FULL`、`LCTN:FULL`、`CELL:FULL`、`CONT:FULL`、`MISC:FULL`、`ALCH:FULL`、`RACE:FULL`、`INGR:FULL`、`SHOU:FULL` の 13 種別とする。
- 翻訳入力に出現する翻訳レコードまたは翻訳フィールドのうち、原語が空でなく、REC が 13 種別のいずれかに該当する場合だけ、単語翻訳フェーズの候補にする。
- 単語翻訳フェーズと XML 辞書取り込みは、同一の単語翻訳対象 REC 判定を共有して対象判定を行う。

追加後の仕様:
- 単語翻訳フェーズの候補判定は、翻訳入力の由来が会話配列か会話以外の配列かで変わらない。
- xEdit 抽出 JSON の会話以外の配列から取り込まれた翻訳フィールドも、REC が 13 種別のいずれかであり、原語が空でない場合は単語翻訳フェーズの候補にする。
- `NPC_:FULL` と `NPC_:SHRT` は、会話以外の配列から取り込まれた場合も別 REC として扱う。
- 13 種別に該当しない REC は、会話以外の配列から取り込まれた場合も単語翻訳フェーズの候補外として扱う。
- 単語翻訳対象 REC 集合は、XML 辞書取り込み対象 REC 集合と同一の 13 種別を維持する。

未決:
- なし

回答:
- なし

## 追加要件

- なし。今回の差分は、既存の `translation-input-intake-REQ-002` と `term-translation-phase-REQ-008` を満たすための詳細化として扱う。

## 実装へ渡す仕様境界

- 画面変更はない。表示 layout、文言、style、表示構造、Wails DTO の変更は本差分の範囲に含めない。
- 取り込み対象の外部形状は、実機 JSON の `target_plugin` と最上位 key を根拠にする。
- 最上位 key は、`cells`、`dialogue_groups`、`items`、`load_screens`、`locations`、`magic`、`messages`、`npcs`、`quests`、`system`、`target_plugin` を確認済みの形状として扱う。
- `npcs` は配列ではなく、FormID を key にした object として扱う。
- `cells` は Lucien と Dawnguard の両方で空配列だったが、取り込み可能な要素が出現した場合は会話以外の翻訳レコードとして扱う。
- `quests` の `stages` と `objectives` は、最上位配列ではなく `quests` の子要素として扱う。
- REC は `RECORD:FIELD` 形式で扱う。JSON 内の `type` が `RECORD FIELD` 形式である場合は、空白で分かれる前半を RECORD、後半を FIELD として扱う。
- 取り込み処理は、会話以外の翻訳レコードを登録結果の件数、カテゴリ、代表例、警告に反映できる必要がある。
- 単語翻訳フェーズは、取り込み済み翻訳フィールドの REC と原語から候補を判断する。候補判定は、取り込み元の JSON 配列名に依存しない。

## 根拠

- `source`: `docs/exec-plans/active/translation-input-import-non-dialogue/plan.md`
- `source`: `docs/detail-specs/translation-input-intake.md`
- `source`: `docs/detail-specs/term-translation-phase.md`
- `source`: `internal/service/translation_input_import_service.go`
- `source`: `internal/recclassification/term_target.go`
- `source`: `dictionaries/Lucien.esp_Export.json`
- `source`: `dictionaries/Dawnguard.esm_Export.json`
- `review`: 人間レビュー前の task-local 差分として作成した。
- `validation`: 静的確認のみ実行した。`sed` で指定 docs と Go ファイルを確認し、`jq` で Lucien と Dawnguard の最上位 key、配列件数、REC 分布、`npcs` 形状、`quests` 子要素を確認した。
