# 次 task plan draft: 既存境界の結合テスト展開（boundary-integration-test-rollout）

- `status`: draft（本 task の active plan folder 内、user が preparation-module で正式起動する候補）
- 起源: 本 task `master-dictionary-boundary-test` の user 指示（2026-06-06）「既存境界の結合テスト作成のタスクも作っておいてほしい。まず、既存実装と画面設計の差異を埋める。次に差異を埋めた画面設計と実境界を突合して、結合テストを作っていく」
- 起点 spec: `boundary-integration-test.spec.md`（本 task で確定した境界結合テスト仕様書）
- 作成日: 2026-06-06

---

## 1. 親 task の目的

repo 内の既存境界（Wails 公開 method 群）全てについて、境界結合テストを揃える。
本 task の master-dictionary pilot で確立した枠組み（contract.ts + golden + 両側 boundary test）を全 feature に展開する。

---

## 2. 進行順序（user 指示）

### Phase 1: 既存実装と画面設計の差異を埋める

feature ごとに次を行う。

1. 現実装の境界 DTO を観察する
   - `internal/controller/wails/<feature>_controller.go`（または該当 file）の公開 method と request / response DTO
   - `frontend/src/controller/wails/<feature>.gateway.ts`（または該当 file）の type 定義
2. 画面設計（`docs/screen-design/screens/<feature>.md`）と突合する
3. 差異リストを作成する
4. user 判断で「画面設計が古い」か「実装が古い」かを feature 単位 / 差異単位で確定する
5. 古い方を修正する
   - 画面設計が古い場合: `storybook-module` 起点で画面設計を更新
   - 実装が古い場合: `design-module` → `implementation-module` で実装変更

Phase 1 の出力: 各 feature の「整合済み画面設計 + 整合済み実装境界」

### Phase 2: 差異を埋めた画面設計と実境界を突合して、結合テストを作る

feature ごとに次を行う。

1. 整合済み境界から `<feature>.contract.ts` を導出する
   - 配置: `frontend/src/controller/wails/<feature>.contract.ts`
   - 形式の真として位置付ける（boundary-integration-test.spec.md §2 に従う）
2. 必要に応じて shared golden を作成する
   - 配置: `internal/apitest/testdata/boundary/<feature>/`
3. backend / frontend の境界結合テストを作成する
   - backend: `internal/apitest/<feature>_boundary_test.go`
   - frontend: `frontend/src/controller/wails/<feature>.boundary.test.ts`
4. 片側書き換え検出を手動確認する
   - boundary-integration-test.spec.md §6.3 に従う

Phase 2 の出力: 各 feature の contract.ts + golden + boundary test 2 file

---

## 3. 対象 feature 一覧

`docs/screen-design/screens/` 配下の画面と、それぞれが対応する境界 controller / gateway をリストアップする。

優先順位は「基盤データ → 業務中核 → 実行系 → 出力系 → 集計表示」の順に並べる。

| # | feature | 画面設計 file | 境界 controller / gateway 候補 | 優先 |
|---|---|---|---|---|
| 0 | master-dictionary | `master-dictionary.md` | 完了済み（本 task で pilot 実施）。**追従 task 別途必要**（実装と画面を contract に追従させる、本 task plan.md の別 task 候補 2） | pilot |
| 1 | provider-settings | `provider-settings.md` | AI provider 設定。基盤 | 高 |
| 2 | master-persona | `master-persona.md` | 基盤データ | 高 |
| 3 | translation-job-setup | （`translation-management.md` 内 or 個別） | ジョブ生成 | 中 |
| 4 | translation-job-management | `translation-job-management.md` | ジョブ管理 | 中 |
| 5 | translation-input-review | `translation-input-review.md` | 入力レビュー | 中 |
| 6 | job-run | `job-run.md` | ジョブ実行のシェル | 中 |
| 7 | persona-generation-phase | `persona-generation-phase.md` | 実行 phase | 中 |
| 8 | term-translation-phase | `term-translation-phase.md` | 実行 phase | 中 |
| 9 | body-translation-phase | `body-translation-phase.md` | 実行 phase | 中 |
| 10 | translation-complete | `translation-complete.md` | 完了表示 | 低 |
| 11 | output-management | `output-management.md` | 出力管理 | 低 |
| 12 | dashboard | `dashboard.md` | 集計表示。境界 method を持たない可能性あり | 低 |
| 13 | translation-management | `translation-management.md` | 集計表示 | 低 |

`example-incomplete-job-list.md` は example file のため対象外。

`*-phase` 3 件は共通の phase 抽象に乗っている可能性が高く、まとめて 1 task として扱う案もある。preparation-module 段階で再判断する。

