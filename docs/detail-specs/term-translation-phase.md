# 詳細仕様: 単語翻訳フェーズ

- `detail_spec_id`: `term-translation-phase`
- `status`: `approved`
- `source_artifacts`: `docs/exec-plans/completed/term-translation-phase/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/plan.md`, `docs/exec-plans/active/2026-05-10-translation-job-state-machine-redesign/plan.md`, `docs/exec-plans/active/translation-job-step-target-list-panel/detail-spec-diff.md`, `docs/exec-plans/active/2026-05-25-phase-prompt-builder-boundary/detail-spec-diff.md`
- `implementation_artifacts`: `docs/exec-plans/completed/term-translation-phase/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/final-validation.md`
- `review_artifacts`: `docs/exec-plans/completed/term-translation-phase/reviewback.*.yaml`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/review-summary.md`

## 親要件と仕様

### `term-translation-phase-REQ-001` 単語翻訳フェーズを開始できる

親要件:
利用者は `Ready` の翻訳ジョブから単語翻訳フェーズを開始し、フェーズ状態と進行状況を判断できる。

仕様:
- 単語翻訳フェーズは、`Ready` の翻訳ジョブ、終端状態以外の翻訳ジョブ、実行中の翻訳段階なしの場合だけ開始する。
- `Ready` の翻訳ジョブは、単語翻訳フェーズ開始前には実行中の翻訳段階を持たない。
- 単語翻訳フェーズ開始が許可された時だけ、単語翻訳の実行単位を開始する。
- 翻訳ジョブは `Running` のまま維持し、単語翻訳フェーズの状態で完了、中断、回復可能失敗を区別する。
- フェーズ操作可否は、共通操作規則とフェーズ状態から決まる。

### `term-translation-phase-REQ-002` AI 設定を開始時に再解決する

親要件:
フェーズ開始と再試行は最新の AIサービス設定を参照し、秘密値を利用者向け情報から分離する。

仕様:
- フェーズ開始と再試行は、AIサービス設定から最新の接続先と認証状態を再解決する。
- 実行時に利用者が判断できる接続情報は、AIサービス、モデル、認証状態、実行方式、一括処理設定である。
- 秘密値、認証キー平文、復号可能な値、認証参照の実値、接続先、外部サービスとの生データ、翻訳本文全文は利用者向け情報の対象外にする。
- 利用者向け情報は、AIサービス、モデル、認証状態、実行方式、入力件数、出力件数、失敗分類、処理要約を対象にする。
- AIサービスへ渡す生成指示の全文は、障害調査用の同一性情報へ要約して扱う。
- `PromptDigest` は、AIサービスへ渡す生成指示の同一性を示す内部情報として扱う。
- `PromptDigest` は、生成指示の全文を復元できない値として扱う。
- `TERM_TRANSLATION_REQUEST_V1` は、単語翻訳フェーズの AIサービス要求形状を識別する印として扱う。
- `TERM_TRANSLATION_REQUEST_V1` は、利用者が選択する生成規則の版として扱わない。
- 運用上必要な要約は、保存済みの状態事実から導出する。

### `term-translation-phase-REQ-003` 共通辞書一致で AI 翻訳対象を決める

親要件:
共通辞書に完全一致する語は、フェーズ開始時の辞書参照に基づく置換対象として扱う。

仕様:
- 共通辞書はフェーズ開始時の辞書参照で固定する。
- 共通辞書に完全一致する語は辞書置換対象にする。
- 共通辞書の完全一致条件外の語は AI 翻訳対象にする。
- 共通辞書の対象範囲を判定した後に AI 翻訳対象語が 0 件でも、単語翻訳フェーズは `Completed` として扱う。
- AI 翻訳対象語が 0 件の場合は、AIサービス未実行を結果として判断できる。

### `term-translation-phase-REQ-004` AIサービス応答をジョブ内辞書へ反映する

親要件:
共通辞書対象外の用語と固有名詞は AI 翻訳対象とし、確定訳語を対象の翻訳ジョブ内辞書へ保存する。

仕様:
- 共通辞書対象外の用語と固有名詞は AIサービスへ送る。
- AIサービスへの実行単位は 1 対象語を基本とする。
- 一括処理を使う場合も、1 項目は 1 対象語に対応する。
- AIサービスへ渡す生成指示は、対象語、原文言語、訳文言語、応答対応に使う識別子を同じ実行単位に固定する。
- AIサービス応答は、対象語ごとに、原語と訳語の対応として検査する。
- 有効な応答は、要求した対象語と同じ原語を持ち、空ではない訳語を持つ応答である。
- AIサービスの有効な応答は、原語と訳語の対応を保持し、自動で確定訳語として扱う。
- 確定訳語は対象の翻訳ジョブ内辞書として保存する。
- 確定訳語は、対象の翻訳ジョブと単語翻訳の実行単位を追跡できる状態で保持する。
- 同一翻訳ジョブ、同一レコード種別、同一原語では一意の辞書項目として扱う。
- 別レコード種別の同一原語は別辞書項目として扱える。

### `term-translation-phase-REQ-005` 再開、リトライ、失敗を安全に扱う

親要件:
再開、再試行、開始再送では重複作成を避け、途中失敗を回復可能な状態として扱う。

仕様:
- 再開、再試行、開始再送では、同じ実行単位を継続する。
- 既存の辞書項目は維持し、未処理の用語だけ AI 翻訳対象にする。
- 再試行可能な失敗では、最新の失敗理由と進行状況を更新する。
- AIサービス失敗、応答不正、保存失敗がある場合、`RecoverableFailed` または `Failed` の対象にする。
- 不正応答、応答欠落、余分な応答、空訳語、対象語との不一致は対象語単位の失敗として扱う。
- 保存途中失敗では、単語翻訳フェーズを `RecoverableFailed` として扱う。
- AIサービスはジョブ設定で固定したサービスを使う。
- 終端状態の翻訳ジョブへの後書きは拒否結果にする。

### `term-translation-phase-REQ-006` 後続フェーズの開始条件を固定する

親要件:
単語翻訳フェーズ完了後だけ、後続フェーズの入力条件が成立する。

仕様:
- 単語翻訳フェーズ未完了、失敗中、辞書参照不能は、後続フェーズの開始拒否理由にする。
- 単語翻訳フェーズ完了後だけ、後続フェーズの入力条件が成立する。
- 利用者は現在の翻訳段階、進行状況、対象語件数、共通辞書一致件数、AI 翻訳対象語件数を判断できる。
- 単語翻訳結果は、確定訳語件数、ジョブ内辞書反映件数、置換対象件数、未一致件数、AIサービス設定の扱いを判断できる状態にする。
- 後続フェーズへ進む操作は、単語翻訳フェーズ完了と辞書参照成立後だけ成立する。

### `term-translation-phase-REQ-007` 操作と状態を利用者が判別できる

親要件:
利用者は開始、一時停止、再開、再試行、取り消しの可否と理由を判断できる。

仕様:
- 開始操作は `Ready` の翻訳ジョブかつ実行中の単語翻訳がない時だけ成立する。
- 一時停止は `Running` の時だけ成立する。
- 再開は `Paused` の時だけ成立する。
- 再試行は `RecoverableFailed` かつ再試行可能な失敗の時だけ成立する。
- 取り消しは `Paused` の時だけ成立する。
- 利用者は待機中、実行中、空完了、完了、一時停止中、回復可能失敗、阻害中を区別できる。
- 更新中は既存の状態情報を保持し、更新中であることを区別できる。
- 利用者は翻訳段階の状態と操作不可理由を判断できる。
- 利用者は、単語翻訳フェーズの処理対象名が共通辞書対象外の用語と固有名詞であり、AIサービスへ送り、確定訳語として翻訳ジョブ内辞書へ保存するものであることを判断できる。

## 根拠

- `term-translation-phase` の plan は `workflow_state: implementation-review-passed` である。
- design bundle は human approved であり、`implementation-scope.md` は `human_review_status: approved` である。
- 最終検証では旧設計 gate、frontend、backend、全体検証が pass している。
- 5 観点の reviewback はすべて `review_status: no_issue`、`must_fix_open: false`、`max_level: none` である。
- 翻訳ジョブステップ処理対象一覧表示パネルの詳細仕様差分は、2026-05-23 の人間設計レビュー承認と 2026-05-24 の Storybook フロント実装承認に基づいて反映済みである。
- 生成指示境界の詳細仕様差分は、2026-05-25 の人間設計レビュー承認に基づいて反映済みである。
