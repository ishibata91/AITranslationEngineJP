---
name: implement
description: Codex implementation lane 側の product code 実装の共通作業プロトコル。承認済み owned_scope を実装する判断基準を提供する。
---
# Implement

## 目的

`implement` は作業プロトコルである。
`implementation_implementer` agent が、承認済み `implementation-scope` の handoff 1 件を owned_scope 内へ実装する時の共通判断を提供する。

tool policy は [implementation_implementer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_implementer.toml) が持ち、handoff は skill に従う。

## 対応ロール

- `implementation_implementer` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- owner artifact は `implement` の出力規約で固定する。

## 入力規約

- owned_scope 内の product code を実装する時
- lane_context_packet に基づいて product code を実装する時
- `APIテスト` 先行時の implementation_tester output を product code 実装へ反映する時
- lane-local validation の扱いを確認する時
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。
- 必須入力: single_handoff_packet, lane_context_packet, approval_record, implementation_target, owned_scope, depends_on_resolved
- 任意入力: implementation_subscope, knowledge_focus, lane_local_validation_commands, implementation_tester_output
- input_notes: {"single_handoff_packet": "implementation-scope から抽出済みの handoff 1 件だけ。full implementation-scope、active work plan 全文、source artifacts、後続 handoff は入力に含めない。", "lane_context_packet": "implementation_distiller が single_handoff_packet 1 件だけから作った実装前整理。fix_ingredients、distracting_context、first_action、change_targets、requirements_policy_decisions、required_reading、symbol / line number 付き related_code_pointers を含む。full implementation-scope、active work plan 全文、source artifacts、後続 handoff 由来の情報を含めない。", "implementation_subscope": "implement_lane が context 枯渇時に同一 handoff 内で狭めた implementation_implementer 用 sub-scope。completion_signal clause、public seam / API boundary、change target / symbol、validation command のいずれか 1 軸で切られる。completion_signal を削るものではない。", "implementation_tester_output": "`APIテスト` 先行 handoff で implementation_tester が先に返した product test result。通常、unit、原因未確定の regression handoff では入力に含めない。", "knowledge_focus": "implement-backend、implement-frontend、implement-mixed、implement-fix-lane の参照ヒント。共通規約と完了条件は変えない。implement-mixed は API / Wails / DTO / gateway など接合点 scope に限定する。"}
- insufficient_context_criteria: {"gate": "structural_gate", "return_insufficient_context_when": ["lane_context_packet に fix_ingredients、first_action、change_targets、requirements_policy_decisions、existing pattern、validation_entry のいずれかが欠けている", "first_action が 1 completion_signal clause に固定されていない、line / symbol / public seam が不明、または必要な closure chain がない", "実装に owned_scope 拡張、product test / fixture / snapshot / test helper 変更、docs / workflow 変更、新規設計判断、broad refactor が必要になる"], "not_insufficient_context_when": ["listed required_reading 内の局所確認だけで first edit に入れる", "既存 pattern への通常追従で実装できる", "lane-local validation failure を product code 内の scope で修正できる"], "required_when_true": ["reason", "needed_context", "suggested_narrowing_axis", "implemented_subscope", "remaining_implementation_subscopes"]}

## 外部参照規約

- agent runtime と tool policy は [implementation_implementer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_implementer.toml) の `allowed_write_paths` / `allowed_commands` とする。
- agent runtime: [implementation_implementer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_implementer.toml)
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-backend/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-frontend/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-mixed/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-fix-lane/SKILL.md

## 内部参照規約

### 拘束観点

- owned_scope を超えない実装判断
- handoff 資料のスコープ粒度に合わせる判断
- coding guidelines と既存 pattern の確認
- lint policy と architecture constraint の局所確認
- boundary、error path、test surface の実装品質判断
- validation result と residual risk の返し方
- focused skill の選び方

- 参照 pattern は [implementation-quality-patterns.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement/references/patterns/implementation-quality-patterns.md) とする。

## 判断規約

