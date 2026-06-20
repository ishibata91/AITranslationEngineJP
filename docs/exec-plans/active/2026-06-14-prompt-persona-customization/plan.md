# Task Plan: 2026-06-14-prompt-persona-customization

- `workflow`: work
- `status`: implementation-module 通過（縦切り 3 段の backend 本番接続・frontend 配線・実画面確認まで完了。finalization 待ち。2026-06-20）
- `task_id`: 2026-06-14-prompt-persona-customization
- `task_mode`: 重 task（画面が動く: Y、`docs/architecture.md` 反映: 要・小。プロンプト構築の所在を §3 で 1 行明確化、完了時に §8 更新）
- `request_summary`: プロンプトテンプレートの編集、実際に翻訳 AI へ投げるプロンプトの参照、翻訳後の機械置換内訳の表示、口調指示テンプレートの精緻化をできるようにする。
- `goal`: ユーザーがプロンプトテンプレートを翻訳前に編集し、実プロンプトと機械置換内訳を目視確認でき、口調指示テンプレートが実プロンプトへ反映される。
- `constraints`: 翻訳エンジン本体の手続き（T1〜T3）の入出力は変えない。engine が機械置換語（used）を保持する拡張は許す。編集と参照の面を足す。
- `close_conditions`:
  - プロンプトテンプレート編集 UI で雛形を編集でき、保存後の翻訳で反映される（実画面 + 実プロンプト参照）。
  - 実際に翻訳 AI へ投げたプロンプトを UI で目視確認できる（実画面）。
  - 翻訳実行後、各結果行で原語 → 確定訳語の機械置換内訳を確認できる（実画面）。取得時に `dict.Apply` で内訳を再構成する純粋ロジックを単体テストで検証する。
  - 口調指示テンプレートの構造を見直し、口調指示が実プロンプトへ合成される（実プロンプト参照）。
- `source_branch`: `master`
- `source_commit`: `003a8968`
- `target_branch`: `master`

## Scope（含む / 含まない）

含む:
- プロンプトテンプレート（base 翻訳指示文・口調指示テンプレート）の編集 UI。テンプレートを永続化する。
- 実際に翻訳 AI へ投げるプロンプトの参照（目視確認）。
- 翻訳実行後の機械置換内訳（原語 → 確定訳語）の表示。
- 口調指示テンプレートの精緻化（テンプレート構造の見直し）。
- 翻訳実行画面とテンプレート編集画面を行き来する画面間ナビゲーション。
- 画面が動くため storybook-module 経由。

含まない:
- 属性（種族・声型）→ 性質文・ペルソナのルール編集。台詞群からのペルソナ生成と一体で設計し直すため、`2026-06-20-character-persona-from-dialogue` へ切り出す。種族・声型 → 性質文の対応は現状ハードコード（`persona_rule.go`）のまま使う。
- AI によるペルソナ生成・会話履歴解析。新 plan で扱う。
- 実 LLM 出力の自然さの定量評価。口調精緻化は実プロンプトへの口調指示の合成までを完了線とし、訳文の品質は実データ目視の補助観測にとどめる。

## 完了定義（preparation-module で固定。scope 縮小を 2026-06-20 に反映）

システム上どこまで動かすかを以下に固定する。差込点や空テーブルを置くだけで振る舞いが観測できない状態を「動く」とは扱わない。

1. プロンプトテンプレート編集。base 翻訳指示文と口調指示テンプレートを UI で編集でき、保存後の翻訳でその雛形が使われる。
   - 観測点: 実画面（編集 → 保存）＋実プロンプト参照（#2）で雛形反映を確認。
2. 実プロンプト参照。実際に翻訳 AI へ投げたプロンプトを UI で目視確認できる。
   - 観測点: 実画面（翻訳実行後に実プロンプトを表示）。
3. 機械置換内訳。翻訳実行後、各結果行で原語 → 確定訳語の機械置換内訳を確認できる。結果取得時（`ListResultsPage`）に各行の原文へ `dict.Apply` を当て直して内訳を再構成し、API が `ResultView.terms` へ供給する。
   - 観測点: 実画面（結果行に内訳表示）＋単体テスト（`dict.Apply` の純粋ロジックで原文 → 置換内訳を検証）。
4. 口調指示の精緻化。口調指示テンプレートの構造を見直し、種族・声質などの性質文が口調指示として実プロンプトへ正しく合成される。性質文の中身（属性 → 性質文の対応）はハードコードのまま使い、その編集は新 plan で扱う。
   - 観測点: 実プロンプト参照（#2）で口調指示の合成を確認。実 LLM での訳文反映は実データ（LM Studio）目視の補助観測にとどめる。

