# generic-voice-tone-fallback（汎用台詞・PC 発話への口調付与）

## 状態

完了（finalization-module で master へ local merge、completed へ移動）。

- 作業 commit: `ed540ff`（branch `claude/generic-voice-tone-fallback`）
- merge commit: `015fe48`（master へ `--no-ff`）
- 設計は `./design.md`（brainstorming 成果物）、確定した実装範囲・テスト設計は `./implementation-scope.md`。

## 分岐情報

- 作業 branch: `claude/generic-voice-tone-fallback`
- 分岐元 branch: `master`
- 分岐元 commit: `cb3aaa2c`

## 依頼要約

会話レコードの口調付与を、話者を解決できた台詞だけでなく全会話レコードへ広げる。親要件は 3 つ。

- 話者を特定できない汎用台詞（衛兵・装備反応など）へ口調を付ける。
- プレイヤー発話（DIAL の選択肢）へ口調を付ける。これを「PC 口調」と呼ぶ。
- どの会話レコードも口調が決まり、「口調なし」が残らない状態にする。

設計判断（brainstorming 確定分）は次のとおり。口調 = 性別（条件由来）＋ 既定の対人段階（利用者指定・固定）＋ 感情段階（本文 1 行から）。対人態度は推測せず固定する。役どころ・勢力 prior・皮肉検出は scope 外。詳細は `./design.md`。

## 背景（観測事実）

backend-violation-cleanup の実画面確認（`dictionaries/Data/Innocence Lost - Quest Expansion.esp`、実 LLM hy-mt2-7b）で観測した。

- 台詞 151 件のうち話者解決は 113 件で、これらには口調が付く。未解決は 38 件で口調なし。
  - DIAL（プレイヤー発話）26 件: 当初は対象外としたが、本 task で PC 口調を付ける対象へ含める。
  - INFO（NPC 応答）で未解決 12 件: すべて衛兵の汎用台詞（護送・逮捕・「他の衛兵の問題だ」など）。汎用台詞 fallback の対象。
- 口調分類器 `tone.Classify(lines, voice)` は voice 気質 prior へ畳む fallback（`DecisionPath="voice"`）を持つが、voice は話者経由でしか渡らず、話者未解決の台詞には届かない。
- spike（`--probe-generic`）で確認: 汎用 INFO のうち話者属性を一切持たない「真の汎用」が Skyrim.esm で 63%、Outfit Recognition Framework.esp で 97%。気質 prior の根拠は薄く、本文 1 行の感情と条件由来の性別が確実に使える信号。

## goal

全会話レコードに対人段階と感情段階を付け、「口調なし」を残さない。汎用台詞・PC 発話・名指し話者のいずれにも口調が決まる。名指し話者の口調は従来どおりで非劣化に保つ。

## 完了定義（動かす範囲と観測点）

### 動かす範囲

- 汎用台詞（INFO・話者なし）に口調が付く。既定の対人段階（固定）＋ 本文 1 行の感情段階 ＋ 条件由来の性別で組み立てる。
- PC 発話（DIAL の選択肢）に PC 既定の口調が付く。PC 既定の対人段階（固定）＋ 本文 1 行の感情段階で組み立てる。
- 条件由来の性別が抽出器から中心 DB へ永続し、汎用台詞の一人称・語尾へ効く。
- どの会話レコードも対人段階と感情段階が埋まり、「口調なし」が残らない。
- 名指し話者の口調（`persona_character`）は従来どおりで非劣化。

### 観測点

- 単体テスト: 1 台詞の `Features` と固定の対人段階から口調を出す純粋ルール（`internal/core/tone`）。守るべき不変ルールとしてユニットテスト 100% カバレッジ。
- 実データ: 抽出器が条件由来の性別を中心 DB へ書くことを、`--sqlite` 取込後の DB 内容で確認する。
- 実画面（`http://localhost:34115`、起動 `npm run dev:wails:run`）:
  - 衛兵の汎用台詞（`Innocence Lost - Quest Expansion.esp` の護送・逮捕など）に口調が付く。
  - OutfitRecognition の装備反応台詞（真の汎用、性別条件のみ一部）に既定の対人段階＋本文由来の感情段階で口調が付く。
  - PC 発話（DIAL の選択肢）に PC 既定の口調が付く。
  - 名指し話者の口調が従来どおりで非劣化。

## close_conditions

- `internal/core/tone` の per-line 口調ルールのユニットテストが通り、カバレッジ 100% を満たす（観測点: go test カバレッジ）。
- 抽出器 `--sqlite` 後の中心 DB に、汎用 INFO の条件由来の性別が入っている（観測点: DB クエリ）。
- 実 app で汎用台詞・PC 発話・名指し話者の口調表示を目視し、「口調なし」が残らないこと、名指し話者が非劣化なことを確認する（観測点: 実画面 `http://localhost:34115`）。

## 軽 / 重判定