- `implementation-scope` と owned_scope を超えない
- handoff 資料のスコープ粒度で実装する
- lane_context_packet に合わせて product code だけを変更する
- `APIテスト` 先行時だけ implementation_tester output も確認する
- implementation_subscope が渡された場合はその sub-scope 内だけを実装する
- 実装完了後、handoff を終える前に touched layer に対応する local validation を実行する
- fix_ingredients に対応する code path を優先し、distracting_context へ寄り道しない
- first_action と change_targets から着手する
- insufficient_context_criteria は structural gate とし、fix_ingredients、first_action、change_targets、requirements_policy_decisions、existing pattern、validation_entry の不足時に返す
- first_action が 1 clause に固定されていない、line / symbol / public seam が不明、closure chain がない場合は insufficient_context を返す
- listed required_reading 内の局所確認、既存 pattern への通常追従、lane-local validation failure は not_insufficient_context として扱う
- 既存 pattern、naming、layer に合わせる
- broad refactor を混ぜない
- product test、fixture、snapshot、test helper は implementation_tester が扱う
- docs 正本化をしない

- 実装前に [coding-guidelines.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/coding-guidelines.md) を読む
- 実装前に [lint-policy.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/lint-policy.md) を読み、handoff に効く静的 check の責務を確認する
- lane_context_packet の fix_ingredients、distracting_context、first_action、change_targets、requirements_policy_decisions、related_code_pointers を確認する
- requirements_policy_decisions に architecture constraint がある場合は、その範囲だけ [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) を局所確認する
- implementation_subscope があれば completion_signal clause、public seam、change target / symbol、validation command を確認する
- insufficient_context を返す場合は reason、needed_context、suggested_narrowing_axis、remaining_implementation_subscopes を structural gate に対応づける
- entry point、call site、data flow、error path、test surface を確認する
- 既存 pattern に naming、constructor、DI、error return を合わせる
- generated import、layer 依存、boundary rule、format 逸脱など、touched layer で踏みやすい lint 観点を先に確認する
- lane-local validation 結果または未実行理由を返す
- backend handoff は `python3 scripts/harness/run.py --suite backend-local`、frontend handoff は `python3 scripts/harness/run.py --suite frontend-local` を使う
- mixed handoff は touched layer に応じて両方を実行する
- touched files は product code だけにする
- active 規約 は agent 1:1。backend / frontend / mixed / fix-lane の差分は focused skill で扱い、output obligation はこの 規約 に固定する。implementation_implementer は distiller 後の product code 実装を扱い、`APIテスト` 先行時だけ implementation_tester output を受け取る。product test / fixture / snapshot / test helper は変更しない。

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力に tool policy、agent runtime、product code の変更義務を含めない。
- 返却先: implement_lane
- 必須出力: touched_files, implemented_scope, implemented_subscope, remaining_implementation_subscopes, validation_results, closeout_notes, insufficient_context, reason, needed_context, suggested_narrowing_axis, residual_risks
- 出力 field 要件: {"touched_files": "single_handoff_packet の owned_scope に対応する product code だけを返す。product test、fixture、snapshot、test helper は含めない", "implemented_scope": "single_handoff_packet 1 件、lane_context_packet、implementation_subscope に対応づけ、複数 handoff を束ねない", "implemented_subscope": "実際に実装した completion_signal clause、public seam / API boundary、change target / symbol、validation command を返す。implementation_subscope が入力された場合はそれに対応づける", "remaining_implementation_subscopes": "同じ handoff 内で未実装の sub-scope を返す。completion_signal は削らず、未処理分を明示する", "validation_results": "実装完了後、handoff を終える前に touched layer に対応する local validation 結果を返す。backend は `python3 scripts/harness/run.py --suite backend-local`、frontend は `python3 scripts/harness/run.py --suite frontend-local`、mixed は touched layer に応じて両方を実行する。未実行なら blocked reason を返す。coverage、Sonar、harness all は implementation_implementer の必須 closeout にしない", "closeout_notes": "entry point、call site、data flow、error path、test surface、既存 pattern への整合を簡潔に返す。mixed の場合は接合点 contract を明記する", "insufficient_context": "insufficient_context_criteria の structural_gate に一致する場合だけ true とし、reason、needed_context、suggested_narrowing_axis、implemented_subscope、remaining_implementation_subscopes を返す。自力で広く調査して埋めない。criteria に一致しない不安、通常の局所確認、lane-local validation failure だけでは true にしない。問題がなければ false または none", "reason": "insufficient_context true 時は insufficient_context_criteria のどの structural gate に一致したかを返す。false 時は none または未使用にする", "needed_context": "insufficient_context 時に不足している fix_ingredients、first_action、change target、public seam、existing pattern、validation command を列挙する", "suggested_narrowing_axis": "insufficient_context 時に orchestrator が次に狭めるべき軸を completion_signal clause、public seam / API boundary、change target / symbol、validation command のいずれかで返す", "residual_risks": "未実行 validation、scope 超過、設計不足、test / fixture 変更が必要になった場合の blocked reason を分ける"}

