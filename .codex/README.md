# .codex

このディレクトリは、AITranslationEngineJp の live なマルチエージェント作業フローの正本です。
プロダクト仕様と設計は `docs/` を正本とし、lane、skill、agent の役割と handoff は `.codex/` を正本とします。
この repo の workflow は proposal lane (`proposing-implementation`)、execution / fix lane (`directing-implementation` / `directing-fixes`) で動かし、過去 repo 固有の packet や review loop は live 契約に戻しません。

## Naming Rule

- workflow 文書では、論理名と実名を分離しない
- 初出または重要な参照は `論理名 (`actual-name`)` を優先する
- 人間 review で意味が先に読めて、actual skill / agent name でも検索できる記述を優先する
- 例: implementation proposal lane owner (`proposing-implementation`)、implementation execution lane owner (`directing-implementation`)、fix lane owner (`directing-fixes`)、task-local design skill (`designing-implementation`)

## 入口

- 実装 proposal の入口: `skills/proposing-implementation/SKILL.md`
- 実装 execution の入口: `skills/directing-implementation/SKILL.md`
- バグ修正の入口: `skills/directing-fixes/SKILL.md`
- workflow 鳥瞰図: `workflow.md`
- 補助 skill:
  - `skills/designing-implementation/SKILL.md`
  - `skills/distilling-implementation/SKILL.md`
  - `skills/planning-implementation/SKILL.md`
  - `skills/architecting-tests/SKILL.md`
  - `skills/reviewing-implementation/SKILL.md`
  - `skills/implementing-frontend/SKILL.md`
  - `skills/implementing-backend/SKILL.md`
  - `skills/distilling-fixes/SKILL.md`
  - `skills/tracing-fixes/SKILL.md`
  - `skills/analyzing-fixes/SKILL.md`
  - `skills/logging-fixes/SKILL.md`
  - `skills/implementing-fixes/SKILL.md`
  - `skills/reviewing-fixes/SKILL.md`
  - `skills/reporting-risks/SKILL.md`
  - `skills/diagramming-d2/SKILL.md`
  - `skills/diagramming-structure-diff/SKILL.md`
  - `skills/diagramming-plantuml/SKILL.md`
  - `skills/explore/SKILL.md`
  - `skills/skill-modification/SKILL.md`
  - `skills/updating-docs/SKILL.md`
  - `skills/working-light/SKILL.md`
- agent 契約:
  - `agents/structure_diagrammer.toml`
  - `agents/task_designer.toml`
  - `agents/ctx_loader.toml`
  - `agents/workplan_builder.toml`
  - `agents/test_architect.toml`
  - `agents/implementer.toml`
  - `agents/fault_tracer.toml`
  - `agents/log_instrumenter.toml`
  - `agents/review_cycler.toml`

## 標準フロー

### Impl lane

`User -> implementation proposal lane owner (`proposing-implementation`) -> MCP memory recall (`repo_conventions` / `recurring_pitfalls`) -> implementation distill skill (`distilling-implementation`) -> task-local design skill (`designing-implementation`) -> structure diagram agent (`structure_diagrammer`) + structure diagram diff skill (`diagramming-structure-diff`) -> human LGTM -> implementation execution lane owner (`directing-implementation`) -> implementation workplan skill (`planning-implementation`) -> test architecture skill (`architecting-tests`) -> frontend implementer (`implementing-frontend`) or backend implementer (`implementing-backend`) + assigned lint suite -> sonar-scanner + Sonar MCP open issue gate -> implementation review skill (`reviewing-implementation`) -> full harness -> 4humans sync + MCP memory distill + implementation execution lane owner (`directing-implementation`) close`

