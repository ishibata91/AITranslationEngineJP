# バックエンド コーディング規約

関連文書: [`coding-guidelines.md`](./coding-guidelines.md), [`architecture.md`](./architecture.md), [`lint-policy.md`](./lint-policy.md), [`observability-logging.md`](./observability-logging.md)

本書は、`internal/` と Wails backend 起点の Go 実装規約を定義する。
usecase、service、repository、adapter、transport boundary の責務を対象にする。
frontend と test の規約は別文書を正本にする。

## 1. Go 基本

- `gofmt` と `goimports` の出力を正とし、整形方針を個別判断で変えない
- public に見える名前は省略しすぎず、利用側が文脈なしで意味を追える粒度にする
- interface は使う側で定義し、1 から 3 method 程度の小さい境界に保つ
- 実装は concrete struct を返し、依存を受ける側では interface を受ける
- constructor で依存を注入し、usecase や service 内で concrete 実装を直接生成しない

## 2. Go ファイル分割

- controller、usecase、service、repository、adapter は責務が異なるため同じファイルへ混ぜない
- request / response DTO 変換、業務手順、永続化 SQL、外部入出力、validation は責務単位で分ける
- 小さい private helper は、親処理だけを補助し、外部境界や独立した検証単位を持たない場合だけ同じファイルに置く
- 1 ファイルに複数の public entrypoint がある場合は、同じ port、同じ状態、同じ検証意図を共有する時だけ許容する
- ファイル名は層名より責務名を優先し、検索時に対象処理を見つけられる名前にする

## 3. 層ごとの責務分離

- `Controller` は Wails Bind の入出力変換だけを持ち、業務処理や repository concrete を持たない
- `Backend UseCase` は操作単位の手順を持ち、adapter concrete や DB driver を直接扱わない
- `Service` は実処理と port 呼び出しを持ち、filesystem、Wails runtime、SQL driver の concrete API を直接扱わない
- `Repository` と adapter は外部入出力の concrete 実装を持ち、usecase の手順判断を持たない
- `internal/bootstrap/` だけが production graph を組み立てる

## 4. エラー処理

- `error` は無視せず、呼び出し側が判断できる形で返す
- 失敗には処理名、対象、境界などの文脈を付けて wrap する
- 通常フローの失敗は戻り値や明示的な失敗表現で扱い、`panic` を制御フローに使わない
- validation error、外部依存の実行失敗、想定外障害を同じ重さで混ぜない
- cleanup が必要な処理では、途中失敗時にも後始末が漏れないようにする

## 5. 境界と永続化

- SQL、filesystem、HTTP、外部プロセス呼び出しは責務をまとめ、呼び出し元へ散らさない
- file path、外部 URL、provider 設定値、外部プロセス入力は使用直前に再検証する
- migration や schema 更新のような初期化処理は通常の request 経路へ混ぜない
- DTO や struct がある契約を、文字列連結や ad-hoc な map で表現しない
- Wails `Bind` する public method は transport boundary として扱い、業務処理や永続化詳細を直書きしない

## 6. ログと機密情報

- 観測ログの出力先、payload、禁止事項は [`observability-logging.md`](./observability-logging.md) に従う
- ログは原因追跡に必要な情報を残しつつ、機密値を無加工で出さない
- user-facing message と internal diagnostic を分ける
- ログ message は検索しやすい語彙を使い、同じ失敗を複数の曖昧な表現で記録しない
- 観測のための一時ログは恒久仕様へしない

## 7. 禁止事項

- controller、usecase、service で concrete 実装を直接 `new` する実装
- service core から filesystem、Wails 実行定義、DB driver の concrete API を直接呼ぶ実装
- controller、usecase、service、repository、adapter の責務を 1 ファイルへまとめる実装
- 失敗を無視して処理を継続する実装
- 初期化、migration、schema 更新を通常 request 経路へ混ぜる実装

## 8. 参照元

- Wails official docs:
  [`Application Development`](https://wails.io/docs/guides/application-development),
  [`Project Config`](https://wails.io/docs/reference/project-config)
- 輸入元: `../everything-claude-code/rules/golang/coding-style.md`
- 輸入元: `../everything-claude-code/.cursor/rules/golang-patterns.md`
- 輸入元: `../everything-claude-code/rules/common/coding-style.md`
