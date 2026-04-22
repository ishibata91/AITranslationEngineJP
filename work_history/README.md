# work_history

## 目的

`work_history/` は、1ランごとの問題点と改善点を残す場所です。
Codex と Copilot の報告を同じラン単位で並べ、次回の設計、handoff、実装、検証を改善します。

記録では、事実、時間配分、詰まり、無駄、改善案を優先します。
長い経緯説明や感想は避け、次のランで使える判断材料に絞ります。

## 配置

- 実レポートの唯一の置き場所は `work_history/runs/YYYY-MM-DD-<task-id>-run/` とする。
- 複製元は `work_history/templates/run/` に置く。
- 1ランの folder には `README.md`、`codex.md`、`copilot.md` を置く。
- `README.md` は全体 index、`codex.md` と `copilot.md` は役割別報告にする。

## 配置判断

- 既存の同一 run folder がある場合は、そこへ追記または更新する。
- run folder がない場合は、`work_history/templates/run/` を複製して作る。
- Codex の報告は `codex.md`、Copilot の報告は `copilot.md` だけに書く。
- 両者の比較、重複、遅延、次回改善は run folder の `README.md` に集約する。
- `docs/exec-plans/`、`.codex/history/`、handoff file には run report を置かない。

## 命名

- 日付はラン開始日を `YYYY-MM-DD` で書く。
- `<task-id>` は exec plan、issue、handoff の名前に合わせる。
- 同日に同じ task を複数回走らせる場合は末尾に `-2` などを足す。
- 例: `work_history/runs/2026-04-22-master-persona-run/`

## 書き方

- `改善すべきこと`、`時間がかかったこと`、`無駄だったこと`、`困ったこと` は必ず書く。
- 追加で、曖昧だった前提、reroute 原因、検証不足、次回の prompt / handoff 改善を書く。
- 分からない項目は空欄にせず、`なし`、`未確認`、`不明` のどれかで明示する。
- 各 template の末尾にある `SUMMARY` は消さず、短く更新する。

## 運用

- ラン終了直後に、Codex と Copilot の両方の報告を埋める。
- Codex / Copilot のオーケストレーターは、最後に必ず該当 lane のレポートを作る、または作らせる。
- 片方だけ実行した場合も、未実行側には `未実行` と書く。
- 比較はラン folder の `README.md` に集約する。
- product code、product test、docs 正本、workflow contract の代わりには使わない。

## SUMMARY

- `変更ファイル`: `work_history/README.md`
- `重要エラー`: なし
- `次に見るべき場所`: `work_history/templates/run/`
- `再実行コマンド`: `find work_history -maxdepth 4 -type f | sort`
