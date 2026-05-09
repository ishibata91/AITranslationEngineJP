# 修正前調査: 再起動あり再現調査

## 判断結果

- 判定: 部分完了
- 調査 mode: `修正前調査`, `UI 根拠`
- 引き継ぎ先: `fix_lane`
- 判断: 明示再起動後の初回操作と再実行操作の証跡は取得した。
- 判断: 人間観測にある `ジョブの進み方` 表示と下部 panel `ジョブ未選択` 状態は、今回の `agent-browser` 操作では再現できなかった。
- 判断: 初回操作と再実行操作の route、DOM、console、network、backend log は、今回採取分では差分を確認できなかった。

## 根拠参照

- 人間観測: `./human-observation.md`
- 起動入力: `./investigation-input.md`
- 既存調査:
  - `./pre-fix-investigation.md`
  - `./pre-fix-investigation.supplemental-term-ui.md`
- 再起動証跡:
  - `tmp/logs/2026-05-09-job-phase-first-open-blank.restart-reproduction.command.log`
  - `tmp/logs/2026-05-09-job-phase-first-open-blank.restart-reproduction.startup.log`
  - `tmp/logs/2026-05-09-job-phase-first-open-blank.restart-reproduction.doctor.log`
- UI 証跡:
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/initial.url.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/initial.snapshot.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/initial.png`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/first-open.url.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/first-open.snapshot.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/after-first-attempt.png`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/reopen.url.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/reopen.snapshot.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/reopen.png`
- frontend log:
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/initial.console.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/after-first-attempt.console.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/reopen.console.txt`
- network:
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/initial.network.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/after-first-attempt.network.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/reopen.network.txt`
- backend log:
  - `tmp/logs/2026-05-09-job-phase-first-open-blank.restart-reproduction.before-first-open.backend.log`
  - `tmp/logs/2026-05-09-job-phase-first-open-blank.restart-reproduction.first-open.backend.log`
  - `tmp/logs/2026-05-09-job-phase-first-open-blank.restart-reproduction.after-first-attempt.backend.log`
  - `tmp/logs/2026-05-09-job-phase-first-open-blank.restart-reproduction.reopen.backend.log`
  - `tmp/logs/2026-05-09-job-phase-first-open-blank.restart-reproduction.after-reopen.backend.log`

## 再起動条件

- 観測事実: 調査開始前の `curl http://127.0.0.1:34115` は接続失敗だった。根拠: `tmp/logs/2026-05-09-job-phase-first-open-blank.restart-reproduction.command.log`
- 観測事実: `sh ./scripts/dev/run-wails-agent-browser.sh` を 2026-05-09 23:49 JST に起動した。根拠: `tmp/logs/2026-05-09-job-phase-first-open-blank.restart-reproduction.command.log`
- 観測事実: 起動 log には `Using DevServer URL: http://0.0.0.0:34115` が残っている。根拠: `tmp/logs/2026-05-09-job-phase-first-open-blank.restart-reproduction.startup.log`
- 観測事実: `agent-browser doctor --offline --quick` は pass だった。根拠: `tmp/logs/2026-05-09-job-phase-first-open-blank.restart-reproduction.doctor.log`

## 観測点

- 入口: `http://127.0.0.1:34115/#translation-management`
- 対象 job: `jobID1`, 画面表示では `ジョブ 1`
- 初回操作:
  - `button "現在の翻訳段階へ進む"` を押下
  - 補助観測として `link "ジョブ 1 を選択して現在の翻訳段階へ進む"` も押下
- 再実行操作:
  - 一覧画面のまま同じ `button "現在の翻訳段階へ進む"` を再押下

## 観測事実

### 初回表示

- 観測事実: 再起動直後の画面には `ジョブの進み方`、`未完了ジョブ一覧`、`ジョブ 1 を選択して現在の翻訳段階へ進む`、`現在の翻訳段階へ進む` が表示されていた。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/initial.snapshot.txt`
- 観測事実: 再起動直後 URL は `http://127.0.0.1:34115/#translation-management` だった。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/initial.url.txt`

### 初回操作後

