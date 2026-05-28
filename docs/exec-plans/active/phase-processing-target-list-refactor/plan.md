# Phase Processing Target List Refactor

## リファクタ目的

翻訳管理から新規翻訳を開始し、`dictionaries/Lucien.esp_Export.json` をロードして単語翻訳へ進むと、処理対象一覧は 0 件だが、進行状況は対象語件数 4930 件を表示する。
この task は、処理対象一覧、進行状況、検索、件数表示のずれを、既存仕様、既存実装、既存テストのずれとして整理し、各フェーズの処理対象一覧を証明できる状態へ直す。

## task 枠

- task-id: `phase-processing-target-list-refactor`
- 作業計画フォルダ: `docs/exec-plans/active/phase-processing-target-list-refactor/`
- 作業場所: `/Users/iorishibata/Repositories/AITranslationEngineJP`
- 作業ブランチ: `codex/phase-processing-target-list-refactor`
- 統合先ブランチ: `master`
- 希望実行型: `コード併走型`

## 人間観測

- 操作: 翻訳管理を開く。
- 操作: 新規翻訳を開始する。
- 操作: `dictionaries/Lucien.esp_Export.json` をロードする。
- 操作: 単語翻訳へ進む。
- 現象: 処理対象一覧は 0 件である。
- 現象: 進行状況は対象語件数 4930 件を表示する。
- 人間判断: 処理対象一覧 test と UC が不足している可能性がある。
- 必要観点: 各フェーズで、一覧の件数が進行状況と合う。
- 必要観点: 各フェーズで、処理対象一覧が表示される。
- 必要観点: 各フェーズで、処理対象一覧を検索できる。

## 対象仕様参照

- `docs/usecases/uc-translation-management.md`
- `docs/e2e-test-design/test-design.csv`
- `docs/detail-specs/term-translation-phase.md`
- `docs/detail-specs/persona-generation-phase.md`
- `docs/detail-specs/body-translation-phase.md`
- `docs/screen-design/screens/term-translation-phase.md`
- `docs/screen-design/screens/persona-generation-phase.md`
- `docs/screen-design/screens/body-translation-phase.md`
- `docs/scenario-tests/`

## 対象実装範囲

- `frontend/src/ui/screens/term-translation-phase/`
- `frontend/src/ui/screens/persona-generation-phase/`
- `frontend/src/ui/screens/body-translation-phase/`
- `frontend/src/ui/components/ProcessingTargetListPanel.svelte`
- `frontend/src/ui/components/ProcessingTargetListWrapper.svelte`
- `frontend/src/application/`
- `frontend/src/controller/`
- `internal/`
- `scripts/test/seed-system-test-db/`
- `tests/system/`

## 対象テスト範囲

- `tests/system/translation-phases.spec.ts`
- `tests/system/job-run-shell.spec.ts`
- `tests/system/support/translation-phase-pages.ts`
- `tests/system/support/scenario-wails-mocks.ts`
- `docs/e2e-test-design/test-design.csv`
- phase 関連 frontend unit test。

## 変更禁止範囲

- remote repository。
- `docs/exec-plans/completed/`。
- 既存の未 commit `.codex/` 差分。
- `docs/exec-plans/active/e2e-test-design-maintenance/`。
- 処理対象一覧、進行状況、検索、件数整合に関係しない product code。
- 実外部 API、実 secret、実利用者データへ到達する経路。
- 仕様乖離整理で `判断保留` または `対象外` になった項目。

## 検証要件

- `python3 scripts/harness/run.py --suite system-test`
- 必要に応じて `python3 scripts/harness/run.py --suite frontend-local`
- 必要に応じて `python3 scripts/harness/run.py --suite backend-local`

`system-test` が sandbox 内で Wails readiness または Chromium 権限により止まる場合は、環境起因と product failure を分ける。
repo rules が許す場合は、同じ command を承認済み sandbox 外実行で再確認する。

## 成果物DAG状態

