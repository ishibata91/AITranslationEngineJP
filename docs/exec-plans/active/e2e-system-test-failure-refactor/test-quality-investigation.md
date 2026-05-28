# テスト品質調査

## 判定

- 判断結果: 完了
- 調査 mode: `テスト品質調査`
- 引き継ぎ先: `designer`
- 対象: `EF-001` から `EF-005`

## 根拠参照

- `docs/exec-plans/active/e2e-system-test-failure-refactor/plan.md`
- `docs/exec-plans/active/e2e-test-design-maintenance/scenario-test-implementation-result.md:61-116`
- `docs/e2e-test-design/test-design.csv:14-34,38,51-53`
- `docs/e2e-test-guidelines.md:5-66`
- `docs/coding-guidelines-tests.md:7-59`
- `tests/system/frontend-backend-connection.spec.ts:67-156`
- `tests/system/master-persona.spec.ts:19-37`
- `tests/system/output-management.spec.ts:10-94`
- `tests/system/translation-job-management.spec.ts:184-199`
- `tests/system/job-run-shell.spec.ts:13-90`
- `tests/system/translation-phases.spec.ts:16-67`
- `tests/system/support/provider-settings-page.ts:10-105`
- `tests/system/support/master-persona-page.ts:10-189`
- `tests/system/support/output-management-page.ts:10-93`
- `tests/system/support/translation-job-management-page.ts:10-139`
- `tests/system/support/job-run-shell-page.ts:10-69`
- `tests/system/support/translation-phase-pages.ts:20-95`
- `tests/system/support/scenario-wails-mocks.ts:91-123,315-510`
- `tests/fixtures/master-persona/system-test-persona.json:1-10`
- `scripts/test/seed-system-test-db/main.go:59-220`

## テスト規約観点別結果

- 仕様根拠: 各失敗は `test-design.csv` の該当観点に対応している。観点欠落よりも、前提固定不足または期待値の絞り込み不足が混在している。
- 前提明示: `EF-002` は AI 設定ロード後の model select 活性化待機が未明示である。`EF-003` は再開後の feedback を前提確認せず画面遷移だけを見る。`EF-004` と `EF-005` は scenario mock 依存で、mock 接続成立の確認が不足する。
- 決定性: `EF-002`、`EF-004`、`EF-005` の fixture と mock 自体は決定的である。`EF-003` の seed も job 状態と runtime snapshot を固定している。
- mock 境界: `scenario-wails-mocks.ts` は外部境界代替に限定されているが、`EF-004` と `EF-005` は mock を使う spec 群と実 backend を使う spec 群が同じ操作語彙を共有しており、接続成立確認がないまま product failure と見分けにくい。
- Page Object: 指定どおり画面操作と selector 解決に留まっている。今回の 5 件で、Page Object 自体が seed や assertion を持ち込んだ形跡はない。
- 失敗診断: `EF-001`、`EF-003`、`EF-005` は失敗時に「入力検証不備」「遷移不備」「mock 接続不備」を一意に切り分けにくい。観測点追加余地がある。

## EF別結果

### EF-001 AIサービス設定

- 対象観点: `E2E-UC-028`
- 観測事実: spec は不正 endpoint 保存時に `provider-settings-screen-summary-region` の入力不正表示と、`provider-settings-settings-detail-region` の未更新を要求する。`tests/system/frontend-backend-connection.spec.ts:145-156` は summary のみを確認し、保存済み詳細の未更新を確認していない。
- test 前提: `openProviderSettingsThroughProductionApp` は実 backend の画面初期表示だけを確認する。`invalid-endpoint` の事前保存状態や、保存前の詳細状態は固定していない。
- fixture / mock: 専用 fixture と mock は使っていない。実 backend 前提のため、観測点が少ないと product validation failure と UI 通知差異を切り分けにくい。
- Page Object 妥当性: `ProviderSettingsPage` は `endpointInput`、`saveButton`、`settingsDetailRegion` を露出しており不足は見えない。Page Object 不足を product failure と混同する根拠はない。
- assertion 妥当性: summary のみでは `test-design.csv:38` の後半期待値を満たさない。エラー表示が出ても保存成功状態へ更新された場合を検出できない。
- 仕様整合性: 部分一致。入力不正表示の確認は仕様に沿うが、保存済み状態へ更新されない確認が欠ける。
- 変更不要テスト範囲: `E2E-UC-003`、`E2E-UC-004`、`E2E-UC-005`、`E2E-UC-006`、`E2E-UC-027` の selector 利用と Page Object 責務分離は維持でよい。
- 未確認事項: 実 backend で不正入力時に summary と detail のどちらへ主要エラーを出すかは未確認である。

