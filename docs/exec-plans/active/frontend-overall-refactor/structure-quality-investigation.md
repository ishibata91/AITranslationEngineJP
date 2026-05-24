# 構造品質調査: frontend-overall-refactor

## 作業ログ

- 担当 role: `investigator`
- 使用 skill / agent: `investigate` / `investigator`
- 編集対象: `docs/exec-plans/active/frontend-overall-refactor/structure-quality-investigation.md`
- 書き換え許可範囲: `docs/exec-plans/active/frontend-overall-refactor/structure-quality-investigation.md`
- 編集禁止範囲: `frontend/src/`, `internal/`, プロダクトテスト, docs 正本本文, `.codex/`, `docs/exec-plans/completed/`, `docs/exec-plans/active/frontend-overall-refactor/refactor-classification.md`, `docs/exec-plans/active/frontend-overall-refactor/plan.md`, remote repository
- 人間承認要否: この成果物の記載は承認不要。後続の `リファクタ範囲確認` と `implementation-scope` は人間承認が必要

## 調査条件

- 調査 mode: `構造品質調査`
- 仕様前提:
  - `FSD-005`: `実装が正`。入力確認から直接ジョブ作成して `job-run` へ進む導線は code 修正候補にしない。
  - `FSD-006`: `JobSetupPage.svelte` は廃止済み dead code として扱う。
  - 旧 `FSD-001` から `FSD-004`: 仕様乖離ではなく構造品質論点として扱う。
- 設計参照: `docs/architecture.md`, `docs/coding-guidelines-frontend.md`
- 実装参照: `frontend/src/main.ts`, `frontend/src/bootstrap/`, `frontend/src/application/`, `frontend/src/controller/`, `frontend/src/ui/`
- 判断方針: 責務過多、責務分離不足、コーディング規約逸脱、構造設計不整合、未使用コードを分ける。`FSD-005` で正とされた導線は変更不要範囲へ退避する。

## 判断結果

- 完了判定: 完了
- 根拠参照:
  - `docs/architecture.md:96-103`
  - `docs/architecture.md:122-148`
  - `docs/architecture.md:263-279`
  - `docs/coding-guidelines-frontend.md:18-24`
  - `docs/coding-guidelines-frontend.md:30-43`
  - `docs/coding-guidelines-frontend.md:61-79`
  - `docs/exec-plans/active/frontend-overall-refactor/refactor-classification.md`
  - `docs/exec-plans/active/frontend-overall-refactor/spec-drift-investigation.md`
- 引き継ぎ先: `designer`
- 推奨 next step: `リファクタ範囲確認`

## 観測点

- `frontend/src/main.ts` は `createProductionAppFactories()` の戻り値から `createTranslationJobSetupScreenController` を `App` へ渡している。
- `frontend/src/ui/App.svelte` は fallback wiring と controller factory / gateway 生成を内部で持ち、`main.ts` から渡された `createTranslationJobSetupScreenController` prop を受け取っていない。
- `frontend/src/ui/views/AppShell.svelte` と `frontend/src/ui/stores/shell-state.ts` は `translation-job-setup` view を持たず、翻訳管理 shell では `job-management`、`input-review`、`job-run` だけを切り替える。
- `frontend/src/bootstrap/app-screen-controller-factories.ts` は `translation-job-setup` を含む全 screen controller factory と gateway を production wiring として生成している。
- `frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts` は 1310 行、`frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts` は 821 行である。
- `frontend/src/ui/components/AIModelSelectionCard.svelte` は 526 行、`frontend/src/ui/components/ProcessingTargetListPanel.svelte` は 525 行である。
- `frontend/src/controller/runtime/` には `master-dictionary` の event adapter と `master-persona` の polling adapter だけがある。runtime 連携 adapter の正本配置 `frontend/src/controller/wails/` とは離れている。

## 責務過多

