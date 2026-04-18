# Distill Compression Patterns

## 採用する考え方

`distill` は Codex 側の設計前 context だけを圧縮する。
実装前の context 圧縮は [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-distill/SKILL.md) が扱う。

- `context-budget`: 先に棚卸しし、重い入力と重複を見つける
- `strategic-compact`: phase 境界で残す情報を選び、途中の探索ノイズを捨てる
- `rules-distill`: 機械的な収集と LLM 判断を分ける
- `agent-compress`: `catalog`、`summary`、`full` の粒度を使い分ける

## 粒度

- `catalog`: path、名前、役割、状態だけ。最初の棚卸しに使う
- `summary`: 最初の意味のある段落、主要見出し、重要な箇条書きだけ。次工程の判断に使う
- `full`: downstream が直接依存する正本、contract、schema、観測入口だけ。広く読まない

## 判断順

1. task mode と lane owner を固定する
2. 読む候補を catalog 化する
3. 正本、重複、任意参照、未確認を分類する
4. 必要な候補だけ summary / full に上げる
5. contract の field に圧縮する

## 重複除去

同じ制約が複数の場所にある場合は、正本だけを facts / constraints に残す。
重複元は required_reading ではなく、必要なら source note として path だけ残す。

## 品質条件

- 重要な fact は根拠 path を持つ
- 推測は `inferred` として facts から分ける
- 未確認事項は gaps に残し、事実として書かない
- downstream が最初に読む順番が分かる
- product code と product test の変更に進まない
