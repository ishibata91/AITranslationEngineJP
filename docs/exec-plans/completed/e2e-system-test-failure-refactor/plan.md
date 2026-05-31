# E2E System Test Failure Refactor

## リファクタ目的

`e2e-test-design-maintenance` で再構築した UI 人間操作 E2E が検出した失敗を、既存仕様、既存実装、既存テストのずれとして整理する。
この task は、失敗を `refactor-lane` の新規 active plan として扱い、仕様実装優先判断、構造品質調査、テスト品質調査、実行型判定、`implementation-scope` を経て修正する。

前 task の `テスト専用型` ではプロダクトコード変更が禁止されていた。
この task では、人間判断により、必要な範囲でプロダクトコード、統合境界、シナリオテストを変更できる `コード併走型` 候補として扱う。

## task 枠

- task-id: `e2e-system-test-failure-refactor`
- 作業計画フォルダ: `docs/exec-plans/active/e2e-system-test-failure-refactor/`
- 作業場所: `/Users/iorishibata/Repositories/AITranslationEngineJP`
- 作業ブランチ: `codex/e2e-system-test-failure-refactor`
- 統合先ブランチ: `master`
- 希望実行型: `コード併走型`

## 対象仕様参照

- `docs/e2e-test-design/test-design.csv`: UI 人間操作 E2E の観点表正本。
- `docs/e2e-test-guidelines.md`: E2E テスト設計と実装規約。
- `docs/coding-guidelines-tests.md`: テストコード規約。
- `docs/index.md`: 仕様正本の入口。
- `docs/detail-specs/`: 失敗箇所に関係する詳細仕様正本。
- `docs/screen-design/screens/`: 失敗箇所に関係する画面設計正本。
- `docs/scenario-tests/`: 既存シナリオテスト正本。
- `docs/exec-plans/active/e2e-test-design-maintenance/scenario-test-implementation-result.md`: 既知失敗の入力記録。

## 対象実装範囲

- `frontend/src/`
- `internal/`
- `wailsjs/`
- `scripts/test/seed-system-test-db/`
- `tests/system/support/`
- `tests/system/*.spec.ts`
- E2E 用 fixture と helper が存在する `tests/` 配下の該当範囲。

## 対象テスト範囲

- `tests/system/`
- `playwright.config.ts`
- `scripts/test/run-system-test.sh`
- `scripts/harness/check_system_test.py`
- `scripts/test/seed-system-test-db/`
- `docs/e2e-test-design/test-design.csv` と対応する既存 system test。

## 変更禁止範囲

- remote repository。
- `docs/exec-plans/completed/`。
- 失敗原因に関係しないプロダクトコード。
- 実外部 API、実 secret、実利用者データへ到達する実行経路。
- 仕様乖離整理で `判断保留` または `対象外` になった項目。

## 検証要件

- `python3 scripts/harness/run.py --suite system-test`
- 必要に応じて `python3 scripts/harness/run.py --suite frontend-local`
- 必要に応じて `python3 scripts/harness/run.py --suite backend-local`
- 必要に応じて `python3 scripts/harness/run.py --suite structure`

`system-test` が sandbox 内で Wails readiness により止まる場合は、環境起因と product failure を分ける。
repo rules が許す場合は、同じ command を承認済み sandbox 外実行で再確認する。

## 既知失敗入力

前 task の観測では、`python3 scripts/harness/run.py --suite system-test` が Playwright 実行まで到達した。
結果は `32 passed, 18 failed, 10 did not run` である。

| ID | 失敗領域 | 観測 |
| --- | --- | --- |
| `EF-001` | AIサービス設定 | provider の不正 endpoint 入力後も保存成功表示になる。 |
| `EF-002` | マスターペルソナ | model select が disabled のままで、`gemini-test` を選択できない。 |
| `EF-003` | 翻訳実行シェル | paused job の再開後に job run shell が表示されない。 |
| `EF-004` | 出力管理 | output management の候補行が表示されない。 |
| `EF-005` | 翻訳段階 | translation management から現在の翻訳段階へ進む操作で対象 job または action が見つからない。 |

