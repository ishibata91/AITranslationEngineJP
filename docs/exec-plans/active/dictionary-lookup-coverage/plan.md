# dictionary-lookup-coverage 実験計画

翻訳の前に行う固有名の機械置換で、辞書に訳があるのに置換されない語を拾う実験。

## 事象

人間が `dialogue-tone-naturalness` の回 25 の訳文を読み、固有名が公式訳と違う形で訳される行を見つけた。
原因を辿ると、辞書に訳が登録されているのに機械置換で当たっていない語があった。

実測した 1 件を書く。原文と、置換後にモデルへ渡した文面を並べる。

```
原文        : My husband Eorlund tends the Skyforge up at Jorrvaskr, the Companions mead hall.
モデルへ渡した: My husband エオルンド tends the スカイフォージ up at ジョルバスクル, the Companions mead hall.
回 25 の訳文 : 私の夫、エオルンドは、…あそこは「コンパニオンズ」の酒場にあるんだ。
公式訳      : 夫のエオルンドは同胞団の酒場、ジョルバスクルにある、スカイフォージで働いている。
```

`Eorlund`、`Skyforge`、`Jorrvaskr` は置換された。`the Companions` だけ置換されず、モデルが自力で「コンパニオンズ」と訳した。
`master_term` には `FACT` 区分で `The Companions` → `同胞団` が登録されている。

置換の実装を読んで、当たらない理由が 2 つあると分かった（2026-07-28 に `internal/core/dictionary/dictionary.go` を読んで確認）。

- 照合は大小を区別する。`dictionary.go:59` が `\b(?:<原語を長い順に並べた選択肢>)\b` の正規表現を組み、`regexp.MustCompile` に大小無視の指定を付けていない。登録は `The Companions`（大文字の The）で、台詞は `the Companions`（小文字の the）なので当たらない。
- 照合は登録された文字列の全体に対する完全一致である。`Companions` 単体の登録は無く、`The Companions` としてだけ登録されている。冠詞を除いた形では引けない。

辞書の登録の形を数えた（2026-07-28 の `db/aitranslation.dev.sqlite3`）。

| 数えたもの | 件数 |
| --- | --- |
| 訳が入っている登録（`master_term` と `proper_noun` の合計） | 17,492 |
| うち空白を含む登録（複合語） | 15,050（86.0%） |
| うち `The` で始まる登録 | 353 |
| うち `A` または `An` で始まる登録 | 50 |

複合語が 86% を占める。登録の多くは Skyrim のレコードの FULL 名（`Companions Hunting (Skjor Aela Njada)` のような形）で、台詞の中で全体が完全一致することはまれである。
台詞に出るのはその一部（`Companions`）なので、現在の照合では当たらない。

取りこぼしの件数はまだ測っていない。本 task の準備の段で測る道具を作って測る。

## 対象

### 変える対象

要因を 2 つに絞る。どちらも照合の緩め方であり、緩めると誤爆（一般語への置換）が増えるので、拾える数と誤爆の数の両方を測る。

| 要因 | 変える対象 | 変える中身 |
| --- | --- | --- |
| A | `internal/core/dictionary` の照合 | 大小を区別する（現状）／原語の 2 文字目以降だけ大小を無視する／全体で大小を無視する |
| B | 辞書へ積む語の作り方（`Engine.LoadDictionary`、`internal/engine/engine.go:457`） | 登録をそのまま積む（現状）／先頭の冠詞（`The`、`A`、`An`）を落とした形も積む／複合語の一部も積む |

要因 B の 3 つ目（複合語の一部も積む）は誤爆が最も増えると見込まれる。`Dark Brotherhood Sanctuary` から `Dark` を積むと、一般語の `dark` に当たる余地が出る。
`internal/engine/engine.go:475` の `translationVocabulary` が一般語 stoplist の選別を通しているので、拾う語を増やすときは stoplist との関係も測る。

### 変えない対象

| 変えない対象 | 理由 |
| --- | --- |
| 辞書の中身（`master_term`、`proper_noun` の行） | 本 task は登録済みの語を引けるようにすることを扱う。登録を足すのは別の話で、要因が混ざる |
| 翻訳のプロンプト（例文、役割語、性質文、`directive` の `口調`） | `dialogue-tone-naturalness` が作り込んだ。変えると訳文の質が動き、置換の効果と分けられない |
| モデルと接続先 | `translategemma:12b`、ollama の `http://localhost:11434` に固定する |
| 標本 | `dialogue-tone-naturalness` の凍結標本 1,196 件（開発用 598・評価用 598）をそのまま使う |

## 砂場

`tmp/dictionary-lookup-coverage/` に道具と出力を置く。
`.gitignore` の 14 行目が `tmp` を除外することを、2026-07-28 に `git check-ignore -v` で確かめた。

`dialogue-tone-naturalness` の道具のうち次を流用する。書き換えずに `-data` と `-out` で場所だけ差し替える。

- `tmp/dialogue-tone-naturalness/translate`: 標本を訳して jsonl へ出す。出力の `prompt_user` に置換後の文面が入るので、何が置換されたかを後から読める
- `tmp/dialogue-tone-naturalness/dataset/`: 凍結標本

本 task で新しく要る道具は 2 つある。準備の段で作る。

- 置換の結果を数える道具。標本の原文へ辞書を当て、置換された語と置換されなかった固有名候補を出す。`internal/core/dictionary` をそのまま呼ぶ Go の道具にする。python で正規表現を組み直すと本番の照合と食い違う。
- 誤爆を数える道具。置換された語のうち、一般語へ当たったものを人が読んで判定する形にする。

## 接続先

| 項目 | 値 |
| --- | --- |
| モデル | `translategemma:12b` |
| 接続先 | ollama の OpenAI 互換エンドポイント `http://localhost:11434` |
| 中心 DB | `db/aitranslation.dev.sqlite3` |
| 辞書 | `master_term` 17,346 件、`proper_noun` 13,471 件（訳あり 146 件）。回の間で行を足さない |

## branch 情報

- 作業 branch: `claude/dictionary-lookup-coverage`
- 統合先 branch: `master`
- 分岐元 commit: `ca4b6ba3`

## やらないこと

- 辞書へ登録を足すこと。本 task は登録済みの語を引けるようにすることを扱う。
- 公式訳（`reference_translation` 72,599 件）から新しい対訳を作ること。`DeriveMasterTerms`（`internal/engine/engine.go:501`）が既にその役を持つ。効いていないなら別 task で扱う。
- 訳文の口調（敬体、一人称、二人称、文字数比）を良くすること。`dialogue-tone-naturalness` が達成済み。本 task では悪化していないことだけを診断として見る。
- 意味の取り違えを減らすこと。別 task `dialogue-meaning-accuracy` が扱う。
- 言及検出（`internal/core/mention`）の挙動を変えること。`translationVocabulary` を共有するので、照合を緩めると言及側にも効く。効いた結果は診断として見るが、言及側を目的にした変更はしない。
- プロダクト正本（`docs/architecture.md` など）への反映。
