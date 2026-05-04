# Fix Handoff: frontend-ui-visual-alignment.reviewfix

- `workflow`: fix-lane
- `status`: ready-for-implementation_implementer
- `task_id`: `ai-provider-settings-management`
- `handoff_kind`: レビュー修正実行入力
- `implementation_skill`: `implement-frontend`
- `source_review`: `reviewback.behavior.yaml`
- `created_by`: `fix_lane`

## 追加人間観測

- 2026-05-04: 共通コンポーネント化は必要である。
- 2026-05-04: ただし現状の Master Persona はプロトタイプと見た目が全然違う。
- 2026-05-04: 合格条件は、共通 component を使うことではなく、プロトタイプのモデルカード確認と同じ見た目、同じ情報密度、同じ配置に揃うことである。

## 修正対象指摘

- `behavior-001`: Master Persona のモデル一覧更新導線がない。
- level: `major`
- status: `open`

## 問題

`MasterPersonaPage` は共有 component を使うが、プロトタイプのモデルカード確認と同じ見た目になっていない。
共通 component 化だけで合格にしてはいけない。

`MasterPersonaPage` は `onRefresh` も渡していない。
そのため、モデル未取得状態で利用者が同じカード内から更新操作へ進めない。

## 根拠

- `docs/exec-plans/active/ai-provider-settings-management/reviewback.behavior.yaml`
- `frontend/src/ui/components/AIModelSelectionCard.svelte`
- `frontend/src/ui/screens/master-persona/MasterPersonaPage.svelte`
- `frontend/src/application/contract/master-persona/master-persona-screen-contract.ts`
- `frontend/src/controller/master-persona/master-persona-screen-controller.ts`
- `frontend/src/application/usecase/master-persona/master-persona.usecase.ts`

## 修正範囲

- frontend プロダクトコードだけを変更する。
- `AIModelSelectionCard` の DOM、CSS、props を、`prototype/index.svelte` の `モデルカード確認` の見た目へ寄せる。
- `JobSetupPage` と `MasterPersonaPage` は、同じ `AIModelSelectionCard` から同じ見た目のモデルカードを表示する。
- Master Persona のカード内に、モデル一覧更新 icon または同等の回復操作を表示する。
- 既存 `loadAISettings` 相当の再読み込みで閉じられる場合は、backend API と Wails gateway は変更しない。
- controller contract へ frontend-local の更新操作を追加してよい。
- usecase に既存 gateway 呼び出しを再利用する薄い操作を追加してよい。

## 禁止変更範囲

- backend と `internal/` を変更しない。
- generated `wailsjs` を変更しない。
- docs 正本本文を変更しない。
- プロダクトテスト、test helper、fixture、snapshot を変更しない。
- 保存項目、route、DTO、backend gateway contract を追加しない。
- APIキー、raw request、raw response、raw prompt を DOM text、DTO、log へ出さない。

## 回帰確認観点

- `prototype/index.svelte` の `モデルカード確認` と `MasterPersonaPage` のモデル設定カードを見比べ、カード外形、余白、見出し、ラベル、status pill、select、更新 icon、処理方式、警告の位置が一致している。
- `prototype/index.svelte` の `モデルカード確認` と `JobSetupPage` のモデル設定カードを見比べ、同じ component 由来の見た目になっている。
- Master Persona だけ、縦長 form、保存 panel、別系統の余白、別系統の色になっていない。
- Master Persona のモデル設定カードに、モデル一覧更新 icon または同等の回復操作が表示される。
- モデル未取得状態の文言と操作が矛盾しない。
- Job Setup のモデル一覧更新 icon、Batch API checkbox、provider / model 選択が退行しない。
- AIサービス設定画面の amber 系配色統一が維持される。
- controller / presenter / view model contract 変更が frontend 内で完結する。

## 検証コマンド

- `npm --prefix frontend run check`
- `npm --prefix frontend run test -- provider-settings AppShell translation-job-setup master-persona`
- `npm --prefix frontend run build`
- `python3 scripts/harness/run.py --suite frontend-local`

## レビュー再判定対象

- `reviewback.behavior.yaml`
- `reviewback.contract.yaml`
- `reviewback.state-invariant.yaml`
- `reviewback.responsibility-boundary.yaml`

## 戻し条件

- backend DTO または Wails gateway の変更が必要な場合は停止する。
- 新規モデル一覧取得 API が必要な場合は停止する。
- 承認済み UI 設計にない保存項目や追加画面が必要な場合は停止する。
