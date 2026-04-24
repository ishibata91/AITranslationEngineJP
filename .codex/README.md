# .codex

このディレクトリは Codex 側 workflow の正本です。
Codex は設計 workflow を進め、GitHub Copilot は人間から渡された承認済み scope から実装します。

プロダクト仕様と設計判断の正本は `docs/` です。
workflow、skill、agent、handoff 契約の正本は `.codex/` と `.github/` に分けます。
live workflow の説明本文と判断基準の正本はこの `README.md` とします。
`.codex/workflow.md` は補助図であり、live 判断を上書きしません。

## Live Skills

- Codex workflow 進行: `skills/propose-plans/SKILL.md`
- 設計壁打ち: `skills/wall-discussion/SKILL.md`
- 設計用文脈整理: `skills/distill/SKILL.md`
- design bundle 進行: `skills/design-bundle/SKILL.md`
- 設計前調査: `skills/investigate/SKILL.md`
- UI 設計 (`ui-design`): `skills/ui-design/SKILL.md`
- シナリオ設計 (`scenario-design`): `skills/scenario-design/SKILL.md`
- 実装スコープ (`implementation-scope`): `skills/implementation-scope/SKILL.md`
- docs 正本化: `skills/updating-docs/SKILL.md`
- workflow 契約変更: `skills/skill-modification/SKILL.md`
- 編集前 gate: `skills/gateguard/SKILL.md`
- Codex run report: `skills/codex-work-reporting/SKILL.md`
- 実装後 review conductor: `skills/codex-review-conductor/SKILL.md`
- 実装後 review 観点: `skills/codex-review-behavior/SKILL.md`、`skills/codex-review-contract/SKILL.md`、`skills/codex-review-trust-boundary/SKILL.md`、`skills/codex-review-state-invariant/SKILL.md`

## Agent / Skill Contract

- live Codex agent は workflow orchestrator (`propose_plans`)、design artifact agent (`designer`)、文脈圧縮 agent (`distiller`)、設計前調査 agent (`investigator`)、docs 更新 agent (`docs_updater`)、実装後 review conductor (`review_conductor`)、観点別 review agent にする
- `propose_plans` は必要判定、task folder、agent spawn、human review、人間向け Copilot handoff、Copilot 完了後の正本化入口を進める
- `designer` は scenario を必須要件の固定点として作り、UI 変更がある時だけ `ui-design` を追加し、human review 後に `implementation-scope` を固定する
- `distiller`、`designer`、`investigator`、`docs_updater` は context を引き継がず、handoff packet だけで動く
- `review_conductor` は Copilot 完了後の `codex exec` 入口として、観点別 review agent を context 継承なしで並列 spawn する
- agent は runtime binding と機械契約の owner として扱い、`agents/<agent>.toml`、permissions、agent 1:1 contract を持つ
- skill は knowledge package であり、人間可読な実行説明の正本として扱い、判断基準、標準 pattern、DO / DON'T、checklist、handoff、stop / reroute を持つ
- Codex agent の人間可読な実行説明は対応する `skills/*/SKILL.md` に置き、binding は `agents/<agent>.toml`、contract は `agents/references/<agent>/contracts/<agent>.contract.json`、permissions は `agents/references/<agent>/permissions.json` に置く
- `.agent.md` は使わない

## 責務境界

- `propose_plans` は Codex workflow の進行役として必要判定、plan、spawn packet、human review、human handoff を扱う
- `propose_plans` は run の closeout、停止、reroute 時に `codex-work-reporting` を参照し、最後に必ず `work_history` 記録材料を作る
- `distiller` は task に関連する事柄と必要資料の判断材料を集める
- `designer` は design bundle と implementation-scope の task-local artifact を作る
- `investigator` は必要な場合だけ実画面や観測対象を確認し、観測事実と risk を返す
- `docs_updater` は Copilot の修正完了が分かった後、human 承認済み scope だけを正本化する
- `review_conductor` は Copilot の全 implementation handoff と final validation 完了後、diff から取得した実コードを観点グループ別に score 化する
- 観点別 review agent は挙動正しさ、契約・互換性、権限・信頼境界、状態・データ不変条件のいずれか 1 つだけを扱う
- Codex は product code と product test を変更しない
- Codex は Copilot へ直接 handoff しない。最後に人間へ Copilot handoff packet を返し、人間が Copilot へ引き渡す
- Copilot は `.github/skills/implementation-orchestrate/SKILL.md` から実装、実装時調査、final validation、Codex review 呼び出しを進める
- Copilot の実装前文脈整理は `.github/skills/implementation-distill/SKILL.md` が扱う
- Copilot は docs 正本、`.codex/`、`.github/skills`、`.github/agents` を変更しない