- 画面が動くか: **Y**。口調メタ表示の決定経路（DecisionPath）へ汎用・PC の新しい値が出る。既定の対人段階を利用者が指定する編集面を足す場合、表示コンポーネントも増える可能性がある。
- `docs/architecture.md` 反映が要るか: 現時点 **N** 寄り。層構成・依存方向・Bootstrap・Wails 境界の構造は不変。`internal/core/tone` への純粋ルール追加、抽出器の永続追加、取込の橋渡し、engine の routing 追加、store の列/橋渡しテーブル追加はいずれも既存層の中の変更。新しい Wails Bind を足す場合も method 追加に留まる。公開境界が design-module で確定した時に再評価する。
- 判定結果: 画面が動く（Y）のため **重 task**。後続は `design-module` →（画面が動くため）`storybook-module` → `implementation-module` → `finalization-module`。

## 未確定事項の確定（design-module、人間設計レビュー承認済み）

5 件は `./implementation-scope.md` に確定済み。#2・#4 は人間設計レビューで段階選択から自由記述方式へ変更した。要点だけ再掲する。

- 性別の経路: staging `extracted_info_condition` → domain `line_condition`（`line_speaker` と対称）。
- 汎用・PC の口調指定: 段階選択でなく自由記述。`prompt_template` へ汎用口調・PC 口調・PC 性別を列追加。汎用・PC は対人段階・セル名を持たない。
- 激情しきい値: `tone.Classifier` へ 1 行用の渋い値を新設、具体値は本文較正で確定。
- 公開境界: 口調メタ DTO へ汎用・PC の決定経路・感情・性別を載せる。新 Wails メソッドなし、prompt_template payload 拡張のみ。
- PC 判定: PC 発話 = (DIAL, FULL) と (INFO, RNAM)、NPC 返答 = (INFO, NAM1)。

## 合意済み frontend 保護（storybook-module、2026-06-30 承認）

承認済みの story と svelte コンポーネントが画面の正本。後続実装は表示構造・props 形を変えずに配線だけ行う。

### 承認済み画面と story（通常分類）

- 結果行 `UI Components/TranslationResultRow`: 名指し話者は従来どおり（セル＋性質文＋対人/感情/印）。汎用・PC は見出し（汎用台詞／PC発話）＋「感情 ・ 性別」だけを出し、対人段階・セル・印を出さない。
- 編集画面 `Screens/プロンプトテンプレート`: レコード別タブ先頭に「話者なし台詞の口調」節（汎用・PC の自由記述欄 2 つ＋PC 性別の選択 男性／女性／未指定）。

### 表示規則・変更禁止範囲

- 汎用・PC の口調メタは対人段階・セル・印を出さない。決定経路の見出しは `DECISION_PATH_LABEL`（汎用→汎用台詞、PC→PC発話）。
- 口調メタ根拠の組み立ては `personaMetaParts`（名指し＝決定経路・対人・感情・性別・印、汎用/PC＝感情・性別）。
- PC 性別の選択肢は `PC_SEX_OPTIONS`（未指定／男性／女性）。自由記述欄の field は `TONE_TEXT_FIELDS`。

### 反映先 frontend ファイル

- `translation-run-view.ts`（`PersonaMeta` の汎用/PC 用任意フィールド・`DecisionPath` 拡張）。
- `translation-run-presentation.ts`（`DECISION_PATH_LABEL`・`DECISION_PATH_HINT`・`personaMetaParts`）。
- `TranslationResultRow.svelte`、`translation-run.fixtures.ts`、`TranslationResultRow.stories.ts`。
- `template-editor-view.ts`（`PromptTemplateForm` の `genericToneText`・`pcToneText`・`pcSex`、`PcSex`）。
- `template-editor-presentation.ts`、`ToneDefaultPane.svelte`、`TemplateEditorScreen.svelte`、`template-editor.fixtures.ts`、`TemplateEditorScreen.stories.ts`。

### implementation-module へ渡す frontend ロジック残課題

- 編集画面 Container（`TemplateEditorContainer.svelte`）と gateway（`template-gateway.ts`）へ、汎用口調・PC 口調・PC 性別の state 保持・取得・保存配線を足す。`onFieldInput` の handler へ新フィールドの分岐を足す。
- 結果取得（api・gateway）で、汎用・PC 台詞の口調メタ（決定経路・感情・性別）を `PersonaMeta` へ載せる。
- 新フィールドは表示型で任意にしてあるため、配線時に実値を供給する。

## 実装と最終検証（implementation-module、2026-06-30）

### 実装した範囲（1 文脈で縦通し）

