---
name: requirements-design
description: Codex 側の要件設計知識 package。product intent を capability、制約、不変条件、非目標、未決事項へ分ける判断基準を提供する。
---

# Requirements Design

## 目的

`requirements-design` は知識 package である。
`designer` agent が実装前に capability と制約を固定するための、要件分解と判断記録の見方を提供する。

実行権限、agent contract、handoff、stop / reroute は [designer.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.agent.md) が持つ。

## いつ参照するか

- human の依頼を capability と制約へ分ける時
- acceptance basis を scenario-design へ渡す時
- 既存正本と依頼の衝突を整理する時

## 参照しない場合

- UI visual system の判断が主目的の時
- product code の owned_scope を確定する時
- docs 正本を直接更新する時

## 知識範囲

- capability、in scope、out of scope
- business rule、invariant、data ownership、state transition
- decision point と open question
- acceptance basis の作り方

## 原則

- 確認済み事実、推測、未決事項を分ける
- user-visible promise と実装詳細を混ぜない
- 既存正本と衝突する場合は open question に戻す
- implementation scope へ踏み込まない

## 標準パターン

1. capability を対象者、可能になること、成果で 1 文にする。
2. 制約を business rule、scope boundary、invariant へ分ける。
3. 論点ごとに issue、options、recommendation、risk を整理する。
4. scenario-design が読める acceptance basis に圧縮する。
5. human 判断が必要な点を open question に残す。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は `designer` agent contract に従う。

## DO / DON'T

DO:
- product truth の根拠 path を残す
- non-goal を明示する
- acceptance に使える粒度で書く

DON'T:
- product truth を発明しない
- 実装対象 file や owned_scope を書かない
- 背景説明で未決事項を隠さない

## Checklist

- [requirements-design-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/requirements-design/references/checklists/requirements-design-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は `designer` agent contract が決める。

## References

- template: [requirements-design.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/templates/task-folder/requirements-design.md)
- agent spec: [designer.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.agent.md)
- agent contract: [designer.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/designer/contracts/designer.contract.json)

## Maintenance

- 権限、write scope、output obligation を skill 本体へ戻さない。
- implementation-scope の粒度判断をこの skill に戻さない。
- 長い判断表は references に分離する。
