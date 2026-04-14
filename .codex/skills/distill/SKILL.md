---
name: distill
description: task の入口情報を圧縮し、facts / constraints / gaps / required reading を次工程へ渡す role skill。
---

# Distill

## Goal

- active plan と入口情報から次工程に必要な repo 文脈だけを抽出する
- facts、constraints、gaps、required_reading を downstream が使える粒度で返す
- implement / fix / investigate / refactor の入口差分を mode で吸収する

## What Stays Here

- facts と推測の分離
- required reading の最小化
- 関連 code pointer の明示
- reproduction context の有無判断
- task を再設計せず、次工程が動ける最小 packet を返すこと

## Outputs

- `facts`
- `constraints`
- `gaps`
- `required_reading`
- `reproduction_status`
- `related_code_pointers`
- `recommended_next_skill`

## Rules

- 実装コード、test、diagram source を変更しない
- `design`、`implement`、`tests`、`review` の成果物を先回りで作らない
- 入口で渡された path だけで不足する時に限って最小限の追加探索を行う
- fix / investigate では事実と推測を分ける
- `packet file` を作らない
- `changes/` や `context_board` を前提にしない
- 役割を再確定せず、呼び出し元で確定した `task_mode` を前提に整理する

## Mode Notes

- `implement`: requirements / affected surface / validation entry を拾う
- `fix`: 再現条件、既知症状、関連ログ経路、近傍コードを拾う
- `refactor`: 境界、依存方向、非目的の振る舞い変更を拾う
- `investigate`: 仮説形成に必要な観測点だけを返す

## Detailed Guides

- `references/mode-guides/implement.md`
- `references/mode-guides/fix.md`
- `references/mode-guides/refactor.md`
- `references/mode-guides/investigate.md`

## Reference Use

- quick overview は `../orchestrate/references/orchestrate.to.distill.json` を使う
- mode 別 contract は `../orchestrate/references/contracts/orchestrate.to.distill.<mode>.json` を正本とする
- 返却 contract は `references/contracts/distill.to.orchestrate.<mode>.json` を正本とする
