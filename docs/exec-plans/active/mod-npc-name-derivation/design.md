# Design: mod-npc-name-derivation

`design.md` は「どう実装し、どう直すか」を人間が読んで判断するための説明を持つ。要求は `plan.md`、確定仕様は `spec.md`、再現確認・原因究明は `investigation.md` が持つ。
`design.md` と `spec.md` が食い違う場合は `spec.md` を優先する。

---

## R-1 Data フォルダ全体の英日対を既訳の辞書にする

### 現況の理解

既訳を受ける器は 2 つある。どちらも入力が `extracted_field` に閉じている。

- `master_term` は横断辞書で、`source` を一意キーに持つ。固有名の箱に割り当てられた record の `FULL` から完全形を書き、`termderive` で人名の部分形を足す（`DeriveMasterTerms`）。
- `reference_translation` は既訳の照合表で、`(rec, field, source)` を一意キーに持つ。同一原文でも record 種別で訳し分けるため `form_id` では絞らない、と doc コメントが明記している。取込（`LoadReferenceTranslations`）は `extracted_field` の `dest` 非空行を写すだけで、変換をしていない。

`extracted_field` は利用者が選んだ plugin の翻訳対象の原文を持つ器で、書く経路は `internal/api/app.go` の `DotnetExtractor.Extract` 1 箇所である。C# 抽出器は `--data` と `--plugin` を必須引数に取り、1 回の起動で 1 plugin を扱う。日本語は record の string 参照を english と japanese の 2 言語で解決して得る。mod の record は公式日本語 Strings に無いので日本語側が必ず空になり、供給が 0 件になる。

C# と Go の受け渡しは中心 DB の table だけである。C# は `--sqlite` で同じ DB を開き、schema を ensure してから書く。

要求が扱う対象の単位と、受け皿が持つキーを並べる。

| | 単位 |
| --- | --- |
| 要求が扱う対象 | Data フォルダにある plugin 1 本の、record 1 件の field 1 つについての英日対 |
| 受け皿が持つキー | `reference_translation` は `(rec, field, source)`、`master_term` は `source` |

どちらの受け皿も record 種別を含むキーを持つ。Strings ファイルは文字列 ID から文字列への表なので、ファイルだけを読んでも `rec` と `field` は取れない。英日対は plugin ごとの record 走査から取る。単位の食い違いは無い。

供給が選んだ plugin に閉じている一方で、画面警告の判定（`internal/api/strings_presence.go` の `GetStringsPresence`）は選んだ plugin の親フォルダ配下の `strings/` を見る。Data フォルダに公式の英日 Strings があれば警告は出ない。警告が立っている前提（Data フォルダ全体を供給源とみなす）と、供給の実装（選んだ plugin の所有分のみ）が食い違っている。

Data フォルダ全 10 本の走査は 18 秒で、英日対は 87,380 件だった（`investigation.md`）。翻訳本体に対して無視できる。

### あるべき形

翻訳対象の原文と既訳を、既にある 2 つの器へ分ける。table は足さない。

- `extracted_field` は、利用者が選んだ plugin の翻訳対象の原文だけを持つ。日本語の列を持たない。
- `reference_translation` は、Data フォルダにある全 plugin の英日対を持つ。english と japanese の両方を解決できた field だけを入れる。`(rec, field, source)` の一意キーは変えない。既訳は plugin をまたいで再利用するので、どの plugin 由来かを持たない。
- 英日対を作る走査は C# 抽出器が行い、`reference_translation` へ直接書く。`extracted_field` を経由しない。同じ Data フォルダで 2 回走らせても件数が増えない。
- 横断辞書の派生（`DeriveMasterTerms`）は、人名の対を `reference_translation` から取る。用法の集計は既訳の有無に依らない材料なので、これまでどおり翻訳対象の台詞の英語原文から取る。

適用範囲と永続性は次のとおり。`reference_translation` は Data フォルダの英日対の写しなので、対象 plugin を変えても中身が変わらず、実行をまたいで残る。`master_term` は横断辞書として残る。どちらも実行内の AI 訳を入れない。

