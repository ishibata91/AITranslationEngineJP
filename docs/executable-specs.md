# Executable Specs

関連文書: [`index.md`](./index.md), [`spec.md`](./spec.md), [`architecture.md`](./architecture.md), [`tech-selection.md`](./tech-selection.md)

この文書は、細かな仕様や制約を「後で実行して確かめられる形」に寄せるための入口とする。
詳細仕様は長い説明文で増やすのではなく、テスト、acceptance checks、fixture、検証コマンドへ落とす。

## Principles

- 細かな振る舞いは、可能な限りテストで分かる形にする
- 文書はテストや acceptance checks を作るための最小限の契約だけを書く
- 仕様変更では、必要なら対応する test case や acceptance checks も同時に更新する
- 実装がまだない領域では、先に期待結果、失敗条件、観測点を記録する

## Record Here

- どの種類のテストで何を担保するか
- acceptance checks に必ず入れるべき観点
- fixture や sample input / output の扱い
- 実行可能仕様に昇格させるべき制約

## Current Policy

- UI、実行、DTO、状態遷移の細かな制約は、将来的に対応する test と validation command で表現する
- plan には `Acceptance Checks` を必須で持たせ、詳細仕様の一時的な置き場にする
- 永続ルールだけを `spec.md` と `architecture.md` に残し、細かな分岐条件はここからテストへ寄せる
