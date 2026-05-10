# レビュー通過証跡

## 判定

状態: 通過。

未解決の must fix は無い。

## 初回レビュー結果

- 挙動正しさ: `no_issue`
- 契約・互換性: `no_issue`
- 権限・信頼境界: `no_issue`
- 責務境界: `no_issue`
- 状態・データ不変条件: `issues_open`

状態・データ不変条件レビューは `state-invariant-001` を `major` とした。

問題は、未開始 `JOB_PHASE_RUN` を実行状態へ昇格する更新が、DB 更新時点で期待 state 条件を持たないことだった。

## 追加修正

- `JobPhaseRunUpdateDraft` に `ExpectedState` を追加した。
- SQLite の `JOB_PHASE_RUN` 更新条件へ期待 state 条件を追加した。
- 更新件数 0 件の場合、既存 row を確認して状態不一致を `ErrConflict` として返すようにした。
- 単語翻訳、ペルソナ生成、本文翻訳の開始処理は、未開始 row 昇格時に期待 state を渡すようにした。

## 追加テスト

- SQLite repository test は、期待 state 一致時の更新成功を確認した。
- SQLite repository test は、期待 state 不一致時に `ErrConflict` となり、既存 state が変わらないことを確認した。
- 単語翻訳 phase service test は、昇格競合時に provider 実行と辞書保存副作用へ進まないことを確認した。
- ペルソナ生成 phase service test は、昇格競合時に provider 実行へ進まないことを確認した。
- 本文翻訳 phase service test は、昇格競合時に provider 実行と output 保存副作用へ進まないことを確認した。

## 再レビュー結果

`reviewback.state-invariant-rereview.yaml` は `no_issue`、`must_fix_open: false`、`max_level: none` と判定した。

初回 issue `state-invariant-001` は `resolved` になった。

## 検証

- `go test ./internal/repository ./internal/service`: 成功
- `python3 scripts/harness/run.py --suite backend-local`: 成功
- `python3 scripts/harness/run.py --suite coverage`: 成功
- coverage summary: `70.8%`
- `ruby -e 'require "yaml"; YAML.load_file("docs/exec-plans/completed/2026-05-10-job-setup-unstarted-phase-run-fix/reviewback.state-invariant-rereview.yaml"); puts "yaml ok"'`: 成功

## 残留未確認

実 DB で同一 job の同一 phase を同時開始する統合再現は未実施。

今回の証明は、DB の状態条件付き更新と単体テストの競合注入で行った。