| 成果物ID | 状態 | 根拠 |
| --- | --- | --- |
| `task 枠` | 完了 | この `plan.md` に目的、対象、変更禁止範囲、検証要件、人間観測を固定した。 |
| `branch 準備` | 完了 | `codex/phase-processing-target-list-refactor` を作成し、現在ブランチとして確認した。 |
| `仕様乖離整理` | 完了 | `spec-drift-investigation.md` を `investigator` が作成した。 |
| `仕様実装優先判断` | 完了 | 人間入力により、件数一致、一覧表示、検索を 3 フェーズ共通の必要観点として扱う。 |
| `構造品質調査` | 完了 | `structure-quality-investigation.md` を `investigator` が作成した。 |
| `テスト品質調査` | 完了 | `test-quality-investigation.md` を `investigator` が作成した。 |
| `リファクタ範囲確認` | 完了 | 人間入力により、3 フェーズの件数一致、一覧表示、検索を実装範囲へ含める。 |
| `実行型判定` | 完了 | backend read model、frontend 表示、system-test を扱うため `コード併走型` とする。 |
| `実装範囲` | 完了 | `implementation-scope.md` を `designer` が作成し、人間継続指示により `approved` にした。 |
| `実装引き継ぎ入力` | 完了 | `implementation-handoff-input.md` に共通入力、禁止事項、wave-1 起動入力を固定した。 |
| `frontend リファクタ` | 完了 | `frontend-count-subject` と `frontend-search-subject` は `frontend_implementer` が完了した。 |
| `backend リファクタ` | 完了 | `backend-count-read-model` と `backend-search-read-model` は `backend_implementer` が完了した。 |
| `統合境界リファクタ` | 完了 | `integration-processing-target-seam` は `integration_implementer` が完了した。 |
| `単体テスト` | 完了 | `unit-count-subject` と `unit-search-subject` は `implementation_unit_tester` が完了した。 |
| `シナリオテスト` | 完了 | `scenario-page-object`、`scenario-fixture`、`scenario-phase-list-search` は `implementation_scenario_tester` が完了した。 |
| `最終検証` | 完了 | `frontend-local`、`backend-local`、sandbox 外 `system-test` は通過した。sandbox 内 `system-test` は Wails readiness で失敗し、環境起因として切り分けた。 |
| `実装後ブラウザ確認` | 環境制約で未完了 | Storybook 応答と 3 phase story の存在は確認した。Playwright 起動が mac 権限で失敗したため、画面操作証跡は未取得。sandbox 外 system-test で UI 経路は通過済み。 |
| `レビュー通過根拠` | 完了 | 挙動正しさ、契約互換性、権限・信頼境界、状態不変条件、責務境界のレビューは findings none になった。 |
| `シナリオテスト差し戻し` | 完了 | fixture 件数と `E2E-UC-045/046/047` の件数一致 assertion を修正した。 |
| `frontend 差し戻し` | 完了 | load / summary 再取得内の一覧 response も最新要求判定へ含めた。 |
| `system-test 補正` | 完了 | `E2E-UC-033` の不足 AI 設定 mock を isolated に登録し、full system-test を通過させた。 |
| `docs正本化判断` | 完了 | docs 正本文は更新しない。task-local `usecase-diff.md` を今回の UC 差分成果物として残す。 |
| `作業 commit` | 完了 | 今回 task 差分だけを commit した。既存 `.codex/`、`docs/e2e-test-design/test-design.csv`、`e2e-test-design-maintenance` は除外した。 |
| `マージ準備入力` | 完了 | この active plan、作業ブランチ、commit hash、検証結果、レビュー結果、残留リスクを merge lane へ渡せる。 |

## 仕様乖離整理 起動入力

- 対象成果物: `仕様乖離整理`
- 起動先 agent: `investigator`
- 使用 skill: `investigate`
- 読むファイル:
  - `docs/exec-plans/active/phase-processing-target-list-refactor/plan.md`
  - `docs/usecases/uc-translation-management.md`
  - `docs/e2e-test-design/test-design.csv`
  - `docs/detail-specs/term-translation-phase.md`
  - `docs/detail-specs/persona-generation-phase.md`
  - `docs/detail-specs/body-translation-phase.md`
  - `docs/screen-design/screens/term-translation-phase.md`
  - `docs/screen-design/screens/persona-generation-phase.md`
  - `docs/screen-design/screens/body-translation-phase.md`
  - `tests/system/translation-phases.spec.ts`
  - `tests/system/job-run-shell.spec.ts`
  - `tests/system/support/translation-phase-pages.ts`
  - `frontend/src/ui/components/ProcessingTargetListPanel.svelte`
  - `frontend/src/ui/components/ProcessingTargetListWrapper.svelte`
  - phase 関連 frontend、controller、backend 実装。
- 禁止事項:
  - product code、product test、docs 正本文を変更しない。
  - 仕様と実装のどちらが正しいかを AI だけで決めない。
  - 人間観測を再現せずに原因を断定しない。
- 期待する成果物:
  - `docs/exec-plans/active/phase-processing-target-list-refactor/spec-drift-investigation.md`
  - 各フェーズの処理対象一覧について、仕様、実装、system-test の差分を記録する。
  - 件数一致、一覧表示、検索の各観点が UC と test にあるかをフェーズ別に記録する。

## 仕様乖離整理 結果

