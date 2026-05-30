---
name: implement-lane
description: 新規実装と機能拡張を、必要な成果物だけを依存順に作る DAG として進め、人間ゲートと責務境界を不変条件として守る作業プロトコル。
---
# Implement Lane

## 目的

`implement-lane` は、新規実装と機能拡張の進行判断を、成果物 DAG と引き継ぎへ固定する作業プロトコルである。
`implement_lane` agent が、設計、実装、テスト、検証、正本化、commit の各成果物を、担当 agent への引き継ぎとして進める時に使う。

## 進め方の原則

この skill は固定手順ではなく、必要な成果物だけを依存順に作るための判断材料である。
進める時は次の 3 原則に従う。

- 必要なものだけ作る: 各成果物の `成立条件` を満たす時だけ着手する。満たさない成果物は省略し、省略理由を `plan.md` に 1 行残す。
- 依存を待つ: `依存対象` の成果物が承認済みまたは停止として記録されるまで着手しない。
- 迷ったら緩めない: 判断に迷う場合は「不変条件」を優先し、人間ゲート、責務境界、安全境界を緩める方向に倒さない。

## 対応ロール

- `implement_lane` が使う。
- 呼び出し元は人間とする。
- 返却先は人間とする。
- 起動担当 agent は `designer`、`diagrammer`、`test_designer`、`backend_implementer`、`frontend_implementer`、`integration_implementer`、`observability_implementer`、`browser_confirmation`、`docs_updater` とする。

## 呼び出し元から渡される情報

- 依頼要約: 新規実装または機能拡張として扱う依頼内容。
- 作業計画フォルダ: task 内成果物を置く `docs/exec-plans/active/<task-id>/`。
- 作業branch: 既定名 `codex/<task-id>` の local branch。
- 統合先branch: 既定名 `master` の local branch。
- 既存成果物: 作業計画フォルダに既にある task 内成果物。
- 人間介入状態: 人間レビュー、承認、差し戻し、追加質問の記録。

## 成果物 DAG

各成果物は `依存対象` が揃い、かつ `成立条件` を満たす時だけ着手する。
`起動先 agent` が「なし」の行は `implement_lane` が直接作る。
`起動先 agent` がある行は、その agent を Task ツールで起動して作らせ、`implement_lane` は本文を代筆しない。

