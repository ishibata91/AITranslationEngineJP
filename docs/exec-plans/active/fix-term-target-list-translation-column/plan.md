# fix-term-target-list-translation-column

## 依頼要約

単語翻訳フェーズの処理対象一覧で、列見出しが「原語」「訳語」なのに「訳語」列へレコード種別フィールドが表示されている。
「訳語」列には訳語を表示し、登録時に訳語が存在しない場合は空表示にする。
新しく「レコード種別」列を追加し、`NPC_:FULL` のような `レコードタイプ:レコードフィールド` を表示する。

## 分岐元

- branch: `master`
- commit: `5d2fb6b97c25763613fca5632ac7bdbdeec9620a`

## 想定 Y/N 評価

- 仕様変更または仕様追加がある: N。`docs/detail-specs/term-translation-phase.md` は訳語と REC を別概念として扱っている。
- 画面変更がある: Y。処理対象一覧の列構造を変更する。
- 内部構造変更がある: 可能性がある。処理対象一覧 item の metadata 生成または表示変換を確認する必要がある。
- 画面の表示変更がある: Y。「レコード種別」列を追加し、「訳語」列の値を変える。
- frontend ロジック変更がある: 可能性がある。metadata を列へ変換する処理が frontend にある場合は変更対象になる。
- backend 変更がある: 可能性がある。処理対象一覧 item の metadata を backend が生成している場合は変更対象になる。
- frontend と backend を接続する: 可能性がある。既存 DTO の範囲で足りるかを確認する。
- 実装済み責務を独立に証明したい: Y。表示変換または read model の回帰をテストで証明する。
- 実行時にしか確定しない値または原因分離が要る分岐がある: N。対象値は処理対象一覧の read model と fixture で確認できる。

## 人間観測記録

- 確認済み問題: 単語翻訳フェーズの処理対象一覧で、「訳語」列にレコード種別フィールドが表示されている。
- 期待状態: 「訳語」列は訳語を表示し、登録直後など訳語がない場合は空表示になる。
- 期待状態: 「レコード種別」列を追加し、`NPC_:FULL` のような `レコードタイプ:レコードフィールド` を表示する。

## 確定原因

- backend の単語翻訳処理対象 SQL が `title_part_2` に REC、`title_part_3` に訳語を入れていた。
- backend の title part 変換処理が空文字を削除していたため、訳語が空の場合に REC が 2 列目へ詰められていた。
- frontend の処理対象一覧 component が title part を最大 2 列として扱っていたため、3 列目のレコード種別を表示できなかった。

## 修正内容

- `internal/repository/processing_target_sqlite_repository.go`: 単語翻訳の title part を原語、訳語、レコード種別の順に変更し、中間の空 title part を保持する。
- `frontend/src/ui/components/ProcessingTargetListPanel.svelte`: title part の列数を可変にし、空 title part は表示セルとして保持する。
- `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte`: 単語翻訳の列名を「原語」「訳語」「レコード種別」に変更する。
- `frontend/src/ui/screens/job-run/__fixtures__/job-run-shell-fixtures.ts`: 単語翻訳 fixture にレコード種別列を追加する。
- `tests/system/support/scenario-wails-mocks.ts`: system test mock の単語翻訳 title part にレコード種別を追加し、Lucien 登録直後の訳語を空にする。
- `frontend/src/ui/screens/job-run/ProcessingTargetListPanel.test.ts`: 空訳語でもレコード種別が訳語列へ詰められないことを追加する。
- `internal/repository/processing_target_sqlite_repository_test.go`: backend read model が原語、訳語、レコード種別の順で title part を返すことを追加する。
- `tests/system/job-run-shell.spec.ts`: Lucien データロード経路でレコード種別列と空訳語を確認する assertion を追加する。
- `docs/exec-plans/active/fix-term-target-list-translation-column/test-design.csv`: system test 変更に対応する task-local E2E 観点を固定する。
- `docs/e2e-test-design/test-design.csv`: 承認済み task-local E2E 観点を正本へ反映する。

## 検証結果

