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
- `docs-only` 以外は入口を `distill` とし、次工程を最小構成で選ぶ
- `docs-only` は human 承認済みの時だけ `updating-docs` へ handoff する
- `HITL`、required evidence、required validation、close 条件を管理する
- 存在する task-local artifact だけを `docs/` 正本へ反映する close summary を残す
- 広い task は downstream skill に投げる前に分割する

## Task Modes

- `implement`: 新機能、既存機能拡張、明確な振る舞い追加
- `fix`: bug、regression、narrow scope の恒久修正
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

- 1 downstream skill には 1 primary agent だけを割り当てる
- `implement` と `refactor` は `design-review` と human review の後に `implementation-scope` を確定し、狭い `owned_scope` で `implement` を実行する。`review` が `pass` を返した後に正本同期し `close` する
- `fix` は `reproduce` で不具合を再現し、`trace` で原因を解析し、必要なら `temporary-logging` を使って観測点を補強し、`reobserve` と `review` で修正を確認する
- `investigate` は evidence だけで close してよい
- `docs-only` は `distill` を通さず、`approval_record` を確認してから `updating-docs` を起動する
- frontend を含む task は close 前に `review_mode: ui-check` を必須とする
- 全 task で `review_mode: implementation-review` を必須とする

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
- 実装前の scope freeze は design の `implementation-scope` で行う
- 各 handoff には `owned_scope`、対象ファイル、完了条件、依存、validation を明示する
- depends_on が未解消の task は handoff しない
- compact 後も呼び出し元で確定済みの役割を引き継ぎ、配下 skill に再判定させない

## Stop Conditions

- plan が破綻している
- user 承認済み判断と衝突する
- skill 権限境界を超える
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
- downstream skill は `fork_context: false` で呼ぶ
- primary skill-agent mapping を複数 skill で共有しない
- 別 skill を増やさない

## Reference Use

- quick overview は `references/orchestrate.to.<skill>.json` を使う
- mode 別 contract は `references/contracts/` を正本とする
- downstream からの返却は各 skill 側 `references/contracts/<skill>.to.orchestrate.<mode>.json` または single-role skill の `references/contracts/<skill>.to.orchestrate.json` を正本とする
- 旧名対応は `.codex/README.md` と `.codex/workflow.md` の対応表だけに残す
