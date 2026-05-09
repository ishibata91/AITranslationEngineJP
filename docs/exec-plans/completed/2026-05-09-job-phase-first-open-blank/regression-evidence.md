# 回帰テスト証跡

## 判断結果

- 判定: 完了
- 対象成果物: `回帰テスト証跡`
- 対象シナリオ: 初回の `現在の翻訳段階へ進む` 操作後に、`job-run` route で `ジョブ #1` と `単語翻訳` UI が表示される。
- テスト種別: frontend component scenario test

## 変更ファイル

- `frontend/src/ui/views/AppShell.test.ts`

## 追加した証明

- 入力開始点: 未完了ジョブ一覧で `jobID1` の `現在の翻訳段階へ進む` button を押す。
- 主要観測点: `window.location.hash` が `#translation-management/job-run` になる。
- 主要観測点: `ジョブ #1` heading が表示される。
- 主要観測点: `単語翻訳` heading が表示される。
- 主要観測点: `未完了ジョブ一覧でジョブを選んでください` が表示されない。
- 再実行観点: `未完了一覧へ戻る` で一覧へ戻り、同じ button を押しても `ジョブ #1` と `単語翻訳` UI が表示される。

## テスト実装内容

- `AppShell.test.ts` に、未完了ジョブ一覧の job card 操作から `JobRunPage` 表示までを通すシナリオテストを追加した。
- テスト補助は、`TranslationJobManagementPage` が一覧 card の `jobRunTarget` を渡した直後、detail loading 相当の `jobRunTarget = null` を通知する fake controller に限定した。
- `JobRunPage` の描画に必要な単語翻訳段階 controller とペルソナ生成段階 controller は、表示に必要な最小 fake に限定した。
- プロダクトコード、UI 表示、画面文言、layout、style は変更していない。

## 検証結果

- `npm --prefix frontend run test -- AppShell.test.ts`: 成功
- `python3 scripts/harness/run.py --suite frontend-local`: 成功
- frontend lint harness: 成功
- frontend test harness: 成功
- Vitest: 58 files passed, 518 tests passed
- boundary test: 1 file passed, 23 tests passed
- 既知 warning: Vite が `optimizeDeps.esbuildOptions` の非推奨 warning を出した。今回のテスト変更とは無関係。

## 未証明小範囲

- 実ブラウザ上の screenshot と操作証跡は、この回帰テスト証跡では新規取得していない。
- 実ブラウザ確認は、既存の実装証跡に証跡がある。fix_lane が browser confirmation へ進む場合は、同じ入力開始点を実ブラウザで再確認する。

## 戻し先判断

- fix_lane は browser confirmation へ進める。
- 理由: frontend component scenario test と `frontend-local` が通過し、証明対象の初回操作、主要観測点、再実行観点をテストで固定した。
