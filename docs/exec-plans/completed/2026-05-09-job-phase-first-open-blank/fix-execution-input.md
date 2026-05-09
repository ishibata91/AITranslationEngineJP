# 修正実行入力

## 対象成果物

- `実装証跡`

## 実装 agent

- agent: `frontend_implementer`
- skill: `implement-frontend`

## 人間観測

- 初回に `jobID1` で `現在の翻訳段階へ進む` を押す。
- route は `#translation-management/job-run` へ進む。
- 進行カードでは単語翻訳が現在作業として表示される。
- 下部 panel は `ジョブ未選択` になる。
- 一覧へ戻って同じ操作を再実行すると、`ジョブ #1` と単語翻訳フェーズ UI が表示される。

## 修正前調査

- 根拠: `./pre-fix-investigation.manual-reproduction.md`
- 根拠: `./cause-sequence.md`
- 初回操作後は route が `job-run` へ進むが、job target が `JobRunPage` へ渡らない。
- 再実行後は同じ route で `ジョブ #1` と `単語翻訳` UI が表示される。
- 差分は route ではなく、`JobRunPage` に渡る selected job target の保持にある。

## 原因箇所シーケンス図

- source: `./cause-sequence.puml`
- 説明: `./cause-sequence.md`
- svg: `./cause-sequence.svg`

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

## 実装対象

- frontend プロダクトコードだけを変更する。
- 修正対象は、detail loading 中の `viewModel.jobRunTarget` 生成に限定する。
- `detailPhase` が `loading` で `selectedJobId` が存在する場合だけ、一覧 summary から `jobRunTarget` を生成する。
- `selectedJobDetail` が存在する場合は detail を優先し、stale 状態では summary から target を復元しない。

## 禁止変更範囲

- `internal/` を変更しない。
- プロダクトテストを変更しない。
- docs 正本本文を変更しない。
- `.codex/` を変更しない。
- UI 表示、画面文言、layout、style を変更しない。
- Wails bridge 境界を迂回しない。

## 回帰確認観点

- 再起動後の初回操作で、`現在の翻訳段階へ進む` から `#translation-management/job-run` へ進む。
- 初回操作後に `ジョブ #1` と `単語翻訳` UI が表示される。
- 初回操作後に `未完了ジョブ一覧でジョブを選んでください` が表示されない。
- 一覧へ戻って同じ操作を再実行しても、`ジョブ #1` と `単語翻訳` UI が表示される。

## 検証コマンド

- `python3 scripts/harness/run.py --suite frontend-local`

## 期待する実装証跡

- 変更ファイル一覧
- 実装内容
- 検証結果
- UI 証跡または未取得理由
- 残留リスク
