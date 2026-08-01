# Spec: translation-paraphrase-prompt-default

## 採用するinstruction

key `口調`が持つ`instruction`の全文を次に固定する。

```text
この台詞の話者の人物像:
{traits}
この人物像に合う口調と人称で訳すこと。
台詞は話し手と聞き手が決まっているので、主語を書かなくても誰の話かは伝わる。英語の I と you を日本語の主語へ置き換えず、原則として書かない。
台詞は常体で書く。ただし、丁寧な依頼や断りを命令へ変えない。「～してくれないか」「すまないが～」など、常体のまま原文の丁寧さを保つ。
原文の内容をすべて訳す。日本語で不要な主語と重複だけを省く。依頼、遠慮、推量、程度、話者の態度は短縮しない。
台詞は文末に句点を打たない。原文が疑問符・感嘆符で終わる時だけ ？ ！ を置く。
1 つの台詞に文が 2 つ以上ある時は、句点でつなげず全角空白で区切る。

英語の文法上の形ではなく、台詞全体が伝える内容と話者の意図を日本語で表す。

次の変更は、本文だけから意味と働きを確定できる場合に限って行う。

- 修辞疑問は、答えを求める質問として訳さない。原文が示す断定、皮肉、呆れ、非難を日本語で表す。例: Why am I not surprised? → やっぱり驚かない
- if you'll excuse me、I'm afraid、would you mind などの定型的な丁寧表現は、条件、恐怖、質問として直訳しない。退出の断り、残念の表明、依頼、拒否など、その箇所が果たす働きを日本語で表す。例: if you'll excuse me, I need to leave → すまないが、もう行かなければ。例: would you mind helping me? → 手伝ってくれないか
- 比喩や否定構文に対応する一般的な日本語表現が明確にある場合は、英語と同じ構造を残さず、その意味を直接表す。例: nothing to joke about → 軽く考えていい話ではない。例: the last word I received → 最後の便り
- 原文の複数の節が同じ内容を説明している場合は、日本語で同じ内容を重ねて言わない。ただし、条件、理由、程度、話者の態度は省かない。

次の内容は必ず保つ。

- 誰が行動するか。
- 誰に行動を求めているか。
- 質問、依頼、勧誘、拒否、警告、抗議、断言の区別。
- 丁寧さ、親しさ、敵意、皮肉の強さ。
- 肯定と否定、可能、義務、推量、条件、理由。
- 原文に明示された情報量。

本文だけでは意味を確定できない省略や指示対象は推測しない。対応する日本語表現に確信がない場合は、無理に意訳せず意味を保つ訳を選ぶ。

意訳する場合も、原文にない態度、経緯、指示対象を作らない。肯定と否定、頻度、程度、行為者と対象の関係を変えない。文をつなぎ直す場合も必要な述語と補語を残し、日本語として意味が完結する形にする。
```

---

## R-1 SQLの`directive`でkey `口調`が持つ`instruction`を人間が採用した全文へ更新し、人物像の`{traits}`を埋めた台詞では翻訳AIへ送る指示文に反映し、人物像を作れない台詞と台詞以外の翻訳対象には反映しない

- R-1-1（正常系）: 人物像を作れた台詞では、key `口調`の`{traits}`へ人物像を埋め、採用する`instruction`を翻訳AIへ送る指示文に含めること
    - 前提条件: `INFO:NAM1`、`INFO:RNAM`、`DIAL:FULL`の同期翻訳またはbatch翻訳で人物像を作れた場合
    - 確かめ方: 翻訳AIへ送る指示文に人物像と、上記の採用する`instruction`の残りの文面があることを見る
    - 対応する実テスト: `internal/engine/engine_test.go` の `TestRunTranslatesLinesWithPersonaDirective`、`internal/core/personatone/personatone_test.go` の `TestBuildToneDirective`
- R-1-2（対象に入る側の境界）: key `口調`の`instruction`が`{traits}`を1つだけ持ち、採用する全文と一致すること
    - 前提条件: 全migrationを適用した新しいDBからkey `口調`を読む場合
    - 確かめ方: 読み出した`instruction`と上記の採用する全文を比較する
    - 対応する実テスト: `internal/store/seed_test.go` の `TestFreshDatabaseUsesAcceptedToneDirective`
- R-1-3（対象に入らない側の境界）: 人物像を作れない台詞と台詞以外の翻訳対象へ、key `口調`の`instruction`を送らないこと
    - 前提条件: 人物像を作れない台詞、または台詞以外の翻訳対象を同期翻訳またはbatch翻訳する場合
    - 確かめ方: 翻訳AIへ送る指示文に採用する意訳指示と安全指示がないことを見る
    - 対応する実テスト: `internal/engine/engine_test.go` の `TestRunTranslatesLineWithoutSpeaker` と `TestRunTranslatesUntranslatedAsProvisional`、`internal/engine/batch_integration_test.go` の `TestBatchMatchesSyncEndToEnd`

---

## R-2 新しいDBとmigration 18まで適用済みの既存DBへmigration 19を適用し、既存DBの未編集の既定値を更新して、利用者が編集した`instruction`、key `口調`以外の`instruction`、`prompt_template.base_directive`は保持する

- R-2-1（正常系）: migration 18まで適用済みでkey `口調`がmigration 17の既定値と一致する既存DBを、採用する`instruction`へ更新すること
    - 前提条件: 未編集のkey `口調`を持つversion 18の既存DBへmigration 19を適用する場合
    - 確かめ方: migration適用後のkey `口調`が上記の採用する全文と一致することを見る
    - 対応する実テスト: `internal/store/seed_test.go` の `TestMigration19UpdatesUneditedToneDirective`
- R-2-2（対象に入る側の境界）: migration 18まで適用済みでkey `口調`がmigration 17の既定値と一致しない編集済みの値を持つ既存DBでは、空文字を含む編集済みの`instruction`を保持すること
    - 前提条件: migration 17の既定値と一致しないkey `口調`を持つversion 18の既存DBへmigration 19を適用する場合
    - 確かめ方: migration適用後のkey `口調`が適用前の編集済み文面と一致することを見る
    - 対応する実テスト: `internal/store/seed_test.go` の `TestMigration19PreservesEditedToneDirective`
- R-2-3（対象に入らない側の境界）: key `口調`以外の`instruction`と`prompt_template.base_directive`を変更しないこと
    - 前提条件: 新しいDBまたはversion 18の既存DBへmigration 19を適用する場合
    - 確かめ方: key `口調`以外の`instruction`と`prompt_template.base_directive`がmigration 19の適用前後で一致することを見る
    - 対応する実テスト: `internal/store/seed_test.go` の `TestMigration19LeavesOtherDirectivesAndBaseDirectiveUntouched`

**満たさない部分**: migration 17以前で止まっている既存DBの編集内容を保持する移行は扱わない。
