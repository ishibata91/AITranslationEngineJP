---
name: implement-frontend
description: Codex implementation レーン 側の frontend 実装作業プロトコル。画面導線、状態、Wails bridge の判断基準を提供する。
---
# Implement Frontend

## 目的

この skill は作業プロトコルである。
`frontend_implementer` agent が frontend 承認済み実装範囲 を実装する時に、画面導線、状態 反映、Wails bridge 呼び出しを守る判断基準を提供する。

## 対応ロール

- `frontend_implementer` が使う。
- 呼び出し元は `implement_lane` または `light_change_lane` とする。
- 返却先は呼び出し元とする。
- 担当成果物は `implement-frontend` の出力規約で固定する。

## 入力規約

- frontend 実行入力: `implementation-scope` から切り出された frontend 実装用 引き継ぎ 1 件、または 軽量変更レーンの `軽量変更計画`。
- UI 根拠: 承認済み `ui-design.md`、関連する `screen-design-diff.<screen-id>.md`、または 軽量変更レーンの軽量変更計画と人間確認観点。
- 実行中タスク成果物場所: 実装結果、検証結果、停止理由を書き戻す作業計画フォルダまたは run 成果物フォルダ。
- 実装対象: 変更してよい frontend ファイル、symbol、公開接点。
- 対象変更範囲: 実装してよい frontend プロダクトコード範囲。
- 依存完了情報: 着手前に完了している必要がある依存対象の完了結果。
- 検証コマンド: 実行を許可された frontend-local の harness command。

## 外部参照規約

- エージェント実行定義と実行境界は [frontend_implementer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/frontend_implementer.toml) に従う。
- コーディング規約: [coding-guidelines-frontend.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/coding-guidelines-frontend.md) とする。
- lint 規約: [lint-policy.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/lint-policy.md) とする。
- architecture 規約: [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) の frontend 境界だけを参照する。
- UX 観点正本: [UX-standard.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/UX-standard.md) とする。
- `agent-browser` 利用規約: [agent-browser.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/references/agent-browser.md) とする。
- UI 設計規約: `ui-design.md` を受け取る場合は [ui-design](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/ui-design/SKILL.md) の UI 要件契約に従う。
- 画面設計差分: `screen-design-diff.<screen-id>.md` を受け取る場合は画面ID、画面要素、表示条件、操作、結果、セレクタ（`aria-label`）の根拠にする。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

- なし。UX 観点は `docs/UX-standard.md` と task-local `ui-design.md` の `UX Standard Review` を正本にする。

## 判断規約

- 画面導線と 状態 反映を frontend 実行入力 に合わせる
- Wails bridge 呼び出しの境界を守る
- generated `wailsjs` は gateway 境界に閉じ込める
- 承認済み UI 根拠と `docs/UX-standard.md` に従い、実画面と UI 根拠 の差分を確認する
- セレクタは、画面設計差分の 画面ID と セレクタ（`aria-label`） が実装と一致しているか確認する
- UI 根拠確認の差分は、承認済み実装範囲 で直す差分と、入力不足または人間承認が必要な差分へ分ける
- frontend 実行入力 と 承認済み実装範囲 を確認して プロダクトコード だけを変更する
- `frontend-local` の失敗原因が承認済み実装範囲 外にある場合でも、generated file、生成元、公開境界、または検証を壊した frontend プロダクトコードに限り 影響範囲修正 として直す
- 影響範囲修正は、UI 表示、画面文言、layout、style、承認済み UI 根拠を越える変更に使わない
- 明確なブロッカーがない限りはレーンを中断せずに成果物の生成を継続すること。
- UI 状態 の初期値と更新条件を確認する
- [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) の frontend 境界に従い、View、ScreenController、Frontend UseCase、Gateway の責務を跨がない
- generated `wailsjs` と backend DTO の import は `frontend/src/controller/wails/` に閉じ込める
- generated file が `frontend-local` の失敗原因である場合は、generated file を直接編集せず、生成元または公開境界を直す

## 非対象規約

- backend だけの変更、design mock 作成、UI check だけの作業は扱わない。
- UI 根拠にない改善は追加しない。
- プロダクトテスト、検証データ、スナップショット、test helper は変更しない。
- Wails bridge と backend DTO の境界を迂回しない。
- docs や作業流れ文書は変更しない。
- coverage、harness all、repo-local Sonar issue 判定条件は必須完了条件にしない。

