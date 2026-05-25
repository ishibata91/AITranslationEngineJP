# 詳細仕様差分: 2026-05-25-phase-prompt-builder-boundary

- `skill`: detail-spec-design
- `status`: ready-for-human-review
- `source_plan`: `./plan.md`
- `detail_spec_target`: `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`
- `screen_design_diff`: `N/A`
- `component_diagram`: `./design-diff.phase-prompt-builder-boundary.md`

## 詳細仕様差分

### `term-translation-phase-REQ-002` AI 設定を開始時に再解決する

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/term-translation-phase.md`

親要件:
フェーズ開始と再試行は最新の AIサービス設定を参照し、秘密値を利用者向け情報から分離する。

仕様:
- 単語翻訳フェーズの利用者向け情報は、AIサービス、モデル、認証状態、実行方式、入力件数、出力件数、失敗分類、処理要約を対象にする。
- AIサービスへ渡す生成指示の全文は、障害調査用の同一性情報へ要約して扱う。
- `PromptDigest` は、AIサービスへ渡す生成指示の同一性を示す内部情報として扱う。
- `PromptDigest` は、生成指示の全文を復元できない値として扱う。
- `TERM_TRANSLATION_REQUEST_V1` は、単語翻訳フェーズの AIサービス要求形状を識別する印として扱う。
- `TERM_TRANSLATION_REQUEST_V1` は、利用者が選択する生成規則の版として扱わない。

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
- AIサービスへ渡す単語翻訳の生成指示は、1 対象語を 1 実行単位として扱う。
- 生成指示は、対象語、原文言語、訳文言語、応答対応に使う識別子を同じ実行単位に固定する。
- AIサービス応答は、対象語ごとに、原語と訳語の対応として検査する。
- 有効な応答は、要求した対象語と同じ原語を持ち、空ではない訳語を持つ応答である。
- 有効な応答は、確定訳語として対象の翻訳ジョブ内辞書へ反映する。
- 応答欠落、余分な応答、空訳語、対象語との不一致は、対象語単位の失敗分類として扱う。

未決:
- なし

回答:
- なし

### `persona-generation-phase-REQ-002` 生成対象と入力条件を固定する

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/persona-generation-phase.md`

親要件:
フェーズは NPC の原文発話、NPC 属性メタデータ、会話文脈、共通ペルソナ参照からペルソナ参照情報を作る。

仕様:
- NPC ペルソナ生成フェーズは、1 NPC を 1 実行単位として扱う。
- AIサービスへ渡すペルソナ生成の生成指示は、NPC の対応識別子、NPC 表示名、NPC 属性、原文発話、会話文脈、共通ペルソナ要約を同じ実行単位に固定する。
- 生成対象は、NPC 件数、対象件数、共通ペルソナ一致有無、対象外理由を判断できる状態にする。
- 生成対象 0 件は `Completed` とし、AIサービス未実行と空のペルソナ参照を判断できる状態にする。

未決:
- なし

回答:
- なし

### `persona-generation-phase-REQ-004` AIサービス結果をペルソナ参照へ採用する

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/persona-generation-phase.md`

親要件:
有効な AIサービス出力は本文翻訳フェーズが参照するペルソナ情報として採用される。

仕様:
- AIサービス応答は、NPC ごとに、要求単位、NPC 対応識別子、ペルソナ本文の対応として検査する。
- 有効な応答は、要求した NPC と同じ対応識別子を持ち、空ではないペルソナ本文を持つ応答である。
- 有効な応答は、翻訳ジョブ内ペルソナまたはペルソナ参照として採用する。
- 応答欠落、余分な応答、NPC 対応識別子との不一致、空のペルソナ本文は、NPC 単位の失敗分類として扱う。
- 利用者は、ペルソナ参照の成立状態、欠落件数、本文翻訳フェーズの開始可否を判断できる。

未決:
- なし

回答:
- なし

### `persona-generation-phase-REQ-008` 保護対象を利用者向け情報から分離する

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/persona-generation-phase.md`

親要件:
NPC ペルソナ生成フェーズは秘密値、生成指示の原文、原文発話全文、会話文脈全文を利用者向け情報から分離する。

仕様:
- NPC ペルソナ生成フェーズの利用者向け情報は、AIサービス、モデル、実行方式、認証状態、入力件数、出力件数、失敗分類、保護対象を含まない結果要約を対象にする。
- AIサービスへ渡す生成指示の全文、原文発話全文、会話文脈全文は、障害調査用の同一性情報へ要約して扱う。
- `PromptDigest` は、AIサービスへ渡す生成指示の同一性を示す内部情報として扱う。
- `PromptDigest` は、生成指示の全文、原文発話全文、会話文脈全文を復元できない値として扱う。
- `PERSONA_GENERATION_REQUEST_V1` は、NPC ペルソナ生成フェーズの AIサービス要求形状を識別する印として扱う。
- `PERSONA_GENERATION_REQUEST_V1` は、利用者が選択する生成規則の版として扱わない。

未決:
- なし

