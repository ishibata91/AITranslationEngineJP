# Mocks

この directory は、主要導線と状態変化をある程度再現する page mock の正本を置く。
共通の visual design は `docs/screen-design/` を参照し、この directory では page ごとの挙動つき mock を管理する。

## Naming

- page mock は `docs/mocks/<page-id>/index.html` を正本とする
- 補助 asset が必要な時は同じ `docs/mocks/<page-id>/` 配下へ置く

## Pages

- [`dashboard-and-app-shell/index.html`](./dashboard-and-app-shell/index.html)
- [`master-dictionary/index.html`](./master-dictionary/index.html)
- [`master-persona/index.html`](./master-persona/index.html)

## Notes

- framework 記法や component 名は持ち込まない
- HTML / CSS / 必要最小限の素の JavaScript だけで主要導線と状態変化を再現する
- 実装前の working copy は `docs/exec-plans/active/<task-id>.ui.html` に置き、完了時にこの directory へ移す
