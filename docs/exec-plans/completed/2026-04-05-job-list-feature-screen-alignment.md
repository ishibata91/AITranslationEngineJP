# Impl Plan Template

- workflow: impl
- status: completed
- lane_owner: directing-implementation
- scope: align job-list screen usecase with feature-screen reusable base and remove duplicated orchestration
- task_id:
- task_catalog_ref:
- parent_phase:

## Request Summary

- `job-list` だけを対象に、既存の `feature-screen` 基盤へ寄せて duplicate を減らす。
- observe 系 3 usecase は対象外とし、将来の仕様差分を先に観測する。

## Decision Basis

- Sonar duplication は `src/application/usecases/job-list/index.ts` と `src/application/usecases/feature-screen/index.ts` の orchestration 重複を検出している。
- `docs/exec-plans/completed/2026-03-29-feature-template-scaffold.md` は `feature-screen` を reusable base として扱う方針を正本化している。
- `job-list` は `initialize` / `refresh` / `retry` / `select` / `updateFilters` と `loadCurrent` の流れが `feature-screen` とほぼ一致しており、既存 generic へ寄せても仕様差分を作らない。

## Owned Scope

- `src/application/usecases/job-list/`
- `src/application/usecases/feature-screen/`
- `src/application/usecases/job-list/index.test.ts`
- `src/main.ts`

## Out Of Scope

- `dictionary-observe` / `persona-observe` / `translation-preview` の抽象化
- `docs/` 正本更新
- backend 実装変更

## Dependencies / Blockers

- `feature-screen` generic が `job-list` の request なし load と selection reconcile を表現できること

## Parallel Safety Notes

- `src/main.ts` は他 task と競合しやすいので差分最小で扱う。
- `job-list` の screen store factory の配置が変わる場合は import path 影響を確認する。

## UI

- UI の表示構造と interaction は変更しない。

## Scenario

- `job-list` は初期表示、`refresh`、`retry` の 3 entry を、`feature-screen` の同一 load orchestration へ委譲する。`job-list` 側で分岐を増やさず、従来どおり毎回同じ一覧再取得を行う。
- `updateFilters` は既存 contract を維持し、store へ `filters` を反映した上で `reload === true` の時だけ同じ load orchestration を再実行する。`job-list` の `filters` は引き続き `undefined` のままとし、画面挙動は増やさない。
- selection reconcile は `job-list` 固有ロジックとして残す。再取得結果に現在の `jobId` が含まれる時だけ選択を維持し、含まれない時は先頭 job、結果が空の時は `null` へ戻す。
- 失敗時の公開挙動は変えない。load 開始時は既存 data / selection を保持したまま loading に入り、失敗時は最後の成功データを残したまま generic な job-list 向けエラーメッセージだけを更新する。

## Logic

- `createJobListScreenUsecase` は thin wrapper に縮小し、内部実装を `createFeatureScreenUsecase` 呼び出しへ置き換える。公開型 `JobListScreenInput`、`JobListScreenStore`、`createJobListScreenStore` はそのまま残し、`job-list` の import 面を変えない。
- `executor: () => Promise<JobListResult>` は `job-list` 内でだけ `FeatureScreenGateway<undefined, JobListResult>` へ包む。`createRequest` は `() => undefined` を返し、request なし load を generic 既存 API で表現する。
- `reconcileSelection` は `job-list` ローカル関数として維持し、`feature-screen` へ callback として渡す。選択維持 / 先頭 fallback / empty success の判断は gateway や store へ移さない。
- `toErrorMessage` の既定値は `job-list` 側の generic 文言を維持し、transport 由来の生エラー文字列を外へ出さない。`feature-screen` の default error mapping は変更しない。
- 今回は `feature-screen` の generic API 拡張や helper 追加を行わない。既存の `createRequest + gateway.load + reconcileSelection` で `job-list` 要件を満たせる前提で進め、必要な型調整が出ても `job-list` ローカル alias に閉じ込める。
- duplicate 削減の対象は `job-list` の `loadCurrent` orchestration のみとし、`dictionary-observe` / `persona-observe` / `translation-preview` へ同じ整理を波及させない。

## Implementation Plan

- `job-list` の current orchestration と `feature-screen` generic の対応点を確定する。
- `job-list` usecase を generic 委譲へ置き換え、必要なら request/gateway bridge を最小追加する。
- 既存 tests を generic 委譲後の振る舞いに合わせて更新し、重複退避と回帰防止を確認する。
- frontend lint suite、Sonar issue gate、single-pass review、full harness を通して close する。

## Acceptance Checks

- `job-list` の公開挙動が変わらず、generic usecase 経由で初期 load / refresh / retry / filter reload / selection reconcile が通る。
- `src/application/usecases/job-list/index.ts` から duplicate orchestration が除去される。
- touched scope で OPEN Sonar issue を残さない。

## Required Evidence

- `python3 scripts/harness/run.py --suite frontend-lint`
- `sonar-scanner`
- `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths ...`
- `python3 scripts/harness/run.py --suite all`

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`
- D2 update 不要見込み。構造図とプロセス図の責務分割、実行順序、依存方向は変えない。

## Outcome

- `src/application/usecases/job-list/index.ts` の `loadCurrent` orchestration を除去し、`createFeatureScreenUsecase` への thin wrapper に置き換えた。
- `job-list` ローカル wrapper で latest selection を `setLoaded` 時に再評価し、empty success と in-flight selection change の既存挙動を維持した。
- `src/application/usecases/job-list/index.test.ts` に `updateFilters(reload)` と refresh / reload 中 selection change の回帰 test を追加した。
- `feature-screen` API、`main.ts` wiring、observe 系 usecase、docs 正本、D2 は変更しなかった。
- Validation passed: `python3 scripts/harness/run.py --suite structure`, `npm run test -- src/application/usecases/job-list/index.test.ts`, `python3 scripts/harness/run.py --suite frontend-lint`, `sonar-scanner`, `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths src/application/usecases/job-list/index.ts src/application/usecases/job-list/index.test.ts`, `python3 scripts/harness/run.py --suite all`.
