# Fix Handoff: frontend-ui-visual-alignment.reviewfix2

- `workflow`: fix-lane
- `status`: ready-for-implementation_implementer
- `task_id`: `ai-provider-settings-management`
- `handoff_kind`: レビュー修正実行入力
- `implementation_skill`: `implement-frontend`
- `source_review`: `reviewback.state-invariant.yaml`
- `created_by`: `fix_lane`

## 修正対象指摘

- `state-invariant-001`: モデル一覧更新失敗で AI 設定選択が default に戻る。
- level: `major`
- status: `open`

## 問題

Master Persona の `refreshAISettings()` は `loadAISettings()` に直結している。
`loadAISettings()` 失敗時の `catch` が `aiSettings` を default に戻し、`modelOptions` を空配列にする。
そのため、利用者が選択中の provider、model、executionMethod を更新失敗だけで失う。

## 修正範囲

- frontend プロダクトコードだけを変更する。
- `frontend/src/application/usecase/master-persona/master-persona.usecase.ts` を主対象にする。
- AI 設定再読み込み失敗時は、既存 `aiSettings` と `modelOptions` を保持する。
- 初回ロード失敗時の扱いが必要な場合も、既存選択がある時は上書きしない。
- error message は表示してよい。

## 禁止変更範囲

- backend と `internal/` を変更しない。
- generated `wailsjs` を変更しない。
- route、保存項目、backend DTO、Wails gateway contract を追加しない。
- プロダクトテスト、test helper、fixture、snapshot を変更しない。
- APIキー、raw request、raw response、raw prompt を DOM text、DTO、log へ出さない。

## 回帰確認観点

- 更新失敗時に、既存 provider、model、executionMethod、modelOptions が保持される。
- 更新成功時は、保存済み AI 設定と modelOptions が反映される。
- `isAISettingsRefreshing` は失敗時も解除される。
- Master Persona のモデルカード外観修正を退行させない。

## 検証コマンド

- `npm --prefix frontend run check`
- `npm --prefix frontend run test -- provider-settings AppShell translation-job-setup master-persona`
- `npm --prefix frontend run build`
- `python3 scripts/harness/run.py --suite frontend-local`

## レビュー再判定対象

- `reviewback.state-invariant.yaml`
- `reviewback.behavior.yaml`
- `reviewback.contract.yaml`
- `reviewback.trust-boundary.yaml`