### EF-002 マスターペルソナ

- 対象観点: `E2E-UC-013`
- 観測事実: spec は AI サービス、モデル、実行方法が選択可能状態であることを前提にする。`tests/system/master-persona.spec.ts:26-33` は `setJsonFile` 後に provider を選んで直ちに model を選択する。`MasterPersonaPage.selectAISettings` は model select の活性化や option 出現を待たない。`scenario-wails-mocks.ts:91-123,384-408` は `gemini-test` を返す決定的 mock を持つ。
- test 前提: `tests/fixtures/master-persona/system-test-persona.json:1-10` は最小 fixture で妥当である。問題は JSON 内容より、provider 選択後に model list が反映されるまでの待機欠如にある。
- fixture / mock: fixture は最小で代表性がある。mock も `configured` 時に `gemini-test` を返すため、テスト前提は再現可能である。
- Page Object 妥当性: `MasterPersonaPage` は `modelSelect` を test id で直接取得している。selector 不足ではない。待機責務を Page Object に置くか test 本体に置くかは未確定だが、現状は活性化待機が存在しない。
- assertion 妥当性: `generateButton` の有効化確認はあるが、model select が有効になったこと自体の確認がない。失敗時に mock 不足か UI 遷移待機不足かを分離しにくい。
- 仕様整合性: 仕様前提のうち「モデルが選択可能状態」の証明が不足する。`test-design.csv:14` は生成中表示と結果表示まで求めるが、現行 test は model 選択直後に生成へ進むため前提失敗に弱い。
- 変更不要テスト範囲: `E2E-UC-014`、`E2E-UC-015`、`E2E-UC-016`、`E2E-UC-033`、`E2E-UC-034`、`E2E-UC-035` の fixture 最小性と mock 決定性は維持でよい。
- 未確認事項: provider 選択後に model select が活性化するまでの UI 契機が、即時更新か非同期更新かは未確認である。

### EF-003 翻訳実行シェル

- 対象観点: `E2E-UC-019`
- 観測事実: spec は paused job 再開後に `feedback-notification` の再開結果表示と、必要に応じた `job-run-job-run-shell` 表示を求める。`tests/system/translation-job-management.spec.ts:184-199` は shell 表示だけを見て、feedback を確認しない。seed は paused job と body translation phase、runtime snapshot を固定している。
- test 前提: `scripts/test/seed-system-test-db/main.go:89-220` の paused job は `JobState=paused`、`PhaseState=paused`、runtime snapshot 保存済みで、再開対象として十分に明示されている。
- fixture / mock: mock は使っていない。実 backend seed に依存するため、失敗時の切り分けには「再開 API 成功」「一覧再描画」「画面遷移」の観測分離が必要である。
- Page Object 妥当性: `TranslationJobManagementPage.resumeButton` と `JobRunShellPage.shell` / `phaseScreenRegion` は selector 取得として妥当である。Page Object 不足の根拠はない。
- assertion 妥当性: shell 可視だけでは `test-design.csv:20` の主要期待値のうち feedback を欠く。再開結果は出たが遷移しない場合と、再開自体が失敗した場合を分離できない。
- 仕様整合性: 部分一致。再開後に job run shell を見る点は合うが、spec が先に要求する再開結果通知の観測が欠ける。
- 変更不要テスト範囲: `SCN-TJM-001`、`E2E-UC-017`、`E2E-UC-036`、`SCN-TJM-003`、`SCN-TJM-005/007`、`E2E-UC-018`、`E2E-UC-037`、`E2E-UC-038`、`E2E-UC-039` の seed 利用方針と操作 selector は維持でよい。
- 未確認事項: 実 backend が再開成功時に常に job run へ遷移すべきか、一覧に留まり feedback のみ出すことがあるかは未確認である。

