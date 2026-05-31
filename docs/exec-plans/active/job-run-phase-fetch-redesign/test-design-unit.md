# 単体テスト観点: job-run-phase-fetch-redesign

## 概要

- 対象: 取得・反映フロー作り直しと可否判断の責務再配置を、画面 E2E ではなく責務境界の単体テストで証明する観点を固定する。
- 判断: 反映取りこぼし防止・取得回数・開き直し時の旧取得破棄・可否導出の等価性は、frontend application 層（usecase / presenter）と backend service の公開振る舞いで観測できる。これらは画面操作 E2E では決定的に再現しにくいため単体テストへ分ける。
- 根拠: `./detail-spec-diff.md` の `term-translation-phase-REQ-007` / `REQ-008`、`persona-generation-phase-REQ-007` / `REQ-008`、`body-translation-phase-REQ-006` / `REQ-007`。`./design-diff.job-run-phase-fetch-redesign.md` の図2・図3・図4・図6・図7・図8。`./screen-design-diff.job-run.md` 差分6・差分7。
- 注意: 本書はテスト観点だけを固定する。テスト実装、検証コマンド、実装手順は含めない。

## 観点分類の対応

| 設計の観点（引き継ぎ） | 主な証明場所 | 単体観点 ID | シナリオ観点 ID |
| --- | --- | --- | --- |
| 反映取りこぼし防止（summary 独立反映・processingTarget 連番ガードの非対称是正） | frontend usecase | UT-ASYM-001〜003 | E2E-LTLE-001/002 |
| 取得回数・起動範囲（初回は表示中段階のみ最大2本、切り替え時に切り替え先取得、全段階同時取得しない） | frontend usecase / JobRunPage 取得起動 | UT-FETCH-001〜003 | （E2E では取得本数を直接観測しない） |
| 開き直し時の再取得と旧取得破棄（sequence 無効化） | frontend usecase | UT-REOPEN-001〜002 | E2E-LTLE-002 |
| 初回取得中ローディングレイヤーの操作排他（フェーズ画面全体を覆うオーバーレイによる全体操作排他） | viewModel/state（initialFetchDone）+ 画面 | UT-LOAD-001 | E2E-LTLE-004 |
| 可否導出の等価性（再配置後の frontend 導出が再配置前の backend 導出と一致） | frontend application 層 + backend 基準値 | UT-EQV-001〜010 | E2E-BTTL-002/003 |
| body 専用取得廃止後、段階要約の事実から出力可否をフロント導出 | frontend application 層 | UT-EQV-008〜010 | E2E-BTTL-002/003 |
| 表示規則（母数1以上で件数分・母数0だけ空状態、進捗母数と総件数は独立値） | viewModel/derived + 画面 | UT-DISP-001〜002 | E2E-LTLE-001/003 |

## 不足テスト一覧

