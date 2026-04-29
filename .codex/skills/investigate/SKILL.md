---
name: investigate
description: Codex 側の設計前調査作業プロトコル。再現、UI 証跡、trace、risk-report を evidence first で扱う判断基準を提供する。
---
# Investigate

## 目的

`investigate` は作業プロトコルである。
`investigator` agent が設計前に必要な証拠を集めるための、観測事実、UI 証跡、仮説、remaining gap の分け方を提供する。

UI check 専用 skill / agent は置かない。
設計前の UI evidence は `investigator` が `investigate` の一部として扱う。

## 対応ロール

- `investigator` が使う。
- 返却先は caller または次 agent とする。
- owner artifact は `investigate` の出力規約で固定する。

## 入力規約

- 設計前に再現可否を確認する時
- UI evidence、console、画面状態を設計判断の証跡として確認する時
- trace の観測点と不足情報を整理する時
- design continuation の risk を短く返す時
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。
- 必須入力: caller, investigation_goal, known_context
- 任意入力: investigation_mode, reproduction_steps, candidate_paths
- selectors: {"investigation_mode": ["reproduce", "ui-evidence", "trace", "risk-report"]}
- 必須 artifact: active task context or caller-provided investigation context

## 外部参照規約

- エージェント実行定義とツール権限は [investigator.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/investigator.toml) の `allowed_write_paths` / `allowed_commands` とする。
- binding: [investigator.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/investigator.toml)
- エージェント実行定義: [investigator.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/investigator.toml)
- ツール権限: エージェント実行定義の `allowed_write_paths` / `allowed_commands` に従う
- binding: [investigator.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/investigator.toml)
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/investigate/SKILL.md

## 内部参照規約

### 拘束観点

- `reproduce`、`ui-evidence`、`trace`、`risk-report` の観点
- observed fact、UI evidence、hypothesis の分離
- evidence path と再現条件の残し方
- 設計を止める residual risk の表現

## 判断規約

- evidence のない結論を書かない
- 観測事実と仮説を混ぜない
- UI evidence は画面状態、console、screenshot、操作条件を分けて残す
- 実装 lane の調査は Codex implementation laneへ戻す

- observed、UI evidence、inferred を分ける
- 証跡 path と再現条件を優先する
- 設計継続可否に効く gap を残す
- active 規約 は agent に対して 1 ファイルだけ置く。調査種別は selector で扱う。

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。

### Handoff

- handoff 先: `designer`
- 渡す scope: observed facts、hypotheses、remaining gaps、residual risks
- 調査 mode: 実施した調査の種類を返す。
- 観測事実: 観測済み事実だけを返す。
- UI 証跡: UI を確認した場合は証跡と参照先を返す。
- 仮説: 事実と分けて原因候補を返す。
- 観測点: 確認した入口、経路、対象を返す。
- 残り gap: 未確認事項と理由を返す。
- 残留リスク: 設計判断に残る risk を返す。
- 推奨 next step: 設計継続、追加調査、停止のどれが妥当かを返す。

## 完了規約

- 出力規約を満たし、次の actor が再解釈なしで判断できる。
- 不足情報または停止理由がある場合は明示されている。
- 観測事実、UI evidence、仮説、未観測 gap を分けた。
- evidence path、再現条件、UI check scope を残した。
- design continuation に必要な risk を返した。
- 必須 evidence: observed fact evidence, UI evidence when mode is ui-evidence, reproduction condition, source path when used
- 完了判断材料: designer が設計継続か停止かを判断できる。
- 残留リスク: 設計判断に残る risk が返っている。

## 停止規約

- implementation-scope 承認後の再現や再観測を扱う時
- 恒久修正や プロダクトテスト 追加が必要な時
- implementation review が主目的の時
- 観測条件が不足する場合は停止する。
- 恒久修正が必要なら `designer` へ戻す。
- 実装時調査なら、Codex implementation lane [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-investigate/SKILL.md) を使う前提で `designer` へ戻す。
- 恒久修正を始めない
- implementation-time investigation を扱わない
- owned_scope や対象 file を確定しない
- 停止時は不足項目、衝突箇所、戻し先を返す。
- evidence なしの結論を書かなかった場合は停止する。
- UI check 専用 agent 前提にしなかった場合は停止する。
- implementation-time investigation を扱わなかった場合は停止する。
- 拒否条件: implementation-time investigation
- 拒否条件: permanent fix request
- 拒否条件: source artifact missing