## 完了規約

- 承認済み owned_scope 内の成果だけが返却されている。
- validation、未実行項目、residual risk が source_ref 付きで整理されている。
- owned_scope と implementation target を確認した。
- single_handoff_packet と lane_context_packet を確認した。
- `APIテスト` 先行時だけ implementation_tester output を確認した。
- implementation_subscope がある場合はその範囲だけを実装した。
- fix_ingredients と distracting_context を確認した。
- insufficient_context_criteria の structural gate に一致する場合だけ insufficient_context、needed_context、suggested_narrowing_axis を返した。
- not_insufficient_context に該当する局所確認、既存 pattern 追従、lane-local validation failure を停止理由にしなかった。
- first_action と change_targets から着手した。
- coding guidelines、lint policy、lane-local validation commands を確認した。
- requirements_policy_decisions にある architecture constraint を局所確認した。
- focused skill の知識だけを追加で参照した。
- touched files が product code だけであることを確認した。
- 必須 evidence: single_handoff_packet id, lane_context_packet id, fix_ingredients, distracting_context, first_action, change_targets, requirements_policy_decisions, owned_scope, implementation_tester_output when API test pre-implementation test exists, implemented_subscope or insufficient_context reason, entry point, call site, data flow or boundary, error path, test surface, touched-layer local validation result or blocked reason
- completion signal: implement_lane が review へ進める product code 実装結果と touched-layer local validation 結果が返っている
- residual risk key: residual_risks

## 停止規約

- 実装前 context 整理だけを行う時
- UI check や implementation review を行う時
- docs や workflow 文書を変更する時
- 要件や設計を追加しない
- fix_ingredients がないまま実装を始めない
- insufficient_context を返さず広い調査で不足 context を埋めない
- criteria mismatch になる不安や通常の局所確認を insufficient_context にしない
- implementation_subscope 外へ実装を広げない
- distracting_context を実装対象に混ぜない
- first_action がないまま広い調査を始めない
- lint を知らないまま実装して local validation で初めて境界違反を知る進め方をしない
- config、lint、test、coverage 設定を変更して gate を回避しない
- product test、fixture、snapshot、test helper を変更しない
- coverage、harness all、repo-local Sonar issue gate を implementation_implementer の必須 closeout にしない
- owned_scope 外の cleanup、rename、format を混ぜない
- docs、`.codex`、`.codex/skills`、`.codex/agents` を変更しない
- 停止時は不足項目、衝突箇所、reroute 先を返す。
- broad refactor を混ぜなかった場合は停止する。
- insufficient_context を広い調査で埋めなかった場合は停止する。
- criteria mismatch になる insufficient_context を返さなかった場合は停止する。
- implementation_subscope 外へ実装を広げなかった場合は停止する。
- distracting_context を実装対象に混ぜなかった場合は停止する。
- first_action 不足を広い調査で埋めなかった場合は停止する。
- product test、fixture、snapshot、test helper を変更しなかった場合は停止する。
- docs / workflow 文書を変更しなかった場合は停止する。
- mode 別 個別 JSON 規約 を使わなかった場合は停止する。
- 拒否条件: missing single_handoff_packet
- 拒否条件: missing lane_context_packet
- 拒否条件: missing fix_ingredients
- 拒否条件: missing distracting_context
- 拒否条件: missing first_action
- 拒否条件: missing change_targets
- 拒否条件: missing requirements_policy_decisions
- 拒否条件: missing approval_record
- 拒否条件: missing implementation_tester_output for API test pre-implementation handoff
- 拒否条件: unclear owned_scope
- 停止条件: 設計判断が不足している
- 停止条件: docs または workflow の変更が必要になる
- 停止条件: broad refactor なしでは実装できない
- 停止条件: product test、fixture、snapshot、test helper の変更が必要になる
- 規約違反条件: touched_test_files を返す
- 規約違反条件: product test、fixture、snapshot、test helper を変更する
- 規約違反条件: full implementation-scope または後続 handoff を入力として要求する
- 規約違反条件: insufficient_context を返さず広く調査して不足 context を埋める
- 規約違反条件: criteria mismatch: insufficient_context_criteria に一致しない理由で insufficient_context を返す
- 規約違反条件: remaining_implementation_subscopes を隠して完了扱いにする
- 規約違反条件: 実装完了後に touched layer の local validation 結果または未実行理由を返さない
