# Investigation: proper-noun-case-insensitive-replacement

`investigation.md` は不具合の再現と原因究明だけを持つ。
どう直すかの設計は `design.md` が持つ。

## 観測済み問題

- 人間は、機械置換辞書による固有名置換が大文字小文字を区別するため、原文の表記によって置換結果が安定しないと報告した。
- `internal/core/dictionary/dictionary_test.go` の `TestDictionaryApply/大小を区別し小文字の一般語は置換しない` は、`Storm` だけを置換し、`storm` を置換しない結果を現在の正常動作として検証している。
- `GOCACHE=/tmp/proper-noun-case-insensitive-replacement-go-cache go test ./internal/core/dictionary -run 'TestDictionaryApply$' -count=1 -v` は通過した。
- 開発用 database のコピーで、原語と訳語が空でない `master_term` と `proper_noun` を `NewDictionary` へ入る語彙として合流すると、大文字小文字だけが異なる組は22組あった。
- 機械置換辞書へ入る22組で、訳語が異なる組は0件だった。
- 観測した原語は ASCII であり、集計には SQLite の `lower(source)` を使った。

## 画面再現確認

- Wails 接続対象は `http://localhost:34115` とする。
- 診断用 database は開発用 database のコピーを別 worktree に置き、元の database を変更しない形で使った。
- `proper_noun` に登録済みの `Inigo` から、大小だけを変えた `I asked inigo to wait here.` と対照用の `I asked Inigo to wait here.` を診断用 result として追加した。
- `chrome-devtools` で Wails 画面へ接続し、画面が使う `window.go.api.App.ListResultsPage('inigo.esp', 'n:5862', 2, false)` を実行した。
- `I asked inigo to wait here.` は `terms` が空で、prompt の本文も原文のままだった。
- `I asked Inigo to wait here.` は `terms` に `Inigo` から `イニーゴ` への置換を持ち、prompt の本文が `I asked イニーゴ to wait here.` になった。
- 診断用 result 2 件は確認後に診断用 database から削除した。

## 原因仮説

1. `internal/core/dictionary.NewDictionary` が大小を区別する正規表現を作り、`Dictionary.Apply` が一致した表記をそのまま `bySource` の key に使うため、大文字小文字が異なる同じ固有名を置換できない可能性がある。
2. `internal/core/mention.NewDetector` も大小を区別する正規表現と完全一致の `bySource` を使うため、本文から固有名を見つける処理が原文の表記によって変わる可能性がある。
3. 横断辞書または plugin 内訳語で、大文字小文字だけが異なる同じ固有名に異なる訳語が登録され、原文の表記ごとに別の訳語が選ばれる可能性がある。
4. `internal/core/runtimetag.Mask` による実行時タグの退避後だけ固有名の境界が変わり、正規表現の語境界に一致しない可能性がある。

## 観測ログ検証

- `internal/core/dictionary.Dictionary.Apply` に、診断用原文だけを対象とする一時観測ログを追加した。
- 小文字の原文は `exact_match_count=0`、`case_fold_candidate_count=1` だった。
- 登録表記の原文は `exact_match_count=1`、`case_fold_candidate_count=1` だった。
- 同じ dictionary と同じ句読点を使い、大小だけを変えたため、実行時タグを原因とする仮説4は再現条件に当たらない。
- 機械置換辞書へ入る原語と訳語が空でない語彙の集計では、大文字小文字だけが異なる22組に異なる訳語はなかったため、仮説3は観測した問題の原因ではない。
- `internal/core/mention.Detector` は `Dictionary` と同じ大小区別、語境界、最長一致を不変条件として持つため、仮説2もコード上で成立する。
- 一時観測ログは検証後に product code から削除した。

## 確定原因

- 原因仮説1と原因仮説2を確定原因とする。
- `NewDictionary` が作る正規表現には大小を区別しない指定がないため、登録済みの `Inigo` は `inigo` に一致しない。
- `Dictionary.Apply` は正規表現が返した表記を `bySource` の key として完全一致で取得するため、正規表現だけを大小無視にすると取得も成立しない。
- `mention.NewDetector` と `Detector.Detect` も同じ正規表現と key の構造を持つため、本文から固有名を見つける結果が原文の大文字小文字で変わる。
- 大文字小文字だけが異なる辞書データの訳語は同じだったため、訳語の競合は確定原因に含めない。
- 既存テストは大小区別を正常動作として固定しているため、要求と既存仕様が食い違っている。

## 禁止する修正

- 正規表現だけを大小無視にし、`bySource` の key の扱いを変えない修正を禁止する。大小が異なる一致で確定訳語を取得できないためである。
- `Dictionary` だけを変更し、同じ照合規則を不変条件とする `mention.Detector` を変更しない修正を禁止する。置換した固有名と言及として記録する固有名がずれるためである。
- 語境界または最長一致を外す修正を禁止する。今回の確定原因に関係せず、部分一致の範囲を広げるためである。
- 原文または保存済み辞書を一律に小文字へ書き換える修正を禁止する。表示用の原語と既存データを失うためである。
