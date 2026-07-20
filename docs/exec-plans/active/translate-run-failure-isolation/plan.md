# Task Plan: translate-run-failure-isolation

`plan.md` は branch 情報と、この task でやること・やらないことの要点を持つ。
設計判断、判断履歴、検証結果、実装結果は持たない。設計は `design.md`、恒久的に残す判断は `docs/changelog.md` に書く。

## やること

- known-issues #7 に対応する。台詞・叙述文の翻訳中に 1 リクエストの応答が失敗した場合（応答エンベロープの decode 失敗、`choices` 無し、非 200 応答、通信失敗）に、その行だけを飛ばさず翻訳 run 全体が中断する挙動を直す。
- 応答 1 件の失敗が残りの未訳行の翻訳をすべて止めないようにする。失敗種別ごとに「run 全体を止める（fail-fast）」か「その行を未訳のまま飛ばして続行する」かの方針を確定してから直す。

## branch 情報

- `execution_branch`: `claude/translate-run-failure-isolation`
- `target_branch`: `master`
- `source_commit`: `ff289309`

## やらないこと

- クラウド AI プロバイダ（Gemini・xAI・Claude）の新規実装は扱わない（known-issues #3）。対象は既存の OpenAI 互換経路の失敗処理だけとする。
- 翻訳結果表示画面の編集・絞り込み機能（known-issues #4）は扱わない。失敗行を人間が手で直す UI は対象外とする。
- 固有名一貫性の事後検証（known-issues #2）は扱わない。
