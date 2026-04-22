---
name: review-implementation
description: GitHub Copilot 側の implementation review 知識 package。
---

# Review Implementation

## 目的

この skill は知識 package である。
`reviewer` agent が implementation review を行う時に、差分が owned_scope と一致するか、test と validation が十分かを確認する判断基準を提供する。

## いつ参照するか

- 差分が owned_scope に収まっているか確認する時
- 必要な product test と validation command を確認する時
- backend を含む場合に repo-local Sonar issue gate を確認する時

## 参照しない場合

- UI check が主目的の時
- design review が必要な時
- 修正を同時に行う時

## 原則

- single_handoff_packet、lane_context_packet、review target diff を照合する
- 好みや将来改善で reroute しない
- finding は再現できる形で返す
- 修正は行わない

## DO / DON'T

DO:
- scope 外 diff、missing test、failed validation を分ける
- repo-local Sonar issue gate の該当可否を明示する
- pass の場合も未実行 validation を残す

DON'T:
- design review をしない
- 新しい要件解釈を追加しない
- active contract をこの skill に置かない

## Sonar Issue Gate 確認方法

backend を含む handoff の review では、Sonar サーバ側の Quality Gate ではなく repo-local gate を確認する。
repo-local gate は coverage 目標と Sonar issue 数量制限で構成する。

1. `mcp_mcp_docker_mcp-exec` で `python3 scripts/harness/run.py --suite all` を実行する。
2. `--suite all` は `check_coverage.py` を含み、内部で `sonar-scanner` を起動する。
3. `sonar-scanner` 完了後、SonarQube API をポーリングして以下の repo-local gate 条件を確認する:
   - coverage >= 70.0%
   - Security issue 0 件
   - Reliability issue 0 件
   - Maintainability HIGH/BLOCKER issue 0 件
4. issue 確認は `check_coverage.py` が `api/issues/search` エンドポイントで自動チェックする。
5. `check_coverage.py` が PASS を返した場合に repo-local gate PASS とみなす。
6. `sonar-project.properties` が repo root に存在することが前提条件。

`npm run test:system` または `harness all` が Wails、sandbox、OS 権限で止まる場合は `FAIL_ENVIRONMENT` として扱う。
`FAIL_ENVIRONMENT` は product failure として reroute せず、blocked reason、再実行環境、再実行コマンドを residual risk に残す。

## Checklist

- [review-implementation-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/review-implementation/references/checklists/review-implementation-checklist.md) を参照する。
