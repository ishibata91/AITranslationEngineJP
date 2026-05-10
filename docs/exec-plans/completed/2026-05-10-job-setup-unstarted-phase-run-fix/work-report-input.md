# work report 入力

## task

`2026-05-10-job-setup-unstarted-phase-run-fix`

## 目的

ジョブセットアップ完了時点で、既存の `JOB_PHASE_RUN` に未開始 phase row を作る。

新しい中間 table は作らない。

`pending` は未開始 state として使わない。

## 変更概要

- Job Setup は、作成済み job に紐づく未開始 `JOB_PHASE_RUN` を 4 件作る。
- 未開始 row は `translation`、`term_translation`、`persona_generation`、`body_translation` を表す。
- `translation`、`term_translation`、`body_translation` は `idle_ready` を使う。
- `persona_generation` は `not_started` を使う。
- phase start は先置き済み row を `running` へ更新する。
- 旧 DB の phase run 不在 fallback は維持する。
- delete guard は `pending` を unsafe のまま扱い、`idle_ready` と `not_started` は unsafe にしない。
- Job Management は未開始 row を実行中 phase として誤判定しない。
- 未開始 row の開始昇格は期待 state 条件付きで更新する。
- 昇格競合時は `ErrConflict` とし、provider 実行や保存副作用へ進まない。

## 成果物

- `human-observation.md`
- `investigation.md`
- `cause-sequence.puml`
- `cause-sequence.svg`
- `fix-execution-input.md`
- `implementation-evidence.md`
- `regression-test-evidence.md`
- `browser-confirmation-input.md`
- `browser-confirmation.md`
- `review-input.md`
- `reviewback.behavior.yaml`
- `reviewback.contract.yaml`
- `reviewback.trust-boundary.yaml`
- `reviewback.responsibility-boundary.yaml`
- `reviewback.state-invariant.yaml`
- `reviewback.state-invariant-rereview.yaml`
- `review-pass-evidence.md`

## 検証

- `go test ./internal/repository ./internal/service`: 成功
- `python3 scripts/harness/run.py --suite backend-local`: 成功
- `python3 scripts/harness/run.py --suite coverage`: 成功
- coverage summary: `70.8%`
- `agent-browser doctor --offline --quick`: 成功
- ブラウザ確認: 翻訳管理で job が `開始待ち` と表示され、summary 取得失敗と `pending` 表示が無いことを確認した。
- PlantUML 構文確認: 成功
- PlantUML SVG 生成: 成功

## レビュー

- 挙動正しさ: `no_issue`
- 契約・互換性: `no_issue`
- 権限・信頼境界: `no_issue`
- 責務境界: `no_issue`
- 状態・データ不変条件: 初回 `major`、追加修正後 `no_issue`

## 既知の残留

実 DB で同一 job の同一 phase を同時開始する統合再現は未実施。

追加した単体テストは、競合を注入して副作用停止を確認している。

## report 出力先

`work_history/runs/2026-05-10-job-setup-unstarted-phase-run-fix-run/`
