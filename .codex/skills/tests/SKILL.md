---
name: tests
description: Codex implementation lane 側の プロダクトテスト 共通作業プロトコル。承認済み owned_scope を test で証明する判断基準を提供する。
---
# Tests

## 目的

`tests` は作業プロトコルである。
`implementation_tester` agent が、single_handoff_packet と owned_scope を プロダクトテスト で証明する時の共通判断を提供する。

ツール権限 は [implementation_tester.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_tester.toml) が持ち、handoff は skill に従う。

## 対応ロール

- `implementation_tester` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- owner artifact は `tests` の出力規約で固定する。

## 入力規約

- 承認済み handoff または実装済み scope を プロダクトテスト で証明する時
- scenario artifact または unit responsibility を test に落とす時
- fake provider や DI seam で paid real AI API を避ける時
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。
- 必須入力: single_handoff_packet, approval_record, owned_scope, test_target, validation_commands
- 任意入力: test_subscope, knowledge_focus, scenario_artifact, implemented_scope, reproduction_evidence, existing_test_pointers
- input_notes: {"single_handoff_packet": "implementation-scope から抽出済みの handoff 1 件だけ。full implementation-scope、active work plan 全文、source artifacts、後続 handoff は入力に含めない。", "test_subscope": "implement_lane が context 枯渇時に同一 handoff 内で狭めた implementation_tester 用 sub-scope。完了条件 clause、public seam / API boundary、test target file、validation command のいずれか 1 軸で切られる。", "knowledge_focus": "tests-scenario または tests-unit の参照ヒント。共通規約と完了条件は変えない。"}
- insufficient_context_criteria: {"gate": "structural_gate", "return_insufficient_context_when": ["single_handoff_packet に behavior_to_prove、public seam、test target、assertion focus、fixture/helper 方針、focused validation のいずれかが欠けている", "`UI人間操作E2E` で開始操作、検証対象の入口、入力模倣方針のいずれかが欠けている", "`APIテスト` で public seam、request / response contract、入力開始点、主要観測点のいずれかが欠けている", "test_subscope が 完了条件 clause、public seam / API boundary、test target file、validation command のいずれにも対応していない", "test 作成に owned_scope 外探索、プロダクトコード 変更、paid API 呼び出しが必要になる"], "not_insufficient_context_when": ["承認済み scenario を元に期待どおり fail する プロダクトテスト を追加できる", "局所的な import 修正または既存 test file 内の軽微な確認だけで進められる", "single_handoff_packet の listed files / symbols 内で test seam を確認できる"], "required_when_true": ["reason", "needed_context", "tested_subscope", "remaining_test_subscopes"]}

## 外部参照規約

- エージェント実行定義とツール権限は [implementation_tester.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_tester.toml) の `allowed_write_paths` / `allowed_commands` とする。
- エージェント実行定義: [implementation_tester.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_tester.toml)
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/tests/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/tests-scenario/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/tests-unit/SKILL.md

## 内部参照規約

### 拘束観点

- test の責務分割
- single_handoff_packet の読み順
- insufficient_context の返し方
- deterministic setup
- final validation lane へ defer する coverage / harness all の扱い
- Arrange / Act / Assert の読みやすさ
- focused skill の選び方

- 参照 pattern は [test-patterns.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/tests/references/patterns/test-patterns.md) とする。

## 判断規約

- 各 test は 1 つの振る舞いを扱う
- single_handoff_packet の完了条件、public seam、test target、validation command の順に読む
- listed files / symbols 外を探索して context 不足を埋めない
- insufficient_context_criteria は structural gate とし、behavior_to_prove、public seam、test target、assertion focus、fixture/helper 方針、focused validation の不足時に返す
- `UI人間操作E2E` で開始操作、検証対象の入口、入力模倣方針が不足する場合は insufficient_context を返す
- `APIテスト` で public seam、request / response contract、入力開始点、主要観測点が不足する場合は insufficient_context を返す
- test_subscope が 完了条件 clause、public seam、test target file、validation command のいずれにも対応しない場合は insufficient_context を返す
- 承認済み scenario を元に期待どおり fail する test、局所的 import 修正、既存 test file 内の軽微な確認は not_insufficient_context として扱う
- 原因未確定の regression test は実装前に書かない
- setup は決定的にする
- test 追加または更新後、handoff を終える前に touched layer に対応する local validation を実行する
- paid real AI API を呼ばない
- 新しい要件解釈を足さない

