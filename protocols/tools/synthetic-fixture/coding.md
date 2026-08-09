# 合成 fixture

- 実データへ依存せず、検証対象の record と field だけを含む plugin を作る。
- 一つの検証条件を識別できる EDID と値を与える。
- random、実行時刻、ローカル環境で出力を変えない。
- fixture 変更時だけ手動生成し、生成物と対応する oracle を同時に確認する。
- 抽出規則と期待結果を generator 内へ複製しない。
