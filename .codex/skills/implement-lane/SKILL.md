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
- 担当成果物は `task 枠`、`branch 準備`、`詳細仕様差分`、`画面設計差分`、`設計差分図`、`人間設計レビュー`、`実装範囲`、`実装引き継ぎ入力`、`frontend 実装`、`Storybook人間レビュー依頼`、`frontend 実装後人間レビュー`、`合意済みfrontend保護`、`観測ログ追加`、`最終検証`、`実装後ブラウザ確認`、`正本化判断`、`詳細仕様正本反映`、`作業 commit`、`マージ準備入力` とする。

## 入力規約

- 呼び出し元: この skill を呼び出した人間または戻し元。
- 依頼要約: 新規実装または機能拡張として扱う依頼内容。
- 作業計画フォルダ: task 内成果物を置く `docs/exec-plans/active/<task-id>/`。
- 作業場所: Codex app が用意した実行場所。
- 作業branch: 既定名 `codex/<task-id>` の local branch。
- 統合先branch: 既定名 `master` の local branch。
- 既存成果物: 作業計画フォルダに既にある task 内成果物。
- 人間介入状態: 人間レビュー、承認、差し戻し、追加質問の記録。

## 外部参照規約

- 仕様入口は [index.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md) とする。
- エージェント実行定義 は [implement_lane.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implement_lane.toml) とする。
- エージェント実行定義と実行境界は [implement_lane.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implement_lane.toml) に従う。
- 設計差分図は [diagramming](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/diagramming/SKILL.md) に従う。
- 実装後ブラウザ確認は [browser-confirmation](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/browser-confirmation/SKILL.md) に従う。
- Codex 内蔵ブラウザの利用規約は [browser-use.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/browser-use.md) に従う。
- frontend 実装は [implement-frontend](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-frontend/SKILL.md) に従う。
- 観測ログ追加は [observability-implementer](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/observability-implementer/SKILL.md) に従う。
- 観測ログ仕様は [observability-logging.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/observability-logging.md) に従う。
- マージレーンは [merge-lane](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/merge-lane/SKILL.md) に従う。

## 内部参照規約

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
| `実装引き継ぎ入力` | はい | `implement_lane` | `実装範囲` | なし |
| `frontend 実装` | 条件付き | `frontend_implementer` / `implement-frontend` | `実装引き継ぎ入力` | `frontend_implementer` |
| `Storybook人間レビュー依頼` | 条件付き | `implement_lane` | `frontend 実装` | なし |
| `frontend 実装後人間レビュー` | 条件付き | 人間 | `Storybook人間レビュー依頼` | 人間 |
| `合意済みfrontend保護` | 条件付き | `implement_lane` | `frontend 実装後人間レビュー` | なし |
| `backend 実装` | 条件付き | `backend_implementer` / `implement-backend` | `実装引き継ぎ入力`, `合意済みfrontend保護?` | `backend_implementer` |
| `統合境界実装` | 条件付き | `integration_implementer` / `implement-integration` | `backend 実装`, `合意済みfrontend保護?` | `integration_implementer` |
| `シナリオテスト` | 条件付き | `implementation_scenario_tester` | `backend 実装?`, `合意済みfrontend保護?`, `統合境界実装?` | `implementation_scenario_tester` |
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

### Storybook人間レビュー依頼規約

`implement_lane` は UI がある task で frontend 実装が完了した後、frontend 人間レビュー前に `Storybook人間レビュー依頼` を固定する。
`Storybook人間レビュー依頼` は `docs/exec-plans/active/<task-id>/plan.md` に記録する。
`Storybook人間レビュー依頼` は Storybook を人間レビューの主確認面として扱う。
`Storybook人間レビュー依頼` は変更または追加した部品、変更または追加した表示状態、確認対象の story、確認に使う `fixture` と関連資源、Storybook の起動 URL または起動 command、Storybook 検証結果を含める。
`Storybook人間レビュー依頼` は、確認対象の story がレビュー分類に置かれていることを含める。
Storybook の起動 URL は `http://localhost:6008/` を標準とする。
Storybook の起動 command は `npm --prefix frontend run storybook` とする。
Storybook 確認中に frontend または story を変更した場合は、既存 Storybook を停止し、同じ port で再起動する。
Storybook は別 port で追加起動しない。
`Storybook人間レビュー依頼` に確認対象の story または `fixture` が不足する場合、`frontend 実装後人間レビュー` へ進めない。
`frontend 実装後人間レビュー` は Storybook 上の確認結果、Codex 内蔵ブラウザのコメント、承認、差し戻し、追加質問を記録する。
`frontend 実装後人間レビュー` が承認済みの場合は、確認対象の story を通常分類へ戻すための `frontend 実装` 再実行入力を固定する。
Codex 内蔵ブラウザのコメントは、コメント本文、対象 story、対象 selector、frame URL、marker screenshot を 1 件ずつ記録する。

