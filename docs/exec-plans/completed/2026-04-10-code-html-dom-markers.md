# 実装計画

- workflow: impl
- status: completed
- lane_owner: codex
- scope: docs/screen-design/code.html, docs/exec-plans
- task_id: 2026-04-10-code-html-dom-markers
- task_catalog_ref: N/A
- parent_phase: N/A

## 要求要約

- `docs/screen-design/code.html` で、AI に固定させたくない未確定要素を DOM 上で識別できるようにする
- 未確定要素を `remove` と `placeholder` に分けて属性で表現する

## 判断根拠

- `code.html` は product UI ではなく design mock として扱う
- decorative image、AI provider status、ブランド名、見出し、sample data は確定仕様として扱うと誤解を生む
- class 名や comment ではなく data 属性を正本にすると downstream の DOM 処理が安定する

## 対象範囲

- `docs/screen-design/code.html`
- `docs/exec-plans/active/2026-04-10-code-html-dom-markers.md`

## 対象外

- product code
- `.codex/`
- `docs/spec.md` の恒久要件変更

## 依存関係・ブロッカー

- なし

## 並行安全メモ

- 既存の未コミット変更があるため、本タスクは `code.html` と plan file だけに限定する

## UI

- `data-dom-action="remove"` と `data-remove-reason` を、decorative image、AI provider status、右下 icon、背景 artifact に付ける
- `data-design-status="placeholder"` と `data-placeholder-kind` を、ブランド名、nav 文言、検索 placeholder、CTA 文言、hero copy、job sample data、引用文に付ける
- 画面骨格は維持し、未確定部分だけを DOM で識別可能にする

## Scenario

- downstream が `[data-dom-action="remove"]` を除去すると、未確定の装飾や mock status が消えた状態の骨格だけが残る
- downstream が `[data-design-status="placeholder"]` を読むと、文言や sample data を確定情報として使わない判断ができる

## Logic

- `remove` は DOM から除去する対象に限定する
- `placeholder` は DOM に残すが、AI や抽出処理で確定文言として扱わない対象に限定する
- `remove` と `placeholder` は同一要素に併用しない

## 実装計画

- `parallel_task_groups`:
  - `group_id`: docs-marking
  - `can_run_in_parallel_with`: none
  - `blocked_by`: none
  - `completion_signal`: `code.html` に属性契約が反映され、plan に結果が記録される
- `tasks`:
  - `task_id`: mark-dom-targets
  - `owned_scope`: `docs/screen-design/code.html`
  - `depends_on`: none
  - `parallel_group`: docs-marking
  - `required_reading`: `docs/spec.md`, `docs/screen-design/design-system-ethereal-archive.md`
  - `validation_commands`: `python3 scripts/harness/run.py --suite structure`
  - `task_id`: close-plan
  - `owned_scope`: `docs/exec-plans/active/2026-04-10-code-html-dom-markers.md`
  - `depends_on`: mark-dom-targets
  - `parallel_group`: docs-marking
  - `required_reading`: `docs/exec-plans/active/README.md`
  - `validation_commands`: `python3 scripts/harness/run.py --suite structure`

## 受け入れ確認

- `remove` 対象に `data-dom-action` と `data-remove-reason` が付与されている
- `placeholder` 対象に `data-design-status` と `data-placeholder-kind` が付与されている
- `remove` と `placeholder` の役割が混在していない

## 必要な証跡

- `python3 scripts/harness/run.py --suite structure` の通過結果

## HITL 状態

- human が `docs/screen-design/code.html` の更新を明示済み

## 承認記録

- user request at 2026-04-10

## review 用差分図

- N/A

## 差分正本適用先

- `docs/screen-design/code.html`

## Closeout Notes

- 完了後は plan を `docs/exec-plans/completed/` へ移動する

## 結果

- `docs/screen-design/code.html` に DOM annotation contract の comment を追加した
- `data-dom-action="remove"` と `data-remove-reason` を、AI provider status block、装飾画像 block、右下 icon、背景 artifact に付与した
- `data-design-status="placeholder"` と `data-placeholder-kind` を、ブランド名、nav 文言、検索 placeholder、CTA 文言、hero title / copy、job sample data、基盤データ card 内の英語補助ラベルに付与した
- `python3 scripts/harness/run.py --suite structure` が通過した
