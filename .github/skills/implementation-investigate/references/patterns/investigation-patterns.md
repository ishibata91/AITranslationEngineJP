# Investigation Patterns

## 目的

`investigator` が実装前後の evidence を集めるための判断パターンをまとめる。
agent contract の権限や output obligation は上書きしない。

## 採用する考え方

- entry point から completion まで call chain を追う。
- observed facts、hypotheses、remaining gaps を分ける。
- silent failure、empty catch、dangerous fallback、lost stack trace を重点的に探す。
- temporary observation は目的、path、除去確認を記録する。
- build / runtime error は最小再現と最小観測点から切り分ける。

## 適用ルール

- Wails binding、frontend gateway、backend service、infra adapter のどこで失敗したかを分ける。
- console、backend log、test output、UI state を evidence として区別する。
- paid real AI API を調査で呼ばない。fake / DI seam / test mode を使う。
- 一時観測点は返却前に除去し、cleanup_status を必ず返す。

## 赤旗

- `catch {}`、`.catch(() => [])`、原因を隠す default value がある。
- 再現条件を変えたまま pass と判断している。
- 仮説を fact として implementer へ渡している。
- temporary change が残ったまま completion している。
