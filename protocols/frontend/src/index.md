# frontend

frontend の実装は `frontend/src/` 配下に置く。

- `App.svelte` は画面遷移と画面 container の選択を担当する。
- `main.ts` は frontend の起動点を担当する。
- `ui/` は画面、表示部品、画面状態、Storybook story、fixture を持つ。
- `gateway/` は Wails generated bindings を frontend の型へ変換する。
- `application/diagnostic/` は frontend の観測 logger を持つ。
- `test/` は frontend test の共通 harness を持つ。
