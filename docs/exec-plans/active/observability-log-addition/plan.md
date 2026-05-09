# コード全体 観測ログ追加 plan

## 状態

- task-id: `observability-log-addition`
- 依頼要約: コード全体を対象に、原因候補を分離する恒久的な観測ログを導入する。
- 現在成果物: `レビュー通過根拠`
- 次成果物: `close`
- 停止条件: なし。

## 決定

- 本 task は新規実装レーンで扱う。
- 観測ログ仕様は `docs/observability-logging.md` に従う。
- backend は `slog` を使い、frontend は `pino` を使う。
- trace ID、全 command の start / finish log、全文入力、巨大 payload、secret は追加しない。
- 対象はコード全体を一括変更せず、原因分離価値が高い境界ごとに分ける。

## 理由

- 観測ログは実行後に消える中間状態、分岐理由、外部境界の失敗分類を残すために使う。
- コード全体への一括追加は、同種ログの大量発生と責務境界の混在を起こす可能性がある。
- `observability_implementer` は完成済み実装成果物と変更ファイルが揃った後にだけ着手できる。

## 影響

- backend、frontend、統合境界の各実装範囲を分けて固定する必要がある。
- UI 表示、画面文言、layout、style は観測ログ追加の対象にしない。
- 既存の `observability-logger-lightweight` の方針は、ログ基盤を重くしない判断材料として参照する。

## 成果物依存表

| 成果物ID | 状態 | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- | --- |
| `task 枠` | 完了 | `implement_lane` | `[]` | なし |
| `scenario_candidates` | 完了 | シナリオ候補 生成 agent | `task 枠` | `scenario_actor_goal_generator`, `scenario_lifecycle_generator`, `scenario_state_transition_generator`, `scenario_failure_generator`, `scenario_external_integration_generator`, `scenario_operation_audit_generator` |
| `シナリオ設計` | 完了 | `designer` | `scenario_candidates` | `designer` |
| `UI設計` | 省略 | `designer` | `シナリオ設計` | なし |
| `設計差分図` | 完了 | `diagrammer` | `シナリオ設計`, `UI設計?` | `diagrammer` |
| `人間設計レビュー` | 完了 | 人間 | `シナリオ設計`, `UI設計?`, `設計差分図` | 人間 |
| `実装範囲` | 完了 | `designer` | `人間設計レビュー` | `designer` |
| `実装引き継ぎ入力` | 完了 | `implement_lane` | `実装範囲` | なし |
| `frontend 実装` | 完了 | `frontend_implementer` / `implement-frontend` | `実装引き継ぎ入力` | `frontend_implementer` |
| `観測ログ追加` | 完了 | `observability_implementer` / `observability-implementer` | `実装引き継ぎ入力` | `observability_implementer` |
| `単体テスト` | 完了 | `implementation_unit_tester` | `観測ログ追加` | `implementation_unit_tester` |
| `シナリオテスト` | 完了 | `implementation_scenario_tester` | `観測ログ追加` | `implementation_scenario_tester` |
| `最終検証` | 完了 | `implement_lane` | `観測ログ追加`, `単体テスト`, `シナリオテスト` | なし |
| `実装後ブラウザ確認` | 完了 | `browser_confirmation` | `最終検証` | `browser_confirmation` |
| `レビュー通過根拠` | 完了 | `implement_lane` | `最終検証`, `実装後ブラウザ確認` | `review_behavior`, `review_contract`, `review_trust_boundary`, `review_state_invariant`, `review_responsibility_boundary` |

## 未決事項

- どの backend 境界を最初の対象にするか。
- どの frontend runtime event を最初の対象にするか。
- UI を伴う観測確認が必要か。
- 既存の `observability-logger-lightweight` を完了または統合するか。

## scenario_candidates

- `scenario-candidates.actor-goal.md`: アクター目的観点の候補。
- `scenario-candidates.lifecycle.md`: lifecycle 観点の候補。
- `scenario-candidates.state-transition.md`: 状態遷移観点の候補。
- `scenario-candidates.failure.md`: 異常系観点の候補。
- `scenario-candidates.external-integration.md`: 外部連携観点の候補。
- `scenario-candidates.operation-audit.md`: 運用・監査観点の候補。

## シナリオ設計

- `scenario-design.md`: 6 観点を 6 シナリオへ統合済み。
- `scenario-design.candidate-coverage.json`: 候補統合と競合解消を記録済み。
- `scenario-design.requirement-coverage.json`: 詳細要求タイプを記録済み。
- `scenario-design.requirement-gate.md`: pass。`finding_count` は 0、`question_count` は 0。
- `scenario-design.questions.md`: 未回答質問なし。
- UI設計: 画面表示、画面文言、layout、style を変更しないため省略。