- 作成 agent: `investigator`
- 作成ファイル: `docs/exec-plans/active/phase-processing-target-list-refactor/spec-drift-investigation.md`
- 結論: 単語翻訳は処理対象一覧の取得元と進行件数の母集団が分かれており、人間観測の `一覧 0 件 / 対象語件数 4930 件` を説明できる可能性が高い。
- 結論: NPC ペルソナ生成と本文翻訳は、一覧表示、検索、画面上の件数一致を system-test が証明していない。
- 結論: `docs/scenario-tests/` に phase 系 scenario test 正本がない。

## 仕様実装優先判断

- 判断者: 人間。
- 判断内容: 各フェーズで、一覧の件数が他の件数表示と合うことを必要観点にする。
- 判断内容: 各フェーズで、処理対象一覧が表示されることを必要観点にする。
- 判断内容: 各フェーズで、処理対象一覧を検索できることを必要観点にする。
- 扱い: 上記 3 観点は `仕様が正` として code / test 修正候補にする。
- 扱い: 単語翻訳の `処理対象一覧` は、進行状況の AI 翻訳対象または表示上の対象件数と矛盾しない母集団へ合わせる。
- 扱い: NPC ペルソナ生成と本文翻訳は、現行実装の一覧取得方針を尊重しつつ、件数一致、一覧表示、検索を system-test で証明する。

## 構造品質調査 起動入力

- 対象成果物: `構造品質調査`
- 起動先 agent: `investigator`
- 使用 skill: `investigate`
- 読むファイル:
  - `docs/exec-plans/active/phase-processing-target-list-refactor/plan.md`
  - `docs/exec-plans/active/phase-processing-target-list-refactor/spec-drift-investigation.md`
  - `frontend/src/ui/components/ProcessingTargetListPanel.svelte`
  - `frontend/src/ui/components/ProcessingTargetListWrapper.svelte`
  - `frontend/src/ui/screens/term-translation-phase/`
  - `frontend/src/ui/screens/persona-generation-phase/`
  - `frontend/src/ui/screens/body-translation-phase/`
  - `frontend/src/application/`
  - `frontend/src/controller/`
  - `internal/service/`
  - `internal/repository/`
  - `internal/controller/wails/`
- 禁止事項:
  - product code、product test、docs 正本文を変更しない。
  - 仕様実装優先判断外へ範囲を広げない。
  - 件数の母集団を根拠なしに統合しない。
- 期待する成果物:
  - `docs/exec-plans/active/phase-processing-target-list-refactor/structure-quality-investigation.md`
  - 母集団不一致、frontend state、backend read model、repository query、Wails DTO の責務境界をフェーズ別に記録する。
  - 変更不要範囲を明示する。

## テスト品質調査 起動入力

- 対象成果物: `テスト品質調査`
- 起動先 agent: `investigator`
- 使用 skill: `investigate`
- 読むファイル:
  - `docs/exec-plans/active/phase-processing-target-list-refactor/plan.md`
  - `docs/exec-plans/active/phase-processing-target-list-refactor/spec-drift-investigation.md`
  - `docs/e2e-test-design/test-design.csv`
  - `tests/system/translation-phases.spec.ts`
  - `tests/system/job-run-shell.spec.ts`
  - `tests/system/support/translation-phase-pages.ts`
  - `tests/system/support/scenario-wails-mocks.ts`
  - phase 関連 frontend unit test。
- 禁止事項:
  - product code、product test、docs 正本文を変更しない。
  - 現行 test が通っていることを証明十分とみなさない。
  - Page Object 不足と product failure を混同しない。
- 期待する成果物:
  - `docs/exec-plans/active/phase-processing-target-list-refactor/test-quality-investigation.md`
  - 件数一致、一覧表示、検索の system-test 欠落をフェーズ別に記録する。
  - 追加すべき scenario test と Page Object 変更候補を記録する。

## 構造品質調査 結果

- 作成 agent: `investigator`
- 作成ファイル: `docs/exec-plans/active/phase-processing-target-list-refactor/structure-quality-investigation.md`
- 結論: 単語翻訳は summary と一覧の主語が一致していない。
- 結論: 本文翻訳は `TargetCount` と `ProviderTargetCount` の二重主語が UI と一覧の不一致を生みうる。
- 結論: NPC ペルソナ生成は件数母集団が近いが、検索文言と query 対象がずれている。
- 変更候補: phase ごとの count 主語、検索主語、processing target read model 主語を固定する。
- 変更不要範囲: `translation_complete`、`ProcessingTargetListPanel.svelte` のページング本体、phase 以外の job lifecycle 全体。

## テスト品質調査 結果

