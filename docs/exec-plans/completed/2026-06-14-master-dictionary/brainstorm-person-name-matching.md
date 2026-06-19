# ブレスト資料: 固有名の機械置換における人名の形ゆれ問題

本書は、別セッションでのブレインストーミング用に、現在詰まっている論点を自己完結でまとめた資料である。この会話の文脈が無くても読めるようにしてある。

## 1 行で言うと

固有名を辞書で機械置換して訳を一貫させる仕組みを作ったが、人名だけは本文での出方が「フルネーム＋二つ名 / 名のみ / 短名」と揺れるため、辞書（レコードのフルネーム）とかみ合わず、名のみで出た人名が AI 任せの訳になって一貫性が崩れる。

## 背景: 何を作っているか

- task: `2026-06-14-master-dictionary`（Mod 横断マスター辞書）。
- goal: 同一原語へ常に同一訳語を当て、Mod 横断で一貫性を保つ。
- 趣旨: 固有名（地名・武器・人名など）の訳をあらかじめ辞書化し、叙述文・台詞を翻訳するとき、本文中に出る固有名を辞書の確定訳語へ**機械置換**してから AI 翻訳する。AI は周りの英語だけを訳し、差し込んだ日本語固有名はそのまま残る。
- 供給源: 既訳流用主軸。公式日本語版の既訳を、xTranslator 英日辞書 XML（`dictionaries/xTranslatorXMLs/*.xml`）から取り込む。

## 今の実装（要点）

- 辞書テーブル `master_term(source 原語, dest 確定訳語, category 種別)`、一意キー `(category, source)`。`db/migrations/0003_master_term.sql`。
- 構築: C# extractor が xTranslator XML を解析し、`REC` が `:FULL` で終わる行（龍語綴り `WOOP:FULL` は除く）の `Source`→`Dest` を登録する。`tools/extractor/MasterTermXmlWriter.cs`。実 6 XML から 24,554 件。
- 置換器: `internal/engine/dictionary.go` の `Dictionary`。
  - 貪欲最長一致（長い原語を先に当てる。`Iron Sword` を `Iron` より先）。
  - 語境界 `\b` で囲む（`Sword` が `Swordsman` の内側に当たらない）。
  - 大小を区別（原語は大文字始まりの固有名。小文字の一般語に当たりにくい）。
  - 同一原語に複数訳が来たら最初の 1 つを保つ。
- 適用: `engine.Run` が叙述文・台詞の原文へ `Dictionary.Apply` を当て、置換後の原文を AI へ渡す。`internal/engine/engine.go`。

## 効いている例（地名・施設・組織は問題ない）

実 mod（`Innocence Lost - Quest Expansion.esp`）を LM Studio（gemma-4-12b）で翻訳した実測:

- `It's Honorhall Orphanage, in Riften.` → `リフテンにある、オナーホール孤児院だよ。`
  - `Riften`→`リフテン`、`Honorhall Orphanage`→`オナーホール孤児院`（辞書の確定訳語）。
- `Riff-raff! That's all you Riften people are!` → `掃き溜れの屑どもめ！リフテンの連中なんて…`
  - `Riften`→`リフテン`。`Riff-raff` は語境界で別語として残す（誤置換しない）。
- 一貫性: `Honorhall Orphanage`・`Riften` は複数行で同一訳に揃う。

地名・武器・施設・組織は、本文がレコードのフルネームと同じ文字列で出るため、完全一致でかみ合う。

## 困りごと: 人名の形ゆれ

人名は本文での出方が揺れるため、辞書とかみ合わない。

実例（同じ実行から）:

- 辞書にあるのは `Grelod the Kind` → `親切者のグレロッド`（`NPC_:FULL`、フルネーム＋二つ名）。
- だが台詞は名のみ `Grelod` で出る: `Grelod is the headmistress of Honorhall Orphanage.`
- `Grelod` 単独は辞書に無い（`Grelod the Kind` とは別文字列）ので置換されず、AI が自前で `グレロド` と訳した（official は `グレロッド`、促音 `ッ` が落ちた）。
- 結果: 同じ人物が、フルネームで出れば `…グレロッド`、名のみで出れば `グレロド` と揺れる。

つまり、辞書に人名は入っている（`NPC_:FULL` 1699 件）が、本文が名・短名で出ると、その「形」が辞書に無く落ちる。

## なぜ人名だけ難しいか

- 地名・物・組織: 本文はレコードのフルネームと同じ文字列で参照する。完全一致で足りる。
- 人名: 本文はフルネーム（`Grelod the Kind`）/ 名のみ（`Grelod`）/ 短名（SHRT）/ 所有格（`Grelod's`）/ 称号つき（`old Grelod`）など、レコードのフルネームと違う形で出る。完全一致が外れる。
- 辞書側も、`NPC_:FULL` は「名＋二つ名」を 1 文字列で持つ（`Grelod the Kind`→`親切者のグレロッド`）。名の部分だけ（`Grelod`→`グレロッド`）を取り出した entry は無い。

## 検討した選択肢と論点

### 案A: 短名 `NPC_:SHRT` も辞書へ取り込む