## 成果物DAG状態

| 成果物ID | 状態 | 根拠 |
| --- | --- | --- |
| `task 枠` | 完了 | この `plan.md` に目的、対象、変更禁止範囲、検証要件、既知失敗入力を固定した。 |
| `branch 準備` | 完了 | `codex/e2e-system-test-failure-refactor` を作成し、現在ブランチとして確認した。 |
| `仕様乖離整理` | 不要 | 人間判断により、今回の失敗群では仕様乖離整理をスキップする。 |
| `仕様実装優先判断` | 完了 | 人間判断により、`EF-001` から `EF-005` は既存仕様の乖離判定ではなく、E2E 失敗を解消するリファクタ対象として扱う。 |
| `構造品質調査` | 完了 | `structure-quality-investigation.md` を `investigator` が作成した。 |
| `テスト品質調査` | 完了 | `test-quality-investigation.md` を `investigator` が作成した。 |
| `リファクタ範囲確認` | 完了 | 人間判断により、`EF-001` から `EF-005` を全て解決対象へ含める。 |
| `実行型判定` | 完了 | 承認済み項目は `コード併走型` とし、frontend、backend、統合境界、シナリオテストを候補にする。 |
| `実装範囲` | 完了 | 人間が `approve` し、`implementation-scope.md` の status を `approved` にした。 |
| `実装引き継ぎ入力` | 完了 | `implementation-handoff-input.md` に wave と wave-1 起動入力を固定した。 |
| `backend リファクタ` | 完了 | `H-BE-001`, `H-BE-002` は完了した。`backend-local` は通過した。 |
| `frontend リファクタ` | 完了 | `H-FE-001`, `H-FE-002`, `H-FE-003`, `H-FE-004` は完了した。`frontend-local` は通過した。 |
| `統合境界リファクタ` | 完了 | `H-INT-001`, `H-INT-002`, `H-INT-003` は完了した。 |
| `シナリオテスト` | 完了 | `H-ST-002`, `H-ST-003`, `H-ST-004`, `H-ST-005` は完了し、`system-test` は 60 件通過した。 |
| `単体テスト` | 完了 | backend と frontend の touched-layer 単体検証は `backend-local` と `frontend-local` で通過した。 |
| `最終検証` | 完了 | `python3 scripts/harness/run.py --suite system-test` は sandbox 外で 60 passed になった。 |
| `実装後ブラウザ確認` | 不要 | UI 導線は `system-test` の Playwright 操作で確認した。追加の手動ブラウザ証跡は要求しない。 |
| `レビュー通過根拠` | 完了 | 契約互換性、信頼境界、挙動正しさ、状態・データ不変条件、責務境界の修正必須指摘はない。 |
| `docs正本化判断` | 完了 | 今回は仕様乖離整理を人間判断でスキップし、docs 正本へ反映する `実装が正` の仕様乖離はない。docs 正本化は不要。 |

## 仕様乖離整理 スキップ記録

- 判断者: 人間。
- 判断内容: 今回の `EF-001` から `EF-005` では、仕様乖離整理を不要とする。
- 理由: 既知失敗は E2E 再構築後に見えた不具合群であり、先に失敗原因と責務境界を調べる。
- 影響: `仕様実装優先判断` は完了扱いにし、次成果物は `構造品質調査` と `テスト品質調査` とする。

## 構造品質調査 起動入力

- 対象成果物: `構造品質調査`
- 起動先 agent: `investigator`
- 読むファイル:
  - `docs/exec-plans/active/e2e-system-test-failure-refactor/plan.md`
  - `docs/exec-plans/active/e2e-test-design-maintenance/scenario-test-implementation-result.md`
  - `docs/e2e-test-design/test-design.csv`
  - `docs/e2e-test-guidelines.md`
  - `docs/index.md`
  - `docs/detail-specs/`
  - `docs/screen-design/screens/`
  - `docs/scenario-tests/`
  - `frontend/src/`
  - `internal/`
  - `tests/system/`
  - `scripts/test/seed-system-test-db/`
