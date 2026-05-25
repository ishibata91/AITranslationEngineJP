# マージ準備入力

- `task_id`: `2026-05-25-phase-prompt-builder-boundary`
- `handoff_to`: `merge_lane`
- `active_plan_folder`: `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/`
- `worktree_path`: `/Users/iorishibata/.codex/worktrees/991e/AITranslationEngineJP`
- `source_branch`: `codex/2026-05-25-phase-prompt-builder-boundary`
- `target_branch`: `master`
- `work_commit_hash`: source branch HEAD at merge lane handoff

## マージ準備確認

- source branch は作業場所に checkout 済みである。
- target branch は local `master` として扱う。
- remote repository を変更する command は実行していない。
- completed 移動は実行していない。
- local merge は実行していない。

## 実装結果

- 3 フェーズ共通の prompt 受け渡し単位として `PromptEnvelope` を追加した。
- `PromptDigest` は生成指示全文を復元できない内部同一性情報として扱う。
- `TERM_TRANSLATION_REQUEST_V1`、`PERSONA_GENERATION_REQUEST_V1`、`BODY_TRANSLATION_REQUEST_V1` は AIサービス要求形状識別子として扱う。
- 単語翻訳、NPC ペルソナ生成、本文翻訳の各フェーズは専用の prompt 生成境界から provider adapter へ接続する。
- provider adapter は AIサービス接続差異の吸収を担当し、生成指示本文の組み立てを各フェーズの builder 側へ寄せた。
- 利用者向け情報と運用要約は raw prompt ではなく、件数、識別子、digest、失敗分類、安全な要約を扱う。

## docs 正本化結果

- `docs_updater/DOC-01` が完了した。
- 更新済み正本: `docs/detail-specs/term-translation-phase.md`
- 更新済み正本: `docs/detail-specs/persona-generation-phase.md`
- 更新済み正本: `docs/detail-specs/body-translation-phase.md`
- 結果記録: `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/docs-canonicalization-result.DOC-01.md`

## 検証結果

- pass: `env GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./internal/service ./internal/usecase ./internal/controller/wails`
- pass: `git diff --check`
- pass: `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite backend-local`
- pass: `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite structure`
- pass: `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite coverage`
- pass: `GOCACHE=/private/tmp/aitranslationenginejp-go-build go test ./scripts/test/seed-system-test-db`
- pass: sandbox 外 `python3 scripts/harness/run.py --suite system-test`
- pass: sandbox 外 `python3 scripts/harness/run.py --suite all`
- coverage: Sonar coverage `71.1%`, line `73.0%`, branch `57.1%`
- Sonar issues: security `0`, reliability `0`, maintainability HIGH `0`
- fail: sandbox 内 `env GOCACHE=/private/tmp/aitranslationenginejp-go-build python3 scripts/harness/run.py --suite system-test`
- fail reason: Wails dev server が `http://127.0.0.1:34115` で ready にならなかった。
- system-test result: sandbox 外で 10 件中 10 件 pass。
- suite all result: sandbox 外で structure、execution、system-test、coverage が pass。

## 実装後ブラウザ確認

- 専用の `browser_confirmation` は未実行である。
- 理由は、承認済み実装範囲が backend service 層の prompt 生成境界であり、確認すべき画面操作または URL を持たないためである。
- Wails と browser 起動を伴う確認は sandbox 外 system-test で実行し、10 件中 10 件 pass した。

## 残留リスク

- sandbox 内 system-test は Wails dev server ready timeout で失敗する。
- sandbox 外 `suite all` は通過しているため、残留リスクは sandbox 実行環境に限定される。
- `frontend/dist` と `frontend/wailsjs` は Wails が生成した ignored 検証生成物である。
- `scripts/harness/__pycache__/harness_common.cpython-314.pyc` は検証副作用で変更された tracked file であり、作業 commit には含めない。

## merge lane への依頼

- `work_commit_hash` と source branch の対応を確認する。
- active plan を completed へ移動する。
- local merge を実施する。
- merge 後検証を実施する。
- merge 結果 commit を作成する。
