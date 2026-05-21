# 詳細仕様: AIサービス設定管理

- `detail_spec_id`: `ai-provider-settings-management`
- `status`: `approved`
- `source_artifacts`: `docs/exec-plans/completed/ai-provider-settings-management/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/plan.md`
- `implementation_artifacts`: `docs/exec-plans/completed/ai-provider-settings-management/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/final-validation.md`
- `review_artifacts`: `docs/exec-plans/completed/ai-provider-settings-management/reviewback.*.yaml`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/review-summary.md`

## 親要件と仕様

### `ai-provider-settings-management-REQ-001` AIサービスごとの接続設定を管理できる

親要件:
利用者は Gemini、xAI、LM Studio の接続設定を AIサービスごとに管理できる。

仕様:
- 管理対象の AIサービスは Gemini、LM Studio、xAI である。
- AIサービス設定は、サービスごとの接続先と認証状態を扱う。
- モデル、実行方式、一括処理の選択は、翻訳ジョブ作成時の設定として扱う。
- 利用者は、接続先、認証状態、接続確認状態、設定変更の結果をサービスごとに判断できる。
- 利用者は、接続先変更、認証キー保存、接続確認、未設定化をサービスごとに実行できる。

### `ai-provider-settings-management-REQ-002` 秘密値本体を利用者向け情報から分離する

親要件:
利用者は認証キー本体や外部サービスとの生データではなく、認証状態を判断できる。

仕様:
- 認証キー本体は、利用者が参照する設定情報から分離する。
- 利用者が判断できる認証情報は、認証キーの有無と利用可能性である。
- 認証キー、外部サービスへ送った内容、外部サービスから返った生データ、生成指示の原文は利用者向け情報の対象外にする。
- 接続先は、外部サービスへ到達するための設定値として扱う。
- 接続確認失敗、設定保存失敗、接続先参照不能、外部サービス不正応答は、失敗分類として区別できる。

### `ai-provider-settings-management-REQ-003` 設定変更時の確認状態を正しく扱う

親要件:
利用者は接続先や認証キーの変更後に、現在設定の接続確認状態を判断できる。

仕様:
- 接続先変更後の接続確認状態は未確定になる。
- 接続先または認証キーを変更した設定は、再確認待ちとして扱う。
- 未設定化した AIサービスは、管理対象として残る。
- 未設定化した AIサービスは、接続先と認証状態が未設定になる。
- 設定変更結果と接続確認結果は、成功、失敗、未確認を区別できる。

### `ai-provider-settings-management-REQ-004` 参照側は AIサービス設定を唯一の接続情報にする

親要件:
ジョブ設定、翻訳フェーズ、共通ペルソナ管理は AIサービス設定を唯一の接続情報として参照する。

仕様:
- ジョブ設定と共通ペルソナ管理は、AIサービス設定を接続情報の参照元にする。
- `Ready` の翻訳ジョブを開始または再試行する時は、最新の接続先と認証状態を参照する。
- 翻訳実行時に利用者が判断できる接続情報は、AIサービス、モデル、認証状態、実行方式、一括処理設定である。
- 翻訳実行時に解決した接続先と秘密値は、利用者向け情報の対象外にする。

### `ai-provider-settings-management-REQ-005` 実行前の AI サービス設定を判断できる

親要件:
利用者は実行前に AIサービス接続の成立状態を判断できる。

仕様:
- 接続確認はサービスごとの接続先と認証状態を対象にする。
- 接続確認結果は成功、失敗、未確認の状態として扱う。
- モデルと一括処理の切り替えは、翻訳ジョブ作成時の翻訳段階設定で扱う。

## 根拠

- `Q-AIPSM-001` から `Q-AIPSM-006` は人間回答済みである。
- 5 観点 reviewback は plan 上で `no_issue` として記録されている。
