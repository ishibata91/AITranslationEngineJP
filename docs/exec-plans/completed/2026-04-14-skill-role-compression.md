# 実装計画

- workflow: work
- status: completed
- lane_owner: skill-modification
- scope: .codex/skills, .codex/README.md, .codex/workflow.md, .codex/agents, docs/exec-plans/templates, docs/exec-plans/active
- task_id: 2026-04-14-skill-role-compression
- task_catalog_ref: N/A
- parent_phase: N/A

## 要求要約

- skill を role-based に圧縮する
- 入口 orchestrator は `orchestrate` だけにする
- `diagramming-*` を `diagramming` に統合する
- `working-light` を廃止する
- `implement-*` を `implement` に統合する
- `ui-check` は `review` に統合する

## 判断根拠

- 現行 skill は lane / phase / role が混在し、同種責務が複数 directory に分散している
- live workflow は `.codex/skills/`、contract JSON、`.codex/README.md`、`.codex/workflow.md`、template を同時にそろえないと崩れる
- fix lane と implementation lane の後段はすでに共通化が進んでおり、role 名への圧縮が可能である
- diagramming は D2 / PlantUML / structure diff に分かれているが、入口としては同一 role にまとめられる

## 対象範囲

- `.codex/skills/`
- `.codex/README.md`
- `.codex/workflow.md`
- `.codex/agents/*.toml`
- `docs/exec-plans/templates/`
- `docs/exec-plans/active/`

## 対象外

- product code
- `docs/` の恒久仕様変更
- completed exec-plan の履歴書き換え

## 依存関係・ブロッカー

- `diagramming` の D2 記法は Context7 で確認済み
- 旧 skill 削除時は rename 後の参照整合を先に揃える必要がある

## 並行安全メモ

- 変更対象は workflow 契約と template に限定する
- active plan のうち live 参照を持つものだけ新 skill 名へ同期する

## 実装計画

- `parallel_task_groups`:
  - `group_id`: role-skill-compression
  - `can_run_in_parallel_with`: none
  - `blocked_by`: none
  - `completion_signal`: 新 skill、contract、workflow docs、template、active plan の参照が同じ命名体系でそろう
- `tasks`:
  - `task_id`: create-role-skills
  - `owned_scope`: `.codex/skills/`, `.codex/agents/`
  - `depends_on`: none
  - `parallel_group`: role-skill-compression
  - `required_reading`: `.codex/README.md`, `.codex/workflow.md`, `docs/exec-plans/completed/2026-04-10-phase-2-skill-split.md`, `docs/exec-plans/completed/2026-04-12-fix-lane-implementation-reuse.md`
  - `validation_commands`: `python3 scripts/harness/run.py --suite structure`
  - `task_id`: sync-docs-and-templates
  - `owned_scope`: `.codex/README.md`, `.codex/workflow.md`, `docs/exec-plans/templates/`, `docs/exec-plans/active/`
  - `depends_on`: create-role-skills
  - `parallel_group`: role-skill-compression
  - `required_reading`: `docs/index.md`, `docs/exec-plans/templates/work-plan.md`, `docs/exec-plans/templates/scenario-tests.md`
  - `validation_commands`: `python3 scripts/harness/run.py --suite structure`, `python3 scripts/harness/run.py --suite all`

## 受け入れ確認

- live skill 一覧が `orchestrate`、`distill`、`investigate`、`design`、`implement`、`tests`、`review`、`diagramming`、`skill-modification`、`updating-docs` の 10 個に揃っている
- `review` が `design-review`、`ui-check`、`implementation-review` を mode として扱える
- `implement` が frontend / backend / mixed を `implementation_target` で扱える
- `work-plan.md` が `impl-plan.md` と `fix-plan.md` を代替し、active plan の記録先として使える
- `.codex/README.md` と `.codex/workflow.md` から legacy orchestrator と phase skill の live 参照が消えている

## 必要な証跡

- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite all`

## HITL 状態

- user が role-based skill 圧縮方針を承認済み

## 承認記録

- user request at 2026-04-14

## review 用差分図

- N/A

## 差分正本適用先

- `.codex/skills/`
- `.codex/README.md`
- `.codex/workflow.md`
- `docs/exec-plans/templates/work-plan.md`

## Closeout Notes

- 旧 template と旧 skill directory は live 参照の除去後に整理する
- active plan の関連記述は completed 履歴には波及させない

## 結果

- `.codex/skills/` を role-based の 10 skill 構成へ再編し、`orchestrate`、`distill`、`investigate`、`design`、`implement`、`tests`、`review`、`diagramming` を live skill として作り直した
- `orchestrate.to.<skill>.json` と `<skill>.to.orchestrate.json` の contract へ統一し、各 skill の `SKILL.md`、`agents/openai.yaml`、`references/permissions.json` を同期した
- `.codex/README.md`、`.codex/workflow.md`、`.codex/agents/*.toml` を新命名と mode 分岐へ更新した
- `docs/exec-plans/templates/work-plan.md` を追加し、`impl-plan.md` と `fix-plan.md` を廃止した
- live 参照から外れた旧 skill directory と旧 contract JSON を整理し、active plan の lane owner と主要 handoff 名だけを新体系へ同期した
- `find .codex/skills -name '*.json' -print0 | xargs -0 -n1 python3 -m json.tool >/dev/null` が通過した
- `python3 scripts/harness/run.py --suite structure` が通過した
- `python3 scripts/harness/run.py --suite all` は skill 再編差分自体ではなく既存 backend build failure により失敗した。失敗点は `internal/usecase/master_dictionary_usecase.go` と `internal/service/master_dictionary_service_test.go` の未解消シンボル不整合である一方、frontend lint / test と `sonar-scanner` は通過した