| ID | 関連UC | 対象画面 | 不足観点 | 理由 | 追加候補 |
| --- | --- | --- | --- | --- | --- |
| UT-ASYM-001 | 処理対象を確認する | 単語翻訳 | 先行取得だけが完了した後に後発取得（より新しい sequence）が processingTarget を解決した場合、後発結果で一覧を更新し、空状態のまま残さない | term-REQ-007 の主軸（反映取りこぼし防止）。fail-test として作り直し前の非対称ガード（先行取得だけ完了で空状態残留）を検出する | term usecase の取得・反映で、後発取得結果が反映され items.length>0 になることを証明する単体テスト |
| UT-ASYM-002 | 処理対象を確認する | 単語翻訳 | 後発取得が sequence を進めた後に先発取得が遅れて解決しても、先発（古い sequence）の processingTarget では一覧を上書きせず取りこぼさない | seq-guard-asymmetry 再発防止。古い取得結果による上書きを許可しない | 古い sequence の解決を破棄し、保持済みの新しい一覧を維持することを証明する単体テスト |
| UT-ASYM-003 | 処理対象を確認する | 単語翻訳 | summary は processingTarget の反映可否と独立に反映され、processingTarget が未反映でも summary が反映される。逆に summary 未反映を理由に processingTarget の反映を止めない | term-REQ-007「進捗の要約は別経路で独立反映」。図2(after) の独立反映 | summary と processingTarget を別経路として、片方の反映が他方の反映可否に依存しないことを証明する単体テスト（term/persona/body 同型） |
| UT-FETCH-001 | 処理対象を確認する | 翻訳実行画面 | 翻訳ジョブを開いた初回表示で、表示中段階の取得だけが起動し（summary + processingTarget の最大2本）、非表示段階の取得が起動しない | translation-job-management-REQ-006「初回は表示中段階だけ取得」。図1(after)。全段階同時取得（最大9本）の廃止 | 初回起動で表示中段階の取得呼び出しだけが発火し、他段階の取得呼び出しが発火しないことを証明する単体テスト |
| UT-FETCH-002 | 処理対象を確認する | 翻訳実行画面 | 表示中段階を切り替えた時、切り替え先段階の取得が切り替え時点で起動する | translation-job-management-REQ-006「段階切り替え時に切り替え先段階を取得」。図1(after) | 段階切り替え操作で切り替え先段階の取得呼び出しだけが発火することを証明する単体テスト |
| UT-FETCH-003 | 処理対象を確認する | 翻訳実行画面 | 取得経路から可否判断 DTO（readiness 等）の取得が発火しない（可否はフロント導出のため取得物ではない） | term-REQ-008・図6(after)。backend は事実だけ返す。可否取得 DTO は取得物でなくなる | 取得起動で可否専用取得（GetTermTranslationNextPhaseReadiness 等）が呼ばれないことを証明する単体テスト |
| UT-REOPEN-001 | 処理対象を確認する | 翻訳実行画面 | 取得遅延中に画面を閉じて開き直した時、開き直し後の新 sequence で再取得し、閉じる前に開始した旧取得（古い sequence）の遅延応答を破棄する | translation-job-management-REQ-006「開き直し時に再取得・旧取得結果を破棄」。図3 | 旧 sequence の遅延応答が store に反映されず、新 sequence の結果が反映されることを証明する単体テスト |
| UT-REOPEN-002 | 処理対象を確認する | 翻訳実行画面 | 開き直し後、総件数1以上の段階で再取得結果が一覧へ反映され、空状態のまま残らない | 図3 と E2E-LTLE-002 の単体側裏付け | 再取得結果反映後に items.length>0 になることを証明する単体テスト |
| UT-LOAD-001 | 処理対象を確認する | 翻訳実行画面 | 初回取得（summary と processingTarget の両方）が完了するまで initialFetchDone=false を保ち、両取得完了で initialFetchDone=true になる | screen-design-diff 規則4〜7・図4。フェーズ画面全体を覆うオーバーレイによる全体操作排他の根拠状態 | initialFetchDone が両取得完了時にだけ true になることを viewModel/state で証明する単体テスト |
| UT-DISP-001 | 処理対象を確認する | 単語翻訳 | 処理対象一覧の総件数が1以上のとき、初回取得完了後に空状態を表示判定にしない（件数分表示の derived 判定が成立する） | screen-design-diff 規則1・規則2。currentProcessingTargetPageState の derived 評価（図5） | initialFetchDone=true かつ総件数1以上で空状態判定にならないことを証明する単体テスト |
| UT-DISP-002 | 処理対象を確認する | 単語翻訳 | 進捗母数（summary.aiTargetCount）と処理対象一覧の総件数は独立値として保持・表示され、一致しなくても表示が破綻しない | translation-job-management-REQ-006 の注意・screen-design-diff 規則3・図5 | 進捗母数と総件数が別フィールドとして保持され、不一致でも空状態へ落ちないことを証明する単体テスト |
| UT-EQV-001 | 処理対象を確認する | 単語翻訳 | 次段階開始可否の等価性: 終端ジョブでない・term フェーズ完了・確認済み件数>=対象件数のとき canStartNextPhase=true、それ以外は false で再配置前の backend 導出と一致 | term-REQ-008 等価性条件。基準は `internal/service/term_translation_phase_service.go` の `readinessFromState` | 同一事実入力に対し frontend 導出の可否が backend readinessFromState と一致することを table-driven で証明する単体テスト |
| UT-EQV-002 | 処理対象を確認する | 単語翻訳 | 次段階開始不可理由の等価性: 終端ジョブは terminal_job 相当、未完了は term_phase_incomplete 相当の理由を事実から区別して導出し再配置前と一致 | term-REQ-008。基準は readinessFromState の BlockedReason | 不可理由の区別が backend 基準と一致することを証明する単体テスト |
| UT-EQV-003 | 処理対象を確認する | 単語翻訳 | 開始可否（canStart）と不可理由の等価性: 終端ジョブ・実行中フェーズ存在・Ready 以外・実行設定未構成の各事実で再配置前の termTranslationStartBlockedReason と一致 | REQ-008 操作可否。基準は `termTranslationStartBlockedReason`・`termTranslationExecutionConfigured` | 4 つの不可理由分岐と活性を backend 基準と一致させる table-driven 単体テスト |
| UT-EQV-004 | 処理対象を確認する | 単語翻訳 | 一時停止可否（canPause）と理由の等価性: 終端ジョブ不可、実行中なら可、それ以外は phase_not_running 相当で再配置前と一致 | REQ-008。基準は `termTranslationPauseBlockedReason` | pause の活性・理由を backend 基準と一致させる単体テスト |
| UT-EQV-005 | 処理対象を確認する | 単語翻訳 | 再開可否（canResume）と理由の等価性: 終端ジョブ不可、paused または recoverable_fail なら可、それ以外は phase_not_resumable 相当で再配置前と一致 | REQ-008。基準は `termTranslationResumeBlockedReason` | resume の活性・理由を backend 基準と一致させる単体テスト |
| UT-EQV-006 | 処理対象を確認する | 単語翻訳 | 再試行可否（canRetry）と理由の等価性: 終端ジョブ不可、recoverable_fail なら可、それ以外は phase_not_retryable 相当で再配置前と一致 | REQ-008。基準は `termTranslationRetryBlockedReason` | retry の活性・理由を backend 基準と一致させる単体テスト |
| UT-EQV-007 | 処理対象を確認する | 単語翻訳 | 取り消し可否（canCancel）と理由の等価性: 再配置前の backend cancel 導出と同じ事実入力で一致 | REQ-008。基準は backend の cancel 可否導出 | cancel の活性・理由を backend 基準と一致させる単体テスト |
| UT-EQV-008 | 処理対象を確認する | NPC ペルソナ生成 | 本文翻訳段階開始可否と各操作可否の等価性: persona 段階データ事実状態（生成対象件数・生成完了件数・ペルソナ参照状態・フェーズ状態・実行設定など）から再配置前の backend 導出と一致 | persona-REQ-008 等価性条件。term と同型。基準は persona service の対応導出関数 | persona の遷移可否・操作可否を backend 基準と一致させる table-driven 単体テスト |
| UT-EQV-009 | 成果物出力を確認する | 本文翻訳 | 成果物出力確認可否（ready 相当）と不可理由の等価性: 段階要約の事実（completedFieldCount・statusConsistent・outputCount・失敗件数・エラー種別）から再配置前の専用取得 BodyTranslationOutputReadiness の ready/blockedReason と一致 | body-REQ-007 等価性条件・Q-005。専用取得廃止後も段階要約事実から同結果を導く。基準は廃止前の body 出力確認導出 | 段階要約事実から導出した出力可否・理由が廃止前 backend 基準と一致することを table-driven で証明する単体テスト |
| UT-EQV-010 | 成果物出力を確認する | 本文翻訳 | body の各操作可否（canStart/canPause/canResume/canRetry/canCancel/canCheckOutputReadiness）と理由の等価性: 段階要約の事実から再配置前の backend 導出と一致 | body-REQ-007 操作可否・Q-004。term/persona と同型。基準は body service の操作可否導出関数 | body の操作可否を backend 基準と一致させる table-driven 単体テスト |

