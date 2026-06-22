# 口調抽出 PoC 結果（確定記録、2026-06-21）

本書は、口調の機械抽出が実データで成立するかを確かめた PoC の確定記録である。確定した判断と、各判断を自分で再現して事実だと確かめる手順をまとめる。途中で捨てた版（v1 の対人スコア潰れ、龍語マーカーの過学習、汎用 7 話者サンプルでの計画）は残さない。概念は `tone-concept-model.md`、信号戦略は `persona-signal-map.md`、設計は `persona-design.md`、コードは `cmd/poc-tone/main.go` を参照する。

## 前提

- データ: skyrim.esm を extractor（mutagen）で抽出した `db/poc-skyrim.sqlite3`（line 34427 件）。実 mod 確認は `db/poc-inigo.sqlite3`。
- 実装: `cmd/poc-tone`（純 Go、`go run ./cmd/poc-tone` で再現）。純粋関数で 2 軸スコア・印・分散・語彙マーカー・voice 気質融合を出す。特徴量の定義（丁寧定型・命令文・罵倒・感嘆・引き伸ばし・NRC 感情語）はコードが正本。
- 軸: 対人態度軸（丁寧定型・命令文・罵倒）と感情表出軸（感嘆符・引き伸ばし・NRC 感情語）。命令文は prose の品詞解析で文頭の動詞原形を検出する。

## 信頼度の基準（印）

「印」はマーカー（命令文・罵倒・丁寧定型）を含む台詞数で、対人スコアの信頼度を表す。対人軸は印が 10 以上のキャラだけを本文で評価する。印が少ないキャラ（1〜9）の対人スコアは ±1.00 へ振れやすいため、本文では参考外にし、voice 気質 prior で代替する（後述「メタと本文の融合」）。感情軸は印に依らず常時測れる。これは PoC で得た運用基準である。

## ライセンス注意

NRC EmoLex は研究用途のみのライセンスである。本番組み込み時は商用ライセンス確認か、代替辞書（Empath は MIT）への差し替えが要る。

## 無作為検証（2026-06-21）

### 目的と方法

手選びの 13 NPC（枠・破れ疑い）に潜む選択バイアスを外し、対人軸の判別が一般のキャラへ及ぶかを確かめた。台詞 100 以上の話者を部分抽出せず信頼層 31 人を悉皆（しっかい）で対象にした。これは無作為抽出よりも選択バイアスが小さい（母集団全部を取るため、私の好みが入らない）。実行は `go run ./cmd/poc-tone verify`。

### 結果（対人軸で昇順、尊大→丁寧）

| 話者 | 台詞 | 印 | 基底口調 | 対人 | 感情 | 感情σ |
|---|---|---|---|---|---|---|
| Ralof | 142 | 11 | ぞんざい | −0.64 | 0.50 | 0.554 |
| Ulfric | 408 | 50 | ぞんざい | −0.52 | 0.50 | 0.671 |
| Hadvar | 154 | 8 | ぞんざい | −0.50 | 0.50 | 0.636 |
| Delphine | 389 | 28 | 冷然・見下し | −0.43 | 0.00 | 0.561 |
| Astrid | 312 | 41 | 冷然・見下し | −0.32 | 0.40 | 0.595 |
| MercerFrey | 148 | 14 | 冷然・見下し | −0.29 | 0.33 | 0.743 |
| GeneralTullius | 281 | 44 | ぞんざい | −0.27 | 0.50 | 0.767 |
| Enthir | 129 | 13 | 冷然・見下し | −0.23 | 0.00 | 0.631 |
| Brynjolf | 436 | 44 | 淡々・実務 | −0.09 | 0.33 | 0.607 |
| AelaTheHuntress | 239 | 20 | 淡々・実務 | 0.00 | 0.00 | 0.530 |
| Maven | 158 | 16 | 淡々・実務 | 0.00 | 0.00 | 0.615 |
| Arngeir | 248 | 23 | 慇懃・端正 | 0.22 | 0.33 | 0.820 |
| Mjoll | 117 | 15 | 物腰やわ | 0.33 | 0.50 | 0.692 |
| Paarthurnax | 160 | 19 | 慇懃・端正 | 0.37 | 0.00 | 0.341 |
| Erandur | 212 | 37 | 物腰やわ | 0.41 | 0.50 | 0.646 |
| Esbern | 222 | 21 | 慇懃・端正 | 0.43 | 0.33 | 0.694 |
| KodlakWhitemane | 111 | 14 | 物腰やわ | 0.43 | 0.50 | 0.714 |
| MirabelleErvine | 114 | 15 | 慇懃・端正 | 0.47 | 0.00 | 0.604 |
| Karliah | 256 | 33 | 物腰やわ | 0.52 | 0.50 | 0.766 |
| Nazir | 203 | 19 | 物腰やわ | 0.68 | 0.50 | 0.617 |
| Tolfdir | 238 | 51 | 物腰やわ | 0.69 | 0.45 | 0.615 |
| ArnielGane | 121 | 21 | 慇懃・端正 | 0.81 | 0.00 | 0.523 |
| Faralda | 134 | 22 | 慇懃・端正 | 0.82 | 0.33 | 0.675 |
| UraggroShub | 105 | 11 | 慇懃・端正 | 0.82 | 0.00 | 0.414 |

