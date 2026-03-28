# Executable Specs

関連文書: [`index.md`](./index.md), [`spec.md`](./spec.md), [`architecture.md`](./architecture.md), [`tech-selection.md`](./tech-selection.md)

この文書は、細かな仕様や制約を「後で実行して確かめられる形」に寄せるための入口とする。
詳細仕様は長い説明文で増やすのではなく、テスト、acceptance checks、fixture、検証コマンドへ落とす。

## Principles

- 細かな振る舞いは、可能な限りテストで分かる形にする
- 文書はテストや acceptance checks を作るための最小限の契約だけを書く
- 仕様変更では、必要なら対応する test case や acceptance checks も同時に更新する
- 実装がまだない領域では、先に期待結果、失敗条件、観測点を記録する
- gate や review で繰り返し見つかる指摘は、可能な限り harness や acceptance checks に昇格する

## Record Here

- どの種類のテストで何を担保するか
- acceptance checks に必ず入れるべき観点
- fixture や sample input / output の扱い
- 実行可能仕様に昇格させるべき制約

## Current Policy

- UI、実行、DTO、状態遷移の細かな制約は、将来的に対応する test と validation command で表現する
- plan には `Acceptance Checks` を必須で持たせ、詳細仕様の一時的な置き場にする
- plan には `Required Evidence` と `Reroute Trigger` を持たせ、workflow gate が pass/fail を判定できるようにする
- 永続ルールだけを `spec.md` と `architecture.md` に残し、細かな分岐条件はここからテストへ寄せる
- ドメイン / アプリケーション層のルールは Rust の `cargo test` を基本入口として表現する
- UI コンポーネントと画面内の振る舞いは `Vitest` と `@testing-library/svelte` で表現する
- デスクトップ統合の acceptance checks は `tauri-driver` と `WebdriverIO` を基本入口とする
- workflow gate は runtime 品質そのものではなく、plan 適合性、evidence 充足、docs 同期漏れ、reroute 要否を判定する
- xEdit importer の acceptance checks では、`extractData.pas` の raw JSON に含まれる `cells` 空配列と `voicetype` 互換項目を読み込んでも canonical DB モデルが崩れないことを確認する
- 複数入力ファイルの acceptance checks では、1 つの `TRANSLATION_JOB` が複数 `PLUGIN_EXPORT` を参照しつつ、出自情報を保持したまま `TRANSLATION_UNIT` を生成できることを確認する
- 出力 writer の acceptance checks では、複数フィールドを持つ 1 レコードから `TRANSLATION_UNIT` を field 単位で生成し、xTranslator XML の `EDID` / `REC` / `FIELD` / `FORMID` / `Source` / `Dest` / `Status` を lossless に再構成できることを確認する
- ペルソナ生成の acceptance checks では、`MASTER_PERSONA` と `JOB_PERSONA_ENTRY` が混線せず、mod 追加 NPC のジョブ内ペルソナを UI 観測用 DTO へ分けて出せることを確認する