goal 整合: goal の「テンプレート編集・参照・置換内訳・口調反映」を上記 1〜4 が実際に動くことで満たす。最小実装・空テーブルだけで満たしたことにしない。
除外整合: 属性 → 性質文のルール編集と AI ペルソナ生成は新 plan へ切り出したため、本 plan の「含まない」と矛盾しない。

## 軽 / 重判定（preparation-module で固定。2026-06-20）

- 画面が動くか: Y。テンプレート編集 UI、実プロンプト表示、機械置換内訳表示、画面間ナビゲーションの表示構造と svelte コンポーネントを変える。
- `docs/architecture.md` 反映が要るか: Y（小）。テンプレートの永続化層と Wails binding の新メソッド、engine が used を保持し API が `ResultView.terms` を返す経路、プロンプト構築の provider → engine 移設で、Wails 境界と層構成の一部が変わる。
- 判定結果: 重 task。後続経路は preparation → design-module → storybook-module（画面が動く）→ implementation-module → finalization-module。

## 設計決定（design-module、2026-06-20、人間設計レビュー承認済み）

- テンプレートの保存先: 中心 DB 内の専用テーブルに分離する。抽出データと別にし、起動ごとの中心 DB 消去の対象から外す。
- プロンプト構築の所在: AI クライアント層（`provider`）から翻訳手続き層（`engine`）へ移す。`provider` は完成プロンプトを受け取り送るだけにする。
- 機械置換内訳と実プロンプト: 結果取得時に `dict.Apply` を当て直して再構成する（保存しない）。
- 実プロンプト参照の範囲: 翻訳実行後の結果行で確認する。翻訳前プレビュー画面は本 task に含めない。
- 画面間ナビゲーション: 翻訳実行画面とテンプレート編集画面を行き来する導線を足す。
- 進め方: 縦切り 3 段で進める（① 機械置換内訳供給 → ② 実プロンプト参照・プロンプト構築移設 → ③ テンプレート編集・口調精緻化）。各段は単体で画面から観測できる成果にする。
- `docs/architecture.md` 反映: 層・依存・Wails 境界の構成は概ね不変。§3 のプロンプト構築の所在を 1 行明確化し、§8「現在の状態」を完了時に更新する。finalization-module で反映する。

## 依存

- T2（ペルソナ機構）と T3（マスター辞書）。

## Routing Notes

- `required_reading`:
  - `docs/system_requirements.md`（§3 テンプレートの永続化と翻訳前編集）
  - `docs/concept-model.md`（素材の性質と口調の合成）
  - `docs/UX-standard.md`（UI 設計の正本）
  - `docs/references/storybook.md`（画面実装の作法）

## Outcome

implementation-module で縦切り 3 段の backend 本番接続と frontend 配線を Claude 本体が 1 文脈で実装し、最終検証まで通過した。

### 実装（変更ファイル 1 行 1 ファイル）

backend:
- `internal/engine/engine.go`: 辞書構築を `LoadDictionary` で公開。`Run` が prompt_template を 1 度読み、`ComposePrompt` で完成プロンプトを組む。`LinePersonas` が口調指示テンプレートを引数で受ける。
- `internal/engine/prompt.go`: プロンプト構築の純粋関数 `ComposePrompt`（base 指示＋口調指示＋機械置換済み原文）と表示用の `RenderPrompt` を新設。provider から文面構築を移設。
- `internal/engine/persona.go`: 口調指示を編集可能テンプレートの `{traits}` 差し込み駆動へ見直し（`personaTraits`＋`buildPersonaDirective`）。
- `internal/provider/openai_compatible.go`: `Translate` を完成 `Prompt`（System/User）受け取りへ変更。base 指示定数と `systemPrompt` を撤去。
- `internal/api/app.go`: `ResultView` に `Terms`・`Prompt` を追加。`buildResultsPage` が辞書とテンプレートをページ単位で読み、各行の置換内訳と実プロンプトを再構成。`GetPromptTemplate`・`SavePromptTemplate` を追加。
- `internal/model/prompt_template.go`: `PromptTemplate` 構造体を新設。
- `internal/store/prompt_template.go`: テンプレートの取得・UPSERT を新設。
- `db/migrations/0004_prompt_template.sql`: 単一行（id=1）の prompt_template テーブルと既定値 seed を追加。

frontend:
- `frontend/src/gateway/translation-gateway.ts`: `ResultRow` へ `terms`・`prompt` を追加し写像。
- `frontend/src/gateway/template-gateway.ts`: テンプレート取得・保存のラッパを新設。
- `frontend/src/ui/screens/template-editor/TemplateEditorContainer.svelte`: 編集 state・dirty/saving・保存/戻すの配線を新設。
- `frontend/src/App.svelte`: AppShell＋AppNav で翻訳実行とテンプレート編集を切り替えるルーティングを配線。
- `frontend/package.json`: story 専用 2 ファイル（`TranslationRunNavPreview.svelte`・`template-editor.fixtures.ts`）を knip ignore へ追加。