| 成果物ID | 担当 | 依存対象 | 成立条件 | 起動先 agent |
| --- | --- | --- | --- | --- |
| `task 枠` | `implement_lane` | `[]` | 常に | なし |
| `着手計画` | `implement_lane` | `task 枠` | 常に | なし |
| `着手計画人間確認` | 人間 | `着手計画` | 常に | 人間 |
| `branch 準備` | `implement_lane` | `task 枠` | 常に | なし |
| `詳細仕様差分` | `designer` | `着手計画人間確認` | 仕様変更または仕様追加がある | `designer` |
| `画面設計差分` | `designer` | `詳細仕様差分?` | 画面変更がある | `designer` |
| `設計差分図` | `diagrammer` | `詳細仕様差分?`, `画面設計差分?` | 構造変更または画面変更がある（揃える場合は人間設計レビュー前に必須） | `diagrammer` |
| `人間設計レビュー` | 人間 | `詳細仕様差分?`, `画面設計差分?`, `設計差分図?` | 常に | 人間 |
| `実装範囲` | `designer` | `人間設計レビュー` | 常に | `designer` |
| `テスト設計` | `test_designer` | `人間設計レビュー` | 常に（`実装範囲` と並列可） | `test_designer` |
| `実装引き継ぎ入力` | `implement_lane` | `実装範囲` | 常に | なし |
| `frontend 実装` | `frontend_implementer` | `実装引き継ぎ入力` | UI がある | `frontend_implementer` |
| `Storybook 入力確認` | `implement_lane` | `frontend 実装` | UX デザインの変更がある（UI の裏側だけの変更では着手しない） | なし |
| `frontend 実装後人間レビュー` | `implement_lane` | `Storybook 入力確認?`, `storybook-review-loop.md` | UX デザインの変更がある | なし |
| `Storybook 後画面設計差分整合` | `designer` | `frontend 実装後人間レビュー?` | Storybook レビューループ後に画面仕様が変わった | `designer` |
| `合意済み frontend 保護` | `implement_lane` | `frontend 実装後人間レビュー?`, `Storybook 後画面設計差分整合?` | UX デザインの変更がある | なし |
| `backend 実装` | `backend_implementer` | `実装引き継ぎ入力`, `合意済み frontend 保護?` | backend 変更がある | `backend_implementer` |
| `統合境界実装` | `integration_implementer` | `backend 実装`, `合意済み frontend 保護?` | frontend と backend を接続する | `integration_implementer` |
| `シナリオテスト` | `implementation_scenario_tester` | `テスト設計`, `backend 実装?`, `合意済み frontend 保護?`, `統合境界実装?` | 利用者経路を証明する | `implementation_scenario_tester` |
| `単体テスト` | `implementation_unit_tester` | `backend 実装?`, `合意済み frontend 保護?`, `統合境界実装?` | 実装済み責務を証明する | `implementation_unit_tester` |
| `観測ログ追加` | `observability_implementer` | 完了済み実装・テスト成果物 | 実行時にしか確定しない値または原因分離が要る分岐がある（`最終検証` の前） | `observability_implementer` |
| `最終検証` | `implement_lane` | `観測ログ追加?` | 常に | なし |
| `実装後ブラウザ確認` | `browser_confirmation` | `最終検証` | UI 経路または実画面挙動の変更がある（`最終検証` の後） | `browser_confirmation` |
| `正本化判断` | `implement_lane` | `最終検証`, `実装後ブラウザ確認?` | 仕様変更または仕様追加がある | `docs_updater?` |
| `詳細仕様正本反映` | `docs_updater` | `正本化判断` | 人間承認済みの恒久仕様がある | `docs_updater` |
| `作業 commit` | `implement_lane` | 全完了または停止済み成果物 | 常に | なし |
| `マージ準備入力` | `implement_lane` | `作業 commit` | 常に | `merge_lane` |

## 不変条件

次は「必要だからやる」判断の対象にせず、常に守る。緩めそうな場合は停止して人間へ返す。

### 人間ゲート

- `着手計画人間確認` を経ずに `詳細仕様差分` 以降の設計成果物へ進めない。
- `設計差分図` の成立条件を満たす場合は、揃えずに `人間設計レビュー` へ進めない。省略する場合は省略理由を `plan.md` に残す。
- 設計成果物（詳細仕様差分、画面設計差分、実装範囲）と UI 承認は AI だけで確定しない。
- 設計レビューが差し戻しまたは追加質問の場合は、新規 agent を起こさず、差し戻し対象を作った同じ agent を再起動し、会話文脈の代わりに差し戻し内容を引き継ぎ入力へ明示する。

### 責務境界

- `designer`、`diagrammer`、`test_designer`、`docs_updater`、実装 agent、テスト agent、`browser_confirmation` の担当成果物の本文を `implement_lane` が代筆しない。担当 agent の返却結果または人間介入状態の転記だけ行う。
- backend、frontend、統合境界 は別成果物として扱い、単一の実装成果物に束ねない。
- 起動先 agent には会話文脈を引き継がず、必要情報を引き継ぎ入力へ明示する。起動先 agent に下位 agent を起動させない。
- 終わったサブエージェントは逐次閉じる。設計中 agent（`designer`、`diagrammer`）は担当成果物が承認済みまたは停止として記録されるまで閉じない。

### UI 順序ゲート（UX デザイン変更がある task）

