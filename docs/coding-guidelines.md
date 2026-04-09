# コーディング規約

関連文書: [`index.md`](./index.md), [`architecture.md`](./architecture.md), [`tech-selection.md`](./tech-selection.md), [`lint-policy.md`](./lint-policy.md)

本書は、`AITranslationEngineJp` の実装時に守るべき横断的なコーディング規約を定義する。
本書では構造や directory 構成の説明は扱わず、実装時に常に効く判断基準だけを記録する。
構造責務の正本は [`architecture.md`](./architecture.md) とする。

## 1. 基本原則

- 実装は採用技術と正本仕様に合わせ、外部テンプレートや一般論を無検証で持ち込まない
- コメントは `何をしているか` ではなく、`なぜその判断が必要か` を短く補足する
- 命名は略語より意味を優先し、役割と責務が読める名前にする

## 2. 命名と記述

- `Go`、`TypeScript`、`Svelte` それぞれの標準的な命名慣習を崩さない
- public に見える名前は省略しすぎず、利用側が文脈なしで意味を追える粒度にする
- boolean は状態が読める名前にし、否定の二重表現を避ける
- 一時変数は寿命を短く保ち、意味のない短縮名を使い回さない
- 1 つの関数や component に複数の責務を詰め込まず、読む側が追える単位へ分ける

## 3. エラー処理

- 通常フローの失敗は戻り値や明示的な失敗表現で扱い、`panic` を制御フローに使わない
- 失敗を握りつぶさず、必要な文脈を付けて上位へ返す
- validation error、外部依存の実行失敗、想定外障害を同じ重さで混ぜず、原因の層が分かる形で扱う
- user-facing message と internal diagnostic は分け、内部詳細を UI や外部境界へそのまま出さない
- 再試行、フォールバック、握りつぶしを入れる場合は、理由と境界をコード上で読める状態にする
- cleanup が必要な処理では、途中失敗時にも後始末が漏れないようにする

## 4. 入出力と validation

- request、response、設定値、外部入力は境界で形を固定し、optional field の意味を曖昧にしない
- frontend で整形や制約をかけても、最終 validation は backend で再実行する
- file path、外部 URL、provider 設定値、外部プロセス入力は使用直前に再検証する
- 無検証の type assertion、`any`、暗黙変換に依存しすぎず、失敗可能性を型または validation で表現する
- 文字列連結や ad-hoc な map 構築で契約を表現せず、契約があるデータは DTO や struct に寄せる

## 5. 技術別ルール

### 5.1 Wails

- `Bind` する public method は transport boundary として扱い、重い業務処理や永続化詳細を直書きしない
- frontend から backend を呼ぶ入口は generated `wailsjs` を経由し、generated output を hand-edit しない
- backend から frontend への通知は `runtime.EventsEmit` を使ってよいが、通常の query / command を event へ逃がさない
- lifecycle hook (`OnStartup`, `OnShutdown`) は薄く保ち、重い初期化や長い判断を抱え込まない

### 5.2 TypeScript / Svelte

- `Svelte 5` と `TypeScript` を前提にし、`any` と無検証の type assertion を常用しない
- component state は `\$state`、派生値は `\$derived`、副作用は `\$effect` を基本にする
- component event は callback prop を優先し、`createEventDispatcher` の新規採用は避ける
- event handler は `onclick` などの標準 event 属性を優先する
- `.svelte` は表示とイベント配線に集中させ、副作用や取得判断を template 内へ散らさない

### 5.3 Go

- `error` は無視せず、呼び出し側が判断できる形で返す
- SQL、filesystem、HTTP、外部プロセス呼び出しは文字列直書きや場当たり実装を増やさず、責務をまとめて扱う
- migration や schema 更新のような初期化処理は通常の request 経路へ混ぜない

## 6. ログと機密情報

- ログは原因追跡に必要な情報を残しつつ、機密値、API key、token、ローカル絶対パスを無加工で出さない
- debug 用の詳細と user-facing な表示文言を混同しない
- ログ message は検索しやすい語彙を使い、同じ失敗を複数の曖昧な表現で記録しない
- 観測のための一時ログは恒久仕様へしない


## 8. 禁止事項

- 失敗を無視して処理を継続する実装
- テストを更新せずに仕様変更をコメントや口頭説明だけで補う実装
- generated file を hand-edit する実装
- 機密値や内部診断情報を UI、ログ、外部境界へ無加工で出す実装
- 採用技術、正本仕様、validation を無視して外部テンプレートをそのまま流用する実装

## 9. 参照元

- Wails official docs:
  [`Application Development`](https://wails.io/docs/guides/application-development),
  [`How Does It Work`](https://wails.io/docs/howdoesitwork),
  [`Project Config`](https://wails.io/docs/reference/project-config)
- Svelte official docs:
  [`Svelte 5 Migration Guide`](https://svelte.dev/docs/svelte/v5-migration-guide),
  [`TypeScript`](https://svelte.dev/docs/svelte/typescript)
- repo 固有の正本: [`architecture.md`](./architecture.md), [`tech-selection.md`](./tech-selection.md), [`lint-policy.md`](./lint-policy.md)