### 変更点

- C# 抽出器: Data フォルダにある plugin を列挙して 1 本ずつ走査し、2 言語を解決できた field を `reference_translation` へ書く起動経路を足す。翻訳対象の原文を `extracted_field` へ書く既存の起動は残し、そこから日本語の解決を外す。
- `db/migrations`: `extracted_field` の `dest` 列を消す。既訳を持つのは `reference_translation` だけになる。
- `internal/api/app.go` の `prepareForTranslation`: 既訳を取り込む段を、Data フォルダ全 plugin を走査する段へ置き換える。対象 plugin の抽出より前に置く。既存の観測ログ `reference_supply_built` へ、走査した plugin 数を足す。
- `internal/engine/reference.go` の `LoadReferenceTranslations`: 消す。`extracted_field` から `reference_translation` へ写す経路が要らなくなる。
- `internal/engine/engine.go` の `DeriveMasterTerms`: 人名の対の入力を `extracted_field` から `reference_translation` へ移す。用法の集計の入力は `extracted_field` の台詞の英語原文のままにする。`record_type_master` による箱の判定と `termderive` の呼び出しは変えない。
- `internal/api/strings_presence.go` の `GetStringsPresence`: 判定の対象を Data フォルダ全体へ合わせる。選んだ plugin の英日 Strings が無くても、Data フォルダに公式の英日 Strings があれば既訳が立つので、警告が指す状態が「選んだ plugin に英日 Strings が無い」から「Data フォルダに英日 Strings が 1 本も無い」へ変わる。
- `internal/harness/fixture.go` と `internal/harness/extractor.go`: `extracted_field` の日本語で英日対を作っている合成 fixture を、`reference_translation` を直接置く形へ移す。

---

## R-2 実行内で確定した NPC 名から部分形の訳を作る

### 現況の理解

人名の部分形を作る規則は `internal/core/termderive` にある純粋ルールで、入力は NPC の英日対・用法の集計・既出原語の集合・設定である。名のみ（shrt）、二つ名の前部（byname）、姓名分割（two）の 3 種を作り、`safePair` と `landmine` で一般語との衝突を落とす。姓名分割は `NamePair.BaseGame` が真の対にだけ働く設計だが、`DeriveMasterTerms` は全ての対へ真を渡しており、doc コメントと実態が食い違っている。

派生を呼ぶ `DeriveMasterTerms` は、抽出の直後・固有名を AI で訳す段より前にある。入力は `extracted_field` の日本語側が非空の行だけなので、mod NPC は入力に入らない。mod NPC の氏名は固有名フェーズで AI 訳として確定するが、そこから部分形を作る経路が無い。R-1 を直しても、公式訳を持たない mod NPC については残る。

実行内の AI 訳を受ける器は `proper_noun` で、「同一実行内で AI 訳を留め、横断辞書 `master_term` へは昇格しない」を doc コメントで明記している。`internal/engine/proper_noun.go` の `SelectSupply` が、AI 訳の書き込み先を `proper_noun` に固定する純粋関数として不変条件を持つ。

機械置換辞書を組む `LoadDictionary` は `master_term` と `proper_noun` を合流する。呼ばれる箇所は 3 つで、`Engine.Run`（固有名フェーズの後）、`BatchRunner.planBodyRequests`、`internal/api/app.go` の結果表示である。batch の送信経路は 2 つあるが、どちらも `planBodyRequests` を通る。

| | 単位 |
| --- | --- |
| 要求が扱う対象 | 実行内で訳が確定した mod NPC 1 件の氏名と、そこから作る部分形 |
| 受け皿が持つキー | `proper_noun` は `(plugin, category, source)`。`category` に入るのは `rec` だけで、氏名と短縮名を分ける列が無い |

`proper_noun` は氏名（`FULL`）と短縮名（`SHRT`）を区別できない。区別は `extracted_field` と `(plugin, category = rec, source)` で結んで取り戻す。同じ結合を `internal/store/mention.go` の `LinkNarrationDescribed` と `internal/store/export.go` の `ProperNounPlacementsForExport` が既に使っている。