- single_handoff_packet の 完了条件 clause、behavior_to_prove、public seam、assertion_focus に沿って test を作る
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
- active 規約 は agent 1:1。scenario / unit の差分は focused skill で扱い、output obligation はこの 規約 に固定する。

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。
- 返却先: implement_lane
- handoff 資料のスコープ粒度に対応する プロダクトテスト と必要最小限の fixture / helper だけを返す
- single_handoff_packet 1 件と test_subscope の観点と対応づけ、複数 handoff を束ねない
- 実際に test で証明した 完了条件 clause、public seam / API boundary、test target file、validation command を返す。test_subscope が入力された場合はそれに対応づける
- 同じ handoff 内で未証明の implementation_tester sub-scope を返す。完了条件は削らず、未処理分を明示する
- test 追加または更新後、handoff を終える前に touched layer に対応する local validation 結果を返す。backend は `python3 scripts/harness/run.py --suite backend-local`、frontend は `python3 scripts/harness/run.py --suite frontend-local`、mixed は touched layer に応じて両方を実行する。未実行なら blocked reason を返す
- insufficient_context_criteria の structural_gate に一致する場合だけ true とし、reason、needed_context、tested_subscope、remaining_test_subscopes を返す。自力で広く調査して埋めない。criteria に一致しない不安や通常の局所確認では true にしない。問題がなければ false または none
- insufficient_context true 時は insufficient_context_criteria のどの structural gate に一致したかを返す。false 時は none または未使用にする
- insufficient_context 時に不足している public seam、request / response contract、existing test target、fixture/helper、assertion focus、validation command、`UI人間操作E2E` の開始操作、検証対象の入口、入力模倣方針、`APIテスト` の入力開始点、主要観測点を列挙する
- final validation lane に defer する。実行した場合だけ結果を返し、未実行なら final_validation_deferred と理由を返す
- final validation lane に defer する。実行した場合だけ結果を返し、未実行なら final_validation_deferred と理由を返す
- 未証明の振る舞い、未実行 validation、blocked reason を分ける

## 完了規約

- 承認済み owned_scope 内の成果だけが返却されている。
- validation、未実行項目、residual risk が source_ref 付きで整理されている。
- implemented scope と owned_scope を確認した。
- single_handoff_packet の完了条件、public seam、test target、validation command を確認した。
- `UI人間操作E2E` では開始操作、検証対象の入口、入力模倣方針を確認した。
- `APIテスト` では public seam、request / response contract、入力開始点、主要観測点を確認した。
- test_subscope がある場合はその範囲だけを証明した。
- insufficient_context_criteria の structural gate に一致する場合だけ insufficient_context、needed_context、remaining_test_subscopes を返した。
- not_insufficient_context に該当する局所確認や承認済み scenario 由来の fail test を停止理由にしなかった。
- 原因未確定の regression test を実装前に書かなかった。
- deterministic setup にした。
- focused skill の知識だけを追加で参照した。
- 必須 evidence: single_handoff_packet id, handoff scope granularity, owned_scope, test_target, tested_subscope or insufficient_context reason, touched-layer local validation result or blocked reason
- 完了判断材料: implement_lane が implementation_implementer、review、reroute の次 action を touched-layer local validation 結果込みで判断できる
- 残留リスク: remaining_gaps

## 停止規約

- プロダクトコードの恒久修正が主目的の時
- design や scenario artifact を新規に作る時
- review だけを行う時
- insufficient_context を返さず広く調査しない
- UI 入口の `UI人間操作E2E` を裏側の直接呼び出しだけで代替しない
- criteria mismatch になる不安や通常の局所確認を insufficient_context にしない
- プロダクトコード を広く直さない
- docs、`.codex`、`.codex/skills`、`.codex/agents` を変更しない
- 停止時は不足項目、衝突箇所、戻し先を返す。
- 新しい要件解釈を足さなかった場合は停止する。
- insufficient_context を広い調査で埋めなかった場合は停止する。
- UI 入口の `UI人間操作E2E` を裏側の直接呼び出しだけで代替しなかった場合は停止する。
- criteria mismatch になる insufficient_context を返さなかった場合は停止する。
- paid real AI API を呼ばなかった場合は停止する。
- mode 別 個別 JSON 規約 を使わなかった場合は停止する。
- 拒否条件: missing single_handoff_packet
- 拒否条件: missing test_target
- 拒否条件: missing owned_scope
- 拒否条件: missing required context for UI人間操作E2E or APIテスト
- 拒否条件: paid real AI API risk
- 停止条件: 設計や ownership の整理が先に必要である
- 停止条件: プロダクトコードの広い変更が必要になる
- 停止条件: paid real AI API を呼ぶ危険がある
- 規約違反条件: insufficient_context を返さず広く調査する
- 規約違反条件: criteria mismatch: insufficient_context_criteria に一致しない理由で insufficient_context を返す
- 規約違反条件: 原因未確定の regression test を実装前に書く
- 規約違反条件: remaining_test_subscopes を隠して完了扱いにする
- 規約違反条件: test 追加または更新後に touched layer の local validation 結果または未実行理由を返さない
