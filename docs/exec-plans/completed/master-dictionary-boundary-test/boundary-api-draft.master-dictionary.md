# 境界 API 補助文書: MasterDictionary（pilot 6 method）

- `status`: supplementary（契約の真は `boundary-api.master-dictionary.contract.ts`）
- `handoff_id`: handoff-boundary-api-draft
- `spec_basis`: detail-spec-diff.md の `boundary-test-REQ-005` および `master-dictionary-REQ-pilot-001`
- 正本反映先: `docs/detail-specs/master-dictionary.md` 内「境界 API」節（補助文書として）
- 作成日: 2026-06-05
- 書き直し履歴:
  - 2026-06-05 初回: 現実装観測ベースで draft 作成
  - 2026-06-05 1 回目書き直し: 現実装ベースの API 仕様書として体裁化
  - 2026-06-05 2 回目書き直し: あるべき仕様として書き直し（`note`、`origin` 削除）
  - 2026-06-06 3 回目書き直し: 形式の真を `boundary-api.master-dictionary.contract.ts` に移譲。本 file は「画面設計 / 現実装との突合結果」のみを残す補助文書に縮退

---

## 1. 本書の位置付け

本書は MasterDictionary 機能の境界 API の **形式の真ではなく** 、画面設計と現実装との **突合結果** を記録する補助文書である。

責務分離:

| 持ち場 | 真とする file |
|---|---|
| 形式（型、必須、null 許容、値域） | `boundary-api.master-dictionary.contract.ts` |
| 意味（field が業務的に何を指すか） | 画面仕様 / UC（`docs/usecases/uc-master-dictionary.md`） |
| 状態遷移規約 | UC |
| 表示文言、UI 構造 | 画面設計（`docs/screen-design/screens/master-dictionary.md`） |

本書は上記の「持ち場」を跨ぐ過渡情報（突合差分、追従候補）を一時的に記録する。

---

## 2. 突合結果（あるべき形式 vs 画面設計 / 現実装）

形式の真である `boundary-api.master-dictionary.contract.ts` と、画面設計および現実装の双方を突合した結果を記述する。

突合対象:
- 画面設計: `docs/screen-design/screens/master-dictionary.md`（読み込み日: 2026-06-05）
- 現実装: `internal/controller/wails/`、`internal/usecase/`、`internal/service/` 配下の MasterDictionary 関連 Go DTO、`frontend/src/controller/wails/master-dictionary.gateway.ts`

### 2.1 突合の方針

- 形式の真である `boundary-api.master-dictionary.contract.ts` を真とする
- 画面設計と現実装の両方を本契約に追従させる方向で乖離を記録する
- presentation 層の責務（gateway 接続状態、進捗バー、状態文の文言など）は契約の範囲外であり、差分判定の対象外とする
- 意味の乖離は本書では扱わない（画面仕様 / UC の責務）

### 2.2 あるべき形式 vs 画面設計（画面設計が古い）

#### 差分 2.2.1: 辞書一覧パネル metadata「メモ」「登録元」表示

- 画面設計 `[6] 辞書一覧パネル`: `metadata は 辞書 ID、カテゴリ、登録元、最終更新、メモ を表示する`
- 契約: `MasterDictionaryEntry`（一覧で使う）に `note`（メモ）と `origin`（登録元）の field は存在しない
- 判定: **画面設計が古い**
- 是正方向: 画面設計から「メモ」と「登録元」を一覧 metadata から外す

#### 差分 2.2.2: 詳細パネル「由来」表示

- 画面設計 `[10] 詳細パネル`: 表示項目に `カテゴリ、由来、原文、訳語、ID、最終更新、詳細状態文` を挙げる
- 契約: `MasterDictionaryEntry`（詳細で使う）に `origin`（由来 / 登録元）の field は存在しない
- 判定: **画面設計が古い**
- 是正方向: 画面設計の詳細パネル表示項目から「由来」を外す

#### 差分 2.2.3: 辞書一覧各行の「由来」表示

- 画面設計 `[8] 辞書一覧`: 各行に `訳語、原文、カテゴリ、由来、ID` の 5 列を表示する旨を記載
- 契約: 各行に表示できる field は `MasterDictionaryEntry` の `id、source、translation、category、updatedAt` のみ
- 判定: **画面設計が古い**
- 是正方向: 各行から「由来」を外し、見出しと同じ 4 列構成に揃える

#### 差分 2.2.4: 新規登録・更新モーダルの「由来」入力

- 画面設計 `[11] 新規登録・更新モーダル`: 入力項目として `原文、カテゴリ、由来、訳語` を表示する
- 契約: `MasterDictionaryEntryPayload` の field は `source、translation、category` の 3 つで、`origin` を含まない
- 判定: **画面設計が古い**
- 是正方向: モーダル入力項目から「由来」を外す。E2E 固定 selector の `master-dictionary-entry-origin-input` も併せて削除候補

### 2.3 あるべき形式 vs 現実装（現実装が古い）