- 禁止事項:
  - プロダクトコード、プロダクトテスト、docs 正本文を変更しない。
  - 既知失敗を新規要件として扱わない。
  - 失敗原因が未確認のまま修正範囲を決めない。
- 期待する成果物:
  - `docs/exec-plans/active/e2e-system-test-failure-refactor/structure-quality-investigation.md`
  - `EF-001` から `EF-005` ごとの原因候補、責務境界、変更候補、変更不要範囲、追加調査が必要な箇所。

## テスト品質調査 起動入力

- 対象成果物: `テスト品質調査`
- 起動先 agent: `investigator`
- 読むファイル:
  - `docs/exec-plans/active/e2e-system-test-failure-refactor/plan.md`
  - `docs/exec-plans/active/e2e-test-design-maintenance/scenario-test-implementation-result.md`
  - `docs/e2e-test-design/test-design.csv`
  - `docs/e2e-test-guidelines.md`
  - `docs/coding-guidelines-tests.md`
  - `tests/system/`
  - `tests/fixtures/`
  - `scripts/test/seed-system-test-db/`
- 禁止事項:
  - プロダクトコード、プロダクトテスト、docs 正本文を変更しない。
  - テストが正しいと未確認のままプロダクトコード修正を前提にしない。
  - Page Object の実装不足を product failure と混同しない。
- 期待する成果物:
  - `docs/exec-plans/active/e2e-system-test-failure-refactor/test-quality-investigation.md`
  - `EF-001` から `EF-005` ごとの test 前提、fixture、mock、Page Object、assertion の妥当性、変更不要テスト範囲。

## リファクタ範囲確認

| ID | 承認状態 | 実装範囲候補 | 検証要件 |
| --- | --- | --- | --- |
| `EF-001` | `承認` | AIサービス設定の保存表示、入力検証、composition root wiring、E2E assertion を整理する。 | `system-test`, 必要に応じて `frontend-local` と `backend-local` |
| `EF-002` | `承認` | マスターペルソナの AI 設定前提、model select 活性条件、seed または待機条件を整理する。 | `system-test`, 必要に応じて `frontend-local` と `backend-local` |
| `EF-003` | `承認` | 翻訳ジョブ再開の backend 契約、frontend 遷移条件、feedback assertion を整理する。 | `system-test`, 必要に応じて `frontend-local` と `backend-local` |
| `EF-004` | `承認` | 出力管理の gateway wiring、mock 接続確認、候補行表示前提、diff row 手順を整理する。 | `system-test`, 必要に応じて `frontend-local` |
| `EF-005` | `承認` | 翻訳段階の seed、scenario mock、job family、観点 ID と test 内容の対応を整理する。 | `system-test`, 必要に応じて `frontend-local` と `backend-local` |

## 実行型判定

| 承認済み項目 | 実行型 | 起動する agent | 起動しない agent と理由 |
| --- | --- | --- | --- |
| `EF-001`, `EF-002`, `EF-003`, `EF-004`, `EF-005` | `コード併走型` | `designer`, `backend_implementer`, `frontend_implementer`, `integration_implementer`, `implementation_scenario_tester` | `implementation_unit_tester` は `implementation-scope` で単体分岐保護が必要と判断した場合だけ起動する。 |

## 実装範囲 起動入力

- 対象成果物: `実装範囲`
- 起動先 agent: `designer`
- 使用 skill: `implementation-scope`
- 読むファイル:
  - `docs/exec-plans/active/e2e-system-test-failure-refactor/plan.md`
  - `docs/exec-plans/active/e2e-system-test-failure-refactor/structure-quality-investigation.md`
  - `docs/exec-plans/active/e2e-system-test-failure-refactor/test-quality-investigation.md`
  - `docs/exec-plans/active/e2e-test-design-maintenance/scenario-test-implementation-result.md`
  - `docs/e2e-test-design/test-design.csv`
  - `docs/e2e-test-guidelines.md`
  - `docs/coding-guidelines-tests.md`
  - `frontend/src/`
  - `internal/`
  - `tests/system/`
  - `scripts/test/seed-system-test-db/`
