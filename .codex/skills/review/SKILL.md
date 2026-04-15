---
name: review
description: design review、UI check、implementation review を mode 分岐で扱い、単発 review 結果を返す role skill。
---

# Review

## Goal

- 役割ごとに review する対象を限定する
- design、UI、implementation の判定軸を混ぜない
- finding を返し、設計差分か実装差分かを切り分ける

## Modes

- `design-review`: task-local requirements、存在する UI / scenario artifact、review 用差分図を照合する
- `ui-check`: Playwright MCP で主要導線、画面状態、console、必要証跡を確認する
- `implementation-review`: 実装差分が `implementation-scope` か accepted fix scope と整合するか確認する

## Common Rules

- review は mode ごとに独立して行う
- `design-review` は design bundle 全体を見て、human review 1 回で判断できる材料がそろっているかを確認する
- `design-review` は要件取りこぼし、存在する artifact の設計不整合、検証不足、構造差分の不整合だけを見る
- architecture 対象がある時は `docs/architecture.md` と対象 D2 の整合を確認する
- `ui-check` は `http://host.docker.internal:34115` を使い、主要導線と画面状態の証跡だけを返す
- `implementation-review` は `implementation-scope` を AI handoff 正本として扱う
- `implementation-review` は `implementation-scope` や accepted fix scope にない仕様や好みで判定しない
- backend を含む `implementation-review` は Sonar MCP で open `HIGH` / `BLOCKER`、open reliability、open security を確認する
- backend を含む `implementation-review` は `sonar_gate_result` に件数、対象 project、supporting issue を含める
- backend を含む `implementation-review` は いずれかが非 0 件なら `reroute` とする
- UI 逸脱や導線不整合は `reroute` とし、恒久修正は `implement` へ戻す
- 新しい改善提案や新しい要件解釈は追加しない
- 役割を再確定せず、呼び出し元で確定した `review_mode` を前提に進める

## Output

- `decision`: `pass` | `reroute`
- `findings`
- `recheck`
- `sonar_gate_result`
- `closeout_notes`
- `human_open_questions`
- `ui_evidence`

## Detailed Guides

- `references/mode-guides/design-review.md`
- `references/mode-guides/ui-check.md`
- `references/mode-guides/implementation-review.md`

## Reference Use

- quick overview は `../orchestrate/references/orchestrate.to.review.json` を使う
- mode 別 contract は `../orchestrate/references/contracts/orchestrate.to.review.<mode>.json` を正本とする
- 返却 contract は `references/contracts/review.to.orchestrate.<mode>.json` を正本とする
