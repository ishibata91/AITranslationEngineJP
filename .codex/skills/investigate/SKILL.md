---
name: investigate
description: Codex 側の設計前調査作業プロトコル。再現、UI 証跡、trace、リスク報告 を 根拠 first で扱う判断基準を提供する。
---
# Investigate

## 目的

`investigate` は作業プロトコルである。
`investigator` agent が設計前に必要な証拠を集めるための、観測事実、UI 証跡、仮説、残り 不足 の分け方を提供する。

UI check 専用 skill / agent は置かない。
設計前の UI 根拠 は `investigator` が `investigate` の一部として扱う。

## 対応ロール

- `investigator` が使う。
- 返却先は 呼び出し元 または次 agent とする。
- 担当成果物は `investigate` の出力規約で固定する。

## 入力規約

- 設計前に再現可否を確認する時
- UI 根拠、console、画面状態を設計判断の証跡として確認する時
- trace の観測点と不足情報を整理する時
- design continuation の リスク を短く返す時
- 入力に 根拠参照、担当者、承認状態が不足する場合は推測で補わない。
- 必須入力: 呼び出し元, investigation_goal, known_context
- 任意入力: investigation_mode, reproduction_steps, candidate_paths
- selector: {"investigation_mode": ["再現", "UI 根拠", "trace", "リスク報告"]}
- 必須 成果物: active task 文脈 or 呼び出し元提供 investigation 文脈

## 外部参照規約

- エージェント実行定義とツール権限は [investigator.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/investigator.toml) の 書き込み許可 / 実行許可 とする。
- エージェント実行定義: [investigator.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/investigator.toml)
- ツール権限: エージェント実行定義の 書き込み許可 / 実行許可 に従う
- Codex in-app browser の操作規約は [browser-use skill](/Users/iorishibata/.codex/plugins/cache/openai-bundled/browser-use/0.1.0-alpha1/skills/browser/SKILL.md) とする。
- `allowed_commands = []` は shell command 不使用を意味し、`browser-use` は Codex runtime のブラウザ操作として扱う。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/investigate/SKILL.md

## 内部参照規約

### 拘束観点

- `再現`、`UI 根拠`、`trace`、`リスク報告` の観点
- 観測済み事実、UI 根拠、仮説 の分離
- 根拠 path と再現条件の残し方
- 設計を止める 残留リスク の表現

## 判断規約

- 根拠 のない結論を書かない
- 観測事実と仮説を混ぜない
- 設計前の UI 根拠 は Codex in-app browser の `browser-use` で確認する
- UI 根拠 は画面状態、console、screenshot、操作条件を分けて残す
- `agent-browser` CLI は Codex implementation レーン の実装時調査でだけ使う
- 実装 レーン の調査は Codex implementation レーンへ戻す

- observed、UI 根拠、inferred を分ける
- 証跡 path と再現条件を優先する
- 設計継続可否に効く 不足 を残す
- active 規約 は agent に対して 1 ファイルだけ置く。調査種別は selector で扱う。

## 出力規約

- 出力は判断結果、根拠参照、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。

### Handoff

- 引き継ぎ先: `designer`
- 渡す対象範囲: 観測済み事実、仮説、残り 不足、残留 risks
- 調査 mode: 実施した調査の種類を返す。
- 観測事実: 観測済み事実だけを返す。
- UI 証跡: UI を確認した場合は証跡と参照先を返す。
- 仮説: 事実と分けて原因候補を返す。
- 観測点: 確認した入口、経路、対象を返す。
- 残り 不足: 未確認事項と理由を返す。
- 残留リスク: 設計判断に残る リスク を返す。
- 推奨 next step: 設計継続、追加調査、停止のどれが妥当かを返す。

## 完了規約

- 出力規約を満たし、次の 実行者 が再解釈なしで判断できる。
- 不足情報または停止理由がある場合は明示されている。
- 観測事実、UI 根拠、仮説、未観測 不足 を分けた。
- 根拠 path、再現条件、UI check 対象範囲 を残した。
- design continuation に必要な リスク を返した。
- 必須 根拠: 観測済み事実 根拠, UI 根拠 when mode is UI 根拠, reproduction condition, 根拠 path when used
- 完了判断材料: designer が設計継続か停止かを判断できる。
- 残留リスク: 設計判断に残る リスク が返っている。

## 停止規約

- implementation-scope 承認後の再現や再観測を扱う時
- 恒久修正や プロダクトテスト 追加が必要な時
- implementation レビュー が主目的の時
- 観測条件が不足する場合は停止する。
- 恒久修正が必要なら `designer` へ戻す。
- 実装時調査なら、Codex implementation レーン [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-investigate/SKILL.md) を使う前提で `designer` へ戻す。
- 恒久修正を始めない
- implementation-time investigation を扱わない
- 承認済み実装範囲 や対象 file を確定しない
- 停止時は不足項目、衝突箇所、戻し先を返す。
- 根拠 なしの結論を書く必要がある場合は停止する。
- UI check 専用 agent を前提にする場合は停止する。
- implementation-time investigation を扱う場合は停止する。
- 拒否条件: implementation-time investigation
- 拒否条件: permanent fix request
- 拒否条件: 根拠成果物 不足