- 作成 agent: `investigator`
- 作成ファイル: `docs/exec-plans/active/phase-processing-target-list-refactor/test-quality-investigation.md`
- 結論: `E2E-UC-045/046/047` は現行 system-test で「phase を開けること」へ読み替えられており、件数一致、一覧表示、検索を証明していない。
- 結論: `translation-phase-pages.ts` は検索 input、件数表示、空状態、ページ移動を観測できない。
- 結論: `scenario-wails-mocks.ts` は processing target list を 1 件固定で返すため、件数差や検索差を診断できない。
- 変更候補: phase Page Object、scenario mock、`E2E-UC-045/046/047` または追加 scenario test の整備。
- 変更不要範囲: `E2E-UC-048/049/050` の次フェーズ遷移、`E2E-UC-051/052/053` の AI 設定不足例外。

## リファクタ範囲確認

| 項目 | 承認状態 | 実装範囲候補 | 検証要件 |
| --- | --- | --- | --- |
| 単語翻訳の処理対象一覧 | `承認` | 一覧 total、行表示、検索、進行状況の対象件数を同じ処理対象母集団へそろえる。 | `frontend-local`, `backend-local`, `system-test` |
| NPC ペルソナ生成の処理対象一覧 | `承認` | 一覧 total、行表示、検索を system-test で証明し、検索文言と query 対象を矛盾しない表現へそろえる。 | `frontend-local`, `backend-local`, `system-test` |
| 本文翻訳の処理対象一覧 | `承認` | 一覧 total、行表示、検索、進行状況の対象件数を同じ処理対象母集団へそろえる。 | `frontend-local`, `backend-local`, `system-test` |
| phase system-test coverage | `承認` | `E2E-UC-045/046/047` または追加 test で、件数一致、一覧表示、検索を 3 フェーズ分証明する。 | `system-test` |

## 実行型判定

| 承認済み項目 | 実行型 | 起動する agent | 起動しない agent と理由 |
| --- | --- | --- | --- |
| 3 フェーズの処理対象一覧 | `コード併走型` | `designer`, `backend_implementer`, `frontend_implementer`, `integration_implementer`, `implementation_scenario_tester`, `implementation_unit_tester?` | `docs_updater` は docs 正本化判断まで起動しない。 |

## 実装範囲 起動入力

- 対象成果物: `実装範囲`
- 起動先 agent: `designer`
- 使用 skill: `implementation-scope`
- 読むファイル:
  - `docs/exec-plans/active/phase-processing-target-list-refactor/plan.md`
  - `docs/exec-plans/active/phase-processing-target-list-refactor/spec-drift-investigation.md`
  - `docs/exec-plans/active/phase-processing-target-list-refactor/structure-quality-investigation.md`
  - `docs/exec-plans/active/phase-processing-target-list-refactor/test-quality-investigation.md`
  - `docs/e2e-test-design/test-design.csv`
  - `docs/detail-specs/term-translation-phase.md`
  - `docs/detail-specs/persona-generation-phase.md`
  - `docs/detail-specs/body-translation-phase.md`
  - `docs/screen-design/screens/term-translation-phase.md`
  - `docs/screen-design/screens/persona-generation-phase.md`
  - `docs/screen-design/screens/body-translation-phase.md`
  - `frontend/src/ui/components/ProcessingTargetListPanel.svelte`
  - `frontend/src/ui/components/ProcessingTargetListWrapper.svelte`
  - phase 関連 frontend、controller、backend、system-test。
- 固定判断:
  - 3 フェーズで、一覧 total と画面上の処理対象件数は同じ母集団を示す。
  - 3 フェーズで、一覧行が表示される。
  - 3 フェーズで、検索できる。
  - 単語翻訳と本文翻訳は、処理対象一覧の母集団を進行状況の対象件数と矛盾しない主語へそろえる。
  - NPC ペルソナ生成は現行の対象件数母集団を尊重し、検索表示と query 対象を矛盾しない主語へそろえる。
- 禁止事項:
  - 未承認の新規要件を作らない。
  - docs 正本文を更新しない。
  - `.codex/` を更新しない。
  - `translation_complete`、phase 以外の job lifecycle、実外部 API、実 secret、実利用者データへ広げない。
- 期待する成果物:
  - `docs/exec-plans/active/phase-processing-target-list-refactor/implementation-scope.md`
  - backend、frontend、統合境界、シナリオテスト、必要な単体テストを別 handoff に分ける。
  - 件数主語、検索主語、Page Object / fixture / scenario test の変更単位を分ける。

## 実装引き継ぎ入力 作成結果

- 作成ファイル: `docs/exec-plans/active/phase-processing-target-list-refactor/implementation-handoff-input.md`
- 起動済み handoff: `frontend-count-subject`, `frontend-search-subject`, `backend-count-read-model`, `backend-search-read-model`, `integration-processing-target-seam`
- 起動済み handoff: `unit-count-subject`, `unit-search-subject`, `scenario-page-object`, `scenario-fixture`
- 起動可能 handoff: なし。
- 起動しない handoff: なし。

## Frontend Count Subject 結果

