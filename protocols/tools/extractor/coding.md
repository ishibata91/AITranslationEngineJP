# 抽出と保存

`tools/extractor/` は plugin に記録された事実を取得して SQLite へ保存する。

- 抽出結果には record、field、参照、条件、感情などの観測可能な事実だけを入れる。
- 口調、翻訳方針、辞書候補、翻訳要否の最終判断を抽出機に置かない。
- 翻訳対象の英語と既訳の英日対は用途が異なるため、同じ table や model field に混ぜない。
- target plugin の翻訳所有判定は `PluginEnvironment` に集約する。
- master と同一の override は参照用 stub として扱い、翻訳対象の列挙へ重複して出さない。
- SQLite へ結ぶ識別子は plugin 名、form ID、record 種別、field、出現順を保つ。
- writer は schema を適用してから transaction 内で保存し、同じ入力を再実行しても重複を作らない。
- schema、抽出 model、SQLite writer を同時に変える場合は、保存される値と一意性を検証する。

# 境界

- Mutagen による plugin 読込と record の解決は `PluginEnvironment` と `PluginExtractor` に置く。
- `*SqliteWriter` は保存だけを担当し、抽出規則を追加しない。
- SQLite に保存した事実の分類、辞書化、翻訳は Go 側の処理に委ねる。
