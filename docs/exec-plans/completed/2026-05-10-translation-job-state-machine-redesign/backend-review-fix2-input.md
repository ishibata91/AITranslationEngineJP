# Backend レビュー修正 2 入力

- `task_id`: `2026-05-10-translation-job-state-machine-redesign`
- `handoff_to`: `backend_implementer`
- `skill`: `implement-backend`
- `status`: `ready`
- `created_at`: `2026-05-14`

## 修正理由

5 観点レビュー再実行で `behavior-003` が open になった。
本文翻訳 phase と単語翻訳 phase の read model が、terminal job でも phase run state だけで操作可能表示を返せる。

## 対象指摘

- file: `reviewback.behavior.yaml`
- issue id: `behavior-003`
- level: `major`
- status: `open`

## 修正対象

product code:
- `internal/service/body_translation_phase_service.go`
- `internal/service/term_translation_phase_service.go`

product test:
- 必要なら `internal/service/body_translation_phase_service_test.go`
- 必要なら `internal/service/term_translation_phase_service_test.go`

## 期待する修正

本文翻訳 phase と単語翻訳 phase の summary 操作可否を terminal job で false にする。
Service 実処理 guard と read model 表示値を一致させる。

保持する条件:
- 非 terminal job の `running` は pause 可能のままにする。
- 非 terminal job の `paused` は resume 可能のままにする。
- 非 terminal job の `paused` は cancel 可能のままにする。
- 非 terminal job の `recoverable_failed` は retry 可能のままにする。
- `RecoverableFailed` の resume は禁止のままにする。
- Service から `translationjobpolicy` を import しない。
- provider raw payload、secret、API key、prompt 全文、翻訳本文全文をログへ追加しない。

## 期待するテスト

本文翻訳 phase と単語翻訳 phase の service summary test に、terminal job かつ active phase run の組み合わせで操作不可になる確認を追加する。

確認対象:
- terminal job + running phase run: pause false。
- terminal job + paused phase run: resume false。
- terminal job + recoverable_failed phase run: retry false。
- 本文翻訳 phase は terminal job + paused phase run で cancel false。

## 禁止範囲

- frontend、Wails DTO、DB schema、migration を変更しない。
- docs 正本、`.codex`、作業計画文書を変更しない。
- 別 task `2026-05-13-notification-module-dependency-separation` を変更しない。
- 初回レビューで resolved になった修正を戻さない。

## 検証コマンド

- `gofmt -l internal/service`
- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite coverage`

## 返却内容

- backend 修正の完了、未完了、停止の判定。
- 変更ファイル。
- terminal job の read model 操作可否を false にした方法。
- 追加または更新したテスト。
- 実行した検証コマンドと結果。
- 未実行検証がある場合は未実行理由。
