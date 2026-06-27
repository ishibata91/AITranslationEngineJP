# ハーネス整理（棚卸し資料）

関連: [`plan.md`](./plan.md) の動かす範囲 ⑦。

本書は、過去の作り直しで使ってきた検証スクリプト群を棚卸しし、区分を人間が確定するための MECE 資料である。
事実（live 参照元・現状）は Claude が埋める。`人間判断` 列は空欄にしてあり、人間が記入する。
判定できない要素は `未確定` と書く。該当が無い欄は `N/A` と書く。

## 区分の定義

- `いる`: 現在も使う。残す。
- `いらない`: 廃止済みの残骸。参照が無いことを確認して削除する。
- `使ってない`: 存在するがどこからも参照されない。使うか消すかを判定する。
- `新規追加`: 本 task で新たに足す。
- `未確定`: 事実だけでは区分を決められず、人間判断が要る。

## 1. scripts/ 配下ファイル（MECE: 配下の全ファイルを 1 回ずつ）

| path | 種別 | live 参照元 | 現状の事実 | 暫定区分 | 人間判断 |
|---|---|---|---|---|---|
| `scripts/dev/run-wails.sh` | sh | `package.json` `dev:wails:run` | 実 app 起動 command。CLAUDE.md 実画面確認の既定 | いる |  |
| `scripts/lint/run-go-backend-lint.sh` | sh | `package.json` `lint:backend:format/vet/static/arch/module` | backend lint 5 種。**arch-lint を含む**（`arch` 引数） | いる |  |
| `scripts/go/run.sh` | sh | `lint/run-go-backend-lint.sh`・`test/run-go-backend-test.sh`・`test/run-go-backend-coverage.sh` が内部利用 | go ツール実行の共通ラッパ | いる |  |
| `scripts/test/run-go-backend-test.sh` | sh | `package.json` `test:backend` | backend test 経路 | いる |  |
| `scripts/test/run-go-backend-coverage.sh` | sh | `package.json` `test:backend:coverage` | backend coverage 経路 | いる |  |
| `scripts/test/run-frontend-coverage.sh` | sh | `package.json` `test:frontend:coverage` | frontend coverage 経路 | いる | いらない |
| `scripts/node/run-sonar-scanner.mjs` | mjs | `package.json` `scan:sonar` | Sonar scanner 起動 | 未確定（Sonar 運用継続するか） | いらない |
| `scripts/harness/run.py` | py | npm 未参照。`docs/coding-guidelines-tests.md` が古い suite 名で言及 | 旧 Python harness。suite は frontend-lint/frontend-local/frontend-test/structure/all のみ。**backend 無し** | 未確定（退役するか正本に戻すか） | メンテする |
| `scripts/harness/check_structure.py` | py | `run.py` `structure`/`all` | 構造検査 | 未確定（run.py の去就に従属） | いらない |
| `scripts/harness/check_frontend_lint.py` | py | `run.py` `frontend-lint`/`frontend-local`/`all` | frontend lint。`npm frontendlint` と二重の疑い | 未確定（run.py の去就に従属） | いらない |
| `scripts/harness/check_frontend_test.py` | py | `run.py` `frontend-test`/`frontend-local`/`all` | frontend test。`npm test:frontend` と二重の疑い | 未確定（run.py の去就に従属） | いらない |
| `scripts/harness/harness_common.py` | py | `run.py`・各 `check_*.py` | harness 共通処理 | 未確定（run.py の去就に従属） | いらない |
| `scripts/harness/README.md` | md | `scripts/harness/*` | harness 説明文書 | 未確定（run.py の去就に従属） | claude.mdにかえる |
| `scripts/dict/derive-master-terms/main.go` | go | live 参照なし（`docs/changelog.md` のみ） | 固有名派生 CLI。`engine.DeriveMasterTerms` と機能重複の疑い | 未確定（重複か独立用途か要確認） | いらない |
| `scripts/dict/derive-master-terms/main_test.go` | go | 同上（`main.go` のテスト） | 上記 CLI のユニットテスト | 未確定（main.go の去就に従属） | いらない |
| `scripts/lint/run-oxlint.mjs` | mjs | live 参照なし | oxlint 起動。frontend lint は別経路で稼働 | 使ってない | いらない |
| `scripts/node/disable-vite-windows-net-use.cjs` | cjs | live 参照なし | Windows 専用の vite net-use 回避 | 未確定（Windows 出荷で要るか） | いる？ |

