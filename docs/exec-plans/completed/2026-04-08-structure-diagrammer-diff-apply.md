# 実装計画テンプレート

- workflow: impl
- status: completed
- lane_owner: skill-modification
- scope: structure-diagrammer-diff-apply
- task_id: 2026-04-08-structure-diagrammer-diff-apply
- task_catalog_ref:
- parent_phase: workflow-contract-sync

## 要求要約

- `diagrammer` を proposal / execution workflow から外し、`structure_diagrammer` だけで構造差分図を扱う。
- `structure_diagrammer` は proposal で差分を生成し、`directing-implementation` close 時に承認済み差分を `diagrams/backend/` 正本へ適用する。

## 判断根拠

<!-- Decision Basis -->

- proposal 時点で正本を先に変更せず、review では差分だけを見せたい。
- 実装完了後に承認済み差分を正本へ反映する方が、proposal と execution の責務分離が明確になる。
- `diagrammer` の専用役を残すより、構造図専用の `structure_diagrammer` に責務を集約した方が単純になる。

## 対象範囲

- `.codex/agents/`
- `.codex/skills/proposing-implementation/`
- `.codex/skills/directing-implementation/`
- `.codex/skills/directing-fixes/`
- `.codex/skills/working-light/`
- `.codex/README.md`
- `.codex/workflow.md`
- `.codex/workflow_activity_diagram.puml`
- `docs/exec-plans/templates/impl-plan.md`

## 対象外

- 実際の `diagrams/backend/` 正本作成
- `docs/` 正本の更新
- product code の変更

## 依存関係・ブロッカー

- `design` / `all` harness は着手前から `docs/core-beliefs.md`、`docs/index.md`、`4humans/diagrams/overview-manifest.json` 由来の failure を含む。

## 並行安全メモ

- shared scope は `.codex` workflow 契約と impl plan template に限る。

## UI

- N/A

## Scenario

- proposal lane は `designing-implementation` の後に `structure_diagrammer` だけを実行し、active plan 配下へ review 用構造差分図を置く。
- execution close は、承認済み差分を `diagrams/backend/` 正本へ適用してから plan close へ進む。

## Logic

- `structure_diagrammer` は `proposal_diff` と `apply_to_source` の 2 mode を持つ。
- active plan は review 用差分図と差分正本適用先を記録する。
- `diagrammer` は workflow から削除し、`4humans` 更新が必要な時は `diagramming-d2` を直接使う。

## 実装計画

<!-- Implementation Plan -->

- `structure_diagrammer` 契約を差分生成 / 正本適用の二段階に更新する。
- `proposing-implementation` と `directing-implementation` の handoff を `review_diff_diagrams` と `diagram_source_targets` へ整理する。
- `diagrammer` 依存を `.codex` workflow docs と関連 skill から取り除く。

## 受け入れ確認

- structure harness が通る。
- 変更対象文書間で `structure_diagrammer` の責務、diff artifact 名、正本適用タイミングが一致する。
- 既存 design harness failure 以外の新規 failure を増やさない。

## 必要な証跡

<!-- Required Evidence -->

- 更新後の workflow / skill / agent / reference files
- harness 実行結果

## HITL 状態

- N/A

## 承認記録

- user requested workflow simplification on 2026-04-08

## review 用差分図

- active exec-plan 配下の review D2 / SVG

## 差分正本適用先

- `diagrams/backend/components.d2`
- `diagrams/backend/<component>/<component>.d2`

## 4humans Sync

- N/A

## 結果

<!-- Outcome -->

- `diagrammer` agent を削除し、implementation proposal lane の図作成役を `structure_diagrammer` に一本化した。
- `structure_diagrammer` を `proposal_diff` と `apply_to_source` の 2 mode 契約へ更新した。
- `proposing-implementation` は review 用差分図と差分正本適用先を handoff し、`directing-implementation` close は承認済み差分を `diagrams/backend/` 正本へ適用する契約に変更した。
- `directing-fixes` と `working-light` の `diagrammer` 依存を外し、`diagramming-d2` 直接利用へ変更した。
- `python3 scripts/harness/run.py --suite structure` は通過した。
- `python3 scripts/harness/run.py --suite design` と `python3 scripts/harness/run.py --suite all` は、着手前から存在する `docs/core-beliefs.md` / `docs/index.md` の pattern 不足と `4humans/diagrams/overview-manifest.json` 欠如で失敗した。