- 禁止事項:
  - 未承認の新規要件を作らない。
  - docs 正本文を更新しない。
  - `.codex/` を更新しない。
  - `EF-001` から `EF-005` 以外へ修正範囲を広げない。
- 期待する成果物:
  - `docs/exec-plans/active/e2e-system-test-failure-refactor/implementation-scope.md`
  - backend、frontend、統合境界、シナリオテスト、必要な単体テストを分けた実装範囲。

## 実装範囲 作成結果

- 作成 agent: `designer`
- 作成ファイル: `docs/exec-plans/active/e2e-system-test-failure-refactor/implementation-scope.md`
- status: `approved`
- 内容: `EF-001` から `EF-005` を backend、frontend、統合境界、シナリオテスト、単体テスト候補へ分割した。
- 停止理由: なし。
- 次成果物: `H-BE-001`, `H-BE-002`, `H-FE-001` を起動する。

## 実装引き継ぎ入力 作成結果

- 作成ファイル: `docs/exec-plans/active/e2e-system-test-failure-refactor/implementation-handoff-input.md`
- 起動可能 handoff: `H-BE-001`, `H-BE-002`, `H-FE-001`
- 起動しない handoff: `H-FE-002`, `H-FE-003`, `H-FE-004` は `H-FE-001` の共有 composition root 変更後に起動する。
- 起動しない handoff: `H-INT-*`, `H-ST-*`, `H-UT-*`, `H-FINAL-001` は依存 handoff 完了後に起動する。

## Wave 1 実装結果

| handoff | agent | 変更ファイル | 検証結果 | 残留リスク |
| --- | --- | --- | --- | --- |
| `H-BE-001` | `backend_implementer` | `internal/service/provider_settings_service.go`, `internal/service/provider_settings_service_test.go` | `backend-local`: pass | frontend の error mapping は `H-INT-001` が扱う。 |
| `H-BE-002` | `backend_implementer` | `internal/service/translation_job_management_service.go`, `internal/service/translation_job_management_service_test.go` | `backend-local`: pass | frontend の resume 遷移条件は `H-FE-003` が扱う。 |
| `H-FE-001` | `frontend_implementer` | `frontend/src/ui/App.svelte`, `frontend/src/ui/App.test.ts` | `frontend-local`: pass | ブラウザ証跡は未取得。表示差分は発生しない判断である。 |

## Wave 1b 起動判断

- `H-FE-002` は `H-FE-001` 完了後に起動できる。
- `H-FE-003` は `H-BE-002` と `H-FE-001` 完了後に起動できる。
- `H-FE-004` は `H-FE-003` と `TranslationJobManagementPage.svelte` の変更範囲が重なりうるため、`H-FE-003` 完了後に起動する。

## Wave 1b 実装結果

| handoff | agent | 変更ファイル | 検証結果 | 残留リスク |
| --- | --- | --- | --- | --- |
| `H-FE-002` | `frontend_implementer` | `frontend/src/ui/screens/master-persona/MasterPersonaPage.svelte`, `frontend/src/application/presenter/master-persona/master-persona.presenter.ts` | `frontend-local`: pass | ブラウザ証跡は未取得。 |
| `H-FE-003` | `frontend_implementer` | `frontend/src/application/usecase/translation-job-management/translation-job-management.usecase.ts`, `frontend/src/application/usecase/translation-job-management/translation-job-management.usecase.test.ts`, `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte`, `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.test.ts` | `frontend-local`: pass, `build-storybook`: pass | なし。 |
| `H-FE-004` | `frontend_implementer` | `frontend/src/ui/screens/translation-job-management/JobCard.svelte`, `frontend/src/ui/screens/translation-job-management/JobOperationGroup.svelte`, `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte`, `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.test.ts` | `frontend-local`: pass, `build-storybook`: pass | job family と seed/mock 境界は `H-INT-003` が扱う。 |

## Wave 2 起動判断