（印 10 未満で参考外: Hadvar 8、Rikke 7、Galmar 9、Vilkas 9。Vex・DelvinMallory・Farkas・Balgruuf は印 10 以上で中位のため割愛。全 31 話者の生表は実行出力にある。）

所見の集計（印 10 以上の 27 話者）:
- 対人帯: 尊大 7 / 中立 6 / 丁寧 14。中立に潰れたのは 22%。v1 は全員潰れていたので、無作為標本でも判別が成立した。
- 基底口調セル: 物腰やわ 7・慇懃端正 7・淡々実務 5・冷然見下し 4・ぞんざい 3・平明 1。6 セルへ分散し、1 マスへ偏らなかった。

### 妥当だった点（サニティ合格）

- 丁寧側の妥当性: 学者・神官・賢者が丁寧へ寄った。Tolfdir（学院教師）・Arngeir（グレイビアード）・Esbern（ブレイズの古老）・Paarthurnax（龍の賢者）・Erandur（マラ神官）・Kodlak（同胞団の長）・Karliah が慇懃・端正か物腰やわ。
- 冷血側の妥当性: Delphine・Astrid・MercerFrey が冷然・見下しに揃った。三者とも統率・裏切り役で、人を信用しない口ぶりが負へ出た。

### 無作為検証で見つけた新しい失敗

- 誘導命令の尊大誤判定（新規・重要）: Hadvar（−0.50）・Ralof（−0.64）が「ぞんざい」へ落ちた。実体は導入部の護衛 NPC で、`Come on, this way.`・`Come on! We need to get inside before that dragon comes back!` のような道案内の命令が多い。命令文＝尊大という規則が、誘導の命令を威圧と取り違えた。手選びでは護衛役を選ばなかったため隠れていた失敗である。
- 含意の尊大（既知・継続）: Maven は印 16 で中立 0.00 のまま。無作為標本でも未解決を確認した。
- Khajiit 三人称マーカーの誤爆（新規）: `this one` の部分一致が普通の英語を拾い、非 Khajiit の 10 話者で発火した。Delphine（種族 Breton）の `This one's mine.` を Khajiit 自己言及と誤検出した。種族条件、または主語位置の自己言及への限定が要る。

### 検証の結論

対人軸の判別は、恣意選択を外した母集団でも一般へ及んだ（主目的を達成）。一方で、命令文の意味区別（威圧の命令と誘導の命令を分ける）と、本文マーカーの精度（部分一致の誤爆を防ぐ）に新しい課題が出た。

## 修正と再検証（失敗 1・失敗 2、2026-06-21）

無作為検証で見つけた失敗 2 つを直し、同じ母集団で再検証した。

### 失敗 1: 誘導命令を威圧から外した

包括命令（`let's`・`let us` で始まる協調の命令）と、道案内の定型（`come on`・`this way`・`follow me`・`stay close`・`keep an eye`・`look around` など）で始まる文を、威圧の命令から除外した。実装は `isGuidanceImperative`。

| 話者 | 修正前 | 修正後 | 評価 |
|---|---|---|---|
| Hadvar | −0.50 ぞんざい | +0.20 物腰やわ | 護衛 NPC が尊大から外れた |
| Ralof | −0.64 ぞんざい | −0.14 平明 | 同上 |
| Ulfric | −0.52 ぞんざい | −0.48 ぞんざい | 本物の威圧は保持 |
| GeneralTullius | −0.27 ぞんざい | −0.22 ぞんざい | 同上 |
| Astrid | −0.32 冷然 | −0.30 冷然 | 同上 |

