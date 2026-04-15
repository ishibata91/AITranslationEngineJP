# Design: scenario

## Goal

- Scenario 一覧を task-local artifact に固定する
- human review と実装者 handoff の両方で、状態遷移と acceptance check を追えるようにする

## 必須の書き方

- 正常系だけでなく、最低限の失敗系と再実行系を含める
- 各 scenario は開始条件、操作、期待結果、失敗時、再実行条件を読める粒度にする
- 状態遷移と acceptance check を接続する
- 固有名詞、既存 field 名、既存 contract 名、mode 名を除き、日本語優先で書く

## Rules

- path は `docs/exec-plans/active/<task-id>.scenario.md` に固定する
- 最終適用先は `docs/scenario-tests/<topic-id>.md` として plan に記録する
- `canonicalization_targets` に scenario 正本の反映先を追加する
- 実装前に test 責務を過不足なく見える化する

## Work Plan Mapping

- scenario 一覧は `Acceptance Checks` の骨子にする
- 観測点は `Required Evidence` に短く接続する
- 未解消の状態遷移や運用条件は `open_questions` と `residual_risks` に分ける

## Avoid

- happy path のみの列挙
- 観測点がない scenario
- human review 前に implementation-scope 相当の断定を書くこと