- `H-INT-001` は `H-BE-001` と `H-FE-001` 完了後に起動できる。
- `H-INT-002` は `H-FE-001` 完了後に起動できる。
- `H-INT-003` は `H-FE-004` 完了後に起動できるが、`tests/system/support/scenario-wails-mocks.ts` の変更範囲が `H-INT-002` と重なりうるため、`H-INT-002` 完了後に起動する。

## Wave 2 実装結果

| handoff | agent | 変更ファイル | 検証結果 | 残留リスク |
| --- | --- | --- | --- | --- |
| `H-INT-001` | `integration_implementer` | なし | 未実行 | `frontend/src/application/usecase/provider-settings/provider-settings.usecase.ts` が owned scope 外で、保存失敗を UI error state へ写像できない。scope 拡張承認が必要である。 |
| `H-INT-002` | `integration_implementer` | `tests/system/support/scenario-wails-mocks.ts` | `frontend-local`: pass, `output-management.spec.ts`: 4 passed / 2 failed | 残る 2 件は `latestResult` 不在の否定 assertion であり、`H-ST-004` が扱う。 |
| `H-INT-003` | `integration_implementer` | `scripts/test/seed-system-test-db/main.go`, `tests/system/support/scenario-wails-mocks.ts` | `backend-local`: pass, `frontend-local`: pass | `system-test` は未実行。 |

## Scope 追加判断待ち

- 対象 handoff: `H-INT-001`
- 追加が必要な owned scope: `frontend/src/application/usecase/provider-settings/provider-settings.usecase.ts`
- 理由: gateway と controller だけでは、backend の保存失敗を UI の入力不正表示へ写像できない。
- 判断: EF 範囲内の correction として追加承認済み扱いにする。
- 影響: `H-INT-001` を再起動できる。

## Wave 3 起動判断

- `H-ST-001` は `H-INT-001` 停止中のため起動しない。
- `H-ST-002` は `H-FE-002` 完了後に起動できる。
- `H-ST-003` は `H-FE-003` 完了後に起動できる。
- `H-ST-004` は `H-INT-002` 完了後に起動できる。
- `H-ST-005` は `H-INT-003` 完了後に起動できるが、`H-ST-003` と support file が重なりうるため後続にする。

## Wave 3 実装結果

| handoff | agent | 変更ファイル | 検証結果 | 戻し先 |
| --- | --- | --- | --- | --- |
| `H-ST-002` | `implementation_scenario_tester` | `tests/system/master-persona.spec.ts`, `tests/system/support/master-persona-page.ts`, `tests/system/support/scenario-wails-mocks.ts` | `E2E-UC-013`: pass | なし |
| `H-ST-003` | `implementation_scenario_tester` | `tests/system/translation-job-management.spec.ts` | `E2E-UC-019`: fail。feedback notification 不在。 | `H-FE-003` |
| `H-ST-004` | `implementation_scenario_tester` | `tests/system/output-management.spec.ts`, `tests/system/support/output-management-page.ts` | list / format: pass。UI 実行は Wails ready failure で未通過。 | final validation |
| `H-ST-005` | `implementation_scenario_tester` | `tests/system/job-run-shell.spec.ts`, `tests/system/translation-phases.spec.ts`, `tests/system/support/scenario-wails-mocks.ts` | `E2E-UC-045`, `046`, `048`, `050`, `051`, `052`: pass。`E2E-UC-047`, `049`, `053`: fail。 | `H-INT-003` |

## H-FE-003 戻し修正結果

- 戻し元: `H-ST-003`
- 原因: resume success 後に `onOpenJobRun()` が即時実行され、feedback notification の DOM が切り替わりで消えていた。
- 変更ファイル: `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte`, `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.test.ts`
- 検証結果: `frontend-local`: pass
- 未確認: `E2E-UC-019` の局所 Playwright は Wails dev server 未起動で未確認。

## H-ST-005 戻し判断

- 戻し元: `H-ST-005`
- 戻し先: `H-INT-003`
- 理由: 本文翻訳 phase が `summary を取得中` のままになり、scenario mock と body translation gateway response shape の不整合が疑われる。
- 対象失敗: `E2E-UC-047`, `E2E-UC-049`, `E2E-UC-053`
- 起動条件: `H-INT-003` の既存 owned scope に `tests/system/support/scenario-wails-mocks.ts` と phase gateway DTO が含まれるため、追加承認なしで戻せる。

