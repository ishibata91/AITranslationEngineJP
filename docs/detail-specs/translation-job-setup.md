# 詳細仕様: 翻訳ジョブ設定

- `detail_spec_id`: `translation-job-setup`
- `status`: `approved`
- `source_artifacts`: `docs/exec-plans/completed/translation-job-setup/plan.md`, `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/plan.md`, `docs/exec-plans/completed/2026-05-06-job-setup-input-cards/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/plan.md`
- `implementation_artifacts`: `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/plan.md`, `docs/exec-plans/completed/2026-05-06-job-setup-input-cards/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/final-validation.md`
- `review_artifacts`: `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/reviewback.*.yaml`, `docs/exec-plans/completed/2026-05-06-job-setup-input-cards/reviewback.*.yaml`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/review-summary.md`

## 親要件と仕様

### `translation-job-setup-REQ-001` 登録済み入力データから翻訳ジョブを作成できる

親要件:
利用者は登録済み入力データを選び、`Ready` の翻訳ジョブを作成できる。

仕様:
- 入力データ候補は、名称、出自、登録時点、規模を判断できる単位で扱う。
- 利用者はジョブ設定の対象入力データを選択できる。
- 既存の翻訳ジョブが参照している入力データも、ジョブ設定の入力候補として扱える。
- 既存の翻訳ジョブ情報は、ジョブ作成可否とは独立した参考情報として扱う。
- 同じ入力データの既存翻訳ジョブは参考情報であり、ジョブ作成の許可条件とは分離する。
- 3 つの翻訳段階で認証状態とモデル選択が充足している時だけ、翻訳ジョブを作成できる。

### `translation-job-setup-REQ-002` 翻訳ジョブ未作成の入力データを削除できる

親要件:
利用者はジョブ設定で、翻訳ジョブ未作成の入力データを削除できる。

仕様:
- 入力データ削除は、翻訳ジョブ未作成の入力データだけを対象にする。
- 入力データ削除の結果は、対象入力データと関連する入力キャッシュの削除である。
- 削除中の対象入力データは、重複操作の対象外にする。
- 削除成功後の対象入力データは、ジョブ設定の入力候補から外れる。
- 翻訳ジョブ削除は入力データを残すため、ジョブ設定の入力データ削除とは別の操作として扱う。

### `translation-job-setup-REQ-003` 翻訳段階ごとの AI 設定を固定できる

親要件:
利用者は単語翻訳、NPC ペルソナ生成、本文翻訳の 3 段階について AI 設定を判断し、翻訳ジョブ作成時に固定できる。

仕様:
- ジョブ設定は単語翻訳、NPC ペルソナ生成、本文翻訳の 3 つの翻訳段階を扱う。
- 各翻訳段階は AIサービス、モデル、実行方式、一括処理、認証状態を扱う。
- 作成後の設定内容には、翻訳段階ごとの AIサービス、モデル、認証状態、実行方式、一括処理の有無を含める。
- ジョブ設定の AI 設定と共通ペルソナ管理の AI 設定は分離する。
- 利用者は入力候補、共通辞書、共通ペルソナ、単語翻訳設定、NPC ペルソナ生成設定、本文翻訳設定、作成後の設定内容を判断できる。

### `translation-job-setup-REQ-004` モデル候補と古い選択を安全に扱う

親要件:
利用者は AIサービスごとのモデル候補を更新し、古い候補に基づく翻訳ジョブ作成を避けられる。

仕様:
- モデル候補は AIサービスごとに取得する。
- モデル候補の取得時点は、認証値や秘密値から分離した鮮度情報として扱う。
- 作成前確認と作成処理は、モデル候補の取得時点と利用者の選択が一致する時だけ成立する。
- 古いモデル候補に基づく選択は拒否結果にする。
- 認証キーが必要な AIサービスでは、認証キー設定済みの場合だけモデル候補を取得できる。
- モデル選択は、モデル候補取得成功後に成立する。
- モデル候補は、未更新、更新中、取得済み、取得失敗、認証不足で更新不可を別状態として扱う。

### `translation-job-setup-REQ-005` AIサービス固有の入力規則を扱う

親要件:
利用者は AIサービスごとの認証キー要否と一括処理可否に沿って翻訳ジョブ設定を選べる。

仕様:
- LM Studio は認証キーなしで利用できる AIサービスとして扱う。
- LM Studio の認証状態は、認証不足の警告対象外にする。
- 一括処理は利用者が明示した設定値として扱う。
- 一括処理の対象 AIサービスは Gemini と xAI である。
- 認証キー登録は、認証キーが必要で未設定の AIサービスを選んだ場合だけ成立する。
- 認証キー登録後のモデル候補は、再更新が必要な状態として扱う。

### `translation-job-setup-REQ-006` 秘密値と内部値を利用者向け情報から分離する

親要件:
利用者は AI 設定の状態を確認でき、秘密値本体や内部識別子は利用者向け情報から分離される。

仕様:
- ジョブ設定は、認証参照の実値、秘密値参照、接続先、認証キー本文を利用者向け情報の対象外にする。
- 認証キー文字列、秘密値、外部サービスの生データ、内部識別子、モデル候補の鮮度情報は利用者向け情報の対象外にする。
- 設定済み、認証キー未設定、モデル未選択、モデル候補取得失敗、モデル候補更新中は別状態として扱う。
- 入力確認、共通基盤、翻訳段階別設定、作成前確認、作成実行は、ジョブ設定の利用順序として扱う。

## 根拠

- `translation-job-setup-phase-provider-settings` は 2026-05-04 に人間設計レビュー承認済みである。
- `translation-job-setup-phase-provider-settings` の最終検証は通過済みである。
- `2026-05-06-job-setup-input-cards` の最終検証は通過済みである。
- `2026-05-06-job-setup-input-cards` の 5 観点 reviewback はすべて `must_fix_open: false` である。