Hadvar（印 8→5）・Ralof（印 11→7）は誘導が印から外れ、印 10 未満で参考外になった。これは妥当である。護衛 NPC は道案内が主で、対人態度を測れるほどの威圧・丁寧の台詞を持たない。

### 失敗 2: Khajiit 三人称を種族レコードで限定した

`this one` の本文一致だけでは Breton 等の指示語を拾うため、種族レコードが Khajiit のときだけ Khajiit 三人称マーカーを出すよう gate した。実装は `lexicalFromText` に種族を渡す変更。

- 非 Khajiit の 10 話者（Delphine ら）から Khajiit 三人称が全て消えた。Delphine（種族 Breton）の `This one's mine.` は除外された。
- 実 Khajiit は保持を確認した。Jzargo・DA13PeryiteMonk（ともに KhajiitRace）は `this one` でマーカーを残す。

### 再検証の所見

- 印 10 以上の 25 話者で、対人帯は尊大 4 / 中立 5 / 丁寧 16。尊大 4 は Ulfric・Astrid・Tullius・Delphine で、全て統率・冷血の役。誤って尊大に入っていた護衛は抜けた。
- 中立に潰れた割合は 20%（修正前 20%）。判別の質は保たれた。
- マーカー列は 老人(種族)・龍(種族)・古英語(本文) など正当なものだけになった。

## メタデータ信号の棚卸し（2026-06-21）

種族 gate の修正を機に、Skyrim メタデータで判定できる信号を DB スキーマから棚卸しした。マーカーは現状 race と voice の一部しか使っておらず、信頼できる信号を大量に取りこぼしている。

### DB に入っていて未活用

- voice_type の EditorID: 全 912 話者の 786（86%）が汎用 voice を持ち、制作側が付けた temperament が EditorID に符号化されている（`EvenToned` 平静・`Condescending` 横柄・`Brute` 粗暴・`SlyCynical` 皮肉・`Coward` 臆病・`Commander` 指揮官・`Sultry` 色気・`OldKindly` 温厚老人・`OldGrumpy` 気難老人・`YoungEager` 若く前のめり 等）。本文と独立した源のため融合しても過学習にならない。Male/Female 接頭辞から性別も復元できる（`speaker.sex` は空でも voice から取れる）。
- 注意: 「汎用ボイス 86%」は汎用 NPC の割合ではない。台詞数の多い汎用ボイス話者は上位が全員固有名（Aela・Tolfdir・Farkas・Festus・Calcelmo）で、固有キャラが汎用ボイスを再利用しているのが大半である。よって voice_type の役割は「汎用 NPC の穴埋め」ではなく、固有キャラへの第 2 の独立信号である。用途は、本文 2 軸の相互検証・精緻化と、台詞が少なすぎて本文で測れない固有キャラ（1 台詞の Vilod・Edith 等）の分類の 2 つ。
- 真の汎用・テンプレ NPC の所在: line 34427 のうち話者リンクありは 21326（62%）で、残り 38%（13101 行）は名前付き話者に紐づかない共有・プール台詞である。真の汎用 NPC はこの未リンク行に散り、個別 speaker として表に現れにくい。本文でもボイスでも個別には拾えない。
- ユニーク/テンプレの判別: 正しい判別軸は template 継承だが、`template_speaker_id` も `speaker_kind` も空（0/912）で未抽出である。archetype 風 edid は 155 件あるが `DB09Noble1` 等の固有名も拾う粗い proxy にとどまる。正確な分離には extractor が template を埋める必要がある。
- 独立な 2 源の相互検証例: Tolfdir（`MaleOldKindly` ↔ 物腰やわ +0.76）、Enthir（`MaleSlyCynical` ↔ 冷然 −0.23）、ArnielGane（`MaleCoward` ↔ 慇懃 +0.81）が一致した。Farkas（`MaleBrute` ↔ 慇懃 +0.20）は不一致で、voice は粗い prior、本文が精緻化する関係を示す。
- faction 所属: `speaker_faction` に 843 話者の所属が入る（Brynjolf → `ThievesGuildFaction` 等）。役割・ギルドの信号で、未使用。

### esm にあるが extractor が未抽出（抽出側の作業が要る）

- `speaker.occupation` は空。class(CLAS) から取れるが抽出していない（faction で部分代替は可能）。
- AI データ（aggression・confidence・morality・assistance）はスキーマに無い。aggression は対人軸（威圧）を直接符号化する最強候補だが、extractor 拡張が要る。

### 設計上の含意

