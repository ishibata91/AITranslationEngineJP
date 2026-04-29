---
name: implementation-distill
description: Codex implementation lane 側の実装前文脈整理の共通作業プロトコル。single_handoff_packet 1 件から実装に必要な facts を圧縮する判断基準を提供する。
---
# Implementation Distill

## 目的

`implementation-distill` は作業プロトコルである。
`implementation_distiller` agent が、`single_handoff_packet` 1 件から lane_context_packet を作る時の共通判断を提供する。

ツール権限 は [implementation_distiller.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_distiller.toml) が持ち、handoff は skill に従う。

## 対応ロール

- `implementation_distiller` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- owner artifact は `implementation-distill` の出力規約で固定する。

## 入力規約

- 実装前に facts、constraints、gaps を圧縮する時
- required reading と code pointer を owned_scope に絞る時
- validation entry を実装前に明示する時
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。
- 必須入力: single_handoff_packet, approval_record, owned_scope, depends_on_resolved, validation_commands
- 任意入力: knowledge_focus, reproduction_evidence, trace_or_analysis_result
- input_notes: {"single_handoff_packet": "implementation-scope から抽出済みの handoff 1 件だけ。full implementation-scope、active work plan 全文、source artifacts、後続 handoff は入力に含めない。", "knowledge_focus": "implementation-distill-implement、implementation-distill-fix、implementation-distill-refactor の参照ヒント。共通規約と完了条件は変えない。"}

## 外部参照規約

- エージェント実行定義とツール権限は [implementation_distiller.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_distiller.toml) の `allowed_write_paths` / `allowed_commands` とする。
- エージェント実行定義: [implementation_distiller.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_distiller.toml)
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-distill/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-distill-implement/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-distill-fix/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-distill-refactor/SKILL.md

## 内部参照規約

### 拘束観点

- path catalog から必要 file だけを summary / full に上げる圧縮
- fix_ingredients と distracting_context の分離
- first_action、change_targets、related_code_pointers の具体化
- 要件、実装方針、決定事項の要約
- facts、inferred、gap の分離
- single_handoff_packet と owned_scope の対応づけ
- focused skill の選び方

- 参照 pattern は [context-compression-patterns.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-distill/references/patterns/context-compression-patterns.md) とする。

## 判断規約

- `single_handoff_packet` 1 件を唯一の source scope にする
- owned_scope に関係する code pointer を優先する
- patch 生成に必要な fix ingredients を file / function / block の構造単位で残す
- 類似していても修正に不要な context は distracting_context として分ける
- repository method、interface、field の追加が必要そうに見えても、実 code で present / absent を確認するまで fact にしない
- first_action は 1 完了条件 clause に限定し、partial や複数 clause にしない
- 1 edit で clause が閉じない場合は、同じ clause の最小 closure chain を上流 symbol から leaf まで残す
- existing_patterns が none なら、探索範囲と実装判断への影響を残す
- validation entry は最初に試せる cheap check を優先する
- 実 code を読んでから first_action を返す
- 実装開始点は path、symbol、line number、変更種別で返す
- 要件、実装方針、決定事項は distiller が要約し、implementation_implementer の再読を原則不要にする
- 実装案を増やさず、実装に必要な制約だけを残す
- 設計不足は実装せず戻す

