# docs 正本化結果: DOC-01

- `task_id`: `2026-05-25-phase-prompt-builder-boundary`
- `agent`: `docs_updater`
- `skill`: `updating-docs`
- `status`: completed
- `canonicalization_input`: `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/docs-canonicalization-input.DOC-01.md`
- `approval_record`: 人間が 2026-05-25 に `approve` と返信し、`detail-spec-diff.md` と `design-diff.phase-prompt-builder-boundary.md` を承認した。
- `return_to`: `implement_lane`

## 更新 docs

- `docs/detail-specs/term-translation-phase.md`
- `docs/detail-specs/persona-generation-phase.md`
- `docs/detail-specs/body-translation-phase.md`

## 反映内容

- 単語翻訳フェーズに、1 対象語を 1 実行単位にする生成指示と応答検査を反映した。
- 単語翻訳フェーズに、対象語、原文言語、訳文言語、応答対応識別子を同じ実行単位へ固定する仕様を反映した。
- 単語翻訳フェーズに、応答欠落、余分な応答、空訳語、対象語との不一致を対象語単位の失敗分類として反映した。
- NPC ペルソナ生成フェーズに、1 NPC を 1 実行単位にする生成指示と応答検査を反映した。
- NPC ペルソナ生成フェーズに、NPC 対応識別子、表示名、属性、原文発話、会話文脈、共通ペルソナ要約を同じ実行単位へ固定する仕様を反映した。
- NPC ペルソナ生成フェーズに、応答欠落、余分な応答、NPC 対応識別子との不一致、空のペルソナ本文を NPC 単位の失敗分類として反映した。
- 本文翻訳フェーズに、1 翻訳項目を 1 実行単位にする生成指示と応答検査を反映した。
- 本文翻訳フェーズに、翻訳項目識別子、レコード種別、フィールド種別、原文、ペルソナ要約、翻訳補助情報、辞書制約、保持要素を同じ実行単位へ固定する仕様を反映した。
- 本文翻訳フェーズに、応答欠落、余分な応答、翻訳項目との不一致、空訳文、保持要素不整合を翻訳項目単位の失敗分類として反映した。
- 3 フェーズに、生成指示全文、外部サービス生データ、秘密値、原文発話全文、会話文脈全文を利用者向け情報から分離する仕様を反映した。
- 3 フェーズに、`PromptDigest` を復元不能な内部同一性情報として扱う仕様を反映した。
- 3 フェーズに、`TERM_TRANSLATION_REQUEST_V1`、`PERSONA_GENERATION_REQUEST_V1`、`BODY_TRANSLATION_REQUEST_V1` を AIサービス要求形状の識別子として扱う仕様を反映した。

## 反映しなかった内容

- `PromptBuilder`、`PromptInput`、`provider adapter` などの実装部品名は、正本仕様本文へ昇格しなかった。
- implementation-scope の実装手順、handoff、検証手順は、正本へ昇格しなかった。
- 画面設計、Storybook、frontend gateway、Wails DTO の意味拡張は、今回の正本化対象外として扱った。
- プロダクトコード、プロダクトテスト、`.codex/`、skill、agent、workflow 契約は変更しなかった。

## 検証結果

- `python3 scripts/harness/run.py --suite structure`: pass

## 残留不足

- なし。

## 戻し先

- 残留不足がある場合の戻し先は `implement_lane` とする。
