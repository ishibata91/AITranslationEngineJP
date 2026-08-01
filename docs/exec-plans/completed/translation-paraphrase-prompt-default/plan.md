# Plan: translation-paraphrase-prompt-default

## branch情報

- 作業branch: `codex/translation-paraphrase-prompt-default`
- 統合先branch: `master`
- 分岐元commit: `b2a3bda4ba04135e3634293479ce48e1edfa4d98`

## やることの要点

- R-1: SQLの`directive`でkey `口調`が持つ`instruction`を人間が採用した全文へ更新し、人物像の`{traits}`を埋めた台詞では翻訳AIへ送る指示文に反映し、人物像を作れない台詞と台詞以外の翻訳対象には反映しない。
- R-2: 新しいDBとmigration 18まで適用済みの既存DBへmigration 19を適用し、既存DBの未編集の既定値を更新して、利用者が編集した`instruction`、key `口調`以外の`instruction`、`prompt_template.base_directive`は保持する。

## やらないことの要点

- 台詞以外の翻訳指示は変更しない。
- 翻訳処理のプロンプト合成方法、画面、外部APIとの接続方法は変更しない。
- 今回採用した指示文の再評価は行わない。
- migration 17以前で止まっている既存DBの編集内容を保持する移行は扱わない。

## finalization-module 結果

- 正本化判断: `docs/architecture.md` は反映しない。層構成、依存方向、Wails 境界を変更しないためである。人間は2026-08-01に承認した。
- 作業 commit: `776eea7891f86b17b7b090d18509729d221d6931`。
- local merge: `master` へ `git merge --no-ff codex/translation-paraphrase-prompt-default` を実行した。conflict はない。merge commit は `6294ba67f2332f49b29533ef016af047803fd4ec`。
- merge 後検証: `env GOCACHE=/private/tmp/aitranslationenginejp-go-build npm run verify:backend` は通過した。Go試験、architecture検査、境界検査が通過した。
- 実アプリ確認: 未実施。現在のmacOS環境でWailsのlinkが`UTType`を解決できず、アプリ本体を起動できなかった。backendのmigration検証は通過している。
- 残留リスク: 実アプリのdirective編集画面で採用文面が表示されることは未確認である。