| ID | 観測事実 | 根拠参照 | 対象範囲 | 変更不要範囲 | 修正候補 | リファクタ範囲確認へ回すか |
| --- | --- | --- | --- | --- | --- | --- |
| `SQ-001` | `translation-job-setup` usecase と presenter は、dead code になった page 専用の状態整形と provider / validation / summary の整合まで 1 ファイルへ集約している。`translation-job-setup.usecase.ts` は phase 選択、モデル一覧更新、validation、job 作成の全更新手順を 1 クラスで持つ。`translation-job-setup.presenter.ts` は page 用の extended view model、phase card、summary card、validation label を 1 ファイルで構築する。 | `frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts:1-260`, `frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts:1-260`, `frontend/src/controller/translation-job-setup/translation-job-setup-screen-controller-factory.ts:1-29`, `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte:34-205` | `frontend/src/application/usecase/translation-job-setup/`, `frontend/src/application/presenter/translation-job-setup/`, `frontend/src/controller/translation-job-setup/` | `InputReviewPage.svelte` から `job-run` へ直接進む現行導線。`FSD-005` の正しい振る舞いは変更不要。 | dead code cleanup を優先し、残す部品がある場合だけ `runtime option / phase card / summary` の単位で再配置候補を切り分ける。 | 回す |
| `SQ-002` | shared component 2 件が画面専用ルールまで抱えている。`AIModelSelectionCard.svelte` は provider 選択、credential 表示、モデル更新、batch toggle、execution mode、action button まで 1 props 集合で扱う。`ProcessingTargetListPanel.svelte` は pagination、expanded row、tooltip、metadata 詳細、action 群まで 1 component で扱う。どちらも `docs/architecture.md` の shared component 条件より、画面固有分岐と props 集合が肥大化している。 | `docs/architecture.md:96-114`, `docs/coding-guidelines-frontend.md:30-43`, `frontend/src/ui/components/AIModelSelectionCard.svelte:1-260`, `frontend/src/ui/components/ProcessingTargetListPanel.svelte:1-260` | `frontend/src/ui/components/AIModelSelectionCard.svelte`, `frontend/src/ui/components/ProcessingTargetListPanel.svelte` | 現在の表示導線と画面文言。`shared component` を即分割すると design 変更に波及する部分は今回の調査では固定しない。 | screen local 化すべき枝と shared のまま残す最小核を分ける候補。`AIModelSelectionCard` は provider/model 選択と secondary control を分離候補、`ProcessingTargetListPanel` は pagination shell と detail row を分離候補。 | 回す |

## 責務分離不足

| ID | 観測事実 | 根拠参照 | 対象範囲 | 変更不要範囲 | 修正候補 | リファクタ範囲確認へ回すか |
| --- | --- | --- | --- | --- | --- | --- |
| `SQ-003` | `App.svelte` が View でありながら production wiring fallback を内部で持つ。`@controller/*` factory import と `@controller/wails/*` gateway import を持ち、`resolve*Factory()` で不足 prop を補う。`docs/coding-guidelines-frontend.md` は production gateway と controller factory の生成を composition root に置くと定義している。`main.ts` と `bootstrap/app-screen-controller-factories.ts` が既に composition root を持つため、View 側 fallback は責務の二重化になっている。 | `docs/architecture.md:67-76`, `docs/coding-guidelines-frontend.md:61-79`, `frontend/src/main.ts:1-35`, `frontend/src/bootstrap/app-screen-controller-factories.ts:1-106`, `frontend/src/ui/App.svelte:1-110` | `frontend/src/ui/App.svelte`, `frontend/src/main.ts`, `frontend/src/bootstrap/app-screen-controller-factories.ts` | `AppShell.svelte` が持つ route / hash / view 切替の画面責務。画面文言と navigation の振る舞いは今回の調査対象外。 | `App.svelte` を pure View に寄せ、production fallback を `main.ts` または `bootstrap/` に集約する候補。 | 回す |

## コーディング規約逸脱

