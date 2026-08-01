# Design: translation-paraphrase-prompt-default

人間が採用した指示文は、`prompt_template.base_directive`ではなく、`directive`のkey `口調`が持つ`instruction`へ入れる。台詞の人物像と一緒にだけ送ることで、台詞以外の翻訳指示へ影響させない。

---

## R-1 SQLの`directive`でkey `口調`が持つ`instruction`を人間が採用した全文へ更新し、人物像の`{traits}`を埋めた台詞では翻訳AIへ送る指示文に反映し、人物像を作れない台詞と台詞以外の翻訳対象には反映しない

### 現況の理解

`db/migrations/0006_record_type_translation.sql`は、翻訳AIへ送る指示文を`directive`へkey単位で保存する。`INFO:NAM1`、`INFO:RNAM`、`DIAL:FULL`はkey `口調`を参照する。

`db/migrations/0017_dialogue_tone_directive.sql`は、key `口調`の既定値を現在の文面へ更新する。現在の文面は主語、常体、句点、全角空白を指定するが、採用した意訳指示と安全指示を持たない。

同期翻訳は`internal/engine/engine.go`の`Run`からkey `口調`を読み、`LinePersonas`が人物像を作れた台詞だけ、`{traits}`を埋めた指示文をbase指示の後へ合成する。人物像を作れない台詞は空の指示文を使う。

batch翻訳は`internal/engine/batch.go`の計画作成でも同じkey `口調`を読み、人物像を作れた台詞の送信内容へ合成する。

| | 単位 |
| --- | --- |
| 要求が扱う対象 | 人物像を作れた台詞へ送る指示文 |
| 受け皿が持つkey | `directive.key = '口調'` |

### あるべき形

key `口調`は、現在の人物像、主語、句読点の指示に加え、人間が採用した意訳指示と安全指示を持つ。

同期翻訳とbatch翻訳は、人物像を作れた`INFO:NAM1`、`INFO:RNAM`、`DIAL:FULL`へ同じ指示文を送る。人物像を作れない台詞と、台詞以外の翻訳対象は今回の指示文を送らない。

採用する全文は`spec.md`の「採用するinstruction」に固定する。

### 変更点

`db/migrations/0019_dialogue_paraphrase_directive.sql`を追加する。migrationはkey `口調`の`instruction`を、`spec.md`に固定した全文へ更新する。

`internal/store/seed_test.go`で、全migrationを適用した新しいDBのkey `口調`が、採用する全文と一致することを確認できるようにする。

プロンプト合成の流れ、関係、責務は変わらないため、変更図は置かない。

---

## R-2 新しいDBとmigration 18まで適用済みの既存DBへmigration 19を適用し、既存DBの未編集の既定値を更新して、利用者が編集した`instruction`、key `口調`以外の`instruction`、`prompt_template.base_directive`は保持する

### 現況の理解

`db/migrate.go`の`Apply`と`tools/extractor/SchemaMigrator.cs`の`Ensure`は、`PRAGMA user_version`より後のmigrationだけをファイル名順に1回適用する。現在の最新versionは18であり、開発DBもversion 18である。

`internal/store/directive.go`の`SaveDirectiveInstruction`は、画面で編集された`instruction`を同じkeyへ保存する。`directive`は、抽出データと別に永続させる器である。

| | 単位 |
| --- | --- |
| 要求が扱う対象 | 新しいDB、またはversion 18の既存DBにあるkey `口調`の1行 |
| 受け皿が持つkey | `directive.key = '口調'` |

### あるべき形

新しいDBには、採用する`instruction`が入る。

version 18の既存DBでは、key `口調`がmigration 17の既定値と一致する場合だけ採用する`instruction`へ更新する。利用者が編集した値はmigration 17の既定値と一致しないため保持する。

migration 17以前で止まっている既存DBでは、migration 17の既存仕様によって編集内容が上書きされる可能性がある。今回の要求はversion 18の既存DBからの移行に限定する。

### 変更点

`db/migrations/0019_dialogue_paraphrase_directive.sql`の`UPDATE`へ、key `口調`とmigration 17の既定値との一致条件を置く。新しいDBではmigration 17の値から採用する値へ進み、version 18の編集済みDBでは更新件数が0件になる。

`internal/store/seed_test.go`へ、次の3状態をmigrationへ通す検証を追加する。

- 新しいDBへ全migrationを適用すると採用する値になる。
- version 18のDBが未編集の既定値を持つ場合は採用する値になる。
- version 18のDBが編集済みの値を持つ場合は編集済みの値を保持する。

既存のmigrationは変更しない。適用済みmigrationの内容と`PRAGMA user_version`の対応を保つ。

---

## 検討が必要なこと

- なし。