### あるべき形

人名の部分形を作る段を、固有名の訳が確定した後にもう 1 度置く。

- 入力は、実行内で訳が確定した対象 plugin の固有名のうち、NPC の氏名と短縮名である。既訳流用で確定した行と AI 訳で確定した行を区別しない。
- 部分形を作る規則は `termderive` をそのまま使う。既出原語の集合には横断辞書と、その実行で確定済みの固有名の両方を入れ、同じ原語に 2 つの訳が立たないようにする。
- 作った部分形は `proper_noun` へ、対象 plugin と同じ実行の行として書く。横断辞書へは書かない。`proper_noun` の不変条件を保つ。
- 部分形が機械置換辞書へ載る時点は、本文（叙述文・台詞）を訳す段より前とする。同期翻訳と batch の両方で満たす。

姓名分割は、訳の出どころが base ゲームかでは決めず、語形の防御だけで決める。要求 R-2 が挙げる「名だけ・苗字だけ」は、空白 2 語の氏名を割って初めて作れる。実行内の AI 訳を割らない扱いにすると、二つ名を持つ NPC と短縮名を持つ NPC しか救えず、要求の主な場合が残る。割ってよいかは、英語が空白 2 語であること、訳が中黒区切りで語数が一致すること、訳に漢字が無いこと、`safePair` と `landmine` を通ることで判定する。この判定は訳の作者が人間か AI かに依らない。

`NamePair.BaseGame` は姓名分割の可否を絞る旗だが、現況では呼び出し側が全ての対へ真を渡しており、絞りとして働いていない。旗を消して、判定を語形の防御へ一本化する。振る舞いは現況と変わらない。

### 変更点

- `internal/core/termderive/termderive.go`: `NamePair` から `BaseGame` を消し、`deriveTwo` の base ゲーム限定の条件を外す。姓名分割の可否は語形の防御（語数一致、中黒区切り、漢字なし、`safePair`、`landmine`）だけで決める。doc コメントを同じ内容へ書き直す。
- `internal/engine/engine.go` の `DeriveMasterTerms`: `NamePair` を組む箇所から旗の指定を外す。派生の入力と手順は変えない。
- `internal/engine/proper_noun.go`: 固有名フェーズの後に、確定した固有名から部分形を作って `proper_noun` へ書く関数を足す。
- `internal/engine/engine.go` の `Engine.Run`: 固有名フェーズの後・`LoadDictionary` の前に、上の関数を呼ぶ。
- `internal/engine/batch.go` の `planBodyRequests`: `LoadDictionary` の前に、同じ関数を呼ぶ。固有名 batch の反映後に本文を組む経路でも部分形が載る。

---

## R-3 台詞の口調指示から句点の禁止を外す

### 現況の理解

台詞の訳文の形を決めるのは、`directive` table の 口調 の指示文である。migration `db/migrations/0006_record_type_translation.sql` が既定値を seed し、`{traits}` へ話者の性質を差し込む。指示文のうち次の 2 文が句点を禁じている。

- 台詞は文末に句点を打たない。原文が疑問符・感嘆符で終わる時だけ ？ ！ を置く。
- 1 つの台詞に文が 2 つ以上ある時は、句点でつなげず全角空白で区切る。

`assets/role-speech-examples.tsv` の訳文 57 行は、句点を持つ行 0、読点を持つ行 24、全角空白で文を区切る行 20 である。例文は `- 例: 原文 → 訳文` の行として口調指示に載るので、指示文と同じ強さで効く。

commit `b4e06664` が、公式既訳の集計（台詞の末尾句点 0%）を根拠に意図して入れた既定値であり、実装の誤りではない。人間の判定は「公式訳の実態以前に日本語として不自然」である。

| | 単位 |
| --- | --- |
| 要求が扱う対象 | 台詞 1 件の訳文 |
| 受け皿が持つキー | `directive` の 口調 1 行のテキスト |