## fail-test 戦略

- 再現可能な観点を fail-test として固定する。具体的には、反映取りこぼし是正（UT-ASYM-001〜003、E2E-LTLE-001/002）、取得回数・起動範囲（UT-FETCH-001〜003）、開き直し時の旧取得破棄（UT-REOPEN-001〜002）、可否導出の等価性（UT-EQV-001〜010）を、作り直し前のコードでは失敗し作り直し後に成功する観点として置く。
- 等価性テストは backend の現状導出（`internal/service/term_translation_phase_service.go` の `readinessFromState`・`termTranslationStartBlockedReason` 等、persona/body の同型 service）を基準値にし、frontend 導出が同じ事実入力で同じ可否・理由を返すことを証明する。基準値は再配置で backend から削除されるため、等価性テストの実装時に基準値（期待値）を仕様根拠から固定し直す。
- モック E2E の限界を明示する。モック E2E（scenario-wails-mocks）は応答が同期解決のため、実機で観測された約15秒の IPC 飽和遅延を再現しない。実機固有症状（15秒遅延・遅延中の unmount による初回0件）の検出は本テスト設計の対象外とし、実装後ブラウザ確認（browser_confirmation）で扱う前提とする。
- 初回取得中ローディングレイヤーの操作排他（E2E-LTLE-004。ローディングはフェーズ画面全体を覆うオーバーレイで、検索・ページ操作・処理対象行展開に加え操作行全体を排他する）は、モックの同期解決によりオーバーレイ出現が一瞬になりうる限界がある。確実な検出は initialFetchDone 状態の単体テスト（UT-LOAD-001）を主にし、E2E は補助とする。