| ID | 観測事実 | 根拠参照 | 対象範囲 | 変更不要範囲 | 修正候補 | リファクタ範囲確認へ回すか |
| --- | --- | --- | --- | --- | --- | --- |
| `SQ-006` | screen component props と fixture が presenter concrete と gateway-contract concrete に直接依存している。`JobSetupPage.svelte`、`PhaseSettingsPanel.svelte`、`job-setup-panel-props.ts`、`translation-job-setup-panel-fixtures.ts` は `@application/presenter/translation-job-setup/...presenter` の型を直接 import する。`InputReviewPage.svelte` も `@application/presenter/translation-input` から `canOpenJobSetup`、`ERROR_LABELS`、`STATUS_LABELS` を直接参照する。`coding-guidelines-frontend.md` は `.svelte` を表示とイベント配線に集中させ、表示整形は Presenter へ分けるとしているため、screen component が presenter concrete の exported helper へ寄りかかっている。 | `docs/coding-guidelines-frontend.md:18-24`, `docs/coding-guidelines-frontend.md:30-43`, `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte:1-105`, `frontend/src/ui/screens/translation-job-setup/PhaseSettingsPanel.svelte:1-56`, `frontend/src/ui/screens/translation-job-setup/job-setup-panel-props.ts:1-86`, `frontend/src/ui/screens/translation-job-setup/__fixtures__/translation-job-setup-panel-fixtures.ts:1-155`, `frontend/src/ui/screens/translation-input/InputReviewPage.svelte:4-170` | `frontend/src/ui/screens/translation-job-setup/`, `frontend/src/ui/screens/translation-input/`, `frontend/src/ui/screens/master-persona/`, `frontend/src/ui/screens/provider-settings/`, `frontend/src/ui/screens/translation-output-artifact/` | `@application/contract/*` から受ける controller contract と screen types。view model を読むだけの screen component 自体は残し得る。 | presenter export に依存している view helper と型を `application/contract` 側の screen types へ寄せる候補。screen local props 型は UI 層で閉じる候補。 | 回す |

## 構造設計不整合

| ID | 観測事実 | 根拠参照 | 対象範囲 | 変更不要範囲 | 修正候補 | リファクタ範囲確認へ回すか |
| --- | --- | --- | --- | --- | --- | --- |
| `SQ-004` | `docs/architecture.md` の directory 正本は、screen local な controller / usecase / presenter / store を `frontend/src/ui/` に置き、`frontend/src/application/` は shared contract、`frontend/src/controller/wails/` は Wails adapter だけに分けると定義している。実装は `frontend/src/application/usecase/translation-input/translation-input.usecase.ts`、`frontend/src/application/presenter/translation-input/translation-input.presenter.ts`、`frontend/src/controller/translation-input/translation-input-screen-controller.ts` のように screen local object を `application/` と `controller/` へ展開している。`architecture.md` の current canonical path と実 directory がずれている。 | `docs/architecture.md:263-279`, `frontend/src/application/usecase/translation-input/translation-input.usecase.ts:1-80`, `frontend/src/application/presenter/translation-input/translation-input.presenter.ts:1-80`, `frontend/src/controller/translation-input/translation-input-screen-controller.ts:1-120`, `frontend/src/ui/screens/README.md`, `frontend/src/application/README.md` | `frontend/src/application/`, `frontend/src/controller/`, `frontend/src/ui/screens/` | shared contract、gateway contract、`frontend/src/controller/wails/` の adapter 境界そのもの。`FSD-005` の導線は変更不要。 | directory 正本に実装を寄せるか、docs 正本を更新するかの人間判断が必要。code refactor を始める前に正本優先の確認が必要。 | 回す |
| `SQ-005` | Wails 境界の正本は `frontend/src/controller/wails/` に gateway、DTO、runtime 連携 adapter を置く構成である。一方で実装は `frontend/src/controller/runtime/master-dictionary/master-dictionary-runtime-event-adapter.ts` と `frontend/src/controller/runtime/master-persona/master-persona-runtime-polling-adapter.ts` を別 tree に置き、`master-persona` では Wails event ではなく polling を screen controller factory から直接注入している。runtime adapter が Wails 境界から分離され、`docs/architecture.md` の `RuntimeEventAdapter` 主語とも `controller/wails` 配置とも揃っていない。 | `docs/architecture.md:34-37`, `docs/architecture.md:141-148`, `docs/architecture.md:263-288`, `frontend/src/controller/wails/README.md:1`, `frontend/src/controller/runtime/master-dictionary/master-dictionary-runtime-event-adapter.ts:1-209`, `frontend/src/controller/runtime/master-persona/master-persona-runtime-polling-adapter.ts:1-40`, `frontend/src/controller/master-dictionary/master-dictionary-screen-controller-factory.ts:22-48`, `frontend/src/controller/master-persona/master-persona-screen-controller-factory.ts:14-29` | `frontend/src/controller/runtime/`, `frontend/src/controller/wails/`, `frontend/src/controller/master-dictionary/`, `frontend/src/controller/master-persona/` | `master-dictionary` の import progress push 通知自体。`master-persona` の定期更新が必要かどうかの機能要否は今回の調査対象外。 | runtime adapter の正本配置と transport 種別を先に確定する候補。polling 継続可否は別判断に分け、少なくとも adapter 境界の配置不整合を解く候補。 | 回す |

