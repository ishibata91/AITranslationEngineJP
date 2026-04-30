# フロントエンド コーディング規約

関連文書: [`coding-guidelines.md`](./coding-guidelines.md), [`architecture.md`](./architecture.md), [`lint-policy.md`](./lint-policy.md)

本書は、`frontend/src/` の TypeScript / Svelte 実装規約を定義する。
Wails bridge、画面状態、表示イベントの責務を対象にする。
backend と test の規約は別文書を正本にする。

## 1. 型と境界

- public API、共有 utility、component props は引数と戻り値の型を明示する
- 外部入力や Wails bridge の戻り値は `unknown` 相当として扱い、使用前に絞り込む
- `any` と無検証の type assertion を常用しない
- 繰り返す object shape は named type または interface へ切り出す
- string literal union を優先し、相互運用が必要な場合だけ enum を使う

## 2. Svelte と状態

- `Svelte 5` と `TypeScript` を前提にする
- component state は `$state`、派生値は `$derived`、副作用は `$effect` を基本にする
- component event は callback prop を優先し、`createEventDispatcher` の新規採用は避ける
- event handler は `onclick` などの標準 event 属性を優先する
- `.svelte` は表示とイベント配線に集中させ、副作用や取得判断を template 内へ散らさない

## 3. 更新と副作用

- 状態更新は既存 object の破壊的変更ではなく、新しい値を返す形を優先する
- 長い条件分岐は早期 return や小さい関数へ分け、template に深い判断を置かない
- magic number、遅延時間、件数上限、表示閾値は名前付き定数にする
- `console.log` を本番コードへ残さず、必要な観測は規約化された logger または一時観測として扱う

## 4. Wails 境界

- frontend から backend を呼ぶ入口は generated `wailsjs` を経由する
- generated output を hand-edit しない
- generated `wailsjs` と backend DTO の import は gateway 境界に閉じ込める
- View、ScreenController、Frontend UseCase から generated `wailsjs` を直接参照しない
- user-facing message と internal diagnostic を分ける

## 5. 禁止事項

- generated file を hand-edit する実装
- UI へ内部診断や機密値を無加工で出す実装
- 無検証の `any`、type assertion、暗黙変換へ依存する実装
- design にない表示改善や導線変更を実装判断だけで追加する実装

## 6. 参照元

- Svelte official docs:
  [`Svelte 5 Migration Guide`](https://svelte.dev/docs/svelte/v5-migration-guide),
  [`TypeScript`](https://svelte.dev/docs/svelte/typescript)
- Wails official docs:
  [`How Does It Work`](https://wails.io/docs/howdoesitwork)
- 輸入元: `../everything-claude-code/rules/typescript/coding-style.md`
- 輸入元: `../everything-claude-code/rules/common/coding-style.md`