#### 差分 2.3.1: 境界 DTO に `origin` field が含まれる

- 現実装: `MasterDictionaryEntrySummary` / `MasterDictionaryEntryDetail` の Go DTO と TypeScript gateway type に `origin` field を持つ。payload にも `origin` field が含まれる
- 契約: `origin` は境界 API から外す
- 判定: **現実装が古い**
- 是正方向: backend Go DTO、frontend gateway type、payload 系 DTO から `origin` field を削除。DB schema の `origin` 列の扱いは別途検討

#### 差分 2.3.2: 境界 DTO に `note` field が含まれる

- 現実装: `MasterDictionaryEntryDetail` の Go DTO と TypeScript gateway type に `note` field を持つ。`toEntryDetailDTO` が固定文言 `"マスター辞書エントリ"` を設定する
- 契約: `note` は境界 API から外す
- 判定: **現実装が古い**
- 是正方向: backend Go DTO、frontend gateway type、`toEntryDetailDTO` 内の固定文言設定を削除。DB schema の `note` 列の扱いは別途検討

#### 差分 2.3.3: 境界結合テスト 20 件と golden 11 件が現実装ベース

- 現実装: wave-1 で固定した golden 11 件は `note`、`origin` field を含む値で書かれている。wave-2 で固定した backend / frontend 境界結合テスト 20 件はその golden を assert する
- 契約: golden と test code は本契約に追従させる
- 判定: **現実装ベースの test が残る**（あるべき形式への追従が必要）
- 是正方向: golden 11 件から `note`、`origin` field を削除、test code 20 件から該当 field の assert を削除

#### 差分 2.3.4: Summary と Detail が別型として実装されている

- 現実装: `MasterDictionaryEntrySummary` と `MasterDictionaryEntryDetail` を別 type として持つ
- 契約: `MasterDictionaryEntry` の 1 type にまとめる（あるべき仕様で `note` を外したため、Summary と Detail が同形になる）
- 判定: **現実装が古い**
- 是正方向: 2 type を 1 type に統合する。または別 type のまま定義を揃える運用にする（実装判断）

### 2.4 差分判定対象外（契約の責務範囲外）

| 画面設計の要素 | 内容 | 判定 |
|---|---|---|
| `[2]` 上部状態パネルの `gateway 接続状態` | presentation 層で frontend 内部から計算する状態 | 対象外 |
| `[2]` 上部状態パネルの `エラーメッセージ` | エラー応答の message を frontend が整形して表示する | 対象外 |
| `[3]` XML 入力パネルの `XML 待機説明` | UI 上の説明文 | 対象外 |
| `[4]` XML 取り込み進行状況パネルの `進捗バー` | API は同期的呼び出しで進捗イベントを発しない | 対象外 |
| `[6]`〜`[10]` の `選択サマリ`、`詳細状態文` などの文言 | presentation 層が状態から導出する表示文 | 対象外 |
| `[9]` ページ操作領域の `前の30件`、`次の30件` | `pageSize=30` default と整合 | 整合（差分なし） |

### 2.5 突合まとめ

| 系統 | 差分件数 | 是正対象 |
|---|---|---|
| あるべき形式 vs 画面設計 | 4 件（2.2.1〜2.2.4） | `master-dictionary.md` から「メモ」「登録元」「由来」を全箇所外す |
| あるべき形式 vs 現実装 | 4 件（2.3.1〜2.3.4） | backend Go DTO / frontend gateway / `toEntryDetailDTO` から `note` / `origin` を外す。Summary / Detail を統合。本 task の golden 11 + test 20 も追従 |

是正は本 task では行わない。別 task として:
1. 画面設計の更新（`storybook-module` 起点、4 件）
2. 現実装と golden / test の追従（`design-module` → `implementation-module` 起点、4 件）

---

## 3. 観察項目（別 task 入力）

### 3.1 結合テストの導出ルール

`boundary-api.master-dictionary.contract.ts` を「形式の真」として、境界結合テストを下位とする「導出関係」を skill 上に固定するかを別 task で判断する。
observation-record 観察項目 7 で固定済みの暫定推奨と整合する。

### 3.2 意味の置き場の責務分離

本 task で「意味は画面仕様 / UC で確定し、API 仕様書は形式のみ」という方針が user 承認のもとで確定した（2026-06-06）。
この方針を skill 設計（design-module / implementation-module / detail-spec-design）に反映するかを `workflow-contract-maintenance` 起動時に判断する。

### 3.3 Wails 生成 TypeScript binding と契約 file の関係

Wails は backend Go の bind 対象から `frontend/wailsjs/go/...` に TypeScript binding を自動生成する。
本 task では `frontend/src/controller/wails/master-dictionary.contract.ts` を手書きの真として上位に置き、Wails 生成は実装側として位置付ける方針を採用した。
この上下関係を skill / 規約として固定するかは `workflow-contract-maintenance` 起動時に判断する。