- 空の `nature` 列（voice_type 99・race 33・faction 695）は、各メタの含意を人手注釈する curation slot である。99 個の voice_type をトーン含意へ一度だけ辞書化して `voice_type.nature` に入れれば、コードに文字列照合を散らさず voice をトーン源にできる（master 辞書と同じ作り）。
- persona/tone の源は 2 チャネルに分かれる。メタデータ系（voice_type 86%＋faction＋race）は汎用 NPC に強い temperament prior、本文系（2 軸抽出）は固有・多弁 NPC に強い。メタを prior、本文が精緻化する融合が妥当である。

## 含意の尊大の機械検出を試した結論（2026-06-21）

Maven 型（丁寧な言葉で尊大、印 16 でも中立 0.00）を機械で取れるか検証した。採用条件は、過学習しないこと、他の NPC を尊大にしないこと、の 2 つである。Maven の実台詞から含意のパターンを 3 種に分け、各々を候補マーカーにして全 31 話者で副作用を測った。

### 候補 3 種と判定

- 含意脅し（条件付き脅しの定型。`i'd advise`・`keep that in mind`・`or else`・`make me angry` 等）: Maven を射抜くが、温厚な老教師 Tolfdir を 3 件誤爆した（`Eagerness must be tempered with caution, or else disaster is inevitable.` は教育的注意で脅しではない）。他 NPC を尊大にするため却下する。
- 軽蔑語（`stupid`・`pitiful`・`pathetic`・`beggar` 等）: 部分一致が描写用法を過大計上した。Brynjolf の 6 件中、真の見下しは `You're pathetic!` の 1 件のみで、残りは `scam a beggar`（描写）・`foolish request`（依頼の形容）だった。雑音が多く限界的である。
- 支配言明（権力・所有の定型。`in my pocket`・`the jarl's ear`・`my approval`・`to me of all people` 等）: 31 話者で Maven 以外ゼロの安全な信号だった。しかし全 912 話者でも 8 人しか出ず、Elisif・Rune など尊大でない者も混じる。一般信号として成立せず、実質 Maven 専用ルールになるため、過学習に当たる。

### 結論: 表層の機械手法では取れない

3 条件（Maven を射抜く・過学習しない・他 NPC を無傷に保つ）を同時に満たす表層特徴は無い。再現性のある信号（含意脅し）は善意の文を脅しと誤爆し、安全な信号（支配言明）は Maven 専用で過学習になる。採用条件に照らし、機械ルールは採用しない。

含意の尊大は、語・文構造の表層では届かず、文脈と発話意図の理解を要する領域である。LLM で特別対応する案も検討したが、採らない。理由は 2 つある。第 1 に、効果が小さい（中立へ分類されても翻訳口調が少し弱まるだけで、対象も少数）。第 2 に、分岐の判定が循環する。LLM へ回す条件を「含意の尊大か」にすると、その判定自体が含意の検出を要し、機械が決められない当のものを機械で評価することになる。Maven は印 16 で見かけ上は中立と確信的に出るため、信頼度の代理指標でも選り分けられない。よって特別対応はせず、Maven 型は中立へ分類される既知の許容誤差として受け入れる。

## メタと本文の融合（voice 気質 prior、2026-06-21）

`persona-signal-map.md` の縮退合成を PoC へ実装し、汎用ボイス込みの分類を実証した。

### 実装

- voice 気質辞書: voice_type の EditorID の気質部分（Male/Female を除く）18 種を、対人軸の事前値（prior）と表示名へ写した（`Condescending` −0.6 横柄、`Brute` −0.5 粗暴、`SlyCynical` −0.4 皮肉、`OldGrumpy` −0.3 気難老、`Coward` +0.4 臆病、`OldKindly` +0.4 温厚老、`EvenToned`・`Sultry`・`Commoner` 等は 0 中立）。
- 融合（`fuseAttitude`）: 印 10 以上なら本文 2 軸で対人を決め、未満なら voice 気質 prior で代替し、固有 voice で prior も無ければ保留する。感情軸は全経路で本文から測る。
- 実行: `go run ./cmd/poc-tone sparse`（voice prior 経路）と `go run ./cmd/poc-tone verify`（多弁の本文経路）。

### voice prior 経路の実証（台詞少の固有キャラ）

