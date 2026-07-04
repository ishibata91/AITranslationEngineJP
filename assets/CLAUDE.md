# assets/ マップ

実行時に bootstrap（composition root）が os 読みする参照データを置く。純粋な解析は `internal/core/*` が担う。

- `role-speech.tsv` — 注入時に引く一人称・語尾テンプレート。実画面確認で中身を見直すため外部ファイルに置く。解析は `internal/core/rolespeech`。
- `stopwords-en.txt` — 機械置換辞書・言及語彙の供給から除く一般語リスト（stopwords-iso 配布、MIT）。上流と byte 一致を保ち、出典・checksum・ライセンス全文は `stopwords-en.LICENSE` に記録する。解析は `internal/core/dictionary`（ParseStoplist）。