- `go test ./internal/repository -run 'TestSQLiteProcessingTargetListTermTranslationUsesAITargetTermPopulation|TestProcessingTargetTermSQL_excludesByDictionaryScopeRECORDFIELD|TestProcessingTargetTermSQL_separatesNPCFullAndSHRT|TestProcessingTargetTermSQL_filtersBy13RECTypes' -count=1`: pass。
- `npm --prefix frontend run test -- --run src/ui/screens/job-run/ProcessingTargetListPanel.test.ts`: pass。15 tests。
- `npm --prefix frontend run build-storybook`: pass。
- `python3 scripts/harness/run.py --suite frontend-local`: pass。frontend test は 556 tests。
- `python3 scripts/harness/run.py --suite backend-local`: pass。
- `go test ./internal/service -run Scenario_ProcessingTargetAndExecutionRECParity -count=1`: pass。
- `PLAYWRIGHT_BASE_URL=http://0.0.0.0:34115 npx playwright test tests/system/job-run-shell.spec.ts -g 'E2E-DIFF-LUCIEN-001'`: pass。
- `PLAYWRIGHT_BASE_URL=http://0.0.0.0:34115 npx playwright test tests/system/job-run-shell.spec.ts -g 'E2E-UC-045|E2E-DIFF-LUCIEN-001'`: `E2E-DIFF-LUCIEN-001` は pass。`E2E-UC-045` は列 assertion の前に「現在の翻訳段階へ進む」ボタンが見つからず fail。失敗 snapshot では対象 job card の操作名が「再開」であり、今回の列表示変更とは別の既存入口 selector 差分として扱う。
- `chrome-devtools` 実画面確認: `http://localhost:34115/#translation-management/job-run` で単語翻訳の処理対象一覧に「原語」「訳語」「レコード種別」列が表示されることを確認した。現在の実 app データでは処理対象が 0 件だったため、実 app の行セルは未確認である。

## E2E テスト観点差分

- 分類: 追加候補あり。
- 成果物: `docs/exec-plans/active/fix-term-target-list-translation-column/test-design.csv`
- 判断: 今回の system test 変更は、処理対象一覧で「訳語」列へレコード種別が詰まらないことを証明する回帰観点である。人間が正本反映を承認したため、`docs/e2e-test-design/test-design.csv` へ `E2E-UC-051` と `E2E-UC-052` を追加する。
- 関連する system test: `tests/system/job-run-shell.spec.ts` の `E2E-DIFF-LUCIEN-001`、`E2E-UC-045`。

## 正本化判断

- 仕様変更または仕様追加の対象: 詳細仕様正本は対象外である。UI 人間操作 E2E の正本観点は対象である。
- 対象 docs パス候補: `docs/e2e-test-design/test-design.csv`。
- 判断結果: 詳細仕様正本反映は不要。E2E テスト観点正本へ、単語翻訳の処理対象一覧で「原語」「訳語」「レコード種別」を確認する 2 観点を反映する。
- 人間承認状態: 承認済み。承認記録はユーザー指示「$finalization-module 正本反映」である。

## 作業 commit

- source branch: `claude/fix-term-target-list-translation-column`
- target branch: `master`
- active plan folder: `docs/exec-plans/active/fix-term-target-list-translation-column/`
- 実装修正 commit hash: `85f29e4f64c0411b86ce8e117d1362aa0c8f9fc1`
- E2E 観点正本反映 commit hash: `5c710fea6cefad3334a61b31055dca47340bddc2`
- 変更ファイル一覧:
  - `frontend/src/ui/components/ProcessingTargetListPanel.svelte`
  - `frontend/src/ui/screens/job-run/ProcessingTargetListPanel.test.ts`
  - `frontend/src/ui/screens/job-run/__fixtures__/job-run-shell-fixtures.ts`
  - `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte`
  - `internal/repository/processing_target_sqlite_repository.go`
  - `internal/repository/processing_target_sqlite_repository_test.go`
  - `tests/system/job-run-shell.spec.ts`
  - `tests/system/support/scenario-wails-mocks.ts`
  - `docs/exec-plans/active/fix-term-target-list-translation-column/test-design.csv`
  - `docs/e2e-test-design/test-design.csv`
  - `docs/exec-plans/active/fix-term-target-list-translation-column/plan.md`
- 検証結果: 上記「検証結果」を参照する。
- 残留リスク: `E2E-UC-045` は既存 job card 操作名差分で停止したため、今回追加した列 assertion までは到達していない。Lucien データロード経路の system test では対象列と空訳語を確認済みである。system test 変更に対応する task-local `test-design.csv` と正本 `docs/e2e-test-design/test-design.csv` は追加済みである。

## マージ準備入力

- active plan folder: `docs/exec-plans/active/fix-term-target-list-translation-column/`
- source branch: `claude/fix-term-target-list-translation-column`
- target branch: `master`
- 作業 commit hash: `85f29e4f64c0411b86ce8e117d1362aa0c8f9fc1`
- E2E 観点正本反映 commit hash: `5c710fea6cefad3334a61b31055dca47340bddc2`
- 最終検証結果: `frontend-local` pass、`backend-local` pass、`build-storybook` pass、`E2E-DIFF-LUCIEN-001` pass、`chrome-devtools` で列見出し確認済み。
- 残留リスク: `E2E-UC-045` の入口操作名差分は別件として残る。task-local E2E 観点と正本 E2E 観点は追加済みである。