## 2. npm scripts（package.json）の検証 wiring（参考・事実）

| npm script | 実体 | 区分 |
|---|---|---|
| `dev:wails:run` | `scripts/dev/run-wails.sh` | いる |
| `lint:backend`（format/vet/static/arch/module） | `scripts/lint/run-go-backend-lint.sh` | いる |
| `test:backend` | `scripts/test/run-go-backend-test.sh` | いる |
| `test:backend:coverage` | `scripts/test/run-go-backend-coverage.sh` | user:いらない |
| `lint:frontend` / `frontendlint` | `npm --prefix frontend run lint` | いる |
| `test:frontend` | `npm --prefix frontend run test` | いる |
| `test:frontend:coverage` | `scripts/test/run-frontend-coverage.sh` | いらないuser |
| `scan:sonar` | `scripts/node/run-sonar-scanner.mjs` | 未確定 |

## 3. scripts/harness/run.py の suite（事実）

| suite | 呼ぶ script | 実体の有無 | 区分 |
|---|---|---|---|
| `frontend-lint` | `check_frontend_lint.py` | あり | npmに固める(user) |
| `frontend-local` | `check_frontend_lint.py`＋`check_frontend_test.py` | あり | 未確定（同上） |
| `frontend-test` | `check_frontend_test.py` | あり | 未確定（同上） |
| `structure` | `check_structure.py` | あり | 未確定（同上） |
| `all` | structure＋frontend 2 種 | あり | 未確定（同上） |
| `backend-local` | — | **無し（docs のみ参照）** | いらない（dangling） |
| `coverage` | — | **無し（legacy plan のみ参照）** | いらない（dangling） |
| `system-test` | — | **無し（legacy plan のみ参照）** | いらない（dangling） |

## 4. docs の dangling 参照（事実）

| 参照元 | 参照先 | 現状 | 区分 |
|---|---|---|---|
| `docs/coding-guidelines-tests.md` §6 | `run.py --suite backend-local` | suite 実体なし | いらない（記述を直す） |
| `docs/coding-guidelines-tests.md` §6 | `run.py --suite frontend-local` | suite 実体あり | いる |

## 5. .github/.trash（事実）

| 対象 | 現状 | 区分 |
|---|---|---|
| `.github/.trash/2026-04-18-*`（旧 multi-agent github 一式） | `.trash` 配下で退役済み。`.github/workflows/` に live workflow 無し | いらない（削除可否は人間判断）消す |

## 6. 横断論点（人間判断が要る所）

検証が 2 系統に割れている。backend・arch-lint は npm 系で稼働、Python harness（run.py）は frontend のみで宙吊り。
どちらへ寄せるかを決めると、§1・§3 の `未確定` の多くが連動して決まる。

| 論点 | 選択肢 | Claude 私見 | 人間判断 |
|---|---|---|---|
| backend 検証をどの系統へ寄せるか | 案ア: npm へ一本化、run.py 退役 / 案イ: run.py を正本に戻し backend suite 追加、npm は委譲 | 案ア（backend・arch-lint が既に npm で稼働、run.py は二重実装で drift 源） | npm一本 |
| Sonar 運用を続けるか | 続ける / やめる | 未確定（運用実態を要確認） | やめる |
| `derive-master-terms` CLI の去就 | 残す / 消す（`engine.DeriveMasterTerms` へ一本化） | 未確定（重複か独立用途か要確認） | 消す |
| Windows 専用 `disable-vite-windows-net-use.cjs` | 残す / 消す | 未確定（Windows 出荷方針に依存） | わからん |
| `.github/.trash` を物理削除するか | 削除 / 保持 | 削除可（退役済み） | 消す |