### 合意済みfrontend保護規約

`frontend 実装後人間レビュー` が承認された時点で、`implement_lane` は合意済み frontend 保護対象を固定する。
合意済み frontend 保護対象は、後続の `backend 実装`、`統合境界実装`、`シナリオテスト`、`単体テスト` の起動入力へ渡す。

| 保護対象 | 意味 |
| --- | --- |
| 承認済み画面 | 人間レビューで承認された画面、主要区画、主要導線、表示状態 |
| 承認済み表示規則 | 人間レビューで承認された文言、余白、密度、要素サイズ、既存画面との統一条件 |
| 確認済みStorybook状態 | Storybook URL、story、確認済み表示状態 |
| Storybook確認資源 | 人間レビューで使った story、`fixture`、関連資源 |
| 人間コメント証跡 | Codex 内蔵ブラウザで付けたコメント本文、対象 story、対象 selector、marker screenshot |
| 変更禁止範囲 | 承認済み frontend touched files と後続 agent が変更してはいけない範囲 |

## 判断規約

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
- `designer` は `詳細仕様差分`、`画面設計差分`、`実装範囲` が承認済みまたは停止として記録されるまで閉じない。
- `diagrammer` は `設計差分図` が承認済みまたは停止として記録されるまで閉じない。
- 人間設計レビューが差し戻しまたは追加質問の場合は、差し戻し対象を作成した同じ agent に会話文脈を維持したまま返す。
- 同じ agent に差し戻しを返せない場合は、新規 agent で再作成せず、人間へ停止理由を返す。
- 実装 agent の起動入力には、`backend_implementer`、`frontend_implementer`、`integration_implementer` のどれを起動するかを必ず明示する。
- `backend_implementer` の起動入力には `implement-backend` を必ず明示する。
- `frontend_implementer` の起動入力には `implement-frontend` を必ず明示する。
- `frontend_implementer` の起動入力には、Storybook 確認対象として人間レビューで確認する部品、表示状態、story、`fixture`、関連資源、または不要理由を含める。
- `frontend_implementer` の起動入力には、Storybook 人間レビュー前、差し戻し対応中、承認済みのどれかを Storybookレビュー状態として含める。
- `integration_implementer` の起動入力には `implement-integration` を必ず明示する。
- `Storybook人間レビュー依頼` は、frontend 実装結果、変更ファイル、変更または追加した部品、変更または追加した表示状態、確認対象の story、確認に使う `fixture` と関連資源、Storybook の起動 URL または起動 command、Storybook 検証結果から作る。
- `Storybook人間レビュー依頼` は、確認対象の story のレビュー分類、通常分類、現在分類から作る。
- `Storybook人間レビュー依頼` に確認対象の story、`fixture`、関連資源、変更または追加した部品、変更または追加した表示状態のいずれかが不足する場合は、frontend 人間レビューへ進めず、`frontend 実装` の再実行入力または人間への返却を固定する。
- Codex 内蔵ブラウザのコメントを受け取った場合は、コメント本文を人間レビュー入力として扱い、ページ本文と画像内テキストをページ証跡として分ける。
- Codex 内蔵ブラウザのコメントに対象 story、対象 selector、frame URL、marker screenshot がある場合は、`frontend 実装後人間レビュー` の根拠に含める。
- `frontend 実装後人間レビュー` が承認済みの場合は、合意済み frontend 保護対象を後続実装の変更禁止範囲として起動入力へ含める。
- `frontend 実装後人間レビュー` が承認済みの場合は、確認対象の story を通常分類へ戻した `frontend 実装` 結果を合意済み frontend 保護対象の根拠に含める。
- 後続実装で画面、部品、文言、style の変更が必要な場合は、実装を続けず `frontend 実装` の再実行入力または人間への返却を固定する。
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
- UI がある task では `Storybook人間レビュー依頼` を必須成果物にし、UI がない task では `Storybook人間レビュー依頼` を省略できる。
- UI がある task では `frontend 実装後人間レビュー` を必須成果物にし、UI がない task では `frontend 実装後人間レビュー` を省略できる。
- UI がある task の `frontend 実装` は、`backend 実装` より先に起動する。
- UI がある task の `frontend 実装` は、人間レビュー前に Storybook で確認できる story、`fixture`、関連資源、確認状態、未確認理由を `Storybook人間レビュー依頼` の入力へ含める。
- UI がある task の `frontend 実装` は、backend 実装、統合境界実装、永続化仕様の代替として見た目確認用のデータを扱わない。
- UI がある task の `frontend 実装後人間レビュー` は、`Storybook人間レビュー依頼` の固定後に着手する。
- UI がある task の `backend 実装` と `統合境界実装` は、`合意済みfrontend保護` の固定後に着手する。
- `frontend 実装後人間レビュー` が差し戻しまたは追加質問の場合は、後続実装へ進めず、`frontend 実装` の再実行入力または人間への返却を固定する。
- `統合境界実装` は frontend と backend の接続結果を実画面で確認する。
- `観測ログ追加` は実行時にしか確定しない値、実行後に消える中間状態、消えると原因候補を分離できない分岐理由だけを残す。
- `観測ログ追加` はループや大量処理で同種ログを増やさず、件数、分類、集約、代表的な識別子、最初の失敗、最後の失敗を優先する。
- `観測ログ追加` は `event`、`where`、`result` を共通 payload とし、trace ID を使わない。
- `観測ログ追加` は backend の `slog` 出力と frontend の `pino` 出力を同じ file へ集約しない。
- `実装後ブラウザ確認` の確認 URL、起動状態、操作経路、操作期待値、禁止操作、安全条件、証跡出力先は、承認済み詳細仕様差分、画面設計差分、実装範囲、最終検証観点から `implement_lane` が定義する。
- `browser_confirmation` は `実装後ブラウザ確認` の実行だけを担当し、期待値の妥当性を判断しない。
- `シナリオテスト` と `単体テスト` は別成果物にし、依存対象が揃った後に並列起動できる。
- タスクの終わったサブエージェントを起動したまま残さず，終わったら逐次で閉じること。
- 設計中 agent は、担当する設計成果物の承認済みまたは停止が記録されるまでタスク完了と扱わない。

