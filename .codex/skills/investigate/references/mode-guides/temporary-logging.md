# Investigate: temporary-logging

## Goal

- trace に必要な一時ログだけを add / remove する

## Rules

- 恒久修正や test 追加を混ぜない
- ログ tag は `[tracing-fixes]` を使う
- 観測完了後に remove できる形で最小差分にする
- log point は `observation_points` に対応させる

## Return

- `observed_facts`
- `remaining_gaps`
- `recommended_next_step`
