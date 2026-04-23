---
name: updating-docs
description: Codex 側の docs 正本化知識 package。Copilot 修正完了後に、human 承認済み docs-only artifact を正本へ反映する判断基準を提供する。
---

# Updating Docs

## 目的

`updating-docs` は知識 package である。
`docs_updater` agent が Copilot 修正完了後に human 承認済み artifact を docs 正本へ反映するための、source of truth、承認確認、validation の見方を提供する。

人間可読な実行境界、handoff、stop / reroute はこの skill を正本にする。

## いつ参照するか

- Copilot の修正完了が分かっている時
- human 承認済み docs-only artifact を docs 正本へ移す時
- canonicalization target と validation を整理する時
- task-local artifact と docs source of truth の対応を確認する時

## 参照しない場合

- Copilot の修正完了が未確認の時
- workflow contract や skill / agent を変更する時
- product code や product test の変更が必要な時
- human approval が不足している時

## 知識範囲

- Copilot completion report の確認
- docs source of truth の選び方
- human approval record の確認
- approved artifact と canonical target の対応
- validation と remaining gaps の記録

## 原則

- Copilot 修正完了後にだけ正本化へ進む
- human 承認済み artifact だけを反映する
- docs-only scope を超えない
- implementation-scope を docs 正本へ自動昇格しない
- 未確定仕様を独断で補完しない

## Runtime Boundary

- binding: [docs_updater.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/docs_updater.toml)
- permissions: [permissions.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/docs_updater/permissions.json)
- contract: [docs_updater.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/docs_updater/contracts/docs_updater.contract.json)
- allowed: approved docs-only scope の docs 更新
- forbidden: product code、product test、workflow contract の変更
- write scope: `docs/` の承認済み正本だけ

## 標準パターン

1. Copilot completion report を確認する。
2. approval record と docs-only scope を確認する。
3. `docs/index.md` から canonical target を選ぶ。
4. approved artifact と target の差分だけを反映する。
5. validation を実行または記録する。
6. remaining gaps と reroute reason を返す。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は `docs_updater` agent contract に従う。

## Stop / Reroute

- Copilot の修正完了が分からない場合は停止する。
- approval がない場合は停止する。
- workflow 変更なら `propose_plans` へ戻す。
- product 実装が必要なら `propose_plans` へ戻す。

## Handoff

- handoff 先: `propose_plans`
- 渡す contract: [docs_updater.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/docs_updater/contracts/docs_updater.contract.json)
- 渡す scope: docs 更新結果、validation、remaining gaps

## DO / DON'T

DO:
- Copilot completion report を根拠として残す
- approval record を根拠として残す
- source of truth と task-local artifact を分ける
- validation 結果を残す

DON'T:
- Copilot 修正完了前に正本化しない
- 未承認 draft を正本化しない
- workflow 変更を docs 更新に混ぜない
- product implementation を同時に進めない

## Checklist

- [updating-docs-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/updating-docs/references/checklists/updating-docs-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は `docs_updater` agent contract が決める。

## References

- docs index: [index.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md)
- binding: [docs_updater.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/docs_updater.toml)
- agent contract: [docs_updater.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/docs_updater/contracts/docs_updater.contract.json)

## Maintenance

- 権限、write scope、output obligation を skill 本体へ戻さない。
- workflow 変更は `skill-modification` へ分ける。
- Copilot 実装 workflow からは使わない。
