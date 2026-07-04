# Task Plan: xtranslator-export

## task 枠

- `task_id`: `xtranslator-export`
- 依頼要約: xTranslator 互換 XML への翻訳結果書き出し機能を新規実装する。`docs/roadmap.md` 項目5、`docs/known-issues.md` 4番が指す残作業であり、業務要件4（`docs/requirements.md` 4番「xTranslator 形式で出力したい」）に対応する。
- 分岐元 branch: `master`
- 分岐元 commit: `ae94c1433429b3542c8103be0dabf78bd7b94daf`
- 作業 branch: `claude/xtranslator-export`

## 現状事実（調査済み）

- `internal/core/termxml` は xTranslator XML の読み込み専用であり、書き出し処理は存在しない（`docs/known-issues.md` 4番）。
- 翻訳結果は `narration`（叙述文・定型句を収容、`db/migrations/0001_init.sql`）と `line`（台詞）に `plugin`・`form_id`・`edid`・`rec`・`field`・`ordinal`・`source`・`dest`・`status` を自己完結で持つ。
- `proper_noun`（固有名、`db/migrations/0006_record_type_translation.sql`）は `UNIQUE(category, source)` の辞書テーブルで、`plugin`・`form_id` 等のレコード位置を持たない。固有名の原文出現位置は `extracted_field`（box='固有名'）側にあり、`proper_noun` と結んで解決する必要がある（`docs/er.md` 332行目に同種の結合パターンの先例あり）。
- frontend（`frontend/src/ui/screens/translation-run/` 配下）を検索したが、xTranslator 書き出しを起動する UI 要素は無い。
- `docs/architecture.md` の component 図（49行目・65行目）は既に `Engine -->|出力| XML` を描いている。

## 完了定義

### 動かす範囲

- 翻訳ジョブ完了後、利用者は画面操作で対象 job の翻訳結果を xTranslator 互換 XML として書き出せる。
- 書き出し対象は `narration`（叙述文・定型句を収容）と `line`（台詞）の各行、および固有名（box='固有名' の `extracted_field` 行を `proper_noun` の確定訳で解決した行）である。
- 出力される XML は plugin ごとに 1 ファイルとし、各 `<String>` 行が `EDID`・`REC`・`FIELD`・`FORMID`・`Source`・`Dest`・`Status` を持つ（`docs/references/xtranslator_ref.md` 3節）。
- `Status` は DB 内部の訳状態（`narration.status` / `line.status` / `proper_noun.status`）を xTranslator の Status コード（0〜4、同 3.4節）へ写像した値である。
- 書き出し実行後、利用者は画面上で成功（出力先パス）または失敗の結果を確認できる。

### 観測点

- 単体テスト: xTranslator XML 生成の純粋ロジック（`<String>` 行の組み立て、Status 写像）をユニットテストで検証する。
- 実画面: chrome-devtools で書き出し操作を実行し、生成された XML ファイルの中身が対象 job の narration・line・固有名の翻訳結果と一致することを確認する。

### goal 整合

- goal（業務要件4）は「既存の翻訳ワークフローにそのまま乗せられる」ことを求めるため、完了定義はダミー出力や空ファイルでの達成を認めない。実データ（narration・line・固有名の実際の翻訳結果）が XML へ反映されることを要求する。
- goal と除外範囲（含まない）の間に矛盾は無い。

## 軽 / 重判定

- 画面が動くか: `Y`。書き出し操作の起点と結果表示を画面へ追加する必要がある。既存 UI に該当要素は無いことを確認済み。
- `docs/architecture.md` への反映が要るか: 現時点の判断は `N`。書き出し処理を `internal/core/` 配下の新規純粋 package として追加しても、層構成・依存方向・Wails 境界は変わらない見込みである。最終判断は `design-module` で確定する。
- 判定結果: 重 task。`design-module` → `storybook-module`（画面が動くため）→ `implementation-module` → `finalization-module` の順で進める。

