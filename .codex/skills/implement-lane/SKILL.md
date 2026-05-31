---
name: implement-lane
description: 新規実装レーンで task 内成果物依存表、人間介入、引き継ぎ、終了条件を固定する作業プロトコル。
---
# Implement Lane

## 目的

`implement-lane` は、新規実装と機能拡張の進行判断を task 内成果物依存表 と 引き継ぎ へ固定する作業プロトコルである。

## 対応ロール

- `implement_lane` が使う。
- 呼び出し元は人間とする。
- 返却先は人間とする。
- 担当成果物は `task 枠`、`branch 準備`、`詳細仕様差分`、`画面設計差分`、`設計差分図`、`人間設計レビュー`、`実装範囲`、`テスト設計`、`実装引き継ぎ入力`、`frontend 実装`、`Storybookレビューループ入力確認`、`frontend 実装後人間レビュー`、`Storybook後画面設計差分整合`、`合意済みfrontend保護`、`観測ログ追加`、`最終検証`、`実装後ブラウザ確認`、`正本化判断`、`詳細仕様正本反映`、`作業 commit`、`マージ準備入力` とする。

## 呼び出し元から渡される情報

- 呼び出し元: この skill を呼び出した人間または戻し元。
- 依頼要約: 新規実装または機能拡張として扱う依頼内容。
- 作業計画フォルダ: task 内成果物を置く `docs/exec-plans/active/<task-id>/`。
- 作業場所: Codex app が用意した実行場所。
- 作業branch: 既定名 `codex/<task-id>` の local branch。
- 統合先branch: 既定名 `master` の local branch。
- 既存成果物: 作業計画フォルダに既にある task 内成果物。
- 人間介入状態: 人間レビュー、承認、差し戻し、追加質問の記録。

## 作業前に読む正本

- 仕様入口は [index.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md) とする。
- エージェント実行定義 は [implement_lane.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implement_lane.toml) とする。
- エージェント実行定義と実行境界は [implement_lane.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implement_lane.toml) に従う。
- 設計差分図は [diagramming](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/diagramming/SKILL.md) に従う。
- テスト設計は [test-design](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/test-design/SKILL.md) に従う。
- 実装後ブラウザ確認は [browser-confirmation](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/browser-confirmation/SKILL.md) に従う。
- Chrome DevTools MCP の利用規約は [chrome-devtools-mcp.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/chrome-devtools-mcp.md) に従う。
- Storybook レビューループは [story-book-review-loop](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/story-book-review-loop/SKILL.md) に従う。
- frontend 実装は [implement-frontend](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-frontend/SKILL.md) に従う。
- 観測ログ追加は [observability-implementer](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/observability-implementer/SKILL.md) に従う。
- 観測ログ仕様は [observability-logging.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/observability-logging.md) に従う。
- サンドボックス外実行の許可は [default.rules](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/rules/default.rules) に従う。
- マージレーンは [merge-lane](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/merge-lane/SKILL.md) に従う。

## skill 内の拘束条件

新規実装レーンの 成果物依存表 は次を必ず持つ。
各 成果物 は、`依存対象` の 成果物 が揃った時だけ着手できる。
`次 agent` は、その 成果物 を揃えるために 引き継ぎ入力 を渡す相手を示す。
`次 agent` が複数ある行は、依存対象が満たされ、ツール権限 が衝突しない場合に並列 起動 できる候補を示す。
当スキルは，この`次エージェント`をコンテキスト継承なしでサブエージェントとしてスポーンすることでDAGの成果物を作っていく。
設計中 agent は、担当する設計成果物が承認済みまたは停止として記録されるまで閉じない。
人間設計レビューの差し戻しは、新規 agent を起動せず、会話文脈を維持した同じ agent へ戻す。
設計中 agent は、`designer` と `diagrammer` を指す。