## 7. 新規追加（本 task で足すもの）

| 要素 | 内容 | 区分 |
|---|---|---|
| backend 検証 1 コマンド | `go test` ＋ arch-lint ＋ 境界違反走査を束ねて 1 コマンドで回す（寄せ先は §6 の判断に従う） | 新規追加 |
| 境界違反走査 | arch-lint で表せない責務違反（runtime handle 漏れ、禁止 import 等）の走査 | 新規追加 |

## 8. 未確定・N/A 集約

- 未確定（保留）: `disable-vite-windows-net-use.cjs`（Windows 出荷判断まで据え置き）。
- N/A: 現時点で該当なし。新たに判定不能・対象外の要素が出たらここへ追記する。

## 9. 確定仕分け（人間判断を反映した正本）

§1〜§6 の人間判断と、その後の事実確認（`check_structure.py` は docs 索引検査でコード構造でない／Sonar 全廃で関連設定も死ぬ）を踏まえた最終仕分け。矛盾していた §1「run.py メンテする」は §6「npm 一本」に従い退役へ統一する。

### 9.1 残す（いる）

| 対象 | 理由 |
|---|---|
| `scripts/dev/run-wails.sh` | 実 app 起動。実画面確認の既定 |
| `scripts/lint/run-go-backend-lint.sh` | backend lint 5 種。arch-lint を含む |
| `scripts/go/run.sh` | go ツール実行の共通ラッパ |
| `scripts/test/run-go-backend-test.sh` | backend test |
| npm: `dev:wails:run`・`lint:backend*`・`test:backend`・`lint:frontend`/`frontendlint`・`test:frontend` | 検証の正本系統 |

### 9.2 消す（いらない）

| 対象 | 連動・備考 |
|---|---|
| `scripts/harness/` 一式（`run.py`・`check_frontend_lint.py`・`check_frontend_test.py`・`check_structure.py`・`harness_common.py`・`README.md`） | npm 一本化で退役。frontend 検査は npm と二重。docs 索引検査（`check_structure.py`）は「docs 整合をやめる」決定で廃止 |
| `scripts/test/run-go-backend-coverage.sh` ＋ npm `test:backend:coverage` | 常設 coverage 廃止。100% は `go test -cover` 手元確認へ |
| `scripts/test/run-frontend-coverage.sh` ＋ npm `test:frontend:coverage` | 同上 |
| `scripts/node/run-sonar-scanner.mjs` ＋ npm `scan:sonar` | Sonar 全廃 |
| `sonar-project.properties`（tracked） | Sonar 全廃で死ぬ設定。削除 |
| `scripts/lint/run-oxlint.mjs` | live 参照なし |
| `scripts/dict/derive-master-terms/`（`main.go`・`main_test.go`） | `engine.DeriveMasterTerms` へ一本化 |
| `.github/.trash/` 一式 | 退役済み multi-agent github |
| `docs/coding-guidelines-tests.md` §6 の `run.py --suite` 参照 | 実体消滅。npm 検証へ記述を書き換え |

### 9.3 保留（未確定）

| 対象 | 起動条件 |
|---|---|
| `scripts/node/disable-vite-windows-net-use.cjs` | Windows 出荷を判断する時に要否再評価。それまで据え置き |

### 9.4 新規追加

| 対象 | 内容 |
|---|---|
| backend 検証 1 コマンド | npm script で `go test` ＋ arch-lint ＋ 境界違反走査を束ねる |
| 境界違反走査 | arch-lint で表せない責務違反（runtime handle 漏れ、禁止 import 等）の走査 |
| `scripts/CLAUDE.md`（scoped 指示） | 退役する `harness/README.md` の代わりに、npm 検証手順を scoped 指示として置く（`dictionaries/CLAUDE.md` と同じ scoped CLAUDE.md パターン） |