## selector 整合

- 残置 E2E（E2E-LTLE-001/002/003）は `tests/system/fix-lucien-target-list-empty.spec.ts` と page object（`tests/system/support/translation-phase-pages.ts`）が `term-translation-phase-processing-target-empty` / `-total` / `-search-input` / `-row` を参照する。screen-design-diff 差分5 の固定 selector と一致する。
- 初回取得中ローディングレイヤー（`<phase-prefix>-processing-target-loading`）は新規 selector であり現行実装に対応要素がない。E2E-LTLE-004 と各段階の同型観点が依存する。`./data-testid-gaps.md` に記録する。
- persona / body の同型 selector（`persona-generation-phase-processing-target-*`、`body-translation-phase-processing-target-*`）は page object の prefix 切り替えで参照する。

## 根拠

- `source`: `./detail-spec-diff.md`、`./screen-design-diff.job-run.md`、`./design-diff.job-run-phase-fetch-redesign.md`、`./plan.md`（不足テスト節）。
- `source`: `tests/system/fix-lucien-target-list-empty.spec.ts`、`tests/system/support/translation-phase-pages.ts`、`tests/system/support/scenario-wails-mocks.ts`（termZeroAITargetJobId）。
- `source`: `internal/service/term_translation_phase_service.go`（`readinessFromState`・`termTranslationStartBlockedReason`・`termTranslationPauseBlockedReason`・`termTranslationResumeBlockedReason`・`termTranslationRetryBlockedReason`）を等価性の基準値とする。
- `source`: `docs/exec-plans/completed/2026-05-30-fix-lucien-target-list-empty/test-design.csv`（残置 E2E 観点の由来）。
- `validation`: 文書整形と CSV header 整合の確認のみ。実機検証は実装後ブラウザ確認で扱う。
