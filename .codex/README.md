# .codex

このディレクトリは、AITranslationEngineJp の live workflow の正本です。
プロダクト仕様と設計は `docs/` を正本とし、lane、skill、agent の役割と handoff は `.codex/` を正本とします。
実装レーンは `workflow.md` の段階番号に合わせた `phase-*` skill と `orchestrating-*` skill で進めます。過去運用の独自 packet や独自 loop は持ち込みません。

## Naming Rule

- workflow 文書では、論理名と実名を分離しない
- 初出または重要な参照は `論理名 (`actual-name`)` を優先する
- 人間 review で意味が先に読めて、actual skill / agent name でも検索できる記述を優先する

## 入口

- 実装レーンの入口: `skills/orchestrating-implementation/SKILL.md`
- バグ修正の入口: `skills/orchestrating-fixes/SKILL.md`
- workflow 鳥瞰図: `workflow.md`

## 標準フロー

### Impl lane

`User -> implementation orchestrator (`orchestrating-implementation`) -> phase-1-distill -> phase-2-design -> phase-2.5-design-review -> human approval -> phase-4-plan -> phase-5-test-implementation -> phase-6-implement-frontend / phase-6-implement-backend -> phase-7-unit-test -> phase-8-review -> full harness -> close`

- implementation orchestrator (`orchestrating-implementation`) は active plan を起点に、各 phase を `workflow.md` の順序で進める
- `phase-1-distill` は要求整理として facts / constraints / gaps / required reading を返す
- `phase-2-design` は `UI` / `Scenario` / `Logic` を task-local design として固定し、必要な時は review 用差分図もこの段階で揃える
- `phase-2.5-design-review` は詳細設計全体を単発 review し、`pass` または `reroute` を返す
- 人間確認が必要な論点は第3段階で active plan に固定する
- `phase-4-plan` は承認済み作業計画を実装順、担当範囲、検証コマンドを持つ implementation brief に変える
- `phase-5-test-implementation` は `Scenario` を tests / fixtures / acceptance checks / validation commands へ適用し、必要な test / fixture を最小範囲で実装する
- `phase-6-implement-frontend` / `phase-6-implement-backend` は担当範囲だけを実装し、local validation を返す
- `phase-7-unit-test` は unit test と coverage gap を補う
- `phase-8-review` は実装が詳細設計と整合しているかだけを単発で確認する
- review 用差分図は `phase-2-design` が必要時に `diagramming-structure-diff` で用意し、完了時に承認済み差分を正本へ適用する

### Fix lane

`User -> fix orchestrator (`orchestrating-fixes`) -> distilling-fixes -> tracing-fixes -> (必要時 logging-fixes / analyzing-fixes) -> phase-5-test-implementation -> implementing-fixes -> reviewing-fixes -> reporting-risks -> close`

- fix orchestrator (`orchestrating-fixes`) は active plan を起点に、`workflow.md` の修正レーンを順番に進める
- `distilling-fixes` は既知事実、再現条件、関連仕様、関連コードを整理する
- `tracing-fixes` は原因仮説と観測方針を返す
- `logging-fixes` と `analyzing-fixes` は必要時だけ使う
- `phase-5-test-implementation` は回帰 test / fixture を先に置く
- `implementing-fixes` は承認済み範囲の恒久修正を行う
- `reviewing-fixes` は単発 review を行う
- `reporting-risks` は必要な時だけ残留リスクを短くまとめる

## 設計記録の扱い

- 非自明な変更は `docs/exec-plans/active/` に plan を置く
- 実装 task でだけ必要になる `UI` / `Scenario` / `Logic` は plan の中に section として置く
- 完了後も保持すべき詳細は `docs/` の正本、コード、型、tests、acceptance checks へ昇格する
- active plan を別 artifact 群へ分解しない
- `orchestrating-* -> downstream skill` の handoff contract 例は、各 orchestrating skill 配下の `references/*.json` を参照する
- `downstream skill -> orchestrating-*` の返却 contract 例は、各 downstream skill 配下の `references/*.json` を参照する
- 各 skill の `references/permissions.json` は、その skill が実行してよい操作、してはいけない操作、期待される返却、停止条件を表す role contract として扱う

## Review と reroute

- review は single-pass とする
- 実装レーンの差し戻しは `workflow.md` の戻り先に合わせる
- 修正レーンの review も `pass` または `reroute` だけを返す
- 繰り返し見つかる指摘は review loop に残さず tests、harness、必要なら plan に昇格する
