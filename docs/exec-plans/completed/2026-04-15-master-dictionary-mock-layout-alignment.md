# Work Plan Template

- workflow: work
- status: planned
- lane_owner: orchestrate
- scope: master dictionary frontend layout, spacing, alignment, and sticky toolbar behavior around filter, list, and detail card areas
- task_id: 2026-04-15-master-dictionary-mock-layout-alignment
- task_catalog_ref: N/A
- parent_phase: N/A

## Request Summary

- `docs/mocks/master-dictionary/index.html` に合わせて master dictionary 画面の差分を詰める。
- 特に filter、一覧、詳細 card の余白と配置を mock 寄せにする。
- sticky toolbar の挙動を alignment scope に含める。
- task-local UI artifact は `frontend/src/ui/App.test.ts` の固定契約を近実装レベルで保持する。
- 振る舞い追加ではなく visual alignment を主目的にする。

## Decision Basis

- user request は既存画面と mock の差分解消を明示している。
- human は sticky toolbar を scope に含める判断を済ませている。
- review 指摘は DOM contract の欠落と list-shell 構造の不一致を示している。
- frontend を含むため close 前に `ui-check` と `implementation-review` が必要。

## Task Mode

- `task_mode`: fix
- `goal`: master dictionary 画面の filter、一覧、詳細 card を mock に近い余白と配置へ揃え、sticky toolbar と tested DOM contract を固定する。
- `constraints`: `docs/` 正本は human 先行更新のみ。今回の変更は product 側 UI 実装に限定する。既存の機能仕様は変えない。
- `close_conditions`: mock と主要差分が解消される。sticky toolbar scope と DOM contract が active plan と task-local UI artifact に固定される。`review_mode: ui-check` と `review_mode: implementation-review` が pass する。必要 validation が通る。

## Facts

- user は比較対象として `docs/mocks/master-dictionary/index.html` を指定している。
- 差分対象は filter、一覧、詳細 card に加えて sticky toolbar の挙動に限定されている。
- target mock は `toolbar`、`column-row`、`list-stack`、`pager-shell`、`detail-title` を独立 panel として持つ。
- review 指摘では `list-shell` が `column-row`、`list-stack`、`pager-shell` を内包する構造で揃っていない。
- current product page は `frontend/src/ui/screens/master-dictionary/MasterDictionaryPage.svelte` に集中している。
- `frontend/src/ui/App.test.ts` は `h3` 見出し `辞書一覧`、`searchbox` 名 `検索`、`combobox` 名 `カテゴリ`、`#detailTitle`、主要 action button、pager button、一覧 row button の可視名を参照している。
- `python3 scripts/harness/run.py --suite structure` は 2026-04-15 に通過した。

## Functional Requirements

- `summary`:
  - master dictionary 画面の visual layout を、mock の internal panel 構造に寄せて再整列する。
  - sticky toolbar を alignment scope に含める。
  - task-local UI artifact は `App.test.ts` の契約を追跡できる近実装 DOM を持つ。
- `in_scope`:
  - list 側に mock 相当の `toolbar` wrapper を持たせ、heading、action、filter controls のまとまりを独立させる。
  - desktop 幅では `toolbar` を sticky にし、mobile breakpoint では static に戻す前提を scope に含める。
  - `list-shell` は `column-row`、`list-stack`、`pager-shell` をこの順で内包する一貫構造に揃える。
  - `toolbar` 内では `h3#listHeading`、`#listHeadline`、`button#createButton`、`#pageStatusText`、`label[for="searchInput"]`、`input#searchInput`、`label[for="categorySelect"]`、`select#categorySelect` を追跡可能にする。
  - `list-stack` 内では row selection を `button` として保持し、可視名から選択中エントリを追える前提を残す。
  - detail 側に mock 相当の `detail-title` wrapper を持たせ、`#detailTitle` と translation を metadata card 群から分離する。
  - detail header では `button#editButton` と `button#deleteButton` の文言を固定する。
  - `pager-shell` では `button#prevPageButton` と `button#nextPageButton` の文言を固定する。
  - `content-grid` の左右比率と各 panel 間 gap を、list 優位かつ detail 補助の見え方へ調整する。
  - responsive で 980px 近辺以下に落ちた時の stack、`column-row` 非表示、row の 1 列化を維持する。