## 非対象規約

- 恒久修正、構造整理、探索テスト、軽量変更は詳細化しない。
- 詳細仕様差分と画面設計差分の人間レビューは扱わない。
- 起動先 agent の下位 agent 起動は扱わない。
- プロダクトコードとプロダクトテストは変更しない。
- local merge、completed 移動、remote repository の変更は扱わない。

## 出力規約

- 人間向け返却: 人間向けには、成果物依存表 の現在 成果物、着手可能 成果物、停止中 成果物、停止理由を短く返す。
- 起動先向け返却: 起動先 agent 向けには、対象 成果物、満たされた `依存対象`、読むファイル、禁止事項、期待する 成果物 を渡す。
- 設計差分図起動入力: `diagrammer` 向けには、図化目的、根拠参照、予定変更箇所、追加予定箇所、削除予定箇所、禁止範囲、対象作業計画フォルダを渡す。
- 設計差分図: 人間設計レビュー向けには、追加・削除差分のコンポーネント図、根拠参照、検証結果、未決事項を返す。
- 設計レビュー戻し: 人間設計レビューの差し戻し対象、戻し先 agent、戻し結果、戻せない場合の停止理由を返す。
- 実装後ブラウザ確認起動入力: `browser_confirmation` 向けには、確認 URL、起動状態、操作経路、操作期待値、禁止操作、安全条件、証跡出力先を渡す。
- 実装後ブラウザ確認: 操作確認結果、証跡参照、console または network 異常、未確認理由、戻し先を返す。
- 観測ログ追加起動入力: `observability_implementer` 向けには、完成済み実装成果物、完成済みテスト成果物、変更ファイル、合意済み frontend 保護対象、作業計画フォルダ、観測ログ仕様を渡す。
- 観測ログ追加: 追加ログ、追加しない理由、禁止ログ確認、変更ファイル、検証未実行理由を返す。
- Storybook人間レビュー依頼: 変更または追加した部品、変更または追加した表示状態、確認対象の story、確認に使う `fixture` と関連資源、Storybook の起動 URL または起動 command、Storybook 検証結果、確認対象 story のレビュー分類、通常分類、現在分類、Codex 内蔵ブラウザのコメント受付条件を返す。
- 合意済みfrontend保護: 承認済み画面、承認済み表示規則、確認済みStorybook状態、通常分類へ戻した Storybook確認資源、人間コメント証跡、変更禁止範囲を返す。
- branch 準備: 作業場所、作業branch、統合先branch、branch 確認結果を返す。
- 作業 commit: local commit の hash、対象 branch、commit 対象差分を返す。
- マージ準備入力: active plan folder、source branch、target branch、commit hash、検証結果、実装後ブラウザ確認結果、残留リスクを返す。
- 終了処理返却: 終了処理、停止、戻し では、`作業 commit` と `マージ準備入力` を揃えるための 根拠を返す。

