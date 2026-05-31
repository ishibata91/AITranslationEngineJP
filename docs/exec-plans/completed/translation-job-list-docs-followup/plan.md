# translation-job-list-docs-followup

## 依頼要約

前 task `translation-job-list-fix`（merge commit d2dbf3a0、`docs/exec-plans/completed/translation-job-list-fix/` 配下）の実装は master 反映済みであるが、正本 docs への反映が未実施だった。本 task で docs-only 反映を進める。

反映対象は次の二系統である。

1. 画面設計正本反映: 翻訳管理ジョブカードの操作ボタン構成変更
   - 「再開」ActionButton（resumeOperation 用のボタン）を表示から削除
   - continue-button のラベルを「現在の翻訳段階へ進む」から「再開」へ変更
   - link の aria-label を「ジョブ N を選択して現在の翻訳段階へ進む」から「ジョブ N を選択して再開する」へ
   - ボタン直下に羅列されていた選択不可理由（disabled-reason 表示塊）を削除（hover tooltip は維持）

2. 詳細仕様正本反映: `docs/detail-specs/translation-job-management.md`
   - REQ-004 の理由分類列挙に `runtime_snapshot_missing` を追加
   - snapshot 欠落時の振る舞い表現を「snapshot が無いため再開できません／削除は可能」と直接的な記述に整理

実装は merge 済みのため、本 task ではコード変更は行わない。docs 正本本文だけを更新する。

## 分岐元

- 分岐元 branch: master
- 分岐元 commit: 7ecd340ce1773d139770601705762708fa846e3a
- 作業 branch: claude/translation-job-list-docs-followup

## 入力資料

- 前 task plan.md: `docs/exec-plans/completed/translation-job-list-fix/plan.md`
- merge commit: d2dbf3a0
- 反映実装の参照 commit: 0e16117

## 後続モジュール引き継ぎ

- task-id: translation-job-list-docs-followup
- 後続: docs-only 反映のため `finalization-module` の `正本化判断` と `詳細仕様正本反映` を直接呼ぶ流れを想定する（実装系の `design-module` / `implementation-module` は経由しない）。

## 正本化判断

- 人間承認状態: 承認済み（前 task の人間修正レビュー承認の延長）。
- 仕様変更または仕様追加がある: Y（画面設計と理由分類列挙の更新）。
- 反映対象 1（画面設計正本）: `docs/screen-design/screens/translation-job-management.md`
  - line 152 周辺: link の aria-label 文言を「ジョブ <ジョブID> を選択して現在の翻訳段階へ進む」から「ジョブ <ジョブID> を選択して再開する」へ。
  - line 170-177 周辺: 操作群の構成記述から「再開操作名」を削除し、「現在の翻訳段階へ進む」のラベル名と機能名を「再開」に統一する。「再開はジョブ再開を要求し、必要に応じて翻訳実行画面へ進む」の記述は機能としては実態と一致しないため、「再開」ボタン押下時は翻訳実行画面へ遷移する記述に整理する。
  - line 185-193 周辺: ボタン直下に「停止: <理由>」「再開: <理由>」「削除: <理由>」「翻訳段階: <理由>」を表示する記述を削除し、選択不可理由は hover tooltip でのみ表示する記述に整理する。
  - line 231 周辺: data-testid 表から `translation-job-management-resume-button` 行を削除する。
- 反映対象 2（詳細仕様正本）: `docs/detail-specs/translation-job-management.md`
  - line 60 の理由分類列挙「入力データ参照不能、終端状態、状態不整合」に「snapshot 欠落」を追加する。
  - REQ-005 周辺（AI 設定要約と認証状態の分離）に「runtime snapshot 欠落の時は再開可否を不可と判断する」相当の文を追加するかどうかは `docs_updater` の判断に委ねる（既存記述で説明可能なら追加しない）。
- 判定: `詳細仕様正本反映` を `docs_updater` 起動で進める。

## 詳細仕様正本反映結果

- 反映日: 2026-06-01
- 実行者: `docs_updater`
- 承認根拠: 前 task `translation-job-list-fix` plan.md「人間修正レビュー（論点1）」2026-05-31 承認、「合意済み frontend 保護（論点2, 3）」承認。本 task plan.md「正本化判断」セクションで人間承認状態「承認済み」と明記。

### 変更ファイル

- `docs/screen-design/screens/translation-job-management.md`
  - [8] ジョブ選択領域 の aria-label を「ジョブ <ジョブID> を選択して現在の翻訳段階へ進む」から「ジョブ <ジョブID> を選択して再開する」へ変更した。
  - [10] ジョブ操作 の概要、表示内容、操作、結果から「再開操作名」と「現在の翻訳段階へ進む」の旧表記を削除し、「再開」「再開は翻訳実行画面へ進む」に統一した。
  - [11] 無効理由 の表示内容を「ボタン直下への理由テキスト列挙は表示しない。hover tooltip でのみ表示する。」に整理した。
  - E2E 固定 selector 表から `translation-job-management-resume-button` 行を削除した。
- `docs/detail-specs/translation-job-management.md`
  - REQ-004 の理由分類列挙「入力データ参照不能、終端状態、状態不整合」に「snapshot 欠落」を追加した。

### 未反映理由

- REQ-005 への「snapshot 欠落時の再開不可判断」追記は行わなかった。理由: REQ-004 の理由分類に「snapshot 欠落」を追加したことで、再開不可時の snapshot 欠落理由は REQ-004 の範囲で説明できる。REQ-005 に同一事実を重複追記することは仕様文の冗長化になるため、追加しないと判断した。

### 残留リスク

- なし（E2E テスト観点正本反映は前 task から引き続き後続課題に切り出し中。本 task では docs-only 正本反映のみ担当）。


## 作業 commit

- commit hash: 3e05f53da23dc8f01e0586110d863717b9e1cf26
- 作業 branch: claude/translation-job-list-docs-followup
- 分岐元: master @ 7ecd340ce1773d139770601705762708fa846e3a
- 変更ファイル: 3 files changed, 84 insertions(+), 10 deletions(-)
  - `docs/screen-design/screens/translation-job-management.md`
  - `docs/detail-specs/translation-job-management.md`
  - `docs/exec-plans/active/translation-job-list-docs-followup/plan.md`
- 検証結果: docs-only のため harness 実行は省略（コード変更なし）。
- 残留リスク: なし（E2E 観点正本反映は前 task から継続して別 task に切り出し中）。

## マージ準備入力

- active plan folder: `docs/exec-plans/active/translation-job-list-docs-followup/`
- source branch: `claude/translation-job-list-docs-followup`
- target branch: `master`
- 作業 commit hash: `3e05f53da23dc8f01e0586110d863717b9e1cf26`
- 最終検証結果: docs-only のため harness 省略（コード変更なし）。
- 残留リスク: なし。

