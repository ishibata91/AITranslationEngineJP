# Mocks

この directory は、過去に作成した page mock を保管する legacy area である。
新規 workflow では実装前の見た目 artifact を UI の必須 artifact または正本にしない。

## Naming

- 既存 page mock は履歴資料として `docs/mocks/<page-id>/index.html` に残す
- 補助 asset が必要な時は同じ `docs/mocks/<page-id>/` 配下へ置く

## Pages

- [`dashboard-and-app-shell/index.html`](./dashboard-and-app-shell/index.html)
- [`master-dictionary/index.html`](./master-dictionary/index.html)
- [`master-persona/index.html`](./master-persona/index.html)

## Notes

- 新規 task は `ui-design.md` に UI 要件契約と実装後確認観点を残す
- 実装前の見た目 artifact の新規作成や正本化は live workflow の標準手順に含めない
- 既存 mock を参照する場合も、実装が満たす契約は `ui-design.md` や page-level requirement に書く
