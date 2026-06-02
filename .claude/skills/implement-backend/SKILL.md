---
name: implement-backend
description: "`implementation-module` 内で `backend_implementer` agent が使う backend 実装作業プロトコル。層責務、モジュール内検証の判断基準を提供する。"
---
# Implement Backend

## 目的

この skill は作業プロトコルである。
`backend_implementer` agent が backend 承認済み実装範囲 を実装する時に、usecase、service、repository、adapter の責務整合と 依存方向 を守る判断基準を提供する。

## 対応ロール

- `backend_implementer` が使う。
- 呼び出し元は `implementation-module`、または `backend_implementer` agent を Task ツールで起動した上位 agent とする。
- 返却先は呼び出し元とする。
- 担当成果物は `implement-backend` の出力規約で固定する。

## 呼び出し元から渡される情報

- 単一引き継ぎ入力: `implementation-scope` から切り出された backend 実装用 引き継ぎ 1 件、または `investigation-module` の `修正実行入力`。
- 実行中タスク成果物場所: 実装結果、検証結果、停止理由を書き戻す作業計画フォルダまたは run 成果物フォルダ。
- 実装対象: 変更してよい backend ファイル、symbol、公開接点。
- 対象変更範囲: 実装してよい backend プロダクトコード範囲。
- 依存完了情報: 着手前に完了している必要がある依存対象の完了結果。
- 検証コマンド: 実行を許可された backend-local の harness command。

## 作業前に読む正本

- エージェント実行定義と実行境界は [backend_implementer.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.claude/agents/backend_implementer.md) に従う。
- コーディング規約: [coding-guidelines-backend.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/coding-guidelines-backend.md) とする。
- lint 規約: [lint-policy.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/lint-policy.md) とする。
- architecture 規約: [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) の backend 境界だけを参照する。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。

## skill 内の拘束条件

backend 層責務、依存方向、lint 観点は [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) と [lint-policy.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/lint-policy.md) の backend 章が拘束する。

## 担当ロールが判断してよい範囲

- エラー経路 と 検証 を 承認済み実装範囲 内で閉じる
- 単一引き継ぎ入力 と 承認済み実装範囲 を確認して プロダクトコード だけを変更する
- モジュール内検証 結果 または未実行理由を返す
- usecase / service / repository / adapter の責務を確認する
- `architecture.md` の backend 依存方向に従い、usecase、service、repository、adapter concrete の境界を跨がない
- usecase から repository concrete、実行定義 concrete、driver API を直接参照しない

## skill が扱わない対象

- frontend だけの変更、UI check、backend 境界の再設計は扱わない。
- 承認済み実装範囲外の層 refactor は扱わない。
- プロダクトテスト、検証データ、スナップショット、test helper は変更しない。
- docs や作業流れ文書は変更しない。
- coverage、harness all、repo-local Sonar issue 判定条件は必須終了処理にしない。

## 返す成果物

- 判断結果: backend プロダクトコード実装の完了、未完了、停止の判定を返す。
- 根拠参照: 実装の根拠にした入力、変更箇所、検証結果を返す。
- 不足情報: 実装を完了できない不足項目を返す。
- 次判断材料: 呼び出し元が次を判断できる材料を返す。
- 実装成果物: 単一引き継ぎ入力 の 承認済み実装範囲 に対応する backend プロダクトコードだけを返す。
- 影響範囲修正: 今回変更が直接壊した生成物、公開境界、検証経路、backend 責務内プロダクトコードを修正した場合に、対象、理由、変更結果を返す。
- モジュール内検証結果: `python3 scripts/harness/run.py --suite backend-local` の失敗時は、承認済み実装範囲 または backend 責務内の影響範囲修正で直して再実行し、通過結果または未実行理由を返す。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコード変更の指示を含めない。

## 作業を完了できる条件

- 承認済み実装範囲 または backend 責務内の影響範囲修正 の成果だけが返却されている。
- 検証、未実行項目、残留リスク が 根拠参照 付きで整理されている。
- 単一引き継ぎ入力、実装対象、対象変更範囲、依存完了情報、検証コマンドを確認した。
- `モジュール内検証結果` を返した。

## 作業を止める条件

- frontend だけの変更を実装する時
- UI check を行う時
- backend 境界を設計し直す時
- 単一引き継ぎ入力、実装対象、対象変更範囲、依存完了情報、検証コマンドが不足する場合は停止する。
- `architecture.md` の backend 境界違反が必要になる場合は停止する（controller / usecase / service での concrete 実装 new、service core からの filesystem / Wails 実行定義 / DB driver の concrete API 直接呼び出しを含む）。
- プロダクトテスト、検証データ、スナップショット、test helper の変更が必要になる場合は停止する。
- UI 表示、画面、部品、文言、style、人間承認済み UI 証跡の差分が必要になる場合は停止する。
- secret、trust boundary、API / DTO / DB / schema の意味拡張が必要になる場合は停止する。
- docs 正本化 または `.codex` 作業流れの変更が必要になる場合は停止する。
- 承認済み実装範囲外へ実装を広げる必要があり、今回変更の直接影響または backend 責務内プロダクトコードとして説明できない場合は停止する。
- 停止時は不足項目、衝突箇所、戻し先を返す。
