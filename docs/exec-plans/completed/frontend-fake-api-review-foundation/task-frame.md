# task 枠: フロントエンド fakeAPI レビュー基盤

## 目的

フロントエンドレビュー用の fakeAPI 起動基盤を用意する。
fakeAPI 起動では、Wails バインディング と バックエンド に依存せずに実画面を確認できる状態にする。

## 入力

- `tasks/usecases/frontend-fake-api-review-foundation.yaml`
- [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md)
- [coding-guidelines-frontend.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/coding-guidelines-frontend.md)
- `tmp/code-map/index.json`

## 完了条件

- 起動モードで フロントエンドの API 接続先を fakeAPI に切り替えられる。
- 本番起動では fakeAPI が選ばれない。
- 空状態、読み込み中、成功状態、進行中状態、失敗状態、設定不足状態を fakeAPI で再現できる。
- 画面固有のモックデータを ユースケース task 側で追加できる。
- 実画面を `agent-browser` で開き、状態パターンごとの表示を確認できる。
- fakeAPI と モックデータが本番 API、永続化、本番初期状態に混入しない。
- fakeAPI 起動モードが壊れていないことを局所テストで確認できる。

## 設計前提

- フロントエンドの composition root は `フロントエンド/src/main.ts` である。
- 本番ゲートウェイは Wails バインディング adapter として `フロントエンド/src/controller/wails/` に閉じ込める。
- View、ScreenController、Frontend UseCase は 生成済み `wailsjs` を直接参照しない。
- fakeAPI は provider 選択肢ではなく、レビュー起動時の DI による差し替えとして扱う。

## 非対象

- バックエンド の本番挙動変更。
- 本番初期状態 への モックデータ 注入。
- 生成済み `wailsjs` の手編集。
- docs 正本の直接更新。
