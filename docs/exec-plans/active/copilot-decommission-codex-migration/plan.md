# Copilot 廃止と Codex 移植

## 状態

- status: `staged`
- source of truth: `.codex/README.md`
- migration policy: mirror copy first, prune after parity
- approval: user requested implementation on 2026-04-28

## 目的

GitHub Copilot lane の実装 workflow を Codex 側へ移す。
初回移植では情報削減を避け、skill、checklist、pattern、contract、permissions を保持する。

## 実装方針

- `.github/skills` の実装系 skill を `.codex/skills` へ mirror copy する。
- `.github/agents` の agent 説明と references を `.codex/agents` へ移す。
- Codex 側 agent は `implementation_*` 名にする。
- `.github` 側の削除は Codex 側 parity と validation 後に行う。
- 削除対象は `deletion-rationale.md` に分類してから消す。

## Validation

- mirror file list: `copy-map.md`
- deletion rationale: `deletion-rationale.md`
- structure check: `python3 scripts/harness/run.py --suite structure`
- residual reference check: `rg -n "Copilot|github-copilot|\\.github/skills|\\.github/agents|RunSubagent" .codex docs AGENTS.md`
