# 実装範囲・テスト設計: record-type-translation-expansion

- task-id: 2026-06-23-record-type-translation-expansion
- 由来: design-module 出口（人間設計レビュー承認 2026-06-25）。確定方針は `plan.md` と `docs/changelog.md`（2026-06-25 entry）にある。
- 位置づけ: implementation-module へ渡す scope の境界・依存・検証単位を固定する。画面の見せ方は storybook-module で設計し、本書はその受け口を示す。

## scope の境界

含む（今回の task で動かす）:

- 会話（`INFO:NAM1`）・本（`BOOK:DESC`）以外で、概念モデルに登場する翻訳対象の全 REC:FIELD を翻訳対象へ加える。武器は一例で、対象を武器へ絞らない。
- 4 つの箱（叙述文・固有名・定型句・台詞）の全 REC:FIELD を今回まとめて対象にする。翻訳しない種別（無訳片 `WOOP:FULL` など）は読み込まずマスターに載せない。
- 固有名を本文翻訳より先に確定する固有名フェーズを新設する。既訳は `master_term`（権威訳）、AI 訳は `proper_noun`（実行内）へ分け、本文へ機械置換注入する。
- プロンプトを「Base 指示 ＋ その REC:FIELD に割り当てた指示文（directive、変数を実行時に埋めたもの）」で組み立てる。口調・文体・固有名・定型句をすべて directive に揃える。directive の指示文を画面で編集できるようにする（割り当ては固定）。

含まない（除外・別 task）:

- 固有名辞書の実行をまたぐ永続化。中心 DB は起動ごとに空にする意向があり永続化方針が未決のため、本 task は同一実行内のフェーズ順序に限る（方針 A）。
- 叙述文・台詞の文中固有名を自動検出して固有名へ結ぶ「言及」。固有名注入は辞書の機械置換に限る。
- 既訳 XML の取り込み機構。
- 定型句の重複排除（`set_phrase` 新設）。不一致・コストが実測で問題化した時の後続 task（起動条件付き）。

## 確定した構造

### 新設・変更テーブル（永続層）

- `record_type_master`（新設）: REC:FIELD の分類と directive 割り当ての正本。`rec` PK・`field` PK は Skyrim 仕様で固定行、`box`（叙述文/固有名/定型句/台詞）固定、`directive`（割り当てた指示文キー）固定。翻訳対象だけを載せる（翻訳対象フラグも無訳片の行も持たない。翻訳しない種別は読み込まない）。
- `directive`（新設・MECE モデル）: 指示文の正本。`key` PK（説明体・書物体・日記体・世界観断片・口調・固有名・定型句）→ `instruction`（指示文・編集可能）。`{traits}` のような変数を宣言できる（口調が `{traits}` を持つ）。複数 REC:FIELD が 1 directive を共有する。`prompt_template` の口調（persona）はここの口調行へ畳む。
- `prompt_template`（既存・base）: Base 指示。directive とは別に全種別共通の前提として残す。
- `extracted_field`（新設）: C# 抽出器が素朴吸い出しした原文の受け皿（plugin・form_id・edid・rec・field・ordinal・source）。概念モデルの箱ではなく C#↔Go 受け渡しバッファ。
- `proper_noun`（新設・方針 A）: 固有名の訳の単位（id・source・category・dest・status）。AI 訳はここに留め `master_term` へ昇格しない。
- `narration`（変更）: 叙述文・定型句の全 REC:FIELD の行が増える。`style` カラムは directive キーの保持に転用するか、`record_type_master.directive` を取込段で引く。スキーマ変更は実装段で確定。
- `line`（変更）: `INFO:RNAM`・`DIAL:FULL` の行が増える。スキーマ変更なし。口調は directive の口調行から取る。
- `master_term`（未変更）: 権威訳の横断永続辞書。AI 訳を入れないため構造も使われ方も変えない。
- migration は次の連番（0005 以降）で新設テーブルと seed（`directive`・`record_type_master` 初期データ）を足す。

### 翻訳手続き（engine の段階順序）

- 取込段（新設）: `extracted_field` を読み、`record_type_master` を rec/field で引いて box と directive を得て `narration`・`proper_noun`・`line` へ振り分ける。
- 固有名フェーズ（新設・本文フェーズより前）: 各 `proper_noun` を `master_term` と突合する。既訳あり → 権威訳。既訳なし → AI 訳して `proper_noun` へ書く（`master_term` へは書かない）。
- 本文フェーズ（変更）: 権威訳（`master_term`）＋AI 訳（`proper_noun`）の辞書を叙述文・台詞へ機械置換注入する。プロンプトは Base 指示に、その REC:FIELD の directive（`directive.instruction`、変数を実行時に埋めたもの。台詞は口調 directive の `{traits}` へ話者の性質を埋める）を合成して作る。