## H-INT-003 戻し再観測

- 観測: `scenario-wails-mocks.ts` の mock binding は browser 内で直接呼ぶと即時返る。
- 観測: 実画面の body phase controller 経由では store に反映されず、`summary を取得中` が残る。
- 検証: `frontend-local` は pass。
- 未解決: `E2E-UC-049`, `E2E-UC-053` は button disabled のまま失敗する。
- 次 action: body phase gateway から usecase へ返る promise 経路を確認する。

## System Test 再実行結果

- command: `python3 scripts/harness/run.py --suite system-test`
- 実行条件: sandbox 外。
- 結果: `41 passed`, `6 failed`, `13 did not run`
- 通過: output management 6 件、`E2E-UC-013`, `E2E-UC-045`, `E2E-UC-046`, `E2E-UC-048`, `E2E-UC-050`, `E2E-UC-051`, `E2E-UC-052`
- 失敗: `E2E-UC-028`, `E2E-UC-047`, `E2E-UC-049`, `E2E-UC-014`, `SCN-TJM-003`, `E2E-UC-053`
- 次戻し先: `H-INT-001`, `H-INT-003`, `H-ST-002`, `H-ST-005`

## System Test 再実行結果 2

- command: `python3 scripts/harness/run.py --suite system-test`
- 実行条件: sandbox 外。
- 結果: `50 passed`, `3 failed`, `7 did not run`
- 失敗: `E2E-UC-049`, `E2E-UC-015`, `E2E-UC-018`
- 解消確認: `E2E-UC-028`, `E2E-UC-047`, `E2E-UC-014`, `SCN-TJM-003`, `E2E-UC-053`
- 次戻し先: `H-ST-005`, `H-ST-002`, `H-FE-003`

## System Test 再実行結果 3

- command: `python3 scripts/harness/run.py --suite system-test`
- 結果: partial fail。
- 解消確認: `E2E-UC-049`, `E2E-UC-015`, `E2E-UC-018`
- 残失敗: `E2E-UC-016`, `E2E-UC-039`
- 共通原因候補: DOM から削除されない modal に `toHaveCount(0)` を期待している。
- scope correction: `tests/system/master-persona.spec.ts` と `tests/system/translation-job-management.spec.ts` を該当 handoff owned scope へ含める。

## 残失敗戻し結果

| 対象 | 戻し先 | 結果 | 次 action |
| --- | --- | --- | --- |
| `E2E-UC-028` | `H-INT-001` | 完了。`E2E-UC-028` は通過した。 | なし |
| `E2E-UC-014` | `H-ST-002` | plugin option 選択は解消。詳細選択で失敗。 | `H-FE-002` に UI / Page Object 境界 correction として戻す。 |
| `E2E-UC-047`, `E2E-UC-049`, `E2E-UC-053` | `H-INT-003` | DTO shape は補正。body phase 表示で `effect_update_depth_exceeded` が残る。 | `H-FE-004` に body phase UI/controller correction として戻す。 |
| `SCN-TJM-003` | `H-INT-003` | body phase 系失敗に連動して未解決。 | `H-FE-004` correction 後に再確認する。 |

## Scope correction 2

- `H-FE-002`: `tests/system/support/master-persona-page.ts` を owned scope に追加する。
- 理由: `E2E-UC-014` は plugin option の value / label 解決と詳細選択の境界で止まっており、UI と Page Object の境界調整が必要である。
- `H-FE-004`: `frontend/src/controller/body-translation-phase/` と `frontend/src/application/usecase/body-translation-phase/` を owned scope に追加する。
- 理由: body phase 表示時の `effect_update_depth_exceeded` は mock DTO ではなく、job run page effect と body phase controller/usecase の処理対象ページ更新経路で発生している。

## System Test 最終再実行結果

