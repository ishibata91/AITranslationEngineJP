# docs 正本化起動入力: DOC-01

- `task_id`: `2026-05-25-phase-prompt-builder-boundary`
- `requester_role`: `implement_lane`
- `target_agent`: `docs_updater`
- `target_skill`: `updating-docs`
- `status`: ready

## 正本化判断

- `docs_canonicalization_required`: yes
- `reason`: 人間承認済みの詳細仕様差分が、単語翻訳、NPC ペルソナ生成、本文翻訳の恒久仕様を変更しているため。
- `approval_record`: 人間が 2026-05-25 に `approve` と返信し、`detail-spec-diff.md` と `design-diff.phase-prompt-builder-boundary.md` を承認した。
- `screen_design_canonicalization_required`: no
- `screen_design_reason`: 画面、Svelte component、Storybook、人間操作導線は今回の対象外である。

## 承認済み根拠

- `detail_spec_diff`: `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/detail-spec-diff.md`
- `implementation_scope`: `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/implementation-scope.md`
- `implementation_results`: `backend-implementation-result.BE-01.md`, `backend-implementation-result.BE-02.md`, `backend-implementation-result.BE-03.md`, `backend-implementation-result.BE-04.md`
- `test_results`: `unit-test-result.TU-01.md`, `scenario-test-result.SCN-01.md`
- `observability_result`: `observability-result.OBS-01.md`
- `validation_fix_result`: `backend-validation-fix-result.VAL-BE-01.md`

## 更新してよい正本

- `docs/detail-specs/term-translation-phase.md`
- `docs/detail-specs/persona-generation-phase.md`
- `docs/detail-specs/body-translation-phase.md`

## 反映する恒久仕様

- 単語翻訳フェーズは、1 対象語を 1 実行単位として AIサービスへ渡す生成指示を扱う。
- 単語翻訳フェーズは、対象語、原文言語、訳文言語、応答対応識別子を同じ実行単位に固定する。
- 単語翻訳フェーズは、応答欠落、余分な応答、空訳語、対象語との不一致を対象語単位の失敗分類として扱う。
- NPC ペルソナ生成フェーズは、1 NPC を 1 実行単位として AIサービスへ渡す生成指示を扱う。
- NPC ペルソナ生成フェーズは、NPC 対応識別子、表示名、属性、原文発話、会話文脈、共通ペルソナ要約を同じ実行単位に固定する。
- NPC ペルソナ生成フェーズは、応答欠落、余分な応答、NPC 対応識別子との不一致、空のペルソナ本文を NPC 単位の失敗分類として扱う。
- 本文翻訳フェーズは、1 翻訳項目を 1 実行単位として AIサービスへ渡す生成指示を扱う。
- 本文翻訳フェーズは、翻訳項目識別子、レコード種別、フィールド種別、原文、ペルソナ要約、翻訳補助情報、辞書制約、保持要素を同じ実行単位に固定する。
- 本文翻訳フェーズは、応答欠落、余分な応答、翻訳項目との不一致、空訳文、保持要素不整合を翻訳項目単位の失敗分類として扱う。
- 3 フェーズは、生成指示全文、外部サービス生データ、秘密値、原文発話全文、会話文脈全文を利用者向け情報から分離する。
- `PromptDigest` は生成指示の同一性を示す復元不能な内部情報として扱う。
- `TERM_TRANSLATION_REQUEST_V1`、`PERSONA_GENERATION_REQUEST_V1`、`BODY_TRANSLATION_REQUEST_V1` は AIサービス要求形状の識別子として扱い、利用者が選ぶ生成規則版として扱わない。

## 反映しない事項

- `PromptBuilder`、`PromptInput`、`provider adapter` などの実装部品名は、正本の仕様本文へ自動昇格しない。
- implementation-scope の実装手順、handoff、検証手順は正本へ昇格しない。
- 画面設計、Storybook、frontend gateway、Wails DTO の意味拡張は扱わない。
- `.codex/`、skill、agent、workflow 契約は変更しない。
- プロダクトコードとプロダクトテストは変更しない。

## 検証証跡

- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/service ./internal/usecase ./internal/controller/wails`: pass
- `git diff --check`: pass
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite backend-local`: pass after Wails generated `frontend/dist`.
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite structure`: pass
- `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite coverage`: pass after Wails generated `frontend/wailsjs` and `frontend/dist`.
- `python3 scripts/harness/run.py --suite system-test`: fail after sandbox 外再実行。10 件中 7 件 pass、翻訳ジョブ管理 3 件 fail。失敗画面は未完了ジョブ 0 件を示す。

## 期待する出力

- 更新した docs 正本の一覧。
- 反映した恒久仕様の要約。
- 反映しなかった内容と理由。
- 実行した検証コマンドと結果。
- 残留不足があれば、戻し先を `implement_lane` として明示する。
