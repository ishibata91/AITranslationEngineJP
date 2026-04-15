# 実装スコープ固定

- `task_id`: `2026-04-15-master-dictionary-mock-layout-alignment`
- `task_mode`: `fix`
- `design_review_status`: `pass`
- `hitl_status`: `approved`
- `summary`: `MasterDictionaryPage.svelte` の list / detail layout と page-local style を mock 寄せで再整列し、sticky toolbar を desktop 限定で固定する。`

## 共通ルール

- `App.test.ts` の DOM contract は変更しない。`h3#listHeading`、`#searchInput`、`#categorySelect`、`#detailTitle`、主要 action button、pager button の見出し名、label 名、id anchor を保持する。
- product code の owned scope は `MasterDictionaryPage.svelte` の list panel、detail panel、`content-grid`、page-local `<style>` に限定する。
- `<script>` は controller call、view model contract、routing、CRUD / import behavior を変えない。wrapper 再配置に伴う id / class 維持の最小調整だけ許可する。
- sticky toolbar は desktop で `position: sticky` と `top: 18px` を前提にし、`max-width: 980px` 以下では `position: static` に戻す。
- import shell、modal、global token、shared component、他 page、mock 正本、product test は変更しない。

## 許可ファイル

- `frontend/src/ui/screens/master-dictionary/MasterDictionaryPage.svelte`

## 実装 handoff 一覧

### `frontend-master-dictionary-layout-alignment`

- `implementation_target`: `frontend`
- `owned_scope`:
  - `frontend/src/ui/screens/master-dictionary/MasterDictionaryPage.svelte`
  - list 側の DOM を `toolbar -> list-shell -> column-row -> list-stack -> pager-shell` の順で固定する。
  - detail 側の DOM を `detail-head -> detail-title -> detail-grid -> detail-list` の順で固定する。
  - `content-grid`、panel padding、gap、row grid、sticky toolbar を page-local style で吸収する。
- `depends_on`:
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-master-dictionary-mock-layout-alignment.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-master-dictionary-mock-layout-alignment.ui.html`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/mocks/master-dictionary/index.html`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/App.test.ts`
- `validation_commands`:
  - `cd frontend && npm test -- App.test.ts`
  - `cd frontend && npm run check`
  - `python3 scripts/harness/run.py --suite all`
  - `npm run dev:wails:docker-mcp`
- `completion_signal`: list 側に sticky `toolbar` と `list-shell` 三層構造が入り、detail 側に `detail-title` wrapper が復元され、`App.test.ts` の DOM contract を壊さず mock 差分の主要箇所が説明可能になる。
- `notes`:
  - 許可する変更は markup と style の再編に限定する。
  - `button#createButton`、`button#editButton`、`button#deleteButton`、`button#prevPageButton`、`button#nextPageButton` の可視名は現状維持とする。
  - `searchbox` 名 `検索` と `combobox` 名 `カテゴリ` は label 経由で維持する。

## 実装順

1. list panel の DOM を `toolbar` と `list-shell` の二段に再編し、heading / action / filter を sticky panel に集約する。
2. `list-shell` を `column-row`、`list-stack`、`pager-shell` の順にそろえ、row と pager の基準線を揃える。
3. detail panel に `detail-title` wrapper を追加し、tag、`#detailTitle`、translation を metadata card 群から分離する。
4. `content-grid` の列比率と breakpoint を mock 寄せで調整し、desktop sticky / mobile static を確認する。
5. `App.test.ts`、full harness、Wails runtime で contract と layout を確認する。

## 明示的な非目標

- `frontend/src/ui/App.test.ts` の更新
- import shell、modal、controller wiring、route、CRUD / import behavior の変更
- global spacing token や shared style の調整
- `docs/mocks/master-dictionary/index.html` と `docs/` 正本の更新
