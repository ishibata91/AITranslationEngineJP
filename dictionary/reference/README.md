# 日本語WordNet参照データ

R-3の一般語判定では、日本語WordNet 1.1のSQLite版を参照する。

- 配布元: https://bond-lab.github.io/wnja/eng/downloads.html
- 取得元: https://github.com/bond-lab/wnja/releases/download/v1.1/wnjpn.db.gz
- version: 1.1
- 配布日: 2010-10-22
- 取得日: 2026-08-03
- `wnjpn.db.gz` sha256: `64a14dcfe3ba296566e91a70a2fc0616e85cf2ee7b7fd8cdcbc66c8b12a505a5`
- `wnjpn.db` sha256: `a8e749c4a356bf93d0b5de505bca8b21e13746f5728f76819728e8b4c3305a12`
- license: [JAPANESE_WORDNET_LICENSE.txt](./JAPANESE_WORDNET_LICENSE.txt)

`wnjpn.db`と`wnjpn.db.gz`は検分用の外部データであり、Gitへ追加しない。

日本語WordNetは英語と日本語を同じ`synset`へ結び、品詞と意味の定義を持つ。R-3では、一般辞書に同じ原語が存在するかではなく、同じ意味の`synset`に辞書項目の訳語が存在するかを調べるために使う。

配布元は最大5%の項目に誤りがある可能性を明記している。日本語WordNetとの一致は除外候補を作る根拠に限り、自動除外の根拠にはしない。
