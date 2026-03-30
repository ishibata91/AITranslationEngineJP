# Codex Workflow Overview

この文書は `.codex/workflow_activity_diagram.puml` を文章で補足する鳥瞰図です。
図は流れを示し、このページは lane ごとの目的、分岐条件、close 条件を短く固定します。

## Source Of Truth

- live workflow の正本は `.codex/README.md` と各 `SKILL.md`
- 鳥瞰図の図版正本は `.codex/workflow_activity_diagram.puml`
- このページは diagram を読むための索引であり、diagram と矛盾する独自フローを追加しない

## Overall Shape

`User request` を起点に、workflow は `impl lane` と `fix lane` の 2 本に分岐します。
feature / change は `directing-implementation` から入り、bug / regression は `directing-fixes` から入ります。
どちらの lane も最後に review を 1 回だけ行い、`pass` なら `4humans sync` と commit を済ませて close し、`reroute` なら direction に戻します。

## Impl Lane

標準順序は `directing-implementation -> designing-implementation -> distilling-implementation -> planning-implementation -> architecting-tests -> implementing-frontend or implementing-backend -> sonar-scanner + Sonar CLI open issue gate -> reviewing-implementation -> 4humans sync + commit + close` です。

- `directing-implementation`: 実装要求の入口。active plan を用意し、task-local design が必要なら `designing-implementation` を起動する。
- `designing-implementation`: active plan の `UI` / `Scenario` / `Logic` を task-local design として固める。
- `distilling-implementation`: facts、constraints、gaps、closeout notes を圧縮する。
- `planning-implementation`: 実装順、owned scope、validation を短い brief に落とす。
- `architecting-tests`: 実装前に tests、fixtures、acceptance checks、validation commands を先に固定し、必要な test / fixture を最小範囲で実装する。
- `implementing-frontend` / `implementing-backend`: owned scope に従って実装する。分岐は frontend / backend の責務で決める。
- `sonar-scanner + Sonar CLI open issue gate`: project root で scanner を実行し、`.codex/skills/directing-implementation/scripts/get-open-sonar-issues.ps1` で `status == OPEN` の issue だけを取得して、issue が残る間は implementing skill に戻す。
- `reviewing-implementation`: `仕様逸脱`、`例外処理`、`リソース解放`、`テスト不足` の 4 観点だけを単発で見る。

Sonar issue が解消し、`reviewing-implementation` が `pass` なら `4humans sync`、commit、close に進みます。
`reroute` の場合は `directing-implementation` に戻り、同じ lane の中で plan と実装を更新します。

## Fix Lane

標準順序は `directing-fixes -> distilling-fixes -> tracing-fixes -> (必要時 logging-fixes / analyzing-fixes) -> architecting-tests -> implementing-fixes -> reviewing-fixes -> reporting-risks + 4humans sync + commit + close` です。

- `directing-fixes`: bugfix 要求の入口。事実不足なら distill と trace に進める。
- `distilling-fixes`: 既知事実、再現条件、関連仕様、関連コードを短く整理する。
- `tracing-fixes`: 原因仮説を順位付けし、最小の trace 方針を決める。
- `logging-fixes`: 一時観測ログだけを追加 / 削除する。恒久修正は混ぜない。
- `analyzing-fixes`: 観測結果を事実へ圧縮し、fix 対象か `4humans sync` 対象か、または human-triggered な `updating-docs` 対象かを整理する。
- `architecting-tests`: 再現条件を tests / acceptance checks / validation commands に落とし、必要な回帰 test / fixture を先に実装する。
- `implementing-fixes`: 承認済み scope の恒久修正を行う。
- `reviewing-fixes`: impl lane と同じ 4 観点で単発 review する。
- `reporting-risks`: 必要な時だけ残留リスクを短くまとめる。

diagram 上では `analyzing-fixes` は常に通ります。
`logging-fixes` は temporary logging が必要な時だけ挿入され、不要なら `tracing-fixes` から直接 `analyzing-fixes` に進みます。

## Reroute And Close

- review は score 制の loop にしない
- `pass` なら同じ変更の中で `4humans sync` を完了させ、commit してから close する
- `reroute` なら direction skill に戻し、plan、tests、実装を必要最小限で更新する
- `4humans sync` と commit は別 lane に押し出さず、その lane の close 条件として扱う
- `docs/` 正本更新は通常 lane の close 条件に含めず、human が `updating-docs` を直接起動した時だけ扱う

## Records And Evidence

- 非自明な変更は `docs/exec-plans/active/` に plan を置く
- 完了後は `docs/exec-plans/completed/` へ移す
- 詳細な挙動や制約は docs へ肥大化させず、tests、acceptance checks、validation commands に寄せる
- `directing-* -> downstream skill` の handoff contract 例は、各 directing skill 配下の `references/*.json` を見る
- `downstream skill -> directing-*` の返却 contract 例は、各 downstream skill 配下の `references/*.json` を見る
- harness は `powershell -File scripts/harness/run.ps1 -Suite structure|design|execution|all` を入口にする

