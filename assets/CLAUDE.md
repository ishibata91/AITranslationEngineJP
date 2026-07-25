# assets/ マップ

実行時に bootstrap（composition root）が os 読みする参照データを置く。純粋な解析は `internal/core/*` が担う。

- `role-speech.tsv` — 注入時に引く一人称・語尾テンプレート。実画面確認で中身を見直すため外部ファイルに置く。解析は `internal/core/rolespeech`。
- `role-speech-examples.tsv` — 口調の例文（英語原文と日本語訳文の 1 対）。`role-speech.tsv` と同じキー（役割区分・性別・基底口調セル）で引き、説明文だけでは揺れる一人称・語尾を実例で固定する。1 行が長くなるのを避けるため役割語表と別ファイルに置く。解析は `internal/core/rolespeech`（ParseRoleSpeechExamples）。
- `stopwords-en.txt` — 機械置換辞書・言及語彙の供給から除く一般語リスト（stopwords-iso 配布、MIT）。上流と byte 一致を保ち、出典・checksum・ライセンス全文は `stopwords-en.LICENSE` に記録する。解析は `internal/core/dictionary`（ParseStoplist）。
- `vader_lexicon.txt` — 口調生成の感情辞書。valence 平均の絶対値が閾値以上の語を強感情語として引く（VADER lexicon 配布、MIT）。上流と byte 一致を保ち、出典・checksum・ライセンス全文は `vader_lexicon.LICENSE` に記録する。解析は `internal/lexicon`（LoadVADER）。
