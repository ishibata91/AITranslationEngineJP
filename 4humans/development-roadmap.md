# Development Roadmap

関連文書: [`../docs/index.md`](../docs/index.md), [`../tasks/README.md`](../tasks/README.md), [`../tasks/phase-1/phase.yaml`](../tasks/phase-1/phase.yaml), [`../tasks/phase-1/tasks/P1-C01.yaml`](../tasks/phase-1/tasks/P1-C01.yaml), [`quality-score.md`](./quality-score.md), [`tech-debt-tracker.md`](./tech-debt-tracker.md)

このファイルは、2026-03-30 時点の repository 状態をもとに、人間向けの開発順序、現在地、直近 batch を整理する。
詳細な task 分解、依存関係、`owned_scope` は [`../tasks/`](../tasks/README.md) の YAML task catalog を正本とする。

## Status Legend

- `完了`: 現物または completed plan で成立を確認できる
- `進行候補`: 次の着手対象として妥当で、直近 batch に載せられる
- `未完了`: 必要性は明確だが、前段 task または検証が不足している

## Current Snapshot

| Area | Status | Notes |
|---|---|---|
| repository 骨格 | 完了 | `src/` と `src-tauri/` の bootstrap 構成がある |
| workflow / role 契約 | 完了 | `.codex/README.md` と workflow skills が正本になっている |
| structure harness | 完了 | required path と markdown link の検査入口がある |
| design harness | 完了 | semantic checks まで含めて成立している |
| xEdit import 入口 | 完了 | importer、入力 validation、raw `PLUGIN_EXPORT` cache が成立している |
| job skeleton | 進行候補 | 正規化 contract、job state contract、job create/list、UI、acceptance anchor が残っている |
| dictionary / persona foundation | 未完了 | task catalog はあるが implementation は未着手である |
| translation flow MVP | 未完了 | task catalog はあるが translation phases と preview は未着手である |
| provider / execution expansion | 未完了 | provider adapters と execution control の拡張は未着手である |
| output / release readiness | 未完了 | writer、cleanup、business-flow harness 統合は未着手である |

## Roadmap Policy

- `4humans/development-roadmap.md` は人間向け summary と immediate batch の正本である
- `tasks/phase-*/phase.yaml` は phase metadata と `parallel_batches` の正本である
- `tasks/phase-*/tasks/*.yaml` は task ID、依存関係、`owned_scope` の正本である
- batch は `contract -> verification -> impl -> integ` の順で固定し、同一 batch 内の `owned_scope` は重複させない
- `integ` は `composition / wiring / scenario proof` だけを持ち、新しい仕様判断を持ち込まない
- 完了判定は文書宣言ではなく、現物、tests、acceptance checks、validation commands、completed plan で行う

## Phase Summary

| Phase | Status | Focus | Task Catalog |
|---|---|---|---|
| Phase 0: Foundation Stabilization | 完了 | bootstrap、directory contract、harness、feature template | summary only |
| Phase 1: Input Cache And Job Skeleton | 進行候補 | `TRANSLATION_UNIT` canonical contract、job state contract、job create/list、first acceptance path | [`phase-1/phase.yaml`](../tasks/phase-1/phase.yaml) |
| Phase 2: Dictionary And Persona Foundation | 未完了 | xTranslator import、master dictionary、master persona、observation UI | [`phase-2/phase.yaml`](../tasks/phase-2/phase.yaml) |
| Phase 3: Translation Flow MVP | 未完了 | instruction builder、word/persona/body translation、preview、regression | [`phase-3/phase.yaml`](../tasks/phase-3/phase.yaml) |
| Phase 4: Provider And Execution Expansion | 未完了 | provider adapters、execution control、failure/retry acceptance | [`phase-4/phase.yaml`](../tasks/phase-4/phase.yaml) |
| Phase 5: Output And Release Readiness | 未完了 | writers、artifact registry、cleanup、harness business-flow checks | [`phase-5/phase.yaml`](../tasks/phase-5/phase.yaml) |

## Immediate Next Batches

### Batch P1-B1

- `進行候補`: `P1-C01` `TRANSLATION_UNIT canonical contract`
- `進行候補`: `P1-C02` minimal job state model contract
- `進行候補`: `P1-V01` lossless translation-unit preservation fixture
- `進行候補`: `P1-V02` import-to-job acceptance anchor
- 理由: `contract` と `verification` を先に固定し、後続の backend / frontend 実装 batch が shared decision を持たない状態にする

### Batch P1-B2

- `未完了`: `P1-I04` backend job creation usecase
- `未完了`: `P1-I05` backend job list query
- `未完了`: `P1-I06` job create screen
- `未完了`: `P1-I07` job list screen
- 条件: `P1-B1` 完了後に着手する
- 理由: backend create、backend list、frontend create、frontend list を別 `owned_scope` で平行に進める

### Batch P1-B3

- `未完了`: `P1-G01` import-to-job integrated scenario
- 条件: `P1-B2` 完了後に着手する
- 理由: `integ` を最後に隔離し、shared wiring と scenario proof だけに閉じ込める

## Current Risks

- `未完了`: importer 周辺を除く translation-domain 固有の tests / fixtures / acceptance checks はまだ薄い
- `未完了`: `tasks/` の YAML schema は手作業運用であり、lint や structural validation はまだない
- `未完了`: `contract` と `verification` を飛ばして `impl` から始めると、`integ` が巨大化して並列安全性が崩れる

## Done Definition

- `未完了`: 非自明な active plan は `task_id` と `owned_scope` を持ち、必要なら `tasks/phase-*/phase.yaml` または `tasks/phase-*/tasks/*.yaml` を参照している
- `未完了`: 振る舞いが変わる task は、同じ変更で対応する tests / acceptance checks / validation commands を更新している
- `未完了`: 完了した非自明 task は completed plan に結果が残っている