- 抽出器（C#）: `InfoConditionSqliteWriter` を追加。INFO の条件（GetIsSex 極性畳み／単一声型の Male/Female 接頭／同性のみ FLST）から性別を導き `extracted_info_condition` へ書く。`Program.cs` から起動。
- スキーマ: migration `0007_generic_tone.sql`。`extracted_info_condition`・`line_condition`・`tone_default`（汎用口調・PC 口調・PC 性別の単一行設定）を足す。
- 取込（Go）: `LinkLineConditionsFromStaging` を `engine.Ingest` へ追加（line_speaker と対称の解決）。
- 純粋ルール（tone）: `EmotionBandOfLine`（1 台詞用の渋い激情しきい値 1.5）を追加。`PathGeneric`・`PathPC` を追加。
- 純粋ルール（personatone）: `BuildFreeToneTraits`（自由記述＋感情助言＋性別の一人称・語尾）を追加。
- 注入（engine）: `LinePersonas` を全台詞へ広げ、(rec, field) で ①名指し／②汎用／③PC を振り分け。PC 発話は NPC 話者が結ばれていても PC 既定を優先する（実画面確認で発見した routing バグを修正）。
- 設定（store・api）: `GetPromptTemplate`/`SavePromptTemplate` を prompt_template と tone_default の結合・往復保存へ拡張。`PromptTemplateView`・`PersonaView`（sex 追加）を拡張。
- frontend ロジック: gateway（translation・template）・Container（口調設定の state・save・dirty）・wailsjs models を配線。

### 設計からの逸脱（実装上の修正）

- 口調設定の保存先を prompt_template への列追加でなく専用テーブル `tone_default` にした。C# 抽出器が全 migration SQL を毎回 ensure するため、`ALTER TABLE ADD COLUMN` は再実行で失敗する。`CREATE TABLE IF NOT EXISTS` で冪等にする必要があった。設計意図（口調設定の永続）は保つ。

### 最終検証（通過）

- backend: `npm run verify:backend`（go test 全通過・arch-lint OK・boundary OK）、`npm run lint:backend`（format 0・vet・static・arch・boundary・module 全 OK）。
- frontend: `npm run test:frontend`・`npm run lint:frontend`（svelte-check アプリ 0 エラー・tsc・eslint・knip・boundaries 全通過）、`build-storybook` 成功。
- 純粋ルール: tone カバレッジ 100%（`EmotionBandOfLine` 含む）、personatone 98.2%（新規関数 100%）。
- 非劣化: harness golden は 汎用台詞が 1 行口調を得る差分のみ（名指し・persona_character・master_term 不変）。

### 完了定義の観測（実 app・実データ）

- 実データ: 実 plugin `Innocence Lost - Quest Expansion.esp` の抽出で `extracted_info_condition` 7 件（衛兵の汎用台詞、全て Male）、`line_condition` 6 件が解決。form_id が INFO:NAM1 と一致し結合が成立。
- 実画面（`http://localhost:34115`、translation 段はローカル stub で代替）:
  - 衛兵の汎用台詞に「口調: 汎用台詞」が付き、展開で「感情 X ・ 性別 男性」と、自由記述口調＋感情助言を合成した実プロンプトを確認。
  - PC 発話（DIAL:FULL・INFO:RNAM）に「口調: PC発話」が付く。NPC 話者が結ばれた INFO:RNAM も PC 既定へ振り直すことを確認。
  - 名指し話者の口調（平明・ぞんざい・慇懃・端正・率直・興奮）が従来どおりで非劣化。
  - 「口調なし」は叙述文・定型句だけに残り、台詞には残らない。
  - プロンプトテンプレート画面で汎用口調と PC 性別を編集・保存し、再表示で反映を確認。編集した自由記述が実プロンプトへ流れることを確認。

### 観測できなかった点（環境境界・人間へ引き継ぎ）

- 実 LLM（LM Studio `http://192.168.0.226:1234`）はこの環境の LAN 外で到達不可のため、実訳文への口調反映の目視は未実施。注入される実プロンプト（口調指示）は stub 実行で確認済み。最終の実訳文確認は LM Studio が稼働する環境で人間が行う。
- OutfitRecognition の装備反応台詞は単独では未確認。性別条件を持たない「真の汎用」は Innocence Lost の衛兵台詞（性別なしの行）で等価に確認済み。

## finalization（2026-07-02）

- 正本化判断: `docs/architecture.md` 反映不要。新規 core 関数・store テーブル/メソッド・engine routing・抽出器 writer は既存層内の追加で、層構成・依存方向・Bootstrap・Wails 境界の構造を変えない。新 Wails メソッドなし。arch-lint・boundary-scan 通過。人間承認不要（構造不変）。
- 正本反映: `docs/architecture.md` は変更なし。
- er.md 更新（人間承認済み）: 論理 ER 正本へ `extracted_info_condition`・`line_condition`・`tone_default` の 3 テーブルと `line ||--o| line_condition` の関連を追記。
- spike probe 撤去: `tools/extractor/GenericDialogueProbe.cs` を削除、`Program.cs` の `--probe-generic` を除去。本番 `InfoConditionSqliteWriter` が性別抽出を担う。抽出器の再ビルド通過。

## spike 成果物

- 調査コード: `tools/extractor/GenericDialogueProbe.cs` と `--probe-generic` flag（read-only。中心 DB へ書かない）。
- 出力: `tmp/generic-voice-tone-fallback/probe-skyrim.md`、`tmp/generic-voice-tone-fallback/probe-outfit.md`。
- probe コードの存置/撤去は finalization で判断する。