## Design Flow

1. `propose_plans` が `propose-plans` を参照し、distiller と investigator の要否を判定し、承認済み design bundle がない限り `designer` を使う
2. `propose_plans` が task folder と `plan.md` を作る、または既存 task folder を確認する
3. 必要なら `propose_plans` が `distiller` を context 継承なしで spawn し、task 関連情報と必要資料の判断材料を集める
4. `propose_plans` が `designer` を context 継承なしで spawn し、`scenario-design` を必須で作り、UI 変更がある時だけ `ui-design` を追加する
5. 必要なら `propose_plans` が `investigator` を context 継承なしで spawn し、実画面や観測対象を確認する
6. 各 agent の戻りを `propose_plans` が `plan.md` の workflow state に反映する
7. design bundle 完了後に `propose_plans` が human review で停止する
8. human 承認後に `propose_plans` が `designer` を再度 context 継承なしで spawn し、`implementation-scope` を固定する
9. `propose_plans` が Copilot handoff packet を人間へ返す。人間が Copilot へ引き渡す
10. Copilot は全 implementation handoff 完了後に suite-all と Sonar check を実行し、その後 `codex exec` で `review_conductor` を呼び出す
11. `review_conductor` は 4 つの観点別 review agent を並列 spawn し、各 score が 0.85 より大きい場合だけ pass とする
12. Copilot の修正完了と Codex review 結果が分かったら、`propose_plans` が正本化の必要性を判定する
13. 必要なら `propose_plans` が `docs_updater` を context 継承なしで spawn し、承認済み範囲だけを正本化する
14. closeout、停止、reroute 時は `codex-work-reporting` を参照し、最後に必ず `work_history` の Codex report とラン横断 finding に必要な材料を作る

## Exec Plan Folder

- 新規 task は `docs/exec-plans/active/<task-id>/` に folder として作る
- `plan.md` は索引、状態、HITL、validation、closeout だけを書く
- 各 skill の資料は同じ folder の skill 名つき file に分ける
- AI は最初に `plan.md` だけ読み、必要な資料だけ追加で読む
- 完了後は folder ごと `docs/exec-plans/completed/<task-id>/` へ移す

## Docs 正本化

- docs 正本化は Copilot の修正完了が分かった後に扱う
- docs 正本化は Codex 側だけで扱う
- human 承認済みの artifact だけ `docs_updater` が `updating-docs` を参照して正本へ反映する
- task-local UI 要件契約と scenario は task folder に置く
- UI の細かな visual polish は実装後に人間が実物を確認して直す
- `implementation-scope` は handoff 履歴であり docs 正本へ昇格しない

## 非 live 扱い

- 旧 `orchestrate` は `propose_plans` agent と `propose-plans` skill に置き換えた
- 旧 `design` は `scenario-design`、`ui-design`、`implementation-scope` 中心の design bundle に再整理した
- 旧 flat file 形式の exec-plan は legacy とし、新規 task では使わない
- Codex 側の実装、UI check 専用、log instrumentation agent は live から外した
- GitHub Copilot 側の実装 workflow は `.github/skills` と `.github/agents` を正本にする
- Codex 側 `distill` は design / investigate の入口整理だけを扱い、implement / fix / refactor の実装前整理は `.github/skills/implementation-distill` が扱う
- Codex 側の人間可読な runtime 説明は skill へ集約し、`.codex/agents/*.agent.md` は持たない
- `.codex/workflow.md` は補助図として残し、live workflow の正本にはしない
- 退避した旧 skill / agent は `.codex/.trash` に置き、live workflow では参照しない

## Work Plan

- 非自明な変更は `docs/exec-plans/active/<task-id>/` に置く
- 完了後は `docs/exec-plans/completed/<task-id>/` へ移す
- completed plan は履歴として残す
