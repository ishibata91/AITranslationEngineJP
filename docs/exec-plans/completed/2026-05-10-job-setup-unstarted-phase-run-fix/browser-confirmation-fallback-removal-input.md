# ブラウザ確認入力: fallback 削除後

## 目的

`JOB_PHASE_RUN` 不在 fallback を削除した後の画面挙動を確認する。

## 確認 URL

`http://localhost:34115`

## 前提

- Wails dev server は起動済み。
- fake provider / dev provider を使える前提で確認する。
- 旧データ `JOB ID 1` は救済対象外である。

## 操作 1: 旧 job の失敗確認

1. 翻訳管理を開く。
2. `JOB ID 1` を選ぶ。
3. 現在 phase または単語翻訳 summary が表示される場所を開く。

期待:

- `JOB ID 1` は summary 取得失敗になる。
- `JOB_PHASE_RUN` 不在の旧 job が未開始として正常表示されない。

## 操作 2: 新規 JSON からの setup 確認

使用してよい JSON:

`docs/exec-plans/completed/exploration-normal-flow-20260503/normal-flow-lucien-mini.json`

操作:

1. JSON を入力データとして登録する。
2. 登録した入力データから新規 translation job を setup する。
3. setup 後の job を翻訳管理で選ぶ。
4. 現在 phase または単語翻訳 summary を表示する。

期待:

- 新規 job は summary 取得失敗にならない。
- 新規 job は `開始待ち` または未開始相当として表示される。
- 画面上に `pending` は表示されない。
- `次へ進む` は未開始時点では無効でよい。

## 禁止操作

- 実 provider へ有料 API 到達する操作は避ける。
- 削除、リトライ、再開、中断は押さない。
- 新規 job の phase start は、fake provider であることが画面上または設定上で確認できる場合だけ実施してよい。

## 証跡出力先

`tmp/agent-browser/2026-05-10-job-setup-unstarted-phase-run-fix/fallback-removal/`
