# internal/core

副作用のない決定規則を、変更理由ごとの package に分けて置く。

- `batchplan/` は batch の送信単位と段遷移を決める。
- `dictionary/` は辞書の照合と置換を行う。
- `mention/` は本文中の辞書語を検出する。
- `prompt/` は AI へ渡す prompt を組み立てる。
- `runtimetag/` は実行時タグを保護する。
- `termderive/` は固有名候補を導出する。
- `termusage/` は固有名の用法を集計する。
- `termxml/` は xTranslator XML を直列化する。
- `rolespeech/` は役割語と例文を照合する。
- `linefeatures/` は台詞の特徴を抽出する。
- `tone/` は基底口調を分類する。
- `personatone/` は口調指示を組み立てる。
