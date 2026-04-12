# Fix Plan

- workflow: fix
- status: completed
- lane_owner: orchestrating-fixes
- scope: master-dictionary xml import button and completion state

## Request Summary

- 辞書構築で XML を選択後、取り込み中表示になるが取込完了へ進まず、何も起きないように見える。
- 取込ボタン系にポインターカーソルが出ず、クリック可能性が伝わらない。

## Decision Basis

- active implementation plan `docs/exec-plans/active/2026-04-11-master-dictionary-management.md` では、XML 取込は file-selection-first と same-page refresh を成立させる要求になっている。
- `frontend/src/ui/screens/master-dictionary/master-dictionary.usecase.ts` では runtime event が購読済みの時、import API 応答だけでは完了状態へ遷移しない。
- `frontend/src/ui/screens/master-dictionary/MasterDictionaryPage.svelte` の button 系 style には `cursor: pointer` がない。

## Known Facts

- `python3 scripts/harness/run.py --suite structure` は通過した。
- `MasterDictionaryRuntimeEventAdapter.subscribe()` は runtime bridge がない時だけ `false` を返し、ある時は completed event 到達待ちになる。
- `startStagedXmlImport()` は `waitForRuntimeCompletion === true` の時、API 応答後に `handleImportCompleted()` を呼ばず return する。
- `frontend/src/ui/App.test.ts` には runtime completion event 到達まで完了へ遷移しない test がある。
- 現在の `stageXmlImport()` はファイル選択時に `importStage = "ready"` と `importProgress = 12` を設定する。
- Playwright MCP の DOM 注入経路では、`この XML を取り込む` 押下後も `取込待ち` のままで `取込中` と progress 変化が出ず、console error は `favicon.ico` 404 だけだった。
- 現行 UI では import button 系の `pointer` 表示は確認でき、前回の cursor 不足は再現しなかった。
- user 報告は、`XMLから取り込む` を押しても開始したように見えず、ファイル選択時点で progress/seek bar が途中まで進んで見える点に収束している。
- Wails browser 経路 `http://host.docker.internal:34115/#master-dictionary` でも、XML 選択直後に import bar が `取込待ち` と `width: 12%` を表示した。
- Wails browser 経路でも、`この XML を取り込む` 押下後 1 秒待機しても `取込中` へ変わらず、progress は `12%` のままだった。
- Wails browser console 追加 error は観測できず、継続して見えた error は `favicon.ico` 404 だけだった。
- 一時 tracing では `startStagedXmlImport()` の guard 手前ログだけが出て、`running/50` 更新後ログ、gateway 呼び出し前ログ、catch ログは出なかった。
- `resolveFileReference(file)` は `path ?? webkitRelativePath ?? file.name` を使っており、`path === ""` の時に空文字を返して `selectedFileReference` が falsy になる可能性がある。

## Trace Plan

- `stageXmlImport()` 実行時の `selectedFileName` `selectedFileReference` `importStage` `importProgress` を観測する。
- `startStagedXmlImport()` 冒頭の guard 判定値と early return の有無を観測する。
- `gateway.importMasterDictionaryXml()` 呼び出し有無、失敗時の `errorMessage`、呼び出し前後の `importStage` を観測する。
- 実 file chooser 経路で XML 選択から取込開始までを再確認し、DOM 注入差分を切り分ける。
- trace refresh: 現在の最有力は `selectedFileReference === ""` により guard で early return する frontend 側不整合であり、binding 解決失敗や runtime event 停滞は優先度を下げた。

## Fix Plan

- `resolveFileReference(file)` で空文字の `path` / `webkitRelativePath` を採用せず、`file.name` へ fallback する。
- file 選択直後の ready state では progress を `0%` に固定し、開始前に進行中へ見えないようにする。
- 既存 unit/system test を更新し、empty-string file reference fallback と pre-start progress `0%` を固定する。

## Acceptance Checks

- XML 選択後に `この XML を取り込む` で取込開始できる。
- ファイル選択時点で progress/seek bar が途中まで進んで見えない。
- 空文字の `file.path` / `webkitRelativePath` でも取込開始できる。

## Required Evidence

- `stageXmlImport()` と `startStagedXmlImport()` の開始前後 state を説明できる trace 証跡を残す。
- ファイル選択時点で progress/seek bar が進んで見えないことを UI と test の両方で証明する。
- `XMLから取り込む` 押下直後に開始状態が視覚的に分かることを UI と test の両方で証明する。
- `npx playwright test tests/system/master-dictionary-management.spec.ts --config ./playwright.config.ts` で XML 取込シナリオが通過した。
- `python3 scripts/harness/run.py --suite all` が通過した。

## Closeout Notes

- 2026-04-12 reopen: `phase-6.5-ui-check` と `phase-8-review` を未実施のまま close したため active へ戻した。
- 2026-04-12 follow-up: cursor 不足は再現せず、開始前 progress 表示と start click 時の state 遷移確認を優先する。
- 2026-04-12 reproduction refresh: `XMLから取り込む` の開始不達と、ファイル選択直後の `12%` progress 表示を Wails browser 経路で再現した。
- 2026-04-12 trace refresh: gateway 呼び出し前に guard で止まる証跡が出たため、frontend の file reference 解決と ready state 表示を narrow fix scope とする。
- 2026-04-12 review reroute: runtime fallback と cursor 契約は accepted scope 外として差分から除去した。

## Outcome

- `frontend/src/ui/screens/master-dictionary/master-dictionary-screen-controller.ts` で empty-string `path` / `webkitRelativePath` を除外し、`file.name` へ fallback するようにした。
- `frontend/src/ui/screens/master-dictionary/master-dictionary-screen-controller.ts` と `frontend/src/ui/screens/master-dictionary/master-dictionary.usecase.ts` で pre-start / error 復帰時の progress を `0%` にそろえた。
- `frontend/src/ui/App.test.ts` で empty-string file reference fallback、pre-start progress `0%`、開始表示の可視化を固定した。
- `tests/system/master-dictionary-management.spec.ts` で XML 選択後 `0%`、開始後 `取込中`、完了後 `完了` の導線を固定した。
- `python3 scripts/harness/run.py --suite all` が通過した。
- `phase-6.5-ui-check` pass: pre-start は `width: 0%`、start 後は `取込中` を経由して `完了` へ進むことを Wails browser 経路で確認した。
- `phase-5-test-implementation` pass: unit / system の両方で empty-string file reference fallback と pre-start progress 0% を固定した。
- `phase-8-review` pass: accepted fix scope と差分、test の整合を確認した。
- coverage harness は `test-results/frontend-coverage/.tmp/coverage-0.json` の `ENOENT` で別途再確認が必要。
- Sonar の open issue 1 件と quality gate `NONE` は backend 側の scope 外事項として残る。
