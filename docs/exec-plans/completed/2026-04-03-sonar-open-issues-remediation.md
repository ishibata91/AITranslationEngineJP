# Sonar Open Issues Remediation

- workflow: impl
- status: completed
- lane_owner: directing-implementation
- scope: `src/application/usecases/job-list/index.ts`, `src/application/usecases/job-create/index.ts`
- task_id:
- task_catalog_ref:
- parent_phase:

## Request Summary

- `sonar issueを全部解消して` の要求に対して、repo root の `sonar-scanner` と open issue gate を使って残存する Sonar issue を解消する。

## Decision Basis

- `directing-implementation` は Sonar remediation loop を close 条件として持つ。
- Sonar issue の正本は scanner 実行後の server-side analysis と helper script の `status == OPEN` 出力で確定する。
- 非自明な変更のため、実装前に active plan を置く。

## Owned Scope

- `src/application/usecases/job-list/index.ts`
- `src/application/usecases/job-create/index.ts`
- 関連 test は確認対象とし、変更は不要だった

## Out Of Scope

- human 先行が必要な `docs/` 正本変更
- Sonar project 側の suppression や severity tuning

## Dependencies / Blockers

- `sonar-scanner` と Sonar helper script がローカル環境で実行可能であること
- open issue の path と rule が implementing scope を安全に定義できること

## Parallel Safety Notes

- scope は scanner 実行後の touched path で再固定する
- shared config や root scripts に issue がある場合は downstream brief で ownership を明示する

## UI

- N/A
- 今回の Sonar open issue は `src/application/usecases/job-list/index.ts` の 3 件と `src/application/usecases/job-create/index.ts` の 4 件に限定され、screen 構造や表示仕様の変更は持ち込まない
- remediation は既存 UI から観測される job create / job list の振る舞いを維持したまま内部実装だけを整理する

## Scenario

- repo root の `sonar-scanner` 実行後、open issue は `src/application/usecases/job-list/index.ts` に `typescript:S3863` 2 件と `typescript:S3735` 1 件、`src/application/usecases/job-create/index.ts` に `typescript:S3735` 2 件、`typescript:S7746` 1 件、`typescript:S3776` 1 件ある状態を起点にする
- remediation loop では application usecase の observable behavior を変えずに、duplicate import、`void` operator、`Promise.resolve`、cognitive complexity の指摘を同一変更で解消する
- implementing 後は再度 `sonar-scanner` と open issue 取得を実行し、`status == OPEN` が 0 件になった時だけ review と close に進む

## Logic

- `src/application/usecases/job-list/index.ts` では duplicate import 2 件を 1 import に統合し、未使用引数の無視方法を `void` operator 以外へ置き換える
- `src/application/usecases/job-create/index.ts` では `defaultToErrorMessage` と validation loop の `void` operator 2 箇所を除去し、`initialize()` の `Promise.resolve` を通常 return へ置き換え、認知的複雑度が高い validation 関数を helper 抽出または分岐整理で分解する
- 既存の application port 契約、feature screen input 契約、job create / job list の tests を壊さず、repo の docs 正本変更は行わない

## Implementation Plan

- open Sonar issue を取得して対象 path と rule を確定する
- design、distill、workplan、test architecture を順に固定する
- implementing skill で issue を解消し、scanner と helper script を再実行する
- open issue が 0 件になった後に single-pass review、`4humans sync`、commit、plan close を行う

## Acceptance Checks

- `python3 scripts/harness/run.py --suite structure`
- `sonar-scanner`
- `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP`
- `python3 scripts/harness/run.py --suite all`

## Required Evidence

- remediation 前後の open Sonar issue 一覧
- implementing scope の差分
- full harness 結果

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`

## Outcome

- `src/application/usecases/job-list/index.ts` の duplicate import 2 件と `void` operator 由来の issue を解消した
- `src/application/usecases/job-create/index.ts` の `void` operator 2 件、`Promise.resolve` 1 件、cognitive complexity 1 件を解消した
- `node --require ./scripts/node/disable-vite-windows-net-use.cjs ./node_modules/vitest/vitest.mjs run --pool threads --config ./vite.config.mjs --configLoader native src/application/usecases/job-list/index.test.ts src/application/usecases/job-create/index.test.ts` は 2 file / 11 test passed だった
- `SONARQUBE_CLI_DIR=/tmp/sonarqube-cli SONARQUBE_CLI_TOKEN="$SONAR_TOKEN" SONARQUBE_CLI_SERVER=https://sonarcloud.io SONARQUBE_CLI_ORG=ishibata91 python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths src/application/usecases/job-list src/application/usecases/job-create` で owned-scope `openIssueCount: 0` を確認した
- `python3 scripts/harness/run.py --suite all` は passed だった
- `reviewing-implementation` の single-pass review は `pass` で、`4humans` の同期が必要な新規 debt / quality record は発生しなかった
