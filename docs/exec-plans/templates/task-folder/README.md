# Task Folder Template

新しい exec-plan は task ごとの folder として作る。
`plan.md` は人間と合意した要求を持つ。設計は `design.md`、仕様は `spec.md`、未決定事項は `pending.md`、参照は `references.md` に分ける。

## 作成先

- active: `docs/exec-plans/active/<task-id>/`
- completed: `finalization-module` が local merge 後に移動する `docs/exec-plans/completed/<task-id>/`

## 標準ファイル

- `plan.md`: 人間と合意した要求だけを持つ。設計判断、未決定事項、参照、判断履歴、検証結果を持たない。`design.md` と `spec.md` は要求ごとに節を分ける
- `design.md`: 要求ごとの as-is と to-be の方針を持つ。仕様、未決定事項、source の path・symbol・外部資料を持たない。両フロー共通の 1 テンプレート
- `spec.md`: この task の確定仕様（要求ごとの仕様）だけを持つ。各仕様は確かめ方を併記する。仕様の文はそのまま実テストの test case 名にする。両フロー共通の 1 テンプレート
- `pending.md`: task 固有の未決定事項とブロッカーを持つ。解決した項目は結論を正本へ反映してから消す
- `references.md`: source の path・symbol と外部資料の所在を `REF-<番号>` で持つ。解釈と結論を持たない
- `log.jsonl`: 状態遷移だけを持つ監査用の追記履歴。field は `at`、`actor`、`type`、`documents` と、pending event の `pending_id` だけにする。人間が task の履歴確認を明示的に依頼した場合だけ読む。通常の設計、仕様、実装、レビュー、コンテキスト再取得、状況確認では読まない
- `implementation.md`: 実装差分、仕様との対応、検証結果、未確認事項、人間が直接記入する指摘を持つ
- `investigation.md`: 修正フロー（`fix-workflow`）だけが作る。再現確認と原因究明（観測済み問題、画面再現確認、原因仮説、観測ログ検証、確定原因）。設計フローでは作らない
- `storybook-review-loop.md`: 画面表示の変更がある task で、Storybook レビューループが確定した story、変更後の画面仕様、反映先、現在分類、承認状態を持つ
- `criteria.md`: 実験フロー（`experiment-workflow`）だけが作る。何をどこまで達成したら終わりかを達成条件の表で持ち、達成条件に置かない値を診断の表、`researcher` が集めた手段を選択肢の表、開発用と評価用の分け方を標本の表で持つ。標本を読むときの誤りの種類も持つ。達成の線と数え方は Claude 本体が決めて `達成条件レビュー` が検証し、人間が承認するのは目的への貢献の欄と、資源を前提から外す選択肢の行に限る。`design.md` と `spec.md` は作らない
- `loop-log.md`: 実験フローだけが作る。準備の段の `予備の回`（達成の線を置くために対象を 1 つ変えて測った動き幅）と、各回の段（探索 / 確証）、変えた対象、測った標本、判定、commit hash と、仮説の表、測定手段を変えた記録、採否を持つ。数値の並びは持たない
- `measurements.csv`: 実験フローだけが作る。回ごとの測定値を持つ。数値と判定だけで、標本そのものを持たない（標本は著作物やライセンスの制約を含む場合に commit できないため、砂場の `sample-review.jsonl` が持つ）

## 読み方

- 最初に `plan.md` で要求を読む
- `pending.md` が空でない場合は、design-review、設計HITL、実装へ進まない
- source の path、symbol、外部資料の所在は `references.md` を読む
- 仕様を短時間で確認する時は `spec.md` を読む。要求ごとの仕様が正。`design.md` と食い違う場合は `spec.md` を優先する
- 設計の as-is と to-be の方針は `design.md` を読む
- 修正フローの再現確認・原因究明は `investigation.md` を読む。確定原因が正
- 実験フローの達成の判定は `criteria.md` の達成条件の表を読む。目的への貢献の欄と線の出どころの欄が空欄の値は達成条件に置かない
- 実験フローで達成の線をなぜその値にしたかは `criteria.md` の線の出どころの欄を読む。目指す位置と、`予備の回` で実測した動き幅と、隔たりが動き幅の何回分かが書かれている。線が揺れに埋もれて診断へ移った値は、診断の表の移した理由の欄を読む
- 実験フローの経過は `loop-log.md` の各回と仮説の表を読む。測った値は `measurements.csv` が正
- 実験の判断基準（段の分離、標本の分離、達成の線の置き方、揺れの底の測定、要因の振り方、値の差を改善と呼ぶ条件）は `research-protocol` skill を読む。その根拠と出典は `docs/references/research-methods.md` を読む
- 画面表示の確認時は Storybook の story と svelte コンポーネントを読む
- Storybook の作成、起動、分類、確認資源、`fixture` 種類基準は `docs/references/storybook.md` を読む
- 語の正本は `docs/vocabulary.md` を読む。`spec.md` と `design.md` は、要求の文に現れる語かこの正本にある語だけを使う
- 恒久的に残す判断・変更履歴は `docs/changelog.md` に書く。正本（`docs/architecture.md`）には現在状態だけを書く
- completed 移動は `finalization-module` だけが扱う
