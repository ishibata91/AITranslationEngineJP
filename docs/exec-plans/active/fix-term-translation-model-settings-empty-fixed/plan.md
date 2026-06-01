# fix-term-translation-model-settings-empty-fixed

## 依頼要約

ジョブ実行画面の単語翻訳フェーズで、モデル設定パネルが空欄であるにもかかわらず「固定済み」と表示される。開始ボタンを押すと「実行設定が未構成のため開始できません」というメッセージが出て開始できない。中断も「フェーズが実行中ではありません」と返る。設定変更も不可。投資調査と恒久修正のため、観測ログ駆動で確定原因と修正方針を固定する。

## 分岐元

- 分岐元 branch: master
- 分岐元 commit: 917737e149c89c4c4458cf98ce8cada64e699ecb

## 作業 branch

- claude/fix-term-translation-model-settings-empty-fixed

## 想定 Y/N 評価

| 想定 | Y/N | 根拠 | 参照 |
| --- | --- | --- | --- |
| 仕様変更または仕様追加がある | N | 単語翻訳画面仕様には「AI 設定不足: 認証不足またはモデル未選択の警告を表示する」が既に定義済み。実装がこの分岐を満たしていない。 | docs/screen-design/screens/term-translation-phase.md:107-113 |
| 画面変更がある | N | 仕様で定義済みの状態（AI 設定不足）に正しく分岐させるだけで、画面構造は変えない。 | docs/screen-design/screens/term-translation-phase.md:91-129 |
| 内部構造変更がある | Y | panel または presenter の判定ロジック（modelLabel 空判定、aiSettingsBlockedReason、isExecutionConfigured）を直す必要がある。 | frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte:61-73, frontend/src/application/presenter/term-translation-phase/term-translation-phase.presenter.ts:346-350,554-557 |
| 画面の表示変更がある | N | 仕様どおりの分岐結果として state pill が "固定済み" → "設定未完了" に切り替わるのは仕様準拠であり、新規 layout・文言・style 追加ではない。 | docs/screen-design/screens/term-translation-phase.md:107-113 |
| frontend ロジック変更がある | Y | panel の derived state または presenter の派生ロジックを修正する。 | TermTranslationPhasePanel.svelte:61-73 |
| backend 変更がある | 判断保留 | backend が `execution.model` を空文字で返すのか null で返すのかを実画面で観測してから確定する。 | fix_decider の観測ログ検証で固定 |
| frontend と backend を接続する | N | 既存 gateway 経路を変えない。 | - |
| 実装済み責務を独立に証明したい | Y | panel の判定ロジック、presenter の `isExecutionConfigured` を単体テストで証明する。 | - |
| 実行時にしか確定しない値または原因分離が要る分岐がある | Y | backend 応答の `execution` の null / 空文字状態と、`presenter` 出力との対応を実行時に観測する必要がある。 | - |

「仕様変更または仕様追加がある」が N のため、本モジュールを継続する。

## 人間観測記録

- 対象ジョブ: ジョブ#3
- 対象フェーズ: 単語翻訳フェーズ
- 観測 1: 単語翻訳フェーズのモデル設定パネルが空欄である（モデルが表示されていない）。
- 観測 2: 上記の空欄状態であるにもかかわらず、状態 pill に「固定済み」と表示されている。
- 観測 3: モデル設定パネルから設定変更ができない（操作しても確定しない、または選択肢が出ない状態）。
- 観測 4: 「開始」操作で「実行設定が未構成のため開始できません。」というメッセージが表示される。
- 観測 5: 「中断」操作で「フェーズが実行中ではありません。」というメッセージが表示される。
- 期待との差分: モデルが空欄の場合は仕様の「AI 設定不足: モデル未選択の警告を表示する」状態（state pill: 設定未完了相当）になり、ユーザーがモデルを選択して固定する経路が確保されているべきである。

## Wails 接続対象

- 起動 command: `npm run dev:wails:run`
- 接続先: `http://localhost:34115`
- 接続先単一化: 同一接続先で fix_decider が画面再現確認する。

## 停止判断（investigation-module）

人間レビューで ER 図正本との乖離を確認した結果、本不具合は対症療法では解決できず仕様変更を伴うと判明したため、investigation-module を停止して design-module へ迂回する。