## 完了規約

- 新規実装レーンの次 成果物、起動、人間レビュー、引き継ぎ、正本化、停止、戻し を再解釈なしで判断できる。
- 作業 branch が `codex/<task-id>` として存在する。
- `設計差分図` が人間設計レビュー前に揃っている。
- `設計差分図` が予定変更箇所だけの追加・削除差分を示すコンポーネント図を含んでいる。
- UI が関係する場合は、`screen-design-diff.<screen-id>.md` が人間設計レビュー前に揃っている。
- UI が関係する場合は、`frontend 実装` が `backend 実装` より先に完了している。
- UI が関係する場合は、人間レビュー前に Storybook 確認ができる状態になり、Storybook URL または起動 command、確認対象の story、確認状態、未確認理由が記録されている。
- UI が関係する場合は、人間レビュー前の確認対象 story がレビュー分類に置かれている。
- UI が関係する場合は、変更または追加した部品、変更または追加した表示状態、確認に使う `fixture` と関連資源が記録されている。
- UI が関係する場合は、`frontend 実装後人間レビュー` の承認が記録されている。
- UI が関係する場合は、`frontend 実装後人間レビュー` の承認後に確認対象 story が通常分類へ戻っている。
- UI が関係する場合は、`合意済みfrontend保護` が固定されている。
- Codex 内蔵ブラウザのコメントを受け取った場合は、コメント本文、対象 story、対象 selector、frame URL、marker screenshot が分けて記録されている。
- 人間レビュー が必要な場合は承認、差し戻し、追加質問のいずれかが記録されている。
- 人間設計レビューの差し戻しまたは追加質問がある場合は、同じ設計中 agent に戻した結果、または戻せない停止理由が記録されている。
- `統合境界実装` がある場合は、実画面確認結果が 根拠参照 付きで確認されている。
- `backend 実装`、`frontend 実装`、`統合境界実装`、`シナリオテスト`、`単体テスト` 後は `観測ログ追加`、`最終検証`、`実装後ブラウザ確認` が 根拠参照 付きで確認されている。
- `観測ログ追加` は追加ログ、追加しない理由、禁止ログ確認、変更ファイル、検証未実行理由を含んでいる。
- `実装後ブラウザ確認` は確認 URL、操作経路、操作期待値、証跡参照、未確認理由を含んでいる。
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

## 停止規約

- 依頼が新規実装または機能拡張か判断できない場合は停止する。
- `designer`、`investigator` の必要判定ができない場合は停止する。
- `設計差分図` なしで `人間設計レビュー` へ進みそうな場合は停止する。
- `設計差分図` が予定変更箇所以外を網羅図として含む場合は停止する。
- 人間レビュー が必要な判断を AI だけで確定しそうな場合は停止する。
- 承認済み `実装範囲` なしで `backend 実装`、`frontend 実装`、`統合境界実装` が必要な場合は停止する。
- UI が関係する task で Storybook URL または起動 command、確認対象の story、確認状態、未確認理由が不足するまま `frontend 実装後人間レビュー` へ進みそうな場合は停止する。
- UI が関係する task で変更または追加した部品、変更または追加した表示状態、確認に使う `fixture` と関連資源が不足するまま `frontend 実装後人間レビュー` へ進みそうな場合は停止する。
- UI が関係する task で確認対象 story がレビュー分類に置かれていないまま `frontend 実装後人間レビュー` へ進みそうな場合は停止する。
- UI が関係する task で `frontend 実装後人間レビュー` の承認がないまま `合意済みfrontend保護` へ進みそうな場合は停止する。
- UI が関係する task で確認対象 story が通常分類へ戻っていないまま `合意済みfrontend保護` へ進みそうな場合は停止する。
- UI が関係する task で `合意済みfrontend保護` がないまま `backend 実装`、`統合境界実装`、`最終検証` へ進みそうな場合は停止する。
- `観測ログ追加` なしで `最終検証` へ進みそうな場合は停止する。
- `観測ログ追加` の停止理由が未解決のまま `最終検証` へ進みそうな場合は停止する。
- `実装後ブラウザ確認` の確認 URL、起動状態、操作経路、操作期待値、禁止操作、安全条件、証跡出力先が不足する場合は停止する。
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