- 担当 agent: `frontend_implementer`
- 完了 handoff: `frontend-count-subject`
- 変更内容: 単語翻訳の上部メトリクスと進行詳細を `AI 翻訳対象語` に変更した。
- 変更内容: 本文翻訳の上部メトリクスと進行詳細を `AI 送信対象` に変更した。
- 維持内容: NPC ペルソナ生成の件数主語は既存の `対象件数` を維持した。
- 検証結果: `python3 scripts/harness/run.py --suite frontend-local` は通過した。
- 検証結果: `npm --prefix frontend run build-storybook` は通過した。
- 残留事項: backend read model の total 母集団調整は後続 handoff の責務として残る。

## Frontend Search Subject 結果

- 担当 agent: `frontend_implementer`
- 完了 handoff: `frontend-search-subject`
- 変更内容: 3 フェーズの検索 input に phase 固有の `data-testid` を追加した。
- 変更内容: 一覧件数表示と 0 件空状態へ system-test 用の観測点を追加した。
- 変更内容: NPC ペルソナ生成の検索 placeholder を、名前、FormID、EditorID、NPC 属性と矛盾しない文言へ変更した。
- 維持内容: `ProcessingTargetListPanel.svelte` のページング計算本体は変更していない。
- 検証結果: `python3 scripts/harness/run.py --suite frontend-local` は通過した。
- 検証結果: `npm --prefix frontend run build-storybook` は通過した。
- 残留事項: system-test 側 Page Object と fixture が追加 test id を読む修正は後続 handoff の責務として残る。

## Backend Count Read Model 結果

- 担当 agent: `backend_implementer`
- 完了 handoff: `backend-count-read-model`
- 変更内容: 単語翻訳の一覧 total を、翻訳入力の用語候補から共通辞書 `master` 一致を除外した母集団へ変更した。
- 変更内容: 単語翻訳の一覧 row は、AI 翻訳対象語の source term、record type、ジョブ内訳語候補を返すようにした。
- 維持内容: NPC ペルソナ生成の母集団は変更していない。
- 維持内容: 本文翻訳の `dictionary_exact_match` 除外母集団は変更していない。
- 維持内容: `translation_complete` の SQL と page state は変更していない。
- 検証結果: `go test ./internal/repository` は通過した。
- 検証結果: `python3 scripts/harness/run.py --suite backend-local` は通過した。
- 残留事項: 検索条件の意味整理は `backend-search-read-model` の責務として残る。

## Backend Search Read Model 結果

- 担当 agent: `backend_implementer`
- 完了 handoff: `backend-search-read-model`
- 変更内容: 単語翻訳の検索条件を、処理対象名と訳語候補相当に限定した。
- 変更内容: NPC ペルソナ生成の検索条件は、名前、FormID、EditorID、NPC 属性を対象にする既存 query を test で固定した。
- 変更内容: 本文翻訳の検索条件を、行名、原文、訳語に限定した。
- 変更内容: repository test で、検索語なし、検索一致 1 件、検索結果 0 件を 3 phase ごとに区別した。
- 維持内容: `translation_complete` の query と page state は変更していない。
- 検証結果: `go test ./internal/repository -run 'TestSQLiteProcessingTargetList(Searches|TermTranslationUses|PersonaGenerationKeeps|BodyTranslationExcludes)'` は通過した。
- 検証結果: `python3 scripts/harness/run.py --suite backend-local` は通過した。
- 残留事項: Wails DTO、frontend gateway、phase usecase の接続境界確認は `integration-processing-target-seam` の責務として残る。

## Integration Processing Target Seam 結果

- 担当 agent: `integration_implementer`
- 完了 handoff: `integration-processing-target-seam`
- 判断内容: DTO shape は追加しない。既存の `totalCount`、`items`、`searchQuery` で目的を満たせるため。
- 変更内容: Wails controller test で `phase`、`page`、`pageSize`、`searchQuery` の usecase 受け渡しを固定した。
- 変更内容: 3 フェーズ gateway test で `GetProcessingTargetList` の request と `totalCount` response を固定した。
- 変更内容: persona gateway に processing target list response の runtime validation を追加した。
- 変更内容: 3 フェーズ usecase test で `totalCount`、`items`、`searchQuery` が page state へ届くことを固定した。
- 検証結果: `GOCACHE=/Users/iorishibata/Repositories/AITranslationEngineJP/.cache/go-build python3 scripts/harness/run.py --suite backend-local` は通過した。
- 検証結果: `python3 scripts/harness/run.py --suite frontend-local` は通過した。
- 残留事項: 実画面確認と system-test は後続成果物で扱う。

## Unit Count Subject 結果

