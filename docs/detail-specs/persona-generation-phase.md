# 詳細仕様: NPC ペルソナ生成フェーズ

- `detail_spec_id`: `persona-generation-phase`
- `status`: `approved`
- `source_artifacts`: `docs/exec-plans/completed/persona-generation-phase/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/plan.md`, `docs/exec-plans/active/2026-05-10-translation-job-state-machine-redesign/plan.md`
- `implementation_artifacts`: `docs/exec-plans/completed/persona-generation-phase/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/final-validation.md`
- `review_artifacts`: `docs/exec-plans/completed/persona-generation-phase/reviewback.*.yaml`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/review-summary.md`

## 親要件と仕様

### `persona-generation-phase-REQ-001` NPC ペルソナ生成フェーズを開始できる

親要件:
利用者は単語翻訳フェーズ完了後に NPC ペルソナ生成フェーズを開始できる。

仕様:
- NPC ペルソナ生成フェーズは、単語翻訳フェーズ完了後に開始する。
- NPC ペルソナ生成フェーズは、単語翻訳フェーズ `Completed`、終端状態以外の翻訳ジョブ、実行中の翻訳段階なしの場合だけ開始する。
- フェーズ操作可否は、共通操作規則とフェーズ状態から決まる。

### `persona-generation-phase-REQ-002` 生成対象と入力条件を固定する

親要件:
フェーズは NPC の原文発話、NPC 属性メタデータ、会話文脈、共通ペルソナ参照からペルソナ参照情報を作る。

仕様:
- 生成対象は、NPC 件数、入力種類、対象件数、共通ペルソナ一致有無、対象外理由を判断できる状態にする。
- 共通ペルソナ一致時は、翻訳ジョブのペルソナ参照を固定する。
- ペルソナ生成は 1 NPC を 1 実行単位とし、NPC 属性と会話文脈を同じ実行単位で扱う。
- 生成対象 0 件は `Completed` とし、対象 0 件、AIサービス未実行、空のペルソナ参照を判断できる状態にする。
- 生成対象は NPC レコード、翻訳対象項目、会話文脈、共通ペルソナ参照、ペルソナ参照情報で構成する。

### `persona-generation-phase-REQ-003` AI 設定を開始時に再解決する

親要件:
フェーズ開始と再試行は最新の AIサービス設定を参照し、秘密値を利用者向け情報から分離する。

仕様:
- AIサービス、モデル、実行方式、一括処理は、ジョブ設定の NPC ペルソナ生成用設定を継承する。
- フェーズ開始と再試行は、AIサービス設定から最新の接続先と認証状態を再解決する。
- 実行時に利用者が判断できる接続情報は、AIサービス、モデル、認証状態、実行方式、一括処理設定である。
- 利用者は入力件数、出力件数、失敗分類を判断できる。

### `persona-generation-phase-REQ-004` AIサービス結果をペルソナ参照へ採用する

親要件:
有効な AIサービス出力は本文翻訳フェーズが参照するペルソナ情報として採用される。

仕様:
- 有効な AIサービス出力は、翻訳ジョブ内ペルソナまたはペルソナ参照へ自動採用する。
- 利用者はペルソナ参照の成立状態、欠落件数、本文翻訳フェーズの開始可否を判断できる。
- 結果要約は、保護対象を含まない範囲で再判断できる情報として扱う。

### `persona-generation-phase-REQ-005` 失敗、再試行、終端状態の翻訳ジョブを安全に扱う

親要件:
AIサービス失敗、入力不備、保存失敗、途中状態は回復可能な失敗として扱える。

仕様:
- AIサービス失敗、不正応答、入力欠落、保存失敗がある場合、`RecoverableFailed` または `Failed` の対象にする。
- 一部 NPC 失敗時は成功分を維持し、フェーズは `RecoverableFailed` として未処理 NPC だけ再試行対象にする。
- 再送、再開、再試行では同じ実行単位を継続する。
- 成功済みペルソナは、同じ実行単位内で一意に扱う。
- 終端状態の翻訳ジョブでは、ペルソナ生成開始、ペルソナ採用、本文翻訳開始可否更新を拒否結果にする。

### `persona-generation-phase-REQ-006` 本文翻訳フェーズの開始条件を固定する

親要件:
本文翻訳フェーズは NPC ペルソナ生成フェーズ完了かつペルソナ参照成立の時だけ開始できる。

仕様:
- 本文翻訳フェーズの開始条件は、NPC ペルソナ生成フェーズ `Completed` かつペルソナ参照成立の時だけ成立する。
- ペルソナ未完了、失敗、参照不能は本文翻訳フェーズの開始拒否理由にする。
- 本文翻訳フェーズ開始は、NPC ペルソナ生成フェーズ `Completed` とペルソナ参照成立が両方満たされた時だけ成立する。

### `persona-generation-phase-REQ-007` 操作と状態を利用者が判別できる

親要件:
利用者はフェーズの進行、生成対象、生成結果、ペルソナ参照状態、本文翻訳フェーズの開始可否を判断できる。

仕様:
- 利用者は現在の翻訳段階が NPC ペルソナ生成であることを判断できる。
- 利用者は進行状況、対象件数、生成済み件数、失敗件数、対象外件数を判断できる。
- 一時停止は `Running` の時だけ成立する。
- 再開は `Paused` の時だけ成立する。
- 再試行は `RecoverableFailed` かつ再試行可能な失敗の時だけ成立する。
- 取り消しは `Paused` の時だけ成立する。
- 利用者は翻訳段階の状態、再試行可否、本文翻訳開始可否、進行状況を判断できる。

### `persona-generation-phase-REQ-008` 保護対象を利用者向け情報から分離する

親要件:
NPC ペルソナ生成フェーズは秘密値、生成指示の原文、原文発話全文、会話文脈全文を利用者向け情報から分離する。

仕様:
- 秘密値、認証キー平文、認証参照の実値、接続先、外部サービスとの生データ、生成指示の原文、原文発話全文、会話文脈全文は利用者向け情報の対象外にする。
- 運用上必要な要約は、保存済みの状態事実から導出する。
- 導出する要約は、識別子、件数、根拠参照、保護対象を含まない結果要約を対象にする。
- 障害調査用の要約では、AIサービス、モデル、実行方式、一括処理、認証状態、入出力件数、失敗分類を判断できる。

## 根拠

- 検証結果は pass である。
- 最終検証は通過済みである。
- 5 観点レビューはすべて `review_status: no_issue`、`must_fix_open: false`、`max_level: none` である。
