# Codex Workflow Overview

この文書は `.codex/workflow_activity_diagram.puml` を文章で補足する鳥瞰図です。
図は流れを示し、このページは lane ごとの目的、分岐条件、close 条件を短く固定します。

## Source Of Truth

- live workflow の正本は `.codex/README.md` と各 `SKILL.md`
- 鳥瞰図の図版正本は `.codex/workflow_activity_diagram.puml`
- このページは diagram を読むための索引であり、diagram と矛盾する独自フローを追加しない

## Naming Rule

- workflow 記述では、論理名と実名をできるだけ同じ行に置く
- 初出または重要な参照は `論理名 (`actual-name`)` を優先する
- 人間 review で意味が先に分かり、actual name でも検索できる記述を優先する

## Overall Shape

`User request` を起点に、workflow は `implementation proposal lane`、`implementation execution lane`、`fix lane` の 3 役割に分かれます。
feature / change は implementation proposal lane owner (`proposing-implementation`) から入り、human LGTM の後に implementation execution lane owner (`directing-implementation`) へ渡り、bug / regression は fix lane owner (`directing-fixes`) から入ります。
どちらの lane も最後に review を 1 回だけ行い、`pass` なら `4humans sync` と commit を済ませて close し、コードベース境界や実行フローが変わる時は diagramming D2 skill (`diagramming-d2`) で `4humans/class-diagrams/` と `4humans/sequence-diagrams/` の `.d2` / `.svg` も同一変更で更新します。new detail diagram を追加する時は `4humans/diagrams/overview-manifest.json` と、manifest で紐づく overview `.d2` / `.svg` も同一変更で更新します。`reroute` なら direction に戻します。

## Impl Lane

標準順序は `implementation proposal lane owner (`proposing-implementation`) -> implementation distill skill (`distilling-implementation`) -> task-local design skill (`designing-implementation`) -> diagrammer (`diagrammer`) + diagramming D2 skill (`diagramming-d2`) -> human LGTM -> implementation execution lane owner (`directing-implementation`) -> implementation workplan skill (`planning-implementation`) -> test architecture skill (`architecting-tests`) -> frontend implementer (`implementing-frontend`) or backend implementer (`implementing-backend`) + assigned lint suite -> sonar-scanner + Sonar MCP open issue gate -> implementation review skill (`reviewing-implementation`) -> full harness -> 4humans sync + commit + close` です。

- implementation proposal lane owner (`proposing-implementation`): 実装要求の入口。日本語の active plan を用意し、task-local design と human review に必要な最小限の情報を整える。
- task-local design skill (`designing-implementation`): active plan の `UI` / `Scenario` / `Logic` を task-local design として固める。
- implementation distill skill (`distilling-implementation`): facts、constraints、gaps、closeout notes を圧縮する。
- diagrammer (`diagrammer`): review 用差分図を担当する。diagramming D2 skill (`diagramming-d2`) を使って、追加は緑、削除は赤で読める D2 / SVG を作る。
- human LGTM: active plan の `承認記録` と `HITL 状態` に記録し、承認前は execution lane を起動しない。
- implementation execution lane owner (`directing-implementation`): 承認済み active plan を受け取り、planning 以降の execution、gate、close を管理する。
- implementation workplan skill (`planning-implementation`): 実装順、owned scope、validation を短い brief に落とす。
- test architecture skill (`architecting-tests`): 実装前に tests、fixtures、acceptance checks、validation commands を先に固定し、必要な test / fixture を最小範囲で実装する。
- frontend implementer (`implementing-frontend`) / backend implementer (`implementing-backend`): owned scope に従って実装し、frontend は `python3 scripts/harness/run.py --suite frontend-lint`、backend は `python3 scripts/harness/run.py --suite backend-lint` だけを local validation として実行する。分岐は frontend / backend の責務で決める。
- `sonar-scanner + Sonar MCP open issue gate`: implementation execution lane owner (`directing-implementation`) が project root で scanner を実行し、その後に Sonar MCP の `search_sonar_issues_in_projects` を直接使って `project == ishibata91_AITranslationEngineJP` かつ `status == OPEN` の issue だけを取得して、issue が残る間は implementing skill に戻す。
- implementation review skill (`reviewing-implementation`): `仕様逸脱`、`例外処理`、`リソース解放`、`テスト不足`、`4humans` D2 sync 要否と実施有無 の 5 観点だけを単発で見る。
- `full harness`: implementation review skill (`reviewing-implementation`) が `pass` を返した後に、implementation execution lane owner (`directing-implementation`) が `python3 scripts/harness/run.py --suite all` を実行する。

