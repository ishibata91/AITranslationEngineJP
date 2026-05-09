# 修正前調査: 手動再現証跡

## 判断結果

- 判定: 再現済み
- 調査 mode: `修正前調査`, `UI 根拠`
- 引き継ぎ先: `fix_lane`
- 判断: 明示的に開発環境を起動し直した後、初回操作で人間観測と同じ下部 panel の `ジョブ未選択` 状態を再現した。
- 判断: 一覧へ戻って同じ job の操作を再実行すると、単語翻訳フェーズ UI が表示されることを確認した。
- 判断: 既存の investigator 調査で再現できなかった主因は、job card と操作 button が viewport 外にあり、`agent-browser click` が実アプリ操作として届いていなかった可能性が高い。

## 再起動条件

- 開発環境は `sh ./scripts/dev/run-wails-agent-browser.sh` で起動した。
- 初回 open は `agent-browser open http://127.0.0.1:34115/#translation-management` で実行した。
- 初回 open 時の URL は `http://127.0.0.1:34115/#translation-management` だった。

## 操作条件

- 初回表示では、snapshot 上に `link "ジョブ 1 を選択して現在の翻訳段階へ進む"` と `button "現在の翻訳段階へ進む"` が存在した。
- DOM 診断では、job card と `現在の翻訳段階へ進む` button は viewport 外の y=1083 付近にあった。
- `agent-browser press End` で job card を viewport 内へ入れた。
- viewport 内に入れた後、`agent-browser click @e18` で `現在の翻訳段階へ進む` を押した。

## 初回操作後の観測事実

- URL は `http://127.0.0.1:34115/#translation-management/job-run` になった。
- 画面には `ジョブの進み方` が表示された。
- 下部 panel は `未完了ジョブ一覧でジョブを選んでください` を表示した。
- snapshot では `ジョブ #1` 見出しと単語翻訳フェーズ UI は表示されなかった。
- screenshot: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/manual-restart-reproduction/after-scroll-click-button.png`

## 再実行後の観測事実

- `一覧へ戻る` を押して一覧へ戻った。
- 再度 `agent-browser press End` で job card を viewport 内へ入れた。
- 再度 `agent-browser click @e18` で `現在の翻訳段階へ進む` を押した。
- URL は `http://127.0.0.1:34115/#translation-management/job-run` になった。
- `ジョブ #1` 見出しが表示された。
- `単語翻訳` 見出しが表示された。
- `操作`、`進行状況`、`実行設定`、`結果 summary`、`失敗情報`、`単語翻訳の次の作業` の region が表示された。
- screenshot: `tmp/agent-browser/2026-05-09-job-phase-first-open-blank/manual-restart-reproduction/after-second-scroll-click-button.png`

## 初回と再実行の差分

- 初回操作後は route が `job-run` へ進んだが、job target が下部 panel へ渡らず `未完了ジョブ一覧でジョブを選んでください` が表示された。
- 再実行後は同じ route で `ジョブ #1` と `単語翻訳` UI が表示された。
- 差分は、route 遷移ではなく job target の受け渡しにある可能性が高い。

## 影響ファイル候補

- `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte`
  - 理由: 一覧 card の `jobRunTarget` で open し、その後 `selectJob` を実行する。
- `frontend/src/ui/views/AppShell.svelte`
  - 理由: `selectedJobRunTarget` と `selectedTranslationManagementViewId` を持ち、job-run 画面へ target を渡す。
- `frontend/src/application/usecase/translation-job-management/translation-job-management.usecase.ts`
  - 理由: `selectJob` の loading 期間に `selectedJobDetail` が未設定になる。
- `frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.ts`
  - 理由: 一覧 card target と画面全体 target の生成元が分かれている。
- `frontend/src/ui/screens/job-run/JobRunPage.svelte`
  - 理由: `selectedJobTarget` がない場合に `ジョブ未選択` 分岐を表示する。

## 残り不足

- 未確認: `selectedJobRunTarget` が初回操作後に null へ戻る時系列。
- 未確認: `TranslationJobManagementPage` の subscription が `viewModel.jobRunTarget = null` を親へ返す時点。
- 未確認: `GetJobDetail` 成功後に再実行で `selectedJobDetail` が保持され、`viewModel.jobRunTarget` が non-null になる時系列。

## 推奨 next step

- 推奨: `fix_lane` は原因箇所シーケンス図へ進める前に、今回の再現事実を根拠に `diagrammer` 起動入力を作る。
- 推奨: 図では `TranslationJobManagementPage.handleOpenJob`、`selectJob` loading 更新、presenter の `viewModel.jobRunTarget`、`AppShell.selectedJobRunTarget`、`JobRunPage` の `ジョブ未選択` 分岐を扱う。