- 担当 agent: `implementation_unit_tester`
- 完了 handoff: `unit-count-subject`
- 変更内容: NPC ペルソナ生成の presenter test に、`progress.totalCount` と `progress.targetCount` の表示主語差を検出する test を追加した。
- 検証結果: `python3 scripts/harness/run.py --suite backend-local` は通過した。
- 一時失敗: `python3 scripts/harness/run.py --suite frontend-local` は検索 page reset 期待の不一致で失敗した。
- 解消結果: 上記失敗は `unit-search-subject` で期待値を仕様へ合わせ、frontend-local 通過で解消した。

## Unit Search Subject 結果

- 担当 agent: `implementation_unit_tester`
- 完了 handoff: `unit-search-subject`
- 変更内容: 3 phase それぞれで `setProcessingTargetSearchQuery()` が `GetProcessingTargetList` request を `page: 1` で送る test を追加した。
- 確認内容: repository test は、検索語なし、検索一致 1 件、検索結果 0 件を 3 phase 分証明済みである。
- 確認内容: gateway test は、`phase`、`page`、`pageSize`、`searchQuery` の転送を 3 phase で証明済みである。
- 検証結果: `python3 scripts/harness/run.py --suite frontend-local` は通過した。
- 検証結果: `python3 scripts/harness/run.py --suite backend-local` は通過した。

## Scenario Page Object 結果

- 担当 agent: `implementation_scenario_tester`
- 完了 handoff: `scenario-page-object`
- 変更内容: Page Object に `processingTargetSearchInput`、`processingTargetListRegion`、`processingTargetTotalCount`、`processingTargetEmptyState` を追加した。
- 変更内容: `searchProcessingTargets(query)` を追加した。
- 観測方法: 一覧領域は phase screen 配下の `aria-label=\"処理対象一覧\"` で参照する。
- 検証結果: `npx playwright test --config ./playwright.config.ts --list` は成功した。

## Scenario Fixture 結果

- 担当 agent: `implementation_scenario_tester`
- 完了 handoff: `scenario-fixture`
- 変更内容: 3 フェーズ別に、単語翻訳、NPC ペルソナ生成、本文翻訳の処理対象行を分けた。
- 変更内容: `request.searchQuery` により検索一致 1 件と検索結果 0 件を返せるようにした。
- 変更内容: `request.phase`、`request.page`、`request.pageSize`、`request.searchQuery` を fixture 応答へ反映した。
- 検証結果: `npx playwright test --config ./playwright.config.ts --list` は成功した。
- 残留事項: 実シナリオ実行は `scenario-phase-list-search` の責務として残る。

## Scenario Phase List Search 結果

- 担当 agent: `implementation_scenario_tester`
- 完了 handoff: `scenario-phase-list-search`
- 変更内容: `E2E-UC-045` は単語翻訳で 3 行、`1-3 / 3 件`、`Dragonborn` 検索一致、検索 0 件を検証する。
- 変更内容: `E2E-UC-046` は NPC ペルソナ生成で 2 行、`1-2 / 2 件`、`FemaleEvenToned` 検索一致、検索 0 件を検証する。
- 変更内容: `E2E-UC-047` は本文翻訳で 4 行、`1-4 / 4 件`、`burdens` 検索一致、検索 0 件を検証する。
- 維持内容: `E2E-UC-048/049/050` の次フェーズ遷移、完了画面遷移、遷移不可条件の意図は維持した。
- 維持内容: `E2E-UC-051/052/053` の AI 設定不足で未開始を維持する意図は維持した。
- 検証結果: `git diff --check -- tests/system/job-run-shell.spec.ts` は通過した。
- 検証結果: sandbox 内 `python3 scripts/harness/run.py --suite system-test` は Wails readiness で失敗した。
- 切り分け: sandbox 外 `python3 scripts/harness/run.py --suite system-test` は `60 passed` で通過した。
- 扱い: sandbox 内失敗は、product failure ではなく実行環境起因として扱う。

## 最終検証 結果

- 実行者: `refactor_lane`
- 検証結果: `python3 scripts/harness/run.py --suite frontend-local` は通過した。`53 passed` test files、`512 passed` tests。
- 検証結果: `python3 scripts/harness/run.py --suite backend-local` は通過した。
- 検証結果: sandbox 内 `python3 scripts/harness/run.py --suite system-test` は Wails readiness で失敗した。
- 切り分け: 承認付き sandbox 外 `python3 scripts/harness/run.py --suite system-test` は `60 passed` で通過した。
- 扱い: sandbox 内 system-test 失敗は環境起因とし、product failure ではない。

## 実装後ブラウザ確認 結果

