# 境界結合テスト仕様書

- `status`: spec-as-should-be（本 task の主成果物。`finalization-module` で `docs/detail-specs/boundary-integration-test.md` に正本反映する）
- 作成日: 2026-06-06

---

## 1. 目的と位置付け

frontend と backend を Wails 境界で接続する repo において、境界 API 形式の整合を「片側変更で他方が落ちる検出経路」として固定する。

scenario E2E は UI 起点のため、UI に露出しない境界 DTO field 値の変化を検出できない。
単体テストは単一責務の検証であり、境界全体の整合を保証しない。
本テスト種別は両者の間の隙間を埋める。

**本書は複数 skill が共通の認識として参照する資料である**。
test-design、tests-scenario、tests-unit、design-module、implementation-module、investigation-module のいずれの skill から参照しても、同じ前提でテストを書ける形を保つ。

**境界結合テストは独立に設計されない**。
上位入力（次節）からの導出物として作る。観点表や case 表を白紙から起こすのではなく、上位入力に書かれた形式・意味・表示の規約から機械的に展開する。

---

## 2. 上位入力（責務分離）

境界結合テストは次の 3 系統から導出される。

| 系統 | 真とする file | 境界結合テストへの寄与 |
|---|---|---|
| 形式 | `frontend/src/controller/wails/<feature>.contract.ts` | request / response の型、必須、null 許容、値域 |
| 意味 | `docs/usecases/uc-<feature>.md` | 状態遷移と業務意味 |
| 表示 | `docs/screen-design/screens/<feature>.md` | 境界に必要な field の最小集合 |

本テスト種別は「形式」の整合だけを検証対象とする。意味と表示の検証は本テスト種別の責務外である。

---

## 3. 責務範囲

### 3.1 扱う

- 境界 contract の型、必須、null 許容、値域の固定
- 状態遷移の前後で contract 上の field がどう変化するか（UC が定める遷移を境界実装側で証明）
- 片側書き換え検出（contract または共有固定値の片側変更で他方が落ちる経路）

### 3.2 扱わない

| 対象 | 扱う test 種別 |
|---|---|
| 業務的意味（field が業務的に何を指すか） | UC docs の正本確認 |
| 表示文言、UI 構造、style | 画面設計の正本確認、storybook-module |
| 単一関数 / メソッドの内部実装 | unit test |
| UI 起点の業務シナリオ証明 | scenario E2E |

---

## 4. 構成要素（最小規約）

1 feature あたり次を持つ。

| 役割 | 配置先（最小規約） | 必要性 |
|---|---|---|
| 形式 contract | `frontend/src/controller/wails/` 配下 | 必須 |
| backend 境界結合テスト | `internal/apitest/` 配下 | 必須 |
| frontend 境界結合テスト | `frontend/src/controller/wails/` 配下 | 必須 |
| shared 固定値（golden 等） | `internal/apitest/testdata/` 配下 | 任意 |

backend と frontend は **同一 contract を真として参照する**。両側が異なる contract を参照すると片側書き換え検出経路が成立しない。

shared 固定値を使う場合、`assert 期待値を test source 内のリテラルにハードコードし、固定値ファイルは mock 応答にのみ流用する` 設計とする。
固定値を mock と assert 期待値の両方に使うと、固定値の片側変更で両側が同時に変わり検出経路が成立しない。

---

## 5. 既存テスト種別との責務境界

| テスト種別 | 入口 | 目的 |
|---|---|---|
| unit | 関数 / メソッド | 単一責務の実装可否 |
| 境界結合 | Wails 境界 | contract と両側実装の整合 |
| scenario E2E | UI 操作（playwright） | 業務シナリオの完走 |
| backend apitest（既存） | 公開 method | 受け入れ条件の system-level 証明 |

境界結合テストは「scenario E2E と単体テストの中間層」として位置付く。
UI 起点でないため scenario E2E ではなく、単一メソッドの内部実装ではないため単体テストでもない。
