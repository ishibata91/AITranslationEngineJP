# 不足セレクタ（論点1 削除不可バグ E2E テスト観点差分）

## 概要

- 対象: 未完了ジョブ一覧の `[data-testid=translation-job-management-disabled-reason]` 内に並ぶ理由種別を、E2E から区別して検証できる selector が未決である。
- 判断: 追加候補テスト観点 E2E-UC-TJL-FIX-2 と E2E-UC-TJL-FIX-3 は「削除不可理由」と「再開不可理由」を独立に観察する必要があり、現状の単一 `disabled-reason` selector では区別できない。
- 根拠: 既存 E2E-UC-037, E2E-UC-038 も「停止できない理由」「再開不可理由」を期待値に含むが、子 selector 未決が共通の未決事項として残っている（`docs/e2e-test-design/test-design.csv` 既存行の備考でも子 selector 未決が明示済み）。

## 不足セレクタ一覧

| ID | 対象画面 | 対象要素 | 必要 selector | 関連テスト ID | 理由 |
| --- | --- | --- | --- | --- | --- |
| GAP-TJL-FIX-1 | 未完了ジョブ一覧 | ジョブカード内の削除不可理由表示要素 | `[data-testid=translation-job-management-delete-disabled-reason]` または `[data-testid=translation-job-management-disabled-reason][data-action=delete]` 相当 | E2E-UC-TJL-FIX-2 | 「真の `state_projection_inconsistent` warning が削除をブロックする」ことを、再開不可理由と独立に検証する必要があるため。 |
| GAP-TJL-FIX-2 | 未完了ジョブ一覧 | ジョブカード内の再開不可理由表示要素 | `[data-testid=translation-job-management-resume-disabled-reason]` または `[data-testid=translation-job-management-disabled-reason][data-action=resume]` 相当 | E2E-UC-TJL-FIX-3 | 「`runtime_snapshot_missing` を再開ブロック理由として引き続き考慮する」ことを、削除可否と独立に検証する必要があるため。 |
| GAP-TJL-FIX-3 | 未完了ジョブ一覧 | ジョブカード内の削除ボタン単体 | `[data-testid=translation-job-management-delete-button]` 相当 | E2E-UC-TJL-FIX-1, E2E-UC-TJL-FIX-2 | 既存 E2E-UC-020, E2E-UC-039 でも「削除ボタン単体の selector は未決」と記録済みであり、`button[name=削除]` 名前マッチに依存している現状を引き継いだまま追加観点を起こすため、本 task では既存未決事項として継続記録する。 |