| 成果物ID | 必須 | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- | --- |
| `task 枠` | はい | `implement_lane` | `[]` | なし |
| `branch 準備` | はい | `implement_lane` | `task 枠` | なし |
| `詳細仕様差分` | はい | `designer` | `task 枠` | `designer` |
| `画面設計差分` | 条件付き | `designer` | `詳細仕様差分` | `designer` |
| `設計差分図` | はい | `diagrammer` | `詳細仕様差分`, `画面設計差分?` | `diagrammer` |
| `人間設計レビュー` | はい | 人間 | `詳細仕様差分`, `画面設計差分?`, `設計差分図` | 人間 |
| `実装範囲` | はい | `designer` | `人間設計レビュー` | `designer` |
| `テスト設計` | はい | `test_designer` | `人間設計レビュー` | `test_designer` |
| `実装引き継ぎ入力` | はい | `implement_lane` | `実装範囲` | なし |
| `frontend 実装` | 条件付き | `frontend_implementer` / `implement-frontend` | `実装引き継ぎ入力` | `frontend_implementer` |
| `Storybookレビューループ入力確認` | 条件付き | `implement_lane` | `frontend 実装` | なし |
| `Storybookレビューループ完了証跡` | 条件付き | 人間が立てた別セッション / `story-book-review-loop` | `Storybookレビューループ入力確認` | なし |
| `frontend 実装後人間レビュー` | 条件付き | `implement_lane` | `Storybookレビューループ完了証跡` | なし |
| `Storybook後画面設計差分整合` | Storybook レビューループ後に画面仕様変更あり | `designer` | `frontend 実装後人間レビュー` | `designer` |
| `合意済みfrontend保護` | 条件付き | `implement_lane` | `frontend 実装後人間レビュー`, `Storybook後画面設計差分整合?` | なし |
| `backend 実装` | 条件付き | `backend_implementer` / `implement-backend` | `実装引き継ぎ入力`, `合意済みfrontend保護?` | `backend_implementer` |
| `統合境界実装` | 条件付き | `integration_implementer` / `implement-integration` | `backend 実装`, `合意済みfrontend保護?` | `integration_implementer` |
| `シナリオテスト` | 条件付き | `implementation_scenario_tester` | `テスト設計`, `backend 実装?`, `合意済みfrontend保護?`, `統合境界実装?` | `implementation_scenario_tester` |
| `単体テスト` | 条件付き | `implementation_unit_tester` | `backend 実装?`, `合意済みfrontend保護?`, `統合境界実装?` | `implementation_unit_tester` |
| `観測ログ追加` | はい | `observability_implementer` / `observability-implementer` | `backend 実装?`, `frontend 実装?`, `合意済みfrontend保護?`, `統合境界実装?`, `シナリオテスト?`, `単体テスト?` | `observability_implementer` |
| `最終検証` | 条件付き | `implement_lane` | `観測ログ追加` | なし |
| `実装後ブラウザ確認` | はい | `browser_confirmation` | `最終検証` | `browser_confirmation` |
| `正本化判断` | 仕様変更または仕様追加あり | `implement_lane` | `最終検証`, `実装後ブラウザ確認` | `docs_updater?` |
| `詳細仕様正本反映` | 仕様変更または仕様追加あり | `docs_updater` | `正本化判断` | `docs_updater?` |
| `作業 commit` | はい | `implement_lane` | 全完了または停止済み 成果物 | なし |
| `マージ準備入力` | はい | `implement_lane` | `作業 commit` | `merge_lane` |

検証証跡 は次をすべて含む。

- 実行コマンド: 呼び出し元が実行した検証コマンド。
- 証跡位置: 実行日時または run 内の証跡位置。
- 成否: pass または fail。
- coverage 値: coverage を測定した場合の値。
- issue 数: security、reliability、maintainability の issue 数。
- system test 件数: system test の実行件数、成功件数、失敗件数。
- 失敗箇所: fail の場合に原因箇所または失敗した検証名。

### Storybookレビューループ入力確認規約

`implement_lane` は UI がある task で frontend 実装が完了した後、Storybook レビューループへ渡せる入力が揃っているか確認する。
`Storybookレビューループ入力確認` はレビュー資料ではない。
`Storybookレビューループ入力確認` は `docs/exec-plans/active/<task-id>/plan.md` に別セッションへ渡す入力の所在と不足項目だけを記録する。
`Storybookレビューループ入力確認` は、作業計画フォルダ、frontend 実装結果、frontend 実装境界だけを対象にする。
Storybook の起動、分類、確認資源、`fixture` の妥当性は扱わない。
作業計画フォルダ、frontend 実装結果、frontend 実装境界の所在が不足する場合、`frontend 実装後人間レビュー` へ進めない。
Storybook レビューループは、人間が立てた別セッションで `story-book-review-loop` に従って実行する。
`implement_lane` は Storybook レビューループを起動または直接実行しない。
`implement_lane` は `Storybookレビューループ入力確認` を固定した時点で停止し、人間が立てる別セッションへ渡す入力の所在と不足項目を返す。
`implement_lane` は Storybook レビューループ中の Chrome DevTools MCP 操作、コメント収集、コメント解釈、frontend 修正、修正結果判定を扱わない。
`Storybookレビューループ完了証跡` は、作業計画フォルダの `storybook-review-loop.md` が存在し、確定した story、変更された画面仕様、反映先、承認状態を持つ状態を指す。
`frontend 実装後人間レビュー` は、作業計画フォルダの `storybook-review-loop.md` から承認状態、frontend レビュー修正成果物、Storybook レビューループ画面仕様だけを記録する。
`frontend 実装後人間レビュー` が承認済みの場合は、`storybook-review-loop.md` に記録された確定済み story と分類を記録する。
Storybook レビューループ後に画面仕様が変わった場合は、`designer` に戻して plan 内の `screen-design-diff.<screen-id>.md` など該当する画面設計成果物を更新させる。

