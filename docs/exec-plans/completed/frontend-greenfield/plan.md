# plan: frontend-greenfield

## 依頼要約

backend greenfield reset の続き。frontend 側の domain 汚染を除去し、視覚資産（Svelte screen、storybook）は残す。同時に `.claude/skills/implement/` を `coding-protocol/` に rename し、`implementation-module` との 2 段構造の命名衝突を解消する。

## 作業 branch

- 作業 branch: `claude/frontend-greenfield`
- 分岐元 branch: `master`
- 分岐元 commit: `956fa1af`

## 決定方針（user 回答 2026-06-06）

| # | 論点 | 確定 |
|---|---|---|
| 1 | frontend 削除範囲 | 案 W: Svelte UI と storybook 残す、配線層だけ削除 |
| 2 | skill rename | `.claude/skills/implement/` → `.claude/skills/coding-protocol/` |

## 残す（決定）

| 区分 | 対象 |
|---|---|
| frontend code | `frontend/src/ui/`（screens、components、views、stores の visual 資産） |
| frontend code | `frontend/src/test/`（generic test infra、domain 汚染なし） |
| frontend code | `frontend/src/application/diagnostic/`（generic logger、domain 汚染なし） |
| storybook | `frontend/.storybook/` などの設定 |
| docs | `docs/screen-design/` |
| skill | `.claude/skills/coding-protocol/`（rename 後） |

## 削除（決定）

| 区分 | 対象 |
|---|---|
| frontend code | `frontend/src/controller/wails/`（gateway 全部） |
| frontend code | `frontend/src/controller/<domain>/` 9 domain |
| frontend code | `frontend/src/application/contract/<domain>/` |
| frontend code | `frontend/src/application/usecase/`、`store/`、`presenter/`、`gateway-contract/` 配下の domain key |
| frontend code | `frontend/src/bootstrap/`（domain 汚染、composition root 全体） |
| frontend code | `frontend/wailsjs/`（generated bindings、domain 汚染 snapshot） |
| frontend code | `frontend/src/main.ts`（bootstrap を呼ぶ entry、新 architecture で再構築） |

## 確認後に判断する file

- `frontend/src/ui/screens/<domain>/` 内 svelte が import 経由で domain 配線層に依存している場合、import を残したまま該当層を消す。**TypeScript / lint / test が fail するのは想定通り**。次 task の再配線で復旧する。

## 進め方

1. branch + plan 初期化（本 file）
2. `.claude/skills/implement/` → `.claude/skills/coding-protocol/` rename
3. rename に伴う参照更新（implementation-module、CLAUDE.md、他 skill）
4. frontend 削除対象の特定（ui/screens 配下の domain ディレクトリは visual 資産として残す方針を確認）
5. frontend 削除実行
6. `docs/architecture.md` の現状記述を更新（frontend も削減済みと記録）
7. `docs/index.md` の参照を更新（壊れた link を除去）
8. 検証
    - `structure` suite で docs index 整合のみ確認（frontend-local は fail 想定、論点 1 の代償として記録）
9. work commit
10. local merge to master（`--no-ff`）
11. active → completed 移動
12. merge 結果 commit

## 検証方針

- `structure` suite は通す
- `frontend-local` は domain 配線層削除に伴い fail する。**次 task の再配線まで一時的に fail を受容**する方針を本 task で明文化する。代替として:
    - 削除前後で `frontend-local` 結果を `plan.md` に記録
    - 削除後の lint / typescript error 件数を 1 行で残す（次 task の入力）

## 実施内容（2026-06-06）

### Phase 1: branch + plan 初期化

- branch `claude/frontend-greenfield` を `956fa1af` から作成
- active plan folder 作成

### Phase 2: skill rename

- `.claude/skills/implement/` → `.claude/skills/coding-protocol/`（`git mv`）
- SKILL.md frontmatter `name:` 更新
- 参照更新: `.claude/skills/implementation-module/SKILL.md`（5 箇所）
- `coding-protocol/SKILL.md` を骨格化（backend 削除済 path、deleted docs 参照、domain ルールを除去。新 architecture 確定後に再構築する旨を明記）

### Phase 3: frontend 配線層削除

削除:

- `frontend/src/controller/wails/`（gateway 全部）
- `frontend/src/controller/<domain>/` 9 dir（screen-controller）
- `frontend/src/controller/runtime/`（runtime adapter）
- `frontend/src/application/contract/`、`usecase/`、`store/`、`presenter/`、`gateway-contract/`
- `frontend/src/bootstrap/`、`frontend/src/main.ts`
- `frontend/wailsjs/`（gitignored generated 含む）
- `scripts/eslint/repository-boundary-plugin.mjs`、`frontend/repository-boundary-plugin.test.mjs`
- `docs/diagrams/`（frontend-architecture、components/frontend）

残存:

- `frontend/src/ui/`（screens、components、views、stores、storybook stories と fixtures、test）
- `frontend/src/application/diagnostic/`（generic logger）
- `frontend/src/application/README.md`（greenfield 残存メモへ更新）
- `frontend/src/test/setup.ts`

### Phase 4: 設定ファイル整理

- `frontend/tsconfig.json` から `wailsjs/**/*.ts` 削除
- `frontend/eslint.config.js` から `wailsjs/**` ignore と `repository-boundary-plugin` import / `local` plugin / `local/no-commented-out-code` / `local/enforce-layer-boundaries` を削除
- `frontend/index.html` から `main.ts` script tag 削除
- `frontend/package.json` の `ignoreIssues` を空に

### Phase 5: docs 更新

- `docs/architecture.md` の Wails 境界・ディレクトリ正本・現在の状態 section を更新（frontend 削減も反映）
- `docs/index.md` から削除済み diagrams への link を除去

### 検証結果（2026-06-06）

- `python3 scripts/harness/run.py --suite structure`: PASS
- `python3 scripts/harness/run.py --suite frontend-local`: FAIL（**想定通り、論点 1 代償**）
    - lint:frontend エラー件数: 1617（domain 配線層削除に伴う未解決 import 連鎖）
    - test:frontend: 5 file fail / 5 file pass、27 個別 test pass
    - 次 task で新 architecture に合わせて配線復活すれば解消する想定