- 確認事実 1: ER 図正本 `docs/diagrams/er/combined-data-model-er.puml` は、フェーズ別 AI 設定を `JOB_PHASE_RUN`（id, ai_provider, model_name, execution_mode, credential_ref ほか）に保持する設計である。`TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT` という独立テーブルは ER 図に定義されていない。
- 確認事実 2: ER 仕様 `docs/er.md:67-68` は「`Ready` job には `JOB_PHASE_RUN` を事前作成しない。フェーズ開始が許可された時だけ作成する」と定める。よって ER 仕様としては Ready 状態のフェーズに「保存済み AI 設定」を保持する場所が無い。
- 確認事実 3: 実装は ER 図に無い `TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT` テーブルを追加し、`SaveAISettings` が開始前のユーザー保存値と実行時固定値を同居して書き込む構造になっている。空文字で record を作成する経路が「設定なし」状態を歪めている。
- 人間判断: design-module へ迂回して ER 仕様変更を含めて検討する。Ready 状態のジョブで AI 設定を保存・編集できる場所を ER 上で定義し直すか、provider-settings 正本からの解決へリファクタするかは design-module で判断する。
- 影響範囲: 単語翻訳のみに見える不具合だが、persona-generation / body-translation も同じ snapshot テーブルを共有する。設計判断の波及範囲は design-module で固定する。

## 後続モジュールへの引き継ぎ

- 入口: design-module。
- 引き継ぐ事実: 本 plan.md の「想定 Y/N 評価」「人間観測記録」、`fix-decision.md`（snapshot 直書き欠陥の確定原因まで）、`missing-usecases.md`、`missing-tests.md`、`data-testid-gaps.md`。
- 検討対象: ER 仕様の修正、または provider-settings 正本からの解決方式、フェーズ AI 設定の永続化場所、SaveAISettings の保存先、`applyTermTranslationRuntimeSnapshot` の責務再定義。

## 想定 Y/N 評価（design-module 入口で再評価）

investigation-module で固定した想定は backend 側を判断保留としていたため、ER 仕様変更を含む設計検討が必要と確定した時点で再評価する。

| 想定 | Y/N | 根拠 | 参照 |
| --- | --- | --- | --- |
| 仕様変更または仕様追加がある | Y | ER 仕様正本（`docs/er.md:23-26`、`docs/diagrams/er/combined-data-model-er.puml`）と現状実装の乖離を解消するため、ER 修正方向（snapshot テーブルの ER 正式化／Ready 状態の AI 設定保持場所の追加／provider-settings 正本への一本化、のいずれか）を選び、対応する詳細仕様差分を固定する必要がある。 | `docs/er.md`, `docs/detail-specs/term-translation-phase.md` |
| 画面変更がある | N | 画面設計正本 `docs/screen-design/screens/term-translation-phase.md:107-113` には「AI 設定不足」「AI 設定準備済み」「実行中」の状態別表示が既に定義済み。本 task は仕様どおりに分岐するための内部構造の修正であり、画面領域・操作・状態カードの構造変更は行わない。 | `docs/screen-design/screens/term-translation-phase.md` |
| 内部構造変更がある | Y | snapshot / JOB_PHASE_RUN / provider-settings の責務再定義、SaveAISettings の保存先変更、applyTermTranslationRuntimeSnapshot 相当の責務再定義を含むため、内部構造の差分図と設計判断を固定する必要がある。 | `internal/service/term_translation_phase_service.go:1467-1510`, `internal/service/provider_execution_snapshot.go:105-178` |
| 画面の表示変更がある | N | 状態 pill の文言・layout・style・表示構造は仕様どおりで、新規 UI は要らない。 | 同上（画面設計正本） |
| frontend ロジック変更がある | Y弱 | backend 応答の null / 値の表現変更に追随して、presenter の正規化、panel の判定経路に微修正が入る可能性。範囲は実装範囲で再判断する。 | `frontend/src/application/presenter/term-translation-phase/term-translation-phase.presenter.ts`, `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte` |
| backend 変更がある | Y | 永続化責務、SaveAISettings の保存先、表示用 query、phase 開始時の固定値転写経路の修正を含む。 | `internal/service/term_translation_phase_service.go`, `internal/service/provider_execution_snapshot.go`, `internal/repository/` 配下 |
| frontend と backend を接続する | Y弱 | 「設定なし」状態の境界 DTO 表現（execution が null か、field が null か、または record 自体が無いか）を確定する必要があるため、境界の DTO 変更が入る可能性。 | gateway / wails bridge 経路 |
| 実装済み責務を独立に証明したい | Y | snapshot 責務、phase AI 設定の永続化責務、表示経路の責務をそれぞれ単体テストで証明する。 | - |
| 実行時にしか確定しない値または原因分離が要る分岐がある | N | investigation-module で確定原因は固定済み（snapshot 直書きの仕様逸脱）。設計確定後の経路は静的に検証できる。 | `fix-decision.md` |