## 設計差分図

- `design-diff.component.puml`: 追加予定の観測点と追加しない接続を示すコンポーネント図。
- `design-diff.sequence.puml`: 観測ログ出力の代表的な流れを示すシーケンス図。
- `design-diff.md`: 根拠参照、追加予定、削除予定、変更しない接続先、検証結果、未決事項。
- 検証: PlantUML syntax check は 2 件とも成功。
- 差分確認: `git diff --check -- docs/exec-plans/active/observability-log-addition` は通過。

## 停止理由

- なし。

## 人間設計レビュー

- `human-design-review.md`: 2026-05-09 に人間が `approve` と回答した。
- 承認対象: `scenario-design.md`、`scenario-design.candidate-coverage.json`、`scenario-design.requirement-coverage.json`、`scenario-design.requirement-gate.md`、`scenario-design.questions.md`、`design-diff.component.puml`、`design-diff.sequence.puml`、`design-diff.md`。

## 実装範囲

- `implementation-scope.md`: 承認済み実装範囲を作成済み。
- `wave-1`: `OBS-FE-001` と `OBS-BE-001`。
- `OBS-FE-001`: frontend runtime event の `pino` log 配線。
- `OBS-BE-001`: backend 状態遷移 log。
- backend logger 基盤追加: 不要。既存 `slog` JSON handler を使う。
- frontend logger 基盤追加: 不要。既存 `pino` diagnostic logger を使う。

## 実装進行

- `implementation-progress.md`: `wave-1` から `wave-4` の実装結果、影響範囲修正、最終検証結果を記録済み。
- `OBS-FE-001`: 実装完了。影響範囲修正後、`frontend-local` は pass。
- `OBS-BE-001`: 実装完了。`backend-local` は pass。
- `wave-2`: `OBS-BE-002`、`OBS-BE-003`、`OBS-UNIT-FE-001` は完了。
- `wave-3`: `OBS-BE-004` は完了。
- `wave-4`: `OBS-UNIT-BE-001` と `OBS-SCN-BE-001` は完了。
- coverage 影響範囲修正: service 層と usecase 層の Sonar maintainability HIGH issue は解消済み。
- 最終検証: `backend-local`、`frontend-local`、`coverage` は pass。

## 実装後ブラウザ確認

- `browser-confirmation/2026-05-09-browser-confirmation.md`: 実装後ブラウザ確認を記録済み。
- 初期画面: blank ではない。runtime error overlay は確認していない。
- マスター辞書画面: 到達可能。辞書一覧と詳細領域を表示した。
- console: `runtime_event_subscribe` を含む frontend log を確認した。
- network: 初期読込の 200 応答を確認した。
- 未確認: master dictionary runtime event の詳細イベント列、更新操作、削除操作。

## レビュー指摘修正

- `behavior-001`: completed runtime event の malformed payload を `dropped` / `payload_parse_failed` に変更済み。
- `behavior-002`: phase state log の固定前提 `before_state` を除去済み。
- `contract-001`: UI 契約外の影響範囲修正は停止または人間承認へ戻す契約に修正済み。
- 修正後検証: `backend-local`、`frontend-local`、`coverage` は pass。
- `state-invariant-001`: completed runtime event の空 object と未知 key object を `invalid_payload` として dropped に変更済み。
- `behavior-002` 再修正: service の command read model から `BeforePhaseState` / `AfterPhaseState` を usecase log へ渡すように変更済み。
- coverage 影響範囲修正: 本文翻訳 service の複雑度指摘を helper 分離で解消済み。
- 再修正後検証: `backend-local`、`frontend-local`、`coverage` は pass。Sonar maintainability HIGH issue は `0 <= 0`。
- `trust-boundary-001`: 承認済み範囲外の影響範囲修正一般許可を、各 skill と agent TOML で限定条件付きに修正済み。
- `contract-002`: review agent 起動入力の検証証跡契約と非対象規約の衝突を解消済み。

## レビュー通過根拠

- `reviewback.behavior.yaml`: `no_issue`、`must_fix_open: false`、`max_level: none`。
- `reviewback.contract.yaml`: `no_issue`、`must_fix_open: false`、`max_level: none`。
- `reviewback.trust-boundary.yaml`: `no_issue`、`must_fix_open: false`、`max_level: none`。
- `reviewback.state-invariant.yaml`: `no_issue`、`must_fix_open: false`、`max_level: none`。
- `reviewback.responsibility-boundary.yaml`: `issues_open`、`must_fix_open: false`、`max_level: minor`。
- `implementation_action`: `close`。
- 残留 minor: provider settings service に横断 provider log helper と persona generation 固有変換が同居している。修正必須ではない。
