# dictionary-false-positive-guard（機械置換辞書の誤爆対策: 一般語 stoplist と内部フラグ勢力の除外）

## 分岐情報

- 作業 branch: `claude/dictionary-false-positive-guard`
- 分岐元 branch: `master`
- 分岐元 commit: `1aeb2f18`

## 背景（観測済みの誤爆）

機械置換辞書（`master_term` ∪ `proper_noun`）は、レコードの FULL 名を語境界・大小区別の文字列一致で本文へ当てる。この供給源に短い一般語と同綴りの固有名が入ると、本文の通常の語に誤って当たる。

実測（inigo.esp、完了 plan `narration-line-mention-linking` の残留リスクとして記録済み）:

- inigo.esp は FULL 名が "Yes"・"No" の勢力（FACT）レコードを持つ。Mod 制作では勢力を対話状態の内部フラグに使う慣行があり、この 2 つも画面に出ない管理用勢力である可能性がある。
- 取込段は FACT:FULL を固有名 box へ振り分けるため、"Yes"・"No" が `proper_noun` に登録され、固有名フェーズで AI 訳が付き、機械置換辞書に載る。
- 本文フェーズで文頭の大文字 "No"/"Yes" に当たり、`No matter what the weather...` → `訳030 matter what the weather...` のような壊れた入力が AI へ渡る。分岐元挙動の実データ golden で置換入りプロンプトは 386 行あった。
- 言及テーブル（`narration_mention`・`line_mention`）は注入の実態を忠実に写すため、同じ誤爆が偽陽性の言及として記録される。

既存の緩和（大小区別・語境界）で小文字 `no` や `nothing` の内側には当たらない。誤爆は文頭・文中の大文字 Yes/No 等に限られる。

## 依頼要約

誤爆を辞書の供給源側で止める 2 層の対策を入れる。

1. **一般語 stoplist**: 本文の通常語と同綴りになりやすい一般語（"Yes"・"No" など）を、機械置換辞書と言及検出語彙の両方の供給から除く。stop-word の語集合は自作リストでなく、外部で配布されている確立した stopword リストを利用する（利用者指示）。配布元は次節のとおり stopwords-iso（MIT）に確定済み。
2. **内部フラグ勢力の除外**: 画面に出ない管理用勢力の FULL 名を、翻訳対象（固有名の供給）から除く。機械的な判定基準の候補は FACT レコードの Hidden from PC flag（Mutagen で読める想定）で、実データでの成立確認を経てから固定する。

不変条件: 注入（`LoadDictionary`）と言及検出（言及語彙）は同じ供給源選別を通す。片側だけに効く除外を作らない（言及レコード＝注入の忠実な記録、という前 task の不変条件を保つ）。

## stopword リストの選定（確定、2026-07-04）

利用者指示「外部配布・commit できるライセンス」に基づき、plan 作成時に選定を確定した。

- 採用: **stopwords-iso の英語リスト（stopwords-en）**
  - 配布元: `https://github.com/stopwords-iso/stopwords-en`（raw: `https://raw.githubusercontent.com/stopwords-iso/stopwords-en/master/stopwords-en.txt`）
  - ライセンス: MIT（Copyright (c) 2016 Gene Diaz）。repo へ commit できる。
  - 内容: 1 行 1 語・小文字・1297 語。誤爆語 "yes"・"no" をどちらも含むことを実物で確認した。
  - 取得時 checksum（sha256）: `883fc88a14aa980677c80d485e6c16863fa40ff87743b1f46c6970845ed8f0a5`
- 選外: Snowball 英語 stop list（BSD だが実物確認で "yes" を含まない）。SMART stop list（"yes" を含むが、一次配布元の明示ライセンスを確認できない）。
- 衝突調査（Skyrim の固有名になりうる一般語がリストに入っているか、実物 grep）:
  - 含まれない（除外されない）: gold・key・chest・guard・ice・storm・fear・calm・courage
  - 含まれる（除外される）: open・home・well・will・fire・alone
  - 評価: 除外されるのは「名前全体が一般語 1 語」の固有名だけで、複数語の名前（Fire Damage、Iron Sword 等）は影響を受けない。名前全体が一般語 1 語の固有名は Yes/No と同じ誤爆クラスであり、盲目的な本文置換から外すことはむしろ本 task の目的に合う。実データでの影響件数は実装時の観測点で確認する。

## 完了定義

### 動かす範囲