- 観測事実: `button "現在の翻訳段階へ進む"` 押下後も URL は `http://127.0.0.1:34115/#translation-management` のままだった。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/first-open.url.txt`
- 観測事実: 初回操作後 snapshot でも `未完了ジョブ一覧` が残り、`ジョブ未選択`、`未完了ジョブ一覧でジョブを選んでください`、`単語翻訳`、`現在の作業` 強調は観測できなかった。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/first-open.snapshot.txt`
- 観測事実: 補助観測として一覧 card の link を押下しても URL は変化しなかった。根拠:
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/first-open-link.url.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/first-open-link.snapshot.txt`
- 観測事実: `agent-browser find role link click --name "ジョブ 1 を選択して現在の翻訳段階へ進む"` は `Element not found` で失敗した。根拠: 作業実行結果

### 再実行後

- 観測事実: 再実行の button 押下後も URL は `http://127.0.0.1:34115/#translation-management` のままだった。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/reopen.url.txt`
- 観測事実: 再実行後 snapshot でも `未完了ジョブ一覧` が残り、`ジョブ未選択`、`未完了ジョブ一覧でジョブを選んでください`、単語翻訳フェーズ UI は観測できなかった。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/reopen.snapshot.txt`

### console

- 観測事実: 初回操作後 console には `Queueing: runtime:ready`、`Connected to backend`、`sending queued message: runtime:ready`、`[vite] connected.` が残っていた。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/after-first-attempt.console.txt`
- 観測事実: 再実行後 console も同じ内容だった。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/reopen.console.txt`
- 観測事実: 初回操作後と再実行後の console には、`GetJobDetail`、`translation-management/job-run`、単語翻訳 phase 読込を示す行は観測できなかった。根拠:
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/after-first-attempt.console.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/reopen.console.txt`

### network

- 観測事実: 初回操作後 network は初期画面読込の script 群と `favicon.ico 404` が中心で、job open 後の追加 request は観測できなかった。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/after-first-attempt.network.txt`
- 観測事実: 再実行後 network も同じ内容だった。根拠: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/reopen.network.txt`

### backend log

- 観測事実: backend log には起動後 `runtime:ready -> Unknown message from front end: runtime:ready` が残っていた。根拠:
  - `tmp/logs/2026-05-09-job-phase-first-open-blank.restart-reproduction.after-first-attempt.backend.log`
  - `tmp/logs/2026-05-09-job-phase-first-open-blank.restart-reproduction.after-reopen.backend.log`
- 観測事実: backend log には、初回操作後と再実行後のどちらでも `GetJobDetail`、job open failure、panic は観測できなかった。根拠:
  - `tmp/logs/2026-05-09-job-phase-first-open-blank.restart-reproduction.after-first-attempt.backend.log`
  - `tmp/logs/2026-05-09-job-phase-first-open-blank.restart-reproduction.after-reopen.backend.log`

## 初回と再実行の差分

- 観測事実: 初回操作後 URL と再実行後 URL は同一だった。根拠:
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/first-open.url.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/reopen.url.txt`
- 観測事実: 初回操作後 DOM と再実行後 DOM は、どちらも `未完了ジョブ一覧` のままだった。根拠:
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/first-open.snapshot.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/reopen.snapshot.txt`
- 観測事実: 初回操作後 console と再実行後 console は同一だった。根拠:
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/after-first-attempt.console.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/reopen.console.txt`
- 観測事実: 初回操作後 network と再実行後 network は同一だった。根拠:
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/after-first-attempt.network.txt`
  - `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/restart-reproduction/reopen.network.txt`
- 観測事実: 初回操作後 backend log と再実行後 backend log では、`runtime:ready` 以外の差分を確認できなかった。根拠:
  - `tmp/logs/2026-05-09-job-phase-first-open-blank.restart-reproduction.after-first-attempt.backend.log`
  - `tmp/logs/2026-05-09-job-phase-first-open-blank.restart-reproduction.after-reopen.backend.log`

## 未確認事項

- 未確認: 人間観測にある「初回操作後に `ジョブの進み方` は表示されるが、下部 panel は `ジョブ未選択` になる」状態。
- 未確認: 一覧へ戻った後の再実行で単語翻訳フェーズ UI が表示される状態。
- 未確認: `agent-browser` の click が今回の画面で app 側操作と同等に扱われているかどうか。
- 未確認: 初回操作時に app 内部で route 遷移が一瞬でも成立したかどうか。

## 推奨 next step

- 推奨: `fix_lane` は次へ進めない。
- 理由: 必須再現条件のうち再起動証跡は満たしたが、人間観測の UI 状態そのものを再現できていない。
- 理由: 初回操作と再実行の差分は、今回の採取分では route、DOM、console、network、backend log のいずれにも現れていない。
