# scenario-test-addition-evidence

## 2026-05-30 ユーザー操作基準への修正

- 担当 role: `implementation_scenario_tester`
- 使用 skill: `.codex/skills/tests-scenario/SKILL.md`
- 呼び出し元: `fix_lane`
- 対象シナリオ: `E2E-DIFF-LUCIEN-001`

## 修正内容

- `tests/system/job-run-shell.spec.ts`
  - 開始点を未完了ジョブ一覧の直接選択から、`翻訳管理`、`新規翻訳を開始`、データロード画面での `dictionaries/Lucien.esp_Export.json` 選択、登録、`単語翻訳へ進む` へ変更した。
  - テスト本体は条件分岐を含まない。
  - 期待値は `4930`、`4931` の固定一致にしていない。
- `tests/system/support/scenario-wails-mocks.ts`
  - `Lucien.esp_Export.json` の未完了 job を直接 seed する補助を使わない形にした。
  - `ImportTranslationInput`、`CreateTranslationJobFromInput`、作成後 job の単語翻訳 summary、処理対象一覧を接続した。
  - 処理対象一覧応答は `items` と `totalCount` を使い、`listItems` / `listTotal` と `items` / `totalCount` の人工的な shape 揺れを使わない。

## 証明したユーザー操作

1. `翻訳管理` を開く。
2. `新規翻訳を開始` を押す。
3. データロード画面で `dictionaries/Lucien.esp_Export.json` を選択する。
4. `この JSON を登録` を押す。
5. 登録済み一覧で `Lucien.esp_Export.json` を確認する。
6. `単語翻訳へ進む` を押す。
7. 作成済み job の単語翻訳画面で進捗と処理対象一覧を確認する。

## 検証結果

- 実行コマンド: `PLAYWRIGHT_BASE_URL=http://0.0.0.0:34115 npx playwright test --config ./playwright.config.ts tests/system/job-run-shell.spec.ts -g "E2E-DIFF-LUCIEN-001"`
- 結果: 1 passed
- 観測点:
  - 進捗パネルは `AI 翻訳対象語件数` を非 0 件として表示した。
  - 処理対象一覧の件数は非 0 件を表示した。
  - `処理対象がありません` の空状態だけを表示しなかった。
  - 処理対象行が 1 件以上表示された。
- 補足: 初回実行は macOS sandbox の Chromium 起動制限で失敗したため、通常権限で再実行した。
- 補足: 2 回目は `http://0.0.0.0:34115` の開発サーバー未起動で失敗したため、`VITE_HOST=0.0.0.0 VITE_PORT=34115 npm --prefix frontend run dev:wails` を起動して再実行した。

## 局所検証

- 実行コマンド: `python3 scripts/harness/run.py --suite frontend-local`
- 結果: passed
- 内訳:
  - frontend lint harness passed
  - frontend test harness passed
  - Vitest: 53 files passed, 518 tests passed
