# Task Plan: 2026-05-16-translation-management-ui-subtraction

- `workflow`: work
- `status`: planned
- `lane_owner`: Codex design-bundle
- `task_id`: `2026-05-16-translation-management-ui-subtraction`
- `task_mode`: UI simplification planning
- `request_summary`: 翻訳管理画面が仕様に対して複雑すぎるため、表示情報、検証表示、操作導線を減らす。
- `goal`: 翻訳管理画面を、利用者の主要目的に沿った少数の画面責務へ戻す。
- `constraints`: 旧 `Job Run` 大箱を復活させない。画面上の情報削減を、状態整合性の削除として扱わない。secret や内部状態名を画面へ出さない。
- `close_conditions`: `detail-spec-diff.md` で親要件、仕様、未決、回答が固定される。`screen-design-diff.<screen-id>.md` で表示する情報、隠す情報、操作、停止条件が画面ごとに固定される。
- `worktree_path`: `/Users/iorishibata/Repositories/AITranslationEngineJP`
- `source_branch`: `codex/2026-05-16-translation-management-ui-subtraction`
- `target_branch`: `master`

## Artifact Index

- `ux_task_frame`: `N/A`
- `detail_spec_diff`: `pending`
- `screen_design_diff`: `pending`
- `storybook_review_request`: `pending-after-frontend-implementation`
- `frontend_human_review`: `required-after-frontend-implementation`
- `approved_frontend_protection`: `pending`
- `implementation_scope`: `pending-after-human-review`
- `detail_spec_target`: `N/A`

## Routing Notes

- `required_reading`:
  - `docs/exec-plans/completed/2026-05-08-translation-flow-navigation-overhaul/plan.md`
  - `docs/exec-plans/completed/2026-05-08-translation-flow-navigation-overhaul/ui-design.md`
  - `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte`
  - `frontend/src/ui/screens/job-run/JobRunPage.svelte`
- `canonicalization_targets`:
  - UI 設計成果物。docs 正本反映は human 承認後に別判断とする。
- `detail_spec_id`: `N/A`
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite frontend-local`
  - `npm run dev:wails:agent-browser`
  - `agent-browser open http://localhost:34115/#translation-management`

## Decision Record

### D-01 画面目的を 3 つへ戻す

決定: 翻訳管理画面の利用者目的は、新しい翻訳 job を作る、未完了 job を続きから進める、完了 job を出力管理へ渡す、の 3 つに絞る。

理由: 現行画面は、状態、検証、進捗、操作、内部都合を同時に表示している可能性がある。利用者は次操作より先に停止理由を読む必要がある。

影響: UI 設計では、各表示項目が 3 目的のどれに必要かを明示する。どの目的にも属さない情報は、画面から削る候補にする。

### D-02 一覧は選択と再開可否に限定する

決定: 未完了 job 一覧では、job を選べるか、続けられるか、続けられないなら短い理由だけを表示する。

理由: 一覧で詳細な検証理由や内部状態を広げると、一覧が診断画面になる。利用者の目的は比較と選択であり、詳細診断ではない。

影響: 詳細な検証理由は、選択後の該当フェーズ画面へ移す。検索、フィルタ、状態表示は、選択を助ける最小限に抑える。

### D-03 バリデーションは常時表示しない

決定: バリデーション結果は、操作直後、進行不能時、または詳細展開時だけ表示する。

理由: 常時表示の検証結果は、操作前から画面を失敗状態に見せる。利用者は、押せる操作と押せない理由を必要な時だけ見たい。

影響: `Job Setup` と各フェーズ画面では、検証一覧を初期表示しない。停止理由は対象操作の近くへ 1 つの主要理由として表示する。

### D-04 内部状態名を主語にしない

決定: `Ready`、`Running`、`RecoverableFailed` などの状態値は、画面の主語にしない。必要な場合は日本語の意味を併記する。

理由: 状態値だけを並べると、利用者が何をすべきかを再解釈する必要がある。

影響: presenter と UI 文言の見直しが必要になる。固定名は contract や test では残してよいが、画面では次操作へ変換する。

### D-05 footer は移動だけに戻す

決定: footer は一覧へ戻る、次へ進む、出力管理へ進む、の移動だけを扱う。

理由: footer に検証情報や実行操作を集めると、移動と実行の責務が混ざる。

影響: 実行、停止、再開、再試行、取消は各フェーズ本文の操作領域へ残す。footer の停止理由は、次へ進めない主要理由 1 件を優先する。

## Screen Reduction Targets

- 未完了一覧: job 選択、再開可否、短い理由だけを表示する。
- 新規作成: 入力選択、AI 設定、job 作成の最小条件だけを表示する。
- フェーズ実行: 現在フェーズの実行操作、現在結果、次へ進めるかだけを表示する。
- 完了確認: 原文と訳文の確認、一覧へ戻る、出力管理へ進むだけを表示する。

## Hidden Or Deferred Information

- 内部 repository 名、secret store key、route state、raw provider payload は表示しない。
- 詳細な検証 pass/fail は、初期表示しない。
- 調査用の詳細理由は、UI ではなく structured log または詳細展開へ寄せる。
- Completed job の出力準備詳細は、翻訳管理ではなく出力管理で扱う。

## HITL Status

- `detail_spec_hitl`: `required-after-design-bundle`
- `storybook_review_request`: `required-after-frontend-implementation`
- `frontend_human_review`: `required-after-frontend-implementation`
- `approval_record`: `pending-after-plan`

## Codex Implementation Result

- `completed_handoffs`: 未着手
- `touched_files`: 未着手
- `implemented_scope`: 未着手
- `test_results`: 未着手
- `implementation_investigation`: 未着手
- `ui_evidence`: 未着手
- `storybook_review_request`: 未着手
- `storybook_review_resources`: 未着手
- `approved_frontend_protection`: 未着手
- `codex_review_result`: 未着手
- `sonar_gate_result`: 未着手
- `residual_risks`: 現行画面の壊れ方は未観測である。実装前に agent-browser で現状確認が必要である。
- `docs_changes`: この plan のみ

## Merge Readiness

- `merge_ready`: `pending`
- `source_branch`: `codex/2026-05-16-translation-management-ui-subtraction`
- `target_branch`: `master`
- `commit_hash`: `N/A`
- `validation_evidence`: 未着手
- `review_evidence`: 未着手
- `residual_risks`: UI 削減が必要な状態説明まで削るリスクがある。`detail-spec-diff.md` で停止条件を固定する。

## Merge Result

- `merge_status`: `pending`
- `conflict_resolution`: `N/A`
- `post_merge_validation`: `N/A`
- `completed_move`: `N/A`
- `merge_commit_hash`: `N/A`
- `remote_operation`: `not-performed`

## Closeout Notes

- `canonicalized_artifacts`: 未着手
- `detail_spec_canonicalization`: 未判断
- `follow_up`: fake secret store 計画と合わせ、AI による UI 確認が安定する環境を先に用意する可能性がある。

## Outcome

- 計画作成のみ。
