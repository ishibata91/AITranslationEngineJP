# Job Setup モデル一覧取得失敗の修正

## 状態

- `task_id`: `2026-05-07-job-setup-model-list-token-fix`
- `lane`: `fix-lane`
- `status`: `completed`
- `source`: 人間観測

## 人間観測

- 画面: ジョブセットアップ画面
- 操作: モデルカードのモデル一覧取得ボタンを押す。
- 結果: 必ず「モデルを取得できませんでした」と表示される。
- 期待: 保存済み provider 設定を使い、モデル一覧を取得できる。
- 比較: マスターペルソナのモデルカードではモデル一覧取得が成功する。

## 成果物

- 人間観測記録: [human-observation.md](./human-observation.md)
- 修正前調査: 完了
- 原因箇所シーケンス図: [cause-sequence.md](./cause-sequence.md)
- 修正実行入力: [fix-execution-input.md](./fix-execution-input.md)
- 実装証跡: [implementation-evidence.md](./implementation-evidence.md)
- 回帰テスト証跡: [regression-test-evidence.md](./regression-test-evidence.md)
- レビュー通過根拠: [review-pass-evidence.md](./review-pass-evidence.md)
- 作業レポート入力: [work-report-input.md](./work-report-input.md)

## 検証予定

- `go test ./internal/service -run 'TranslationJobSetupService.*ProviderSettings|ProviderSettingsTestSafe|ModelList'`
- 必要なら `npm --prefix frontend run test -- src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts`
