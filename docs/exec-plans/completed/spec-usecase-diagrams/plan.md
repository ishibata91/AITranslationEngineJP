# Task Plan: spec-usecase-diagrams

- `workflow`: diagramming
- `status`: completed
- `lane_owner`: `diagrammer`
- `task_id`: `spec-usecase-diagrams`
- `task_mode`: diagramming-only
- `request_summary`: 画面設計からユースケースを Markdown と Mermaid で整理する
- `goal`: 画面操作から見える利用者目的、外部ツール、AIサービスとの関係を UC 図と UC 記述で確認できる状態にする
- `constraints`: プロダクトコード、プロダクトテストを変更しない
- `actor_policy`: UC 図と UC 記述のアクターは、画面を操作して目的を達成する人間に限定する。外部ツール、外部ファイル形式、AIサービスは、外部データまたは外部連携先として表す
- `close_conditions`: 画面別 Markdown に、システム境界付き UML風UC図、画面設計由来の個別UC、指定テンプレートの概要、条件、主シナリオ、代替シナリオ、例外シナリオ、完了条件が揃っている。翻訳管理だけは同一 Markdown 内でフェーズ別に分かれている
- `worktree_path`: `/Users/iorishibata/Repositories/AITranslationEngineJP`
- `source_branch`: `N/A`
- `target_branch`: `N/A`

## Artifact Index

- `usecase_dashboard`: `../../../usecases/uc-dashboard.md`
- `usecase_provider_settings`: `../../../usecases/uc-provider-settings.md`
- `usecase_master_dictionary`: `../../../usecases/uc-master-dictionary.md`
- `usecase_master_persona`: `../../../usecases/uc-master-persona.md`
- `usecase_translation_management`: `../../../usecases/uc-translation-management.md`
- `usecase_output_management`: `../../../usecases/uc-output-management.md`
- `detail_spec_diff`: `N/A`
- `screen_design_diff`: `N/A`
- `implementation_scope`: `N/A`
- `detail_spec_target`: `N/A`

## Routing Notes

- `required_reading`: `docs/screen-design/screens/dashboard.md`, `docs/screen-design/screens/provider-settings.md`, `docs/screen-design/screens/master-dictionary.md`, `docs/screen-design/screens/master-persona.md`, `docs/screen-design/screens/translation-management.md`, `docs/screen-design/screens/translation-job-management.md`, `docs/screen-design/screens/translation-input-review.md`, `docs/screen-design/screens/translation-job-setup.md`, `docs/screen-design/screens/job-run.md`, `docs/screen-design/screens/term-translation-phase.md`, `docs/screen-design/screens/persona-generation-phase.md`, `docs/screen-design/screens/body-translation-phase.md`, `docs/screen-design/screens/translation-complete.md`, `docs/screen-design/screens/output-management.md`
- `canonicalization_targets`: `N/A`
- `detail_spec_id`: `N/A`
- `validation_commands`: `git diff --check -- docs/exec-plans/active/spec-usecase-diagrams`

## HITL Status

- `detail_spec_hitl`: `not-required`
- `frontend_human_review`: `not-required`
- `approval_record`: `N/A`

## Outcome

- 画面設計を根拠に、画面単位のユースケース図とユースケース記述を Markdown として作成した。
- 作成した UC 成果物を `docs/usecases/` へ正本化した。
- 人間指摘により、翻訳管理 UC からハルシネーション由来の再構築系 UC と文言を削除した。

## Closeout

- `closed_at`: 2026-05-28
- `close_reason`: UC 正本を `docs/usecases/` に移動し、正本側の再構築仕様も削除済み。
- `validation`: `git diff --check -- docs/usecases docs/index.md docs/spec.md docs/detail-specs docs/screen-design/screens docs/exec-plans/active/spec-usecase-diagrams`
- `remaining_risk`: Mermaid のレンダリング確認は未実行。Markdown 構文と禁止文言の残存は検索で確認する。