- `non_functional_requirements`:
  - 変更は `MasterDictionaryPage.svelte` の layout / style 局所修正に閉じる。
  - controller contract、view model contract、routing、CRUD / import behavior は不変とする。
  - global spacing token の変更は避け、必要なら component local の wrapper と spacing で吸収する。
  - task-local UI artifact は review 用抽象図ではなく、`App.test.ts` の selector と role を追える near-implementation DOM とする。
  - review 時に mock と product の差分説明が active plan と task-local artifact だけで追える状態にする。
- `out_of_scope`:
  - XML 取り込み導線の redesign。
  - hero、gateway status、modal、CRUD 文言の仕様変更。
  - 新規 field 追加、情報構造変更、interaction 追加。
  - docs 正本更新、mock canonical source 更新、product test 更新。
- `open_questions`:
  - Wails 実画面と mock の side-by-side runtime comparison は未実施であり、`ui-check` で最終確認が必要。
  - sticky toolbar の停止位置と top offset は task-local artifact で `top: 18px` として固定するが、実装時に header overlap が出ないかを review で確認する必要がある。
  - global token の影響で数 px の残差が残る場合、component local override を優先してよいかを review で最終確認する。
- `required_reading`:
  - /Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-master-dictionary-mock-layout-alignment.md
  - /Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-master-dictionary-mock-layout-alignment.ui.html
  - /Users/iorishibata/Repositories/AITranslationEngineJP/docs/mocks/master-dictionary/index.html
  - /Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-dictionary/MasterDictionaryPage.svelte
  - /Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/App.test.ts

## Fixed DOM Contract

- `h3#listHeading` の見出し名は `辞書一覧` で固定する。
- `label[for="searchInput"]` と `input#searchInput[type="search"]` の組み合わせで、role `searchbox` の名前 `検索` を固定する。
- `label[for="categorySelect"]` と `select#categorySelect` の組み合わせで、role `combobox` の名前 `カテゴリ` を固定する。
- `#detailTitle` は選択中エントリの source を描画する固定 anchor とする。
- `button#createButton` の可視名 `新規登録` を固定する。
- `button#editButton` の可視名 `更新` を固定する。
- `button#deleteButton` の可視名 `削除` を固定する。
- `button#prevPageButton` の可視名 `前の30件` を固定する。
- `button#nextPageButton` の可視名 `次の30件` を固定する。
- 一覧 row は `button` のまま保持し、可視名から対象訳語を辿れる前提を維持する。

## Artifacts

- `ui_artifact_path`: /Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-master-dictionary-mock-layout-alignment.ui.html
- `final_mock_path`: /Users/iorishibata/Repositories/AITranslationEngineJP/docs/mocks/master-dictionary/index.html
- `scenario_artifact_path`: N/A
- `final_scenario_path`: N/A
- `implementation_scope_artifact_path`: /Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-master-dictionary-mock-layout-alignment.implementation-scope.md
- `review_diff_diagrams`: N/A
- `source_diagram_targets`: N/A
- `canonicalization_targets`:
  - /Users/iorishibata/Repositories/AITranslationEngineJP/docs/mocks/master-dictionary/index.html
  - この task では canonical source を更新しない。target path は review 時の参照先としてのみ固定する。

## Work Brief

- `implementation_target`: /Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-dictionary/MasterDictionaryPage.svelte
- `accepted_scope`:
  - list 領域の heading / action / filter を一体に見せる `toolbar` wrapper を追加または再編する。
  - `toolbar` の sticky behavior を desktop scope に含め、mobile breakpoint で static へ戻す。
  - list 領域の `list-shell` を `column-row`、`list-stack`、`pager-shell` の三層で再編する。
  - detail 領域の tag、`#detailTitle`、translation を `detail-title` 相当 block へ寄せ、metadata card と detail list の前段に置く。
  - 左右カラム幅、panel padding、内部 gap、border radius、背景レイヤーを mock 寄せで局所調整する。
  - `App.test.ts` が参照する heading level、label 名、button 名、id anchor は変更対象に含めず維持する。
  - import section、controller wiring、product test は変更対象に含めない。