## 設計確定（design-module、人間設計レビュー承認済み）

### decision table 結果

- `設計差分図`: 不要。`docs/architecture.md` 反映が不要なため。
- architecture.md 反映: 不要（`N` 確定）。書き出しは既存層（`engine`・`core/termxml`・`store`・`api`・frontend）への機能追加のみで、層構成・依存方向・Wails 境界が変わらない。§2 図（65行目 `Engine -->|出力| XML`）と §3 engine 責務（91行目「xTranslator XML を出力する」）に既に設計目標として描かれている。

### 人間が確定した出力仕様

- 出力対象行: 訳出済みのみ。`narration` と `line` は `status != 0`、固有名は `proper_noun.status` が確定（`1` 訳済 または `3` 仮訳）の行だけを出力する。未訳（`status = 0`）の行は書き出さない。
- 書き出し単位: 単一 plugin で完結する。複数 plugin の一括書き出しは今回スコープ外とする。翻訳フローが plugin 単位で、開発起動では中心 DB が 1 回の実行分に保たれるため、対象は DB にある単一 plugin とする。
- Status（訳状態）: DB の status 値をそのまま xTranslator の Status コードへ出力する（恒等）。DB status は設計上すでに xTranslator コード（`0` 未訳・`1` 訳済・`3` 仮訳）として定義済み（`internal/engine/engine.go:34` `statusProvisional = 3`、`internal/engine/proper_noun.go:13` `statusTranslated = 1`）。範囲外の値は防御的に検証する。訳出済みのみ出力するため、実際に現れる値は `1` と `3` になる。
- 出力ファイル名: `<plugin 名>_english_japanese.xml`（既存 base 辞書 `dictionaries/xTranslatorXMLs/*_english_japanese.xml` の命名に合わせる）。
- ルート要素: `SSETranslator`（Skyrim SE 向け）。
- 出力先: フォルダ選択ダイアログで出力先フォルダを 1 回選び、その中へファイルを書く（`internal/api/app.go:224` `SelectPluginFile` と同じ Wails runtime ダイアログ方式）。

### 固有名の位置解決（確定）

- `proper_noun`（訳の辞書、レコード位置なし）と `extracted_field`（box='固有名' の行、レコード位置あり）を結ぶ。
- 結合キーは `proper_noun.category = extracted_field.rec` かつ `source`（原文）一致。`category` を `rec` で絞るため同綴り異義（人名と地名など）を取り違えない。
- 先例は `internal/store/mention.go:58` `LinkNarrationDescribed`（同一の結合キー構造）。box='固有名' 判定は `record_type_master`（`db/migrations/0006_record_type_translation.sql:70-110`）を `rec`・`field` で引く。
- 対象種別は box='固有名' 全件（一部種別へ絞らない）。

## 実装範囲（implementation-module へ渡す）

scope の境界と依存だけを固定する。詳細な仕様列挙はしない（Claude 本体が文脈を保つため）。

