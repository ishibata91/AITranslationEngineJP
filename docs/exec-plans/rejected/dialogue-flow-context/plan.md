# Plan: dialogue-flow-context

`task_id`: `dialogue-flow-context`

`状態`: 棄却（2026-07-20。上流 `dialogue-graph` の棄却に連鎖）

## 棄却の要旨

会話往復の文脈を翻訳へ供給する本 task は、基盤となる上流 `dialogue-graph` が実測で棄却されたため、同時に棄却する。棄却根拠（PNAM/TCLT の実データ形と費用対効果）は `dialogue-graph/plan.md` に集約する。

## branch 情報

- 作業 branch: `claude/dialogue-flow-context`
- 統合先 branch: `master`
- 分岐元 commit: `eee7bb82`

## 経緯（旧プランの狙い）

会話の往復（プレイヤー選択肢 → NPC 応答 → 次）の文脈を翻訳へ渡し前後で噛み合った訳にする狙いだった。`dialogue-graph` が保存する会話グラフを辿り root path を翻訳リクエストへ添える計画だったが、上流棄却により着手しない。
