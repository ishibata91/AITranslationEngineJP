# UI

- `.svelte` は表示、DOM event、props と callback の配線に集中させる。
- Wails binding と backend DTO を直接扱わない。
- 読込中、空、error、完了、未保存を画面上で区別する。
- 状態 label は短い名詞句、button は操作名、説明は利用者が判断できる自然な文章にする。
- error message は失敗内容と次の操作を示し、内部実装名を出さない。
- 色だけで状態と操作を区別しない。
- 承認済みの story と Svelte component を画面表示の正本とする。