### Storybook後画面設計差分整合規約

`Storybook後画面設計差分整合` は、Storybook レビューループ後に画面仕様が変わった場合だけ必要な成果物である。
`implement_lane` は `storybook-review-loop.md` の変更された画面仕様、反映先、承認状態を根拠に `designer` を起動する。
`designer` は plan 内の該当する `screen-design-diff.<screen-id>.md` へ変更後の画面仕様を反映する。
`implement_lane` は `screen-design-diff.<screen-id>.md` の本文を直接作成または更新しない。
`designer` を起動できる入力が揃っている場合、`implement_lane` は人間への停止返却だけで終わらず `designer` を起動する。
`designer` を起動できない場合、`implement_lane` は不足項目、衝突箇所、戻し先を記録して停止する。

### 合意済みfrontend保護規約

`frontend 実装後人間レビュー` が承認された時点で、`implement_lane` は合意済み frontend 保護対象を固定する。
合意済み frontend 保護対象は、後続の `backend 実装`、`統合境界実装`、`シナリオテスト`、`単体テスト` の起動入力へ渡す。

| 保護対象 | 意味 |
| --- | --- |
| 承認済み画面 | 人間レビューで承認された画面、主要区画、主要導線、表示状態 |
| 承認済み表示規則 | 人間レビューで承認された文言、余白、密度、要素サイズ、既存画面との統一条件 |
| 確認済みStorybook状態 | Storybook URL、story、確認済み表示状態 |
| Storybook確認資源 | 人間レビューで使った story、`fixture`、関連資源 |
| Storybook画面仕様 | Storybook レビューループで確定した変更後の画面仕様 |
| 変更禁止範囲 | 承認済み frontend touched files と後続 agent が変更してはいけない範囲 |

### 成果物編集主体規約

`implement_lane` は、`docs/exec-plans/active/<task-id>/` の中でも、担当者が `implement_lane` の成果物だけを直接作成または更新する。
担当者が `designer`、`diagrammer`、`test_designer`、`docs_updater`、実装 agent、テスト agent、`browser_confirmation` の成果物は、担当 agent の返却結果または人間介入状態を記録する場合だけ参照または転記する。
担当 agent の成果物本文を `implement_lane` が代筆しない。

| 成果物分類 | `implement_lane` の扱い |
| --- | --- |
| `task 枠` | 直接作成または更新できる。 |
| `branch 準備` | branch 確認結果を直接記録できる。 |
| `詳細仕様差分` | `designer` の担当成果物として扱い、本文を直接作成または更新しない。 |
| `画面設計差分` | `designer` の担当成果物として扱い、本文を直接作成または更新しない。 |
| `設計差分図` | `diagrammer` の担当成果物として扱い、本文を直接作成または更新しない。 |
| `人間設計レビュー` | 人間の承認、差し戻し、追加質問を記録できる。 |
| `実装範囲` | `designer` の担当成果物として扱い、本文を直接作成または更新しない。 |
| `テスト設計` | `test_designer` の担当成果物として扱い、本文を直接作成または更新しない。 |
| `実装引き継ぎ入力` | 承認済み `実装範囲` から直接作成または更新できる。 |
| `Storybookレビューループ入力確認` | frontend 実装結果と既存成果物から入力の所在と不足項目だけを記録できる。 |
| `Storybookレビューループ完了証跡` | 作業計画フォルダの `storybook-review-loop.md` の存在と完了状態だけを確認できる。 |
| `frontend 実装後人間レビュー` | `storybook-review-loop.md` の完了状態を記録できる。 |
| `Storybook後画面設計差分整合` | `designer` の担当成果物として扱い、本文を直接作成または更新しない。 |
| `合意済みfrontend保護` | 承認済み frontend レビュー結果から直接作成または更新できる。 |
| `観測ログ追加` | `observability_implementer` の返却結果を記録できる。 |
| `最終検証` | 実行結果を直接記録できる。 |
| `実装後ブラウザ確認` | `browser_confirmation` の返却結果を記録できる。 |
| `正本化判断` | 最終検証と実装後ブラウザ確認の結果から直接作成または更新できる。 |
| `詳細仕様正本反映` | `docs_updater` の担当成果物として扱い、本文を直接作成または更新しない。 |
| `作業 commit` | local commit 結果を直接記録できる。 |
| `マージ準備入力` | 作業 commit と検証結果から直接作成または更新できる。 |

## 担当ロールが判断してよい範囲