- implementation proposal lane owner (`proposing-implementation`) は実装要求を受け、日本語の active plan を作成し、重複確認と handoff に必要な最小限の入口情報だけを整える
- implementation proposal lane owner (`proposing-implementation`) は MCP memory bucket (`repo_conventions`, `recurring_pitfalls`) を recall 用に読み、今回の task に関係する項目だけを context summary へ持ち込む。MCP memory は repo 作法と再発失敗の recall に限定し、仕様や設計の正本代替には使わない
- implementation distill skill (`distilling-implementation`) は入口情報を起点に必要最小限の repo 文脈を探索し、facts、constraints、gaps、closeout notes、required reading を返す
- task-local design skill (`designing-implementation`) は distill 結果を前提に active plan の `UI` / `Scenario` / `Logic` だけを task-local design として固める
- task-local な設計は `docs/exec-plans/active/*.md` の中だけに置き、`changes/` や `context_board` は live 正本にしない
- structure diagram agent (`structure_diagrammer`) は GPT-5.4 / high で proposal の構造差分図を担当し、structure diagram diff skill (`diagramming-structure-diff`) を使って active exec-plan 配下に追加を緑、削除を赤で読める D2 / SVG を作る。active exec-plan の task-local design と既存 `diagrams/backend/` を照合し、既存図更新か new detail 図作成かを判断する。execution close では承認済み差分を `diagrams/backend/components.d2` と `diagrams/backend/<component>/<component>.d2` へ適用する
- human LGTM が active plan に記録されるまで implementation execution lane へ進めない
- implementation execution lane owner (`directing-implementation`) は承認済み active plan を受け取り、review 用差分図と差分正本適用先を参照しながら execution の handoff、gate、close を管理する
- implementation workplan skill (`planning-implementation`) は実装順、owned scope、validation、実装前に確認すべき relevant な `repo_conventions` / `recurring_pitfalls` を短い brief に落とす
- test architecture skill (`architecting-tests`) は active plan と関連仕様から、実装前に必要な failing tests、fixtures、validation commands を先に固定し、必要な test / fixture を最小範囲で実装する
- frontend implementer (`implementing-frontend`) / backend implementer (`implementing-backend`) は brief と plan に従って実装し、frontend では `python3 scripts/harness/run.py --suite frontend-lint`、backend では `python3 scripts/harness/run.py --suite backend-lint` だけを local validation として実行する
- `sonar-scanner + Sonar MCP open issue gate` は implementation execution lane owner (`directing-implementation`) が server-side analysis を更新し、その後に Sonar MCP の `search_sonar_issues_in_projects` を直接使って `project == ishibata91_AITranslationEngineJP` かつ `status == OPEN` の issue だけを gate 対象にして、issue が残る限り implementing skill へ差し戻す
- Sonar issue read の前提設定は Sonar CLI 認証ではなく、`codexmcps` profile に入った `mcp/sonarqube` の secret / config とする
- implementation review skill (`reviewing-implementation`) は単発で `仕様逸脱`、`例外処理`、`リソース解放`、`テスト不足`、`4humans` D2 sync 要否と実施有無 だけを見る
- review が `reroute` を返したら lane に差し戻すが、score 制の自動 review loop は持たない
- implementation execution lane owner (`directing-implementation`) は review が `pass` の後に `python3 scripts/harness/run.py --suite all` を final harness として実行する
- Sonar issue remediation loop は review の前段で、final harness は review の後段で implementation execution lane owner (`directing-implementation`) が扱い、どちらも close 条件に含める
- review が `pass` の時は `4humans sync` を整理し、backend 構造の変更または追加があった時は structure diagram agent (`structure_diagrammer`) を structure diagram diff skill (`diagramming-structure-diff`) で起動して承認済み差分を `diagrams/backend/` 正本へ適用する。処理の変更または追加があった時は `diagramming-d2` で `4humans/diagrams/processes/` の relevant `.d2` / `.svg` を更新し、構造の変更または追加があった時は `4humans/diagrams/structures/` の relevant `.d2` / `.svg` を更新してから close する
- `4humans/diagrams/processes/` または `4humans/diagrams/structures/` に new detail `.d2` を追加する時は、`4humans/diagrams/overview-manifest.json` を同じ変更で更新し、manifest で紐づいた overview `.d2` / `.svg` も同じ変更で更新する
- implementation execution lane owner (`directing-implementation`) は close 前に completed work から task-local ではない知識だけを MCP memory bucket (`repo_conventions` または `recurring_pitfalls`) へ蒸留し、次回 task で recall できる MCP memory を更新する
- review 用に active exec-plan 配下へ置いた差分 D2 / SVG は、`diagrams/backend/` 正本適用と `4humans` 正本同期が終わったら削除し、completed plan へ持ち越さない

### Fix lane

`User -> fix lane owner (`directing-fixes`) -> fix distill skill (`distilling-fixes`) -> fault trace skill (`tracing-fixes`) -> (必要時 logging skill (`logging-fixes`) / fix analysis skill (`analyzing-fixes`)) -> test architecture skill (`architecting-tests`) -> fix implementer (`implementing-fixes`) -> fix review skill (`reviewing-fixes`) -> risk reporting skill (`reporting-risks`) + 4humans sync + fix lane owner (`directing-fixes`) close`

