# Test Patterns

## 目的

`implementation_tester` が single_handoff_packet 1 件を product test で証明するための判断パターンをまとめる。
agent TOML の tool policy と skill の出力規約は上書きしない。

## 採用する考え方

- 実装前 implementation_tester は、承認済み受け入れ条件、public seam、入力開始点、主要観測点、期待 outcome が固定済みの `APIテスト` だけで使う。
- unit test と原因未確定の regression test は実装後に追加または更新する。
- single_handoff_packet は完了条件、public seam、test target、validation command の順に読む。
- test は behavior を証明し、implementation detail を固定しない。
- null、empty、invalid、boundary、error path、concurrency をリスクに応じて含める。
- E2E は `UI人間操作E2E` として critical user flow と browser surface の証跡に絞る。
- flaky test は arbitrary wait ではなく、明確な condition wait へ直す。

## 適用ルール

- paid real AI API は test で呼ばない。fake provider、DI seam、test bootstrap を使う。
- owned_scope 外を直接読んで test scope を広げない。
- test_subscope が渡された場合は、その sub-scope の completion_signal clause、public seam、test target file、validation command だけを証明する。
- insufficient_context_criteria は structural gate とし、behavior_to_prove、public seam、test target、assertion focus、fixture/helper 方針、focused validation の不足時だけ insufficient_context を返す。
- test_subscope が completion_signal clause、public seam、test target file、validation command のいずれにも対応しない場合は insufficient_context を返す。
- not_insufficient_context: 承認済み scenario を元に期待どおり fail する test、局所的 import 修正、既存 test file 内の軽微な確認は停止理由にしない。
- 原因未確定の regression test を実装前に書く必要がある場合は停止し、post-implementation test として orchestrator へ返す。
- backend service / usecase / controller と frontend gateway / screen controller を主戦場にする。
- UI 操作証跡は `agent-browser` CLI を使う。
- product test runner としての Playwright は、既存 test が必要とする場合だけ使う。
- coverage と `python3 scripts/harness/run.py --suite all` は final validation lane へ defer する。
- AAA を守り、1 test は 1 behavior / branch / scenario outcome を証明する。

## 赤旗

- test name が `works` など曖昧で、何を証明したか分からない。
- test body に条件分岐がある。
- external AI / network / clock / random に依存している。
- coverage / harness all の未実行理由を final validation lane へ渡していない。
- assertion が弱く、壊れても pass する。
- listed files / symbols 外を探索して context を膨らませている。
- insufficient_context を返さず broad investigation を始めている。
- criteria mismatch になる不安や通常の局所確認を insufficient_context にしている。