- 担当 agent: `browser_confirmation`
- confirmation_url: `http://localhost:6008/`
- 確認内容: Storybook の `index.json` に、単語翻訳、NPC ペルソナ生成、本文翻訳の 3 phase story が存在することを確認した。
- 未確認内容: 画面操作による一覧行、件数表示、検索一致、検索 0 件の確認は未実施。
- 失敗理由: Playwright 起動時に `bootstrap_check_in ... Permission denied (1100)` が発生した。
- 追加警告: Storybook 起動時に `Watchpack Error (watcher): EMFILE: too many open files, watch` が発生した。
- 代替根拠: sandbox 外 `python3 scripts/harness/run.py --suite system-test` は `60 passed` で通過し、`E2E-UC-045/046/047` が一覧行、件数表示、検索一致、検索 0 件を検証した。
- 扱い: ブラウザ確認の未完了は環境制約として残し、product failure ではない。

## レビュー通過根拠 結果

- 契約互換性レビュー: findings は `none`。
- 権限・信頼境界レビュー: findings は `none`。
- 挙動正しさレビュー: `major` 指摘あり。
- 挙動正しさレビュー指摘: `E2E-UC-045/046/047` が画面上の処理対象件数と一覧 total の一致を証明していない。
- 挙動正しさレビュー指摘: fixture は単語翻訳の `aiTargetCount` と本文翻訳の `providerTargetCount` が一覧 total と異なる。
- 必要修正: fixture の件数を検索前の一覧母集団と同じ主語にそろえ、`E2E-UC-045/046/047` で初期表示時の画面上の処理対象件数と一覧 total の一致を検証する。

## シナリオテスト差し戻し 結果

- 担当 agent: `implementation_scenario_tester`
- 完了 handoff: `scenario-phase-list-search-review-fix`
- 変更内容: 単語翻訳 fixture は、画面上の対象件数と一覧 total を `3 件` に統一した。
- 変更内容: NPC ペルソナ生成 fixture は、画面上の対象件数と一覧 total を `2 件` に統一した。
- 変更内容: 本文翻訳 fixture は、画面上の対象件数と一覧 total を `4 件` に統一した。
- 変更内容: `E2E-UC-045/046/047` は、初期表示時の一覧 total と画面上の処理対象件数の一致を検証する。
- 変更内容: `E2E-UC-045/046/047` は、検索 0 件後も一覧 total だけが `0 件` になり、画面上の処理対象件数が検索前母集団として維持されることを検証する。
- 検証結果: sandbox 内 `python3 scripts/harness/run.py --suite system-test` は Wails readiness で失敗した。
- 切り分け: sandbox 外 `python3 scripts/harness/run.py --suite system-test` は `60 passed` で通過した。
- 扱い: sandbox 内失敗は環境起因であり、product failure ではない。

## レビュー再確認 結果

- 挙動正しさ再レビュー: findings は `none`。前回 major 指摘は解消済み。
- 状態不変条件レビュー: `major` 指摘あり。
- 状態不変条件レビュー指摘: 3 フェーズ usecase で、検索応答の到着順により古い検索結果が `processingTargetPageState` を上書きする可能性がある。
- 必要修正: 最新検索要求だけが `page`、`pageSize`、`totalCount`、`searchQuery`、`items` を更新できる状態不変条件を追加する。
- 必要修正: 検索 A と検索 B の応答順序を逆転させる usecase test を追加する。
- 責務境界レビュー: `major` 指摘あり。ただし指摘対象の `.codex/` と `docs/e2e-test-design/test-design.csv` は今回 task 開始前からの未 commit 差分であり、今回 task の commit 対象から除外する。

## Frontend 差し戻し 結果

- 担当 agent: `frontend_implementer`
- 完了 handoff: `frontend-search-response-order-review-fix`
- 変更内容: 3 フェーズ usecase に最新要求 guard を追加した。
- 変更内容: 単語翻訳と NPC ペルソナ生成は usecase 単位の request sequence で古い応答を破棄する。
- 変更内容: 本文翻訳は phase 別 request sequence で古い応答を破棄する。
- 変更内容: 3 フェーズ test に、検索 A と検索 B の応答順序を逆転させる検証を追加した。
- 維持内容: `setProcessingTargetSearchQuery()` が `page: 1` request を送る既存 test は維持した。
- 検証結果: `npm --prefix frontend run test -- term-translation-phase.usecase.test.ts persona-generation-phase.usecase.test.ts body-translation-phase.usecase.test.ts` は通過した。
- 検証結果: `python3 scripts/harness/run.py --suite frontend-local` は通過した。
- 状態不変条件再レビュー: 検索 A/B 応答順序逆転は解消済み。
- 状態不変条件再レビュー指摘: `load` / summary 再取得内の一覧 response が最新検索結果を上書きできる。
- 必要修正: `load` / summary 再取得内の処理対象一覧 response も、検索と同じ最新要求判定に含める。
- 必要修正: 3 usecase test に、`load` の一覧 response と検索 response の解決順逆転を追加する。

