# Screen

- `*Container.svelte` は画面 state、gateway 呼出、非同期処理、購読の lifecycle を持つ。
- `*Screen.svelte` は表示用 props を描画し、操作を callback prop で返す。
- `*-view.ts` は画面表示に使う型を持つ。
- `*-presentation.ts` は状態から表示値を作る純粋関数と定数を持つ。
- `*.fixtures.ts` は story と表示テストで共有する決定的な入力を持つ。
- `*.stories.ts` は主要状態を列挙し、各状態の前提条件と画面仕様を示す。
- `*-screen-specs.ts` の仕様 ID と `*.spec.test.ts` の検証を一対一で対応させる。
- test は component の外から見える表示、活性、callback を検証する。
- timeout、event subscription、非同期 request は破棄時と入力変更時に古い処理を無効化する。
