# .codex

このディレクトリは、AITranslationEngineJp の live なマルチエージェント作業フローの正本です。
プロダクト仕様と設計は `docs/` を正本とし、lane、skill、agent の役割と handoff は `.codex/` を正本とします。
この repo の workflow は `directing-implementation` と `directing-fixes` の 2 lane で動かし、過去 repo 固有の packet や review loop は live 契約に戻しません。

## Naming Rule

- workflow 文書では、論理名と実名を分離しない
- 初出または重要な参照は `論理名 (`actual-name`)` を優先する
- 人間 review で意味が先に読めて、actual skill / agent name でも検索できる記述を優先する
- 例: implementation lane owner (`directing-implementation`)、fix lane owner (`directing-fixes`)、task-local design skill (`designing-implementation`)

## 入口

- 実装・設計内包の入口: `skills/directing-implementation/SKILL.md`
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
  - `skills/diagramming-plantuml/SKILL.md`
  - `skills/explore/SKILL.md`
  - `skills/skill-modification/SKILL.md`
  - `skills/updating-docs/SKILL.md`
  - `skills/working-light/SKILL.md`
- agent 契約:
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

`User -> implementation lane owner (`directing-implementation`) -> implementation distill skill (`distilling-implementation`) -> task-local design skill (`designing-implementation`) -> implementation workplan skill (`planning-implementation`) -> test architecture skill (`architecting-tests`) -> frontend implementer (`implementing-frontend`) or backend implementer (`implementing-backend`) -> sonar-scanner + Docker MCP Sonar open issue gate -> implementation review skill (`reviewing-implementation`) -> 4humans sync + implementation lane owner (`directing-implementation`) close`

- implementation lane owner (`directing-implementation`) は実装要求を受け、active plan を作成し、重複確認と handoff に必要な最小限の入口情報だけを整える
- implementation distill skill (`distilling-implementation`) は入口情報を起点に必要最小限の repo 文脈を探索し、facts、constraints、gaps、closeout notes、required reading を返す
- task-local design skill (`designing-implementation`) は distill 結果を前提に active plan の `UI` / `Scenario` / `Logic` だけを task-local design として固める
- task-local な設計は `docs/exec-plans/active/*.md` の中だけに置き、`changes/` や `context_board` は live 正本にしない
- implementation workplan skill (`planning-implementation`) は実装順、owned scope、validation を短い brief に落とす
- test architecture skill (`architecting-tests`) は active plan と関連仕様から、実装前に必要な failing tests、fixtures、validation commands を先に固定し、必要な test / fixture を最小範囲で実装する
- frontend implementer (`implementing-frontend`) / backend implementer (`implementing-backend`) は brief と plan に従って実装する
- `sonar-scanner + Docker MCP Sonar open issue gate` は server-side analysis を更新し、`docker mcp tools call search_sonar_issues_in_projects --gateway-arg=--profile --gateway-arg=codexmcps` を使う helper script から `project == ishibata91_AITranslationEngineJP` かつ `status == OPEN` の issue だけを gate 対象にして、issue が残る限り implementing skill へ差し戻す
- Sonar issue read の前提設定は Sonar CLI 認証ではなく、`codexmcps` profile に入った `mcp/sonarqube` の secret / config とする
- implementation review skill (`reviewing-implementation`) は単発で `仕様逸脱`、`例外処理`、`リソース解放`、`テスト不足` だけを見る
- review が `reroute` を返したら lane に差し戻すが、score 制の自動 review loop は持たない
- Sonar issue remediation loop は review の前段で扱い、close 条件に含める
- review が `pass` の時は `4humans sync` を整理し、コードベース境界や実行フローが変わる時は diagramming D2 skill (`diagramming-d2`) で `4humans/class-diagrams/` と `4humans/sequence-diagrams/` の `.d2` / `.svg` を更新してから close する

### Fix lane

`User -> fix lane owner (`directing-fixes`) -> fix distill skill (`distilling-fixes`) -> fault trace skill (`tracing-fixes`) -> (必要時 logging skill (`logging-fixes`) / fix analysis skill (`analyzing-fixes`)) -> test architecture skill (`architecting-tests`) -> fix implementer (`implementing-fixes`) -> fix review skill (`reviewing-fixes`) -> risk reporting skill (`reporting-risks`) + 4humans sync + fix lane owner (`directing-fixes`) close`

- fix lane owner (`directing-fixes`) は bugfix 要求を受け、active plan を作成し、重複確認と handoff に必要な最小限の入口情報だけを整える
- fix distill skill (`distilling-fixes`) は入口情報を起点に必要最小限の repo 文脈を探索し、known facts、reproduction status、related constraints、related code pointers、open gaps、required reading を返す
- fault trace skill (`tracing-fixes`) は direction が整えた known facts と reproduction status を前提に、最小の trace 計画を返す
- logging skill (`logging-fixes`) は一時観測だけを追加 / 削除し、恒久修正を混ぜない
- fix analysis skill (`analyzing-fixes`) は観測結果を事実に圧縮し、fix 対象か `4humans sync` 対象か、または human-triggered な docs sync skill (`updating-docs`) 対象かを整理する
- test architecture skill (`architecting-tests`) は再現条件を tests / acceptance checks / validation commands に落とし、修正前に必要な回帰 test / fixture を最小範囲で実装する
- fix review skill (`reviewing-fixes`) も単発で `仕様逸脱`、`例外処理`、`リソース解放`、`テスト不足` だけを見る
- risk reporting skill (`reporting-risks`) は残留リスクを短くまとめる補助 skill として扱う
- review が `pass` の時は residual risk と `4humans sync` を整理し、コードベース境界や実行フローが変わる時は diagramming D2 skill (`diagramming-d2`) で `4humans/class-diagrams/` と `4humans/sequence-diagrams/` の `.d2` / `.svg` を更新してから close する

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
- review の正式観点は `仕様逸脱`、`例外処理`、`リソース解放`、`テスト不足` の 4 つだけ
- `pass` か `reroute` の判定を返し、必要な再実行は lane 側で扱う
- 繰り返し見つかる指摘は review loop に残さず tests、harness、必要なら plan に昇格する

## 守ること

- live workflow に `architect-direction`、`light-direction`、`gating-workflow`、`context_board`、`tasks.md` を戻さない
- 過去 repo 由来で今の repo に合わない skill / agent / artifact 前提は、互換維持より削除を優先する
- 通常 lane の close 条件は `4humans sync` を含めて扱い、必要な `4humans/class-diagrams/` と `4humans/sequence-diagrams/` の `.d2` / `.svg` 更新も同一変更で完了させる
- `docs/` 正本更新は human が直接起動した `updating-docs` に限定する
- harness は repo-owned files だけを検査対象とし、`node_modules`、`dist`、`coverage`、`target`、生成物を含めない