- 次の実行判断は 成果物依存表 の未完了 成果物、満たされた `依存対象`、既存 成果物、対象 skill の完了規約で決める。
- 既存 成果物 がある場合は、対象 skill の完了規約を満たすか確認してから後続 成果物 へ進む。
- `branch 準備` は active plan ごとの `codex/<task-id>` の local branch を作成または確認する。
- 作業 branch を特定できない場合は、後続成果物へ進めない。
- 起動先 agent の 起動入力 は、対象 skill の入力規約、完了規約、停止規約に合わせて作る。
- `設計差分図` は `diagrammer` を起動して作る。
- `設計差分図` の標準形は、予定変更箇所だけの追加・削除差分を示すコンポーネント図に限定する。
- シーケンス図、状態遷移図、その他の図は、ユーザー要求または複雑性がある場合だけ作る。
- `設計差分図` は、全体構成図、正本図、変更しない箇所の網羅図として作らない。
- `設計差分図` の起動入力には、詳細仕様差分、画面設計差分がある場合の画面設計差分、予定変更箇所、追加予定箇所、削除予定箇所、禁止範囲、出力先を含める。
- `詳細仕様差分` または `画面設計差分` が必要な場合は、`designer` の起動入力を作る。
- `designer` を起動できない場合は、`detail-spec-diff.md` と `screen-design-diff.<screen-id>.md` を `implement_lane` が代筆せず、人間へ停止理由を返す。
- 既存の `detail-spec-diff.md` または `screen-design-diff.<screen-id>.md` が不足している場合は、`designer` への戻し入力または人間への停止理由を固定する。
- `designer` は `詳細仕様差分`、`画面設計差分`、`実装範囲` が承認済みまたは停止として記録されるまで閉じない。
- `diagrammer` は `設計差分図` が承認済みまたは停止として記録されるまで閉じない。
- `テスト設計` は `人間設計レビュー` 後、`実装範囲` と並列起動できる成果物として扱う。
- `テスト設計` の起動入力には、作業計画フォルダ、承認済み詳細仕様差分、承認済み画面設計差分、関連ユースケース、承認記録、出力先 `test-design.csv` を含める。
- `テスト設計` を起動できない場合は、`test-design.csv` を `implement_lane` が代筆せず、人間へ停止理由を返す。
- `シナリオテスト` の起動入力には、`test-design.csv` を根拠として含める。
- 人間設計レビューが差し戻しまたは追加質問の場合は、差し戻し対象を作成した同じ agent に会話文脈を維持したまま返す。
- 同じ agent に差し戻しを返せない場合は、新規 agent で再作成せず、人間へ停止理由を返す。
- 実装 agent の起動入力には、`backend_implementer`、`frontend_implementer`、`integration_implementer` のどれを起動するかを必ず明示する。
- `backend_implementer` の起動入力には `implement-backend` を必ず明示する。
- `frontend_implementer` の起動入力には `implement-frontend` を必ず明示する。
- `frontend_implementer` の起動入力には、Storybook 確認対象として人間レビューで確認するコンポーネント、画面、表示状態、story、`fixture`、関連資源、または不要理由を含める。
- `frontend_implementer` の起動入力には、Storybook 人間レビュー前、差し戻し対応中、承認済みのどれかを Storybookレビュー状態として含める。
- `integration_implementer` の起動入力には `implement-integration` を必ず明示する。
- `Storybookレビューループ入力確認` は、作業計画フォルダ、frontend 実装結果、frontend 実装境界の所在に不足がないかだけを確認する。
- `Storybookレビューループ入力確認` は、Storybook の起動状態、分類、確認資源、`fixture` の妥当性を判断しない。
- `Storybookレビューループ入力確認` に作業計画フォルダ、frontend 実装結果、frontend 実装境界のいずれかが不足する場合は、frontend 人間レビューへ進めず、人間への返却を固定する。
- Storybook レビューループは人間が立てた別セッションで `story-book-review-loop` に従って実行するため、`implement_lane` は `Storybookレビューループ入力確認` の固定後に停止する。
- `implement_lane` は Storybook レビューループ中の Chrome DevTools MCP 操作、コメント収集、コメント解釈、frontend 修正、修正結果判定を行わない。
- 作業計画フォルダに `storybook-review-loop.md` が出来上がっている場合だけ、変更された画面仕様、反映先、現在分類、承認状態を `frontend 実装後人間レビュー` の根拠に含める。
- `story-book-review-loop` から設計整合入力が返った場合は、`designer` へ戻し、plan 内の `screen-design-diff.<screen-id>.md` など該当する画面設計成果物を更新させる。
- Storybook レビューループ後に画面仕様が変わり、`designer` を起動できる入力が揃っている場合は、`Storybook後画面設計差分整合` を着手可能成果物として扱う。
- Storybook レビューループ後に画面仕様が変わり、`designer` を起動できる入力が不足する場合は、`合意済みfrontend保護` へ進めず不足項目を固定する。
- `frontend 実装後人間レビュー` が承認済みの場合は、合意済み frontend 保護対象を後続実装の変更禁止範囲として起動入力へ含める。
- `frontend 実装後人間レビュー` が承認済みの場合は、`storybook-review-loop.md` に記録された確定済み story と分類を合意済み frontend 保護対象の根拠に含める。
- 後続実装で画面、コンポーネント、文言、style の変更が必要な場合は、実装を続けず `frontend 実装` の再実行入力または人間への返却を固定する。
- `観測ログ追加` の起動入力には、完成済み実装成果物、完成済みテスト成果物、変更ファイル、合意済み frontend 保護対象、作業計画フォルダ、観測ログ仕様を含める。
- `観測ログ追加` は `backend 実装`、`frontend 実装`、`統合境界実装`、`シナリオテスト`、`単体テスト` が必要分だけ揃った後、`最終検証` の前に起動する。
- `観測ログ追加` が停止した場合は、`最終検証` へ進めず、人間または該当 実装 agent への戻しを固定する。
- frontend 検証失敗の再実行入力は、承認済み frontend 実装範囲、generated file、生成元、公開境界、今回変更が直接壊した frontend 責務内プロダクトコードに限り、`frontend_implementer` へ渡す。
- シナリオテスト検証失敗の再実行入力は、承認済み詳細仕様差分、承認済み実装範囲、軽量変更レーンの `task 枠`、今回のテスト変更が直接壊したテスト補助、fixture、検証経路、担当シナリオテスト成果物に限り、`implementation_scenario_tester` へ渡す。
- 単体テスト検証失敗の再実行入力は、仕様根拠、承認済み実装範囲、軽量変更レーンの `task 枠`、今回のテスト変更が直接壊したテスト補助、fixture、検証経路、担当単体テスト成果物に限り、`implementation_unit_tester` へ渡す。
- generated file が検証失敗原因である場合は、generated file の直接編集ではなく、生成元または公開境界の修正として担当 agent に渡す。
- UI 表示、画面文言、layout、style、承認済み画面設計差分を越える変更が必要な場合は、実装を続けず人間への返却を固定する。
- 完了判断をする場合は、全必須成果物と残留リスクを確認した後に local commit を作る。
- `マージ準備入力` は active plan folder、source branch、target branch、commit hash、検証結果、実装後ブラウザ確認結果、残留リスクを含める。
- 作業計画フォルダの completed 移動と local merge は `merge_lane` に渡す。
- `detail-spec-diff.md`、`screen-design-diff.<screen-id>.md`、実装結果のいずれかに仕様変更または仕様追加が少しでも含まれる場合は、`正本化判断` を必須成果物にする。
- 仕様変更または仕様追加が human 承認済みの恒久仕様である場合は、`詳細仕様正本反映` を必須成果物にする。
- `詳細仕様正本反映` は `docs/detail-specs/` の詳細仕様正本へ、human 承認済みの恒久仕様だけを反映する。
- `詳細仕様正本反映` の入力は、`detail-spec-diff.md`、`screen-design-diff.<screen-id>.md`、実装結果、承認記録のうち正本化判断で承認済みとされた成果物に限定する。
- 起動先 agent には 文脈 を引き継がず、必要情報を 引き継ぎ入力 に明示する。
- 人間介入 が必要な 成果物 は AI だけで完了にしない。
- 恒久修正、構造整理、探索テスト、軽量変更はこの skill で詳細化しない。
- backend、frontend、統合境界 は別 成果物 として扱い、単一の実装成果物に束ねない。
- UI がある task では `frontend 実装` を必須成果物にし、UI がない task では `frontend 実装` を省略できる。
- UI がある task では `Storybookレビューループ入力確認` を必須成果物にし、UI がない task では `Storybookレビューループ入力確認` を省略できる。
- UI がある task では `frontend 実装後人間レビュー` を必須成果物にし、UI がない task では `frontend 実装後人間レビュー` を省略できる。
- UI がある task の `frontend 実装` は、`backend 実装` より先に起動する。
- UI がある task の `frontend 実装` は、人間が別セッションで Storybook レビューループを実行できるように、作業計画フォルダと frontend 実装結果の所在を `Storybookレビューループ入力確認` の対象へ含める。
- UI がある task の `frontend 実装` は、backend 実装、統合境界実装、永続化仕様の代替として見た目確認用のデータを扱わない。
- UI がある task の `frontend 実装後人間レビュー` は、`Storybookレビューループ入力確認` を経て作業計画フォルダに `storybook-review-loop.md` が出来上がった後に記録する。
- UI がある task の `backend 実装` と `統合境界実装` は、`合意済みfrontend保護` の固定後に着手する。
- `frontend 実装後人間レビュー` が差し戻しまたは追加質問の場合は、後続実装へ進めず、人間が立てた別セッションで Storybook レビューループを続けるために停止する。
- Storybook レビューループ後の UI 変更が画面設計差分と不整合な場合は、後続実装へ進めず `designer` へ戻して plan 内の画面設計成果物を更新させる。
- `統合境界実装` は frontend と backend の接続結果を実画面で確認する。
- `観測ログ追加` は実行時にしか確定しない値、実行後に消える中間状態、消えると原因候補を分離できない分岐理由だけを残す。
- `観測ログ追加` はループや大量処理で同種ログを増やさず、件数、分類、集約、代表的な識別子、最初の失敗、最後の失敗を優先する。
- `観測ログ追加` は `event`、`where`、`result` を共通 payload とし、trace ID を使わない。
- `観測ログ追加` は backend の `slog` 出力と frontend の `pino` 出力を同じ file へ集約しない。
- `実装後ブラウザ確認` の確認 URL、起動状態、操作経路、操作期待値、禁止操作、安全条件、証跡出力先は、承認済み詳細仕様差分、画面設計差分、実装範囲、最終検証観点から `implement_lane` が定義する。
- `実装後ブラウザ確認` の起動状態は、確認 URL を開くための app または server が listen していること、または HTTP 到達できることを指す。
- `実装後ブラウザ確認` の app または server 起動 command が `default.rules` の許可対象である場合は、`implement_lane` が `browser_confirmation` へ渡す前に起動状態を用意する。
- `default.rules` の許可対象 command が通常のサンドボックス実行で起動不能になった場合は、`implement_lane` が同じ command をサンドボックス外実行で再実行してから起動不能を判断する。
- `実装後ブラウザ確認` の起動不能を記録する場合は、通常のサンドボックス実行結果、サンドボックス外実行結果、listen または HTTP 到達確認結果を根拠に含める。
- `browser_confirmation` は `実装後ブラウザ確認` の実行だけを担当し、期待値の妥当性を判断しない。
- `シナリオテスト` と `単体テスト` は別成果物にし、依存対象が揃った後に並列起動できる。
- タスクの終わったサブエージェントを起動したまま残さず，終わったら逐次で閉じること。
- 設計中 agent は、担当する設計成果物の承認済みまたは停止が記録されるまでタスク完了と扱わない。

