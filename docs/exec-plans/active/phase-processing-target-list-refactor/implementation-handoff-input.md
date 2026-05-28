# Implementation Handoff Input: phase-processing-target-list-refactor

## 起動元

- role: `refactor_lane`
- skill: `refactor-lane`
- task folder: `docs/exec-plans/active/phase-processing-target-list-refactor/`
- approved scope: `docs/exec-plans/active/phase-processing-target-list-refactor/implementation-scope.md`
- approval: 人間が「この問題のテスト追加，UC差分，修正までを行うこと」と指示した。

## 共通入力

- `docs/exec-plans/active/phase-processing-target-list-refactor/plan.md`
- `docs/exec-plans/active/phase-processing-target-list-refactor/spec-drift-investigation.md`
- `docs/exec-plans/active/phase-processing-target-list-refactor/structure-quality-investigation.md`
- `docs/exec-plans/active/phase-processing-target-list-refactor/test-quality-investigation.md`
- `docs/exec-plans/active/phase-processing-target-list-refactor/implementation-scope.md`
- `docs/exec-plans/active/phase-processing-target-list-refactor/usecase-diff.md`
- `docs/e2e-test-design/test-design.csv`

## 共通禁止事項

- docs 正本文を更新しない。
- `.codex/` を更新しない。
- `translation_complete`、phase 以外の job lifecycle、実外部 API、実 secret、実利用者データへ広げない。
- 既存未 commit 差分を revert しない。
- `docs/exec-plans/active/e2e-test-design-maintenance/` を変更しない。

## 実装順序

1. `wave-1`: `frontend-count-subject`
2. `wave-2`: `frontend-search-subject`
3. `wave-3`: `backend-count-read-model`
4. `wave-4`: `backend-search-read-model`
5. `wave-5`: `integration-processing-target-seam`
6. `wave-6`: `unit-count-subject`, `unit-search-subject`, `scenario-page-object`, `scenario-fixture`
7. `wave-7`: `scenario-phase-list-search`

## Wave 1 起動入力

### `frontend-count-subject`

- agent: `frontend_implementer`
- skill: `implement-frontend`
- 目的: 3 フェーズ画面の件数主語を処理対象一覧 total と矛盾しない表示へそろえる。
- owned scope:
  - `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte`
  - `frontend/src/ui/screens/persona-generation-phase/PersonaGenerationPhasePanel.svelte`
  - `frontend/src/ui/screens/body-translation-phase/BodyTranslationPhasePanel.svelte`
  - phase presenter、phase screen type、phase view-model test。
- first action: 単語翻訳の `phaseMetrics` と `progressDetails` を、処理対象件数が AI 翻訳対象語件数を指す clause へ変更する。
- completion signal:
  - 単語翻訳の処理対象件数は、AI 翻訳対象語件数と同じ主語で表示される。
  - NPC ペルソナ生成の処理対象件数は、現行対象件数の主語を維持する。
  - 本文翻訳の処理対象件数は、AI 送信対象件数と同じ主語で表示される。
  - `ProcessingTargetListPanel.svelte` のページング計算本体は変更されていない。
- validation: `python3 scripts/harness/run.py --suite frontend-local`

## 完了報告

各 agent は次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `blocked_items`
- `residual_risks`
- `docs_changes: none`
