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

## 3. Svelte ファイル分割

- `.svelte` は表示、DOM event、受け取った view model の描画に集中させる
- 画面状態の更新手順は `Frontend UseCase`、表示用の整形は `Presenter`、状態保持は `Store` へ分ける
- Wails 呼び出し、backend DTO 変換、generated `wailsjs` の import は gateway ファイルへ閉じ込める
- template が長い条件分岐、複数の表示領域、複数の操作群を持つ場合は、意味単位で screen local component へ分ける
- 画面専用 component は `frontend/src/ui/screens/<screen>/` に置き、複数画面で使う component だけ `frontend/src/ui/components/` に置く

## 4. component 分割基準

- component 化する対象は、独立した意味、明確な props、明確な event、閉じた内部状態、再利用性のいずれかを持つ表示単位にする
- 親画面の状態を大量に直接読む部品、業務フロー全体を進める部品、画面専用の大きなレイアウトは無理に分けない
- 共有 component は UI 規則を集約する目的に限り、画面固有条件を props の分岐として増やし続けない
- component の出力は callback prop または表示結果に限定し、内部で `Store`、`Gateway`、generated binding を直接扱わない
- 詳細な判断は [`architecture.md`](./architecture.md) の `UI Component` 判断表を正本にする

## 5. 更新と副作用

- 状態更新は既存 object の破壊的変更ではなく、新しい値を返す形を優先する
- 長い条件分岐は早期 return や小さい関数へ分け、template に深い判断を置かない
- magic number、遅延時間、件数上限、表示閾値は名前付き定数にする
- `console.log` を本番コードへ残さず、必要な観測は規約化された logger または一時観測として扱う

## 6. Wails 境界

- frontend から backend を呼ぶ入口は generated `wailsjs` を経由する
- generated output を hand-edit しない
- generated `wailsjs` と backend DTO の import は gateway 境界に閉じ込める
- View、ScreenController、Frontend UseCase から generated `wailsjs` を直接参照しない
- user-facing message と internal diagnostic を分ける

## 7. production wiring の配置責務

- production gateway、controller factory、外部 adapter の生成は composition root に置く
- View component は production wiring を作らない
- View component は受け取った controller、usecase、store、props だけを使う
- 既存 fallback wiring がある場合でも、新規機能で View に production dependency 生成を増やす実装は避ける
- responsibility-boundary review では、production wiring が composition root に寄っているかを確認する

## 8. UX 一般規約

- 主要操作はユーザーの作業順に沿って配置し、確認、実行、取消、戻るの意味を混同させない
- ボタンは操作の重要度に応じて primary、secondary、danger を使い分け、破壊的操作は誤操作を避ける配置にする
- ページ見出し、ラベル、説明、エラーメッセージは、ユーザーが次に何をすべきか分かる文言にする
- 読み込み中、空状態、エラー、完了、未保存などの状態を画面上で区別できるようにする
- 同じ画面内では用語、ボタン文言、並び順、余白、入力単位を一貫させる
- design にない機能追加や装飾追加ではなく、承認済み UI 要件を分かりやすく使える形で実装する

## 9. 禁止事項

- generated file を hand-edit する実装
- UI へ内部診断や機密値を無加工で出す実装
- 無検証の `any`、type assertion、暗黙変換へ依存する実装
- View component で production gateway、controller factory、外部 adapter を生成する実装
- Store、Gateway、generated binding を直接扱う shared component
- design にない表示改善や導線変更を実装判断だけで追加する実装

## 10. 参照元

- Svelte official docs:
  [`Svelte 5 Migration Guide`](https://svelte.dev/docs/svelte/v5-migration-guide),
  [`TypeScript`](https://svelte.dev/docs/svelte/typescript)
- Wails official docs:
  [`How Does It Work`](https://wails.io/docs/howdoesitwork)
- 輸入元: `../everything-claude-code/rules/typescript/coding-style.md`
- 輸入元: `../everything-claude-code/rules/common/coding-style.md`
