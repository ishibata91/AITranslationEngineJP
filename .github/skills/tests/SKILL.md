---
name: tests
description: GitHub Copilot 側の product test 共通知識 package。承認済み owned_scope を test で証明する判断基準を提供する。
---

# Tests

## 目的

`tests` は知識 package である。
`tester` agent が、single_handoff_packet と tester_context_packet の owned_scope を product test で証明する時の共通判断を提供する。

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
- tester_context_packet の読み順
- insufficient_context の返し方
- deterministic setup
- final validation lane へ defer する coverage / harness all の扱い
- Arrange / Act / Assert の読みやすさ
- focused skill の選び方

## 原則

- 各 test は 1 つの振る舞いを扱う
- tester_context_packet の test_ingredients、test_required_reading、requirements_policy_decisions の test impact、test_validation_entry の順に読む
- listed files / symbols 外を探索して context 不足を埋めない
- insufficient_context_criteria は structural gate とし、behavior_to_prove、public seam、test target、assertion focus、fixture/helper 方針、focused validation の不足時に返す
- `UI人間操作E2E` で開始操作、検証対象の入口、入力模倣方針が不足する場合は insufficient_context を返す
- `APIテスト` で public seam、request / response contract、入力開始点、主要観測点が不足する場合は insufficient_context を返す
- test_subscope が completion_signal clause、public seam、test target file、validation command のいずれにも対応しない場合は insufficient_context を返す
- 承認済み scenario を元に期待どおり fail する test、局所的 import 修正、既存 test file 内の軽微な確認は not_insufficient_context として扱う
- 原因未確定の regression test は実装前に書かない
- setup は決定的にする
- test 追加または更新後、handoff を終える前に touched layer に対応する local validation を実行する
- paid real AI API を呼ばない
- 新しい要件解釈を足さない

## Focused Skills

- [tests-scenario](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/tests-scenario/SKILL.md): scenario artifact の product test 化
- [tests-unit](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/tests-unit/SKILL.md): unit test による責務補強

## DO / DON'T

DO:
- test_ingredients の completion_signal clause、behavior_to_prove、public seam、assertion_focus に沿って test を作る
- `UI人間操作E2E` では、承認済みシナリオの開始操作と入力模倣方針に沿って試験を作る
- `APIテスト` では、承認済み受け入れ条件、public seam、request / response contract、入力開始点、主要観測点に沿って試験を作る
- test_subscope が渡された場合はその sub-scope だけを証明し、残りを remaining_test_subscopes に残す
- insufficient_context を返す場合は reason、needed_context、remaining_test_subscopes を structural gate に対応づける
- fixture の入力値、clock、runtime 応答、seed を固定する
- test body に条件分岐を入れない
- coverage、harness all、repo-local Sonar issue gate は final validation lane へ defer する
- backend handoff は `python3 scripts/harness/run.py --suite backend-local`、frontend handoff は `python3 scripts/harness/run.py --suite frontend-local` を使う
- mixed handoff は touched layer に応じて両方を実行する
- validation result と remaining gaps を返す

DON'T:
- full lane_context_packet、fix_ingredients、change_targets、broad related_code_pointers を直接追わない
- insufficient_context を返さず広く調査しない
- UI 入口の `UI人間操作E2E` を裏側の直接呼び出しだけで代替しない
- criteria mismatch になる不安や通常の局所確認を insufficient_context にしない
- product code を広く直さない
- docs、`.codex`、`.github/skills`、`.github/agents` を変更しない
- mode 別 active contract を使わない

## 参照パターン

- [test-patterns.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/tests/references/patterns/test-patterns.md) を参照する。
- 対象は APIテスト先行、post-implementation unit / regression、edge-case inventory、AAA、flaky test avoidance、E2E artifact handling である。
- coverage は final validation lane が repo の `MINIMUM_COVERAGE = 70.0` に従って確認する。

## Checklist

- [tests-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/tests/references/checklists/tests-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。

## Agent が持つもの

- active contract: [tester.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/references/tester/contracts/tester.contract.json)
- permissions: [permissions.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/references/tester/permissions.json)

## Maintenance

- scenario / unit の知識差分は focused skill に置く。
- output obligation を skill 本体へ戻さない。
