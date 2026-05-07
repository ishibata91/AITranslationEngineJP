# Frontend Human Review Input

## 状態

- `status`: approved
- `source_input`: `./frontend-implementation-input.md`
- `source_result`: `./frontend-implementation-result.md`
- `next_on_approved`: `backend 実装`
- `next_on_rejected`: `frontend 実装`
- `approval_record`: 人間が「フロントレビュー終わり」と回答したため、frontend 実装後人間レビューを承認済みとする。
- `review_fix_record`: 人間レビュー中の指摘により、Job Setup の不足理由表示を修正済み。

## 確認してほしい結果

- マスターペルソナ画面で、provider、model、model list、保存状態がモデル設定カードの共有制御で扱われる。
- Job Setup の 3 翻訳段階で、同じモデル設定カード制御が使われる。
- provider 変更後に、旧 provider の model list と model が現在 provider へ混入しない。
- APIキー未設定時に、共有カード内には AIサービス設定を開く導線を出さず、更新不可状態だけが表示される。
- `AIModelSelectionCard.svelte` は表示部品として維持され、fake mode 固有分岐が UI に見えない。

## 実装者の確認結果

- `npm --prefix frontend run check`: 通過。
- `python3 scripts/harness/run.py --suite frontend-local`: 通過。
- 代替対象 frontend test: 10 files / 121 tests 通過。
- agent-browser でマスターペルソナと Job Setup を確認済み。
- fakeAPI 付きでマスターペルソナと Job Setup を確認済み。

## UI 証跡

- 起動 command: `npm run dev:wails:agent-browser`
- マスターペルソナ URL: `http://localhost:34115/#master-persona`
- マスターペルソナ screenshot: `/Users/iorishibata/.agent-browser/tmp/screenshots/screenshot-1778112299317.png`
- Job Setup 経路: `http://localhost:34115/#translation-management` から `セットアップ` tab
- Job Setup screenshot: `/Users/iorishibata/.agent-browser/tmp/screenshots/screenshot-1778112335271.png`
- fakeAPI マスターペルソナ URL: `http://localhost:34115/?fakeApi=1&fakeScenario=success#master-persona`
- fakeAPI マスターペルソナ screenshot: `/Users/iorishibata/.agent-browser/tmp/screenshots/screenshot-1778128013316.png`
- fakeAPI Job Setup URL: `http://localhost:34115/?fakeApi=1&fakeScenario=success#translation-management` から `セットアップ` tab
- fakeAPI Job Setup screenshot: `/Users/iorishibata/.agent-browser/tmp/screenshots/screenshot-1778127990067.png`
- fakeAPI error マスターペルソナ URL: `http://localhost:34115/?fakeApi=1&fakeScenario=error#master-persona`
- fakeAPI error Job Setup URL: `http://localhost:34115/?fakeApi=1&fakeScenario=error#translation-management` から `セットアップ` tab
- fakeAPI error 確認結果: マスターペルソナは「モデル一覧を取得できませんでした。」を表示し、Job Setup は「Job Setup の確認に失敗しました。」を表示した。
- fakeAPI config-missing 確認結果: Job Setup の不足理由は `ほか 2 件` となり、「3 つの翻訳段階が揃うと作成前確認を実行します。」は不足理由から消えた。
- console errors: 詳細なしの空結果。

## 残留リスク

- backend / Wails gateway は未変更である。
- 実 provider の model list 保存取得は後続 wave の接続結果に依存する。
- success / error 以外の状態 variant の幅別 UI 確認は未実行である。
