---
name: implement-frontend
description: Codex implementation レーン 側の frontend 実装作業プロトコル。画面導線、状態、Wails bridge の判断基準を提供する。
---
# Implement Frontend

## 目的

この skill は作業プロトコルである。
`implementation_implementer` agent が frontend 承認済み実装範囲 を実装する時に、画面導線、状態 反映、Wails bridge 呼び出しを守る判断基準を提供する。

## 対応ロール

- `implementation_implementer` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- 担当成果物は `implement-frontend` の出力規約で固定する。

## 入力規約

- 不足時の扱い: 入力に 根拠参照、担当者、承認状態が不足する場合は推測で補わない。
- 単一引き継ぎ入力: implementation-scope から抽出済みの 引き継ぎ 1 件。
- 承認記録: 人間が承認した実装範囲の根拠参照。
- 実装対象: 変更してよい frontend ファイル、symbol、公開接点。
- 承認済み実装範囲: 実装してよい frontend プロダクトコード範囲。
- 依存解消状態: 依存対象が完了済みかを示す状態。
- 任意入力: 実装小範囲、レーン内検証コマンド、implementation_scenario_tester 出力。

## 外部参照規約

- エージェント実行定義とツール権限は [implementation_implementer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_implementer.toml) の 書き込み許可 / 実行許可 とする。
- コーディング規約: [coding-guidelines-frontend.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/coding-guidelines-frontend.md) とする。
- lint 規約: [lint-policy.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/lint-policy.md) とする。
- architecture 規約: [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) の frontend 境界だけを参照する。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

## 判断規約

- 画面導線と 状態 反映を implementation-scope に合わせる
- Wails bridge 呼び出しの境界を守る
- generated `wailsjs` は gateway 境界に閉じ込める
- affected UI の manual flow を確認できる状態にする
- UI check に必要な 根拠 を残す
- ユーザーが利用することを前提に、ページデザイン、ボタン配置、表示テキストを承認済み UI 要件と frontend コーディング規約に合わせる
- 単一引き継ぎ入力 と 承認済み実装範囲 を確認して プロダクトコード だけを変更する
- `APIテスト` 先行時だけ implementation_scenario_tester 出力 も確認する

- 単一引き継ぎ入力、affected UI flow を確認する
- `APIテスト` 先行時だけ implementation_scenario_tester 出力 を確認する
- console エラー の有無を 終了処理 に残す
- UI 状態 の初期値と更新条件を確認する
- 主要操作、取消、戻る、破壊的操作のボタン配置と文言がユーザーの作業順に沿っていることを確認する
- ページ見出し、ラベル、説明、エラー、空状態、完了状態のテキストがユーザーの次の行動を示すことを確認する
- [lint-policy.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/lint-policy.md) の frontend lint 内訳を確認し、`npm run lint` と `format:check` で拾われる観点を先に意識する
- [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) の frontend 境界に従い、View、ScreenController、Frontend UseCase、Gateway の責務を跨がない
- generated `wailsjs` と backend DTO の import は `frontend/src/controller/wails/` に閉じ込める

## 非対象規約

- backend だけの変更、design mock 作成、UI check だけの作業は扱わない。
- design にない改善は追加しない。
- プロダクトテスト、検証データ、スナップショット、test helper は変更しない。
- Wails bridge と backend DTO の境界を迂回しない。
- docs や作業流れ文書は変更しない。
- coverage、harness all、repo-local Sonar issue 判定条件は必須終了処理にしない。

## 出力規約

- 判断結果: frontend プロダクトコード実装の完了、未完了、停止の判定を返す。
- 根拠参照: 実装の根拠にした入力、変更箇所、検証結果を返す。
- 不足情報: 実装を完了できない不足項目を返す。
- 次判断材料: `implement_lane` が次を判断できる材料を返す。
- 実装成果物: 単一引き継ぎ入力 の 承認済み実装範囲 に対応する frontend プロダクトコードだけを返す。
- レーン内検証結果: `python3 scripts/harness/run.py --suite frontend-local` の失敗時はその場で直して再実行し、通過結果または未実行理由を返す。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコード変更の指示を含めない。

## 完了規約

- 承認済み実装範囲 内の成果だけが返却されている。
- 検証、未実行項目、残留リスク が 根拠参照 付きで整理されている。
- 単一引き継ぎ入力、承認記録、実装対象、承認済み実装範囲を確認した。
- 画面導線と 状態 反映を確認した。
- Wails bridge 境界を確認した。
- generated `wailsjs` を gateway 境界に閉じ込めた。
- affected UI flow と console エラー を確認した。
- ユーザーが利用することを意識し、ページデザイン、ボタン配置、表示テキストを承認済み UI 要件と frontend コーディング規約に合わせて実装した。
- frontend lint と format:check で拾われる境界違反を確認した。
- frontend 変更として `python3 scripts/harness/run.py --suite frontend-local` を実行し、失敗した場合は承認済み実装範囲 内でその場で直して再実行し、通過結果または未実行理由を返した。

## 停止規約

- backend だけの変更を実装する時
- design mock を作る時
- UI check だけを行う時
- 単一引き継ぎ入力、実装対象、承認記録、承認済み実装範囲が不足する場合は停止する。
- 通信境界を迂回する必要がある場合は停止する。
- View、ScreenController、Frontend UseCase から generated `wailsjs` を直接 import する必要がある場合は停止する。
- gateway 以外で backend DTO 変換が必要な場合は停止する。
- プロダクトテスト、検証データ、スナップショット、test helper の変更が必要になる場合は停止する。
- `python3 scripts/harness/run.py --suite frontend-local` の失敗原因が承認済み実装範囲 外にある場合は停止する。
- 承認済み実装範囲外へ実装を広げる必要がある場合は停止する。
- 停止時は不足項目、衝突箇所、戻し先を返す。
- UI check に必要な 根拠 を残した。
