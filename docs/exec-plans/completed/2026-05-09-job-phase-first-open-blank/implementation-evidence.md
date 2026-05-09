# 実装証跡

## 判断結果

- 判定: 完了
- 対象成果物: `実装証跡`
- 実装 agent: `frontend_implementer`
- 変更範囲: frontend presenter の detail loading 中 target 生成

## 変更ファイル一覧

- `frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.ts`
- `frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.test.ts`
- `frontend/src/ui/views/AppShell.test.ts`

## 実装内容

- `translation-job-management.presenter.ts` に `resolveJobRunTarget` を追加した。
- `selectedJobDetail` が存在する場合は、従来どおり detail から `jobRunTarget` を作る。
- `detailPhase` が `loading` で `selectedJobId` が存在する場合だけ、一覧 summary から `jobRunTarget` を作る。
- `detailPhase` が `stale`、`idle`、または選択 job が一覧に存在しない場合は、`jobRunTarget` を `null` のままにする。
- `AppShell.syncJobRunTarget` の `null` 握りつぶしは削除した。

## 検証結果

- `python3 scripts/harness/run.py --suite frontend-local`: 成功
- `npm --prefix frontend run test -- translation-job-management.presenter.test.ts AppShell.test.ts`: 成功
- frontend lint harness: 成功
- frontend test harness: 成功
- Vitest: 58 files passed, 517 tests passed
- boundary test: 1 file passed, 23 tests passed
- 既知 warning: Vite が `optimizeDeps.esbuildOptions` の非推奨 warning を出した。今回変更とは無関係。

## UI 証跡

- `agent-browser doctor --offline --quick`: 6 pass, 0 warn, 0 fail
- 初回 open: `agent-browser open http://127.0.0.1:34115/#translation-management`
- 初回操作前 screenshot: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/frontend-implementation/initial-management.png`
- 初回操作後 URL: `http://127.0.0.1:34115/#translation-management/job-run`
- 初回操作後 snapshot: `ジョブ #1` と `単語翻訳` を確認した。
- 初回操作後 screenshot: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/frontend-implementation/after-first-open.png`
- 一覧へ戻った後 screenshot: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/frontend-implementation/after-return-list.png`
- 再実行後 URL: `http://127.0.0.1:34115/#translation-management/job-run`
- 再実行後 snapshot: `ジョブ #1` と `単語翻訳` を確認した。
- 再実行後 screenshot: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/frontend-implementation/after-second-open.png`
- `agent-browser errors`: 初回操作後、再実行後とも出力なし。
- 追加確認 screenshot: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/null-source-fix-after-first-open.png`
- 追加確認: presenter 起点修正後も、初回操作後に `ジョブ #1` と `単語翻訳` UI を確認した。
- 追加確認: `agent-browser errors` は出力なし。

## UI 根拠確認結果

- 回帰確認観点: 初回操作で `#translation-management/job-run` へ進むことを確認した。
- 回帰確認観点: 初回操作後に `ジョブ #1` と `単語翻訳` UI が表示されることを確認した。
- 回帰確認観点: 初回操作後に `未完了ジョブ一覧でジョブを選んでください` が表示されないことを snapshot で確認した。
- 回帰確認観点: 一覧へ戻って再実行しても `ジョブ #1` と `単語翻訳` UI が表示されることを確認した。
- UI 表示、画面文言、layout、style は変更していない。

## 残留リスク

- `job-run` 直リンク時は既存実装どおり `#translation-management` へ戻す。今回の修正対象ではない。
- `GetJobDetail` が失敗して stale になった場合は、一覧 summary から target を復元しない。
