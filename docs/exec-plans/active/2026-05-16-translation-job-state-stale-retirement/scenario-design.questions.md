# Scenario Design Questions

## 状態

- `status`: `waiting_human_decision`
- `source`: `./scenario-design.md`

### [Q-001] `JobIOService` の扱い

決める仕様:
状態事実の取得と保存を、現行 architecture 正本の `JobIOService` として実体化するか、architecture 正本から外して既存境界へ寄せるかを決める。

決定済み:
`JobIOService` は、状態遷移可否、terminal guard、provider 応答検証、UI 表示文言を判断しない。現時点の実体 package は `doc.go` だけである。

未確定:
`JobIOService` を今回の stale 廃止で外すか、別 task で実体化するかが未確定である。

選択肢:
1. `JobIOService` を architecture 正本から外し、状態事実の取得と保存の境界を既存の usecase、service、repository へ寄せる。
2. `JobIOService` を別 task で実体化し、今回の stale 廃止では残留理由と後続 task 条件だけを記録する。
3. 今回は判断保留とし、`JobIOService` 参照を既知の残留参照として残す。
4. その他

AI 推奨:
1 を推奨する。実体 package がない構造主語を残すと、状態保存境界と状態遷移判断の誤読が続く。

### [Q-002] active `observability-log-addition` 旧名参照の扱い

決める仕様:
active `observability-log-addition` に残る `StateMachine` / `JobIOService` 旧名参照を、今回の task-local 更新に含めるかを決める。

決定済み:
`docs/exec-plans/completed/**` は履歴として変更しない。product code から `StateMachine` 旧名は外れている。

未確定:
active task-local の旧名参照を今回更新するか、observability task 再開時に更新するかが未確定である。

選択肢:
1. 今回の stale 廃止に含め、active task-local の旧名参照を現在の責務名または残留理由へ更新する。
2. observability task 再開時に更新し、今回の task では残留参照と後続条件だけを記録する。
3. 今回は検索結果だけを記録し、更新 task はまだ起票しない。
4. その他

AI 推奨:
1 を推奨する。active task-local は後続作業の入力になるため、旧名参照を残すと古い責務境界が再注入される可能性がある。

### [Q-003] `cancelled` fixture spelling の今回範囲

決める仕様:
正本 spelling と異なる `cancelled` fixture spelling を、今回の stale 廃止で修正するかを決める。

決定済み:
正本仕様と service 実装は `Canceled` / `canceled` を使う。`cancelled` は正本 state として扱わない。

未確定:
`PersonaGenerationPhaseContractStub` の fixture spelling 修正を今回の実装範囲へ含めるか、別 task に送るかが未確定である。

選択肢:
1. 今回の stale 廃止に含め、fixture spelling を `canceled` へそろえる。
2. 別 task に送り、今回の task では検索漏れ防止の残留参照として記録する。
3. fixture だけの差分として現状維持し、正本 state との差分を検証対象にしない。
4. その他

AI 推奨:
1 を推奨する。正本 spelling が明示済みであり、検索漏れと terminal 判定の誤読を小さな範囲で減らせる。