## skill が扱わない対象

- 恒久修正、構造整理、探索テスト、軽量変更は詳細化しない。
- 詳細仕様差分と画面設計差分の人間レビューは扱わない。
- Storybook レビューループの起動と直接実行は扱わない。
- Storybook レビューループ中の Chrome DevTools MCP 操作、コメント収集、コメント解釈、frontend 修正、修正結果判定は扱わない。
- 起動先 agent の下位 agent 起動は扱わない。
- プロダクトコードとプロダクトテストは変更しない。
- local merge、completed 移動、remote repository の変更は扱わない。

## 返す成果物

- 人間向け返却: 人間向けには、成果物依存表 の現在 成果物、着手可能 成果物、停止中 成果物、停止理由を短く返す。
- 起動先向け返却: 起動先 agent 向けには、対象 成果物、満たされた `依存対象`、読むファイル、禁止事項、期待する 成果物 を渡す。
- 設計差分図起動入力: `diagrammer` 向けには、図化目的、根拠参照、予定変更箇所、追加予定箇所、削除予定箇所、禁止範囲、対象作業計画フォルダを渡す。
- 設計差分図: 人間設計レビュー向けには、追加・削除差分のコンポーネント図、根拠参照、検証結果、未決事項を返す。
- 設計レビュー戻し: 人間設計レビューの差し戻し対象、戻し先 agent、戻し結果、戻せない場合の停止理由を返す。
- 実装後ブラウザ確認起動入力: `browser_confirmation` 向けには、確認 URL、起動 command、起動状態、listen または HTTP 到達確認結果、操作経路、操作期待値、禁止操作、安全条件、証跡出力先を渡す。
- 実装後ブラウザ確認: 操作確認結果、証跡参照、console または network 異常、未確認理由、戻し先を返す。
- 観測ログ追加起動入力: `observability_implementer` 向けには、完成済み実装成果物、完成済みテスト成果物、変更ファイル、合意済み frontend 保護対象、作業計画フォルダ、観測ログ仕様を渡す。
- 観測ログ追加: 追加ログ、追加しない理由、禁止ログ確認、変更ファイル、検証未実行理由を返す。
- Storybookレビューループ入力確認: 作業計画フォルダ、frontend 実装結果、frontend 実装境界の所在と不足項目を返す。
- Storybookレビューループ完了証跡: 作業計画フォルダの `storybook-review-loop.md` の有無、承認状態、Storybook レビューループ画面仕様、frontend レビュー修正成果物、設計整合入力を返す。
- Storybook後画面設計差分整合: `designer` の起動入力、更新対象の `screen-design-diff.<screen-id>.md`、更新結果、起動不能理由を返す。
- テスト設計起動入力: `test_designer` 向けには、作業計画フォルダ、承認済み詳細仕様差分、承認済み画面設計差分、関連ユースケース、承認記録、出力先を渡す。
- テスト設計: `test-design.csv`、根拠参照、不足情報、未決 selector、戻し先を返す。
- 合意済みfrontend保護: 承認済み画面、承認済み表示規則、確認済みStorybook状態、通常分類へ戻した Storybook確認資源、Storybook画面仕様、変更禁止範囲を返す。
- branch 準備: 作業場所、作業branch、統合先branch、branch 確認結果を返す。
- 作業 commit: local commit の hash、対象 branch、commit 対象差分を返す。
- マージ準備入力: active plan folder、source branch、target branch、commit hash、検証結果、実装後ブラウザ確認結果、残留リスクを返す。
- 終了処理返却: 終了処理、停止、戻し では、`作業 commit` と `マージ準備入力` を揃えるための 根拠を返す。

