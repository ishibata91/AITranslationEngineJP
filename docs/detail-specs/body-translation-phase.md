# 詳細仕様: 本文翻訳フェーズ

- `detail_spec_id`: `body-translation-phase`
- `status`: `approved`
- `source_artifacts`: `docs/exec-plans/completed/body-translation-phase/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/plan.md`, `docs/exec-plans/active/2026-05-10-translation-job-state-machine-redesign/plan.md`, `docs/exec-plans/active/translation-job-step-target-list-panel/detail-spec-diff.md`, `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/detail-spec-diff.md`
- `implementation_artifacts`: `docs/exec-plans/completed/body-translation-phase/implementation-scope.md`, `docs/exec-plans/completed/body-translation-phase/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/final-validation.md`
- `review_artifacts`: `docs/exec-plans/completed/body-translation-phase/reviewback.behavior.yaml`, `docs/exec-plans/completed/body-translation-phase/reviewback.contract.yaml`, `docs/exec-plans/completed/body-translation-phase/reviewback.trust-boundary.yaml`, `docs/exec-plans/completed/body-translation-phase/reviewback.state-invariant.yaml`, `docs/exec-plans/completed/body-translation-phase/reviewback.responsibility-boundary.yaml`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/review-summary.md`

## 親要件と仕様

### `body-translation-phase-REQ-001` 本文翻訳フェーズを開始できる

親要件:
利用者は NPC ペルソナ生成フェーズ完了後に、本文翻訳フェーズを開始できる。

仕様:
- 本文翻訳フェーズは、ジョブ設定で固定した本文翻訳用の AIサービス、モデル、実行方式、一括処理設定を使う。
- 本文翻訳フェーズは、開始条件が成立した時だけ開始する。
- 開始後のフェーズ操作可否は、共通操作規則とフェーズ状態から決まる。
- フェーズ開始と再試行は、AIサービス設定から最新の接続先と認証状態を再解決する。
- 実行時に利用者が判断できる接続情報は、AIサービス、モデル、認証状態、実行方式、一括処理設定である。

### `body-translation-phase-REQ-002` 翻訳入力を固定する

親要件:
本文翻訳フェーズは確定訳語、ジョブ内辞書、ジョブ内ペルソナ、翻訳補助情報を同一実行単位の入力として固定する。

仕様:
- 入力条件は対象翻訳項目件数、辞書参照、ペルソナ参照、翻訳補助情報、生成指示の同一性を持つ。
- 完全一致した辞書 hit は辞書置換対象にする。
- AI 翻訳対象は辞書置換対象外の翻訳項目にする。
- 部分一致は訳語固定制約として AIサービスへ渡す。
- AIサービスへ渡す本文翻訳の生成指示は、1 翻訳項目を 1 実行単位として扱う。
- 生成指示は、翻訳項目の対応識別子、レコード種別、フィールド種別、原文、ペルソナ要約、翻訳補助情報、完全一致の辞書置換結果、部分一致の訳語固定制約、保持要素を同じ実行単位に固定する。
- 翻訳レコード種別と翻訳項目種別に応じて翻訳指示を構成する。
- 翻訳項目の対応関係と保護要素の同一性を失わず AIサービス境界へ渡す。
- 生成指示は、翻訳項目の対応関係と保持要素の同一性を AIサービス境界で判断できる状態にする。
- `PromptDigest` は、AIサービスへ渡す生成指示の同一性を示す内部情報として扱う。

### `body-translation-phase-REQ-003` 翻訳結果を翻訳項目単位で保存できる

親要件:
本文翻訳フェーズは訳文、出力状態、保護要素検証結果を翻訳項目単位で保持する。

仕様:
- 訳文、出力状態、保護要素検証結果は同一翻訳項目に対応付ける。
- AIサービス応答は、翻訳項目ごとに、要求単位、翻訳項目の対応識別子、レコード種別、フィールド種別、訳文、保持要素の同一性として検査する。
- 有効な応答は、要求した翻訳項目と同じ対応識別子を持ち、空ではない訳文を持ち、保持要素の同一性を満たす応答である。
- 有効な応答は、翻訳項目単位の訳文と出力状態の候補として採用する。
- 保護要素検証に失敗した訳文は、採用結果の対象外にする。
- 保持要素検証に失敗した訳文は、翻訳項目単位の失敗分類として扱う。
- 応答欠落、余分な応答、翻訳項目との不一致、空訳文、保持要素の欠落、変更、重複、順序変更、余分な追加は、翻訳項目単位の失敗分類として扱う。
- 保護要素検証に失敗した翻訳項目は再試行対象にする。
- 保存失敗または検証失敗がある翻訳段階は、`RecoverableFailed` または `Failed` の対象にする。
- 本文翻訳フェーズ `Completed`、翻訳結果整合、出力状態整合を満たす時だけ成果物出力条件が成立する。
- 本文翻訳対象 0 件は Completed として扱う。
- AIサービス未実行でも、単語だけの翻訳対象は成果物出力へ進める。

### `body-translation-phase-REQ-004` 失敗、再試行、取り消しを安全に扱う

親要件:
AIサービス失敗、応答不正、対応関係不整合、保存失敗、保護要素検証失敗は回復可能な失敗として扱える。

仕様:
- AIサービス失敗、応答不正、対応関係不整合、保存失敗、保護要素検証失敗がある場合、`RecoverableFailed` または `Failed` の対象にする。
- 部分失敗では成功済み翻訳結果を保持し、翻訳段階全体は `RecoverableFailed` として扱う。
- 再試行、再開、開始再送は同じ実行単位を継続する。
- 同一翻訳項目の結果は、同じ実行単位内で一意に扱う。
- 取り消しは `Paused` の時だけ成立する。
- `Canceled` 後はフェーズ終端とし、途中成功結果は成果物出力条件の対象外にする。
- 終端状態の翻訳ジョブでは、本文翻訳フェーズの開始、翻訳結果採用、成果物出力条件更新、遅延応答採用を拒否結果にする。

### `body-translation-phase-REQ-005` 翻訳ジョブ全体の完了と成果物出力条件を作る

親要件:
本文翻訳フェーズの完了時点で翻訳ジョブ全体は `Completed` になり、完了済み翻訳ジョブから成果物を出力できる。

仕様:
- 本文翻訳フェーズの完了時点で翻訳ジョブ全体は `Completed` になる。
- 本文翻訳フェーズ `Completed`、翻訳結果整合、出力状態整合を満たす時だけ成果物出力条件が成立する。
- 成果物出力は、成果物出力条件が成立した完了済み翻訳ジョブだけを候補にする。
- 利用者は、翻訳完了確認の処理対象名が翻訳項目単位の訳文であり、本文翻訳で保持された訳文として出力管理へ進む前に確認するものであることを判断できる。

### `body-translation-phase-REQ-006` 進行と結果を判断できる

親要件:
利用者は本文翻訳フェーズの進行、結果、失敗理由、成果物出力条件を判断できる。

仕様:
- 利用者は現在の翻訳段階、進行状況、対象件数、処理済み件数、未処理件数、AIサービス設定の扱い、出力件数を判断できる。
- 利用者は辞書適用件数、ペルソナ参照件数、翻訳補助情報、訳文、出力状態、保護要素検証結果、成果物出力条件を判断できる。
- 利用者は失敗状態、失敗分類、再試行可否、影響件数、保護対象を含まない失敗要約を判断できる。
- 利用者は本文翻訳フェーズ開始、一時停止、再開、再試行、取り消し、翻訳結果確認、保護要素検証結果確認、成果物出力条件確認を行える。
- 状態差分は、未準備、準備完了、開始中、実行中、一時停止中、回復可能失敗、検証失敗、空完了、完了、取り消し済み、失敗として扱う。
- AIサービス未実行、AIサービス実行中、AIサービス部分失敗、保存失敗、遅延応答拒否を区別できる。
- 利用者は、本文翻訳フェーズの処理対象名が辞書置換対象外の翻訳項目であり、辞書とペルソナを参照して訳文を作るものであることを判断できる。

### `body-translation-phase-REQ-007` 操作可否と状態理由を判断できる

親要件:
利用者は操作可否と理由を判別できる。

仕様:
- 開始は開始条件が成立した時だけ成立する。
- 一時停止は本文翻訳フェーズ `Running` の時だけ成立する。
- 再開は本文翻訳フェーズ `Paused` の時だけ成立する。
- 再試行は本文翻訳フェーズ `RecoverableFailed` かつ再試行可能な失敗がある時だけ成立する。
- 取り消しは本文翻訳フェーズ `Paused` の時だけ成立する。
- 本文翻訳フェーズ `Running` の取り消し要求は拒否結果になる。
- 成果物出力条件は本文翻訳フェーズ `Completed` かつ翻訳結果整合時だけ成立する。
- 利用者は翻訳段階の状態、検証結果、再試行可否、進行状況を判断できる。

### `body-translation-phase-REQ-008` 秘密値と生データを利用者向け情報から分離する

親要件:
本文翻訳フェーズは秘密値、外部サービスとの生データ、生成指示の原文を利用者向け情報から分離する。

仕様:
- 秘密値、認証キー平文、復号可能値、認証参照の実値、接続先、外部サービスとの生データ、生成指示の原文は利用者向け情報の対象外にする。
- 運用上必要な要約は、保存済みの状態事実から導出する。
- 原文と訳文は、翻訳結果を確認するための情報として扱える。
- 利用者向け情報は、AIサービス、モデル、実行方式、認証状態、入力件数、出力件数、失敗分類、保持要素の件数、保持要素検証結果、保護対象を含まない結果要約を対象にする。
- AIサービスへ渡す生成指示の全文と外部サービスとの生データは、障害調査用の同一性情報へ要約して扱う。
- `PromptDigest` は、AIサービスへ渡す生成指示の同一性を示す内部情報として扱う。
- `PromptDigest` は、生成指示の全文または外部サービスとの生データを復元できない値として扱う。
- `BODY_TRANSLATION_REQUEST_V1` は、本文翻訳フェーズの AIサービス要求形状を識別する印として扱う。
- `BODY_TRANSLATION_REQUEST_V1` は、利用者が選択する生成規則の版として扱わない。

## 根拠

- human decision は plan の `human_review_status: approved-after-design-bundle` と人間設計レビュー結果 `approved` に記録されている。
- 最終検証は plan の最終検証通過結果で確認済みである。
- 5 観点 reviewback はすべて `review_status: no_issue`、`must_fix_open: false`、`max_level: none` である。
- 翻訳ジョブステップ処理対象一覧表示パネルの詳細仕様差分は、2026-05-23 の人間設計レビュー承認と 2026-05-24 の Storybook フロント実装承認に基づいて反映済みである。
- 生成指示境界の詳細仕様差分は、2026-05-25 の人間設計レビュー承認に基づいて反映済みである。
