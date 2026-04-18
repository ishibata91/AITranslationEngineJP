---
name: implementation-scope
description: Codex 側の実装スコープ知識 package。human review 後に、人間が Copilot へ渡せる handoff packet を owned_scope、依存、検証単位へ分ける判断基準を提供する。
---

# Implementation Scope

## 目的

`implementation-scope` は知識 package である。
`designer` agent が human review 後に、人間向け Copilot handoff packet を固定するための、分割粒度、依存、validation、completion signal の見方を提供する。

実行権限、agent contract、handoff、stop / reroute は [designer.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.agent.md) が持つ。

## いつ参照するか

- design bundle が human review 済みになった時
- 人間が Copilot に渡せる owned_scope を作る時
- handoff ごとの depends_on と validation を固定する時

## 参照しない場合

- human review 前に実装 scope を決める時
- Codex から Copilot へ直接 handoff する時
- product code を直接実装する時
- 実装時の再現、trace、review 補助を扱う時

## 知識範囲

- `implementation-scope.md` の構成
- owned_scope、depends_on、validation_commands、completion_signal
- 人間向け Copilot handoff packet の構成
- docs 正本化を handoff に混ぜない境界

## 原則

- human review 後にだけ作る
- 1 handoff は独立検証可能な粒度にする
- scope、依存、validation、done condition を必ず揃える
- Codex は Copilot へ直接渡さず、人間へ handoff packet を返す
- Copilot に docs 正本化や workflow 変更を渡さない

## 標準パターン

1. human review status と approval record を確認する。
2. source artifact を列挙する。
3. handoff を risk と validation 単位で分割する。
4. 各 handoff に owned_scope、depends_on、validation_commands を書く。
5. 人間が Copilot に渡す entry、禁止事項、期待される完了報告を明示する。
6. Copilot 修正完了後に正本化が必要なら `propose_plans` へ戻す前提を残す。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は `designer` agent contract に従う。

## DO / DON'T

DO:
- 承認済み artifact だけを source にする
- implementation handoff を小さく分ける
- validation command と completion signal を揃える
- 人間がそのまま Copilot に渡せる packet にする

DON'T:
- human review 前に owned_scope を確定しない
- Codex から Copilot へ直接 handoff しない
- docs 正本化を Copilot handoff に含めない
- implementation-time investigation を Codex 側へ戻さない

## Checklist

- [implementation-scope-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-scope/references/checklists/implementation-scope-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は `designer` agent contract が決める。

## References

- template: [implementation-scope.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/templates/task-folder/implementation-scope.md)
- Copilot entry: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-orchestrate/SKILL.md)
- agent spec: [designer.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.agent.md)
- agent contract: [designer.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/designer/contracts/designer.contract.json)

## Maintenance

- 権限、write scope、output obligation を skill 本体へ戻さない。
- Copilot 実装 workflow の詳細は [.github/skills](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/) に置く。
- handoff 粒度の長い例は references に分離する。
