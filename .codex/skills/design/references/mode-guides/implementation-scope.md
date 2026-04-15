# Design: implementation-scope

## Goal

- human review 後の AI handoff を独立 artifact に固定する
- 実装者向け仕様書である `implementation-brief` を再説明せず、AI が再解釈なしで動ける入力だけを残す

## Rules

- path は `docs/exec-plans/active/<task-id>.implementation-scope.md` に固定する
- 本文は英語で短く書く
- `requirements`、存在する `ui-mock` / `scenario`、`implementation-brief`、`design-review` 結果、HITL を前提にする
- handoff ごとに `implementation_target`、`owned_scope`、`depends_on`、`validation_commands` を明示する
- 実装 task があるものはすべて `owned_scope` を先に固定する
- work plan には path と接続情報だけを残し、本体を埋め込まない
- orchestrate が再解釈しなくてよい粒度で分割する

## Include

- short summary
- handoff order
- owned scope per handoff
- dependency edges
- validation commands
- completion signal
- notes only when the handoff would be ambiguous without them

## Avoid

- 日本語 prose
- 背景や判断理由の長い再説明
- human review 前の確定
- 実装 scope と将来 scope の混在
