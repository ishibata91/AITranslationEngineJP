# Lint 方針

関連文書: [`index.md`](./index.md), [`tech-selection.md`](./tech-selection.md), [`architecture.md`](./architecture.md)

この文書は、lint と静的チェックが何を管理し、何を管理しないかをまとめるための正本とする。
個別ツールの採用自体は `tech-selection.md`、検証入口と失敗条件は対応する tests / acceptance checks / validation commands を正本とし、本書は責務範囲を一覧化する。

## 1. Lint と静的チェックが管理するもの

- import 境界: UI が generated binding や backend 内部構造へ直接依存しないこと
- import の健全性: 未使用 import、危険な import パターン、雑な path alias の増殖
- 型と構文の破綻: `TypeScript` / `Svelte` の静的に検出できる問題
- backend の静的品質: format 逸脱、未使用要素、危険な記述、境界違反、禁止 import、非推奨依存、test と production の依存混入
- format の逸脱: repo が採用した formatter の出力との不一致

## 2. Lint と静的チェックが管理しないもの

- 実行時の正しさや振る舞い: `go test`、`Vitest`、acceptance checks で担保する
- UI の表示仕様や操作結果: component test、screen-level test、end-to-end check で担保する
- 業務フローや受け入れ条件の成立: 対応する tests / acceptance checks / validation commands で担保する
- 仕様変更の正当性: human review と正本 docs が担保する

## 3. ツールの担当範囲

- `ESLint`: TypeScript / Svelte の通常 lint を担当し、未使用 import、到達不能コード、無効化コメントの取り残し、コメントアウトされたコード、frontend の依存境界違反を検出する
- `TypeScript compiler`: `tsc --noEmit` により、未使用 local / parameter と到達不能コードを含む型レベルの静的診断を担当する
- `knip`: 未使用 export、未使用 file、未使用 dependency を担当し、type export も検出対象に含める
- `svelte-check`: Svelte component と TypeScript 境界の診断を担当する
- backend format check: backend package の formatter 出力との不一致を担当する
- backend static check: backend package の未使用、危険な記述、error 処理不備、非推奨 API、禁止 import を担当する
- backend architecture / dependency check: backend package の依存方向、横方向依存、test と production の依存分離、禁止 module / local replace の混入を担当する
- `go test`: backend の executable spec を担当する
- `Vitest`: frontend の executable spec を担当する

## 4. 初期設定の配置

- frontend の lint / test 設定は `frontend/` を基準に置く
- repo root には frontend lint の共通入口だけを置き、frontend 実体の lint 定義は `frontend/` 側へ寄せる
- Wails project の設定は repo root の `wails.json` を正本にする
- Go module と backend の検証は repo root の `go.mod` を基準に置く
- generated `wailsjs` は lint の主対象から外してよいが、gateway 以外からの import は許可しない

## 5. 初期ゲートの分割

- frontend のゲート:
  - repo root の `lint:frontend` または `frontendlint`
  - `npm run lint`
  - `npm run check`
  - `npm run test`
  - `npm run build`
- backend のゲート:
  - backend package に対する format 確認
  - backend package に対する静的品質確認
  - backend package に対する architecture / dependency 確認
  - backend package に対する executable spec
- desktop packaging:
  - `wails build`

初期ゲートの責務は `ESLint`、`svelte-check`、backend format check、backend static check、backend architecture / dependency check、`Vitest`、`go test` に固定する。
追加ツールは、繰り返し同種の失敗が出た時だけ導入を検討する。

### 5.2 backend lint の内訳

- backend lint は、format、静的品質、architecture / dependency の 3 系統を分けて実行する
- format では backend package の整形逸脱だけを検出する
- 静的品質では backend package の未使用、危険な記述、error 処理不備、非推奨 API を検出する
- architecture / dependency では backend package の依存方向、横方向依存、禁止 import、非推奨依存、test と production の依存混入を検出する
- backend lint の主対象は backend package に限定し、frontend dependency directory や generated directory を巻き込まない

### 5.1 frontend lint の内訳

- frontend lint は、構文 lint、型検査、未使用検出を分けて実行する
- 構文 lint では `ESLint` を使い、`TypeScript` / `Svelte` の通常 lint と frontend 固有の境界制約を検証する
- 型検査では `tsc --noEmit` を使い、emit なしで未使用 local / parameter と到達不能コードを検証する
- 未使用検出では `knip` を使い、value export だけでなく type export も含めて検証する
- generated directory、build artifact、dependency directory は lint の主対象から除外してよい

## 6. クリーンアップ方針

- 未使用 file / export / dependency は、検出した変更の中で削除する
- 削除しない場合は、なぜ保持する必要があるかを plan か review note に残す
- generated code の除外は許可するが、hand-written code の恒久逃げ道にしない
- lint の都合で architecture を曲げず、必要なら architecture と tests を先に見直す

## 7. 記録の分担

- 採用ツールと品質基盤の決定: [`tech-selection.md`](./tech-selection.md)
- import 境界の設計原則: [`architecture.md`](./architecture.md)
- 実装規約: [`coding-guidelines.md`](./coding-guidelines.md)
- 検証 command と失敗条件: 対応する tests / acceptance checks / validation commands
- task-local な一時判断: [`exec-plans/`](./exec-plans/active/README.md)