| 話者 | 台詞 | 印 | 経路 | 基底口調 | 対人 | voice 気質 |
| --- | --- | --- | --- | --- | --- | --- |
| Nazeem | 5 | 2 | voice | 冷然・見下し | −0.60 | 横柄 |
| Grelod | 14 | 1 | voice | ぞんざい | −0.30（感情 1.00） | 気難老 |
| Edorfin | 1 | 0 | voice | 冷然・見下し | −0.50 | 高慢 |
| ConstanceMichel | 11 | 4 | voice | 淡々・実務 | 0.00 | 前のめり |
| Curwe | 1 | 0 | voice | 淡々・実務 | 0.00 | 平静 |

横柄・高慢・気難の気質 voice は対人を妥当に与えた（Nazeem は皮肉屋の農夫、Edorfin は 1 台詞だけで voice が唯一の信号）。気質を持たない voice（前のめり・平静）は背伸びせず中立にした。感情軸は本文から取れ、Grelod の激情（1.00）を 14 台詞でも拾えた。

### 多弁の本文経路（無回帰と相互検証）

- 無回帰: 多弁 31 人の本文経路の結果は融合前と同一だった（Ulfric −0.48、所見も尊大 4 / 中立 5 / 丁寧 16 のまま）。融合は多弁側を壊さない。
- 相互検証（本文経路で voice 気質を並記）: 一致は Tolfdir（温厚老 ↔ 物腰やわ）・ArnielGane（臆病 ↔ 慇懃）。不一致は Farkas（粗暴 ↔ 慇懃）・AelaTheHuntress（指揮 ↔ 慇懃）・Faralda（色気 ↔ 慇懃）・Enthir（皮肉 ↔ 淡々）。
- 所見: 多弁キャラでは voice 気質と本文がしばしば食い違い、本文の方が精緻である。voice 気質は粗い prior で、印が十分なら本文が正しく上書きする（融合設計どおり）。voice prior は台詞が無い時の代役として有効で、台詞がある時の精度には及ばない。

### 残課題（融合）

- 閾値の崖: 印 7〜9 で経路が変わる。気質を持つ voice（Rikke=指揮）は voice prior で妥当だが、気質を持たない voice（Vilkas・Ralof=Nord）は中立 0.00 になり、弱い本文証拠を捨てる。固有 voice（Mercer・Galmar）は保留で薄い本文値を低信頼で保持する。境界の連続化は設計で詰める。
- prior 値の較正: 多弁での不一致は voice prior がやや尊大寄りに振れることを示す。本実装で prior の符号・強さを較正する。

## 実 mod での動作確認（inigo.esp、2026-06-21）

フォロワー mod `inigo.esp` を extractor で抽出し（`db/poc-inigo.sqlite3`、line 7696 件、master 連鎖で skyrim.esm の話者属性を解決）、融合分類をかけた。狙いは 2 つ。追加 NPC の Inigo を判定できるか、既存 NPC の persona を出せるか（追加台詞だけならセリフ不足になるはず）。再現は `dotnet run --project tools/extractor -- --data dictionaries/Data --plugin inigo.esp --sqlite db/poc-inigo.sqlite3` と `POC_DB=db/poc-inigo.sqlite3 go run ./cmd/poc-tone all`。

### 結果

| 話者 | 種別 | 台詞 | 印 | 経路 | 基底口調 | 対人 | voice 気質 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Inigounique | 追加 NPC | 3980 | 452 | 本文 | 物腰やわ | +0.69 | 固有(mod) |
| Langleyunique | 追加 NPC | 402 | 44 | 本文 | 物腰やわ | +0.50 | 固有(mod) |
| Nazeem | 既存 | 12 | 4 | voice | 冷然・見下し | −0.60 | 横柄 |
| Mjoll | 既存 | 45 | 8 | voice | 淡々・実務 | 0.00 | —（FemaleNord） |
| Jzargo | 既存 | 42 | 7 | voice | 淡々・実務 | 0.00 | —（MaleKhajiit） |

### 読み取れたこと

- 追加 NPC は本文で出せる: Inigo は 3980 台詞・印 452 で物腰やわ（温厚）と分類された。Khajiit 種族・三人称マーカーも正しく発火した。新規 NPC の判定は対話で成立する。
- 既存 NPC は voice が気質を持てば、セリフ不足でも一致する: Nazeem は inigo 追加の 12 台詞（印 4）だけだが、master 連鎖で解決した voice 横柄から 冷然・見下しになり、skyrim.esm 版（冷然・見下し）と一致した。
- 既存 NPC は voice が中立だと、セリフ不足で persona が落ちる: Mjoll は voice が FemaleNord（気質なし）で、追加 45 台詞では印 8 と不足し中立になった。skyrim.esm 版は物腰やわ +0.54（印 13）だった。彼女の温厚さは台詞に宿るが、その台詞が skyrim.esm 側にあり inigo 側に無い。これがセリフ不足そのものである。
- creature の雑音: 馬・犬（Cr 声型）は 1〜数台詞で保留経路に落ち、物腰やわ +1.00 等の極端値を出した。creature は Cr 声型で対象外にすべきである。

