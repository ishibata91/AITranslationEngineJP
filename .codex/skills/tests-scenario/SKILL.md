---
name: tests-scenario
description: Codex implementation lane 側の scenario artifact を プロダクトテスト に反映する作業プロトコル。
---
# Tests Scenario

## 目的

この skill は作業プロトコルである。
`implementation_tester` agent が承認済みシステムテスト設計を プロダクトテスト に落とす時に、主要 outcome を決定的に証明する判断基準を提供する。
この skill の主対象は `UI人間操作E2E` と `APIテスト` である。

## 対応ロール

- `implementation_tester` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- owner artifact は `tests-scenario` の出力規約で固定する。

## 入力規約

- scenario artifact の観点を プロダクトテスト にする時
- happy path と主要 failure path を整理する時
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。

## 外部参照規約

- エージェント実行定義とツール権限は [implementation_tester.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_tester.toml) の `allowed_write_paths` / `allowed_commands` とする。
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

## 判断規約

- 各 test method は 1 つの scenario outcome だけを証明する
- setup は決定的にする
- test body に条件分岐を入れない
- runtime event 完了は completion event を観測点にする
- `UI人間操作E2E` は、承認済みシナリオの開始操作を模倣する
- UI が入口のシナリオでは、画面操作、ファイル選択、フォーム入力などのユーザー入力を開始点にする
- `APIテスト` は、public seam、request / response contract、外部入力開始、主要観測点を開始点にする
- 裏側の直接呼び出しや fixture 直接投入だけの試験は、明示された補助試験でない限り主 `UI人間操作E2E` にしない

- Arrange / Act / Assert が body 構造で読めるようにする
- happy path と failure path を別 test case に分ける
- fixture や helper は scenario を支える範囲に限定する
- UI が入口の場合は、ユーザー入力から得られる値を `UI人間操作E2E` の検証対象にする
- `APIテスト` では request / response contract と external input start を検証対象にする

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。

## 完了規約

- 承認済み owned_scope 内の成果だけが返却されている。
- validation、未実行項目、residual risk が source_ref 付きで整理されている。
- scenario outcome を 1 test 1 outcome に分けた。
- `UI人間操作E2E` は、ユーザー入力の模倣を開始点にした。
- `APIテスト` は、public seam と外部入力開始を開始点にした。
- happy path と failure path を別 test case にした。
- runtime event 完了の観測点を明示した。

## 停止規約

- unit branch だけを補う時
- scenario artifact が未承認の時
- 原因未確定の regression test を書く時
- プロダクトコードの修正が主目的の時
- 新しい要件解釈を足さない
- paid real AI API を呼ばない
- UI 入口の `UI人間操作E2E` を裏側の直接呼び出しだけで代替しない
- 停止時は不足項目、衝突箇所、戻し先を返す。
- 新しい要件解釈を足さなかった場合は停止する。
- test body に条件分岐を入れなかった場合は停止する。
- UI 入口の `UI人間操作E2E` を裏側の直接呼び出しだけで代替しなかった場合は停止する。
- paid real AI API を呼ばなかった場合は停止する。