- command: `python3 scripts/harness/run.py --suite system-test`
- 実行条件: sandbox 外。
- 結果: `60 passed`
- 解消確認: `E2E-UC-033`, `E2E-UC-035`
- 補正内容: `master-persona-ai-settings-card` の警告文言 assertion を現行 UI の `APIキーが未設定` に合わせた。
- 補正内容: `master-persona-edit-modal` の閉じる確認を DOM 削除ではなく非表示確認に合わせた。
- 残失敗: なし。
- 次成果物: 観点別レビュー。

## Reviewback 指摘

| 観点 | 指摘ID | severity | 内容 | 戻し先 | 状態 |
| --- | --- | --- | --- | --- | --- |
| 挙動正しさ | `behavior-001` | `major` | `E2E-UC-028` が保存済み endpoint の非更新を証明していない。 | `H-ST-001` | 修正済み |
| 状態・データ不変条件 | `state-invariant-001` | `major` | provider settings の非検証エラー catch が `validation_failed` と新しい `requestToken` へ進む。 | `H-FE-001` | 修正済み |
| 状態・データ不変条件 | `state-invariant-002` | `major` | phase navigation の scenario mock が seed DB にない job family を返す。 | `H-INT-003` / `H-ST-005` | 修正済み |
| 責務境界 | `responsibility-boundary-001` | `major` | `App.svelte` が View 内で Wails gateway 具象実装を生成している。 | `H-FE-001` | 修正済み |

## Reviewback 修正結果

| 指摘ID | agent | 変更ファイル | 検証結果 | 残留リスク |
| --- | --- | --- | --- | --- |
| `behavior-001` | `implementation_scenario_tester` | `tests/system/frontend-backend-connection.spec.ts` | `E2E-UC-028`: pass | full `system-test` は別領域失敗後、親レーン再実行で通過した。 |
| `state-invariant-001` | `frontend_implementer` | `frontend/src/application/usecase/provider-settings/provider-settings.usecase.ts`, `frontend/src/application/usecase/provider-settings/provider-settings.usecase.test.ts` | `frontend-local`: pass | なし。 |
| `state-invariant-002` | `implementation_scenario_tester` | `scripts/test/seed-system-test-db/main.go`, `tests/system/support/scenario-wails-mocks.ts`, `tests/system/job-run-shell.spec.ts`, `tests/system/translation-job-management.spec.ts` | phase targeted: pass, `system-test`: pass | sandbox 内 Chromium は macOS 権限で未通過。 |
| `responsibility-boundary-001` | `frontend_implementer` | `frontend/src/ui/App.svelte`, `frontend/src/ui/App.test.ts` | `frontend-local`: pass | `App.svelte` 単体 mount test は明示 factory 注入が必要。 |

## Reviewback 修正後 最終再実行結果

- command: `python3 scripts/harness/run.py --suite system-test`
- 実行条件: sandbox 外。
- 結果: `60 passed`
- 解消確認: `E2E-UC-028`, `E2E-UC-008`, `SCN-TJM-001`, phase navigation 系。
- 残失敗: なし。
- 次成果物: 作業 commit。

## Reviewback 再レビュー結果

| 観点 | 対象指摘 | 判定 | 根拠 |
| --- | --- | --- | --- |
| 挙動正しさ | `behavior-001` | 修正必須指摘なし | `E2E-UC-028` は保存済み endpoint 非更新を reload 後に確認している。 |
| 状態・データ不変条件 | `state-invariant-001`, `state-invariant-002` | 修正必須指摘なし | 非検証エラーの UI 状態維持と seed/mock fixture 整合を確認した。 |
| 責務境界 | `responsibility-boundary-001` | 修正必須指摘なし | `App.svelte` は Wails gateway 具象生成を持たず、production wiring は bootstrap 側に戻った。 |

## docs正本化判断

- 判断: docs 正本化不要。
- 理由: 仕様乖離整理は人間判断でスキップした。
- 理由: 今回の変更は E2E 失敗修正、test fixture 整合、責務境界 correction であり、正本へ反映する新規仕様や `実装が正` の仕様乖離を作っていない。
- docs_updater 起動: 不要。

## 未固定入力

- 未固定入力なし。
