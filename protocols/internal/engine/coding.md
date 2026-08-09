# 翻訳手続き

`internal/engine/` は store、provider、filesystem と純粋規則を束ねる。

- 抽出結果の取込、辞書解決、固有名翻訳、本文翻訳、結果保存、XML 出力の順序を持つ。
- 決定的な変換と判定は `internal/core/` へ置く。
- SQL、Wails runtime、provider 固有 request を直接扱わない。
- 外部処理の途中失敗では、完了していない状態を成功として保存しない。
- engine が組み立てた完成 prompt だけを provider へ渡す。