## 作業を完了できる条件

- 新規実装レーンの次 成果物、起動、人間レビュー、引き継ぎ、正本化、停止、戻し を再解釈なしで判断できる。
- 作業 branch が `codex/<task-id>` として存在する。
- `設計差分図` が人間設計レビュー前に揃っている。
- `設計差分図` が予定変更箇所だけの追加・削除差分を示すコンポーネント図を含んでいる。
- UI が関係する場合は、`screen-design-diff.<screen-id>.md` が人間設計レビュー前に揃っている。
- UI が関係する場合は、`frontend 実装` が `backend 実装` より先に完了している。
- UI が関係する場合は、`Storybookレビューループ入力確認` に作業計画フォルダ、frontend 実装結果、frontend 実装境界の所在が記録されている。
- UI が関係する場合は、`frontend 実装後人間レビュー` の承認が記録されている。
- UI が関係する場合は、`frontend 実装後人間レビュー` の承認後に確認対象 story が通常分類へ戻っている。
- UI が関係する場合は、作業計画フォルダに `storybook-review-loop.md` が存在し、Storybook レビューループで変更された画面仕様が記録されている。
- UI が関係する場合は、Storybook レビューループ後の画面仕様変更が plan 内の該当する画面設計成果物へ反映されている。
- UI が関係し、Storybook レビューループ後に画面仕様変更がある場合は、`Storybook後画面設計差分整合` の完了結果または停止理由が記録されている。
- UI が関係する場合は、`合意済みfrontend保護` が固定されている。
- 人間レビュー が必要な場合は承認、差し戻し、追加質問のいずれかが記録されている。
- `テスト設計` が完了しているか、停止理由が記録されている。
- 人間設計レビューの差し戻しまたは追加質問がある場合は、同じ設計中 agent に戻した結果、または戻せない停止理由が記録されている。
- `統合境界実装` がある場合は、実画面確認結果が 根拠参照 付きで確認されている。
- `backend 実装`、`frontend 実装`、`統合境界実装`、`シナリオテスト`、`単体テスト` 後は `観測ログ追加`、`最終検証`、`実装後ブラウザ確認` が 根拠参照 付きで確認されている。
- `観測ログ追加` は追加ログ、追加しない理由、禁止ログ確認、変更ファイル、検証未実行理由を含んでいる。
- `実装後ブラウザ確認` は確認 URL、起動 command、起動状態、操作経路、操作期待値、証跡参照、未確認理由を含んでいる。
- `実装後ブラウザ確認` の起動状態は、listen または HTTP 到達確認結果を根拠参照付きで含んでいる。
- 影響範囲修正 がある場合は、対象、理由、変更結果、再検証結果が 根拠参照 付きで記録されている。
- DAGで必須とされている成果物が全て用意できていること。
- 仕様変更または仕様追加がある場合は、`正本化判断` の結果が 根拠参照 付きで記録されている。
- human 承認済みの恒久仕様がある場合は、`詳細仕様正本反映` の完了結果または停止理由が 根拠参照 付きで記録されている。
- `backend 実装` またはテスト変更に backend 変更が含まれる場合は `python3 scripts/harness/run.py --suite backend-local` を `.codex/rules/default.rules` の許可対象として実行し、失敗時は担当 agent がその場で直して再実行した通過結果または未実行理由が確認されている。
- `frontend 実装` またはテスト変更に frontend 変更が含まれる場合は `python3 scripts/harness/run.py --suite frontend-local` を `.codex/rules/default.rules` の許可対象として実行し、失敗時は担当 agent がその場で直して再実行した通過結果または未実行理由が確認されている。
- 終了処理、停止、戻し のいずれでも 作業 commit と マージ準備入力 を判断できる根拠が作成されている。
- 完了判断済みの場合は、変更が local commit 済みである。
- 完了判断済みの場合は、`マージ準備入力` が active plan folder、source branch、target branch、commit hash、検証結果、実装後ブラウザ確認結果、残留リスクを含んでいる。
- remote repository を変更する command を実行していない。

