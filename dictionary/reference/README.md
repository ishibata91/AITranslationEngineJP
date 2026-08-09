# 日本語WordNet参照データ

R-3の一般語判定では、日本語WordNet 1.1のSQLite版を参照する。

- 配布元: https://bond-lab.github.io/wnja/eng/downloads.html
- version: 1.1
- 配布日: 2010-10-22
- 取得日: 2026-08-03
- `wnjpn.db` sha256: `a8e749c4a356bf93d0b5de505bca8b21e13746f5728f76819728e8b4c3305a12`
- license: [JAPANESE_WORDNET_LICENSE.txt](./JAPANESE_WORDNET_LICENSE.txt)

保持するSQLiteファイルは`dictionary/reference/wnjpn.db`である。

`wnjpn.db`はMCPの`dictionary_classify` toolが参照する外部データであり、Gitへ追加しない。

整合性は`sqlite3 dictionary/reference/wnjpn.db 'PRAGMA integrity_check;'`で確認する。

分類用の参照データだけを保持し、削除済みのgzip配布物は保持しない。

日本語WordNetは英語と日本語を同じ`synset`へ結び、品詞と意味の定義を持つ。R-3では、一般辞書に同じ原語が存在するかではなく、同じ意味の`synset`に辞書項目の訳語が存在するかを調べるために使う。

配布元は最大5%の項目に誤りがある可能性を明記している。日本語WordNetとの一致は除外候補を作る根拠に限り、自動除外の根拠にはしない。
