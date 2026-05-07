# provider settings job decoupling implement lane

- `task_id`: `2026-05-07-provider-settings-job-decoupling-implement`
- `lane`: `implement-lane`
- `status`: `completed`
- `current_artifact`: `作業計画完了移動`
- `source_stop_plan`: `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling/plan.md`

## task 枠

- 人間依頼: クレデンシャル管理を DB レベルで Job から分離する。
- 人間依頼: AIサービス設定で設定できる値を Job から分離する。
- 人間依頼: secret store 情報と endpoint を全 job provider 共通設定にする。
- 対象: Job Setup、provider settings、翻訳フェーズ実行、DB 永続化境界。
- 非対象: 直接のプロダクトコード変更、直接のプロダクトテスト変更、承認なしの docs 正本本文更新。

## 成果物状態

- `task 枠`: 完了。
- `scenario_candidates`: 完了。6 観点、42 候補。
- `シナリオ設計`: 完了。scenario gate は pass。
- `UI設計`: 完了。
- `設計差分図`: 完了。PlantUML 構文検証と SVG 生成は成功。
- `人間設計レビュー`: 承認済み。人間応答は `ok`。
- `実装範囲`: 完了。
- `実装引き継ぎ入力`: 完了。`frontend-implementation-input.md`
- `frontend 実装`: 完了。`frontend-implementation-result.md`
- `単体テスト`: 完了。frontend test 追従は `frontend-test-followup-result.md`、lane 内単体テストは `unit-test-result.md`
- `frontend 実装後人間レビュー`: 承認済み。人間応答は「フロントで見たいものないから先進めていい」。
- `backend 実装`: 完了。`PSJD-BE-01` と `PSJD-BE-02` は完了。
- `統合境界実装`: 完了。`PSJD-INT-01`
- `単体テスト`: 完了。`PSJD-UT-01`
- `シナリオテスト`: 完了。`PSJD-SCN-01`
- `最終検証`: 完了。レビュー修正後の `python3 scripts/harness/run.py --suite all` は通過。
- `レビュー通過根拠`: 完了。5 観点すべて `no_issue`。
- `正本化判断`: 完了。詳細仕様正本反映が必要。
- `詳細仕様正本反映`: 完了。`docs-canonicalization-result.md`
- `作業レポート入力`: 完了。`work-report-input.md`
- `作業レポート`: 完了。`work_history/runs/2026-05-07-provider-settings-job-decoupling-implement-run/`
- `作業計画完了移動`: 完了。`docs/exec-plans/completed/2026-05-07-provider-settings-job-decoupling-implement/`

## 判断根拠

- 軽量変更計画は、永続仕様と公開契約の確定が必要として設計戻しで停止した。
- `PROVIDER_SETTINGS` は provider ごとの endpoint と credential 参照状態を保存済みである。
- `JOB_PHASE_RUN` と `TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT` は credential / endpoint 系 snapshot を持つ。
- 既存詳細仕様は、Ready job は実行前に最新 provider settings を再解決し、Running phase は開始時 snapshot を使うと定義している。

## 検証方針

- 設計段階ではプロダクトコードとプロダクトテストを変更しない。
- 実装範囲承認後、backend 変更は `python3 scripts/harness/run.py --suite backend-local` で確認する。
- frontend 変更が発生する場合は `python3 scripts/harness/run.py --suite frontend-local` と実画面確認を使う。