1. stoplist: 実データ（inigo.esp）の翻訳実行で、stoplist 語（"Yes"・"No" を含む）が本文へ機械置換されず、言及レコードにも記録されない。stoplist に無い固有名（例: Riften、Dragonbane）の置換・言及は従来どおり残る。
2. 内部フラグ勢力の除外: 判定基準が実データで成立する場合、inigo.esp の "Yes"・"No" 勢力の FULL が抽出結果（翻訳対象）から除かれ、`proper_noun` に現れない。基準の成立確認は次の 2 点で行う。
   - inigo.esp の "Yes"・"No" 勢力が基準に該当すること。
   - base ゲーム（Skyrim.esm 等）の実在勢力名で、本文に出る名前（`master_term` に既訳がある名前）が基準により大量に落ちないこと（落ちても注入は `master_term` 側が担い続けることの確認を含む）。
   - 基準が成立しない場合: 観測結果と不成立理由を plan に固定し、2 は不成立で停止する（同じ誤爆は 1 の stoplist が止める）。

### 観測点

- 単体テスト: 供給源選別（stoplist を含む語彙の選別ルール）の純粋 IO 部分。守るべき不変ルールとしてユニットテスト 100% カバレッジを基準にする（`feedback-pure-io-rule-100-coverage`）。注入語彙と言及語彙が同じ選別を通ることをテストで固定する。
- 合成 golden（`internal/harness`）: 合成 fixture へ stoplist 対象語（例: FULL="Yes" のレコード）を追加し、置換も言及も起きないことを golden で凍結する。golden 差分は fixture 追加と stoplist 効果に限ることをレビューで確認する。
- 実データ: inigo.esp で翻訳実行し、(a) プロンプトに Yes/No 由来の置換が無いこと、(b) `narration_mention`・`line_mention` に Yes/No への言及が無いこと、(c) 2 が成立した場合は `proper_noun` に "Yes"・"No"（FACT 由来）が無いことを DB 内容と golden 比較で確認する。

### 意図的な出力変更（非劣化の扱い）

本 task は既存出力を意図的に変える（stoplist 語の置換をやめる）。よって「golden 完全一致」は成立条件にしない。代わりに次で劣化と変更を分離する。

- 分岐元 commit で実データ golden（inigo.esp）を捕獲し、本変更後との差分が「stoplist 語・除外勢力に関する置換と訳の消失」だけであることを確認する。無関係な行（stoplist 外の固有名の置換、話者連関、口調）に差分が無いこと。
- `npm run verify:backend` が通過すること。

### 含まない（除外範囲）

- 注入方式の変更（本文置換をやめ用語対訳ヒントを添付する方式）。誤爆対策の第 3 案だが、プロンプト設計の再検討を伴う別テーマのため対象外とする。
- 注入語の保持検証（`known-issues.md` 2 番）。言及テーブルを使う事後検証は別 task とする。
- 辞書に無い漏れ語の拾い上げ（`known-issues.md` 1 番の第 2 層）。
- stoplist の画面からの編集機能。stoplist は本 task では外部配布リスト由来の固定データとして持ち、編集 UI は作らない。
- stopword リストの自作・独自拡張。外部配布リストに無い語を独自に足す運用は本 task ではしない（リスト選定で賄う）。

## 軽 / 重判定

| 軸 | 判定 | 根拠 |
| --- | --- | --- |
| 画面が動くか | N | 変更は辞書の供給源選別と抽出規則。`Terms`（機械置換内訳）の表示件数は変わるが、表示ロジック・layout・文言・style・svelte コンポーネント・story・fixture は変えない。 |
| `docs/architecture.md` 反映が要るか | N | 供給源選別の純粋ルール追加は `internal/core/<name>` の既存パターン内。C# 抽出器の抽出規則変更（Hidden from PC の除外）は `tools/extractor` 内で完結し、schema（C#↔Go 契約の migration）・層構成・依存方向・Wails 境界を変えない。 |

判定結果: 両方 N のため軽 task。`design-module` と `storybook-module` を bypass し、`preparation-module` → `implementation-module` → `finalization-module` で進める。

実装時に Claude 本体が固定する設計判断（軽 task の範囲）:

- stopwords-en.txt（MIT、選定済み）の repo 内の置き場所（`assets/` 配下を第一候補）と、出典・ライセンス・checksum の記録方法（LICENSE 併置または冒頭コメント）。commit するため合成 harness と CI は決定的に読める。
- stopword の照合規則（原語全体を小文字化して一語一致、複数語の固有名には当てない等）の純粋ルール化と読み込みアダプタの形。
- 供給源選別を注入と言及の両方に通す共通点の切り出し方（`LoadDictionary` と言及語彙構築の共通化）。
- 内部フラグ勢力の判定基準の実装位置（C# 抽出器で書かない、の一択か、書いた上で Go 取込段が除くか。schema 変更を伴わない前者を優先候補とする）。
- 基準の成立確認の手順（Mutagen で FACT flags を読む調査 code の置き場所は一時 runner とし、観測結果だけ plan へ残す）。

## 実装記録（2026-07-04）

### 前提の訂正: Yes/No の供給源は FACT:MNAM（実測）

一時 runner（Mutagen 0.53.1、scratchpad 上の使い捨て console。repo へは入れない）で inigo.esp の全 FACT を読んだ。

- inigo.esp に FULL 名が "Yes"・"No" の勢力は存在しない。"Yes"・"No" は勢力 `InigoFollowOnHorseFaction` の階級称号（FACT:MNAM、ordinal 0="No"・1="Yes"・2="Yes, but waiting for steed"）に由来する。
- `record_type_master` は FACT:MNAM / FNAM も固有名 box へ振り分けるため、誤爆の経路（固有名フェーズで AI 訳 → 機械置換辞書 → 本文の文頭 Yes/No に当たる）は背景の記述どおり成立する。誤りは供給 field の特定（FULL でなく MNAM）だけで、対策の設計は変わらない。

### 内部フラグ勢力の除外（依頼 2）: 基準不成立で停止

Hidden from PC flag 基準の成立確認の観測結果:

- inigo.esp: 勢力 8 件中 hidden=1 は `InigoFollowDistanceFaction`（"Preferred Distance when folloing on foot"）の 1 件だけ。"Yes"・"No" を供給する `InigoFollowOnHorseFaction` と `InigoRideWithoutPlayerFaction` は hidden=0 で、基準に該当しない。
- Skyrim.esm: FULL 名つき勢力 435 件中 hidden=1 は 163 件。うち 154 件は xTranslator 辞書（master_term の供給源 XML）に既訳がある。hidden=1 には "Thieves Guild"・"Khajiit Traders" のような本文に出る実在勢力名が含まれる。
- 判定: 完了定義 2 の 2 条件（対象勢力が基準に該当する、実在勢力名が大量に落ちない）の両方が不成立。管理用勢力の判別に Hidden from PC flag は使えないため、依頼 2 は不成立で停止する。抽出器（tools/extractor）は変更しない。同じ誤爆は依頼 1 の stoplist が止める（"Yes"・"No" は stoplist に含まれることを確認済み）。

### 実装で固定した設計判断

- stopwords-en.txt の置き場所: `assets/stopwords-en.txt`。上流と byte 一致を保ち（取得時 sha256 は選定節の値と一致）、出典・取得日・checksum・MIT 全文は併置の `assets/stopwords-en.LICENSE` に記録した。ディレクトリマップ `assets/CLAUDE.md` を新設して両ファイルの役割を記した。
- 照合規則の実装位置: 新 package を作らず、置換コア `internal/core/dictionary` へ `Stoplist`（`ParseStoplist`・`Blocks`）として統合した（利用者指示）。規則は「原語全体を小文字化して一語一致。複数語の固有名には当てない。nil は選別なし」。dictionary パッケージのカバレッジは 100% を維持。
- 供給源選別の共通化: `engine.translationVocabulary` を新設し、`LoadDictionary`（注入）と `mentionDetector`（言及）の両方がこの 1 つの読み出しを通る。単体テスト `TestStoplistFiltersInjectionAndMentionAlike` が「stoplist 語は置換も言及もされず、stoplist 外は両方に残る」を同一 store で固定する。
- 配線: stoplist の asset 読み込みは composition root（`internal/bootstrap`）が行い、`engine.New` の引数で注入する。合成 harness は実ファイルから切り離した最小リスト（yes・no）を使い（役割語 `syntheticRoleSpeech` と同じ方針）、goldcap は `-stopwords` flag（既定 `assets/stopwords-en.txt`）で実リストを読む。
- 観測ログ: 追加しない。選別は純粋ルール＋共通供給点の単体テストと golden（プロンプト・言及テーブル）で観測でき、実行時にしか確定しない分岐が残らないため。
- golden の観測点拡張: `internal/harness` の DB 最終状態ダンプへ `narration_mention`・`line_mention` を追加した（自然キーと相手の原語で捕獲）。分岐元 golden にはこの 2 節が無いため、実データ差分ではこの節追加を既知差分として除いて比較する。

