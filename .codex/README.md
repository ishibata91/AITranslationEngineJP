# .codex

このディレクトリは、AITranslationEngineJp の live なマルチエージェント作業フローの正本です。
プロダクト仕様と設計は `docs/` を正本とし、lane、skill、agent の役割と handoff は `.codex/` を正本とします。
この repo の workflow は `directing-implementation` と `directing-fixes` の 2 lane で動かし、過去 repo 固有の packet や review loop は live 契約に戻しません。

## 入口

- 実装・設計内包の入口: `skills/directing-implementation/SKILL.md`
- バグ修正の入口: `skills/directing-fixes/SKILL.md`
- workflow 鳥瞰図: `workflow.md`
- 補助 skill:
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
  - `skills/diagramming-plantuml/SKILL.md`
- agent 契約:
  - `agents/ctx_loader.toml`
  - `agents/workplan_builder.toml`
  - `agents/test_architect.toml`
  - `agents/implementer.toml`
  - `agents/fault_tracer.toml`
  - `agents/log_instrumenter.toml`
  - `agents/review_cycler.toml`

## 標準フロー

### Impl lane

`User -> directing-implementation -> distilling-implementation -> planning-implementation -> architecting-tests -> implementing-frontend or implementing-backend -> reviewing-implementation -> directing-implementation close`

- `directing-implementation` は実装要求を受け、必要なら active plan の中に `UI` / `Scenario` / `Logic` を埋める
- task-local な設計は `docs/exec-plans/active/*.md` の中だけに置き、`changes/` や `context_board` は live 正本にしない
- `distilling-implementation` は facts、constraints、gaps、docs sync 候補を整理する
- `planning-implementation` は実装順、owned scope、validation を短い brief に落とす
- `architecting-tests` は active plan と関連仕様から、実装前に必要な failing tests、fixtures、validation commands を先に固定する
- `implementing-frontend` / `implementing-backend` は brief と plan に従って実装する
- `reviewing-implementation` は単発で `仕様逸脱`、`例外処理`、`リソース解放`、`テスト不足` だけを見る
- review が `reroute` を返したら lane に差し戻すが、score 制の自動 review loop は持たない

### Fix lane

`User -> directing-fixes -> distilling-fixes -> tracing-fixes -> (必要時 logging-fixes / analyzing-fixes) -> architecting-tests -> implementing-fixes -> reviewing-fixes -> directing-fixes close`

- `directing-fixes` は bugfix 要求を受け、事実不足なら `distilling-fixes` と `tracing-fixes` で scope を狭める
- `logging-fixes` は一時観測だけを追加 / 削除し、恒久修正を混ぜない
- `analyzing-fixes` は観測結果を事実に圧縮し、fix 対象か docs sync 対象かを整理する
- `architecting-tests` は再現条件を tests / acceptance checks / validation commands に落とし、修正前に回帰テストを準備する
- `reviewing-fixes` も単発で `仕様逸脱`、`例外処理`、`リソース解放`、`テスト不足` だけを見る
- `reporting-risks` は残留リスクを短くまとめる補助 skill として扱う

## 設計記録の扱い

- 非自明な変更は `docs/exec-plans/active/` に plan を置く
- 実装 task でだけ必要になる `UI` / `Scenario` / `Logic` は plan の中に section として置く
- 完了後も保持すべき詳細は `docs/` の正本、コード、型、tests、acceptance checks へ昇格する
- active plan を別 artifact 群へ分解しない

## Review と reroute

- review は single-pass で、主観レビューや好みの改善提案を主目的にしない
- review の正式観点は `仕様逸脱`、`例外処理`、`リソース解放`、`テスト不足` の 4 つだけ
- `pass` か `reroute` の判定を返し、必要な再実行は lane 側で扱う
- 繰り返し見つかる指摘は review loop に残さず tests、harness、必要なら plan に昇格する

## 守ること

- live workflow に `architect-direction`、`light-direction`、`gating-workflow`、`context_board`、`tasks.md` を戻さない
- 過去 repo 由来で今の repo に合わない skill / agent / artifact 前提は、互換維持より削除を優先する
- docs sync は lane の close 条件として扱い、別の人手前提 lane に押し戻さない
- harness は repo-owned files だけを検査対象とし、`node_modules`、`dist`、`coverage`、`target`、生成物を含めない


