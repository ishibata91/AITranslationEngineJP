# Test Patterns

## 目的

`tester` が single_handoff_packet 1 件と tester_context_packet を product test で証明するための判断パターンをまとめる。
agent contract の権限や output obligation は上書きしない。

## 採用する考え方

- Red / Green / Refactor の考え方を、single_handoff_packet、tester_context_packet、owned_scope の範囲で使う。
- tester_context_packet は test_ingredients、test_required_reading、requirements_policy_decisions の test impact、test_validation_entry の順に読む。
- test は behavior を証明し、implementation detail を固定しない。
- null、empty、invalid、boundary、error path、concurrency をリスクに応じて含める。
- E2E は critical user flow と browser surface の証跡に絞る。
- flaky test は arbitrary wait ではなく、明確な condition wait へ直す。

## 適用ルール

- paid real AI API は test で呼ばない。fake provider、DI seam、test bootstrap を使う。
- full lane_context_packet、fix_ingredients、change_targets、broad related_code_pointers を直接読んで test scope を広げない。
- test_subscope が渡された場合は、その sub-scope の completion_signal clause、public seam、test target file、validation command だけを証明する。
- insufficient_context_criteria は structural gate とし、behavior_to_prove、public seam、test target、assertion focus、fixture/helper 方針、focused validation の不足時だけ insufficient_context を返す。
- test_subscope が completion_signal clause、public seam、test target file、validation command のいずれにも対応しない場合は insufficient_context を返す。
- not_insufficient_context: 期待どおり fail する test、局所的 import 修正、既存 test file 内の軽微な確認は停止理由にしない。
- backend service / usecase / controller と frontend gateway / screen controller を主戦場にする。
- Playwright は必要最小限の UI evidence に使い、mock 不能な Wails binding は専用 bootstrap を前提にする。
- coverage は `python3 scripts/harness/run.py --suite coverage` で Sonar-compatible coverage 70% 以上を確認する。
- closeout では `python3 scripts/harness/run.py --suite all` の evidence を残す。
- AAA を守り、1 test は 1 behavior / branch / scenario outcome を証明する。

## 赤旗

- test name が `works` など曖昧で、何を証明したか分からない。
- test body に条件分岐がある。
- external AI / network / clock / random に依存している。
- coverage 70% 未満を test gap として扱っていない。
- harness coverage / all の失敗または未実行理由がない。
- assertion が弱く、壊れても pass する。
- listed files / symbols 外を探索して context を膨らませている。
- insufficient_context を返さず broad investigation を始めている。
- criteria mismatch になる不安や通常の局所確認を insufficient_context にしている。
