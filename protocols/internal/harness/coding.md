# 統合検証 harness

`internal/harness/` は api、engine、store、provider を決定的な入力で接続する test 用 composition root である。

- AI provider、抽出結果、辞書、外部資源を固定した fake へ置き換える。
- product の公開入口から実行し、prompt、保存結果、件数を観測する。
- test ごとに一時 database を作り、実行順と実データへ依存させない。
- golden を変える場合は、仕様変更と期待値の対応を確認する。