## Frontend Load Response Order 差し戻し 結果

- 担当 agent: `frontend_implementer`
- 完了 handoff: `frontend-load-response-order-review-fix`
- 変更内容: 単語翻訳、NPC ペルソナ生成、本文翻訳で、load 側の処理対象一覧 request を既存の最新 request sequence に参加させた。
- 変更内容: summary / readiness の更新は維持し、古い一覧 response の page state 反映だけを抑止した。
- 変更内容: 3 usecase test に、load の一覧 response と検索 response の解決順逆転 test を追加した。
- 維持内容: 既存の検索 A/B 逆順 test と `page: 1` request test は維持した。
- 検証結果: `npm --prefix frontend test -- --run src/application/usecase/term-translation-phase/term-translation-phase.usecase.test.ts src/application/usecase/persona-generation-phase/persona-generation-phase.usecase.test.ts src/application/usecase/body-translation-phase/body-translation-phase.usecase.test.ts` は通過した。
- 検証結果: `python3 scripts/harness/run.py --suite frontend-local` は通過した。

## System-test 補正 結果

- 担当 agent: `implementation_scenario_tester`
- 完了 handoff: `system-test-master-persona-ai-settings-isolation-fix`
- 変更内容: `E2E-UC-033` は不足 AI 設定 mock だけを登録してから画面を開く。
- 変更内容: `master-persona.spec.ts` の他 scenario は共通 helper で設定済み mock を登録する。
- 理由: `beforeEach` の設定済み mock と不足 mock の二重登録により、AI 設定不足 scenario が設定済み状態へ汚染されていた。
- 検証結果: `npx playwright test --config ./playwright.config.ts tests/system/master-persona.spec.ts -g "E2E-UC-033"` は通過した。
- 検証結果: sandbox 内 `python3 scripts/harness/run.py --suite system-test` は Wails readiness で失敗した。
- 切り分け: 承認付き sandbox 外 `python3 scripts/harness/run.py --suite system-test` は `60 passed` で通過した。
- 確認内容: `E2E-UC-045/046/047` は full system-test 内で通過した。

## レビュー最終結果

- 挙動正しさレビュー: findings は `none`。前回 major 指摘は解消済み。
- 契約互換性レビュー: findings は `none`。
- 権限・信頼境界レビュー: findings は `none`。
- 状態不変条件レビュー: findings は `none`。検索 A/B 応答順序逆転と load/search 応答順序逆転は解消済み。
- 責務境界レビュー: findings は `none`。`.codex/`、`docs/e2e-test-design/test-design.csv`、`e2e-test-design-maintenance` は今回 task レビュー対象外として除外した。

## docs正本化判断

- 判断: docs 正本文は更新しない。
- 理由: ユーザー要望は UC 差分であり、今回 task では `docs/exec-plans/active/phase-processing-target-list-refactor/usecase-diff.md` を task-local 成果物として作成済み。
- 理由: `docs/e2e-test-design/test-design.csv` には今回 task 開始前からの未 commit 差分があり、今回 task で混ぜない。
- 扱い: docs 正本化が必要な場合は、別途 `updating-docs` 起動対象とする。

## 作業 commit 結果

- 作業ブランチ: `codex/phase-processing-target-list-refactor`
- commit: `HEAD Fix phase processing target lists`
- stage 対象: 今回 task の product code、product test、system-test、task-local 成果物。
- stage 除外: `.codex/`、`docs/e2e-test-design/test-design.csv`、`docs/exec-plans/active/e2e-test-design-maintenance/`。
- 検証結果: `python3 scripts/harness/run.py --suite frontend-local` は通過した。
- 検証結果: `python3 scripts/harness/run.py --suite backend-local` は通過した。
- 検証結果: sandbox 内 `python3 scripts/harness/run.py --suite system-test` は Wails readiness で失敗した。
- 切り分け: 承認付き sandbox 外 `python3 scripts/harness/run.py --suite system-test` は `60 passed` で通過した。

## マージ準備入力

- active plan folder: `docs/exec-plans/active/phase-processing-target-list-refactor/`
- 作業ブランチ: `codex/phase-processing-target-list-refactor`
- 統合先ブランチ: `master`
- commit hash: `HEAD`
- docs 正本化結果: docs 正本文更新なし。UC 差分は task-local `usecase-diff.md` に固定済み。
- 残留リスク: 実ブラウザ操作証跡は mac 権限で未取得。sandbox 外 system-test が UI 経路を通過済み。
- 残留差分: `.codex/`、`docs/e2e-test-design/test-design.csv`、`docs/exec-plans/active/e2e-test-design-maintenance/` は今回 task 外の未 commit 差分として残る。

## 未固定入力

- 未固定入力なし。