### テスト（単体）

- `internal/engine/prompt_test.go`: `ComposePrompt`・`RenderPrompt` の純粋関数テストを追加。
- `internal/engine/dictionary_test.go`: 既存の `dict.Apply` テストを terms 供給の検証単位として再利用。
- 既存テスト（engine・provider・api・persona）をシグネチャ変更へ追従。

### 最終検証

- backend: `go build ./internal/... ./db/... .`・`go vet ./internal/...`・`go test ./internal/...` 全通過（engine/api/provider/store）。harness に backend suite が無いため go ツール直実行。
- frontend: `python3 scripts/harness/run.py --suite frontend-local` 通過（eslint・tsc・knip・boundaries・vitest）。`npm run build`（vite production）通過。
- 実画面（`http://localhost:34115`、dev DB の既存 121 台詞・master_term 25071 件を使用）:
  - 機械置換内訳: 結果行に `Grelod → グレロッド`、`Honorhall Orphanage → オナーホール孤児院`、`Grelod the Kind → 親切者のグレロッド` の内訳と固有名チップを確認。
  - 実プロンプト参照: 結果行展開で system（base 指示＋口調指示）と user（機械置換済み原文）の全文を確認。
  - 口調精緻化: system に口調指示テンプレートの `{traits}` へ性質列（声質・種族の気質）が差し込まれることを確認。
  - テンプレート編集→保存→反映: base 指示へ目印文字列を保存し、結果行の実プロンプトへ反映、再読み込みで永続も確認。検証後に既定文へ復元。
  - 画面間ナビゲーション: 翻訳実行とテンプレート編集の往復を確認（往復時に翻訳画面の container は再マウントして先頭ページを再取得する）。

### finalization への引き継ぎ

- 仕様追加・仕様変更: なし（plan の close_conditions の範囲内）。
- `docs/architecture.md` 反映（finalization-module で実施）: §3 のプロンプト構築の所在を provider → engine に 1 行明確化。§8「現在の状態」へ prompt_template の永続化と Wails 新メソッド（GetPromptTemplate/SavePromptTemplate）、ResultView の terms/prompt 供給を追記。
- `docs/changelog.md` への変更・判断履歴の追記（finalization-module で実施）。
- 生成物 `frontend/wailsjs/` は gitignore 対象（`wails generate module` で再生成）。

## Finalization（finalization-module、2026-06-20）

### 正本化判断

- 反映対象: `docs/architecture.md`。
- 判断: §3（プロンプト構築の所在を `provider` → `engine` に明確化、`store` にテンプレート CRUD を追記、`provider` は完成プロンプトを送るだけと明記）と §8（T4 の現在状態：構築移設・`prompt_template` 永続化・Wails 新メソッド・`ResultView` の terms/prompt 供給）を反映。
- 根拠: design-module の設計決定（plan「設計決定」节、人間設計レビュー承認済み 2026-06-20）で「§3 を 1 行明確化し §8 を完了時に更新、finalization で反映」と固定済み。provider port（`Translate`）の契約変更と Wails 公開面への新メソッド追加という実際の境界変化があるため反映は妥当。
- 人間承認状態: 設計レビューで反映方針を承認済み。本反映は承認範囲内の最小反映のため、新規承認サイクル（一時 summary.md）は設けない。一時資料 `summary.md` はレビュー後削除の明記どおり削除する。
- 恒久仕様の追加・廃案: なし（close_conditions の範囲内）。

### 正本反映（docs/architecture.md）

- §3 `engine`: 「翻訳プロンプトの組み立て（base 指示・口調指示・機械置換済み原文の合成）は engine の純粋関数が持ち、完成プロンプトを provider へ渡す」を追記。
- §3 `store`: 「プロンプトテンプレートの CRUD を含む」を追記。
- §3 `provider`: 「engine が組んだ完成プロンプトを受け取って送るだけで、文面構築はしない」を追記。
- §8: T4（2026-06-20）の現在状態（構築移設・`prompt_template` 単一行永続化・Wails 新メソッド・取得時の terms/prompt 再構成・口調テンプレートの `{traits}` 差し込み）を 1 段追記。
- 層・依存・ディレクトリ正本（§4・§7）は不変のため変更なし。

### 作業 commit / local merge / completed 移動 / merge 結果 commit

（以下は実行後に hash を追記する）