Sonar issue が解消し、implementation review skill (`reviewing-implementation`) が `pass` で、さらに full harness が通った時だけ `4humans sync`、必要な `4humans/class-diagrams/` と `4humans/sequence-diagrams/` の `.d2` / `.svg` 更新、new detail diagram 追加時の `4humans/diagrams/overview-manifest.json` と対応 overview 更新、review 用差分図削除、commit、close に進みます。
`reroute` の場合は implementation execution lane owner (`directing-implementation`) に戻り、proposal のやり直しが必要な時だけ implementation proposal lane owner (`proposing-implementation`) に戻します。

## Fix Lane

標準順序は `fix lane owner (`directing-fixes`) -> fix distill skill (`distilling-fixes`) -> fault trace skill (`tracing-fixes`) -> (必要時 logging skill (`logging-fixes`) / fix analysis skill (`analyzing-fixes`)) -> test architecture skill (`architecting-tests`) -> fix implementer (`implementing-fixes`) -> fix review skill (`reviewing-fixes`) -> risk reporting skill (`reporting-risks`) + 4humans sync + commit + close` です。

- fix lane owner (`directing-fixes`): bugfix 要求の入口。事実不足なら fix distill skill (`distilling-fixes`) と fault trace skill (`tracing-fixes`) に進める。
- fix distill skill (`distilling-fixes`): 既知事実、再現条件、関連仕様、関連コードを短く整理する。
- fault trace skill (`tracing-fixes`): 原因仮説を順位付けし、最小の trace 方針を決める。
- logging skill (`logging-fixes`): 一時観測ログだけを追加 / 削除する。恒久修正は混ぜない。
- fix analysis skill (`analyzing-fixes`): 観測結果を事実へ圧縮し、fix 対象か `4humans sync` 対象か、または human-triggered な docs sync skill (`updating-docs`) 対象かを整理する。
- test architecture skill (`architecting-tests`): 再現条件を tests / acceptance checks / validation commands に落とし、必要な回帰 test / fixture を先に実装する。
- fix implementer (`implementing-fixes`): 承認済み scope の恒久修正を行う。
- fix review skill (`reviewing-fixes`): impl lane と同じ 5 観点で単発 review する。
- risk reporting skill (`reporting-risks`): 必要な時だけ残留リスクを短くまとめる。

diagram 上では fix analysis skill (`analyzing-fixes`) は常に通ります。
logging skill (`logging-fixes`) は temporary logging が必要な時だけ挿入され、不要なら fault trace skill (`tracing-fixes`) から直接 fix analysis skill (`analyzing-fixes`) に進みます。

## Reroute And Close

- review は score 制の loop にしない
- `pass` なら同じ変更の中で `4humans sync` を完了させ、必要な `4humans/class-diagrams/` と `4humans/sequence-diagrams/` の `.d2` / `.svg` 更新も済ませてから commit して close する
- new detail diagram 追加時は `4humans/diagrams/overview-manifest.json` を更新し、manifest で紐づく overview `.d2` / `.svg` も同じ変更で更新する
- `reroute` なら direction skill に戻し、plan、tests、実装を必要最小限で更新する
- `4humans sync` と commit は別 lane に押し出さず、その lane の close 条件として扱う
- `docs/` 正本更新は通常 lane の close 条件に含めず、human が `updating-docs` を直接起動した時だけ扱う

## Records And Evidence

- 非自明な変更は `docs/exec-plans/active/` に plan を置く
- 完了後は `docs/exec-plans/completed/` へ移す
- 詳細な挙動や制約は docs へ肥大化させず、tests、acceptance checks、validation commands に寄せる
- `directing-* -> downstream skill` の handoff contract 例は、各 directing skill 配下の `references/*.json` を見る
- `downstream skill -> directing-*` の返却 contract 例は、各 downstream skill 配下の `references/*.json` を見る
- harness は `python3 scripts/harness/run.py --suite structure|design|execution|all` を入口にする
