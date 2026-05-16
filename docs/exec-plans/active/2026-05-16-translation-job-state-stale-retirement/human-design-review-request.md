# Human Design Review Request

- `skill`: `implement-lane`
- `status`: `approved`
- `source`: `scenario-design.md`, `design-diff.md`
- `return_to`: `implement_lane`

## レビュー対象

- `scenario-design.md`: 人間回答反映済みのシナリオ設計。
- `scenario-design.candidate-coverage.json`: シナリオ候補の採否と人間回答反映。
- `scenario-design.requirement-coverage.json`: 詳細要求タイプの充足状態。
- `design-diff.md`: 予定変更箇所だけの設計差分メモ。
- `design-diff.component.puml`: 追加、削除、変更しない接続先のコンポーネント差分図。
- `design-diff.sequence.puml`: state facts と操作可否導出のシーケンス差分図。

## 人間回答反映済み判断

- `JobIOService`: stale として architecture 正本から外す。
- `observability-log-addition`: completed へ移動済みなので、今回の active task-local 更新対象にしない。
- `cancelled`: fixture spelling を今回の stale 廃止に含め、`canceled` へそろえる。

## レビューしてほしい点

- `pending` を canonical state へ昇格しない判断でよいか。
- `JobIOService` を stale として削除予定にする差分でよいか。
- read model の操作可否を、`TranslationJobPolicy` と同じ state 事実から導く設計でよいか。
- completed archive の `observability-log-addition` を変更しない扱いでよいか。
- `cancelled` fixture spelling を今回範囲へ含めてよいか。

## 設計レビュー後の進行

承認された場合、`designer` を起動して `implementation-scope.md` を作成する。
差し戻しの場合、`scenario-design.md` と `design-diff.*` を修正してから再レビューする。

## 人間レビュー結果

- `review_status`: `approved`
- `approved_at`: `2026-05-16`
- `approved_by`: 人間
- `comment`: `approve`

## 次成果物

`implementation-scope.md` を作成する。

## 検証

- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/scenario-design.md --json`: pass
- `plantuml --check-syntax docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/design-diff.component.puml`: pass
- `plantuml --check-syntax docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/design-diff.sequence.puml`: pass
- `plantuml -tsvg docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/design-diff.component.puml docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/design-diff.sequence.puml`: pass
