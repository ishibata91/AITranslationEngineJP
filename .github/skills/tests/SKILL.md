---
name: tests
description: GitHub Copilot 側の product test 共通知識 package。承認済み owned_scope を test で証明する判断基準を提供する。
---

# Tests

## 目的

`tests` は知識 package である。
`tester` agent が、single_handoff_packet と lane_context_packet の owned_scope を product test で証明する時の共通判断を提供する。

実行権限、write scope、active contract、handoff は [tester.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/tester.agent.md) が持つ。

## いつ参照するか

- 承認済み handoff または実装済み scope を product test で証明する時
- scenario artifact または unit responsibility を test に落とす時
- fake provider や DI seam で paid real AI API を避ける時

## 参照しない場合

- product code の恒久修正が主目的の時
- design や scenario artifact を新規に作る時
- review だけを行う時

## 知識範囲

- test の責務分割
- deterministic setup
- coverage 70% と harness coverage suite の扱い
- Arrange / Act / Assert の読みやすさ
- focused skill の選び方

## 原則

- 各 test は 1 つの振る舞いを扱う
- setup は決定的にする
- paid real AI API を呼ばない
- 新しい要件解釈を足さない

## Focused Skills

- [tests-scenario](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/tests-scenario/SKILL.md): scenario artifact の product test 化
- [tests-unit](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/tests-unit/SKILL.md): unit test による責務補強

## DO / DON'T

DO:
- fixture の入力値、clock、runtime 応答、seed を固定する
- test body に条件分岐を入れない
- `python3 scripts/harness/run.py --suite coverage` で Sonar-compatible coverage 70% 以上を確認する
- validation result と remaining gaps を返す

DON'T:
- product code を広く直さない
- docs、`.codex`、`.github/skills`、`.github/agents` を変更しない
- mode 別 active contract を使わない

## 参照パターン

- [test-patterns.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/tests/references/patterns/test-patterns.md) を参照する。
- 対象は Red / Green / Refactor、edge-case inventory、AAA、flaky test avoidance、E2E artifact handling である。
- coverage は repo の `MINIMUM_COVERAGE = 70.0`、validation command、Wails test seam に従う。

## Checklist

- [tests-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/tests/references/checklists/tests-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。

## Agent が持つもの

- active contract: [tester.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/references/tester/contracts/tester.contract.json)
- permissions: [permissions.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/references/tester/permissions.json)

## Maintenance

- scenario / unit の知識差分は focused skill に置く。
- output obligation を skill 本体へ戻さない。
- 旧 mode contract は active 正本として扱わない。
