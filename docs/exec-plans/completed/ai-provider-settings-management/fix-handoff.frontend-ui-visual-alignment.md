# Fix Handoff: frontend-ui-visual-alignment

- `workflow`: fix-lane
- `status`: ready-for-implementation_implementer
- `task_id`: `ai-provider-settings-management`
- `handoff_kind`: 修正実行入力
- `implementation_skill`: `implement-frontend`
- `created_by`: `fix_lane`

## 人間観測記録

- 観測 1: 既存プランの UIプロトタイプとモデル選択 UI に大きな差分がある。
- 観測 1 の期待: モデル選択 UI をコンポーネント化し、一貫性のあるデザインにする。
- 観測 2: AIサービス設定画面の色が承認済み UIプロトタイプと異なる。
- 観測 2 の期待: AIサービス設定画面を既存 AppShell と UIプロトタイプの配色へ合わせる。

## 修正前調査

- 判断結果: `修正前調査` は完了である。
- 根拠: `investigator` が UIプロトタイプ、本番画面、Svelte 実装の 3 系統で差分を確認した。
- UI 証跡: `tmp/agent-browser/ai-provider-settings-management/prototype-provider-settings.png`
- UI 証跡: `tmp/agent-browser/ai-provider-settings-management/prototype-model-cards.png`
- UI 証跡: `tmp/agent-browser/ai-provider-settings-management/app-provider-settings.png`
- UI 証跡: `tmp/agent-browser/ai-provider-settings-management/app-translation-job-setup.png`
- UI 証跡: `tmp/agent-browser/ai-provider-settings-management/app-master-persona.png`

## 観測事実

- 承認済み UI 設計は、参照側に `モデルカード確認` を置き、3 枚のモデルカードで `AIサービス`、`モデル`、更新 icon button、`Batch API` checkbox、`?` tooltip を揃える前提である。
- UIプロトタイプは、`モデルカード確認` を単独画面として持ち、3 枚のカードが同じ構造で `AIサービス`、`モデル`、更新、`Batch API` を持つ。
- 本番 `JobSetupPage` はカード型で、更新 icon、`Batch API` checkbox、警告文を持つ。
- 本番 `MasterPersonaPage` は縦並びの select 3 個と保存ボタンだけで、更新 icon と `?` tooltip を持たない。
- 本番 `ProviderSettingsPage` は blue 系 gradient と blue 系 label color を持ち、UIプロトタイプと AppShell の amber 系 panel と異なる。

## 影響ファイル候補

- `frontend/src/ui/screens/provider-settings/ProviderSettingsPage.svelte`
- `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`
- `frontend/src/ui/screens/master-persona/MasterPersonaPage.svelte`
- `frontend/src/ui/views/AppShell.svelte`

## 実装対象

- frontend プロダクトコードだけを変更する。
- `ProviderSettingsPage` の配色を、承認済み UIプロトタイプと AppShell の dark shell / amber 系 token に合わせる。
- 参照側モデル選択 UI の共有 Svelte component を `frontend/src/ui` 配下へ追加する。
- `JobSetupPage` と `MasterPersonaPage` のモデル選択 UI を共有 component 経由に寄せる。
- 既存 controller、usecase、gateway の責務境界を維持する。

## 禁止変更範囲

- `docs/` 正本本文を変更しない。
- `.codex/` と `.codex/skills` を変更しない。
- backend と `internal/` を変更しない。
- プロダクトテスト、test helper、fixture、snapshot を変更しない。
- UIプロトタイプの sample data を product code へ移植しない。
- generated `wailsjs` を View、ScreenController、Frontend UseCase へ直接 import しない。

## UI 設計根拠

- `docs/exec-plans/active/ai-provider-settings-management/ui-design.md`
- `docs/exec-plans/active/ai-provider-settings-management/prototype/index.svelte`
- `docs/exec-plans/active/ai-provider-settings-management/implementation-scope.md`

## 回帰確認観点

- `AIサービス設定` 画面の panel、status pill、label、button、input が AppShell と UIプロトタイプの色系統から外れない。
- `Job Setup` と `Master Persona` のモデル選択 UI が同一 component 由来の構造と余白で表示される。
- 390px 幅で provider list、設定詳細、モデル選択 UI が横にはみ出さない。
- モデル一覧更新中はボタン外形を動かさず、icon だけが状態を示す。
- `Batch API` の説明は `?` tooltip に閉じ、通常表示を圧迫しない。
- APIキー、raw request、raw response、raw prompt は DOM text、DTO、log に出ない。

## 検証コマンド

- `npm --prefix frontend run check`
- `npm --prefix frontend run test -- provider-settings AppShell translation-job-setup master-persona`
- `npm --prefix frontend run build`
- `python3 scripts/harness/run.py --suite frontend-local`

## UI 確認

- UIプロトタイプ確認: `npm --prefix frontend run dev:prototype -- --task ai-provider-settings-management --port 34116`
- UIプロトタイプ URL: `http://127.0.0.1:34116/prototype`
- 本番確認: `npm run dev:wails:agent-browser`
- 本番 URL: `http://localhost:34115`
- 確認対象: `#provider-settings`、`#translation-management`、`#master-persona`

## 戻し条件

- 共有 component 化が controller contract 変更、backend DTO 変更、Wails gateway 変更を必要とする場合は停止する。
- 承認済み UI 設計にない新規状態、追加操作、追加保存項目が必要な場合は停止する。
- プロダクトテスト変更が必要な場合は停止する。
