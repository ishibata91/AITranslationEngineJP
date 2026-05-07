# 詳細仕様正本反映入力

- `task_id`: `2026-05-07-provider-settings-job-decoupling-implement`
- `created_at`: `2026-05-07T21:58:26+0900`
- `target_agent`: `docs_updater`

## 入力成果物

- `scenario-design.md`
- `ui-design.md`
- `implementation-scope.md`
- `final-validation.md`
- `review-summary.md`
- `canonicalization-decision.md`

## 更新対象

- `docs/detail-specs/translation-job-setup.md`
- `docs/detail-specs/ai-provider-settings-management.md`
- `docs/detail-specs/translation-job-management.md`
- `docs/detail-specs/persona-generation-phase.md`
- `docs/detail-specs/body-translation-phase.md`
- `docs/detail-specs/term-translation-phase.md`

## 更新方針

- human 承認済みの恒久仕様だけを反映する。
- Job Setup の公開境界から credential 参照実値、endpoint、secret store key、APIキー本文、外部 provider raw payload を除外する。
- AIサービス設定を endpoint と secret store 参照の共通正本として明示する。
- Ready job 実行開始時の provider settings 再解決を明示する。
- Running phase の job 側保存値を非 secret 要約へ限定する。
- `credential 参照` という表現は、表示可能な状態分類と、実値の非公開を区別して書く。

## 検証証跡

- `python3 scripts/harness/run.py --suite all`: `passed`
- 実行時刻: `2026-05-07T21:52:39+0900`
- system test: `9 passed`, `0 failed`
- frontend coverage: statements `68.1%`, lines `68.3%`
- backend coverage: statements `68.9%`, lines `68.5%`
- Sonar coverage: `70.6%`, line `71.7%`, branch `63.2%`
- Sonar security issues: `0`
- Sonar reliability issues: `0`
- Sonar maintainability HIGH issues: `0`