## 出力規約

- 判断結果: frontend プロダクトコード実装の完了、未完了、停止の判定を返す。
- 根拠参照: 実装の根拠にした入力、変更箇所、検証結果を返す。
- 不足情報: 実装を完了できない不足項目を返す。
- 次判断材料: 呼び出し元が次を判断できる材料を返す。
- 実装成果物: frontend 実行入力 の 承認済み実装範囲 に対応する frontend プロダクトコードだけを返す。
- 影響範囲修正: generated file、生成元、公開境界、または検証を壊した frontend プロダクトコードを修正した場合に、対象、理由、変更結果を返す。
- レーン内検証結果: `python3 scripts/harness/run.py --suite frontend-local` の失敗時はその場で直して再実行し、通過結果または未実行理由を返す。
- UI 根拠確認結果: 実画面と UI 根拠 の一致、差分、画面ID と セレクタ（`aria-label`）差分、未確認理由、`docs/UX-standard.md` との対応を返す。
- UI証跡参照: `agent-browser` の snapshot、screenshot、console、errors の参照または未取得理由を返す。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコード変更の指示を含めない。

## 完了規約

- 承認済み実装範囲 または許可された 影響範囲修正 の成果だけが返却されている。
- 検証、未実行項目、残留リスク が 根拠参照 付きで整理されている。
- frontend 実行入力、実装対象、対象変更範囲、依存完了情報、検証コマンドを確認した。
- UI 根拠 を確認した。
- 画面導線と 状態 反映を確認した。
- Wails bridge 境界を確認した。
- generated `wailsjs` を gateway 境界に閉じ込めた。
- UI 状態 の初期値と更新条件を確認した。
- 承認済み UI 根拠、`docs/UX-standard.md`、frontend コーディング規約に合わせて実装した。
- 実画面と UI 根拠 の一致確認結果を返した。
- 画面設計差分がある場合は、画面ID と セレクタ（`aria-label`） の実装差分を確認した。
- UI 根拠確認結果は、`agent-browser` の snapshot、screenshot、console、errors の根拠または未取得理由を含んでいる。
- frontend lint と format:check で拾われる境界違反を確認した。
- frontend 変更として `python3 scripts/harness/run.py --suite frontend-local` を実行し、失敗した場合は承認済み実装範囲 または許可された 影響範囲修正 でその場で直して再実行し、通過結果または未実行理由を返した。

## 停止規約

- backend だけの変更を実装する時
- design mock を作る時
- UI check だけを行う時
- frontend 実行入力、UI 根拠、実装対象、対象変更範囲、依存完了情報、検証コマンドが不足する場合は停止する。
- 通信境界を迂回する必要がある場合は停止する。
- View、ScreenController、Frontend UseCase から generated `wailsjs` を直接 import する必要がある場合は停止する。
- gateway 以外で backend DTO 変換が必要な場合は停止する。
- プロダクトテスト、検証データ、スナップショット、test helper の変更が必要になる場合は停止する。
- 実画面と UI 根拠 の差分が承認済み実装範囲 外の修正を必要とする場合は停止し、人間承認へ戻す。
- UI 根拠確認に必要な実画面確認根拠を取得できない場合は停止し、未取得理由と戻し先を返す。
- `python3 scripts/harness/run.py --suite frontend-local` の失敗原因が承認済み実装範囲 外にある場合は、generated file、生成元、公開境界、または検証を壊した frontend プロダクトコードに限り 影響範囲修正 として直す。
- `python3 scripts/harness/run.py --suite frontend-local` の失敗原因が generated file にある場合は、generated file を直接編集せず、生成元または公開境界を 影響範囲修正 として直す。
- `python3 scripts/harness/run.py --suite frontend-local` の失敗原因を直すために UI 表示、画面文言、layout、style、承認済み UI 根拠を越える変更が必要な場合は停止し、人間承認へ戻す。
- 承認済み実装範囲外へ実装を広げる必要があり、許可された 影響範囲修正 に該当しない場合は停止し、人間承認へ戻す。
- 停止時は不足項目、衝突箇所、戻し先を返す。