---

## 4. task 分割の選択肢

### 案 A: 親 task 1 つで全 feature を進める

- 親 task 名: `boundary-integration-test-rollout`
- 範囲: 13 feature 全て
- 利点: 一元管理。仕様変更時の方針調整が 1 箇所
- 欠点: 範囲が広すぎ、design-module / implementation-module の不変条件を満たしにくい。1 commit に収まらず、branch 寿命が長期化する

### 案 B: feature ごとに別 task を立てる（推奨）

- 親 task: 本 plan を「展開計画 docs」として正本化（`docs/exec-plans/templates/` 配下に置く案、または `docs/detail-specs/boundary-integration-test-rollout-plan.md` として正本化）
- 子 task: feature ごとに別 active plan folder を立てる
  - 例: `boundary-test-provider-settings`、`boundary-test-master-persona` など
  - 各 task が Phase 1 → Phase 2 を完結させる
- 利点: 1 task の範囲が pilot と同等で予測可能。順次進行できる
- 欠点: 共通方針の維持に親 plan の参照が要る

### 案 C: Phase 1 と Phase 2 を別 task として立てる

- task 1: 全 feature の画面設計 vs 実装 差異リスト作成と修正
- task 2: 全 feature の contract + boundary test 作成
- 利点: 段階が明確
- 欠点: task 2 の前提として task 1 が全 feature 分完了する必要があり、長期間ブロックされる

**推奨: 案 B**。pilot の規模（master-dictionary で約 1 日）を基準に、feature 単位で順次進める。

---

## 5. 各子 task の入口条件

各 feature 単位の子 task は次の入口条件を満たして preparation-module を起動する。

- 親方針 docs（本 plan の正本化版）を参照可能
- 対象 feature の境界 controller / gateway が現存する
- 画面設計（`<feature>.md`）が現存する
- UC（`docs/usecases/uc-<feature>.md`）が現存する（無い場合は新規作成または「UC が無い境界」として別途扱う判断）

---

## 6. 各子 task の進行（pilot に基づくテンプレ）

1. preparation-module: 作業 branch、active plan folder、plan.md を固定
2. Phase 1 観察: 現実装と画面設計の突合
   - investigation-module（差分観察）または design-module（仕様変更扱い）を判断
   - 差分が観察結果として整理されるだけなら investigation-module
   - 差分の解消に画面設計の修正が含まれるなら storybook-module も並走
   - 差分の解消に実装変更が含まれるなら design-module + implementation-module
3. Phase 1 整合: user 判断で真を確定、画面 / 実装を追従させる
4. Phase 2 contract 作成: design-module で contract.ts を確定（detail-spec-diff 起こし）
5. Phase 2 test 作成: implementation-module で golden + boundary test を作成
6. finalization-module: 正本化 + commit + merge

---

## 7. 親方針の正本化先候補

本 plan を正本化する場合の候補:

- `docs/detail-specs/boundary-integration-test-rollout-plan.md`（rollout 計画として）
- 上記とは別に `docs/exec-plans/templates/boundary-test-task-template/` 配下に子 task テンプレートを置く（preparation-module 起動時に参照）

両者を本 task の `finalization-module` で同時に正本化するか、別 task で正本化するかは user 判断とする。

---

## 8. master-dictionary 追従 task との関係

本 task plan.md の「別 task 候補 2: 現実装と golden / test の追従」は、master-dictionary 個別の追従 task。
本 rollout plan の対象（13 feature 全体）の一部として位置付ける。

順序:
1. master-dictionary 追従 task が先行（既に contract.ts と突合結果がある）
2. その他 feature は Phase 1 → Phase 2 の順で並列または順次進行

---

## 9. 観察項目（rollout 開始前に検討すべき）

1. 全 feature を pilot 規模で進められるか、feature ごとに規模が異なるか
2. UC（`docs/usecases/`）が未整備の境界があるか。ある場合は UC 整備を先行させる
3. 画面設計 file 単位 / 実装 controller 単位の対応関係（1:1 でない feature があるか）
4. `*-phase` 3 件をまとめるか分けるか
5. dashboard / translation-management の境界 method の有無（無い場合は対象外）

---

## 10. 起動手順（user が次にやること）

1. 本 task を `finalization-module` で完了させる
2. 本 plan を `docs/detail-specs/boundary-integration-test-rollout-plan.md` に正本化するかを判断
3. 最初の子 task として「master-dictionary 追従」または「provider-settings 境界結合テスト」を `preparation-module` で立ち上げる
4. 子 task のテンプレが必要なら、本 plan の §6 を参照する形で進める