- `internal/core/termxml`（純粋・新規関数）: xTranslator の `<String>` 行構造体、行の組み立て、XML 直列化（ルート要素 `SSETranslator`、要素 `EDID`/`REC`/`FIELD`/`FORMID`/`Source`/`Dest`/`Status`）、Status 値の範囲検証。ファイル入出力は持たない（既存 termxml 方針を踏襲）。新規 package は作らず既存 `termxml` に統合する。
- `internal/store`（新規読み出し）: 単一 plugin の `narration`・`line` を訳出済み（`status != 0`）で読む関数、固有名の位置解決クエリ（`extracted_field` box='固有名' × `proper_noun` を上記結合キーで結び、確定訳のみ返す）。既存読み出しは全件のみで plugin 絞りを持たないため新規に足す。
- `internal/engine`（新規手続き）: 書き出しの orchestration。store から 3 系統を読み、`core/termxml` で `<String>` 行へ組み立て・直列化し、`os` で plugin のファイルへ書く。engine は store・provider・os の IO を束ねる層のため os 書き込みは責務内。
- `internal/api`（新規 Bind）: `ExportXTranslatorXml`（仮名）。出力先フォルダ選択ダイアログ（`OpenDirectoryDialog`、新規）を開き、engine の書き出し手続きを呼び、結果（出力先パスまたは失敗）を返す。
- `internal/bootstrap`: `api.New` のシグネチャ変更が要る場合のみ配線更新。engine に手続きを足す形なら既存の engine 参照で足りる見込み。
- frontend（storybook-module と implementation-module で分担）: `ResultsPanel.svelte` の結果一覧ヘッダに書き出しボタン（結果 1 件以上で有効）。`TranslationRunContainer.svelte` に書き出し状態とハンドラ。`frontend/src/gateway/translation-gateway.ts` に薄いラッパ関数。`ResultsPanel.stories.ts` と `translation-run.fixtures.ts` に story・fixture。

## テスト設計（implementation-module へ渡す）

- 単体テスト（純粋・カバレッジ 100% 基準）: `core/termxml` の `<String>` 行組み立て、XML 直列化（要素順・XML エスケープ・ルート要素）、Status 検証。副作用なしの純粋関数として書く。
- 単体テストで書かない対象: `store` の固有名解決クエリと plugin 読み出し（SQLite アクセス込みのため）、`engine` の orchestration、`api` の Bind とダイアログ、frontend の状態・配線。いずれも実画面 E2E に委ねる。
- 実画面（chrome-devtools）: 書き出し操作を実行し、生成 XML の中身が対象 plugin の narration・line・固有名の訳と一致することを確認する。ネイティブのフォルダ選択ダイアログ操作は人間へ依頼する。

## Storybook 実装（storybook-module、人間レビュー承認済み）

### 合意済み frontend 保護

- 承認済み画面表示: 結果一覧パネル（`ResultsPanel`）のヘッダ右に「xTranslator へ書き出し」ボタンを追加。
- 表示規則:
    - ボタンは結果が 1 件以上（`results.length > 0`）のときだけ表示する。空状態（未実行・実行中）では出さない。
    - 書き出し中（`exporting = true`）はボタンを無効化し、スピナーと「書き出し中…」を出す。通常時の文言は「xTranslator へ書き出し」。
    - 配置は「結果一覧」見出しの右、件数バッジの隣。
- 反映先ファイル:
    - `frontend/src/ui/screens/translation-run/ResultsPanel.svelte`（props `onExport`・`exporting` を追加、ヘッダにボタンを追加）。
    - `frontend/src/ui/screens/translation-run/TranslationRunScreen.svelte`（props `onExportXml`・`exporting` を受け `ResultsPanel` へ流す配線のみ）。
    - `frontend/src/ui/screens/translation-run/ResultsPanel.stories.ts`（story `結果あり（書き出し可能）`・`書き出し中` を追加）。
- fixture 追加はしない。書き出しボタンの表示状態はコンポーネント story（`ResultsPanel.stories.ts`）で確認でき、画面 fixture（`translation-run.fixtures.ts`）への重複追加は Storybook 規約（コンポーネント story で確認できる状態は画面へ重複させない）に沿って見送った。
- 通常分類へ戻した story: `UI Components/ResultsPanel`（レビュー中は作業中分類 `Review/Changed Components/ResultsPanel`）。
- 確認済み Storybook 状態・承認・検証結果は `storybook-review-loop.md` のとおり。
- 変更禁止範囲（後続の実装で表示を変えない境界）: 上記 svelte コンポーネントの表示構造・props 形・ボタン文言・style。implementation-module は表示を変えず、`onExport`／`onExportXml`／`exporting` に state と処理を接続するだけにする。

## 実装（implementation-module）

### 変更ファイル

backend:

