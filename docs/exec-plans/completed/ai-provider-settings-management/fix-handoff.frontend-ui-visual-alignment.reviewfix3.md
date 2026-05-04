# Fix Handoff: frontend-ui-visual-alignment.reviewfix3

- `workflow`: fix-lane
- `status`: ready-for-implementation_implementer
- `task_id`: `ai-provider-settings-management`
- `handoff_kind`: レビュー修正実行入力
- `implementation_skill`: `implement-frontend`
- `source_review`: `reviewback.state-invariant.yaml`
- `created_by`: `fix_lane`

## 修正対象指摘

- `state-invariant-002`: 再読み込み中の UI 操作が失敗時 snapshot 復元と競合する。
- level: `major`
- status: `open`

## 問題

Master Persona の AI 設定再読み込み中も、provider、model、executionMethod の select が操作可能である。
再読み込み失敗時は開始前 snapshot を復元するため、更新中に利用者が変更した最新選択が戻る。

## 修正範囲

- frontend プロダクトコードだけを変更する。
- `frontend/src/ui/screens/master-persona/MasterPersonaPage.svelte` を主対象にする。
- `isAISettingsRefreshing` が true の間は、provider、model、executionMethod、保存操作を無効化する。
- 更新 icon は既存どおり更新中に disabled と spinning を表示する。
- 必要なら `AIModelSelectionCard` の既存 disabled props を使う。

## 禁止変更範囲

- backend と `internal/` を変更しない。
- generated `wailsjs` を変更しない。
- route、保存項目、backend DTO、Wails gateway contract を追加しない。
- プロダクトテスト、test helper、fixture、snapshot を変更しない。
- APIキー、raw request、raw response、raw prompt を DOM text、DTO、log へ出さない。

## 回帰確認観点

- 更新中は provider、model、executionMethod を変更できない。
- 更新中は保存操作を実行できない。
- 更新失敗時に既存 provider、model、executionMethod、modelOptions が保持される。
- 更新成功時は保存済み AI 設定と modelOptions が反映される。
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