## 未使用コード

| ID | 観測事実 | 根拠参照 | 対象範囲 | 変更不要範囲 | 修正候補 | リファクタ範囲確認へ回すか |
| --- | --- | --- | --- | --- | --- | --- |
| `SQ-007` | `translation-job-setup` の本番 wiring は残っているが、実使用経路は消えている。`main.ts` は `createTranslationJobSetupScreenController` を `App` へ渡す。`bootstrap/app-screen-controller-factories.ts` も gateway と screen controller factory を生成する。だが `App.svelte` はその prop を受け取らず、`AppShell.svelte` と `shell-state.ts` は `translation-job-setup` view を持たない。`InputReviewPage.svelte` は入力確認から `createTranslationJobFromSelected()` を呼び、直接 `job-run` へ進む。`JobSetupPage.svelte` 本体に加え、controller factory、gateway、test、fixture、story が本番導線外に残っている。`AppShell.test.ts` は存在しない prop `createTranslationJobSetupScreenController` を渡しており、dead wiring が test にも残っている。 | `frontend/src/main.ts:12-33`, `frontend/src/bootstrap/app-screen-controller-factories.ts:9-23`, `frontend/src/bootstrap/app-screen-controller-factories.ts:48-102`, `frontend/src/ui/App.svelte:30-52`, `frontend/src/ui/views/AppShell.svelte:294-339`, `frontend/src/ui/stores/shell-state.ts:9-18`, `frontend/src/ui/screens/translation-input/InputReviewPage.svelte:149-170`, `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte:1-205`, `frontend/src/ui/views/AppShell.test.ts:296-329` | `frontend/src/ui/screens/translation-job-setup/`, `frontend/src/controller/translation-job-setup/`, `frontend/src/controller/wails/translation-job-setup.gateway.ts`, `frontend/src/application/usecase/translation-job-setup/`, `frontend/src/application/presenter/translation-job-setup/`, `frontend/src/application/store/translation-job-setup/`, `frontend/src/application/contract/translation-job-setup/`, `frontend/src/application/gateway-contract/translation-job-setup/`, `frontend/src/ui/views/AppShell.test.ts`, `frontend/src/main.ts`, `frontend/src/bootstrap/app-screen-controller-factories.ts` | `InputReviewPage.svelte` から直接 job を作成する現行導線。`job-run`、`translation-job-management`、`shell-state` の現行 view 切替は変更不要。 | dead code cleanup 候補。削除範囲が広いため、`App` / `main.ts` / `bootstrap` の未使用 wiring、screen / controller / gateway 群、story / fixture / test 群を同じ slice で扱うか人間確認が必要。 | 回す |

## 変更不要範囲

- `FSD-005` で人間が正とした `InputReviewPage.svelte` から `job-run` へ直接進む導線。
- `AppShell.svelte` と `shell-state.ts` が表している現行の翻訳管理 step 並び。
- `TranslationJobManagementPage.svelte`、`JobRunPage.svelte`、`TranslationOutputArtifactPage.svelte` のユーザー向け振る舞い。
- `master-dictionary` の import progress 通知と `master-persona` の定期更新が必要かどうかの機能要否判断。

## 残り不足

- `SQ-004` は docs 正本を守るか、実装配置を正として docs 正本化候補へ戻すかの人間判断が未確定である。
- `SQ-005` は polling 継続可否そのものをこの成果物では判断しない。transport 種別の正本判断は後続で必要になる可能性がある。
- `SQ-007` は dead code cleanup の削除単位を 1 回で扱うか、wiring 削除と screen 群削除に分けるかの人間判断が未確定である。

## 残留リスク

- `SQ-004` を判断せずに実装へ進むと、directory 正本に寄せる refactor と docs 正本化候補が混在する可能性がある。
- `SQ-006` と `SQ-007` は密接であり、dead code cleanup だけ先に進めると presenter 型境界の改善対象が変わる可能性がある。
- `SQ-003` と `SQ-007` は同じ wiring 面を触るため、fallback wiring の解消と dead wiring の削除を別々に実施すると二度手間になる可能性がある。
