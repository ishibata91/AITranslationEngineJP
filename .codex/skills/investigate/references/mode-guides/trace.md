# Investigate: trace

## Goal

- `hypotheses` と `observation_points` だけを返す

## Procedure

- `observed_facts` を時系列に並べる
- 仮説は 1 から 3 個に絞る
- 各仮説ごとに必要な `observation_points` を定義する
- 恒久修正案は返さない

## Return

- `hypotheses`
- `observation_points`
- `recommended_next_step`
