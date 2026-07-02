# 実装範囲とテスト設計（design-module 確定）

人間設計レビュー承認済み。確定した設計判断は design.md と本書に従う。詳細仕様は列挙せず、scope の境界・依存・検証単位を固定する。

## 確定した設計判断（5 件、design-module 人間設計レビューで #2・#4 を自由記述方式へ変更）

- #1 性別の経路: staging `extracted_info_condition`(info_plugin, info_form_id, sex) → domain `line_condition`(line_id, sex)。`extracted_info_speaker`→`line_speaker` と対称。
- #2 汎用・PC の口調指定: 段階選択でなく自由記述。`generic_tone_text`（汎用口調の自由記述）・`pc_tone_text`（PC 口調の自由記述）・`pc_sex`（PC の性別、Female/Male/空）を持つ。汎用・PC は対人段階・セル名を持たない。
  - 実装注記: 保存先は prompt_template への列追加でなく専用テーブル `tone_default`（単一行）にした。C# 抽出器が全 migration SQL を毎回 ensure するため ALTER は冪等にできない。取得・保存は prompt_template と tone_default の結合で 1 つの DTO に束ねる（公開境界は不変）。
- #3 激情しきい値: `tone.Classifier` へ 1 行用の渋い激情しきい値を新設。具体値は本文較正で確定。
- #4 公開境界: 口調メタ DTO へ汎用・PC 用の決定経路・感情・性別を載せる（対人段階・セル・印は汎用・PC では持たない）。新 Wails メソッドなし。prompt_template の取得・保存 payload を拡張。
- #5 PC 判定: PC 発話 = (DIAL, FULL) と (INFO, RNAM)。NPC 返答 = (INFO, NAM1)。

## 実装範囲（component と依存）

### C# 抽出器（tools/extractor）

- INFO の条件から性別を判定し `extracted_info_condition` へ書く writer を足す。判定は `GetIsSex`（極性畳み）・単一声型 EDID の `Male`/`Female` 接頭・同性のみ FLST 展開。GenericDialogueProbe の Classify 相当を本番用に整える。
- schema ensure に新テーブルを足す。
- 依存: 既存の `SpeakerSqliteWriter`・`PluginExtractor` の条件解析。
- 検証単位: `--sqlite` 取込後の DB に `extracted_info_condition` 行が入る（実データ確認）。

### スキーマ（db/migrations）

- 新 migration: `extracted_info_condition`（staging）・`line_condition`（domain）・`prompt_template` へ 3 列追加（`generic_tone_text`・`pc_tone_text`・`pc_sex`、既定の文面を seed）。
- `docs/er.md` へ domain の `line_condition` を反映（docs 正本）。
- 依存: 0002（line・prompt_template）・0006（extracted_info_speaker）。
- 検証単位: migration 適用と seed 整合。

### Go 取込（internal/store）

- `LinkLineConditionsFromStaging`: `line`(rec='INFO') と `extracted_info_condition` を (plugin, form_id) で結び `line_condition` へ。`LinkLineSpeakersFromStaging` と対称。
- 依存: line 投入後に 1 度呼ぶ。
- 検証単位: 取込テスト（既存 ingest テストの形に合わせる）。

### 純粋ルール（internal/core/tone）

- 1 台詞の特徴量から感情段階だけを出す純粋関数を足す（対人段階は持たない）。激情しきい値は 1 行用に渋め（新フィールド）。
- 感情段階の助言文（抑制／中／激情の 3 段）と、性別の一人称・語尾（自由記述へ重ねる用）の組み立ては純粋ルールへ閉じる。性別の一人称・語尾は `rolespeech` を性別主体で引く（汎用・PC はセルを持たないため、セルをワイルドカードで照合する）。
- 依存: 既存の `scoreAxes`・`arousalBand`・`rolespeech` を再利用。
- 検証単位: ユニットテスト 100%（境界条件）。

### 注入（internal/engine）

- `LinePersonas` を全未訳台詞へ広げる。話者解決なしの台詞へ、自由記述の口調指示文（汎用／PC）に本文 1 行の感情段階と性別の一人称・語尾を重ねた口調指示を組む。(rec, field) で経路①②③へ振り分ける。
- `line_analysis` を汎用・PC 台詞の本文へも広げる（現状は話者台詞のみ）。
- 汎用・PC の自由記述指示文と PC 性別を `prompt_template` から読む。汎用の性別は `line_condition` から読む。
- 依存: tone の新関数、store の `line_condition` 読み出し、`prompt_template` の新列。
- 検証単位: engine テスト（routing の 3 経路、感情・性別の重ね）。

### 設定の保存・取得（internal/store・internal/api）

- `model.PromptTemplate` と `GetPromptTemplate`/`SavePromptTemplate` へ汎用口調・PC 口調・PC 性別を足す。app の prompt-template Bind payload を拡張する（新メソッドなし）。
- 検証単位: 実画面（編集→保存→翻訳反映）。

### 表示（storybook-module 担当）

- 口調メタ DTO へ汎用・PC の決定経路（「汎用」「PC」）・感情段階・性別を載せる。汎用・PC は対人段階・セル名・印を持たないため、結果行はそれらを出さず、決定経路・感情段階・性別だけを出す。汎用は条件由来の性別、PC は利用者選択の性別。
- プロンプトテンプレート画面へ「汎用台詞の口調」「PC 発話の口調」の自由記述欄 2 つと、「PC の性別」の選択（男性／女性／未指定）を足す。
- 画面の正本は Storybook の story と svelte コンポーネント。本書では doc 化しない。

## scope 外（design.md §8 再掲）

- 役どころ（衛兵・兵士）の register tag。
- 名指し話者の口調を 1 行感情へ作り替えること。
- 勢力・種族・声型の気質 prior。
- 皮肉の検出。

## テスト設計

### 単体テストで書く

- `internal/core/tone` の per-line 感情段階関数。1 行の感情段階化、激情しきい値の境界（単発の感嘆符で激情にしない）、感情助言文の引き、性別の一人称・語尾の引き（セルをワイルドカード照合）。カバレッジ 100%。
- 取込 `LinkLineConditionsFromStaging` の解決ロジック（既存 store テストの慣行に合わせる）。

### 単体テストで書かない（実画面・E2E に任せる）

- Wails 境界（prompt_template payload 拡張）。
- 画面ロジック（自由記述欄・PC 性別選択、根拠表示）。
- DB・LLM 込みの注入経路全体。

### 実画面確認（plan.md の完了定義の観測点）

- 衛兵の汎用台詞（`Innocence Lost - Quest Expansion.esp`）に口調が付く。
- OutfitRecognition の装備反応台詞に自由記述の汎用口調＋本文感情で口調が付く。
- PC 発話（DIAL の選択肢）に自由記述の PC 口調が付き、選択した PC 性別の一人称・語尾が反映される。
- 名指し話者の口調が従来どおりで非劣化。
- 汎用口調・PC 口調・PC 性別を画面で変更し、翻訳プロンプトへ反映される。