### 設計への含意

既存 NPC の persona は、解決した話者の同一性（master の NPC）で全プラグインの台詞を束ねて出すべきである。mod 単体の追加台詞だけでは Mjoll 型がセリフ不足になる。実運用では、既存 NPC は base game で算出済みの persona を引き、新規 NPC（Inigo）だけ新たに算出する形になる。voice prior は束ねても薄い時の最後の代役である。

## 本実装で残る較正（確定済みは除く）

設計（`persona-design.md`）と実装計画（`implementation-scope.md`）で決着した項目（voice 気質の辞書化、含意の尊大の不対応、extractor 拡張の対象外）は除く。本実装の較正で詰める残課題だけを記す。

- 閾値の崖: 印 7〜9 で本文経路と voice prior 経路が切り替わる。気質を持たない voice（Nord 等）は中立へ落ち、弱い本文証拠を捨てる。境界の連続化を設計で詰める。
- prior 値の較正: 多弁での不一致は voice prior がやや尊大寄りに振れることを示す。prior の符号・強さを較正する。
- 頑健統計: 感情σを中央絶対偏差へ替え、外れ値（絶叫 1 件）への過敏を抑える。
- 語彙マーカー拡張: 龍語は Dovahzul 辞書、役割語の出力類型（一人称・語尾）を外部資料から。

## 再現と事実確認

各判断は次のコマンドで再現でき、報告が事実だと自分で確認できる。コマンドはリポジトリ直下で実行する。

| 主張 | 再現コマンド | 確認できる事実 |
| --- | --- | --- |
| 無作為検証は 31 人、中立潰れ 20% | `go run ./cmd/poc-tone verify` | 末尾の所見に「印10以上で評価可能: 25 話者」「中立 …20%」 |
| 護衛 NPC の誤判定は道案内が原因 | `sqlite3 db/poc-skyrim.sqlite3 "SELECT l.source FROM line l JOIN line_speaker ls ON ls.line_id=l.id JOIN speaker s ON s.id=ls.speaker_id WHERE s.edid='Hadvar' AND l.source LIKE 'Come %'"` | `Come on, this way.` 等の道案内が返る |
| Khajiit 誤爆は Breton の指示語 | `sqlite3 db/poc-skyrim.sqlite3 "SELECT r.edid FROM speaker s JOIN race r ON s.race_id=r.id WHERE s.edid='Delphine'"` | `This one's mine.` を話す Delphine の種族が `BretonRace` |
| voice 86%が汎用気質 | `sqlite3 db/poc-skyrim.sqlite3 "SELECT CASE WHEN v.edid LIKE '%Unique%' OR v.edid LIKE 'Cr%' OR v.edid LIKE 'SPECIAL%' OR v.edid='' OR v.edid IS NULL THEN '固有/特殊/無' ELSE '汎用気質' END k, count(*) FROM speaker s LEFT JOIN voice_type v ON s.voice_type_id=v.id GROUP BY k"` | 汎用 786 / 固有等 126（voice 無し 49 を `IS NULL` で含める） |
| 台詞リンクは 62% | `sqlite3 db/poc-skyrim.sqlite3 "SELECT (SELECT count(*) FROM line), (SELECT count(DISTINCT line_id) FROM line_speaker)"` | 34427 と 21326 |
| 含意脅しは善意を誤爆 | `sqlite3 db/poc-skyrim.sqlite3 "SELECT l.source FROM line l JOIN line_speaker ls ON ls.line_id=l.id JOIN speaker s ON s.id=ls.speaker_id WHERE s.edid='Tolfdir' AND LOWER(l.source) LIKE '%or else%'"` | 温厚な老教師の善意の文が返る |
| voice prior が端役に効く | `go run ./cmd/poc-tone sparse` | Nazeem −0.60 横柄、Edorfin −0.50 高慢 |
| 融合は多弁側を壊さない | `go run ./cmd/poc-tone verify` | 多弁の数値が修正後と同一 |
| 実 mod の追加 NPC を分類 | `POC_DB=db/poc-inigo.sqlite3 go run ./cmd/poc-tone all` | Inigo +0.69 物腰やわ、Nazeem −0.60 横柄 |
