# Frontend の境界

- `main.ts` と `App.svelte` は起動、画面遷移、production の画面構成だけを持つ。
- Wails generated bindings と backend DTO は `gateway/` の外から import しない。
- backend が返す状態、種別、条件から、表示と操作可否を frontend で導出する。
- user-facing message と内部診断を分け、secret、API key、内部 payload を画面へ出さない。
- frontend の変更後は `npm run test:frontend` と `npm run lint:frontend` を実行する。