- fix_ingredients に path、symbol/type/function、line number、why_needed_for_patch を残す
- distracting_context に why_excluded と risk_if_read を残す
- 実装者が最初に触る file、symbol、line number、変更種別を残す
- requirements_policy_decisions に要件、実装方針、決定事項、out of scope、禁止事項を要約する
- repository method、interface、field の有無を present / absent の code fact として残す
- first_action が閉じる clause を 1 つだけ明示する
- existing_patterns がない場合は searched scope と impact を添えて none とする
- required_reading は読む目的と symbol を添えて順序づける
- 要件や決定事項の原文は、要約では判断できない時だけ required_reading に残す
- related_code_pointers は path、symbol/type/function、line number、読み取った事実を残す
- validation entry を明示する
- gap と residual risk を分ける
- active 規約 は agent 1:1。implement / fix / refactor の差分は focused skill で扱い、output obligation はこの 規約 に固定する。

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。
- 返却先: implement_lane
- single_handoff_packet 1 件だけから作る。implementation_implementer に渡せる fix_ingredients、distracting_context、first_action、change_targets、requirements_policy_decisions、required_reading、related_code_pointers、existing_patterns、facts、constraints、gaps、validation_entry を含める。implementation_tester 向けには full lane_context_packet ではなく implementation_tester_context_packet を別に返す
- single_handoff_packet 1 件だけから作る implementation_tester 専用 context。test_ingredients、test_required_reading、test_existing_patterns、test_distracting_context、test_validation_entry、requirements_policy_decisions の test impact を含める。fix_ingredients、change_targets、broad related_code_pointers をそのまま渡さない
- patch 生成に必要な最小十分 context を path、symbol/type/function、line number、structural_unit、why_needed_for_patch で返す。first_action と change_targets は fix_ingredients から導く
- 類似しているが今回の修正に使わない context を path、symbol/type/function、line number、why_excluded、risk_if_read で返す。なければ none と探索範囲を返す
- implementation_implementer が最初に行う 1 手を返す。file path、symbol name、line number、変更種別、対応する fix_ingredients、なぜそこから始めるかを必ず含める。完了条件 clause は 1 件だけに固定し、partial、複数 clause、曖昧な advance 表現は禁止する。1 edit で clause を閉じられない場合は、同じ clause の最小 clause_closure_chain を change_targets に明記し、first_action がその最上流 symbol であることを示す
- 変更候補を path、symbol/type/function、line number、structural_unit、change_kind、reason、対応する fix_ingredients、clause_link で返す。first_action が 1 edit で clause を閉じられない場合は、同じ clause の最小 clause_closure_chain を上流から leaf まで列挙する。実装対象が未特定なら gaps に blocker として返し、空のまま実装可能に見せない
- 要件、実装方針、決定事項、out of scope、禁止事項を distiller が要約して返す。各項目は source、summary、implementation_implementer impact を含める。implementation_implementer が原文を再読しなくても実装判断できる粒度にする
- 関連 code pointer は path、symbol/type/function、line number、structural_unit、読み取った事実、fix_ingredients との関係を含める。ファイル名だけの列挙は禁止。存在しない method / field / contract を新規追加候補にする場合は、interface または type 定義を読んだ事実として present / absent を返し、推測を fact にしない
- 踏襲する既存実装を path、symbol/type/function、line number、structural_unit、踏襲する理由で返す。見つからない場合は none と探索範囲、探索した layer、none が実装判断に与える影響を返す。owned_scope 内の同 layer pattern を読まずに none としない
- facts、inferred、gap を混ぜずに書く
- 実装者が読む順番を返す。各項目は path、symbol/type/function、line number、structural_unit、読む目的を含める。最初の項目は first_action と対応する実 code にする。要件、実装方針、決定事項の文書は requirements_policy_decisions に要約し、原文確認が必要な場合だけ required_reading に残す
- 最初に実行する cheap validation command と完了条件を返す。validation_commands から lane-local に最も狭い command を優先し、広い command しか返せない場合は narrower command を見つけられなかった理由と final validation との役割差を明記する
- implementation_tester が証明すべき最小 context を 完了条件 clause、behavior_to_prove、public seam/API boundary、request / response contract、入力開始点、主要観測点、`UI人間操作E2E` の開始操作、入力模倣方針、existing test target、fixture/helper、assertion_focus、validation command で返す。implementation 内部の fix_ingredients は test seam 特定に必要な場合だけ参照元として含める
- implementation_tester が読む最小 file / symbol / line を返す。既存 プロダクトテスト、public contract、gateway/controller/usecase seam を優先し、implementation 内部の broad code path を読ませない
- 踏襲する既存 test pattern を path、test name / helper、line number、踏襲する理由で返す。見つからない場合は none と探索範囲、test 実装への影響を返す
- test 作成に不要な implementation detail、similar but unrelated test、後続 handoff context を why_excluded、risk_if_read 付きで返す
- implementation_tester が最初に実行する focused validation command と完了条件を返す。broad coverage、harness all、final validation は final validation lane へ defer する理由を書く

## 完了規約