- `internal/core/termxml/export.go` 追加: xTranslator `<String>` 行の型、XML 直列化（SSETranslator ルート）、Status 範囲検証。純粋・ファイル入出力なし。
- `internal/core/termxml/export_test.go` 追加: 書き出しの単体テスト（要素順・エスケープ・Status 検証・空）。
- `internal/model/proper_noun_placement.go` 追加: 固有名の位置解決結果（extracted_field × proper_noun）の型。
- `internal/store/export.go` 追加: plugin 単位の訳出済み読み出し（narration・line）と固有名の位置解決クエリ、書き出し対象 plugin の列挙。
- `internal/engine/export.go` 追加: 書き出し手続き（store から読み→`<String>` 組立→termxml で直列化→os でファイル出力）。
- `internal/engine/engine.go` 変更: `Store` interface に `ExportStore` を追加。
- `internal/engine/engine_test.go` 変更: `fakeStore` に `ExportStore` の空実装を追加（書き出し経路は DB 込みのため単体では検証せず E2E に委ねる）。
- `internal/api/export.go` 追加: `ExportXTranslatorXml` Bind（出力先フォルダ選択と結果通知をネイティブダイアログで行う）。

frontend:

- `frontend/src/gateway/translation-gateway.ts` 変更: `exportXTranslatorXml` ラッパを追加。
- `frontend/src/ui/screens/translation-run/TranslationRunContainer.svelte` 変更: `exporting` 状態と `onExportXml` ハンドラを追加し `TranslationRunScreen` へ配線。
- `frontend/wailsjs/go/api/App.d.ts`・`App.js`・`models.ts` 変更: `wails generate module` による再生成（`ExportXTranslatorXml` を追加）。

### 実装中の設計判断（結果通知の場所）

- 書き出し結果（成功の出力先・件数、失敗の原因）はネイティブの MessageDialog で通知する。出力先フォルダ選択がネイティブダイアログである流れと一貫し、svelte 表示コンポーネントを一切変えないため storybook-module の承認範囲（書き出しボタンのみ）を侵さない。frontend は `exporting` 状態だけ管理する。
- 完了定義「画面上で成功（出力先パス）または失敗を確認できる」を、このネイティブダイアログで満たす。

### 最終検証

- backend: `npm run verify:backend`（go test ＋ arch-lint ＋ 境界走査）通過（exit 0）。
- frontend: `npm --prefix frontend run test`（vitest）通過（exit 0）、`npm --prefix frontend run lint`（eslint＋tsc＋knip＋境界）通過（exit 0）。`npm --prefix frontend run check`（svelte-check）は既存の `node_modules/@storybook/svelte` 型宣言エラー 1 件のみで、本 task の変更部分にエラーは無い。
- 単体（純粋ルール）: `internal/core/termxml` の `MarshalStrings` は到達可能な全分岐（Status 検証成功/失敗・正常直列化・空）をテストでカバー。未到達は `xml.MarshalIndent` の防御的エラー返し 1 行のみで、string/int フィールドだけの struct では技術的に到達不能。

### 実画面確認（完了）

実 app（`dev:wails:run`、中心 DB は起動時 wipe）で「Innocence Lost - Quest Expansion.esp」を LM Studio（`hy-mt2-7b`）で翻訳（197 件）し、書き出しボタンから `dictionaries` フォルダへ書き出して確認した。

- 生成 XML の中身が対象 plugin の翻訳結果と一致することを確認: 叙述文（QUST:CNAM/NNAM・MGEF:DNAM・MESG:DESC）、台詞（INFO:NAM1・DIAL:FULL・INFO:RNAM）、固有名（CELL:FULL「Honorhall Orphanage → オナーホール孤児院」・SPEL:FULL・QUST:FULL・MGEF:FULL・ACTI:FULL・WRLD:FULL・ALCH:FULL）が訳文つきで出力された。
- 書き出しボタンの表示条件（結果 0 件で非表示、結果ありで表示）が実画面でも正しく動くことを確認。
- 訳状態: 確定訳（権威訳流用）7 件は Partial 属性なし、AI 仮訳 190 件は Partial="1"。
- 生成 XML が本物の base 辞書（`dictionaries/xTranslatorXMLs/*_english_japanese.xml`）と同一形式（SSTXMLRessources、BOM、Params/Content、String[List/sID/Partial]）であることをファイル比較で確認。

