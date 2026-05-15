# Implement Lane Human Decision Request

- `skill`: `implement-lane`
- `status`: `waiting_human_decision`
- `source`: `scenario-design.questions.md`
- `return_to`: `implement_lane`

## 現在成果物

| 成果物ID | 状態 | 根拠 |
| --- | --- | --- |
| `task 枠` | 完了 | `implement-lane-task-frame.md` |
| `scenario_candidates` | 完了 | `scenario-candidates.*.md` 6 件 |
| `シナリオ設計` | 停止中 | `scenario-design.md`, `scenario-design.questions.md` |
| `UI設計` | 該当なし | UI 変更なし |
| `設計差分図` | 未着手 | `シナリオ設計` の人間判断待ち |
| `人間設計レビュー` | 未着手 | `設計差分図` 未完了 |
| `実装範囲` | 停止中 | `人間設計レビュー` 未完了 |

## 停止理由

`scenario-design` は、人間判断 3 件が残るため完了できない。
`設計差分図`、`人間設計レビュー`、`implementation-scope` は、質問回答後に進める。

## 人間判断

### Q-001 `JobIOService` の扱い

決める仕様:
状態事実の取得と保存を、現行 architecture 正本の `JobIOService` として実体化するか、architecture 正本から外して既存境界へ寄せるかを決める。

選択肢:

1. `JobIOService` を architecture 正本から外し、状態事実の取得と保存の境界を既存の usecase、service、repository へ寄せる。
2. `JobIOService` を別 task で実体化し、今回の stale 廃止では残留理由と後続 task 条件だけを記録する。
3. 今回は判断保留とし、`JobIOService` 参照を既知の残留参照として残す。

AI 推奨:
1 を推奨する。実体 package がない構造主語を残すと、状態保存境界と状態遷移判断の誤読が続く。

### Q-002 active `observability-log-addition` 旧名参照の扱い

決める仕様:
active `observability-log-addition` に残る `StateMachine` / `JobIOService` 旧名参照を、今回の task-local 更新に含めるかを決める。

選択肢:

1. 今回の stale 廃止に含め、active task-local の旧名参照を現在の責務名または残留理由へ更新する。
2. observability task 再開時に更新し、今回の task では残留参照と後続条件だけを記録する。
3. 今回は検索結果だけを記録し、更新 task はまだ起票しない。

AI 推奨:
1 を推奨する。active task-local は後続作業の入力になるため、旧名参照を残すと古い責務境界が再注入される可能性がある。

### Q-003 `cancelled` fixture spelling の今回範囲

決める仕様:
正本 spelling と異なる `cancelled` fixture spelling を、今回の stale 廃止で修正するかを決める。

選択肢:

1. 今回の stale 廃止に含め、fixture spelling を `canceled` へそろえる。
2. 別 task に送り、今回の task では検索漏れ防止の残留参照として記録する。
3. fixture だけの差分として現状維持し、正本 state との差分を検証対象にしない。

AI 推奨:
1 を推奨する。正本 spelling が明示済みであり、検索漏れと terminal 判定の誤読を小さな範囲で減らせる。

## 再開条件

人間が Q-001、Q-002、Q-003 へ回答したら、`designer` を再起動して `scenario-design` を更新する。
`scenario-design` が通過した後に、`diagrammer` で設計差分図を作る。
設計差分図が揃った後に、人間設計レビューへ進む。

## 検証

- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/scenario-design.md --json`: exit code `1`
- 結果: `status: fail`, `finding_count: 20`, `question_count: 3`
- 判断: 人間判断待ちによる想定停止。
