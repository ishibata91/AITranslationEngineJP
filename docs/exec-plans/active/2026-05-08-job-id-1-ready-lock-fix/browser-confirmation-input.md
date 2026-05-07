# 実装後ブラウザ確認起動入力

## 確認 URL

`http://localhost:34115`

## 起動状態

- `npm run dev:wails:agent-browser` を起動済み。
- 既存 port が使用中だったが、`http://localhost:34115` は既存 Wails dev server として利用可能である。
- DB migration 適用後、job 1 の `translation/pending` placeholder は 0 件である。
- `TRANSLATION_JOB.id=1` は `state=ready`、`progress_percent=0` である。

## 操作経路

- `agent-browser open http://localhost:34115` で開く。
- 未完了 job 一覧または Job Management 相当の画面へ移動する。
- job ID 1 の表示を探す。
- job ID 1 の削除可否表示と理由文言を確認する。
- Job Run 表示へ移動できる導線がある場合は、job ID 1 を表示対象にする。
- 実行開始ボタンが有料 API や外部送信へ到達しそうな場合は押さない。

## 操作期待値

- job ID 1 は Ready または未実行相当として表示される。
- job ID 1 の削除不可理由に `phase 実行状態と job 状態が不整合` が表示されない。
- job ID 1 の削除ボタンが見える場合、状態不整合理由では無効化されない。
- Job Run 表示に進める場合、表示だけで job が Running に暗黙遷移しない。

## 禁止操作

- job ID 1 を実際に削除しない。
- 実行開始、再開、リトライなど、外部 provider に到達し得る操作を押さない。
- API key、credential、secret store key を表示または保存しない。

## 安全条件

- 有料 API 到達と外部送信を避ける。
- 破壊的操作を避ける。
- 確認は表示状態、ボタン有効状態、理由文言、console/network 異常の確認までに限定する。

## 証跡出力先

- `tmp/agent-browser/2026-05-08-job-id-1-ready-lock-fix/snapshot.txt`
- `tmp/agent-browser/2026-05-08-job-id-1-ready-lock-fix/errors.txt`
- `tmp/agent-browser/2026-05-08-job-id-1-ready-lock-fix/screenshot.png`
