# TypeScript / Svelte

- formatter と linter の出力を正とする。
- public API、共有関数、component props は引数と戻り値の型を明示する。
- 外部入力は `unknown` として受け、使用前に型を絞り込む。
- `any` と無検証の type assertion を使わない。
- 繰り返す object shape は named type にする。
- 閉じた値の集合は string literal union で表す。
- Svelte 5 の状態は `$state`、派生値は `$derived`、副作用は `$effect` で表す。
- component の出力は callback prop で返し、DOM event は標準 event 属性で受ける。
- 状態更新は既存値の破壊的変更より、新しい値への置換を優先する。
- 遅延時間、件数上限、表示閾値は名前付き定数にする。
