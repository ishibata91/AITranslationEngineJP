# ジョブID1 単語翻訳 summary 取得失敗 実装後ブラウザ確認入力

## 対象

- 呼び出し元: `fix_lane`
- 担当 agent: `browser_confirmation`
- 使用 skill: `browser-confirmation`
- 確認 URL: `http://localhost:34115`
- 証跡出力先: `tmp/agent-browser/2026-05-10-job-1-term-summary-resume-fix/`

## 起動状態

- `agent-browser doctor --offline --quick` は成功した。
- `npm run dev:wails:agent-browser` で Wails を起動中である。
- `curl -I --max-time 5 http://localhost:34115` は `HTTP/1.1 405 Method Not Allowed` を返した。
- 405 は `HEAD` への応答であり、localhost への接続は成立している。

## 操作経路

- `agent-browser open http://localhost:34115` でアプリを開く。
- ダッシュボードから翻訳管理を開く。
- 未完了ジョブ一覧でジョブID1を探す。
- ジョブID1の「現在の翻訳段階へ進む」操作入口を開く。
- 単語翻訳段階の画面状態を確認する。

## 操作期待値

- 未完了ジョブ一覧にジョブID1が表示される。
- ジョブID1は実行前の job として扱われる。
- ジョブID1の単語翻訳段階を開いても「単語翻訳段階の summary 取得に失敗しました。」は表示されない。
- 単語翻訳段階は実行前状態として表示される。
- next phase readiness は service error ではなく、未完了理由で blocked として扱われる。

## 禁止操作

- 翻訳実行を開始しない。
- AI provider への外部送信を発生させない。
- job、phase run、dictionary、snapshot を変更する操作を実行しない。
- DB を直接変更しない。

## 安全条件

- 確認は表示確認に限定する。
- 実行開始ボタン、retry、resume、cancel、delete は押さない。
- 有料 API 到達リスクがある操作は行わない。
- 外部送信リスクがある操作は行わない。
- 破壊的操作リスクがある操作は行わない。

## 必須証跡

- `snapshot`
- `errors`
- 必要な `screenshot`
- backend log: `tmp/logs/wails-dev.log`

## 根拠参照

- `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/human-observation.md`
- `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/investigation.md`
- `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/fix-execution-input.md`
- `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/implementation-evidence.md`
- `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/regression-test-evidence.md`