### EF-004 出力管理

- 対象観点: `E2E-UC-023`、`E2E-UC-024`、`E2E-UC-025`、`E2E-UC-042`、`E2E-UC-043`、`E2E-UC-044`
- 観測事実: `scenario-test-implementation-result.md:107-115` には、output management は `App.svelte` 側で gateway が `null` のため scenario mock が画面へ接続されないとある。`output-management.spec.ts` は全件で `installScenarioWailsMocks` を使うが、候補行表示確認なしに `selectCandidate` へ進む。mock 自体は `scenario-wails-mocks.ts:343-378,468-505` で completed job と diff rows を固定している。
- test 前提: `test-design.csv:23-25,51-53` は候補行表示を前提にする。現行 test は `open()` 後に summary だけ待ち、candidate list の接続成立を待たない。
- fixture / mock: mock データは job 401, 402, 403 と diff row 1 件を決定的に返す。fixture 不足ではなく、mock 接続成立の観測不足が大きい。
- Page Object 妥当性: `OutputManagementPage` は candidate row、target game、path、diff row を test id で取得している。責務逸脱はない。
- assertion 妥当性: `E2E-UC-023` と `E2E-UC-024` は結果表示に寄り、候補一覧前提の崩れを診断しにくい。`E2E-UC-025` は spec が求める差分行 click を行わず、表示確認だけで終わるため仕様整合が弱い。`E2E-UC-042` から `E2E-UC-044` も候補行前提が崩れると assertion へ到達しない。
- 仕様整合性: `E2E-UC-025` は不一致である。spec は diff row click 後の維持確認を求めるが、test は click しない。`E2E-UC-023`、`E2E-UC-024`、`E2E-UC-042`、`E2E-UC-043`、`E2E-UC-044` は大枠一致だが、候補一覧表示の前提観測が不足する。
- 変更不要テスト範囲: output management Page Object の selector 方針、`job #401` から `job #403` を使う mock データの最小性、外部出力を fake 境界へ閉じる方針は維持でよい。
- 未確認事項: `open()` 完了時に summary だけ表示され candidate list が遅延表示される設計か、gateway 未接続時でも summary が見える設計かは未確認である。

### EF-005 翻訳段階

