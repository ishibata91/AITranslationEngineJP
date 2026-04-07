# Implementation Proposal Lane Split

- workflow: impl
- status: review-ready
- lane_owner: `skill-modification`
- scope:
  - `.codex/skills/proposing-implementation/`
  - `.codex/skills/directing-implementation/`
  - `.codex/skills/distilling-implementation/`
  - `.codex/skills/designing-implementation/`
  - `.codex/agents/diagrammer.toml`
  - `.codex/README.md`
  - `.codex/workflow.md`
  - `.codex/workflow_activity_diagram.puml`
  - `docs/index.md`
  - `docs/exec-plans/templates/impl-plan.md`
- task_id:
- task_catalog_ref:
- parent_phase:

## 要求要約

- implementation lane を proposal と execution に分割する。
- proposal lane owner (`proposing-implementation`) を新設し、active exec-plan、`distilling-implementation`、`designing-implementation`、review 用差分図、HITL を担当させる。
- implementation lane owner (`directing-implementation`) は `planning-implementation` 以降だけを担当させる。
- review 用差分図は active exec-plan と同じフォルダに置き、human review で差分が読めるようにする。
- review 用差分図の作成は diagrammer (`diagrammer`) が担当し、diagramming D2 skill (`diagramming-d2`) へ handoff する。
- impl plan template と今後の exec-plan は日本語で扱う。

## 判断根拠

- 既存 `directing-implementation` は proposal と execution の責務が混在している。
- `distilling-implementation` と `designing-implementation` は proposal 側の責務に寄せた方が境界が明確になる。
- human review 前に execution へ進ませない契約を skill と plan の両方に持たせる必要がある。
- 現在の `4humans/diagrams/` には implementation workflow を直接表す正本図がないため、この変更では review 用差分図を active exec-plan 配下に置く。

## 対象範囲

- implementation lane の入口 skill と downstream handoff 契約
- workflow 文書と overview 図
- impl plan template の日本語化
- diagrammer (`diagrammer`) agent 契約

## 対象外

- fix lane の責務変更
- product code や product docs の仕様変更
- `4humans/diagrams/` 正本への新規図追加

## 依存関係・ブロッカー

- human review 用差分図の validate / render が通ること
- implementation lane の入口説明を参照する文書を同期すること

## 並行安全メモ

- 本変更は `.codex/` と workflow docs が中心で、product code とは競合しない。
- review 用差分図は active exec-plan 配下に閉じ、`4humans/diagrams/` 正本を直接編集しない。

## UI

- N/A

## Scenario

- human は `proposing-implementation` が作成した日本語 active exec-plan と review 用差分 SVG を確認して LGTM を出す。
- implementation lane owner (`directing-implementation`) は承認済み active exec-plan を受けた時だけ execution を開始する。

## Logic

- `proposing-implementation` は active exec-plan 作成、implementation distill、task-local design、diagrammer handoff、HITL 停止を担当する。
- `directing-implementation` は planning、test architecture、implementation、Sonar、review、final harness、`4humans sync` を担当する。
- diagrammer (`diagrammer`) は GPT-5.4 high で review 用差分図だけを担当し、diagramming D2 skill (`diagramming-d2`) の validate / render 契約を使う。

## 実装計画

- 新しい `proposing-implementation` skill と handoff 契約を追加する。
- `directing-implementation` から proposal 責務を外し、execution 開始条件を明記する。
- `distilling-implementation` と `designing-implementation` の入力 / 返却契約を proposal lane に付け替える。
- workflow docs と activity diagram を proposal -> execution の順に更新する。
- 日本語の impl plan template と review 用差分図 artifact を追加する。

## 受け入れ確認

- `proposing-implementation` の新設と `directing-implementation` の責務分割が `.codex/README.md`、各 `SKILL.md`、reference JSON で矛盾しない。
- `diagrammer` agent 契約が存在し、diagramming D2 skill (`diagramming-d2`) への handoff を文書化できている。
- `docs/exec-plans/templates/impl-plan.md` が日本語見出しと HITL 記録欄を持つ。
- review 用差分 D2 / SVG が active exec-plan 配下に存在し、追加は緑、削除は赤で読める。

## 必要な証跡

- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite design`
- `d2 validate docs/exec-plans/active/2026-04-07-proposing-implementation-lane-split-review-diff.d2`
- `d2 -t 201 docs/exec-plans/active/2026-04-07-proposing-implementation-lane-split-review-diff.d2 docs/exec-plans/active/2026-04-07-proposing-implementation-lane-split-review-diff.svg`

## 4humans Sync

- この変更では `4humans/diagrams/` 正本を更新しない。
- review 用差分図は `docs/exec-plans/active/` 配下の一時 artifact として扱う。
- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`

## HITL 状態

- 承認済み

## 承認記録

- 2026-04-07 human LGTM

## review 用差分図

- `docs/exec-plans/active/2026-04-07-proposing-implementation-lane-split-review-diff.d2`
- `docs/exec-plans/active/2026-04-07-proposing-implementation-lane-split-review-diff.svg`

## 結果

- `proposing-implementation`、`diagrammer`、関連 handoff 契約、workflow 文書、impl plan template を更新した。
- `python3 scripts/harness/run.py --suite structure`、`python3 scripts/harness/run.py --suite design`、`d2 validate`、`d2 -t 201` が通った。
- human LGTM を反映し、この plan を completed へ移す。