### 最終検証（2026-07-04、通過）

- 単体テスト・lint: `npm run verify:backend` 通過（go test 16 package・arch-lint 違反 0・境界走査違反 0）。`npm run lint:backend`（format・vet・static・module を含む）も通過。`internal/core/dictionary` はカバレッジ 100%。
- 合成 golden（`internal/harness/testdata/synthetic.golden`、再生成済み）: fixture へ FACT:MNAM の "Yes"・"No"、文頭 Yes の台詞、文頭 No の叙述文を追加した。proper_noun に AI 訳（訳002・訳003）が付いても、プロンプトでは "Yes, the road is clear now."・"No matter what the weather..." が原文のままで、言及 2 節に Yes/No が現れない。Aventus・Grelod の置換・言及は従来どおり残る。差分は fixture 追加・stoplist 効果・言及節追加・訳番号のずれに限ることを目視で確認した。
- 実データ（inigo.esp）: 分岐元 commit `1aeb2f18` の golden を `scripts/golden/capture.sh` で捕獲し、変更後の捕獲と行・単語レベルで機械分類した。
  - `translated_count` は 8803 で不変。行の増減・不揃いな差分は 0。
  - 変更行内の置換対はすべて stoplist 語の還元だけ: 訳031→Yes（約 205 箇所）、訳030→No（約 172 箇所）、やあ→Hello（34）、鉱山→Mine（7）、炎→Fire（6）、吸血鬼化→Turn（3）、下へ→Down（1）。後ろ 5 対は master_term 側の「名前全体が一般語 1 語」誤爆（選定節の想定クラス）で、"Fire an arrow"（矢を放て）を「炎 an arrow」にする類の壊れた入力が解消した。
  - 挿入行は新設の言及 2 節（1601 行）だけで、Yes/No への言及は 0 件。Riften の言及は 26 件、リフテン の注入は 142 箇所で正例は維持。
  - proper_noun には "Yes"（訳031）・"No"（訳030）が FACT 由来のまま残る。依頼 2 が不成立で停止したため抽出は変わらず、訳は付くが注入・言及はされない（完了定義 2(c) は不成立時の扱いに従う）。

## 出口記録（finalization、2026-07-04）

### 作業 commit

- commit hash: `9ce4cdb1`（`claude/dictionary-false-positive-guard`、17 files、feat(dictionary)）。
- 検証結果: `npm run verify:backend`・`npm run lint:backend` 通過（最終検証節のとおり）。
- 残留リスク: 残留リスク節のとおり（stopwords-en に無い一般語 1 語の固有名、複数語の管理用文字列）。

### 正本化判断

- 反映対象候補: `docs/architecture.md` §4.1（`.go-arch-lint.yml` へ bootstrap・harness・goldcap → dictionary の import 許可を追加した）。
- 判断: 反映不要。層構成・Wails 境界・C#↔Go schema は不変で、import の追加は「core の純粋ルールを上位層が一方向に import する」既存パターン内に収まる。依存の機械正本は `.go-arch-lint.yml` で固定済み（違反 0）。前 task narration-line-mention-linking が engine → mention の追加で §4.1 の列挙を更新しなかった前例に従う。
- 人間承認状態: 承認依頼なし（構造不変のため正本反映を起こさない。恒久仕様の承認対象なし）。

### local merge・merge 後検証

- merge: `master` へ `git merge --no-ff claude/dictionary-false-positive-guard`（merge commit `a38f4bf2`）。conflict なし。
- merge 後検証: master 上で `npm run verify:backend` 通過（exit 0）。conflict が無く code は作業 branch と同一のため、実データ golden の再捕獲は行わない。

### 残留リスク

- 「名前全体が一般語 1 語」でも stopwords-en に無い語（例: inigo.esp の階級称号 Summonable・Close・Ranged・Health、Skyrim 由来の Chest・Door 等）は従来どおり辞書へ載り、文頭の大文字出現に当たり得る。リスト独自拡張は本 task の除外範囲（利用者指示でリスト選定で賄う）。
- 複数語の管理用文字列（例: "Yes, but waiting for steed"・"Not Summonable"）は stoplist の対象外だが、本文に全文一致で現れる可能性が低く、語境界・大小区別の既存緩和で実害は観測されていない。