- 承認済み owned_scope 内の成果だけが返却されている。
- validation、未実行項目、residual risk が source_ref 付きで整理されている。
- facts、inferred、gap を分けた。
- single_handoff_packet 1 件だけを source scope にした。
- owned_scope の実 code を読んだ。
- fix_ingredients に path、symbol/type/function、line number、why_needed_for_patch を入れた。
- distracting_context に why_excluded と risk_if_read を入れた。
- repository method、interface、field の present / absent を実 code fact として確認した。
- first_action に path、symbol、line number、変更種別、1 つの clause_closed を入れた。
- first_action が 1 edit で clause を閉じない場合は、同じ clause の最小 closure chain を change_targets に入れた。
- change_targets と related_code_pointers に symbol/type/function と line number を入れた。
- existing_patterns または none の探索範囲と影響を入れた。
- validation_entry は cheap check を優先した。
- requirements_policy_decisions に要件、実装方針、決定事項を要約した。
- required reading と code pointer を owned_scope に絞った。
- focused skill の知識だけを追加で参照した。
- 必須 evidence: single_handoff_packet, owned_scope, validation_commands, requirements_policy_decisions summary, actual code file read, fix_ingredients with file path, symbol, and line number, distracting_context or none with searched scope, first_action with file path and symbol, first_action tied to exactly one 完了条件 clause, change target with line number, existing_patterns or none with searched scope and impact, cheap validation entry or reason broader command is unavoidable, implementation_tester_context_packet with test_ingredients, test_required_reading, and test_validation_entry
- 完了判断材料: implement_lane が implement、tests、review の次 handoff を判断できる
- 残留リスク: gaps

## 停止規約

- プロダクトコードまたはプロダクトテスト を変更する時
- review や UI check を行う時
- design 不足を実装側で補う時
- fix_ingredients を特定せず first_action だけを返さない
- 類似 context を required_reading に混ぜない
- 存在確認していない repository method、interface、field を追加前提にしない
- first_action に partial、複数 clause、曖昧な advance を書かない
- existing_patterns の none を探索範囲なしで返さない
- cheap validation を検討せず広い command だけを返さない
- 実 code を読まず handoff の文章を言い換えない
- 要件、実装方針、決定事項を要約せず required_reading に丸投げしない
- required_reading をファイル名だけの列挙にしない
- implementation_implementer に「どこから調べるか」を委ねない
- プロダクトコード / プロダクトテスト を変更しない
- 要件や設計を追加しない
- full implementation-scope、active work plan 全文、source artifacts、後続 handoff を要求しない
- 停止時は不足項目、衝突箇所、戻し先を返す。
- プロダクトコード / プロダクトテスト を変更しなかった場合は停止する。
- 類似 context を required_reading に混ぜなかった場合は停止する。
- 推測だけで repository method、interface、field の追加を前提にしなかった場合は停止する。
- first_action に partial、複数 clause、曖昧な advance を書かなかった場合は停止する。
- existing_patterns の none を探索範囲なしで返さなかった場合は停止する。
- cheap validation を検討せず広い command だけを返さなかった場合は停止する。
- handoff 文章の言い換えだけで返さなかった場合は停止する。
- 要件、実装方針、決定事項を required_reading に丸投げしなかった場合は停止する。
- required_reading をファイル名だけの列挙にしなかった場合は停止する。
- 要件や設計を追加しなかった場合は停止する。
- full implementation-scope、active work plan 全文、source artifacts、後続 handoff を要求しなかった場合は停止する。
- mode 別 個別 JSON 規約 を使わなかった場合は停止する。
- 拒否条件: missing single_handoff_packet
- 拒否条件: missing approval_record
- 拒否条件: unclear owned_scope
- 停止条件: 設計判断が不足している
- 停止条件: 変更が docs や workflow へ広がる
- 規約違反条件: full implementation-scope を入力として要求する
- 規約違反条件: active work plan 全文を入力として要求する
- 規約違反条件: source_artifacts または後続 handoff を入力として要求する
- 規約違反条件: 実 code を読まず handoff 文面の言い換えだけを返す
- 規約違反条件: fix_ingredients がない、または path / symbol / line number / why_needed_for_patch がない
- 規約違反条件: distracting_context を確認せず類似 context を required_reading に混ぜる
- 規約違反条件: implementation_tester_context_packet がない、または test_ingredients / test_required_reading / test_validation_entry がない
- 規約違反条件: implementation_tester_context_packet に fix_ingredients、change_targets、broad related_code_pointers を丸ごと混ぜる
- 規約違反条件: 要件、実装方針、決定事項を要約せず required_reading に丸投げする
- 規約違反条件: required_reading がファイル名の列挙だけである
- 規約違反条件: related_code_pointers に symbol/type/function または line number がない
- 規約違反条件: 存在確認していない repository method / interface / field を fact として扱う
- 規約違反条件: first_action がなく implementation_implementer に調査開始を委ねる
- 規約違反条件: first_action の clause_closed が partial、複数 clause、または曖昧な advance 表現である
- 規約違反条件: first_action が 1 clause の最小 closure chain を示さない
- 規約違反条件: existing_patterns が none なのに探索範囲と実装判断への影響がない
- 規約違反条件: validation_entry が広い command だけで cheap check の探索理由がない
- 規約違反条件: change_targets がなく実装可能に見せる