### 実画面確認で判明した仕様相違と是正（当初形式が誤り）

- 当初は `docs/references/xtranslator_ref.md`（xTranslator README 由来）の SSETranslator 形式（ルート SSETranslator、要素 EDID/REC/FIELD/FORMID/Source/Dest/Status）で実装したが、実物の xTranslator エクスポート（この repo が読み込む base 辞書）と食い違うことが実画面確認で判明した。
- xTranslator ソース `TESVT_XMLFunc.pas`（GitHub MGuffin/xTranslator）で実物形式を確認し、次へ是正した:
    - ルート `SSTXMLRessources`、`Params`（Addon・Source=english・Dest=japanese・Version=2）、`Content` ラッパ。
    - `String` 属性 `List="0"`・`sID`（6桁hex連番）・`Partial`（訳状態: 属性なし=確定訳、"1"=未完了訳、"2"=ロック訳）。
    - 子要素 `EDID`・`REC`（"REC:FIELD" 結合）・`Source`・`Dest`。`FIELD`・`FORMID`・`Status` 要素は廃止。
    - 先頭 BOM ＋ `encoding="UTF-8" standalone="yes"`。
- 人間確認事項の是正: 「Status を DB 値そのまま出力」（SSETranslator 前提）→「DB status を Partial 属性へ写像」（確定訳=省略、仮訳="1"）。ルート要素は SSETranslator → SSTXMLRessources。
- `docs/references/xtranslator_ref.md` の SSETranslator 記述は実物と異なるため、実物形式（SSTXMLRessources）へ訂正する。

## 後続モジュールへの引き継ぎ

- 実装・自動検証・実画面確認はすべて完了。次は `finalization-module` へ進む。
- `docs/architecture.md` への反映は不要（`N` 確定）。層・依存・Wails 境界は不変。
- ただし `docs/references/xtranslator_ref.md` は README 由来の誤形式（SSETranslator）だったため、実物形式（SSTXMLRessources）へ訂正する（finalization で正本反映）。
- 書き出した実データ XML（`dictionaries/Innocence Lost - Quest Expansion_english_japanese.xml`）は成果物で、commit 対象に含めない（削除済み）。

## 正本化判断（finalization-module）

- `docs/architecture.md` への反映: 不要（design-module で `N` 確定、構造不変）。層構成・依存方向・Wails 境界は変わらない。§2 図・§3 engine 責務に既に「xTranslator XML を出力する」が描かれている。
- 恒久仕様の architecture 正本反映: なし。
- `docs/references/xtranslator_ref.md` の訂正（README 由来の誤形式 → 実物 SSTXMLRessources 形式）は参照資料の是正であり、architecture 正本反映の対象外。作業 commit に含める。
- 人間承認: architecture 反映が無いため正本反映の人間承認は不要（構造不変）。

## finalization 結果

- 作業 commit: `5b032cba`（16 files changed、branch `claude/xtranslator-export`）。
- local merge: `master` へ `git merge --no-ff`。merge commit `06b077c7`。conflict なし。
- merge 後検証: `npm run verify:backend` 通過（exit 0）、`npm --prefix frontend run test` ＋ `npm --prefix frontend run lint` 通過（exit 0）。
- completed 移動: `docs/exec-plans/active/xtranslator-export/` → `docs/exec-plans/completed/xtranslator-export/`。
- merge 結果 commit: 本 completed 移動と finalization 記録を master へ commit する。
- remote への push は行わない。