- 内容: 現在の `:FULL` だけでなく `NPC_:SHRT`（短名、Skyrim XML に 260 件、例 `Elda`→`エルダ`）も登録する。
- 利点: 短名で出る人名の穴が減る。実装は filter に `NPC_:SHRT` を足すだけ。
- 欠点: SHRT には一般語に近い語が混じる（`Conjurer`→`コンジュラー`、`Addict`→`中毒者`）。一般語の誤置換が増える。なお `NPC_:FULL` 側にも既にこの傾向はある（`Addict` 等が NPC 名として入っている）。
- 限界: SHRT は「名のみ」とは限らない。`Grelod` のように FULL から名だけを取った形は SHRT に無いこともある。

### 案B: 人名だけ部分一致を許す（利用者が floated した案）

- 内容: NPC 名に限り、完全一致でなく部分一致（本文の語がフルネームの一部に当たれば置換）を許す。
- 最大の難点: 「何へ置換するか」。英語側で `Grelod the Kind` から名 `Grelod` を切り出すのは比較的容易（` the ` の前など）。だが日本語側 `親切者のグレロッド` から名の部分 `グレロッド` を機械的に取り出すのが難しい（二つ名 `親切者の` が接頭、語順も英日で違う）。置換先の日本語が確定できない。
- 過剰一致の懸念: 部分一致は誤爆しやすい。短い名や一般語と同綴りの名で、無関係な箇所を置換する恐れ。

### 案C: 人名は置換でなく参考語注入（speaker 経由）

- 内容: 台詞は話者（speaker）を解決済み（ペルソナ生成で使用）。その話者の名（FULL＋短名）の訳を、その台詞の AI プロンプトへ参考語として渡す。AI が語形変化を吸収しつつ official 訳を使う。
- 利点: 名の形ゆれ・助詞・敬称を AI が自然に処理できる。
- 欠点: 注入は助言で、AI が必ず従う保証は無い（本 task は一貫性のため置換主軸にした経緯がある）。話者名を中心 DB へ持つ配線が要る（現状 speaker テーブルは名前列を持たない。下記参照）。

### 案D: 据え置き（既知の弱点として許容）

- 内容: 完全一致で当たる形（フルネーム・地名・物）は辞書で揃え、名のみ等の部分形は AI 任せとして許容する。概念モデルの弱点（同綴り異義の誤統合許容）と同種のトレードオフと整理する。

## ブレストで決めたい問い

- 部分一致（案B）を採るなら、英語の部分名に対応する**日本語の置換先**をどう決めるか。`親切者のグレロッド` から `グレロッド` を取り出す機械的な手立てはあるか（名・姓・二つ名の構造データ、または分割規則）。
- 名の構造（first / last / byname）を与える供給源はあるか。xTranslator XML や Skyrim レコードに、名の分割や短名の対応はあるか。
- 案C（注入）と置換主軸は両立できるか。人名だけ注入、その他は置換、の併用は妥当か。一貫性は落ちないか。
- 過剰一致をどう抑えるか（大文字始まり・語境界・最小長・話者文脈での限定など）。
- そもそも人名の一貫性は、地名ほど厳密に要るか。許容できる揺れの範囲はどこか。

## 参考: データと現状の事実

- xTranslator XML 1 件の形（`dictionaries/xTranslatorXMLs/Skyrim_english_japanese.xml` ほか）:

  ```xml
  <String List="0" sID="..." Partial="1">
    <EDID>Grelod</EDID>
    <REC>NPC_:FULL</REC>
    <Source>Grelod the Kind</Source>
    <Dest>親切者のグレロッド</Dest>
  </String>
  ```

- `master_term` の `category` 別件数（上位）: `DIAL 6230 / ARMO 3213 / WEAP 3150 / NPC_ 1699 / QUST 1266 / BOOK 997 / MGEF 938 / SPEL 731 / LCTN 709 / CELL 673 …`。
- `NPC_:SHRT` は Skyrim XML に 260 件あるが、現状 `master_term` は `:FULL` のみ取り込み、SHRT は未取込。
- 中心 DB の `speaker` テーブルは名前列を持たない（列は speaker_kind / sex / occupation / person / tone / background / race_id / voice_type_id / template_speaker_id / plugin / form_id / edid）。話者の名は中心 DB に無い（T2 で e8 話者→名は対象外にした）。
- 置換は本文翻訳の前に原文（英語）へ当て、置換後の原文を AI に渡す。辞書の照合キーは原語文字列（種別は区別に保持するが照合には使っていない）。

## 関連ファイル

- 辞書テーブル: `db/migrations/0003_master_term.sql`
- 辞書構築（C#）: `tools/extractor/MasterTermXmlWriter.cs`、`tools/extractor/Program.cs`（`--sqlite` 経路）
- 置換器（Go）: `internal/engine/dictionary.go`（＋ `dictionary_test.go`）
- 置換の適用: `internal/engine/engine.go`（`Run` の差込点）
- 表示（結果行への「置換した固有名」併記、storybook 承認済み・実データ配線は未）: `frontend/src/ui/screens/translation-run/TranslationResultRow.svelte`
- 概念モデルの固有名・名前: `docs/concept-model.md`（e8 話者→固有名、弱点1 同綴り異義）
