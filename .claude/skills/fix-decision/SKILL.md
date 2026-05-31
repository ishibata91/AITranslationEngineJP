---
name: fix-decision
description: "`investigation-module` 内で `fix_decider` agent が使う作業プロトコル。観測記録から仮説、観測ログ検証、確定原因、修正方針、禁止修正を固定する。"
---
# Fix Decision

## 目的

`fix-decision` は、人間が確認した不具合、レビュー非通過、検証失敗の観測記録から、恒久修正へ渡せる修正方針判断へ変換する作業プロトコルである。
`fix_decider` が複数の原因仮説、観測ログによる仮説検証、確定原因、採用する修正方針、禁止する対症療法を task 内成果物として固定する時に使う。

## 対応ロール

- `fix_decider` が使う。
- 呼び出し元は `investigation-module`、または `fix_decider` agent を Task ツールで起動した上位 agent とする。
- 返却先は呼び出し元とする。
- 担当成果物は `修正方針判断` とする。

## 呼び出し元から渡される情報

- 必須人間観測記録: 人間が見た画面、操作、ログ、失敗、期待との差分。
- 必須作業計画フォルダ: `修正方針判断` を置く `docs/exec-plans/active/<task-id>/`。
- 必須Wails接続対象: `investigation-module` が事前準備で単一化した Wails process または接続先。

## 作業前に読む正本

- エージェント実行定義と実行境界は [fix_decider.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.claude/agents/fix_decider.md) に従う。
- 修正モジュールの進行境界は [investigation-module](/Users/iorishibata/Repositories/AITranslationEngineJP/.claude/skills/investigation-module/SKILL.md) に従う。
- 修正方針判断の報告形式は [fix-decision-report-template.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.claude/skills/fix-decision/fix-decision-report-template.md) に従う。
- 画面設計書正本は [screen-design](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/screen-design/README.md) に従う。
- ユースケース正本は [usecases](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/usecases/README.md) に従う。
- 観測ログ仕様は [observability-logging.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/observability-logging.md) に従う。
- ブラウザ操作は `chrome-devtools` MCP ツール群（`mcp__plugin_chrome-devtools-mcp_chrome-devtools__*`）を MCP ツールとして実行する。
- 外部成果物が不足または衝突する場合は停止し、衝突箇所を返す。

## skill 内の拘束条件

### 修正方針判断の必須観点

| 観点 | 拘束する判断 |
| --- | --- |
| `観測済み問題` | 人間観測記録から確認できる問題だけを固定する。 |
| `画面再現確認` | 画面設計のセレクタに従い、人間観測記録のユーザー操作を `chrome-devtools` MCP ツールで再現する。 |
| `原因仮説` | 画面操作結果、DB 状態、ログの観測事実から複数の原因候補を立て、検証する順序と根拠を固定する。 |
| `観測ログ検証` | 仮説を否定または支持するために追加した一時ログ、観測結果、削除確認を固定する。 |
| `確定原因` | 観測で確定した原因だけを固定する。 |
| `採用する修正方針` | 仕様が不足していない場合だけ恒久修正を固定する。 |
| `禁止する修正` | 新しい状態値の追加、症状だけを隠す分岐などを具体的に固定する。 |

## 担当ロールが判断してよい範囲

- 観測事実と仮説を分ける。
- 人間観測記録の画面操作は、原因仮説を固定する前に `chrome-devtools` MCP ツールで再現確認する。
- 画面再現確認では、呼び出し元から渡された Wails process または接続先へアクセスする。
- 画面再現確認では、画面設計書の selector または `data-testid` に従ってユーザー操作を再現する。
- 画面再現確認では、`chrome-devtools` MCP ツール群を MCP ツールとして実行する。
- 画面操作結果、DB 状態、ログを観測事実として分けて整理する。
- 画面操作結果、DB 状態、ログの差分から、どの層で期待と実際が分かれたかを特定する。
- 特定した分岐点ごとに、複数の原因仮説を立てる。
- 画面操作結果だけで原因仮説を固定せず、DB 状態またはログで仮説を検証する。
- 仮説ごとに観測ログを仕込み、観測結果で否定または支持を判断する。
- 観測ログで否定された仮説は、確定原因から除外する。
- 観測ログで検証できていない仮説は、確定原因として固定しない。
- 追加した一時観測ログは、修正方針判断を返す前に削除する。
- 原因が仮説に留まる場合は修正方針を採用しない。
- 表面症状ではなく、観測で確定した原因を扱う。
- 既存状態モデルで説明できる場合は、新しい状態値を採用しない。
- 修正方針が仕様変更、機能追加、受け入れ条件の新規判断を含む可能性がある場合は停止する。
- 既存期待を説明するユースケース記述が不足している可能性がある場合は停止する。
- 影響ファイルは候補として返してよいが、実装 agent の変更ファイルを確定しない。

## skill が扱わない対象

- 修正実行入力の作成は扱わない。
- 一時観測ログ以外のプロダクトコード変更は扱わない。
- プロダクトテストの変更は扱わない。
- docs 正本本文の更新は扱わない。

## 返す成果物

- 判断結果: `修正方針判断` の完了、未完了、停止の判定を返す。
- 観測済み問題: 根拠から確認できる問題を返す。
- 画面再現確認: `chrome-devtools` MCP ツールで再現に使った再現手順、操作結果、画面状態、証跡参照を返す。
- 確定原因: 観測で確定した原因を返す。
- 採用する修正方針: 恒久修正として採用する方針を返す。
- 禁止する修正: 実装 agent に許可しない対症療法を返す。
- 影響ファイル候補: 観測事実に基づく候補を返す。

## 作業を完了できる条件

- 観測済み問題、画面再現確認、確定原因、採用する修正方針、禁止する修正が返っている。
- 追加した一時観測ログが削除されている。
- 仕様が不足していると判断した場合は停止している。

## 作業を止める条件

- 人間観測記録が不足する場合は停止する。
- Wails 接続対象が渡されていない場合は停止する。
- 画面設計の selector に従ったユーザー操作を `chrome-devtools` MCP ツールで再現確認できない場合は停止する。
- `chrome-devtools` MCP ツールが利用できない場合は停止する。
- 観測ログを追加または確認できず、仮説を検証できない場合は停止する。
- 追加した一時観測ログを削除できない場合は停止する。
- 原因が仮説に留まり、採用する修正方針を固定できない場合は停止する。
- 人間観測記録と観測ログ検証から、既存期待の記述不足か新規判断が必要な変更かを分けられない場合は停止する。
- 修正方針が対症療法か恒久修正か判断できない場合は停止する。
- 停止時は不足項目、衝突箇所、固定できない判断、戻し先を返す。