台詞は全て同じ 口調 の指示文を通るので、種別 1 行を変えれば台詞全件に効く。

### あるべき形

口調の指示文は、句点の有無を指定しない。翻訳モデルが本文から判断する。

- 句点を禁じる 2 文を置かない。
- 句点を求める文も置かない。指定を足さずに外すだけにする。
- 疑問符・感嘆符の扱い、一人称の表記、人物像に合う口調で訳すという他の指定は残す。
- 例文の訳文は、日本語として句点を打つ形へ戻す。例文は指示文と同じ強さで効くので、指示文だけを外しても例文が句点なしの形を教え続ける。

読点は原因が確定していないので扱わない。

### 変更点

- `db/migrations/0006_record_type_translation.sql`: 口調 の指示文から、句点を禁じる 2 文を消す。commit `b4e06664` と同じく seed のテキストを直接書き換える。
- `assets/role-speech-examples.tsv`: 訳文 57 行のうち、文として終わる行の末尾へ句点を戻す。全角空白で文を区切っている 20 行は、区切りを句点へ戻す。疑問符・感嘆符で終わる行はそのまま残す。
- `internal/harness/oracle_test.go`: 指示文の文面へ文字列一致しているアンカーが消した 2 文に掛かる場合、残る文へ移す。

---

## R-4 話者の性別を口調指示の行として出す

### 現況の理解

口調指示を組む純粋ルールは `internal/core/personatone` にある。組み立て点は 2 つで、`internal/engine/engine.go` の `namedPersona` と `freeTonePersona` が呼び分ける。

- `BuildToneTraits`: 名指し話者。性別は `model.LinePersonaInput` の `Sex` から取る。
- `BuildFreeToneTraits`: 話者を解決できない汎用台詞とプレイヤーの選択肢。性別は引数 `sex` から取る。

どちらも性別を役割語テンプレート（`assets/role-speech.tsv`）を引くキーとしてだけ使い、プロンプトへ載るのは引いた結果の一人称と言い回しである。`assets/role-speech.tsv` の `adult male *` 行と `adult * *` 行は「一人称は私・言い回し空」で完全に一致するため、成人男性と性別不明の出力が区別できない。

口調テンプレートの分布に欠陥は無いと判定済みである。`inigo.esp` の台詞 8545 件は、名指し話者 5217 件・既定値へ落ちた汎用 2472 件・プレイヤーの選択肢 856 件だった。名指し話者の 92.6% が同じ口調セルに寄るのは、`Inigounique` 単独で 4377 件を占める mod 側の話者構成による。分布は変更しない。

| | 単位 |
| --- | --- |
| 要求が扱う対象 | 台詞 1 件の話者の性別 |
| 受け皿が持つキー | 口調指示の箇条書きの 1 行 |

性別は箇条書きの行として出ていない。

### あるべき形

話者の性別を、口調指示の箇条書きの 1 行として出す。

- 性別が取れる話者は、男性か女性かを示す行を持つ。役割語から引いた一人称・言い回しとは別の行にする。
- 性別を取れない話者は、性別の行を持たない。取れない話者へ既定の性別を当てない。
- 名指し話者と、汎用台詞・プレイヤーの選択肢の両方で同じ形の行を出す。プレイヤーの性別が未設定なら、性別を取れない話者と同じ扱いにする。
- 役割語テンプレートの引き方は変えない。分布も変えない。

### 変更点

- `internal/core/personatone/personatone.go`: 性別から口調指示の 1 行を作る純粋関数を足し、`BuildToneTraits` と `BuildFreeToneTraits` の両方から呼ぶ。行を置く位置は性質文の後、役割語の前とする。
- `internal/core/personatone/personatone_test.go`: 性別あり・性別なしの両方で行の有無を確かめる単体テストを足す。

呼び出し側（`namedPersona`、`freeTonePersona`）は既に性別を渡しているので変えない。

---

## 検討が必要なこと

なし。