## 作業を止める条件

- 依頼が新規実装または機能拡張か判断できない場合は停止する。
- `designer`、`investigator` の必要判定ができない場合は停止する。
- `設計差分図` なしで `人間設計レビュー` へ進みそうな場合は停止する。
- `設計差分図` が予定変更箇所以外を網羅図として含む場合は停止する。
- 人間レビュー が必要な判断を AI だけで確定しそうな場合は停止する。
- `designer` 担当成果物の本文を `implement_lane` が直接作成または更新しそうな場合は停止する。
- `diagrammer` 担当成果物の本文を `implement_lane` が直接作成または更新しそうな場合は停止する。
- `test_designer` 担当成果物の本文を `implement_lane` が直接作成または更新しそうな場合は停止する。
- `docs_updater` 担当成果物の本文を `implement_lane` が直接作成または更新しそうな場合は停止する。
- 承認済み `実装範囲` なしで `backend 実装`、`frontend 実装`、`統合境界実装` が必要な場合は停止する。
- UI が関係する task で作業計画フォルダ、frontend 実装結果、frontend 実装境界の所在が不足するまま `frontend 実装後人間レビュー` へ進みそうな場合は停止する。
- UI が関係する task で `Storybookレビューループ入力確認` を固定した後に、同じ `implement_lane` セッション内で Storybook レビューループを起動または実行しそうな場合は停止する。
- UI が関係する task で作業計画フォルダに `storybook-review-loop.md` が出来上がっていないまま後続成果物へ進みそうな場合は停止する。
- UI が関係する task で `frontend 実装後人間レビュー` の承認がないまま `合意済みfrontend保護` へ進みそうな場合は停止する。
- UI が関係する task で確認対象 story が通常分類へ戻っていないまま `合意済みfrontend保護` へ進みそうな場合は停止する。
- UI が関係する task で Storybook レビューループ後の画面仕様変更が plan 内の該当する画面設計成果物へ反映されず、`Storybook後画面設計差分整合` の完了結果または起動不能理由もないまま `合意済みfrontend保護` へ進みそうな場合は停止する。
- UI が関係する task で `合意済みfrontend保護` がないまま `backend 実装`、`統合境界実装`、`最終検証` へ進みそうな場合は停止する。
- `観測ログ追加` なしで `最終検証` へ進みそうな場合は停止する。
- `観測ログ追加` の停止理由が未解決のまま `最終検証` へ進みそうな場合は停止する。
- `実装後ブラウザ確認` の確認 URL、起動状態、操作経路、操作期待値、禁止操作、安全条件、証跡出力先が不足する場合は停止する。
- `default.rules` の許可対象 command をサンドボックス外実行で再実行しないまま、サンドボックス起因の起動不能として `実装後ブラウザ確認` を停止しそうな場合は停止する。
- `実装後ブラウザ確認` の app または server 起動 command が `default.rules` の許可対象ではなく、人間承認もない場合は停止する。
- `python3 scripts/harness/run.py --suite all` の失敗原因が、対象 agent の許可された 影響範囲修正 に該当しない承認済み実装範囲 外の変更を必要とする場合は停止する。
- `python3 scripts/harness/run.py --suite backend-local` または `python3 scripts/harness/run.py --suite frontend-local` の失敗原因が、対象 agent の許可された 影響範囲修正 に該当しない承認済み実装範囲 外の変更を必要とする場合は停止する。
- 検証失敗原因が generated file にある場合は、generated file の直接編集を求めず、生成元または公開境界の修正として担当 agent に渡す。
- UI 表示、画面文言、layout、style、承認済み画面設計差分を越える変更が必要な場合は停止し、人間へ返す。
- 最終検証 または `実装後ブラウザ確認` が不明なまま正本化が必要な場合は停止する。
- 仕様変更または仕様追加があるのに `正本化判断` が不足する場合は終了不可とする。
- human 承認済みの恒久仕様があるのに `詳細仕様正本反映` が不足する場合は終了不可とする。
- 完了判断済みの状態で local commit を作成できない場合は終了不可とする。
- `マージ準備入力` が不足する場合は終了不可とする。
- local merge または completed 移動を実行しそうな場合は停止する。
- `push`、tag push、remote branch delete など remote repository を変更しそうな場合は停止する。
- 作業 commit または マージ準備入力 の根拠が不足する場合は終了不可とする。
