# Distill: fix

## Focus

- known facts と unknowns を切り分ける
- reproduce / trace に必要な入口だけを返す

## Capture

- 症状
- 再現条件の有無
- 関連 UI、console、Wails log の入口
- 近傍コードと最近の変更点

## Avoid

- 原因断定
- 恒久修正の設計
- logging を先回りで足すこと