- fix lane owner (`directing-fixes`) は bugfix 要求を受け、active plan を作成し、重複確認と handoff に必要な最小限の入口情報だけを整える
- fix distill skill (`distilling-fixes`) は入口情報を起点に必要最小限の repo 文脈を探索し、known facts、reproduction status、related constraints、related code pointers、open gaps、required reading を返す
- fault trace skill (`tracing-fixes`) は direction が整えた known facts と reproduction status を前提に、最小の trace 計画を返す
- logging skill (`logging-fixes`) は一時観測だけを追加 / 削除し、恒久修正を混ぜない
- fix analysis skill (`analyzing-fixes`) は観測結果を事実に圧縮し、fix 対象か `4humans sync` 対象か、または human-triggered な docs sync skill (`updating-docs`) 対象かを整理する
- test architecture skill (`architecting-tests`) は再現条件を tests / acceptance checks / validation commands に落とし、修正前に必要な回帰 test / fixture を最小範囲で実装する
- fix review skill (`reviewing-fixes`) も単発で `仕様逸脱`、`例外処理`、`リソース解放`、`テスト不足`、`4humans` D2 sync 要否と実施有無 だけを見る
- risk reporting skill (`reporting-risks`) は残留リスクを短くまとめる補助 skill として扱う
- review が `pass` の時は residual risk と `4humans sync` を整理し、実装の変更または追加があった時は diagramming D2 skill (`diagramming-d2`) で `4humans/diagrams/processes/` の relevant `.d2` / `.svg` を更新し、構造の変更または追加があった時は `4humans/diagrams/structures/` の relevant `.d2` / `.svg` を更新してから close する
- `4humans/diagrams/processes/` または `4humans/diagrams/structures/` に new detail `.d2` を追加する時は、`4humans/diagrams/overview-manifest.json` を同じ変更で更新し、manifest で紐づいた overview `.d2` / `.svg` も同じ変更で更新する

## 設計記録の扱い

- 非自明な変更は `docs/exec-plans/active/` に plan を置く
- 実装 task でだけ必要になる `UI` / `Scenario` / `Logic` は plan の中に section として置く
- 完了後も保持すべき詳細は `docs/` の正本、コード、型、tests、acceptance checks へ昇格する
- active plan を別 artifact 群へ分解しない
- `directing-* -> downstream skill` の handoff contract 例は、各 directing skill 配下の `references/*.json` を参照する
- `downstream skill -> directing-*` の返却 contract 例は、各 downstream skill 配下の `references/*.json` を参照する
- 各 skill の `references/permissions.json` は、その skill が実行してよい操作、してはいけない操作、期待される返却、停止条件を表す role contract として扱う
- skill の権限にない操作や解釈が曖昧な依頼は続行せず、stop and handoff を選ぶ
- reference JSON は説明用であり、live workflow の正本を packet 契約へ戻さない

## Review と reroute

- review は single-pass で、主観レビューや好みの改善提案を主目的にしない
- review の正式観点は `仕様逸脱`、`例外処理`、`リソース解放`、`テスト不足`、`4humans` D2 sync 要否と実施有無 の 5 つだけ
- `pass` か `reroute` の判定を返し、必要な再実行は lane 側で扱う
- 繰り返し見つかる指摘は review loop に残さず tests、harness、必要なら plan に昇格する

## 守ること

- live workflow に `architect-direction`、`light-direction`、`gating-workflow`、`context_board`、`tasks.md` を戻さない
- 過去 repo 由来で今の repo に合わない skill / agent / artifact 前提は、互換維持より削除を優先する
- 通常 lane の close 条件は `4humans sync` を含めて扱い、実装の変更または追加に伴う `4humans/diagrams/processes/` と構造の変更または追加に伴う `4humans/diagrams/structures/` の relevant `.d2` / `.svg` 更新も同一変更で完了させる
- new detail diagram 追加時の overview 更新要否は推測で決めず、`4humans/diagrams/overview-manifest.json` を正本として扱う
- `docs/` 正本更新は human が直接起動した `updating-docs` に限定する
- harness は repo-owned files だけを検査対象とし、`node_modules`、`dist`、`coverage`、`target`、生成物を含めない
