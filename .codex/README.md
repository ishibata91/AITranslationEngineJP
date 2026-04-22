# .codex

このディレクトリは Codex 側 workflow の正本です。
Codex は設計 workflow を進め、GitHub Copilot は人間から渡された承認済み scope から実装します。

プロダクト仕様と設計判断の正本は `docs/` です。
workflow、skill、agent、handoff 契約の正本は `.codex/` と `.github/` に分けます。

## Live Skills

- Codex workflow 進行: `skills/propose-plans/SKILL.md`
- 設計壁打ち: `skills/wall-discussion/SKILL.md`
- 設計用文脈整理: `skills/distill/SKILL.md`
- 設計前調査: `skills/investigate/SKILL.md`
- 要件設計 (`requirements-design`): `skills/requirements-design/SKILL.md`
- UI 設計 (`ui-design`): `skills/ui-design/SKILL.md`
- シナリオ設計 (`scenario-design`): `skills/scenario-design/SKILL.md`
- 実装スコープ (`implementation-scope`): `skills/implementation-scope/SKILL.md`
- 図: `skills/diagramming/SKILL.md`
- docs 正本化: `skills/updating-docs/SKILL.md`
- workflow 契約変更: `skills/skill-modification/SKILL.md`
- 編集前 gate: `skills/gateguard/SKILL.md`
- Codex run report: `skills/codex-work-reporting/SKILL.md`

## Agent / Skill Contract

- live Codex agent は workflow orchestrator (`propose_plans`)、design artifact agent (`designer`)、文脈圧縮 agent (`distiller`)、設計前調査 agent (`investigator`)、docs 更新 agent (`docs_updater`) にする
- `propose_plans` は必要判定、task folder、agent spawn、human review、人間向け Copilot handoff、Copilot 完了後の正本化入口を進める
- `designer` は requirements、UI、scenario、implementation-scope、diagram などの design artifact を task-local に固定する
- `distiller`、`designer`、`investigator`、`docs_updater` は context を引き継がず、handoff packet だけで動く
- agent は実行主体として扱い、permissions、agent 1:1 contract、handoff、stop / reroute を持つ
- skill は knowledge package として扱い、判断基準、標準 pattern、DO / DON'T、checklist を持つ
- Codex agent の詳細 spec は `agents/<agent>.agent.md`、contract は `agents/references/<agent>/contracts/<agent>.contract.json` に置く
- skill 側の旧 `references/permissions.json` や contract slice は互換用であり、新しい active 正本にしない

## 責務境界

- `propose_plans` は Codex workflow の進行役として必要判定、plan、spawn packet、human review、human handoff を扱う
- `propose_plans` は run の closeout、停止、reroute 時に `codex-work-reporting` を参照し、最後に必ず `work_history` 記録材料を作る
- `distiller` は task に関連する事柄と必要資料の判断材料を集める
- `designer` は diagram を含む design bundle と implementation-scope の task-local artifact を作る
- `investigator` は必要な場合だけ実画面や観測対象を確認し、観測事実と risk を返す
- `docs_updater` は Copilot の修正完了が分かった後、human 承認済み scope だけを正本化する
- Codex は product code と product test を変更しない
- Codex は Copilot へ直接 handoff しない。最後に人間へ Copilot handoff packet を返し、人間が Copilot へ引き渡す
- Copilot は `.github/skills/implementation-orchestrate/SKILL.md` から実装と実装時調査を進める
- Copilot の実装前文脈整理は `.github/skills/implementation-distill/SKILL.md` が扱う
- Copilot は docs 正本、`.codex/`、`.github/skills`、`.github/agents` を変更しない

## Design Flow

1. `propose_plans` が `propose-plans` を参照し、distiller、designer、investigator が必要か最初に判定する
2. `propose_plans` が task folder と `plan.md` を作る、または既存 task folder を確認する
3. 必要なら `propose_plans` が `distiller` を context 継承なしで spawn し、task 関連情報と必要資料の判断材料を集める
4. 必要なら `propose_plans` が `designer` を context 継承なしで spawn し、requirements、UI、scenario、implementation-scope、diagram を含む必要資料を作る
5. 必要なら `propose_plans` が `investigator` を context 継承なしで spawn し、実画面や観測対象を確認する
6. 各 agent の戻りを `propose_plans` が `plan.md` の workflow state に反映する
7. design bundle 完了後に `propose_plans` が human review で停止する
8. human 承認後に `designer` が `implementation-scope` を参照し、人間向け Copilot handoff packet の材料を固定する
9. `propose_plans` が Copilot handoff packet を人間へ返す。人間が Copilot へ引き渡す
10. Copilot の修正完了が分かったら、`propose_plans` が正本化の必要性を判定する
11. 必要なら `propose_plans` が `docs_updater` を context 継承なしで spawn し、承認済み範囲だけを正本化する
12. closeout、停止、reroute 時は `codex-work-reporting` を参照し、最後に必ず `work_history` の Codex report とラン横断 finding に必要な材料を作る

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
- task-local UI 判断と scenario は task folder に置く
- UI の構造正本は task-local HTML mock とし、完了時に必要なら `docs/mocks/<page-id>/index.html` へ正本化する
- `implementation-scope` は handoff 履歴であり docs 正本へ昇格しない

## 非 live 扱い

- 旧 `orchestrate` は `propose_plans` agent と `propose-plans` skill に置き換えた
- 旧 `design` は `requirements-design`、`ui-design`、`scenario-design`、`implementation-scope` に分割した
- 旧 flat file 形式の exec-plan は legacy とし、新規 task では使わない
- Codex 側の実装、UI check 専用、log instrumentation agent は live から外した
- diagram 専用 agent は標準 flow では spawn せず、diagram は `designer` が `diagramming` を参照して必要資料として扱う
- GitHub Copilot 側の実装 workflow は `.github/skills` と `.github/agents` を正本にする
- Codex 側 `distill` は design / investigate の入口整理だけを扱い、implement / fix / refactor の実装前整理は `.github/skills/implementation-distill` が扱う
- 退避した旧 skill / agent は `.codex/.trash` に置き、live workflow では参照しない

## Work Plan

- 非自明な変更は `docs/exec-plans/active/<task-id>/` に置く
- 完了後は `docs/exec-plans/completed/<task-id>/` へ移す
- completed plan は履歴として残す