- `parallel_task_groups`: none
- `tasks`:
  - mock と現行 UI の panel 構造差分を task-local UI mock で固定する。
  - `MasterDictionaryPage.svelte` の DOM を wrapper 単位で再配置する。
  - component local style を使って filter、list、detail、sticky toolbar の spacing / alignment を再調整する。
  - responsive 崩れと tested selector 維持を確認したうえで `ui-check` と `implementation-review` へ渡す。
- `implementation_brief`:
  - list 側は `shell-card` 直下に要素を並べる構成をやめ、mock の `toolbar` に相当する内側 panel を先頭に置く。
  - `toolbar` は `position: sticky` と `top: 18px` を持つ想定で扱い、980px 以下では static に戻す。
  - filter controls は heading / action と視覚的に同じ panel に収める。
  - `list-shell` は `column-row`、`list-stack`、`pager-shell` の三層を内包し、header と rows の左右基準線を揃える。
  - detail 側は tags、`#detailTitle`、translation を `detail-title` 相当 block にまとめ、その下に metadata card と detail list を続ける。
  - `content-grid` は 1:1 ではなく list を広く、detail を狭く見せる列比率へ変更する。
  - `h3#listHeading`、`#searchInput`、`#categorySelect`、`#detailTitle`、主要 action button、pager button の文言と role は固定する。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite structure`
  - `python3 scripts/harness/run.py --suite all`
  - `npm run dev:wails:docker-mcp`

## Investigation

- `reproduction_status`: static structure comparison completed; side-by-side runtime comparison pending
- `trace_hypotheses`:
  - current frontend implementation の `content-grid` が 1:1 であり、mock の list 優位レイアウトとずれている可能性が高い。
  - current frontend implementation は `filter-grid` と detail title 周辺に dedicated panel がなく、gap の基準点が mock とずれている。
  - current frontend implementation は `list-shell` 内の block 分離が不足し、header、rows、pager の alignment が mock より曖昧に見える。
  - sticky toolbar の scope を明記しないままだと、implementation-review で top offset と breakpoint 条件が再解釈される可能性がある。
- `observation_points`:
  - filter wrapper と heading / action の結合方法。
  - sticky toolbar の top offset と mobile での解除条件。
  - `list-shell` 内の header row、row stack、pager の上下関係。
  - detail title block、metadata card、detail list の間隔。
  - mobile breakpoint での column collapse と row readability。
- `residual_risks`:
  - global token の影響で mock との差が完全一致にならない可能性がある。
  - wrapper 追加時に tested selector の近接関係が変わり、想定外の DOM query へ影響する可能性がある。
  - sticky toolbar を追加すると、狭い desktop 幅で detail 側の先頭との視覚基準がずれる可能性がある。

## Acceptance Checks

- `toolbar` が heading / action / filter を同じ内部 panel に収め、desktop では sticky、mobile では static であることが task-local artifact 上で読める。
- `list-shell` が `column-row`、`list-stack`、`pager-shell` を内包する構造で統一されている。
- `h3#listHeading`、`#searchInput`、`#categorySelect`、`#detailTitle`、`#createButton`、`#editButton`、`#deleteButton`、`#prevPageButton`、`#nextPageButton` が task-local artifact 上で追跡できる。
- detail 領域に `detail-title` 相当 block があり、tags、title、translation が metadata card と分離される。
- list が detail より広い比率で表示され、desktop と mobile 相当で破綻しない。

## Required Evidence

- task-local UI mock と active plan の整合。
- `review_mode: ui-check` の確認結果。
- `review_mode: implementation-review` の確認結果。
- `python3 scripts/harness/run.py --suite all`
- 必要に応じて `http://host.docker.internal:34115` 上の side-by-side visual comparison 結果。

## HITL Status

- `functional_or_design_hitl`: required-after-plan
- `approval_record`: approved by human on 2026-04-15 to proceed with implementation after design bundle review

## Closeout Notes

- `canonicalized_artifacts`: N/A

## Outcome

- review 指摘に合わせて requirements mode と implementation-brief mode を修正した。
- sticky toolbar を alignment scope に含めることを active plan に固定した。
- `App.test.ts` の fixed DOM contract を active plan に明文化した。
- scenario は N/A のまま維持し、implementation-scope artifact を active path へ固定した。
