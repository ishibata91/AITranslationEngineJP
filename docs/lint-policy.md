# Lint 方針

関連文書: [`index.md`](./index.md), [`tech-selection.md`](./tech-selection.md), [`architecture.md`](./architecture.md)

この文書は、lint と static checks が何を管理し、何を管理しないかをまとめるための正本とする。
個別ツールの採用自体は `tech-selection.md`、検証入口と失敗条件は対応する tests / acceptance checks / validation commands を正本とし、本書は責務範囲を一覧化する。

## 1. Lint と static checks が管理するもの

- import 境界: UI が generated binding や backend 内部構造へ直接依存しないこと
- import hygiene: 未使用 import、危険な import パターン、雑な path alias の増殖
- 型と構文の破綻: `TypeScript` / `Svelte` の静的に検出できる問題
- Go の基本静的品質: compile error、未使用要素、`go vet` で検出できる明確な問題
- format の逸脱: repo が採用した formatter 出力との不一致

## 2. Lint と static checks が管理しないもの

- runtime の正しさや振る舞い: `go test`、`Vitest`、acceptance checks で担保する
- UI の表示仕様や操作結果: component test、screen-level test、end-to-end check で担保する
- 業務フローや受け入れ条件の成立: 対応する tests / acceptance checks / validation commands で担保する
- 仕様変更の正当性: human review と source-of-truth docs が担保する

## 3. Tool Ownership

- `ESLint`: TypeScript / Svelte の通常 lint を担当する
- `svelte-check`: Svelte component と TypeScript 境界の診断を担当する
- `gofmt`: Go の formatter 出力を担当する
- `go vet`: Go の基本 static check を担当する
- `go test`: backend の executable spec を担当する
- `Vitest`: frontend の executable spec を担当する

## 4. 初期 config 配置

- frontend lint / test config は `frontend/` を基準に置く
- Wails project config は repo root の `wails.json` を正本にする
- Go module と backend validation は repo root の `go.mod` を基準に置く
- generated `wailsjs` は lint の主対象から外してよいが、gateway 以外からの import は許可しない

## 5. 初期 gate split

- frontend gate:
  - `npm run lint`
  - `npm run check`
  - `npm run test`
  - `npm run build`
- backend gate:
  - `gofmt` による format check
  - `go vet ./...`
  - `go test ./...`
- desktop packaging:
  - `wails build`

初期 gate の責務は `ESLint`、`svelte-check`、`gofmt`、`go vet`、`Vitest`、`go test` に固定する。
追加ツールは、繰り返し同種の失敗が出た時だけ導入を検討する。

## 6. Cleanup Policy

- 未使用 file / export / dependency は、検出した変更の中で削除する
- 削除しない場合は、なぜ保持する必要があるかを plan か review note に残す
- generated code の ignore は許可するが、hand-written code の恒久逃げ道にしない
- lint の都合で architecture を曲げず、必要なら architecture と tests を先に見直す

## 7. Record Split

- 採用ツールと品質基盤の決定: [`tech-selection.md`](./tech-selection.md)
- import 境界の設計原則: [`architecture.md`](./architecture.md)
- 実装規約: [`coding-guidelines.md`](./coding-guidelines.md)
- validation command と失敗条件: 対応する tests / acceptance checks / validation commands
- task-local な一時判断: [`exec-plans/`](./exec-plans/active/README.md)
