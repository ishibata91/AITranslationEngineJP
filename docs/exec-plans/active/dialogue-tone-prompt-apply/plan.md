# dialogue-tone-prompt-apply 作業計画

実験 task `dialogue-tone-naturalness` で確定した口調指示を、製品の正本へ入れる。

## branch 情報

- 作業 branch: `claude/dialogue-tone-prompt-apply`
- 統合先 branch: `master`
- 分岐元 commit: `ee6607b0`（実験 branch `claude/dialogue-tone-naturalness` の最後の commit）

分岐元を `master` ではなく実験 branch にした理由を書く。
口調指示は例文・役割語・基底口調の性質文・`directive` の 4 つの組み合わせで決まり、実験 branch はそのうち 3 つを既に変えてある。
`master` から切ると同じ変更を作り直すことになり、実験で測った値と一致する保証が無くなる。

## やることの要点

実験の回 25 の設定を製品へ入れる。回 25 は達成条件 6 行すべてを開発用 598 件と評価用 598 件の両方で満たし、`作り込みの検査` も通過している。

### 1. 中心 DB の `directive` テーブル key `口調` を正本へ入れる

実験中は画面から手で編集していたため、変更が repo に入っていない。これが唯一の未反映の変更である。

- `db/migrations/0017_*.sql` を足し、`key='口調'` の `instruction` を回 25 の 280 文字へ無条件で `UPDATE` する。既存の DB を確実に新しい指示文へ揃えるため。まだ配布していないので、人間が画面で独自に編集した内容を上書きする恐れは無いと人間が 2026-07-28 に判断した。
- `db/migrations/0006_record_type_translation.sql:37` の `INSERT OR IGNORE` の seed も同じ 280 文字へ書き換える。新しく作る DB が最初から新しい指示文を持つようにするため。
- migration は `db/migrations/*.sql` を `embed` して名前順に適用する（`db/migrations.go`）ので、ファイルを置くだけでよい。登録の追記は要らない。
- 同じ SQL を C# の extractor も冪等に ensure する（`db/migrations.go` の冒頭が定める契約）。`UPDATE` 文が C# 側の適用でも問題を起こさないことを確かめる。

回 25 の 280 文字は次のとおり。

```
この台詞の話者の人物像:
{traits}
この人物像に合う口調と人称で訳すこと。
台詞は話し手と聞き手が決まっているので、主語を書かなくても誰の話かは伝わる。英語の I と you を日本語の主語へ置き換えず、原則として書かない。
台詞は常体で書く。話者が女性でも年配でも礼儀正しくても、です・ます体にしない。
原文の内容をすべて訳す。そのうえで英語の語順と品詞をなぞらず、説明を足さずに短く言い切る。
台詞は文末に句点を打たない。原文が疑問符・感嘆符で終わる時だけ ？ ！ を置く。
1 つの台詞に文が 2 つ以上ある時は、句点でつなげず全角空白で区切る。
```

### 2. 実験 branch から引き継いだ変更を確かめる

次の 5 file は分岐元の commit に既に入っている。作り直さず、製品の正本として妥当かを確かめる。

| file | 変更 |
| --- | --- |
| `assets/role-speech-examples.tsv` | 訳例 57 行。主語を書かず、敬体を 4 行に限り、原文の 0.32 倍前後の長さにした |
| `assets/role-speech.tsv` | 一人称の欄を全区分で空にし、言い回しの傾向へ主語を書かない旨を足した |
| `internal/core/personatone/personatone.go` | 対人段階が丁寧の 3 セルの性質文から敬語を求める語を外し、文末の形を具体的に示した。`sexTrait` を短縮した |
| `internal/core/rolespeech/assets_test.go` | 例文が一人称を含むことを要求する検査を、含む場合に役割語と食い違わないことだけ見る形へ変えた |
| `internal/engine/engine_test.go` | 性質文の変更に追従して期待値の文字列を直した |

### 3. 画面の表示を確かめる

`directive` の `口調` は画面から編集できる。280 文字が画面で正しく表示され、編集して保存し直せることを実 app で確かめる。
起動は `npm run dev:wails:run`、接続先は `http://localhost:34115`。

## やらないことの要点

- 口調指示の中身をさらに良くすること。実験は回 25 で止めてある。次に良くするなら別の実験 task を立てる。
- 意味の取り違えを減らすこと。別 task `dialogue-meaning-accuracy` が扱う。
- 固有名が置換されない問題を直すこと。別 task `dictionary-lookup-coverage` が扱う。
- モデルを替えること。`translategemma:12b` を前提にした指示文であり、モデルの選択は本 task の範囲外。
- 実験の砂場（`tmp/dialogue-tone-naturalness/`）を commit へ入れること。
- 実験の作業計画フォルダ（`docs/exec-plans/active/dialogue-tone-naturalness/`）を `completed/` へ移すこと。実験 task 側の出口で人間が決める。