- 対象観点: `E2E-UC-045` から `E2E-UC-053` の一部
- 観測事実: `scenario-test-implementation-result.md:108` は translation management から現在段階へ進む操作で対象 job または action が見つからないと記録する。`job-run-shell.spec.ts:13-20` と `translation-phases.spec.ts:16-25` は共通で `management.openCurrentPhase(management.jobCard(jobText))` を使う。`TranslationJobManagementPage.openCurrentPhaseButton` は button 名 `現在の翻訳段階へ進む` に固定依存する。mock job 一覧は `scenario-wails-mocks.ts:434-442` の 5 件だけで、paused job は含まない。
- test 前提: `E2E-UC-045` から `E2E-UC-050` は job run shell 表示を、`E2E-UC-051` から `E2E-UC-053` は phase 画面表示を前提にする。現行 test は candidate card 可視や open action 可用性を明示確認せず直ちに click する。
- fixture / mock: phase 系 mock は `termSummary`、`personaSummary`、`bodySummary` を決定的に返し妥当である。失敗の中心は fixture 不足より、translation management から job run への導線観測不足にある可能性が高い。
- Page Object 妥当性: `TranslationJobManagementPage` と `JobRunShellPage`、`TranslationPhasePage` は test id と role button による薄い wrapper であり、責務逸脱はない。ただし button 名固定依存のため、card 表示不備と action 名差異が同じ失敗に見えやすい。
- assertion 妥当性: `E2E-UC-045` は spec `test-design.csv:26` の「開始 click で実行中から完了へ進む」ではなく、現在段階表示の確認に留まるため、観点 ID と証明対象がずれている。`E2E-UC-046` と `E2E-UC-047` も spec `27`、`28` の開始操作ではなく、失敗段階表示と再実行へ置き換わっている。`E2E-UC-048` から `E2E-UC-050` は shell 系 spec と概ね整合する。`E2E-UC-051` から `E2E-UC-053` は AI 設定不足の例外 spec と整合するが、入口の openCurrentPhase 成立確認が弱い。
- 仕様整合性: `E2E-UC-045` から `E2E-UC-047` は ID と期待値の対応が `test-design.csv:26-28` と不一致である。`E2E-UC-048` から `E2E-UC-053` は概ね整合する。
- 変更不要テスト範囲: `JobRunShellPage` と `TranslationPhasePage` の selector 分離、phase 例外系 mock の決定性、`E2E-UC-051` から `E2E-UC-053` の開始ボタン click 後に未開始維持を確認する観測点は維持でよい。
- 未確認事項: `E2E-UC-045` から `E2E-UC-047` を現行 shell 系振る舞いへ読み替える人間判断が既にあるかは未確認である。

## 仕様整合性まとめ

- `EF-001`: 期待値の後半が未実装で、仕様部分欠落である。
- `EF-002`: 仕様前提は正しいが、前提成立待機が不足する。
- `EF-003`: 仕様主要期待値の一部だけを見ている。
- `EF-004`: `E2E-UC-025` は spec 手順と不一致である。ほかは前提観測不足である。
- `EF-005`: `E2E-UC-045` から `E2E-UC-047` は spec ID と test 内容の対応が崩れている。`E2E-UC-048` から `E2E-UC-053` は大枠整合する。

## 変更不要テスト範囲

- `tests/system/support/` 配下の Page Object は、selector 解決と操作に責務を限定しており、この調査範囲では構造変更不要である。
- `tests/fixtures/master-persona/system-test-persona.json` は最小で代表性があり、`EF-002` の主因とは扱わない。
- `scripts/test/seed-system-test-db/main.go` の paused / failed / running job 固定は、`EF-003` の前提データとして妥当である。
- `tests/system/support/scenario-wails-mocks.ts` の master persona と phase 系 mock は決定的であり、外部境界隔離として維持でよい。

## 残り不足

- `EF-001` の実 backend が不正入力時に detail region をどう維持するかの UI 観測がない。
- `EF-002` の model select 活性化契機の UI 観測がない。
- `EF-003` の再開成功時に feedback のみか画面遷移も必須かの仕様根拠が不足する。
- `EF-004` の output management gateway 接続成立時系列を示す UI 証跡がない。
- `EF-005` の `E2E-UC-045` から `E2E-UC-047` を現行 shell 観点へ読み替えた判断記録がない。

## 残留リスク

- `EF-004` と `EF-005` は mock 接続不備と test 手順不一致が混在している可能性がある。
- `EF-002` は product failure でも test 待機不足でも同じ失敗文言になりやすい。
- `EF-003` は再開後遷移を必須とみなすか任意とみなすかで評価が変わる可能性がある。

## 推奨 next step

- 推奨 next step: 追加調査
- 理由: `EF-001`、`EF-002`、`EF-003` は test quality の不足を根拠付きで特定できたが、`EF-004` と `EF-005` は mock 接続成立と spec ID 対応の追加確認なしにテスト修正方針を確定できない。
