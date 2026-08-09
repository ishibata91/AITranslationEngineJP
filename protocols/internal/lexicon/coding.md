# 外部辞書 adapter

`internal/lexicon/` は外部辞書ファイルを読み、core が要求する照合契約を実装する。

- file format の解析と照合用データの保持を担当する。
- 口調分類と翻訳手順を置かない。
- 読込失敗は path と処理が分かる error として返す。
