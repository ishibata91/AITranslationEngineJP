# 人間観測記録: fix-lucien-target-list-empty

## 確認された不具合

- 操作: `dictionaries/Lucien.esp_Export.json` をデータロード画面で読み込む。
- 結果1: 単語翻訳画面へ遷移できる。
- 結果2: 進捗パネルは「4900件ある」と表示する。
- 結果3: 処理対象の一覧パネルは「0件」と表示する。
- 期待: 処理対象一覧パネルに、読み込んだ翻訳対象（約4900件相当）が表示される。

## 期待との差分

- 進捗パネルが認識する件数（4900件）と、処理対象一覧パネルが表示する件数（0件）が一致しない。
- 利用者は読み込み済みの対象を一覧で操作できない。

## 観測された操作または条件

- 入力ファイル: `dictionaries/Lucien.esp_Export.json`。
- 画面遷移: データロード画面 → 単語翻訳画面。

## 関連資料

- `docs/exec-plans/active/fix-lucien-target-list-empty/attempts/round-1/notes.md`（codex 着手前の素データ）。