### 純粋判定（`internal/engine/` 配下・横の層を作らない）

- 供給源選別: 入力＝固有名の既訳の有無、出力＝書き込み宛先（権威訳を使う / AI 訳→`proper_noun`）。不変ルール＝AI 訳の宛先に `master_term` を返さない。
- プロンプト合成: 入力＝Base 指示・directive の指示文・変数の値（台詞は話者の性質）、出力＝合成済みプロンプト。directive の取得（`directive`・`record_type_master` 読み出し）は外側の I/O に置き、合成と変数差し込みは引数で受け取る純粋関数にする。変数なしの directive（固有名・定型句など）は Base ＋ 指示文のみになる経路を含む。

## 依存（実装の順序）

各段は前段の出力に依存する。横の層を完了線にせず、各段は次段が観測できる出力まで作る。

1. migration（新設テーブル＋seed）。`directive` に 7 指示文、`record_type_master` 初期データで全 translatable REC:FIELD を directive へ割り当てる。
2. C# 抽出器を `extracted_field` への素朴吸い出しへ変える（箱・directive の判定を持たない）。
3. Go 取込段（マスター参照で box 振り分け・directive 引き当て）。
4. 固有名フェーズ＋供給源選別（純粋）。
5. 本文フェーズの directive 合成＋プロンプト合成（純粋）。
6. プロンプトテンプレート画面のレコード別タブ（directive 編集）と結果表示の更新（種別バッジ）の配線。画面は storybook-module で承認済み（`Screens/プロンプトテンプレート`・`UI Components/DirectiveEditor`・`UI Components/TranslationResultRow`）、状態・gateway・タブ・directive 保存の配線は implementation-module。`TemplateEditorContainer` を新 props（activeTab・directives・assignments・onInstructionInput・onTabChange）へ配線し直す。

## テスト設計

単体テスト（純粋判定はカバレッジ 100% 基準）:

- 供給源選別: 既訳ありは AI 翻訳しない、既訳なしだけ AI 訳対象、AI 訳の宛先は必ず `proper_noun` で `master_term` を返さない（不変ルール）。境界条件を網羅。
- プロンプト合成: Base 指示と directive の指示文の合成、変数の差し込み（口調の `{traits}` に話者の性質）。変数なしの directive で Base ＋ 指示文のみになる経路を含む。
- `record_type_master` 初期データ: 全 translatable REC:FIELD が directive を 1 つだけ持ち（排他）、全 REC:FIELD が directive へ割り当たる（網羅）ことを初期データの整合テストで確かめる。
- 固有名・定型句の本文注入は既存 `Dictionary` 単体テストの範囲を流用する。

実画面・実データの観測点（DB アクセス・AI 呼び出し込みは E2E に任せる）:

- 取込段の振り分け: 抽出→翻訳実行後、会話・本以外の各 REC:FIELD の行が中心 DB に入り訳文が書き戻る（実データ）。
- 固有名フェーズ: 固有名が本文翻訳より先に確定し、既訳の無いものが AI 訳されて `proper_noun` へ入り、`master_term` には権威訳だけが残る（実データ）。
- directive 編集とプロンプト合成: 実 app（`npm run dev:wails:run`、`http://localhost:34115`）で、プロンプトテンプレートのレコード別タブで directive の指示文を編集し、複数の代表種別（武器・呪文・任務・地名など）で指示文に応じた訳と固有名の一貫注入を目視（実画面）。
- 退行防止: 既存の本（`BOOK:DESC`）・会話（`INFO:NAM1`）翻訳が壊れない（実画面）。

検証コマンド:

- backend 単体: `go test ./internal/engine/...`（純粋判定・初期データ整合）、ビルド `go build ./...`。
- 実画面: `npm run dev:wails:run` で抽出→翻訳を実行し観測点を目視。

## architecture.md 反映（design-module 出口で確定）

- 反映は要（構造が変わる）。取込段・固有名フェーズの追加（段階順序の制約）、`record_type_master`・`directive`・`extracted_field`・`proper_noun` の新設、プロンプト合成（Base ＋ directive）への変更、プロンプトテンプレートのレコード別タブの Wails 境界を §8 へ反映する。横レイヤー追加は含めない。
- 実際の §8 追記は、構造が実在する implementation-module / finalization-module で行う（実在しない構造を先に書いて churn させない）。
