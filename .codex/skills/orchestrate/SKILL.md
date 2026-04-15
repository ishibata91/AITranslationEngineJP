---
name: orchestrate
description: AITranslationEngineJp 専用。唯一入口として task mode、primary skill-agent mapping、close 条件を固定し、role-based skill へ handoff する orchestrator。
---

# Orchestrate

この skill は live workflow の唯一入口です。
自身では product 実装、恒久修正、詳細 trace、docs 正本更新を抱えません。
役割は routing、primary handoff、gate 管理に限定します。
何があっても自身で調査・実装を始めないこと。必ず downstream skill に解決させること。

## 役割

- active work plan を作成または更新する
- `task_mode` を決め、判断根拠を plan に残す
- design bundle として `requirements`、`ui-mock`、`scenario`、`implementation-brief` を揃える
- design bundle 完了後に `functional_or_design_hitl: required-after-design-bundle` と `approval_record: pending-after-design-bundle` を記録し、human review 完了まで停止する
- `docs-only` 以外は human review 完了後に `implementation-scope` 以降の次工程を最小構成で選ぶ
- `docs-only` は human 承認済みの時だけ `updating-docs` へ handoff する
- `HITL`、required evidence、required validation、close 条件を管理する
- 存在する task-local artifact だけを `docs/` 正本へ反映する close summary を残す
- 広い task は downstream skill に投げる前に分割する

## Task Modes

- `implement`: 新機能、既存機能拡張、明確な振る舞い追加
- `fix`: bug、regressionの恒久修正
- `refactor`: 主目的が構造改善で、要件追加が主ではない変更
- `investigate`: まず evidence を集めるべき調査
- `docs-only`: human 承認済みの docs 正本変更

## Primary Skill-Agent Mapping

- `distill` -> `distiller`
- `design` -> `designer`
- `investigate` -> `investigator`
- `implement` -> `implementer`
- `tests` -> `tester`
- `review` -> `reviewer`
- `diagramming` -> `diagrammer`
- `updating-docs` -> `docs_updater`

## Routing Rules

- 共通ルール: 1 downstream skill には 1 primary agent だけを割り当てる
- 共通ルール: frontend を含む task は close 前に `review_mode: ui-check` を必須とする
- 共通ルール: 全 task で `review_mode: implementation-review` を必須とする

### `implement` / `refactor`

1. `distill` (`distiller`) に渡し、入口文脈を最小化する。
2. `design` (`designer`) で design bundle として `requirements`、`ui-mock`、`scenario`、`implementation-brief` を揃える。
3. design bundle 完了後に human review gate を立てる。
4. `review` (`reviewer`) の `design-review` と human review を通した後、`design` (`designer`) で `implementation-scope` を確定する。
5. 狭い `owned_scope` で `implement` (`implementer`) を実行する。
6. 必須 review が `pass` を返した後に正本同期し `close` する。

### `fix`

1. `distill` (`distiller`) に渡し、参照物と入口文脈を確定する。
2. `investigate` (`investigator`) で `reproduce` を行い、不具合を再現する。
3. 同じく `investigate` (`investigator`) で `trace` を行い、原因を解析する。
4. 必要な時だけ `temporary-logging` で観測点を補強する。
5. 修正後に `reobserve` と `review` で結果を確認する。
6. 修正方針が確定したら`implement`で修正し、`review`が`pass`を返したら`close`とする

### `investigate`

1. `distill` (`distiller`) で入口文脈を整理する。
2. `investigate` (`investigator`) で evidence を集める。
3. evidence のみで close してよい。

### `docs-only`

1. `distill` は通さない。
2. `approval_record` を確認する。
3. `updating-docs` (`docs_updater`) を起動する。

## Downstream Selection

- `distill`: implement、fix、refactor、investigate の入口整理
- `design`: requirements、ui-mock、scenario、implementation-brief、implementation-scope
- `investigate`: reproduce、trace、temporary-logging、reobserve、risk-report
- `implement`: frontend、backend、mixed の実装
- `tests`: scenario-implementation、unit
- `review`: design-review、ui-check、implementation-review
- `diagramming`: structure-diff、d2、plantuml
- `updating-docs`: human 承認済み docs-only の docs 正本更新

## Scope Rules

- 広い変更は orchestrate 側で frontend / backend / docs / review 単位へ分割する
- human review に必要な判断材料は design bundle 完了時点で active work plan へ固定し、human review 完了前に `implementation-scope` 以降へ渡さない
- 実装前の scope freeze は design の `implementation-scope` で行う
- 各 handoff には `owned_scope`、対象ファイル、完了条件、依存、validation を明示する
- depends_on が未解消の task は handoff しない
- compact 後も呼び出し元で確定済みの役割を引き継ぎ、配下 skill に再判定させない

## Stop Conditions

- plan が破綻している
- user 承認済み判断と衝突する
- skill 権限境界を超える
- design bundle 完了後に `functional_or_design_hitl` が `required-after-design-bundle` のまま、または `approval_record` が `pending-after-design-bundle` のままになっている
- narrow scope を安全に定義できない
- docs-only で `approval_record` がない

## close条件

- review が `pass` を返すこと
- backend を含む task は implement と review の両方で Sonar 件数ゲートを確認すること
- `HIGH` / `BLOCKER` の open issue が 0 件であること
- open reliability issue が 0 件であること
- open security issue が 0 件であること
- `ui_artifact_path` がある時だけ `docs/mocks/<page-id>/index.html` への反映を確認すること
- `scenario_artifact_path` がある時だけ `docs/scenario-tests/<topic-id>.md` への反映を確認すること
- `source_diagram_targets` がある時だけ `docs/architecture.md` と対象 D2 正本への反映を確認すること
- close summary に `canonicalized_artifacts` を残すこと

## Rules

- orchestrate 自身でコードを書かない
- orchestrate 自身で詳細調査を抱え込まない
- human review 未完了の design bundle を `implementation-scope` 以降の downstream handoff で迂回しない
- downstream skill は `fork_context: false` で呼ぶ
- primary skill-agent mapping を複数 skill で共有しない
- 別 skill を増やさない

## Reference Use

- quick overview は `references/orchestrate.to.<skill>.json` を使う
- mode 別 contract は `references/contracts/` を正本とする
- downstream からの返却は各 skill 側 `references/contracts/<skill>.to.orchestrate.<mode>.json` または single-role skill の `references/contracts/<skill>.to.orchestrate.json` を正本とする
- 旧名対応は `.codex/README.md` と `.codex/workflow.md` の対応表だけに残す
