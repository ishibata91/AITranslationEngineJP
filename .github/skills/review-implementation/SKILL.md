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
- backend を含む場合に Sonar gate を確認する時

## 参照しない場合

- UI check が主目的の時
- design review が必要な時
- 修正を同時に行う時

## 原則

- implementation-scope と review target diff を照合する
- 好みや将来改善で reroute しない
- finding は再現できる形で返す
- 修正は行わない

## DO / DON'T

DO:
- scope 外 diff、missing test、failed validation を分ける
- Sonar gate の該当可否を明示する
- pass の場合も未実行 validation を残す

DON'T:
- design review をしない
- 新しい要件解釈を追加しない
- active contract をこの skill に置かない

## Sonar Gate 確認方法

backend を含む handoff の review では、次の手順で Sonar gate を確認する。

1. `mcp_mcp_docker_mcp-exec` で `python3 scripts/harness/run.py --suite all` を実行する。
2. `--suite all` は `check_coverage.py` を含み、内部で `sonar-scanner` を起動する。
3. `sonar-scanner` 完了後、SonarQube API をポーリングして以下のゲート条件を確認する:
   - coverage >= 70.0%
   - Security issue 0 件
   - Reliability issue 0 件
   - Maintainability HIGH issue 0 件
4. issue 確認は `check_coverage.py` が `api/issues/search` エンドポイントで自動チェックする。
5. `check_coverage.py` が PASS を返した場合に Sonar gate PASS とみなす。
6. `sonar-project.properties` が repo root に存在することが前提条件。

## Checklist

- [review-implementation-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/review-implementation/references/checklists/review-implementation-checklist.md) を参照する。
