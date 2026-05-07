# 詳細仕様正本反映結果

- `task_id`: `2026-05-07-provider-settings-job-decoupling-implement`
- `agent`: `docs_updater`
- `skill`: `updating-docs`
- `result`: `completed`

## 根拠

- 承認記録: `canonicalization-decision.md`
- 実装完了根拠: `final-validation.md`
- レビュー完了根拠: `review-summary.md`
- 実装範囲根拠: `implementation-scope.md`

## 変更ファイル

- `docs/detail-specs/translation-job-setup.md`
- `docs/detail-specs/ai-provider-settings-management.md`
- `docs/detail-specs/translation-job-management.md`
- `docs/detail-specs/persona-generation-phase.md`
- `docs/detail-specs/body-translation-phase.md`
- `docs/detail-specs/term-translation-phase.md`

## 反映内容

- Job Setup の各翻訳段階が provider、model、execution mode、batch mode、APIキー状態を扱うことを反映した。
- Job Setup の公開 DTO、UI、作成後 summary へ credential 参照実値、secret store key、endpoint、APIキー本文を出さないことを反映した。
- provider model list の `sourceToken` は非 secret 鮮度 token とし、credential 参照実値、secret store key、endpoint、secret を含めないことを反映した。
- 作成前検証と作成処理が非 secret 鮮度 token を使い、古いモデル一覧由来の選択を stale として拒否することを反映した。
- Ready job の実行開始と retry が AIサービス設定から最新 endpoint と credential 参照状態を再解決することを反映した。
- Running phase の job 側 runtime snapshot は provider、model、credential 状態分類、execution mode、batch mode だけを保存することを反映した。

## 未反映理由

- 有料の実 AI API 呼び出しは必須検証にしない方針のため、詳細仕様正本へ追加していない。
- provider SDK 実装方式、migration 番号、repository owner は実装詳細のため、詳細仕様正本へ追加していない。
- review URL、fakeAPI 実行手順、検証ログ時刻は task-local 証跡のため、詳細仕様正本へ追加していない。

## 検証

- `python3 scripts/harness/run.py --suite structure`: `passed`
- 未解決の検証失敗: `なし`
