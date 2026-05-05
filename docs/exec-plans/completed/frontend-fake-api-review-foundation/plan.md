# フロントエンド fakeAPI レビュー基盤

## 状態

- `task_id`: `frontend-fake-api-review-foundation`
- `lane`: `implement-lane`
- `task_mode`: `new-feature`
- `current_artifact`: `完了`
- `source`: `tasks/usecases/frontend-fake-api-review-foundation.yaml`

## 依頼要約

起動モードで、フロントエンドの API 接続先を fakeAPI に切り替える。
実フロントを Wails バインディングなしで起動し、レビュー用の状態をモックデータで再現できる基盤を作る。

## 成果物DAG

- `task 枠`: 完了
- `scenario_candidates`: 完了
- `シナリオ設計`: 完了
- `UI設計`: 省略
- `人間設計レビュー`: 完了
- `実装範囲`: 完了
- `実装引き継ぎ入力`: 完了
- `frontend 実装`: 完了
- `frontend 実装後人間レビュー`: 省略
- `単体テスト`: 完了
- `シナリオテスト`: 完了
- `最終検証`: 完了
- `レビュー通過根拠`: 完了
- `正本化判断`: 完了
- `作業レポート入力`: 完了
- `作業計画完了移動`: 完了

## 境界

- 変更対象候補は、フロントエンドの composition root、ゲートウェイ DI、fakeAPI データ、レビュー起動モードに限定する。
- fakeAPI は本番 API、永続化、本番初期状態に混入させない。
- Wails バインディングなしのレビュー起動では、実画面と presenter / usecase 経路を使う。
- 状態パターンの既定値はレビュー起動条件で指定し、実画面確認では URL パラメータで上書きできる。
- 状態パターン指定は fakeAPI 起動中だけ有効にし、本番起動では無視する。
- バックエンド実装は原則対象外とし、バックエンドが必要な場合は実装範囲で別成果物として分ける。
- レビュー専用 UI、状態パターン選択 UI、表示文言設計は作らない。

## 検証

- `python3 scripts/harness/run.py --suite frontend-local`
- fakeAPI 起動モードの局所テスト
- `agent-browser open http://localhost:34115/?fakeScenario=<状態パターンID>` からの状態パターン確認

## 人間レビュー

- シナリオ設計: 承認済み
- 承認者: human
- 承認指示: 「設計OK，先へ進んで」
- UI 設計: 省略承認済み

## 検証結果

- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/completed/frontend-fake-api-review-foundation/scenario-design.md --coverage docs/exec-plans/completed/frontend-fake-api-review-foundation/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/completed/frontend-fake-api-review-foundation/scenario-design.candidate-coverage.json --json`: pass
- `finding_count`: 0
- `question_count`: 0
- `npm --prefix frontend run test -- src/controller/review-fake-api`: pass, 2 files, 8 tests
- `npm --prefix frontend run test -- src/ui`: pass, 7 files, 104 tests
- `python3 scripts/harness/run.py --suite frontend-local`: pass, frontend lint harness pass, frontend test harness pass
- `agent-browser open http://localhost:34115/?fakeApi=1&fakeScenario=success#provider-settings`: pass
- `agent-browser open http://localhost:34115/?fakeApi=1&fakeScenario=error#provider-settings`: pass
- `agent-browser open http://localhost:34115/?fakeApi=1&fakeScenario=config-missing#provider-settings`: pass

## レビュー集約

- `implementation_action`: `close`
- `reviewback.behavior.yaml`: `no_issue`
- `reviewback.contract.yaml`: `no_issue`
- `reviewback.trust-boundary.yaml`: `no_issue`
- `reviewback.state-invariant.yaml`: `no_issue`
- `reviewback.responsibility-boundary.yaml`: `no_issue`

## 正本化判断

- docs 正本化: 省略
- 理由: 本 task は fakeAPI レビュー基盤の実装であり、詳細仕様正本へ昇格する human 承認済み恒久仕様は未指定である。

## 作業レポート

- `work_history/runs/2026-05-05-frontend-fake-api-review-foundation-run/README.md`
- `work_history/runs/2026-05-05-frontend-fake-api-review-foundation-run/codex.md`
- `work_history/runs/2026-05-05-frontend-fake-api-review-foundation-run/transcript_refs.json`
- `work_history/runs/2026-05-05-frontend-fake-api-review-foundation-run/workflow-improvement-log.jsonl`
