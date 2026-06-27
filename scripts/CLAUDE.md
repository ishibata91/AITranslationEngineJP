# scripts/ マップ

検証と開発起動の実行スクリプトを置く。検証は npm から呼ぶ 1 系統へ一本化し、旧 Python harness（`scripts/harness/`）は退役した。

## 検証コマンド（npm 経由で呼ぶ）

- `npm run verify:backend` — backend 検証 1 コマンド。`go test`（`./` と `./internal/...`）・arch-lint・境界違反走査を束ねて回す。backend を触ったらこれを実行する。
- `npm run test:backend` — backend の `go test` だけを回す。
- `npm run lint:backend` — backend lint 一式（format・vet・static・arch・boundary・module）。
- `npm run lint:backend:arch` — import 方向の検査（`.go-arch-lint.yml`、architecture.md §4）。
- `npm run lint:backend:boundary` — 境界違反走査。arch-lint で表せない責務違反（Wails runtime handle の漏れ、SQLite driver の漏れ）を検出する。
- `npm run test:frontend` / `npm run lint:frontend` — frontend のテストと lint（`frontend/` 配下へ委譲）。
- `npm run dev:wails:run` — 実 app を dev 起動する。実画面確認の既定。起動ごとに中心 DB を空から始める。

純粋ルールの 100% は `go test -cover` の手元確認で見る。常設の coverage ゲートは持たない。

## スクリプト構成（直下1段）

- `dev/run-wails.sh` — 実 app の dev 起動。devserver ポートの既存 process を停止し、中心 DB を消してから `wails dev` を起動する。
- `lint/run-go-backend-lint.sh` — backend lint の入口（`format-check`・`vet`・`static`・`arch`・`boundary`・`module`）。
- `lint/run-boundary-scan.sh` — 境界違反走査の本体。禁止 import を層ごとに固定し、許可層以外での出現を検出する。
- `test/run-go-backend-test.sh` — backend の `go test`。
- `go/run.sh` — go ツール実行の共通ラッパ（GOCACHE 等の環境を整える）。
- `golden/capture.sh` — 実 `.esm` 由来の非劣化 golden を捕獲する local 限定ツール（完了定義③）。
- `node/disable-vite-windows-net-use.cjs` — Windows 出荷判断までの保留物。据え置く。

## 関連

- arch 設定: `.go-arch-lint.yml`
- 静的解析設定: `.golangci.yml`
- 恒久ユニットテスト方針: `docs/coding-guidelines-tests.md`
