# Implementation: replace-extraction-to-prebuilt-dictionary

## 変更したfile

- `internal/store/prebuilt_dictionary.go` に事前作成済み辞書の読み取り専用readerを追加した。
- `internal/model/translation_reference.go` にreader内部候補と本文用参考語を分けて追加した。
- `db/migrations/0022_translation_reference_snapshot.sql` に送信時参考語snapshotを追加した。
- `internal/engine/` を本文置換から参考語promptへ移行した。
- `internal/api/app.go` を結果の再辞書照合からsnapshot参照へ移行した。
- `db/dictionary.sqlite3` へ事前作成済み辞書DBを移動した。

## 仕様との対応

- R-1-1: 本文promptは英語本文とmeaningを含まない参考語を使う。
- R-1-8: 同期本文送信時とbatch結果反映時のsnapshot保存、結果再構成、prompt hash検証を追加した。
- R-1-18: 辞書DBを`db/dictionary.sqlite3`へ移動し、SQLite integrity checkが`ok`を返した。

## 検証結果

- `sqlite3 db/dictionary.sqlite3 'PRAGMA integrity_check;'`: `ok`
- `npm run verify:backend`: 成功
- `go test ./internal/harness`: 成功。本文置換を前提にした統合oracleを参考語注入の観測へ更新した。
- `go test ./internal/engine`: 成功。本文batch送信失敗後の本文段再送を追加で確認した。

## 未確認事項と停止理由

- batch送信前仮状態と失敗再送は実装済みである。
- Wails生成型の再生成、frontend unit test、TypeScript型検査、ESLintは実施済みである。Knipは既存の未使用file 6件を報告する。
- fresh実装レビューは未実行である。

## 人間の指摘

<人間が直接記入する。メインエージェントは記入内容を削除または言い換えない。>
