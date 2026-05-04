# reviewfix behavior-001

## 状態

- `fix_status`: 完了
- `reviewback`: [reviewback.behavior.yaml](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/reviewback.behavior.yaml)
- `issue_id`: `behavior-001`

## 修正内容

- [RunStatusPanel.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/RunStatusPanel.svelte)
- `中止` ボタンから旧 `aria-label="停止"` を削除した。
- [PersonaReviewPanel.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte)
- `編集` ボタンから旧 `aria-label="更新"` を削除した。
- [App.test.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/App.test.ts)
- master-persona 画面の期待 accessible name を `編集` と `中止` に更新した。

## 検証

- `git diff --check -- frontend/src/ui/screens/master-persona/RunStatusPanel.svelte frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte frontend/src/ui/App.test.ts`: 成功
- `python3 scripts/harness/run.py --suite frontend-local`: 成功
- frontend lint: 成功
- frontend test: `42 passed (files) / 421 passed (tests)`

## 残留確認

- master-persona 画面配下に `aria-label="更新"` と `aria-label="停止"` は残っていない。
- `App.test.ts:1881` の `更新` は別画面の既存テストであり、今回の master-persona UX 修正対象外。
