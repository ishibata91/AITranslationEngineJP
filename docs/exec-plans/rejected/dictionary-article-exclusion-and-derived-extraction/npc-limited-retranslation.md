# 人名限定再翻訳: dictionary-article-exclusion-and-derived-extraction

## 入力

`INFO:NAM1` の10語以上の会話文を本文SHA-256の安定順で500件抽出した。

完全一致原語に加え、`NPC_` 出現を持つ語義だけから冠詞除去、二つ名、二語姓名分割の別名を作った。

冠詞除去候補は2件、二つ名候補は3件、二語姓名分割候補は13件だった。

## 実行

OpenAI Batch `batch_6a78732f8cc081909bb49ed75864c5c5` を `gpt-5.6-luna` へ送信した。

500件すべて成功した。

参考語の一致は完全一致240件と二語姓名分割4件だった。

冠詞除去と二つ名派生は500件に一致しなかった。

## 人名派生の観察

`Ulfric Stormcloak → Stormcloak → ストームクローク` は本文の `Stormcloak traitors` で参照された。

`Tiber Septim → Tiber → タイバー` は本文の `Tiber Septim` で完全一致との最長一致により `タイバー・セプティム` と訳された。

`Dremora Lord → Lord → ロード` と `Captain Aquilius → Captain → キャプテン` は一般語を含むため、参考語として提示されたが訳文は人名へ機械固定されなかった。

500件の標本だけでは、冠詞除去と二つ名派生の再翻訳効果を判定できない。