- `frontend 実装` を `backend 実装` より先に着手する。
- Storybook レビューループは人間が別セッションで `story-book-review-loop` に従って実行する。`implement_lane` は `Storybook 入力確認` を固定した時点で停止し、起動も直接実行もしない。
- `frontend 実装後人間レビュー` の承認と、確認対象 story の通常分類への復帰がないまま `合意済み frontend 保護` へ進めない。
- Storybook レビューループ後に画面仕様が変わった場合は、`designer` に戻して plan 内の `screen-design-diff.<screen-id>.md` を更新させてから先へ進む。
- `合意済み frontend 保護`（承認済み画面、表示規則、確認済み Storybook 状態、変更禁止範囲）を固定してから `backend 実装`、`統合境界実装` へ進める。
- 後続実装で UI 表示、画面文言、layout、style、承認済み画面設計差分を越える変更が必要になった場合は、続けず `frontend 実装` の再実行入力または人間への返却を固定する。

### 検証順序ゲート

- `観測ログ追加` の成立条件を満たす場合は、追加を経ずに `最終検証` へ進めない。観測ログが停止した場合も `最終検証` へ進めない。省略する場合は省略理由を `plan.md` に残す。
- `実装後ブラウザ確認` の成立条件を満たす場合は `最終検証` の後に行う。確認 URL、起動状態、操作経路、操作期待値、禁止操作、安全条件、証跡出力先が揃わない場合は進めない。省略する場合は省略理由を `plan.md` に残す。
- backend 変更があれば `python3 scripts/harness/run.py --suite backend-local`、frontend 変更があれば `--suite frontend-local` を許可済みコマンドとして実行し、失敗時は担当 agent がその場で直して再実行する。

### 安全境界

- `implement_lane` はプロダクトコード、プロダクトテスト、人間承認なしの docs 正本を直接変更しない。
- 検証失敗原因が generated file にある場合は、直接編集せず生成元または公開境界の修正として担当 agent に渡す。
- 仕様変更または仕様追加があれば `正本化判断` を必須にし、人間承認済みの恒久仕様だけ `詳細仕様正本反映` で `docs/detail-specs/` へ反映する。
- local merge、`docs/exec-plans/completed/` への移動、remote repository の変更（push、tag push、remote branch delete）は行わない。完了後は `マージ準備入力` を `merge_lane` へ渡す。

## 返す成果物

- 人間向け返却: 現在成果物、着手可能成果物、停止中成果物、停止理由を短く返す。
- 着手計画返却: 今回作る成果物のリスト、各成果物の要否と理由（仕様変更、画面変更、構造変更、UI、backend 変更の有無）、新規か作り直しかの区別を人間へ返す。
- 起動先向け返却: 対象成果物、満たされた依存対象、読むファイル、禁止事項、期待する成果物を渡す。
- 設計レビュー戻し: 差し戻し対象、戻し先 agent、戻し結果、戻せない場合の停止理由を返す。
- 終了処理返却: `作業 commit` と `マージ準備入力`（active plan folder、source branch、target branch、commit hash、検証結果、実装後ブラウザ確認結果、残留リスク）を返す。

## 完了の目安

- `着手計画` が人間確認済みで、今回作る成果物の要否が固定されている。
- DAG で `成立条件` を満たした成果物が、承認済みまたは停止理由付きで揃っている。
- 不変条件（人間ゲート、責務境界、UI 順序、検証順序、安全境界）に違反がない。
- 作業 branch が `codex/<task-id>` として存在し、変更が local commit 済みである。
- `マージ準備入力` が揃い、remote repository を変更していない。

## 作業を止める条件

- 依頼が新規実装または機能拡張か判断できない。
- `着手計画` の人間確認が得られない。
- 不変条件を緩めないと先へ進めない（人間ゲート、責務境界、UI 順序、検証順序、安全境界のいずれかに抵触する）。
- 承認済み `実装範囲` なしで実装成果物が必要になる。
- 検証失敗が、担当 agent の差し戻し範囲を超える承認済み実装範囲外の変更を必要とする。
- 仕様変更があるのに `正本化判断` を固定できない、または恒久仕様承認があるのに `詳細仕様正本反映` を固定できない。
- 停止時は、不足項目、衝突箇所、固定できない判断、戻し先を返す。