### decision table 適用

| 想定 | 詳細仕様差分 | 画面設計差分 | 設計差分図 |
| --- | --- | --- | --- |
| 仕様変更または仕様追加がある: Y | 要 | - | - |
| 画面変更がある: N | - | -（省略） | -（推奨対象外） |
| 内部構造変更がある: Y | - | - | 要 |

- 作成対象: `詳細仕様差分`、`設計差分図`、`実装範囲`、`テスト設計`、`人間設計レビュー`。
- 省略対象: `画面設計差分`。省略理由: 画面設計正本に既存の状態別表示で覆える。仕様どおりの分岐を満たすための内部構造修正であり、画面構造は変えない。

## implementation-module decision table 適用

| 想定 | 結果 |
| --- | --- |
| frontend ロジック変更がある: Y弱 | frontend ロジック実装: 要 |
| backend 変更がある: Y | backend 実装: 要、単体テスト: 要 |
| frontend と backend を接続する: Y弱 | 統合境界実装: 要、シナリオテスト: 要 |
| 実装済み責務を独立に証明したい: Y | 単体テスト: 要 |
| 実行時にしか確定しない値: N | 観測ログ追加: 不要（省略理由: 実行時分岐は無く、確定原因と修正経路が静的に検証可能） |

要 artifact: backend 実装、統合境界実装、frontend ロジック実装、単体テスト、シナリオテスト、最終検証。

## 実装引き継ぎ入力

- 実装範囲: `implementation-scope.md` の 9 handoff / 6 wave に従う。
- 順序: wave-1 (`H-MIG-ER`) → wave-2 (`H-BE-TERM` / `H-BE-PERSONA` / `H-BE-BODY` 並列) → wave-3 (`H-INT-PHASE-AI-SETTINGS`) → wave-4 (`H-FE-PRESENTER`) → wave-5 (`H-FE-SELECTOR`) → wave-6 (`H-TEST-UNIT` / `H-TEST-SCENARIO` 並列)。
- 画面表示変更なし、storybook-module 経由不要、合意済み frontend 保護は不要。
- 各 handoff は別 agent で起動する。`H-MIG-ER`/`H-BE-*` は `backend_implementer`、`H-INT-*` は `integration_implementer`、`H-FE-*` は `frontend_implementer`、`H-TEST-UNIT` は `implementation_tester`（単体）、`H-TEST-SCENARIO` は `implementation_tester`（シナリオ）。




## 最終検証結果（2026-06-02）

- 実装済み 6 wave 全 handoff 完了。
- backend-local / frontend-local harness 通過済み（前回 task 内で記録済み、本 finalize 時点での再検証は別 task 後の merge 後検証で行う）。
- 実画面確認:
  - 単語翻訳フェーズの AI 設定 pill が「設定未完了」と表示されるべきケースで仕様どおり分岐するようになった。
  - 開始ボタンが「実行設定が未構成のため開始できません」と表示され、AI 設定固定前は disabled で正しく動作する。
- 残課題は別 task `fix-phase-ai-model-list-empty`（fake モデル一覧取得経路）として切り出し、本 task 完了後に同 branch 上で実施した。
- さらに残った課題（モデル選択後の pill 更新と開始ボタン有効化）も別 task として切り出す予定。

## 正本化判断

- 仕様変更または仕様追加: なし（想定 Y/N 評価でも N）。
- 詳細仕様正本反映: 不要。
- 人間承認状態: 不要（仕様変更なし）。

