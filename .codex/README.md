# .codex

このディレクトリは、AITranslationEngineJp の live なマルチエージェント作業フローの正本です。
プロダクト仕様と設計は `docs/` を正本とし、lane、skill、agent の役割と handoff は `.codex/` を正本とします。
この repo の workflow は `impl-direction` と `fix-direction` の 2 lane で動かし、過去 repo 固有の packet や review loop は live 契約に戻しません。

## 入口

- 実装・設計内包の入口: `skills/impl-direction/SKILL.md`
- バグ修正の入口: `skills/fix-direction/SKILL.md`
- 補助 skill:
  - `skills/impl-distill/SKILL.md`
  - `skills/impl-workplan/SKILL.md`
  - `skills/test-architect/SKILL.md`
  - `skills/impl-review/SKILL.md`
  - `skills/impl-frontend-work/SKILL.md`
  - `skills/impl-backend-work/SKILL.md`
  - `skills/fix-distill/SKILL.md`
  - `skills/fix-trace/SKILL.md`
  - `skills/fix-analysis/SKILL.md`
  - `skills/fix-logging/SKILL.md`
  - `skills/fix-work/SKILL.md`
  - `skills/fix-review/SKILL.md`
  - `skills/risk-report/SKILL.md`
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

`User -> impl-direction -> impl-distill -> impl-workplan -> test-architect -> impl-work -> impl-review -> impl-direction close`

- `impl-direction` は実装要求を受け、必要なら active plan の中に `UI` / `Scenario` / `Logic` を埋める
- task-local な設計は `docs/exec-plans/active/*.md` の中だけに置き、`changes/` や `context_board` は live 正本にしない
- `impl-distill` は facts、constraints、gaps、docs sync 候補を整理する
- `impl-workplan` は実装順、owned scope、validation を短い brief に落とす
- `test-architect` は active plan と関連仕様から、実装前に必要な failing tests、fixtures、validation commands を先に固定する
- `impl-frontend-work` / `impl-backend-work` は brief と plan に従って実装する
- `impl-review` は単発で `仕様逸脱`、`例外処理`、`リソース解放`、`テスト不足` だけを見る
- review が `reroute` を返したら lane に差し戻すが、score 制の自動 review loop は持たない

### Fix lane

`User -> fix-direction -> fix-distill -> fix-trace -> (必要時 fix-logging / fix-analysis) -> test-architect -> fix-work -> fix-review -> fix-direction close`

- `fix-direction` は bugfix 要求を受け、事実不足なら `fix-distill` と `fix-trace` で scope を狭める
- `fix-logging` は一時観測だけを追加 / 削除し、恒久修正を混ぜない
- `fix-analysis` は観測結果を事実に圧縮し、fix 対象か docs sync 対象かを整理する
- `test-architect` は再現条件を tests / acceptance checks / validation commands に落とし、修正前に回帰テストを準備する
- `fix-review` も単発で `仕様逸脱`、`例外処理`、`リソース解放`、`テスト不足` だけを見る
- `risk-report` は残留リスクを短くまとめる補助 skill として扱う

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

- live workflow に `architect-direction`、`light-direction`、`workflow-gate`、`context_board`、`tasks.md` を戻さない
- 過去 repo 由来で今の repo に合わない skill / agent / artifact 前提は、互換維持より削除を優先する
- docs sync は lane の close 条件として扱い、別の人手前提 lane に押し戻さない
- harness は repo-owned files だけを検査対象とし、`node_modules`、`dist`、`coverage`、`target`、生成物を含めない
