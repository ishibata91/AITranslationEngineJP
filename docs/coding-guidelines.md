# コーディング規約

関連文書: [`index.md`](./index.md), [`architecture.md`](./architecture.md), [`tech-selection.md`](./tech-selection.md), [`lint-policy.md`](./lint-policy.md)

本書は、実装規約の入口である。
実装時に守る規約は、変更対象に合わせて frontend、backend、test の正本を読む。
構造責務の正本は [`architecture.md`](./architecture.md) とする。

## 1. 分割

- frontend 実装規約: [`coding-guidelines-frontend.md`](./coding-guidelines-frontend.md)
- backend 実装規約: [`coding-guidelines-backend.md`](./coding-guidelines-backend.md)
- test 実装規約: [`coding-guidelines-tests.md`](./coding-guidelines-tests.md)

## 2. ファイル分割方針

- 1 ファイルへ置く処理は、同じ責務、同じ変更理由、同じ検証単位でそろえる
- 境界処理、状態管理、永続化、外部入出力、表示、変換処理を 1 ファイルへ混ぜない
- private helper は、親責務を補助し、単独の境界や検証単位を持たない場合だけ同じファイルに置く
- 同じファイル内で複数の public entrypoint が別々の責務を持つ場合は、責務単位で分ける
- ファイル分割は層境界を増やすためではなく、既存の `architecture.md` の責務境界を保つために行う

## 3. 共通原則

- 実装は採用技術と正本仕様に合わせ、外部テンプレートや一般論を無検証で持ち込まない
- コメントは `何をしているか` ではなく、`なぜその判断が必要か` を短く補足する
- 命名は略語より意味を優先し、役割と責務が読める名前にする
- 境界入力は使用直前に検証し、失敗を握りつぶさない
- production wiring は `architecture.md` の composition root に置き、View や中核層で具象依存を生成しない
- 機密値、API key、token、ローカル絶対パスを UI、ログ、外部境界へ無加工で出さない

## 4. 参照元

- repo 固有の正本: [`architecture.md`](./architecture.md), [`tech-selection.md`](./tech-selection.md), [`lint-policy.md`](./lint-policy.md)
- 輸入元: `../everything-claude-code/rules/common/coding-style.md`
- 輸入元: `../everything-claude-code/rules/golang/coding-style.md`
- 輸入元: `../everything-claude-code/rules/typescript/coding-style.md`
- 輸入元: `../everything-claude-code/rules/common/testing.md`
- 輸入元: `../everything-claude-code/rules/golang/testing.md`
- 輸入元: `../everything-claude-code/rules/typescript/testing.md`
