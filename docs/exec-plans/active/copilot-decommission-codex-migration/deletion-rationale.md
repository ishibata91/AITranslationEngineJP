# Deletion Rationale

## 削除分類

- Copilot 固有: `.github/agents/*.agent.md` と `.github/agents/references/*`
- VSCode tool 固有: `target: vscode`、`RunSubagent`、Copilot tool 名を含む runtime binding
- Codex で重複: `.github/skills/*` の mirror 済み skill
- live workflow 外: `.github` 側の実装 lane 正本

## 削除条件

- `.codex` 側に mirror copy が存在する。
- path、runtime、tool 名の最小置換が完了している。
- contract / permissions の JSON 構文検証が通る。
- `python3 scripts/harness/run.py --suite structure` の結果を確認済みである。
- 削除後に live workflow の正本が `.codex` だけになる。

## 削除しないもの

- `.github/workflows` が存在する場合は対象外にする。
- issue template や pull request template が存在する場合は対象外にする。
