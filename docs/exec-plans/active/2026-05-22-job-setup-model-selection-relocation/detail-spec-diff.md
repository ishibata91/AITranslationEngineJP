# 詳細仕様差分: 2026-05-22-job-setup-model-selection-relocation

- `skill`: detail-spec-design
- `status`: ready-for-human-review
- `source_plan`: `./plan.md`
- `detail_spec_target`: `docs/detail-specs/translation-job-setup.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`
- `screen_design_diff`: `./screen-design-diff.translation-management.md`, `./screen-design-diff.translation-input-review.md`, `./screen-design-diff.translation-job-setup.md`, `./screen-design-diff.term-translation-phase.md`, `./screen-design-diff.persona-generation-phase.md`, `./screen-design-diff.body-translation-phase.md`
- `component_diagram`: `./design-diff.job-setup-model-selection-relocation.md`

## 詳細仕様差分

### `translation-job-setup-REQ-001` 登録済み入力データから翻訳ジョブを作成できる

- `変更種別`: 変更
- `正本反映先`: `docs/detail-specs/translation-job-setup.md`

親要件:
利用者は登録済み入力データを選び、AI モデル選択を先に完了しなくても翻訳ジョブを作成できる。

仕様:
- 翻訳ジョブ作成は、入力データ候補の選択と入力データの利用可能状態で成立する。
- 翻訳ジョブ作成は、単語翻訳、NPC ペルソナ生成、本文翻訳のすべての AI モデル選択完了を条件にしない。
- 作成された翻訳ジョブは、単語翻訳を開始できる前段階として扱う。
- 同じ入力データの既存翻訳ジョブは参考情報として扱い、ジョブ作成の許可条件とは分離する。

未決:
- `Q-001`: なし

回答:
- `Q-001`: なし
  - `回答者`: N/A
  - `回答日`: N/A
  - `反映先`: N/A

### `translation-job-setup-REQ-003` 翻訳段階ごとの AI 設定を固定できる

- `変更種別`: 変更
- `正本反映先`: `docs/detail-specs/translation-job-setup.md`

親要件:
利用者は翻訳段階ごとの AI 設定を、ジョブ作成時ではなく対象段階の開始前に固定できる。

仕様:
- ジョブ作成は、翻訳段階ごとの AI サービス、モデル、実行方式、一括処理をまとめて固定しない。
- 単語翻訳の AI 設定は、単語翻訳を開始する前に固定する。
- NPC ペルソナ生成の AI 設定は、NPC ペルソナ生成を開始する前に固定する。
- 本文翻訳の AI 設定は、本文翻訳を開始する前に固定する。
- 作成後のジョブ情報は、未設定の翻訳段階と設定済みの翻訳段階を区別できる。

未決:
- `Q-002`: なし

回答:
- `Q-002`: なし
  - `回答者`: N/A
  - `回答日`: N/A
  - `反映先`: N/A

### `term-translation-phase-REQ-002` AI 設定を開始時に再解決する

- `変更種別`: 変更
- `正本反映先`: `docs/detail-specs/term-translation-phase.md`

親要件:
利用者は単語翻訳フェーズの開始前に AI 設定を選び、開始時と再試行時に最新の AI サービス設定を参照できる。

仕様:
- 単語翻訳フェーズは、単語翻訳用の AI サービス、モデル、実行方式、一括処理設定が固定された場合だけ開始できる。
- 単語翻訳フェーズ開始と再試行は、単語翻訳用の AI 設定を使う。
- 単語翻訳フェーズ開始と再試行は、AI サービス設定から最新の接続先と認証状態を再解決する。
- 利用者は単語翻訳用の AI サービス、モデル、認証状態、実行方式、一括処理設定を判断できる。
- 秘密値、認証キー平文、復号可能な値、認証参照の実値、接続先、外部サービスとの生データ、翻訳本文全文は利用者向け情報の対象外にする。

未決:
- `Q-004`: なし

回答:
- `Q-004`: なし
  - `回答者`: N/A
  - `回答日`: N/A
  - `反映先`: N/A

### `persona-generation-phase-REQ-003` AI 設定を開始時に再解決する

- `変更種別`: 変更
- `正本反映先`: `docs/detail-specs/persona-generation-phase.md`

親要件:
利用者は NPC ペルソナ生成フェーズの開始前に AI 設定を選び、開始時と再試行時に最新の AI サービス設定を参照できる。

仕様:
- NPC ペルソナ生成フェーズは、NPC ペルソナ生成用の AI サービス、モデル、実行方式、一括処理設定が固定された場合だけ開始できる。
- NPC ペルソナ生成フェーズ開始と再試行は、NPC ペルソナ生成用の AI 設定を使う。
- NPC ペルソナ生成フェーズ開始と再試行は、AI サービス設定から最新の接続先と認証状態を再解決する。
- 利用者は NPC ペルソナ生成用の AI サービス、モデル、認証状態、実行方式、一括処理設定を判断できる。
- 秘密値、認証キー平文、認証参照の実値、接続先、外部サービスとの生データ、生成指示の原文、原文発話全文、会話文脈全文は利用者向け情報の対象外にする。

未決:
- `Q-005`: なし

回答:
- `Q-005`: なし
  - `回答者`: N/A
  - `回答日`: N/A
  - `反映先`: N/A

### `body-translation-phase-REQ-001` 本文翻訳フェーズを開始できる

- `変更種別`: 変更
- `正本反映先`: `docs/detail-specs/body-translation-phase.md`

親要件:
利用者は本文翻訳フェーズの開始前に AI 設定を選び、開始時と再試行時に最新の AI サービス設定を参照できる。

仕様:
- 本文翻訳フェーズは、本文翻訳用の AI サービス、モデル、実行方式、一括処理設定が固定された場合だけ開始できる。
- 本文翻訳フェーズ開始と再試行は、本文翻訳用の AI 設定を使う。
- 本文翻訳フェーズ開始と再試行は、AI サービス設定から最新の接続先と認証状態を再解決する。
- 利用者は本文翻訳用の AI サービス、モデル、認証状態、実行方式、一括処理設定を判断できる。
- 秘密値、認証キー平文、復号可能値、認証参照の実値、接続先、外部サービスとの生データ、生成指示の原文は利用者向け情報の対象外にする。

未決:
- `Q-006`: なし

回答:
- `Q-006`: なし
  - `回答者`: N/A
  - `回答日`: N/A
  - `反映先`: N/A

## 根拠

- `source`: `./task-frame.md`, `./plan.md`
- `source`: `docs/detail-specs/translation-job-setup.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`
- `source`: `docs/screen-design/screens/translation-input-review.md`, `docs/screen-design/screens/translation-job-setup.md`, `docs/screen-design/screens/term-translation-phase.md`, `docs/screen-design/screens/persona-generation-phase.md`, `docs/screen-design/screens/body-translation-phase.md`
- `review`: 人間設計レビュー待ち
- `validation`: プロダクト検証は未実行。設計成果物のみ作成した。
