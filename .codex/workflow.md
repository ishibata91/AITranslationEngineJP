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
どちらの lane も最後に review を 1 回だけ行い、`pass` なら close、`reroute` なら direction に戻します。

## Impl Lane

標準順序は `directing-implementation -> distilling-implementation -> planning-implementation -> architecting-tests -> implementing-frontend or implementing-backend -> reviewing-implementation -> docs sync + close` です。

- `directing-implementation`: 実装要求の入口。必要なら active plan の `UI` / `Scenario` / `Logic` を埋める。
- `distilling-implementation`: facts、constraints、gaps、docs sync 候補を圧縮する。
- `planning-implementation`: 実装順、owned scope、validation を短い brief に落とす。
- `architecting-tests`: 実装前に tests、fixtures、acceptance checks、validation commands を先に固定する。
- `implementing-frontend` / `implementing-backend`: owned scope に従って実装する。分岐は frontend / backend の責務で決める。
- `reviewing-implementation`: `仕様逸脱`、`例外処理`、`リソース解放`、`テスト不足` の 4 観点だけを単発で見る。

`reviewing-implementation` が `pass` なら docs sync と close に進みます。
`reroute` の場合は `directing-implementation` に戻り、同じ lane の中で plan と実装を更新します。

## Fix Lane

標準順序は `directing-fixes -> distilling-fixes -> tracing-fixes -> (必要時 logging-fixes / analyzing-fixes) -> architecting-tests -> implementing-fixes -> reviewing-fixes -> reporting-risks + docs sync + close` です。

- `directing-fixes`: bugfix 要求の入口。事実不足なら distill と trace に進める。
- `distilling-fixes`: 既知事実、再現条件、関連仕様、関連コードを短く整理する。
- `tracing-fixes`: 原因仮説を順位付けし、最小の trace 方針を決める。
- `logging-fixes`: 一時観測ログだけを追加 / 削除する。恒久修正は混ぜない。
- `analyzing-fixes`: 観測結果を事実へ圧縮し、fix 対象か docs sync 対象かを整理する。
- `architecting-tests`: 再現条件を tests / acceptance checks / validation commands に落とし、回帰確認を先に固める。
- `implementing-fixes`: 承認済み scope の恒久修正を行う。
- `reviewing-fixes`: impl lane と同じ 4 観点で単発 review する。
- `reporting-risks`: 必要な時だけ残留リスクを短くまとめる。

diagram 上では `analyzing-fixes` は常に通ります。
`logging-fixes` は temporary logging が必要な時だけ挿入され、不要なら `tracing-fixes` から直接 `analyzing-fixes` に進みます。

## Reroute And Close

- review は score 制の loop にしない
- `pass` なら同じ変更の中で docs sync を完了させて close する
- `reroute` なら direction skill に戻し、plan、tests、実装を必要最小限で更新する
- docs sync は別 lane に押し出さず、その lane の close 条件として扱う

## Records And Evidence

- 非自明な変更は `docs/exec-plans/active/` に plan を置く
- 完了後は `docs/exec-plans/completed/` へ移す
- 詳細な挙動や制約は docs へ肥大化させず、tests、acceptance checks、validation commands に寄せる
- harness は `powershell -File scripts/harness/run.ps1 -Suite structure|design|execution|all` を入口にする

