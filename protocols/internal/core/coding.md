# Functional core

- 入力から出力を決める純粋関数として実装する。
- filesystem、database、network、environment、clock、random を参照しない。
- `engine/`、`store/`、`lexicon/` を import しない。
- package は一つの決定規則と、規則を証明する単体テストを持つ。
- 同じ語彙へ異なる処理を適用する規則は、照合順、大小文字、語境界、重複時の優先順位を一致させる。
- 境界値、順序、重複、空入力、失敗を deterministic test で固定する。
