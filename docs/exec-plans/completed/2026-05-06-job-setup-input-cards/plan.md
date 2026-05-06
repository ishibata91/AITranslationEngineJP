# 2026-05-06 job setup input cards plan

## 状態

- `task_id`: `2026-05-06-job-setup-input-cards`
- `workflow_state`: `implementation-reviewed`
- `lane_owner`: `implement_lane`
- `task_mode`: Job Setup UX 改善
- `source_request`: Job Setup の入力データを、プルダウンではなく削除ボタン付きカードにする。

## 目的

Job Setup で登録済み input の存在に気づけるようにする。
Data Load で登録だけして job 未作成の input を、Job Setup から見つけて選択または削除できるようにする。

## 固定判断

- 入力データ選択はプルダウンではなくカード一覧にする。
- カードには input 名、出自、登録日時、レコード数を表示する。
- カード選択で Job Setup の対象 input を切り替える。
- 削除ボタンは job 未作成 input だけを削除する。
- 既存 job が参照している input は Job Setup の入力候補に表示しない。
- Job Management の job 削除は input を残す。今回の input 削除と混ぜない。
- input 削除は `X_EDIT_EXTRACTED_DATA` の親行削除を正本にし、関連行は DB の `ON DELETE CASCADE` に任せる。
- 既存 DB 内の input、job、cache、出力データは migration で完全破棄してよい。
- input 削除中は対象カードを操作不能にし、カード内に `削除中...` を表示する。
- 削除成功後は options 全再読込ではなく、画面状態から対象候補だけ除去する。

## 成果物DAG

- `task 枠`: completed。この `plan.md`。
- `frontend 実装`: completed。カード単位の削除中表示と局所更新へ変更した。
- `backend 実装`: completed。canonical ER tables を cascade reset し、input 削除を親行削除へ変更した。
- `単体テスト`: completed。cascade 削除、削除拒否、削除中表示、局所更新を証明した。
- `最終検証`: completed。frontend-local、backend-local、all を通過した。
- `レビュー通過根拠`: completed。5 観点レビューはすべて `must_fix_open: false`。

## 実装入力

- `frontend_target`: `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`
- `frontend_test`: `frontend/src/ui/screens/translation-job-setup/JobSetupPage.test.ts`
- `backend_target`: canonical ER migration、Job Setup の input 削除境界、SQLite repository。
- `frontend_target_2`: input 削除中状態、カード単位の無効化、options 局所更新。
- `guard`: `TRANSLATION_JOB` が参照している `X_EDIT_EXTRACTED_DATA` は入力候補に返さず、削除もしない。既存 DB 内の canonical ER データは reset で破棄してよい。

## 検証予定

- `python3 scripts/harness/run.py --suite frontend-local`
- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite all`

## 過去検証結果

- `python3 scripts/harness/run.py --suite frontend-local`: pass。57 files / 486 tests passed。
- `python3 scripts/harness/run.py --suite backend-local`: pass。backend lint と backend test passed。
- `python3 scripts/harness/run.py --suite all`: 初回 fail。Sonar maintainability HIGH が `DeleteInputSource` の認知的複雑度を検出した。
- `reviewback.contract.yaml`: 初回 fail。`existingJob` 応答 field の意味互換性破壊を検出した。
- `reviewback.state_invariant.yaml`: 初回 minor。末尾候補削除時に候補が残っても `selectedInputSourceId` が `null` になる可能性を検出した。

## 実装結果

- 既存 job 参照 input は `internal/service/translation_job_setup_service.go:1058` で判定し、候補一覧から除外する。
- 削除処理は `internal/service/translation_job_setup_service.go:466` で既存 job 参照を判定し、参照中なら削除しない。
- 入力候補 UI は `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte` でカード一覧へ変更した。
- canonical ER reset migration は `internal/infra/sqlite/migrations/014_canonical_er_cascade_reset.sql` と `internal/infra/sqlite/dbinit/migrations/014_canonical_er_cascade_reset.sql` に追加した。
- input 削除は `DeleteInputSource` から `DeleteXEditExtractedDataByID` だけを呼び、子行削除は DB の cascade に任せる。
- frontend は `deletingInputSourceId` を保持し、削除成功後に options 全再読込ではなく候補一覧を局所更新する。
- Sonar 指摘後、`DeleteInputSource` を transaction orchestration と helper に分割した。
- 契約レビュー指摘後、`existingJob` は候補除外と独立して返すようにした。
- 状態レビュー指摘後、削除対象が末尾でも残存候補を選択するようにした。

## 追加実装方針

- `014_canonical_er_cascade_reset.sql` を migration と dbinit migration の両方に追加し、canonical ER tables を drop して cascade 付きで再作成する。
- fresh DB 用の `003_canonical_er_v1_tables.sql` も cascade 付きへ更新する。
- Job Setup の `DeleteInputSource` は、既存 job 参照 guard の後に `DeleteXEditExtractedDataByID` だけを呼ぶ。
- frontend は削除中 input id を状態に持ち、削除成功時に対象候補だけを取り除く。

## 最終検証結果

- `python3 scripts/harness/run.py --suite frontend-local`: pass。57 files / 491 tests passed。
- `python3 scripts/harness/run.py --suite backend-local`: pass。
- `python3 scripts/harness/run.py --suite all`: pass。
- system test: pass。9 tests passed。
- coverage: pass。Sonar coverage 70.5%、line 71.3%、branch 63.9%。
- Sonar issues: pass。security 0、reliability 0、maintainability HIGH 0。

## レビュー通過根拠

- `reviewback.behavior.yaml`: `review_status: no_issue`、`must_fix_open: false`、`max_level: none`。
- `reviewback.contract.yaml`: `review_status: no_issue`、`must_fix_open: false`、`max_level: none`。
- `reviewback.trust-boundary.yaml`: `review_status: no_issue`、`must_fix_open: false`、`max_level: none`、`hard_gate: true`。
- `reviewback.state_invariant.yaml`: `review_status: no_issue`、`must_fix_open: false`、`max_level: none`。
- `reviewback.responsibility_boundary.yaml`: `review_status: no_issue`、`must_fix_open: false`、`max_level: none`。
- `implementation_action`: `close`。

## 停止条件

- input 削除 API の追加が既存責務境界と衝突する場合は、削除ボタンを実装せず戻す。
- 参照 job あり input を安全に判定できない場合は、削除実装へ進めず戻す。