回答:
- なし

### `body-translation-phase-REQ-002` 翻訳入力を固定する

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/body-translation-phase.md`

親要件:
本文翻訳フェーズは確定訳語、ジョブ内辞書、ジョブ内ペルソナ、翻訳補助情報を同一実行単位の入力として固定する。

仕様:
- 本文翻訳フェーズは、辞書置換対象外の翻訳項目を AI 翻訳対象として扱う。
- AIサービスへ渡す本文翻訳の生成指示は、1 翻訳項目を 1 実行単位として扱う。
- 生成指示は、翻訳項目の対応識別子、レコード種別、フィールド種別、原文、ペルソナ要約、翻訳補助情報、完全一致の辞書置換結果、部分一致の訳語固定制約、保持要素を同じ実行単位に固定する。
- 生成指示は、翻訳項目の対応関係と保持要素の同一性を AIサービス境界で判断できる状態にする。
- `PromptDigest` は、AIサービスへ渡す生成指示の同一性を示す内部情報として扱う。

未決:
- なし

回答:
- なし

### `body-translation-phase-REQ-003` 翻訳結果を翻訳項目単位で保存できる

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/body-translation-phase.md`

親要件:
本文翻訳フェーズは訳文、出力状態、保護要素検証結果を翻訳項目単位で保持する。

仕様:
- AIサービス応答は、翻訳項目ごとに、要求単位、翻訳項目の対応識別子、レコード種別、フィールド種別、訳文、保持要素の同一性として検査する。
- 有効な応答は、要求した翻訳項目と同じ対応識別子を持ち、空ではない訳文を持ち、保持要素の同一性を満たす応答である。
- 有効な応答は、翻訳項目単位の訳文と出力状態の候補として採用する。
- 保持要素検証に失敗した訳文は、翻訳項目単位の失敗分類として扱う。
- 応答欠落、余分な応答、翻訳項目との不一致、空訳文、保持要素の欠落、変更、重複、順序変更、余分な追加は、翻訳項目単位の失敗分類として扱う。

未決:
- なし

回答:
- なし

### `body-translation-phase-REQ-008` 秘密値と生データを利用者向け情報から分離する

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/body-translation-phase.md`

親要件:
本文翻訳フェーズは秘密値、外部サービスとの生データ、生成指示の原文を利用者向け情報から分離する。

仕様:
- 本文翻訳フェーズの利用者向け情報は、AIサービス、モデル、実行方式、認証状態、入力件数、出力件数、失敗分類、保持要素の件数、保持要素検証結果、保護対象を含まない結果要約を対象にする。
- AIサービスへ渡す生成指示の全文と外部サービスとの生データは、障害調査用の同一性情報へ要約して扱う。
- `PromptDigest` は、AIサービスへ渡す生成指示の同一性を示す内部情報として扱う。
- `PromptDigest` は、生成指示の全文または外部サービスとの生データを復元できない値として扱う。
- `BODY_TRANSLATION_REQUEST_V1` は、本文翻訳フェーズの AIサービス要求形状を識別する印として扱う。
- `BODY_TRANSLATION_REQUEST_V1` は、利用者が選択する生成規則の版として扱わない。

未決:
- なし

回答:
- なし

## 実装範囲で確認する事項

- `PromptBuilder`、`PromptInput`、`provider adapter`、`response parser`、`validation`、`adoption` は、上記仕様を満たすための実装判断材料として扱う。
- 実装範囲では、生成指示を作る責務、AIサービス接続差異を吸収する責務、応答を候補へ解釈する責務、有効性を検査する責務、有効結果を成果へ採用する責務を分けて確認する。
- 実装範囲では、`PromptDigest` と `REQUEST_V1` 系の値が利用者向け情報の中心にならず、内部同一性情報と要求形状識別子に留まることを確認する。

## 図成果物との整合

- `design-diff.phase-prompt-builder-boundary.md` は実装構造を示す補助図であり、詳細仕様差分の正本ではない。
- 詳細仕様差分は、入力固定、生成指示の意味、応答検査単位、有効結果の採用成果、公開情報と内部情報の境界を正本化対象にする。
- 補助図は実装判断材料としては利用できるが、人間レビューで詳細仕様差分と同じ粒度へ寄せる必要がある場合は、diagrammer へ戻す。

## 根拠

- `source`: `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/plan.md`
- `source`: `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/design-diff.phase-prompt-builder-boundary.md`
- `source`: `docs/detail-specs/term-translation-phase.md`
- `source`: `docs/detail-specs/persona-generation-phase.md`
- `source`: `docs/detail-specs/body-translation-phase.md`
- `source`: `docs/spec.md`
- `review`: 人間設計レビュー差し戻し。指摘内容は「実装用の判断軸としてはありだが、仕様差分ではなく実装計画に寄っている」。
- `validation`: `rg -n "PromptBuilder|PromptInput|provider adapter|response parser|validation|adoption" detail-spec-diff.md` で、実装部品名が仕様本文ではなく `実装範囲で確認する事項` にあることを確認する。
